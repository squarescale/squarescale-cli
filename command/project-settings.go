package command

import (
	"errors"
	"flag"
	"fmt"
	"strings"

	"github.com/squarescale/squarescale-cli/squarescale"
)

type ProjectSettingsCommand struct {
	Meta
	flagSet *flag.FlagSet
}

func (cmd *ProjectSettingsCommand) Run(args []string) int {
	cmd.flagSet = newFlagSet(cmd, cmd.Ui)
	endpoint := endpointFlag(cmd.flagSet)
	projectUUID := projectUUIDFlag(cmd.flagSet)
	projectName := projectNameFlag(cmd.flagSet)
	hybridCluster := projectHybridClusterFlag(cmd.flagSet)

	if err := cmd.flagSet.Parse(args); err != nil {
		return 1
	}

	if *projectUUID == "" && *projectName == "" {
		return cmd.errorWithUsage(errors.New("Project name or uuid is mandatory"))
	}

	if cmd.flagSet.NArg() > 0 {
		return cmd.errorWithUsage(fmt.Errorf("Unparsed arguments on the command line: %v", cmd.flagSet.Args()))
	}

	return cmd.runWithSpinner("change project settings", endpoint.String(), func(client *squarescale.Client) (string, error) {
		var err error
		var UUID string
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

		msg := fmt.Sprintf("Successfully change project settings for project '%s'", projectToShow)

		err = client.ConfigProjectSettings(UUID, squarescale.Project{
			HybridClusterEnabled: *hybridCluster,
		})

		return msg, err
	})
}

// Synopsis is part of cli.Command implementation.
func (cmd *ProjectSettingsCommand) Synopsis() string {
	return "Change settings of a project."
}

// Help is part of cli.Command implementation.
func (cmd *ProjectSettingsCommand) Help() string {
	helpText := `
usage: sqsc settings [options]

  Change settings of a project.
`
	return strings.TrimSpace(helpText + optionsFromFlags(cmd.flagSet))
}
