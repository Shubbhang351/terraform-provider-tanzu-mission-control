/*
Copyright 2022 VMware, Inc. All Rights Reserved.
SPDX-License-Identifier: MPL-2.0
*/

package ekscluster

const (
	ResourceName = "tanzu-mission-control_ekscluster"

	CredentialNameKey          = "credential_name" //nolint:gosec
	RegionKey                  = "region"
	NameKey                    = "name"
	specKey                    = "spec"
	StatusKey                  = "status"
	waitKey                    = "ready_wait_timeout"
	clusterGroupKey            = "cluster_group"
	clusterGroupDefaultValue   = "default"
	proxyNameKey               = "proxy"
	configKey                  = "config"
	nodepoolKey                = "nodepool"
	roleArnKey                 = "role_arn"
	kubernetesVersionKey       = "kubernetes_version"
	tagsKey                    = "tags"
	kubernetesNetworkConfigKey = "kubernetes_network_config"
	serviceCidrKey             = "service_cidr"
	loggingKey                 = "logging"
	apiServerKey               = "api_server"
	auditKey                   = "audit"
	authenticatorKey           = "authenticator"
	controllerManagerKey       = "controller_manager"
	schedulerKey               = "scheduler"
	vpcKey                     = "vpc"
	enablePrivateAccessKey     = "enable_private_access"
	enablePublicAccessKey      = "enable_public_access"
	publicAccessCidrsKey       = "public_access_cidrs"
	securityGroupsKey          = "security_groups"
	subnetIdsKey               = "subnet_ids"

	infoKey                     = "info"
	amiTypeKey                  = "ami_type"
	amiInfoKey                  = "ami_info"
	amiIDKey                    = "ami_id"
	overrideBootstrapCmdKey     = "override_bootstrap_cmd"
	capacityTypeKey             = "capacity_type"
	rootDiskSizeKey             = "root_disk_size"
	nodeLabelsKey               = "node_labels"
	launchTemplateKey           = "launch_template"
	idKey                       = "id"
	nameKey                     = "name"
	versionKey                  = "version"
	remoteAccessKey             = "remote_access"
	sshKeyKey                   = "ssh_key"
	scalingConfigKey            = "scaling_config"
	desiredSizeKey              = "desired_size"
	maxSizeKey                  = "max_size"
	minSizeKey                  = "min_size"
	updateConfigKey             = "update_config"
	maxUnavailableNodesKey      = "max_unavailable_nodes"
	maxUnavailablePercentageKey = "max_unavailable_percentage"
	taintsKey                   = "taints"
	effectKey                   = "effect"
	keyKey                      = "key"
	valueKey                    = "value"
	instanceTypesKey            = "instance_types"

	readyCondition = "Ready"
	errorSeverity  = "ERROR"
)
