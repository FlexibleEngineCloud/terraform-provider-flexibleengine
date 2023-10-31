---
subcategory: "Relational Database Service (RDS)"
---

# flexibleengine_rds_backups

Use this data source to get the list of RDS backups.

## Example Usage

```hcl
variable "instance_id" {}

data "flexibleengine_rds_backups" "test" {
  instance_id = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the DB instance ID.

* `name` - (Optional, String) Specifies the backup name.

* `backup_id` - (Optional, String) Specifies the backup ID.

* `backup_type` - (Optional, String) Specifies the backup type. The options are as follows:
  - **auto**: Automated full backup.
  - **manual**: Manual full backup.
  - **fragment**: Differential full backup.
  - **incremental**: Automated incremental backup.

* `begin_time` - (Optional, String) Specifies the start time for obtaining the backup list.
  The format of the start time is "yyyy-mm-ddThh:mm:ssZ".

* `end_time` - (Optional, String) Specifies the end time for obtaining the backup list.
  The format of the end time is "yyyy-mm-ddThh:mm:ssZ" and the end time must be later than the start time.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `backups` - Backup list. For details, see Data structure of the Backup field.
  The [backups](#rds_backups) structure is documented below.

<a name="rds_backups"></a>
The `backups` block supports:

* `id` - Backup ID.

* `instance_id` - RDS instance ID.

* `name` - Backup name.

* `type` - Backup type. The options are as follows:
  - **auto**: Automated full backup.
  - **manual**: Manual full backup.
  - **fragment**: Differential full backup.
  - **incremental**: Automated incremental backup.

* `size` - Backup size in KB.

* `status` - Backup status. The options are as follows:
  - **BUILDING**: Backup in progress.
  - **COMPLETED**: Backup completed.
  - **FAILED**: Backup failed.
  - **DELETING**: Backup being deleted.

* `begin_time` - Backup start time in the "yyyy-mm-ddThh:mm:ssZ" format.

* `end_time` - Backup end time in the "yyyy-mm-ddThh:mm:ssZ" format.

* `associated_with_ddm` - Whether a DDM instance has been associated.

* `datastore` - The database information. The [datastore](#rds_datastore) structure is documented below.

* `databases` - Database been backed up. The [databases](#rds_databases) structure is documented below.

<a name="rds_datastore"></a>
The `datastore` block supports:

* `type` - DB engine. The value can be: **MySQL**, **PostgreSQL**, **SQL Server**, **MariaDB**.

* `version` - DB engine version.

<a name="rds_databases"></a>
The `rds_databases` block supports:

* `name` - Database to be backed up for Microsoft SQL Server.
