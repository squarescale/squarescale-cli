package command

import (
	"errors"
	"flag"
	"fmt"
	"net/url"
	"strings"

	"github.com/squarescale/squarescale-cli/squarescale"
)

// ProjectCreateCommand is a cli.Command implementation for creating project.
type ProjectCreateCommand struct {
	Meta
	flagSet *flag.FlagSet
}

// Run is part of cli.Command implementation.
func (cmd *ProjectCreateCommand) Run(args []string) int {
	// Parse flags
	cmd.flagSet = newFlagSet(cmd, cmd.Ui)
	endpoint := endpointFlag(cmd.flagSet)
	alwaysYes := yesFlag(cmd.flagSet)
	name := projectNameFlag(cmd.flagSet)
	uuid := cmd.flagSet.String("uuid", "", "set the uuid of the project")
	organization := cmd.flagSet.String("organization", "", "set the organization the project will belongs to")
	provider := cmd.flagSet.String("provider", "", "set the cloud provider (aws or azure or outscale)")
	region := cmd.flagSet.String("region", "", "set the cloud provider region (eu-west-1)")
	credential := cmd.flagSet.String("credential", "", "set the credential used to build the infrastructure")
	infraType := cmd.flagSet.String("infra-type", "high-availability", "Set the infrastructure configuration (high-availability or single-node)")
	monitoringEngine := cmd.flagSet.String("monitoring", "", "Set the monitoring configuration (netdata)")
	nodeSize := cmd.flagSet.String("node-size", "", "Set the cluster node size")
	nodeCount := cmd.flagSet.Int("node-count", 3, "Set the cluster node count (1 for single-node, odd number for high-availability)")
	rootDiskSizeGB := cmd.flagSet.Int("root-disk-size", 20, "Set the root filesystem size (in GB)")
	slackURL := cmd.flagSet.String("slackbot", "", "Set the Slack webhook URL")
	hybridClusterEnabled := cmd.flagSet.Bool("hybrid-cluster-enabled", false, "Enable Hybrid Cluster")
	externalESURL := cmd.flagSet.String("external-elasticsearch", "", "Set the external ElasticSearch URL")

	dbEngine := cmd.flagSet.String("db-engine", "", "Select database engine")
	dbSize := cmd.flagSet.String("db-size", "", "Select database size")
	dbVersion := cmd.flagSet.String("db-version", "", "Select database version")
	dbBackupEnabled := cmd.flagSet.Bool("db-backup", false, "Enable database automated backups")
	dbBackupRetention := cmd.flagSet.Int("db-backup-retention", 0, "Database automated backups retention days")

	consulEnabled := cmd.flagSet.Bool("consul-enabled", false, "Enable Consul")
	nomadEnabled := cmd.flagSet.Bool("nomad-enabled", false, "Enable Nomad")
	vaultEnabled := cmd.flagSet.Bool("vault-enabled", false, "Enable Vault")
	observabilityEnabled := cmd.flagSet.Bool("observability-enabled", false, "Enable Observability")
	elasticsearchEnabled := cmd.flagSet.Bool("elasticsearch-enabled", false, "Enable ElasticSearch")

	consulBasicAuth := cmd.flagSet.String("consul-basic-auth", "", "Set Consul Basic Authentication credentials")
	nomadBasicAuth := cmd.flagSet.String("nomad-basic-auth", "", "Set Nomad Basic Authentication credentials")
	vaultBasicAuth := cmd.flagSet.String("vault-basic-auth", "", "Set Vault Basic Authentication credentials")
	observabilityBasicAuth := cmd.flagSet.String("observability-basic-auth", "", "Set Observability Basic Authentication credentials")
	elasticsearchBasicAuth := cmd.flagSet.String("elasticsearch-basic-auth", "", "Set ElasticSearch Basic Authentication credentials")

	consulPrefix := cmd.flagSet.String("consul-prefix", "", "Set Consul HTTP Prefix")
	nomadPrefix := cmd.flagSet.String("nomad-prefix", "", "Set Nomad HTTP Prefix")
	vaultPrefix := cmd.flagSet.String("vault-prefix", "", "Set Vault HTTP Prefix")
	observabilityPrefix := cmd.flagSet.String("observability-prefix", "", "Set Observability HTTP Prefix")
	elasticsearchPrefix := cmd.flagSet.String("elasticsearch-prefix", "", "Set ElasticSearch HTTP Prefix")

	consulIpWhiteList := cmd.flagSet.String("consul-ip-whitelist", "", "Set Consul IP whitelist")
	nomadIpWhiteList := cmd.flagSet.String("nomad-ip-whitelist", "", "Set Nomad IP whitelist")
	vaultIpWhiteList := cmd.flagSet.String("vault-ip-whitelist", "", "Set Vault IP whitelist")
	observabilityIpWhiteList := cmd.flagSet.String("observability-ip-whitelist", "", "Set Observability IP whitelist")
	elasticsearchIpWhiteList := cmd.flagSet.String("elasticsearch-ip-whitelist", "", "Set ElasticSearch IP whitelist")

	nowait := nowaitFlag(cmd.flagSet)

	payload := squarescale.JSONObject{}

	if err := cmd.flagSet.Parse(args); err != nil {
		return 1
	}

	if *name == "" && cmd.flagSet.Arg(0) != "" {
		nameArg := cmd.flagSet.Arg(0)
		name = &nameArg
	}

	if cmd.flagSet.NArg() > 1 {
		return cmd.errorWithUsage(fmt.Errorf("Unparsed arguments on the command line: %v", cmd.flagSet.Args()[1:]))
	}

	if *name == "" {
		return cmd.errorWithUsage(errors.New("Project name is mandatory"))
	}

	payload["name"] = *name
	payload["hybrid_cluster_enabled"] = *hybridClusterEnabled

	if *uuid != "" {
		payload["uuid"] = *uuid
	}

	if *organization != "" {
		payload["organization_name"] = *organization
	}

	if *provider == "" {
		return cmd.errorWithUsage(errors.New("Cloud provider is mandatory"))
	}

	payload["provider"] = *provider

	if *region == "" {
		return cmd.errorWithUsage(errors.New("Cloud provider region is mandatory"))
	}

	payload["region"] = *region

	if *credential != "" {
		payload["credential_name"] = *credential
	}

	if *nodeSize == "" {
		return cmd.errorWithUsage(errors.New("node size is mandatory"))
	}

	payload["node_size"] = *nodeSize
	payload["desired_size"] = *nodeCount
	payload["root_disk_size_gb"] = *rootDiskSizeGB

	if *infraType != "high-availability" && *infraType != "single-node" {
		return cmd.errorWithUsage(fmt.Errorf("Unknown infrastructure type: %v. Correct values are high-availability or single-node", *infraType))
	} else if *infraType == "high-availability" {
		payload["infra_type"] = "high_availability"
		if *nodeCount < 3 {
			return cmd.error(fmt.Errorf("Infrastructure type %v can not have less than 3 nodes (%d)", payload["infra_type"], *nodeCount))
		}
		if *nodeCount % 2 != 1 {
			return cmd.error(fmt.Errorf("Infrastructure type %v must have odd number of nodes (%d)", payload["infra_type"], *nodeCount))
		}
	} else {
		payload["infra_type"] = "single_node"
		if *nodeCount != 1 {
			return cmd.error(fmt.Errorf("Infrastructure type %v can not have more than 1 node (%d)", payload["infra_type"], *nodeCount))
		}
	}

	if *monitoringEngine != "" && *monitoringEngine != "netdata" {
		return cmd.errorWithUsage(fmt.Errorf("Unknown monitoring engine: %v. Correct values are empty-string (aka '') or netdata", *monitoringEngine))
	} else if *monitoringEngine == "netdata" {
		payload["monitoring"] = map[string]string{"engine": "netdata", "enabled": "true"}
	} else {
		payload["monitoring"] = map[string]string{"engine": "", "enabled": "false"}
	}

	if *dbEngine != "" && *dbSize == "" {
		return cmd.errorWithUsage(errors.New("if db engine is present, db size must be set"))
	} else if *dbEngine == "" && *dbSize != "" {
		return cmd.errorWithUsage(errors.New("if db size is present, db engine must be set"))
	} else if *dbEngine != "" && *dbSize != "" {
		if *dbBackupRetention < 0 {
			return cmd.errorWithUsage(errors.New("retention value can not be negative"))
		}
		if *dbBackupEnabled && *dbBackupRetention == 0 {
			return cmd.errorWithUsage(errors.New("retention value can only be strictly positive when db backup enabled"))
		}
		if !*dbBackupEnabled && *dbBackupRetention != 0 {
			return cmd.errorWithUsage(errors.New("retention value can only be 0 when db backup disabled"))
		}
		payload["databases"] = []map[string]interface{}{{"engine": *dbEngine, "size": *dbSize, "version": *dbVersion, "backup_enabled": *dbBackupEnabled, "backup_retention_days": *dbBackupRetention}}
	}

	if *slackURL != "" {
		payload["slack_webhook"] = *slackURL
	}

	if *externalESURL != "" {
		url, err := url.Parse(*externalESURL)
		if err != nil {
			return cmd.errorWithUsage(errors.New(fmt.Sprintf("URL format error on %s: %q", *externalESURL, err)))
		}
		// "host" could be valid entry but it results in empty scheme
		if url.Scheme != "http" && url.Scheme != "https" {
			return cmd.errorWithUsage(errors.New(fmt.Sprintf("URL scheme '%s' not supported for external ElasticSearch endpoint (should be either 'http' or 'https')", url.Scheme)))
		}
		payload["external_elasticsearch"] = *externalESURL
	}

	var integrated_services []map[string]string

	if *consulEnabled {
		integrated_services = append(integrated_services, map[string]string{"name": "consul", "enabled": "true", "basicauth": *consulBasicAuth, "prefix": *consulPrefix, "ipwhitelist": *consulIpWhiteList})
	}

	if *nomadEnabled {
		integrated_services = append(integrated_services, map[string]string{"name": "nomad", "enabled": "true", "basicauth": *nomadBasicAuth, "prefix": *nomadPrefix, "ipwhitelist": *nomadIpWhiteList})
	}

	if *vaultEnabled {
		integrated_services = append(integrated_services, map[string]string{"name": "vault", "enabled": "true", "basicauth": *vaultBasicAuth, "prefix": *vaultPrefix, "ipwhitelist": *vaultIpWhiteList})
	}

	if *observabilityEnabled {
		integrated_services = append(integrated_services, map[string]string{"name": "observability", "enabled": "true", "basicauth": *observabilityBasicAuth, "prefix": *observabilityPrefix, "ipwhitelist": *observabilityIpWhiteList})
	}

	if *elasticsearchEnabled {
		integrated_services = append(integrated_services, map[string]string{"name": "elasticsearch", "enabled": "true", "basicauth": *elasticsearchBasicAuth, "prefix": *elasticsearchPrefix, "ipwhitelist": *elasticsearchIpWhiteList})
	}

	// TODO: add check to prevent both external ES and integrated ES

	payload["integrated_services"] = integrated_services
	// ask confirmation
	cmd.Ui.Warn("Project will be created with the following configuration :")
	cmd.Ui.Warn(fmt.Sprintf("name : %s", *name))
	if *uuid != "" {
		cmd.Ui.Warn(fmt.Sprintf("uuid : %s", *uuid))
	}
	if *organization != "" {
		cmd.Ui.Warn(fmt.Sprintf("organization : %s", *organization))
	}
	cmd.Ui.Warn(fmt.Sprintf("cloud provider : %s", *provider))
	cmd.Ui.Warn(fmt.Sprintf("cloud provider region : %s", *region))
	cmd.Ui.Warn(fmt.Sprintf("credential : %s", *credential))
	cmd.Ui.Warn(fmt.Sprintf("node size : %s", *nodeSize))
	cmd.Ui.Warn(fmt.Sprintf("node count : %d", *nodeCount))
	cmd.Ui.Warn(fmt.Sprintf("infra type : %s", *infraType))
	cmd.Ui.Warn(fmt.Sprintf("monitoring engine : %s", *monitoringEngine))
	if *slackURL != "" {
		cmd.Ui.Warn(fmt.Sprintf("Slack webhook : %s", *slackURL))
	}
	if *dbEngine != "" && *dbSize != "" {
		cmd.Ui.Warn(fmt.Sprintf("database engine : %s", *dbEngine))
		cmd.Ui.Warn(fmt.Sprintf("database size : %s", *dbSize))
		if *dbVersion != "" {
			cmd.Ui.Warn(fmt.Sprintf("database version : %s", *dbVersion))
		}
		cmd.Ui.Warn(fmt.Sprintf("database backup : %v", *dbBackupEnabled))
		cmd.Ui.Warn(fmt.Sprintf("database backup retention : %v", *dbBackupRetention))
	}

	ok, err := AskYesNo(cmd.Ui, alwaysYes, "Proceed ?", true)
	if err != nil {
		return cmd.error(err)
	} else if !ok {
		return cmd.error(CancelledError)
	}

	var project squarescale.Project

	res := cmd.runWithSpinner("create project", endpoint.String(), func(client *squarescale.Client) (string, error) {
		var err error
		project, err = client.CreateProject(&payload)
		return fmt.Sprintf("Created project '%s' with UUID '%s'", project.Name, project.UUID), err
	})

	if res != 0 {
		return res
	}

	if !*nowait {
		res = cmd.runWithSpinner("wait for project creation", endpoint.String(), func(client *squarescale.Client) (string, error) {
			projectStatus, err := client.WaitProject(project.UUID, 5)
			if err != nil {
				return projectStatus, err
			} else {
				return projectStatus, nil
			}
		})
	}

	return res
}

// Synopsis is part of cli.Command implementation.
func (cmd *ProjectCreateCommand) Synopsis() string {
	return "Create project"
}

// Help is part of cli.Command implementation.
func (cmd *ProjectCreateCommand) Help() string {
	helpText := `
usage: sqsc project create [options]

  Creates a new project using the provided options.

`
	return strings.TrimSpace(helpText + optionsFromFlags(cmd.flagSet))
}

func (cmd *ProjectCreateCommand) askConfirmation(alwaysYes *bool, name string) error {
	cmd.pauseSpinner()

	cmd.startSpinner()
	return nil
}
