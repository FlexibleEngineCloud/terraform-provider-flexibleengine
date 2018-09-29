---
layout: "flexibleengine"
page_title: "FlexibleEngine: flexibleengine_vbs_backup_policy_v2"
sidebar_current: "docs-flexibleengine-resource-vbs-backup-policy-v2"
description: |-
  Provides an VBS Backup Policy resource.
---

# flexibleengine_vbs_backup_policy_v2

Provides an VBS Backup Policy resource.

# Example Usage

 ```hcl
resource "flexibleengine_vbs_backup_policy_v2" "vbs" {
  name = "policy_002"
  start_time  = "12:00"
  status  = "ON"
  retain_first_backup = "N"
  rentention_num = 2
  frequency = 1
}
 ```

# Argument Reference

The following arguments are supported:

* `name` (Required) - Specifies the policy name. The value is a string of 1 to 64 characters that can contain letters, digits, underscores (_), and hyphens (-). It cannot start with default.

* `start_time` (Required) - Specifies the start time of the backup job.The value is in the HH:mm format.                                                         

* `status` (Required) - Specifies the backup policy status. The value can ON or OFF.

* `retain_first_backup` (Required) - Specifies whether to retain the first backup in the current month. Possible values are Y or N. 

* `rentention_num` (Required) - Specifies number of retained backups. Minimum value is 2.

* `frequency` (Required) - Specifies the backup interval. The value is in the range of 1 to 14 days.


# Attributes Reference

All of the argument attributes are also exported as
result attributes:

* `id` - Specifies a backup policy ID.
 
* `policy_resource_count` - Specifies the number of volumes associated with the backup policy.

# Import

Backup Policy can be imported using the `id`, e.g.

```
$ terraform import flexibleengine_vbs_backup_policy_v2.vbs 4779ab1c-7c1a-44b1-a02e-93dfc361b32d
```