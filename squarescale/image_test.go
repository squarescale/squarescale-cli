package squarescale_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/squarescale/squarescale-cli/squarescale"
)

func TestAddImage(t *testing.T) {
	// AddImage
	t.Run("Nominal case on AddImage", nominalCaseOnAddImage)

	// Error cases
	t.Run("Test HTTP client error on image methods (add)", ClientHTTPErrorOnAddImage)
	t.Run("Test internal server error on image methods (add)", InternalServerErrorOnAddImage)
}

func nominalCaseOnAddImage(t *testing.T) {
	// given
	token := "some-token"
	projectName := "my-project"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		checkPath(t, "/projects/"+projectName+"/docker_images", r.URL.Path)
		checkAuthorization(t, r.Header.Get("Authorization"), token)

		resBody := `
		{
			"name": "nginx2",
			"docker_image": {
				"name": "nginx"
			},
			"container": {
				"limits": {
					"mem": 256,
					"cpu": 100,
					"iops": 0,
					"net": 1
				}
			},
			"size": 1,
			"volumes_to_bind": [
				{
					"volume_name": "vol02b",
					"read_only": false,
					"mount_point": "/usr/share/nginx/html"
				}
			]
		}
		`

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(201)

		w.Write([]byte(resBody))
	}))

	defer server.Close()
	cli := squarescale.NewClient(server.URL, token)

	// when
	volumesToBind := make([]squarescale.VolumeToBind, 1)
	volumesToBind[0].Name = "vol02b"
	volumesToBind[0].MountPoint = "/usr/share/nginx/html"
	volumesToBind[0].ReadOnly = false

	dockerCapabilities := []string{"AUDIT_WRITE", "CHOWN"}
	dockerDevices := []squarescale.DockerDevice{{SRC: "src", DST: "dst"}}

	errPublic := cli.AddImage(projectName, "nginx", "", "", "", "", 1, "nginx", volumesToBind, dockerCapabilities, dockerDevices, true)
	errPrivate := cli.AddImage(projectName, "nginx", "login", "pass", "", "", 1, "nginx", volumesToBind, dockerCapabilities, dockerDevices, true)

	// then
	if errPublic != nil {
		t.Fatalf("Expect no error, got `%s`", errPublic)
	}

	if errPrivate != nil {
		t.Fatalf("Expect no error, got `%s`", errPrivate)
	}
}

func ClientHTTPErrorOnAddImage(t *testing.T) {
	// given
	token := "some-token"
	projectName := "my-project"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	}))

	defer server.Close()
	cli := squarescale.NewClient(server.URL, token)

	// when
	volumesToBind := make([]squarescale.VolumeToBind, 1)
	volumesToBind[0].Name = "vol02b"
	volumesToBind[0].MountPoint = "/usr/share/nginx/html"
	volumesToBind[0].ReadOnly = false

	dockerCapabilities := []string{"AUDIT_WRITE", "CHOWN"}
	dockerDevices := []squarescale.DockerDevice{{SRC: "src", DST: "dst"}}

	errOnAdd := cli.AddImage(projectName, "nginx", "", "", "", "", 1, "nginx", volumesToBind, dockerCapabilities, dockerDevices, true)

	// then
	if errOnAdd == nil {
		t.Errorf("Error is not raised on image AddImage")
	}
}

func InternalServerErrorOnAddImage(t *testing.T) {
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
	volumesToBind := make([]squarescale.VolumeToBind, 1)
	volumesToBind[0].Name = "vol02b"
	volumesToBind[0].MountPoint = "/usr/share/nginx/html"
	volumesToBind[0].ReadOnly = false

	dockerCapabilities := []string{"AUDIT_WRITE", "CHOWN"}
	dockerDevices := []squarescale.DockerDevice{{SRC: "src", DST: "dst"}}

	errOnAdd := cli.AddImage(projectName, "nginx", "", "", "", "", 1, "nginx", volumesToBind, dockerCapabilities, dockerDevices, true)

	// then
	expectedError := "An unexpected error occurred (code: 500)"
	if errOnAdd == nil {
		t.Errorf("Error is not raised with `%s`", expectedError)
	}

	if fmt.Sprintf("%s", errOnAdd) != expectedError {
		t.Errorf("Expected error message:\n`%s`\nGot:\n`%s`", expectedError, errOnAdd)
	}
}
