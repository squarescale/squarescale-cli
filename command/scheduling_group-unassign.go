package command

import (
	"errors"
	"flag"
	"fmt"
	"strings"

	"github.com/squarescale/squarescale-cli/squarescale"
)

// SchedulingGroupUnAssignCommand is a cli.Command implementation for assigning hosts to scheduling group.
type SchedulingGroupUnAssignCommand struct {
	Meta
	flagSet *flag.FlagSet
}

func (cmd *SchedulingGroupUnAssignCommand) Run(args []string) int {
	cmd.flagSet = newFlagSet(cmd, cmd.Ui)
	endpoint := endpointFlag(cmd.flagSet)
	projectUUID := projectUUIDFlag(cmd.flagSet)
	projectName := projectNameFlag(cmd.flagSet)

	if err := cmd.flagSet.Parse(args); err != nil {
		return 1
	}

	if *projectUUID == "" && *projectName == "" {
		return cmd.errorWithUsage(errors.New("Project name or uuid is mandatory"))
	}

	schedulingGroupName, err := schedulingGroupNameArg(cmd.flagSet, 0)
	if err != nil {
		return cmd.errorWithUsage(err)
	}

	schedulingGroupHosts, err := schedulingGroupHostsArg(cmd.flagSet, 1)
	if err != nil {
		return cmd.errorWithUsage(err)
	}

	if cmd.flagSet.NArg() > 2 {
		return cmd.errorWithUsage(fmt.Errorf("Unparsed arguments on the command line: %v", cmd.flagSet.Args()))
	}

	return cmd.runWithSpinner("unassign host(s) from scheduling group", endpoint.String(), func(client *squarescale.Client) (string, error) {
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
		err = client.PutSchedulingGroupsNodes(UUID, schedulingGroupName, "cluster_member_to_remove", ids)
		if err != nil {
			return "", err
		}
		return "Success", nil
	})
}

// Synopsis is part of cli.Command implementation.
func (cmd *SchedulingGroupUnAssignCommand) Synopsis() string {
	return "Unassign host(s) from scheduling group of project"
}

// Help is part of cli.Command implementation.
func (cmd *SchedulingGroupUnAssignCommand) Help() string {
	helpText := `
usage: sqsc scheduling-group unassign [options] <scheduling_group_name> <host1[,host2...,[hostn]>

  Unassign host(s) from scheduling group of project.
`
	return strings.TrimSpace(helpText + optionsFromFlags(cmd.flagSet))
}
