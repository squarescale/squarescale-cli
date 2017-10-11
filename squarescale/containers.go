package squarescale

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

// Container describes a project container as returned by the Squarescale API
type Container struct {
	ID                   int                  `json:"id"`
	ShortName            string               `json:"short_url"`
	PreCommand           []string             `json:"pre_command"`
	RunCommand           []string             `json:"run_command"`
	Running              int                  `json:"running"`
	Size                 int                  `json:"size"`
	Web                  bool                 `json:"web"`
	WebPort              int                  `json:"web_port"`
	Type                 string               `json:"type"`
	BuildStatus          string               `json:"build_status"`
	BuildOutOfDate       bool                 `json:"build_out_of_date"`
	Scheduled            bool                 `json:"scheduled"`
	RepositoryConfigured bool                 `json:"repository_configured"`
	PreCommandStatus     string               `json:"pre_command_status"`
	RefreshCallbacks     []string             `json:"refresh_callbacks"`
	BuildCallbacks       []string             `json:"build_callbacks"`
	BuildService         string               `json:"build_service"`
	Limits               ContainerLimits      `json:"limits"`
	NomadStatus          ContainerNomadStatus `json:"nomad_status"`
}

type ContainerNomadStatus struct {
	Service    JobNomadStatus `json:"service"`
	PreCommand JobNomadStatus `json:"pre_command"`
}

type JobNomadStatus struct {
	Name             string
	PlacementFailure bool
	Placements       []JobNomadPlacement
}

func (placements *JobNomadStatus) Errors(kind string) []string {
	var errors []string
	for _, placement := range placements.Placements {
		msg := "Missing resources to schedule " + kind
		n := 0
		for dimension, num := range placement.DimensionExhausted {
			if n == 0 {
				msg += ": "
			} else {
				msg += ", "
			}
			msg += fmt.Sprintf("%s on %d node(s)", dimension, num)
			n++
		}
		if n == 0 {
			msg += fmt.Sprintf(": %d exhausted node out of %d", placement.NodesExhausted, placement.NodesEvaluated)
		} else if n > 1 {
			msg += fmt.Sprintf(" (total of %d exhausted node out of %d)", placement.NodesExhausted, placement.NodesEvaluated)
		}
		errors = append(errors, msg)
	}
	return errors
}

type JobNomadPlacement struct {
	Errors             []JobNomadPlacementError
	NodesEvaluated     int
	NodesExhausted     int
	ClassFiltered      map[string]int
	ConstraintFiltered map[string]int
	ClassExhausted     map[string]int
	DimensionExhausted map[string]int
}

type JobNomadPlacementError struct {
	Message string
}

type ContainerLimits struct {
	Memory int `json:"mem"`
	CPU    int `json:"cpu"`
	IOPS   int `json:"iops"`
	Net    int `json:"net"`
}

func (c *Container) Status() (string, string) {
	if !c.RepositoryConfigured {
		return "no_ci", "Travis not configured"
	} else if c.BuildStatus == "no_build" {
		return "no_build", "Build not started"
	} else if c.BuildStatus == "error" {
		return "build_error", "Build error"
	} else if c.BuildStatus == "running" {
		return "build_progress", "Build running"
	} else if c.BuildOutOfDate {
		return "build_too_old", "Build too old"
	} else if !c.Scheduled {
		return "not_scheduled", "Job not scheduled yet"
	} else if c.PreCommandStatus == "error" {
		return "update_error", "Update command failed"
	} else if c.PreCommandStatus != "ok" {
		return "update_progress", "Update command in progress"
	} else {
		return "running", "Running"
	}
}

// GetContainers gets all the containers attached to a Project
func (c *Client) GetContainers(project string) ([]Container, error) {
	code, body, err := c.get("/projects/" + project + "/containers")
	if err != nil {
		return []Container{}, err
	}

	switch code {
	case http.StatusOK:
	case http.StatusNotFound:
		return []Container{}, fmt.Errorf("Project '%s' does not exist", project)
	default:
		return []Container{}, unexpectedHTTPError(code, body)
	}

	var containersByID []Container

	if err := json.Unmarshal(body, &containersByID); err != nil {
		return []Container{}, err
	}

	return containersByID, nil
}

// GetContainerInfo gets the container of a project based on its short name.
func (c *Client) GetContainerInfo(project, shortName string) (Container, error) {
	containers, err := c.GetContainers(project)
	if err != nil {
		return Container{}, err
	}

	for _, container := range containers {
		if container.ShortName == shortName {
			return container, nil
		}
	}

	return Container{}, fmt.Errorf("Container '%s' not found for project '%s'", shortName, project)
}

// ConfigContainer calls the API to update the number of instances and update command.
func (c *Client) ConfigContainer(container Container) error {
	cont := jsonObject{}
	if container.PreCommand != nil {
		cont["pre_command"] = container.PreCommand
	}
	if container.RunCommand != nil {
		cont["run_command"] = container.RunCommand
	}
	if container.Size > 0 {
		cont["size"] = container.Size
	}
	if container.BuildService != "" {
		cont["build_service"] = container.BuildService
	}
	limits := jsonObject{}
	if container.Limits.Memory >= 0 {
		limits["mem"] = container.Limits.Memory
	}
	if container.Limits.CPU >= 0 {
		limits["cpu"] = container.Limits.CPU
	}
	if container.Limits.IOPS >= 0 {
		limits["iops"] = container.Limits.IOPS
	}
	if container.Limits.Net >= 0 {
		limits["net"] = container.Limits.Net
	}
	cont["limits"] = limits

	payload := &jsonObject{"container": cont}
	code, body, err := c.put(fmt.Sprintf("/containers/%d", container.ID), payload)
	if err != nil {
		return err
	}

	switch code {
	case http.StatusOK:
		return nil
	case http.StatusNotFound:
		return errors.New("Container does not exist")
	default:
		return unexpectedHTTPError(code, body)
	}
}

func (c *Client) DeleteContainer(container Container) error {
	url := fmt.Sprintf("/containers/%d", container.ID)
	code, body, err := c.delete(url)
	if err != nil {
		return err
	}

	switch code {
	case http.StatusOK:
	case http.StatusNotFound:
		return fmt.Errorf("Container with id '%d' does not exist", container.ID)
	default:
		return unexpectedHTTPError(code, body)
	}

	return nil
}
