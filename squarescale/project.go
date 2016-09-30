package squarescale

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type payloadWithName struct {
	Name string `json:"name"`
}

type freeNameResponse struct {
}

// CheckProjectName asks the Squarescale service to validate a given project name.
func CheckProjectName(sqscURL, token, name string) (valid bool, same bool, fmtName string, err error) {
	payload := struct {
		Name string `json:"name"`
	}{Name: name}

	payloadBytes, err := json.Marshal(&payload)
	if err != nil {
		return
	}

	req := SqscRequest{
		Method: "POST",
		URL:    sqscURL + "/free_name",
		Token:  token,
		Body:   payloadBytes,
	}

	res, err := doRequest(req)
	if err != nil {
		return
	}

	var resBody struct {
		Valid bool   `json:"valid"`
		Same  bool   `json:"same"`
		Name  string `json:"name"`
	}

	err = json.Unmarshal(res.Body, &resBody)
	if err != nil {
		return
	}

	return resBody.Valid, resBody.Same, resBody.Name, nil
}

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
func CreateProject(sqscURL, token, name string) error {
	var payload struct {
		Project struct {
			Name string `json:"name"`
		} `json:"project"`
	}

	payload.Project.Name = name
	payloadBytes, err := json.Marshal(&payload)
	if err != nil {
		return err
	}

	req := SqscRequest{
		Method: "POST",
		URL:    sqscURL + "/projects",
		Token:  token,
		Body:   payloadBytes,
	}

	res, err := doRequest(req)
	if err != nil {
		return err
	}

	var response struct {
		Error string `json:"error"`
	}
	err = json.Unmarshal(res.Body, &response)
	if err != nil {
		return err
	}

	if res.Code != http.StatusCreated {
		return fmt.Errorf("'%s %s' return code: %d", req.Method, req.URL, res.Code)
	}

	return nil
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
