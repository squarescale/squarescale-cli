package squarescale_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/squarescale/squarescale-cli/squarescale"
)

func TestProject(t *testing.T) {
	// CreateProject
	t.Run("Nominal case on CreateProject", nominalCaseCreateProject)

	t.Run("Credential is mandatory case on CreateProject", credentialMandatoryCaseCreateProject)
	t.Run("Badly http error code case on CreateProject", badHttpErrrorCoreCaseCreateProject)
	t.Run("Error in task case on CreateProject", errorInTaskCaseCreateProject)

	// ProvisionProject
	t.Run("Nominal case on ProvisionProject", nominalCaseProvisionProject)
	t.Run("Badly http error code case on ProvisionProject", badHttpErrrorCoreCaseProvisionProject)

	// UnprovisionProject
	t.Run("Nominal case on UnprovisionProject", nominalCaseUnprovisionProject)
	t.Run("Badly http error code case on UNProvisionProject", badHttpErrrorCoreCaseUNProvisionProject)

	// ListProjects
	t.Run("Nominal case on ListProjects", nominalCaseListProjects)
	t.Run("Badly http error code case on ListProjects", badHttpErrrorCoreCaseListProjects)

	// GetProject
	t.Run("Nominal case on GetProject", nominalCaseGetProject)

	t.Run("Project not found on GetProject", UnknownProjectOnGetProject)
	t.Run("Badly http error code case on GetProject", badHttpErrrorCoreCaseGetProject)

	// WaitProject
	t.Run("Nominal case on WaitProject", nominalCaseWaitProject)

	// ProjectUnprovision
	t.Run("Nominal case on ProjectUnprovision", nominalCaseProjectUnprovision)

	t.Run("Project not found on ProjectUnprovision", UnknownProjectOnProjectUnprovision)
	t.Run("Badly http error code case on ProjectUnprovision", badHttpErrrorCoreCaseOnProjectUnprovision)

	// ProjectDelete
	t.Run("Nominal case on ProjectDelete", nominalCaseProjectDelete)

	t.Run("Project not found on ProjectDelete", UnknownProjectOnProjectDelete)
	t.Run("Badly http error code case on ProjectDelete", badHttpErrrorCoreCaseOnProjectDelete)

	// Error cases
	t.Run("Test HTTP client error on project methods (create, provision, unprovision, get)", ClientHTTPErrorOnProjectMethods)
	// t.Run("Test internal server error on project methods ()", InternalServerErrorOnProjectMethods)
	t.Run("Test badly JSON on project methods (create)", CantUnmarshalOnProjectMethods)
}

func nominalCaseCreateProject(t *testing.T) {
	// given
	token := "some-token"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		checkPath(t, "/projects", r.URL.Path)
		checkAuthorization(t, r.Header.Get("Authorization"), token)

		resBody := `
			{
				"id": 10,
				"name": "project-test",
				"uuid": "ba90e5fe-f520-4275-897b-49a95c1157a3",
				"next": "https://cloud.squarescale.io/projects/ba90e5fe-f520-4275-897b-49a95c1157a3",
				"load_balancer": true,
				"db_enabled": false,
				"db_engine": "postgres",
				"db_size": "dev",
				"db_version": null,
				"database": {
					"enabled": false,
					"engine": "postgres",
					"size": "dev"
				},
				"infra_status": "no_infra",
				"cluster_size": 1,
				"task": 12
			}
		`

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(201)

		w.Write([]byte(resBody))
	}))

	defer server.Close()
	cli := squarescale.NewClient(server.URL, token)

	// when
	payload := squarescale.JSONObject{}
	payload["name"] = "project-test"
	payload["credential_name"] = "aws-production"
	payload["provider"] = "aws"
	payload["region"] = "eu-west-1"
	payload["infra_type"] = "single_node"
	payload["node_size"] = "dev"
	payload["slack_webhook"] = ""
	payload["databases"] = "[]"
	taskID, err := cli.CreateProject(&payload)

	// then
	if err != nil {
		t.Fatalf("Expect no error, got `%s`", err)
	}

	expectedTaskID := 12
	if taskID != expectedTaskID {
		t.Errorf("Expect volume.ID `%d`, got `%d`", expectedTaskID, taskID)
	}
}

