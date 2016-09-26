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
		var response struct {
			URL      []string `json:"url"`
			ShortURL []string `json:"short_url"`
		}

		err = json.Unmarshal(res.Body, &response)
		if err != nil {
			return []string{}, err
		}

		errors := append(response.URL, response.ShortURL...)
		return errors, fmt.Errorf("Cannot attach repository '%s' to project '%s'", repoURL, project)
	}

	return []string{}, nil
}
