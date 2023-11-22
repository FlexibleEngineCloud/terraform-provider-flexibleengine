---
subcategory: "Relational Database Service (RDS)"
---

# flexibleengine_rds_backup

Manages a RDS manual backup resource within FlexibleEngine.  

## Example Usage

```hcl
variable "instance_id" {}
variable "backup_name" {}

resource "flexibleengine_rds_backup" "test" {
  instance_id = var.instance_id
  name        = var.backup_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `name` - (Required, String, ForceNew) Specifies the name of the resource backup.
  The valid length is limited from can contain 4 to 64 characters long, start with a letter, and contain only
  letters (case-sensitive),digits, hyphens (-), and underscores (_).
  Changing this parameter will create a new resource.

* `instance_id` - (Required, String, ForceNew) Specifies the instance id.
  Changing this parameter will create a new resource.

* `description` - (Optional, String, ForceNew) The description about the backup.  
  It contains a maximum of 256 characters and cannot contain the following special characters: **> ! < " & ' =**.
  Changing this parameter will create a new resource.

* `databases` - (Optional, List, ForceNew) List of self-built Microsoft SQL Server databases that are partially
  backed up. (Only Microsoft SQL Server supports partial backups). The [databases](#rds_databases) structure is
  documented below. Changing this parameter will create a new resource.

<a name="rds_databases"></a>
The `databases` block supports:

* `name` - (Required, String, ForceNew) Database to be backed up for Microsoft SQL Server.
  Changing this parameter will create a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `begin_time` - Backup start time in the "yyyy-mm-ddThh:mm:ssZ" format.

* `end_time` - Backup end time in the "yyyy-mm-ddThh:mm:ssZ" format.

* `status` - Backup status. The options are as follows:
  + **BUILDING**: Backup in progress.
  + **COMPLETED**: Backup completed.
  + **FAILED**: Backup failed.
  + **DELETING**: Backup being deleted.

* `size` - Backup size in KB.

* `associated_with_ddm` - Whether a DDM instance has been associated.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 30 minutes.
* `delete` - Default is 30 minutes.

## Import

The rds manual backup can be imported using the instance ID and the backup ID separated by a slash, e.g.:

```shell
terraform import flexibleengine_rds_backup.test 1ce123456a00f2591fabc00385ff1235/0ce123456a00f2591fabc00385ff1234
```
