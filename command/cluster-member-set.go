package command

import (
	"errors"
	"flag"
	"fmt"
	"strings"

	"github.com/squarescale/squarescale-cli/squarescale"
)

type ClusterMemberSetCommand struct {
	Meta
	flagSet *flag.FlagSet
}

// Run is part of cli.Command implementation.
func (c *ClusterMemberSetCommand) Run(args []string) int {
	c.flagSet = newFlagSet(c, c.Ui)
	endpoint := endpointFlag(c.flagSet)
	projectUUID := c.flagSet.String("project-uuid", "", "set the uuid of the project")
	projectName := c.flagSet.String("project-name", "", "set the name of the project")
	nameArg := c.flagSet.String("name", "", "The node name.")
	schedulingGroupArg := c.flagSet.String("scheduling-group", "", "The scheduling group to attach to the node.")

	if err := c.flagSet.Parse(args); err != nil {
		return 1
	}

	if c.flagSet.NArg() > 0 {
		return c.errorWithUsage(fmt.Errorf("Unparsed arguments on the command line: %v", c.flagSet.Args()))
	}

	if *nameArg == "" {
		return c.errorWithUsage(errors.New("Node name is mandatory."))
	}

	return c.runWithSpinner("configure node settings", endpoint.String(), func(client *squarescale.Client) (string, error) {
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

		clusterMember, err := client.GetClusterMemberInfo(UUID, *nameArg)
		if err != nil {
			return "", err
		}

		schedulingGroup, err := client.GetSchedulingGroupInfo(UUID, *schedulingGroupArg)
		if err != nil {
			return "", err
		}

		newClusterMember := squarescale.ClusterMember{
			SchedulingGroup: schedulingGroup,
		}

		msg := fmt.Sprintf("Successfully configured cluster node '%s'", clusterMember.Name)
		err = client.ConfigClusterMember(UUID, clusterMember, newClusterMember)

		return msg, err
	})
}

// Synopsis is part of cli.Command implementation.
func (c *ClusterMemberSetCommand) Synopsis() string {
	return "Set cluster node settings attached to a project cluster."
}

// Help is part of cli.Command implementation.
func (c *ClusterMemberSetCommand) Help() string {
	helpText := `
usage: sqsc cluster node set [options]

  Set cluster node settings attached to a project cluster.
`
	return strings.TrimSpace(helpText + optionsFromFlags(c.flagSet))
}
