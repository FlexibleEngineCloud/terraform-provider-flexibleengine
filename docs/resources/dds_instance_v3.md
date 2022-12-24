---
subcategory: "Document Database Service (DDS)"
description: ""
page_title: "flexibleengine_dds_instance_v3"
---

# flexibleengine_dds_instance_v3

Manages dds instance resource within FlexibleEngine

## Example Usage: Creating a Cluster Community Edition

```hcl
resource "flexibleengine_vpc_v1" "example_vpc" {
  name = "example-vpc"
  cidr = "192.168.0.0/16"
}

resource "flexibleengine_vpc_subnet_v1" "example_subnet" {
  name       = "example-vpc-subnet"
  cidr       = "192.168.0.0/24"
  gateway_ip = "192.168.0.1"
  vpc_id     = flexibleengine_vpc_v1.example_vpc.id
}

resource "flexibleengine_networking_secgroup_v2" "example_secgroup" {
  name        = "example-secgroup"
  description = "My neutron security group"
}

resource "flexibleengine_dds_instance_v3" "instance" {
  name              = "dds-instance"
  region            = "eu-west-0"
  availability_zone = "eu-west-0a"
  vpc_id            = flexibleengine_vpc_v1.example_vpc.id
  subnet_id         = flexibleengine_vpc_subnet_v1.example_subnet.id
  security_group_id = flexibleengine_networking_secgroup_v2.example_secgroup.id
  password          = "Test@123"
  mode              = "Sharding"

  datastore {
    type           = "DDS-Community"
    version        = "3.4"
    storage_engine = "wiredTiger"
  }
  flavor {
    type      = "mongos"
    num       = 2
    spec_code = "dds.mongodb.s3.medium.4.mongos"
  }
  flavor {
    type      = "shard"
    num       = 2
    storage   = "ULTRAHIGH"
    size      = 20
    spec_code = "dds.mongodb.s3.medium.4.shard"
  }
  flavor {
    type      = "config"
    num       = 1
    storage   = "ULTRAHIGH"
    size      = 20
    spec_code = "dds.mongodb.s3.large.2.config"
  }
  backup_strategy {
    start_time = "08:00-09:00"
    keep_days  = "8"
  }
}
```

## Example Usage: Creating a Replica Set

