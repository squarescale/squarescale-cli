package squarescale

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Environment struct {
	Preset     map[string]string            `json:"default"`
	Global     map[string]string            `json:"global"`
	PerService map[string]map[string]string `json:"per_service"`
}

func NewEnvironment(c APIClient, project string) (*Environment, error) {
	code, body, err := c.Get("/projects/" + project + "/environment")
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

	var env Environment
	err = json.Unmarshal(body, &env)
	if err != nil {
		return nil, err
	}

	return &env, nil
}

// SetEnvironmentVariables sets all the environment variables specified for the project.
func (c *Client) SetEnvironmentVariables(project string, env *Environment) error {
	code, body, err := c.put("/projects/"+project+"/environment/custom", &jsonObject{"environment": env, "format": "json"})
	if err != nil {
		return err
	}

	switch code {
	case http.StatusNoContent:
		return nil
	case http.StatusNotFound:
		return fmt.Errorf("Project '%s' not found", project)
	case http.StatusUnprocessableEntity:
		return fmt.Errorf("%s", body)
	default:
		return unexpectedHTTPError(code, body)
	}
}
