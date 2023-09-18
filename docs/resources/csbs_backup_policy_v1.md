---
subcategory: "Cloud Server Backup Service (CSBS)"
description: ""
page_title: "flexibleengine_csbs_backup_policy_v1"
---

# flexibleengine_csbs_backup_policy_v1

Provides a FlexibleEngine Backup Policy of Resources.

## Example Usage

 ```hcl
 variable "name" { }
 variable "id" { }
 variable "resource_name" { }
 
 resource "flexibleengine_csbs_backup_policy_v1" "backup_policy_v1" {
   name  = var.name
   resource {
     id   = var.id
     type = "OS::Nova::Server"
     name = var.resource_name
   }
   scheduled_operation {
     enabled         = true
     operation_type  = "backup"
     trigger_pattern = "BEGIN:VCALENDAR\r\nBEGIN:VEVENT\r\nRRULE:FREQ=WEEKLY;BYDAY=TH;BYHOUR=12;BYMINUTE=27\r\nEND:VEVENT\r\nEND:VCALENDAR\r\n"
   }
 }
 ```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the CSBS backup policy resource.
  If omitted, the provider-level region will be used. Changing this will create a new CSBS backup policy resource.

* `name` - (Required, String) Specifies the name of backup policy. The value consists of 1 to 255 characters and
  can contain only letters, digits, underscores (_), and hyphens (-).

* `description` - (Optional, String) Backup policy description. The value consists of 0 to 255 characters and
  must not contain a greater-than sign (>) or less-than sign (<).

* `provider_id` - (Optional, String, ForceNew) Specifies backup provider ID. Default value is
  **fc4d5750-22e7-4798-8a46-f48f62c4c1da**

* `common` - (Optional, Map) General backup policy parameters, which are blank by default.

* `scheduled_operation` - (Required, Set)  Backup plan information.

    + `name` - (Optional, String) Specifies Scheduling period name.The value consists of 1 to 255 characters and
      can contain only letters, digits, underscores (_), and hyphens (-).

    + `description` - (Optional, String) Specifies Scheduling period description.The value consists of 0 to 255
      characters and must not contain a greater-than sign (>) or less-than sign (<).

    + `enabled` - (Optional, Bool) Specifies whether the scheduling period is enabled. Default value is **true**.

    + `max_backups` - (Optional, Int) Specifies maximum number of backups that can be automatically created for a
      backup object.

    + `retention_duration_days` - (Optional, Int) Specifies duration of retaining a backup, in days.

    + `permanent` - (Optional, Bool) Specifies whether backups are permanently retained.

    + `trigger_pattern` - (Required, String) Specifies Scheduling policy of the scheduler.

    + `operation_type` - (Required, String) Specifies Operation type, which can be backup.

* `resource` - (Required, List) Backup Object.

    + `id` - (Required, String) Specifies the ID of the object to be backed up.

    + `type` - (Required, String) Entity object type of the backup object.
      If the type is VMs, the value is **OS::Nova::Server**.

    + `name` - (Required, String) Specifies backup object name.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `status` - Status of Backup Policy.

* `id` - Backup Policy ID.

* `scheduled_operation` -  Backup plan information.
  The [scheduled_operation](#csbs_scheduled_operation) object structure is documented below.

<a name="csbs_scheduled_operation"></a>
The `scheduled_operation` block supports:

* `id` -  Specifies Scheduling period ID.

* `trigger_id` -  Specifies Scheduler ID.

* `trigger_name` -  Specifies Scheduler name.

* `trigger_type` -  Specifies Scheduler type.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.
* `delete` - Default is 10 minutes.

## Import

Backup Policy can be imported using  `id`, e.g.

```shell
terraform import flexibleengine_csbs_backup_policy_v1.backup_policy_v1 7056d636-ac60-4663-8a6c-82d3c32c1c64
```
