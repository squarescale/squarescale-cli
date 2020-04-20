package squarescale

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Common

type BatchCommon struct {
	Name           string      `json:"name"`
	Periodic       bool        `json:"periodic"`
	CronExpression string      `json:"cron_expression"`
	TimeZoneName   string      `json:"time_zone_name"`
	Limits         BatchLimits `json:"limits"`
}

type BatchLimits struct {
	Memory int `json:"mem"`
	CPU    int `json:"cpu"`
	NET    int `json:"net"`
	IOPS   int `json:"iops"`
}

// Create Batch part

type BatchOrder struct {
	BatchCommon
	DockerImage DockerImageInfos   `json:"docker_image"`
	VolumesBind []BatchVolumesBind `json:"volumes_to_bind"`
}

type DockerImageInfos struct {
	Name     string `json:"name"`
	Private  bool   `json:"private"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type BatchVolumesBind struct {
	Name string `json:"name"`
}

func (c *Client) CreateBatch(project string, batchOrderContent BatchOrder) ([]BatchOrder, error) {

	payload := &JSONObject{
		"name":            batchOrderContent.BatchCommon.Name,
		"docker_image":    batchOrderContent.DockerImage,
		"periodic":        batchOrderContent.BatchCommon.Periodic,
		"cron_expression": batchOrderContent.BatchCommon.CronExpression,
		"time_zone_name":  batchOrderContent.BatchCommon.TimeZoneName,
		"limits":          batchOrderContent.BatchCommon.Limits,
		"volumes_bind":    batchOrderContent.VolumesBind,
	}

	code, body, err := c.post("/projects/"+project+"/batches", payload)
	if err != nil {
		return []BatchOrder{}, err
	}

	fmt.Printf("\npayload: \n`%s`\n", payload)
	fmt.Printf("\ncode: \n`%d`\n", code)

	switch code {
	case http.StatusCreated:
	case http.StatusNotFound:
		return []BatchOrder{}, fmt.Errorf("Project '%s' does not exist", project)
	case http.StatusConflict:
		return []BatchOrder{}, fmt.Errorf("Batch already exist on project '%s'", project)
	default:
		return []BatchOrder{}, unexpectedHTTPError(code, body)
	}

	var batchesByID []BatchOrder

	if err := json.Unmarshal(body, &batchesByID); err != nil {
		return []BatchOrder{}, err
	}

	return batchesByID, err
}

// Get Batch part

type RunningBatch struct {
	BatchCommon
	CustomEnvironment  BatchCustomEnvironment  `json:"custom_environment"`
	DefaultEnvironment BatchDefaultEnvironment `json:"default_environment"`
	RunCommand         string                  `json:"run_command"`
	Status             BatchStatus             `json:"status"`
	DockerImage        BatchDockerImage        `json:"docker_image"`
	RefreshUrl         string                  `json:"refresh_url"`
	Volumes            string                  `json:"mounted_volumes"`
}

type BatchCustomEnvironment struct {
	CustomEnvironmentValue string `json:"custom_environment"`
}

type BatchDefaultEnvironment struct {
	DefaultEnvironmentValue string `json:"default_environment"`
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

func (c *Client) GetBatches(project string) ([]RunningBatch, error) {
	code, body, err := c.get("/projects/" + project + "/batches")
	if err != nil {
		return []RunningBatch{}, err
	}

	switch code {
	case http.StatusOK:
	case http.StatusNotFound:
		return []RunningBatch{}, fmt.Errorf("Project '%s' does not exist", project)
	default:
		return []RunningBatch{}, unexpectedHTTPError(code, body)
	}

	var batchesByID []RunningBatch

	if err := json.Unmarshal(body, &batchesByID); err != nil {
		return []RunningBatch{}, err
	}

	return batchesByID, nil

}
