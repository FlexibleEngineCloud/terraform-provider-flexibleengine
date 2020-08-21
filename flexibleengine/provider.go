package flexibleengine

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/mutexkv"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

// This is a global MutexKV for use within this plugin.
var osMutexKV = mutexkv.NewMutexKV()

// Provider returns a schema.Provider for FlexibleEngine.
func Provider() terraform.ResourceProvider {
	provider := &schema.Provider{
		Schema: map[string]*schema.Schema{
			"access_key": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_ACCESS_KEY", ""),
				Description: descriptions["access_key"],
			},

			"secret_key": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_SECRET_KEY", ""),
				Description: descriptions["secret_key"],
			},

			"auth_url": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_AUTH_URL", nil),
				Description: descriptions["auth_url"],
			},

			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: descriptions["region"],
				DefaultFunc: schema.EnvDefaultFunc("OS_REGION_NAME", ""),
			},

			"user_name": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_USERNAME", ""),
				Description: descriptions["user_name"],
			},

			"user_id": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_USER_ID", ""),
				Description: descriptions["user_name"],
			},

			"tenant_id": {
				Type:     schema.TypeString,
				Optional: true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"OS_TENANT_ID",
					"OS_PROJECT_ID",
				}, ""),
				Description: descriptions["tenant_id"],
			},

			"tenant_name": {
				Type:     schema.TypeString,
				Optional: true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"OS_TENANT_NAME",
					"OS_PROJECT_NAME",
				}, ""),
				Description: descriptions["tenant_name"],
			},

			"password": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("OS_PASSWORD", ""),
				Description: descriptions["password"],
			},

			"token": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_AUTH_TOKEN", ""),
				Description: descriptions["token"],
			},

			"security_token": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: descriptions["security_token"],
			},

			"domain_id": {
				Type:     schema.TypeString,
				Optional: true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"OS_USER_DOMAIN_ID",
					"OS_PROJECT_DOMAIN_ID",
					"OS_DOMAIN_ID",
				}, ""),
				Description: descriptions["domain_id"],
			},

			"domain_name": {
				Type:     schema.TypeString,
				Optional: true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"OS_USER_DOMAIN_NAME",
					"OS_PROJECT_DOMAIN_NAME",
					"OS_DOMAIN_NAME",
					"OS_DEFAULT_DOMAIN",
				}, ""),
				Description: descriptions["domain_name"],
			},

			"insecure": {
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_INSECURE", false),
				Description: descriptions["insecure"],
			},

			"endpoint_type": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_ENDPOINT_TYPE", ""),
			},

			"cacert_file": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_CACERT", ""),
				Description: descriptions["cacert_file"],
			},

			"cert": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_CERT", ""),
				Description: descriptions["cert"],
			},

			"key": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_KEY", ""),
				Description: descriptions["key"],
			},

			"swauth": {
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_SWAUTH", false),
				Description: descriptions["swauth"],
			},
		},

		DataSourcesMap: map[string]*schema.Resource{
			"flexibleengine_blockstorage_availability_zones_v3": dataSourceBlockStorageAvailabilityZonesV3(),
			"flexibleengine_blockstorage_volume_v2":             dataSourceBlockStorageVolumeV2(),
			"flexibleengine_compute_availability_zones_v2":      dataSourceComputeAvailabilityZonesV2(),
			"flexibleengine_images_image_v2":                    dataSourceImagesImageV2(),
			"flexibleengine_networking_network_v2":              dataSourceNetworkingNetworkV2(),
			"flexibleengine_networking_secgroup_v2":             dataSourceNetworkingSecGroupV2(),
			"flexibleengine_s3_bucket_object":                   dataSourceS3BucketObject(),
			"flexibleengine_kms_key_v1":                         dataSourceKmsKeyV1(),
			"flexibleengine_kms_data_key_v1":                    dataSourceKmsDataKeyV1(),
			"flexibleengine_rds_flavors_v1":                     dataSourceRdsFlavorV1(),
			"flexibleengine_rds_flavors_v3":                     dataSourceRdsFlavorV3(),
			"flexibleengine_vpc_v1":                             dataSourceVirtualPrivateCloudVpcV1(),
			"flexibleengine_vpc_subnet_v1":                      dataSourceVpcSubnetV1(),
			"flexibleengine_vpc_subnet_ids_v1":                  dataSourceVpcSubnetIdsV1(),
			"flexibleengine_vpc_route_v2":                       dataSourceVPCRouteV2(),
			"flexibleengine_vpc_route_ids_v2":                   dataSourceVPCRouteIdsV2(),
			"flexibleengine_vpc_peering_connection_v2":          dataSourceVpcPeeringConnectionV2(),
			"flexibleengine_sfs_file_system_v2":                 dataSourceSFSFileSystemV2(),
			"flexibleengine_compute_bms_flavors_v2":             dataSourceBMSFlavorV2(),
			"flexibleengine_compute_bms_nic_v2":                 dataSourceBMSNicV2(),
			"flexibleengine_compute_bms_server_v2":              dataSourceBMSServersV2(),
			"flexibleengine_compute_bms_keypairs_v2":            dataSourceBMSKeyPairV2(),
			"flexibleengine_rts_software_config_v1":             dataSourceRtsSoftwareConfigV1(),
			"flexibleengine_rts_stack_v1":                       dataSourceRTSStackV1(),
			"flexibleengine_rts_stack_resource_v1":              dataSourceRTSStackResourcesV1(),
			"flexibleengine_csbs_backup_v1":                     dataSourceCSBSBackupV1(),
			"flexibleengine_csbs_backup_policy_v1":              dataSourceCSBSBackupPolicyV1(),
			"flexibleengine_vbs_backup_policy_v2":               dataSourceVBSBackupPolicyV2(),
			"flexibleengine_vbs_backup_v2":                      dataSourceVBSBackupV2(),
			"flexibleengine_cts_tracker_v1":                     dataSourceCTSTrackerV1(),
			"flexibleengine_dcs_az_v1":                          dataSourceDcsAZV1(),
			"flexibleengine_dcs_maintainwindow_v1":              dataSourceDcsMaintainWindowV1(),
			"flexibleengine_dcs_product_v1":                     dataSourceDcsProductV1(),
			"flexibleengine_cce_node_v3":                        dataSourceCceNodesV3(),
			"flexibleengine_cce_node_ids_v3":                    dataSourceCceNodeIdsV3(),
			"flexibleengine_cce_cluster_v3":                     dataSourceCCEClusterV3(),
			"flexibleengine_dns_zone_v2":                        dataSourceDNSZoneV2(),
			"flexibleengine_dds_flavor_v3":                      dataSourceDDSFlavorV3(),
			"flexibleengine_lb_certificate_v2":                  dataSourceCertificateV2(),
			"flexibleengine_sdrs_domain_v1":                     dataSourceSdrsDomainV1(),
			"flexibleengine_identity_project_v3":                dataSourceIdentityProjectV3(),
			"flexibleengine_identity_role_v3":                   dataSourceIdentityRoleV3(),
		},

		ResourcesMap: map[string]*schema.Resource{
			"flexibleengine_blockstorage_volume_v2":             resourceBlockStorageVolumeV2(),
			"flexibleengine_compute_instance_v2":                resourceComputeInstanceV2(),
			"flexibleengine_compute_interface_attach_v2":        resourceComputeInterfaceAttachV2(),
			"flexibleengine_compute_keypair_v2":                 resourceComputeKeypairV2(),
			"flexibleengine_compute_servergroup_v2":             resourceComputeServerGroupV2(),
			"flexibleengine_compute_floatingip_v2":              resourceComputeFloatingIPV2(),
			"flexibleengine_compute_floatingip_associate_v2":    resourceComputeFloatingIPAssociateV2(),
			"flexibleengine_compute_volume_attach_v2":           resourceComputeVolumeAttachV2(),
			"flexibleengine_dns_recordset_v2":                   resourceDNSRecordSetV2(),
			"flexibleengine_dns_zone_v2":                        resourceDNSZoneV2(),
			"flexibleengine_dcs_instance_v1":                    resourceDcsInstanceV1(),
			"flexibleengine_fw_firewall_group_v2":               resourceFWFirewallGroupV2(),
			"flexibleengine_fw_policy_v2":                       resourceFWPolicyV2(),
			"flexibleengine_fw_rule_v2":                         resourceFWRuleV2(),
			"flexibleengine_images_image_v2":                    resourceImagesImageV2(),
			"flexibleengine_kms_key_v1":                         resourceKmsKeyV1(),
			"flexibleengine_lb_certificate_v2":                  resourceCertificateV2(),
			"flexibleengine_lb_loadbalancer_v2":                 resourceLoadBalancerV2(),
			"flexibleengine_lb_listener_v2":                     resourceListenerV2(),
			"flexibleengine_lb_pool_v2":                         resourcePoolV2(),
			"flexibleengine_lb_member_v2":                       resourceMemberV2(),
			"flexibleengine_lb_monitor_v2":                      resourceMonitorV2(),
			"flexibleengine_lb_whitelist_v2":                    resourceWhitelistV2(),
			"flexibleengine_lb_l7policy_v2":                     resourceL7PolicyV2(),
			"flexibleengine_lb_l7rule_v2":                       resourceL7RuleV2(),
			"flexibleengine_mrs_hybrid_cluster_v1":              resourceMRSHybridClusterV1(),
			"flexibleengine_mrs_cluster_v1":                     resourceMRSClusterV1(),
			"flexibleengine_mrs_job_v1":                         resourceMRSJobV1(),
			"flexibleengine_mls_instance_v1":                    resourceMlsInstanceV1(),
			"flexibleengine_networking_network_v2":              resourceNetworkingNetworkV2(),
			"flexibleengine_networking_subnet_v2":               resourceNetworkingSubnetV2(),
			"flexibleengine_networking_floatingip_v2":           resourceNetworkingFloatingIPV2(),
			"flexibleengine_networking_floatingip_associate_v2": resourceNetworkingFloatingIPAssociateV2(),
			"flexibleengine_networking_port_v2":                 resourceNetworkingPortV2(),
			"flexibleengine_networking_router_v2":               resourceNetworkingRouterV2(),
			"flexibleengine_networking_router_interface_v2":     resourceNetworkingRouterInterfaceV2(),
			"flexibleengine_networking_router_route_v2":         resourceNetworkingRouterRouteV2(),
			"flexibleengine_networking_secgroup_v2":             resourceNetworkingSecGroupV2(),
			"flexibleengine_networking_secgroup_rule_v2":        resourceNetworkingSecGroupRuleV2(),
			"flexibleengine_identity_agency_v3":                 resourceIdentityAgencyV3(),
			"flexibleengine_identity_group_v3":                  resourceIdentityGroupV3(),
			"flexibleengine_identity_group_membership_v3":       resourceIdentityGroupMembershipV3(),
			"flexibleengine_identity_role_v3":                   resourceIdentityRoleV3(),
			"flexibleengine_identity_role_assignment_v3":        resourceIdentityRoleAssignmentV3(),
			"flexibleengine_identity_user_v3":                   resourceIdentityUserV3(),
			"flexibleengine_s3_bucket":                          resourceS3Bucket(),
			"flexibleengine_s3_bucket_policy":                   resourceS3BucketPolicy(),
			"flexibleengine_s3_bucket_object":                   resourceS3BucketObject(),
			"flexibleengine_obs_bucket":                         resourceObsBucket(),
			"flexibleengine_obs_bucket_object":                  resourceObsBucketObject(),
			"flexibleengine_elb_loadbalancer":                   resourceELoadBalancer(),
			"flexibleengine_elb_listener":                       resourceEListener(),
			"flexibleengine_elb_backend":                        resourceBackend(),
			"flexibleengine_elb_health":                         resourceHealth(),
			"flexibleengine_as_group_v1":                        resourceASGroup(),
			"flexibleengine_as_configuration_v1":                resourceASConfiguration(),
			"flexibleengine_as_policy_v1":                       resourceASPolicy(),
			"flexibleengine_smn_topic_v2":                       resourceTopic(),
			"flexibleengine_smn_subscription_v2":                resourceSubscription(),
			"flexibleengine_rds_instance_v1":                    resourceRdsInstance(),
			"flexibleengine_rds_instance_v3":                    resourceRdsInstanceV3(),
			"flexibleengine_rds_read_replica_v3":                resourceReplicaRdsInstance(),
			"flexibleengine_rds_parametergroup_v3":              resourceRdsConfigurationV3(),
			"flexibleengine_networking_vip_v2":                  resourceNetworkingVIPV2(),
			"flexibleengine_networking_vip_associate_v2":        resourceNetworkingVIPAssociateV2(),
			"flexibleengine_drs_replication_v2":                 resourceReplication(),
			"flexibleengine_drs_replicationconsistencygroup_v2": resourceReplicationConsistencyGroup(),
			"flexibleengine_nat_dnat_rule_v2":                   resourceNatDnatRuleV2(),
			"flexibleengine_nat_gateway_v2":                     resourceNatGatewayV2(),
			"flexibleengine_nat_snat_rule_v2":                   resourceNatSnatRuleV2(),
			"flexibleengine_vpc_eip_v1":                         resourceVpcEIPV1(),
			"flexibleengine_vpc_v1":                             resourceVirtualPrivateCloudV1(),
			"flexibleengine_vpc_subnet_v1":                      resourceVpcSubnetV1(),
			"flexibleengine_vpc_route_v2":                       resourceVPCRouteV2(),
			"flexibleengine_vpc_peering_connection_v2":          resourceVpcPeeringConnectionV2(),
			"flexibleengine_vpc_peering_connection_accepter_v2": resourceVpcPeeringConnectionAccepterV2(),
			"flexibleengine_sfs_file_system_v2":                 resourceSFSFileSystemV2(),
			"flexibleengine_sfs_access_rule_v2":                 resourceSFSAccessRuleV2(),
			"flexibleengine_sfs_turbo":                          resourceSFSTurbo(),
			"flexibleengine_rts_software_config_v1":             resourceSoftwareConfigV1(),
			"flexibleengine_rts_stack_v1":                       resourceRTSStackV1(),
			"flexibleengine_compute_bms_server_v2":              resourceComputeBMSInstanceV2(),
			"flexibleengine_ces_alarmrule":                      resourceAlarmRule(),
			"flexibleengine_dws_cluster_v1":                     resourceDWSClusterV1(),
			"flexibleengine_csbs_backup_v1":                     resourceCSBSBackupV1(),
			"flexibleengine_csbs_backup_policy_v1":              resourceCSBSBackupPolicyV1(),
			"flexibleengine_vbs_backup_policy_v2":               resourceVBSBackupPolicyV2(),
			"flexibleengine_vbs_backup_v2":                      resourceVBSBackupV2(),
			"flexibleengine_antiddos_v1":                        resourceAntiDdosV1(),
			"flexibleengine_css_cluster_v1":                     resourceCssClusterV1(),
			"flexibleengine_css_snapshot_v1":                    resourceCssSnapshotV1(),
			"flexibleengine_cts_tracker_v1":                     resourceCTSTrackerV1(),
			"flexibleengine_cce_node_v3":                        resourceCCENodeV3(),
			"flexibleengine_cce_cluster_v3":                     resourceCCEClusterV3(),
			"flexibleengine_dds_instance_v3":                    resourceDdsInstanceV3(),
			"flexibleengine_sdrs_drill_v1":                      resourceSdrsDrillV1(),
			"flexibleengine_sdrs_protectiongroup_v1":            resourceSdrsProtectiongroupV1(),
			"flexibleengine_sdrs_protectedinstance_v1":          resourceSdrsProtectedInstanceV1(),
			"flexibleengine_sdrs_replication_pair_v1":           resourceSdrsReplicationPairV1(),
			"flexibleengine_sdrs_replication_attach_v1":         resourceSdrsReplicationAttachV1(),
		},
	}

	provider.ConfigureFunc = func(d *schema.ResourceData) (interface{}, error) {
		terraformVersion := provider.TerraformVersion
		if terraformVersion == "" {
			// Terraform 0.12 introduced this field to the protocol
			// We can therefore assume that if it's missing it's 0.10 or 0.11
			terraformVersion = "0.11+compatible"
		}
		return configureProvider(d, terraformVersion)
	}

	return provider
}

