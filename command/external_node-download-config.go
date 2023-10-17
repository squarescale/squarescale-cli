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

func (b *ExternalNodeDownloadConfigCommand) Run(args []string) int {
	b.flagSet = newFlagSet(b, b.Ui)
	endpoint := endpointFlag(b.flagSet)
	projectUUID := projectUUIDFlag(b.flagSet)
	projectName := projectNameFlag(b.flagSet)
	configName := externalNodeConfigNameFlag(b.flagSet)

	if err := b.flagSet.Parse(args); err != nil {
		return 1
	}

	if *projectUUID == "" && *projectName == "" {
		return b.errorWithUsage(errors.New("Project name or uuid is mandatory"))
	}

	if *configName == "" || (*configName != "all" && *configName != "openvpn" && *configName != "consul" && *configName != "nomad") {
		return b.errorWithUsage(errors.New(fmt.Sprintf("Invalid service configuration name: %s", *configName)))
	}

	externalNodeName, err := externalNodeNameArg(b.flagSet, 0)
	if err != nil {
		return b.errorWithUsage(err)
	}

	if b.flagSet.NArg() > 1 {
		return b.errorWithUsage(fmt.Errorf("Unparsed arguments on the command line: %v", b.flagSet.Args()))
	}

	return b.runWithSpinner("downloading external node service configuration file(s)", endpoint.String(), func(client *squarescale.Client) (string, error) {
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
func (b *ExternalNodeDownloadConfigCommand) Synopsis() string {
	return "Download service configuration file(s) for external node of project"
}

// Help is part of cli.Command implementation.
func (b *ExternalNodeDownloadConfigCommand) Help() string {
	helpText := `
usage: sqsc external-node download-config [options] <external_node_name>

  Download service configuration file(s) for external node of project.
`
	return strings.TrimSpace(helpText + optionsFromFlags(b.flagSet))
}