func credentialMandatoryCaseCreateProject(t *testing.T) {
	// given
	token := "some-token"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		checkPath(t, "/projects", r.URL.Path)
		checkAuthorization(t, r.Header.Get("Authorization"), token)

		resBody := ``

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(500)

		w.Write([]byte(resBody))
	}))

	defer server.Close()
	cli := squarescale.NewClient(server.URL, token)

	// when
	payload := squarescale.JSONObject{}
	payload["name"] = "project-test"
	payload["provider"] = "aws"
	payload["region"] = "eu-west-1"
	payload["infra_type"] = "single_node"
	payload["node_size"] = "dev"
	payload["slack_webhook"] = ""
	payload["databases"] = "[]"
	taskID, err := cli.CreateProject(&payload)

	// then
	expectedError := "Credential is mandatory"
	if fmt.Sprintf("%s", err) != expectedError {
		t.Fatalf("Error is not raised with `%s`", expectedError)
	}

	expectedTaskID := 0
	if taskID != expectedTaskID {
		t.Errorf("Expect volume.ID `%d`, got `%d`", expectedTaskID, taskID)
	}
}

func badHttpErrrorCoreCaseCreateProject(t *testing.T) {
	// given
	token := "some-token"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		checkPath(t, "/projects", r.URL.Path)
		checkAuthorization(t, r.Header.Get("Authorization"), token)

		resBody := ``

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(500)

		w.Write([]byte(resBody))
	}))

	defer server.Close()
	cli := squarescale.NewClient(server.URL, token)

	// when
	payload := squarescale.JSONObject{}
	payload["name"] = "project-test"
	payload["credential_name"] = "aws-production"
	payload["provider"] = "aws"
	payload["region"] = "eu-west-1"
	payload["infra_type"] = "single_node"
	payload["node_size"] = "dev"
	payload["slack_webhook"] = ""
	payload["databases"] = "[]"
	taskID, err := cli.CreateProject(&payload)

	// then
	expectedError := "An unexpected error occurred (code: 500)"
	if fmt.Sprintf("%s", err) != expectedError {
		t.Fatalf("Error is not raised with `%s`", expectedError)
	}

	expectedTaskID := 0
	if taskID != expectedTaskID {
		t.Errorf("Expect volume.ID `%d`, got `%d`", expectedTaskID, taskID)
	}
}

func errorInTaskCaseCreateProject(t *testing.T) {
	// given
	token := "some-token"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		checkPath(t, "/projects", r.URL.Path)
		checkAuthorization(t, r.Header.Get("Authorization"), token)

		resBody := `
			{
				"id": 10,
				"name": "project-test",
				"uuid": "ba90e5fe-f520-4275-897b-49a95c1157a3",
				"next": "https://cloud.squarescale.io/projects/ba90e5fe-f520-4275-897b-49a95c1157a3",
				"load_balancer": true,
				"db_enabled": false,
				"db_engine": "postgres",
				"db_size": "dev",
				"db_version": null,
				"database": {
					"enabled": false,
					"engine": "postgres",
					"size": "dev"
				},
				"infra_status": "no_infra",
				"cluster_size": 1,
				"error": "Error in creation",
				"task": 12
			}
		`

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(201)

		w.Write([]byte(resBody))
	}))

	defer server.Close()
	cli := squarescale.NewClient(server.URL, token)

	// when
	payload := squarescale.JSONObject{}
	payload["name"] = "project-test"
	payload["credential_name"] = "aws-production"
	payload["provider"] = "aws"
	payload["region"] = "eu-west-1"
	payload["infra_type"] = "single_node"
	payload["node_size"] = "dev"
	payload["slack_webhook"] = ""
	payload["databases"] = "[]"
	taskID, err := cli.CreateProject(&payload)

	// then
	expectedError := "Error in creation"
	if err == nil {
		t.Fatalf("Error is not raised with `%s`", expectedError)
	}

	if fmt.Sprintf("%s", err) != expectedError {
		t.Fatalf("Expected error message:\n`%s`\nGot:\n`%s`", expectedError, err)
	}

	expectedTaskID := 12
	if taskID != expectedTaskID {
		t.Errorf("Expect volume.ID `%d`, got `%d`", expectedTaskID, taskID)
	}
}

