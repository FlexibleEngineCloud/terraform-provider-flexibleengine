---
subcategory: "Relational Database Service (RDS)"
---

# flexibleengine_rds_instances

Use this data source to list all available RDS instances.

## Example Usage

```hcl
data "flexibleengine_rds_instances" "test" {
  name = "rds-instance"
}
```

## Argument Reference

* `region` - (Optional, String) The region in which query obtain the instances. If omitted, the provider-level region
  will be used.

* `name` - (Optional, String) Specifies the name of the instance.

* `type` - (Optional, String) Specifies the type of the instance. Valid values are:
  **Single**, **Ha**, **Replica**, and **Enterprise**.

* `datastore_type` - (Optional, String) Specifies the type of the database. Valid values are:
  **MySQL**, **PostgreSQL**, **SQLServer** and **MariaDB**.

* `vpc_id` - (Optional, String) Specifies the VPC ID.

* `subnet_id` - (Optional, String) Specifies the network ID of a subnet.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project id.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the data source.

* `instances` - An array of available instances. The [instances](#rds_instances) object structure is documented below.

<a name="rds_instances"></a>
The `instances` block supports:

* `region` - The region of the instance.

* `name` - Indicates the name of the instance.

* `availability_zone` - Indicates the availability zone name.

* `flavor` - Indicates the instance specifications.

* `vpc_id` - Indicates the VPC ID.

* `subnet_id` - Indicates the network ID of a subnet.

* `id` - Indicates the ID of the instance.

* `security_group_id` - Indicates the security group ID.

* `param_group_id` - Indicates the configuration ID.

* `enterprise_project_id` - Indicates the enterprise project id.

* `fixed_ip` - Indicates the intranet floating IP address of the instance.

* `ssl_enable` - Indicates whether to enable SSL.

* `tags` - Indicates the tags of the instance.

* `ha_replication_mode` - Indicates the replication mode for the standby DB instance.

* `time_zone` - Indicates the time zone.

* `private_ips` - Indicates the private ips in list.

* `public_ips` - Indicates the public ips in list.

* `status` - Indicates the DB instance status.

* `created` - Indicates the creation time.

* `db` - Indicates the database information. The [db](#rds_db) object structure is documented below.

* `volume` - Indicates the volume information. The [volume](#rds_volume) object structure is documented below.

* `backup_strategy` - Indicates the advanced backup policy. The [backup_strategy](#rds_backup_strategy) object structure
  is documented below.

* `nodes` - Indicates the instance nodes information. The [nodes](#rds_nodes) object structure is documented below.

<a name="rds_db"></a>
The `db` block supports:

* `type` - Indicates the database engine.

* `version` - Indicates the database version.

* `port` - Indicates the database port.

* `user_name` - Indicates the database username.

<a name="rds_volume"></a>
The `volume` block supports:

* `size` - Indicates the volume size.

* `type` - Indicates the volume type.

* `disk_encryption_id` - Indicates the kms key id.

<a name="rds_backup_strategy"></a>
The `backup_strategy` block supports:

* `start_time` - Indicates the backup time window.

* `keep_days` - Indicates the number of days to retain the generated.

<a name="rds_nodes"></a>
The `nodes` block supports:

* `id` - Indicates the node ID.

* `name` - Indicates the node name.

* `status` - Indicates the node status.

* `role` - Indicates the node type.

* `availability_zone` - Indicates the availability zone where the node resides.
