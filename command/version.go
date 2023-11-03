package command

import (
	"bytes"
	"fmt"
)

// VersionCommand is a cli.Command implementation for displaying the cli version.
type VersionCommand struct {
	Meta
	Name     string
	Version  string
	Arch     string
	Os       string
	Revision string
	Branch   string
	Date     string
}

// Run is part of cli.Command implementation.
func (cmd *VersionCommand) Run(args []string) int {
	var versionString bytes.Buffer

	fmt.Fprintf(&versionString, "%s version %s", cmd.Name, cmd.Version)
	if cmd.Branch != "" {
		fmt.Fprintf(&versionString, " %s", cmd.Branch)
	}
	if cmd.Revision != "" {
		fmt.Fprintf(&versionString, " %s", cmd.Revision)
	}
	if cmd.Date != "" {
		fmt.Fprintf(&versionString, " %s", cmd.Date)
	}
	if cmd.Os != "" {
		fmt.Fprintf(&versionString, " %s", cmd.Os)
	}
	if cmd.Arch != "" {
		fmt.Fprintf(&versionString, " %s", cmd.Arch)
	}

	cmd.Ui.Output(versionString.String())
	return 0
}

// Synopsis is part of cli.Command implementation.
func (cmd *VersionCommand) Synopsis() string {
	return fmt.Sprintf("Print %s version and quit", cmd.Name)
}

// Help is part of cli.Command implementation.
func (cmd *VersionCommand) Help() string {
	return ""
}
