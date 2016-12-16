package squarescale

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// AddRepository asks the Squarescale service to attach the provided repository to the project.
func (c *Client) AddRepository(project, repoURL string) error {
	payload := &jsonObject{
		"repository": jsonObject{
			"url": repoURL,
		},
	}

	code, body, err := c.post("/projects/"+project+"/repositories", payload)
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
func (c *Client) ListRepositories(project string) ([]string, error) {
	code, body, err := c.get("/projects/" + project + "/repositories")
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
