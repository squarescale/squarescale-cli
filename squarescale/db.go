package squarescale

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// GetDBConfig asks the Squarescale API for the database config of a project.
// Returns, in this order:
// - the db engine in use (string)
// - the db instance in use (string)
func GetDBConfig(sqscURL, token, project string) (string, string, error) {
	req := SqscRequest{
		Method: "GET",
		URL:    sqscURL + "/projects/" + project,
		Token:  token,
	}

	res, err := doRequest(req)
	if err != nil {
		return "", "", err
	}

	var body struct {
		Engine        string `json:"db_engine"`
		InstanceClass string `json:"db_instance_class"`
	}

	err = json.Unmarshal(res.Body, &body)
	if err != nil {
		return "", "", err
	}

	if res.Code != http.StatusOK {
		return "", "", fmt.Errorf("'%s %s' return code: %d", req.Method, req.URL, res.Code)
	}

	return body.Engine, body.InstanceClass, nil
}

// ConfigDB calls the Squarescale API to update database scale options for a given project.
func ConfigDB(sqscURL, token, project, engine, instance string) error {
	var payload struct {
		Project struct {
			Enabled       bool   `json:"db_enabled"`
			Engine        string `json:"db_engine"`
			InstanceClass string `json:"db_instance_class"`
		} `json:"project"`
	}

	payload.Project.Enabled = true
	payload.Project.Engine = engine
	payload.Project.InstanceClass = instance
	payloadBytes, err := json.Marshal(&payload)
	if err != nil {
		return err
	}

	req := SqscRequest{
		Method: "POST",
		URL:    sqscURL + "/projects/" + project + "/cluster",
		Token:  token,
		Body:   payloadBytes,
	}

	res, err := doRequest(req)
	if err != nil {
		return err
	}

	if res.Code != http.StatusNoContent {
		return fmt.Errorf("'%s %s' return code: %d", req.Method, req.URL, res.Code)
	}

	return nil
}
