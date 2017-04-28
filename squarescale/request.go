package squarescale

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	actioncable "github.com/bgentry/actioncable-go"
	multierr "github.com/hashicorp/go-multierror"
)

const supportedAPI string = "1"

type jsonObject map[string]interface{}

// Client is the basic structure to make API calls to Squarescale services
type Client struct {
	httpClient        http.Client // http.client is concurrently safe and should be reused across multiple connections
	cachedCableClient *actioncable.Client
	endpoint          string
	token             string
}

// NewClient creates a new Squarescale client
func NewClient(endpoint, token string) *Client {
	c := &Client{
		endpoint: endpoint,
		token:    token,
	}
	return c
}

func (c *Client) cableClient() *actioncable.Client {
	if c.cachedCableClient == nil {
		c.cachedCableClient = actioncable.NewClient(strings.Replace(c.endpoint, "http", "ws", 1)+"/cable", c.cableHeaders)
	}
	return c.cachedCableClient
}

func (c *Client) cableHeaders() (*http.Header, error) {
	h := make(http.Header, 2)
	h.Set("Authorization", "bearer "+c.token)
	h.Set("API-Version", supportedAPI)
	h.Set("Origin", c.endpoint)
	return &h, nil
}

func (c *Client) get(path string) (int, []byte, error) {
	return c.request("GET", path, nil)
}

func (c *Client) post(path string, payload *jsonObject) (int, []byte, error) {
	return c.request("POST", path, payload)
}

func (c *Client) patch(path string, payload *jsonObject) (int, []byte, error) {
	return c.request("PATCH", path, payload)
}

func (c *Client) delete(path string) (int, []byte, error) {
	return c.request("DELETE", path, nil)
}

func (c *Client) put(path string, payload *jsonObject) (int, []byte, error) {
	return c.request("PUT", path, payload)
}

func (c *Client) request(method, path string, payload *jsonObject) (int, []byte, error) {
	var bodyReader io.Reader
	if payload != nil {
		payloadBytes, err := json.Marshal(payload)
		if err != nil {
			return 0, []byte{}, err
		}

		bodyReader = bytes.NewReader(payloadBytes)
	}

	req, err := http.NewRequest(method, c.endpoint+path, bodyReader)
	if err != nil {
		return 0, []byte{}, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "bearer "+c.token)
	req.Header.Set("API-Version", supportedAPI)

	res, err := c.httpClient.Do(req)
	if err != nil {
		return 0, []byte{}, err
	}

	defer res.Body.Close()
	bytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return 0, []byte{}, err
	}

	return res.StatusCode, bytes, nil
}

func unexpectedHTTPError(code int, body []byte) error {
	// try to decode the body with { "error": "<description>" }
	var description struct {
		Error string `json:"error"`
	}

	json.Unmarshal(body, description)
	if description.Error != "" {
		return errors.New("Error: " + description.Error)
	}

	var data map[string]interface{}
	json.Unmarshal(body, &data)
	err := decodeAnyError(data, "error", "")

	if err == nil {
		err = fmt.Errorf("An unexpected error occurred (code: %d)", code)
	}

	return err
}

func decodeAnyError(errs map[string]interface{}, exclude, prefix string) error {
	var err *multierr.Error
	for k, v := range errs {
		var pre string
		if prefix == "" {
			pre = k
		} else {
			pre = prefix + " " + k
		}
		switch val := v.(type) {
		case string:
			err = multierr.Append(err, fmt.Errorf("%s: %s", pre, val))
		case []interface{}:
			for _, value := range val {
				if s, ok := value.(string); ok {
					err = multierr.Append(err, fmt.Errorf("%s: %s", pre, s))
				}
			}
		case map[string]interface{}:
			err = multierr.Append(err, decodeAnyError(val, "", pre))
		}
	}
	return err.ErrorOrNil()
}

func readErrors(body []byte, header string) error {
	var response struct {
		URL      []string `json:"url"`
		ShortURL []string `json:"short_url"`
		Project  []string `json:"project"`
	}

	err := json.Unmarshal(body, &response)
	if err != nil {
		return err
	}

	var messages []string
	messages = append(messages, response.URL...)
	messages = append(messages, response.ShortURL...)
	messages = append(messages, response.Project...)

	var errorMessage string
	for _, message := range messages {
		errorMessage += "  " + message + "\n"
	}

	if len(errorMessage) == 0 {
		errorMessage = header
	} else {
		errorMessage = header + "\n\n" + errorMessage
	}

	return errors.New(errorMessage)
}
