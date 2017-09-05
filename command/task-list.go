package command

import (
	"encoding/json"
	"flag"
	"fmt"
	"strings"
	"time"

	"github.com/squarescale/squarescale-cli/squarescale"
)

// TaskListCommand is a cli.Command implementation for scaling the db of a project.
type TaskListCommand struct {
	Meta
	flagSet *flag.FlagSet
}

// Run is part of cli.Command implementation.
func (c *TaskListCommand) Run(args []string) int {
	c.flagSet = newFlagSet(c, c.Ui)
	endpoint := endpointFlag(c.flagSet)
	debug := c.flagSet.Bool("debug-details", false, "Show debug info about each task")
	projectName := projectFlag(c.flagSet)

	if err := c.flagSet.Parse(args); err != nil {
		return 1
	}

	if c.flagSet.NArg() > 0 {
		return c.errorWithUsage(fmt.Errorf("Unparsed arguments on the command line: %v", c.flagSet.Args()))
	}

	return c.runWithSpinner("fetch tasks", endpoint.String(), func(client *squarescale.Client) (string, error) {
		tasks, err := client.GetTasks(*projectName)
		if err != nil {
			return "", err
		}

		var out string
		var table string
		table += "ID\tTYPE\tSTATUS\tTIME\tCOMPLETE EVENT\tWAITING\n"
		for _, t := range tasks {
			table += fmt.Sprintf("%d\t%s\t%s\t%s\t%s\t%s\n",
				t.Id,
				t.TaskType,
				t.StatusWithHold(),
				t.LatestTime(time.Stamp),
				t.CompletedBy,
				strings.Join(t.WaitingEvents, " "))
		}

		out += c.FormatTable(table, true)

		if *debug {
			for _, t := range tasks {
				out += fmt.Sprintf("\n\nTask %d %s\n", t.Id, t.TaskType)
				data, err := json.MarshalIndent(t, "", "\t")
				if err != nil {
					return out, err
				}
				out += string(data)
			}
		}

		return out, nil
	})
}

// Synopsis is part of cli.Command implementation.
func (c *TaskListCommand) Synopsis() string {
	return "Show the list of tasks in a project"
}

// Help is part of cli.Command implementation.
func (c *TaskListCommand) Help() string {
	helpText := `
usage: sqsc task list [options]

  Show the list of tasks in a project.

`
	return strings.TrimSpace(helpText + optionsFromFlags(c.flagSet))
}
