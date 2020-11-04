package squarescale_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/squarescale/squarescale-cli/squarescale"
)

func TestLoadBalancerGet(t *testing.T) {
	// Get LoadBalancers
	t.Run("Nominal case on LoadBalancerGet", nominalCaseOnLoadBalancerGet)
	t.Run("Nominal case on LoadBalancerEnable", nominalCaseOnLoadBalancerEnable)
	t.Run("Nominal case on LoadBalancerDisable", nominalCaseOnLoadBalancerDisable)

	t.Run("Test project not found on LoadBalancerGet", UnknownProjectOnLoadBalancerGet)
	t.Run("Test project not found on LoadBalancerEnable", UnknownProjectOnLoadBalancerEnable)
	t.Run("Test project not found on LoadBalancerDisable", UnknownProjectOnLoadBalancerDisable)

	// Error cases
	t.Run("Test HTTP client error on load_balancers methods (get, configure)", ClientHTTPErrorOnLoadBalancersMethods)
	t.Run("Test internal server error on load_balancers methods (get, configure)", InternalServerErrorOnLoadBalancersMethods)
	t.Run("Test badly JSON on load_balancers methods (get)", CantUnmarshalOnLoadBalancersMethods)
}

// Get LoadBalancers
func nominalCaseOnLoadBalancerGet(t *testing.T) {
	token := "some-token"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		checkPath(t, "/projects/61d3bd06-0b9b-4e3a-9916-69ce6257cc3f/load_balancers", r.URL.Path)
		checkAuthorization(t, r.Header.Get("Authorization"), token)

		resBody := `[{
			"id":3,
			"active":true,
			"certificate_body":null,
			"certificate_chain":"",
			"https":false,
			"public_url":"https://www.mu9nius7aiBiengo.projects.squarescale.io"
			}]`

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(resBody))
	}))

	defer server.Close()

	client := squarescale.NewClient(server.URL, token)

	//when
	loadBalancers, err := client.LoadBalancerGet("61d3bd06-0b9b-4e3a-9916-69ce6257cc3f")

	//then
	expectedPublicURL := "https://www.mu9nius7aiBiengo.projects.squarescale.io"

	if err != nil {
		t.Fatalf("Expect no error, got `%s`", err)
	}

	if loadBalancers[0].PublicURL != expectedPublicURL {
		t.Errorf("Expect loadBalancers.PublicURL `%s`, got `%s`", expectedPublicURL, loadBalancers[0].PublicURL)
	}
}

// Enable LoadBalancers
func nominalCaseOnLoadBalancerEnable(t *testing.T) {
	token := "some-token"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		checkPath(t, "/projects/61d3bd06-0b9b-4e3a-9916-69ce6257cc3f/load_balancers/42", r.URL.Path)
		checkAuthorization(t, r.Header.Get("Authorization"), token)

		resBody := `null`

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(resBody))
	}))

	defer server.Close()

	client := squarescale.NewClient(server.URL, token)
	cert := "-----BEGIN CERTIFICATE-----\nMIIGbDCC\n-----END CERTIFICATE-----\n"
	certChain := "-----BEGIN CERTIFICATE-----\nMIIGbDCCBVSg\n-----END CERTIFICATE-----\n"
	secretKey := "-----BEGIN RSA PRIVATE KEY-----\nMIIJK\n-----END RSA PRIVATE KEY-----\n"

	//when
	err := client.LoadBalancerEnable("61d3bd06-0b9b-4e3a-9916-69ce6257cc3f", 42, cert, certChain, secretKey)

	if err != nil {
		t.Fatalf("Expect no error, got `%s`", err)
	}
}

// Disable LoadBalancers
func nominalCaseOnLoadBalancerDisable(t *testing.T) {
	token := "some-token"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		checkPath(t, "/projects/61d3bd06-0b9b-4e3a-9916-69ce6257cc3f/load_balancers/42", r.URL.Path)
		checkAuthorization(t, r.Header.Get("Authorization"), token)

		resBody := `null`

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(resBody))
	}))

	defer server.Close()

	client := squarescale.NewClient(server.URL, token)

	//when
	err := client.LoadBalancerDisable("61d3bd06-0b9b-4e3a-9916-69ce6257cc3f", 42)

	if err != nil {
		t.Fatalf("Expect no error, got `%s`", err)
	}
}

func UnknownProjectOnLoadBalancerGet(t *testing.T) {
	// given
	token := "some-token"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resBody := `{"error":"Couldn't find Project with [WHERE \"projects\".\"uuid\" = $1]"}`

		w.Header().Set("Content-Type", "application/json")

		w.WriteHeader(404)
		w.Write([]byte(resBody))
	}))

	defer server.Close()
	client := squarescale.NewClient(server.URL, token)

	// when
	_, err := client.LoadBalancerGet("missing-project-uuid")

	// then
	expectedError := "Project 'missing-project-uuid' not found"
	if err == nil {
		t.Fatalf("Error is not raised with `%s`", expectedError)
	}

	if fmt.Sprintf("%s", err) != expectedError {
		t.Fatalf("Expected error message:\n`%s`\nGot:\n`%s`", expectedError, err)
	}
}

