---
subcategory: "Scalable File Service (SFS)"
---

# flexibleengine_sfs_turbos

Use this data source to get the list of the available SFS turbos.

## Example Usage

```hcl
variable "sfs_turbo_name" {}

data "flexibleengine_sfs_turbos" "test" {
  name = var.sfs_turbo_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) The region in which to query the SFS turbo file systems. If omitted, the provider-level
  region will be used.

* `name` - (Optional, String) Specifies the name of the SFS turbo file system.

* `size` - (Optional, Int) Specifies the capacity of the SFS turbo file system, in GB.
  The value ranges from `500` to `32,768`, and must be larger than `10,240` for an enhanced file system.

* `share_proto` - (Optional, String) Specifies the protocol of the SFS turbo file system. The valid value is **NFS**.

* `share_type` - (Optional, String) Specifies the type of the SFS turbo file system.
  The valid values are **STANDARD** and **PERFORMANCE**.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project id of the SFS turbo file systems
  resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `turbos` - The list of the SFS turbo file systems. The [turbos](#sfs_turbo) object structure is documented below.

<a name="sfs_turbo"></a>
The `turbos` block supports:

* `id` - The resource ID of the SFS turbo file system.

* `name` - The name of the SFS turbo file system.

* `size` - The capacity of the SFS turbo file system.

* `share_proto` - The protocol of the SFS turbo file system.

* `share_type` - The type of the SFS turbo file system.

* `version` - The version of the SFS turbo file system.

* `enhanced` - Whether the SFS turbo file system is enhanced.

* `availability_zone` - The availability zone where the SFS turbo file system is located.

* `available_capacity` - The available capacity of the SFS turbo file system, in GB.

* `export_location` - The mount point of the SFS turbo file system.

* `crypt_key_id` - The ID of a KMS key to encrypt the SFS turbo file system.

* `vpc_id` - The ID of the VPC to which the SFS turbo belongs.

* `subnet_id` - The ID of the VPC Subnet to which the SFS turbo belongs.

* `security_group_id` - The ID of the security group to which the SFS turbo belongs.

* `enterprise_project_id` - The enterprise project id to which the SFS turbo belongs.
