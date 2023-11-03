# CHANGELOG

## 1.43.0 (November 3, 2023)

FEATURES:

* **New Resource:**
  - `flexibleengine_dli_global_variable` [GH-1034]
  - `flexibleengine_aom_service_discovery_rule` [GH-1036]

* **New Data Source:**
  - `flexibleengine_compute_servergroups` [GH-1037]
  - `flexibleengine_gaussdb_cassandra_flavors` [GH-1040]
  - `flexibleengine_cce_nodes` [GH-1060]

BUG FIXES:

* `resource/flexibleengine_cce_node_pool_v3`: `initial_node_count` should not trigger scale down [GH-975]

## 1.42.0 (September 28, 2023)

FEATURES:

* **New Resource:**
  - `flexibleengine_ddm_instance` [GH-998]
  - `flexibleengine_ddm_schema` [GH-998]
  - `flexibleengine_ddm_account` [GH-998]

* **New Data Source:**
  - `flexibleengine_ddm_accounts` [GH-998]
  - `flexibleengine_ddm_engines` [GH-998]
  - `flexibleengine_ddm_flavors` [GH-998]
  - `flexibleengine_ddm_instance_nodes` [GH-998]
  - `flexibleengine_ddm_instances` [GH-998]
  - `flexibleengine_ddm_schemas` [GH-998]

ENHANCEMENTS:

* `resource/flexibleengine_lb_loadbalancer_v3`: Support `autoscaling_enabled` and `min_l7_flavor_id` parameters [GH-1019]
* `resource/flexibleengine_cce_cluster`: Support cluster hibernation feature [GH-1022]

DEPRECATED:

* `resource/flexibleengine_mrs_cluster_v1` [GH-1028]
* `resource/flexibleengine_mrs_hybrid_cluster_v1` [GH-1028]
* `resource/flexibleengine_mrs_job_v1` [GH-1028]

## 1.41.0 (September 6, 2023)

FEATURES:

* **New Resource:**
  - `flexibleengine_nat_private_dnat_rule` [GH-971]
  - `flexibleengine_nat_private_gateway` [GH-971]
  - `flexibleengine_nat_private_snat_rule` [GH-971]
  - `flexibleengine_nat_private_transit_ip` [GH-971]

ENHANCEMENTS:

* `resource/flexibleengine_sfs_turbo`: Support enhanced feature  [GH-980]

## 1.40.0 (August 2, 2023)

FEATURES:

* **New Resource:**
  - `flexibleengine_kms_grant` [GH-965]
  - `flexibleengine_as_instance_attach` [GH-966]
  - `flexibleengine_as_notification` [GH-966]

* **New Data Source:**
  - `flexibleengine_as_configurations` [GH-966]
  - `flexibleengine_as_groups` [GH-966]

## 1.39.0 (July 7, 2023)

FEATURES:

* **New Resource:**
  - `flexibleengine_smn_message_template` [GH-958]
  - `flexibleengine_images_image` [GH-959]

* **New Data Source:**
  - `flexibleengine_images_image` [GH-959]
  - `flexibleengine_images_images` [GH-959]
  - `flexibleengine_cbr_backup` [GH-960]

ENHANCEMENTS:

* `resource/flexibleengine_cce_cluster_v3`: Support custom_san feature  [GH-950]
* `resource/flexibleengine_cce_node_v3`: Support kms_key_id feature [GH-954]
* `resource/flexibleengine_cce_node_pool_v3`: Support kms_key_id feature [GH-954]
* `data/flexibleengine_cce_cluster_v3`: Support custom_san feature  [GH-950]

DEPRECATED:

* `resource/flexibleengine_images_image_v2` [GH-959]
* `data/flexibleengine_images_image_v2` [GH-959]

## 1.38.0 (June 1, 2023)

FEATURES:

* **New Resource:**
  - `flexibleengine_obs_bucket_acl` [GH-943]
  - `flexibleengine_images_image_copy` [GH-944]
  - `flexibleengine_images_image_share` [GH-947]
  - `flexibleengine_images_image_share_accepter` [GH-947]

ENHANCEMENTS:

* `resource/flexibleengine_dds_instance_v3`: Support MongoDB 4.2 [GH-930]
* `resource/flexibleengine_obs_bucket`: Support parallel_fs feature [GH-942]

## 1.37.0 (April 14, 2023)

FEATURES:

* **New Resource:**
  - `flexibleengine_gaussdb_cassandra_instance` [GH-934]
  - `flexibleengine_gaussdb_influx_instance` [GH-934]

* **New Data Source:**
  - `flexibleengine_dms_kafka_instances` [GH-932]
  - `flexibleengine_gaussdb_cassandra_instances` [GH-934]
  - `flexibleengine_gaussdb_nosql_flavors` [GH-934]

ENHANCEMENTS:

* `provider`: Add security_token support [GH-931]

BUG FIXES:

* `resource/flexibleengine_waf_dedicated_domain`: Fix the issue that certificate can't be used in dedicated domain [GH-933]
* `resource/flexibleengine_smn_subscription_v2`: Delete the resource if it does not exist [GH-936]

## 1.36.1 (March 3, 2023)

FEATURES:

* **New Resource:**
  - `flexibleengine_obs_bucket_notifications` [GH-922]

ENHANCEMENTS:

* `resource/flexibleengine_dms_kafka_instance`: Add tags support [GH-925]

BUG FIXES:

* `resource/flexibleengine_dds_database_role`: Fix import issue when the role does not exist [GH-885]
* `resource/flexibleengine_dms_rocketmq_instance`: Fix ForceNew issue caused by the order of `availability_zones` [GH-906]
* `resource/flexibleengine_dms_rocketmq_user`: Fix import issue when the user does not exist [GH-907]

## 1.36.0 (January 13, 2023)

FEATURES:

* **New Resource:**
  - `flexibleengine_identity_acl` [GH-469]
  - `flexibleengine_sms_server_template` [GH-886]
  - `flexibleengine_sms_task` [GH-886]
  - `flexibleengine_dms_rocketmq_instance` [GH-880]
  - `flexibleengine_dms_rocketmq_consumer_group` [GH-897]
  - `flexibleengine_dms_rocketmq_topic` [GH-897]
  - `flexibleengine_dms_rocketmq_user` [GH-897]

* **New Data Source:**
  - `flexibleengine_sms_source_servers` [GH-886]
  - `flexibleengine_dws_flavors` [GH-894]
  - `flexibleengine_dms_rocketmq_instances` [GH-897]
  - `flexibleengine_dms_rocketmq_broker` [GH-897]

ENHANCEMENTS:

