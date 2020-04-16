package squarescale

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Volume describes a project container as returned by the Squarescale API
type Volume struct {
	Name              string `json:"name"`
	Size              int    `json:"size"`
	Type              string `json:"type"`
	Zone              string `json:"zone"`
	StatefullNodeName string `json:"statefull_node_name"`
	Status            string `json:"status"`
}

// GetVolume gets all the volumes attached to a Project
func (c *Client) GetVolumes(project string) ([]Volume, error) {
	code, body, err := c.get("/projects/" + project + "/volumes")
	if err != nil {
		return []Volume{}, err
	}

	switch code {
	case http.StatusOK:
	case http.StatusNotFound:
		return []Volume{}, fmt.Errorf("Project '%s' does not exist", project)
	default:
		return []Volume{}, unexpectedHTTPError(code, body)
	}

	var volumesByID []Volume

	if err := json.Unmarshal(body, &volumesByID); err != nil {
		return []Volume{}, err
	}

	return volumesByID, nil

}
