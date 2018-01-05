package squarescale

import (
	"encoding/json"
	"net/http"
)

// NodeSizes allow to display and validate client side the
// node size the user wants
type NodeSizes interface {
	CheckSize(size, infraType string) bool
	ListSizes(infraType string) []string
}

type nodeSizesResponse struct {
	SingleNode       nodeSizesWithAdditional `json:"single_node"`
	HighAvailability nodeSizesWithAdditional `json:"high_availability"`
}

type nodeSizesWithAdditional struct {
	Default    []string `json:"default"`
	Additional []string `json:"additional"`
}

type nodeSizes struct {
	ByInfraType map[string][]string
}

func (ns *nodeSizesResponse) flattenSizes() map[string][]string {
	res := make(map[string][]string)
	res["single-node"] = append(ns.SingleNode.Default, ns.SingleNode.Additional...)
	res["high-availability"] = append(ns.HighAvailability.Default, ns.HighAvailability.Additional...)

	return res
}

func (ns *nodeSizes) CheckSize(size, infraType string) bool {
	for _, v := range ns.ByInfraType[infraType] {
		if v == size {
			return true
		}
	}
	return false
}

func (ns *nodeSizes) ListSizes(infraType string) []string {
	return ns.ByInfraType[infraType]
}

// GetClusterNodeSizes asks the Squarescale API for available cluster node sizes
func (c *Client) GetClusterNodeSizes() (NodeSizes, error) {
	code, body, err := c.get("/infra/node/sizes?all=true")
	if err != nil {
		return nil, err
	}

	if code == http.StatusForbidden {
		return nil, nil
	}

	if code != http.StatusOK {
		return nil, unexpectedHTTPError(code, body)
	}

	var response nodeSizesResponse

	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	var sizes nodeSizes
	sizes.ByInfraType = response.flattenSizes()

	return &sizes, nil
}
