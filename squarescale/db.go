package squarescale

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// ScaleDB calls the Squarescale API to update database scale options for a given project.
func ScaleDB(sqscURL, token, project, instance string) error {
	var payload struct {
		Project struct {
			Enabled       bool   `json:"db_enabled"`
			InstanceClass string `json:"db_instance_class"`
		} `json:"project"`
	}

	payload.Project.Enabled = true
	payload.Project.InstanceClass = "db.t2." + instance
	payloadBytes, err := json.Marshal(&payload)
	if err != nil {
		return err
	}

	req := SqscRequest{
		Method: "POST",
		URL:    sqscURL + "/projects/" + project + "/cluster",
		Token:  token,
		Body:   payloadBytes,
	}

	res, err := doRequest(req)
	if err != nil {
		return err
	}

	if res.Code != http.StatusNoContent {
		return fmt.Errorf("'%s %s' return code: %d", req.Method, req.URL, res.Code)
	}

	return nil
}
