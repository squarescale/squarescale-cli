package command

import (
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
	project := projectFlag(c.flagSet)
	container := containerFlag(c.flagSet)
	all := c.flagSet.Bool("all", false, "Print all variables")

	if err := c.flagSet.Parse(args); err != nil {
		return 1
	}

	if c.flagSet.NArg() > 1 {
		return c.errorWithUsage(fmt.Errorf("Unparsed arguments on the command line: %v", c.flagSet.Args()))
	}

	if err := validateProjectName(*project); err != nil {
		return c.errorWithUsage(err)
	}

	queryVar := c.flagSet.Arg(0)

	return c.runWithSpinner("list environment variables", endpoint.String(), func(client *squarescale.Client) (string, error) {
		env, err := squarescale.NewEnvironment(client, *project)
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
	return "Display project's environment variables"
}

// Help is part of cli.Command implementation.
func (c *EnvGetCommand) Help() string {
	helpText := `
usage: sqsc env get [options] [VARNAME]

  Display project environment variables.

`
	return strings.TrimSpace(helpText + optionsFromFlags(c.flagSet))
}
