package command

import (
	"flag"

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
