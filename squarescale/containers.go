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
	req := SqscRequest{
		Method: "GET",
		URL:    fmt.Sprintf("%s/projects/%s/containers", sqscURL, project),
		Token:  token,
	}

	res, err := doRequest(req)
	if err != nil {
		return []Container{}, err
	}

	if res.Code == http.StatusNotFound {
		return []Container{}, fmt.Errorf("Project '%s' does not exist", project)
	}

	if res.Code != http.StatusOK {
		return []Container{}, fmt.Errorf("'%s %s' return code: %d", req.Method, req.URL, res.Code)
	}

	var containersByID []struct {
		ID       int      `json:"id"`
		ShortURL string   `json:"short_url"`
		Command  []string `json:"pre_command"`
		Size     int      `json:"size"`
		Web      bool     `json:"web"`
		WebPort  int      `json:"web_port"`
	}

	if err := json.Unmarshal(res.Body, &containersByID); err != nil {
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
	var payload struct {
		Container struct {
			Size    int    `json:"size"`
			Command string `json:"pre_command"`
		} `json:"container"`
	}

	payload.Container.Size = container.Size
	payload.Container.Command = container.Command
	payloadBytes, err := json.Marshal(&payload)
	if err != nil {
		return err
	}

	req := SqscRequest{
		Method: "PUT",
		URL:    fmt.Sprintf("%s/containers/%d", sqscURL, container.ID),
		Token:  token,
		Body:   payloadBytes,
	}

	res, err := doRequest(req)
	if err != nil {
		return err
	}

	if res.Code == http.StatusNotFound {
		return errors.New("Container does not exist")
	}

	if res.Code != http.StatusOK {
		return fmt.Errorf("'%s %s' return code: %d", req.Method, req.URL, res.Code)
	}

	return nil
}
