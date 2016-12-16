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
	code, _, err := c.put("/projects/"+project+"/environment/custom", &jsonObject{"environment": vars})
	if err != nil {
		return err
	}

	if code == http.StatusNotFound {
		return fmt.Errorf("Project '%s' not found", project)
	}

	if code != http.StatusNoContent {
		return unexpectedError(code)
	}

	return nil
}

func (c *Client) envVariables(project, category string) (map[string]string, error) {
	code, body, err := c.get("/projects/" + project + "/environment/" + category)
	if err != nil {
		return map[string]string{}, err
	}

	if code == http.StatusNotFound {
		return map[string]string{}, fmt.Errorf("Project '%s' not found", project)
	}

	if code != http.StatusOK {
		return map[string]string{}, unexpectedError(code)
	}

	var variables map[string]string
	err = json.Unmarshal(body, &variables)
	if err != nil {
		return map[string]string{}, err
	}

	return variables, nil
}