var descriptions map[string]string

func init() {
	descriptions = map[string]string{
		"access_key": "The access key for API operations. You can retrieve this\n" +
			"from the 'My Credential' section of the console.",

		"secret_key": "The secret key for API operations. You can retrieve this\n" +
			"from the 'My Credential' section of the console.",

		"auth_url": "The Identity authentication URL.",

		"region": "The FlexibleEngine region to connect to.",

		"user_name": "Username to login with.",

		"user_id": "User ID to login with.",

		"tenant_id": "The ID of the Tenant (Identity v2) or Project (Identity v3)\n" +
			"to login with.",

		"tenant_name": "The name of the Tenant (Identity v2) or Project (Identity v3)\n" +
			"to login with.",

		"password": "Password to login with.",

		"token": "Authentication token to use as an alternative to username/password.",

		"security_token": "Security token to use for OBS federated authentication.",

		"domain_id": "The ID of the Domain to scope to (Identity v3).",

		"domain_name": "The name of the Domain to scope to (Identity v3).",

		"insecure": "Trust self-signed certificates.",

		"cacert_file": "A Custom CA certificate.",

		"endpoint_type": "The catalog endpoint type to use.",

		"cert": "A client certificate to authenticate with.",

		"key": "A client private key to authenticate with.",

		"swauth": "Use Swift's authentication system instead of Keystone. Only used for\n" +
			"interaction with Swift.",
	}
}

func configureProvider(d *schema.ResourceData, terraformVersion string) (interface{}, error) {
	config := Config{
		AccessKey:        d.Get("access_key").(string),
		SecretKey:        d.Get("secret_key").(string),
		CACertFile:       d.Get("cacert_file").(string),
		ClientCertFile:   d.Get("cert").(string),
		ClientKeyFile:    d.Get("key").(string),
		DomainID:         d.Get("domain_id").(string),
		DomainName:       d.Get("domain_name").(string),
		EndpointType:     d.Get("endpoint_type").(string),
		IdentityEndpoint: d.Get("auth_url").(string),
		Insecure:         d.Get("insecure").(bool),
		Password:         d.Get("password").(string),
		Region:           d.Get("region").(string),
		Swauth:           d.Get("swauth").(bool),
		Token:            d.Get("token").(string),
		SecurityToken:    d.Get("security_token").(string),
		TenantID:         d.Get("tenant_id").(string),
		TenantName:       d.Get("tenant_name").(string),
		Username:         d.Get("user_name").(string),
		UserID:           d.Get("user_id").(string),
		terraformVersion: terraformVersion,
	}

	if err := config.LoadAndValidate(); err != nil {
		return nil, err
	}

	return &config, nil
}
