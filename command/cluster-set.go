package command

import (
	"errors"
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
	projectUUID := c.flagSet.String("project-uuid", "", "set the uuid of the project")
	c.flagSet.UintVar(&c.Cluster.Size, "size", 0, "Cluster Size")
	alwaysYes := yesFlag(c.flagSet)
	if err := c.flagSet.Parse(args); err != nil {
		return 1
	}

	if c.flagSet.NArg() > 0 {
		return c.errorWithUsage(fmt.Errorf("Unparsed arguments on the command line: %v", c.flagSet.Args()))
	}

	if *projectUUID == "" {
		return c.errorWithUsage(errors.New("Project uuid is mandatory"))
	}

	c.Ui.Warn(fmt.Sprintf("Changing cluster settings for project '%s' may cause a downtime.", *projectUUID))
	ok, err := AskYesNo(c.Ui, alwaysYes, "Is this ok?", false)
	if err != nil {
		return c.error(err)
	} else if !ok {
		return c.cancelled()
	}

	var taskId int

	res := c.runWithSpinner("scale project cluster", endpoint.String(), func(client *squarescale.Client) (string, error) {
		var err error
		cluster, e := client.GetClusterConfig(*projectUUID)
		if e != nil {
			return "", e
		}

		if c.Cluster.Size == cluster.Size {
			*nowait = true
			return fmt.Sprintf("Cluster for project '%s' is already configured with these parameters", *projectUUID), nil
		}

		cluster.Update(c.Cluster)

		taskId, err = client.ConfigCluster(*projectUUID, cluster)

		msg := fmt.Sprintf(
			"[#%d] Successfully set cluster for project '%s': %s",
			taskId, *projectUUID, cluster.String())

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
	return "Set and scale up/down cluster attached to project"
}

// Help is part of cli.Command implementation.
func (c *ClusterSetCommand) Help() string {
	helpText := `
usage: sqsc cluster set [options]

  Set and scale up/down cluster attached to project.
`
	return strings.TrimSpace(helpText + optionsFromFlags(c.flagSet))
}

func validateClusterSetCommandArgs(project string, cluster squarescale.ClusterConfig) error {
	if err := validateProjectName(project); err != nil {
		return err
	}

	if cluster.Size == 0 {
		return errors.New("Size cannot be empty or 0.")
	}

	return nil
}
