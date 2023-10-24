---
subcategory: "Relational Database Service (RDS)"
description: ""
page_title: "flexibleengine_rds_read_replica_v3"
---

# flexibleengine_rds_read_replica_v3

RDS read replica management

## Example Usage

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
  name        = "terraform_test_security_group"
  description = "terraform security group acceptance test"
}

resource "flexibleengine_rds_instance_v3" "instance_1" {
  name              = "terraform_test_rds_instance"
  flavor            = "rds.pg.s1.medium"
  availability_zone = [var.primary_az]
  security_group_id = flexibleengine_networking_secgroup_v2.example_secgroup.id
  vpc_id            = flexibleengine_vpc_v1.example_vpc.id
  subnet_id         = flexibleengine_vpc_subnet_v1.example_subnet.id

  db {
    password = var.db_password
    type     = "PostgreSQL"
    version  = "11"
    port     = "8635"
  }
  volume {
    type = "ULTRAHIGH"
    size = 100
  }
  backup_strategy {
    start_time = "08:00-09:00"
    keep_days  = 1
  }
}

resource "flexibleengine_rds_read_replica_v3" "instance_2" {
  name              = "replica_instance"
  flavor            = "rds.pg.c2.large.rr"
  replica_of_id     = flexibleengine_rds_instance_v3.instance_1.id
  availability_zone = var.primary_az

  volume {
    type = "ULTRAHIGH"
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the RDS read replica resource.
  If omitted, the provider-level region will be used. Changing this will create a new RDS read replica resource.
  Currently, read replicas can be created only in the same region as that of the primary DB instance.

* `name` - (Required, String, ForceNew) Specifies the DB instance name. The DB instance name of the same type must be
  unique for the same tenant. The value must be 4 to 64 characters in length and start with a letter. It is
  case-sensitive and can contain only letters, digits, hyphens (-), and underscores (_).
  Changing this parameter will create a new resource.

* `flavor` - (Required, String) Specifies the specification code.

* `replica_of_id` - (Required, String, ForceNew) Specifies the DB instance ID, which is used to create a read replica.
  Changing this parameter will create a new resource.

* `volume` - (Required, List, ForceNew) Specifies the volume information. The [volume](#rds_volume) object structure is
  documented below. Changing this parameter will create a new resource.

* `availability_zone` - (Required, String, ForceNew) Specifies the AZ name.
  Changing this parameter will create a new resource.

* `tags` - (Optional, Map) A mapping of tags to assign to the RDS read replica instance.
  Each tag is represented by one key-value pair.

<a name="rds_volume"></a>
The `volume` block supports:

* `type` - (Required, String, ForceNew) Specifies the volume type. Its value can be any of the following
  and is case-sensitive:
  - ULTRAHIGH: indicates the SSD type.
  - ULTRAHIGHPRO: indicates the ultra-high I/O.

  Changing this parameter will create a new resource.

* `disk_encryption_id` -  (Optional, String, ForceNew) Specifies the key ID for disk encryption.
  Changing this parameter will create a new resource.

* `size` - (Optional, Int) Specifies the volume size. Its value range is from **40** GB to **4000** GB. The value must
  be a multiple of 10 and greater than the original size.

## Attribute Reference

In addition to the arguments listed above, the following computed attributes are exported:

* `id` - Indicates the instance ID.

* `status` - Indicates the instance status.

* `type` -  Indicates the type of the read replica instance.

* `db` - Indicates the database information. The [db](#rds_db) object structure is documented below.

* `private_ips` - Indicates the private IP address list.

* `public_ips` - Indicates the public IP address list.

* `security_group_id` - Indicates the security group which the RDS DB instance belongs to.

* `subnet_id` - Indicates the subnet id.

* `vpc_id` - Indicates the VPC ID.

<a name="rds_db"></a>
The `db` block supports:

* `port` - Indicates the database port information.

* `type` - Indicates the DB engine. Value: MySQL, PostgreSQL, SQLServer.

* `user_name` - Indicates the default username of database.

* `version` - Indicates the database version.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 30 minutes.
* `delete` - Default is 30 minutes.

## Import

RDS instance can be imported using the `id`, e.g.

```shell
terraform import flexibleengine_rds_read_replica_v3.instance_1 7117d38e-4c8f-4624-a505-bd96b97d024c
```
