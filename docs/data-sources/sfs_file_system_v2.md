---
subcategory: "Scalable File Service (SFS)"
---

# flexibleengine_sfs_file_system_v2

Provides information about an Shared File System (SFS).

## Example Usage

```hcl
variable "share_name" {}

data "flexibleengine_sfs_file_system_v2" "shared_file"
{
  name = var.share_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) The region in which to query the data source. If omitted, the provider-level region
  will be used.

* `name` - (Optional, String) The name of the shared file system.

* `id` - (Optional, String) The UUID of the shared file system.

* `status` - (Optional, String) The status of the shared file system.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `availability_zone` - The availability zone name.

* `description` - The description of the shared file system.

* `project_id` - The project ID of the operating user.

* `size` - The size (GB) of the shared file system.

* `share_type` - The storage service type for the shared file system, such as high-performance storage (composed of SSDs)
  or large-capacity storage (composed of SATA disks).

* `host` - The host name of the shared file system.

* `is_public` - The level of visibility for the shared file system.

* `share_proto` - The protocol for sharing file systems.

* `volume_type` - The volume type.

* `metadata` - Metadata key and value pairs as a dictionary of strings.

* `export_location` - The path for accessing the shared file system.

* `export_locations` - The list of mount locations.

* `access_level` - The level of the access rule.

* `access_type` - The type of the share access rule.

* `access_to` - The access that the back end grants or denies.

* `state` - The status of the access rule.

* `share_access_id` - The UUID of the share access rule.

* `mount_id` - The UUID of the mount location of the shared file system.

* `share_instance_id` - The access that the back end grants or denies.

* `preferred` - Identifies which mount locations are most efficient and are used preferentially
  when multiple mount locations exist.
