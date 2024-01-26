---
subcategory: "Document Database Service (DDS)"
---

# flexibleengine_dds_instances

Use this data source to get the list of DDS instances.

## Example Usage

```hcl
variable "vpc_id" {}
variable "subnet_id" {}

data "flexibleengine_dds_instances" "test" {
  name      = "test_name"
  mode      = "Sharding"
  vpc_id    = var.vpc_id
  subnet_id = var.subnet_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source. If omitted, the provider-level
  region will be used.

* `name` - (Optional, String) Specifies the DB instance name.

* `mode` - (Optional, String) Specifies the mode of the database instance.

* `vpc_id` - (Optional, String) Specifies the VPC ID.

* `subnet_id` - (Optional, String) Specifies the subnet Network ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `instances` - Indicates the list of DDS instances.
  The [instances](#dds_instances) object structure is documented below.

<a name="dds_instances"></a>
The `instances` block supports:

* `id` - Indicates the ID of the instance.

* `name` - Indicates the DB instance name.

* `ssl` - Indicates whether to enable or disable SSL.

* `port` - Indicates the database port number. The port range is 2100 to 9500.

* `datastore` - Indicates database information.
  The [datastore](#dds_datastore) object structure is documented below.

* `backup_strategy` - Indicates backup strategy.
  The [backup_strategy](#dds_backup_strategy) object structure is documented below.

* `vpc_id` - Indicates the VPC ID.

* `subnet_id` - Indicates the subnet Network ID.

* `security_group_id` - Indicates the security group ID of the DDS instance.

* `disk_encryption_id` - Indicates the disk encryption ID of the instance.

* `mode` - Specifies the mode of the database instance.

* `db_username` - Indicates the DB Administrator name.

* `status` - Indicates the DB instance status.

* `enterprise_project_id` - Indicates the enterprise project id of the dds instance.

* `nodes` - Indicates the instance nodes information.
  The [nodes](#dds_nodes) object structure is documented below.

* `tags` - Indicates the key/value pairs to associate with the DDS instance.

<a name="dds_datastore"></a>
The `datastore` block supports:

* `type` - Indicates the DB engine.

* `version` - Indicates the DB instance version.

* `storage_engine` - Indicates the storage engine of the DB instance.

<a name="dds_backup_strategy"></a>
The `backup_strategy` block supports:

* `start_time` - Indicates the backup time window.

* `keep_days` - Indicates the number of days to retain the generated backup files.

<a name="dds_nodes"></a>
The `nodes` block supports:

* `id` - Indicates the node ID.

* `name` - Indicates the node name.

* `role` - Indicates the node role.

* `type` - Indicates the node type.

* `private_ip` - Indicates the private IP address of a node.

* `public_ip` - Indicates the EIP that has been bound on a node.

* `status` - Indicates the node status.
