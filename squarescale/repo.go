package squarescale

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// AddRepository asks the Squarescale service to attach the provided repository to the project.
func AddRepository(sqscURL, token, project, repoURL string) ([]string, error) {
	var payload struct {
		Repository struct {
			URL string `json:"url"`
		} `json:"repository"`
	}

	payload.Repository.URL = repoURL
	payloadBytes, err := json.Marshal(&payload)
	if err != nil {
		return []string{}, err
	}

	req := SqscRequest{
		Method: "POST",
		URL:    sqscURL + "/projects/" + project + "/repositories",
		Token:  token,
		Body:   payloadBytes,
	}

	res, err := doRequest(req)
	if err != nil {
		return []string{}, err
	}

	if res.Code != http.StatusCreated {
		errMsgs, err := readErrors(res.Body)
		if err != nil {
			return []string{}, err
		}

		return errMsgs, fmt.Errorf("Cannot attach repository '%s' to project '%s'", repoURL, project)
	}

	return []string{}, nil
}

// ListRepositories asks the Squarescale service to lists all repositories for a given project.
func ListRepositories(sqscURL, token, project string) ([]string, error) {
	req := SqscRequest{
		Method: "GET",
		URL:    sqscURL + "/projects/" + project + "/repositories",
		Token:  token,
	}

	res, err := doRequest(req)
	if err != nil {
		return []string{}, err
	}

	if res.Code != http.StatusOK {
		return []string{}, fmt.Errorf("'%s %s' return code: %d", req.Method, req.URL, res.Code)
	}

	var repositoryJSON []struct {
		URL string `json:"url"`
	}
	err = json.Unmarshal(res.Body, &repositoryJSON)
	if err != nil {
		return []string{}, err
	}

	var repositories []string
	for _, rJSON := range repositoryJSON {
		repositories = append(repositories, rJSON.URL)
	}

	return repositories, nil
}
