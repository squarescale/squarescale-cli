package squarescale

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/squarescale/logger"
)

// StatefulNode describes a project container as returned by the Squarescale API
type StatefulNode struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	NodeType string `json:"node_type"`
	Zone     string `json:"zone"`
	Status   string `json:"status"`
}

// GetStatefulNodes gets all the statefulNodes attached to a Project
func (c *Client) GetStatefulNodes(projectUUID string) ([]StatefulNode, error) {
	code, body, err := c.get("/projects/" + projectUUID + "/statefull_nodes")
	if err != nil {
		return []StatefulNode{}, err
	}

	switch code {
	case http.StatusOK:
	case http.StatusNotFound:
		return []StatefulNode{}, fmt.Errorf("Project '%s' does not exist", projectUUID)
	default:
		return []StatefulNode{}, unexpectedHTTPError(code, body)
	}

	var statefulNodesByID []StatefulNode

	if err := json.Unmarshal(body, &statefulNodesByID); err != nil {
		return []StatefulNode{}, err
	}

	return statefulNodesByID, nil
}

// AddStatefulNode add a new stateful-node
func (c *Client) AddStatefulNode(projectUUID string, name string, nodeType string, zone string) (StatefulNode, error) {
	var newStatefulNode StatefulNode

	payload := JSONObject{
		"name":      name,
		"node_type": nodeType,
		"zone":      zone,
	}
	code, body, err := c.post("/projects/"+projectUUID+"/statefull_nodes", &payload)
	if err != nil {
		return newStatefulNode, err
	}

	switch code {
	case http.StatusCreated:
	case http.StatusNotFound:
		return newStatefulNode, fmt.Errorf("Project '%s' does not exist", projectUUID)
	case http.StatusConflict:
		return newStatefulNode, fmt.Errorf("Stateful-node already exist on project '%s': %s", projectUUID, name)
	default:
		return newStatefulNode, unexpectedHTTPError(code, body)
	}

	if err := json.Unmarshal(body, &newStatefulNode); err != nil {
		return newStatefulNode, err
	}

	return newStatefulNode, nil
}

// GetStatefulNodeInfo get information for a stateful-node
func (c *Client) GetStatefulNodeInfo(projectUUID string, name string) (StatefulNode, error) {
	statefulNodes, err := c.GetStatefulNodes(projectUUID)
	if err != nil {
		return StatefulNode{}, err
	}

	for _, statefulNode := range statefulNodes {
		if statefulNode.Name == name {
			return statefulNode, nil
		}
	}

	return StatefulNode{}, fmt.Errorf("Stateful-node '%s' not found for project '%s'", name, projectUUID)
}

// WaitStatefulNode wait a new stateful-node
func (c *Client) WaitStatefulNode(projectUUID string, name string, timeToWait int64) (StatefulNode, error) {
	statefulNode, err := c.GetStatefulNodeInfo(projectUUID, name)
	if err != nil {
		return statefulNode, err
	}

	for statefulNode.Status != "provisionned" && err == nil {
		time.Sleep(time.Duration(timeToWait) * time.Second)
		statefulNode, err = c.GetStatefulNodeInfo(projectUUID, name)
		logger.Debug.Println("statefulNode status update: ", statefulNode.Name)
	}

	return statefulNode, err
}

// DeleteStatefulNode delete a existing stateful-node
func (c *Client) DeleteStatefulNode(projectUUID string, name string) error {
	code, body, err := c.delete("/projects/" + projectUUID + "/statefull_nodes/" + name)
	if err != nil {
		return err
	}

	switch code {
	case http.StatusOK:
	case http.StatusNotFound:
		if fmt.Sprintf("%s", body) == `{"error":"Couldn't find StatefulNode with [WHERE \"statefull_nodes\".\"cluster_id\" = $1 AND \"statefull_nodes\".\"name\" = $2]"}` {
			return fmt.Errorf("Stateful-node '%s' does not exist", name)
		}
		return fmt.Errorf("Project '%s' does not exist", projectUUID)
	case http.StatusBadRequest:
		return fmt.Errorf("Deploy probably in progress")
	default:
		return unexpectedHTTPError(code, body)
	}

	return nil
}

// BindVolumeOnStatefulNode bind a volume to a stateful-node
func (c *Client) BindVolumeOnStatefulNode(projectUUID string, name string, volumeName string) error {
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
