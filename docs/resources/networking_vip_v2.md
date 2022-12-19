---
subcategory: "Virtual Private Cloud (VPC)"
description: ""
page_title: "flexibleengine_networking_vip_v2"
---

# flexibleengine_networking_vip_v2

Manages a V2 vip resource within FlexibleEngine.

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

resource "flexibleengine_networking_vip_v2" "vip_1" {
  network_id = flexibleengine_vpc_subnet_v1.example_subnet.id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the VIP.
  If omitted, the provider-level region will be used. Changing this will create a new VIP resource.

* `network_id` - (Required, String, ForceNew) Specifies the ID of the VPC Subnet to which the VIP belongs.
  Changing this will create a new VIP resource.

* `ip_version` - (Optional, Int, ForceNew) Specifies the IP version, either `4` (default) or `6`.
  Changing this will create a new VIP resource.

* `ip_address` - (Optional, String, ForceNew) Specifies the IP address desired in the subnet for this VIP.
  Changing this will create a new VIP resource.

* `name` - (Optional, String) Specifies a unique name for the VIP.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The VIP ID.

* `mac_address` - The MAC address of the VIP.

* `status` - The VIP status.

* `device_owner` - The device owner of the VIP.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 2 minute.
* `delete` - Default is 2 minute.

## Import

Network VIPs can be imported using their `id`, e.g.:

```shell
terraform import flexibleengine_networking_vip_v2.test ce595799-da26-4015-8db5-7733c6db292e
```
