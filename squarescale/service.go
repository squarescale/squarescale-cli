package squarescale

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/squarescale/logger"
)

type ServiceEnv struct {
	Key        string `json:"key"`
	Value      string `json:"value"`
	Predefined bool   `json:"predefined"`
}

// Service describes a project container as returned by the SquareScale API
type Service struct {
	ID                  int               `json:"container_id"`
	Name                string            `json:"name"`
	RunCommand          string            `json:"run_command"`
	Entrypoint          string            `json:"entrypoint"`
	Running             int               `json:"running"`
	Size                int               `json:"size"`
	WebPort             int               `json:"web_port"`
	RefreshCallbacks    []string          `json:"refresh_callbacks"`
	Limits              ServiceLimits     `json:"limits"`
	CustomEnv           []ServiceEnv      `json:"custom_environment"`
	SchedulingGroups    []SchedulingGroup `json:"scheduling_groups"`
	DockerCapabilities  []string          `json:"docker_capabilities"`
	DockerDevices       []DockerDevice    `json:"docker_devices"`
	AutoStart           bool              `json:"auto_start"`
	MaxClientDisconnect string            `json:"max_client_disconnect"`
	Volumes             []VolumeToBind    `json:"volumes"`
}

type ServiceBody struct {
	ID                  int               `json:"container_id"`
	Name                string            `json:"name"`
	RunCommand          []string          `json:"run_command"`
	Entrypoint          string            `json:"entrypoint"`
	Running             int               `json:"running"`
	Size                int               `json:"size"`
	WebPort             int               `json:"web_port"`
	RefreshCallbacks    []string          `json:"refresh_callbacks"`
	Limits              ServiceLimits     `json:"limits"`
	CustomEnv           []ServiceEnv      `json:"custom_environment"`
	SchedulingGroups    []SchedulingGroup `json:"scheduling_groups"`
	DockerCapabilities  []string          `json:"docker_capabilities"`
	DockerDevices       []DockerDevice    `json:"docker_devices"`
	AutoStart           bool              `json:"auto_start"`
	MaxClientDisconnect int               `json:"max_client_disconnect"`
	Volumes             []VolumeToBind    `json:"volumes"`
}

func (c *Service) SetEnv(path string) error {

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

	c.CustomEnv = make([]ServiceEnv, len(env), len(env))
	i := 0

	for k, v := range env {
		c.CustomEnv[i] = ServiceEnv{Key: k, Value: v}
		i++
	}
	return nil
}

func (c *Service) SetEnvParams(params []string) error {
	for _, p := range params {
		v := strings.Split(p, "=")
		if len(v) != 2 {
			return errors.New(fmt.Sprintf("environment parameter %v not in the form param=value", p))
		}
		found := -1
		for i, curParam := range c.CustomEnv {
			if curParam.Key == v[0] {
				found = i
				break
			}
		}
		if found >= 0 {
			c.CustomEnv[found].Value = v[1]
		} else {
			c.CustomEnv = append(c.CustomEnv, ServiceEnv{Key: v[0], Value: v[1]})
		}
	}
	return nil
}

type ServiceLimits struct {
	Memory int `json:"mem"`
	CPU    int `json:"cpu"`
	IOPS   int `json:"iops"`
}

// GetContainers gets all the services attached to a Project
func (c *Client) GetServices(projectUUID string) ([]Service, error) {
	code, body, err := c.get("/projects/" + projectUUID + "/services")
	if err != nil {
		return []Service{}, err
	}

	switch code {
	case http.StatusOK:
	case http.StatusNotFound:
		return []Service{}, fmt.Errorf("Project '%s' does not exist", projectUUID)
	default:
		return []Service{}, unexpectedHTTPError(code, body)
	}

	var servicesBody []ServiceBody

	if err := json.Unmarshal(body, &servicesBody); err != nil {
		return []Service{}, err
	}

	var services []Service

	for _, c := range servicesBody {
		service := &Service{
			ID:                  c.ID,
			Name:                c.Name,
			RunCommand:          strings.Join(c.RunCommand, " "),
			Entrypoint:          c.Entrypoint,
			Running:             c.Running,
			Size:                c.Size,
			WebPort:             c.WebPort,
			RefreshCallbacks:    c.RefreshCallbacks,
			Limits:              c.Limits,
			CustomEnv:           c.CustomEnv,
			SchedulingGroups:    c.SchedulingGroups,
			DockerCapabilities:  c.DockerCapabilities,
			DockerDevices:       c.DockerDevices,
			AutoStart:           c.AutoStart,
			MaxClientDisconnect: strconv.Itoa(c.MaxClientDisconnect),
			Volumes:             c.Volumes,
		}
		services = append(services, *service)
	}

	return services, nil
}

