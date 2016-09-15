package github

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Bowery/prompt"
	"io/ioutil"
	"net/http"
	"strings"
)

type AuthorizationRequestParams struct {
	Scopes []string `json:"scopes"`
	Note   string   `json:"note"`
}

type AuthorizationResponse struct {
	Errors []struct {
		Code string `json:"code"`
	} `json:"errors"`
	HashedToken string `json:"hashed_token"`
	Token       string `json:"token"`
	Url         string `json:"url"`
}

func GeneratePersonalToken(note string) (login, password, otpass, token, token_url string, err error) {
	var c http.Client
	var jsondata []byte

	login, err = prompt.Basic("GitHub Login: ", true)
	if err != nil {
		return
	}
	password, err = prompt.Password("GitHub Password: ")
	if err != nil {
		return
	}

	for i := 1; true; i++ {

		var encodedParams []byte
		var req *http.Request
		var res *http.Response
		params := AuthorizationRequestParams{
			Scopes: []string{"user"},
			Note:   note,
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
				otpass, err = prompt.Basic("GitHub One Time Password ("+strings.Trim(otps[1], " ")+"): ", true)
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

		var response AuthorizationResponse
		err = json.Unmarshal(jsondata, &response)
		if err != nil {
			return
		}

		if res.StatusCode == 422 && len(response.Errors) > 0 && response.Errors[0].Code == "already_exists" {
			fmt.Printf("Token %#v already exists, this is probably a leftover from a previous login attempt. You should remove this token.\n", params.Note)
			continue
		}

		fmt.Printf("Generating token %#s\n", params.Note)

		if res.StatusCode == 201 {
			token = response.Token
			token_url = response.Url
			break
		}

		err = fmt.Errorf("%s error in <%s>: %s", res.Status, req.URL.String(), jsondata)
		return
	}
	return
}

func RevokePersonalToken(token_url, login, pass, otpass string) error {
	var c http.Client
	req, err := http.NewRequest("DELETE", token_url, nil)
	if err != nil {
		return err
	}

	if otpass != "" {
		req.Header.Set("X-GitHub-OTP", otpass)
	}
	req.SetBasicAuth(login, pass)
	res, err := c.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode == 204 {
		return nil
	}

	return fmt.Errorf("%s when removing token <%s>", res.Status, token_url)
}
