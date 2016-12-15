package squarescale

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

// Container describes a project container as returned by the Squarescale API
type Container struct {
	ID        int
	ShortName string
	Command   string
	Size      int
	Web       bool
	WebPort   int
}

// GetContainers gets all the containers attached to a Project
func GetContainers(sqscURL, token, project string) ([]Container, error) {
	url := fmt.Sprintf("%s/projects/%s/containers", sqscURL, project)
	code, body, err := get(url, token)
	if err != nil {
		return []Container{}, err
	}

	if code == http.StatusNotFound {
		return []Container{}, fmt.Errorf("Project '%s' does not exist", project)
	}

	if code != http.StatusOK {
		return []Container{}, unexpectedError(code)
	}

	var containersByID []struct {
		ID       int      `json:"id"`
		ShortURL string   `json:"short_url"`
		Command  []string `json:"pre_command"`
		Size     int      `json:"size"`
		Web      bool     `json:"web"`
		WebPort  int      `json:"web_port"`
	}

	if err := json.Unmarshal(body, &containersByID); err != nil {
		return []Container{}, err
	}

	var containers []Container
	for _, c := range containersByID {
		containers = append(containers, Container{
			ID:        c.ID,
			ShortName: c.ShortURL,
			Size:      c.Size,
			Command:   strings.Join(c.Command, " "),
			Web:       c.Web,
			WebPort:   c.WebPort,
		})
	}

	return containers, nil
}

// GetContainerInfo gets the container of a project based on its short name.
func GetContainerInfo(sqscURL, token, project, shortName string) (Container, error) {
	containers, err := GetContainers(sqscURL, token, project)
	if err != nil {
		return Container{}, err
	}

	for _, container := range containers {
		if container.ShortName == shortName {
			return container, nil
		}
	}

	return Container{}, fmt.Errorf("Container '%s' not found for project '%s'", shortName, project)
}

// ConfigContainer calls the API to update the number of instances and update command.
func ConfigContainer(sqscURL, token string, container Container) error {
	url := fmt.Sprintf("%s/containers/%d", sqscURL, container.ID)
	payload := jsonObject{
		"container": jsonObject{
			"size":        container.Size,
			"pre_command": container.Command,
		},
	}

	code, _, err := put(url, token, &payload)
	if err != nil {
		return err
	}

	if code == http.StatusNotFound {
		return errors.New("Container does not exist")
	}

	if code != http.StatusOK {
		return unexpectedError(code)
	}

	return nil
}
