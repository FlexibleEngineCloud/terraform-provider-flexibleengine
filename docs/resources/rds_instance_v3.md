---
subcategory: "Relational Database Service (RDS)"
description: ""
page_title: "flexibleengine_rds_instance_v3"
---

# flexibleengine_rds_instance_v3

Manage RDS instance resource within FlexibleEngine.

## Example Usage

### create a single db instance

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

resource "flexibleengine_rds_instance_v3" "instance" {
  name              = "terraform_test_rds_instance"
  flavor            = "rds.pg.s3.medium.4"
  availability_zone = [var.primary_az]
  security_group_id = flexibleengine_networking_secgroup_v2.example_secgroup.id
  vpc_id            = flexibleengine_vpc_v1.example_vpc.id
  subnet_id         = flexibleengine_vpc_subnet_v1.example_subnet.id

  db {
    type     = "PostgreSQL"
    version  = "11"
    password = var.db_password
    port     = "8635"
  }
  volume {
    type = "COMMON"
    size = 100
  }
  backup_strategy {
    start_time = "08:00-09:00"
    keep_days  = 1
  }
}
```

### create a primary/standby db instance

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

resource "flexibleengine_rds_instance_v3" "instance" {
  name                = "terraform_test_rds_instance"
  flavor              = "rds.pg.s3.large.4.ha"
  ha_replication_mode = "async"
  availability_zone   = [var.primary_az, var.standby_az]
  security_group_id   = flexibleengine_networking_secgroup_v2.example_secgroup.id
  vpc_id              = flexibleengine_vpc_v1.example_vpc.id
  subnet_id           = flexibleengine_vpc_subnet_v1.example_subnet.id

  db {
    type     = "PostgreSQL"
    version  = "11"
    password = var.db_password
    port     = "8635"
  }
  volume {
    type = "COMMON"
    size = 100
  }
  backup_strategy {
    start_time = "08:00-09:00"
    keep_days  = 1
  }
}
```

### create a single db instance with encrypted volume

