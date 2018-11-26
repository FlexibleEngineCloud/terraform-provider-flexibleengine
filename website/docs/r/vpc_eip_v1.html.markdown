---
layout: "flexibleengine"
page_title: "FlexibleEngine: flexibleengine_vpc_eip_v1"
sidebar_current: "docs-flexibleengine-resource-vpc-eip-v1"
description: |-
  Manages a V1 EIP resource within FlexibleEngine VPC.
---

# flexibleengine\_vpc\_eip_v1

Manages a V1 EIP resource within FlexibleEngine VPC.

## Example Usage

```hcl
resource "flexibleengine_vpc_eip_v1" "eip_1" {
  publicip {
    type = "5_bgp"
  }
  bandwidth {
    name = "test"
    size = 8
    share_type = "PER"
    charge_mode = "traffic"
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

* `type` - (Required) The value must be a type supported by the system. Only
    `5_bgp` supported now. Changing this creates a new eip.

* `ip_address` - (Optional) The value must be a valid IP address in the available
    IP address segment. Changing this creates a new eip.

* `port_id` - (Optional) The port id which this eip will associate with. If the value
    is "" or this not specified, the eip will be in unbind state.


The `bandwidth` block supports:

* `name` - (Required) The bandwidth name, which is a string of 1 to 64 characters
    that contain letters, digits, underscores (_), and hyphens (-).

* `size` - (Required) The bandwidth size. The value ranges from 1 to 300 Mbit/s.

* `share_type` - (Required) Whether the bandwidth is shared or exclusive. Changing
    this creates a new eip.

* `charge_mode` - (Optional) This is a reserved field. If the system supports charging
    by traffic and this field is specified, then you are charged by traffic for elastic
    IP addresses. Changing this creates a new eip.

## Attributes Reference

The following attributes are exported:

* `region` - See Argument Reference above.
* `publicip/type` - See Argument Reference above.
* `publicip/ip_address` - See Argument Reference above.
* `publicip/port_id` - See Argument Reference above.
* `bandwidth/name` - See Argument Reference above.
* `bandwidth/size` - See Argument Reference above.
* `bandwidth/share_type` - See Argument Reference above.
* `bandwidth/charge_mode` - See Argument Reference above.

## Import

EIPs can be imported using the `id`, e.g.

```
$ terraform import flexibleengine_vpc_eip_v1.eip_1 2c7f39f3-702b-48d1-940c-b50384177ee1
```
