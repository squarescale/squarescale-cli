package squarescale

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// LoadBalancer describes a load balancer as returned by the Squarescale API
type LoadBalancer struct {
	ID               int64  `json:"id"`
	Active           bool   `json:"active"`
	CertificateBody  string `json:"certificate_body"`
	CertificateChain string `json:"certificate_chain"`
	HTTPS            bool   `json:"https"`
	PublicURL        string `json:"public_url"`
}

// LoadBalancerGet get details for load balancer
func (c *Client) LoadBalancerGet(projectUUID string) ([]LoadBalancer, error) {
	code, body, err := c.get("/projects/" + projectUUID + "/load_balancers")

	if err != nil {
		return []LoadBalancer{}, err
	}

	switch code {
	case http.StatusOK:
	case http.StatusNotFound:
		return []LoadBalancer{}, fmt.Errorf("Project '%s' not found", projectUUID)
	default:
		return []LoadBalancer{}, unexpectedHTTPError(code, body)
	}

	var loadBalancers []LoadBalancer

	err = json.Unmarshal(body, &loadBalancers)
	if err != nil {
		return []LoadBalancer{}, err
	}

	return loadBalancers, nil
}

// LoadBalancerEnable sets the load balancer configuration for a given project.
func (c *Client) LoadBalancerEnable(projectUUID string, loadBalancerID int64, cert string, certChain string, secretKey string) error {
	payload := &JSONObject{
		"active":            true,
		"certificate_body":  cert,
		"certificate_chain": certChain,
		"secret_key":        secretKey,
		"https":             len(cert) > 0,
	}

	URL := fmt.Sprintf("/projects/%s/load_balancers/%d", projectUUID, loadBalancerID)
	code, body, err := c.put(URL, payload)
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

// LoadBalancerDisable sets the load balancer configuration for a given project.
func (c *Client) LoadBalancerDisable(projectUUID string, loadBalancerID int64) error {
	payload := &JSONObject{
		"active": true,
		"https":  false,
	}

	URL := fmt.Sprintf("/projects/%s/load_balancers/%d", projectUUID, loadBalancerID)
	code, body, err := c.put(URL, payload)
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
