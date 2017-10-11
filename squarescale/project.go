package squarescale

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type ProjectStatus struct {
	InfraStatus string `json:"infra_status"`
}

type UnprovisionError struct {
	Errors struct {
		Unprovision []string `json:unprovision`
	} `json:errors`
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
func (c *Client) CreateProject(name, infraType, databaseEngine, databaseClass string, enableDb bool) (taskId int, err error) {
	dbSettings := jsonObject{
		"enabled": enableDb,
	}
	if databaseEngine != "" {
		dbSettings["engine"] = databaseEngine
	}
	if databaseClass != "" {
		dbSettings["instance_class"] = databaseClass
	}
	payload := &jsonObject{
		"project": jsonObject{
			"name":       name,
			"infra_type": getInfraTypeEnumValue(infraType),
		},
		"database": dbSettings,
	}

	code, body, err := c.post("/projects", payload)
	if err != nil {
		return 0, err
	}

	if code != http.StatusCreated {
		return 0, unexpectedHTTPError(code, body)
	}

	var response struct {
		Error string `json:"error"`
		Task  int    `json:"task"`
	}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return 0, err
	}

	if response.Error != "" {
		err = errors.New(response.Error)
	}

	return response.Task, err
}

func getInfraTypeEnumValue(infraType string) string {
	return strings.Replace(infraType, "-", "_", 1)
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
	case http.StatusUnprocessableEntity:
		var errJSON UnprovisionError
		json.Unmarshal(body, &errJSON)
		return fmt.Errorf("Operation failed: %s", errJSON.Errors.Unprovision)
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
func (c *Client) ProjectLogs(project string, container string, after string) ([]string, string, error) {
	query := ""
	if after != "" {
		query = "?after=" + url.QueryEscape(after)
	}

	code, body, err := c.get("/projects/" + project + "/logs/" + url.QueryEscape(container) + query)
	if err != nil {
		return []string{}, "", err
	}

	switch code {
	case http.StatusOK:
	case http.StatusBadRequest:
		return []string{}, "", fmt.Errorf("Project '%s' not found", project)
	case http.StatusNotFound:
		return []string{}, "", fmt.Errorf("Container '%s' is not found for project '%s'", container, project)
	default:
		return []string{}, "", unexpectedHTTPError(code, body)
	}

	var response []struct {
		Timestamp     string `json:"timestamp"`
		ProjectName   string `json:"project_name"`
		ContainerName string `json:"container_name"`
		Error         bool   `json:"error"`
		Type          string `json:"type"`
		Message       string `json:"message"`
		Level         int    `json:"level"`
	}
	err = json.Unmarshal(body, &response)

	if err != nil {
		return []string{}, "", err
	}

	var messages []string
	for _, log := range response {
		var linePattern string
		var lt string
		var containerWithBrackets string
		var level string
		if log.Type == "docker" {
			lt = "docker"
		} else if log.Type == "nomad" {
			lt = "nomad "
		} else if log.Type == "event" {
			lt = "sqsc  "
		} else {
			lt = ""
		}
		if log.ContainerName == "" {
			containerWithBrackets = ""
		} else {
			containerWithBrackets = "[" + log.ContainerName + "] "
		}
		if log.Level >= 5 {
			level = "INFO "
		} else if log.Level >= 4 {
			level = "WARN "
		} else {
			level = "ERROR"
		}
		linePattern = "%s %s %s-- [%s] %s"
		if log.Error {
			linePattern = "\033[0;33m" + linePattern + "\033[0m"
		}
		t, err := time.Parse(time.RFC3339Nano, log.Timestamp)
		if err == nil {
			formatedTime := t.Format("2006-01-02 15:04:05.999")
			padTime := fmt.Sprintf("%-23s", formatedTime)
			messages = append(messages, fmt.Sprintf(linePattern, padTime, lt, containerWithBrackets, level, log.Message))
		} else {
			messages = append(messages, fmt.Sprintf(linePattern, log.Timestamp, lt, containerWithBrackets, level, log.Message))
		}
	}
	var lastTimestamp string
	if len(response) == 0 {
		lastTimestamp = ""
	} else {
		lastTimestamp = response[len(response)-1].Timestamp
	}

	return messages, lastTimestamp, nil
}
