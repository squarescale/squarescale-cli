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
func (c *ProjectCreateCommand) Run(args []string) int {
	// Parse flags
	c.flagSet = newFlagSet(c, c.Ui)
	endpoint := endpointFlag(c.flagSet)
	alwaysYes := yesFlag(c.flagSet)
	name := projectNameFlag(c.flagSet)
	uuid := c.flagSet.String("uuid", "", "set the uuid of the project")
	organization := c.flagSet.String("organization", "", "set the organization the project will belongs to")
	provider := c.flagSet.String("provider", "", "set the cloud provider (aws or azure or outscale)")
	region := c.flagSet.String("region", "", "set the cloud provider region (eu-west-1)")
	credential := c.flagSet.String("credential", "", "set the credential used to build the infrastructure")
	infraType := c.flagSet.String("infra-type", "high-availability", "Set the infrastructure configuration (high-availability or single-node)")
	monitoringEngine := c.flagSet.String("monitoring", "", "Set the monitoring configuration (netdata)")
	nodeSize := c.flagSet.String("node-size", "", "Set the cluster node size")
	rootDiskSizeGB := c.flagSet.Int("root-disk-size", 20, "Set the root filesystem size (in GB)")
	slackURL := c.flagSet.String("slackbot", "", "Set the Slack webhook URL")
	hybridClusterEnabled := c.flagSet.Bool("hybrid-cluster-enabled", false, "Enable Hybrid Cluster")
	externalESURL := c.flagSet.String("external-elasticsearch", "", "Set the external ElasticSearch URL")

	dbEngine := c.flagSet.String("db-engine", "", "Select database engine")
	dbSize := c.flagSet.String("db-size", "", "Select database size")
	dbVersion := c.flagSet.String("db-version", "", "Select database version")
	dbBackupEnabled := c.flagSet.Bool("db-backup", false, "Enable database automated backups")
	dbBackupRetention := c.flagSet.Int("db-backup-retention", 0, "Database automated backups retention days")

	consulEnabled := c.flagSet.Bool("consul-enabled", false, "Enable Consul")
	nomadEnabled := c.flagSet.Bool("nomad-enabled", false, "Enable Nomad")
	vaultEnabled := c.flagSet.Bool("vault-enabled", false, "Enable Vault")

	consulBasicAuth := c.flagSet.String("consul-basic-auth", "", "Set Consul Basic Authentication credentials")
	nomadBasicAuth := c.flagSet.String("nomad-basic-auth", "", "Set Nomad Basic Authentication credentials")
	vaultBasicAuth := c.flagSet.String("vault-basic-auth", "", "Set Vault Basic Authentication credentials")

	consulPrefix := c.flagSet.String("consul-prefix", "", "Set Consul HTTP Prefix")
	nomadPrefix := c.flagSet.String("nomad-prefix", "", "Set Nomad HTTP Prefix")
	vaultPrefix := c.flagSet.String("vault-prefix", "", "Set Vault HTTP Prefix")

	consulIpWhiteList := c.flagSet.String("consul-ip-whitelist", "", "Set Consul IP whitelist")
	nomadIpWhiteList := c.flagSet.String("nomad-ip-whitelist", "", "Set Nomad IP whitelist")
	vaultIpWhiteList := c.flagSet.String("vault-ip-whitelist", "", "Set Vault IP whitelist")

	nowait := nowaitFlag(c.flagSet)

	payload := squarescale.JSONObject{}

	if err := c.flagSet.Parse(args); err != nil {
		return 1
	}

	if *name == "" && c.flagSet.Arg(0) != "" {
		nameArg := c.flagSet.Arg(0)
		name = &nameArg
	}

	if c.flagSet.NArg() > 1 {
		return c.errorWithUsage(fmt.Errorf("Unparsed arguments on the command line: %v", c.flagSet.Args()[1:]))
	}

	if *name == "" {
		return c.errorWithUsage(errors.New("Project name is mandatory"))
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
		return c.errorWithUsage(errors.New("Cloud provider is mandatory"))
	}

	payload["provider"] = *provider

	if *region == "" {
		return c.errorWithUsage(errors.New("Cloud provider region is mandatory"))
	}

	payload["region"] = *region

	if *credential != "" {
		payload["credential_name"] = *credential
	}

	if *nodeSize == "" {
		return c.errorWithUsage(errors.New("node size is mandatory"))
	}

	payload["node_size"] = *nodeSize
	payload["root_disk_size_gb"] = *rootDiskSizeGB

	if *infraType != "high-availability" && *infraType != "single-node" {
		return c.errorWithUsage(fmt.Errorf("Unknown infrastructure type: %v. Correct values are high-availability or single-node", *infraType))
	} else if *infraType == "high-availability" {
		payload["infra_type"] = "high_availability"
	} else {
		payload["infra_type"] = "single_node"
	}

	if *monitoringEngine != "" && *monitoringEngine != "netdata" {
		return c.errorWithUsage(fmt.Errorf("Unknown monitoring engine: %v. Correct values are empty-string (aka '') or netdata", *monitoringEngine))
	} else if *monitoringEngine == "netdata" {
		payload["monitoring"] = map[string]string{"engine": "netdata", "enabled": "true"}
	} else {
		payload["monitoring"] = map[string]string{"engine": "", "enabled": "false"}
	}

	if *dbEngine != "" && *dbSize == "" {
		return c.errorWithUsage(errors.New("if db engine is present, db size must be set"))
	} else if *dbEngine == "" && *dbSize != "" {
		return c.errorWithUsage(errors.New("if db size is present, db engine must be set"))
	} else if *dbEngine != "" && *dbSize != "" {
		if *dbBackupRetention < 0 {
			return c.errorWithUsage(errors.New("retention value can not be negative"))
		}
		if *dbBackupEnabled && *dbBackupRetention == 0 {
			return c.errorWithUsage(errors.New("retention value can only be strictly positive when db backup enabled"))
		}
		if !*dbBackupEnabled && *dbBackupRetention != 0 {
			return c.errorWithUsage(errors.New("retention value can only be 0 when db backup disabled"))
		}
		payload["databases"] = []map[string]interface{}{{"engine": *dbEngine, "size": *dbSize, "version": *dbVersion, "backup_enabled": *dbBackupEnabled, "backup_retention_days": *dbBackupRetention}}
	}

	if *slackURL != "" {
		payload["slack_webhook"] = *slackURL
	}

	if *externalESURL != "" {
		url, err := url.Parse(*externalESURL)
		if err != nil {
			return c.errorWithUsage(errors.New(fmt.Sprintf("URL format error on %s: %q", *externalESURL, err)))
		}
		// "host" could be valid entry but it results in empty scheme
		if url.Scheme != "http" && url.Scheme != "https" {
			return c.errorWithUsage(errors.New(fmt.Sprintf("URL scheme '%s' not supported for external ElasticSearch endpoint (should be either 'http' or 'https')", url.Scheme)))
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

	// TODO: add observability and ElasticSearch integrated services

	// TODO: add check to prevent both external ES and integrated ES

	payload["integrated_services"] = integrated_services
	// ask confirmation
	c.Ui.Warn("Project will be created with the following configuration :")
	c.Ui.Warn(fmt.Sprintf("name : %s", *name))
	if *uuid != "" {
		c.Ui.Warn(fmt.Sprintf("uuid : %s", *uuid))
	}
	if *organization != "" {
		c.Ui.Warn(fmt.Sprintf("organization : %s", *organization))
	}
	c.Ui.Warn(fmt.Sprintf("cloud provider : %s", *provider))
	c.Ui.Warn(fmt.Sprintf("cloud provider region : %s", *region))
	c.Ui.Warn(fmt.Sprintf("credential : %s", *credential))
	c.Ui.Warn(fmt.Sprintf("node size : %s", *nodeSize))
	c.Ui.Warn(fmt.Sprintf("infra type : %s", *infraType))
	c.Ui.Warn(fmt.Sprintf("monitoring engine : %s", *monitoringEngine))
	if *slackURL != "" {
		c.Ui.Warn(fmt.Sprintf("Slack webhook : %s", *slackURL))
	}
	if *dbEngine != "" && *dbSize != "" {
		c.Ui.Warn(fmt.Sprintf("database engine : %s", *dbEngine))
		c.Ui.Warn(fmt.Sprintf("database size : %s", *dbSize))
		if *dbVersion != "" {
			c.Ui.Warn(fmt.Sprintf("database version : %s", *dbVersion))
		}
		c.Ui.Warn(fmt.Sprintf("database backup : %v", *dbBackupEnabled))
		c.Ui.Warn(fmt.Sprintf("database backup retention : %v", *dbBackupRetention))
	}

	ok, err := AskYesNo(c.Ui, alwaysYes, "Proceed ?", true)
	if err != nil {
		return c.error(err)
	} else if !ok {
		return c.error(CancelledError)
	}

	var project squarescale.Project

	res := c.runWithSpinner("create project", endpoint.String(), func(client *squarescale.Client) (string, error) {
		var err error
		project, err = client.CreateProject(&payload)
		return fmt.Sprintf("Created project '%s' with UUID '%s'", project.Name, project.UUID), err
	})

	if res != 0 {
		return res
	}

	if !*nowait {
		res = c.runWithSpinner("wait for project creation", endpoint.String(), func(client *squarescale.Client) (string, error) {
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
func (c *ProjectCreateCommand) Synopsis() string {
	return "Create project"
}

// Help is part of cli.Command implementation.
func (c *ProjectCreateCommand) Help() string {
	helpText := `
usage: sqsc project create [options]

  Creates a new project using the provided options.

`
	return strings.TrimSpace(helpText + optionsFromFlags(c.flagSet))
}

func (c *ProjectCreateCommand) askConfirmation(alwaysYes *bool, name string) error {
	c.pauseSpinner()

	c.startSpinner()
	return nil
}
