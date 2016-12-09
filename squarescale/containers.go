package squarescale

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

// GetContainerInfo gets the id of a container based on its short name.
func GetContainerInfo(sqscURL, token, project, shortName string) (int, int, string, error) {
	req := SqscRequest{
		Method: "GET",
		URL:    fmt.Sprintf("%s/projects/%s/containers", sqscURL, project),
		Token:  token,
	}

	res, err := doRequest(req)
	if err != nil {
		return 0, 0, "", err
	}

	if res.Code == http.StatusNotFound {
		return 0, 0, "", fmt.Errorf("Project '%s' does not exist", project)
	}

	if res.Code != http.StatusOK {
		return 0, 0, "", fmt.Errorf("'%s %s' return code: %d", req.Method, req.URL, res.Code)
	}

	var containersByID []struct {
		ID       int      `json:"id"`
		ShortURL string   `json:"short_url"`
		Command  []string `json:"pre_command"`
		Size     int      `json:"size"`
	}
	if err := json.Unmarshal(res.Body, &containersByID); err != nil {
		return 0, 0, "", err
	}

	for _, container := range containersByID {
		if container.ShortURL == shortName {
			return container.ID, container.Size, strings.Join(container.Command, " "), nil
		}
	}

	return 0, 0, "", fmt.Errorf("Container '%s' not found for project '%s'", shortName, project)
}

// ConfigContainer calls the API to update the number of instances and update command.
func ConfigContainer(sqscURL, token string, containerID, nInstances int, command string) error {
	var payload struct {
		Container struct {
			Size    int    `json:"size"`
			Command string `json:"pre_command"`
		} `json:"container"`
	}

	payload.Container.Size = nInstances
	payload.Container.Command = command
	payloadBytes, err := json.Marshal(&payload)
	if err != nil {
		return err
	}

	req := SqscRequest{
		Method: "PUT",
		URL:    fmt.Sprintf("%s/containers/%d", sqscURL, containerID),
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
