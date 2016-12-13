package squarescale

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// EnvironmentVariables gets all the environment variables specified for the project.
func EnvironmentVariables(sqscURL, token, project string) (map[string]string, error) {
	req := SqscRequest{
		Method: "GET",
		URL:    fmt.Sprintf("%s/projects/%s/environment", sqscURL, project),
		Token:  token,
	}

	res, err := doRequest(req)
	if err != nil {
		return map[string]string{}, err
	}

	if res.Code == http.StatusNotFound {
		return map[string]string{}, fmt.Errorf("Project '%s' not found", project)
	}

	if res.Code != http.StatusOK {
		return map[string]string{}, fmt.Errorf("'%s %s' return code: %d", req.Method, req.URL, res.Code)
	}

	var variables map[string]string
	err = json.Unmarshal(res.Body, &variables)
	if err != nil {
		return map[string]string{}, err
	}

	return variables, nil
}

// SetEnvironmentVariables sets all the environment variables specified for the project.
func SetEnvironmentVariables(sqscURL, token, project string, vars map[string]string) error {
	payload := jsonObject{"environment": vars}
	payloadBytes, err := json.Marshal(&payload)
	if err != nil {
		return err
	}

	req := SqscRequest{
		Method: "PUT",
		URL:    fmt.Sprintf("%s/projects/%s/environment/custom", sqscURL, project),
		Token:  token,
		Body:   payloadBytes,
	}

	res, err := doRequest(req)
	if err != nil {
		return err
	}

	if res.Code == http.StatusNotFound {
		return fmt.Errorf("Project '%s' not found", project)
	}

	if res.Code != http.StatusOK {
		return fmt.Errorf("'%s %s' return code: %d", req.Method, req.URL, res.Code)
	}

	return nil
}
