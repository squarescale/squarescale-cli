package squarescale

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// http.client is concurrently safe and should be reused across multiple connections
var client http.Client

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

func readErrors(body []byte) ([]string, error) {
	var response struct {
		URL      []string `json:"url"`
		ShortURL []string `json:"short_url"`
		Project  []string `json:"project"`
	}

	err := json.Unmarshal(body, &response)
	if err != nil {
		return []string{}, err
	}

	var messages []string
	messages = append(messages, response.URL...)
	messages = append(messages, response.ShortURL...)
	messages = append(messages, response.Project...)
	return messages, nil
}
