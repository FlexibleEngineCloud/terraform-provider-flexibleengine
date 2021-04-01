---
subcategory: "Virtual Private Cloud (VPC)"
---

# flexibleengine\_networking\_port_v2

Manages a V2 port resource within FlexibleEngine.

## Example Usage

### Basic Usage
```hcl
resource "flexibleengine_networking_network_v2" "network_1" {
  name           = "network_1"
  admin_state_up = "true"
}

resource "flexibleengine_networking_port_v2" "port_1" {
  name           = "port_1"
  network_id     = flexibleengine_networking_network_v2.network_1.id
  admin_state_up = "true"
}
```

### Port With allowed_address_pairs
```hcl
resource "flexibleengine_networking_network_v2" "network_1" {
  name           = "network_1"
  admin_state_up = "true"
}

resource "flexibleengine_networking_port_v2" "port_1" {
  name           = "port_1"
  network_id     = flexibleengine_networking_network_v2.net.id
  admin_state_up = "true"

  allowed_address_pairs {
    ip_address = "192.168.0.0/24"
  }
}

```

## Argument Reference

The following arguments are supported:

* `region` - (Optional) The region in which to obtain the V2 networking client.
    A networking client is needed to create a port. If omitted, the
    `region` argument of the provider is used. Changing this creates a new
    port.

* `name` - (Optional) A unique name for the port. Changing this
    updates the `name` of an existing port.

* `network_id` - (Required) The ID of the network to attach the port to. Changing
    this creates a new port.

* `admin_state_up` - (Optional) Administrative up/down status for the port
    (must be "true" or "false" if provided). Changing this updates the
    `admin_state_up` of an existing port.

* `mac_address` - (Optional) Specify a specific MAC address for the port. Changing
    this creates a new port.

* `tenant_id` - (Optional) The owner of the Port. Required if admin wants
    to create a port for another tenant. Changing this creates a new port.

* `device_owner` - (Optional) The device owner of the Port. Changing this creates
    a new port.

* `security_group_ids` - (Optional) A list of security group IDs to apply to the
    port. The security groups must be specified by ID and not name (as opposed
    to how they are configured with the Compute Instance).

* `device_id` - (Optional) The ID of the device attached to the port. Changing this
    creates a new port.

* `fixed_ip` - (Optional) An array of desired IPs for this port. The structure is
    described below.

* `allowed_address_pairs` - (Optional) An array of IP/MAC Address pairs of additional IP
    addresses that can be active on this port. The structure is described below.

* `value_specs` - (Optional) Map of additional options.

The `fixed_ip` block supports:

* `subnet_id` - (Required) Subnet in which to allocate IP address for this port.

* `ip_address` - (Optional) IP address desired in the subnet for this port. If
    you don't specify `ip_address`, an available IP address from the specified
    subnet will be allocated to this port.

The `allowed_address_pairs` block supports:

* `ip_address` - (Required) The additional IP address. The value can be an IP Address or a CIDR,
    and can not be *0.0.0.0*. A server connected to the port can send a packet with source address
    which matches one of the specified allowed address pairs.
    It is recommended to configure an independent security group for the port if a large CIDR
    block (subnet mask less than 24) is configured.

* `mac_address` - (Optional) The additional MAC address.

## Attributes Reference

The following attributes are exported:

* `region` - See Argument Reference above.
* `admin_state_up` - See Argument Reference above.
* `mac_address` - See Argument Reference above.
* `tenant_id` - See Argument Reference above.
* `device_owner` - See Argument Reference above.
* `security_group_ids` - See Argument Reference above.
* `device_id` - See Argument Reference above.
* `fixed_ip` - See Argument Reference above.
* `all fixed_ips` - The collection of Fixed IP addresses on the port in the
  order returned by the Network v2 API.

## Import

Ports can be imported using the `id`, e.g.

```
$ terraform import flexibleengine_networking_port_v2.port_1 eae26a3e-1c33-4cc1-9c31-0cd729c438a1
```

## Notes

### Ports and Instances

There are some notes to consider when connecting Instances to networks using
Ports. Please see the `flexibleengine_compute_instance_v2` documentation for further
documentation.
