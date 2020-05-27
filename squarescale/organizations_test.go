package squarescale_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/squarescale/squarescale-cli/squarescale"
)

func TestGetOrganization(t *testing.T) {
	// Get Organization
	t.Run("Nominal case on GetOrganizationInfo", nominalCaseOnGetOrganizationInfo)
	t.Run("Test organization not found on GetOrganizationInfo", UnknownOrganizationOnGetOrganizationInfo)

	// Add Organization
	t.Run("Nominal case on AddOrganization", nominalCaseOnAddOrganization)
	t.Run("Test duplicate organization name on AddOrganization", DuplicateOrganizationErrorCaseOnAddOrganization)

	// Delete Organization
	t.Run("Nominal case on DeleteOrganization", nominalCaseOnDeleteOrganization)
	t.Run("Test organization not found on DeleteOrganization", UnknownOrganizationOnDeleteOrganization)

	// List Organization
	t.Run("Nominal case on ListOrganizations", nominalCaseOnListOrganizations)

	// Error cases
	t.Run("Test HTTP client error on organization methods (get, add, delete and wait)", ClientHTTPErrorOnOrganizationMethods)
	t.Run("Test internal server error on organization methods (get, add, delete and wait)", InternalServerErrorOnOrganizationMethods)
	t.Run("Test badly JSON on organization methods (get)", CantUnmarshalOnOrganizationMethods)
}

// Get Organization
func nominalCaseOnGetOrganizationInfo(t *testing.T) {
	token := "some-token"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		checkPath(t, "/organizations/Sqsc", r.URL.Path)
		checkAuthorization(t, r.Header.Get("Authorization"), token)

		resBody := `{
			"id":6,
			"name":"Sqsc",
			"collaborators":[
				{
					"id":1,
					"email":"no-reply@squarescale.com",
					"remember_created_at":null,
					"sign_in_count":0,
					"current_sign_in_at":null,
					"last_sign_in_at":null,
					"current_sign_in_ip":null,
					"last_sign_in_ip":null,
					"created_at":"2020-04-27T15:08:38.293Z",
					"updated_at":"2020-05-18T12:43:08.043Z",
					"provider":null,
					"uid":"154c1a74-ddd8-46b6-9263-44a9e6a508d2",
					"name":"User 1",
					"image":null,
					"username":"bdumas",
					"admin":true,
					"max_projects":2,
					"stripe_source_id":null,
					"stripe_customer_id":null,
					"stripe_client_secret":null,
					"legal_entity":true,
					"first_name":"",
					"last_name":"",
					"address":"",
					"zip_code":"",
					"state":"",
					"city":"",
					"phone":"",
					"company_name":"",
					"registration_number":"",
					"vat":"",
					"country":"",
					"voucher_code":null,
					"voucher_remaining_discount":0,
					"country_code":"",
					"voucher_remaining_discount_bill_period":null
				}
			],
			"projects":[
				{
					"id":2,
					"name":"sub-mariner-aerified",
					"created_at":"2020-05-12T13:09:44.625Z",
					"infra_status":"no_infra"
				}
			]
		}`

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(resBody))
	}))

	defer server.Close()

	client := squarescale.NewClient(server.URL, token)

	//when
	organization, err := client.GetOrganizationInfo("Sqsc")

	//then
	expectedName := "Sqsc"
	expectedCollaborators := []squarescale.Collaborator{{Email: "no-reply@squarescale.com", Name: "User 1"}}
	expectedProjects := []squarescale.Project{{Name: "sub-mariner-aerified", InfraStatus: "no_infra"}}

	if err != nil {
		t.Fatalf("Expect no error, got `%s`", err)
	}

	if organization.Name != expectedName {
		t.Errorf("Expect organization.Name `%s`, got `%s`", expectedName, organization.Name)
	}

	if !reflect.DeepEqual(organization.Collaborators, expectedCollaborators) {
		t.Errorf("Expect organization.Collaborators `%v`, got `%v`", expectedCollaborators, organization.Collaborators)
	}

	if !reflect.DeepEqual(organization.Projects, expectedProjects) {
		t.Errorf("Expect organization.Projects `%v`, got `%v`", expectedProjects, organization.Projects)
	}
}

// Add Organization
func nominalCaseOnAddOrganization(t *testing.T) {
	//given
	token := "some-token"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		checkPath(t, "/organizations", r.URL.Path)
		checkAuthorization(t, r.Header.Get("Authorization"), token)

		resBody := `
		{
			"id": 1,
			"name": "Sqsc"
		}
		`

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(201)

		w.Write([]byte(resBody))
	}))

	defer server.Close()
	client := squarescale.NewClient(server.URL, token)

	//when
	err := client.AddOrganization("Sqsc")

	//then
	if err != nil {
		t.Fatalf("Expect no error, got `%s`", err)
	}
}

