package flexibleengine

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/mutexkv"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/apig"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/cbr"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/cce"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/cse"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/drs"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/eip"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/elb"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/eps"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/fgs"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/iam"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/modelarts"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/rds"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/sfs"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/smn"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/swr"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/vpc"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

const (
	defaultCloud     = "prod-cloud-ocb.orange-business.com"
	terraformVersion = "0.12+compatible"
)

// This is a global MutexKV for use within this plugin.
var osMutexKV = mutexkv.NewMutexKV()

func init() {
	utils.PackageName = "FlexibleEngine"
	cse.DefaultVersion = "CSE"
}

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
				Description: descriptions["user_id"],
			},

			"user_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: descriptions["user_name"],
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"OS_USER_NAME",
					"OS_USERNAME",
				}, ""),
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
				DefaultFunc: schema.EnvDefaultFunc("OS_SECURITY_TOKEN", nil),
			},

			"auth_url": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: descriptions["auth_url"],
				DefaultFunc: schema.EnvDefaultFunc("OS_AUTH_URL", nil),
			},

			"cloud": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: descriptions["cloud"],
				DefaultFunc: schema.EnvDefaultFunc("OS_CLOUD", defaultCloud),
			},

			"endpoints": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Experimental Feature: the custom endpoints used to override the default endpoint URL",
				Elem:        &schema.Schema{Type: schema.TypeString},
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
		},

		DataSourcesMap: map[string]*schema.Resource{
			"flexibleengine_availability_zones":        dataSourceAvailabilityZones(),
			"flexibleengine_blockstorage_volume_v2":    dataSourceBlockStorageVolumeV2(),
			"flexibleengine_compute_instance_v2":       dataSourceComputeInstance(),
			"flexibleengine_compute_instances":         dataSourceComputeInstances(),
			"flexibleengine_compute_flavors_v2":        dataSourceEcsFlavors(),
			"flexibleengine_images_image_v2":           dataSourceImagesImageV2(),
			"flexibleengine_networking_network_v2":     dataSourceNetworkingNetworkV2(),
			"flexibleengine_networking_secgroup_v2":    dataSourceNetworkingSecGroupV2(),
			"flexibleengine_s3_bucket_object":          dataSourceS3BucketObject(),
			"flexibleengine_kms_key_v1":                dataSourceKmsKeyV1(),
			"flexibleengine_kms_data_key_v1":           dataSourceKmsDataKeyV1(),
			"flexibleengine_rds_flavors_v3":            dataSourceRdsFlavorV3(),
			"flexibleengine_vpc_v1":                    dataSourceVirtualPrivateCloudVpcV1(),
			"flexibleengine_vpc_subnet_v1":             dataSourceVpcSubnetV1(),
			"flexibleengine_vpc_subnet_ids_v1":         dataSourceVpcSubnetIdsV1(),
			"flexibleengine_vpc_route_v2":              dataSourceVPCRouteV2(),
			"flexibleengine_vpc_route_ids_v2":          dataSourceVPCRouteIdsV2(),
			"flexibleengine_vpc_peering_connection_v2": dataSourceVpcPeeringConnectionV2(),
			"flexibleengine_vpc_eip":                   dataSourceVpcEipV1(),
			"flexibleengine_nat_gateway_v2":            dataSourceNatGatewayV2(),
			"flexibleengine_sfs_file_system_v2":        dataSourceSFSFileSystemV2(),
			"flexibleengine_compute_bms_flavors_v2":    dataSourceBMSFlavorV2(),
			"flexibleengine_compute_bms_nic_v2":        dataSourceBMSNicV2(),
			"flexibleengine_compute_bms_server_v2":     dataSourceBMSServersV2(),
			"flexibleengine_compute_bms_keypairs_v2":   dataSourceBMSKeyPairV2(),
			"flexibleengine_rts_software_config_v1":    dataSourceRtsSoftwareConfigV1(),
			"flexibleengine_rts_stack_v1":              dataSourceRTSStackV1(),
			"flexibleengine_rts_stack_resource_v1":     dataSourceRTSStackResourcesV1(),
			"flexibleengine_csbs_backup_v1":            dataSourceCSBSBackupV1(),
			"flexibleengine_csbs_backup_policy_v1":     dataSourceCSBSBackupPolicyV1(),
			"flexibleengine_vbs_backup_policy_v2":      dataSourceVBSBackupPolicyV2(),
			"flexibleengine_vbs_backup_v2":             dataSourceVBSBackupV2(),
			"flexibleengine_cts_tracker_v1":            dataSourceCTSTrackerV1(),
			"flexibleengine_dcs_maintainwindow_v1":     dataSourceDcsMaintainWindowV1(),
			"flexibleengine_dcs_product_v1":            dataSourceDcsProductV1(),
			"flexibleengine_dms_product":               dataSourceDmsProduct(),
			"flexibleengine_cce_node_v3":               dataSourceCceNodesV3(),
			"flexibleengine_cce_node_ids_v3":           dataSourceCceNodeIdsV3(),
			"flexibleengine_cce_cluster_v3":            dataSourceCCEClusterV3(),
			"flexibleengine_cce_addon_template":        dataSourceCCEAddonTemplate(),
			"flexibleengine_dns_zone_v2":               dataSourceDNSZoneV2(),
			"flexibleengine_dds_flavors_v3":            dataSourceDDSFlavorsV3(),
			"flexibleengine_lb_certificate_v2":         dataSourceCertificateV2(),
			"flexibleengine_lb_loadbalancer_v2":        dataSourceELBV2Loadbalancer(),
			"flexibleengine_sdrs_domain_v1":            dataSourceSdrsDomainV1(),
			"flexibleengine_identity_project_v3":       dataSourceIdentityProjectV3(),
			"flexibleengine_identity_role_v3":          dataSourceIdentityRoleV3(),
			"flexibleengine_identity_custom_role_v3":   dataSourceIdentityCustomRoleV3(),
			"flexibleengine_vpcep_public_services":     dataSourceVPCEPPublicServices(),
			"flexibleengine_vpcep_endpoints":           dataSourceVPCEPEndpoints(),
			"flexibleengine_elb_flavors":               dataSourceElbFlavorsV3(),

			// importing data source
			"flexibleengine_apig_environments":  apig.DataSourceEnvironments(),
			"flexibleengine_enterprise_project": eps.DataSourceEnterpriseProject(),
			"flexibleengine_cbr_vaults":         cbr.DataSourceCbrVaultsV3(),
			"flexibleengine_cce_clusters":       cce.DataSourceCCEClusters(),
			"flexibleengine_elb_certificate":    elb.DataSourceELBCertificateV3(),
			"flexibleengine_fgs_dependencies":   fgs.DataSourceFunctionGraphDependencies(),
			"flexibleengine_networking_port":    vpc.DataSourceNetworkingPortV2(),
			"flexibleengine_identity_group":     iam.DataSourceIdentityGroup(),
			"flexibleengine_identity_users":     iam.DataSourceIdentityUsers(),
			"flexibleengine_sfs_turbos":         sfs.DataSourceTurbos(),
			"flexibleengine_smn_topics":         smn.DataSourceTopics(),

			"flexibleengine_modelarts_datasets":         modelarts.DataSourceDatasets(),
			"flexibleengine_modelarts_dataset_versions": modelarts.DataSourceDatasetVerions(),

			// Deprecated data source
			"flexibleengine_compute_availability_zones_v2":      dataSourceAvailabilityZones(),
			"flexibleengine_blockstorage_availability_zones_v3": dataSourceBlockStorageAvailabilityZonesV3(),

			"flexibleengine_vpc_eip_v1": dataSourceVpcEipV1(),

			"flexibleengine_dcs_az_v1":      dataSourceDcsAZV1(),
			"flexibleengine_dds_flavor_v3":  dataSourceDDSFlavorV3(),
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
			"flexibleengine_dms_kafka_instance":                 resourceDmsKafkaInstances(),
			"flexibleengine_dms_kafka_topic":                    resourceDmsKafkaTopic(),
			"flexibleengine_dis_stream":                         resourceDisStreamV2(),
			"flexibleengine_fw_firewall_group_v2":               resourceFWFirewallGroupV2(),
			"flexibleengine_fw_policy_v2":                       resourceFWPolicyV2(),
			"flexibleengine_fw_rule_v2":                         resourceFWRuleV2(),
			"flexibleengine_images_image_v2":                    resourceImagesImageV2(),
			"flexibleengine_kms_key_v1":                         resourceKmsKeyV1(),
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
			"flexibleengine_mrs_cluster_v2":                     resourceMRSClusterV2(),
			"flexibleengine_mrs_job_v1":                         resourceMRSJobV1(),
			"flexibleengine_mrs_job_v2":                         resourceMRSJobV2(),
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
			"flexibleengine_identity_project_v3":                resourceIdentityProjectV3(),
			"flexibleengine_identity_role_v3":                   resourceIdentityRoleV3(),
			"flexibleengine_identity_role_assignment_v3":        resourceIdentityRoleAssignmentV3(),
			"flexibleengine_identity_user_v3":                   resourceIdentityUserV3(),
			"flexibleengine_identity_provider":                  resourceIdentityProvider(),
			"flexibleengine_identity_provider_conversion":       resourceIAMProviderConversion(),
			"flexibleengine_lts_group":                          resourceLTSGroupV2(),
			"flexibleengine_lts_topic":                          resourceLTSTopicV2(),
			"flexibleengine_s3_bucket":                          resourceS3Bucket(),
			"flexibleengine_s3_bucket_policy":                   resourceS3BucketPolicy(),
			"flexibleengine_s3_bucket_object":                   resourceS3BucketObject(),
			"flexibleengine_obs_bucket":                         resourceObsBucket(),
			"flexibleengine_obs_bucket_object":                  resourceObsBucketObject(),
			"flexibleengine_obs_bucket_replication":             resourceObsBucketReplication(),
			"flexibleengine_as_group_v1":                        resourceASGroup(),
			"flexibleengine_as_configuration_v1":                resourceASConfiguration(),
			"flexibleengine_as_policy_v1":                       resourceASPolicy(),
			"flexibleengine_as_lifecycle_hook_v1":               resourceASLifecycleHook(),
			"flexibleengine_smn_topic_v2":                       resourceTopic(),
			"flexibleengine_smn_subscription_v2":                resourceSubscription(),
			"flexibleengine_rds_read_replica_v3":                resourceRdsReadReplicaInstance(),
			"flexibleengine_rds_parametergroup_v3":              resourceRdsConfigurationV3(),
			"flexibleengine_networking_vip_v2":                  resourceNetworkingVIPV2(),
			"flexibleengine_networking_vip_associate_v2":        resourceNetworkingVIPAssociateV2(),
			"flexibleengine_drs_replication_v2":                 resourceReplication(),
			"flexibleengine_drs_replicationconsistencygroup_v2": resourceReplicationConsistencyGroup(),
			"flexibleengine_nat_dnat_rule_v2":                   resourceNatDnatRuleV2(),
			"flexibleengine_nat_gateway_v2":                     resourceNatGatewayV2(),
			"flexibleengine_nat_snat_rule_v2":                   resourceNatSnatRuleV2(),
			"flexibleengine_vpc_eip":                            resourceVpcEIPV1(),
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
			"flexibleengine_cce_cluster_v3":                     resourceCCEClusterV3(),
			"flexibleengine_cce_node_v3":                        resourceCCENodeV3(),
			"flexibleengine_cce_node_pool_v3":                   resourceCCENodePool(),
			"flexibleengine_cce_addon_v3":                       resourceCCEAddon(),
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

			// importing resource
			"flexibleengine_apig_api":                         apig.ResourceApigAPIV2(),
			"flexibleengine_apig_api_publishment":             apig.ResourceApigApiPublishment(),
			"flexibleengine_apig_instance":                    apig.ResourceApigInstanceV2(),
			"flexibleengine_apig_application":                 apig.ResourceApigApplicationV2(),
			"flexibleengine_apig_custom_authorizer":           apig.ResourceApigCustomAuthorizerV2(),
			"flexibleengine_apig_environment":                 apig.ResourceApigEnvironmentV2(),
			"flexibleengine_apig_group":                       apig.ResourceApigGroupV2(),
			"flexibleengine_apig_response":                    apig.ResourceApigResponseV2(),
			"flexibleengine_apig_throttling_policy_associate": apig.ResourceThrottlingPolicyAssociate(),
			"flexibleengine_apig_throttling_policy":           apig.ResourceApigThrottlingPolicyV2(),
			"flexibleengine_apig_vpc_channel":                 apig.ResourceApigVpcChannelV2(),

			"flexibleengine_api_gateway_api":   huaweicloud.ResourceAPIGatewayAPI(),
			"flexibleengine_api_gateway_group": huaweicloud.ResourceAPIGatewayGroup(),

			"flexibleengine_enterprise_project":        eps.ResourceEnterpriseProject(),
			"flexibleengine_cbr_policy":                cbr.ResourceCBRPolicyV3(),
			"flexibleengine_cbr_vault":                 cbr.ResourceVault(),
			"flexibleengine_cce_namespace":             cce.ResourceCCENamespaceV1(),
			"flexibleengine_cce_pvc":                   cce.ResourceCcePersistentVolumeClaimsV1(),
			"flexibleengine_cse_microservice":          cse.ResourceMicroservice(),
			"flexibleengine_cse_microservice_engine":   cse.ResourceMicroserviceEngine(),
			"flexibleengine_cse_microservice_instance": cse.ResourceMicroserviceInstance(),
			"flexibleengine_drs_job":                   drs.ResourceDrsJob(),
			"flexibleengine_fgs_dependency":            fgs.ResourceFgsDependency(),
			"flexibleengine_fgs_function":              fgs.ResourceFgsFunctionV2(),
			"flexibleengine_fgs_trigger":               fgs.ResourceFunctionGraphTrigger(),
			"flexibleengine_rds_account":               rds.ResourceRdsAccount(),
			"flexibleengine_rds_database":              rds.ResourceRdsDatabase(),
			"flexibleengine_rds_instance_v3":           rds.ResourceRdsInstance(),
			"flexibleengine_swr_organization":          swr.ResourceSWROrganization(),
			"flexibleengine_swr_organization_users":    swr.ResourceSWROrganizationPermissions(),
			"flexibleengine_swr_repository":            swr.ResourceSWRRepository(),
			"flexibleengine_swr_repository_sharing":    swr.ResourceSWRRepositorySharing(),

			"flexibleengine_vpc_v1":        vpc.ResourceVirtualPrivateCloudV1(),
			"flexibleengine_vpc_subnet_v1": vpc.ResourceVpcSubnetV1(),

			"flexibleengine_vpc_eip_associate": eip.ResourceEIPAssociate(),

			"flexibleengine_lb_loadbalancer_v3": elb.ResourceLoadBalancerV3(),
			"flexibleengine_lb_listener_v3":     elb.ResourceListenerV3(),
			"flexibleengine_elb_certificate":    elb.ResourceCertificateV3(),
			"flexibleengine_elb_ipgroup":        elb.ResourceIpGroupV3(),

			"flexibleengine_modelarts_dataset":         modelarts.ResourceDataset(),
			"flexibleengine_modelarts_dataset_version": modelarts.ResourceDatasetVersion(),

			// Deprecated resource
			"flexibleengine_elb_loadbalancer":  resourceELoadBalancer(),
			"flexibleengine_elb_listener":      resourceEListener(),
			"flexibleengine_elb_backend":       resourceBackend(),
			"flexibleengine_elb_health":        resourceHealth(),
			"flexibleengine_lb_certificate_v2": resourceCertificateV2(),
			"flexibleengine_rds_instance_v1":   resourceRdsInstance(),
			"flexibleengine_vpc_eip_v1":        resourceVpcEIPV1(),
		},
		// configuring the provider
		ConfigureContextFunc: configureProvider,
	}

	return provider
}

