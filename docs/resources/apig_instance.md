---
subcategory: "API Gateway (Dedicated APIG)"
description: ""
page_title: "flexibleengine_apig_instance"
---

# flexibleengine_apig_instance

Manages an APIG dedicated instance resource within Flexibleengine.

## Example Usage

```hcl
variable "instance_name" {}

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

resource "flexibleengine_networking_secgroup_v2" "example_secgroup" {
  name        = "example-secgroup"
  description = "My neutron security group"
}

resource "flexibleengine_vpc_eip" "example_eip" {
  publicip {
    type = "5_bgp"
  }
  bandwidth {
    name       = "test"
    size       = 10
    share_type = "PER"
  }
}

data "flexibleengine_availability_zones" "az" {}

resource "flexibleengine_apig_instance" "example_apig_instance" {
  name                  = var.instance_name
  edition               = "BASIC"
  vpc_id                = flexibleengine_vpc_v1.example_vpc.id
  subnet_id             = flexibleengine_vpc_subnet_v1.example_subnet.id
  security_group_id     = flexibleengine_networking_secgroup_v2.example_secgroup.id
  maintain_begin        = "06:00:00"
  description           = "Created by script"
  bandwidth_size        = 3
  eip_id                = flexibleengine_vpc_eip.example_eip.id

  available_zones = [
    data.flexibleengine_availability_zones.az.names[0],
    data.flexibleengine_availability_zones.az.names[1],
  ]
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the APIG dedicated instance resource.
  If omitted, the provider-level region will be used. Changing this will create a new APIG dedicated instance resource.

* `name` - (Required, String) Specifies the name of the API dedicated instance. The API group name consists of 3 to 64
  characters, starting with a letter. Only letters, digits, and underscores (_) are allowed.

* `edition` - (Required, String, ForceNew) Specifies the edition of the APIG dedicated instance. The supported editions
  are as follows: BASIC, PROFESSIONAL, ENTERPRISE, PLATINUM. Changing this will create a new APIG dedicated instance
  resource.

* `vpc_id` - (Required, String, ForceNew) Specifies the ID of the VPC used to create the APIG dedicated instance.
  Changing this will create a new APIG dedicated instance resource.

* `subnet_id` - (Required, String, ForceNew) Specifies the ID of the VPC Subnet used to create the APIG dedicated
  instance. Changing this will create a new APIG dedicated instance resource.

* `security_group_id` - (Required, String) Specifies the ID of the security group to which the APIG dedicated instance
  belongs to.

* `available_zones` - (Required, List, ForceNew) Specifies an array of available zone names for the APIG dedicated
  instance. Changing this will create a new APIG dedicated instance resource.

* `description` - (Optional, String) Specifies the description about the APIG dedicated instance. The description
  contain a maximum of 255 characters and the angle brackets (< and >) are not allowed.

* `maintain_begin` - (Optional, String) Specifies a start time of the maintenance time window in the format 'xx:00:00'.
  The value of xx can be 02, 06, 10, 14, 18 or 22.

* `bandwidth_size` - (Optional, Int) Specifies the egress bandwidth size of the APIG dedicated instance. The range of
  valid value is from 1 to 2000.

* `eip_id` - (Optional, String) Specifies the eip ID associated with the APIG dedicated instance.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the APIG dedicated instance.
* `maintain_end` - End time of the maintenance time window, 4-hour difference between the start time and end time.
* `create_time` - Time when the APIG instance is created, in RFC-3339 format.
* `status` - Status of the APIG dedicated instance.
* `supported_features` - The supported features of the APIG dedicated instance.
* `egress_address` - The egress (nat) public ip address.
* `ingress_address` - The ingress eip address.
* `vpc_ingress_address` - The ingress private ip address of vpc.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 40 minute.
* `update` - Default is 10 minute.
* `delete` - Default is 10 minute.

## Import

APIG Dedicated Instances can be imported by their `id`, e.g.

```shell
terraform import flexibleengine_apig_instance.test de379eed30aa4d31a84f426ea3c7ef4e
```
