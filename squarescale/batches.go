package squarescale

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Batch describes a project container as returned by the Squarescale API
type Batch struct {
	Name               string           `json:"name"`
	Periodic           bool             `json:"periodic"`
	CronExpression     string           `json:"cron_expression"`
	TimeZoneName       string           `json:"time_zone_name"`
	Limits             BatchLimits      `json:"limits"`
	CustomEnvironment  string           `json:"custom_environment"`
	DefaultEnvironment string           `json:"default_environment"`
	RunCommand         []string         `json:"run_command"`
	Status             BatchStatus      `json:"status"`
	DockerImage        BatchDockerImage `json:"docker_image"`
	RefreshUrl         []string         `json:"refresh_url"`
	Volumes            string           `json:"mounted_volumes"`
}

type BatchLimits struct {
	Memory int `json:"mem"`
	CPU    int `json:"cpu"`
	NET    int `json:"net"`
	IOPS   int `json:"iops"`
}

type BatchStatus struct {
	Infra    string        `json:"display_infra_status"`
	Schedule BatchSchedule `json:"schedule"`
}

type BatchSchedule struct {
	Running   int    `json:"running"`
	Instances int    `json:"instances"`
	Level     string `json:"level"`
	Message   string `json:"message"`
}

type BatchDockerImage struct {
	Name string `json:"name"`
}

// GetBatch gets all the batches attached to a Project
func (c *Client) GetBatches(project string) ([]Batch, error) {
	code, body, err := c.get("/projects/" + project + "/batches")
	if err != nil {
		return []Batch{}, err
	}

	switch code {
	case http.StatusOK:
	case http.StatusNotFound:
		return []Batch{}, fmt.Errorf("Project '%s' does not exist", project)
	default:
		return []Batch{}, unexpectedHTTPError(code, body)
	}

	var batchesByID []Batch

	if err := json.Unmarshal(body, &batchesByID); err != nil {
		return []Batch{}, err
	}

	return batchesByID, nil

}
