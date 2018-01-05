package squarescale

import (
	"encoding/json"
	"fmt"
	"net/http"
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
	Size uint `json:"cluster_size"`
}

func (cluster *ClusterConfig) Update(other ClusterConfig) {
	if other.Size != 0 {
		cluster.Size = other.Size
	}
}

func (cluster *ClusterConfig) ConfigSettings() jsonObject {
	clusterSettings := jsonObject{
		"desired_size": cluster.Size,
	}
	return clusterSettings
}

func (cluster *ClusterConfig) String() string {
	return fmt.Sprintf("%d nodes", cluster.Size)
}
