package squarescale

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// GetAvailableDBInstances returns all the database instances available for use in Squarescale.
func (c *Client) GetAvailableDBInstances() ([]string, error) {
	code, body, err := c.get("/db/instances")
	if err != nil {
		return []string{}, err
	}

	if code != http.StatusOK {
		return []string{}, unexpectedHTTPError(code, body)
	}

	var instancesList []string
	err = json.Unmarshal(body, &instancesList)
	if err != nil {
		return []string{}, err
	}

	return instancesList, nil
}

// GetAvailableDBEngines returns all the database engines available for use in Squarescale.
func (c *Client) GetAvailableDBEngines() ([]string, error) {
	code, body, err := c.get("/db/engines")
	if err != nil {
		return []string{}, err
	}

	if code != http.StatusOK {
		return []string{}, unexpectedHTTPError(code, body)
	}

	var enginesList []string
	err = json.Unmarshal(body, &enginesList)
	if err != nil {
		return []string{}, err
	}

	return enginesList, nil
}

// GetDBConfig asks the Squarescale API for the database config of a project.
// Returns, in this order:
// - if the db is enabled
// - the db engine in use (string)
// - the db instance in use (string)
func (c *Client) GetDBConfig(project string) (bool, string, string, error) {
	code, body, err := c.get("/projects/" + project)
	if err != nil {
		return false, "", "", err
	}

	if code != http.StatusOK {
		return false, "", "", unexpectedHTTPError(code, body)
	}

	var resp struct {
		Enabled       bool   `json:"db_enabled"`
		Engine        string `json:"db_engine"`
		InstanceClass string `json:"db_instance_class"`
	}

	err = json.Unmarshal(body, &resp)
	if err != nil {
		return false, "", "", err
	}

	return resp.Enabled, resp.Engine, resp.InstanceClass, nil
}

// ConfigDB calls the Squarescale API to update database scale options for a given project.
func (c *Client) ConfigDB(project string, enabled bool, engine, instance string) (taskId int, err error) {
	payload := &jsonObject{
		"project": jsonObject{
			"db_enabled":        enabled,
			"db_engine":         engine,
			"db_instance_class": instance,
		},
	}

	code, body, err := c.post("/projects/"+project+"/cluster", payload)
	if err != nil {
		return 0, err
	}

	switch code {
	case http.StatusAccepted:
		fallthrough
	case http.StatusOK:
		break
	case http.StatusNoContent:
		return 0, nil
	case http.StatusUnprocessableEntity:
		return 0, fmt.Errorf("Invalid value for either database engine ('%s') or instance ('%s')", engine, instance)
	default:
		return 0, unexpectedHTTPError(code, body)
	}

	var resp struct {
		DBTask int `json:"db_task"`
	}

	err = json.Unmarshal(body, &resp)
	if err != nil {
		return 0, err
	}

	return resp.DBTask, nil
}
