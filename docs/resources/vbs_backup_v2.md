---
subcategory: "Volume Backup Service (VBS)"
---

# flexibleengine_vbs_backup_v2

Provides an VBS Backup resource.
 
# Example Usage

 ```hcl
variable "backup_name" {}

variable "volume_id" {}
 
resource "flexibleengine_vbs_backup_v2" "mybackup" {
  volume_id = "${var.volume_id}"
  name = "${var.backup_name}"
}
 ```

# Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the vbs backup. Changing the parameter will create new resource.

* `volume_id` - (Required) The id of the disk to be backed up. Changing the parameter will create new resource.

* `snapshot_id` - (Optional) The snapshot id of the disk to be backed up. Changing the parameter will create new resource.

* `description` - (Optional) The description of the vbs backup. Changing the parameter will create new resource.


# Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The id of the vbs backup.

* `container` - The container of the backup.

* `status` - The status of the VBS backup.

* `availability_zone` - The AZ where the backup resides.

* `size` - The size of the vbs backup.

* `service_metadata` - The metadata of the vbs backup.

# Import

VBS Backup can be imported using the `backup id`, e.g.

```
 $ terraform import flexibleengine_vbs_backup_v2.mybackup 4779ab1c-7c1a-44b1-a02e-93dfc361b32d
```