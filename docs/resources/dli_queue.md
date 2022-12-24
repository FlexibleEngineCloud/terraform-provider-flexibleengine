---
subcategory: "Data Lake Insight (DLI)"
description: ""
page_title: "flexibleengine_dli_queue"
---

# flexibleengine_dli_queue

DLI Queue management
Allows you to create a queue. The queue will be bound to specified compute resources.

## Example Usage

### create a queue

```hcl
resource "flexibleengine_dli_queue" "queue" {
  name     = "terraform_dli_queue_test"
  cu_count = 16
  tags     = {
    foo = "bar"
    key = "value"
  }
}
```

## Argument Reference

The following arguments are supported:

* `cu_count` - (Required, Int) Minimum number of CUs that are bound to a queue. Initial value can be `16`,
  `64`, or `256`. When scale_out or scale_in, the number must be a multiple of 16

* `name` - (Required, String, ForceNew) Name of a queue. Name of a newly created resource queue.
    The name can contain only digits, letters, and underscores (\_),
    but cannot contain only digits or start with an underscore (_).
    Length range: 1 to 128 characters. Changing this parameter will create a new resource.

* `description` - (Optional, String, ForceNew) Description of a queue.
    Changing this parameter will create a new resource.

* `queue_type` - (Optional, String, ForceNew) Indicates the queue type.
    Changing this parameter will create a new resource. The options are as follows:
    - sql,
    - general

    The default value is `sql`.

* `resource_mode` - (Optional, String, ForceNew) Queue resource mode.
  Changing this parameter will create a new resource.
  The options are as follows:
  - 0: indicates the shared resource mode.
  - 1: indicates the exclusive resource mode.

* `tags` - (Optional, Map, ForceNew) Label of a queue. Changing this parameter will create a new resource.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Specifies a resource ID in UUID format.

* `create_time` -  Time when a queue is created.

## Timeouts

This resource provides the following timeouts configuration options:

* `update` - Default is 45 minute.

## Import

DLI queue can be imported by  `id`. For example,

```shell
terraform import flexibleengine_dli_queue.example  abc123
```
