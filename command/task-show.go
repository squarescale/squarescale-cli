package command

import (
	"flag"
	"fmt"
	"strconv"
	"strings"

	"github.com/squarescale/squarescale-cli/squarescale"
)

// TaskShowCommand is a cli.Command implementation for scaling the db of a project.
type TaskShowCommand struct {
	Meta
	flagSet *flag.FlagSet
}

// Run is part of cli.Command implementation.
func (c *TaskShowCommand) Run(args []string) int {
	c.flagSet = newFlagSet(c, c.Ui)
	endpoint := endpointFlag(c.flagSet)
	if err := c.flagSet.Parse(args); err != nil {
		return 1
	}

	taskId, err := strconv.Atoi(c.flagSet.Arg(0))
	if err != nil {
		return c.errorWithUsage(err)
	} else if taskId <= 0 {
		return c.errorWithUsage(fmt.Errorf("No task id provided"))
	}

	if c.flagSet.NArg() > 1 {
		return c.errorWithUsage(fmt.Errorf("Unparsed arguments on the command line: %v", c.flagSet.Args()[1:]))
	}

	return c.runWithSpinner("fetch task", endpoint.String(), func(client *squarescale.Client) (string, error) {
		t, err := client.GetTask(taskId)
		if err != nil {
			return "", err
		}

		status := [][]string{
			[]string{"Id", strconv.Itoa(t.Id)},
			[]string{"Project Id", strconv.Itoa(t.ProjectId)},
			[]string{"Status", t.Status},
			[]string{"Params", string(t.Params)},
			[]string{"Completed By", t.CompletedBy},
			[]string{"Created At", t.CreatedAt},
			[]string{"Updated At", t.UpdatedAt},
			[]string{"Completed At", t.CompletedAt},
			[]string{"Waiting Events", strconv.Itoa(len(t.WaitingEvents))},
		}
		for _, ev := range t.WaitingEvents {
			status = append(status, []string{"Waiting Event", ev})
		}

		var lines []string
		for _, l := range status {
			lines = append(lines, strings.Join(l, "\t"))
		}
		return c.FormatTable(strings.Join(lines, "\n"), false), nil
	})
}

// Synopsis is part of cli.Command implementation.
func (c *TaskShowCommand) Synopsis() string {
	return "Show the status of a task"
}

// Help is part of cli.Command implementation.
func (c *TaskShowCommand) Help() string {
	helpText := `
usage: sqsc task show [options] <task-id>

  Show the status of a task.

`
	return strings.TrimSpace(helpText + optionsFromFlags(c.flagSet))
}
