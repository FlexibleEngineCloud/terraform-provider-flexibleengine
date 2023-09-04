---
subcategory: "Elastic Volume Service (EVS)"
---

# flexibleengine_blockstorage_volume_v2

Use this data source to get the ID of an available FlexibleEngine volume.

## Example Usage

```hcl
data "flexibleengine_blockstorage_volume_v2" "volume" {
  name = "test_volume"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to obtain the V2 Volume client.
  If omitted, the `region` argument of the provider is used.

* `name` - (Optional) The name of the volume.

* `status` - (Optional) The status of the volume.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the volume.
