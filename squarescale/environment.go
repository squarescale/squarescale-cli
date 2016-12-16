package squarescale

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// CustomEnvironmentVariables gets all the custom environment variables specified for the project.
func (c *Client) CustomEnvironmentVariables(project string) (map[string]string, error) {
	return c.envVariables(project, "custom")
}

// EnvironmentVariables gets all the environment variables specified for the project.
func (c *Client) EnvironmentVariables(project string) (map[string]string, error) {
	return c.envVariables(project, "")
}

// SetEnvironmentVariables sets all the environment variables specified for the project.
func (c *Client) SetEnvironmentVariables(project string, vars map[string]string) error {
	code, body, err := c.put("/projects/"+project+"/environment/custom", &jsonObject{"environment": vars})
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

func (c *Client) envVariables(project, category string) (map[string]string, error) {
	code, body, err := c.get("/projects/" + project + "/environment/" + category)
	if err != nil {
		return map[string]string{}, err
	}

	switch code {
	case http.StatusOK:
	case http.StatusNotFound:
		return map[string]string{}, fmt.Errorf("Project '%s' not found", project)
	default:
		return map[string]string{}, unexpectedHTTPError(code, body)
	}

	var variables map[string]string
	err = json.Unmarshal(body, &variables)
	if err != nil {
		return map[string]string{}, err
	}

	return variables, nil
}
