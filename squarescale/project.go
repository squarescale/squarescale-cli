package squarescale

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

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

// ListProjects asks the Squarescale service for available projects.
func (c *Client) ListProjects() ([]string, error) {
	code, body, err := c.get("/projects")
	if err != nil {
		return []string{}, err
	}

	if code != http.StatusOK {
		return []string{}, unexpectedHTTPError(code, body)
	}

	var projectsJSON []struct {
		Name string `json:"name"`
	}
	err = json.Unmarshal(body, &projectsJSON)
	if err != nil {
		return []string{}, err
	}

	var projects []string
	for _, pJSON := range projectsJSON {
		projects = append(projects, pJSON.Name)
	}

	return projects, nil
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

// ProjectLogs gets the logs for a project container.
func (c *Client) ProjectLogs(project, container, after string) ([]string, string, error) {
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
		Environment   string `json:"environment"`
		Error         bool   `json:"error"`
		LogType       string `json:"log_type"`
		Message       string `json:"message"`
	}
	err = json.Unmarshal(body, &response)

	if err != nil {
		return []string{}, "", err
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
	var lastTimestamp string
	if len(response) == 0 {
		lastTimestamp = ""
	} else {
		lastTimestamp = response[len(response)-1].Timestamp
	}

	return messages, lastTimestamp, nil
}