func nominalCaseProvisionProject(t *testing.T) {
	// given
	token := "some-token"
	projectUUID := "ba90e5fe-f520-4275-897b-49a95c1157a3"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		checkPath(t, "/projects/"+projectUUID+"/provision", r.URL.Path)
		checkAuthorization(t, r.Header.Get("Authorization"), token)

		resBody := ``

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(204)

		w.Write([]byte(resBody))
	}))

	defer server.Close()
	cli := squarescale.NewClient(server.URL, token)

	// when
	err := cli.ProvisionProject(projectUUID)

	// then
	if err != nil {
		t.Fatalf("Expect no error, got `%s`", err)
	}
}

func badHttpErrrorCoreCaseProvisionProject(t *testing.T) {
	// given
	token := "some-token"
	projectUUID := "ba90e5fe-f520-4275-897b-49a95c1157a3"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		checkPath(t, "/projects/"+projectUUID+"/provision", r.URL.Path)
		checkAuthorization(t, r.Header.Get("Authorization"), token)

		resBody := ``

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(500)

		w.Write([]byte(resBody))
	}))

	defer server.Close()
	cli := squarescale.NewClient(server.URL, token)

	// when
	err := cli.ProvisionProject(projectUUID)

	// then
	expectedError := "An unexpected error occurred (code: 500)"
	if fmt.Sprintf("%s", err) != expectedError {
		t.Fatalf("Error is not raised with `%s`", expectedError)
	}
}

func nominalCaseUnprovisionProject(t *testing.T) {
	// given
	token := "some-token"
	projectUUID := "ba90e5fe-f520-4275-897b-49a95c1157a3"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		checkPath(t, "/projects/"+projectUUID+"/unprovision", r.URL.Path)
		checkAuthorization(t, r.Header.Get("Authorization"), token)

		resBody := ``

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(204)

		w.Write([]byte(resBody))
	}))

	defer server.Close()
	cli := squarescale.NewClient(server.URL, token)

	// when
	err := cli.UNProvisionProject(projectUUID)

	// then
	if err != nil {
		t.Fatalf("Expect no error, got `%s`", err)
	}
}

func badHttpErrrorCoreCaseUNProvisionProject(t *testing.T) {
	// given
	token := "some-token"
	projectUUID := "ba90e5fe-f520-4275-897b-49a95c1157a3"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		checkPath(t, "/projects/"+projectUUID+"/unprovision", r.URL.Path)
		checkAuthorization(t, r.Header.Get("Authorization"), token)

		resBody := ``

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(500)

		w.Write([]byte(resBody))
	}))

	defer server.Close()
	cli := squarescale.NewClient(server.URL, token)

	// when
	err := cli.UNProvisionProject(projectUUID)

	// then
	expectedError := "An unexpected error occurred (code: 500)"
	if fmt.Sprintf("%s", err) != expectedError {
		t.Fatalf("Error is not raised with `%s`", expectedError)
	}
}

