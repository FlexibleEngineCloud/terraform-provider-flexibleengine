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

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the dli queue resource. If omitted,
  the provider-level region will be used. Changing this will create a new VPC channel resource.

* `cu_count` - (Required, Int) Minimum number of CUs that are bound to a queue. Initial value can be `16`,
  `64`, or `256`. When scale_out or scale_in, the number must be a multiple of 16.

* `name` - (Required, String, ForceNew) Name of a queue. Name of a newly created resource queue.
  The name can contain only digits, letters, and underscores (\_),
  but cannot contain only digits or start with an underscore (_).
  Length range: 1 to 128 characters. Changing this parameter will create a new resource.

* `description` - (Optional, String, ForceNew) Description of a queue.
  Changing this parameter will create a new resource.

* `queue_type` - (Optional, String, ForceNew) Indicates the queue type.
  Changing this parameter will create a new resource. The options are as follows:
    - **sql**. SQL queue type, used to run SQL jobs.
    - **general**. Universal queue type for running Flink, Spark Jar jobs.

    The default value is `sql`.

* `resource_mode` - (Optional, String, ForceNew) Queue resource mode.
  Changing this parameter will create a new resource.
  The options are as follows:
    - **0**: indicates the shared resource mode.
    - **1**: indicates the exclusive resource mode.

* `feature` - (Optional, String, ForceNew)Indicates the queue feature. Changing this parameter will create a new
  resource. The options are as follows:
    - **basic**: basic type (default value)

* `platform` - (Optional, String, ForceNew) CPU architecture of queue compute resources. Changing this parameter will
  create a new resource. The options are as follows:
    - **x86_64** : default value.
    - **aarch64**.

* `vpc_cidr` - (Optional, String) The CIDR block of a queue. If use DLI enhanced datasource connections, the CIDR block
  cannot be the same as that of the data source.
  The CIDR blocks supported by different CU specifications:

    - When `cu_count` is `16` or `64`: 10.0.0.0~10.255.0.0/8~24, 172.16.0.0~172.31.0.0/12~24,
      192.168.0.0~192.168.0.0/16~24.
    - When `cu_count` is `256`: 10.0.0.0~10.255.0.0/8~22, 172.16.0.0~172.31.0.0/12~22, 192.168.0.0~192.168.0.0/16~22.

* `tags` - (Optional, Map, ForceNew) Label of a queue. Changing this parameter will create a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Specifies a resource ID in UUID format.

* `create_time` -  Time when a queue is created.

## Timeouts

This resource provides the following timeouts configuration options:

* `update` - Default is 45 minutes.

## Import

DLI queue can be imported by  `id`. For example,

```shell
terraform import flexibleengine_dli_queue.example abc123
```
