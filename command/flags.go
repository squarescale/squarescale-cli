package command

import (
	"errors"
	"flag"
	"fmt"
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

func optionsFromFlags(fs *flag.FlagSet) string {
	res := "Options\n\n"
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
