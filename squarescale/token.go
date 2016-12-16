package squarescale

import (
	"fmt"
	"net/http"
)

// ValidateToken asks Squarescale service for token validity. Returns nil if user is authorized.
func (c *Client) ValidateToken() error {
	code, _, err := c.get("/status")
	if err != nil {
		return err
	}

	switch code {
	case http.StatusOK:
		return nil
	default:
		return fmt.Errorf("token rejected")
	}
}
