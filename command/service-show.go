package command

import (
	"errors"
	"flag"
	"fmt"
	"strings"

	"github.com/squarescale/squarescale-cli/squarescale"
)

// ServiceShowCommand is a cli.Command implementation for listing all projects.
type ServiceShowCommand struct {
	Meta
	flagSet *flag.FlagSet
}

// Run is part of cli.Command implementation.
func (cmd *ServiceShowCommand) Run(args []string) int {
	cmd.flagSet = newFlagSet(cmd, cmd.Ui)
	endpoint := endpointFlag(cmd.flagSet)
	projectUUID := projectUUIDFlag(cmd.flagSet)
	projectName := projectNameFlag(cmd.flagSet)
	containerArg := filterServiceFlag(cmd.flagSet)
	if err := cmd.flagSet.Parse(args); err != nil {
		return 1
	}

	if cmd.flagSet.NArg() > 0 {
		return cmd.errorWithUsage(fmt.Errorf("Unparsed arguments on the command line: %v", cmd.flagSet.Args()))
	}

	if *projectUUID == "" && *projectName == "" {
		return cmd.errorWithUsage(errors.New("Project name or uuid is mandatory"))
	}

	return cmd.runWithSpinner("showing service", endpoint.String(), func(client *squarescale.Client) (string, error) {
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

		containers, err := client.GetServices(UUID)
		if err != nil {
			return "", err
		}

		var msg string
		found := false
		for _, co := range containers {
			if *containerArg != "" && *containerArg != co.Name {
				continue
			}
			found = true
			if msg != "" {
				msg += "\n-----------\n\n"
			}
			tbl := ""
			tbl += fmt.Sprintf("Name:\t%s\n", co.Name)
			tbl += fmt.Sprintf("Auto start: \t%v \n", co.AutoStart)
			tbl += fmt.Sprintf("Size:\t%d/%d\n", co.Running, co.Size)
			tbl += fmt.Sprintf("Run Command:\t%s\n", co.RunCommand)
			tbl += fmt.Sprintf("Web Port:\t%d\n", co.WebPort)
			tbl += fmt.Sprintf("Memory limit:\t%d MB\n", co.Limits.Memory)
			tbl += fmt.Sprintf("CPU limit:\t%d MHz\n", co.Limits.CPU)
			tbl += fmt.Sprintf("Max Client Disconnect: \t%s\n", co.MaxClientDisconnect)
			tbl += fmt.Sprintf("Enabled capabilities: \t%s \n", strings.Join(co.DockerCapabilities, ","))
			tbl += fmt.Sprintf("Volumes:\n")
			for _, vol := range co.Volumes {
				tbl += fmt.Sprintf("\tname:%s mount:%s R/O:%v\n", vol.Name, vol.MountPoint, vol.ReadOnly)
			}
			tbl += fmt.Sprintf("Docker devices:\n")
			for _, dev := range co.DockerDevices {
				tbl += fmt.Sprintf("\tsrc:%s dst:%s opts:%s\n", dev.SRC, dev.DST, dev.OPT)
			}
			msg += cmd.FormatTable(tbl, false)
			msg += "\n"
		}

		if len(containers) == 0 {
			msg = fmt.Sprintf("\nNo service defined in project %s %s\n", *projectName, *projectUUID)
		} else if !found {
			msg = fmt.Sprintf("\nService %s not found in project %s %s\n", *containerArg, *projectName, *projectUUID)
		}

		return msg, nil
	})
}

// Synopsis is part of cli.Command implementation.
func (cmd *ServiceShowCommand) Synopsis() string {
	return "Show service aka Docker container of project"
}

// Help is part of cli.Command implementation.
func (cmd *ServiceShowCommand) Help() string {
	helpText := `
usage: sqsc service show [options]

  Show service aka Docker container of project.
`
	return strings.TrimSpace(helpText + optionsFromFlags(cmd.flagSet))
}
