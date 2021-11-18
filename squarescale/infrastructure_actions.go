package squarescale

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// InfrastructureAction describes an infrastructure action item as returned by the SquareScale API
type InfrastructureAction struct {
	UUID       string    `json:"uuid"`
	ActionType string    `json:"action_type"`
	Status     string    `json:"status"`
	ProjectId  int       `json:"project_id"`
	StartedAt  time.Time `json:"started_at"`
	EndAt      time.Time `json:"end_at"`
	Log        string    `json:"log"`
}

// GetInfrastructureActions gets all the infrastructure Actions attached to a Project
func (c *Client) GetInfrastructureActions(projectUUID string) ([]InfrastructureAction, error) {
	code, body, err := c.get("/projects/" + projectUUID + "/infrastructure_actions")
	if err != nil {
		return []InfrastructureAction{}, err
	}

	switch code {
	case http.StatusOK:
	case http.StatusNotFound:
		return []InfrastructureAction{}, fmt.Errorf("Project '%s' does not exist", projectUUID)
	default:
		return []InfrastructureAction{}, unexpectedHTTPError(code, body)
	}

	var infrastructureActionsByMostRecent []InfrastructureAction

	if err := json.Unmarshal(body, &infrastructureActionsByMostRecent); err != nil {
		return []InfrastructureAction{}, err
	}

	return infrastructureActionsByMostRecent, nil
}
