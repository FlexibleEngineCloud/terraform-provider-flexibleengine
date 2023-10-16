---
subcategory: "Elastic Cloud Server (ECS)"
description: ""
page_title: "flexibleengine_compute_interface_attach_v2"
---

# flexibleengine_compute_interface_attach_v2

Attaches a Network Interface (a Port) to an Instance using the FlexibleEngine
Compute (Nova) v2 API.

## Example Usage

### Basic Attachment

```hcl
resource "flexibleengine_vpc_v1" "example_vpc" {
  name = "example-vpc"
  cidr = "192.168.0.0/16"
}

resource "flexibleengine_vpc_subnet_v1" "example_subnet" {
  name       = "example-vpc-subnet"
  cidr       = "192.168.0.0/24"
  gateway_ip = "192.168.0.1"
  vpc_id     = flexibleengine_vpc_v1.example_vpc.id
}

resource "flexibleengine_compute_instance_v2" "example_instance" {
  name            = "example-instance"
  security_groups = ["default"]
}

resource "flexibleengine_compute_interface_attach_v2" "example_interface_attach" {
  instance_id = flexibleengine_compute_instance_v2.example_instance.id
  network_id  = flexibleengine_vpc_subnet_v1.example_subnet.id
}

```

### Attachment Specifying a Fixed IP

```hcl
resource "flexibleengine_compute_interface_attach_v2" "example_interface_attach" {
  instance_id = flexibleengine_compute_instance_v2.example_instance.id
  network_id  = flexibleengine_vpc_subnet_v1.example_subnet.id
  fixed_ip    = "10.0.10.10"
}

```

### Attachment Using an Existing Port

```hcl
resource "flexibleengine_networking_port_v2" "example_port" {
  name           = "port_1"
  network_id     = flexibleengine_vpc_subnet_v1.example_subnet.id
  admin_state_up = "true"
}

resource "flexibleengine_compute_interface_attach_v2" "example_interface_attach" {
  instance_id = flexibleengine_compute_instance_v2.example_instance.id
  port_id     = flexibleengine_networking_port_v2.example_port.id
}

```

### Attaching Multiple Interfaces

```hcl
resource "flexibleengine_networking_port_v2" "example_ports" {
  count          = 2
  name           = format("port-%02d", count.index + 1)
  network_id     = flexibleengine_vpc_subnet_v1.example_subnet.id
  admin_state_up = "true"
}

resource "flexibleengine_compute_interface_attach_v2" "example_attachments" {
  count          = 2
  instance_id = flexibleengine_compute_instance_v2.example_instance.id
  port_id     = flexibleengine_networking_port_v2.example_ports.*.id[count.index]
}
```

Note that the above example will not guarantee that the ports are attached in
a deterministic manner. The ports will be attached in a seemingly random
order.

If you want to ensure that the ports are attached in a given order, create
explicit dependencies between the ports, such as:

```hcl
resource "flexibleengine_networking_port_v2" "example_ports" {
  count          = 2
  name           = format("port-%02d", count.index + 1)
  network_id     = flexibleengine_vpc_subnet_v1.example_subnet.id
  admin_state_up = "true"
}

resource "flexibleengine_compute_interface_attach_v2" "example_interface_attach_1" {
  instance_id = flexibleengine_compute_instance_v2.instance_1.id
  port_id     = flexibleengine_networking_port_v2.example_ports.*.id[0]
}

resource "flexibleengine_compute_interface_attach_v2" "example_interface_attach_2" {
  instance_id = flexibleengine_compute_instance_v2.instance_1.id
  port_id     = flexibleengine_networking_port_v2.example_ports.*.id[1]
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the interface attachment.
  If omitted, the `region` argument of the provider is used. Changing this creates a new attachment.

* `instance_id` - (Required, String, ForceNew) The ID of the Instance to attach the Port or Network to.

* `port_id` - (Optional, String, ForceNew) The ID of the Port to attach to an Instance.
  This option and `network_id` are mutually exclusive.

* `network_id` - (Optional, String, ForceNew) The ID of the Network to attach to an Instance.
  A port will be created automatically. This option and `port_id` are mutually exclusive.

* `fixed_ip` - (Optional, String, ForceNew) An IP address to associate with the port.
  This option cannot be used with port_id. You must specify a network_id.
  The IP address must lie in a range on the supplied network.

## Attribute Reference

All the arguments above can also be exported attributes.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.
* `delete` - Default is 10 minutes.

## Import

Interface Attachments can be imported using the Instance ID and Port ID
separated by a slash, e.g.

```shell
terraform import flexibleengine_compute_interface_attach_v2.ai_1 89c60255-9bd6-460c-822a-e2b959ede9d2/45670584-225f-46c3-b33e-6707b589b666
```