func DuplicateOrganizationErrorCaseOnAddOrganization(t *testing.T) {
	// given
	token := "some-token"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resBody := `
		{"error":"PG::UniqueViolation: ERROR:  duplicate key value violates unique constraint \"index_organizations_on_name\"\nDETAIL:  Key (name)=(orga1) already exists.\n: INSERT INTO \"organizations\" (\"name\) VALUES ($1) RETURNING \"id\""}
		`

		w.Header().Set("Content-Type", "application/json")

		w.WriteHeader(409)
		w.Write([]byte(resBody))
	}))

	defer server.Close()
	cli := squarescale.NewClient(server.URL, token)

	// when
	err := cli.AddOrganization("Sqsc")

	// then
	expectedError := "Organization already exist: Sqsc"

	if err == nil {
		t.Fatalf("Error is not raised with `%s`", expectedError)
	}

	if fmt.Sprintf("%s", err) != expectedError {
		t.Fatalf("Expected error message:\n`%s`\nGot:\n`%s`", expectedError, err)
	}
}

// Delete Organization
func nominalCaseOnDeleteOrganization(t *testing.T) {
	// give
	token := "some-token"
	organizationName := "organization-test"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		checkPath(t, "/organizations/"+organizationName, r.URL.Path)
		checkAuthorization(t, r.Header.Get("Authorization"), token)

		resBody := `
		null
		`

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)

		w.Write([]byte(resBody))
	}))

	defer server.Close()
	client := squarescale.NewClient(server.URL, token)

	// when
	err := client.DeleteOrganization(organizationName)

	// then
	if err != nil {
		t.Fatalf("Expect no error, got `%s`", err)
	}
}

func UnknownOrganizationOnDeleteOrganization(t *testing.T) {
	// given
	token := "some-token"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resBody := `{"error":"Organization not found"}`

		w.Header().Set("Content-Type", "application/json")

		w.WriteHeader(404)
		w.Write([]byte(resBody))
	}))

	defer server.Close()
	cli := squarescale.NewClient(server.URL, token)

	// when
	err := cli.DeleteOrganization("toto")

	// then
	expectedError := `No organization found for name: toto`
	if err == nil {
		t.Fatalf("Error is not raised with `%s`", expectedError)
	}

	if fmt.Sprintf("%s", err) != expectedError {
		t.Fatalf("Expected error message:\n`%s`\nGot:\n`%s`", expectedError, err)
	}
}

// List Organizations
func nominalCaseOnListOrganizations(t *testing.T) {
	// given
	token := "some-token"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resBody := `
		[
			{
				"id": 6,
				"name": "Sqsc",
				"collaborators": [
					{
						"id": 1,
						"email": "user1@sqsc.fr",
						"name": "User 1"
					},
					{
						"id": 2,
						"email": "user2@sqsc.fr",
						"name": "User 2"
					}
				],
				"projects": [
					{
						"id": 2,
						"name": "sub-mariner-aerified",
						"created_at": "2020-05-12T13:09:44.625Z",
						"infra_status": "no_infra"
					},
					{
						"id": 3,
						"name": "toto",
						"created_at": "2020-05-13T13:09:44.625Z",
						"infra_status": "no_infra"
					}
				]
			}
		]
		`

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(resBody))
	}))

	defer server.Close()

	client := squarescale.NewClient(server.URL, token)

	// when
	organizations, err := client.ListOrganizations()

	fmt.Println(organizations)

	// then
	expectedCollaborators := []squarescale.Collaborator{
		{Email: "user1@sqsc.fr", Name: "User 1"},
		{Email: "user2@sqsc.fr", Name: "User 2"},
	}
	expectedProjects := []squarescale.Project{
		{Name: "sub-mariner-aerified", InfraStatus: "no_infra"},
		{Name: "toto", InfraStatus: "no_infra"},
	}

	if err != nil {
		t.Fatalf("Expect no error, got `%s`", err)
	}

	if !reflect.DeepEqual(organizations[0].Collaborators, expectedCollaborators) {
		t.Errorf("Expect organization.Collaborators `%v`, got `%v`", expectedCollaborators, organizations[0].Collaborators)
	}

	if !reflect.DeepEqual(organizations[0].Projects, expectedProjects) {
		t.Errorf("Expect organization.Projects `%v`, got `%v`", expectedProjects, organizations[0].Projects)
	}
}

