package squarescale

import (
	"fmt"
	"net/http"
)

// AddImage asks the Squarescale service to attach an image to the project.
func (c *Client) AddImage(project, name, username, password string, instances int, serviceName string, volumeToBind []VolumeToBind) error {
	payload := JSONObject{
		"docker_image": c.dockerImage(
			name,
			username,
			password,
		),

		"size":            instances,
		"volumes_to_bind": volumeToBind,
	}

	if len(serviceName) > 0 {
		payload["name"] = serviceName
	}

	code, body, err := c.post("/projects/"+project+"/docker_images", &payload)
	if err != nil {
		return fmt.Errorf("Cannot add docker image '%s' to project '%s' (%d %s)\n\t%s", name, project, code, http.StatusText(code), err)
	}

	switch code {
	case http.StatusCreated:
		return nil
	default:
		return readErrors(body, fmt.Sprintf("Cannot add docker image '%s' to project '%s' (%d %s)", name, project, code, http.StatusText(code)))
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
