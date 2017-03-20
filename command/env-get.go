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
	if err := c.flagSet.Parse(args); err != nil {
		return 1
	}

	if c.flagSet.NArg() > 0 {
		return c.errorWithUsage(fmt.Errorf("Unparsed arguments on the command line: %v", c.flagSet.Args()))
	}

	if err := validateProjectName(*project); err != nil {
		return c.errorWithUsage(err)
	}

	return c.runWithSpinner("list environment variables", *endpoint, func(client *squarescale.Client) (string, error) {
		vars, err := client.EnvironmentVariables(*project)
		if err != nil {
			return "", err
		}

		var lines []string
		for k, v := range vars {
			lines = append(lines, fmt.Sprintf("%s=\"%s\"", k, v))
		}

		var msg string
		if len(lines) > 0 {
			msg = strings.Join(lines, "\n")
		} else {
			msg = fmt.Sprintf("No environment variables found for project '%s'", *project)
		}

		return msg, nil
	})
}

// Synopsis is part of cli.Command implementation.
func (c *EnvGetCommand) Synopsis() string {
	return "Display project's environment variables"
}

// Help is part of cli.Command implementation.
func (c *EnvGetCommand) Help() string {
	helpText := `
usage: sqsc env get [options]

  Display project environment variables.

`
	return strings.TrimSpace(helpText + optionsFromFlags(c.flagSet))
}
