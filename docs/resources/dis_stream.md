---
subcategory: "Data Ingestion Service (DIS)"
---

# flexibleengine_dis_stream

Manages DIS Stream resource within FlexibleEngine.

## Example Usage

```hcl
resource "flexibleengine_dis_stream" "stream" {
  name            = "dis-demo"
  partition_count = 3
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, String, ForceNew) Specifies the name of the DIS stream to be created.
  Changing this will create a new resource.

* `partition_count` - (Required, Int, ForceNew) Specifies the number of the expect partitions.
  Changing this will create a new resource.

* `type` - (Optional, String, ForceNew) Specifies the Stream type. The value can be *COMMON* or *ADVANCED*.
  Defaults to *COMMON*. Changing this will create a new resource.

  + **COMMON stream:**
    Each partition supports a read speed of up to 2 MB/s and a write speed of up to 1000 records/s and 1 MB/s.

  + **ADVANCED stream:**
    Each partition supports a read speed of up to 10 MB/s and a write speed of up to 2000 records/s and 5 MB/s.

* `retention_period` - (Optional, Int, ForceNew) Specifies the number of hours for which data from the stream
  will be retained in DIS. The value ranges from 24 to 168 and defaults to 24. Changing this will create a new resource.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID which equals to stream name.

* `status` - Status of stream: `CREATING`,`RUNNING`,`TERMINATING`,`TERMINATED`,`FROZEN`.

* `partitions` - The information of stream partitions. Structure is documented below.

The `partitions` block contains:

* `id` - The ID of the partition.

* `status` - The status of the partition.

* `hash_range` - Possible value range of the hash key used by each partition.

* `sequence_number_range` - Sequence number range of each partition.

## Import

Dis stream can be imported by `name`. For example,

```
terraform import flexibleengine_dis_stream.example dis-demo
```
