package squarescale_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/squarescale/squarescale-cli/squarescale"
)

func TestGetVolumes(t *testing.T) {

	// nominal case
	t.Run("nominal get volumes", nominalCaseForVolumes)
	t.Run("nominal volume info", nominalCaseForVolumeInfo)
	t.Run("nominal volume add", nominalCaseForVolumeAdd)
	t.Run("nominal volume delete", nominalCaseForVolumeDelete)
	t.Run("nominal volume wait", nominalCaseForWaitVolume)

	// other cases
	t.Run("test badly JSON", CantUnmarshal)                                 // cli.GetVolumes
	t.Run("test Internal Server Error ", InternalServerErrorCaseForVolumes) // cli.GetVolumes
	t.Run("test HTTP error", UnexpectedErrorOnGet)                          // cli.GetVolumes

	t.Run("test Not Found Page", NotFoundCaseForVolume)  // cli.GetVolumeInfo
	t.Run("test Not Found Page", NotFoundCaseForProject) // cli.GetVolumeInfo

	t.Run("test Duplicate volume name on VolumeAdd", DuplicateVolumeErrorCaseForVolumeAdd) // cli.AddVolume
	t.Run("test Unknown project on VolumeAdd", UnknownProjectErrorCaseForVolumeAdd)        // cli.AddVolume
	t.Run("test Internal Server Error on VolumeAdd", InternalServerErrorCaseForVolumeAdd)  // cli.AddVolume
	t.Run("test HTTP Error on VolumeAdd", UnexpectedHTTPErrorVolumeAdd)                    // cli.AddVolume

	t.Run("test Unknown project on VolumeDelete", UnknownProjectErrorCaseForVolumeDelete)       // cli.DeleteVolume
	t.Run("test Unknown project on VolumeDelete", UnknownVolumeErrorCaseForVolumeDelete)        // cli.DeleteVolume
	t.Run("test Internal Server Error on VolumeDelete", InternalServerErrorCaseForVolumeDelete) // cli.DeleteVolume
	t.Run("test HTTP Error on VolumeDelete", UnexpectedHTTPErrorVolumeDelete)                   // cli.DeleteVolume

	t.Run("test Internal Server Error on WaitVolume", InternalServerErrorCaseForWaitVolume) // cli.DeleteVolume
}

func nominalCaseForVolumes(t *testing.T) {
	// given
	token := "some-token"
	projectName := "my-project"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var path string = r.URL.Path

		if path != "/projects/"+projectName+"/volumes" {
			t.Fatalf("Wrong path ! Expected %s, got %s", "/projects/my-project/volumes", path)
		}

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
			t.Fatalf("Wrong path ! Expected %s, got %s", "bearer some-token", r.Header.Get("Authorization"))
		}

		w.Write([]byte(resBody))
	}))

	defer server.Close()
	cli := squarescale.NewClient(server.URL, token)

	// when
	volumes, err := cli.GetVolumes(projectName)

	if err != nil {
		t.Fatalf("Expect no error, got %s", err)
	}

	if len(volumes) != 2 {
		t.Fatalf("Expect volumes to contain one element %d, but got actually %d", 1, len(volumes))
	}

	if volumes[0].Name != "vol02a" {
		t.Errorf("Expect volumeName %s , got %s", "vol02a", volumes[0].Name)
	}

	if volumes[0].Size != 1 {
		t.Errorf("Expect volumeSize %d , got %d", 1, volumes[0].Size)
	}

	if volumes[0].Type != "gp2" {
		t.Errorf("Expect volumeType %s , got %s", "gp2", volumes[0].Type)
	}

	if volumes[0].Zone != "eu-west-1a" {
		t.Errorf("Expect volumeZone %s , got %s", "eu-west-1a", volumes[0].Zone)
	}

	if volumes[0].StatefullNodeName != "node02" {
		t.Errorf("Expect volumeStatefullNodeName %s , got %s", "node02", volumes[0].StatefullNodeName)
	}

	if volumes[0].Status != "provisionned" {
		t.Errorf("Expect volumeStatus %s , got %s", "provisionned", volumes[0].Status)
	}

	if volumes[1].Name != "vol02b" {
		t.Errorf("Expect volumeName %s , got %s", "vol02b", volumes[1].Name)
	}

	if volumes[1].Size != 3 {
		t.Errorf("Expect volumeSize %d , got %d", 3, volumes[1].Size)
	}

	if volumes[1].Type != "io1" {
		t.Errorf("Expect volumeType %s , got %s", "io1", volumes[1].Type)
	}

	if volumes[1].Zone != "eu-west-1b" {
		t.Errorf("Expect volumeZone %s , got %s", "eu-west-1b", volumes[1].Zone)
	}

	if volumes[1].StatefullNodeName != "" {
		t.Errorf("Expect volumeStatefullNodeName %s , got %s", "", volumes[1].StatefullNodeName)
	}

	if volumes[1].Status != "not_provisionned" {
		t.Errorf("Expect volumeStatus %s , got %s", "not_provisionned", volumes[1].Status)
	}
}

