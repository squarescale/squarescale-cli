package squarescale

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Collaborator describes a collaborator as returned by the SquareScale API
type Collaborator struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

// Organization describes a organization as returned by the SquareScale API
type Organization struct {
	Name          string         `json:"name"`
	Collaborators []Collaborator `json:"collaborators"`
	Projects      []Project      `json:"projects"`
	RootUser      struct {
		Email string `json:"email"`
	} `json:"root_user"`
}

// AddOrganization add organization
func (c *Client) AddOrganization(name, email string) error {
	payload := JSONObject{
		"name":          name,
		"contact_email": email,
	}

	url := fmt.Sprintf("/organizations")
	code, body, err := c.post(url, &payload)

	if err != nil {
		return err
	}

	switch code {
	case http.StatusCreated:
		return nil
	case http.StatusConflict:
		return fmt.Errorf("Organization already exist: %s", name)
	default:
		return unexpectedHTTPError(code, body)
	}
}

// GetOrganizationInfo gets the organization based on its name.
func (c *Client) GetOrganizationInfo(name string) (Organization, error) {
	url := fmt.Sprintf("/organizations/%s", name)
	code, body, err := c.get(url)

	if err != nil {
		return Organization{}, err
	}

	switch code {
	case http.StatusOK:
	case http.StatusNotFound:
		return Organization{}, fmt.Errorf("Organization '%s' not found", name)
	default:
		return Organization{}, unexpectedHTTPError(code, body)
	}

	var organizationByName Organization

	if err := json.Unmarshal(body, &organizationByName); err != nil {
		return Organization{}, err
	}

	return organizationByName, nil
}

// ListOrganizations asks the SquareScale service for available organizations.
func (c *Client) ListOrganizations() ([]Organization, error) {
	code, body, err := c.get("/organizations")

	if err != nil {
		return nil, err
	}

	if code != http.StatusOK {
		return nil, unexpectedHTTPError(code, body)
	}

	var organizationsJSON []Organization
	err = json.Unmarshal(body, &organizationsJSON)
	if err != nil {
		return nil, err
	}

	return organizationsJSON, nil
}

// DeleteOrganization delete organization based on its name.
func (c *Client) DeleteOrganization(name string) error {
	url := fmt.Sprintf("/organizations/%s", name)
	code, body, err := c.delete(url)

	if err != nil {
		return err
	}

	switch code {
	case http.StatusOK:
		return nil
	case http.StatusNotFound:
		return fmt.Errorf("No organization found for name: %s", name)
	default:
		return unexpectedHTTPError(code, body)
	}
}
