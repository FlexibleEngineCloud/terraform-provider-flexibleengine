---
subcategory: "Data Lake Insight (DLI)"
---

# flexibleengine_dli_queue

DLI Queue management
This is an alternative to `flexibleengine_dli_queue`

## Example Usage

### create a queue

```hcl
resource "flexibleengine_dli_queue" "queue" {
  name     = "terraform_dli_queue_test"
  cu_count = 16
}
```

## Argument Reference

The following arguments are supported:

* `cu_count` - (Required, Int, ForceNew) Minimum number of CUs that are bound to a queue. The value can be 16,
  64, or 256. Changing this parameter will create a new resource.

* `name` - (Required, String, ForceNew) Name of a queue. Name of a newly created resource queue. 
    The name can contain only digits, letters, and underscores (_), 
    but cannot contain only digits or start with an underscore (_).
    Length range: 1 to 128 characters. Changing this parameter will create a new resource.

* `description` - (Optional, String, ForceNew) Description of a queue. 
    Changing this parameter will create a new resource.

* `queue_type` - (Optional, String, ForceNew) Indicates the queue type. 
    Changing this parameter will create a new resource. The options are as follows:
    - sql,
    - general
    - all
    > NOTE: If the type is not specified, the default value sql is used. 

* `subnet_cidr` - (Optional, String, ForceNew) Subnet CIDR. Changing this parameter will create a new resource.

* `enterprise_project_id` - (Optional, String, ForceNew) Enterprise project ID. 
    The value 0 indicates the default enterprise project. Changing this parameter will create a new resource.

* `platform` - (Optional, String, ForceNew) CPU architecture of queue compute resources. The value can be x86_64. 
    Changing this parameter will create a new resource.

* `resource_mode` - (Optional, String, ForceNew) Queue resource mode. 
  Changing this parameter will create a new resource. 
  The options are as follows: 
  - 0: indicates the shared resource mode.
  - 1: indicates the exclusive resource mode. 

* `tags` - (Optional, String, ForceNew) Label of a queue. Changing this parameter will create a new resource.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `create_time` -  Time when a queue is created.
