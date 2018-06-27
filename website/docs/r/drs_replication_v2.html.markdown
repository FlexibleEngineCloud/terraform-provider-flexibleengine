---
layout: "flexibleengine"
page_title: "FlexibleEngine: flexibleengine_drs_replication_v2"
sidebar_current: "docs-flexibleengine-resource-drs-replication-v2"
description: |-
  Manages a V2 replication resource within FlexibleEngine.
---

# flexibleengine\_drs\_replication\_v2

Manages a V2 replication resource within FlexibleEngine.

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
```

## Argument Reference

The following arguments are supported:

* `name` - (Optional) The name of the EVS replication pair. The name can contain a maximum of 255 bytes.

* `description` - (Optional) The description of the EVS replication pair. The description can contain a maximum of 255 bytes.

* `volume_ids` - (Required) An array of one or more IDs of the EVS disks used to create the EVS replication pair.

* `priority_station` - (Required) The primary AZ of the EVS replication pair. That is the AZ where the production disk belongs.

* `replication_model` - (Optional) The type of the EVS replication pair. Currently only type hypermetro is supported.

## Attributes Reference

The following attributes are exported:

* `name` - See Argument Reference above.
* `description` - See Argument Reference above.
* `volume_ids` - See Argument Reference above.
* `priority_station` - See Argument Reference above.
* `replication_model` - See Argument Reference above.
* `status` - The status of the EVS replication pair.
* `replication_consistency_group_id` - The ID of the replication consistency group where the EVS replication pair belongs.
* `created_at` - The creation time of the EVS replication pair.
* `updated_at` - The update time of the EVS replication pair.
* `replication_status` - The replication status of the EVS replication pair.
* `progress` - The synchronization progress of the EVS replication pair. Unit: %.
* `failure_detail` - The returned error code if the EVS replication pair status is error.
* `record_metadata` - The metadata of the EVS replication pair.
* `fault_level` - The fault level of the EVS replication pair.
