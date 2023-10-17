package squarescale

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/squarescale/logger"
)

type ExternalNode struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	PublicIP       string `json:"public_ip"`
	Status         string `json:"status"`
	PrivateNetwork string `json:"private_network"`
}

// GetExternalNodes gets all the external nodes attached to a Project
func (c *Client) GetExternalNodes(projectUUID string) ([]ExternalNode, error) {
	code, body, err := c.get("/projects/" + projectUUID + "/external_nodes")
	if err != nil {
		return []ExternalNode{}, err
	}

	switch code {
	case http.StatusOK:
	case http.StatusNotFound:
		return []ExternalNode{}, fmt.Errorf("Project '%s' does not exist", projectUUID)
	default:
		return []ExternalNode{}, unexpectedHTTPError(code, body)
	}

	var externalNodesByID []ExternalNode

	if err := json.Unmarshal(body, &externalNodesByID); err != nil {
		return []ExternalNode{}, err
	}

	return externalNodesByID, nil
}

// GetExternalNodeInfo get information for a external node
func (c *Client) GetExternalNodeInfo(projectUUID string, name string) (ExternalNode, error) {
	externalNodes, err := c.GetExternalNodes(projectUUID)
	if err != nil {
		return ExternalNode{}, err
	}

	for _, externalNode := range externalNodes {
		if externalNode.Name == name {
			return externalNode, nil
		}
	}

	return ExternalNode{}, fmt.Errorf("External node '%s' not found for project '%s'", name, projectUUID)
}

// AddExternalNode add a new scheduling group
func (c *Client) AddExternalNode(projectUUID string, name string, public_ip string) (ExternalNode, error) {
	var newExternalNode ExternalNode

	payload := JSONObject{
		"name":      name,
		"public_ip": public_ip,
	}
	code, body, err := c.post("/projects/"+projectUUID+"/external_nodes", &payload)
	if err != nil {
		return newExternalNode, err
	}

	switch code {
	case http.StatusCreated:
	case http.StatusNotFound:
		return newExternalNode, fmt.Errorf("Project '%s' does not exist", projectUUID)
	default:
		return newExternalNode, unexpectedHTTPError(code, body)
	}

	if err := json.Unmarshal(body, &newExternalNode); err != nil {
		return newExternalNode, err
	}

	return newExternalNode, nil
}

// TODO: potentially add a real timeout on this wait condition
// WaitExternalNode wait a new external-node to reach one status
func (c *Client) WaitExternalNode(projectUUID string, name string, timeToWait int64, targetStatuses []string) (ExternalNode, error) {
	externalNode, err := c.GetExternalNodeInfo(projectUUID, name)
	if err != nil {
		return externalNode, err
	}

	if len(targetStatuses) > 0 {
		found := false
		logger.Debug.Println("externalNode status: ", externalNode.Name, " Status: ", externalNode.Status)
		for err == nil && !found {
			for _, v := range targetStatuses {
				if externalNode.Status == v {
					found = true
					break
				}
			}
			if !found {
				time.Sleep(time.Duration(timeToWait) * time.Second)
				externalNode, err = c.GetExternalNodeInfo(projectUUID, name)
				logger.Debug.Println("externalNode status update: ", externalNode.Name, " Status: ", externalNode.Status)
			}
		}
	}

	return externalNode, err
}

// DownloadConfigExternalNode download service configuration file(s) for a external node
func (c *Client) DownloadConfigExternalNode(projectUUID, name, configName string) (error) {
	externalNodes, err := c.GetExternalNodes(projectUUID)
	if err != nil {
		return err
	}

	for _, externalNode := range externalNodes {
		if externalNode.Name == name {
			for _, cfg := range []string{"openvpn", "consul", "nomad"} {
				if configName == "all" || configName == cfg {
					code, err := c.download(fmt.Sprintf("/projects/%s/external_nodes/%d/%s_client", projectUUID, externalNode.ID, cfg), name)
					if err != nil {
						return err
					}

					switch code {
					case http.StatusOK:
					case http.StatusNotFound:
						return fmt.Errorf("Error retrieving service configuration %s for external node %s in project '%s'", configName, name, projectUUID)
					default:
						return unexpectedHTTPError(code, []byte{})
					}
				}
			}
			return nil
		}
	}
	return fmt.Errorf("Unable to find external node '%s' in project '%s'", name, projectUUID)
}