func UnknownProjectOnLoadBalancerEnable(t *testing.T) {
	// given
	token := "some-token"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resBody := `{"error":"Couldn't find Project with [WHERE \"projects\".\"uuid\" = $1]"}`

		w.Header().Set("Content-Type", "application/json")

		w.WriteHeader(404)
		w.Write([]byte(resBody))
	}))

	defer server.Close()
	client := squarescale.NewClient(server.URL, token)
	cert := "-----BEGIN CERTIFICATE-----\nMIIGbDCC\n-----END CERTIFICATE-----\n"
	certChain := "-----BEGIN CERTIFICATE-----\nMIIGbDCCBVSg\n-----END CERTIFICATE-----\n"
	secretKey := "-----BEGIN RSA PRIVATE KEY-----\nMIIJK\n-----END RSA PRIVATE KEY-----\n"

	//when
	err := client.LoadBalancerEnable("missing-project-uuid", 42, cert, certChain, secretKey)

	// then
	expectedError := "Project 'missing-project-uuid' not found"
	if err == nil {
		t.Fatalf("Error is not raised with `%s`", expectedError)
	}

	if fmt.Sprintf("%s", err) != expectedError {
		t.Fatalf("Expected error message:\n`%s`\nGot:\n`%s`", expectedError, err)
	}
}

func UnknownProjectOnLoadBalancerDisable(t *testing.T) {
	// given
	token := "some-token"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resBody := `{"error":"Couldn't find Project with [WHERE \"projects\".\"uuid\" = $1]"}`

		w.Header().Set("Content-Type", "application/json")

		w.WriteHeader(404)
		w.Write([]byte(resBody))
	}))

	defer server.Close()
	client := squarescale.NewClient(server.URL, token)

	//when
	err := client.LoadBalancerDisable("missing-project-uuid", 42)

	// then
	expectedError := "Project 'missing-project-uuid' not found"
	if err == nil {
		t.Fatalf("Error is not raised with `%s`", expectedError)
	}

	if fmt.Sprintf("%s", err) != expectedError {
		t.Fatalf("Expected error message:\n`%s`\nGot:\n`%s`", expectedError, err)
	}
}

// Error cases
func CantUnmarshalOnLoadBalancersMethods(t *testing.T) {
	// given
	token := "some-token"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resBody := `{]`

		w.Header().Set("Content-Type", "application/json")

		w.Write([]byte(resBody))
	}))

	defer server.Close()
	client := squarescale.NewClient(server.URL, token)

	// when
	_, errOnLoadBalancerGet := client.LoadBalancerGet("61d3bd06-0b9b-4e3a-9916-69ce6257cc3f")

	// then
	expectedError := "invalid character ']' looking for beginning of object key string"

	if errOnLoadBalancerGet == nil {
		t.Fatalf("Error is not raised with `%s`", expectedError)
	}

	if fmt.Sprintf("%s", errOnLoadBalancerGet) != expectedError {
		t.Fatalf("Expected error message:\n`%s`\nGot:\n`%s`", expectedError, errOnLoadBalancerGet)
	}
}

func ClientHTTPErrorOnLoadBalancersMethods(t *testing.T) {
	// given
	token := "some-token"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	}))

	defer server.Close()
	client := squarescale.NewClient(server.URL, token)

	// when
	_, errOnGet := client.LoadBalancerGet("61d3bd06-0b9b-4e3a-9916-69ce6257cc3f")
	errOnEnable := client.LoadBalancerEnable("61d3bd06-0b9b-4e3a-9916-69ce6257cc3f", 42, "", "", "")
	errOnDisable := client.LoadBalancerDisable("61d3bd06-0b9b-4e3a-9916-69ce6257cc3f", 42)

	// then
	if errOnGet == nil {
		t.Errorf("Error is not raised on LoadBalancerGet")
	}

	if errOnEnable == nil {
		t.Errorf("Error is not raised on LoadBalancerEnable")
	}

	if errOnDisable == nil {
		t.Errorf("Error is not raised on LoadBalancerDisable")
	}
}

func InternalServerErrorOnLoadBalancersMethods(t *testing.T) {
	// given
	token := "some-token"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(500)
	}))

	defer server.Close()
	client := squarescale.NewClient(server.URL, token)

	// when
	_, errOnLoadBalancerGet := client.LoadBalancerGet("61d3bd06-0b9b-4e3a-9916-69ce6257cc3f")
	errOnLoadBalancerEnable := client.LoadBalancerEnable("61d3bd06-0b9b-4e3a-9916-69ce6257cc3f", 42, "", "", "")
	errOnLoadBalancerDisable := client.LoadBalancerDisable("61d3bd06-0b9b-4e3a-9916-69ce6257cc3f", 42)

	// then
	expectedError := "An unexpected error occurred (code: 500)"

	if errOnLoadBalancerGet == nil {
		t.Errorf("Error is not raised with `%s`", expectedError)
	}

	if fmt.Sprintf("%s", errOnLoadBalancerGet) != expectedError {
		t.Errorf("Expected error message:\n`%s`\nGot:\n`%s`", expectedError, errOnLoadBalancerGet)
	}

	if errOnLoadBalancerEnable == nil {
		t.Errorf("Error is not raised with `%s`", expectedError)
	}

	if fmt.Sprintf("%s", errOnLoadBalancerEnable) != expectedError {
		t.Errorf("Expected error message:\n`%s`\nGot:\n`%s`", expectedError, errOnLoadBalancerEnable)
	}

	if errOnLoadBalancerDisable == nil {
		t.Errorf("Error is not raised with `%s`", expectedError)
	}

	if fmt.Sprintf("%s", errOnLoadBalancerDisable) != expectedError {
		t.Errorf("Expected error message:\n`%s`\nGot:\n`%s`", expectedError, errOnLoadBalancerDisable)
	}
}
