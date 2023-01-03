package command

import (
	"fmt"
	"strings"

	"github.com/squarescale/squarescale-cli/squarescale"
)

func getSchedulingGroupsArray(UUID string, client *squarescale.Client, schedulingGroups string) []squarescale.SchedulingGroup {
	if schedulingGroups == "" {
		return []squarescale.SchedulingGroup{}
	}

	schedulingGroupsNameSplitted := strings.Split(schedulingGroups, ",")
	var schedulingGroupsArray []squarescale.SchedulingGroup

	for _, name := range schedulingGroupsNameSplitted {
		schedulingGroup, err := client.GetSchedulingGroupInfo(UUID, name)
		if err != nil {
			fmt.Println(err)
			continue
		}
		schedulingGroupsArray = append(schedulingGroupsArray, squarescale.SchedulingGroup{
			ID:   schedulingGroup.ID,
			Name: name,
		})
	}

	return schedulingGroupsArray
}
