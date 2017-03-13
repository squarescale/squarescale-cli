package squarescale

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type ProjectStatus struct {
	InfraStatus string `json:"infra_status"`
}

// CheckProjectName asks the Squarescale service to validate a given project name.
func (c *Client) CheckProjectName(name string) (valid bool, same bool, fmtName string, err error) {
	code, body, err := c.post("/free_name", &jsonObject{"name": name})
	if err != nil {
		return
	}

	if code != http.StatusAccepted {
		err = unexpectedHTTPError(code, body)
		return
	}

	var resBody struct {
		Valid bool   `json:"valid"`
		Same  bool   `json:"same"`
		Name  string `json:"name"`
	}

	err = json.Unmarshal(body, &resBody)
	if err != nil {
		return
	}

	return resBody.Valid, resBody.Same, resBody.Name, nil
}

// FindProjectName asks the Squarescale service for a project name, using the provided token.
func (c *Client) FindProjectName() (string, error) {
	code, body, err := c.get("/free_name")
	if code != http.StatusOK {
		return "", unexpectedHTTPError(code, body)
	}

	var response struct {
		Name string `json:"name"`
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", err
	}

	return response.Name, nil
}

// CreateProject asks the Squarescale platform to create a new project, using the provided name and user token.
func (c *Client) CreateProject(name string) error {
	payload := &jsonObject{
		"project": jsonObject{
			"name": name,
		},
	}

	code, body, err := c.post("/projects", payload)
	if err != nil {
		return err
	}

	if code != http.StatusCreated {
		return unexpectedHTTPError(code, body)
	}

	var response struct {
		Error string `json:"error"`
	}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return err
	}

	return nil
}

type Project struct {
	Name            string `json:"name"`
	InfraStatus     string `json:"infra_status"`
	ClusterSize     int    `json:"cluster_size"`
	NomadNodesReady int    `json:"nomad_nodes_ready"`
}

// ListProjects asks the Squarescale service for available projects.
func (c *Client) ListProjects() ([]Project, error) {
	code, body, err := c.get("/projects")
	if err != nil {
		return nil, err
	}

	if code != http.StatusOK {
		return nil, unexpectedHTTPError(code, body)
	}

	var projectsJSON []Project
	err = json.Unmarshal(body, &projectsJSON)
	if err != nil {
		return nil, err
	}

	return projectsJSON, nil
}

// ProjectURL asks the Squarescale service the url of the project if available, using the provided token.
func (c *Client) ProjectSlackURL(project string) (string, error) {
	code, body, err := c.get("/projects/" + project + "/slack")
	if err != nil {
		return "", err
	}

	switch code {
	case http.StatusOK:
	case http.StatusPreconditionFailed:
		return "", fmt.Errorf("Project '%s' not found", project)
	case http.StatusNotFound:
		return "", fmt.Errorf("Project '%s' is not available on the web", project)
	default:
		return "", unexpectedHTTPError(code, body)
	}

	var response struct {
		URL string `json:"slack_webhook"`
	}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", err
	}

	return response.URL, nil
}

// ProjectURL asks the Squarescale service the url of the project if available, using the provided token.
func (c *Client) SetProjectSlackURL(project, url string) error {
	code, body, err := c.post("/projects/"+project+"/slack", &jsonObject{"project": &jsonObject{"slack_webhook": url}})
	if err != nil {
		return err
	}

	switch code {
	case http.StatusCreated:
		break
	default:
		err = unexpectedHTTPError(code, body)
		return err
	}

	var resBody struct {
		Valid bool   `json:"valid"`
		Same  bool   `json:"same"`
		Name  string `json:"name"`
	}

	err = json.Unmarshal(body, &resBody)
	if err != nil {
		return err
	}

	return nil
}

// ProjectURL asks the Squarescale service the url of the project if available, using the provided token.
func (c *Client) ProjectURL(project string) (string, error) {
	code, body, err := c.get("/projects/" + project + "/url")
	if err != nil {
		return "", err
	}

	switch code {
	case http.StatusOK:
	case http.StatusPreconditionFailed:
		return "", fmt.Errorf("Project '%s' not found", project)
	case http.StatusNotFound:
		return "", fmt.Errorf("Project '%s' is not available on the web", project)
	default:
		return "", unexpectedHTTPError(code, body)
	}

	var response struct {
		URL string `json:"url"`
	}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", err
	}

	return response.URL, nil
}

