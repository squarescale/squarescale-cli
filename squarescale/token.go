package squarescale

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type User struct {
	UUID            string `json:"uid"`
	Email           string `json:"email"`
	FirstName       string `json:"first_name"`
	LastName        string `string:"last_name"`
	FullName        string `json:"name"`
	IsAdmin         bool   `json:"is_admin"`
}

// ValidateToken asks SquareScale service for token validity. Returns nil if user is authorized.
func (c *Client) ValidateToken() (*User, error) {
	code, body, err := c.get("/me")
	if err != nil {
		return nil, err
	}

	if code != http.StatusOK {
		return nil, fmt.Errorf("You're not logged in, please run login command")
		//return nil, unexpectedHTTPError(code, body)
	}

	var userJSON User
	err = json.Unmarshal(body, &userJSON)
	if err != nil {
		return nil, err
	}

	return &userJSON, nil
}
