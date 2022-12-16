---
subcategory: "Elastic Cloud Server (ECS)"
description: ""
page_title: "flexibleengine_compute_floatingip_v2"
---

# flexibleengine_compute_floatingip_v2

Manages a V2 floating IP resource within FlexibleEngine Nova (compute)
that can be used for compute instances.

!> **Warning:** It will be deprecated, using `flexibleengine_vpc_eip` instead.

## Example Usage

```hcl
resource "flexibleengine_compute_floatingip_v2" "floatip_1" {
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional) The region in which to obtain the V2 Compute client.
    A Compute client is needed to create a floating IP that can be used with
    a compute instance. If omitted, the `region` argument of the provider
    is used. Changing this creates a new floating IP (which may or may not
    have a different address).

* `pool` - (Optional) The name of the pool from which to obtain the floating
    IP. Default value is admin_external_net. Changing this creates a new floating IP.

## Attributes Reference

The following attributes are exported:

* `region` - See Argument Reference above.
* `pool` - See Argument Reference above.
* `address` - The actual floating IP address itself.
* `fixed_ip` - The fixed IP address corresponding to the floating IP.
* `instance_id` - UUID of the compute instance associated with the floating IP.

## Import

Floating IPs can be imported using the `id`, e.g.

```shell
terraform import flexibleengine_compute_floatingip_v2.floatip_1 89c60255-9bd6-460c-822a-e2b959ede9d2
```
