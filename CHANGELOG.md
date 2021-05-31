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
