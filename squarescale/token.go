package squarescale

import (
	"fmt"
	"net/http"
)

// ValidateToken asks Squarescale service for token validity. Returns nil if user is authorized.
func ValidateToken(sqscURL, token string) error {
	url := sqscURL + "/status"
	code, _, err := get(url, token)
	if err != nil {
		return err
	}

	if code == http.StatusOK {
		return nil
	}

	return fmt.Errorf("token rejected")
}
