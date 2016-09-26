package command

import (
	"errors"
	"flag"
	"fmt"

	"github.com/mitchellh/cli"
)

// Meta contain the meta-option that nearly all subcommand inherited.
type Meta struct {
	Ui cli.Ui
}

// EndpointFlag returns a pointer to a string that will be populated
// when the given flagset is parsed with the Squarescale endpoint.
func EndpointFlag(f *flag.FlagSet) *string {
	return f.String("endpoint", "http://www.staging.sqsc.squarely.io", "Squarescale endpoint")
}

// ProjectFlag returns a pointer to a string that will be populated
// when the given flagset is parsed with the Squarescale project.
func ProjectFlag(f *flag.FlagSet) *string {
	return f.String("project", "", "Squarescale project")
}

func validateArgs(project string) error {
	if project == "" {
		return errors.New("Project name must be specified")
	}

	return nil
}

// Error is a shortcut for Meta.Ui.Error with special formatting for errors.
func (m *Meta) Error(err error) {
	m.Ui.Error(fmt.Sprintf("Error: %v", err))
}

// ErrorWithUsage is a shortcut for displaying an error and the usage.
func (m *Meta) ErrorWithUsage(err error, usage string) {
	m.Error(err)
	m.Ui.Output("")
	m.Ui.Output(usage)
}
