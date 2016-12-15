package squarescale

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

// DisableLB asks the Squarescale service to deactivate the load balancer.
func DisableLB(sqscURL, token, project string) error {
	payload := &jsonObject{
		"project": jsonObject{
			"load_balancer": false,
		},
	}

	return updateLBConfig(sqscURL, token, project, payload)
}

// ConfigLB sets the load balancer configuration for a given project.
func ConfigLB(sqscURL, token, project string, container, port int) error {
	payload := &jsonObject{
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

func updateLBConfig(sqscURL, token, project string, payload *jsonObject) error {
	url := sqscURL + "/projects/" + project + "/web-ports"
	code, _, err := post(url, token, payload)
	if err != nil {
		return err
	}

	if code == http.StatusNotFound {
		return fmt.Errorf("Project '%s' not found", project)
	}

	if code != http.StatusOK {
		return unexpectedError(code)
	}

	return nil
}

// LoadBalancerEnabled asks if the project load balancer is enabled.
func LoadBalancerEnabled(sqscURL, token, project string) (bool, error) {
	url := sqscURL + "/projects/" + project
	code, body, err := get(url, token)
	if err != nil {
		return false, err
	}

	if code == http.StatusNotFound {
		return false, fmt.Errorf("Project '%s' not found", project)
	}

	if code != http.StatusOK {
		return false, unexpectedError(code)
	}

	var response struct {
		Enabled bool `json:"load_balancer"`
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		return false, err
	}

	return response.Enabled, nil
}