func nominalCaseListProjects(t *testing.T) {
	// given
	token := "some-token"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		checkPath(t, "/projects", r.URL.Path)
		checkAuthorization(t, r.Header.Get("Authorization"), token)

		resBody := `
		[
			{
				"id": 11,
				"uuid": "8cfe8f68-cad5-4157-b8a6-d9efa12caf0e",
				"name": "tera-project",
				"user": {
					"id": 1,
					"email": "john.doe@example.com",
					"uid": "eb3e2663-6438-4e7f-9e5d-b486d6c37083",
					"name": "John Doe"
				},
				"provider_label": "Amazon Web Service",
				"region_label": "Europe (Paris)",
				"credentials_name": "aws-dev",
				"created_at": "2020-11-13T14:55:07.880Z",
				"infra_status": "ok"
			},
			{
				"id": 10,
				"uuid": "ba90e5fe-f520-4275-897b-49a95c1157a3",
				"name": "nova-project",
				"user": {
					"id": 1,
					"email": "john.doe@example.com",
					"uid": "eb3e2663-6438-4e7f-9e5d-b486d6c37083",
					"name": "John Doe"
				},
				"provider_label": "Amazon Web Service",
				"region_label": "Europe (Irlande)",
				"credentials_name": "aws-prod",
				"created_at": "2020-11-12T15:01:41.310Z",
				"infra_status": "provisionning"
			}
		]`

		w.Header().Set("Content-Type", "application/json")

		w.Write([]byte(resBody))
	}))

	defer server.Close()
	cli := squarescale.NewClient(server.URL, token)

	// when
	projectsList, err := cli.ListProjects()

	// then
	var expectedInt int
	var expectedString string

	if err != nil {
		t.Fatalf("Expect no error, got `%s`", err)
	}

	expectedInt = 2
	if len(projectsList) != expectedInt {
		t.Fatalf("Expect projectsList to contain one element `%d`, but got actually `%d`", expectedInt, len(projectsList))
	}

	expectedString = "8cfe8f68-cad5-4157-b8a6-d9efa12caf0e"
	if projectsList[0].UUID != expectedString {
		t.Fatalf("Expect projectsList[0].UUID `%s`, got `%s`", expectedString, projectsList[0].UUID)
	}

	expectedString = "ba90e5fe-f520-4275-897b-49a95c1157a3"
	if projectsList[1].UUID != expectedString {
		t.Fatalf("Expect projectsList[1].UUID `%s`, got `%s`", expectedString, projectsList[1].UUID)
	}

	expectedString = "ok"
	if projectsList[0].InfraStatus != expectedString {
		t.Fatalf("Expect projectsList[0].InfraStatus`%s`, got `%s`", expectedString, projectsList[0].InfraStatus)
	}

	expectedString = "provisionning"
	if projectsList[1].InfraStatus != expectedString {
		t.Fatalf("Expect projectsList[1].InfraStatus`%s`, got `%s`", expectedString, projectsList[1].InfraStatus)
	}
}

func badHttpErrrorCoreCaseListProjects(t *testing.T) {
	// given
	token := "some-token"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		checkPath(t, "/projects", r.URL.Path)
		checkAuthorization(t, r.Header.Get("Authorization"), token)

		resBody := ``

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(500)

		w.Write([]byte(resBody))
	}))

	defer server.Close()
	cli := squarescale.NewClient(server.URL, token)

	// when
	projectsList, err := cli.ListProjects()

	// then
	expectedError := "An unexpected error occurred (code: 500)"
	if fmt.Sprintf("%s", err) != expectedError {
		t.Fatalf("Error is not raised with `%s`", expectedError)
	}

	if projectsList != nil {
		t.Fatalf("Expected projectsList nil, but got another result")
	}
}

func nominalCaseGetProject(t *testing.T) {
	// given
	token := "some-token"
	projectUUID := "8cfe8f68-cad5-4157-b8a6-d9efa12caf0e"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		checkPath(t, "/projects/"+projectUUID, r.URL.Path)
		checkAuthorization(t, r.Header.Get("Authorization"), token)

		resBody := `
		{
			"id": 11,
			"name": "tera-project",
			"uuid": "8cfe8f68-cad5-4157-b8a6-d9efa12caf0e",
			"next": "https://cloud.squarescale.io/projects/ba90e5fe-f520-4275-897b-49a95c1157a3",
			"load_balancer": true,
			"db_enabled": false,
			"db_engine": "postgres",
			"db_size": "dev",
			"db_version": null,
			"database": {
				"enabled": false,
				"engine": "postgres",
				"size": "dev"
			},
			"infra_status": "ok",
			"cluster_size": 1,
			"task": {}
		}
		`

		w.Header().Set("Content-Type", "application/json")

		w.Write([]byte(resBody))
	}))

	defer server.Close()
	cli := squarescale.NewClient(server.URL, token)

	// when
	project, err := cli.GetProject(projectUUID)

	// then
	var expectedInt int
	var expectedString string

	if err != nil {
		t.Fatalf("Expect no error, got `%s`", err)
	}

	expectedString = projectUUID
	if project.UUID != expectedString {
		t.Fatalf("Expect project.UUID `%s`, got `%s`", expectedString, project.UUID)
	}

	expectedString = "tera-project"
	if project.Name != expectedString {
		t.Fatalf("Expect project.Name `%s`, got `%s`", expectedString, project.Name)
	}

	expectedString = "ok"
	if project.InfraStatus != expectedString {
		t.Fatalf("Expect project.InfraStatus `%s`, got `%s`", expectedString, project.InfraStatus)
	}

	expectedInt = 1
	if project.ClusterSize != expectedInt {
		t.Fatalf("Expect project.ClusterSize `%d`, got `%d`", expectedInt, project.ClusterSize)
	}
}

