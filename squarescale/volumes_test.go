package squarescale_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/squarescale/squarescale-cli/squarescale"
)

func TestGetVolumes(t *testing.T) {
	// GetVolumes
	t.Run("Nominal case on GetVolumes", nominalCaseOnGetVolumes)

	// GetVolumeInfo
	t.Run("Nominal case on GetVolumeInfo", nominalCaseOnGetVolumeInfo)

	t.Run("Test volume not found on GetVolumeInfo", VolumeNotFoundOnGetVolumeInfo)
	t.Run("Test project not found on GetVolumeInfo", ProjectNotFoundOnGetVolumeInfo)

	// AddVolume
	t.Run("Nominal case on AddVolume", nominalCaseOnAddVolume)

	t.Run("Test duplicate volume name on AddVolume", DuplicateVolumeErrorCaseOnAddVolume)
	t.Run("Test project not found on AddVolume", ProjectNotFoundOnAddVolume)

	// DeleteVolume
	t.Run("Nominal case on DeleteVolume", nominalCaseOnDeleteVolume)

	t.Run("Test volume not found on DeleteVolume", VolumeNotFoundOnDeleteVolume)
	t.Run("Test project not found on DeleteVolume", ProjectNotFoundOnDeleteVolume)

	// WaitVolume
	t.Run("Nominal case on WaitVolume", nominalCaseOnWaitVolume)

	// Error cases
	t.Run("Test HTTP client error on volume methods (get, add, delete and wait)", ClientHTTPErrorOnVolumeMethods)
	t.Run("Test internal server error on volume methods (get, add, delete and wait)", InternalServerErrorOnVolumeMethods)
	t.Run("Test badly JSON on volume methods (get)", CantUnmarshalOnVolumeMethods)
}

func nominalCaseOnGetVolumes(t *testing.T) {
	// given
	token := "some-token"
	projectName := "my-project"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		checkPath(t, "/projects/"+projectName+"/volumes", r.URL.Path)
		checkAuthorization(t, r.Header.Get("Authorization"), token)

		resBody := `
		[
			{
				"id": 30,
				"name": "vol02a",
				"size": 1,
				"type": "gp2",
				"zone": "eu-west-1a",
				"statefull_node_name": "node02",
				"status": "provisionned"
			},
			{
				"id": 31,
				"name": "vol02b",
				"size": 3,
				"type": "io1",
				"zone": "eu-west-1b",
				"statefull_node_name": null,
				"status": "not_provisionned"
			}
		]
		`

		w.Header().Set("Content-Type", "application/json")

		w.Write([]byte(resBody))
	}))

	defer server.Close()
	cli := squarescale.NewClient(server.URL, token)

	// when
	volumes, err := cli.GetVolumes(projectName)

	// then
	var expectedInt int
	var expectedString string

	if err != nil {
		t.Fatalf("Expect no error, got `%s`", err)
	}

	expectedInt = 2
	if len(volumes) != expectedInt {
		t.Fatalf("Expect volumes to contain one element `%d`, but got actually `%d`", expectedInt, len(volumes))
	}

	expectedInt = 30
	if volumes[0].ID != expectedInt {
		t.Errorf("Expect volume.ID `%d`, got `%d`", expectedInt, volumes[0].ID)
	}

	expectedString = "vol02a"
	if volumes[0].Name != expectedString {
		t.Errorf("Expect volume.Name `%s`, got `%s`", expectedString, volumes[0].Name)
	}

	expectedInt = 1
	if volumes[0].Size != expectedInt {
		t.Errorf("Expect volume.Size %d , got %d", expectedInt, volumes[0].Size)
	}

	expectedString = "gp2"
	if volumes[0].Type != expectedString {
		t.Errorf("Expect volume.Type `%s`, got `%s`", expectedString, volumes[0].Type)
	}

	expectedString = "eu-west-1a"
	if volumes[0].Zone != expectedString {
		t.Errorf("Expect volume.Zone `%s`, got `%s`", expectedString, volumes[0].Zone)
	}

	expectedString = "node02"
	if volumes[0].StatefullNodeName != expectedString {
		t.Errorf("Expect volume.StatefullNodeName `%s`, got `%s`", expectedString, volumes[0].StatefullNodeName)
	}

	expectedString = "provisionned"
	if volumes[0].Status != expectedString {
		t.Errorf("Expect volume.Status `%s`, got `%s`", expectedString, volumes[0].Status)
	}

	expectedInt = 31
	if volumes[1].ID != expectedInt {
		t.Errorf("Expect volume.ID %d, got %d", expectedInt, volumes[1].ID)
	}

	expectedString = "vol02b"
	if volumes[1].Name != expectedString {
		t.Errorf("Expect volume.Name `%s`, got `%s`", expectedString, volumes[1].Name)
	}

	expectedInt = 3
	if volumes[1].Size != expectedInt {
		t.Errorf("Expect volume.Size %d , got %d", expectedInt, volumes[1].Size)
	}

	expectedString = "io1"
	if volumes[1].Type != expectedString {
		t.Errorf("Expect volume.Type `%s`, got `%s`", expectedString, volumes[1].Type)
	}

	expectedString = "eu-west-1b"
	if volumes[1].Zone != expectedString {
		t.Errorf("Expect volume.Zone `%s`, got `%s`", expectedString, volumes[1].Zone)
	}

	expectedString = ""
	if volumes[1].StatefullNodeName != expectedString {
		t.Errorf("Expect volume.StatefullNodeName `%s`, got `%s`", expectedString, volumes[1].StatefullNodeName)
	}

	expectedString = "not_provisionned"
	if volumes[1].Status != expectedString {
		t.Errorf("Expect volume.Status `%s`, got `%s`", expectedString, volumes[1].Status)
	}
}