func nominalCaseForVolumeInfo(t *testing.T) {
	// given
	token := "some-token"
	projectName := "my-project"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var path string = r.URL.Path

		if path != "/projects/"+projectName+"/volumes" {
			t.Fatalf("Wrong path ! Expected %s, got %s", "/projects/my-project/volumes", path)
		}

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
			t.Fatalf("Wrong path ! Expected %s, got %s", "bearer some-token", r.Header.Get("Authorization"))
		}

		w.Write([]byte(resBody))
	}))

	defer server.Close()
	cli := squarescale.NewClient(server.URL, token)

	// when
	volume, err := cli.GetVolumeInfo(projectName, "vol02a")

	if err != nil {
		t.Fatalf("Expect no error, got %s", err)
	}

	if volume.Name != "vol02a" {
		t.Errorf("Expect volumeName %s , got %s", "vol02a", volume.Name)
	}

	if volume.Size != 1 {
		t.Errorf("Expect volumeSize %d , got %d", 1, volume.Size)
	}

	if volume.Type != "gp2" {
		t.Errorf("Expect volumeType %s , got %s", "gp2", volume.Type)
	}

	if volume.Zone != "eu-west-1a" {
		t.Errorf("Expect volumeZone %s , got %s", "eu-west-1a", volume.Zone)
	}

	if volume.StatefullNodeName != "node02" {
		t.Errorf("Expect volumeStatefullNodeName %s , got %s", "node02", volume.StatefullNodeName)
	}

	if volume.Status != "provisionned" {
		t.Errorf("Expect volumeStatus %s , got %s", "provisionned", volume.Status)
	}
}

func nominalCaseForVolumeAdd(t *testing.T) {
	// given
	token := "some-token"
	projectName := "my-project"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var path string = r.URL.Path

		if path != "/projects/"+projectName+"/volumes" {
			t.Fatalf("Wrong path ! Expected %s, got %s", "/projects/my-project/volumes", path)
		}

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

		if (r.Header.Get("Authorization")) != "bearer some-token" {
			t.Fatalf("Wrong path ! Expected %s, got %s", "bearer some-token", r.Header.Get("Authorization"))
		}

		w.Write([]byte(resBody))
	}))

	defer server.Close()
	cli := squarescale.NewClient(server.URL, token)

	// when
	err := cli.AddVolume(projectName, "vol02c1", 1, "gp2", "eu-west-1c")

	if err != nil {
		t.Fatalf("Expect no error, got %s", err)
	}
}

func nominalCaseForVolumeDelete(t *testing.T) {
	// given
	token := "some-token"
	projectName := "my-project"
	volumeName := "my-volume"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var path string = r.URL.Path

		if path != "/projects/"+projectName+"/volumes/"+volumeName {
			t.Fatalf("Wrong path ! Expected %s, got %s", "/projects/my-project/volumes/my-volume", path)
		}

		resBody := `
		null
		`

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)

		if (r.Header.Get("Authorization")) != "bearer some-token" {
			t.Fatalf("Wrong path ! Expected %s, got %s", "bearer some-token", r.Header.Get("Authorization"))
		}

		w.Write([]byte(resBody))
	}))

	defer server.Close()
	cli := squarescale.NewClient(server.URL, token)

	// when
	err := cli.DeleteVolume(projectName, volumeName)

	if err != nil {
		t.Fatalf("Expect no error, got %s", err)
	}
}

func nominalCaseForWaitVolume(t *testing.T) {
	// given
	token := "some-token"
	projectName := "my-project"
	volumeName := "my-volume"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var path string = r.URL.Path

		if path != "/projects/"+projectName+"/volumes" {
			t.Fatalf("Wrong path ! Expected %s, got %s", "/projects/my-project/volumes", path)
		}

		resBody := `
		[
			{
				"id": 30,
				"name": "my-volume",
				"size": 1,
				"type": "gp2",
				"zone": "eu-west-1a",
				"statefull_node_name": "node02",
				"status": "provisionned"
			}
		]
		`

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)

		if (r.Header.Get("Authorization")) != "bearer some-token" {
			t.Fatalf("Wrong path ! Expected %s, got %s", "bearer some-token", r.Header.Get("Authorization"))
		}

		w.Write([]byte(resBody))
	}))

	defer server.Close()
	cli := squarescale.NewClient(server.URL, token)

	// when
	_, err := cli.WaitVolume(projectName, volumeName)

	if err != nil {
		t.Fatalf("Expect no error, got %s", err)
	}
}

