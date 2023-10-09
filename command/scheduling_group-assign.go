package command

import (
	"errors"
	"flag"
	"fmt"
	"strings"

	"github.com/squarescale/squarescale-cli/squarescale"
)

// SchedulingGroupAssignCommand is a cli.Command implementation for assigning hosts to scheduling group.
type SchedulingGroupAssignCommand struct {
	Meta
	flagSet *flag.FlagSet
}

func (b *SchedulingGroupAssignCommand) Run(args []string) int {
	b.flagSet = newFlagSet(b, b.Ui)
	endpoint := endpointFlag(b.flagSet)
	projectUUID := projectUUIDFlag(b.flagSet)
	projectName := projectNameFlag(b.flagSet)

	if err := b.flagSet.Parse(args); err != nil {
		return 1
	}

	if *projectUUID == "" && *projectName == "" {
		return b.errorWithUsage(errors.New("Project name or uuid is mandatory"))
	}

	schedulingGroupName, err := schedulingGroupNameArg(b.flagSet, 0)
	if err != nil {
		return b.errorWithUsage(err)
	}

	schedulingGroupHosts, err := schedulingGroupHostsArg(b.flagSet, 1)
	if err != nil {
		return b.errorWithUsage(err)
	}

	if b.flagSet.NArg() > 2 {
		return b.errorWithUsage(fmt.Errorf("Unparsed arguments on the command line: %v", b.flagSet.Args()))
	}

	return b.runWithSpinner("assign host(s) to scheduling group", endpoint.String(), func(client *squarescale.Client) (string, error) {
		var UUID string
		var err error
		if *projectUUID == "" {
			UUID, err = client.ProjectByName(*projectName)
			if err != nil {
				return "", err
			}
		} else {
			UUID = *projectUUID
		}
		_, err = client.GetSchedulingGroupInfo(UUID, schedulingGroupName)
		if err != nil {
			return "", err
		}

		projectDetails, err := client.GetProjectDetails(UUID)
		if err != nil {
			return "", err
		}
		if projectDetails == nil {
			return "", fmt.Errorf("No project details found")
		}
		var ids []int
		hosts := strings.Split(schedulingGroupHosts, ",")
		for _, h := range hosts {
			found := false
			for _, c := range projectDetails.Project.Infrastructure.Cluster.ClusterMembersDetails {
				if c.Name == h {
					ids = append(ids, c.ID)
					found = true
					break
				}
			}
			if !found {
				for _, c := range projectDetails.Project.Infrastructure.Cluster.ExternalNodes {
					if c.Name == h {
						ids = append(ids, c.ID)
						found = true
						break
					}
				}
				if !found {
					return "", fmt.Errorf("Can not find host %s", h)
				}
			}
		}
		err = client.PutSchedulingGroupsNodes(UUID, schedulingGroupName, "cluster_members_to_add", ids)
		if err != nil {
			return "", err
		}
		return "Success", nil
	})
}

// Synopsis is part of cli.Command implementation.
func (b *SchedulingGroupAssignCommand) Synopsis() string {
	return "Assign host(s) to scheduling group of project"
}

// Help is part of cli.Command implementation.
func (b *SchedulingGroupAssignCommand) Help() string {
	helpText := `
usage: sqsc scheduling-group assign [options] <scheduling_group_name> <host1[,host2...,[hostn]>

  Assign host(s) to scheduling group of project.
`
	return strings.TrimSpace(helpText + optionsFromFlags(b.flagSet))
}
