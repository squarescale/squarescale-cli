package squarescale

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// FindProjectName asks the Squarescale service for a project name, using the provided token.
func FindProjectName(sqscURL, token string) (string, error) {
	req := SqscRequest{
		Method: "GET",
		URL:    sqscURL + "/free_name",
		Token:  token,
	}

	res, err := doRequest(req)
	var response struct {
		Name string `json:"name"`
	}

	err = json.Unmarshal(res.Body, &response)
	if err != nil {
		return "", err
	}

	if res.Code != http.StatusOK {
		return "", fmt.Errorf("'%s %s' return code: %d", req.Method, req.URL, res.Code)
	}
	return response.Name, nil
}

// CreateProject asks the Squarescale platform to create a new project, using the provided name and user token.
func CreateProject(sqscURL, token, wantedName string) (string, error) {
	if wantedName == "" {
		var err error
		wantedName, err = FindProjectName(sqscURL, token)
		if err != nil {
			return "", err
		}
	}

	var payload struct {
		Project struct {
			Name string `json:"name"`
		} `json:"project"`
	}

	payload.Project.Name = wantedName
	payloadBytes, err := json.Marshal(&payload)
	if err != nil {
		return "", err
	}

	req := SqscRequest{
		Method: "POST",
		URL:    sqscURL + "/projects",
		Token:  token,
		Body:   payloadBytes,
	}

	res, err := doRequest(req)
	if err != nil {
		return "", err
	}

	var response struct {
		Error string `json:"error"`
	}
	err = json.Unmarshal(res.Body, &response)
	if err != nil {
		return "", err
	}

	if res.Code != http.StatusCreated {
		return "", fmt.Errorf("'%s %s' return code: %d", req.Method, req.URL, res.Code)
	}

	return wantedName, nil
}

// ListProjects asks the Squarescale service for available projects.
func ListProjects(sqscURL, token string) ([]string, error) {
	req := SqscRequest{
		Method: "GET",
		URL:    sqscURL + "/projects",
		Token:  token,
	}

	res, err := doRequest(req)
	if err != nil {
		return []string{}, err
	}

	if res.Code != http.StatusOK {
		return []string{}, fmt.Errorf("'%s %s' return code: %d", req.Method, req.URL, res.Code)
	}

	var projectsJSON []struct {
		Name string `json:"name"`
	}
	err = json.Unmarshal(res.Body, &projectsJSON)
	if err != nil {
		return []string{}, err
	}

	var projects []string
	for _, pJSON := range projectsJSON {
		projects = append(projects, pJSON.Name)
	}

	return projects, nil
}