* `resource/flexibleengine_dws_cluster_v1`: support `public_ip` parameter [GH-891]
* `resource/flexibleengine_cce_node_pool_v3`: add `ecs_group_id` parameter [GH-899]
* `resource/flexibleengine_cce_node_v3`: add `subnet_id` and `ecs_group_id` parameters [GH-900]

BUG FIXES:

* `resource/flexibleengine_dms_kafka_instance`: change `availability_zones` to set [GH-895]

## 1.35.1 (December 19, 2022)

ENHANCEMENTS:

* `resource/flexibleengine_vpc_subnet_v1`: add `ipv4_subnet_id` and deprecate `subnet_id` attribute [GH-858]
* `resource/flexibleengine_networking_vip_v2`: use ip_version instead of subnet_id [GH-868]

BUG FIXES:

* `resource/flexibleengine_lts_group`: update the API version [GH-859]

DEPRECATED:

* `data.flexibleengine_networking_network_v2` [GH-869]
* `flexibleengine_networking_floatingip_associate_v2` [GH-869]
* `flexibleengine_networking_floatingip_v2` [GH-869]
* `flexibleengine_networking_network_v2` [GH-869]
* `flexibleengine_networking_subnet_v2` [GH-869]
* `flexibleengine_networking_router_interface_v2` [GH-869]
* `flexibleengine_networking_router_v2` [GH-869]

## 1.35.0 (December 1, 2022)

FEATURES:

* **New Resource:**
  - `flexibleengine_vpc_route_table` [GH-843]
  - `flexibleengine_vpc_route` [GH-843]
  - `flexibleengine_dli_database` [GH-844]
  - `flexibleengine_dli_table` [GH-844]
  - `flexibleengine_dli_package` [GH-844]
  - `flexibleengine_dli_flinksql_job` [GH-844]
  - `flexibleengine_dli_spark_job` [GH-844]
  - `flexibleengine_waf_dedicated_instance` [GH-846]
  - `flexibleengine_waf_dedicated_domain` [GH-846]
  - `flexibleengine_waf_dedicated_certificate` [GH-849]
  - `flexibleengine_waf_dedicated_policy` [GH-849]
  - `flexibleengine_lb_pool_v3` [GH-848]
  - `flexibleengine_lb_member_v3` [GH-848]
  - `flexibleengine_lb_monitor_v3` [GH-848]

* **New Data Source:**
  - `flexibleengine_vpc_route_table` [GH-843]
  - `flexibleengine_waf_dedicated_instances` [GH-846]

ENHANCEMENTS:

* `resource/flexibleengine_lb_pool_v2`: support setting persistence timeout [GH-821]
* `resource/flexibleengine_lb_listener_v3`: support advanced forwarding [GH-802]

BUG FIXES:

* `provider`: auth_url should not use region var [GH-831]

DEPRECATED:

* `flexibleengine_vpc_route_v2` [GH-845]
* `flexibleengine_networking_router_route_v2` [GH-845]

## 1.34.0 (October 29, 2022)

FEATURES:

* **New Resource:**
  - `flexibleengine_tms_tags` [GH-822]
  - `flexibleengine_rds_database_privilege` [GH-823]
  - `flexibleengine_dds_database_user` [GH-827]
  - `flexibleengine_dds_database_role` [GH-827]
  - `flexibleengine_dms_kafka_user` [GH-829]

ENHANCEMENTS:

* `resource/flexibleengine_cce_cluster_v3`: support binding or unbinding EIP without rebuild [GH-818]
* `resource/flexibleengine_dcs_instance_v1`: support port customization [GH-830]
* `resource/flexibleengine_dms_kafka_instance`: support enable_auto_topic in kafka instance [GH-832]
* `data/flexibleengine_dcs_product_v1`: support more filter parameters [GH-824]
* `data/flexibleengine_rds_flavors_v3`: support filtering flavors by group_type [GH-825]

## 1.33.0 (September 30, 2022)

FEATURES:

* **New Resource:**
  - `flexibleengine_modelarts_dataset` [GH-810]
  - `flexibleengine_modelarts_dataset_version` [GH-810]
  - `flexibleengine_drs_job` [GH-811]
  - `flexibleengine_apig_api_publishment` [GH-815]
  - `flexibleengine_apig_api` [GH-815]
  - `flexibleengine_apig_application` [GH-815]
  - `flexibleengine_apig_custom_authorizer` [GH-815]
  - `flexibleengine_apig_environment` [GH-815]
  - `flexibleengine_apig_group` [GH-815]
  - `flexibleengine_apig_instance` [GH-815]
  - `flexibleengine_apig_response` [GH-815]
  - `flexibleengine_apig_throttling_policy_associate` [GH-815]
  - `flexibleengine_apig_throttling_policy` [GH-815]
  - `flexibleengine_apig_vpc_channel` [GH-815]

* **New Data Source:**
  - `flexibleengine_smn_topics` [GH-806]
  - `flexibleengine_sfs_turbos` [GH-807]
  - `flexibleengine_modelarts_datasets` [GH-810]
  - `flexibleengine_modelarts_dataset_versions` [GH-810]
  - `flexibleengine_cce_clusters` [GH-812]
  - `flexibleengine_apig_environments` [GH-815]

BUG FIXES:

* `resource/flexibleengine_cce_cluster_v3`: set internal_endpoint and external_endpoint correctly [GH-808]
* `resource/flexibleengine_rts_stack_v1`: fix misuse of reflect.StringHeader [GH-813]

## 1.32.0 (August 16, 2022)

FEATURES:

* **New Resource:**
  - `flexibleengine_rds_account` [GH-784]
  - `flexibleengine_rds_database` [GH-784]
  - `flexibleengine_cce_namespace` [GH-785]
  - `flexibleengine_cce_pvc` [GH-785]
  - `flexibleengine_vpc_eip_associate` [GH-792]
  - `flexibleengine_cse_microservice` [GH-794]
  - `flexibleengine_cse_microservice_engine` [GH-794]
  - `flexibleengine_cse_microservice_instance` [GH-794]

* **New Data Source:**
  - `flexibleengine_networking_port` [GH-790]
  - `flexibleengine_identity_group` [GH-793]
  - `flexibleengine_identity_users` [GH-793]

ENHANCEMENTS:

* `resource/flexibleengine_dms_kafka_instance`: support manager user and password [GH-786]

## 1.31.1 (Jul 22, 2022)

BUG FIXES:

* `resource/flexibleengine_compute_instance_v2`: ignore tags validation [GH-781]

## 1.31.0 (June 30, 2022)

FEATURES:

