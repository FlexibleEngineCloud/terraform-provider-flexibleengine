---
subcategory: "Elastic Cloud Server (ECS)"
description: ""
page_title: "flexibleengine_compute_floatingip_associate_v2"
---

# flexibleengine_compute_floatingip_associate_v2

Associate a floating IP to an instance. This can be used instead of the
`floating_ip` options in `flexibleengine_compute_instance_v2`.

## Example Usage

### Associate with EIP

```hcl
resource "flexibleengine_compute_instance_v2" "instance_1" {
  name            = "instance_1"
  image_id        = "ad091b52-742f-469e-8f3c-fd81cadf0743"
  flavor_id       = 3
  key_pair        = "my_key_pair_name"
  security_groups = ["default"]
}

resource "flexibleengine_vpc_eip" "eip_1" {
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

resource "flexibleengine_compute_floatingip_associate_v2" "fip_1" {
  floating_ip = flexibleengine_vpc_eip.eip_1.publicip.0.ip_address
  instance_id = flexibleengine_compute_instance_v2.instance_1.id
}
```

### Explicitly set the network to attach to

```hcl
resource "flexibleengine_compute_instance_v2" "instance_1" {
  name            = "instance_1"
  image_id        = "ad091b52-742f-469e-8f3c-fd81cadf0743"
  flavor_id       = 3
  key_pair        = "my_key_pair_name"
  security_groups = ["default"]

  network {
    name = "my_network"
  }

  network {
    name = "default"
  }
}

esource "flexibleengine_vpc_eip" "eip_1" {
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

resource "flexibleengine_compute_floatingip_associate_v2" "fip_1" {
  floating_ip = flexibleengine_vpc_eip.eip_1.publicip.0.ip_address
  instance_id = flexibleengine_compute_instance_v2.instance_1.id
  fixed_ip    = flexibleengine_compute_instance_v2.instance_1.network.1.fixed_ip_v4
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional) The region in which to obtain the V2 Compute client.
    Keypairs are associated with accounts, but a Compute client is needed to
    create one. If omitted, the `region` argument of the provider is used.
    Changing this creates a new floatingip_associate.

* `floating_ip` - (Required) The floating IP to associate.

* `instance_id` - (Required) The instance to associte the floating IP with.

* `fixed_ip` - (Optional) The specific IP address to direct traffic to.

## Attributes Reference

The following attributes are exported:

* `region` - See Argument Reference above.
* `floating_ip` - See Argument Reference above.
* `instance_id` - See Argument Reference above.
* `fixed_ip` - See Argument Reference above.

## Import

This resource can be imported by specifying `floating_ip`, `instance_id` and `fixed_ip`, separated
by a forward slash:

```shell
terraform import flexibleengine_compute_floatingip_associate_v2.fip_1 <floating_ip>/<instance_id>/<fixed_ip>
```