// ProjectStatus return the status of the project
func (c *Client) ProjectStatus(project string) (*ProjectStatus, error) {
	code, body, err := c.get("/projects/" + project)
	if err != nil {
		return nil, err
	}

	switch code {
	case http.StatusOK:
	case http.StatusNotFound:
		return nil, fmt.Errorf("Project '%s' not found", project)
	default:
		return nil, unexpectedHTTPError(code, body)
	}

	var status ProjectStatus
	err = json.Unmarshal(body, &status)
	if err != nil {
		return nil, err
	}

	return &status, nil
}

// ProjectUnprovision unprovisions a project
func (c *Client) ProjectUnprovision(project string) error {
	code, body, err := c.post("/projects/"+project+"/unprovision", nil)
	if err != nil {
		return err
	}

	switch code {
	case http.StatusOK:
	case http.StatusAccepted:
	case http.StatusNoContent:
	case http.StatusNotFound:
		return fmt.Errorf("Project '%s' not found", project)
	default:
		return unexpectedHTTPError(code, body)
	}

	return nil
}

// ProjectDelete deletes an unprovisionned project
func (c *Client) ProjectDelete(project string) error {
	code, body, err := c.delete("/projects/" + project)
	if err != nil {
		return err
	}

	switch code {
	case http.StatusOK:
	case http.StatusNoContent:
	case http.StatusNotFound:
		return fmt.Errorf("Project '%s' not found", project)
	default:
		return unexpectedHTTPError(code, body)
	}

	return nil
}

// ProjectLogs gets the logs for a project container.
func (c *Client) ProjectLogs(project string, container string, after int) ([]string, int, error) {
	query := ""
	if after > 0 {
		query = "?after=" + url.QueryEscape(strconv.Itoa(after))
	}

	code, body, err := c.get("/projects/" + project + "/logs/" + url.QueryEscape(container) + query)
	if err != nil {
		return []string{}, -1, err
	}

	switch code {
	case http.StatusOK:
	case http.StatusBadRequest:
		return []string{}, -1, fmt.Errorf("Project '%s' not found", project)
	case http.StatusNotFound:
		return []string{}, -1, fmt.Errorf("Container '%s' is not found for project '%s'", container, project)
	default:
		return []string{}, -1, unexpectedHTTPError(code, body)
	}

	var response []struct {
		Timestamp     string `json:"timestamp"`
		ProjectName   string `json:"project_name"`
		ContainerName string `json:"container_name"`
		Environment   string `json:"environment"`
		Error         bool   `json:"error"`
		LogType       string `json:"log_type"`
		Message       string `json:"message"`
		Offset        int    `json:"offset"`
	}
	err = json.Unmarshal(body, &response)

	if err != nil {
		return []string{}, -1, err
	}

	var messages []string
	for _, log := range response {
		var linePattern string
		if log.LogType == "stderr" {
			linePattern = "[%s][%s] E| %s"
		} else if log.LogType == "stdout" {
			linePattern = "[%s][%s] I| %s"
		} else {
			linePattern = "[%s][%s] S| %s"
		}
		if log.Error {
			linePattern = "\033[0;33m" + linePattern + "\033[0m"
		}
		t, err := time.Parse(time.RFC3339Nano, log.Timestamp)
		if err == nil {
			messages = append(messages, fmt.Sprintf(linePattern, t.Format("2006-01-02T15:04:05.999Z07:00"), log.ContainerName, log.Message))
		} else {
			messages = append(messages, fmt.Sprintf(linePattern, log.Timestamp, log.ContainerName, log.Message))
		}
	}
	var lastOffset int
	if len(response) == 0 {
		lastOffset = -1
	} else {
		lastOffset = response[len(response)-1].Offset
	}

	return messages, lastOffset, nil
}
