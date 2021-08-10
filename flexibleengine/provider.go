package flexibleengine

import (
	"context"
	"fmt"
	"sync"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/mutexkv"
)

const (
	defaultCloud     = "prod-cloud-ocb.orange-business.com"
	terraformVersion = "0.12+compatible"
)

// This is a global MutexKV for use within this plugin.
var osMutexKV = mutexkv.NewMutexKV()

// Provider returns a schema.Provider for FlexibleEngine.
func Provider() *schema.Provider {
	provider := &schema.Provider{
		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Required:    true,
				Description: descriptions["region"],
				DefaultFunc: schema.EnvDefaultFunc("OS_REGION_NAME", nil),
			},

			"access_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: descriptions["access_key"],
				DefaultFunc: schema.EnvDefaultFunc("OS_ACCESS_KEY", nil),
			},

			"secret_key": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  descriptions["secret_key"],
				RequiredWith: []string{"access_key"},
				DefaultFunc:  schema.EnvDefaultFunc("OS_SECRET_KEY", nil),
			},

			"domain_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: descriptions["domain_id"],
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"OS_USER_DOMAIN_ID",
					"OS_PROJECT_DOMAIN_ID",
					"OS_DOMAIN_ID",
				}, ""),
			},

			"domain_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: descriptions["domain_name"],
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"OS_USER_DOMAIN_NAME",
					"OS_PROJECT_DOMAIN_NAME",
					"OS_DOMAIN_NAME",
					"OS_DEFAULT_DOMAIN",
				}, ""),
			},

			"user_id": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_USER_ID", ""),
				Description: descriptions["user_name"],
			},

			"user_name": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_USERNAME", ""),
				Description: descriptions["user_name"],
			},

			"password": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("OS_PASSWORD", ""),
				Description: descriptions["password"],
			},

			"tenant_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: descriptions["tenant_id"],
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"OS_TENANT_ID",
					"OS_PROJECT_ID",
				}, ""),
			},

			"tenant_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: descriptions["tenant_name"],
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"OS_TENANT_NAME",
					"OS_PROJECT_NAME",
				}, ""),
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

			"auth_url": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: descriptions["auth_url"],
				DefaultFunc: schema.EnvDefaultFunc("OS_AUTH_URL", nil),
			},

			"insecure": {
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_INSECURE", false),
				Description: descriptions["insecure"],
			},

			"max_retries": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: descriptions["max_retries"],
				DefaultFunc: schema.EnvDefaultFunc("OS_MAX_RETRIES", 5),
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

			"endpoint_type": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_ENDPOINT_TYPE", nil),
				Deprecated:  "endpoint_type is deprecated",
				ValidateFunc: validation.StringInSlice([]string{
					"public", "publicURL", "admin", "adminURL", "internal", "internalURL",
				}, false),
			},
		},

		DataSourcesMap: map[string]*schema.Resource{
			"flexibleengine_blockstorage_availability_zones_v3": dataSourceBlockStorageAvailabilityZonesV3(),
			"flexibleengine_blockstorage_volume_v2":             dataSourceBlockStorageVolumeV2(),
			"flexibleengine_compute_availability_zones_v2":      dataSourceComputeAvailabilityZonesV2(),
			"flexibleengine_compute_instance_v2":                dataSourceComputeInstance(),
			"flexibleengine_images_image_v2":                    dataSourceImagesImageV2(),
			"flexibleengine_networking_network_v2":              dataSourceNetworkingNetworkV2(),
			"flexibleengine_networking_secgroup_v2":             dataSourceNetworkingSecGroupV2(),
			"flexibleengine_s3_bucket_object":                   dataSourceS3BucketObject(),
			"flexibleengine_kms_key_v1":                         dataSourceKmsKeyV1(),
			"flexibleengine_kms_data_key_v1":                    dataSourceKmsDataKeyV1(),
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
			"flexibleengine_lb_loadbalancer_v2":                 dataSourceELBV2Loadbalancer(),
			"flexibleengine_sdrs_domain_v1":                     dataSourceSdrsDomainV1(),
			"flexibleengine_identity_project_v3":                dataSourceIdentityProjectV3(),
			"flexibleengine_identity_role_v3":                   dataSourceIdentityRoleV3(),
			"flexibleengine_identity_custom_role_v3":            dataSourceIdentityCustomRoleV3(),
			"flexibleengine_vpcep_public_services":              dataSourceVPCEPPublicServices(),
			"flexibleengine_vpcep_endpoints":                    dataSourceVPCEPEndpoints(),

			// Deprecated data source
			"flexibleengine_rds_flavors_v1": dataSourceRdsFlavorV1(),
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
			"flexibleengine_dns_ptrrecord_v2":                   resourceDNSPtrRecordV2(),
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
			"flexibleengine_network_acl":                        resourceNetworkACL(),
			"flexibleengine_network_acl_rule":                   resourceNetworkACLRule(),
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
			"flexibleengine_lts_group":                          resourceLTSGroupV2(),
			"flexibleengine_lts_topic":                          resourceLTSTopicV2(),
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
			"flexibleengine_as_lifecycle_hook_v1":               resourceASLifecycleHook(),
			"flexibleengine_smn_topic_v2":                       resourceTopic(),
			"flexibleengine_smn_subscription_v2":                resourceSubscription(),
			"flexibleengine_rds_instance_v3":                    resourceRdsInstanceV3(),
			"flexibleengine_rds_read_replica_v3":                resourceRdsReadReplicaInstance(),
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
			"flexibleengine_vpc_flow_log_v1":                    resourceVpcFlowLogV1(),
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
			"flexibleengine_cce_node_pool_v3":                   resourceCCENodePool(),
			"flexibleengine_dds_instance_v3":                    resourceDdsInstanceV3(),
			"flexibleengine_sdrs_drill_v1":                      resourceSdrsDrillV1(),
			"flexibleengine_sdrs_protectiongroup_v1":            resourceSdrsProtectiongroupV1(),
			"flexibleengine_sdrs_protectedinstance_v1":          resourceSdrsProtectedInstanceV1(),
			"flexibleengine_sdrs_replication_pair_v1":           resourceSdrsReplicationPairV1(),
			"flexibleengine_sdrs_replication_attach_v1":         resourceSdrsReplicationAttachV1(),
			"flexibleengine_vpcep_approval":                     resourceVPCEndpointApproval(),
			"flexibleengine_vpcep_endpoint":                     resourceVPCEndpoint(),
			"flexibleengine_vpcep_service":                      resourceVPCEndpointService(),
			"flexibleengine_waf_certificate":                    resourceWafCertificateV1(),
			"flexibleengine_waf_domain":                         resourceWafDomainV1(),
			"flexibleengine_waf_policy":                         resourceWafPolicyV1(),
			"flexibleengine_waf_rule_blacklist":                 resourceWafRuleBlackList(),
			"flexibleengine_waf_rule_alarm_masking":             resourceWafRuleAlarmMasking(),
			"flexibleengine_waf_rule_data_masking":              resourceWafRuleDataMasking(),
			"flexibleengine_waf_rule_cc_protection":             resourceWafRuleCCAttackProtection(),
			"flexibleengine_waf_rule_precise_protection":        resourceWafRulePreciseProtection(),
			"flexibleengine_waf_rule_web_tamper_protection":     resourceWafRuleWebTamperProtection(),
			"flexibleengine_dli_queue":                          ResourceDliQueueV1(),
			// Deprecated resource
			"flexibleengine_rds_instance_v1": resourceRdsInstance(),
		},
		// configuring the provider
		ConfigureContextFunc: configureProvider,
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

		"max_retries": "How many times HTTP connection should be retried until giving up.",

		"cacert_file": "A Custom CA certificate.",

		"endpoint_type": "The catalog endpoint type to use.",

		"cert": "A client certificate to authenticate with.",

		"key": "A client private key to authenticate with.",
	}
}

