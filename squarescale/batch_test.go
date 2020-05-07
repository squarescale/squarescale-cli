package squarescale_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/squarescale/squarescale-cli/squarescale"
)

func TestBatches(t *testing.T) {

	// getBatches
	t.Run("Nominal case on getBatches", nominalCaseOnGetBatches)

	t.Run("Test project not found on getBatches", ProjectNotFoundOnGetBatches)

	//Error Cases
	t.Run("Test HTTP client error on batch methods (get)", ClientHTTPErrorOnBatchMethods)
	t.Run("Test internal server error on batch methods (get)", InternalServerErrorOnBatchMethods)
	t.Run("Test badly JSON on batch methods (get)", CantUnmarshalOnBatchMethods)
}

func nominalCaseOnGetBatches(t *testing.T) {

	// given
	token := "some-token"
	projectName := "my-project"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		checkPath(t, "/projects/"+projectName+"/batches", r.URL.Path)
		checkAuthorization(t, r.Header.Get("Authorization"), token)

		resBody := `
			[
				{
					"name": "my-little-batch",
					"periodic": true,
					"cron_expression": "* * * * *",
					"time_zone_name": "Europe/Paris",
					"limits": {
							"mem": 256,
							"cpu": 100,
							"iops": 0,
							"net": 1
					},
					"custom_environment": {
						"custom_environment": "final-environment"
					},
					"default_environment": {
						"default_environment": "basic-environment"
					},
					"run_command": "command1",
					"status": {
						"display_infra_status": "ok",
						"schedule": {
							"running": 0,
							"instances": 0,
							"level": "ok",
							"message": "Scheduled"
						}
					},
					"docker_image": {
						"name": "my-little-image"
					},
					"refresh_url": "new-url",
					"mounted_volumes": "volume1"
				}
			]
		`

		w.Header().Set("Content-Type", "application/json")

		w.Write([]byte(resBody))
	}))

	defer server.Close()

	cli := squarescale.NewClient(server.URL, token)

	// when
	batches, err := cli.GetBatches(projectName)

	// then
	var expectedInt int
	var expectedString string
	var expectedBool bool

	if err != nil {
		t.Fatalf("Expect no error, got %s", err)
	}

	expectedInt = 1
	if len(batches) != expectedInt {
		t.Fatalf("Expected batches to contain %d element, but got actually %d", expectedInt, len(batches))
	}

	expectedString = "my-little-batch"
	if batches[0].Name != expectedString {
		t.Errorf("Expected batches.Name to be '%s', but got '%s' instead", expectedString, batches[0].Name)
	}

	expectedBool = true
	if batches[0].Periodic != expectedBool {
		t.Errorf("Expected batches.Periodic to be '%t', but got '%t' instead", expectedBool, batches[0].Periodic)
	}

	expectedString = "* * * * *"
	if batches[0].CronExpression != expectedString {
		t.Errorf("Expected batches.CronExpression to be '%s', but got '%s' instead", expectedString, batches[0].CronExpression)
	}

	expectedString = "Europe/Paris"
	if batches[0].TimeZoneName != expectedString {
		t.Errorf("Expected batches.TimeZoneName to be '%s', but got '%s' instead", expectedString, batches[0].TimeZoneName)
	}

	expectedInt = 256
	if batches[0].Limits.Memory != expectedInt {
		t.Errorf("Expected batches.Limits.Memory to be %d, but got %d instead", expectedInt, batches[0].Limits.Memory)
	}

	expectedInt = 100
	if batches[0].Limits.CPU != expectedInt {
		t.Errorf("Expected batches.Limits.CPU to be %d, but got %d instead", expectedInt, batches[0].Limits.CPU)
	}

	expectedInt = 1
	if batches[0].Limits.NET != expectedInt {
		t.Errorf("Expected batches.Limits.NET to be %d, but got %d instead", expectedInt, batches[0].Limits.NET)
	}

	expectedInt = 0
	if batches[0].Limits.IOPS != expectedInt {
		t.Errorf("Expected batches.Limits.IOPS to be %d, but got %d instead", expectedInt, batches[0].Limits.IOPS)
	}

	expectedString = "final-environment"
	if batches[0].CustomEnvironment.CustomEnvironmentValue != expectedString {
		t.Errorf("Expected batches.CustomEnvironment to be '%s', but got '%s' instead", expectedString, batches[0].CustomEnvironment.CustomEnvironmentValue)
	}

	expectedString = "basic-environment"
	if batches[0].DefaultEnvironment.DefaultEnvironmentValue != expectedString {
		t.Errorf("Expected batches.DefaultEnvironment to be '%s', but got '%s' instead", expectedString, batches[0].DefaultEnvironment.DefaultEnvironmentValue)
	}

	expectedString = "command1"
	if batches[0].RunCommand != expectedString {
		t.Errorf("Expected batches.RunCommand to be '%s', but got '%s' instead", expectedString, batches[0].RunCommand)
	}

	expectedString = "ok"
	if batches[0].Status.Infra != expectedString {
		t.Errorf("Expected batches.Status.Infra to be '%s', but got '%s' instead", expectedString, batches[0].Status.Infra)
	}

	expectedInt = 0
	if batches[0].Status.Schedule.Running != expectedInt {
		t.Errorf("Expected batches.Status.Schedule.Running to be %d, but got %d instead", expectedInt, batches[0].Status.Schedule.Running)
	}

	expectedInt = 0
	if batches[0].Status.Schedule.Instances != expectedInt {
		t.Errorf("Expected batches.Status.Schedule.Instances to be %d, but got %d instead", expectedInt, batches[0].Status.Schedule.Instances)
	}

	expectedString = "ok"
	if batches[0].Status.Schedule.Level != expectedString {
		t.Errorf("Expected batches.Status.Schedule.Level to be '%s', but got '%s' instead", expectedString, batches[0].Status.Schedule.Level)
	}

	expectedString = "Scheduled"
	if batches[0].Status.Schedule.Message != expectedString {
		t.Errorf("Expected batches.Status.Schedule.Message to be '%s', but got '%s' instead", expectedString, batches[0].Status.Schedule.Message)
	}

	expectedString = "my-little-image"
	if batches[0].DockerImage.Name != expectedString {
		t.Errorf("Expected batches.DockerImage.Name to be '%s', but got '%s' instead", expectedString, batches[0].DockerImage.Name)
	}

	expectedString = "new-url"
	if batches[0].RefreshUrl != expectedString {
		t.Errorf("Expected batches.RefreshUrl to be '%s', but got '%s' instead", expectedString, batches[0].RefreshUrl)
	}

	expectedString = "volume1"
	if batches[0].Volumes != expectedString {
		t.Errorf("Expected batches.Volumes to be '%s', but got '%s' instead", expectedString, batches[0].Volumes)
	}

}