func NotFoundCaseForVolume(t *testing.T) {
	// given
	token := "some-token"
	projectName := "my-project"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var path string = r.URL.Path

		if path != "/projects/"+projectName+"/volumes" {
			t.Fatalf("Wrong path ! Expected %s, got %s", "/projects/my-project/volumes", path)
		}

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
			t.Fatalf("Wrong path ! Expected %s, got %s", "bearer some-token", r.Header.Get("Authorization"))
		}

		w.Write([]byte(resBody))
	}))

	defer server.Close()
	cli := squarescale.NewClient(server.URL, token)

	// when
	_, err := cli.GetVolumeInfo(projectName, "missing-volume")

	if err == nil {
		t.Fatalf("Error is not raised with `Volume 'missing-volume' not found for project 'my-project'`")
	}
}

func NotFoundCaseForProject(t *testing.T) {
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
			t.Fatalf("Wrong path ! Expected %s, got %s", "bearer some-token", r.Header.Get("Authorization"))
		}

		w.WriteHeader(404)
		w.Write([]byte(resBody))
	}))

	defer server.Close()
	cli := squarescale.NewClient(server.URL, token)

	// when
	_, err := cli.GetVolumeInfo("/projects/not-a-project/volumes", "a-volume")

	if err == nil {
		t.Fatalf("Error is not raised with `Volume 'missing-volume' not found for project 'my-project'`")
	}
}

func DuplicateVolumeErrorCaseForVolumeAdd(t *testing.T) {
	// given
	token := "some-token"
	projectName := "my-project"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resBody := `
		{"error":"PG::UniqueViolation: ERROR:  duplicate key value violates unique constraint \"index_volumes_on_cluster_id_and_name\"\nDETAIL:  Key (cluster_id, name)=(6163, vol02c1) already exists.\n: INSERT INTO \"volumes\" (\"name\", \"size\", \"type\", \"zone\", \"cluster_id\", \"status\") VALUES ($1, $2, $3, $4, $5, $6) RETURNING \"id\""}
		`

		w.Header().Set("Content-Type", "application/json")

		if (r.Header.Get("Authorization")) != "bearer some-token" {
			t.Fatalf("Wrong path ! Expected %s, got %s", "bearer some-token", r.Header.Get("Authorization"))
		}

		w.WriteHeader(409)
		w.Write([]byte(resBody))
	}))

	defer server.Close()
	cli := squarescale.NewClient(server.URL, token)

	// when
	err := cli.AddVolume(projectName, "vol02c", 1, "gp2", "eu-west-1c")

	if err == nil {
		t.Fatalf("Error is not raised")
	}
}

func UnknownProjectErrorCaseForVolumeAdd(t *testing.T) {
	// given
	token := "some-token"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resBody := `
		{"error":"No project found for config name: another-project"}
		`

		w.Header().Set("Content-Type", "application/json")

		if (r.Header.Get("Authorization")) != "bearer some-token" {
			t.Fatalf("Wrong path ! Expected %s, got %s", "bearer some-token", r.Header.Get("Authorization"))
		}

		w.WriteHeader(404)
		w.Write([]byte(resBody))
	}))

	defer server.Close()
	cli := squarescale.NewClient(server.URL, token)

	// when
	err := cli.AddVolume("another-project", "vol02c1", 1, "gp2", "eu-west-1c")

	if err == nil {
		t.Fatalf("Error is not raised %s", err)
	}
}

func InternalServerErrorCaseForVolumeAdd(t *testing.T) {
	// given
	token := "some-token"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(500)
	}))

	defer server.Close()
	cli := squarescale.NewClient(server.URL, token)

	// when
	err := cli.AddVolume("another-project", "vol02c1", 1, "gp2", "eu-west-1c")

	if err == nil {
		t.Fatalf("Error is not raised")
	}
}

func UnexpectedHTTPErrorVolumeAdd(t *testing.T) {
	// given
	token := "some-token"
	projectName := "my-project"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if (r.Header.Get("Authorization")) != "bearer some-token" {
			t.Fatalf("Wrong path ! Expected %s, got %s", "bearer some-token", r.Header.Get("Authorization"))
		}

		w.WriteHeader(666)
	}))

	defer server.Close()
	cli := squarescale.NewClient(server.URL, token)

	err := cli.AddVolume(projectName, "vol02c1", 1, "gp2", "eu-west-1c")

	// when
	if err == nil {
		t.Fatalf("Error is not raised")
	}
}

