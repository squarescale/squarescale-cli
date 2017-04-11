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
	preSet := c.flagSet.Bool("preset", false, "Print pre-set variables")
	global := c.flagSet.Bool("global", false, "Print global variables")

	if err := c.flagSet.Parse(args); err != nil {
		return 1
	}

	noFlag := !*global && *container == "" && !*preSet
	printGlobal := *global || noFlag
	printPreset := *preSet || noFlag

	if c.flagSet.NArg() > 0 {
		return c.errorWithUsage(fmt.Errorf("Unparsed arguments on the command line: %v", c.flagSet.Args()))
	}

	if err := validateProjectName(*project); err != nil {
		return c.errorWithUsage(err)
	}

	return c.runWithSpinner("list environment variables", *endpoint, func(client *squarescale.Client) (string, error) {
		env, err := client.EnvironmentVariables(*project)
		if err != nil {
			return "", err
		}

		var lines []string

		if printPreset {
			printPresets(&lines, env)
		}

		if printGlobal || *container != "" || noFlag {
			lines = append(lines, "CUSTOM VARIABLES")
		}

		if printGlobal {
			printGlobals(&lines, env)
		}

		if noFlag || *container != "" {
			printPerService(&lines, env, container)
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

func printPresets(lines *[]string, env *squarescale.Environment) {
	*lines = append(*lines, "PRESET VARIABLES")
	printVars(lines, &env.Preset)
	*lines = append(*lines, "")
}

func printGlobals(lines *[]string, env *squarescale.Environment) {
	*lines = append(*lines, "|- GLOBAL")
	printVars(lines, &env.Custom.Global)
}

func printPerService(lines *[]string, env *squarescale.Environment, container *string) {
	for serviceName, vars := range env.Custom.PerService {
		if serviceName == *container || *container == "" {
			*lines = append(*lines, fmt.Sprintf("|- %s", serviceName))
			printVars(lines, &vars)
		}
	}
}

func printVars(lines *[]string, vars *map[string]string) {
	for k, v := range *vars {
		*lines = append(*lines, fmt.Sprintf("|-- %s=\"%s\"", k, v))
	}
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
