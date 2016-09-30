package command

import (
	"errors"
	"flag"
	"fmt"
	"time"

	"github.com/briandowns/spinner"
	"github.com/mitchellh/cli"
	"github.com/squarescale/squarescale-cli/squarescale"
	"github.com/squarescale/squarescale-cli/tokenstore"
)

// Meta contain the meta-option that nearly all subcommand inherited.
type Meta struct {
	Ui   cli.Ui
	spin *spinner.Spinner
}

// DefaultMeta returns a default meta object with an initialized spinner.
func DefaultMeta(ui cli.Ui) *Meta {
	s := spinner.New(spinner.CharSets[24], 100*time.Millisecond)
	s.FinalMSG = "... done\n"

	return &Meta{
		Ui:   ui,
		spin: s,
	}
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

func (m *Meta) runWithSpinner(text, endpoint string, action func(token string) error) error {
	m.spin.Start()
	m.spin.Suffix = " " + text
	defer m.spin.Stop()

	token, err := m.ensureLogin(endpoint)
	if err != nil {
		return err
	}

	return action(token)
}

func (m *Meta) ensureLogin(endpoint string) (string, error) {
	token, err := tokenstore.GetToken(endpoint)
	if err != nil {
		return "", err
	}

	err = squarescale.ValidateToken(endpoint, token)
	if err != nil {
		return "", fmt.Errorf("You're not authenticated, please login first (%v)", err)
	}

	return token, nil
}
