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
func (cmd *ClusterSetCommand) Run(args []string) int {
	cmd.flagSet = newFlagSet(cmd, cmd.Ui)
	endpoint := endpointFlag(cmd.flagSet)
	nowait := nowaitFlag(cmd.flagSet)
	projectUUID := projectUUIDFlag(cmd.flagSet)
	projectName := projectNameFlag(cmd.flagSet)
	cmd.Cluster.Size = *clusterSizeFlag(cmd.flagSet)
	alwaysYes := yesFlag(cmd.flagSet)
	if err := cmd.flagSet.Parse(args); err != nil {
		return 1
	}

	if cmd.flagSet.NArg() > 0 {
		return cmd.errorWithUsage(fmt.Errorf("Unparsed arguments on the command line: %v", cmd.flagSet.Args()))
	}

	if *projectUUID == "" && *projectName == "" {
		return cmd.errorWithUsage(errors.New("Project name or uuid is mandatory"))
	}

	var projectToShow string
	if *projectUUID == "" {
		projectToShow = *projectName
	} else {
		projectToShow = *projectUUID
	}

	cmd.Ui.Warn(fmt.Sprintf("Changing cluster settings for project '%s' may cause a downtime.", projectToShow))
	ok, err := AskYesNo(cmd.Ui, alwaysYes, "Is this ok?", false)
	if err != nil {
		return cmd.error(err)
	} else if !ok {
		return cmd.cancelled()
	}

	var UUID string

	res := cmd.runWithSpinner("scale project cluster", endpoint.String(), func(client *squarescale.Client) (string, error) {
		var err error
		if *projectUUID == "" {
			UUID, err = client.ProjectByName(*projectName)
			if err != nil {
				return "", err
			}
		} else {
			UUID = *projectUUID
		}

		cluster, e := client.GetClusterConfig(UUID)
		if e != nil {
			return "", e
		}

		if cmd.Cluster.Size == cluster.Size {
			*nowait = true
			return fmt.Sprintf("Cluster for project '%s' is already configured with these parameters", projectToShow), nil
		}

		cluster.Update(cmd.Cluster)

		_, err = client.ConfigCluster(UUID, cluster)

		msg := fmt.Sprintf(
			"Successfully set cluster for project '%s': %s",
			projectToShow, cluster.String())

		return msg, err
	})
	if res != 0 {
		return res
	}

	if !*nowait {
		res = cmd.runWithSpinner("wait for cluster change", endpoint.String(), func(client *squarescale.Client) (string, error) {
			projectStatus, err := client.WaitProject(UUID, 5)
			if err != nil {
				return projectStatus, err
			} else {
				return projectStatus, nil
			}
		})
	}

	return res
}

// Synopsis is part of cli.Command implementation.
func (cmd *ClusterSetCommand) Synopsis() string {
	return "Set and scale up/down cluster attached to project"
}

// Help is part of cli.Command implementation.
func (cmd *ClusterSetCommand) Help() string {
	helpText := `
usage: sqsc cluster set [options]

  Set and scale up/down cluster attached to project.
`
	return strings.TrimSpace(helpText + optionsFromFlags(cmd.flagSet))
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
