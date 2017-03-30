package squarescale

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

// DisableLB asks the Squarescale service to deactivate the load balancer.
func (c *Client) DisableLB(project string) (int, error) {
	payload := &jsonObject{
		"load_balancer": jsonObject{
			"active": false,
		},
		"project": jsonObject{
			"load_balancer": false,
		},
	}

	return c.updateLBConfig(project, payload)
}

// ConfigLB sets the load balancer configuration for a given project.
func (c *Client) ConfigLB(project string, container, port int, https bool, cert string, cert_chain []string, secret_key string) (int, error) {
	if cert_chain == nil {
		cert_chain = []string{}
	}
	payload := &jsonObject{
		"load_balancer": jsonObject{
			"active":            true,
			"container_id":      container,
			"https":             https,
			"certificate_body":  cert,
			"certificate_chain": cert_chain,
			"secret_key":        secret_key,
		},
		"containers": jsonObject{
			strconv.Itoa(container): jsonObject{
				"web_port": port,
			},
		},
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
	code, body, err := c.post("/projects/"+project+"/load_balancer", payload)
	if err != nil {
		return 0, err
	}

	switch code {
	case http.StatusOK:
		break
	case http.StatusBadRequest:
		break
	case http.StatusNotFound:
		return 0, fmt.Errorf("Project '%s' not found", project)
	default:
		return 0, unexpectedHTTPError(code, body)
	}

	var response struct {
		Ok              bool                           `json:"ok"`
		Task            int                            `json:"task"`
		Errors          map[string][]string            `json:"errors"`
		ContainerErrors map[string]map[string][]string `json:"container_errors"`
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		return 0, err
	}

	if !response.Ok {
		var errs []string
		for field, v := range response.Errors {
			for _, msg := range v {
				errs = append(errs, fmt.Sprintf("%s: %s", field, msg))
			}
		}
		for _, c := range response.ContainerErrors {
			for field, v := range c {
				for _, msg := range v {
					errs = append(errs, fmt.Sprintf("container %s: %s", field, msg))
				}
			}
		}
		return 0, fmt.Errorf("Bad request: %v", strings.Join(errs, "; "))
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
