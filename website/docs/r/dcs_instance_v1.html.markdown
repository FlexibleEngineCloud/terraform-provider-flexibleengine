---
layout: "flexibleengine"
page_title: "flexibleengine: flexibleengine_dcs_instance_v1"
sidebar_current: "docs-flexibleengine-resource-dcs-instance-v1"
description: |-
  Manages a DCS instance in the flexibleengine DCS Service
---

# flexibleengine\_dcs\_instance_v1

Manages a DCS instance in the flexibleengine DCS Service.

## Example Usage

### DCS instance for Redis 3.0

```hcl
resource "flexibleengine_networking_secgroup_v2" "secgroup_1" {
  name        = "secgroup_1"
  description = "secgroup_1"
}

resource "flexibleengine_vpc_v1" "vpc_1" {
  name = "test_vpc1"
  cidr = "192.168.0.0/16"
}

resource "flexibleengine_vpc_subnet_v1" "subnet_1" {
  name       = "test_subnet1"
  cidr       = "192.168.0.0/24"
  gateway_ip = "192.168.0.1"
  vpc_id     = flexibleengine_vpc_v1.vpc_1.id
}

resource "flexibleengine_dcs_instance_v1" "instance_1" {
  name              = "test_dcs_instance"
  engine            = "Redis"
  engine_version    = "3.0"
  password          = var.my_password
  product_id        = "dcs.master_standby-h"
  capacity          = 2
  vpc_id            = flexibleengine_vpc_v1.vpc_1.id
  network_id        = flexibleengine_vpc_subnet_v1.subnet_1.id
  security_group_id = flexibleengine_networking_secgroup_v2.secgroup_1.id
  available_zones   = ["eu-west-0a"]
  save_days         = 1
  backup_type       = "manual"
  begin_at          = "00:00-01:00"
  period_type       = "weekly"
  backup_at         = [1]
}
```

### DCS instance for Redis 5.0

