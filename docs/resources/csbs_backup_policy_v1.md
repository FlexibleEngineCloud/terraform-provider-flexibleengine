---
subcategory: "Cloud Server Backup Service (CSBS)"
---

# flexibleengine_csbs_backup_policy_v1

Provides an FlexibleEngine Backup Policy of Resources.

## Example Usage

 ```hcl
 variable "name" { }
 variable "id" { }
 variable "resource_name" { }
 
 resource "flexibleengine_csbs_backup_policy_v1" "backup_policy_v1" {
   name  = "${var.name}"
   resource {
     id = "${var.id}"
     type = "OS::Nova::Server"
     name = "${var.resource_name}"
   }
   scheduled_operation {
     enabled = true
     operation_type = "backup"
     trigger_pattern = "BEGIN:VCALENDAR\r\nBEGIN:VEVENT\r\nRRULE:FREQ=WEEKLY;BYDAY=TH;BYHOUR=12;BYMINUTE=27\r\nEND:VEVENT\r\nEND:VCALENDAR\r\n"
   }
 }

 ```
## Argument Reference
The following arguments are supported:

* `name` - (Required) Specifies the name of backup policy. The value consists of 1 to 255 characters and can contain only letters, digits, underscores (_), and hyphens (-).

* `description` - (Optional) Backup policy description. The value consists of 0 to 255 characters and must not contain a greater-than sign (>) or less-than sign (<).

* `provider_id` - (Required) Specifies backup provider ID. Default value is **fc4d5750-22e7-4798-8a46-f48f62c4c1da**

* `common` - (Optional) General backup policy parameters, which are blank by default.

* `scheduled_operation` block supports the following arguments:

    * `name` - (Optional) Specifies Scheduling period name.The value consists of 1 to 255 characters and can contain only letters, digits, underscores (_), and hyphens (-).
    
    * `description` - (Optional) Specifies Scheduling period description.The value consists of 0 to 255 characters and must not contain a greater-than sign (>) or less-than sign (<).

    * `enabled` - (Optional) Specifies whether the scheduling period is enabled. Default value is **true**

    * `max_backups` - (Optional) Specifies maximum number of backups that can be automatically created for a backup object.

    * `retention_duration_days` - (Optional) Specifies duration of retaining a backup, in days.

    * `permanent` - (Optional) Specifies whether backups are permanently retained.

    * `trigger_pattern` - (Required) Specifies Scheduling policy of the scheduler.

    * `operation_type` - (Required) Specifies Operation type, which can be backup.

* `resource` block supports the following arguments:

    * `id` - (Required) Specifies the ID of the object to be backed up.
    
    * `type` - (Required) Entity object type of the backup object. If the type is VMs, the value is **OS::Nova::Server**.

    * `name` - (Required) Specifies backup object name.

## Attributes Reference
In addition to all arguments above, the following attributes are exported:

* `status` - Status of Backup Policy.

* `id` - Backup Policy ID.

* scheduled_operation - Backup plan information

    * `id` -  Specifies Scheduling period ID.

    * `trigger_id` -  Specifies Scheduler ID.

    * `trigger_name` -  Specifies Scheduler name.

    * `trigger_type` -  Specifies Scheduler type.


## Import

Backup Policy can be imported using  `id`, e.g.

```
$ terraform import flexibleengine_csbs_backup_policy_v1.backup_policy_v1 7056d636-ac60-4663-8a6c-82d3c32c1c64
```




