package squarescale

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

// DisableLB asks the Squarescale service to deactivate the load balancer.
func (c *Client) DisableLB(project string) error {
	payload := &jsonObject{
		"project": jsonObject{
			"load_balancer": false,
		},
	}

	return c.updateLBConfig(project, payload)
}

// ConfigLB sets the load balancer configuration for a given project.
func (c *Client) ConfigLB(project string, container, port int) error {
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

func (c *Client) updateLBConfig(project string, payload *jsonObject) error {
	code, _, err := c.post("/projects/"+project+"/web-ports", payload)
	if err != nil {
		return err
	}

	if code == http.StatusNotFound {
		return fmt.Errorf("Project '%s' not found", project)
	}

	if code != http.StatusOK {
		return unexpectedError(code)
	}

	return nil
}

// LoadBalancerEnabled asks if the project load balancer is enabled.
func (c *Client) LoadBalancerEnabled(project string) (bool, error) {
	code, body, err := c.get("/projects/" + project)
	if err != nil {
		return false, err
	}

	if code == http.StatusNotFound {
		return false, fmt.Errorf("Project '%s' not found", project)
	}

	if code != http.StatusOK {
		return false, unexpectedError(code)
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
