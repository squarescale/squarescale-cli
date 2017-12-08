package squarescale

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type DbSizes struct {
	Sizes map[string]string `json:"sizes"`
}

func (s *DbSizes) ListHuman() []string {
	var res []string
	for k, v := range s.Sizes {
		res = append(res, k+": "+v)
	}
	return res
}

func (s *DbSizes) CheckId(size string) bool {
	for k, _ := range s.Sizes {
		if k == size {
			return true
		}
	}
	return false
}

func (s *DbSizes) ListIds() []string {
	var res []string
	for k, _ := range s.Sizes {
		res = append(res, k)
	}
	return res
}

// GetAvailableDBInstances returns all the database instances available for use in Squarescale.
func (c *Client) GetAvailableDBSizes() (*DbSizes, error) {
	code, body, err := c.get("/db/sizes")
	if err != nil {
		return nil, err
	}

	if code != http.StatusOK {
		return nil, unexpectedHTTPError(code, body)
	}

	var sizeList DbSizes
	sizeList.Sizes = map[string]string{}

	err = json.Unmarshal(body, &sizeList)
	if err != nil {
		return nil, err
	}

	return &sizeList, nil
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

type DbConfig struct {
	Enabled bool   `json:"db_enabled"`
	Engine  string `json:"db_engine"`
	Size    string `json:"db_size"`
}

func (db *DbConfig) String() string {
	if db.Enabled {
		return db.Size + " " + db.Engine
	} else {
		return "none"
	}
}

func (db *DbConfig) Update(other DbConfig) {
	db.Enabled = other.Enabled
	if other.Size != "" {
		db.Size = other.Size
	}
	if other.Engine != "" {
		db.Engine = other.Engine
	}
}

func (db *DbConfig) ProjectCreationSettings() jsonObject {
	dbSettings := jsonObject{
		"enabled": db.Enabled,
	}
	if db.Engine != "" {
		dbSettings["engine"] = db.Engine
	}
	if db.Size != "" {
		dbSettings["size"] = db.Size
	}
	return dbSettings
}

func (db *DbConfig) ConfigSettings() jsonObject {
	dbSettings := jsonObject{
		"db_enabled": db.Enabled,
	}
	if db.Engine != "" {
		dbSettings["db_engine"] = db.Engine
	}
	if db.Size != "" {
		dbSettings["db_size"] = db.Size
	}
	return dbSettings
}

// GetDBConfig asks the Squarescale API for the database config of a project.
// Returns, in this order:
// - if the db is enabled
// - the db engine in use (string)
// - the db instance in use (string)
func (c *Client) GetDBConfig(project string) (*DbConfig, error) {
	code, body, err := c.get("/projects/" + project)
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

// ConfigDB calls the Squarescale API to update database scale options for a given project.
func (c *Client) ConfigDB(project string, db *DbConfig) (taskId int, err error) {
	payload := &jsonObject{
		"project": db.ConfigSettings(),
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
		return 0, fmt.Errorf("Invalid value for either database engine ('%s') or size ('%s')", db.Engine, db.Size)
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
