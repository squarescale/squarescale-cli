package squarescale

import (
	"encoding/json"
	"net/http"
)

// NodeSizes allow to display and validate client side the
// node size the user wants
type NodeSizes interface {
	CheckID(size, infraType string) bool
	ListIds(infraType string) []string
}

type nodeSizes struct {
	SingleNode       nodeSizesWithAdditional `json:"single_node"`
	HighAvailability nodeSizesWithAdditional `json:"high_availability"`
}

type nodeSizesWithAdditional struct {
	Default    []string `json:"default"`
	Additional []string `json:"additional"`
}

func (ns *nodeSizes) allSizes() map[string]map[string]bool {
	res := make(map[string]map[string]bool)
	singleNode := make(map[string]bool)
	ha := make(map[string]bool)
	for _, v := range ns.SingleNode.Default {
		singleNode[v] = true
	}
	for _, v := range ns.HighAvailability.Default {
		ha[v] = true
	}
	for _, v := range ns.SingleNode.Additional {
		singleNode[v] = true
	}
	for _, v := range ns.HighAvailability.Additional {
		ha[v] = true
	}
	res["single-node"] = singleNode
	res["high-availability"] = ha

	return res
}

func (ns *nodeSizes) CheckID(size, infraType string) bool {
	_, ok := ns.allSizes()[infraType][size]
	return ok
}

func (ns *nodeSizes) ListIds(infraType string) []string {
	var res []string
	for id := range ns.allSizes()[infraType] {
		res = append(res, id)
	}
	return res
}

// GetClusterNodeSizes asks the Squarescale API for available cluster node sizes
func (c *Client) GetClusterNodeSizes() (NodeSizes, error) {
	code, body, err := c.get("/infra/node/sizes?all=true")
	if err != nil {
		return nil, err
	}

	if code != http.StatusOK {
		return nil, unexpectedHTTPError(code, body)
	}

	var sizes nodeSizes

	err = json.Unmarshal(body, &sizes)
	if err != nil {
		return nil, err
	}

	return &sizes, nil
}

// HasNodeSize is true if user can use node sizes
func (c *Client) HasNodeSize() bool {
	code, _, err := c.get("/infra/node/sizes?all=true")
	if err != nil {
		return false
	}

	if code == http.StatusForbidden {
		return false
	}

	return true
}
