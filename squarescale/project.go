package squarescale

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/squarescale/logger"
)

// Project defined how project could see
type Project struct {
	Name              string `json:"name"`
	UUID              string `json:"uuid"`
	Provider          string `json:"provider"`
	Region            string `json:"region"`
	Organization      string `string:"organization"`
	InfraStatus       string `json:"infra_status"`
	ClusterSize       int    `json:"cluster_size"`
	NomadNodesReady   int    `json:"nomad_nodes_ready"`
	MonitoringEnabled bool   `json:"monitoring_enabled"`
	MonitoringEngine  string `json:"monitoring_engine"`
	SlackWebHook      string `json:"slack_webhook"`
}

/*
  json.organization p.organization&.name
*/

// UnprovisionError defined how export provision errors
type UnprovisionError struct {
	Errors struct {
		Unprovision []string `json:unprovision`
	} `json:errors`
}

func projectSettings(name string, cluster ClusterConfig) JSONObject {
	projectSettings := JSONObject{
		"name": name,
	}
	for attr, value := range cluster.ProjectCreationSettings() {
		projectSettings[attr] = value
	}
	return projectSettings
}

// CreateProject asks the SquareScale platform to create a new project
func (c *Client) CreateProject(payload *JSONObject) (taskId int, err error) {

	_, ok := (*payload)["credential_name"]
	if !c.user.IsAdmin && !ok {
		return 0, errors.New("Credential is mandatory")
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

// ProvisionProject asks the SquareScale platform to provision the project
func (c *Client) ProvisionProject(projectUUID string) (err error) {

	code, body, err := c.post(fmt.Sprintf("/projects/%s/provision", projectUUID), nil)
	if err != nil {
		return err
	}

	if code != http.StatusNoContent {
		return unexpectedHTTPError(code, body)
	}

	return nil
}

// UNProvisionProject asks the SquareScale platform to provision the project
func (c *Client) UNProvisionProject(projectUUID string) (err error) {

	code, body, err := c.post(fmt.Sprintf("/projects/%s/unprovision", projectUUID), nil)
	if err != nil {
		return err
	}

	if code != http.StatusNoContent {
		return unexpectedHTTPError(code, body)
	}

	return nil
}

// ListProjects asks the SquareScale service for available projects.
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

// FullListProjects asks the SquareScale service for available projects.
func (c *Client) FullListProjects() ([]Project, error) {
	projects, err := c.ListProjects()
	if err != nil {
		return nil, err
	}

	organizations, err := c.ListOrganizations()
	if err != nil {
		return nil, err
	}

	projectCount := len(projects)

	for _, organization := range organizations {
		projectCount += len(organization.Projects)
	}

	allProjects := make([]Project, projectCount)
	projectCount = 0

	for _, project := range projects {
		allProjects[projectCount].Name = project.Name
		allProjects[projectCount].UUID = project.UUID
		projectCount++
	}

	for _, organization := range organizations {
		for _, project := range organization.Projects {
			allProjects[projectCount].Name = organization.Name + "/" + project.Name
			allProjects[projectCount].UUID = project.UUID
			projectCount++
		}
	}

	return allProjects, nil
}

// ProjectByName get UUID for a project name
func (c *Client) ProjectByName(projectName string) (string, error) {
	projects, err := c.FullListProjects()
	if err != nil {
		return "", err
	}

	for _, project := range projects {
		if project.Name == projectName {
			return project.UUID, nil
		}
	}

	return "", fmt.Errorf("Project '%s' not found", projectName)
}

// GetProject return the status of the project
func (c *Client) GetProject(project string) (*Project, error) {
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

	var status Project
	err = json.Unmarshal(body, &status)
	if err != nil {
		return nil, err
	}

	return &status, nil
}

// WaitProject wait project provisioning
func (c *Client) WaitProject(projectUUID string, timeToWait int64) (string, error) {
	project, err := c.GetProject(projectUUID)
	if err != nil {
		return "", err
	}

	logger.Info.Println("wait for project : ", projectUUID)

	for project.InfraStatus != "ok" && err == nil {
		time.Sleep(time.Duration(timeToWait) * time.Second)
		project, err = c.GetProject(projectUUID)
		logger.Debug.Println("project status update: ", projectUUID)
	}

	return project.InfraStatus, err
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
		err = json.Unmarshal(body, &errJSON)
		if err != nil {
			return err
		}
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