func nominalCaseOnGetVolumeInfo(t *testing.T) {
	// given
	token := "some-token"
	projectName := "my-project"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		checkPath(t, "/projects/"+projectName+"/volumes", r.URL.Path)
		checkAuthorization(t, r.Header.Get("Authorization"), token)

		resBody := `
		[
			{
				"id": 30,
				"name": "vol02a",
				"size": 1,
				"type": "gp2",
				"zone": "eu-west-1a",
				"statefull_node_name": "node02",
				"status": "provisionned"
			},
			{
				"id": 31,
				"name": "vol02b",
				"size": 3,
				"type": "io1",
				"zone": "eu-west-1b",
				"statefull_node_name": null,
				"status": "not_provisionned"
			}
		]
		`

		w.Header().Set("Content-Type", "application/json")

		w.Write([]byte(resBody))
	}))

	defer server.Close()
	cli := squarescale.NewClient(server.URL, token)

	// when
	volume, err := cli.GetVolumeInfo(projectName, "vol02a")

	// then
	var expectedInt int
	var expectedString string

	if err != nil {
		t.Fatalf("Expect no error, got `%s`", err)
	}

	expectedString = "vol02a"
	if volume.Name != expectedString {
		t.Errorf("Expect volume.Name `%s`, got `%s`", expectedString, volume.Name)
	}

	expectedInt = 1
	if volume.Size != expectedInt {
		t.Errorf("Expect volume.Size `%d` , got `%d`", expectedInt, volume.Size)
	}

	expectedString = "gp2"
	if volume.Type != expectedString {
		t.Errorf("Expect volume.Type `%s`, got `%s`", expectedString, volume.Type)
	}

	expectedString = "eu-west-1a"
	if volume.Zone != expectedString {
		t.Errorf("Expect volume.Zone `%s`, got `%s`", expectedString, volume.Zone)
	}

	expectedString = "node02"
	if volume.StatefullNodeName != expectedString {
		t.Errorf("Expect volume.StatefullNodeName `%s`, got `%s`", expectedString, volume.StatefullNodeName)
	}

	expectedString = "provisionned"
	if volume.Status != expectedString {
		t.Errorf("Expect volume.Status `%s`, got `%s`", expectedString, volume.Status)
	}
}

func nominalCaseOnAddVolume(t *testing.T) {
	// given
	token := "some-token"
	projectName := "my-project"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		checkPath(t, "/projects/"+projectName+"/volumes", r.URL.Path)
		checkAuthorization(t, r.Header.Get("Authorization"), token)

		resBody := `
		{
			"id": 61,
			"name": "vol02c1",
			"size": 1,
			"type": "gp2",
			"zone": "eu-west-1c",
			"statefull_node_name": null,
			"status": "not_provisionned"
		}
		`

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(201)

		w.Write([]byte(resBody))
	}))

	defer server.Close()
	cli := squarescale.NewClient(server.URL, token)

	// when
	err := cli.AddVolume(projectName, "vol02c1", 1, "gp2", "eu-west-1c")

	// then
	if err != nil {
		t.Fatalf("Expect no error, got `%s`", err)
	}
}

