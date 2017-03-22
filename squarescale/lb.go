package squarescale

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

// DisableLB asks the Squarescale service to deactivate the load balancer.
func (c *Client) DisableLB(project string) (int, error) {
	payload := &jsonObject{
		"project": jsonObject{
			"load_balancer": false,
		},
	}

	return c.updateLBConfig(project, payload)
}

// ConfigLB sets the load balancer configuration for a given project.
func (c *Client) ConfigLB(project string, container, port int) (int, error) {
	payload := &jsonObject{
		"project": jsonObject{
			"load_balancer": true,
		},
		"web-container": strconv.Itoa(container),
		"container": jsonObject{
			strconv.Itoa(container): jsonObject{
				"web-port": strconv.Itoa(port),
			},
		},
	}

	return c.updateLBConfig(project, payload)
}

func (c *Client) updateLBConfig(project string, payload *jsonObject) (int, error) {
	code, body, err := c.post("/projects/"+project+"/web-ports", payload)
	if err != nil {
		return 0, err
	}

	switch code {
	case http.StatusOK:
		break
	case http.StatusNotFound:
		return 0, fmt.Errorf("Project '%s' not found", project)
	default:
		return 0, unexpectedHTTPError(code, body)
	}

	var response struct {
		Task int `json:"task"`
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		return 0, err
	}

	return response.Task, nil
}

// LoadBalancerEnabled asks if the project load balancer is enabled.
func (c *Client) LoadBalancerEnabled(project string) (bool, error) {
	code, body, err := c.get("/projects/" + project)
	if err != nil {
		return false, err
	}

	switch code {
	case http.StatusOK:
	case http.StatusNotFound:
		return false, fmt.Errorf("Project '%s' not found", project)
	default:
		return false, unexpectedHTTPError(code, body)
	}

	var response struct {
		Enabled bool `json:"load_balancer"`
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		return false, err
	}

	return response.Enabled, nil
}
