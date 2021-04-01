---
subcategory: "Virtual Private Cloud (VPC)"
---

# flexibleengine\_networking\_network_v2

Manages a V2 Neutron network resource within FlexibleEngine.

## Example Usage

```hcl
resource "flexibleengine_networking_network_v2" "network_1" {
  name           = "network_1"
  admin_state_up = "true"
}

resource "flexibleengine_networking_subnet_v2" "subnet_1" {
  name       = "subnet_1"
  network_id = "${flexibleengine_networking_network_v2.network_1.id}"
  cidr       = "192.168.199.0/24"
  ip_version = 4
}

resource "flexibleengine_compute_secgroup_v2" "secgroup_1" {
  name        = "secgroup_1"
  description = "a security group"

  rule {
    from_port   = 22
    to_port     = 22
    ip_protocol = "tcp"
    cidr        = "0.0.0.0/0"
  }
}

resource "flexibleengine_networking_port_v2" "port_1" {
  name               = "port_1"
  network_id         = "${flexibleengine_networking_network_v2.network_1.id}"
  admin_state_up     = "true"
  security_group_ids = ["${flexibleengine_compute_secgroup_v2.secgroup_1.id}"]

  fixed_ip {
    "subnet_id"  = "${flexibleengine_networking_subnet_v2.subnet_1.id}"
    "ip_address" = "192.168.199.10"
  }
}

resource "flexibleengine_compute_instance_v2" "instance_1" {
  name            = "instance_1"
  security_groups = ["${flexibleengine_compute_secgroup_v2.secgroup_1.name}"]

  network {
    port = "${flexibleengine_networking_port_v2.port_1.id}"
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional) The region in which to obtain the V2 Networking client.
    A Networking client is needed to create a Neutron network. If omitted, the
    `region` argument of the provider is used. Changing this creates a new
    network.

* `name` - (Optional) The name of the network. Changing this updates the name of
    the existing network.

* `shared` - (Optional)  Specifies whether the network resource can be accessed
    by any tenant or not. Changing this updates the sharing capabalities of the
    existing network.

* `tenant_id` - (Optional) The owner of the network. Required if admin wants to
    create a network for another tenant. Changing this creates a new network.

* `admin_state_up` - (Optional) The administrative state of the network.
    Acceptable values are "true" and "false". Changing this value updates the
    state of the existing network.

* `segments` - (Optional) An array of one or more provider segment objects.

* `value_specs` - (Optional) Map of additional options.

The `segments` block supports:

* `physical_network` - The physical network where this network is implemented.
* `segmentation_id` - An isolated segment on the physical network.
* `network_type` - The type of physical network.

## Attributes Reference

The following attributes are exported:

* `region` - See Argument Reference above.
* `name` - See Argument Reference above.
* `shared` - See Argument Reference above.
* `tenant_id` - See Argument Reference above.
* `admin_state_up` - See Argument Reference above.

## Import

Networks can be imported using the `id`, e.g.

```
$ terraform import flexibleengine_networking_network_v2.network_1 d90ce693-5ccf-4136-a0ed-152ce412b6b9
```
