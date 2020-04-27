package squarescale

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/squarescale/logger"
)

// StatefullNode describes a project container as returned by the Squarescale API
type StatefullNode struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	NodeType string `json:"node_type"`
	Zone     string `json:"zone"`
	Status   string `json:"status"`
}

// GetStatefullNodes gets all the statefullNodes attached to a Project
func (c *Client) GetStatefullNodes(project string) ([]StatefullNode, error) {
	code, body, err := c.get("/projects/" + project + "/statefull_nodes")
	if err != nil {
		return []StatefullNode{}, err
	}

	switch code {
	case http.StatusOK:
	case http.StatusNotFound:
		return []StatefullNode{}, fmt.Errorf("Project '%s' does not exist", project)
	default:
		return []StatefullNode{}, unexpectedHTTPError(code, body)
	}

	var statefullNodesByID []StatefullNode

	if err := json.Unmarshal(body, &statefullNodesByID); err != nil {
		return []StatefullNode{}, err
	}

	return statefullNodesByID, nil
}

// AddStatefullNode add a new statefull node
func (c *Client) AddStatefullNode(project string, name string, nodeType string, zone string) (StatefullNode, error) {
	var newStatefullNode StatefullNode

	payload := JSONObject{
		"name":      name,
		"node_type": nodeType,
		"zone":      zone,
	}
	code, body, err := c.post("/projects/"+project+"/statefull_nodes", &payload)
	if err != nil {
		return newStatefullNode, err
	}

	switch code {
	case http.StatusCreated:
	case http.StatusNotFound:
		return newStatefullNode, fmt.Errorf("Project '%s' does not exist", project)
	case http.StatusConflict:
		return newStatefullNode, fmt.Errorf("Statefull node already exist on project '%s': %s", project, name)
	default:
		return newStatefullNode, unexpectedHTTPError(code, body)
	}

	if err := json.Unmarshal(body, &newStatefullNode); err != nil {
		return newStatefullNode, err
	}

	return newStatefullNode, nil
}

// GetStatefullNodeInfo get information for a statefull node
func (c *Client) GetStatefullNodeInfo(project string, name string) (StatefullNode, error) {
	statefullNodes, err := c.GetStatefullNodes(project)
	if err != nil {
		return StatefullNode{}, err
	}

	for _, statefullNode := range statefullNodes {
		if statefullNode.Name == name {
			return statefullNode, nil
		}
	}

	return StatefullNode{}, fmt.Errorf("Statefull node '%s' not found for project '%s'", name, project)
}

// WaitStatefullNode wait a new statefull node
func (c *Client) WaitStatefullNode(project string, name string, timeToWait int64) (StatefullNode, error) {
	statefullNode, err := c.GetStatefullNodeInfo(project, name)
	if err != nil {
		return statefullNode, err
	}

	logger.Info.Println("wait for statefullNode : ", statefullNode.Name)

	for statefullNode.Status != "provisionned" && err == nil {
		time.Sleep(time.Duration(timeToWait) * time.Second)
		statefullNode, err = c.GetStatefullNodeInfo(project, name)
		logger.Debug.Println("statefullNode status update: ", statefullNode.Name)
	}

	return statefullNode, err
}

// DeleteStatefullNode delete a existing statefull node
func (c *Client) DeleteStatefullNode(project string, name string) error {
	// fmt.Printf("/projects/" + project + "/statefull_nodes/" + name + "\n")
	code, body, err := c.delete("/projects/" + project + "/statefull_nodes/" + name)
	if err != nil {
		return err
	}

	switch code {
	case http.StatusOK:
	case http.StatusNotFound:
		if fmt.Sprintf("%s", body) == `{"error":"Couldn't find StatefullNode with [WHERE \"statefull_nodes\".\"cluster_id\" = $1 AND \"statefull_nodes\".\"name\" = $2]"}` {
			return fmt.Errorf("Statefull node '%s' does not exist", name)
		}
		return fmt.Errorf("Project '%s' does not exist", project)
	case http.StatusBadRequest:
		return fmt.Errorf("Deploy probably in progress")
	default:
		return unexpectedHTTPError(code, body)
	}

	return nil
}

// BindVolumeOnStatefullNode bind a volume to a statefull node
func (c *Client) BindVolumeOnStatefullNode(project string, name string, volumeName string) error {
	var volumeToBind [1]string
	var volumeToUnbind [0]string

	volumeToBind[0] = volumeName

	payload := JSONObject{
		"volumes_to_bind":   volumeToBind,
		"volumes_to_unbind": volumeToUnbind,
	}
	fmt.Printf("PUT /projects/%s/statefull_nodes/%s\n", project, name)
	fmt.Println(payload)
	code, body, err := c.put("/projects/"+project+"/statefull_nodes/"+name, &payload)
	if err != nil {
		return err
	}

	switch code {
	case http.StatusOK:
		fmt.Printf("http.StatusOK\n")
	case http.StatusNotFound:
		if fmt.Sprintf("%s", body) == `{"error":"Couldn't find Volume with [WHERE \"volumes\".\"cluster_id\" = $1 AND \"volumes\".\"name\" = $2]"}` {
			return fmt.Errorf("Volume '%s' does not exist", volumeName)
		}
		return fmt.Errorf("Project '%s' does not exist", project)
	case http.StatusBadRequest:
		fmt.Printf("http.StatusBadRequest\n")
		return fmt.Errorf("Volume %s already bound with %s", volumeName, name)
	default:
		return unexpectedHTTPError(code, body)
	}

	return nil
}
