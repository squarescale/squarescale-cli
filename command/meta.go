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

func (m *Meta) info(message string) int {
	m.Ui.Info(message)
	return 0
}

func (m *Meta) error(err error) int {
	m.Ui.Error(fmt.Sprintf("Error: %v", err))
	return 1
}

func (m *Meta) errorWithUsage(err error, usage string) int {
	retCode := m.error(err)
	m.Ui.Output("")
	m.Ui.Output(usage)
	return retCode
}
