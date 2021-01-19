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

	// createBatches
	t.Run("Nominal case on createBatches", nominalCaseOnCreateBatches)

	t.Run("Test project not found on createBatches", ProjectNotFoundOnCreateBatches)
	t.Run("Test to create a duplicate on CreateBatches", DuplicateBatchOnCreateBatches)

	// executeBatches
	t.Run("Nominal case on executeBatches", nominalCaseOnExecuteBatches)

	//DeleteBatch
	t.Run("Nominal case on DeleteBatch", nominalCaseOnDeleteBatch)

	t.Run("Test project not found on DeleteBatch", UnknownProjectOnDeleteBatch)
	t.Run("Test to delete a missing batch on DeleteBatch", UnknownBatchOnDeleteBatch)
	t.Run("Test to delete a volume when deploy is in progress on DeleteBatch", DeployInProgressOnDeleteBatch)

	//Error Cases
	t.Run("Test HTTP client error on batch methods (get, create, delete)", ClientHTTPErrorOnBatchMethods)
	t.Run("Test internal server error on batch methods (get, create, delete)", InternalServerErrorOnBatchMethods)
	t.Run("Test badly JSON on batch methods (get, create, delete)", CantUnmarshalOnBatchMethods)
}

func nominalCaseOnGetBatches(t *testing.T) {

	// given
	token := "some-token"
	projectUuid := "d5dde1ad-64ce-4e7c-ade9-43cfd16596d5"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		checkPath(t, "/projects/"+projectUuid+"/batches", r.URL.Path)
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
						"default_environment": "basic-environment",
						"SQSC_SERVICE_NAME": "my-little-batch"
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
					"volumes": [
						{
							"id": 62,
							"volume_name": "vol04b",
							"batch_name": "my-little-image",
							"service_name": null,
							"mount_point": "/var/data",
							"read_only": false
						}
					]

				}
			]
		`

		w.Header().Set("Content-Type", "application/json")

		w.Write([]byte(resBody))
	}))

	defer server.Close()

	cli := squarescale.NewClient(server.URL, token)

	// when
	batches, err := cli.GetBatches(projectUuid)

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

	expectedString = "vol04b"
	if batches[0].Volumes[0].Name != expectedString {
		t.Errorf("Expected batches[0].Volumes[0].Name to be '%s', but got '%s' instead", expectedString, batches[0].Volumes[0].Name)
	}

	expectedString = "/var/data"
	if batches[0].Volumes[0].MountPoint != expectedString {
		t.Errorf("Expected batches[0].Volumes[0].MountPoint to be '%s', but got '%s' instead", expectedString, batches[0].Volumes[0].MountPoint)
	}

	expectedBool = false
	if batches[0].Volumes[0].ReadOnly != expectedBool {
		t.Errorf("Expected batches[0].Volumes[0].ReadOnly to be '%t', but got '%t' instead", expectedBool, batches[0].Volumes[0].ReadOnly)
	}

}

func ProjectNotFoundOnGetBatches(t *testing.T) {

	// given
	token := "some-token"
	projectUuid := "d5dde1ad-64ce-4e7c-ade9-43cfd16596d5"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		resBody := `{"error":"No project found for config name: unknown-project"}`

		w.Header().Set("Content-Type", "application/json")

		w.WriteHeader(404)

		w.Write([]byte(resBody))

	}))

	defer server.Close()

	cli := squarescale.NewClient(server.URL, token)

	// when
	_, errOnGet := cli.GetBatches(projectUuid)

	// then

	expectedError := `1 error occurred: No project found for config name: unknow-project`
	if errOnGet == nil {
		t.Fatalf("Error is not raised with `%s`", expectedError)
	}

	if fmt.Sprintf("%s", errOnGet) == expectedError {
		t.Fatalf("Expected error message:\n`%s`\nGot:\n`%s`", expectedError, errOnGet)
	}

}

func nominalCaseOnExecuteBatches(t *testing.T) {

	// Given
	token := "some-token"
	batchName := "my-little-batch"
	projectUuid := "d5dde1ad-64ce-4e7c-ade9-43cfd16596d5"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		checkPath(t, "/projects/"+projectUuid+"/batches/"+batchName+"/execute", r.URL.Path)
		checkAuthorization(t, r.Header.Get("Authorization"), token)

		resBody := `null`

		w.Header().Set("Content-Type", "application/json")

		w.Write([]byte(resBody))
	}))

	defer server.Close()
	cli := squarescale.NewClient(server.URL, token)

	// when
	err := cli.ExecuteBatch(projectUuid, batchName)

	// then
	if err != nil {
		t.Fatalf("Expect no error, got `%s`", err)
	}

}

func nominalCaseOnCreateBatches(t *testing.T) {

	//Given
	token := "some-token"
	batchName := "my-little-batch"
	projectUuid := "d5dde1ad-64ce-4e7c-ade9-43cfd16596d5"
	periodicBatch := true
	cronExpression := "* * * * *"
	timeZoneName := "Europe/Paris"
	limitNet := 1
	limitMemory := 256
	limitIOPS := 0
	limitCPU := 100
	dockerImageName := "mydockerimage"
	dockerImagePrivate := false
	dockerImageUsername := "me"
	dockerImagePassword := "pwd"
	volumesToBind := make([]squarescale.VolumeToBind, 2)
	volumesToBind[0].Name = "vol02b"
	volumesToBind[0].MountPoint = "/usr/share/nginx/html"
	volumesToBind[0].ReadOnly = false
	volumesToBind[1].Name = "vol03b"
	volumesToBind[1].MountPoint = "/mnt"
	volumesToBind[1].ReadOnly = true

	dockerImageContent := squarescale.DockerImageInfos{
		Name:     dockerImageName,
		Private:  dockerImagePrivate,
		Username: dockerImageUsername,
		Password: dockerImagePassword,
	}
	batchLimitsContent := squarescale.BatchLimits{
		NET:    limitNet,
		Memory: limitMemory,
		IOPS:   limitIOPS,
		CPU:    limitCPU,
	}
	batchCommonContent := squarescale.BatchCommon{
		Name:           batchName,
		Periodic:       periodicBatch,
		CronExpression: cronExpression,
		TimeZoneName:   timeZoneName,
		Limits:         batchLimitsContent,
	}
	batchOrderContent := squarescale.BatchOrder{
		BatchCommon: batchCommonContent,
		DockerImage: dockerImageContent,
		Volumes:     volumesToBind,
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		checkPath(t, "/projects/"+projectUuid+"/batches", r.URL.Path)
		checkAuthorization(t, r.Header.Get("Authorization"), token)

		resBody := `
			{
				"name": "my-little-batch",
				"periodic": true,
				"cron_expression": "* * * * *",
				"time_zone_name": "Europe/Paris",
				"docker_image": {
					"name": "mydockerimage",
					"private": true,
					"username": "someUserName",
					"password": "somePassword"
				},
				"volumes": [
					{
						"id": 51,
						"volume_name": "vol02b",
						"batch_name": "my-little-batch",
						"service_name": null,
						"mount_point": "/usr/share/nginx/html",
						"read_only": false
					},
					{
						"id": 52,
						"volume_name": "vol03b",
						"batch_name": "my-little-batch",
						"service_name": null,
						"mount_point": "/mnt",
						"read_only": true
					}
				],
				"limits": {
					"mem": 256,
					"cpu": 100,
					"iops": 40,
					"net": 256
				}
			}
		`

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(201)

		w.Write([]byte(resBody))
	}))

	defer server.Close()

	cli := squarescale.NewClient(server.URL, token)

	// when
	batch, err := cli.CreateBatch(projectUuid, batchOrderContent)

	var expectedString string
	var expectedInt int
	var expectedBool bool

	if err != nil {
		t.Fatalf("Expect no error, got %s", err)
	}

	expectedString = "my-little-batch"
	if batch.Name != expectedString {
		t.Errorf("Expected batches.Name to be '%s', but got '%s' instead", expectedString, batch.Name)
	}

	expectedBool = true
	if batch.Periodic != expectedBool {
		t.Errorf("Expected batches.Periodic to be '%t', but got '%t' instead", expectedBool, batch.Periodic)
	}

	expectedString = "* * * * *"
	if batch.CronExpression != expectedString {
		t.Errorf("Expected batches.CronExpression to be '%s', but got '%s' instead", expectedString, batch.CronExpression)
	}

	expectedString = "Europe/Paris"
	if batch.TimeZoneName != expectedString {
		t.Errorf("Expected batches.TimeZoneName to be '%s', but got '%s' instead", expectedString, batch.TimeZoneName)
	}

	expectedString = "mydockerimage"
	if batch.DockerImage.Name != expectedString {
		t.Errorf("Expected batches.DockerImage.Name to be '%s', but got '%s' instead", expectedString, batch.DockerImage.Name)
	}

	expectedBool = true
	if batch.DockerImage.Private != expectedBool {
		t.Errorf("Expected batches.DockerImage.Private to be '%t', but got '%t' instead", expectedBool, batch.DockerImage.Private)
	}

	expectedString = "someUserName"
	if batch.DockerImage.Username != expectedString {
		t.Errorf("Expected batches.DockerImage.Username to be '%s', but got '%s' instead", expectedString, batch.DockerImage.Username)
	}

	expectedString = "somePassword"
	if batch.DockerImage.Password != expectedString {
		t.Errorf("Expected batches.DockerImage.Password to be '%s', but got '%s' instead", expectedString, batch.DockerImage.Password)
	}

	expectedInt = 2
	if len(batch.Volumes) != expectedInt {
		t.Errorf("Expected batches.Volumes to contain %d element, but got actually %d", expectedInt, len(batch.Volumes))
	}

	expectedString = "vol02b"
	if batch.Volumes[0].Name != expectedString {
		t.Errorf("Expected batches.Volumes[0].Name to be '%s', but got '%s' instead", expectedString, batch.Volumes[0].Name)
	}

	expectedString = "/usr/share/nginx/html"
	if batch.Volumes[0].MountPoint != expectedString {
		t.Errorf("Expected batches.Volumes[0].MountPoint to be '%s', but got '%s' instead", expectedString, batch.Volumes[0].MountPoint)
	}

	expectedBool = false
	if batch.Volumes[0].ReadOnly != expectedBool {
		t.Errorf("Expected batches.Volumes[0].ReadOnly to be '%t', but got '%t' instead", expectedBool, batch.Volumes[0].ReadOnly)
	}

	expectedString = "vol03b"
	if batch.Volumes[1].Name != expectedString {
		t.Errorf("Expected batches.Volumes[1].Name to be '%s', but got '%s' instead", expectedString, batch.Volumes[1].Name)
	}

	expectedString = "/mnt"
	if batch.Volumes[1].MountPoint != expectedString {
		t.Errorf("Expected batches.Volumes[1].MountPoint to be '%s', but got '%s' instead", expectedString, batch.Volumes[1].MountPoint)
	}

	expectedBool = true
	if batch.Volumes[1].ReadOnly != expectedBool {
		t.Errorf("Expected batches.Volumes[1].ReadOnly to be '%t', but got '%t' instead", expectedBool, batch.Volumes[1].ReadOnly)
	}

	expectedInt = 256
	if batch.Limits.Memory != expectedInt {
		t.Errorf("Expected batches.Limits.Memory to be %d, but got %d instead", expectedInt, batch.Limits.Memory)
	}

	expectedInt = 100
	if batch.Limits.CPU != expectedInt {
		t.Errorf("Expected batches.Limits.CPU to be %d, but got %d instead", expectedInt, batch.Limits.CPU)
	}

	expectedInt = 40
	if batch.Limits.IOPS != expectedInt {
		t.Errorf("Expected batches.Limits.IOPS to be %d, but got %d instead", expectedInt, batch.Limits.IOPS)
	}

	expectedInt = 256
	if batch.Limits.NET != expectedInt {
		t.Errorf("Expected batches.Limits.NET to be %d, but got %d instead", expectedInt, batch.Limits.NET)
	}

}

func ProjectNotFoundOnCreateBatches(t *testing.T) {

	// given
	token := "some-token"
	projectUuid := "d5dde1ad-64ce-4e7c-ade9-43cfd16596d5"
	batchName := "my-little-batch"
	dockerImageName := "mydockerimage"
	volumesToBind := make([]squarescale.VolumeToBind, 1)
	volumesToBind[0].Name = "vol02b"
	volumesToBind[0].MountPoint = "/usr/share/nginx/html"
	volumesToBind[0].ReadOnly = false

	dockerImageContent := squarescale.DockerImageInfos{
		Name: dockerImageName,
	}
	batchCommonContent := squarescale.BatchCommon{
		Name: batchName,
	}
	batchOrderContent := squarescale.BatchOrder{
		BatchCommon: batchCommonContent,
		DockerImage: dockerImageContent,
		Volumes:     volumesToBind,
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		resBody := `{"error":"No project found for config name: unknown-project"}`

		w.Header().Set("Content-Type", "application/json")

		w.WriteHeader(404)

		w.Write([]byte(resBody))

	}))

	defer server.Close()

	cli := squarescale.NewClient(server.URL, token)

	// when
	_, errOnGet := cli.CreateBatch(projectUuid, batchOrderContent)

	// then

	expectedError := `1 error occurred: No project found for config name: unknow-project`
	if errOnGet == nil {
		t.Fatalf("Error is not raised with `%s`", expectedError)
	}

	if fmt.Sprintf("%s", errOnGet) == expectedError {
		t.Fatalf("Expected error message:\n`%s`\nGot:\n`%s`", expectedError, errOnGet)
	}

}

func DuplicateBatchOnCreateBatches(t *testing.T) {
	// given
	token := "some-token"
	projectUuid := "d5dde1ad-64ce-4e7c-ade9-43cfd16596d5"
	batchName := "my-little-batch"
	dockerImageName := "mydockerimage"
	volumesToBind := make([]squarescale.VolumeToBind, 1)
	volumesToBind[0].Name = "vol02b"
	volumesToBind[0].MountPoint = "/usr/share/nginx/html"
	volumesToBind[0].ReadOnly = false

	dockerImageContent := squarescale.DockerImageInfos{
		Name: dockerImageName,
	}
	batchCommonContent := squarescale.BatchCommon{
		Name: batchName,
	}
	batchOrderContent := squarescale.BatchOrder{
		BatchCommon: batchCommonContent,
		DockerImage: dockerImageContent,
		Volumes:     volumesToBind,
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		resBody := `{"error":"PG::UniqueViolation: ERROR:  duplicate key value violates unique constraint \"batch_name_uq\"\nDETAIL:  Key (name, project_id)=(batch2, 4) already exists.\n: INSERT INTO \"batches\" (\"name\", \"project_id\") VALUES ($1, $2) RETURNING \"id\""}`

		w.Header().Set("Content-Type", "application/json")

		w.WriteHeader(409)

		w.Write([]byte(resBody))
	}))

	defer server.Close()
	cli := squarescale.NewClient(server.URL, token)

	// when
	_, err := cli.CreateBatch(projectUuid, batchOrderContent)

	// then
	expectedError := "Batch already exist on project 'd5dde1ad-64ce-4e7c-ade9-43cfd16596d5'"
	if err == nil {
		t.Fatalf("Error is not raised with `%s`", expectedError)
	}

	if fmt.Sprintf("%s", err) != expectedError {
		t.Fatalf("Expected error message:\n`%s`\nGot:\n`%s`", expectedError, err)
	}
}

func nominalCaseOnDeleteBatch(t *testing.T) {
	// given
	token := "some-token"
	projectUuid := "d5dde1ad-64ce-4e7c-ade9-43cfd16596d5"
	batchName := "my-batch"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		checkPath(t, "/projects/"+projectUuid+"/batches/"+batchName, r.URL.Path)
		checkAuthorization(t, r.Header.Get("Authorization"), token)

		resBody := `null`

		w.Header().Set("Content-Type", "application/json")

		w.Write([]byte(resBody))
	}))

	defer server.Close()
	cli := squarescale.NewClient(server.URL, token)

	// when
	err := cli.DeleteBatch(projectUuid, batchName)

	// then
	if err != nil {
		t.Fatalf("Expect no error, got `%s`", err)
	}
}

func UnknownProjectOnDeleteBatch(t *testing.T) {
	// given
	token := "some-token"
	projectUuid := "d5dde1ad-64ce-4e7c-ade9-43cfd16596d5"
	batchName := "my-batch"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resBody := `{"error":"No project found for config name: not-a-project"}`

		w.Header().Set("Content-Type", "application/json")

		w.WriteHeader(404)
		w.Write([]byte(resBody))
	}))

	defer server.Close()
	cli := squarescale.NewClient(server.URL, token)

	// when
	err := cli.DeleteBatch(projectUuid, batchName)

	// then
	expectedError := "Project 'd5dde1ad-64ce-4e7c-ade9-43cfd16596d5' does not exist"
	if err == nil {
		t.Fatalf("Error is not raised with `%s`", expectedError)
	}

	if fmt.Sprintf("%s", err) != expectedError {
		t.Fatalf("Expected error message:\n`%s`\nGot:\n`%s`", expectedError, err)
	}
}

func UnknownBatchOnDeleteBatch(t *testing.T) {
	// given
	token := "some-token"
	projectUuid := "d5dde1ad-64ce-4e7c-ade9-43cfd16596d5"
	batchName := "missing-batch"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resBody := `{"error":"Couldn't find Batch with [WHERE \"batches\".\"cluster_id\" = $1 AND \"batches\".\"name\" = $2]"}`

		w.Header().Set("Content-Type", "application/json")

		w.WriteHeader(404)
		w.Write([]byte(resBody))
	}))

	defer server.Close()
	cli := squarescale.NewClient(server.URL, token)

	// when
	err := cli.DeleteBatch(projectUuid, batchName)

	// then
	expectedError := "Batch 'missing-batch' does not exist"
	if err == nil {
		t.Fatalf("Error is not raised with `%s`", expectedError)
	}

	if fmt.Sprintf("%s", err) != expectedError {
		t.Fatalf("Expected error message:\n`%s`\nGot:\n`%s`", expectedError, err)
	}
}

func DeployInProgressOnDeleteBatch(t *testing.T) {
	// given
	token := "some-token"
	projectUuid := "d5dde1ad-64ce-4e7c-ade9-43cfd16596d5"
	batchName := "my-batch"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		w.WriteHeader(400)
	}))

	defer server.Close()
	cli := squarescale.NewClient(server.URL, token)

	// when
	err := cli.DeleteBatch(projectUuid, batchName)

	// then
	expectedError := "Deploy probably in progress"
	if err == nil {
		t.Fatalf("Error is not raised with `%s`", expectedError)
	}

	if fmt.Sprintf("%s", err) != expectedError {
		t.Fatalf("Expected error message:\n`%s`\nGot:\n`%s`", expectedError, err)
	}
}

func ClientHTTPErrorOnBatchMethods(t *testing.T) {
	// given
	token := "some-token"
	projectUuid := "d5dde1ad-64ce-4e7c-ade9-43cfd16596d5"
	batchName := "my-little-batch"
	dockerImageName := "mydockerimage"
	volumesToBind := make([]squarescale.VolumeToBind, 1)
	volumesToBind[0].Name = "vol02b"
	volumesToBind[0].MountPoint = "/usr/share/nginx/html"
	volumesToBind[0].ReadOnly = false

	dockerImageContent := squarescale.DockerImageInfos{
		Name: dockerImageName,
	}
	batchCommonContent := squarescale.BatchCommon{
		Name: batchName,
	}
	batchOrderContent := squarescale.BatchOrder{
		BatchCommon: batchCommonContent,
		DockerImage: dockerImageContent,
		Volumes:     volumesToBind,
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	}))

	defer server.Close()

	cli := squarescale.NewClient(server.URL, token)

	// when
	_, errOnGet := cli.GetBatches(projectUuid)
	_, errOnCreate := cli.CreateBatch(projectUuid, batchOrderContent)
	errOnDelete := cli.DeleteBatch(projectUuid, batchName)

	// then
	if errOnGet == nil {
		t.Errorf("Error is not raised on GetBatches")
	}

	if errOnCreate == nil {
		t.Errorf("Error is not raised on CreateBatch")
	}

	if errOnDelete == nil {
		t.Errorf("Error is not raised on DeleteBatch")
	}

}

func InternalServerErrorOnBatchMethods(t *testing.T) {
	// given
	token := "some-token"
	projectUuid := "d5dde1ad-64ce-4e7c-ade9-43cfd16596d5"
	batchName := "my-little-batch"
	dockerImageName := "mydockerimage"
	volumesToBind := make([]squarescale.VolumeToBind, 1)
	volumesToBind[0].Name = "vol02b"
	volumesToBind[0].MountPoint = "/usr/share/nginx/html"
	volumesToBind[0].ReadOnly = false

	dockerImageContent := squarescale.DockerImageInfos{
		Name: dockerImageName,
	}
	batchCommonContent := squarescale.BatchCommon{
		Name: batchName,
	}
	batchOrderContent := squarescale.BatchOrder{
		BatchCommon: batchCommonContent,
		DockerImage: dockerImageContent,
		Volumes:     volumesToBind,
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "application/json")

		w.WriteHeader(500)

	}))

	defer server.Close()

	cli := squarescale.NewClient(server.URL, token)

	// when
	_, errOnGet := cli.GetBatches(projectUuid)
	_, errOnCreate := cli.CreateBatch(projectUuid, batchOrderContent)
	errOnDelete := cli.DeleteBatch(projectUuid, batchName)

	// then
	expectedError := "An unexpected error occurred (code: 500)"
	if errOnGet == nil {
		t.Errorf("Error is not raised with `%s`", expectedError)
	}

	if fmt.Sprintf("%s", errOnGet) != expectedError {
		t.Errorf("Expected error message:\n`%s`\nGot:\n`%s`", expectedError, errOnGet)
	}

	if errOnCreate == nil {
		t.Errorf("Error is not raised with `%s`", expectedError)
	}

	if fmt.Sprintf("%s", errOnCreate) != expectedError {
		t.Errorf("Expected error message:\n`%s`\nGot:\n`%s`", expectedError, errOnGet)
	}

	if errOnDelete == nil {
		t.Errorf("Error is not raised with `%s`", expectedError)
	}

	if fmt.Sprintf("%s", errOnDelete) != expectedError {
		t.Errorf("Expected error message:\n`%s`\nGot:\n`%s`", expectedError, errOnDelete)
	}
}

func CantUnmarshalOnBatchMethods(t *testing.T) {
	// given
	token := "some-token"
	projectUuid := "d5dde1ad-64ce-4e7c-ade9-43cfd16596d5"
	batchName := "my-little-batch"
	dockerImageName := "mydockerimage"
	volumesToBind := make([]squarescale.VolumeToBind, 1)
	volumesToBind[0].Name = "vol02b"
	volumesToBind[0].MountPoint = "/usr/share/nginx/html"
	volumesToBind[0].ReadOnly = false

	dockerImageContent := squarescale.DockerImageInfos{
		Name: dockerImageName,
	}
	batchCommonContent := squarescale.BatchCommon{
		Name: batchName,
	}
	batchOrderContent := squarescale.BatchOrder{
		BatchCommon: batchCommonContent,
		DockerImage: dockerImageContent,
		Volumes:     volumesToBind,
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		resBody := `{]`

		w.Header().Set("Content-Type", "application/json")

		if r.Method == "POST" { // Manage Add
			w.WriteHeader(201)
		} else {
			w.WriteHeader(200)
		}

		w.Write([]byte(resBody))
	}))

	defer server.Close()
	cli := squarescale.NewClient(server.URL, token)

	// when
	_, errOnGet := cli.GetBatches(projectUuid)
	_, errOnCreate := cli.CreateBatch(projectUuid, batchOrderContent)

	// then
	expectedError := "invalid character ']' looking for beginning of object key string"
	if errOnGet == nil {
		t.Fatalf("Error is not raised with `%s`", expectedError)
	}

	if fmt.Sprintf("%s", errOnGet) != expectedError {
		t.Fatalf("Expected error message:\n`%s`\nGot:\n`%s`", expectedError, errOnGet)
	}

	if errOnCreate == nil {
		t.Fatalf("Error is not raised with `%s`", expectedError)
	}

	if fmt.Sprintf("%s", errOnCreate) != expectedError {
		t.Fatalf("Expected error message:\n`%s`\nGot:\n`%s`", expectedError, errOnCreate)
	}
}
