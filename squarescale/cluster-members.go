package squarescale

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type ClusterMember struct {
	ID              int             `json:"id"`
	Name            string          `json:"name"`
	PublicIP        string          `json:"public_ip"`
	PrivateIP       string          `json:"private_ip"`
	Status          string          `json:"nomad_status"`
	SchedulingGroup SchedulingGroup `json:"scheduling_group"`
}

// GetClusterMembers gets all the cluster members attached to a Project
func (c *Client) GetClusterMembers(projectUUID string) ([]ClusterMember, error) {
	code, body, err := c.get("/projects/" + projectUUID + "/cluster_members")
	if err != nil {
		return []ClusterMember{}, err
	}

	switch code {
	case http.StatusOK:
	case http.StatusNotFound:
		return []ClusterMember{}, fmt.Errorf("Project '%s' does not exist", projectUUID)
	default:
		return []ClusterMember{}, unexpectedHTTPError(code, body)
	}

	var clusterMembersByID []ClusterMember

	if err := json.Unmarshal(body, &clusterMembersByID); err != nil {
		return []ClusterMember{}, err
	}

	return clusterMembersByID, nil
}

// GetClusterMemberInfo get information for a cluster member
func (c *Client) GetClusterMemberInfo(projectUUID string, name string) (ClusterMember, error) {
	clusterMembers, err := c.GetClusterMembers(projectUUID)
	if err != nil {
		return ClusterMember{}, err
	}

	for _, clusterMember := range clusterMembers {
		if clusterMember.Name == name {
			return clusterMember, nil
		}
	}

	return ClusterMember{}, fmt.Errorf("Cluster node '%s' not found for project '%s'", name, projectUUID)
}

// ConfigClusterMember set information for a cluster member
func (c *Client) ConfigClusterMember(projectUUID string, clusterMember ClusterMember, newClusterMember ClusterMember) error {
	payload := &JSONObject{"scheduling_group": newClusterMember.SchedulingGroup.ID}

	code, body, err := c.put(fmt.Sprintf("/projects/%s/cluster_members/%d", projectUUID, clusterMember.ID), payload)
	if err != nil {
		return err
	}

	switch code {
	case http.StatusOK:
		return nil
	case http.StatusNotFound:
		return fmt.Errorf("Project '%s' not found", projectUUID)
	default:
		return unexpectedHTTPError(code, body)
	}
}
