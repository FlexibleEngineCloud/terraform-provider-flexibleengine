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
    uuid = var.subnet_id
  }
  ...
}

resource "flexibleengine_compute_instance_v2" "server_2" {
  name = "instance_2"
  network {
    uuid = var.subnet_id
  }
  ...
}

resource "flexibleengine_networking_vip_v2" "vip_1" {
  network_id = var.subnet_id
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

* `vip_id` - (Required) The ID of vip to attach the port to.
    Changing this creates a new vip associate.

* `port_ids` - (Required) An array of one or more IDs of the ports to attach the vip to.
    Changing this creates a new vip associate.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
* `vip_subnet_id` - The ID of the subnet this vip connects to.
* `vip_ip_address` - The IP address in the subnet for this vip.
