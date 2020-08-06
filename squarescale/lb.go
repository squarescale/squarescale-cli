package squarescale

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// DisableLB asks the Squarescale service to deactivate the load balancer.
func (c *Client) DisableLB(projectUUID string) error {
	payload := &JSONObject{
		"active": false,
	}

	return c.updateLBConfig(projectUUID, payload)
}

// ConfigLB sets the load balancer configuration for a given project.
func (c *Client) ConfigLB(projectUUID string, cert string, cert_chain string, secret_key string) error {
	payload := &JSONObject{
		"active":            true,
		"certificate_body":  cert,
		"certificate_chain": cert_chain,
		"secret_key":        secret_key,
		"https":             len(cert) > 0,
	}
	return c.updateLBConfig(projectUUID, payload)
}

func (c *Client) updateLBConfig(projectUUID string, payload *JSONObject) error {
	code, body, err := c.put("/projects/"+projectUUID+"/load_balancers/1", payload) //TODO change it when support more many LB.
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
	PublicUrl string `json:"public_url"`
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
