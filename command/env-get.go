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
	noPP := c.flagSet.Bool("no-pretty-print", false, "Format to allow easier parsing")

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
			printPresets(&lines, env, *noPP)
		}

		if printGlobal || *container != "" || noFlag && !*noPP {
			lines = append(lines, "CUSTOM VARIABLES")
		}

		if printGlobal {
			printGlobals(&lines, env, *noPP)
		}

		if noFlag || *container != "" {
			printPerService(&lines, env, container, *noPP)
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

func printPresets(lines *[]string, env *squarescale.Environment, noPP bool) {
	if !noPP {
		*lines = append(*lines, "PRESET VARIABLES")
	}
	printVars(lines, &env.Preset, noPP)
	if !noPP {
		*lines = append(*lines, "")
	}
}

func printGlobals(lines *[]string, env *squarescale.Environment, noPP bool) {
	if !noPP {
		*lines = append(*lines, "|- GLOBAL")
	}
	printVars(lines, &env.Custom.Global, noPP)
}

func printPerService(lines *[]string, env *squarescale.Environment, container *string, noPP bool) {
	for serviceName, vars := range env.Custom.PerService {
		if serviceName == *container || *container == "" {
			if !noPP {
				*lines = append(*lines, fmt.Sprintf("|- %s", serviceName))
			}
			printVars(lines, &vars, noPP)
		}
	}
}

func printVars(lines *[]string, vars *map[string]string, noPP bool) {
	format := ""
	if noPP {
		format = "%s=%s"
	} else {
		format = "|-- %s=\"%s\""
	}
	for k, v := range *vars {
		*lines = append(*lines, fmt.Sprintf(format, k, v))
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