func nominalCaseOnDeleteVolume(t *testing.T) {
	// given
	token := "some-token"
	projectName := "my-project"
	volumeName := "my-volume"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		checkPath(t, "/projects/"+projectName+"/volumes/"+volumeName, r.URL.Path)
		checkAuthorization(t, r.Header.Get("Authorization"), token)

		resBody := `
		null
		`

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)

		w.Write([]byte(resBody))
	}))

	defer server.Close()
	cli := squarescale.NewClient(server.URL, token)

	// when
	err := cli.DeleteVolume(projectName, volumeName)

	// then
	if err != nil {
		t.Fatalf("Expect no error, got `%s`", err)
	}
}

func nominalCaseOnWaitVolume(t *testing.T) {
	// given
	token := "some-token"
	projectName := "my-project"
	volumeName := "my-volume"
	httptestCount := 0

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		checkPath(t, "/projects/"+projectName+"/volumes", r.URL.Path)
		checkAuthorization(t, r.Header.Get("Authorization"), token)

		var volumeStatus string

		httptestCount++

		if httptestCount < 2 {
			volumeStatus = "not_provisionned"
		} else {
			volumeStatus = "provisionned"
		}
		resBody := fmt.Sprintf(`
		[
			{
				"id": 30,
				"name": "my-volume",
				"size": 1,
				"type": "gp2",
				"zone": "eu-west-1a",
				"statefull_node_name": "node02",
				"status": "%s"
			}
		]
		`, volumeStatus)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)

		w.Write([]byte(resBody))
	}))

	defer server.Close()
	cli := squarescale.NewClient(server.URL, token)

	// when
	_, err := cli.WaitVolume(projectName, volumeName, 0)

	// then
	if err != nil {
		t.Fatalf("Expect no error, got `%s`", err)
	}
}

func VolumeNotFoundOnGetVolumeInfo(t *testing.T) {
	// given
	token := "some-token"
	projectName := "my-project"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resBody := `
		[
			{
				"id": 30,
				"name": "vol02a",
				"size": 1,
				"type": "gp2",
				"zone": "eu-west-1a",
				"statefull_node_name": "node02",
				"status": "provisionned"
			},
			{
				"id": 31,
				"name": "vol02b",
				"size": 3,
				"type": "io1",
				"zone": "eu-west-1b",
				"statefull_node_name": null,
				"status": "not_provisionned"
			}
		]
		`

		w.Header().Set("Content-Type", "application/json")

		if (r.Header.Get("Authorization")) != "bearer some-token" {
			t.Fatalf("Wrong token! Expected `%s`, got `%s`", "bearer some-token", r.Header.Get("Authorization"))
		}

		w.Write([]byte(resBody))
	}))

	defer server.Close()
	cli := squarescale.NewClient(server.URL, token)

	// when
	_, err := cli.GetVolumeInfo(projectName, "missing-volume")

	// then
	expectedError := "Volume 'missing-volume' not found for project 'my-project'"
	if err == nil {
		t.Fatalf("Error is not raised with `%s`", expectedError)
	}

	if fmt.Sprintf("%s", err) != expectedError {
		t.Fatalf("Expected error:\n`%s`\nGot:\n`%s`", expectedError, err)
	}
}

func ProjectNotFoundOnGetVolumeInfo(t *testing.T) {
	// given
	token := "some-token"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resBody := `
		{
			"error": "No project found for config name: not-a-project"
		}
		`

		w.Header().Set("Content-Type", "application/json")

		if (r.Header.Get("Authorization")) != "bearer some-token" {
			t.Fatalf("Wrong token! Expected `%s`, got `%s`", "bearer some-token", r.Header.Get("Authorization"))
		}

		w.WriteHeader(404)
		w.Write([]byte(resBody))
	}))

	defer server.Close()
	cli := squarescale.NewClient(server.URL, token)

	// when
	_, err := cli.GetVolumeInfo("not-a-project", "a-volume")

	// then
	expectedError := "Project 'not-a-project' does not exist"
	if err == nil {
		t.Fatalf("Error is not raised with `%s`", expectedError)
	}

	if fmt.Sprintf("%s", err) != expectedError {
		t.Fatalf("Expected error:\n`%s`\nGot:\n`%s`", expectedError, err)
	}
}

