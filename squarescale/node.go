package squarescale

import (
	"net/http"
)

// HasNodeSize is true if user can use node sizes
func (c *Client) HasNodeSize() bool {
	code, _, err := c.get("/infra/node/sizes")
	if err != nil {
		return false
	}

	if code == http.StatusForbidden {
		return false
	}

	return true
}
