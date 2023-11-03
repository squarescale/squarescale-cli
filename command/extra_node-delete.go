package command

import (
	"errors"
	"flag"
	"fmt"
	"strings"

	"github.com/squarescale/squarescale-cli/squarescale"
)

// ExtraNodeDeleteCommand is a cli.Command implementation for deleting a extra-node.
type ExtraNodeDeleteCommand struct {
	Meta
	flagSet *flag.FlagSet
}

// Run is part of cli.Command implementation.
func (cmd *ExtraNodeDeleteCommand) Run(args []string) int {
	// Parse flags
	cmd.flagSet = newFlagSet(cmd, cmd.Ui)
	alwaysYes := yesFlag(cmd.flagSet)
	endpoint := endpointFlag(cmd.flagSet)
	projectUUID := projectUUIDFlag(cmd.flagSet)
	projectName := projectNameFlag(cmd.flagSet)
	nowait := nowaitFlag(cmd.flagSet)

	if err := cmd.flagSet.Parse(args); err != nil {
		return 1
	}

	extraNodeName, err := extraNodeNameArg(cmd.flagSet, 0)
	if err != nil {
		return cmd.errorWithUsage(err)
	}

	if cmd.flagSet.NArg() > 4 {
		return cmd.errorWithUsage(fmt.Errorf("Unparsed arguments on the command line: %v", cmd.flagSet.Args()[1:]))
	}

	if *projectUUID == "" && *projectName == "" {
		return cmd.errorWithUsage(errors.New("Project name or uuid is mandatory"))
	}

	cmd.Ui.Info("Are you sure you want to delete " + extraNodeName + "?")
	if *alwaysYes {
		cmd.Ui.Info("(approved from command line)")
	} else {
		res, err := cmd.Ui.Ask("y/N")
		if err != nil {
			return cmd.error(err)
		} else if res != "Y" && res != "y" {
			return cmd.cancelled()
		}
	}

	var UUID string

	res := cmd.runWithSpinner("deleting extra-node", endpoint.String(), func(client *squarescale.Client) (string, error) {
		var err error
		var projectToShow string
		if *projectUUID == "" {
			projectToShow = *projectName
			UUID, err = client.ProjectByName(*projectName)
			if err != nil {
				return "", err
			}
		} else {
			projectToShow = *projectUUID
			UUID = *projectUUID
		}

		fmt.Printf("Delete on project `%s` the extra-node `%s`\n", projectToShow, extraNodeName)
		err = client.DeleteExtraNode(UUID, extraNodeName)
		return "", err
	})
	if res != 0 {
		return res
	}

	if !*nowait {
		cmd.runWithSpinner("wait for extra-node delete", endpoint.String(), func(client *squarescale.Client) (string, error) {
			_, err := client.WaitProject(UUID, 5)
			if err != nil {
				return "", err
			} else {
				return "", nil
			}
		})
	}

	return 0
}

// Synopsis is part of cli.Command implementation.
func (cmd *ExtraNodeDeleteCommand) Synopsis() string {
	return "Delete extra-node from project."
}

// Help is part of cli.Command implementation.
func (cmd *ExtraNodeDeleteCommand) Help() string {
	helpText := `
usage: sqsc extra-node delete [options] <extra_node_name>

  Delete extra-node from project.
`
	return strings.TrimSpace(helpText + optionsFromFlags(cmd.flagSet))
}