func ProjectNotFoundOnGetBatches(t *testing.T) {

	// given
	token := "some-token"
	projectName := "unknow-project"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		resBody := `{"error":"No project found for config name: unknown-project"}`

		w.Header().Set("Content-Type", "application/json")

		w.WriteHeader(404)

		w.Write([]byte(resBody))

	}))

	defer server.Close()

	cli := squarescale.NewClient(server.URL, token)

	// when
	_, errOnGet := cli.GetBatches(projectName)

	// then

	expectedError := `1 error occurred: No project found for config name: unknow-project`
	if errOnGet == nil {
		t.Fatalf("Error is not raised with `%s`", expectedError)
	}

	if fmt.Sprintf("%s", errOnGet) == expectedError {
		t.Fatalf("Expected error message:\n`%s`\nGot:\n`%s`", expectedError, errOnGet)
	}

}

func ClientHTTPErrorOnBatchMethods(t *testing.T) {
	// given
	token := "some-token"
	projectName := "my-project"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	}))

	defer server.Close()

	cli := squarescale.NewClient(server.URL, token)

	// when
	_, errOnGet := cli.GetBatches(projectName)

	// then
	if errOnGet == nil {
		t.Errorf("Error is not raised on GetBatches")
	}

}

func InternalServerErrorOnBatchMethods(t *testing.T) {
	// given
	token := "some-token"
	projectName := "bad-project"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "application/json")

		w.WriteHeader(500)

	}))

	defer server.Close()

	cli := squarescale.NewClient(server.URL, token)

	// when
	_, errOnGet := cli.GetBatches(projectName)

	// then
	expectedError := "An unexpected error occurred (code: 500)"
	if errOnGet == nil {
		t.Errorf("Error is not raised with `%s`", expectedError)
	}

	if fmt.Sprintf("%s", errOnGet) != expectedError {
		t.Errorf("Expected error message:\n`%s`\nGot:\n`%s`", expectedError, errOnGet)
	}
}

func CantUnmarshalOnBatchMethods(t *testing.T) {
	// given
	token := "some-token"
	projectName := "my-project"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		resBody := `{]`

		w.Header().Set("Content-Type", "application/json")

		w.Write([]byte(resBody))
	}))

	defer server.Close()
	cli := squarescale.NewClient(server.URL, token)

	// when
	_, errOnGet := cli.GetBatches(projectName)

	// then
	expectedError := "invalid character ']' looking for beginning of object key string"
	if errOnGet == nil {
		t.Fatalf("Error is not raised with `%s`", expectedError)
	}

	if fmt.Sprintf("%s", errOnGet) != expectedError {
		t.Fatalf("Expected error message:\n`%s`\nGot:\n`%s`", expectedError, errOnGet)
	}
}
