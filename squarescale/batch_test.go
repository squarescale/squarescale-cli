package squarescale_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/squarescale/squarescale-cli/squarescale"
)

func TestGetBatches(t *testing.T) {

	// nominal case
	t.Run("nominal get batches", nominalCase)

	// other cases

}

func nominalCase(t *testing.T) {

	/*
		1 OK - completer la payload et la vérifier dans le then (regarder exemples sur infrastructure-scheduler)
		2 OK - vérifier dans le "mock" http que l'url soit correcte
		3 vérifier que l'on ai bien le token passé au bon endroit
		4 tester Not Found (à faire dans un autre cas de tests)
		4 tester Internal server error (à faire dans un autre cas de tests)


	*/

	// given
	token := "some-token"
	projectName := "titi"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var path string = r.URL.Path

		if path != "/projects/"+projectName+"/batches" {
			t.Fatalf("Wrong path ! Expected %s, got %s", "/projects/titi/batches", path)
		}

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
					"custom_environment": "final-environment",
					"default_environment": "basic-environment",
					"run_command": ["command1"],
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
					"refresh_url": ["new-url"],
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

	fmt.Println(batches)

	// then
	if err != nil {
		t.Fatalf("Expect no error, got %s", err)
	}

	if len(batches) != 1 {
		t.Fatalf("Expect %d, got %d", 1, len(batches))
	}

	if batches[0].Name != "my-little-batch" {
		t.Errorf("Expect batchName %s , got %s", "my-little-batch", batches[0].Name)
	}

	if batches[0].Periodic != true {
		t.Errorf("Expect batchPeriodic %s , got %t", "true", batches[0].Periodic)
	}

	if batches[0].CronExpression != "* * * * *" {
		t.Errorf("Expect batchCronExpression %s , got %s", "* * * * *", batches[0].CronExpression)
	}

	if batches[0].TimeZoneName != "Europe/Paris" {
		t.Errorf("Expect batchTimeZoneName %s , got %s", "Europe/Paris", batches[0].TimeZoneName)
	}

	if batches[0].Limits.Memory != 256 {
		t.Errorf("Expect batchesLimitsMemory %d , got %d", 256, batches[0].Limits.Memory)
	}

	if batches[0].Limits.CPU != 100 {
		t.Errorf("Expect batchesLimitsCPU %d , got %d", 100, batches[0].Limits.CPU)
	}

	if batches[0].Limits.NET != 1 {
		t.Errorf("Expect batchesLimitsNet %d , got %d", 1, batches[0].Limits.NET)
	}

	if batches[0].Limits.IOPS != 0 {
		t.Errorf("Expect batchesLimitsIops %d , got %d", 0, batches[0].Limits.IOPS)
	}

	if batches[0].CustomEnvironment != "final-environment" {
		t.Errorf("Expect batchesCustomEnvironment %s , got %s", "final-environment", batches[0].CustomEnvironment)
	}

	if batches[0].DefaultEnvironment != "basic-environment" {
		t.Errorf("Expect batchesDefaultEnvironment %s , got %s", "basic-environment", batches[0].DefaultEnvironment)
	}

	if len(batches[0].RunCommand) != 1 {
		t.Errorf("Expect batchesRunCommand length of %d , got %d", 1, len(batches[0].RunCommand))
	}

	if batches[0].RunCommand[0] != "command1" {
		t.Errorf("Expect batchesRunCommand %s , got %d", "command1", len(batches[0].RunCommand))
	}

	if batches[0].Status.Infra != "ok" {
		t.Errorf("Expect batchesStatusInfrastructures %s , got %s", "ok", batches[0].Status.Infra)
	}

	if batches[0].Status.Schedule.Running != 0 {
		t.Errorf("Expect batchesStatusRunning %d , got %d", 0, batches[0].Status.Schedule.Running)
	}

	if batches[0].Status.Schedule.Instances != 0 {
		t.Errorf("Expect batchesStatusInstances %d , got %d", 0, batches[0].Status.Schedule.Instances)
	}

	if batches[0].Status.Schedule.Level != "ok" {
		t.Errorf("Expect batchesStatusLevel %s , got %s", "ok", batches[0].Status.Schedule.Level)
	}

	if batches[0].Status.Schedule.Message != "Scheduled" {
		t.Errorf("Expect batchesStatusMessage %s , got %s", "Scheduled", batches[0].Status.Schedule.Message)
	}

	if batches[0].DockerImage.Name != "my-little-image" {
		t.Errorf("Expect batchesDockerImageName %s , got %s", "my-little-image", batches[0].DockerImage.Name)
	}

	if len(batches[0].RefreshUrl) != 1 {
		t.Errorf("Expect batchesRefreshUrl length of %d , got %d", 1, len(batches[0].RefreshUrl))
	}

	if batches[0].RefreshUrl[0] != "new-url" {
		t.Errorf("Expect batchesRefreshUrl %s , got %s", "new-url", batches[0].RefreshUrl[0])
	}

	if batches[0].Volumes != "volume1" {
		t.Errorf("Expect batchesRefreshUrl %s , got %s", "volume1", batches[0].Volumes)
	}

}
