package command

import (
	"log"
	"strings"

	"github.com/squarescale/squarescale-cli/squarescale"
)

func parseVolumesToBind(volumes string) []squarescale.VolumeToBind {
	volumesSplited := strings.Split(volumes, ",")
	volumesLength := len(volumesSplited)
	if volumes == "" {
		volumesLength = 0
	}
	volumesToBind := make([]squarescale.VolumeToBind, volumesLength)
	if volumesLength >= 1 {
		for index, volume := range volumesSplited {
			volumeSplited := strings.Split(volume, ":")
			volumeName := volumeSplited[0]
			mountPoint := volumeSplited[1]
			var readOnly bool
			if len(volumeSplited) > 2 {
				switch strings.ToLower(volumeSplited[2]) {
				case "ro":
					readOnly = true
				case "rw":
					readOnly = false
				default:
					log.Fatal("Read only parameter must be RO or RW")
				}
			} else {
				readOnly = false
			}
			volumesToBind[index] = squarescale.VolumeToBind{volumeName, mountPoint, readOnly}
		}
	}
	return volumesToBind
}
