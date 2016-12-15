package squarescale

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// GetAvailableDBInstances returns all the database instances available for use in Squarescale.
func GetAvailableDBInstances(sqscURL, token string) ([]string, error) {
	url := sqscURL + "/db/instances"
	code, body, err := get(url, token)
	if err != nil {
		return []string{}, err
	}

	if code != http.StatusOK {
		return []string{}, unexpectedError(code)
	}

	var instancesList []string
	err = json.Unmarshal(body, &instancesList)
	if err != nil {
		return []string{}, err
	}

	return instancesList, nil
}

// GetAvailableDBEngines returns all the database engines available for use in Squarescale.
func GetAvailableDBEngines(sqscURL, token string) ([]string, error) {
	code, body, err := get(sqscURL+"/db/engines", token)
	if err != nil {
		return []string{}, err
	}

	if code != http.StatusOK {
		return []string{}, unexpectedError(code)
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
func GetDBConfig(sqscURL, token, project string) (bool, string, string, error) {
	code, body, err := get(sqscURL+"/projects/"+project, token)
	if err != nil {
		return false, "", "", err
	}

	if code != http.StatusOK {
		return false, "", "", unexpectedError(code)
	}

	var resp struct {
		Enabled       bool   `json:"db_enabled"`
		Engine        string `json:"db_engine"`
		InstanceClass string `json:"db_instance_class"`
	}

	err = json.Unmarshal(body, &body)
	if err != nil {
		return false, "", "", err
	}

	return resp.Enabled, resp.Engine, resp.InstanceClass, nil
}

// ConfigDB calls the Squarescale API to update database scale options for a given project.
func ConfigDB(sqscURL, token, project string, enabled bool, engine, instance string) error {
	payload := &jsonObject{
		"project": jsonObject{
			"db_enabled":        enabled,
			"db_engine":         engine,
			"db_instance_class": instance,
		},
	}

	code, _, err := post(sqscURL+"/projects/"+project+"/cluster", token, payload)
	if err != nil {
		return err
	}

	if code == http.StatusUnprocessableEntity {
		return fmt.Errorf("Invalid value for either database engine ('%s') or instance ('%s')", engine, instance)
	}

	if code != http.StatusNoContent {
		return unexpectedError(code)
	}

	return nil
}
