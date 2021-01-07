---
subcategory: "VPC Endpoint"
layout: "flexibleengine"
page_title: "FlexibleEngine: flexibleengine_vpcep_endpoint"
description: |-
  Provides a resource to manage a VPC endpoint resource.
---

# flexibleengine\_vpcep\_endpoint

Provides a resource to manage a VPC endpoint resource.

## Example Usage

### Access to the public service

```hcl
variable "vpc_id" {}
variable "network_id" {}

data "flexibleengine_vpcep_public_services" "cloud_service" {
  service_name = "dns"
}

resource "flexibleengine_vpcep_endpoint" "myendpoint" {
  service_id       = data.flexibleengine_vpcep_public_services.cloud_service.services[0].id
  vpc_id           = var.vpc_id
  network_id       = var.network_id
  enable_dns       = true
  enable_whitelist = true
  whitelist        = ["192.168.0.0/24"]
}
```

### Access to the private service

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
}
```

## Argument Reference

The following arguments are supported:

* `service_id` (Required) - Specifies the ID of the VPC endpoint service.
    Changing this creates a new VPC endpoint.

* `vpc_id` (Required) - Specifies the ID of the VPC where the VPC endpoint is to be created.
    Changing this creates a new VPC endpoint.

* `network_id` (Required) - Specifies the network ID of the subnet in the VPC specified by `vpc_id`.
    Changing this creates a new VPC endpoint.

* `ip_address` (Optional) - Specifies the IP address for accessing the associated VPC endpoint service.
    Only IPv4 addresses are supported. Changing this creates a new VPC endpoint.

* `enable_dns` (Optional) - Specifies whether to create a private domain name. The default value is true.
    Changing this creates a new VPC endpoint.

* `enable_whitelist` (Optional) - Specifies whether to enable access control. The default value is false.
    Changing this creates a new VPC endpoint.

* `whitelist` (Optional) - Specifies the list of IP address or CIDR block which can be accessed to the VPC endpoint.
    Changing this creates a new VPC endpoint.

* `tags` - (Optional) The key/value pairs to associate with the VPC endpoint.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The unique ID of the VPC endpoint.

* `region` - The region in which to create the VPC endpoint.

* `status` - The status of the VPC endpoint. The value can be **accepted**, **pendingAcceptance** or **rejected**.

* `service_name` - The name of the VPC endpoint service.

* `service_type` - The type of the VPC endpoint service.

* `packet_id` - The packet ID of the VPC endpoint.

* `private_domain_name` -  The domain name for accessing the associated VPC endpoint service.
    This parameter is only available when enable_dns is set to true.

## Import

VPC endpoint can be imported using the `id`, e.g.

```
$ terraform import flexibleengine_vpcep_endpoint.test 828907cc-40c9-42fe-8206-ecc1bdd30060
```
