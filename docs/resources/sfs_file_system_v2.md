---
subcategory: "Scalable File Service (SFS)"
---

# flexibleengine_sfs_file_system_v2

Provides an Shared File System (SFS) resource.

## Example Usage

### basic example

```hcl
variable "share_name" {}
variable "share_description" {}
variable "vpc_id" {}

resource "flexibleengine_sfs_file_system_v2" "share-file" {
  name         = var.share_name
  size         = 100
  share_proto  = "NFS"
  access_level = "rw"
  access_to    = var.vpc_id
  description  = var.share_description
}
```

### sfs with data encryption

```hcl
variable "share_name" {}
variable "share_description" {}
variable "vpc_id" {}

resource "flexibleengine_kms_key_v1" "mykey" {
  key_alias    = "kms_sfs"
  pending_days = "7"
}

resource "flexibleengine_sfs_file_system_v2" "share-file" {
  name         = var.share_name
  size         = 100
  share_proto  = "NFS"
  access_level = "rw"
  access_to    = var.vpc_id
  description  = var.share_description

  metadata = {
    "#sfs_crypt_key_id"    = flexibleengine_kms_key_v1.mykey.id
    "#sfs_crypt_domain_id" = flexibleengine_kms_key_v1.mykey.domain_id
    "#sfs_crypt_alias"     = flexibleengine_kms_key_v1.mykey.key_alias
  }
}
```

## Argument Reference
The following arguments are supported:

* `size` - (Required) The size (GB) of the shared file system.

* `share_proto` - (Optional) The protocol for sharing file systems. The default value is NFS.

* `name` - (Optional) The name of the shared file system.

* `description` - (Optional) Describes the shared file system.

* `is_public` - (Optional) The level of visibility for the shared file system.

* `metadata` - (Optional) Metadata key and value pairs as a dictionary of strings.
    The supported metadata keys are "#sfs_crypt_key_id", "#sfs_crypt_domain_id" and "#sfs_crypt_alias",
    and the keys should be exist at the same time to enable the data encryption function.
    Changing this will create a new resource.

* `availability_zone` - (Optional) The availability zone name. Changing this parameter will create a new resource.

* `access_level` - (Optional) Specifies the access level of the shared file system. Possible values are *ro* (read-only)
    and *rw* (read-write). The default value is *rw* (read/write). Changing this will create a new access rule.

* `access_type` - (Optional) Specifies the type of the share access rule. The default value is *cert*.
    Changing this will create a new access rule.

* `access_to` - (Optional) Specifies the value that defines the access rule. The value contains 1 to 255 characters.
    Changing this will create a new access rule. The value varies according to the scenario:
    - Set the VPC ID in VPC authorization scenarios.
    - Set this parameter in IP address authorization scenario.

        - For an NFS shared file system, the value in the format of *VPC_ID#IP_address#priority#user_permission*.
        For example, 0157b53f-4974-4e80-91c9-098532bcaf00#2.2.2.2/16#100#all_squash,root_squash.

        - For a CIFS shared file system, the value in the format of *VPC_ID#IP_address#priority*.
        For example, 0157b53f-4974-4e80-91c9-098532bcaf00#2.2.2.2/16#0.

-> **NOTE:** If you want to create more access rules, please using [flexibleengine_sfs_access_rule_v2](https://www.terraform.io/docs/providers/flexibleengine/r/sfs_access_rule_v2.html).

## Attributes Reference
In addition to all arguments above, the following attributes are exported:

* `id` - The UUID of the shared file system.

* `status` - The status of the shared file system.

* `volume_type` - The volume type.

* `export_location` - The address for accessing the shared file system.

* `share_access_id` - The UUID of the share access rule.

* `access_rules_status` - The status of the share access rule.

* `access_rules` - All access rules of the shared file system. The object includes the following:
    - `access_rule_id` - The UUID of the share access rule.
    - `access_level` - The access level of the shared file system
    - `access_type` - The type of the share access rule.
    - `access_to` - The value that defines the access rule.
    - `status` - The status of the share access rule.

## Import

SFS can be imported using the `id`, e.g.

```
$ terraform import flexibleengine_sfs_file_system_v2 4779ab1c-7c1a-44b1-a02e-93dfc361b32d
```

**NOTE:** The `access_to`, `access_type` and `access_level` will not be imported.
Please importing them by [flexibleengine_sfs_access_rule_v2](https://www.terraform.io/docs/providers/flexibleengine/r/sfs_access_rule_v2.html).
