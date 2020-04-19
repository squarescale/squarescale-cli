package squarescale

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// StatefullNode describes a project container as returned by the Squarescale API
type StatefullNode struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Size     int    `json:"size"`
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
