package squarescale

import (
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