func CantUnmarshal(t *testing.T) {
	// given
	token := "some-token"
	projectName := "my-project"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var path string = r.URL.Path

		if path != "/projects/"+projectName+"/volumes" {
			t.Fatalf("Wrong path ! Expected %s, got %s", "/projects/my-project/volumes", path)
		}

		resBody := `
		{]
		`

		w.Header().Set("Content-Type", "application/json")

		if (r.Header.Get("Authorization")) != "bearer some-token" {
			t.Fatalf("Wrong path ! Expected %s, got %s", "bearer some-token", r.Header.Get("Authorization"))
		}

		w.Write([]byte(resBody))
	}))

	defer server.Close()
	cli := squarescale.NewClient(server.URL, token)

	// when
	_, err := cli.GetVolumes(projectName)

	if err == nil {
		t.Fatalf("Expect error`invalid character ']' looking for beginning of object key string`, got %s", err)
	}
}

func InternalServerErrorCaseForVolumes(t *testing.T) {
	// given
	token := "some-token"
	projectName := "my-project"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		w.WriteHeader(500)
	}))

	defer server.Close()
	cli := squarescale.NewClient(server.URL, token)

	_, err := cli.GetVolumes(projectName)

	// when
	if err == nil {
		t.Fatalf("Error is not raised with 500 Answer")
	}
}

func UnexpectedErrorOnGet(t *testing.T) {
	// given
	token := "some-token"
	projectName := "my-project"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	}))

	defer server.Close()
	cli := squarescale.NewClient(server.URL, token)

	_, err := cli.GetVolumes(projectName)

	// when
	if err == nil {
		t.Fatalf("Error is not raised")
	}
}

func UnknownProjectErrorCaseForVolumeDelete(t *testing.T) {
	// given
	token := "some-token"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resBody := `{"error":"No project found for config name: unknown-project"}`

		w.Header().Set("Content-Type", "application/json")

		if (r.Header.Get("Authorization")) != "bearer some-token" {
			t.Fatalf("Wrong path ! Expected %s, got %s", "bearer some-token", r.Header.Get("Authorization"))
		}

		w.WriteHeader(404)
		w.Write([]byte(resBody))
	}))

	defer server.Close()
	cli := squarescale.NewClient(server.URL, token)

	// when
	err := cli.DeleteVolume("unknown-project", "vol02c1")

	if err == nil {
		t.Fatalf("Error is not raised %s", err)
	}

	if fmt.Sprintf("%s", err) != "{\"error\":\"No project found for config name: unknown-project\"}" {
		t.Fatalf("Error raised `%s` is not `{\"error\":\"No project found for config name: unknown-project\"}`", err)
	}
}

func UnknownVolumeErrorCaseForVolumeDelete(t *testing.T) {
	// given
	token := "some-token"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resBody := `{"error":"Couldn't find Volume with [WHERE \"volumes\".\"cluster_id\" = $1 AND \"volumes\".\"name\" = $2]"}`

		w.Header().Set("Content-Type", "application/json")

		if (r.Header.Get("Authorization")) != "bearer some-token" {
			t.Fatalf("Wrong path ! Expected %s, got %s", "bearer some-token", r.Header.Get("Authorization"))
		}

		w.WriteHeader(404)
		w.Write([]byte(resBody))
	}))

	defer server.Close()
	cli := squarescale.NewClient(server.URL, token)

	// when
	err := cli.DeleteVolume("my-project", "vol02c1")

	if err == nil {
		t.Fatalf("Error is not raised %s", err)
	}

	if fmt.Sprintf("%s", err) != "{\"error\":\"No volume found for name: vol02c1\"}" {
		t.Fatalf("Error raised `%s` is not `{\"error\":\"No volume found for name: vol02c1\"}`", err)
	}
}

func InternalServerErrorCaseForVolumeDelete(t *testing.T) {
	// given
	token := "some-token"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(500)
	}))

	defer server.Close()
	cli := squarescale.NewClient(server.URL, token)

	// when
	err := cli.DeleteVolume("my-project", "vol02c1")

	if err == nil {
		t.Fatalf("Error is not raised")
	}
}

func UnexpectedHTTPErrorVolumeDelete(t *testing.T) {
	// given
	token := "some-token"
	projectName := "my-project"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if (r.Header.Get("Authorization")) != "bearer some-token" {
			t.Fatalf("Wrong path ! Expected %s, got %s", "bearer some-token", r.Header.Get("Authorization"))
		}

		w.WriteHeader(666)
	}))

	defer server.Close()
	cli := squarescale.NewClient(server.URL, token)

	err := cli.DeleteVolume(projectName, "vol02c1")

	// when
	if err == nil {
		t.Fatalf("Error is not raised")
	}
}

func InternalServerErrorCaseForWaitVolume(t *testing.T) {
	// given
	token := "some-token"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(500)
	}))

	defer server.Close()
	cli := squarescale.NewClient(server.URL, token)

	// when
	_, err := cli.WaitVolume("my-project", "vol02c1")

	if err == nil {
		t.Fatalf("Error is not raised")
	}
}
