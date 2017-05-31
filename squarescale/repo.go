package squarescale

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// AddRepository asks the Squarescale service to attach the provided repository to the project.
func (c *Client) AddRepository(project, repoURL, buildService string) error {
	payload := &jsonObject{
		"repository": jsonObject{
			"url": repoURL,
		},
		"container": jsonObject{
			"build_service": buildService,
		},
	}

	code, body, err := c.post("/projects/"+project+"/repositories", payload)
	if err != nil {
		return err
	}

	switch code {
	case http.StatusCreated:
		return nil
	default:
		return readErrors(body, fmt.Sprintf("Cannot attach repository '%s' to project '%s'", repoURL, project))
	}
}

// ListRepositories asks the Squarescale service to lists all repositories for a given project.
func (c *Client) ListRepositories(project string) ([]string, error) {
	code, body, err := c.get("/projects/" + project + "/repositories")
	if err != nil {
		return []string{}, err
	}

	switch code {
	case http.StatusOK:
	case http.StatusNotFound:
		return []string{}, readErrors(body, fmt.Sprintf("Cannot list repositories for project '%s'", project))
	default:
		return []string{}, unexpectedHTTPError(code, body)
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
