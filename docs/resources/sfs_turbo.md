---
subcategory: "Scalable File Service (SFS)"
description: ""
page_title: "flexibleengine_sfs_turbo"
---

# flexibleengine_sfs_turbo

Provides an Shared File System (SFS) Turbo resource.

## Example Usage

```hcl
variable "test_az" {}

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

resource "flexibleengine_sfs_turbo" "sfs-turbo-1" {
  name        = "sfs-turbo-1"
  size        = 500
  share_proto = "NFS"
  vpc_id      = flexibleengine_vpc_v1.example_vpc.id
  subnet_id   = flexibleengine_vpc_subnet_v1.example_subnet.id
  security_group_id = flexibleengine_networking_secgroup_v2.example_secgroup.id
  availability_zone = var.test_az
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of an SFS Turbo file system. The value contains 4 to 64
  characters and must start with a letter. Changing this will create a new resource.

* `size` - (Required) Specifies the capacity of a common file system, in GB. The value ranges from 500 to 32768.

* `share_proto` - (Optional) Specifies the protocol for sharing file systems. The valid value is NFS.
  Changing this will create a new resource.

* `share_type` - (Optional) Specifies the file system type. The valid values are STANDARD and PERFORMANCE.
  Changing this will create a new resource.

* `availability_zone` - (Required) Specifies the availability zone where the file system is located.
  Changing this will create a new resource.

* `vpc_id` - (Required) Specifies the VPC ID. Changing this will create a new resource.

* `subnet_id` - (Required) Specifies the ID of the VPC Subnet. Changing this will create a new resource.

* `security_group_id` - (Required) Specifies the security group ID. Changing this will create a new resource.

* `crypt_key_id` - (Optional) Specifies the ID of a KMS key to encrypt the file system.
  Changing this will create a new resource.

-> **NOTE:**
  SFS Turbo will create two private IP addresses and one virtual IP address under the subnet you specified.
  To ensure normal use, SFS Turbo will enable the inbound rules for ports *111*, *445*, *2049*, *2051*, *2052*,
  and *20048* in the security group you specified.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The UUID of the SFS Turbo file system.

* `region` - The region of the SFS Turbo file system.

* `status` - The status of the SFS Turbo file system.

* `version` - The version ID of the SFS Turbo file system.

* `export_location` - Tthe mount point of the SFS Turbo file system.

* `available_capacity` - The available capacity of the SFS Turbo file system in the unit of GB.

## Import

SFS Turbo can be imported using the `id`, e.g.

```shell
terraform import flexibleengine_sfs_turbo 1e3d5306-24c9-4316-9185-70e9787d71ab
```
