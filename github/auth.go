package github

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type ghAuthorizationRequestParams struct {
	Scopes []string `json:"scopes"`
	Note   string   `json:"note"`
}

type ghAuthorizationResponse struct {
	Errors []struct {
		Code string `json:"code"`
	} `json:"errors"`
	HashedToken string `json:"hashed_token"`
	Token       string `json:"token"`
	URL         string `json:"url"`
}

// OneTimePasswordPrompt is a simple interface to ask the user about his one time password.
type OneTimePasswordPrompt interface {
	Ask(string) (string, error)
	Info(string)
}

// GeneratePersonalToken contacts GitHub APIs to generate a token according to provided credentials.
func GeneratePersonalToken(login, password string, prompt OneTimePasswordPrompt) (otpass, token, tokenURL string, err error) {
	var c http.Client
	var jsondata []byte

	for i := 1; ; i++ {

		var encodedParams []byte
		var req *http.Request
		var res *http.Response
		params := ghAuthorizationRequestParams{
			Scopes: []string{"user"},
			Note:   "Squarescale CLI",
		}

		if i > 1 {
			params.Note = fmt.Sprintf("%s %d", params.Note, i)
		}

		encodedParams, err = json.Marshal(params)
		if err != nil {
			return
		}

		req, err = http.NewRequest("POST", "https://api.github.com/authorizations", bytes.NewReader(encodedParams))
		if err != nil {
			return
		}

		req.SetBasicAuth(login, password)
		res, err = c.Do(req)
		if err != nil {
			return
		}

		otp := res.Header.Get("X-Github-Otp")
		otps := strings.SplitN(otp, ";", 2)

		if otps[0] == "required" {
			if otpass == "" {
				otpass, err = prompt.Ask("GitHub One Time Password (" + strings.Trim(otps[1], " ") + "):")
				if err != nil {
					return
				}
			}

			req.Header.Set("X-GitHub-OTP", otpass)
			req.Body = ioutil.NopCloser(bytes.NewReader(encodedParams))
			res, err = c.Do(req)
			if err != nil {
				return
			}
		}

		jsondata, err = ioutil.ReadAll(res.Body)
		if err != nil {
			return
		}

		var response ghAuthorizationResponse
		err = json.Unmarshal(jsondata, &response)
		if err != nil {
			return
		}

		if res.StatusCode == http.StatusUnprocessableEntity && len(response.Errors) > 0 && response.Errors[0].Code == "already_exists" {
			prompt.Info(fmt.Sprintf("Token %v already exists, this is probably a leftover from a previous login attempt. You should remove this token.", params.Note))
			continue
		}

		prompt.Info(fmt.Sprintf("Generating token '%s'", params.Note))

		if res.StatusCode == http.StatusCreated {
			token = response.Token
			tokenURL = response.URL
			return
		}

		err = fmt.Errorf("%s error in '%s'", res.Status, req.URL.String())
		return
	}
}

// RevokePersonalToken revokes the token using the provided token url.
func RevokePersonalToken(tokenURL, login, pass, otpass string) error {
	req, err := http.NewRequest("DELETE", tokenURL, nil)
	if err != nil {
		return err
	}

	if otpass != "" {
		req.Header.Set("X-GitHub-OTP", otpass)
	}

	req.SetBasicAuth(login, pass)
	var c http.Client
	res, err := c.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode == http.StatusNoContent {
		return nil
	}

	return fmt.Errorf("%s when removing token <%s>", res.Status, tokenURL)
}
