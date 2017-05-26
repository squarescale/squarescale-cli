package command

import (
	"bytes"
	"flag"
	"fmt"
	"strings"

	"github.com/olekukonko/tablewriter"
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

	var printFunction func(string, *map[string]string, *bytes.Buffer)
	if c.niceFormat {
		printFunction = prettyPrintVars
	} else {
		printFunction = printVars
	}

	queryvar := c.flagSet.Arg(0)
	queryval := ""
	if c.flagSet.NArg() > 1 {
		return c.errorWithUsage(fmt.Errorf("Unparsed arguments on the command line: %v", c.flagSet.Args()))
	}

	if err := validateProjectName(*project); err != nil {
		return c.errorWithUsage(err)
	}

	return c.runWithSpinner("list environment variables", endpoint.String(), func(client *squarescale.Client) (string, error) {
		env, err := client.EnvironmentVariables(*project)
		if err != nil {
			return "", err
		}

		var linesBuffer bytes.Buffer

		if queryvar != "" {
			queryval = env.Preset[queryvar]
		} else if printPreset {
			printFunction("Presets", &env.Preset, &linesBuffer)
		}

		if queryvar != "" {
			if val, ok := env.Global[queryvar]; ok && !*preSet {
				queryval = val
			}
		} else if printGlobal {
			printFunction("Global", &env.Global, &linesBuffer)
		}

		if noFlag || *container != "" {
			for containerName, vars := range env.PerService {
				if queryvar != "" {
					if val, ok := vars[queryvar]; ok && containerName == *container {
						queryval = val
					}
				} else if containerName == *container || *container == "" {
					printFunction(containerName, &vars, &linesBuffer)
				}
			}
		}

		if queryvar != "" {
			return queryval, nil
		} else {
			return linesBuffer.String(), nil
		}
	})
}

func prettyPrintVars(title string, variables *map[string]string, lines *bytes.Buffer) {
	lines.WriteString(fmt.Sprintf("%s\n", title))

	if len(*variables) == 0 {
		lines.WriteString("  none\n")
	}

	data := make([][]string, 0, len(*variables))
	for k, v := range *variables {
		data = append(data, []string{k, v})
	}

	table := tablewriter.NewWriter(lines)
	table.AppendBulk(data)
	table.SetBorder(false)
	table.SetColumnSeparator("=")
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.Render()
}

func printVars(title string, variables *map[string]string, lines *bytes.Buffer) {
	for k, v := range *variables {
		lines.WriteString(fmt.Sprintf("%s=%s\n", k, v))
	}
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
