---
subcategory: "Auto Scaling (AS)"
---

# flexibleengine_as_lifecycle_hook_v1

Manages an AS Lifecycle Hook resource within FlexibleEngine.

## Example Usage

### Basic Lifecycle Hook

```hcl
variable "hook_name" {}

variable "as_group_id" {}

variable "smn_topic_urn" {}

resource "flexibleengine_as_lifecycle_hook_v1" "test" {
  name                   = var.hook_name
  scaling_group_id       = var.as_group_id
  type                   = "ADD"  
  default_result         = "ABANDON"
  notification_topic_urn = var.smn_topic_urn
  notification_message   = "This is a test message"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, String) Specifies the lifecycle hook name.
  This parameter can contain a maximum of 32 characters, which may consist of letters, digits,
  underscores (_) and hyphens (-).

* `scaling_group_id` - (Required, String, ForceNew) Specifies the ID of the AS group in UUID format.
  Changing this creates a new AS lifecycle hook.

* `type` - (Required, String) Specifies the lifecycle hook type.
  The valid values are following strings:
  * `ADD`: The hook suspends the instance when the instance is started.
  * `REMOVE`: The hook suspends the instance when the instance is terminated.

* `notification_topic_urn` - (Required, String) Specifies a unique topic in SMN.

* `default_result` - (Optional, String) Specifies the default lifecycle hook callback operation.
  This operation is performed when the timeout duration expires.
  The valid values are *ABANDON* and *CONTINUE*, default to *ABANDON*.

* `timeout` - (Optional, Int) Specifies the lifecycle hook timeout duration, which ranges from 300 to 86400 in the
  unit of second, default to 3600.

* `notification_message` - (Optional, String) Specifies a customized notification.
  This parameter can contains a maximum of 256 characters, which cannot contain the following characters: <>&'().

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `notification_topic_name` - The topic name in SMN.

* `create_time` - The server time in UTC format when the lifecycle hook is created.

## Import

Lifecycle hooks can be imported using the AS group ID and hook ID separated by a slash, e.g.

```
$ terraform import flexibleengine_as_lifecycle_hook_v1.test <AS group ID>/<Lifecycle hook ID>
```
