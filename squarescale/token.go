package squarescale

import (
	"errors"
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
	case http.StatusBadRequest:
		return errors.New("update needed, api versions mismatch")
	default:
		return fmt.Errorf("token rejected")
	}
}
