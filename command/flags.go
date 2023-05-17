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

func projectUUIDFlag(f *flag.FlagSet) *string {
	return f.String("project-uuid", "", "Project UUID")
}

func projectNameFlag(f *flag.FlagSet) *string {
	return f.String("project-name", "", "Project name")
}

func projectHybridClusterFlag(f *flag.FlagSet) *bool {
	return f.Bool("hybrid-cluster", false, "Enable/Disable hybrid cluster")
}

func providerFlag(f *flag.FlagSet) *string {
	return f.String("provider", "", "Cloud provider name")
}

func regionFlag(f *flag.FlagSet) *string {
	return f.String("region", "", "Cloud provider region name")
}

func organizationNameFlag(f *flag.FlagSet) *string {
	return f.String("name", "", "organization name")
}

func organizationEmailFlag(f *flag.FlagSet) *string {
	return f.String("email", "", "contact email")
}

func serviceFlag(f *flag.FlagSet) *string {
	return f.String("service", "", "Service aka Docker container to configure")
}

func filterServiceFlag(f *flag.FlagSet) *string {
	return f.String("service", "", "Filter on service name")
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

func containerNoRunCmdFlag(f *flag.FlagSet) *bool {
	return f.Bool("no-command", false, "Disable command override")
}

func containerSchedulingGroupsFlag(f *flag.FlagSet) *string {
	return f.String("scheduling-groups", "", "This is the scheduling groups of your service. Format: ${scheduling_group_name_1},${scheduling_group_name_2}")
}

func containerRunCmdFlag(f *flag.FlagSet) *string {
	return f.String("command", "", "Command to run when starting container (override)")
}

func containerLimitMemoryFlag(f *flag.FlagSet) *int {
	return f.Int("memory", 256, "This is the maximum amount of memory your service will be able to use until it is killed and restarted automatically.")
}

func containerLimitCPUFlag(f *flag.FlagSet) *int {
	return f.Int("cpu", 100, "This is an indicative limit of how much CPU your service requires.")
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
	return f.String("env", "", "JSON file containing all environment variables")
}

func envParameterFlag(f *flag.FlagSet, ret *[]string) {
	f.Func("env-param", "add one/multiple environment parameter(s) in the form param=value (type is `string`)\n\tcan be speficied multiple times and takes precedence over values defined in -env", func(s string) error {
		if len(strings.Split(s, "=")) != 2 {
			return errors.New(fmt.Sprintf("env-param %v not in the form param=value", s))
		}
		*ret = append(*ret, s)
		return nil
	})
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

func schedulingGroupNameArg(f *flag.FlagSet, arg int) (string, error) {
	value := f.Arg(arg)
	if value == "" {
		return "", errors.New("Scheduling group name must be specified")
	} else {
		return value, nil
	}
}

func externalNodeNameArg(f *flag.FlagSet, arg int) (string, error) {
	value := f.Arg(arg)
	if value == "" {
		return "", errors.New("External node name must be specified")
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

	res := "\nOptions:\n\n"
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

func batchNameFlag(f *flag.FlagSet) *string {
	return f.String("batch-name", "", "Batch name")
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

func batchLimitIOPSFlag(f *flag.FlagSet) *int {
	return f.Int("iops", 0, "This is an indicative limit of how many I/O operation per second your service requires.")
}

func batchNoRunCmdFlag(f *flag.FlagSet) *bool {
	return f.Bool("no-command", false, "Disable command override")
}

func dockerCapabilitiesFlag(f *flag.FlagSet) *string {
	return f.String("docker-capabilities", "", "This is a list of enabled docker capabilities")
}

func noDockerCapabilitiesFlag(f *flag.FlagSet) *bool {
	return f.Bool("no-docker-capabilities", false, "Disable all docker capabilities")
}

func dockerDevicesFlag(f *flag.FlagSet) *string {
	return f.String("docker-devices", "", "The list of device mappings if any, separated with commas. format: src:dst[:opt][,src:dst[:opt]]*")
}

func autostartFlag(f *flag.FlagSet) *bool {
	return f.Bool("auto-start", true, "Allow automatic start")
}

func maxclientdisconnect(f *flag.FlagSet) *string {
	return f.String("max-client-disconnect", "",
		"specifies a duration in second (or add a label suffix like '30m' or '1h' or '1d') "+
			"during which the service is considered normally disconnected. "+
			"After this duration the service will be tentatively rescheduled elsewhere. "+
			"To reset to default set to 0 "+
			"By default, this duration is quite short: 2m")
}

func entrypointFlag(f *flag.FlagSet) *string {
	return f.String("entrypoint", "", "This is the script / program that will be executed")
}

func volumeFlag(f *flag.FlagSet) *string {
	return f.String("volumes", "", "Volumes. format: ${VOL2_NAME}:${VOL2_MOUNT_POINT}:ro,${VOL2_NAME}:${VOL2_MOUNT_POINT}:rw")
}

func dockerImageNameFlag(f *flag.FlagSet) *string {
	return f.String("docker-image", "", "Docker image name and potential repository and version (ex: [repoName[:repoPort]/]image[:version])")
}

func dockerImagePrivateFlag(f *flag.FlagSet) *bool {
	return f.Bool("docker-image-private", false, "Docker image is private")
}

func dockerImageUsernameFlag(f *flag.FlagSet) *string {
	return f.String("docker-image-user", "", "Docker image user")
}

func dockerImagePasswordFlag(f *flag.FlagSet) *string {
	return f.String("docker-image-password", "", "Docker image password")
}

func periodicBatchFlag(f *flag.FlagSet) *bool {
	return f.Bool("periodic", false, "batch periodicity")
}

func cronExpressionFlag(f *flag.FlagSet) *string {
	return f.String("cron", "", "cron expression")
}

func timeZoneNameFlag(f *flag.FlagSet) *string {
	return f.String("time", "", "time zone name")
}

func clusterSchedulingGroupsFlag(f *flag.FlagSet) *string {
	return f.String("scheduling-groups", "", "Those are the scheduling groups of your cluster. Format: ${scheduling_group_name_1},${scheduling_group_name_2}")
}

func clusterNameFlag(f *flag.FlagSet) *string {
	return f.String("cluster-name", "", "Cluster name")
}

func clusterSizeFlag(f *flag.FlagSet) *uint {
	return f.Uint("size", 0, "Cluster Size")
}

func dbEngineFlag(f *flag.FlagSet) *string {
	return f.String("engine", "", "Database engine")
}

func dbSizeFlag(f *flag.FlagSet) *string {
	return f.String("size", "", "Database size")
}

func dbVersionFlag(f *flag.FlagSet) *string {
	return f.String("version", "", "Database version")
}

func dbBackupFlag(f *flag.FlagSet) *bool {
	return f.Bool("backup", false, "Enable database backup")
}

func envAllFlag(f *flag.FlagSet) *bool {
	return f.Bool("all", false, "Print all variables")
}

func envRemoveFlag(f *flag.FlagSet) *bool {
	return f.Bool("remove", false, "Remove the key from environment variables")
}

func externalNodePublicIP(f *flag.FlagSet) *string {
	return f.String("public-ip", "", "External node public IP")
}

func loadBalancerIDFlag(f *flag.FlagSet) *int64 {
	return f.Int64("id", 0, "id of the targeted load balancer inside project")
}

func loadBalancerDisableFlag(f *flag.FlagSet) *bool {
	return f.Bool("disable", false, "Disable load balancer")
}

func networkRuleNameFlag(f *flag.FlagSet) *string {
	return f.String("name", "", "name of the rule")
}

func networkServiceNameFlag(f *flag.FlagSet) *string {
	return f.String("service-name", "", "name of the service the rule will be attached")
}

func networkExternalProtocolFlag(f *flag.FlagSet) *string {
	return f.String("external-protocol", "", "name of the externally exposed protocol")
}

func networkInternalProtocolFlag(f *flag.FlagSet) *string {
	return f.String("internal-protocol", "", "name of the internally mapped protocol")
}

func networkInternalPortFlag(f *flag.FlagSet) *int {
	return f.Int("internal-port", 0, "value of the internal port")
}

func networkDomainFlag(f *flag.FlagSet) *string {
	return f.String("domain-expression", "", "custom domain the service is accessed")
}

func networkPathFlag(f *flag.FlagSet) *string {
	return f.String("path", "", "custom path prefix routed to the service")
}

func volumeNameFlag(f *flag.FlagSet) *string {
	return f.String("name", "", "Volume name")
}

func volumeSizeFlag(f *flag.FlagSet) *int {
	return f.Int("size", 1, "Volume size (in Gb)")
}

func volumeTypeFlag(f *flag.FlagSet) *string {
	return f.String("type", "gp2", "Volume type")
}

func volumeZoneFlag(f *flag.FlagSet) *string {
	return f.String("zone", "eu-west-1a", "Volume zone")
}
