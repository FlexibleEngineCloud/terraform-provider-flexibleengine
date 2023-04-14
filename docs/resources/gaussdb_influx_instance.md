---
subcategory: "GaussDB NoSQL"
---

# flexibleengine_gaussdb_influx_instance

GaussDB for influx instance management within FlexibleEngine.

## Example Usage

### create a gaussdb for influx instance with tags

```hcl
resource "flexibleengine_gaussdb_influx_instance" "instance_1" {
  name              = "gaussdb_influx_instance_1"
  password          = var.password
  flavor            = "geminidb.influxdb.large.4"
  volume_size       = 100
  vpc_id            = var.vpc_id
  subnet_id         = var.subnet_id
  security_group_id = var.secgroup_id
  availability_zone = var.availability_zone

  tags = {
    foo = "bar"
    key = "value"
  }
}
```

### create a gaussdb influx instance with backup strategy

```hcl
resource "flexibleengine_gaussdb_influx_instance" "instance_1" {
  name              = "gaussdb_influx_instance_1"
  password          = var.password
  flavor            = "geminidb.influxdb.large.4"
  volume_size       = 100
  vpc_id            = var.vpc_id
  subnet_id         = var.subnet_id
  security_group_id = var.secgroup_id
  availability_zone = var.availability_zone

  backup_strategy {
    start_time = "03:00-04:00"
    keep_days  = 14
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the influx instance resource. If omitted, the
  provider-level region will be used. Changing this creates a new influx instance resource.

* `availability_zone` - (Required, String, ForceNew) Specifies the AZ name. For a three-AZ deployment instance,
  use commas (,) to separate the AZs, for example, `cn-north-4a,cn-north-4b,cn-north-4c`.
  Changing this parameter will create a new resource.

* `name` - (Required, String) Specifies the instance name, which can be the same as an existing instance name. The
  value must be **4** to **64** characters in length and start with a letter. It is case-sensitive and can contain only
  letters, digits, hyphens (-), and underscores (_).

* `flavor` - (Required, String, ForceNew) Specifies the instance specifications. For details,
  see [DB Instance Specifications](https://docs.prod-cloud-ocb.orange-business.com/usermanual/nosql/nosql_05_0045.html)
  Changing this parameter will create a new resource.

* `node_num` - (Optional, Int) Specifies the number of nodes, ranges from **3** to **16**. Defaults to **3**.

* `volume_size` - (Required, Int) Specifies the storage space in GB. The value must be a multiple of **10**. For a
  GaussDB influx instance, the minimum storage space is **100** GB, and the maximum storage space is related to the
  instance performance specifications. For details,
  see [DB Instance Specifications](https://docs.prod-cloud-ocb.orange-business.com/usermanual/nosql/nosql_05_0045.html)

* `password` - (Required, String) Specifies the database password. The value must be **8** to **32** characters in
  length, including uppercase and lowercase letters, digits, and special characters, such as ~!@#%^*-_=+? You are
  advised to enter a strong password to improve security, preventing security risks such as brute force cracking.

* `vpc_id` - (Required, String, ForceNew) Specifies the VPC ID. Changing this parameter will create a new resource.

* `subnet_id` - (Required, String, ForceNew) Specifies the network ID of a subnet. Changing this parameter will create
  a new resource.

* `security_group_id` - (Optional, String) Specifies the security group ID. Required if the selected subnet doesn't
  enable network ACL.

* `configuration_id` - (Optional, String) Specifies the Parameter Template ID.

* `ssl` - (Optional, Bool, ForceNew) Specifies whether to enable or disable SSL. Defaults to **false**. Changing this
  parameter will create a new resource.

* `force_import` - (Optional, Bool) If specified, try to import the instance instead of creating if the name already
  existed.

* `datastore` - (Optional, List, ForceNew) Specifies the database information. Structure is documented below. Changing
  this parameter will create a new resource.

* `backup_strategy` - (Optional, List) Specifies the advanced backup policy. Structure is documented below.

* `tags` - (Optional, Map) The key/value pairs to associate with the instance.

The `datastore` block supports:

* `engine` - (Required, String, ForceNew) Specifies the database engine. Only **influxdb** is supported now.
  Changing this parameter will create a new resource.

* `version` - (Required, String, ForceNew) Specifies the database version.
  Changing this parameter will create a new resource.

* `storage_engine` - (Required, String, ForceNew) Specifies the storage engine. Only **rocksDB** is supported now.
  Changing this parameter will create a new resource.

The `backup_strategy` block supports:

* `start_time` - (Required, String) Specifies the backup time window. Automated backups will be triggered during the
  backup time window. It must be a valid value in the "hh:mm-HH:MM" format. The current time is in the UTC format. The
  HH value must be 1 greater than the hh value. The values of mm and MM must be the same and must be set to 00. Example
  value: 08:00-09:00, 03:00-04:00.

* `keep_days` - (Optional, Int) Specifies the number of days to retain the generated backup files. The value ranges from
  **0** to **35**. If this parameter is set to **0**, the automated backup policy is not set. If this parameter is not
  transferred, the automated backup policy is enabled by default. Backup files are stored for seven days by default.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Indicates a resource ID in UUID format.
* `status` - Indicates the DB instance status.
* `port` - Indicates the database port.
* `mode` - Indicates the instance type.
* `db_user_name` - Indicates the default username.
* `nodes` - Indicates the instance nodes information. Structure is documented below.
* `private_ips` - Indicates the IP address list of the db.
* `lb_ip_address` - Indicates the LB IP address of the db.
* `lb_port` - Indicates the LB port of the db.

The `nodes` block contains:

* `id` - Indicates the node ID.
* `name` - Indicates the node name.
* `status` - Indicates the node status.
* `support_reduce` - Indicates whether the node support reduce or not.
* `private_ip` - Indicates the private IP address of a node.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 60 minute.
* `update` - Default is 120 minute.
* `delete` - Default is 30 minute.

## Import

GaussDB influx instance can be imported using the `id`, e.g.

```bash
terraform import flexibleengine_gaussdb_influx_instance.instance_1 2e045d8b-b226-4aa2-91b9-7e76357655c06
```
