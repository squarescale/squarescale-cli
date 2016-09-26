package squarescale

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type freeNameResponse struct {
	Name string `json:"name"`
}

// FindProjectName asks the Squarescale service for a project name, using the provided token.
func FindProjectName(sqscURL, token string) (string, error) {
	var c http.Client
	req, err := http.NewRequest("GET", sqscURL+"/free_name", nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "bearer "+token)
	res, err := c.Do(req)
	if err != nil {
		return "", fmt.Errorf("Could not send request: %v", err)
	}

	jsondata, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("Could not read response: %v", err)
	}

	var response freeNameResponse
	err = json.Unmarshal(jsondata, &response)
	if err != nil {
		return "", fmt.Errorf("Could not parse JSON result %s: %v", jsondata, err)
	}

	if res.StatusCode != 200 {
		return "", fmt.Errorf("Could not generate a free name %v", jsondata)
	}
	return response.Name, nil
}

type createProjectRequest struct {
	Project struct {
		Name string `json:"name"`
	} `json:"project"`
}

type createProjectResponse struct {
	Error string `json:"error"`
}

// CreateProject asks the Squarescale platform to create a new project, using the provided name and user token.
func CreateProject(sqscURL, token, wantedName string) (projectName string, err error) {
	var c http.Client
	var reqdata createProjectRequest

	if wantedName == "" {
		wantedName, err = FindProjectName(sqscURL, token)
		if err != nil {
			return "", err
		}
	}

	reqdata.Project.Name = wantedName
	reqbytes, err := json.Marshal(&reqdata)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", sqscURL+"/projects", bytes.NewReader(reqbytes))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "bearer "+token)
	res, err := c.Do(req)
	if err != nil {
		return "", fmt.Errorf("Could not send request: %v", err)
	}

	jsondata, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("Could not read response: %v", err)
	}

	var response createProjectResponse
	err = json.Unmarshal(jsondata, &response)
	if err != nil {
		return "", fmt.Errorf("Could not parse JSON result %s: %v", jsondata, err)
	}

	if res.StatusCode != 201 {
		return "", fmt.Errorf("Could not create project: %s", jsondata)
	}
	projectName = wantedName
	return
}

// ListProjects asks the Squarescale service for available projects.
func ListProjects(sqscURL, token string) ([]string, error) {
	endpoint := sqscURL + "/projects"
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return []string{}, fmt.Errorf("Cannot create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "bearer "+token)

	var c http.Client
	res, err := c.Do(req)
	if err != nil {
		return []string{}, fmt.Errorf("Cannot send request: %v", err)
	}

	jsonBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return []string{}, fmt.Errorf("Cannot read response: %v", err)
	}

	if res.StatusCode != http.StatusOK {
		return []string{}, fmt.Errorf("Cannot read response, received status is %d", res.StatusCode)
	}

	type ProjectJSON struct {
		Name string `json:"name"`
	}

	var projectsJSON []ProjectJSON
	err = json.Unmarshal(jsonBody, &projectsJSON)
	if err != nil {
		return []string{}, fmt.Errorf("Cannot parse JSON %s: %v", string(jsonBody), err)
	}

	var projects []string
	for _, pJSON := range projectsJSON {
		projects = append(projects, pJSON.Name)
	}

	return projects, nil
}