func configureProvider(_ context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	config := Config{
		EndpointType:  d.Get("endpoint_type").(string),
		SecurityToken: d.Get("security_token").(string),
	}

	region := d.Get("region").(string)
	config.Region = region

	config.TenantID = d.Get("tenant_id").(string)
	config.TenantName = d.Get("tenant_name").(string)
	// set tenant_name to region when neither `tenant_name` nor `tenant_id` was specified
	if config.TenantID == "" && config.TenantName == "" {
		config.TenantName = region
	}

	if v, ok := d.GetOk("auth_url"); ok {
		config.IdentityEndpoint = v.(string)
	} else {
		config.IdentityEndpoint = fmt.Sprintf("https://iam.%s.%s/v3", region, defaultCloud)
	}

	config.DomainID = d.Get("domain_id").(string)
	config.DomainName = d.Get("domain_name").(string)
	config.UserID = d.Get("user_id").(string)
	config.Username = d.Get("user_name").(string)
	config.Password = d.Get("password").(string)
	config.AccessKey = d.Get("access_key").(string)
	config.SecretKey = d.Get("secret_key").(string)
	config.Token = d.Get("token").(string)

	config.MaxRetries = d.Get("max_retries").(int)
	config.Insecure = d.Get("insecure").(bool)
	config.CACertFile = d.Get("cacert_file").(string)
	config.ClientCertFile = d.Get("cert").(string)
	config.ClientKeyFile = d.Get("key").(string)
	config.TerraformVersion = terraformVersion
	config.Cloud = defaultCloud
	config.RegionClient = true
	config.RegionProjectIDMap = make(map[string]string)
	config.RPLock = new(sync.Mutex)

	if err := config.LoadAndValidate(); err != nil {
		return nil, diag.FromErr(err)
	}

	return &config, nil
}
