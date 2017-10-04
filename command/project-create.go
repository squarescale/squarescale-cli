package command

import (
	"flag"
	"fmt"
	"strings"

	"github.com/squarescale/squarescale-cli/squarescale"
)

// ProjectCreateCommand is a cli.Command implementation for creating a Squarescale project.
type ProjectCreateCommand struct {
	Meta
	flagSet         *flag.FlagSet
	DatabaseEngine  string
	DatabaseClass   string
	DatabaseDisable bool
}

// Run is part of cli.Command implementation.
func (c *ProjectCreateCommand) Run(args []string) int {
	// Parse flags
	c.flagSet = newFlagSet(c, c.Ui)
	alwaysYes := yesFlag(c.flagSet)
	nowait := nowaitFlag(c.flagSet)
	endpoint := endpointFlag(c.flagSet)
	wantedProjectName := projectNameFlag(c.flagSet)
	infraType := infraTypeFlag(c.flagSet)
	c.flagSet.BoolVar(&c.DatabaseDisable, "no-db", false, "Disable database creation")
	c.flagSet.StringVar(&c.DatabaseEngine, "db-engine", "", "Select database engine")
	c.flagSet.StringVar(&c.DatabaseClass, "db-class", "", "Select database instance class")
	if err := c.flagSet.Parse(args); err != nil {
		return 1
	}

	if *wantedProjectName == "" && c.flagSet.Arg(0) != "" {
		name := c.flagSet.Arg(0)
		wantedProjectName = &name
	}

	if c.flagSet.NArg() > 1 {
		return c.errorWithUsage(fmt.Errorf("Unparsed arguments on the command line: %v", c.flagSet.Args()[1:]))
	}

	if *infraType != "high-availability" && *infraType != "single-node" {
		return c.errorWithUsage(fmt.Errorf("Unknown infrastructure type: %v. Correct values are high-availability or single-node", *infraType))
	}

	var taskId int

	res := c.runWithSpinner("create project", endpoint.String(), func(client *squarescale.Client) (string, error) {
		var definitiveName string
		var err error

		if c.DatabaseEngine != "" {
			engines, err := client.GetAvailableDBEngines()
			if err != nil {
				return "", err
			}
			valid := false
			for _, engine := range engines {
				if engine == c.DatabaseEngine {
					valid = true
					break
				}
			}
			if !valid {
				return "", fmt.Errorf("Cannot validate database engine '%s'. Must be one of '%s'", c.DatabaseEngine, strings.Join(engines, "', '"))
			}
		}

		if c.DatabaseClass != "" {
			classes, err := client.GetAvailableDBInstances()
			if err != nil {
				return "", err
			}
			valid := false
			for _, class := range classes {
				if class == c.DatabaseClass {
					valid = true
					break
				}
			}
			if !valid {
				return "", fmt.Errorf("Cannot validate database class '%s'. Must be one of '%s'", c.DatabaseClass, strings.Join(classes, "', '"))
			}
		}

		if *wantedProjectName != "" {
			valid, same, fmtName, err := client.CheckProjectName(*wantedProjectName)
			if err != nil {
				return "", fmt.Errorf("Cannot validate project name '%s'", *wantedProjectName)
			}

			if !valid {
				return "", fmt.Errorf(
					"Project name '%s' is invalid (already taken or not well formed), please choose another one",
					*wantedProjectName)
			}

			if !same {
				if err := c.askConfirmName(alwaysYes, fmtName); err != nil {
					return "", err
				}
			}

			definitiveName = fmtName

		} else {
			generatedName, err := client.FindProjectName()
			if err != nil {
				return "", err
			}

			if err := c.askConfirmName(alwaysYes, generatedName); err != nil {
				return "", err
			}

			definitiveName = generatedName
		}

		taskId, err = client.CreateProject(definitiveName, *infraType, c.DatabaseEngine, c.DatabaseClass, !c.DatabaseDisable)

		return fmt.Sprintf("[#%d] Created project '%s'", taskId, definitiveName), err
	})
	if res != 0 {
		return res
	}

	if !*nowait {
		res = c.runWithSpinner("wait for project creation", endpoint.String(), func(client *squarescale.Client) (string, error) {
			task, err := client.WaitTask(taskId)
			if err != nil {
				return "", err
			} else {
				return task.Status, nil
			}
		})
	}

	return res
}

// Synopsis is part of cli.Command implementation.
func (c *ProjectCreateCommand) Synopsis() string {
	return "Create Project in Squarescale"
}

// Help is part of cli.Command implementation.
func (c *ProjectCreateCommand) Help() string {
	helpText := `
usage: sqsc project create [options] <project_name>

  Creates a new project using the provided project name. If no project name is
  provided, then a new name is automatically generated.

`
	return strings.TrimSpace(helpText + optionsFromFlags(c.flagSet))
}

func (c *ProjectCreateCommand) askConfirmName(alwaysYes *bool, name string) error {
	c.pauseSpinner()
	c.Ui.Warn(fmt.Sprintf("Project will be created as '%s', is this ok?", name))
	ok, err := AskYesNo(c.Ui, alwaysYes, "Is this ok?", true)
	if err != nil {
		return err
	} else if !ok {
		return CancelledError
	}

	c.startSpinner()
	return nil
}

func infraTypeFlag(f *flag.FlagSet) *string {
	return f.String("infra-type", "high-availability", "Set the infrastructure configuration.")
}
