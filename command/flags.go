package command

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/mitchellh/cli"
)

func endpointFlag(f *flag.FlagSet) *string {
	env := os.Getenv("SQSC_ENV")
	if env == "" {
		env = "production"
	}
	defaultValue := os.Getenv("SQSC_ENDPOINT")
	if defaultValue == "" {
		defaultValue = "https://www." + env + ".sqsc.squarely.io"
	}

	endpoint := f.String("endpoint", defaultValue, "Squarescale endpoint")
	normalizedEndpoint := strings.TrimSuffix(*endpoint, "/")
	return &normalizedEndpoint
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

func filterNameFlag(f *flag.FlagSet) *string {
	return f.String("name", "", "Filter on name")
}

func filterTypeFlag(f *flag.FlagSet) *string {
	return f.String("container", "", "Filter on type")
}

func portFlag(f *flag.FlagSet) *int {
	return f.Int("port", -1, "Port number")
}

func containerInstancesFlag(f *flag.FlagSet) *int {
	return f.Int("instances", -1, "Number of container instances")
}

func containerBuildServiceFlag(f *flag.FlagSet) *string {
	return f.String("build-service", "", "Build service to use (internal|travis)")
}

func containerUpdateCmdFlag(f *flag.FlagSet) *string {
	return f.String("update", "", "Update command to run for each build")
}

func clusterSizeFlag(f *flag.FlagSet) *uint {
	return f.Uint("cluster-size", 0, "Cluster Size")
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

func httpsFlag(f *flag.FlagSet) *bool {
	return f.Bool("https", false, "Enable https")
}

func certFlag(f *flag.FlagSet) *string {
	return f.String("cert", "", "Path to a PEM-encoded certificate")
}

func certChainFlag(f *flag.FlagSet) *string {
	return f.String("cert-chain", "", "Path to a PEM-encoded certificate chain")
}

func secretKeyFlag(f *flag.FlagSet) *string {
	return f.String("secret-key", "", "Path to the secret key used for the certificate")
}

func repoUrlFlag(f *flag.FlagSet) *string {
	return f.String("url", "", "Git remote url")
}

func yesFlag(f *flag.FlagSet) *bool {
	return f.Bool("yes", false, "Assume yes")
}

func nowaitFlag(f *flag.FlagSet) *bool {
	return f.Bool("nowait", false, "Don't wait for operation to complete")
}

func projectNameArg(f *flag.FlagSet, arg int) (string, error) {
	value := f.Arg(arg)
	if value == "" {
		return "", errors.New("Project name must be specified")
	} else {
		return value, nil
	}
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
