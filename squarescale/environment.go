package squarescale

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// CustomEnvironmentVariables gets all the custom environment variables specified for the project.
func CustomEnvironmentVariables(sqscURL, token, project string) (map[string]string, error) {
	return envVariables(sqscURL, token, project, "custom")
}

// EnvironmentVariables gets all the environment variables specified for the project.
func EnvironmentVariables(sqscURL, token, project string) (map[string]string, error) {
	return envVariables(sqscURL, token, project, "")
}

// SetEnvironmentVariables sets all the environment variables specified for the project.
func SetEnvironmentVariables(sqscURL, token, project string, vars map[string]string) error {
	payload := jsonObject{"environment": vars}
	url := fmt.Sprintf("%s/projects/%s/environment/custom", sqscURL, project)
	code, _, err := put(url, token, &payload)
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

func envVariables(sqscURL, token, project, category string) (map[string]string, error) {
	url := fmt.Sprintf("%s/projects/%s/environment/%s", sqscURL, project, category)
	code, body, err := get(url, token)
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
