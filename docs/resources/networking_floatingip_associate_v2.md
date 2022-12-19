---
subcategory: "Deprecated"
description: ""
page_title: "flexibleengine_networking_floatingip_associate_v2"
---

# flexibleengine_networking_floatingip_associate_v2

!> **Warning:** It will be deprecated, using `flexibleengine_vpc_eip_associate` instead.

Associates a floating IP to a port. This is useful for situations where you have
a pre-allocated floating IP or are unable to use the `flexibleengine_vpc_eip`
or `flexibleengine_networking_floatingip_v2` resource to create a floating IP.

You can also use `publicip.0.port_id` of `flexibleengine_vpc_eip` resource to associate a port.

## Example Usage

```hcl
resource "flexibleengine_networking_port_v2" "port_1" {
  network_id = "a5bbd213-e1d3-49b6-aed1-9df60ea94b9a"
}

resource "flexibleengine_networking_floatingip_associate_v2" "fip_1" {
  floating_ip = "1.2.3.4"
  port_id     = flexibleengine_networking_port_v2.port_1.id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional) The region in which to obtain the V2 Networking client.
    A Networking client is needed to create a floating IP that can be used with
    another networking resource, such as a load balancer. If omitted, the
    `region` argument of the provider is used. Changing this creates a new
    floating IP (which may or may not have a different address).

* `floating_ip` - (Required) IP Address of an existing floating IP.

* `port_id` - (Required) ID of an existing port with at least one IP address to
    associate with this floating IP.

## Attributes Reference

The following attributes are exported:

* `region` - See Argument Reference above.
* `floating_ip` - See Argument Reference above.
* `port_id` - See Argument Reference above.

## Import

Floating IP associations can be imported using the `id` of the floating IP, e.g.

```shell
terraform import flexibleengine_networking_floatingip_associate_v2.fip 2c7f39f3-702b-48d1-940c-b50384177ee1
```
