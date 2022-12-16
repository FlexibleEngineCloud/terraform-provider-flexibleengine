---
subcategory: "Virtual Private Cloud (VPC)"
description: ""
page_title: "flexibleengine_vpc_v1"
---

# flexibleengine_vpc_v1

Manages a VPC resource within FlexibleEngine.

## Example Usage

```hcl

variable "vpc_name" {
  default = "flexibleengine_vpc"
}

variable "vpc_cidr" {
  default = "192.168.0.0/16"
}

resource "flexibleengine_vpc_v1" "vpc_v1" {
  name = var.vpc_name
  cidr = var.vpc_cidr
}

resource "flexibleengine_vpc_v1" "vpc_with_tags" {
  name = var.vpc_name
  cidr = var.vpc_cidr

  tags = {
    foo = "bar"
    key = "value"
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the VPC. If omitted, the
  provider-level region will be used. Changing this creates a new VPC resource.

* `name` - (Required, String) Specifies the name of the VPC. The name must be unique for a tenant. The value is a string
  of no more than 64 characters and can contain digits, letters, underscores (_), and hyphens (-).

* `cidr` - (Required, String) Specifies the range of available subnets in the VPC. The value ranges from 10.0.0.0/8 to
  10.255.255.0/24, 172.16.0.0/12 to 172.31.255.0/24, or 192.168.0.0/16 to 192.168.255.0/24.

* `description` - (Optional, String) Specifies supplementary information about the VPC. The value is a string of
  no more than 255 characters and cannot contain angle brackets (< or >).

* `secondary_cidr` - (Optional, String) Specifies the secondary CIDR block of the VPC.

  -> The argument cannot be imported into your Terraform state. And the following secondary CIDR blocks cannot be added
  to a VPC: 10.0.0.0/8, 172.16.0.0/12, and 192.168.0.0/16.
  [View the complete list of unsupported CIDR blocks](https://docs.prod-cloud-ocb.orange-business.com/usermanual/vpc/vpc_vpc_0007.html).

* `tags` - (Optional, Map) Specifies the key/value pairs to associate with the VPC.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The VPC ID in UUID format.

* `status` - The current status of the VPC. Possible values are as follows: CREATING, OK or ERROR.

## Import

VPCs can be imported using the `id`, e.g.

```shell
terraform import flexibleengine_vpc_v1.vpc_v1 7117d38e-4c8f-4624-a505-bd96b97d024c
```
