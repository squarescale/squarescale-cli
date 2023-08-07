package squarescale

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/squarescale/logger"
)

// Common
type BatchCommon struct {
	Name               string         `json:"name"`
	Periodic           bool           `json:"periodic"`
	CronExpression     string         `json:"cron_expression"`
	TimeZoneName       string         `json:"time_zone_name"`
	Limits             BatchLimits    `json:"limits"`
	RunCommand         string         `json:"run_command,omitempty"`
	Entrypoint         string         `json:"entrypoint,omitempty"`
	DockerCapabilities []string       `json:"docker_capabilities"`
	DockerDevices      []DockerDevice `json:"docker_devices"`
}

type BatchLimits struct {
	Memory int `json:"mem"`
	CPU    int `json:"cpu"`
	IOPS   int `json:"iops"`
}

type DockerDevice struct {
	SRC string `json:"src_path"`
	DST string `json:"dst_path,omitempty"`
	OPT string `json:"options,omitempty"`
}

// Create Batch part
type BatchOrder struct {
	BatchCommon
	DockerImage DockerImageInfos `json:"docker_image"`
	Volumes     []VolumeToBind   `json:"volumes_to_bind"`
}

type CreatedBatch struct {
	BatchCommon
	DockerImage DockerImageInfos `json:"docker_image"`
	Volumes     []VolumeToBind   `json:"volumes"`
}

