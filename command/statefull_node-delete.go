package command

import (
	"flag"
	"fmt"
	"strings"

	"github.com/squarescale/squarescale-cli/squarescale"
)

// StatefullNodeDeleteCommand is a cli.Command implementation for creating a Squarescale project.
type StatefullNodeDeleteCommand struct {
	Meta
	flagSet *flag.FlagSet
}

// Run is part of cli.Command implementation.
func (c *StatefullNodeDeleteCommand) Run(args []string) int {
	// Parse flags
	c.flagSet = newFlagSet(c, c.Ui)
	alwaysYes := yesFlag(c.flagSet)
	endpoint := endpointFlag(c.flagSet)
	projectArg := projectFlag(c.flagSet)
	nowait := nowaitFlag(c.flagSet)

	if err := c.flagSet.Parse(args); err != nil {
		return 1
	}

	statefullNodeName, err := statefullNodeNameArg(c.flagSet, 0)
	if err != nil {
		return c.errorWithUsage(err)
	}

	if c.flagSet.NArg() > 4 {
		return c.errorWithUsage(fmt.Errorf("Unparsed arguments on the command line: %v", c.flagSet.Args()[1:]))
	}

	if err := validateProjectName(*projectArg); err != nil {
		return c.errorWithUsage(err)
	}

	c.Ui.Info("Are you sure you want to delete " + statefullNodeName + "?")
	if *alwaysYes {
		c.Ui.Info("(approved from command line)")
	} else {
		res, err := c.Ui.Ask("y/N")
		if err != nil {
			return c.error(err)
		} else if res != "Y" && res != "y" {
			return c.cancelled()
		}
	}

	res := c.runWithSpinner("deleting statefull node", endpoint.String(), func(client *squarescale.Client) (string, error) {
		projectName := *projectArg
		fmt.Printf("Delete on project `%s` the statefull node `%s`\n", projectName, statefullNodeName)
		err := client.DeleteStatefullNode(projectName, statefullNodeName)
		return "", err
	})
	if res != 0 {
		return res
	}

	if !*nowait {
		c.runWithSpinner("wait for statefull node delete", endpoint.String(), func(client *squarescale.Client) (string, error) {
			_, err := client.WaitProject(*projectArg)
			if err != nil {
				return "", err
			} else {
				return "", nil
			}
		})
	}

	return 0
}

// Synopsis is part of cli.Command implementation.
func (c *StatefullNodeDeleteCommand) Synopsis() string {
	return "Delete a statefull node from the project."
}

// Help is part of cli.Command implementation.
func (c *StatefullNodeDeleteCommand) Help() string {
	helpText := `
usage: sqsc statefull_node delete [options] <statefull_node_name>

  Delete the statefull node from the project.

`
	return strings.TrimSpace(helpText + optionsFromFlags(c.flagSet))
}
