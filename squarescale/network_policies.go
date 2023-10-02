package squarescale

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

type NetworkPolicyVersions struct {
	Version    string    `json:"version"`
	Name       string    `json:"name"`
	DeployedAt time.Time `json:"deployed_at",omit`
	Active     bool
	Status     string
}

type NetworkPolicy struct {
	loaded     bool
	NumRules   int
	StatusStr  string
	Version    string    `json:"version"`
	Name       string    `json:"name"`
	DeployedAt time.Time `json:"deployed_at"`
	CreatedAt  time.Time `json:"created_at"`
	Validated  bool      `json:"validated"`
	Status     int       `json:"status"`
	Policies   string    `json:"policies"`
}

var NETWORK_POLICY_DEFAULT_STATUS = "stand-by"

var NETWORK_POLICY_STATUS = []string{"pending", "ok", "error", "partial"}

func (p *NetworkPolicy) IsLoaded() bool {
	return p.loaded
}

func (p *NetworkPolicy) DumpYAMLRules() ([]byte, error) {
	var rules []interface{}
	err := json.Unmarshal([]byte(p.Policies), &rules)
	if err != nil {
		return nil, err
	}
	data, err := yaml.Marshal(rules)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (p *NetworkPolicy) ImportPoliciesFromFile(path string) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	var rules []interface{}

	if strings.Contains(path, ".yaml") {
		err = yaml.Unmarshal(data, &rules)
		data, err = json.Marshal(&rules)
		if err != nil {
			return errors.New(fmt.Sprintf("Error when unmarshalling network policy file: %s", err))
		}
	} else {
		err = json.Unmarshal(data, &rules)
		if err != nil {
			return errors.New(fmt.Sprintf("Error when unmarshalling network policy file: %s", err))
		}
	}

	p.Policies = string(data)

	return nil
}

func (p *NetworkPolicy) String() string {
	deployed := "never deployed"
	created_at := fmt.Sprintf("created at: %s", p.CreatedAt)
	if !p.DeployedAt.IsZero() {
		deployed = fmt.Sprintf("deployed at: %s", p.DeployedAt)
	}
	return fmt.Sprintf("version: %s\nstatus: %s\npolicy rules: %d\n%s\n%s\n",
		p.Version, p.StatusStr, p.NumRules, created_at, deployed,
	)
}

func (c *Client) GetNetworkPolicy(projectUUID, version string) (*NetworkPolicy, error) {
	var code int
	var body []byte
	var err error
	var current *NetworkPolicy

	if version != "" {
		code, body, err = c.get("/projects/" + projectUUID + "/network_policies/" + version)
		current, err = c.GetNetworkPolicy(projectUUID, "")
		if err != nil {
			return nil, err
		}
	} else {
		code, body, err = c.get("/projects/" + projectUUID + "/network_policy")
	}

	if err != nil {
		return nil, err
	}

	var policy NetworkPolicy
	switch code {
	case http.StatusOK:
	case http.StatusNotFound:
		return &policy, nil
	default:
		return nil, unexpectedHTTPError(code, body)
	}

	if string(body) == "" || string(body) == "[]" {
		return &policy, nil
	}

	err = json.Unmarshal(body, &policy)
	if err != nil {
		return nil, err
	}

	var rules []interface{}
	err = json.Unmarshal([]byte(policy.Policies), &rules)
	if err != nil {
		return nil, err
	}
	policy.NumRules = len(rules)
	policy.loaded = true

	if version == "" || (current != nil &&
		current.Version == policy.Version) {
		policy.StatusStr = NETWORK_POLICY_STATUS[policy.Status]
	} else {
		policy.StatusStr = NETWORK_POLICY_DEFAULT_STATUS
	}

	return &policy, nil
}

func (c *Client) DeployNetworkPolicy(projectUUID, version string) error {
	code, body, err := c.put("/projects/"+projectUUID+"/network_policies/"+version+"/deploy", nil)
	if err != nil {
		return err
	}
	switch code {
	case http.StatusOK:
	case http.StatusNotFound:
		return errors.New("version not found")
	default:
		return unexpectedHTTPError(code, body)
	}
	return nil
}

func (c *Client) DeleteNetworkPolicy(projectUUID, version string) error {
	code, body, err := c.delete("/projects/" + projectUUID + "/network_policies/" + version)
	if err != nil {
		return err
	}
	switch code {
	case http.StatusNoContent:
	case http.StatusNotFound:
		return errors.New("version not found")
	default:
		return unexpectedHTTPError(code, body)
	}
	return nil
}

func (c *Client) ListNetworkPolicies(projectUUID string) (*NetworkPolicy, []NetworkPolicyVersions, error) {
	code, body, err := c.get("/projects/" + projectUUID + "/network_policies")
	if err != nil {
		return nil, nil, err
	}

	var versions []NetworkPolicyVersions

	switch code {
	case http.StatusOK:
	case http.StatusNotFound:
		return nil, nil, nil
	default:
		return nil, nil, unexpectedHTTPError(code, body)
	}

	if string(body) == "" || string(body) == "[]" {
		return nil, nil, nil
	}

	err = json.Unmarshal(body, &versions)
	if err != nil {
		return nil, nil, err
	}

	currentPolicy, err := c.GetNetworkPolicy(projectUUID, "")
	if err != nil {
		return nil, nil, err
	}

	return currentPolicy, versions, nil
}

func (c *Client) AddNetworkPolicy(projectUUID, name, rules string) (string, error) {
	payload := JSONObject{
		"name":         name,
		"policies":     rules,
		"project_uuid": projectUUID,
	}

	code, body, err := c.post("/projects/"+projectUUID+"/network_policies", &payload)
	if err != nil {
		return "", err
	}

	switch code {
	case http.StatusOK:
	default:
		return "", unexpectedHTTPError(code, body)
	}

	if string(body) == "" || string(body) == "[]" {
		return "unknown", nil
	}

	var policyVersion map[string]string

	err = json.Unmarshal(body, &policyVersion)
	if err != nil {
		return "", err
	}

	if policyVersion["version"] == "" {
		return "", errors.New("Unknown network policy version !")

	}
	return policyVersion["version"], nil
}
