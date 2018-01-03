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
	flagSet   *flag.FlagSet
	DbDisable bool
	Db        squarescale.DbConfig
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
	nodeSize := nodeSizeFlag(c.flagSet)
	c.flagSet.BoolVar(&c.DbDisable, "no-db", false, "Disable database creation")
	c.flagSet.StringVar(&c.Db.Engine, "db-engine", "", "Select database engine")
	c.flagSet.StringVar(&c.Db.Size, "db-size", "", "Select database size")
	if err := c.flagSet.Parse(args); err != nil {
		return 1
	}

	c.Db.Enabled = !c.DbDisable

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

		nodeSizes, err := client.GetClusterNodeSizes()
		if err != nil {
			return "", err
		}
		if nodeSizes != nil {
			if *nodeSize != "" {
				if !nodeSizes.CheckSize(*nodeSize, *infraType) {
					return "", fmt.Errorf("Cannot validate node size '%s'. Must be one of '%s'", *nodeSize, strings.Join(nodeSizes.ListSizes(*infraType), "', '"))
				}
			} else {
				if *infraType == "single-node" {
					*nodeSize = "dev"
				} else {
					*nodeSize = "small"
				}
			}
		}

		if c.Db.Engine != "" {
			engines, err := client.GetAvailableDBEngines()
			if err != nil {
				return "", err
			}
			valid := false
			for _, engine := range engines {
				if engine == c.Db.Engine {
					valid = true
					break
				}
			}
			if !valid {
				return "", fmt.Errorf("Cannot validate database engine '%s'. Must be one of '%s'", c.Db.Engine, strings.Join(engines, "', '"))
			}
		}

		if c.Db.Size != "" {
			sizes, err := client.GetDBSizes()
			if err != nil {
				return "", err
			}
			if !sizes.CheckID(c.Db.Size, *infraType) {
				return "", fmt.Errorf("Cannot validate database size '%s'. Must be one of '%s'", c.Db.Size, strings.Join(sizes.ListIds(*infraType), "', '"))
			}
		}

		if !c.DbDisable && c.Db.Engine == "" && c.Db.Size == "" {
			c.Db.Engine = "postgres"
			if client.HasNewDB() {
				if *infraType == "single-node" {
					c.Db.Size = "dev"
				} else {
					c.Db.Size = "small"
				}
			} else {
				c.Db.Size = "micro"
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

		taskId, err = client.CreateProject(definitiveName, *infraType, *nodeSize, c.Db)

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

func nodeSizeFlag(f *flag.FlagSet) *string {
	return f.String("node-size", "small", "Set the cluster node size.")
}
