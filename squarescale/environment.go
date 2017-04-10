package squarescale

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Environment struct {
    Preset map[string]string `json:"default"`
    Custom struct {
        Global map[string]string `json:"global"`
        PerService map[string]map[string]string `json:"per_service"`
    } `json:"custom"`
}

// EnvironmentVariables gets all the environment variables specified for the project.
func (c *Client) EnvironmentVariables(project string) (*Environment, error) {
    return c.envVariables(project)
}

// SetEnvironmentVariables sets all the environment variables specified for the project.
func (c *Client) SetEnvironmentVariables(project string, env *Environment) error {
	code, body, err := c.put("/projects/"+project+"/environment/custom", &jsonObject{"environment": env.Custom, "format": "json"})
	if err != nil {
		return err
	}

	switch code {
	case http.StatusNoContent:
		return nil
	case http.StatusNotFound:
		return fmt.Errorf("Project '%s' not found", project)
	default:
		return unexpectedHTTPError(code, body)
	}
}

func (c *Client) envVariables(project string) (*Environment, error) {
	code, body, err := c.get("/projects/" + project + "/environment")
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
