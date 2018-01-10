package squarescale

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

// NodeSizes allow to display and validate client side the
// node size the user wants
type NodeSizes interface {
	CheckSize(size, infraType string) (bool, error)
	ListSizes(infraType string) ([]string, error)
	ListHuman() []string
}

type nodeSizes struct {
	SingleNode       nodeSizesWithAdditional `json:"single_node"`
	HighAvailability nodeSizesWithAdditional `json:"high_availability"`
	ByInfraType      map[string][]string
}

type nodeSizesWithAdditional struct {
	Default    []string `json:"default"`
	Additional []string `json:"additional"`
}

func (ns *nodeSizes) CheckSize(size, infraType string) (bool, error) {
	if ns.ByInfraType == nil {
		return false, errors.New("Node sizes were not flattened before calling CheckSize")
	}

	for _, v := range ns.ByInfraType[infraType] {
		if v == size {
			return true, nil
		}
	}
	return false, nil
}

func (ns *nodeSizes) ListSizes(infraType string) ([]string, error) {
	if ns.ByInfraType == nil {
		return nil, errors.New("Node sizes were not flattened before calling ListSizes")
	}

	return ns.ByInfraType[infraType], nil
}

func (ns *nodeSizes) ListHuman() []string {
	var res []string
	res = append(res, "Single Node infrastructure")
	for _, v := range ns.SingleNode.Default {
		res = append(res, "\t"+v)
	}
	for _, v := range ns.SingleNode.Additional {
		res = append(res, "\t["+v+"]")
	}
	res = append(res, "High Availability infrastructure")
	for _, v := range ns.HighAvailability.Default {
		res = append(res, "\t"+v)
	}
	for _, v := range ns.HighAvailability.Additional {
		res = append(res, "\t["+v+"]")
	}

	return res
}

func (ns *nodeSizes) flattenSizes() {
	ns.ByInfraType = make(map[string][]string)
	ns.ByInfraType["single-node"] = append(
		ns.SingleNode.Default,
		ns.SingleNode.Additional...)
	ns.ByInfraType["high-availability"] = append(
		ns.HighAvailability.Default,
		ns.HighAvailability.Additional...)
}

type emptyNodeSizes struct{}

func (ns *emptyNodeSizes) CheckSize(size, infraType string) (bool, error) {
	return false, nil
}

func (ns *emptyNodeSizes) ListSizes(infraType string) ([]string, error) {
	return []string{}, nil
}

func (ns *emptyNodeSizes) ListHuman() []string {
	return []string{"Your cluster will use nodes with 1 vCPU and 2 GiB of RAM"}
}

// GetClusterSize asks the Squarescale API for the cluster size of a project.
func (c *Client) GetClusterConfig(project string) (*ClusterConfig, error) {
	code, body, err := c.get("/projects/" + project)
	if err != nil {
		return nil, err
	}

	if code != http.StatusOK {
		return nil, unexpectedHTTPError(code, body)
	}

	var cluster ClusterConfig
	err = json.Unmarshal(body, &cluster)
	if err != nil {
		return nil, err
	}

	return &cluster, nil
}

// GetClusterNodeSizes asks the Squarescale API for available cluster node sizes
func (c *Client) GetClusterNodeSizes() (NodeSizes, error) {
	code, body, err := c.get("/infra/node/sizes?all=true")
	if err != nil {
		return nil, err
	}

	if code == http.StatusForbidden {
		return &emptyNodeSizes{}, nil
	}

	if code != http.StatusOK {
		return nil, unexpectedHTTPError(code, body)
	}

	var sizes nodeSizes

	err = json.Unmarshal(body, &sizes)
	if err != nil {
		return nil, err
	}

	sizes.flattenSizes()

	return &sizes, nil
}

// ConfigCluster calls the Squarescale API to update the cluster size for a given project.
func (c *Client) ConfigCluster(project string, cluster *ClusterConfig) (taskid int, err error) {
	payload := &jsonObject{
		"cluster": cluster.ConfigSettings(),
	}

	code, body, err := c.post("/projects/"+project+"/cluster", payload)
	if err != nil {
		return 0, err
	}

	switch code {
	case http.StatusAccepted:
		fallthrough
	case http.StatusOK:
		break
	case http.StatusNoContent:
		return 0, nil
	case http.StatusUnprocessableEntity:
		return 0, fmt.Errorf("Invalid value for cluster size ('%d')", cluster.Size)
	default:
		return 0, unexpectedHTTPError(code, body)
	}

	var resp struct {
		ResizeTask int `json:"resize_task"`
		Task       int `json:"task"`
	}

	err = json.Unmarshal(body, &resp)
	if err != nil {
		return 0, err
	}

	if resp.Task != 0 {
		return resp.Task, nil
	} else {
		return resp.ResizeTask, nil
	}
}

type ClusterConfig struct {
	InfraType string `json:"infra_type"`
	Size      uint   `json:"cluster_size"`
	NodeSize  string `json:"node_size"`
}

func (cluster *ClusterConfig) Update(other ClusterConfig) {
	if other.Size != 0 {
		cluster.Size = other.Size
	}
	if other.InfraType != "" {
		cluster.InfraType = other.InfraType
	}
	if other.NodeSize != "" {
		cluster.NodeSize = other.NodeSize
	}
}

func (cluster *ClusterConfig) ProjectCreationSettings() jsonObject {
	clusterSettings := jsonObject{
		"infra_type": getInfraTypeEnumValue(cluster.InfraType),
	}
	if cluster.NodeSize != "" {
		clusterSettings["node_size"] = cluster.NodeSize
	}
	return clusterSettings
}

func (cluster *ClusterConfig) ConfigSettings() jsonObject {
	clusterSettings := jsonObject{
		"desired_size": cluster.Size,
	}
	return clusterSettings
}

func (cluster *ClusterConfig) String() string {
	infraType := strings.ToTitle(strings.Replace(cluster.InfraType, "-", " ", 1))

	return fmt.Sprintf(
		"%s cluster with %d %s nodes",
		infraType,
		cluster.Size,
		cluster.NodeSize)
}

func getInfraTypeEnumValue(infraType string) string {
	return strings.Replace(infraType, "-", "_", 1)
}
