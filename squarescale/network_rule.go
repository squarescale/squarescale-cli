package squarescale

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type NetworkRule struct {
	Name             string `json:"name"`
	InternalPort     int    `json:"internal_port"`
	InternalProtocol string `json:"internal_protocol"`
	ExternalPort     int    `json:"external_port"`
	ExternalProtocol string `json:"external_protocol"`
	DomainExpression string `json:"domain_expression,omitempty"`
	PathPrefix       string `json:"path_prefix,omitempty"`
}

func (c *Client) ListNetworkRules(projectUUID, serviceName string) ([]NetworkRule, error) {

	code, body, err := c.get(fmt.Sprintf("/projects/%s/services/%s/service_network_rules", projectUUID, serviceName))
	if err != nil {
		return nil, err
	}

	switch code {
	case http.StatusOK:
	case http.StatusNotFound:
		return []NetworkRule{}, fmt.Errorf("Project '%s' and/or service container '%s' does not exist", projectUUID, serviceName)
	default:
		return []NetworkRule{}, unexpectedHTTPError(code, body)
	}

	var netRules []NetworkRule

	if err := json.Unmarshal(body, &netRules); err != nil {
		return []NetworkRule{}, err
	}

	return netRules, nil
}

func (c *Client) CreateNetworkRule(projectUUID, serviceName string, newRule NetworkRule) error {

	code, body, err := c.post(fmt.Sprintf("/projects/%s/services/%s/service_network_rules", projectUUID, serviceName), newRule)
	if err != nil {
		return err
	}

	if code != http.StatusCreated {
		return unexpectedHTTPError(code, body)
	}

	return nil
}

func (c *Client) DeleteNetworkRule(projectUUID, serviceName string, ruleToDelete string) error {

	code, body, err := c.delete(fmt.Sprintf("/projects/%s/services/%s/service_network_rules/%s", projectUUID, serviceName, ruleToDelete))
	if err != nil {
		return err
	}

	switch code {
	case http.StatusOK:
	case http.StatusNotFound:
		return fmt.Errorf("Project '%s' and/or service container '%s' or network rule '%s' does not exist", projectUUID, serviceName, ruleToDelete)
	default:
		return unexpectedHTTPError(code, body)
	}

	return nil
}
