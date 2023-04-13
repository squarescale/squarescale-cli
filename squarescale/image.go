package squarescale

import (
	"fmt"
	"net/http"
)

// AddImage asks the SquareScale service to attach an image to the project.
func (c *Client) AddImage(projectUUID, name, username, password, entrypoint, runCommand string, instances int, serviceName string, volumeToBind []VolumeToBind, dockerCapabilities []string, dockerDevices []DockerDevice, autostart bool) error {
	payload := JSONObject{
		"docker_image": c.dockerImage(
			name,
			username,
			password,
		),
		"auto_start":      autostart,
		"size":            instances,
		"volumes_to_bind": volumeToBind,
	}

	if len(serviceName) > 0 {
		payload["name"] = serviceName
	}

	if len(entrypoint) > 0 {
		payload["entrypoint"] = entrypoint
	}

	if len(runCommand) > 0 {
		payload["run_command"] = runCommand
	}

	if len(dockerDevices) > 0 {
		payload["docker_devices"] = dockerDevices
	}

	payload["docker_capabilities"] = dockerCapabilities

	code, body, err := c.post("/projects/"+projectUUID+"/docker_images", &payload)
	if err != nil {
		return fmt.Errorf("Cannot add docker image '%s' to project '%s' (%d %s)\n\t%s", name, projectUUID, code, http.StatusText(code), err)
	}

	switch code {
	case http.StatusCreated:
		return nil
	default:
		return unexpectedHTTPError(code, body)
	}
}
func (c *Client) dockerImage(name, username, password string) JSONObject {
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
