package squarescale

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

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

// ConfigCluster calls the Squarescale API to update the cluster size for a given project.
func (c *Client) ConfigCluster(project string, cluster *ClusterConfig) (taskid int, err error) {
	payload := &JSONObject{
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

// ProjectCreationSettings returns a JSON representation of the ClusterConfig as
// expected by the API for the creation of a cluster.
func (cluster *ClusterConfig) ProjectCreationSettings() JSONObject {
	clusterSettings := JSONObject{
		"infra_type": getInfraTypeEnumValue(cluster.InfraType),
	}
	if cluster.NodeSize != "" {
		clusterSettings["node_size"] = cluster.NodeSize
	}
	return clusterSettings
}

// ConfigSettings returns a JSON representation of the ClusterConfig as
// expected by the API for the update of a cluster's settings.
func (cluster *ClusterConfig) ConfigSettings() JSONObject {
	clusterSettings := JSONObject{
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
