package squarescale

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// SchedulingGroup describes a project scheduling group as returned by the SquareScale API
type SchedulingGroup struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// AddSchedulingGroup add a new scheduling group
func (c *Client) AddSchedulingGroup(projectUUID string, name string) (SchedulingGroup, error) {
	var newSchedulingGroup SchedulingGroup

	payload := JSONObject{
		"name": name,
	}
	code, body, err := c.post("/projects/"+projectUUID+"/scheduling_groups", &payload)
	if err != nil {
		return newSchedulingGroup, err
	}

	switch code {
	case http.StatusCreated:
	case http.StatusNotFound:
		return newSchedulingGroup, fmt.Errorf("Project '%s' does not exist", projectUUID)
	case http.StatusConflict:
		return newSchedulingGroup, fmt.Errorf("Scheduling group already exist on project '%s': %s", projectUUID, name)
	default:
		return newSchedulingGroup, unexpectedHTTPError(code, body)
	}

	if err := json.Unmarshal(body, &newSchedulingGroup); err != nil {
		return newSchedulingGroup, err
	}

	return newSchedulingGroup, nil
}

// GetSchedulingGroups gets all the scheduling groups attached to a Project
func (c *Client) GetSchedulingGroups(projectUUID string) ([]SchedulingGroup, error) {
	code, body, err := c.get("/projects/" + projectUUID + "/scheduling_groups")
	if err != nil {
		return []SchedulingGroup{}, err
	}

	switch code {
	case http.StatusOK:
	case http.StatusNotFound:
		return []SchedulingGroup{}, fmt.Errorf("Project '%s' does not exist", projectUUID)
	default:
		return []SchedulingGroup{}, unexpectedHTTPError(code, body)
	}

	var schedulingGroupsByID []SchedulingGroup

	if err := json.Unmarshal(body, &schedulingGroupsByID); err != nil {
		return []SchedulingGroup{}, err
	}

	return schedulingGroupsByID, nil
}

// GetStatefulNodeInfo get information for a scheduling group
func (c *Client) GetSchedulingGroupInfo(projectUUID string, name string) (SchedulingGroup, error) {
	schedulingGroups, err := c.GetSchedulingGroups(projectUUID)
	if err != nil {
		return SchedulingGroup{}, err
	}

	for _, schedulingGroup := range schedulingGroups {
		if schedulingGroup.Name == name {
			return schedulingGroup, nil
		}
	}

	return SchedulingGroup{}, fmt.Errorf("Scheduling group '%s' not found for project '%s'", name, projectUUID)
}
