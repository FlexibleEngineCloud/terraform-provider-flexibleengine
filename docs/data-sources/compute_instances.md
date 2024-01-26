---
subcategory: "Elastic Cloud Server (ECS)"
description: ""
page_title: "flexibleengine_compute_instances"
---

# flexibleengine_compute_instances

Use this data source to get a list of compute instances.

## Example Usage

```hcl
variable "name_regex" {}

data "flexibleengine_compute_instances" "demo" {
  name = var.name_regex
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) The region in which to obtain the server instance.

* `name` - (Optional, String) Specifies the server name, which can be queried with a regular expression.

* `fixed_ip_v4` - (Optional, String)  Specifies the IPv4 addresses of the server.

* `flavor_id` - (Optional, String) Specifies the flavor ID.

* `status` - (Optional, String) Specifies the status of the instance. The valid values are as follows:
  + **ACTIVE**: The instance is running properly.
  + **SHUTOFF**: The instance has been properly stopped.
  + **ERROR**: An error has occurred on the instance.

* `flavor_name` - (Optional, String) Specifies the flavor name of the instance.

* `image_id` - (Optional, String) Specifies the image ID of the instance.

* `availability_zone` - (Optional, String) Specifies the availability zone where the instance is located.

* `key_pair` - (Optional, String) Specifies the key pair that is used to authenticate the instance.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `instances` - List of ECS instance details. The [instances](#ecs_attr_instances) object structure is documented below.

<a name="ecs_attr_instances"></a>
The `instances` block supports:

* `id` - The instance ID in UUID format.

* `name` - The instance name.

* `availability_zone` - The availability zone where the instance is located.

* `image_id` - The image ID of the instance.

* `image_name` - The image name of the instance.

* `flavor_id` - The flavor ID of the instance.

* `flavor_name` - The flavor name of the instance.

* `key_pair` - The key pair that is used to authenticate the instance.

* `floating_ip` - The EIP address that is associated to the instance.

* `user_data` -  The user data (information after encoding) configured during instance creation.

* `security_groups` - An array of one or more security group names
    to associate with the instance.

* `network` - An array of one or more networks to attach to the instance.
  The [network](#ecs_attr_network) object structure is documented below.

* `volume_attached` - An array of one or more disks to attach to the instance.
  The [volume_attached](#ecs_attr_volume_attached) object structure is documented below.

* `scheduler_hints` - The scheduler with hints on how the instance should be launched.
  The [scheduler_hints](#ecs_attr_scheduler_hints) object structure is documented below.

* `tags` - The tags of the instance in key/value format.

* `metadata` - The metadata of the instance in key/value format.

* `status` - The status of the instance.

<a name="ecs_attr_network"></a>
The `network` block supports:

* `uuid` - The network UUID to attach to the server.

* `port` - The port ID corresponding to the IP address on that network.

* `mac` - The MAC address of the NIC on that network.

* `fixed_ip_v4` - The fixed IPv4 address of the instance on this network.

* `fixed_ip_v6` - The Fixed IPv6 address of the instance on that network.

<a name="ecs_attr_volume_attached"></a>
The `volume_attached` block supports:

* `volume_id` - The volume id on that attachment.

* `is_sys_volume` - Whether the volume is the system disk.

<a name="ecs_attr_scheduler_hints"></a>
The `scheduler_hints` block supports:

* `group` - The UUID of a Server Group where the instance will be placed into.
