package command

import (
	"flag"
	"fmt"
	"strings"
	"time"

	"github.com/squarescale/squarescale-cli/squarescale"
)

// ClusterSetCommand is a cli.Command implementation for changing the cluster of
// a project.
type ClusterSetCommand struct {
	Meta
	flagSet *flag.FlagSet
}

// Run is part of cli.Command implementation.
func (c *ClusterSetCommand) Run(args []string) int {
	c.flagSet = newFlagSet(c, c.Ui)
	endpoint := endpointFlag(c.flagSet)
	nowait := nowaitFlag(c.flagSet)
	projectNameArg := projectFlag(c.flagSet)
	clusterSizeArg := clusterSizeFlag(c.flagSet)
	alwaysYes := yesFlag(c.flagSet)
	if err := c.flagSet.Parse(args); err != nil {
		return 1
	}

	if *clusterSizeArg == 0 {
		return c.errorWithUsage(fmt.Errorf("You must specify a cluster size"))
	}

	c.Ui.Warn(fmt.Sprintf("Changing cluster settings for project '%s' may cause a downtime.", *projectNameArg))
	ok, err := AskYesNo(c.Ui, alwaysYes, "Is this ok?", false)
	if err != nil {
		return c.error(err)
	} else if !ok {
		return c.cancelled()
	}

	var taskId int

	res := c.runWithSpinner("change cluster size", *endpoint, func(client *squarescale.Client) (string, error) {
		var err error
		currentSize, e := client.GetClusterSize(*projectNameArg)
		if e != nil {
			return "", e
		}

		if currentSize == *clusterSizeArg {
			return fmt.Sprintf("cluster is already size %d", currentSize), nil
		}

		taskId, err = client.SetClusterSize(*projectNameArg, *clusterSizeArg)
		return fmt.Sprintf("[#%d] changed cluster size from %d to %d", taskId, currentSize, *clusterSizeArg), err
	})
	if res != 0 {
		return res
	}

	if !*nowait {
		res = c.runWithSpinner("wait for cluster change", *endpoint, func(client *squarescale.Client) (string, error) {
			return client.WaitTask(taskId, time.Second)
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