type DockerImageInfos struct {
	Name     string `json:"name"`
	Private  bool   `json:"private"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func (c *Client) CreateBatch(uuid string, batchOrderContent BatchOrder) (CreatedBatch, error) {

	payload := &JSONObject{
		"name":                batchOrderContent.BatchCommon.Name,
		"docker_image":        batchOrderContent.DockerImage,
		"periodic":            batchOrderContent.BatchCommon.Periodic,
		"cron_expression":     batchOrderContent.BatchCommon.CronExpression,
		"time_zone_name":      batchOrderContent.BatchCommon.TimeZoneName,
		"limits":              batchOrderContent.BatchCommon.Limits,
		"volumes_to_bind":     batchOrderContent.Volumes,
		"run_command":         batchOrderContent.BatchCommon.RunCommand,
		"entrypoint":          batchOrderContent.BatchCommon.Entrypoint,
		"docker_capabilities": batchOrderContent.BatchCommon.DockerCapabilities,
		"docker_devices":      batchOrderContent.BatchCommon.DockerDevices,
	}

	code, body, err := c.post("/projects/"+uuid+"/batches", payload)
	if err != nil {
		return CreatedBatch{}, err
	}

	switch code {
	case http.StatusCreated:
	case http.StatusNotFound:
		return CreatedBatch{}, fmt.Errorf("Project '%s' does not exist", uuid)
	case http.StatusConflict:
		return CreatedBatch{}, fmt.Errorf("Batch already exist on project '%s'", uuid)
	default:
		return CreatedBatch{}, unexpectedHTTPError(code, body)
	}

	var createdBatch CreatedBatch

	if err := json.Unmarshal(body, &createdBatch); err != nil {
		return CreatedBatch{}, err
	}

	return createdBatch, err
}

// Get Batch part
type RunningBatch struct {
	BatchCommon
	CustomEnvironment  map[string]string `json:"custom_environment"`
	DefaultEnvironment map[string]string `json:"default_environment"`
	RunCommand         string            `json:"run_command"`
	Status             BatchStatus       `json:"status"`
	DockerImage        BatchDockerImage  `json:"docker_image"`
	RefreshUrl         string            `json:"refresh_url"`
	Volumes            []VolumeToBind    `json:"volumes"`
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

func (c *Client) GetBatches(uuid string) ([]RunningBatch, error) {
	code, body, err := c.get("/projects/" + uuid + "/batches")
	if err != nil {
		return []RunningBatch{}, err
	}

	switch code {
	case http.StatusOK:
	case http.StatusNotFound:
		return []RunningBatch{}, fmt.Errorf("Project '%s' does not exist", uuid)
	default:
		return []RunningBatch{}, unexpectedHTTPError(code, body)
	}

	var batchesByID []RunningBatch

	if err := json.Unmarshal(body, &batchesByID); err != nil {
		return []RunningBatch{}, err
	}

	return batchesByID, nil

}

// DeleteBatch delete a existing batch
func (c *Client) DeleteBatch(uuid string, name string) error {
	code, body, err := c.delete("/projects/" + uuid + "/batches/" + name)
	if err != nil {
		return err
	}

	switch code {
	case http.StatusOK:
	case http.StatusNotFound:
		if fmt.Sprintf("%s", body) == `{"error":"Couldn't find Batch with [WHERE \"batches\".\"cluster_id\" = $1 AND \"batches\".\"name\" = $2]"}` {
			return fmt.Errorf("Batch '%s' does not exist", name)
		}
		return fmt.Errorf("Project '%s' does not exist", uuid)
	case http.StatusBadRequest:
		return fmt.Errorf("Deploy probably in progress")
	default:
		return unexpectedHTTPError(code, body)
	}

	return nil
}

// ExecuteBatch execute an existing batch
func (c *Client) ExecuteBatch(projectUUID, name string) error {
	url := fmt.Sprintf("/projects/%s/batches/%s/execute", projectUUID, name)
	code, body, err := c.post(url, nil)
	if err != nil {
		return err
	}

	if code != http.StatusOK {
		return unexpectedHTTPError(code, body)
	}

	return nil
}

func (c *Client) GetBatchesInfo(projectUUID, name string) (RunningBatch, error) {
	batches, err := c.GetBatches(projectUUID)
	if err != nil {
		return RunningBatch{}, err
	}

	for _, batch := range batches {
		if batch.BatchCommon.Name == name {
			return batch, nil
		}
	}

	return RunningBatch{}, fmt.Errorf("Batch %q not found for project %q", name, projectUUID)
}

func (c *RunningBatch) SetEnv(path string) error {
	var env map[string]string

	file, err := os.Open(path)
	if err != nil {
		return errors.New(fmt.Sprintf("Error when reading env file: %s", err))
	}

	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&env)
	if err != nil {
		return errors.New(fmt.Sprintf("Error when unmarshalling env file: %s", err))
	}

	for k, v := range env {
		c.CustomEnvironment[k] = v
	}
	return nil
}

// ConfigBatch calls the API to update command.
func (c *Client) ConfigBatch(batch RunningBatch, projectUUID string) error {
	payload := JSONObject{}
	if len(batch.RunCommand) != 0 {
		payload["run_command"] = batch.RunCommand
	}
	limits := JSONObject{}
	if batch.Limits.Memory >= 0 {
		limits["mem"] = batch.Limits.Memory
	}
	if batch.Limits.CPU >= 0 {
		limits["cpu"] = batch.Limits.CPU
	}
	if batch.Limits.IOPS >= 0 {
		limits["iops"] = batch.Limits.IOPS
	}
	payload["limits"] = limits
	if batch.CustomEnvironment != nil {
		payload["custom_environment"] = batch.CustomEnvironment
	}
	if batch.DockerCapabilities != nil {
		payload["docker_capabilities"] = batch.DockerCapabilities
	}
	if batch.DockerDevices != nil {
		payload["docker_devices"] = batch.DockerDevices
	}

	logger.Debug.Println("Json payload : ", payload)
	code, body, err := c.put(fmt.Sprintf("/projects/%s/batches/%s", projectUUID, batch.BatchCommon.Name), payload)
	if err != nil {
		return err
	}

	switch code {
	case http.StatusOK:
		return nil
	case http.StatusNotFound:
		return errors.New("Batch does not exist")
	default:
		return unexpectedHTTPError(code, body)
	}
}
