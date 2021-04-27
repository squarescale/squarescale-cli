package command

import (
	"fmt"
	"strings"

	"github.com/squarescale/squarescale-cli/squarescale"
)

func parseSchedulingGroupsToAdd(UUID string, client *squarescale.Client, schedulingGroups string) []int {
	schedulingGroupsNameSplitted := strings.Split(schedulingGroups, ",")

	var schedulingGroupsIds []int

	for _, name := range schedulingGroupsNameSplitted {
		schedulingGroup, err := client.GetSchedulingGroupInfo(UUID, name)
		if err != nil {
			fmt.Println(err)
		}
		schedulingGroupsIds = append(schedulingGroupsIds, schedulingGroup.ID)
	}

	return schedulingGroupsIds
}
