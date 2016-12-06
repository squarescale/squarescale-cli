package command

import (
	"errors"
	"flag"
)

func endpointFlag(f *flag.FlagSet) *string {
	return f.String("endpoint", "http://www.staging.sqsc.squarely.io", "Squarescale endpoint")
}

func projectFlag(f *flag.FlagSet) *string {
	return f.String("project", "", "Squarescale project")
}

func dbEngineFlag(f *flag.FlagSet) *string {
	return f.String("engine", "", "Database engine")
}

func dbEngineInstanceFlag(f *flag.FlagSet) *string {
	return f.String("instance", "", "Database engine instance")
}

func dbDisabledFlag(f *flag.FlagSet) *bool {
	return f.Bool("disabled", false, "Deactivate database")
}

func validateProjectName(project string) error {
	if project == "" {
		return errors.New("Project name must be specified")
	}

	return nil
}
