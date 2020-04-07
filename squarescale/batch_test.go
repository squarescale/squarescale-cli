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
		1 completer la payload et la vérifier dans le then (regarder exemples sur infrastructure-scheduler)
		2 vérifier dans le "mock" http que l'url soit correcte
		3 vérifier que l'on ai bien le token passé au bon endroit
		4 tester Not Found (à faire dans un autre cas de tests)
		4 tester Internal server error (à faire dans un autre cas de tests)


	*/

	// given
	token := "some-token"
	projectName := "titi"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		/*
			{

				"docker_image" => {"name" => self.container.docker_image.name},
				"refresh_url" => self.container.container_callbacks.map { |cb| cb.full_url_path(:refresh) }[0],
				"volumes" => self.mounted_volumes
			}

			DEFAULT_LIMITS = {

			}.freeze
		*/
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
						"infra": "ok",
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
}
