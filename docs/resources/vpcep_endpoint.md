---
subcategory: "VPC Endpoint (VPCEP)"
description: ""
page_title: "flexibleengine_vpcep_endpoint"
---

# flexibleengine_vpcep_endpoint

Provides a resource to manage a VPC endpoint resource.

## Example Usage

### Access to the public service

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

data "flexibleengine_vpcep_public_services" "cloud_service" {
  service_name = "dns"
}

resource "flexibleengine_vpcep_endpoint" "myendpoint" {
  service_id       = data.flexibleengine_vpcep_public_services.cloud_service.services[0].id
  vpc_id           = flexibleengine_vpc_v1.example_vpc.id
  network_id       = flexibleengine_vpc_subnet_v1.example_subnet.id
  enable_dns       = true
  enable_whitelist = true
  whitelist        = ["192.168.0.0/24"]
}
```

### Access to the private service

```hcl
variable "vm_port" {}

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

resource "flexibleengine_vpcep_service" "demo" {
  name        = "demo-service"
  server_type = "VM"
  vpc_id      = flexibleengine_vpc_v1.example_vpc_service.id
  port_id     = var.vm_port

  port_mapping {
    service_port  = 8080
    terminal_port = 80
  }
}

resource "flexibleengine_vpcep_endpoint" "demo" {
  service_id  = flexibleengine_vpcep_service.demo.id
  vpc_id      = flexibleengine_vpc_v1.example_vpc.id
  network_id  = flexibleengine_vpc_subnet_v1.example_subnet.id
  enable_dns  = true
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the VPC endpoint.
  If omitted, the provider-level region will be used. Changing this creates a new VPC endpoint.

* `service_id` - (Required, String, ForceNew) Specifies the ID of the VPC endpoint service.
  Changing this creates a new VPC endpoint.

* `vpc_id` - (Required, String, ForceNew) Specifies the ID of the VPC where the VPC endpoint is to be created.
  Changing this creates a new VPC endpoint.

* `network_id` - (Required, String, ForceNew) Specifies the network ID of the subnet in the VPC specified by `vpc_id`.
  Changing this creates a new VPC endpoint.

* `ip_address` - (Optional, String, ForceNew) Specifies the IP address for accessing the associated VPC endpoint service.
  Only IPv4 addresses are supported. Changing this creates a new VPC endpoint.

* `enable_dns` - (Optional, Bool) Specifies whether to create a private domain name. The default value is true.

* `enable_whitelist` (Optional, Bool) - Specifies whether to enable access control. The default value is false.

* `whitelist` (Optional, List, ForceNew) - Specifies the list of IP address or CIDR block,
  which can be accessed to the VPC endpoint. Changing this creates a new VPC endpoint.

* `description` - (Optional, String, ForceNew) Specifies the description of the VPC endpoint.
  Changing this creates a new VPC endpoint.

* `tags` - (Optional, Map) The key/value pairs to associate with the VPC endpoint.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The unique ID of the VPC endpoint.

* `status` - The status of the VPC endpoint. The value can be **accepted**, **pendingAcceptance** or **rejected**.

* `service_name` - The name of the VPC endpoint service.

* `service_type` - The type of the VPC endpoint service.

* `packet_id` - The packet ID of the VPC endpoint.

* `private_domain_name` -  The domain name for accessing the associated VPC endpoint service.
    This parameter is only available when enable_dns is set to true.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.
* `delete` - Default is 10 minutes.

## Import

VPC endpoint can be imported using the `id`, e.g.

```shell
terraform import flexibleengine_vpcep_endpoint.test 828907cc-40c9-42fe-8206-ecc1bdd30060
```
