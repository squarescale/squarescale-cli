package squarescale

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// AddRepository asks the Squarescale service to attach the provided repository to the project.
func AddRepository(sqscURL, token, project, repoURL string) error {
	payload := jsonObject{
		"repository": jsonObject{
			"url": repoURL,
		},
	}

	url := sqscURL + "/projects/" + project + "/repositories"
	code, body, err := post(url, token, &payload)
	if err != nil {
		return err
	}

	if code != http.StatusCreated {
		errorsHeader := fmt.Sprintf("Cannot attach repository '%s' to project '%s'", repoURL, project)
		return readErrors(body, errorsHeader)
	}

	return nil
}

// ListRepositories asks the Squarescale service to lists all repositories for a given project.
func ListRepositories(sqscURL, token, project string) ([]string, error) {
	url := sqscURL + "/projects/" + project + "/repositories"
	code, body, err := get(url, token)
	if err != nil {
		return []string{}, err
	}

	if code == http.StatusNotFound {
		return []string{}, readErrors(body, fmt.Sprintf("Cannot list repositories for project '%s'", project))
	}

	if code != http.StatusOK {
		return []string{}, unexpectedError(code)
	}

	var repositoryJSON []struct {
		URL string `json:"url"`
	}
	err = json.Unmarshal(body, &repositoryJSON)
	if err != nil {
		return []string{}, err
	}

	var repositories []string
	for _, rJSON := range repositoryJSON {
		repositories = append(repositories, rJSON.URL)
	}

	return repositories, nil
}
