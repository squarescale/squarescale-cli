package command

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/mitchellh/cli"
	log "github.com/sirupsen/logrus"
)

type endpoint string

func (e *endpoint) String() string {
	if *e == "" {
		env := os.Getenv("SQSC_ENV")
		if env == "" {
			env = "production"
		}
		defaultValue := os.Getenv("SQSC_ENDPOINT")
		if defaultValue == "" && env == "dev" {
			defaultValue = "http://web.service.consul:3000"
		} else if defaultValue == "" && env == "production" {
			defaultValue = "https://www.squarescale.io"
		} else if defaultValue == "" {
			defaultValue = "https://www." + env + ".sqsc.io"
		}

		log.WithField("endpoint", defaultValue).Debug()
		return defaultValue
	}
	log.WithField("endpoint", *e).Debug()
	return string(*e)
}

func (e *endpoint) Set(value string) error {
	normalizedEndpoint := strings.TrimSuffix(value, "/")
	*e = endpoint(normalizedEndpoint)
	return nil
}

var endPointFlag endpoint

func endpointFlag(f *flag.FlagSet) *endpoint {
	f.Var(&endPointFlag, "endpoint", "SquareScale endpoint")
	return &endPointFlag
}

func projectFlag(f *flag.FlagSet) *string {
	return f.String("project", "", "SquareScale project")
}

func providerFlag(f *flag.FlagSet) *string {
	return f.String("provider", "", "Cloud provider name")
}

func regionFlag(f *flag.FlagSet) *string {
	return f.String("region", "", "Cloud provider region name")
}

func organizationFlag(f *flag.FlagSet) *string {
	return f.String("organization", "", "organization name")
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
	return f.String("type", "", "Filter on type")
}

func portFlag(f *flag.FlagSet) *int {
	return f.Int("port", -1, "Port number")
}

func exprFlag(f *flag.FlagSet) *string {
	return f.String("expr", "", "custom domain expression")
}

func containerInstancesFlag(f *flag.FlagSet) *int {
	return f.Int("instances", -1, "Number of container instances")
}

func containerBuildServiceFlag(f *flag.FlagSet) *string {
	return f.String("build-service", "", "Build service to use (internal|travis)")
}

func containerNoRunCmdFlag(f *flag.FlagSet) *bool {
	return f.Bool("no-command", false, "Disable command override")
}

func containerRunCmdFlag(f *flag.FlagSet) *string {
	return f.String("command", "", "Command to run when starting container (override)")
}

func containerLimitMemoryFlag(f *flag.FlagSet) *int {
	return f.Int("memory", -1, "This is the maximum amount of memory your service will be able to use until it is killed and restarted automatically.")
}

func containerLimitCPUFlag(f *flag.FlagSet) *int {
	return f.Int("cpu", -1, "This is an indicative limit of how much CPU your service requires.")
}

func containerLimitIOPSFlag(f *flag.FlagSet) *int {
	return f.Int("iops", -1, "This is an indicative limit of how many I/O operation per second your service requires.")
}

func containerLimitNetFlag(f *flag.FlagSet) *int {
	return f.Int("net", -1, "This is an indicative limit of how much network bandwidth your service requires.")
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

func envFileFlag(f *flag.FlagSet) *string {
	return f.String("env", "", "json file with all env variable")
}

func containerNameArg(f *flag.FlagSet, arg int) (string, error) {
	value := f.Arg(arg)
	if value == "" {
		return "", errors.New("Container name must be specified")
	} else {
		return value, nil
	}
}

func volumeNameArg(f *flag.FlagSet, arg int) (string, error) {
	value := f.Arg(arg)
	if value == "" {
		return "", errors.New("Volume name must be specified")
	} else {
		return value, nil
	}
}

func statefulNodeNameArg(f *flag.FlagSet, arg int) (string, error) {
	value := f.Arg(arg)
	if value == "" {
		return "", errors.New("Stateful-node name must be specified")
	} else {
		return value, nil
	}
}

func batchNameArg(f *flag.FlagSet, arg int) (string, error) {
	value := f.Arg(arg)
	if value == "" {
		return "", errors.New("Batch name must be specified")
	} else {
		return value, nil
	}
}

func serviceNameArg(f *flag.FlagSet, arg int) (string, error) {
	value := f.Arg(arg)
	if value == "" {
		return "", errors.New("Service name must be specified")
	} else {
		return value, nil
	}
}

func projectNameArg(f *flag.FlagSet, arg int) (string, error) {
	value := f.Arg(arg)
	if value == "" {
		return "", errors.New("Project name must be specified")
	} else {
		return value, nil
	}
}

func redisNameArg(f *flag.FlagSet, arg int) (string, error) {
	value := f.Arg(arg)
	if value == "" {
		return "", errors.New("Redis name must be specified")
	} else {
		return value, nil
	}
}

func organizationNameArg(f *flag.FlagSet, arg int) (string, error) {
	value := f.Arg(arg)
	if value == "" {
		return "", errors.New("Organization name must be specified")
	} else {
		return value, nil
	}
}

func repoOrImageInstancesFlag(f *flag.FlagSet) *int {
	return f.Int("instances", 1, "Number of container instances at creation")
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

func batchRunCmdFlag(f *flag.FlagSet) *string {
	return f.String("command", "", "Command to run when starting batch (override)")
}

func batchLimitMemoryFlag(f *flag.FlagSet) *int {
	return f.Int("memory", -1, "This is the maximum amount of memory your batch will be able to use until it is killed and restarted automatically.")
}

func batchLimitCPUFlag(f *flag.FlagSet) *int {
	return f.Int("cpu", -1, "This is an indicative limit of how much CPU your batch requires.")
}

func batchLimitNetFlag(f *flag.FlagSet) *int {
	return f.Int("net", -1, "This is an indicative limit of how much network bandwidth your batch requires.")
}

func batchNoRunCmdFlag(f *flag.FlagSet) *bool {
	return f.Bool("no-command", false, "Disable command override")
}
