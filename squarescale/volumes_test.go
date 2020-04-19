package squarescale_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/squarescale/squarescale-cli/squarescale"
)

func TestGetVolumes(t *testing.T) {

	// nominal case
	t.Run("nominal get volumes", nominalCaseForVolumes)

	// other cases
	t.Run("test of Not Found Page", NotFoundCaseForVolume)
	t.Run("test of Internal Server Error ", InternalServerErrorCaseForVolumes)

}

func nominalCaseForVolumes(t *testing.T) {
	// given
	token := "some-token"
	projectName := "titi"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var path string = r.URL.Path

		if path != "/projects/"+projectName+"/volumes" {
			t.Fatalf("Wrong path ! Expected %s, got %s", "/projects/titi/volumes", path)
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

	volume, err := cli.GetVolumeInfo(projectName, volumes[0].Name)

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

func NotFoundCaseForVolume(t *testing.T) {
	// given
	token := "some-token"
	projectName := "titi"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var path string = r.URL.Path

		if path != "/projects/"+projectName+"/volumes" {
			t.Fatalf("Wrong path ! Expected %s, got %s", "/projects/titi/volumes", path)
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
		t.Fatalf("Error is not raised with `Volume 'missing-volume' not found for project 'titi'`")
	}
}

func InternalServerErrorCaseForVolumes(t *testing.T) {
	// given
	token := "some-token"
	projectName := "titi"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var path string = r.URL.Path

		if path != "/projects/"+projectName+"/volumes" {
			t.Fatalf("Wrong path ! Expected %s, got %s", "/projects/titi/volumes", path)
		}

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
