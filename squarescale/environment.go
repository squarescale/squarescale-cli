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

	var (
		result map[string]interface{}
		env    Environment
	)

	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	env = Environment{
		Preset:     make(map[string]string),
		Global:     make(map[string]string),
		PerService: make(map[string]map[string]string),
	}
	for name, value := range result["default"].(map[string]interface{}) {
		env.Preset[name] = value.(string)
	}
	for name, value := range result["global"].(map[string]interface{}) {
		env.Global[name] = value.(string)
	}
	for service, vars := range result["per_service"].(map[string]interface{}) {
		env.PerService[service] = make(map[string]string)

		for name, value := range vars.(map[string]interface{}) {
			switch v := value.(type) {
			default:
				fmt.Printf("Unexpected value type %T", v)
			case string:
				env.PerService[service][name] = value.(string)
			case map[string]interface{}:
				continue
			}
		}
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
