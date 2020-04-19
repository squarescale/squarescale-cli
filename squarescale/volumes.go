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

// WaitVolume wait the volume of a project based on its name.
func (c *Client) WaitVolume(project, name string) (Volume, error) {
	volume, err := c.GetVolumeInfo(project, name)
	if err != nil {
		return volume, err
	}

	logger.Info.Println("wait for volume : ", volume.Name)

	for volume.Status != "provisionned" && err == nil {
		time.Sleep(5 * time.Second)
		volume, err = c.GetVolumeInfo(project, name)
		logger.Debug.Println("volume status update: ", volume.Name)
	}

	return volume, err
}

// DeleteVolume delete volume of a project based on its name.
func (c *Client) DeleteVolume(project string, volume string) error {
	url := fmt.Sprintf("/projects/%s/volumes/%s", project, volume)
	code, body, err := c.delete(url)
	if err != nil {
		return err
	}

	switch code {
	case http.StatusOK:
		return nil
	case http.StatusNotFound:
		if fmt.Sprintf("%s", body) == `{"error":"Couldn't find Volume with [WHERE \"volumes\".\"cluster_id\" = $1 AND \"volumes\".\"name\" = $2]"}` {
			return fmt.Errorf("{\"error\":\"No volume found for name: %s\"}", volume)
		}
		return fmt.Errorf("%s", body)
	default:
		return unexpectedHTTPError(code, body)
	}
}

// AddVolume add volume of a project based on its name.
func (c *Client) AddVolume(project string, name string, size int, volumeType string, zone string) error {
	payload := JSONObject{
		"name": name,
		"size": size,
		"type": volumeType,
		"zone": zone,
	}

	url := fmt.Sprintf("/projects/%s/volumes", project)
	code, body, err := c.post(url, &payload)
	if err != nil {
		return err
	}

	switch code {
	case http.StatusCreated:
		return nil
	case http.StatusConflict:
		return fmt.Errorf("Volume already exist: %s", project)
	case http.StatusNotFound:
		return unexpectedHTTPError(code, body)
	default:
		return unexpectedHTTPError(code, body)
	}
}