func nominalCaseWaitProject(t *testing.T) {
	// given
	token := "some-token"
	projectUUID := "8cfe8f68-cad5-4157-b8a6-d9efa12caf0e"
	httptestCount := 0

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		checkPath(t, "/projects/"+projectUUID, r.URL.Path)
		checkAuthorization(t, r.Header.Get("Authorization"), token)

		var projectStatus string

		httptestCount++

		if httptestCount < 2 {
			projectStatus = "provisionning"
		} else {
			projectStatus = "ok"
		}
		resBody := fmt.Sprintf(`
		{
			"id": 11,
			"name": "tera-project",
			"uuid": "8cfe8f68-cad5-4157-b8a6-d9efa12caf0e",
			"next": "https://cloud.squarescale.io/projects/ba90e5fe-f520-4275-897b-49a95c1157a3",
			"load_balancer": true,
			"db_enabled": false,
			"db_engine": "postgres",
			"db_size": "dev",
			"db_version": null,
			"database": {
				"enabled": false,
				"engine": "postgres",
				"size": "dev"
			},
			"infra_status": "%s",
			"cluster_size": 1,
			"task": {}
		}
		`, projectStatus)

		w.Header().Set("Content-Type", "application/json")

		w.Write([]byte(resBody))
	}))

	defer server.Close()
	cli := squarescale.NewClient(server.URL, token)

	// when
	_, err := cli.WaitProject(projectUUID, 0)

	// then
	if err != nil {
		t.Fatalf("Expect no error, got `%s`", err)
	}
}

func UnknownProjectOnGetProject(t *testing.T) {
	// given
	token := "some-token"
	projectUUID := "00000000-0000-0000-0000-000000000000"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		checkPath(t, "/projects/"+projectUUID, r.URL.Path)
		checkAuthorization(t, r.Header.Get("Authorization"), token)

		resBody := `{"error":"Couldn't find Project with [WHERE \"projects\".\"uuid\" = $1]"}`

		w.Header().Set("Content-Type", "application/json")

		w.WriteHeader(404)
		w.Write([]byte(resBody))
	}))

	defer server.Close()
	cli := squarescale.NewClient(server.URL, token)

	// when
	project, err := cli.GetProject(projectUUID)

	// then
	expectedError := "Project '00000000-0000-0000-0000-000000000000' not found"
	if err == nil {
		t.Fatalf("Error is not raised with `%s`", expectedError)
	}

	if fmt.Sprintf("%s", err) != expectedError {
		t.Fatalf("Expected error message:\n`%s`\nGot:\n`%s`", expectedError, err)
	}

	if project != nil {
		t.Fatalf("Project is not nil")
	}
}

func badHttpErrrorCoreCaseGetProject(t *testing.T) {
	// given
	token := "some-token"
	projectUUID := "8cfe8f68-cad5-4157-b8a6-d9efa12caf0e"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		checkPath(t, "/projects/"+projectUUID, r.URL.Path)
		checkAuthorization(t, r.Header.Get("Authorization"), token)

		resBody := ``

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(500)

		w.Write([]byte(resBody))
	}))

	defer server.Close()
	cli := squarescale.NewClient(server.URL, token)

	// when
	project, err := cli.GetProject(projectUUID)

	// then
	expectedError := "An unexpected error occurred (code: 500)"
	if fmt.Sprintf("%s", err) != expectedError {
		t.Fatalf("Error is not raised with `%s`", expectedError)
	}

	if project != nil {
		t.Fatalf("Expected project nil, but got another result")
	}
}

