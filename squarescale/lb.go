package squarescale

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

// DisableLB asks the Squarescale service to deactivate the load balancer.
func (c *Client) DisableLB(projectUUID string) (int, error) {
	payload := &JSONObject{
		"load_balancer": JSONObject{
			"active": false,
		},
		"project": JSONObject{
			"load_balancer": false,
		},
	}

	return c.updateLBConfig(projectUUID, payload)
}

// ConfigLB sets the load balancer configuration for a given project.
func (c *Client) ConfigLB(projectUUID string, container, port int, exprArg string, https bool, cert string, cert_chain []string, secret_key string) (int, error) {
	if cert_chain == nil {
		cert_chain = []string{}
	}
	payload := &JSONObject{
		"load_balancer": JSONObject{
			"active":            true,
			"container_id":      container,
			"https":             https,
			"certificate_body":  cert,
			"certificate_chain": cert_chain,
			"secret_key":        secret_key,
		},
		"containers": JSONObject{
			strconv.Itoa(container): JSONObject{
				"web_port":                 port,
				"custom_domain_expression": exprArg,
			},
		},
		"project": JSONObject{
			"load_balancer": true,
		},
		"web-container": strconv.Itoa(container),
		"container": JSONObject{
			strconv.Itoa(container): JSONObject{
				"web-port": strconv.Itoa(port),
			},
		},
	}

	return c.updateLBConfig(projectUUID, payload)
}

func (c *Client) updateLBConfig(projectUUID string, payload *JSONObject) (int, error) {
	code, body, err := c.put("/projects/"+projectUUID+"/load_balancers/1", payload) //TODO change it when support more many LB.
	if err != nil {
		return 0, err
	}

	switch code {
	case http.StatusOK:
		break
	case http.StatusNotFound:
		return 0, fmt.Errorf("Project '%s' not found", projectUUID)
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
func (c *Client) LoadBalancerEnabled(projectUUID string) (bool, error) {
	code, body, err := c.get("/projects/" + projectUUID)
	if err != nil {
		return false, err
	}

	switch code {
	case http.StatusOK:
	case http.StatusNotFound:
		return false, fmt.Errorf("Project '%s' not found", projectUUID)
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

type LoadBalancer struct {
	Port           int    `json:"Port"`
	ExposedService string `json:"exposed_service"`
	PublicUrl      string `json:"public_url"`
}

func (c *Client) LoadBalancerList(projectUUID string) ([]LoadBalancer, error) {
	code, body, err := c.get("/projects/" + projectUUID + "/load_balancers")
	if err != nil {
		return nil, err
	}

	switch code {
	case http.StatusOK:
	case http.StatusNotFound:
		return nil, fmt.Errorf("Project '%s' not found", projectUUID)
	default:
		return nil, unexpectedHTTPError(code, body)
	}

	var response []LoadBalancer

	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	return response, nil
}
