---
layout: "flexibleengine"
page_title: "FlexibleEngine: flexibleengine_blockstorage_volume_v2"
sidebar_current: "docs-flexibleengine-resource-blockstorage-volume-v2"
description: |-
  Get information on an FlexibleEngine Volume.
---

# flexibleengine\_blockstorage\_volume_v2

Use this data source to get the ID of an available FlexibleEngine volume.

## Example Usage

```hcl
data "flexibleengine_blockstorage_volume_v2" "volume" {
  name = "test_volume"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional) The region in which to obtain the V2 Volume client. If omitted, the `region` argument of the provider is used.

* `name` - (Optional) The name of the volume.

* `status` - (Optional) The status of the volume.


## Attributes Reference

The following attributes are exported:

* `id` - The ID of the volume.
