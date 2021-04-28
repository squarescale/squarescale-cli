package command

import (
	"errors"
	"flag"
	"fmt"
	"strconv"
	"strings"

	"github.com/squarescale/squarescale-cli/squarescale"
)

type ProjectSettingsCommand struct {
	Meta
	flagSet *flag.FlagSet
}

func (p *ProjectSettingsCommand) Run(args []string) int {
	p.flagSet = newFlagSet(p, p.Ui)
	endpoint := endpointFlag(p.flagSet)
	projectUUID := p.flagSet.String("project-uuid", "", "set the uuid of the project")
	projectName := p.flagSet.String("project-name", "", "set the name of the project")
	hybridClusterArg := p.flagSet.String("hybrid-cluster", "", "Enable/Disable hybrid cluster")

	if err := p.flagSet.Parse(args); err != nil {
		return 1
	}

	if *projectUUID == "" && *projectName == "" {
		return p.errorWithUsage(errors.New("Project name or uuid is mandatory"))
	}

	if *hybridClusterArg == "" {
		return p.errorWithUsage(errors.New("Hybrid cluster is mandatory"))
	}

	if p.flagSet.NArg() > 0 {
		return p.errorWithUsage(fmt.Errorf("Unparsed arguments on the command line: %v", p.flagSet.Args()))
	}

	return p.runWithSpinner("change project settings", endpoint.String(), func(client *squarescale.Client) (string, error) {
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

		hybridCluster, parseErr := strconv.ParseBool(*hybridClusterArg)
		if parseErr != nil {
			return "", err
		}

		err = client.ConfigProjectSettings(UUID, squarescale.Project{
			HybridClusterEnabled: hybridCluster,
		})

		return msg, err
	})
}

// Synopsis is part of cli.Command implementation.
func (p *ProjectSettingsCommand) Synopsis() string {
	return "Change settings of a project."
}

// Help is part of cli.Command implementation.
func (p *ProjectSettingsCommand) Help() string {
	helpText := `
usage: sqsc settings [options]

  Change settings of a project.
`
	return strings.TrimSpace(helpText + optionsFromFlags(p.flagSet))
}
