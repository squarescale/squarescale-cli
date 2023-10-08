package squarescale

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/squarescale/logger"
)

// Project defined how project could see
type Project struct {
	Name                 string    `json:"name"`
	UUID                 string    `json:"uuid"`
	Provider             string    `json:"provider"`
	Region               string    `json:"region"`
	Organization         string    `string:"organization"`
	InfraStatus          string    `json:"infra_status"`
	ClusterSize          int       `json:"cluster_size"`
	NomadNodesReady      int       `json:"nomad_nodes_ready"`
	MonitoringEnabled    bool      `json:"monitoring_enabled"`
	MonitoringEngine     string    `json:"monitoring_engine"`
	SlackWebHook         string    `json:"slack_webhook"`
	ExternalES           string    `json:"external_elasticsearch"`
	Error                string    `json:"error"`
	HybridClusterEnabled bool      `json:"hybrid_cluster_enabled"`
	ProviderLabel        string    `json:"provider_label"`
	NodeSize             string    `json:"node_size"`
	RootDiskSizeGB       int       `json:"root_disk_size_gb"`
	RegionLabel          string    `json:"region_label"`
	Credentials          string    `json:"credentials_name"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
	StatusBefore         string    `json:"status_before"`
	DesiredStateless     int       `json:"desired_stateless"`
	ActualStateless      int       `json:"actual_stateless"`
	DesiredStateful      int       `json:"desired_stateful"`
	ActualStateful       int       `json:"actual_stateful"`
	TfCommand            string    `json:"tf_command"`
}

// Need special decoding
// cf https://stackoverflow.com/questions/37782278/fully-parsing-timestamps-in-golang
// or more complete example
// https://dev.to/arshamalh/how-to-unmarshal-json-in-a-custom-way-in-golang-42m5
type Timestamp struct {
	time.Time
}

// UnmarshalJSON decodes an int64 timestamp into a time.Time object
func (p *Timestamp) UnmarshalJSON(bytes []byte) error {
	// 1. Decode the bytes into an int64
	var raw int64
	err := json.Unmarshal(bytes, &raw)

	if err != nil {
		fmt.Printf("error decoding timestamp: %s\n", err)
		return err
	}

	// 2. Parse the unix timestamp
	p.Time = time.Unix(raw, 0)
	return nil
}

type Notification struct {
//	Component         string    `json:"component_name"`
	Level             string    `json:"level"`
	NotificationType  string    `json:"type"`
	Message           string    `json:"message"`
	NotifiedAt        Timestamp `json:"notified_at"`
	ProjectUUID       string    `json:"project_uuid"`
}

type ClusterMemberDetails struct {
	ConsulName        string `json:"consul_name"`
	ConsulVersion     string `json:"consul_version"`
	CPUArch           string `json:"cpu_arch"`
	CPUCores          string `json:"cpu_cores"`
	CPUFrequency      string `json:"cpu_frequency"`
	CPUModel          string `json:"cpu_model_name"`
	// Drivers
	Hostname          string `json:"hostname"`
	ID                int    `json:"id"`
	InstanceID        string `json:"instance_id"`
	InstanceType      string `json:"instance_type"`
	KernelArch        string `json:"kernel_arch"`
	KernelName        string `json:"kernel_name"`
	KernelVersion     string `json:"kernel_version"`
	Memory            string `json:"memory"`
	Name              string `json:"name"`
	NomadEligibility  string `json:"nomad_eligibility"`
	NomadID           string `json:"nomad_id"`
	NomadStatus       string `json:"nomad_status"`
	NomadVersion      string `json:"nomad_version"`
	OSName            string `json:"os_name"`
	OSVersion         string `json:"os_version"`
	PrivateIP         string `json:"private_ip"`
	// PublicIP
	// SchedulingGroup
	// StatefulNode
	StorageBytesFree  string `json:"storage_bytesfree"`
	StorageBytesTotal string `json:"storage_bytestotal"`
	Zone              string `json:"zone"`
}

type Cluster struct {
	ActualExternal        int                    `json:"actual_external"`
	ActualStateful        int                    `json:"actual_stateful"`
	ActualStateless       int                    `json:"actual_stateless"`
	CurrentSize           int                    `json:"current_size"`
	DesiredExternal       int                    `json:"desired_external"`
	DesiredSize           int                    `json:"desired_size"`
	DesiredStateful       int                    `json:"desired_stateful"`
	DesiredStateless      int                    `json:"desired_stateless"`
	ExternalNodes         []ExternalNode         `json:"external_nodes"`
	ClusterMembersDetails []ClusterMemberDetails `json:"members"`
	RootDiskSize          int                    `json:"root_disk_size_gb"`
	SchedulingGroups      []SchedulingGroup      `json:"scheduling_groups"`
	Status                string                 `json:"status"`
}

type IntegratedServices struct {
	IntegratedServices []IntegratedServiceInfo
}

type IntegratedServiceInfo struct {
	BasicAuth   string     `json:"basic_auth"`
	Enabled     bool       `json:"enabled"`
	IPWhiteList string     `json:"ip_whitelist"`
	Name        string     `json:"name"`
	Prefix      string     `json:"prefix"`
	URLs        [][]string `json:"urls"`
}

// UnmarshalJSON decodes an Integrated Service JSON map into a proper object
func (p *IntegratedServices) UnmarshalJSON(bytes []byte) error {
	// 1. Decode the bytes into a raw interface object
	var raw map[string]IntegratedServiceInfo
	err := json.Unmarshal(bytes, &raw)

	if err != nil {
		fmt.Printf("error decoding integrated service: %s\n", err)
		return err
	}

	res := make([]IntegratedServiceInfo, len(raw))
	i := 0
	for _, v := range raw {
		res[i] = v
		i++
	}
	*p = IntegratedServices{
		IntegratedServices: res,
	}
	// 2. Parse the unix timestamp
	//p.Time = time.Unix(raw, 0)
	return nil
}

type Infrastructure struct {
	Action             string             `json:"action"`
	Cluster            Cluster            `json:"cluster"`
	Database           Database           `json:"db"`
	IntegratedServices IntegratedServices `json:"integrated_services"`
	LoadBalancer       LoadBalancer       `json:"lb"`
	MonitoringEngine   string             `json:"monitoring_engine"`
	NodeSize           string             `json:"node_size"`
	CloudProvider      string             `json:"provider"`
	CredentialName     string             `json:"provider_credential_name"`
	CloudProviderLabel string             `json:"provider_label"`
	//`json:"redis_databases"`
	Region             string             `json:"region"`
	RegionLabel        string             `json:"region_label"`
	RootDiskSize       int                `json:"root_disk_size_gb"`
	Status             string             `json:"status"`
	TFRunAt            time.Time          `json:"terraform_run_at"`
	Type               string             `json:"type"`
}

type GenericVariable struct {
	Key        string
	Value      interface{}
	Predefined bool
}

type ProjectDetails struct {
	CreatedAt             time.Time         `json:"created_at"`
	ExternalElasticSearch string            `json:"external_elasticsearch"`
	Environment           []GenericVariable `json:"global_environment"`
	HighAvailability      bool              `json:"high_availability"`
	HybridClusterEnabled  bool              `json:"hybrid_cluster_enabled"`
	Infrastructure        Infrastructure    `json:"infra"`
	//`json:"integrated_services"`
	//`json:"intentions"`
	//`json:"managed_services"`
	Name                  string            `json:"name"`
	Organization          string            `json:"organization_name"`
	//`json:"services"`
	SlackWebHook          string            `json:"slack_webhook"`
	UpdatedAt             time.Time         `json:"updated_at"`
	User                  User              `json:"user"`
	UUID                  string            `json:"uuid"`
}

type ProjectWithAllDetails struct {
	Notifications  []Notification `json:"notifications"`
	Project        ProjectDetails `json:"project"`
}

// UnprovisionError defined how export provision errors
type UnprovisionError struct {
	Errors struct {
		Unprovision []string `json:unprovision`
	} `json:errors`
}

// See front/utils/infraAction
func (p *Project) ProjectAction() string {
	if p.TfCommand == "destroy" {
		return "destroying"
	}
	if p.StatusBefore == "no_infra" {
		return "building"
	}
	return "updating"
}

// See front/src/components/InfraStatusBadge
func (p *Project) ProjectStatus() string {
	if p.InfraStatus == "provisionning" {
		if p.ProjectAction() == "destroying" {
			return "unprovisionning"
		}
		return "provisionning"
	}
	return p.InfraStatus
}

// See front/src/components/ClusterStatusBadge
func (p *Project) ProjectStateLessCount() string {
	dsl := p.ClusterSize
	if p.DesiredStateless > dsl {
		dsl = p.DesiredStateless
	}
	asl := p.NomadNodesReady
	if p.ActualStateless > asl {
		asl = p.ActualStateless
	}
	return fmt.Sprintf("%d/%d", asl, dsl)
}

func (p *Project) ProjectStateFulCount() string {
	return fmt.Sprintf("%d/%d", p.ActualStateful, p.DesiredStateful)
}

func projectSettings(name string, cluster ClusterConfig) JSONObject {
	projectSettings := JSONObject{
		"name": name,
	}
	for attr, value := range cluster.ProjectCreationSettings() {
		projectSettings[attr] = value
	}
	return projectSettings
}

// CreateProject asks the SquareScale platform to create a new project
func (c *Client) CreateProject(payload *JSONObject) (newProject Project, err error) {
	_, ok := (*payload)["credential_name"]
	if !ok {
		return newProject, errors.New("Credential is mandatory")
	}

	code, body, err := c.post("/projects", payload)
	if err != nil {
		return newProject, err
	}

	if code != http.StatusCreated {
		return newProject, unexpectedHTTPError(code, body)
	}

	err = json.Unmarshal(body, &newProject)
	if err != nil {
		return newProject, err
	}

	if newProject.Error != "" {
		err = errors.New(newProject.Error)
	}

	return newProject, err
}

// ProvisionProject asks the SquareScale platform to provision the project
func (c *Client) ProvisionProject(projectUUID string) (err error) {

	code, body, err := c.post(fmt.Sprintf("/projects/%s/provision", projectUUID), nil)
	if err != nil {
		return err
	}

	if code != http.StatusNoContent {
		return unexpectedHTTPError(code, body)
	}

	return nil
}

// UNProvisionProject asks the SquareScale platform to provision the project
func (c *Client) UNProvisionProject(projectUUID string) (err error) {

	code, body, err := c.post(fmt.Sprintf("/projects/%s/unprovision", projectUUID), nil)
	if err != nil {
		return err
	}

	if code != http.StatusNoContent {
		return unexpectedHTTPError(code, body)
	}

	return nil
}

// ListProjects asks the SquareScale service for available projects.
func (c *Client) ListProjects() ([]Project, error) {
	code, body, err := c.get("/projects")
	if err != nil {
		return nil, err
	}

	if code != http.StatusOK {
		return nil, unexpectedHTTPError(code, body)
	}

	var projectsJSON []Project
	err = json.Unmarshal(body, &projectsJSON)
	if err != nil {
		return nil, err
	}

	return projectsJSON, nil
}

// FullListProjects asks the SquareScale service for available projects.
func (c *Client) FullListProjects() ([]Project, error) {
	projects, err := c.ListProjects()
	if err != nil {
		return nil, err
	}

	organizations, err := c.ListOrganizations()
	if err != nil {
		return nil, err
	}

	projectCount := len(projects)

	for _, organization := range organizations {
		projectCount += len(organization.Projects)
	}

	allProjects := make([]Project, projectCount)
	projectCount = 0

	for _, project := range projects {
		allProjects[projectCount].Name = project.Name
		allProjects[projectCount].UUID = project.UUID
		projectCount++
	}

	for _, organization := range organizations {
		for _, project := range organization.Projects {
			allProjects[projectCount].Name = organization.Name + "/" + project.Name
			allProjects[projectCount].UUID = project.UUID
			projectCount++
		}
	}

	return allProjects, nil
}

// ProjectByName get UUID for a project name
func (c *Client) ProjectByName(projectName string) (string, error) {
	projects, err := c.FullListProjects()
	if err != nil {
		return "", err
	}

	for _, project := range projects {
		if project.Name == projectName {
			return project.UUID, nil
		}
	}

	return "", fmt.Errorf("Project '%s' not found", projectName)
}

// GetProjectDetails return the detailed informations of the project
func (c *Client) GetProjectDetails(project string) (*ProjectWithAllDetails, error) {
	code, body, err := c.get("/project_info/" + project)
	if err != nil {
		return nil, err
	}

	switch code {
	case http.StatusOK:
	case http.StatusNotFound:
		return nil, fmt.Errorf("Project '%s' not found", project)
	default:
		return nil, unexpectedHTTPError(code, body)
	}

	var details ProjectWithAllDetails
	err = json.Unmarshal(body, &details)
	if err != nil {
		fmt.Printf("ERROR %+v\n", err)
		return nil, err
	}

	return &details, nil
}

// GetProject return the basic infos of the project
func (c *Client) GetProject(project string) (*Project, error) {
	code, body, err := c.get("/projects/" + project)
	if err != nil {
		return nil, err
	}

	switch code {
	case http.StatusOK:
	case http.StatusNotFound:
		return nil, fmt.Errorf("Project '%s' not found", project)
	default:
		return nil, unexpectedHTTPError(code, body)
	}

	var basicInfos Project
	err = json.Unmarshal(body, &basicInfos)
	if err != nil {
		return nil, err
	}

	return &basicInfos, nil
}

// WaitProject wait project provisioning
func (c *Client) WaitProject(projectUUID string, timeToWait int64) (string, error) {
	project, err := c.GetProject(projectUUID)
	if err != nil {
		return "", err
	}

	logger.Info.Println("wait for project : ", projectUUID)

	for err == nil && project == nil && !(project.InfraStatus == "ok" || project.InfraStatus == "error") {
		time.Sleep(time.Duration(timeToWait) * time.Second)
		project, err = c.GetProject(projectUUID)
		if err != nil {
			return "", err
		}
		logger.Debug.Println("project status update: ", projectUUID)
	}

	if project == nil {
		return "", fmt.Errorf("Project %s not found", projectUUID)
	}

	if project.InfraStatus == "error" {
		actions, err := c.GetInfrastructureActions(project.UUID)
		if err != nil {
			return "", err
		}
		if len(actions) == 0 {
			return project.InfraStatus, errors.New("Unable to retrieve latest project deployment log")
		} else {
			return project.InfraStatus, errors.New(actions[0].Log)
		}
	}

	return project.InfraStatus, err
}

// ProjectUnprovision unprovisions a project
func (c *Client) ProjectUnprovision(project string) error {
	code, body, err := c.post("/projects/"+project+"/unprovision", nil)
	if err != nil {
		return err
	}

	switch code {
	case http.StatusOK:
	case http.StatusAccepted:
	case http.StatusNoContent:
	case http.StatusNotFound:
		return fmt.Errorf("Project '%s' not found", project)
	case http.StatusUnprocessableEntity:
		var errJSON UnprovisionError
		err = json.Unmarshal(body, &errJSON)
		if err != nil {
			return err
		}
		return fmt.Errorf("Operation failed: %s", errJSON.Errors.Unprovision)
	default:
		return unexpectedHTTPError(code, body)
	}

	return nil
}

// ProjectDelete deletes an unprovisionned project
func (c *Client) ProjectDelete(project string) error {
	code, body, err := c.delete("/projects/" + project)
	if err != nil {
		return err
	}

	switch code {
	case http.StatusOK:
	case http.StatusNoContent:
	case http.StatusNotFound:
		return fmt.Errorf("Project '%s' not found", project)
	default:
		return unexpectedHTTPError(code, body)
	}

	return nil
}

// ProjectLogs gets the logs for a project container.
func (c *Client) ProjectLogs(project string, container string, after string) ([]string, string, error) {
	query := ""
	if after != "" {
		query = "?after=" + url.QueryEscape(after)
	}

	code, body, err := c.get("/projects/" + project + "/logs/" + url.QueryEscape(container) + query)
	if err != nil {
		return []string{}, "", err
	}

	switch code {
	case http.StatusOK:
	case http.StatusBadRequest:
		return []string{}, "", fmt.Errorf("Project '%s' not found", project)
	case http.StatusNotFound:
		return []string{}, "", fmt.Errorf("Container '%s' is not found for project '%s'", container, project)
	default:
		return []string{}, "", unexpectedHTTPError(code, body)
	}

	var response []struct {
		Timestamp     string `json:"timestamp"`
		ProjectName   string `json:"project_name"`
		ContainerName string `json:"container_name"`
		Error         bool   `json:"error"`
		Type          string `json:"type"`
		Message       string `json:"message"`
		Level         int    `json:"level"`
	}
	err = json.Unmarshal(body, &response)

	if err != nil {
		return []string{}, "", err
	}

	var messages []string
	for _, log := range response {
		var linePattern string
		var lt string
		var containerWithBrackets string
		var level string
		if log.Type == "docker" {
			lt = "docker"
		} else if log.Type == "nomad" {
			lt = "nomad "
		} else if log.Type == "event" {
			lt = "sqsc  "
		} else {
			lt = ""
		}
		if log.ContainerName == "" {
			containerWithBrackets = ""
		} else {
			containerWithBrackets = "[" + log.ContainerName + "] "
		}
		if log.Level >= 5 {
			level = "INFO "
		} else if log.Level >= 4 {
			level = "WARN "
		} else {
			level = "ERROR"
		}
		linePattern = "%s %s %s-- [%s] %s"
		if log.Error {
			linePattern = "\033[0;33m" + linePattern + "\033[0m"
		}
		t, err := time.Parse(time.RFC3339Nano, log.Timestamp)
		if err == nil {
			formatedTime := t.Format("2006-01-02 15:04:05.999")
			padTime := fmt.Sprintf("%-23s", formatedTime)
			messages = append(messages, fmt.Sprintf(linePattern, padTime, lt, containerWithBrackets, level, log.Message))
		} else {
			messages = append(messages, fmt.Sprintf(linePattern, log.Timestamp, lt, containerWithBrackets, level, log.Message))
		}
	}
	var lastTimestamp string
	if len(response) == 0 {
		lastTimestamp = ""
	} else {
		lastTimestamp = response[len(response)-1].Timestamp
	}

	return messages, lastTimestamp, nil
}

// ConfigProjectSettings configure project settings
func (c *Client) ConfigProjectSettings(projectUUID string, project Project) error {
	payload := &JSONObject{"hybrid_cluster_enabled": project.HybridClusterEnabled}

	code, body, err := c.put(fmt.Sprintf("/projects/%s", projectUUID), payload)
	if err != nil {
		return err
	}

	switch code {
	case http.StatusOK:
		return nil
	case http.StatusNotFound:
		return fmt.Errorf("Project '%s' not found", projectUUID)
	default:
		return unexpectedHTTPError(code, body)
	}
}
