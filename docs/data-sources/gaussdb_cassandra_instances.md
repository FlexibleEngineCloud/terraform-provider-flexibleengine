---
subcategory: "GaussDB NoSQL"
---

# flexibleengine_gaussdb_cassandra_instances

Use this data source to get available FlexibleEngine gaussdb cassandra instances.

## Example Usage

```hcl
data "flexibleengine_gaussdb_cassandra_instances" "this" {
  name = "gaussdb-instance"
}
```

## Argument Reference

* `region` - (Optional, String) The region in which to obtain the instance. If omitted, the provider-level region will
  be used.

* `name` - (Optional, String) Specifies the name of the instance.

* `vpc_id` - (Optional, String) Specifies the VPC ID.

* `subnet_id` - (Optional, String) Specifies the network ID of a subnet.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Indicates the ID of the data source.

* `instances` - An array of available instances. The [instances](#gaussdb_instances) object structure is documented
  below.

<a name="gaussdb_instances"></a>
The `instances` block supports:

* `id` - The id of the instance.

* `region` - The region of the instance.

* `name` - Indicates the name of the instance.

* `vpc_id` - Indicates the VPC ID.

* `subnet_id` - Indicates the network ID of a subnet.

* `status` - Indicates the DB instance status.

* `mode` - Indicates the instance mode.

* `flavor` - Indicates the instance specifications.

* `security_group_id` - Indicates the security group ID.

* `enterprise_project_id` - Indicates the enterprise project id.

* `db_user_name` - Indicates the default username.

* `availability_zone` - Indicates the instance availability zone.

* `port` - Indicates the database port.

* `node_num` - Indicates the count of the nodes.

* `volume_size` - Indicates the size of the volume.

* `private_ips` - Indicates the list of private IP address of the nodes.

* `datastore` - Indicates the database information. The [datastore](#gaussdb_datastore) object structure is documented
  below.

* `backup_strategy` - Indicates the advanced backup policy. The [backup_strategy](#gaussdb_backup_strategy) object
  structure is documented below.

* `nodes` - Indicates the instance nodes information. The [nodes](#gaussdb_nodes) object structure is documented below.

* `tags` - Indicates the key/value tags of the instance.

<a name="gaussdb_datastore"></a>
The `datastore` block supports:

* `engine` - Indicates the database engine.

* `storage_engine` - Indicates the database storage engine.

* `version` - Indicates the database version.

<a name="gaussdb_backup_strategy"></a>
The `backup_strategy` block supports:

* `start_time` - Indicates the backup time window.

* `keep_days` - Indicates the number of days to retain the generated

<a name="gaussdb_nodes"></a>
The `nodes` block supports:

* `id` - Indicates the node ID.

* `name` - Indicates the node name.

* `private_ip` - Indicates the private IP address of a node.

* `status` - Indicates the node status.

* `support_reduce` - Indicates whether the node support reduce.

* `availability_zone` - Indicates the availability zone where the node resides.
