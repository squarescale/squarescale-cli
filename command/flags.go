package command

import (
	"errors"
	"flag"
	"fmt"

	"github.com/mitchellh/cli"
)

func endpointFlag(f *flag.FlagSet) *string {
	return f.String("endpoint", "http://www.staging.sqsc.squarely.io", "Squarescale endpoint")
}

func projectFlag(f *flag.FlagSet) *string {
	return f.String("project", "", "Squarescale project")
}

func projectNameFlag(f *flag.FlagSet) *string {
	return f.String("name", "", "Project name")
}

func containerFlag(f *flag.FlagSet) *string {
	return f.String("container", "", "Container to configure")
}

func portFlag(f *flag.FlagSet) *int {
	return f.Int("port", -1, "Port number")
}

func containerInstancesFlag(f *flag.FlagSet) *int {
	return f.Int("instances", -1, "Number of container instances")
}

func containerUpdateCmdFlag(f *flag.FlagSet) *string {
	return f.String("update", "", "Update command to run for each build")
}

func dbEngineFlag(f *flag.FlagSet) *string {
	return f.String("engine", "", "Database engine")
}

func dbEngineInstanceFlag(f *flag.FlagSet) *string {
	return f.String("instance", "", "Database engine instance")
}

func disabledFlag(f *flag.FlagSet, doc string) *bool {
	return f.Bool("disabled", false, doc)
}

func validateProjectName(project string) error {
	if project == "" {
		return errors.New("Project name must be specified")
	}

	return nil
}

func newFlagSet(cmd cli.Command, ui cli.Ui) *flag.FlagSet {
	return &flag.FlagSet{
		// Bind the usage of the command to the usage of the flag set.
		Usage: func() { ui.Output(cmd.Help()) },
	}
}

func optionsFromFlags(fs *flag.FlagSet) string {
	if fs == nil {
		return ""
	}

	res := "Options:\n\n"
	fs.VisitAll(func(f *flag.Flag) {
		s := fmt.Sprintf("  -%s", f.Name)
		name, usage := flag.UnquoteUsage(f)
		if len(name) > 0 {
			s += " " + name
		}

		s += fmt.Sprintf("\n        %s", usage)
		if name == "string" {
			s += fmt.Sprintf(" (default %q)\n", f.DefValue)
		} else {
			s += fmt.Sprintf(" (default %v)\n", f.DefValue)
		}

		res += s
	})

	return res
}