// Error cases
func UnknownOrganizationOnGetOrganizationInfo(t *testing.T) {
	// given
	token := "some-token"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resBody := `
		{
			"error": "No organization found"
		}
		`

		w.Header().Set("Content-Type", "application/json")

		w.WriteHeader(404)
		w.Write([]byte(resBody))
	}))

	defer server.Close()
	client := squarescale.NewClient(server.URL, token)

	// when
	_, err := client.GetOrganizationInfo("OrgaNotFound")

	// then
	expectedError := "Organization 'OrgaNotFound' not found"
	if err == nil {
		t.Fatalf("Error is not raised with `%s`", expectedError)
	}

	if fmt.Sprintf("%s", err) != expectedError {
		t.Fatalf("Expected error message:\n`%s`\nGot:\n`%s`", expectedError, err)
	}
}

func CantUnmarshalOnOrganizationMethods(t *testing.T) {
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
	_, errOnGetOrganizationInfo := client.GetOrganizationInfo("Sqsc")
	_, errOnListOrganizations := client.ListOrganizations()

	// then
	expectedError := "invalid character ']' looking for beginning of object key string"

	if errOnGetOrganizationInfo == nil {
		t.Fatalf("Error is not raised with `%s`", expectedError)
	}

	if fmt.Sprintf("%s", errOnGetOrganizationInfo) != expectedError {
		t.Fatalf("Expected error message:\n`%s`\nGot:\n`%s`", expectedError, errOnGetOrganizationInfo)
	}

	if errOnListOrganizations == nil {
		t.Fatalf("Error is not raised with `%s`", expectedError)
	}

	if fmt.Sprintf("%s", errOnListOrganizations) != expectedError {
		t.Fatalf("Expected error message:\n`%s`\nGot:\n`%s`", expectedError, errOnListOrganizations)
	}
}

func ClientHTTPErrorOnOrganizationMethods(t *testing.T) {
	// given
	token := "some-token"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	}))

	defer server.Close()
	cli := squarescale.NewClient(server.URL, token)

	// when
	errOnAdd := cli.AddOrganization("Sqsc")
	errOnDelete := cli.DeleteOrganization("Sqsc")
	_, errOnGet := cli.GetOrganizationInfo("Sqsc")
	_, errOnList := cli.ListOrganizations()

	// then
	if errOnAdd == nil {
		t.Errorf("Error is not raised on AddOrganization")
	}

	if errOnDelete == nil {
		t.Errorf("Error is not raised on DeleteOrganization")
	}

	if errOnGet == nil {
		t.Errorf("Error is not raised on GetOrganizationInfo")
	}

	if errOnList == nil {
		t.Errorf("Error is not raised on ListOrganizations")
	}
}

func InternalServerErrorOnOrganizationMethods(t *testing.T) {
	// given
	token := "some-token"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(500)
	}))

	defer server.Close()
	client := squarescale.NewClient(server.URL, token)

	// when
	errOnAddOrganization := client.AddOrganization("Sqsc")
	errOnDeleteOrganization := client.DeleteOrganization("Sqsc")
	_, errOnGetOrganization := client.GetOrganizationInfo("Sqsc")
	_, errOnListOrganizations := client.ListOrganizations()

	// then
	expectedError := "An unexpected error occurred (code: 500)"

	if errOnAddOrganization == nil {
		t.Errorf("Error is not raised with `%s`", expectedError)
	}

	if errOnDeleteOrganization == nil {
		t.Errorf("Error is not raised with `%s`", expectedError)
	}

	if errOnGetOrganization == nil {
		t.Errorf("Error is not raised with `%s`", expectedError)
	}

	if errOnListOrganizations == nil {
		t.Errorf("Error is not raised with `%s`", expectedError)
	}

	if fmt.Sprintf("%s", errOnAddOrganization) != expectedError {
		t.Errorf("Expected error message:\n`%s`\nGot:\n`%s`", expectedError, errOnAddOrganization)
	}

	if fmt.Sprintf("%s", errOnDeleteOrganization) != expectedError {
		t.Errorf("Expected error message:\n`%s`\nGot:\n`%s`", expectedError, errOnDeleteOrganization)
	}

	if fmt.Sprintf("%s", errOnGetOrganization) != expectedError {
		t.Errorf("Expected error message:\n`%s`\nGot:\n`%s`", expectedError, errOnGetOrganization)
	}

	if fmt.Sprintf("%s", errOnListOrganizations) != expectedError {
		t.Errorf("Expected error message:\n`%s`\nGot:\n`%s`", expectedError, errOnListOrganizations)
	}
}
