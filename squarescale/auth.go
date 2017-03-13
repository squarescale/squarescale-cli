package squarescale

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type authTokenResponse struct {
	AuthToken string `json:"auth_token"`
	Error     string `json:"error"`
}

// ObtainTokenFromGitHub gets a token from Squarescale using Github token.
func ObtainTokenFromGitHub(sqscURL, token string) (string, error) {
	var c http.Client
	req, err := http.NewRequest("GET", sqscURL+"/me/token?provider=github&token="+url.QueryEscape(token), nil)
	if err != nil {
		return "", fmt.Errorf("Could not make request: %v", err)
	}
	req.Header.Set("Accept", "application/json")
	res, err := c.Do(req)
	if err != nil {
		return "", fmt.Errorf("Could not send request: %v", err)
	}

	jsondata, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("Could not read response: %v", err)
	}

	var response authTokenResponse
	err = json.Unmarshal(jsondata, &response)
	if err != nil {
		if res.StatusCode != 200 {
			return "", fmt.Errorf("Could not login to Squarescale: %s", res.Status)
		} else {
			return "", fmt.Errorf("Could not unmarshal response: %v", err)
		}
	}

	if res.StatusCode != 200 {
		return "", fmt.Errorf("Could not login to SquareScale: %s", response.Error)
	}

	return response.AuthToken, nil
}
