package squarescale

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

// DisableLB asks the Squarescale service to deactivate the load balancer.
func DisableLB(sqscURL, token, project string) error {
	payload := jsonObject{
		"project": jsonObject{
			"load_balancer": false,
		},
	}

	return updateLBConfig(sqscURL, token, project, payload)
}

// ConfigLB sets the load balancer configuration for a given project.
func ConfigLB(sqscURL, token, project string, container, port int) error {
	payload := jsonObject{
		"project": jsonObject{
			"load_balancer": true,
		},
		"web-container": strconv.Itoa(container),
		"container": jsonObject{
			strconv.Itoa(container): jsonObject{
				"web-port": strconv.Itoa(port),
			},
		},
	}

	return updateLBConfig(sqscURL, token, project, payload)
}

func updateLBConfig(sqscURL, token, project string, payload jsonObject) error {
	payloadBytes, err := json.Marshal(&payload)
	if err != nil {
		return err
	}

	req := SqscRequest{
		Method: "POST",
		URL:    sqscURL + "/projects/" + project + "/web-ports",
		Token:  token,
		Body:   payloadBytes,
	}

	res, err := doRequest(req)
	if err != nil {
		return err
	}

	if res.Code == http.StatusNotFound {
		return fmt.Errorf("Project '%s' not found", project)
	}

	if res.Code != http.StatusOK {
		return fmt.Errorf("'%s %s' return code: %d", req.Method, req.URL, res.Code)
	}

	return nil
}
