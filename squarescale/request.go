package squarescale

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"os"
	"regexp"
	"strings"
	"time"

	multierr "github.com/hashicorp/go-multierror"
	mime "github.com/jpillora/go-mime"
	actioncable "github.com/squarescale/actioncable-go"
	"github.com/squarescale/logger"
)

const supportedAPI string = "1"

// JSONObject is a simple struct matching a JSON object with any kind of values
type JSONObject map[string]interface{}

// Client is the basic structure to make API calls to SquareScale services
type Client struct {
	httpClient        http.Client // http.client is concurrently safe and should be reused across multiple connections
	cachedCableClient *actioncable.Client
	endpoint          string
	token             string
	user              User
}

// NewClient creates a new SquareScale client
func NewClient(endpoint, token string) *Client {
	c := &Client{
		endpoint: endpoint,
		token:    token,
	}
	return c
}

func (c *Client) AddUser(user User) {
	c.user = user
}

func (c *Client) cableClient() *actioncable.Client {
	if c.cachedCableClient == nil {
		c.cachedCableClient = actioncable.NewClient(strings.Replace(c.endpoint, "http", "ws", 1)+"/cable", c.cableHeaders)
	}
	return c.cachedCableClient
}

func (c *Client) cableHeaders() (*http.Header, error) {
	h := make(http.Header, 3)
	h.Set("Authorization", "bearer "+c.token)
	h.Set("API-Version", supportedAPI)
	h.Set("Origin", c.endpoint)
	return &h, nil
}

func (c *Client) get(path string) (int, []byte, error) {
	code, response, err := c.request("GET", path, nil)
	logger.Trace.Println("HTTP Code:", code)
	logger.Trace.Printf("HTTP Response: %s\n", response)
	logger.Trace.Println("Err:", err)
	return code, response, err
}

func (c *Client) post(path string, payload interface{}) (int, []byte, error) {
	code, response, err := c.request("POST", path, payload)
	logger.Trace.Println("HTTP Code:", code)
	logger.Trace.Printf("HTTP Response: %s\n", response)
	logger.Trace.Println("Err:", err)
	return code, response, err
}

func (c *Client) patch(path string, payload interface{}) (int, []byte, error) {
	code, response, err := c.request("PATCH", path, payload)
	logger.Trace.Println("HTTP Code:", code)
	logger.Trace.Printf("HTTP Response: %s\n", response)
	logger.Trace.Println("Err:", err)
	return code, response, err
}

func (c *Client) delete(path string) (int, []byte, error) {
	code, response, err := c.request("DELETE", path, nil)
	logger.Trace.Println("HTTP Code:", code)
	logger.Trace.Printf("HTTP Response: %s\n", response)
	logger.Trace.Println("Err:", err)
	return code, response, err
}

func (c *Client) put(path string, payload interface{}) (int, []byte, error) {
	code, response, err := c.request("PUT", path, payload)
	logger.Trace.Println("HTTP Code:", code)
	logger.Trace.Printf("HTTP Response: %s\n", response)
	logger.Trace.Println("Err:", err)
	return code, response, err
}

func (c *Client) download(path, nodeName string) (int, error) {
	var bodyReader io.Reader
	req, err := http.NewRequest("GET", c.endpoint+path, bodyReader)
	if err != nil {
		return 0, err
	}

	req.Header.Set("Authorization", "bearer "+c.token)
	req.Header.Set("API-Version", supportedAPI)

	reqDump, err := httputil.DumpRequestOut(req, true)
	if err != nil {
		logger.Error.Println("Error dumping HTTP request:", err)
	} else {
		logger.Trace.Printf("REQUEST: %s", string(reqDump))
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return 0, err
	}
	fName := ""
	val, ok := res.Header["Content-Disposition"]
	if ok {
		if len(val) != 1 {
			return 0, fmt.Errorf("Content-Disposition header value: %+v", val)
		}
		re := regexp.MustCompile(`.*filename="([^"]*)".*`)
		fName = re.ReplaceAllString(val[0], "$1")
	}

	defer res.Body.Close()
	// Check server response
	if res.StatusCode != http.StatusOK {
		var reason map[string]interface{}
		rbytes, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return 0, fmt.Errorf("bad status: %s unabled to read body error: %+v", res.Status, err)
		}
		err1 := json.Unmarshal(rbytes, &reason)
		if err1 != nil {
			return 0, fmt.Errorf("bad status: %s unabled to decode error: %+v", res.Status, err1)
		}
		var val1 map[string]interface{}
		val1, ok := reason["errors"].(map[string]interface{})
		if ok {
			val, ok := val1["Error"]
			if ok {
				return 0, fmt.Errorf("Error: %s", val.(string))
			}
		}
		return 0, fmt.Errorf("bad status: %s", res.Status)
	}

	// Create the file
	out, err := os.Create(fmt.Sprintf("%s_%s", nodeName, fName))
	if err != nil  {
		return 0, err
	}
	defer out.Close()

	// Writer the body to file
	_, err = io.Copy(out, res.Body)
	if err != nil  {
		return 0, err
	}

	return res.StatusCode, nil
}

