package command

import (
	"errors"
	"flag"
	"fmt"
	"strings"

	"github.com/squarescale/squarescale-cli/squarescale"
)

// ExternalNodeDownloadConfigCommand is a cli.Command implementation for getting OpenVPN, Consul & Nomad configuration for an external node.
type ExternalNodeDownloadConfigCommand struct {
	Meta
	flagSet *flag.FlagSet
}

func (cmd *ExternalNodeDownloadConfigCommand) Run(args []string) int {
	cmd.flagSet = newFlagSet(cmd, cmd.Ui)
	endpoint := endpointFlag(cmd.flagSet)
	projectUUID := projectUUIDFlag(cmd.flagSet)
	projectName := projectNameFlag(cmd.flagSet)
	configName := externalNodeConfigNameFlag(cmd.flagSet)

	if err := cmd.flagSet.Parse(args); err != nil {
		return 1
	}

	if *projectUUID == "" && *projectName == "" {
		return cmd.errorWithUsage(errors.New("Project name or uuid is mandatory"))
	}

	if *configName == "" || (*configName != "all" && *configName != "openvpn" && *configName != "consul" && *configName != "nomad") {
		return cmd.errorWithUsage(errors.New(fmt.Sprintf("Invalid service configuration name: %s", *configName)))
	}

	externalNodeName, err := externalNodeNameArg(cmd.flagSet, 0)
	if err != nil {
		return cmd.errorWithUsage(err)
	}

	if cmd.flagSet.NArg() > 1 {
		return cmd.errorWithUsage(fmt.Errorf("Unparsed arguments on the command line: %v", cmd.flagSet.Args()))
	}

	return cmd.runWithSpinner("downloading external node service configuration file(s)", endpoint.String(), func(client *squarescale.Client) (string, error) {
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

		err = client.DownloadConfigExternalNode(UUID, externalNodeName, *configName)
		if err != nil {
			return "", err
		}

		return "All done", nil
	})
}

// Synopsis is part of cli.Command implementation.
func (cmd *ExternalNodeDownloadConfigCommand) Synopsis() string {
	return "Download service configuration file(s) for external node of project"
}

// Help is part of cli.Command implementation.
func (cmd *ExternalNodeDownloadConfigCommand) Help() string {
	helpText := `
usage: sqsc external-node download-config [options] <external_node_name>

  Download service configuration file(s) for external node of project.
`
	return strings.TrimSpace(helpText + optionsFromFlags(cmd.flagSet))
}
