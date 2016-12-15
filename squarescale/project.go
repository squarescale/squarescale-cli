package squarescale

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type payloadWithName struct {
	Name string `json:"name"`
}

type freeNameResponse struct {
}

// CheckProjectName asks the Squarescale service to validate a given project name.
func CheckProjectName(sqscURL, token, name string) (valid bool, same bool, fmtName string, err error) {
	code, body, err := post(sqscURL+"/free_name", token, &jsonObject{"name": name})
	if err != nil {
		return
	}

	if code != http.StatusOK {
		err = unexpectedError(code)
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
func FindProjectName(sqscURL, token string) (string, error) {
	code, body, err := get(sqscURL+"/free_name", token)
	if code != http.StatusOK {
		return "", unexpectedError(code)
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
func CreateProject(sqscURL, token, name string) error {
	payload := jsonObject{
		"project": jsonObject{
			"name": name,
		},
	}

	code, body, err := post(sqscURL+"/projects", token, &payload)
	if err != nil {
		return err
	}

	if code != http.StatusCreated {
		return unexpectedError(code)
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
func ListProjects(sqscURL, token string) ([]string, error) {
	code, body, err := get(sqscURL+"/projects", token)
	if err != nil {
		return []string{}, err
	}

	if code != http.StatusOK {
		return []string{}, unexpectedError(code)
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
func ProjectURL(sqscURL, token, project string) (string, error) {
	code, body, err := get(sqscURL+"/projects/"+project+"/url", token)
	if err != nil {
		return "", err
	}

	if code == http.StatusPreconditionFailed {
		return "", fmt.Errorf("Project '%s' not found", project)
	}

	if code == http.StatusNotFound {
		return "", fmt.Errorf("Project '%s' is not available on the web", project)
	}

	if code != http.StatusOK {
		return "", unexpectedError(code)
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
func ProjectLogs(sqscURL, token, project, container, after string) ([]string, string, error) {
	query := ""
	if after != "" {
		query = "?after=" + url.QueryEscape(after)
	}

	url := sqscURL + "/projects/" + project + "/logs/" + url.QueryEscape(container) + query
	code, body, err := get(url, token)
	if err != nil {
		return []string{}, "", err
	}

	if code == http.StatusBadRequest {
		return []string{}, "", fmt.Errorf("Project '%s' not found", project)
	}

	if code == http.StatusNotFound {
		return []string{}, "", fmt.Errorf("Container '%s' is not found for project '%s'", container, project)
	}

	if code != http.StatusOK {
		return []string{}, "", unexpectedError(code)
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
