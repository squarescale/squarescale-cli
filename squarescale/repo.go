package squarescale

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type addRepositoryRequest struct {
	Repository struct {
		URL string `json:"url"`
	} `json:"repository"`
}

type addRepositoryErrorResponse struct {
	URL      []string `json:"url"`
	ShortURL []string `json:"short_url"`
}

// AddRepository asks the Squarescale service to attach the provided repository to the project.
func AddRepository(sqscURL, token, project, repoURL string) ([]string, error) {
	var payload addRepositoryRequest
	payload.Repository.URL = repoURL
	payloadBytes, err := json.Marshal(&payload)
	if err != nil {
		return []string{}, fmt.Errorf("Cannot marshal payload %v: %v", payload, err)
	}

	endpoint := sqscURL + "/projects/" + project + "/repositories"
	req, err := http.NewRequest("POST", endpoint, bytes.NewReader(payloadBytes))
	if err != nil {
		return []string{}, fmt.Errorf("Cannot create request: %v", err)
	}

	setHeaders(req, token)
	var c http.Client
	res, err := c.Do(req)
	if err != nil {
		return []string{}, fmt.Errorf("Cannot send request: %v", err)
	}

	if res.StatusCode != http.StatusCreated {
		jsonBody, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return []string{}, fmt.Errorf("Cannot read response: %v", err)
		}

		var response addRepositoryErrorResponse
		err = json.Unmarshal(jsonBody, &response)
		if err != nil {
			return []string{}, fmt.Errorf("Cannot read JSON: %v", err)
		}

		errors := append(response.URL, response.ShortURL...)
		return errors, fmt.Errorf("Cannot attach repository '%s' to project '%s'", repoURL, project)
	}

	return []string{}, nil
}
