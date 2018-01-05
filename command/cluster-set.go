package command

import (
	"flag"
	"fmt"
	"strings"

	"github.com/squarescale/squarescale-cli/squarescale"
)

// ClusterSetCommand is a cli.Command implementation for changing the cluster of
// a project.
type ClusterSetCommand struct {
	Meta
	flagSet *flag.FlagSet
	Cluster squarescale.ClusterConfig
}

// Run is part of cli.Command implementation.
func (c *ClusterSetCommand) Run(args []string) int {
	c.flagSet = newFlagSet(c, c.Ui)
	endpoint := endpointFlag(c.flagSet)
	nowait := nowaitFlag(c.flagSet)
	projectNameArg := projectFlag(c.flagSet)
	c.flagSet.UintVar(&c.Cluster.Size, "size", 0, "Cluster Size")
	alwaysYes := yesFlag(c.flagSet)
	if err := c.flagSet.Parse(args); err != nil {
		return 1
	}

	if c.Cluster.Size == 0 {
		return c.errorWithUsage(fmt.Errorf("You must specify a cluster size"))
	}

	if c.flagSet.NArg() > 0 {
		return c.errorWithUsage(fmt.Errorf("Unparsed arguments on the command line: %v", c.flagSet.Args()))
	}

	c.Ui.Warn(fmt.Sprintf("Changing cluster settings for project '%s' may cause a downtime.", *projectNameArg))
	ok, err := AskYesNo(c.Ui, alwaysYes, "Is this ok?", false)
	if err != nil {
		return c.error(err)
	} else if !ok {
		return c.cancelled()
	}

	var taskId int

	res := c.runWithSpinner("scale project cluster", endpoint.String(), func(client *squarescale.Client) (string, error) {
		var err error
		cluster, e := client.GetClusterConfig(*projectNameArg)
		if e != nil {
			return "", e
		}

		if c.Cluster.Size == cluster.Size {
			*nowait = true
			return fmt.Sprintf("Cluster for project '%s' is already configured with these parameters", *projectNameArg), nil
		}

		cluster.Update(c.Cluster)

		taskId, err = client.ConfigCluster(*projectNameArg, cluster)

		msg := fmt.Sprintf(
			"[#%d] Successfully set cluster for project '%s': %s",
			taskId, *projectNameArg, cluster.String())

		return msg, err
	})
	if res != 0 {
		return res
	}

	if !*nowait {
		res = c.runWithSpinner("wait for cluster change", endpoint.String(), func(client *squarescale.Client) (string, error) {
			task, err := client.WaitTask(taskId)
			if err != nil {
				return "", err
			} else {
				return task.Status, nil
			}
		})
	}

	return res
}

// Synopsis is part of cli.Command implementation.
func (c *ClusterSetCommand) Synopsis() string {
	return "Set and scale up/down the cluster attached to a Squarescale project"
}

// Help is part of cli.Command implementation.
func (c *ClusterSetCommand) Help() string {
	helpText := `
usage: sqsc cluster set [options]

  Set and scale up/down the cluster attached to a Squarescale project.

`
	return strings.TrimSpace(helpText + optionsFromFlags(c.flagSet))
}
