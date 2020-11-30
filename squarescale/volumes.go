package squarescale

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/squarescale/logger"
)

// Volume describes a project container as returned by the SquareScale API
type Volume struct {
	ID               int    `json:"id"`
	Name             string `json:"name"`
	Size             int    `json:"size"`
	Type             string `json:"type"`
	Zone             string `json:"zone"`
	StatefulNodeName string `json:"statefull_node_name"`
	Status           string `json:"status"`
}

// VolumeToBind describe a volume bind to a container
type VolumeToBind struct {
	Name       string `json:"volume_name"`
	MountPoint string `json:"mount_point"`
	ReadOnly   bool   `json:"read_only"`
}

// GetVolumes gets all the volumes attached to a Project
func (c *Client) GetVolumes(projectUUID string) ([]Volume, error) {
	code, body, err := c.get("/projects/" + projectUUID + "/volumes")
	if err != nil {
		return []Volume{}, err
	}

	switch code {
	case http.StatusOK:
	case http.StatusNotFound:
		return []Volume{}, fmt.Errorf("Project '%s' does not exist", projectUUID)
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
func (c *Client) GetVolumeInfo(projectUUID, name string) (Volume, error) {
	volumes, err := c.GetVolumes(projectUUID)
	if err != nil {
		return Volume{}, err
	}

	for _, volume := range volumes {
		if volume.Name == name {
			return volume, nil
		}
	}

	return Volume{}, fmt.Errorf("Volume '%s' not found for project '%s'", name, projectUUID)
}

// WaitVolume wait the volume of a project based on its name.
func (c *Client) WaitVolume(projectUUID, name string, timeToWait int64) (Volume, error) {
	volume, err := c.GetVolumeInfo(projectUUID, name)
	if err != nil {
		return volume, err
	}

	for volume.Status != "provisionned" && err == nil {
		time.Sleep(time.Duration(timeToWait) * time.Second)
		volume, err = c.GetVolumeInfo(projectUUID, name)
		logger.Debug.Println("volume status update: ", volume.Name)
	}

	return volume, err
}

// DeleteVolume delete volume of a project based on its name.
func (c *Client) DeleteVolume(projectUUID string, volume string) error {
	url := fmt.Sprintf("/projects/%s/volumes/%s", projectUUID, volume)
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
func (c *Client) AddVolume(projectUUID string, name string, size int, volumeType string, zone string) error {
	payload := JSONObject{
		"name": name,
		"size": size,
		"type": volumeType,
		"zone": zone,
	}

	url := fmt.Sprintf("/projects/%s/volumes", projectUUID)
	code, body, err := c.post(url, &payload)
	if err != nil {
		return err
	}

	switch code {
	case http.StatusCreated:
		return nil
	case http.StatusConflict:
		return fmt.Errorf("Volume already exist on project '%s': %s", projectUUID, name)
	case http.StatusNotFound:
		return unexpectedHTTPError(code, body)
	default:
		return unexpectedHTTPError(code, body)
	}
}
