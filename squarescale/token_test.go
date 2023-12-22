package squarescale_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/squarescale/squarescale-cli/squarescale"
	"github.com/stretchr/testify/assert"
)

func TestToken(t *testing.T) {
	// ValidateToken HTTP error
	t.Run("Nominal case on ValidateToken", serverUnauthorizedCaseOnToken)
	// ValidateToken not modified
	t.Run("Nominal case on ValidateToken", notModifiedCaseOnToken)
	// ValidateToken bad JSON
	t.Run("Nominal case on ValidateToken", badJSONCaseOnToken)
	// ValidateToken OK
	t.Run("Nominal case on ValidateToken", everythingOKCaseOnToken)
}

func serverUnauthorizedCaseOnToken(t *testing.T) {
	// given
	token := "some-token"
	expectedErrorMsg := "401 Unauthorized: map[message:Unauthorized], please make sure your SQSC_TOKEN environment variable is properly set"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		checkPath(t, "/me", r.URL.Path)
		checkAuthorization(t, r.Header.Get("Authorization"), token)

		resBody := `
		{
			"message": "Unauthorized"
		}
		`

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(401)
		w.Write([]byte(resBody))
	}))

	defer server.Close()
	cli := squarescale.NewClient(server.URL, token)

	// when
	_, err := cli.ValidateToken()
	if err == nil {
		t.Errorf("Expecting error, got `%+v`", err)
	}
	assert.EqualErrorf(t, err, expectedErrorMsg, "Error should be: %v, got: %v", expectedErrorMsg, err)
}

func notModifiedCaseOnToken(t *testing.T) {
	// given
	token := "some-token"
	expectedErrorMsg := "You're not logged in, please run login command"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		checkPath(t, "/me", r.URL.Path)
		checkAuthorization(t, r.Header.Get("Authorization"), token)

		resBody := `{}`

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(204)
		w.Write([]byte(resBody))
	}))

	defer server.Close()
	cli := squarescale.NewClient(server.URL, token)

	// when
	_, err := cli.ValidateToken()
	if err == nil {
		t.Errorf("Expecting error, got `%+v`", err)
	}
	assert.EqualErrorf(t, err, expectedErrorMsg, "Error should be: %v, got: %v", expectedErrorMsg, err)
}

func badJSONCaseOnToken(t *testing.T) {
	// given
	token := "some-token"
	expectedErrorMsg := "invalid character 'd' looking for beginning of value"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		checkPath(t, "/me", r.URL.Path)
		checkAuthorization(t, r.Header.Get("Authorization"), token)

		resBody := `dummy`

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write([]byte(resBody))
	}))

	defer server.Close()
	cli := squarescale.NewClient(server.URL, token)

	// when
	user, err := cli.ValidateToken()
	fmt.Printf("GOT USER %+v === Err %+v ===\n", user, err)
	if err == nil {
		t.Errorf("Expecting error, got `%+v`", err)
	}
	assert.EqualErrorf(t, err, expectedErrorMsg, "Error should be: %v, got: %v", expectedErrorMsg, err)
}

func everythingOKCaseOnToken(t *testing.T) {
	// given
	token := "some-token"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		checkPath(t, "/me", r.URL.Path)
		checkAuthorization(t, r.Header.Get("Authorization"), token)

		resBody := `{}`

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write([]byte(resBody))
	}))

	defer server.Close()
	cli := squarescale.NewClient(server.URL, token)

	// when
	user, err := cli.ValidateToken()
	fmt.Printf("GOT USER %+v === Err %+v ===\n", user, err)
	if err != nil {
		t.Errorf("Expecting no error, got `%+v`", err)
	}
}