var descriptions map[string]string

func init() {
	descriptions = map[string]string{
		"access_key": "The access key of the FlexibleEngine cloud to use.",

		"secret_key": "The secret key of the FlexibleEngine cloud to use.",

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

		"cert": "A client certificate to authenticate with.",

		"key": "A client private key to authenticate with.",

		"cloud": "The endpoint of cloud provider, defaults to prod-cloud-ocb.orange-business.com",
	}
}

func configureProvider(_ context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	config := Config{}

	region := d.Get("region").(string)
	cloud := d.Get("cloud").(string)
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
		config.IdentityEndpoint = fmt.Sprintf("https://iam.%s.%s/v3", region, cloud)
	}

	config.DomainID = d.Get("domain_id").(string)
	config.DomainName = d.Get("domain_name").(string)
	config.UserID = d.Get("user_id").(string)
	config.Username = d.Get("user_name").(string)
	config.Password = d.Get("password").(string)
	config.AccessKey = d.Get("access_key").(string)
	config.SecretKey = d.Get("secret_key").(string)
	config.SecurityToken = d.Get("security_token").(string)
	config.Token = d.Get("token").(string)

	config.MaxRetries = d.Get("max_retries").(int)
	config.Insecure = d.Get("insecure").(bool)
	config.CACertFile = d.Get("cacert_file").(string)
	config.ClientCertFile = d.Get("cert").(string)
	config.ClientKeyFile = d.Get("key").(string)
	config.TerraformVersion = terraformVersion
	config.Cloud = cloud
	config.RegionClient = true
	config.RegionProjectIDMap = make(map[string]string)
	config.RPLock = new(sync.Mutex)

	// get custom endpoints
	endpoints, err := flattenProviderEndpoints(d)
	if err != nil {
		return nil, diag.FromErr(err)
	}

	// set default endpoints
	if _, ok := endpoints["obs"]; !ok {
		endpoints["obs"] = fmt.Sprintf("https://oss.%s.%s/", region, config.Cloud)
	}
	if _, ok := endpoints["fgs"]; !ok {
		endpoints["fgs"] = fmt.Sprintf("https://fgs.%s.%s/", region, config.Cloud)
	}
	if _, ok := endpoints["dns"]; !ok {
		endpoints["dns"] = fmt.Sprintf("https://dns.%s/", config.Cloud)
	}
	if _, ok := endpoints["eps"]; !ok {
		endpoints["eps"] = fmt.Sprintf("https://eps.%s/", config.Cloud)
	}

	config.Endpoints = endpoints
	if err := LoadAndValidate(&config); err != nil {
		return nil, diag.FromErr(err)
	}
	return &config, nil
}

func flattenProviderEndpoints(d *schema.ResourceData) (map[string]string, error) {
	endpoints := d.Get("endpoints").(map[string]interface{})
	epMap := make(map[string]string)

	for key, val := range endpoints {
		endpoint := strings.TrimSpace(val.(string))
		// check empty string
		if endpoint == "" {
			return nil, fmt.Errorf("the value of customer endpoint %s must be specified", key)
		}

		// add prefix "https://" and suffix "/"
		if !strings.HasPrefix(endpoint, "http") {
			endpoint = fmt.Sprintf("https://%s", endpoint)
		}
		if !strings.HasSuffix(endpoint, "/") {
			endpoint = fmt.Sprintf("%s/", endpoint)
		}
		epMap[key] = endpoint
	}

	// unify the endpoint which has multiple versions
	for key := range endpoints {
		ep, ok := epMap[key]
		if !ok {
			continue
		}

		multiKeys := config.GetServiceDerivedCatalogKeys(key)
		for _, k := range multiKeys {
			epMap[k] = ep
		}
	}

	log.Printf("[DEBUG] customer endpoints: %+v", epMap)
	return epMap, nil
}
