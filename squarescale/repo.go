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

// AddRepository asks the Squarescale service to attach the provided repository to the project.
func AddRepository(sqscURL, token, project, repoURL string) error {
	var payload addRepositoryRequest
	payload.Repository.URL = repoURL
	payloadBytes, err := json.Marshal(&payload)
	if err != nil {
		return fmt.Errorf("Cannot marshal payload %v: %v", payload, err)
	}

	endpoint := sqscURL + "/projects/" + project + "/repositories"
	req, err := http.NewRequest("POST", endpoint, bytes.NewReader(payloadBytes))
	if err != nil {
		return fmt.Errorf("Cannot create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "bearer "+token)

	var c http.Client
	res, err := c.Do(req)
	if err != nil {
		return fmt.Errorf("Cannot send request: %v", err)
	}

	if res.StatusCode != http.StatusCreated {
		jsonBody, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return fmt.Errorf("Could not read response: %v", err)
		}

		return fmt.Errorf("Cannot attach repository %s to project %s, errors: %v", repoURL, project, string(jsonBody))
	}

	return nil
}
