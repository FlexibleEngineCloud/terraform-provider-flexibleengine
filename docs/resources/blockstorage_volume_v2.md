---
subcategory: "Elastic Volume Service (EVS)"
description: ""
page_title: "flexibleengine_blockstorage_volume_v2"
---

# flexibleengine_blockstorage_volume_v2

Manages a V2 volume resource within FlexibleEngine.

## Example Usage

```hcl
resource "flexibleengine_blockstorage_volume_v2" "volume_1" {
  name        = "volume_1"
  description = "first test volume"
  size        = 3
  metadata = {
    __system__encrypted = "1"
    __system__cmkid     = "kms_id"
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the volume.
  If omitted, the `region` argument of the provider is used.
  Changing this creates a new volume.

* `size` - (Required, Int) The size of the volume to create (in gigabytes).

* `availability_zone` - (Optional, String, ForceNew) The availability zone for the volume.
  Changing this creates a new volume.

* `consistency_group_id` - (Optional, String, ForceNew) The consistency group to place the volume in.
  Changing this creates a new volume.

* `description` - (Optional, String) A description of the volume.
  Changing this updates the volume's description.

* `image_id` - (Optional, String, ForceNew) The image ID from which to create the volume.
  Changing this creates a new volume.

* `metadata` - (Optional, Map) Metadata key/value pairs to associate with the volume.
  Changing this updates the existing volume metadata.
  
  The EVS encryption capability with KMS key can be set with the following parameters:
    + `__system__encrypted` - The default value is set to '0', which means
      the volume is not encrypted, the value '1' indicates volume is encrypted.
    + `__system__cmkid` - (Optional) The ID of the kms key.

* `name` - (Optional, String) A unique name for the volume. Changing this updates the volume's name.

* `snapshot_id` - (Optional, String, ForceNew) The snapshot ID from which to create the volume.
  Changing this creates a new volume.

* `source_replica` - (Optional, String, ForceNew) The volume ID to replicate with.
  Changing this creates a new volume.

* `source_vol_id` - (Optional, String, ForceNew) The volume ID from which to create the volume.
  Changing this creates a new volume.

* `volume_type` - (Optional, String, ForceNew) The type of volume to create.
  Changing this creates a new volume.

* `cascade` - (Optional, Bool) Specifies to delete all snapshots associated with the EVS disk, Defaults to false.

* `multiattach` - (Optional, Bool) Specifies whether the EVS disk is shareable.

* `tags` - (Optional, Map) The key/value pairs to associate with the volume.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Specifies a resource ID in UUID format.

* `attachment` - If a volume is attached to an instance, this attribute will
  display the Attachment ID, Instance ID, and the Device as the Instance sees it.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.
* `delete` - Default is 10 minutes.

## Import

Volumes can be imported using the `id`, e.g.

```shell
terraform import flexibleengine_blockstorage_volume_v2.volume_1 ea257959-eeb1-4c10-8d33-26f0409a755d
```
