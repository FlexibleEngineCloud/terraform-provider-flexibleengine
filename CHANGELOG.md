## 1.2.0 (October 01, 2018)

FEATURES:

* **New Data Source:** `flexibleengine_images_image_v2` ([#20](https://github.com/terraform-providers/terraform-provider-flexibleengine/issues/20))
* **New Data Source:** `flexibleengine_sfs_file_system_v2` ([#23](https://github.com/terraform-providers/terraform-provider-flexibleengine/issues/23))
* **New Data Source:** `flexibleengine_compute_bms_flavor_v2` ([#26](https://github.com/terraform-providers/terraform-provider-flexibleengine/issues/26))
* **New Data Source:** `flexibleengine_compute_bms_keypair_v2` ([#26](https://github.com/terraform-providers/terraform-provider-flexibleengine/issues/26))
* **New Data Source:** `flexibleengine_compute_bms_nic_v2` ([#26](https://github.com/terraform-providers/terraform-provider-flexibleengine/issues/26))
* **New Data Source:** `flexibleengine_compute_bms_server_v2` ([#26](https://github.com/terraform-providers/terraform-provider-flexibleengine/issues/26))
* **New Data Source:** `flexibleengine_rts_software_config_v1` ([#28](https://github.com/terraform-providers/terraform-provider-flexibleengine/issues/28))
* **New Data Source:** `flexibleengine_rts_stack_v1` ([#28](https://github.com/terraform-providers/terraform-provider-flexibleengine/issues/28))
* **New Data Source:** `flexibleengine_rts_stack_resource_v1` ([#28](https://github.com/terraform-providers/terraform-provider-flexibleengine/issues/28))
* **New Data Source:** `flexibleengine_csbs_backup_v1` ([#49](https://github.com/terraform-providers/terraform-provider-flexibleengine/issues/49))
* **New Data Source:** `flexibleengine_csbs_backup_policy_v1` ([#49](https://github.com/terraform-providers/terraform-provider-flexibleengine/issues/49))
* **New Data Source:** `flexibleengine_vbs_backup_policy_v2` ([#54](https://github.com/terraform-providers/terraform-provider-flexibleengine/issues/54))
* **New Data Source:** `flexibleengine_vbs_backup_v2` ([#54](https://github.com/terraform-providers/terraform-provider-flexibleengine/issues/54))
* **New Resource:** `flexibleengine_images_image_v2` ([#20](https://github.com/terraform-providers/terraform-provider-flexibleengine/issues/20))
* **New Resource:** `flexibleengine_vpc_eip_v1` ([#21](https://github.com/terraform-providers/terraform-provider-flexibleengine/issues/21))
* **New Resource:** `flexibleengine_lb_loadbalancer_v2` ([#22](https://github.com/terraform-providers/terraform-provider-flexibleengine/issues/22))
* **New Resource:** `flexibleengine_lb_listener_v2` ([#22](https://github.com/terraform-providers/terraform-provider-flexibleengine/issues/22))
* **New Resource:** `flexibleengine_lb_pool_v2` ([#22](https://github.com/terraform-providers/terraform-provider-flexibleengine/issues/22))
* **New Resource:** `flexibleengine_lb_member_v2` ([#22](https://github.com/terraform-providers/terraform-provider-flexibleengine/issues/22))
* **New Resource:** `flexibleengine_lb_monitor_v2` ([#22](https://github.com/terraform-providers/terraform-provider-flexibleengine/issues/22))
* **New Resource:** `flexibleengine_sfs_file_system_v2` ([#23](https://github.com/terraform-providers/terraform-provider-flexibleengine/issues/23))
* **New Resource:** `flexibleengine_rts_software_config_v1` ([#28](https://github.com/terraform-providers/terraform-provider-flexibleengine/issues/28))
* **New Resource:** `flexibleengine_rts_stack_v1` ([#28](https://github.com/terraform-providers/terraform-provider-flexibleengine/issues/28))
* **New Resource:** `flexibleengine_ces_alarmrule` ([#29](https://github.com/terraform-providers/terraform-provider-flexibleengine/issues/29))
* **New Resource:** `flexibleengine_fw_firewall_group_v2` ([#30](https://github.com/terraform-providers/terraform-provider-flexibleengine/issues/30))
* **New Resource:** `flexibleengine_fw_policy_v2` ([#30](https://github.com/terraform-providers/terraform-provider-flexibleengine/issues/30))
* **New Resource:** `flexibleengine_fw_rule_v2` ([#30](https://github.com/terraform-providers/terraform-provider-flexibleengine/issues/30))
* **New Resource:** `flexibleengine_compute_bms_server_v2` ([#31](https://github.com/terraform-providers/terraform-provider-flexibleengine/issues/31))
* **New Resource:** `flexibleengine_mrs_cluster_v1` ([#36](https://github.com/terraform-providers/terraform-provider-flexibleengine/issues/36))
* **New Resource:** `flexibleengine_mrs_job_v1` ([#36](https://github.com/terraform-providers/terraform-provider-flexibleengine/issues/36))
* **New Resource:** `flexibleengine_mls_instance_v1` ([#44](https://github.com/terraform-providers/terraform-provider-flexibleengine/issues/44))
* **New Resource:** `flexibleengine_dws_cluster_v1` ([#47](https://github.com/terraform-providers/terraform-provider-flexibleengine/issues/47))
* **New Resource:** `flexibleengine_lb_certificate_v2` ([#48](https://github.com/terraform-providers/terraform-provider-flexibleengine/issues/48))
* **New Resource:** `flexibleengine_csbs_backup_v1` ([#49](https://github.com/terraform-providers/terraform-provider-flexibleengine/issues/49))
* **New Resource:** `flexibleengine_csbs_backup_policy_v1` ([#49](https://github.com/terraform-providers/terraform-provider-flexibleengine/issues/49))
* **New Resource:** `flexibleengine_vbs_backup_policy_v2` ([#54](https://github.com/terraform-providers/terraform-provider-flexibleengine/issues/54))
* **New Resource:** `flexibleengine_vbs_backup_v2` ([#54](https://github.com/terraform-providers/terraform-provider-flexibleengine/issues/54))

ENHANCEMENTS:

* resource/flexibleengine_vpc_subnet_v1: Add `subnet_id` parameter ([#19](https://github.com/terraform-providers/terraform-provider-flexibleengine/issues/19))
* provider: Add AK/SK authentication support ([#35](https://github.com/terraform-providers/terraform-provider-flexibleengine/issues/35))

## 1.1.0 (July 20, 2018)

FEATURES:

* **New Data Source:** `flexibleengine_vpc_v1` ([#15](https://github.com/terraform-providers/terraform-provider-flexibleengine/issues/15))
* **New Data Source:** `flexibleengine_vpc_subnet_v1` ([#15](https://github.com/terraform-providers/terraform-provider-flexibleengine/issues/15))
* **New Data Source:** `flexibleengine_vpc_subnet_ids_v1` ([#15](https://github.com/terraform-providers/terraform-provider-flexibleengine/issues/15))
* **New Data Source:** `flexibleengine_vpc_route_v2` ([#15](https://github.com/terraform-providers/terraform-provider-flexibleengine/issues/15))
* **New Data Source:** `flexibleengine_vpc_route_ids_v2` ([#15](https://github.com/terraform-providers/terraform-provider-flexibleengine/issues/15))
* **New Data Source:** `flexibleengine_vpc_peering_connection_v2` ([#15](https://github.com/terraform-providers/terraform-provider-flexibleengine/issues/15))
* **New Resource:** `flexibleengine_drs_replication_v2` ([#13](https://github.com/terraform-providers/terraform-provider-flexibleengine/issues/13))
* **New Resource:** `flexibleengine_drs_replicationconsistencygroup_v2` ([#13](https://github.com/terraform-providers/terraform-provider-flexibleengine/issues/13))
* **New Resource:** `flexibleengine_networking_vip_v2` ([#13](https://github.com/terraform-providers/terraform-provider-flexibleengine/issues/13))
* **New Resource:** `flexibleengine_networking_vip_associate_v2` ([#13](https://github.com/terraform-providers/terraform-provider-flexibleengine/issues/13))
* **New Resource:** `flexibleengine_nat_gateway_v2` ([#14](https://github.com/terraform-providers/terraform-provider-flexibleengine/issues/14))
* **New Resource:** `flexibleengine_nat_snat_rule_v2` ([#14](https://github.com/terraform-providers/terraform-provider-flexibleengine/issues/14))
* **New Resource:** `flexibleengine_vpc_v1` ([#15](https://github.com/terraform-providers/terraform-provider-flexibleengine/issues/15))
* **New Resource:** `flexibleengine_vpc_subnet_v1` ([#15](https://github.com/terraform-providers/terraform-provider-flexibleengine/issues/15))
* **New Resource:** `flexibleengine_vpc_route_v2` ([#15](https://github.com/terraform-providers/terraform-provider-flexibleengine/issues/15))
* **New Resource:** `flexibleengine_vpc_peering_connection_v2` ([#15](https://github.com/terraform-providers/terraform-provider-flexibleengine/issues/15))
* **New Resource:** `flexibleengine_vpc_peering_connection_accepter_v2` ([#15](https://github.com/terraform-providers/terraform-provider-flexibleengine/issues/15))

## 1.0.1 (June 08, 2018)

BUG FIXES:

* resource/flexibleengine_elb_backend: Correct ELB Backend parameter names ([#7](https://github.com/terraform-providers/terraform-provider-flexibleengine/issues/7))

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
