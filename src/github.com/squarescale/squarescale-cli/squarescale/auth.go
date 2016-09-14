package squarescale

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type authTokenResponse struct {
	AuthToken string `json:"auth_token"`
	Error     string `json:"error"`
}

func ObtainTokenFromGitHub(sqsc_url, token string) (string, error) {
	var c http.Client
	req, err := http.NewRequest("GET", sqsc_url+"/me/token?provider=github&token="+token, nil)
	if err != nil {
		return "", err
	}
	res, err := c.Do(req)
	if err != nil {
		return "", err
	}

	jsondata, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	var response authTokenResponse
	err = json.Unmarshal(jsondata, &response)
	if err != nil {
		return "", err
	}

	if res.StatusCode != 200 {
		return "", fmt.Errorf("Could not login to squarescale: %s", response.Error)
	}
	return response.AuthToken, nil
}
