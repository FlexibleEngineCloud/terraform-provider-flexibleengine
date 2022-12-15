---
subcategory: "Elastic Cloud Server (ECS)"
---

# flexibleengine\_compute\_instance

Use this data source to get the details of a specified compute instance.

## Example Usage

```hcl
variable "server_name" {}

data "flexibleengine_compute_instance_v2" "demo" {
  name = var.server_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional) The region in which to obtain the server instance.

* `name` - (Optional) Specifies the server name, which can be queried with a regular expression.

* `fixed_ip_v4` - (Optional)  Specifies the IPv4 addresses of the server.

* `flavor_id` - (Optional) Specifies the flavor ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The instance ID in UUID format.
* `availability_zone` - The availability zone where the instance is located.
* `image_id` - The image ID of the instance.
* `image_name` - The image name of the instance.
* `flavor_name` - The flavor name of the instance.
* `key_pair` - The key pair that is used to authenticate the instance.
* `floating_ip` - The EIP address that is associted to the instance.
* `system_disk_id` - The system disk voume ID.
* `user_data` -  The user data (information after encoding) configured during instance creation.
* `security_groups` - An array of one or more security group names
    to associate with the instance.
* `network` - An array of one or more networks to attach to the instance.
    The network object structure is documented below.
* `block_device` - An array of one or more disks to attach to the instance.
    The block_device object structure is documented below.
* `scheduler_hints` - The scheduler with hints on how the instance should be launched.
    The available hints are described below.
* `tags` - The tags of the instance in key/value format.
* `metadata` - The metadata of the instance in key/value format.
* `status` - The status of the instance.

The `network` block supports:

* `uuid` - The network UUID to attach to the server.
* `port` - The port ID corresponding to the IP address on that network.
* `mac` - The MAC address of the NIC on that network.
* `fixed_ip_v4` - The fixed IPv4 address of the instance on this network.
* `fixed_ip_v6` - The Fixed IPv6 address of the instance on that network.

The `block_device` block supports:

* `uuid` - The volume id on that attachment.
* `boot_index` - The volume boot index on that attachment.
* `size` - The volume size on that attachment.
* `type` - The volume type on that attachment.
* `pci_address` - The volume pci address on that attachment.

The `scheduler_hints` block supports:

* `group` - The UUID of a Server Group where the instance will be placed into.
