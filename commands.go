package main

import (
	"github.com/mitchellh/cli"
	"github.com/squarescale/squarescale-cli/command"
)

// Commands returns a map of command factories.
func Commands(meta *command.Meta) map[string]cli.CommandFactory {
	return map[string]cli.CommandFactory{
		"login": func() (cli.Command, error) {
			return &command.LoginCommand{
				Meta: *meta,
			}, nil
		},
		"batch": func() (cli.Command, error) {
			return &command.BatchCommand{}, nil
		},
		"batch list": func() (cli.Command, error) {
			return &command.BatchListCommand{
				Meta: *meta,
			}, nil
		},
		"batch add": func() (cli.Command, error) {
			return &command.BatchAddCommand{
				Meta: *meta,
			}, nil
		},
		"batch delete": func() (cli.Command, error) {
			return &command.BatchDeleteCommand{
				Meta: *meta,
			}, nil
		},
		"batch exec": func() (cli.Command, error) {
			return &command.BatchExecuteCommand{
				Meta: *meta,
			}, nil
		},
		"batch set": func() (cli.Command, error) {
			return &command.BatchSetCommand{
				Meta: *meta,
			}, nil
		},
		"cluster": func() (cli.Command, error) {
			return &command.ClusterCommand{}, nil
		},
		"cluster set": func() (cli.Command, error) {
			return &command.ClusterSetCommand{
				Meta: *meta,
			}, nil
		},
		"cluster node list": func() (cli.Command, error) {
			return &command.ClusterMemberListCommand{
				Meta: *meta,
			}, nil
		},
		"cluster node set": func() (cli.Command, error) {
			return &command.ClusterMemberSetCommand{
				Meta: *meta,
			}, nil
		},
		"volume": func() (cli.Command, error) {
			return &command.VolumeCommand{}, nil
		},
		"volume add": func() (cli.Command, error) {
			return &command.VolumeAddCommand{
				Meta: *meta,
			}, nil
		},
		"volume list": func() (cli.Command, error) {
			return &command.VolumeListCommand{
				Meta: *meta,
			}, nil
		},
		"volume delete": func() (cli.Command, error) {
			return &command.VolumeDeleteCommand{
				Meta: *meta,
			}, nil
		},
		"extra-node": func() (cli.Command, error) {
			return &command.ExtraNodeCommand{}, nil
		},
		"extra-node add": func() (cli.Command, error) {
			return &command.ExtraNodeAddCommand{
				Meta: *meta,
			}, nil
		},
		"extra-node list": func() (cli.Command, error) {
			return &command.ExtraNodeListCommand{
				Meta: *meta,
			}, nil
		},
		"extra-node delete": func() (cli.Command, error) {
			return &command.ExtraNodeDeleteCommand{
				Meta: *meta,
			}, nil
		},
		"extra-node bind": func() (cli.Command, error) {
			return &command.ExtraNodeBindCommand{
				Meta: *meta,
			}, nil
		},
		"scheduling-group": func() (cli.Command, error) {
			return &command.SchedulingGroupCommand{}, nil
		},
		"scheduling-group add": func() (cli.Command, error) {
			return &command.SchedulingGroupAddCommand{
				Meta: *meta,
			}, nil
		},
		"scheduling-group delete": func() (cli.Command, error) {
			return &command.SchedulingGroupDeleteCommand{
				Meta: *meta,
			}, nil
		},
		"scheduling-group list": func() (cli.Command, error) {
			return &command.SchedulingGroupListCommand{
				Meta: *meta,
			}, nil
		},
		"scheduling-group get": func() (cli.Command, error) {
			return &command.SchedulingGroupGetCommand{
				Meta: *meta,
			}, nil
		},
		"scheduling-group assign": func() (cli.Command, error) {
			return &command.SchedulingGroupAssignCommand{
				Meta: *meta,
			}, nil
		},
		"scheduling-group unassign": func() (cli.Command, error) {
			return &command.SchedulingGroupUnAssignCommand{
				Meta: *meta,
			}, nil
		},
		"external-node": func() (cli.Command, error) {
			return &command.ExternalNodeCommand{}, nil
		},
		"external-node list": func() (cli.Command, error) {
			return &command.ExternalNodeListCommand{
				Meta: *meta,
			}, nil
		},
		"external-node add": func() (cli.Command, error) {
			return &command.ExternalNodeAddCommand{
				Meta: *meta,
			}, nil
		},

		"external-node get": func() (cli.Command, error) {
			return &command.ExternalNodeGetCommand{
				Meta: *meta,
			}, nil
		},
		"service": func() (cli.Command, error) {
			return &command.ServiceCommand{}, nil
		},
		"service add": func() (cli.Command, error) {
			return &command.ServiceAddCommand{
				Meta: *meta,
			}, nil
		},
		"service delete": func() (cli.Command, error) {
			return &command.ServiceDeleteCommand{
				Meta: *meta,
			}, nil
		},
		"service list": func() (cli.Command, error) {
			return &command.ServiceListCommand{
				Meta: *meta,
			}, nil
		},
		"service schedule": func() (cli.Command, error) {
			return &command.ServiceScheduleCommand{
				Meta: *meta,
			}, nil
		},
		"service set": func() (cli.Command, error) {
			return &command.ServiceSetCommand{
				Meta: *meta,
			}, nil
		},
		"service show": func() (cli.Command, error) {
			return &command.ServiceShowCommand{
				Meta: *meta,
			}, nil
		},
		"env": func() (cli.Command, error) {
			return &command.EnvCommand{}, nil
		},
		"env get": func() (cli.Command, error) {
			return &command.EnvGetCommand{
				Meta: *meta,
			}, nil
		},
		"env set": func() (cli.Command, error) {
			return &command.EnvSetCommand{
				Meta: *meta,
			}, nil
		},
		"project": func() (cli.Command, error) {
			return &command.ProjectCommand{}, nil
		},
		"project create": func() (cli.Command, error) {
			return &command.ProjectCreateCommand{
				Meta: *meta,
			}, nil
		},
		"project provision": func() (cli.Command, error) {
			return &command.ProjectProvisionCommand{
				Meta: *meta,
			}, nil
		},
		"project unprovision": func() (cli.Command, error) {
			return &command.ProjectUnprovisionCommand{
				Meta: *meta,
			}, nil
		},
		"project list": func() (cli.Command, error) {
			return &command.ProjectListCommand{
				Meta: *meta,
			}, nil
		},
		"project get": func() (cli.Command, error) {
			return &command.ProjectGetCommand{
				Meta: *meta,
			}, nil
		},
		"project details": func() (cli.Command, error) {
			return &command.ProjectDetailsCommand{
				Meta: *meta,
			}, nil
		},
		"project remove": func() (cli.Command, error) {
			return &command.ProjectRemoveCommand{
				Meta: *meta,
			}, nil
		},
		"project settings": func() (cli.Command, error) {
			return &command.ProjectSettingsCommand{
				Meta: *meta,
			}, nil
		},
		"db": func() (cli.Command, error) {
			return &command.DBCommand{}, nil
		},
		"db list": func() (cli.Command, error) {
			return &command.DBListCommand{
				Meta: *meta,
			}, nil
		},
		"db set": func() (cli.Command, error) {
			return &command.DBSetCommand{
				Meta: *meta,
			}, nil
		},
		"db show": func() (cli.Command, error) {
			return &command.DBShowCommand{
				Meta: *meta,
			}, nil
		},
		"redis": func() (cli.Command, error) {
			return &command.RedisCommand{}, nil
		},
		"redis list": func() (cli.Command, error) {
			return &command.RedisListCommand{
				Meta: *meta,
			}, nil
		},
		"redis add": func() (cli.Command, error) {
			return &command.RedisAddCommand{
				Meta: *meta,
			}, nil
		},
		"redis delete": func() (cli.Command, error) {
			return &command.RedisDeleteCommand{
				Meta: *meta,
			}, nil
		},
		"lb": func() (cli.Command, error) {
			return &command.LBCommand{}, nil
		},
		"lb list": func() (cli.Command, error) {
			return &command.LBListCommand{
				Meta: *meta,
			}, nil
		},
		"lb get": func() (cli.Command, error) {
			return &command.LBGetCommand{
				Meta: *meta,
			}, nil
		},
		"lb set": func() (cli.Command, error) {
			return &command.LBSetCommand{
				Meta: *meta,
			}, nil
		},
		"status": func() (cli.Command, error) {
			return &command.StatusCommand{
				Meta: *meta,
			}, nil
		},
		"version": func() (cli.Command, error) {
			return &command.VersionCommand{
				Meta:     *meta,
				Version:  Version,
				Arch:     GoArch,
				Os:       GoOs,
				Revision: GitCommit,
				Branch:   GitBranch,
				Date:     BuildDate,
				Name:     Name,
			}, nil
		},
		"organization": func() (cli.Command, error) {
			return &command.OrganizationCommand{}, nil
		},
		"organization add": func() (cli.Command, error) {
			return &command.OrganizationAddCommand{
				Meta: *meta,
			}, nil
		},
		"organization delete": func() (cli.Command, error) {
			return &command.OrganizationDeleteCommand{
				Meta: *meta,
			}, nil
		},
		"organization list": func() (cli.Command, error) {
			return &command.OrganizationListCommand{
				Meta: *meta,
			}, nil
		},
		"network-rule": func() (cli.Command, error) {
			return &command.NetworkRuleCommand{}, nil
		},
		"network-rule create": func() (cli.Command, error) {
			return &command.NetworkRuleCreateCommand{
				Meta: *meta,
			}, nil
		},
		"network-rule list": func() (cli.Command, error) {
			return &command.NetworkRuleListCommand{
				Meta: *meta,
			}, nil
		},
		"network-rule delete": func() (cli.Command, error) {
			return &command.NetworkRuleDeleteCommand{
				Meta: *meta,
			}, nil
		},
		"network-policy get": func() (cli.Command, error) {
			return &command.NetworkPolicyGetCommand{
				Meta: *meta,
			}, nil
		},
		"network-policy add": func() (cli.Command, error) {
			return &command.NetworkPolicyAddCommand{
				Meta: *meta,
			}, nil
		},
		"network-policy list": func() (cli.Command, error) {
			return &command.NetworkPolicyListCommand{
				Meta: *meta,
			}, nil
		},
		"network-policy delete": func() (cli.Command, error) {
			return &command.NetworkPolicyDeleteCommand{
				Meta: *meta,
			}, nil
		},
		"network-policy deploy": func() (cli.Command, error) {
			return &command.NetworkPolicyDeployCommand{
				Meta: *meta,
			}, nil
		},
	}
}
