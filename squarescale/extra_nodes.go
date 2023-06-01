package squarescale

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/squarescale/logger"
)

// ExtraNode describes a project container as returned by the SquareScale API
type ExtraNode struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	NodeType string `json:"node_type"`
	Zone     string `json:"zone"`
	Status   string `json:"status"`
}

// GetExtraNodes gets all the extraNodes attached to a Project
func (c *Client) GetExtraNodes(projectUUID string) ([]ExtraNode, error) {
	code, body, err := c.get("/projects/" + projectUUID + "/statefull_nodes")
	if err != nil {
		return []ExtraNode{}, err
	}

	switch code {
	case http.StatusOK:
	case http.StatusNotFound:
		return []ExtraNode{}, fmt.Errorf("Project '%s' does not exist", projectUUID)
	default:
		return []ExtraNode{}, unexpectedHTTPError(code, body)
	}

	var extraNodesByID []ExtraNode

	if err := json.Unmarshal(body, &extraNodesByID); err != nil {
		return []ExtraNode{}, err
	}

	return extraNodesByID, nil
}

// AddExtraNode add a new extra-node
func (c *Client) AddExtraNode(projectUUID string, name string, nodeType string, zone string) (ExtraNode, error) {
	var newExtraNode ExtraNode

	payload := JSONObject{
		"name":      name,
		"node_type": nodeType,
		"zone":      zone,
	}
	code, body, err := c.post("/projects/"+projectUUID+"/statefull_nodes", &payload)
	if err != nil {
		return newExtraNode, err
	}

	switch code {
	case http.StatusCreated:
	case http.StatusNotFound:
		return newExtraNode, fmt.Errorf("Project '%s' does not exist", projectUUID)
	case http.StatusConflict:
		return newExtraNode, fmt.Errorf("extra-node already exist on project '%s': %s", projectUUID, name)
	default:
		return newExtraNode, unexpectedHTTPError(code, body)
	}

	if err := json.Unmarshal(body, &newExtraNode); err != nil {
		return newExtraNode, err
	}

	return newExtraNode, nil
}

// GetExtraNodeInfo get information for a extra-node
func (c *Client) GetExtraNodeInfo(projectUUID string, name string) (ExtraNode, error) {
	extraNodes, err := c.GetExtraNodes(projectUUID)
	if err != nil {
		return ExtraNode{}, err
	}

	for _, extraNode := range extraNodes {
		if extraNode.Name == name {
			return extraNode, nil
		}
	}

	return ExtraNode{}, fmt.Errorf("extra-node '%s' not found for project '%s'", name, projectUUID)
}

// WaitExtraNode wait a new extra-node
func (c *Client) WaitExtraNode(projectUUID string, name string, timeToWait int64) (ExtraNode, error) {
	extraNode, err := c.GetExtraNodeInfo(projectUUID, name)
	if err != nil {
		return extraNode, err
	}

	for extraNode.Status != "provisionned" && err == nil {
		time.Sleep(time.Duration(timeToWait) * time.Second)
		extraNode, err = c.GetExtraNodeInfo(projectUUID, name)
		logger.Debug.Println("extraNode status update: ", extraNode.Name)
	}

	return extraNode, err
}

// DeleteExtraNode delete a existing extra-node
func (c *Client) DeleteExtraNode(projectUUID string, name string) error {
	code, body, err := c.delete("/projects/" + projectUUID + "/statefull_nodes/" + name)
	if err != nil {
		return err
	}

	switch code {
	case http.StatusOK:
	case http.StatusNotFound:
		if fmt.Sprintf("%s", body) == `{"error":"Couldn't find ExtraNode with [WHERE \"statefull_nodes\".\"cluster_id\" = $1 AND \"statefull_nodes\".\"name\" = $2]"}` {
			return fmt.Errorf("extra-node '%s' does not exist", name)
		}
		return fmt.Errorf("Project '%s' does not exist", projectUUID)
	case http.StatusBadRequest:
		return fmt.Errorf("Deploy probably in progress")
	default:
		return unexpectedHTTPError(code, body)
	}

	return nil
}

// BindVolumeOnExtraNode bind a volume to a extra-node
func (c *Client) BindVolumeOnExtraNode(projectUUID string, name string, volumeName string) error {
	var volumeToBind [1]string
	var volumeToUnbind [0]string

	volumeToBind[0] = volumeName

	payload := JSONObject{
		"volumes_to_bind":   volumeToBind,
		"volumes_to_unbind": volumeToUnbind,
	}
	code, body, err := c.put("/projects/"+projectUUID+"/statefull_nodes/"+name, &payload)
	if err != nil {
		return err
	}

	switch code {
	case http.StatusOK:
		return nil
	case http.StatusNotFound:
		if fmt.Sprintf("%s", body) == `{"error":"Couldn't find Volume with [WHERE \"volumes\".\"cluster_id\" = $1 AND \"volumes\".\"name\" = $2]"}` {
			return fmt.Errorf("Volume '%s' does not exist", volumeName)
		}
		return fmt.Errorf("Project '%s' does not exist", projectUUID)
	case http.StatusBadRequest:
		return fmt.Errorf("Volume %s already bound with %s", volumeName, name)
	default:
		return unexpectedHTTPError(code, body)
	}
}
