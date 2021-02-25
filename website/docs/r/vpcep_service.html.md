---
subcategory: "VPC Endpoint"
layout: "flexibleengine"
page_title: "FlexibleEngine: flexibleengine_vpcep_service"
description: |-
  Provides a resource to manage a VPC endpoint service resource.
---

# flexibleengine\_vpcep\_service

Provides a resource to manage a VPC endpoint service resource.

## Example Usage

```hcl
variable "vpc_id" {}
variable "vm_port" {}

resource "flexibleengine_vpcep_service" "demo" {
  name        = "demo-service"
  server_type = "VM"
  vpc_id      = var.vpc_id
  port_id     = var.vm_port

  port_mapping {
    service_port  = 8080
    terminal_port = 80
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` (Optional) - Specifies the name of the VPC endpoint service. The value contains a maximum of
    16 characters, including letters, digits, underscores (_), and hyphens (-).

* `vpc_id` (Required) - Specifies the ID of the VPC to which the backend resource of
    the VPC endpoint service belongs. Changing this creates a new VPC endpoint service.

* `server_type` (Required) - Specifies the backend resource type. The value can be **VM**, **VIP** or **LB**.
    Changing this creates a new VPC endpoint service.

* `port_id` (Required) - Specifies the ID for identifying the backend resource of the VPC endpoint service.
    - If the `server_type` is **VM**, the value is the NIC ID of the ECS where the VPC endpoint service is deployed. 
    - If the `server_type` is **VIP**, the value is the NIC ID of the physical server where virtual resources are created.
    - If the `server_type` is **LB**, the value is the ID of the port bound to the private IP address of the load balancer.

* `port_mapping` (Required) - Specified the port mappings opened to the VPC endpoint service.
    Structure is documented below.

* `approval` (Optional) - Specifies whether connection approval is required. The default value is false.

* `permissions` (Optional) - Specifies the list of accounts to access the VPC endpoint service.
    The record is in the `iam:domain::domain_id` format. *iam:domain::\** allows all users to access the VPC endpoint service.

* `tags` - (Optional) The key/value pairs to associate with the VPC endpoint service.

The `port_mapping` block supports:

* `protocol` - (Optional) Specifies the protocol used in port mappings.
    The value can be _TCP_ or _UDP_. The default value is _TCP_.

* `service_port` - (Optional) Specifies the port for accessing the VPC endpoint service.
    This port is provided by the backend service to provide services. The value ranges from 1 to 65535.

* `terminal_port` - (Optional) Specifies the port for accessing the VPC endpoint.
    This port is provided by the VPC endpoint, allowing you to access the VPC endpoint service.
    The value ranges from 1 to 65535.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The unique ID of the VPC endpoint service.

* `region` - The region in which to create the VPC endpoint service.

* `status` - The status of the VPC endpoint service. The value can be **available** or **failed**.

* `service_name` - The full name of the VPC endpoint service in the format: *region.name.id*.

* `service_type` - The type of the VPC endpoint service. Only **interface** can be configured.

* `connections` - An array of VPC endpoints connect to the VPC endpoint service. Structure is documented below.
    - `endpoint_id` - The unique ID of the VPC endpoint.
    - `packet_id` - The packet ID of the VPC endpoint.
    - `domain_id` - The user's domain ID.
    - `status` - The connection status of the VPC endpoint.

## Import

VPC endpoint services can be imported using the `id`, e.g.

```
$ terraform import flexibleengine_vpcep_service.test_service 950cd3ba-9d0e-4451-97c1-3e97dd515d46
```
