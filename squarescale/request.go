package squarescale

import (
	"bytes"
	"io/ioutil"
	"net/http"
)

// http.client is concurrently safe and should be reused across multiple connections
var client http.Client

// SquarescaleRequest stores the data required to contact Squarescale services over HTTP.
type SquarescaleRequest struct {
	Method string
	URL    string
	Token  string // user token for authentication
	Body   []byte // only relevant if Method != "GET"
}

// SquarescaleResponse stores the result of a SquarescaleRequest once performed by the client.
type SquarescaleResponse struct {
	Code int
	Body []byte
}

func setHeaders(req *http.Request, token string) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "bearer "+token)
}

func request(r SquarescaleRequest) (SquarescaleResponse, error) {
	var req *http.Request
	var err error

	if len(r.Body) > 0 {
		req, err = http.NewRequest(r.Method, r.URL, bytes.NewReader(r.Body))
	} else {
		req, err = http.NewRequest(r.Method, r.URL, nil)
	}

	if err != nil {
		return SquarescaleResponse{}, err
	}

	setHeaders(req, r.Token)
	res, err := client.Do(req)
	if err != nil {
		return SquarescaleResponse{}, err
	}

	defer res.Body.Close()
	bytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return SquarescaleResponse{}, err
	}

	return SquarescaleResponse{Code: res.StatusCode, Body: bytes}, nil
}