func DuplicateVolumeErrorCaseOnAddVolume(t *testing.T) {
	// given
	token := "some-token"
	projectName := "my-project"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resBody := `
		{"error":"PG::UniqueViolation: ERROR:  duplicate key value violates unique constraint \"index_volumes_on_cluster_id_and_name\"\nDETAIL:  Key (cluster_id, name)=(6163, vol02c1) already exists.\n: INSERT INTO \"volumes\" (\"name\", \"size\", \"type\", \"zone\", \"cluster_id\", \"status\") VALUES ($1, $2, $3, $4, $5, $6) RETURNING \"id\""}
		`

		w.Header().Set("Content-Type", "application/json")

		if (r.Header.Get("Authorization")) != "bearer some-token" {
			t.Fatalf("Wrong token! Expected `%s`, got `%s`", "bearer some-token", r.Header.Get("Authorization"))
		}

		w.WriteHeader(409)
		w.Write([]byte(resBody))
	}))

	defer server.Close()
	cli := squarescale.NewClient(server.URL, token)

	// when
	err := cli.AddVolume(projectName, "vol02c", 1, "gp2", "eu-west-1c")

	// then
	expectedError := "Volume already exist on project 'my-project': vol02c"
	if err == nil {
		t.Fatalf("Error is not raised with `%s`", expectedError)
	}

	if fmt.Sprintf("%s", err) != expectedError {
		t.Fatalf("Expected error:\n`%s`\nGot:\n`%s`", expectedError, err)
	}
}

func ProjectNotFoundOnAddVolume(t *testing.T) {
	// given
	token := "some-token"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resBody := `
		{"error":"No project found for config name: another-project"}
		`

		w.Header().Set("Content-Type", "application/json")

		if (r.Header.Get("Authorization")) != "bearer some-token" {
			t.Fatalf("Wrong token! Expected `%s`, got `%s`", "bearer some-token", r.Header.Get("Authorization"))
		}

		w.WriteHeader(404)
		w.Write([]byte(resBody))
	}))

	defer server.Close()
	cli := squarescale.NewClient(server.URL, token)

	// when
	err := cli.AddVolume("another-project", "vol02c1", 1, "gp2", "eu-west-1c")

	// then
	expectedError := `1 error occurred:

* error: No project found for config name: another-project`
	if err == nil {
		t.Fatalf("Error is not raised with `%s`", expectedError)
	}

	if fmt.Sprintf("%s", err) != expectedError {
		t.Fatalf("Expected error:\n`%s`\nGot:\n`%s`", expectedError, err)
	}
}

func ProjectNotFoundOnDeleteVolume(t *testing.T) {
	// given
	token := "some-token"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resBody := `{"error":"No project found for config name: unknown-project"}`

		w.Header().Set("Content-Type", "application/json")

		if (r.Header.Get("Authorization")) != "bearer some-token" {
			t.Fatalf("Wrong token! Expected `%s`, got `%s`", "bearer some-token", r.Header.Get("Authorization"))
		}

		w.WriteHeader(404)
		w.Write([]byte(resBody))
	}))

	defer server.Close()
	cli := squarescale.NewClient(server.URL, token)

	// when
	err := cli.DeleteVolume("unknown-project", "vol02c1")

	// then
	expectedError := `{"error":"No project found for config name: unknown-project"}`
	if err == nil {
		t.Fatalf("Error is not raised with `%s`", expectedError)
	}

	if fmt.Sprintf("%s", err) != expectedError {
		t.Fatalf("Expected error:\n`%s`\nGot:\n`%s`", expectedError, err)
	}
}

