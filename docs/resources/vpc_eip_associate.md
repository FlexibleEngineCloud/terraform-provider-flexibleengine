---
subcategory: "Elastic IP (EIP)"
description: ""
page_title: "flexibleengine_vpc_eip_associate"
---

# flexibleengine_vpc_eip_associate

Associates an EIP to a specified IP address or port.

## Example Usage

### Associate with a fixed IP

```hcl
resource "flexibleengine_vpc_eip_associate" "example_eip_associated" {
  public_ip  = flexibleengine_vpc_eip.example_eip.address
  network_id = flexibleengine_vpc_subnet_v1.example_subnet.id
  fixed_ip   = "192.168.0.100"
}
```

### Associate with a port

```hcl
data "flexibleengine_networking_port" "example_port" {
  network_id = flexibleengine_vpc_subnet_v1.example_subnet.id
  fixed_ip   = "192.168.0.100"
}

resource "flexibleengine_vpc_eip" "example_eip" {
  publicip {
    type = "5_bgp"
  }
  bandwidth {
    name        = "test"
    size        = 5
    share_type  = "PER"
    charge_mode = "traffic"
  }
}

resource "flexibleengine_vpc_eip_associate" "associated" {
  public_ip = flexibleengine_vpc_eip.example_eip.address
  port_id   = data.flexibleengine_networking_port.example_port.id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to associate the EIP. If omitted, the provider-level
  region will be used. Changing this creates a new resource.

* `public_ip` - (Required, String, ForceNew) Specifies the EIP address to associate. Changing this creates a new resource.

* `fixed_ip` - (Optional, String, ForceNew) Specifies a private IP address to associate with the EIP.
  Changing this creates a new resource.

* `network_id` - (Optional, String, ForceNew) Specifies the ID of the VPC Subnet to which the **fixed_ip** belongs.
  It is mandatory when `fixed_ip` is set. Changing this creates a new resource.

* `port_id` - (Optional, String, ForceNew) Specifies an existing port ID to associate with the EIP.
  This parameter and `fixed_ip` are alternative. Changing this creates a new resource.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in UUID format.
* `mac_address` - The MAC address of the private IP.
* `status` - The status of EIP, should be **BOUND**.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 5 minute.
* `delete` - Default is 5 minute.

## Import

EIP associations can be imported using the `id` of the EIP, e.g.

```shell
terraform import flexibleengine_vpc_eip_associate.eip 2c7f39f3-702b-48d1-940c-b50384177ee1
```