// GetServiceInfo get the service of a project based on its name.
func (c *Client) ScheduleService(projectUUID, name string) error {
	code, body, err := c.post(fmt.Sprintf("/projects/%s/services/%s/schedule", projectUUID, name), nil)
	if err != nil {
		return err
	}

	if code != http.StatusOK {
		return unexpectedHTTPError(code, body)
	}

	return nil
}

// GetServiceInfo get the service of a project based on its name.
func (c *Client) GetServiceInfo(projectUUID, name string) (Service, error) {
	// TODO: if services are to be retrieved with Docker image informations (like for service add)
	// then GetServices should call GET on project_info/UUID and not project/UUID
	services, err := c.GetServices(projectUUID)
	if err != nil {
		return Service{}, err
	}

	for _, service := range services {
		if service.Name == name {
			return service, nil
		}
	}

	return Service{}, fmt.Errorf("Service '%s' not found for project '%s'", name, projectUUID)
}

// ConfigService calls the API to update the number of instances and update command.
func (c *Client) ConfigService(service Service) error {
	svcConf := JSONObject{}
	if len(service.RunCommand) != 0 {
		svcConf["run_command"] = service.RunCommand
	}
	if service.Size > 0 {
		svcConf["size"] = service.Size
	}
	limits := JSONObject{}
	if service.Limits.Memory >= 0 {
		limits["mem"] = service.Limits.Memory
	}
	if service.Limits.CPU >= 0 {
		limits["cpu"] = service.Limits.CPU
	}
	if service.Limits.IOPS > 0 {
		limits["iops"] = service.Limits.IOPS
	}
	svcConf["limits"] = limits
	if service.CustomEnv != nil {
		svcConf["custom_environment"] = service.CustomEnv
	}
	if len(service.SchedulingGroups) != 0 {
		svcConf["scheduling_groups"] = getSchedulingGroupsIds(service.SchedulingGroups)
	}
	if service.DockerCapabilities != nil {
		svcConf["docker_capabilities"] = service.DockerCapabilities
	}
	if service.DockerDevices != nil {
		svcConf["docker_devices"] = service.DockerDevices
	}
	if service.MaxClientDisconnect != "" {
		svcConf["max_client_disconnect"] = service.MaxClientDisconnect
	} else {
		svcConf["max_client_disconnect"] = "0"
	}
	svcConf["auto_start"] = service.AutoStart

	payload := &JSONObject{"container": svcConf}
	logger.Debug.Println("Json payload : ", payload)
	code, body, err := c.put(fmt.Sprintf("/containers/%d", service.ID), payload)
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

func (c *Client) DeleteService(service Service) error {
	url := fmt.Sprintf("/containers/%d", service.ID)
	code, body, err := c.delete(url)
	if err != nil {
		return err
	}

	switch code {
	// strange reading this but seems like OK does continue without any return
	case http.StatusOK:
	case http.StatusNotFound:
		return fmt.Errorf("service with id '%d' does not exist", service.ID)
	default:
		return unexpectedHTTPError(code, body)
	}

	return nil
}

func getSchedulingGroupsIds(schedulingGroups []SchedulingGroup) []int {
	var schedulingGroupsIds []int

	for _, schedulingGroup := range schedulingGroups {
		schedulingGroupsIds = append(schedulingGroupsIds, schedulingGroup.ID)
	}

	return schedulingGroupsIds
}

// AddService asks the SquareScale service to attach an image to the project.
func (c *Client) AddService(projectUUID string, payload JSONObject) error {
	code, body, err := c.post("/projects/"+projectUUID+"/docker_images", &payload)
	if err != nil {
		return fmt.Errorf("Cannot add docker image '%v' to project '%s' (%d %s)\n\t%s", payload["docker_image"], projectUUID, code, http.StatusText(code), err)
	}

	switch code {
	case http.StatusCreated:
		return nil
	default:
		return unexpectedHTTPError(code, body)
	}
}

func DockerImage(name, username, password string) JSONObject {
	o := JSONObject{
		"name": name,
	}

	if username != "" && password != "" {
		o["private"] = true
		o["username"] = username
		o["password"] = password
	}

	return o
}
