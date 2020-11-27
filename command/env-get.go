package command

import (
	"errors"
	"flag"
	"fmt"
	"strings"

	"github.com/squarescale/squarescale-cli/squarescale"
)

// EnvGetCommand gets the URL of the load balancer associated to a projects and prints it on the standard output.
type EnvGetCommand struct {
	Meta
	flagSet *flag.FlagSet
}

// Run is part of cli.Command implementation.
func (c *EnvGetCommand) Run(args []string) int {
	c.flagSet = newFlagSet(c, c.Ui)
	endpoint := endpointFlag(c.flagSet)
	projectUUID := c.flagSet.String("project-uuid", "", "set the uuid of the project")
	projectName := c.flagSet.String("project-name", "", "set the name of the project")
	container := containerFlag(c.flagSet)
	all := c.flagSet.Bool("all", false, "Print all variables")

	if err := c.flagSet.Parse(args); err != nil {
		return 1
	}

	if c.flagSet.NArg() > 1 {
		return c.errorWithUsage(fmt.Errorf("Unparsed arguments on the command line: %v", c.flagSet.Args()))
	}

	if *projectUUID == "" && *projectName == "" {
		return c.errorWithUsage(errors.New("Project name or uuid is mandatory"))
	}

	queryVar := c.flagSet.Arg(0)

	return c.runWithSpinner("list environment variables", endpoint.String(), func(client *squarescale.Client) (string, error) {
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

		env, err := squarescale.NewEnvironment(client, UUID)
		if err != nil {
			return "", err
		}

		queryOptions := squarescale.QueryOptions{
			DisplayAll:   *all,
			ServiceName:  *container,
			VariableName: queryVar,
		}
		queryResult, err := env.QueryVars(queryOptions)
		if err != nil {
			return "", err
		}

		return queryResult.String(), nil
	})
}

// Synopsis is part of cli.Command implementation.
func (c *EnvGetCommand) Synopsis() string {
	return "Display environment variables of project"
}

// Help is part of cli.Command implementation.
func (c *EnvGetCommand) Help() string {
	helpText := `
usage: sqsc env get [options] [VARNAME]

  Display environment variables of project.
`
	return strings.TrimSpace(helpText + optionsFromFlags(c.flagSet))
}
