---
layout: "flexibleengine"
page_title: "FlexibleEngine: flexibleengine_drs_replicationconsistencygroup_v2"
sidebar_current: "docs-flexibleengine-resource-drs-replicationconsistencygroup-v2"
description: |-
  Manages a V2 replicationconsistencygroup resource within FlexibleEngine.
---

# flexibleengine\_drs\_replicationconsistencygroup\_v2

Manages a V2 replicationconsistencygroup resource within FlexibleEngine.

## Example Usage

```hcl
resource "flexibleengine_blockstorage_volume_v2" "volume_1" {
  name = "volume_1"
  size = 1
  availability_zone = "eu-west-0a"
}

resource "flexibleengine_blockstorage_volume_v2" "volume_2" {
  name = "volume_2"
  size = 1
  availability_zone = "eu-west-0b"
}

resource "flexibleengine_drs_replication_v2" "replication_1" {
  name = "replication_1"
  description = "The description of replication_1"
  volume_ids = ["${flexibleengine_blockstorage_volume_v2.volume_1.id}", "${flexibleengine_blockstorage_volume_v2.volume_2.id}"]
  priority_station = "eu-west-0a"
}

resource "flexibleengine_drs_replicationconsistencygroup_v2" "replicationconsistencygroup_1" {
  name = "replicationconsistencygroup_1"
  description = "The description of replicationconsistencygroup_1"
  replication_ids = ["${flexibleengine_drs_replication_v2.replication_1.id}"]
  priority_station = "eu-west-0a"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Optional) The name of the replication consistency group. The name can contain a maximum of 255 bytes.

* `description` - (Optional) The description of the replication consistency group. The description can contain a maximum of 255 bytes.

* `replication_ids` - (Required) An array of one or more IDs of the EVS replication pairs used to create the replication consistency group.

* `priority_station` - (Required) The primary AZ of the replication consistency group. That is the AZ where the production disk belongs.

* `replication_model` - (Optional) The type of the created replication consistency group. Currently only type hypermetro is supported.

## Attributes Reference

The following attributes are exported:

* `name` - See Argument Reference above.
* `description` - See Argument Reference above.
* `replication_ids` - See Argument Reference above.
* `priority_station` - See Argument Reference above.
* `replication_model` - See Argument Reference above.
* `status` - The status of the replication consistency group.
* `replication_status` - The replication status of the replication consistency group.
* `created_at` - The creation time of the replication consistency group.
* `updated_at` - The update time of the replication consistency group.
* `failure_detail` - The returned error code if the replication consistency group status is error.
* `fault_level` - The fault level of the replication consistency group.
