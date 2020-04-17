package squarescale

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/squarescale/logger"
)

// Volume describes a project container as returned by the Squarescale API
type Volume struct {
	ID                int    `json:"id"`
	Name              string `json:"name"`
	Size              int    `json:"size"`
	Type              string `json:"type"`
	Zone              string `json:"zone"`
	StatefullNodeName string `json:"statefull_node_name"`
	Status            string `json:"status"`
}

// GetVolumes gets all the volumes attached to a Project
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

// GetVolumeInfo gets the volume of a project based on its name.
func (c *Client) GetVolumeInfo(project, name string) (Volume, error) {
	volumes, err := c.GetVolumes(project)
	if err != nil {
		return Volume{}, err
	}

	for _, volume := range volumes {
		if volume.Name == name {
			return volume, nil
		}
	}

	return Volume{}, fmt.Errorf("Volume '%s' not found for project '%s'", name, project)
}

func (c *Client) DeleteVolume(project string, volume Volume) error {
	url := fmt.Sprintf("/projects/%s/volumes/%s", project, volume.Name)
	code, body, err := c.delete(url)
	if err != nil {
		return err
	}

	switch code {
	case http.StatusOK:
	case http.StatusNotFound:
		return fmt.Errorf("Volume '%s' does not exist", volume.Name)
	default:
		return unexpectedHTTPError(code, body)
	}

	return nil
}