func nominalCaseProjectUnprovision(t *testing.T) {
	// given
	token := "some-token"
	projectUUID := "8cfe8f68-cad5-4157-b8a6-d9efa12caf0e"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		checkPath(t, "/projects/"+projectUUID+"/unprovision", r.URL.Path)
		checkAuthorization(t, r.Header.Get("Authorization"), token)

		resBody := ``

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(204)

		w.Write([]byte(resBody))
	}))

	defer server.Close()
	cli := squarescale.NewClient(server.URL, token)

	// when
	err := cli.ProjectUnprovision(projectUUID)

	// then
	if err != nil {
		t.Fatalf("Expect no error, got `%s`", err)
	}
}

func UnknownProjectOnProjectUnprovision(t *testing.T) {
	// given
	token := "some-token"
	projectUUID := "00000000-0000-0000-0000-000000000000"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		checkPath(t, "/projects/"+projectUUID+"/unprovision", r.URL.Path)
		checkAuthorization(t, r.Header.Get("Authorization"), token)

		resBody := `{"error":"Couldn't find Project with [WHERE \"projects\".\"uuid\" = $1]"}`

		w.Header().Set("Content-Type", "application/json")

		w.WriteHeader(404)
		w.Write([]byte(resBody))
	}))

	defer server.Close()
	cli := squarescale.NewClient(server.URL, token)

	// when
	err := cli.ProjectUnprovision(projectUUID)

	// then
	expectedError := "Project '00000000-0000-0000-0000-000000000000' not found"
	if err == nil {
		t.Fatalf("Error is not raised with `%s`", expectedError)
	}
}

func badHttpErrrorCoreCaseOnProjectUnprovision(t *testing.T) {
	// given
	token := "some-token"
	projectUUID := "8cfe8f68-cad5-4157-b8a6-d9efa12caf0e"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		checkPath(t, "/projects/"+projectUUID+"/unprovision", r.URL.Path)
		checkAuthorization(t, r.Header.Get("Authorization"), token)

		resBody := ``

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(500)

		w.Write([]byte(resBody))
	}))

	defer server.Close()
	cli := squarescale.NewClient(server.URL, token)

	// when
	err := cli.ProjectUnprovision(projectUUID)

	// then
	expectedError := "An unexpected error occurred (code: 500)"
	if fmt.Sprintf("%s", err) != expectedError {
		t.Fatalf("Error is not raised with `%s`", expectedError)
	}
}

func nominalCaseProjectDelete(t *testing.T) {
	// given
	token := "some-token"
	projectUUID := "8cfe8f68-cad5-4157-b8a6-d9efa12caf0e"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		checkPath(t, "/projects/"+projectUUID, r.URL.Path)
		checkAuthorization(t, r.Header.Get("Authorization"), token)

		resBody := ``

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(204)

		w.Write([]byte(resBody))
	}))

	defer server.Close()
	cli := squarescale.NewClient(server.URL, token)

	// when
	err := cli.ProjectDelete(projectUUID)

	// then
	if err != nil {
		t.Fatalf("Expect no error, got `%s`", err)
	}
}

func UnknownProjectOnProjectDelete(t *testing.T) {
	// given
	token := "some-token"
	projectUUID := "00000000-0000-0000-0000-000000000000"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		checkPath(t, "/projects/"+projectUUID, r.URL.Path)
		checkAuthorization(t, r.Header.Get("Authorization"), token)

		resBody := `{"error":"Couldn't find Project with [WHERE \"projects\".\"uuid\" = $1]"}`

		w.Header().Set("Content-Type", "application/json")

		w.WriteHeader(404)
		w.Write([]byte(resBody))
	}))

	defer server.Close()
	cli := squarescale.NewClient(server.URL, token)

	// when
	err := cli.ProjectDelete(projectUUID)

	// then
	expectedError := "Project '00000000-0000-0000-0000-000000000000' not found"
	if err == nil {
		t.Fatalf("Error is not raised with `%s`", expectedError)
	}
}

func badHttpErrrorCoreCaseOnProjectDelete(t *testing.T) {
	// given
	token := "some-token"
	projectUUID := "8cfe8f68-cad5-4157-b8a6-d9efa12caf0e"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		checkPath(t, "/projects/"+projectUUID, r.URL.Path)
		checkAuthorization(t, r.Header.Get("Authorization"), token)

		resBody := ``

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(500)

		w.Write([]byte(resBody))
	}))

	defer server.Close()
	cli := squarescale.NewClient(server.URL, token)

	// when
	err := cli.ProjectDelete(projectUUID)

	// then
	expectedError := "An unexpected error occurred (code: 500)"
	if fmt.Sprintf("%s", err) != expectedError {
		t.Fatalf("Error is not raised with `%s`", expectedError)
	}
}

