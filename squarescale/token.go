package squarescale

import (
	"fmt"
	"net/http"
)

// ValidateToken asks Squarescale service for token validity. Returns nil if user is authorized.
func ValidateToken(sqscURL, token string) error {
	res, err := doRequest(SqscRequest{
		Method: "GET",
		URL:    sqscURL + "/status",
		Token:  token,
	})

	if err != nil {
		return err
	}

	if res.Code == http.StatusOK {
		return nil
	}

	return fmt.Errorf("Token validation failed on Squarescale services")
}