func VolumeNotFoundOnDeleteVolume(t *testing.T) {
	// given
	token := "some-token"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resBody := `{"error":"Couldn't find Volume with [WHERE \"volumes\".\"cluster_id\" = $1 AND \"volumes\".\"name\" = $2]"}`

		w.Header().Set("Content-Type", "application/json")

		if (r.Header.Get("Authorization")) != "bearer some-token" {
			t.Fatalf("Wrong token! Expected `%s`, got `%s`", "bearer some-token", r.Header.Get("Authorization"))
		}

		w.WriteHeader(404)
		w.Write([]byte(resBody))
	}))

	defer server.Close()
	cli := squarescale.NewClient(server.URL, token)

	// when
	err := cli.DeleteVolume("my-project", "vol02c1")

	// then
	expectedError := `{"error":"No volume found for name: vol02c1"}`
	if err == nil {
		t.Fatalf("Error is not raised with `%s`", expectedError)
	}

	if fmt.Sprintf("%s", err) != expectedError {
		t.Fatalf("Expected error:\n`%s`\nGot:\n`%s`", expectedError, err)
	}
}

func ClientHTTPErrorOnVolumeMethods(t *testing.T) {
	// given
	token := "some-token"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	}))

	defer server.Close()
	cli := squarescale.NewClient(server.URL, token)

	// when
	_, errOnGet := cli.GetVolumes("another-project")
	errOnAdd := cli.AddVolume("another-project", "vol02c1", 1, "gp2", "eu-west-1c")
	errOnDelete := cli.DeleteVolume("my-project", "vol02c1")
	_, errOnWait := cli.WaitVolume("my-project", "vol02c1", 0)

	// then
	if errOnGet == nil {
		t.Errorf("Error is not raised on GetVolumes")
	}

	if errOnAdd == nil {
		t.Errorf("Error is not raised on AddVolume")
	}

	if errOnDelete == nil {
		t.Errorf("Error is not raised on DeleteVolume")
	}

	if errOnWait == nil {
		t.Errorf("Error is not raised on WaitVolume")
	}
}

func InternalServerErrorOnVolumeMethods(t *testing.T) {
	// given
	token := "some-token"
	projectName := "my-project"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(500)
	}))

	defer server.Close()
	cli := squarescale.NewClient(server.URL, token)

	// when
	_, errOnGetVolumes := cli.GetVolumes(projectName)
	errOnAddVolume := cli.AddVolume(projectName, "vol02c1", 1, "gp2", "eu-west-1c")
	errOnDeleteVolume := cli.DeleteVolume(projectName, "vol02c1")

	// then
	expectedError := "An unexpected error occurred (code: 500)"
	if errOnGetVolumes == nil {
		t.Errorf("Error is not raised with `%s`", expectedError)
	}

	if fmt.Sprintf("%s", errOnGetVolumes) != expectedError {
		t.Errorf("Expected error:\n`%s`\nGot:\n`%s`", expectedError, errOnGetVolumes)
	}

	if errOnAddVolume == nil {
		t.Errorf("Error is not raised with `%s`", expectedError)
	}

	if fmt.Sprintf("%s", errOnAddVolume) != expectedError {
		t.Errorf("Expected error:\n`%s`\nGot:\n`%s`", expectedError, errOnAddVolume)
	}

	if errOnDeleteVolume == nil {
		t.Errorf("Error is not raised with `%s`", expectedError)
	}

	if fmt.Sprintf("%s", errOnDeleteVolume) != expectedError {
		t.Errorf("Expected error:\n`%s`\nGot:\n`%s`", expectedError, errOnDeleteVolume)
	}
}

func CantUnmarshalOnVolumeMethods(t *testing.T) {
	// given
	token := "some-token"
	projectName := "my-project"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resBody := `{]`

		w.Header().Set("Content-Type", "application/json")

		if (r.Header.Get("Authorization")) != "bearer some-token" {
			t.Fatalf("Wrong token! Expected `%s`, got `%s`", "bearer some-token", r.Header.Get("Authorization"))
		}

		w.Write([]byte(resBody))
	}))

	defer server.Close()
	cli := squarescale.NewClient(server.URL, token)

	// when
	_, errOnGetVolumes := cli.GetVolumes(projectName)

	// then
	expectedError := "invalid character ']' looking for beginning of object key string"
	if errOnGetVolumes == nil {
		t.Fatalf("Error is not raised with `%s`", expectedError)
	}

	if fmt.Sprintf("%s", errOnGetVolumes) != expectedError {
		t.Fatalf("Expected error:\n`%s`\nGot:\n`%s`", expectedError, errOnGetVolumes)
	}
}