* **New Resource:**
  - `flexibleengine_swr_organization` [GH-766]
  - `flexibleengine_swr_organization_users` [GH-766]
  - `flexibleengine_swr_repository` [GH-766]
  - `flexibleengine_swr_repository_sharing` [GH-766]
  - `flexibleengine_lb_loadbalancer_v3` [GH-770]
  - `flexibleengine_lb_listener_v3` [GH-772]
  - `flexibleengine_api_gateway_api` [GH-771]
  - `flexibleengine_api_gateway_group` [GH-771]
  - `flexibleengine_enterprise_project` [GH-775]
  - `flexibleengine_elb_certificate` [GH-777]
  - `flexibleengine_elb_ipgroup` [GH-777]

* **New Data Source:**
  - `flexibleengine_availability_zones` [GH-768]
  - `flexibleengine_elb_flavors` [GH-770]
  - `flexibleengine_elb_certificate` [GH-777]
  - `flexibleengine_enterprise_project` [GH-775]

ENHANCEMENTS:

* `resource/flexibleengine_vpc_subnet_v1`: support to enable IPv6 function [GH-769]
* `resource/flexibleengine_compute_instance_v2`: support to deploy ECS to a dedicated host [GH-773]
* `resource/flexibleengine_cce_node_v3`: support data volume encryption [GH-774]

DEPRECATED:

* flexibleengine_lb_certificate_v2 [GH-777]
* data/flexibleengine_compute_availability_zones_v2 [GH-768]
* data/flexibleengine_blockstorage_availability_zones_v3 [GH-768]

## 1.30.0 (June 15, 2022)

FEATURES:

* **New Resource:**
  - `flexibleengine_cbr_vaults` [GH-760]
  - `flexibleengine_cbr_policy` [GH-760]

* **New Data Source:**
  - `flexibleengine_cbr_vaults` [GH-760]

ENHANCEMENTS:

* `resource/flexibleengine_kms_key_v1`: support key rotation management [GH-702]
* `resource/flexibleengine_networking_secgroup_rule_v2`: support `description` field [GH-749]
* `resource/flexibleengine_nat_dnat_rule_v2`: support `description` field [GH-750]
* `resource/flexibleengine_compute_instance_v2`: add warnings when flavor is Xen-based [GH-758]

## 1.29.0 (April 29, 2022)

ENHANCEMENTS:

* `resource/flexibleengine_rds_instance_v3`: support enable SSL for MySQL [GH-744]
* `resource/flexibleengine_vpc_v1`: support description and secondary_cidr [GH-745]
* `resource/flexibleengine_as_group_v1`: support to forcibly delete an AS group [GH-746]

BUG FIXES:

* `resource/flexibleengine_cts_tracker_v1`: only filter CTS trckers by name [GH-738]
* `resource/flexibleengine_networking_secgroup_rule_v2`: add validation for remote_ip_prefix [GH-747]

## 1.28.0 (April 2, 2022)

FEATURES:

* **New Resource:**
  - `flexibleengine_fgs_function` [GH-732]
  - `flexibleengine_fgs_trigger` [GH-732]
  - `flexibleengine_fgs_dependency` [GH-732]

* **New Data Source:**
  - `flexibleengine_fgs_dependencies` [GH-732]

BUG FIXES:

* `resource/flexibleengine_cts_tracker_v1`: check length to avoid index out of range [GH-728]

## 1.27.1 (March 4, 2022)

BUG FIXES:

* Be able to remove a description of ELB loadbalancer and listener [GH-692]
* Do not update other fields when only tags was changed [GH-693]
* `resource/flexibleengine_obs_bucket`: update obs Location when region is not specified [GH-705]

## 1.27.0 (January 29, 2022)

FEATURES:

* **New Resource:**
  - `flexibleengine_dms_kafka_instance` [GH-682]
  - `flexibleengine_dms_kafka_topic` [GH-683]
  - `flexibleengine_obs_bucket_replication` [GH-688]

* **New Data Source:**
  - `flexibleengine_dms_product` [GH-681]

BUG FIXES:

* `resource/flexibleengine_lb_listener_v2`: avoid to request empty body when only update tags [GH-680]

## 1.26.0 (December 18, 2021)

FEATURES:

* **New Resource:**
  - `flexibleengine_identity_provider` [GH-553]
  - `flexibleengine_identity_provider_conversion` [GH-553]
  - `flexibleengine_dis_stream` [GH-620]
  - `flexibleengine_mrs_cluster_v2` [GH-659]

ENHANCEMENTS:

* `resource/flexibleengine_lb_listener_v2` - support to obtain client IP address [GH-609]
* `resource/flexibleengine_compute_keypair_v2` - support `private_key_path` for keypair [GH-665]
* `data/flexibleengine_dcs_product_v1` - support to filter by `cache_mode` [GH-655]

BUG FIXES:

* update endpoint of dnsV2Client [GH-664]

DEPRECATED:

* data/flexibleengine_dcs_az_v1 [GH-655]

## 1.25.1 (November 30, 2021)

FEATURES:

* **New Data Source:**
  - `flexibleengine_compute_instances` [GH-646]

BUG FIXES:

* `resource/flexibleengine_antiddos_v1`: throw an error when got a 403 response [GH-648]
* `resource/flexibleengine_cce_node_v3`: extend delay interval to get node ID [GH-653]

## 1.25.0 (October 28, 2021)

FEATURES:

* **New Resurce:**
  - `flexibleengine_identity_project_v3"` [GH-625]
  - `flexibleengine_mrs_job_v2` [GH-640]

* **New Data Source:**
  - `flexibleengine_nat_gateway_v2` [GH-629]
  - `flexibleengine_vpc_eip_v1` [GH-642]
  - `flexibleengine_dds_flavors_v3` [643]

## 1.24.2 (October 04, 2021)

ENHANCEMENTS:

* `data/flexibleengine_compute_bms_flavors_v2`: support to query an available BMS flavor by vcpus [GH-618]
* `resource/flexibleengine_compute_instance_v2`: add device_type and disk_bus fields [GH-622]

## 1.24.1 (September 18, 2021)

FEATURES:

* **New Resurce:**
  - `flexibleengine_cce_addon_v3` [GH-612]

* **New Data Source:**
  - `flexibleengine_cce_addon_template` [GH-612]
  - `flexibleengine_compute_flavors_v2` [GH-615]

ENHANCEMENTS:

* `data/flexibleengine_rds_flavors_v3`: Add vcpus and memory arguments to get RDS flavors [GH-616]

## 1.24.0 (September 10, 2021)

ENHANCEMENTS:

