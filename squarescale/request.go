package squarescale

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

const supportedAPI string = "1"

type jsonObject map[string]interface{}

// Client is the basic structure to make API calls to Squarescale services
type Client struct {
	httpClient http.Client // http.client is concurrently safe and should be reused across multiple connections
	endpoint   string
	token      string
}

// NewClient creates a new Squarescale client
func NewClient(endpoint, token string) *Client {
	return &Client{
		endpoint: endpoint,
		token:    token,
	}
}

func (c *Client) get(path string) (int, []byte, error) {
	return c.request("GET", path, nil)
}

func (c *Client) post(path string, payload *jsonObject) (int, []byte, error) {
	return c.request("POST", path, payload)
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

	return fmt.Errorf("An unexpected error occurred (code: %d)", code)
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
