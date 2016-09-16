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

func FindProjectName(sqsc_url, token string) (string, error) {
	var c http.Client
	req, err := http.NewRequest("GET", sqsc_url+"/free_name", nil)
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
		return "", fmt.Errorf("Could not generate a free name", jsondata)
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

func CreateProject(sqsc_url, token, wanted_name string) (project_name string, err error) {
	var c http.Client
	var reqdata createProjectRequest

	if wanted_name == "" {
		wanted_name, err = FindProjectName(sqsc_url, token)
		if err != nil {
			return "", err
		}
	}

	reqdata.Project.Name = wanted_name
	reqbytes, err := json.Marshal(&reqdata)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", sqsc_url+"/projects", bytes.NewReader(reqbytes))
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
	project_name = wanted_name
	return
}