```hcl
resource "flexibleengine_dcs_instance_v1" "instance_1" {
  name              = "test_dcs_instance"
  engine            = "Redis"
  engine_version    = "5.0"
  password          = var.my_password
  product_id        = "redis.cluster.xu1.large.r1.8-h"
  capacity          = 8
  vpc_id            = flexibleengine_vpc.vpc_1.id
  subnet_id         = flexibleengine_vpc_subnet.subnet_1.id
  available_zones   = ["eu-west-0a", "eu-west-0b"]
  save_days         = 1
  backup_type       = "manual"
  begin_at          = "00:00-01:00"
  period_type       = "weekly"
  backup_at         = [1]
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Indicates the name of an instance. An instance name starts with a letter,
	consists of 4 to 64 characters, and supports only letters, digits, and hyphens (-).

* `description` - (Optional) Indicates the description of an instance. It is a character
    string containing not more than 1024 characters.

* `engine` - (Required) Indicates a cache engine. Only Redis is supported.
    Changing this creates a new instance.

* `engine_version` - (Required) Indicates the version of a cache engine, which is 3.0.7.
    Changing this creates a new instance.

* `capacity` - (Required) Indicates the Cache capacity. Unit: GB.
    For a DCS Redis or Memcached instance in single-node or master/standby mode, the cache
    capacity can be 2 GB, 4 GB, 8 GB, 16 GB, 32 GB, or 64 GB.
    For a DCS Redis instance in cluster mode, the cache capacity can be 64, 128, 256, 512,
    or 1024 GB. Changing this creates a new instance.

* `access_user` - (Optional) Username used for accessing a DCS instance after password
    authentication. A username starts with a letter, consists of 1 to 64 characters,
    and supports only letters, digits, and hyphens (-).
    Changing this creates a new instance.

* `password` - (Required) Password of a DCS instance.
    The password of a DCS Redis instance must meet the following complexity requirements:
    Changing this creates a new instance.

* `vpc_id` - (Required) Specifies the id of the VPC. Changing this creates a new instance.

* `network_id` - (Required) Specifies the network id of the subnet. Changing this creates a new instance.

* `security_group_id` - (Optional) Specifies the id of the security group which the instance belongs to.
    This parameter is mandatory for Memcached and Redis 3.0 versions.

* `available_zones` - (Required) IDs or Names of the AZs where cache nodes reside. For details
    on how to query AZs, see Querying AZ Information.
    Changing this creates a new instance.

* `product_id` - (Optional) Product ID used to differentiate DCS instance types. For example now there are following values available:

	- dcs.memcached.master_standby-h,
	- dcs.memcached.master_standby-m,
	- dcs.memcached.single_node-h,
	- dcs.memcached.single_node-m,
	- dcs.master_standby-h,
	- dcs.cluster-h,
	- dcs.single_node-h,
	- redis.cluster.xu1.large.r1.4-h,
	- redis.cluster.xu1.large.r2.4-h,
	- redis.cluster.xu1.large.r1.8-h,
	- redis.cluster.xu1.large.r2.16-h,
	- redis.cluster.xu1.large.r1.16-h,
	- redis.cluster.xu1.large.r2.24-h,
	- redis.cluster.xu1.large.r2.32-h,
	- redis.cluster.xu1.large.r1.32-h,
	- redis.cluster.xu1.large.r2.48-h,
	- redis.cluster.xu1.large.r1.48-h 
	- redis.cluster.xu1.large.r2.64-h
	- redis.cluster.xu1.large.r1.64-h
	- redis.cluster.xu1.large.r2.96-h
	- redis.cluster.xu1.large.r1.96-h
	- redis.cluster.xu1.large.r2.128-h
	- redis.cluster.xu1.large.r1.128-h
	- redis.cluster.xu1.large.r1.192-h
	- redis.cluster.xu1.large.r2.192-h
	- redis.cluster.xu1.large.r2.256-h
	- redis.cluster.xu1.large.r1.256-h
	- redis.cluster.xu1.large.r2.384-h
	- redis.cluster.xu1.large.r1.384-h
	- redis.cluster.xu1.large.r1.512-h
	- redis.cluster.xu1.large.r2.512-h

.....

For the whole list and the specification of product id please check the DCS API DOC for querying: https://dcs.eu-west-0.prod-cloud-ocb.orange-business.com/v1.0/products

Changing this creates a new instance.

* `instance_type` - (Deprecated, Optional) DCS instance specification code. For example now there are following values available:

	- dcs.single_node
	- dcs.master_standby
	- dcs.cluster

* `maintain_begin` - (Optional) Indicates the time at which a maintenance time window starts.
    Format: HH:mm:ss.
    The start time and end time of a maintenance time window must indicate the time segment of
	a supported maintenance time window. For details, see section Querying Maintenance Time Windows.
    The start time must be set to 22:00, 02:00, 06:00, 10:00, 14:00, or 18:00.
    Parameters maintain_begin and maintain_end must be set in pairs. If parameter maintain_begin
	is left blank, parameter maintain_end is also blank. In this case, the system automatically
	allocates the default start time 02:00.

* `maintain_end` - (Optional) Indicates the time at which a maintenance time window ends.
    Format: HH:mm:ss.
    The start time and end time of a maintenance time window must indicate the time segment of
	a supported maintenance time window. For details, see section Querying Maintenance Time Windows.
    The end time is four hours later than the start time. For example, if the start time is 22:00,
	the end time is 02:00.
    Parameters maintain_begin and maintain_end must be set in pairs. If parameter maintain_end is left
	blank, parameter maintain_begin is also blank. In this case, the system automatically allocates
	the default end time 06:00.

* `save_days` - (Optional) Retention time. Unit: day. Range: 1–7.
    Changing this creates a new instance.

* `backup_type` - (Optional) Backup type. Options:
    auto: automatic backup.
    manual: manual backup.
    Changing this creates a new instance.

* `begin_at` - (Optional) Time at which backup starts. "00:00-01:00" indicates that backup
    starts at 00:00:00. Changing this creates a new instance.

* `period_type` - (Optional) Interval at which backup is performed. Currently, only weekly
    backup is supported. Changing this creates a new instance.

* `backup_at` - (Optional) Day in a week on which backup starts. Range: 1–7. Where: 1
    indicates Monday; 7 indicates Sunday. Changing this creates a new instance.

## Attributes Reference

The following attributes are exported:


* `name` - See Argument Reference above.
* `description` - See Argument Reference above.
* `engine` - See Argument Reference above.
* `engine_version` - See Argument Reference above.
* `capacity` - See Argument Reference above.
* `access_user` - See Argument Reference above.
* `password` - See Argument Reference above.
* `vpc_id` - See Argument Reference above.
* `vpc_name` - Indicates the name of a vpc.
* `security_group_id` - See Argument Reference above.
* `security_group_name` - Indicates the name of a security group.
* `subnet_name` - Indicates the name of a subnet.
* `available_zones` - See Argument Reference above.
* `product_id` - See Argument Reference above.
* `maintain_begin` - See Argument Reference above.
* `maintain_end` - See Argument Reference above.
* `save_days` - See Argument Reference above.
* `backup_type` - See Argument Reference above.
* `begin_at` - See Argument Reference above.
* `period_type` - See Argument Reference above.
* `backup_at` - See Argument Reference above.
* `order_id` - An order ID is generated only in the monthly or yearly billing mode.
    In other billing modes, no value is returned for this parameter.
* `port` - Port of the cache node.
* `resource_spec_code` - Resource specifications.
    dcs.single_node: indicates a DCS instance in single-node mode.
    dcs.master_standby: indicates a DCS instance in master/standby mode.
    dcs.cluster: indicates a DCS instance in cluster mode.
* `used_memory` - Size of the used memory. Unit: MB.
* `internal_version` - Internal DCS version.
* `max_memory` - Overall memory size. Unit: MB.
* `user_id` - Indicates a user ID.
* `ip` - Cache node's IP address in tenant's VPC.
