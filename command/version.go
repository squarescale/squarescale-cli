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
	Revision string
}

// Run is part of cli.Command implementation.
func (c *VersionCommand) Run(args []string) int {
	var versionString bytes.Buffer

	fmt.Fprintf(&versionString, "%s version %s", c.Name, c.Version)
	if c.Revision != "" {
		fmt.Fprintf(&versionString, " (%s)", c.Revision)
	}

	c.Ui.Output(versionString.String())
	return 0
}

// Synopsis is part of cli.Command implementation.
func (c *VersionCommand) Synopsis() string {
	return fmt.Sprintf("Print %s version and quit", c.Name)
}

// Help is part of cli.Command implementation.
func (c *VersionCommand) Help() string {
	return ""
}