* **provider:** Upgrade to terraform-plugin-sdk v2 [GH-587]
* **config:** Add validation of domain name [GH-607]
* `resource/flexibleengine_identity_user_v3`: Add email and phone fields into identity_user_v3 [GH-474]
* `resource/flexibleengine_waf_policy`: Add ability to set/update protection_status block [GH-592]
* `resource/flexibleengine_dli_queue`: Support to scale out/in dli queues [GH-600]

BUG FIXES:

* `resource/flexibleengine_as_group_v1`: Support up to six load balancers can be added to an AS group [GH-524]
* `resource/flexibleengine_vpc_subnet_v1`: Fix subnet does not belong to the VPC error when deleting [GH-595]
* `resource/flexibleengine_obs_bucket`: Use proxy URL from environment variables by default in obs client [GH-602]

## 1.23.0 (July 27, 2021)

FEATURES:

* **New Resurce:**
  - `flexibleengine_lts_group` ([#583](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/pull/583))
  - `flexibleengine_lts_topic` ([#583](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/pull/583))
  - `flexibleengine_vpc_flow_log_v1` ([#451](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/451))
  - `flexibleengine_dli_queue` ([#588](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/pull/588))

BUG FIXES:

* `resource/flexibleengine_sfs_file_system_v2`: remove unused `host` attribute ([#420](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/420))

## 1.22.1 (July 9, 2021)

ENHANCEMENTS:

* **provider:** Try to request more times when the API call is return 429 error code ([#527](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/527))
* **provider:** Set default value for `tenant_name` and `auth_url` ([#569](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/569)))
* **provider:** Unset *Connection* header when sending API requests ([#577](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/pull/577))
* `resource/flexibleengine_vpc_eip_v1`: Add tags support ([#564](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/564))

BUG FIXES:

* `resource/flexibleengine_dds_instance_v3`: fix to set port attribute of DDS instance
 ([#571](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/571))

## 1.22.0 (June 28, 2021)

FEATURES:

* **New Resurce:**
  - `flexibleengine_waf_certificate` ([#533](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/533))
  - `flexibleengine_waf_domain` ([#533](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/533))
  - `flexibleengine_waf_policy` ([#533](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/533))
  - `flexibleengine_waf_rule_alarm_masking` ([#533](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/533))
  - `flexibleengine_waf_rule_blacklist` ([#533](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/533))
  - `flexibleengine_waf_rule_cc_protection` ([#533](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/533))
  - `flexibleengine_waf_rule_data_masking` ([#533](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/533))
  - `flexibleengine_waf_rule_precise_protection` ([#533](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/533))
  - `flexibleengine_waf_rule_web_tamper_protection` ([#533](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/533))

ENHANCEMENTS:

* **provider:** add max_retries to try more times when an API call is experiencing transient failures ([#566](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/pull/566))
* `flexibleengine_obs_bucket`: Support default encryption ([#419](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/419))
* `flexibleengine_obs_bucket`: Support to enable multi-AZ mode ([#505](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/505))
* `flexibleengine_cce_node_pool_v3`: Add ability to import ([#528](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/pull/528))
* `flexibleengine_cce_node_pool_v3`: Make taints and labels be updatable ([#549](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/549))
* `flexibleengine_rds_instance_v3`: Omit empty if port is not specified ([#550](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/pull/550))

BUG FIXES:

* `flexibleengine_rds_instance_v3`: set RegionProjectIDMap as early as possible in config
 ([#548](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/548))

## 1.21.0 (May 31, 2021)

FEATURES:

* **New Data Source:** `flexibleengine_lb_loadbalancer_v2` ([#476](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/476))
* **New Resurce:** `flexibleengine_as_lifecycle_hook_v1` ([#525](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/525))

ENHANCEMENTS:

* `flexibleengine_lb_loadbalancer_v2`: Add ability to import ([#538](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/pull/538))
* `flexibleengine_cce_node_pool_v3`: Support tags and max_pods ([#519](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/519))
* `flexibleengine_rds_instance_v3`: Add ability to update flavor, volume size and backup_strategy ([#522](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/522))

BUG FIXES:

* `flexibleengine_vpc_subnet_v1`: Allow DNS Nameservers to be cleared ([#544](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/pull/544))

DEPRECATED:

* resource/flexibleengine_rds_instance_v1: ([#542](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/pull/542))
* data/flexibleengine_rds_flavors_v1: ([#542](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/pull/542))

## 1.20.0 (April 30, 2021)

FEATURES:

* **New Data Source:** `flexibleengine_vpcep_endpoints` ([#514](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/pull/514))

ENHANCEMENTS:

* `flexibleengine_ces_alarmrule`: add alarm_level parameter ([#526](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/526))

BUG FIXES:

* `flexibleengine_dns_recordset_v2`: change records type to set ([#492](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/492))
* `flexibleengine_vpcep_approval`: vpcep approval can work cross-project ([#495](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/495))
* `flexibleengine_lb_monitor_v2`: fix crashes if monitor type is HTTP ([#517](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/517))

DEPRECATED:

* deprecate `swauth` in provider: ([#508](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/pull/508))
* deprecate `network/name` in `flexibleengine_compute_instance_v2`: ([#498](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/498))

## 1.19.1 (April 2, 2021)

ENHANCEMENTS:

* `flexibleengine_obs_bucket`: Add MaxItems in website block ([#491](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/491))
* Update structure of website document ([#501](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/501))

## 1.19.0 (March 9, 2021)

ENHANCEMENTS:

* `flexibleengine_identity_agency_v3`: Support Agency type "Cloud service" Management ([#409](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/409))
* `flexibleengine_rds_instance_v3`: Support `fixed_ip` parameter ([#477](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/pull/477))
* `flexibleengine_cce_cluster_v3`: support `authenticating_proxy_ca` and `kube_proxy_mode`([#484](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/pull/484))
* `flexibleengine_cce_node_v3`: Support `extend_param` to specify agency name ([#485](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/pull/485))
* `flexibleengine_dds_instance_v3`:
    - Support version 4.0 ([#473](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/473))
    - Support `tags` and `ssl` parameters ([#478](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/pull/478))

BUG FIXES:

* `flexibleengine_vpcep_service`: fix documentation error ([#479](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/479))
* `flexibleengine_lb_pool_v2`: Update available values of `protocol` ([#487](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/487))

## 1.18.1 (February 8, 2021)

ENHANCEMENTS:

* `flexibleengine_lb_listener_v2`: Support `http2_enable` parameter ([#466](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/pull/466))
* `flexibleengine_cce_cluster_v3`: Support `masters` parameter ([#468](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/pull/468))
* `flexibleengine_dcs_instance_v1`: Support Redis 4.0 and 5.0 instance ([#471](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/pull/471))

BUG FIXES:

* `flexibleengine_compute_servergroup_v2`: Support anti-affinity policy only ([#463](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/pull/463))
* `flexibleengine_sfs_file_system_v2`: Make access_type and access_level to be computed ([#470](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/pull/470))

## 1.18.0 (January 16, 2021)

FEATURES:

* **New Data Source:** `flexibleengine_vpcep_public_services` ([#334](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/334))
* **New Data Source:** `flexibleengine_identity_custom_role_v3` ([#461](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/461))
* **New Resource:** `flexibleengine_vpcep_service` ([#334](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/334))
* **New Resource:** `flexibleengine_vpcep_endpoint` ([#334](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/334))
* **New Resource:** `flexibleengine_vpcep_approval` ([#334](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/334))
* **New Resource:** `flexibleengine_dns_ptrrecord_v2` ([#441](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/pull/441))

ENHANCEMENTS:

* `flexibleengine_compute_instance_v2`: Improve on compute instance to import more attributes ([#459](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/pull/459))
* `flexibleengine_networking_port_v2`: Support import fixed_ip block ([#435](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/pull/435))
* `flexibleengine_sfs_turbo`: Add encryption on SFS Turbo ([#443](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/433))
* `flexibleengine_dns_recordset_v2`: Support CAA type ([#450](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/450))
* `flexibleengine_rds_instance_v3`: Support time_zome ([#457](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/pull/457))
* `flexibleengine_identity_role_v3`: Support custom IAM policy ([#408](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/408))
* **tags support:**
  - `flexibleengine_as_group_v1` ([#453](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/pull/453))
  - `flexibleengine_blockstorage_volume_v2` ([#454](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/pull/454))
  - `flexibleengine_dns_zone_v2` ([#439](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/pull/439))
  - `flexibleengine_dns_recordset_v2` ([#440](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/pull/440))
  - `flexibleengine_lb_loadbalancer_v2` ([#453](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/pull/453))
  - `flexibleengine_lb_listener_v2` ([#453](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/pull/453))
  - `flexibleengine_rds_instance_v3` ([#449](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/pull/449))
  - `flexibleengine_rds_read_replica_v3` ([#449](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/pull/449))
  - `flexibleengine_vpc_v1` ([#448](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/pull/448))
  - `flexibleengine_vpc_subnet_v1` ([#448](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/pull/448))

BUG FIXES:

* `flexibleengine_sfs_file_system_v2`: Fix output metadata is empty ([#421](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/421))
* `flexibleengine_cce_node_pool_v3`: Enable to create node pool with 0 nodes ([#429](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/429))
* `flexibleengine_kms_key_v1`: Fix docs issue ([#452](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/452))

## 1.17.0 (December 18, 2020)

FEATURES:

* **New Data Source:** `flexibleengine_compute_instance_v2` ([#424](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/424))
* **New Resource:** `flexibleengine_cce_node_pool_v3` ([#414](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/414))

ENHANCEMENTS:

* `resource/flexibleengine_cce_node_v3`: Add taints parameter support ([#412](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/412))
* `resource/flexibleengine_compute_instance_v2`: Add import support ([#425](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/425))

## 1.16.2 (December 03, 2020)

ENHANCEMENTS:

* `resource/flexibleengine_cce_cluster_v3`: Bump Create/Delete timeout ([#410](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/410))
* `resource/flexibleengine_cce_node_v3`: Bump Create/Delete timeout ([#410](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/410))

## 1.16.1 (October 19, 2020)

ENHANCEMENTS:

* `resource/flexibleengine_identity_role_assignment_v3`: Clean up unsupported user_id ([#393](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/393))
* `resource/flexibleengine_css_cluster_v1`: Make tags updatable ([#397](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/397))
* `resource/flexibleengine_cce_node_v3`: Add tags and labels support ([#398](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/398))
* `resource/flexibleengine_css_cluster_v1`: Add security_mode support ([#400](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/400))

BUG FIXES:

* `resource/flexibleengine_css_cluster_v1`: Fix enable backup issue ([#396](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/396))

## 1.16.0 (September 04, 2020)

FEATURES:

* **New Resource:** `flexibleengine_network_acl` ([#388](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/388))
* **New Resource:** `flexibleengine_network_acl_rule` ([#388](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/388))

## 1.15.0 (August 31, 2020)

FEATURES:

* **New Resource:** `flexibleengine_sfs_turbo` ([#384](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/384))

ENHANCEMENTS:

* `resource/flexibleengine_ces_alarmrule`: Improve docs and examples ([#382](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/382))
* `resource/flexibleengine_vbs_backup_policy`: Add resources parameter support ([#385](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/385))

## 1.14.0 (August 10, 2020)

FEATURES:

* **New Resource:** `flexibleengine_css_snapshot_v1` ([#368](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/368))
* **New Resource:** `flexibleengine_sfs_access_rule_v2` ([#377](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/377))

ENHANCEMENTS:

* `resource/flexibleengine_css_cluster_v1`: Make backup_strategy updatable ([#367](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/367))
* `resource/flexibleengine_sfs_file_system_v2`: Make access_to parameter optional ([#375](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/375))
* `resource/flexibleengine_cce_node_v3`: Add ability to import cce nodes ([#354](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/354))

BUG FIXES:

* `resource/flexibleengine_compute_instance_v2`: Fix multi nics issue ([#373](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/373))

## 1.13.0 (July 16, 2020)

FEATURES:

* **New Data Source:** `flexibleengine_identity_project_v3` ([#358](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/358))
* **New Data Source:** `flexibleengine_identity_role_v3` ([#358](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/358))
* **New Resource:** `flexibleengine_css_cluster_v1` ([#345](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/345))
* **New Resource:** `flexibleengine_identity_group_v3` ([#356](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/356))
* **New Resource:** `flexibleengine_identity_group_membership_v3` ([#356](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/356))
* **New Resource:** `flexibleengine_identity_user_v3` ([#356](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/356))
* **New Resource:** `flexibleengine_identity_role_v3` ([#358](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/358))
* **New Resource:** `flexibleengine_identity_role_assignment_v3` ([#358](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/358))
* **New Resource:** `flexibleengine_identity_agency_v3` ([#359](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/359))
* **New Resource:** `flexibleengine_mrs_hybrid_cluster_v1` ([#363](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/363))

ENHANCEMENTS:

* `resource/flexibleengine_compute_instance_v2`: Add tags support ([#344](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/344))
* Add sensitive flag to password parameters ([#346](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/346))
* `resource/flexibleengine_dds_instance_v3`: Add port and nodes attributes ([#348](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/348))
* `resource/flexibleengine_cce_cluster_v3`: Update supported version in docs ([#350](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/350))
* `resource/flexibleengine_dcs_instance_v1`: Add instance_type parameter support ([#360](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/360))

BUG FIXES:

* `resource/flexibleengine_obs_bucket`: Fix storage_class validation issue ([#338](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/338))
* `resource/flexibleengine_cce_node_v3`: Fix annotations and lables issue ([#341](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/341))
* `resource/flexibleengine_rds_instance_v3`: Fix docs issue for db.user_name ([#343](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/343))

## 1.12.1 (May 14, 2020)

BUG FIXES:

* `resource/flexibleengine_dds_instance_v3`: Fix instance create issue ([#330](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/330))

## 1.12.0 (April 30, 2020)

FEATURES:

* **New Data Source:** `flexibleengine_compute_availability_zones_v2` ([#311](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/311))
* **New Data Source:** `flexibleengine_blockstorage_availability_zones_v2` ([#316](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/316))
* **New Resource:** `flexibleengine_obs_bucket` ([#318](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/318))
* **New Resource:** `flexibleengine_obs_bucket_object` ([#318](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/318))
* **New Resource:** `flexibleengine_rds_read_replica_v3` ([#320](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/320))

ENHANCEMENTS:

* `resource/flexibleengine_cce_node_v3`: Add preinstall and postinstall support ([#310](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/310))
* `resource/flexibleengine_nat_snat_rule_v2`: Add preinstall and postinstall script support ([#310](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/310))

BUG FIXES:

* `resource/flexibleengine_vpc_eip_v1`: Update max EIP bandwidth ([#309](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/309))
* `data source/flexibleengine_networking_network_v2`: Catch API request error ([#313](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/313))
* `resource/flexibleengine_mrs_cluster_v1`: Fix cluster_type empty issue ([#325](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/325))

## 1.11.1 (February 18, 2020)

ENHANCEMENTS:

* `resource/flexibleengine_sdrs_protectiongroup_v1`: Add protection enable/disable support ([#304](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/304))
* `resource/flexibleengine_cce_node_v3`: Add preinstall/postinstall script support ([#305](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/305))

## 1.11.0 (February 07, 2020)

FEATURES:

* **New Resource:** `flexibleengine_sdrs_replication_attach_v1` ([#300](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/300))

BUG FIXES:

* `resource/flexibleengine_nat_gateway_v2`: Fix argument reference for a better user experience ([#291](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/291))
* Fix provider docs with token issue ([#298](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/298))

## 1.10.0 (January 22, 2020)

FEATURES:

* **New Data Source:** `flexibleengine_lb_certificate_v2` ([#281](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/281))
* **New Data Source:** `flexibleengine_sdrs_domain_v1` ([#287](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/287))
* **New Resource:** `flexibleengine_sdrs_protectiongroup_v1` ([#287](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/287))
* **New Resource:** `flexibleengine_sdrs_replication_pair_v1` ([#289](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/289))
* **New Resource:** `flexibleengine_sdrs_protectedinstance_v1` ([#290](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/290))
* **New Resource:** `flexibleengine_sdrs_drill_v1` ([#292](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/292))

ENHANCEMENTS:

* Add security_token to OBS federated authentication ([#278](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/278))
* `resource/flexibleengine_lb_listener_v2`: Add tls_ciphers_policy support ([#282](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/282))
* `resource/flexibleengine_cce_node_v3`: Add Computed to annotations and labels attributes ([#286](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/286))

## 1.9.0 (October 31, 2019)

FEATURES:

* **New Resource:** `flexibleengine_lb_whitelist_v2` ([#266](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/266))

ENHANCEMENTS:

* `resource/flexibleengine_compute_instance_v2`: Log fault message when building compute instance failed ([#261](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/261))
* `resource/flexibleengine_dcs_instance_v1`: Rename subnet_id parameter to network_id ([#265](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/265))
* `resource/flexibleengine_cce_cluster_v3`: Add certificates information to cce cluster ([#268](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/268))
* `resource/flexibleengine_cce_cluster_v3`: Add eip support to cce cluster ([#269](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/269))

BUG FIXES:

* `resource/flexibleengine_dcs_instance_v1`: Fix ip/port attribute issue ([#267](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/267))

## 1.8.0 (September 30, 2019)

FEATURES:

* **New Data Source:** `flexibleengine_dds_flavors_v3` ([#246](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/246))
* **New Data Source:** `flexibleengine_cce_node_ids_v3` ([#248](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/248))
* **New Resource:** `flexibleengine_dds_instance_v3` ([#245](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/245))

ENHANCEMENTS:

* `resource/flexibleengine_cce_cluster_v3`: Add authentication mode support ([#234](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/234))
* `resource/flexibleengine_cce_node_v3`: Add os parameter support ([#238](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/238))
* `resource/flexibleengine_cce_cluster_v3`: Add security_group_id attribute ([#239](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/239))

BUG FIXES:

* `resource/flexibleengine_cce_cluster_v3`: Fix cluster default version issue ([#235](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/235))
* `resource/flexibleengine_dcs_instance_v1`: Fix DCS parameters issue for single node ([#241](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/241))
* `resource/flexibleengine_lb_monitor_v2`: Fix lb monitor expected_codes issue ([#249](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/249))
* `resource/flexibleengine_rds_instance_v3`: Fix RDS instance db version issue ([#250](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/250))

## 1.7.0 (August 29, 2019)

FEATURES:

* **New Data Source:** `flexibleengine_blockstorage_volume_v2` ([#192](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/192))
* **New Data Source:** `flexibleengine_rds_flavors_v3` ([#209](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/209))
* **New Resource:** `flexibleengine_rds_instance_v3` ([#209](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/209))
* **New Resource:** `flexibleengine_rds_parametergroup_v3` ([#209](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/209))

ENHANCEMENTS:

* `resource/flexibleengine_blockstorage_volume_v2`: Add multiattach support ([#195](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/195))
* `resource/flexibleengine_cce_cluster_v3`: Add endpoints to CCE cluster ([#202](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/202))
* `resource/flexibleengine_cce_node_v3`: Add public_ip/private_ip to CCE node ([#205](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/205))
* `resource/flexibleengine_networking_floatingip_v2`: Add default value for floating_ip pool ([#211](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/211))
* `resource/flexibleengine_lb_listener_v2`: Add example for lb_listener with TERMINATED_HTTPS ([#217](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/217))
* Add detailed error message for 404 ([#225](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/225))

BUG FIXES:

* `resource/flexibleengine_csbs_backup_policy_v1`: Fix CSBS backup policy name issue ([#201](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/201))
* `datasource/flexibleengine_cce_node_v3`: Fix cce_node datasource naming issue ([#222](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/222))
* `resource/flexibleengine_cce_node_v3`: Fix data_volumes type issue ([#229](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/229))
* `resource/flexibleengine_vpc_subnet_v1`: Fix dns_list type issue ([#230](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/230))

## 1.6.0 (June 05, 2019)

FEATURES:

* **New Data Source:** `flexibleengine_dns_zone_v2` ([#190](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/190))

ENHANCEMENTS:

* The provider is now compatible with Terraform v0.12, while retaining compatibility with prior versions.
* `resource/flexibleengine_lb_monitor_v2`: Add health port option support. ([#189](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/189))

## 1.5.0 (April 30, 2019)

FEATURES:

* **New Data Source:** `flexibleengine_kms_key_v1` ([#149](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/149))
* **New Data Source:** `flexibleengine_kms_data_key_v1` ([#149](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/149))
* **New Resource:** `flexibleengine_lb_l7policy_v2` ([#114](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/114))
* **New Resource:** `flexibleengine_lb_l7rule_v2` ([#114](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/114))
* **New Resource:** `flexibleengine_kms_key_v1` ([#149](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/149))
* **New Resource:** `flexibleengine_compute_interface_attach_v2` ([#152](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/152))
* **New Resource:** `flexibleengine_nat_dnat_rule_v2` ([#153](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/153))

BUG FIXES:

* `data_source/flexibleengine_cce_cluster_v3`: Remove wrong attributes internal, external, and external_otc. ([#119](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/119))
* `resource/flexibleengine_smn_topic_v2`: Fix SMN topic parameters issue. ([#126](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/126))
* `resource/flexibleengine_dcs_instance_v1`: Fix DCS instance parameters issue. ([#128](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/128))
* `resource/flexibleengine_cce_node_v3`: Remove Abnormal from create target state. ([#168](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/168))
* `resource/flexibleengine_lb_pool_v2`: Fix LB Pool stuck issue. ([#169](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/169))

ENHANCEMENTS:

* `resource/flexibleengine_dns_zone_v2`: Add support for attaching multi routers to dns zone. ([#143](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/143))
* `resource/flexibleengine_blockstorage_volume_v2`: Add volume extending support. ([#156](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/156))
* `resource/flexibleengine_compute_instance_v2`: Add auto_recovery support. ([#163](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/163))
* `resource/flexibleengine_as_group_v1`: Add lbaas_listeners support. ([#172](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/172))

## 1.4.0 (January 18, 2019)

FEATURES:

* **New Data Source:** `flexibleengine_cce_node_v3` ([#105](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/105))
* **New Data Source:** `flexibleengine_cce_cluster_v3` ([#105](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/105))
* **New Resource:** `flexibleengine_cce_node_v3` ([#105](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/105))
* **New Resource:** `flexibleengine_cce_cluster_v3` ([#105](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/105))

BUG FIXES:

* `resource/flexibleengine_dns_recordset_v2`: Fix dns records update error ([#101](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/101))
* `resource/flexibleengine_dns_recordset_v2`: Fix dns entries re-sort issue ([#103](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/103))

## 1.3.1 (January 08, 2019)

BUG FIXES:

* Fix ak/sk authentication issue ([#102](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/102))

## 1.3.0 (January 07, 2019)

FEATURES:

* **New Data Source:** `flexibleengine_cts_tracker_v1` ([#64](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/64))
* **New Data Source:** `flexibleengine_dcs_az_v1` ([#76](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/76))
* **New Data Source:** `flexibleengine_dcs_maintainwindow_v1` ([#76](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/76))
* **New Data Source:** `flexibleengine_dcs_product_v1` ([#76](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/76))
* **New Resource:** `flexibleengine_cts_tracker_v1` ([#64](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/64))
* **New Resource:** `flexibleengine_antiddos_v1` ([#66](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/66))
* **New Resource:** `flexibleengine_dcs_instance_v1` ([#76](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/76))
* **New Resource:** `flexibleengine_networking_floatingip_associate_v2` ([#83](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/83))

BUG FIXES:

* `resource/flexibleengine_vpc_subnet_v1`: Remove UNKNOWN status to avoid error ([#73](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/73))
* `resource/flexibleengine_rds_instance_v1`: Add PostgreSQL support ([#81](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/81))
* `resource/flexibleengine_rds_instance_v1`: Suppress rds name change ([#82](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/82))
* `resource/flexibleengine_smn_topic_v2`: Fix smn update error ([#84](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/84))
* `resource/flexibleengine_elb_listener`: Add check for elb listener certificate_id ([#85](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/85))
* `all resources`: Expose real error message of BadRequest error ([#91](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/91))
* `resource/flexibleengine_sfs_file_system_v2`: Suppress sfs system metadata ([#98](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/98))

ENHANCEMENTS:

* `resource/flexibleengine_networking_router_v2`: Add enable_snat option support ([#97](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/97))

## 1.2.1 (October 29, 2018)

BUG FIXES:

* `resource/flexibleengine_as_configuration_v1`: Fix AutoScaling client error ([#60](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/60))

## 1.2.0 (October 01, 2018)

FEATURES:

* **New Data Source:** `flexibleengine_images_image_v2` ([#20](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/20))
* **New Data Source:** `flexibleengine_sfs_file_system_v2` ([#23](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/23))
* **New Data Source:** `flexibleengine_compute_bms_flavor_v2` ([#26](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/26))
* **New Data Source:** `flexibleengine_compute_bms_keypair_v2` ([#26](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/26))
* **New Data Source:** `flexibleengine_compute_bms_nic_v2` ([#26](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/26))
* **New Data Source:** `flexibleengine_compute_bms_server_v2` ([#26](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/26))
* **New Data Source:** `flexibleengine_rts_software_config_v1` ([#28](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/28))
* **New Data Source:** `flexibleengine_rts_stack_v1` ([#28](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/28))
* **New Data Source:** `flexibleengine_rts_stack_resource_v1` ([#28](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/28))
* **New Data Source:** `flexibleengine_csbs_backup_v1` ([#49](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/49))
* **New Data Source:** `flexibleengine_csbs_backup_policy_v1` ([#49](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/49))
* **New Data Source:** `flexibleengine_vbs_backup_policy_v2` ([#54](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/54))
* **New Data Source:** `flexibleengine_vbs_backup_v2` ([#54](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/54))
* **New Resource:** `flexibleengine_images_image_v2` ([#20](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/20))
* **New Resource:** `flexibleengine_vpc_eip_v1` ([#21](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/21))
* **New Resource:** `flexibleengine_lb_loadbalancer_v2` ([#22](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/22))
* **New Resource:** `flexibleengine_lb_listener_v2` ([#22](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/22))
* **New Resource:** `flexibleengine_lb_pool_v2` ([#22](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/22))
* **New Resource:** `flexibleengine_lb_member_v2` ([#22](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/22))
* **New Resource:** `flexibleengine_lb_monitor_v2` ([#22](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/22))
* **New Resource:** `flexibleengine_sfs_file_system_v2` ([#23](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/23))
* **New Resource:** `flexibleengine_rts_software_config_v1` ([#28](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/28))
* **New Resource:** `flexibleengine_rts_stack_v1` ([#28](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/28))
* **New Resource:** `flexibleengine_ces_alarmrule` ([#29](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/29))
* **New Resource:** `flexibleengine_fw_firewall_group_v2` ([#30](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/30))
* **New Resource:** `flexibleengine_fw_policy_v2` ([#30](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/30))
* **New Resource:** `flexibleengine_fw_rule_v2` ([#30](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/30))
* **New Resource:** `flexibleengine_compute_bms_server_v2` ([#31](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/31))
* **New Resource:** `flexibleengine_mrs_cluster_v1` ([#36](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/36))
* **New Resource:** `flexibleengine_mrs_job_v1` ([#36](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/36))
* **New Resource:** `flexibleengine_mls_instance_v1` ([#44](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/44))
* **New Resource:** `flexibleengine_dws_cluster_v1` ([#47](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/47))
* **New Resource:** `flexibleengine_lb_certificate_v2` ([#48](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/48))
* **New Resource:** `flexibleengine_csbs_backup_v1` ([#49](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/49))
* **New Resource:** `flexibleengine_csbs_backup_policy_v1` ([#49](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/49))
* **New Resource:** `flexibleengine_vbs_backup_policy_v2` ([#54](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/54))
* **New Resource:** `flexibleengine_vbs_backup_v2` ([#54](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/54))

ENHANCEMENTS:

* resource/flexibleengine_vpc_subnet_v1: Add `subnet_id` parameter ([#19](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/19))
* provider: Add AK/SK authentication support ([#35](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/35))

## 1.1.0 (July 20, 2018)

FEATURES:

* **New Data Source:** `flexibleengine_vpc_v1` ([#15](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/15))
* **New Data Source:** `flexibleengine_vpc_subnet_v1` ([#15](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/15))
* **New Data Source:** `flexibleengine_vpc_subnet_ids_v1` ([#15](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/15))
* **New Data Source:** `flexibleengine_vpc_route_v2` ([#15](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/15))
* **New Data Source:** `flexibleengine_vpc_route_ids_v2` ([#15](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/15))
* **New Data Source:** `flexibleengine_vpc_peering_connection_v2` ([#15](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/15))
* **New Resource:** `flexibleengine_drs_replication_v2` ([#13](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/13))
* **New Resource:** `flexibleengine_drs_replicationconsistencygroup_v2` ([#13](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/13))
* **New Resource:** `flexibleengine_networking_vip_v2` ([#13](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/13))
* **New Resource:** `flexibleengine_networking_vip_associate_v2` ([#13](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/13))
* **New Resource:** `flexibleengine_nat_gateway_v2` ([#14](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/14))
* **New Resource:** `flexibleengine_nat_snat_rule_v2` ([#14](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/14))
* **New Resource:** `flexibleengine_vpc_v1` ([#15](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/15))
* **New Resource:** `flexibleengine_vpc_subnet_v1` ([#15](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/15))
* **New Resource:** `flexibleengine_vpc_route_v2` ([#15](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/15))
* **New Resource:** `flexibleengine_vpc_peering_connection_v2` ([#15](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/15))
* **New Resource:** `flexibleengine_vpc_peering_connection_accepter_v2` ([#15](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/15))

## 1.0.1 (June 08, 2018)

BUG FIXES:

* resource/flexibleengine_elb_backend: Correct ELB Backend parameter names ([#7](https://github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/issues/7))

## 1.0.0 (June 01, 2018)

FEATURES:

* **New Data Source:** `flexibleengine_networking_network_v2`
* **New Data Source:** `flexibleengine_networking_secgroup_v2`
* **New Data Source:** `flexibleengine_s3_bucket_object`
* **New Data Source:** `flexibleengine_rds_flavors_v1`
* **New Resource:** `flexibleengine_blockstorage_volume_v2`
* **New Resource:** `flexibleengine_compute_instance_v2`
* **New Resource:** `flexibleengine_compute_keypair_v2`
* **New Resource:** `flexibleengine_compute_servergroup_v2`
* **New Resource:** `flexibleengine_compute_floatingip_v2`
* **New Resource:** `flexibleengine_compute_volume_attach_v2`
* **New Resource:** `flexibleengine_dns_recordset_v2`
* **New Resource:** `flexibleengine_dns_zone_v2`
* **New Resource:** `flexibleengine_networking_network_v2`
* **New Resource:** `flexibleengine_networking_subnet_v2`
* **New Resource:** `flexibleengine_networking_floatingip_v2`
* **New Resource:** `flexibleengine_networking_port_v2`
* **New Resource:** `flexibleengine_networking_router_v2`
* **New Resource:** `flexibleengine_networking_router_interface_v2`
* **New Resource:** `flexibleengine_networking_router_route_v2`
* **New Resource:** `flexibleengine_networking_secgroup_v2`
* **New Resource:** `flexibleengine_networking_secgroup_rule_v2`
* **New Resource:** `flexibleengine_s3_bucket`
* **New Resource:** `flexibleengine_s3_bucket_policy`
* **New Resource:** `flexibleengine_s3_bucket_object`
* **New Resource:** `flexibleengine_elb_loadbalancer`
* **New Resource:** `flexibleengine_elb_listener`
* **New Resource:** `flexibleengine_elb_backend`
* **New Resource:** `flexibleengine_elb_health`
* **New Resource:** `flexibleengine_as_group_v1`
* **New Resource:** `flexibleengine_as_configuration_v1`
* **New Resource:** `flexibleengine_as_policy_v1`
* **New Resource:** `flexibleengine_smn_topic_v2`
* **New Resource:** `flexibleengine_smn_subscription_v2`
* **New Resource:** `flexibleengine_rds_instance_v1`
