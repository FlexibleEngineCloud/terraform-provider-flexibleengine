---
subcategory: "Image Management Service (IMS)"
---

# flexibleengine_images_image_v2

Use this data source to get the ID of an available FlexibleEngine image.

## Example Usage

```hcl
data "flexibleengine_images_image_v2" "ubuntu" {
  name = "OBS Ubuntu 18.04"
}
```

## Argument Reference

* `region` - (Optional) The region in which to obtain the image.
    If omitted, the `region` argument of the provider is used.

* `name` - (Optional) The name of the image. Exact matching is used.

* `owner` - (Optional) The owner (UUID) of the image.

* `size_min` - (Optional) The minimum size (in bytes) of the image to return.

* `size_max` - (Optional) The maximum size (in bytes) of the image to return.

* `sort_direction` - (Optional) Order the results in either `asc` or `desc`.

* `sort_key` - (Optional) Sort images based on a certain key. Defaults to `name`.

* `tag` - (Optional) Search for images with a specific tag.

* `visibility` - (Optional) The visibility of the image. Must be one of
   "public", "private", "community", or "shared".

* `most_recent` - (Optional) If more than one result is returned, use the most
  recent image.

## Attributes Reference

`id` is set to the ID of the found image. In addition, the following attributes
are exported:

* `checksum` - The checksum of the data associated with the image.
* `container_format`: The format of the image's container.
* `disk_format`: The format of the image's disk.
* `file` - The URL for uploading and downloading the image file.
* `metadata` - The metadata associated with the image.
   Image metadata allow for meaningfully define the image properties and tags.
* `min_disk_gb`: The minimum amount of disk space required to use the image.
* `min_ram_mb`: The minimum amount of ram required to use the image.
* `protected` - Whether or not the image is protected.
* `schema` - The path to the JSON-schema that represent
   the image or image
* `size_bytes` - The size of the image (in bytes).
* `created_at` - The date the image was created.
* `updated_at` - The date the image was last updated.
