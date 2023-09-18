---
subcategory: "Distributed Message Service (DMS)"
---

# flexibleengine_dms_kafka_instances

Use this data source to query the available instances within FlexibleEngine DMS service.

## Example Usage

### Query all instances with the keyword in the name

```hcl
variable "keyword" {}

data "flexibleengine_dms_kafka_instances" "test" {
  name        = var.keyword
  fuzzy_match = true
}
```

### Query the instance with the specified name

```hcl
variable "instance_name" {}

data "flexibleengine_dms_kafka_instances" "test" {
  name = var.instance_name
}
```

## Argument Reference

* `region` - (Optional, String) The region in which to query the kafka instance list.
  If omitted, the provider-level region will be used.

* `instance_id` - (Optional, String) Specifies the kafka instance ID to match exactly.

* `name` - (Optional, String) Specifies the kafka instance name for data-source queries.

* `fuzzy_match` - (Optional, Bool) Specifies whether to match the instance name fuzzily, the default is a exact
  match (`flase`).

* `status` - (Optional, String) Specifies the kafka instance status for data-source queries.

* `include_failure` - (Optional, Bool) Specifies whether the query results contain instances that failed to create.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `instances` - The result of the query's list of kafka instances. The [instances](#dms_instances) object structure
  is documented below.

<a name="dms_instances"></a>
The `instances` block supports:

* `id` - The instance ID.

* `type` - The instance type.

* `name` - The instance name.

* `description` - The instance description.

* `availability_zones` - The list of AZ names.

* `product_id` - The product ID used by the instance.

* `engine_version` - The kafka engine version.

* `storage_spec_code` - The storage I/O specification.

* `storage_space` - The message storage capacity, in GB unit.

* `vpc_id` - The VPC ID to which the instance belongs.

* `network_id` - The subnet ID to which the instance belongs.

* `security_group_id` - The security group ID associated with the instance.

* `manager_user` - The username for logging in to the Kafka Manager.

* `access_user` - The access username.

* `maintain_begin` - The time at which a maintenance time window starts, the format is `HH:mm`.

* `maintain_end` - The time at which a maintenance time window ends, the format is `HH:mm`.

* `enable_public_ip` - Whether public access to the instance is enabled.

* `public_ip_ids` - The IDs of the elastic IP address (EIP).

* `public_conn_addresses` - The instance public access address.
  The format of each connection address is `{IP address}:{port}`.

* `retention_policy` - The action to be taken when the memory usage reaches the disk capacity threshold.

* `dumping` - Whether to dumping is enabled.

* `enable_auto_topic` - Whether to enable automatic topic creation.

* `partition_num` - The maximum number of topics in the DMS kafka instance.

* `ssl_enable` - Whether the Kafka SASL_SSL is enabled.

* `used_storage_space` - The used message storage space, in GB unit.

* `connect_address` - The IP address for instance connection.

* `port` - The port number of the instance.

* `status` - The instance status.

* `resource_spec_code` - The resource specifications identifier.

* `user_id` - The user ID who created the instance.

* `user_name` - The username who created the instance.

* `management_connect_address` - The connection address of the Kafka manager of an instance.

* `tags` - The key/value pairs to associate with the instance.

* `cross_vpc_accesses` - Indicates the Access information of cross-VPC. The [cross_vpc_accesses](#dms_cross_vpc_accesses)
  object structure is documented below.

<a name="dms_cross_vpc_accesses"></a>
The `cross_vpc_accesses` block supports:

* `listener_ip` - The listener IP address.

* `advertised_ip` - The advertised IP Address.

* `port` - The port number.

* `port_id` - The port ID associated with the address.
