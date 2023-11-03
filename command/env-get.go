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
func (cmd *EnvGetCommand) Run(args []string) int {
	cmd.flagSet = newFlagSet(cmd, cmd.Ui)
	endpoint := endpointFlag(cmd.flagSet)
	projectUUID := projectUUIDFlag(cmd.flagSet)
	projectName := projectNameFlag(cmd.flagSet)
	container := serviceFlag(cmd.flagSet)
	all := envAllFlag(cmd.flagSet)

	if err := cmd.flagSet.Parse(args); err != nil {
		return 1
	}

	if cmd.flagSet.NArg() > 1 {
		return cmd.errorWithUsage(fmt.Errorf("Unparsed arguments on the command line: %v", cmd.flagSet.Args()))
	}

	if *projectUUID == "" && *projectName == "" {
		return cmd.errorWithUsage(errors.New("Project name or uuid is mandatory"))
	}

	queryVar := cmd.flagSet.Arg(0)

	return cmd.runWithSpinner("list environment variables", endpoint.String(), func(client *squarescale.Client) (string, error) {
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
func (cmd *EnvGetCommand) Synopsis() string {
	return "Display environment variables of project"
}

// Help is part of cli.Command implementation.
func (cmd *EnvGetCommand) Help() string {
	helpText := `
usage: sqsc env get [options] [VARNAME]

  Display environment variables of project.
`
	return strings.TrimSpace(helpText + optionsFromFlags(cmd.flagSet))
}
