---
subcategory: "Virtual Private Cloud (VPC)"
---

# flexibleengine\_networking\_vip_associate_v2

Manages a V2 vip associate resource within FlexibleEngine.

## Example Usage

```hcl
resource "flexibleengine_networking_network_v2" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "flexibleengine_networking_subnet_v2" "subnet_1" {
  name = "subnet_1"
  cidr = "192.168.199.0/24"
  ip_version = 4
  network_id = "${flexibleengine_networking_network_v2.network_1.id}"
}

resource "flexibleengine_networking_router_interface_v2" "router_interface_1" {
  router_id = "${flexibleengine_networking_router_v2.router_1.id}"
  subnet_id = "${flexibleengine_networking_subnet_v2.subnet_1.id}"
}

resource "flexibleengine_networking_router_v2" "router_1" {
  name = "router_1"
  external_gateway = "0a2228f2-7f8a-45f1-8e09-9039e1d09975"
}

resource "flexibleengine_networking_port_v2" "port_1" {
  name = "port_1"
  admin_state_up = "true"
  network_id = "${flexibleengine_networking_network_v2.network_1.id}"

  fixed_ip {
    subnet_id =  "${flexibleengine_networking_subnet_v2.subnet_1.id}"
  }
}

resource "flexibleengine_compute_instance_v2" "instance_1" {
  name = "instance_1"
  security_groups = ["default"]

  network {
    port = "${flexibleengine_networking_port_v2.port_1.id}"
  }
}

resource "flexibleengine_networking_port_v2" "port_2" {
  name = "port_2"
  admin_state_up = "true"
  network_id = "${flexibleengine_networking_network_v2.network_1.id}"

  fixed_ip {
    subnet_id =  "${flexibleengine_networking_subnet_v2.subnet_1.id}"
  }
}

resource "flexibleengine_compute_instance_v2" "instance_2" {
  name = "instance_2"
  security_groups = ["default"]

  network {
    port = "${flexibleengine_networking_port_v2.port_1.id}"
  }
}

resource "flexibleengine_networking_vip_v2" "vip_1" {
  network_id = "${flexibleengine_networking_network_v2.network_1.id}"
  subnet_id = "${flexibleengine_networking_subnet_v2.subnet_1.id}"
}

resource "flexibleengine_networking_vip_associate_v2" "vip_associate_1" {
  vip_id = "${flexibleengine_networking_vip_v2.vip_1.id}"
  port_ids = ["${flexibleengine_networking_port_v2.port_1.id}", "${flexibleengine_networking_port_v2.port_2.id}"]
}
```

## Argument Reference

The following arguments are supported:

* `vip_id` - (Required) The ID of vip to attach the port to.
    Changing this creates a new vip associate.

* `port_ids` - (Required) An array of one or more IDs of the ports to attach the vip to.
    Changing this creates a new vip associate.

## Attributes Reference

The following attributes are exported:

* `vip_id` - See Argument Reference above.
* `port_ids` - See Argument Reference above.
* `vip_subnet_id` - The ID of the subnet this vip connects to.
* `vip_ip_address` - The IP address in the subnet for this vip.
