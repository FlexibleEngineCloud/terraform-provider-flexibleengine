---
subcategory: "Virtual Private Cloud (VPC)"
description: ""
page_title: "flexibleengine_networking_vip_associate_v2"
---

# flexibleengine_networking_vip_associate_v2

Manages a V2 vip associate resource within FlexibleEngine.

## Example Usage

```hcl
variable subnet_id{}

resource "flexibleengine_compute_instance_v2" "server_1" {
  name = "instance_1"
  network {
    uuid = flexibleengine_vpc_subnet_v1.example_subnet.id
  }
  ...
}

resource "flexibleengine_compute_instance_v2" "server_2" {
  name = "instance_2"
  network {
    uuid = flexibleengine_vpc_subnet_v1.example_subnet.id
  }
  ...
}

resource "flexibleengine_networking_vip_v2" "vip_1" {
  network_id = flexibleengine_vpc_subnet_v1.example_subnet.id
}

resource "flexibleengine_networking_vip_associate_v2" "vip_associate_1" {
  vip_id   = flexibleengine_networking_vip_v2.vip_1.id
  port_ids = [
    flexibleengine_compute_instance_v2.server_1.network.0.port,
    flexibleengine_compute_instance_v2.server_2.network.0.port,
  ]
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the vip associate resource. If omitted, the
  provider-level region will be used.

* `vip_id` - (Required, String, ForceNew) The ID of vip to attach the port to.
  Changing this creates a new vip associate.

* `port_ids` - (Required, List, ForceNew) An array of one or more IDs of the ports to attach the vip to.
  Changing this creates a new vip associate.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `vip_subnet_id` - The ID of the subnet this vip connects to.

* `vip_ip_address` - The IP address in the subnet for this vip.

* `ip_addresses` - The IP addresses of ports to attach the vip to.

## Import

Vip associate can be imported using the `vip_id` and port IDs separated by slashes (no limit on the number of
port IDs), e.g.

```shell
terraform import flexibleengine_networking_vip_associate_v2.vip_associated vip_id/port1_id/port2_id
```
