package squarescale

import (
	"fmt"
	"net/http"
)

// AddImage asks the Squarescale service to attach an image to the project.
func (c *Client) AddImage(project, dockerImage string, instances *int) error {
	payload := JSONObject{
		"docker_image": JSONObject{
			"name": dockerImage,
		},
	}
	if instances != nil {
		payload["size"] = *instances
	}
	code, body, err := c.post("/projects/"+project+"/docker_images", &payload)
	if err != nil {
		return err
	}

	switch code {
	case http.StatusCreated:
		return nil
	default:
		return readErrors(body, fmt.Sprintf("Cannot add docker image '%s' to project '%s'", dockerImage, project))
	}
}
