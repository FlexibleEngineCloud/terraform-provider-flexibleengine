---
subcategory: "Cloud Server Backup Service (CSBS)"
---

# flexibleengine_csbs_backup_policy_v1

The FlexibleEngine CSBS Backup Policy data source allows access of backup Policy resources.

## Example Usage

```hcl
variable "policy_id" {}

data "flexibleengine_csbs_backup_policy_v1" "csbs_policy" {
  id = var.policy_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) The region in which to query the data source. If omitted, the provider-level region
  will be used.

* `id` - (Optional, String) Specifies the ID of backup policy.

* `name` - (Optional, String) Specifies the backup policy name.

* `status` - (Optional, String) Specifies the backup policy status.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `description` - Specifies the backup policy description.

* `provider_id` - Provides the Backup provider ID.

* `common` - General backup policy parameters, which are blank by default.

* `scheduled_operation` -  Backup plan information.
  The [scheduled_operation](#csbs_scheduled_operation) object structure is documented below.

* `resource` - Backup Object. The [resource](#csbs_resource) object structure is documented below.

<a name="csbs_scheduled_operation"></a>
The `scheduled_operation` block supports:

* `name` - Specifies Scheduling period name.

* `description` - Specifies Scheduling period description.

* `enabled` - Specifies whether the scheduling period is enabled.

* `max_backups` - Specifies maximum number of backups that can be automatically created for a backup object.

* `retention_duration_days` - Specifies duration of retaining a backup, in days.

* `permanent` - Specifies whether backups are permanently retained.

* `trigger_pattern` - Specifies Scheduling policy of the scheduler.

* `operation_type` - Specifies Operation type, which can be backup.

* `id` -  Specifies Scheduling period ID.

* `trigger_id` -  Specifies Scheduler ID.

* `trigger_name` -  Specifies Scheduler name.

* `trigger_type` -  Specifies Scheduler type.

<a name="csbs_resource"></a>
The `resource` block supports:

* `id` - Specifies the ID of the object to be backed up.

* `type` - Entity object type of the backup object.

* `name` - Specifies backup object name.
  