```hcl
resource "flexibleengine_dds_instance_v3" "instance" {
  name              = "dds-instance"
  region            = "eu-west-0"
  availability_zone = "eu-west-0a"
  vpc_id            = flexibleengine_vpc_v1.example_vpc.id
  subnet_id         = flexibleengine_vpc_subnet_v1.example_subnet.id
  security_group_id = flexibleengine_networking_secgroup_v2.example_secgroup.id
  password          = "Test@123"
  mode              = "ReplicaSet"

  datastore {
    type           = "DDS-Community"
    version        = "3.4"
    storage_engine = "wiredTiger"
  }
  flavor {
    type      = "replica"
    num       = 1
    storage   = "ULTRAHIGH"
    size      = 30
    spec_code = "dds.mongodb.s3.medium.4.repset"
  }
  backup_strategy {
    start_time = "08:00-09:00"
    keep_days  = "8"
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional) Specifies the region of the DDS instance. Changing this creates
  a new instance.

* `name` - (Required) Specifies the DB instance name. The DB instance name of the same
  type is unique in the same tenant. Changing this creates a new instance.

* `datastore` - (Required) Specifies database information. The structure is described
  below. Changing this creates a new instance.

* `availability_zone` - (Required) Specifies the ID of the availability zone. Changing
  this creates a new instance.

* `vpc_id` - (Required) Specifies the VPC ID. Changing this creates a new instance.

* `subnet_id` - (Required) Specifies the ID of the VPC Subnet. Changing this creates a new instance.

* `security_group_id` - (Required) Specifies the security group ID of the DDS instance.
  Changing this creates a new instance.

* `password` - (Required) Specifies the Administrator password of the database instance.
  Changing this creates a new instance.

* `disk_encryption_id` - (Required) Specifies the disk encryption ID of the instance.
  Changing this creates a new instance.

* `mode` - (Required) Specifies the mode of the database instance. Changing this creates a new instance.

* `flavor` - (Required) Specifies the flavors information. The structure is described below.
  Changing this creates a new instance.

* `backup_strategy` - (Optional) Specifies the advanced backup policy. The structure is
  described below. Changing this creates a new instance.

* `ssl` - (Optional) Specifies whether to enable or disable SSL. Defaults to true.
  Changing this creates a new instance.

* `tags` - (Optional) The key/value pairs to associate with the DDS instance.

The `datastore` block supports:

* `type` - (Required) Specifies the DB engine. Only DDS-Community is supported now.

* `version` - (Required) Specifies the DB instance version. Only 3.4 and 4.0 are supported now.

* `storage_engine` - (Optional) Specifies the storage engine of the DB instance. Only wiredTiger is supported now.

The `flavor` block supports:

* `type` - (Required) Specifies the node type. Valid value: mongos, shard, config, replica.

* `num` - (Required) Specifies the node quantity. Valid value:
  + the number of mongos ranges from 2 to 12.
  + the number of shard ranges from 2 to 12.
  + config: the value is 1.
  + replica: the value is 1.

* `storage` - (Optional) Specifies the disk type. Valid value: ULTRAHIGH which indicates the type SSD.

* `size` - (Optional) Specifies the disk size. The value must be a multiple of 10. The unit is GB. This parameter
  is mandatory for nodes except mongos and invalid for mongos.

* `spec_code` - (Required) Specifies the resource specification code. Valid values:

engine_name | type | vcpus | ram | speccode
---- | --- | ---
DDS-Community | mongos | 1 | 4 | dds.mongodb.s3.medium.4.mongos
DDS-Community | mongos | 2 | 8 | dds.mongodb.s3.large.4.mongos
DDS-Community | mongos | 4 | 16 | dds.mongodb.s3.xlarge.4.mongos
DDS-Community | mongos | 8 | 32 | dds.mongodb.s3.2xlarge.4.mongos
DDS-Community | mongos | 16 | 64 | dds.mongodb.s3.4xlarge.4.mongos
DDS-Community | shard | 1 | 4 | dds.mongodb.s3.medium.4.shard
DDS-Community | shard | 2 | 8 | dds.mongodb.s3.large.4.shard
DDS-Community | shard | 4 | 16 | dds.mongodb.s3.xlarge.4.shard
DDS-Community | shard | 8 | 32 | dds.mongodb.s3.2xlarge.4.shard
DDS-Community | shard | 16 | 64 | dds.mongodb.s3.4xlarge.4.shard
DDS-Community | config | 2 | 4 | dds.mongodb.s3.large.2.config
DDS-Community | replica | 1 | 4 | dds.mongodb.s3.medium.4.repset
DDS-Community | replica | 2 | 8 | dds.mongodb.s3.large.4.repset
DDS-Community | replica | 4 | 16 | dds.mongodb.s3.xlarge.4.repset
DDS-Community | replica | 8 | 32 | dds.mongodb.s3.2xlarge.4.repset
DDS-Community | replica | 16 | 64 | dds.mongodb.s3.4xlarge.4.repset

The `backup_strategy` block supports:

* `start_time` - (Required) Specifies the backup time window. Automated backups will be triggered
  during the backup time window. The value cannot be empty. It must be a valid value in the
  "hh:mm-HH:MM" format. The current time is in the UTC format.
  + The HH value must be 1 greater than the hh value.
  + The values from mm and MM must be the same and must be set to any of the following 00, 15, 30, or 45.

* `keep_days` - (Optional) Specifies the number of days to retain the generated backup files.
  The value range is from 0 to 732.
  + If this parameter is set to 0, the automated backup policy is not set.
  + If this parameter is not transferred, the automated backup policy is enabled by default.
    Backup files are stored for seven days by default.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in UUID format.
* `status` - Indicates the the DB instance status.
* `db_username` - Indicates the DB Administator name.
* `port` - Indicates the database port number. The port range is 2100 to 9500.
* `nodes` - Indicates the instance nodes information. Structure is documented below.

The `nodes` block contains:

* `id` - Indicates the node ID.
* `name` - Indicates the node name.
* `role` - Indicates the node role.
* `type` - Indicates the node type.
* `private_ip` - Indicates the private IP address of a node. This parameter is valid only for
  mongos nodes, replica set instances, and single node instances.
* `public_ip` - Indicates the EIP that has been bound on a node. This parameter is valid only for
  mongos nodes of cluster instances, primary nodes and secondary nodes of replica set instances,
  and single node instances.
* `status` - Indicates the node status.
