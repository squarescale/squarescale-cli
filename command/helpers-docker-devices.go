package command

import (
	"errors"
	"regexp"
	"strings"

	"github.com/squarescale/squarescale-cli/squarescale"
)

func getDockerDevicesArray(devicesString string) ([]squarescale.DockerDevice, error) {
	devicesArray := strings.Split(devicesString, ",")
	err := checkDockerDevicesRegex(devicesArray)
	if err != nil {
		return nil, err
	}
	dockerDevicesArray := []squarescale.DockerDevice{}
	for _, device := range devicesArray {
		mapping := strings.Split(device, ":")
		dockerDevice := squarescale.DockerDevice{SRC: mapping[0]}
		if len(mapping) > 1 && mapping[1] != "" {
			dockerDevice.DST = mapping[1]
		}
		if len(mapping) > 2 && mapping[2] != "" {
			dockerDevice.OPT = mapping[2]
		}
		dockerDevicesArray = append(dockerDevicesArray, dockerDevice)
	}
	return dockerDevicesArray, nil
}

func checkDockerDevicesRegex(devices []string) error {
	pattern, compError := regexp.Compile(`^(([^:\s])+(:[^:\s]+){0,2})$|^([^:\s]+):{2}([^:\s]+)$`)
	if compError != nil {
		return compError
	}
	for _, device := range devices {
		if !pattern.MatchString(device) {
			return errors.New("wrong mapping format for " + device)
		}
	}
	return nil
}
