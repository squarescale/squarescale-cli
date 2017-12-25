package squarescale

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// DbSizes allow to display and validate client side the
// database the user wants
type DbSizes interface {
	ListHuman() []string
	CheckID(size, infraType string) bool
	ListIds(infraType string) []string
}

type basicDbSizes struct {
	Sizes map[string]string `json:"sizes"`
}

type fullDbSizes struct {
	SingleNode       dbSizesWithAdditional `json:"single_node"`
	HighAvailability dbSizesWithAdditional `json:"high_availability"`
}

type dbSizesWithAdditional struct {
	Default    []string `json:"default"`
	Additional []string `json:"additional"`
}

func (s *basicDbSizes) ListHuman() []string {
	var res []string
	for k, v := range s.Sizes {
		res = append(res, k+": "+v)
	}
	return res
}

func (s *basicDbSizes) CheckID(size, _ string) bool {
	for k := range s.Sizes {
		if k == size {
			return true
		}
	}
	return false
}

func (s *basicDbSizes) ListIds(_ string) []string {
	var res []string
	for k := range s.Sizes {
		res = append(res, k)
	}
	return res
}

func (fs *fullDbSizes) ListHuman() []string {
	var res []string
	res = append(res, "Single Node infrastructure")
	for _, v := range fs.SingleNode.Default {
		res = append(res, "\t"+v)
	}
	for _, v := range fs.SingleNode.Additional {
		res = append(res, "\t["+v+"]")
	}
	res = append(res, "High Availability infrastructure")
	for _, v := range fs.HighAvailability.Default {
		res = append(res, "\t"+v)
	}
	for _, v := range fs.HighAvailability.Additional {
		res = append(res, "\t["+v+"]")
	}

	return res
}

func (fs *fullDbSizes) allSizes() map[string]map[string]bool {
	res := make(map[string]map[string]bool)
	singleNode := make(map[string]bool)
	ha := make(map[string]bool)
	for _, v := range fs.SingleNode.Default {
		singleNode[v] = true
	}
	for _, v := range fs.HighAvailability.Default {
		ha[v] = true
	}
	for _, v := range fs.SingleNode.Additional {
		singleNode[v] = true
	}
	for _, v := range fs.HighAvailability.Additional {
		ha[v] = true
	}
	res["single-node"] = singleNode
	res["high-availability"] = ha

	return res
}

func (fs *fullDbSizes) CheckID(size, infraType string) bool {
	_, ok := fs.allSizes()[infraType][size]
	return ok
}

func (fs *fullDbSizes) ListIds(infraType string) []string {
	var res []string
	for id := range fs.allSizes()[infraType] {
		res = append(res, id)
	}
	return res
}

// GetAvailableDBSizes returns all the database instances available for use in Squarescale.
func (c *Client) GetAvailableDBSizes() (DbSizes, error) {
	code, body, err := c.get("/db/sizes")
	if err != nil {
		return nil, err
	}

	if code != http.StatusOK {
		return nil, unexpectedHTTPError(code, body)
	}

	var sizeList basicDbSizes
	sizeList.Sizes = map[string]string{}

	err = json.Unmarshal(body, &sizeList)
	if err != nil {
		return nil, err
	}

	return &sizeList, nil
}

func (c *Client) HasNewDB() bool {
	code, _, err := c.get("/db/sizes/infra")
	if err != nil {
		return false
	}

	if code == http.StatusForbidden {
		return false
	}

	return true
}

func (c *Client) GetDBSizes() (DbSizes, error) {
	code, body, err := c.get("/db/sizes/infra")
	if err != nil {
		return nil, err
	}

	if code == http.StatusForbidden {
		// not authorized to see new db sizes? Fallback to old db sizes
		return c.GetAvailableDBSizes()
	}
	if code != http.StatusOK {
		return nil, unexpectedHTTPError(code, body)
	}

	var sizes fullDbSizes

	err = json.Unmarshal(body, &sizes)
	if err != nil {
		return nil, err
	}

	return &sizes, nil
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
