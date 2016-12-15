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

// http.client is concurrently safe and should be reused across multiple connections
var client http.Client

type jsonObject map[string]interface{}

// SqscRequest stores the data required to contact Squarescale services over HTTP.
type SqscRequest struct {
	Method string
	URL    string
	Token  string // user token for authentication
	Body   []byte // only relevant if Method != "GET"
}

// SqscResponse stores the result of a SquarescaleRequest once performed by the client.
type SqscResponse struct {
	Code int
	Body []byte
}

func get(url, token string) (int, []byte, error) {
	return request("GET", url, token, nil)
}

func post(url, token string, payload *jsonObject) (int, []byte, error) {
	return request("POST", url, token, payload)
}

func put(url, token string, payload *jsonObject) (int, []byte, error) {
	return request("PUT", url, token, payload)
}

func request(method, url, token string, payload *jsonObject) (int, []byte, error) {
	var bodyReader io.Reader
	if payload != nil {
		payloadBytes, err := json.Marshal(payload)
		if err != nil {
			return 0, []byte{}, err
		}

		bodyReader = bytes.NewReader(payloadBytes)
	}

	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		return 0, []byte{}, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "bearer "+token)

	res, err := client.Do(req)
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

func doRequest(r SqscRequest) (SqscResponse, error) {
	var req *http.Request
	var err error

	if len(r.Body) > 0 {
		req, err = http.NewRequest(r.Method, r.URL, bytes.NewReader(r.Body))
	} else {
		req, err = http.NewRequest(r.Method, r.URL, nil)
	}

	if err != nil {
		return SqscResponse{}, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "bearer "+r.Token)

	res, err := client.Do(req)
	if err != nil {
		return SqscResponse{}, err
	}

	defer res.Body.Close()
	bytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return SqscResponse{}, err
	}

	return SqscResponse{Code: res.StatusCode, Body: bytes}, nil
}

func unexpectedError(code int) error {
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