func ClientHTTPErrorOnProjectMethods(t *testing.T) {
	// given
	token := "some-token"
	projectUUID := "ba90e5fe-f520-4275-897b-49a95c1157a3"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	}))

	defer server.Close()
	cli := squarescale.NewClient(server.URL, token)

	// when
	payload := squarescale.JSONObject{}
	payload["name"] = "project-test"
	payload["credential_name"] = "aws-production"
	payload["provider"] = "aws"
	payload["region"] = "eu-west-1"
	payload["infra_type"] = "single_node"
	payload["node_size"] = "dev"
	payload["slack_webhook"] = ""
	payload["databases"] = "[]"
	_, errOnCreate := cli.CreateProject(&payload)
	errOnProvision := cli.ProvisionProject(projectUUID)
	errOnUNProvision := cli.UNProvisionProject(projectUUID)
	_, errOnList := cli.ListProjects()
	_, errOnGet := cli.GetProject(projectUUID)
	_, errOnWait := cli.WaitProject(projectUUID, 0)
	errOnUnprovision := cli.ProjectUnprovision(projectUUID)
	errOnDelete := cli.ProjectDelete(projectUUID)

	// then
	if errOnCreate == nil {
		t.Errorf("Error is not raised on CreateProject")
	}

	if errOnProvision == nil {
		t.Errorf("Error is not raised on ProvisionProject")
	}

	if errOnUNProvision == nil {
		t.Errorf("Error is not raised on UNProvisionProject")
	}

	if errOnList == nil {
		t.Errorf("Error is not raised on ListProjects")
	}

	if errOnGet == nil {
		t.Errorf("Error is not raised on GetProject")
	}

	if errOnWait == nil {
		t.Errorf("Error is not raised on WaitProject")
	}

	if errOnUnprovision == nil {
		t.Errorf("Error is not raised on ProjectUnprovision")
	}

	if errOnDelete == nil {
		t.Errorf("Error is not raised on ProjectDelete")
	}
}

func CantUnmarshalOnProjectMethods(t *testing.T) {
	// given
	token := "some-token"
	projectUUID := "ba90e5fe-f520-4275-897b-49a95c1157a3"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resBody := `{]`

		w.Header().Set("Content-Type", "application/json")
		if r.Method == "POST" && r.URL.Path == "/projects" {
			w.WriteHeader(201)
		}

		w.Write([]byte(resBody))
	}))

	defer server.Close()
	cli := squarescale.NewClient(server.URL, token)

	// when
	payload := squarescale.JSONObject{}
	payload["name"] = "project-test"
	payload["credential_name"] = "aws-production"
	payload["provider"] = "aws"
	payload["region"] = "eu-west-1"
	payload["infra_type"] = "single_node"
	payload["node_size"] = "dev"
	payload["slack_webhook"] = ""
	payload["databases"] = "[]"
	_, errOnCreate := cli.CreateProject(&payload)
	_, errOnList := cli.ListProjects()
	_, errOnGet := cli.GetProject(projectUUID)

	// then
	expectedError := "invalid character ']' looking for beginning of object key string"
	if errOnCreate == nil {
		t.Fatalf("Error is not raised with `%s`", expectedError)
	}

	if fmt.Sprintf("%s", errOnCreate) != expectedError {
		t.Fatalf("Expected error message:\n`%s`\nGot:\n`%s`", expectedError, errOnCreate)
	}

	if fmt.Sprintf("%s", errOnList) != expectedError {
		t.Fatalf("Expected error message:\n`%s`\nGot:\n`%s`", expectedError, errOnList)
	}

	if fmt.Sprintf("%s", errOnGet) != expectedError {
		t.Fatalf("Expected error message:\n`%s`\nGot:\n`%s`", expectedError, errOnGet)
	}
}

// WaitProject(projectUUID string) (*ProjectStatus, error) {
// ProjectUnprovision(project string) error {
// ProjectDelete(project string) error {
// ProjectLogs(project string, container string, after string) ([]string, string, error) {
