---
subcategory: "Virtual Private Cloud (VPC)"
---

# flexibleengine_vpc_eip_v1

Manages an EIP resource within FlexibleEngine VPC.

## Example Usage

```hcl
resource "flexibleengine_vpc_eip_v1" "eip_1" {
  publicip {
    type = "5_bgp"
  }
  bandwidth {
    name       = "test"
    size       = 10
    share_type = "PER"
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional) The region in which to create the eip. If omitted,
    the `region` argument of the provider is used. Changing this creates a new eip.

* `publicip` - (Required) The elastic IP address object.

* `bandwidth` - (Required) The bandwidth object.


The `publicip` block supports:

* `type` - (Required) The value must be a type supported by the system. Only `5_bgp` supported now.
    Changing this creates a new eip.

* `ip_address` - (Optional) The value must be a valid IP address in the available IP address segment.
    Changing this creates a new eip.

* `port_id` - (Optional) The port id which this eip will associate with. If the value
    is "" or this not specified, the eip will be in unbind state.


The `bandwidth` block supports:

* `name` - (Required) The bandwidth name, which is a string of 1 to 64 characters
    that contain letters, digits, underscores (_), and hyphens (-).

* `size` - (Required) The bandwidth size. The value ranges from 1 to 1000 Mbit/s.

* `share_type` - (Required) Specifies the bandwidth type.
    The value is *PER*, indicating that the bandwidth is dedicated.
    Changing this creates a new eip.

* `charge_mode` - (Optional) Specifies whether the bandwidth is billed by traffic or by bandwidth size.
    Only *traffic* supported now. Changing this creates a new eip.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in UUID format.
* `status` - The status of eip.

## Import

EIPs can be imported using the `id`, e.g.

```
$ terraform import flexibleengine_vpc_eip_v1.eip_1 2c7f39f3-702b-48d1-940c-b50384177ee1
```