func (c *Client) request(method, path string, payload interface{}) (int, []byte, error) {
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

	ct := mime.TypeByExtension(".json")

	req.Header.Set("Content-Type", ct)
	req.Header.Set("Accept", ct)
	req.Header.Set("Authorization", "bearer "+c.token)
	req.Header.Set("API-Version", supportedAPI)

	reqDump, err := httputil.DumpRequestOut(req, true)
	if err != nil {
		logger.Error.Println("Error dumping HTTP request:", err)
	} else {
		logger.Trace.Printf("REQUEST: %s", string(reqDump))
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return 0, []byte{}, err
	}

	defer res.Body.Close()
	rbytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return 0, []byte{}, err
	}

	// No content type on 204 !!! Damnit
	if res.StatusCode == http.StatusServiceUnavailable {
		plannedDate := "Unknown"
		retryAfter := res.Header.Get("Retry-After")
		parseTime, err1 := time.Parse(time.RFC1123, retryAfter)
		if err1 == nil {
			deltaTime := parseTime.Sub(time.Now()).Round(1 * time.Minute)
			plannedDate = strings.Trim(fmt.Sprintf("%v", deltaTime), "0s")
		}
		err = fmt.Errorf("%s: Potential date of service availability: %s (aka %s from now)", res.Status, retryAfter, plannedDate)
	} else if res.StatusCode == http.StatusUnauthorized {
		extra := ""
		var reason map[string]interface{}
		err1 := json.Unmarshal(rbytes, &reason)
		if err1 != nil {
			extra = fmt.Sprintf("%v, ", err1)
		} else {
			if val, ok := reason["error"]; ok {
				extra = fmt.Sprintf("%v, ", val)
			} else {
				extra = fmt.Sprintf("%v, ", reason)
			}
		}
		err = fmt.Errorf("%s: %splease make sure your SQSC_TOKEN environment variable is properly set", res.Status, extra)
	} else if res.StatusCode != http.StatusNoContent && !strings.Contains(res.Header.Get("Content-Type"), ct) {
		if res.StatusCode == http.StatusInternalServerError &&
			strings.Contains(res.Header.Get("Content-Type"), "text/html") &&
			(bytes.Contains(rbytes, []byte("<title>Action Controller: Exception caught</title>")) || len(rbytes) == 0) {
			err = fmt.Errorf(
				"%s: Something went wrong on server side. Please report error to support@squarescale.com.",
				res.Status,
			)
		} else {
			if res.StatusCode == http.StatusGatewayTimeout {
				logger.Warn.Printf("Got %q from %s at %s", res.Status, res.Header.Get("Server"), res.Header.Get("Date"))
				err = fmt.Errorf("Got %q from %s at %s", res.Status, res.Header.Get("Server"), res.Header.Get("Date"))
			} else {
				err = fmt.Errorf(
					"Invalid response return code (got %q instead of %+v/%+v) or content type (got %q instead of %q).\nDo you use the right value for -endpoint=%s ?",
					res.Status, http.StatusNoContent, http.StatusInternalServerError, res.Header.Get("Content-Type"), ct, c.endpoint,
				)
			}
		}
	}

	return res.StatusCode, rbytes, err
}

type RequestError struct {
	Error         string   `json:"error"`
	DetailedError []string `json:"detailed_error"`
}

func unexpectedHTTPError(code int, body []byte) error {
	// try to decode the body with { "error": "<description>" }
	var description RequestError

	err := json.Unmarshal(body, &description)
	if err == nil && description.Error != "" {
		return errors.New("Error: " + description.Error)
	}

	var data map[string]interface{}
	json.Unmarshal(body, &data)
	err = decodeAnyError(data, "error", "")

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
