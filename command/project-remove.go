package command

import (
	"errors"
	"flag"
	"fmt"
	"strings"
	"time"

	"github.com/squarescale/squarescale-cli/squarescale"
)

// ProjectRemoveCommand is a cli.Command implementation for creating project.
type ProjectRemoveCommand struct {
	Meta
	flagSet *flag.FlagSet
}

// Run is part of cli.Command implementation.
func (cmd *ProjectRemoveCommand) Run(args []string) int {
	// Parse flags
	cmd.flagSet = newFlagSet(cmd, cmd.Ui)
	alwaysYes := yesFlag(cmd.flagSet)
	endpoint := endpointFlag(cmd.flagSet)
	nowait := nowaitFlag(cmd.flagSet)
	projectUUID := projectUUIDFlag(cmd.flagSet)
	projectName := projectNameFlag(cmd.flagSet)
	if err := cmd.flagSet.Parse(args); err != nil {
		return 1
	}

	if *projectUUID == "" && *projectName == "" {
		return cmd.errorWithUsage(errors.New("Project name or uuid is mandatory"))
	}

	if cmd.flagSet.NArg() > 1 {
		return cmd.errorWithUsage(fmt.Errorf("Unparsed arguments on the command line: %v", cmd.flagSet.Args()[1:]))
	}

	var UUID string
	var err error
	var projectToShow string
	if *projectUUID == "" {
		projectToShow = *projectName
	} else {
		projectToShow = *projectUUID
	}

	cmd.Ui.Info("Destroy infrastructure and configuration for project " + projectToShow + "?")
	ok, err := AskYesNo(cmd.Ui, alwaysYes, "Proceed ?", true)
	if err != nil {
		return cmd.error(err)
	} else if !ok {
		return cmd.error(CancelledError)
	}

	res := cmd.runWithSpinner("unprovision project", endpoint.String(), func(client *squarescale.Client) (string, error) {
		var err error
		if *projectUUID == "" {
			UUID, err = client.ProjectByName(*projectName)
			if err != nil {
				return "", err
			}
		} else {
			UUID = *projectUUID
		}

		err = client.ProjectUnprovision(UUID)
		if err != nil {
			return "", err
		}

		var infraStatus string = "provisionning"

	loop:
		for {
			status, err := client.GetProject(UUID)
			if err != nil {
				return "", err
			}
			if status.InfraStatus == "error" {
				return "", fmt.Errorf("Unknown infrastructure error.")
			} else if status.InfraStatus == "no_infra" {
				break loop
			} else if infraStatus != status.InfraStatus {
				infraStatus = status.InfraStatus
				cmd.Ui.Info("Infrastructure status: " + infraStatus)
			}

			time.Sleep(time.Second)
		}

		return "", nil
	})
	if res != 0 {
		return res
	}

	res = cmd.runWithSpinner("delete project", endpoint.String(), func(client *squarescale.Client) (string, error) {
		err := client.ProjectDelete(UUID)
		return fmt.Sprintf("Remove project '%s'", UUID), err
	})

	if res != 0 {
		return res
	}

	if !*nowait {
		res = cmd.runWithSpinner("wait for project remove", endpoint.String(), func(client *squarescale.Client) (string, error) {
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
func (cmd *ProjectRemoveCommand) Synopsis() string {
	return "Remove project"
}

// Help is part of cli.Command implementation.
func (cmd *ProjectRemoveCommand) Help() string {
	helpText := `
usage: sqsc project remove [options] <project_name>

  Destroy infrastructure and remove project completely.
`
	return strings.TrimSpace(helpText + optionsFromFlags(cmd.flagSet))
}
