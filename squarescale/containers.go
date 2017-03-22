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
	Size      int
	Web       bool
	WebPort   int
	Command              []string
}

// GetContainers gets all the containers attached to a Project
func (c *Client) GetContainers(project string) ([]Container, error) {
	code, body, err := c.get("/projects/" + project + "/containers")
	if err != nil {
		return []Container{}, err
	}

	switch code {
	case http.StatusOK:
	case http.StatusNotFound:
		return []Container{}, fmt.Errorf("Project '%s' does not exist", project)
	default:
		return []Container{}, unexpectedHTTPError(code, body)
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
			Web:       c.Web,
			WebPort:   c.WebPort,
			Command:              c.Command,
		})
	}

	return containers, nil
}

// GetContainerInfo gets the container of a project based on its short name.
func (c *Client) GetContainerInfo(project, shortName string) (Container, error) {
	containers, err := c.GetContainers(project)
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
func (c *Client) ConfigContainer(container Container) error {
	payload := &jsonObject{
		"container": jsonObject{
			"size":        container.Size,
			"pre_command": container.Command,
		},
	}

	code, body, err := c.put(fmt.Sprintf("/containers/%d", container.ID), payload)
	if err != nil {
		return err
	}

	switch code {
	case http.StatusOK:
		return nil
	case http.StatusNotFound:
		return errors.New("Container does not exist")
	default:
		return unexpectedHTTPError(code, body)
	}
}
