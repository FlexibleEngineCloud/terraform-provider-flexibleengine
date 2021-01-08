---
subcategory: "VPC Endpoint"
layout: "flexibleengine"
page_title: "FlexibleEngine: flexibleengine_vpcep_approval"
description: |-
  Provides a resource to manage the VPC endpoint connections.
---

# flexibleengine\_vpcep\_approval

Provides a resource to manage the VPC endpoint connections.

## Example Usage

```hcl
variable "service_vpc_id" {}
variable "vm_port" {}
variable "vpc_id" {}
variable "network_id" {}

resource "flexibleengine_vpcep_service" "demo" {
  name        = "demo-service"
  server_type = "VM"
  vpc_id      = var.service_vpc_id
  port_id     = var.vm_port
  approval    = true

  port_mapping {
    service_port  = 8080
    terminal_port = 80
  }
}

resource "flexibleengine_vpcep_endpoint" "demo" {
  service_id  = flexibleengine_vpcep_service.demo.id
  vpc_id      = var.vpc_id
  network_id  = var.network_id
  enable_dns  = true

  lifecycle {
    # enable_dns and ip_address are not assigned until connecting to the service
    ignore_changes = [enable_dns, ip_address]
  }
}

resource "flexibleengine_vpcep_approval" "approval" {
  service_id = flexibleengine_vpcep_service.demo.id
  endpoints  = [flexibleengine_vpcep_endpoint.demo.id]
}
```

## Argument Reference

The following arguments are supported:

* `service_id` (Required) - Specifies the ID of the VPC endpoint service. Changing this creates a new resource.

* `endpoints` (Required) - Specifies the list of VPC endpoint IDs which accepted to connect to VPC endpoint service.
    The VPC endpoints will be rejected when the resource was destroyed.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The unique ID in UUID format which equals to the ID of the VPC endpoint service.

* `region` - The region in which to obtain the VPC endpoint service.

* `connections` - An array of VPC endpoints connect to the VPC endpoint service. Structure is documented below.
    - `endpoint_id` - The unique ID of the VPC endpoint.
    - `packet_id` - The packet ID of the VPC endpoint.
    - `domain_id` - The user's domain ID.
    - `status` - The connection status of the VPC endpoint.

## Timeouts
This resource provides the following timeouts configuration options:
- `create` - Default is 10 minute.
- `delete` - Default is 3 minute.
