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

func (c *ProjectSettingsCommand) Run(args []string) int {
	c.flagSet = newFlagSet(c, c.Ui)
	endpoint := endpointFlag(c.flagSet)
	projectUUID := projectUUIDFlag(c.flagSet)
	projectName := projectNameFlag(c.flagSet)
	hybridCluster := projectHybridClusterFlag(c.flagSet)

	if err := c.flagSet.Parse(args); err != nil {
		return 1
	}

	if *projectUUID == "" && *projectName == "" {
		return c.errorWithUsage(errors.New("Project name or uuid is mandatory"))
	}

	if c.flagSet.NArg() > 0 {
		return c.errorWithUsage(fmt.Errorf("Unparsed arguments on the command line: %v", c.flagSet.Args()))
	}

	return c.runWithSpinner("change project settings", endpoint.String(), func(client *squarescale.Client) (string, error) {
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
func (c *ProjectSettingsCommand) Synopsis() string {
	return "Change settings of a project."
}

// Help is part of cli.Command implementation.
func (c *ProjectSettingsCommand) Help() string {
	helpText := `
usage: sqsc settings [options]

  Change settings of a project.
`
	return strings.TrimSpace(helpText + optionsFromFlags(c.flagSet))
}
