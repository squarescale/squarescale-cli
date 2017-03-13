package squarescale

import (
	"encoding/json"
	"net/http"
)

// GetClusterSize asks the Squarescale API for the cluster size of a project.
func (c *Client) GetClusterSize(project string) (uint, error) {
	code, body, err := c.get("/projects/" + project)
	if err != nil {
		return 0, err
	}

	if code != http.StatusOK {
		return 0, unexpectedHTTPError(code, body)
	}

	var resp struct {
		ClusterSize uint `json:"cluster_size"`
	}

	err = json.Unmarshal(body, &resp)
	if err != nil {
		return 0, err
	}

	return resp.ClusterSize, nil
}

// ConfigDB calls the Squarescale API to update database scale options for a given project.
func (c *Client) SetClusterSize(project string, clusterSize uint) error {
	payload := &jsonObject{
		"project": jsonObject{
			"cluster_size": clusterSize,
		},
	}

	code, body, err := c.post("/projects/"+project+"/cluster", payload)
	if err != nil {
		return err
	}

	switch code {
	case http.StatusOK:
		fallthrough
	case http.StatusNoContent:
		return nil
	default:
		return unexpectedHTTPError(code, body)
	}
}
