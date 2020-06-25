package squarescale

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type DataseSize struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type DataseEngine struct {
	Name    string `json:"name"`
	Label   string `json:"label"`
	Version string `json:"version"`
}

type DbConfig struct {
	Enabled bool   `json:"db_enabled"`
	Engine  string `json:"db_engine"`
	Size    string `json:"db_size"`
	Version string `json:"db_version"`
}

func (db *DbConfig) String() string {
	if db.Enabled {
		return db.Size + " " + db.Engine
	} else {
		return "none"
	}
}

// GetAvailableDBSizes return the db node size available for a cloud provider
// on a given region
func (c *Client) GetAvailableDBSizes(provider, region string) ([]DataseSize, error) {
	code, body, err := c.get(fmt.Sprintf("/infra/providers/%s/regions/%s/database_sizes", provider, url.PathEscape(region)))
	if err != nil {
		return nil, err
	}

	if code != http.StatusOK {
		return nil, unexpectedHTTPError(code, body)
	}

	var sizes []DataseSize

	err = json.Unmarshal(body, &sizes)
	if err != nil {
		return nil, err
	}

	return sizes, nil
}

// GetAvailableDBEngines returns all the database engines available for a cloud provider
// on a given region
func (c *Client) GetAvailableDBEngines(provider, region string) ([]DataseEngine, error) {
	code, body, err := c.get(fmt.Sprintf("/infra/providers/%s/regions/%s/database_engines", provider, url.PathEscape(region)))
	if err != nil {
		return nil, err
	}

	if code != http.StatusOK {
		return nil, unexpectedHTTPError(code, body)
	}

	var enginesList []DataseEngine
	err = json.Unmarshal(body, &enginesList)
	if err != nil {
		return nil, err
	}

	return enginesList, nil
}

// GetDBConfig asks the Squarescale API for the database config of a project.
// Returns, in this order:
// - if the db is enabled
// - the db engine in use (string)
// - the db instance in use (string)
func (c *Client) GetDBConfig(uuid string) (*DbConfig, error) {
	code, body, err := c.get("/projects/" + uuid)
	if err != nil {
		return nil, err
	}

	if code != http.StatusOK {
		return nil, unexpectedHTTPError(code, body)
	}

	var db DbConfig
	err = json.Unmarshal(body, &db)
	if err != nil {
		return nil, err
	}

	return &db, nil
}

// ConfigDB calls the Squarescale API to update database options for a given project.
func (c *Client) ConfigDB(uuid string, payload *JSONObject) (taskId int, err error) {

	code, body, err := c.put("/projects/"+uuid+"/database", payload)
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
	default:
		return 0, unexpectedHTTPError(code, body)
	}

	var resp struct {
		DBTask int `json:"db_task"`
		Task   int `json:"task"`
	}

	err = json.Unmarshal(body, &resp)
	if err != nil {
		return 0, err
	}

	if resp.Task != 0 {
		return resp.Task, nil
	} else {
		return resp.DBTask, nil
	}
}