```hcl
resource "flexibleengine_kms_key_v1" "key" {
  key_alias       = "key_1"
  key_description = "first test key"
  is_enabled      = true
}

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
resource "flexibleengine_rds_instance_v3" "instance" {
  name              = "terraform_test_rds_instance"
  flavor            = "rds.pg.s3.medium.4"
  availability_zone = [var.primary_az]
  security_group_id = flexibleengine_networking_secgroup_v2.example_secgroup.id
  vpc_id            = flexibleengine_vpc_v1.example_vpc.id
  subnet_id         = flexibleengine_vpc_subnet_v1.example_subnet.id

  db {
    type     = "PostgreSQL"
    version  = "11"
    password = var.db_password
    port     = "8635"
  }
  volume {
    disk_encryption_id = flexibleengine_kms_key_v1.key.id
    type               = "COMMON"
    size               = 100
  }
  backup_strategy {
    start_time = "08:00-09:00"
    keep_days  = 1
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the RDS instance resource.
  If omitted, the provider-level region will be used. Changing this will create a new RDS instance resource.

* `name` - (Required, String) Specifies the DB instance name. The DB instance name of the same type must be unique for
  the same tenant. The value must be 4 to 64 characters in length and start with a letter. It is case-sensitive and can
  contain only letters, digits, hyphens (-), and underscores (_).

* `flavor` - (Required, String) Specifies the specification code.

  -> **NOTE:** Services will be interrupted for 5 to 10 minutes when you change RDS instance flavor.

* `availability_zone` - (Required, List, ForceNew) Specifies the list of AZ name.
  Changing this parameter will create a new resource.

* `db` - (Required, List, ForceNew) Specifies the database information. The [db](#rds_db) object structure is
  documented below. Changing this parameter will create a new resource.

* `vpc_id` - (Required, String, ForceNew) Specifies the VPC ID. Changing this parameter will create a new resource.

* `subnet_id` - (Required, String, ForceNew) Specifies the ID of the VPC Subnet.
  Changing this parameter will create a new resource.

* `security_group_id` - (Required, String) Specifies the security group which the RDS DB instance belongs to.

* `volume` - (Required, List) Specifies the volume information. The [volume](#rds_volume) object structure is
  documented below.

* `fixed_ip` - (Optional, String, ForceNew) Specifies an intranet IP address of RDS DB instance.
  Changing this parameter will create a new resource.

* `backup_strategy` - (Optional, List) Specifies the advanced backup policy. The [backup_strategy](#rds_backup_strategy)
  object structure is documented below.

* `ha_replication_mode` - (Optional, String, ForceNew) Specifies the replication mode for the standby DB instance.
  Changing this parameter will create a new resource.
  + For MySQL, the value is *async* or *semisync*.
  + For PostgreSQL, the value is *async* or *sync*.
  + For SQLServer, the value is *sync*.
  + For MariaDB, the value is *async* or *semisync*.

  -> **NOTE:** async indicates the asynchronous replication mode. semisync indicates the semi-synchronous replication
  mode. sync indicates the synchronous replication mode.

* `param_group_id` - (Optional, String, ForceNew) Specifies the parameter group ID.
  Changing this parameter will create a new resource.

* `time_zone` - (Optional, String, ForceNew) Specifies the UTC time zone. The value ranges from
  UTC-12:00 to UTC+12:00 at the full hour, and defaults to *UTC*.
  Changing this parameter will create a new resource.

* `ssl_enable` - (Optional, Bool) Specifies whether to enable the SSL for **MySQL** database.

* `description` - (Optional, String) Specifies the description of the instance. The value consists of 0 to 64
  characters, including letters, digits, periods (.), underscores (_), and hyphens (-).

* `tags` - (Optional, Map) A mapping of tags to assign to the RDS instance.
  Each tag is represented by one key-value pair.

* `parameters` - (Optional, List) Specify an array of one or more parameters to be set to the RDS instance after
  launched. You can check on console to see which parameters supported. The [parameters](#rds_parameters) object
  structure is documented below.

<a name="rds_db"></a>
The `db` block supports:

* `type` - (Required, String, ForceNew) Specifies the DB engine. Available value are *MySQL*, *PostgreSQL* and
  *MariaDB*. Changing this parameter will create a new resource.

* `version` - (Required, String, ForceNew) Specifies the database version. The supported versions of each database
  are as follows:
  + **MySQL**: MySQL databases support **5.6**, **5.7** and **8.0**.
  + **PostgreSQL**: PostgreSQL databases support **9.5**, **9.6**, **10**, **11**, **12**, **13**, **14** and
    **1.0 (Enhanced Edition)**.
  + **SQLServer**: SQLServer databases support **2014 SE** and **2014 EE**.
  + **MariaDB**: MariaDB databases support **10.5**.

  Changing this parameter will create a new resource.

* `password` - (Required, String, ForceNew) Specifies the database password. The value cannot be
  empty and should contain 8 to 32 characters, including uppercase
  and lowercase letters, digits, and the following special
  characters: ~!@#%^*-_=+? You are advised to enter a strong
  password to improve security, preventing security risks such as
  brute force cracking. Changing this parameter will create a new resource.

* `port` - (Optional, Int) Specifies the database port.
  + The MySQL database port ranges from 1024 to 65535 (excluding 12017 and 33071, which are occupied by the RDS system
    and cannot be used). The default value is 3306.
  + The PostgreSQL database port ranges from 2100 to 9500. The default value is 5432.
  + The Microsoft SQL Server database port can be 1433 or ranges from 2100 to 9500, excluding 5355 and 5985.
    The default value is 1433.
  + The MariaDB database port ranges from 1024 to 65535 (excluding 12017 and 33071, which are occupied by the RDS system
    and cannot be used). The default value is 3306.

<a name="rds_volume"></a>
The `volume` block supports:

* `size` - (Required, Int) Specifies the volume size. Its value range is from 40 GB to 4000 GB.
  The value must be a multiple of 10 and greater than the original size.

* `type` - (Required, String, ForceNew) Specifies the volume type. Its value can be any of the following
  and is case-sensitive:
  + *COMMON*: indicates the SATA type.
  + *ULTRAHIGH*: indicates the SSD type.
  + *CLOUDSSD*: cloud SSD storage. This storage type is supported only with general-purpose and dedicated DB
    instances.

  Changing this parameter will create a new resource.

* `disk_encryption_id` - (Optional, String, ForceNew) Specifies the key ID for disk encryption.
  Changing this parameter will create a new resource.

<a name="rds_backup_strategy"></a>
The `backup_strategy` block supports:

* `keep_days` - (Optional, Int) Specifies the retention days for specific backup files. The value range is from 0 to
  732. If this parameter is not specified or set to 0, the automated backup policy is disabled.

  -> **NOTE:** Primary/standby DB instances of Microsoft SQL Server do not support disabling the automated backup
  policy.

* `start_time` - (Required, String) Specifies the backup time window. Automated backups will be triggered during the
  backup time window. It must be a valid value in the **hh:mm-HH:MM** format.
  The current time is in the UTC format. The HH value must be 1 greater than the hh value. The values of mm and MM must
  be the same and must be set to any of the following: 00, 15, 30, or 45. Example value: 08:15-09:15 23:00-00:00.

<a name="rds_parameters"></a>
The `parameters` block supports:

* `name` - (Required, String) Specifies the parameter name. Some of them needs the instance to be restarted
  to take effect.

* `value` - (Required, String) Specifies the parameter value.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in UUID format.

* `status` - Indicates the DB instance status.

* `created` - Indicates the creation time.

* `nodes` - Indicates the instance nodes information. The [nodes](#rds_nodes) object structure is documented below.

* `private_ips` - Indicates the private IP address list.
  It is a blank string until an ECS is created.

* `public_ips` - Indicates the public IP address list.

* `db` - See Argument Reference above. The [db](#rds_attr_db) object structure is documented below.

<a name="rds_nodes"></a>
The `nodes` block supports:

* `availability_zone` - Indicates the AZ.

* `id` - Indicates the node ID.

* `name` - Indicates the node name.

* `role` - Indicates the node type. The value can be master or slave,
  indicating the primary node or standby node respectively.

* `status` - Indicates the node status.

<a name="rds_attr_db"></a>
The `db` block supports:

* `user_name` - Indicates the default username of database.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 30 minutes.
* `update` - Default is 30 minutes.

## Import

RDS instance can be imported using the `id`, e.g.

```shell
terraform import flexibleengine_rds_instance_v3.instance_1 7117d38e-4c8f-4624-a505-bd96b97d024c
```

But due to some attributes missing from the API response, it's required to ignore changes as below.

```hcl
resource "flexibleengine_rds_instance_v3" "instance_1" {
  ...

  lifecycle {
    ignore_changes = [
      "db",
    ]
  }
}
```
