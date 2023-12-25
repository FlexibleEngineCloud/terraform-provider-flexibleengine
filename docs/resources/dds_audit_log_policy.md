---
subcategory: "Document Database Service (DDS)"
---

# flexibleengine_dds_audit_log_policy

Manages a DDS audit log policy resource within FlexibleEngine.

## Example Usage

```hcl
variable "instance_id" {}
variable "keep_days" {}

resource "flexibleengine_dds_audit_log_policy" "test"{
  instance_id = var.instance_id
  keep_days   = var.keep_days
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `instance_id` - (Required, String, ForceNew) Specifies the ID of the DDS instance.

  Changing this parameter will create a new resource.

* `keep_days` - (Required, Int) Specifies the number of days for storing audit logs. The value ranges from 7 to 732.

* `audit_scope` - (Optional, String) Specifies the audit scope.
  If this parameter is left blank or set to **all**, all audit log policies are enabled.
  You can enter the database or collection name. Use commas (,) to separate multiple databases
  or collections. If the name contains a comma (,), add a dollar sign ($) before the comma
  to distinguish it from the separators. Enter a maximum of 1024 characters. The value
  cannot contain spaces or the following special characters "[]{}():? The dollar sign ($)
  can be used only in escape mode.

* `audit_types` - (Optional, List) Specifies the audit type. The value is **auth**, **insert**, **delete**, **update**,
  **query** or **command**.

* `reserve_auditlogs` - (Optional, String) Specifies whether the historical audit logs are
  retained when SQL audit is disabled.
  + **true**: indicates that historical audit logs are retained when SQL audit is disabled.(default value)
  + **false**: indicates that existing historical audit logs are deleted when SQL audit is disabled.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

## Import

The DDS audit log policy can be imported using the instance ID, e.g.:

```shell
terraform import flexibleengine_dds_audit_log_policy.test <instance_id>
```
