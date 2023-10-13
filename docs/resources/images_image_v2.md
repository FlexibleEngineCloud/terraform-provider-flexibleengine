---
subcategory: "Deprecated"
description: ""
page_title: "flexibleengine_images_image_v2"
---

# flexibleengine_images_image_v2

Manages a V2 Image resource within FlexibleEngine Glance.

!> **WARNING:** It has been deprecated, please use `flexibleengine_images_image` instead.

## Example Usage

```hcl
resource "flexibleengine_images_image_v2" "rancheros" {
  name   = "RancherOS"
  image_source_url = "https://releases.rancher.com/os/latest/rancheros-openstack.img"
  container_format = "bare"
  disk_format = "qcow2"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this creates a new resource.

* `container_format` - (Required, String, ForceNew) The container format. Must be one of
  "ami", "ari", "aki", "bare", "ovf". Changing this creates a new Image.

* `disk_format` - (Required, String, ForceNew) The disk format. Must be one of
  "ami", "ari", "aki", "vhd", "vmdk", "raw", "qcow2", "vdi", "iso". Changing this creates a new Image.

* `local_file_path` - (Optional, String, ForceNew) This is the filepath of the raw image file
  that will be uploaded to Glance. Conflicts with `image_source_url`. Changing this creates a new Image.

* `image_cache_path` - (Optional, String) This is the directory where the images will
  be downloaded. Images will be stored with a filename corresponding to
  the url's md5 hash. Defaults to "$HOME/.terraform/image_cache"

* `image_source_url` - (Optional, String, ForceNew) This is the url of the raw image that will
  be downloaded in the `image_cache_path` before being uploaded to Glance.
  Glance is able to download image from internet but the `gophercloud` library does not yet provide a way to do so.
  Conflicts with `local_file_path`. Changing this creates a new Image.

* `min_disk_gb` - (Optional, Int, ForceNew) Amount of disk space (in GB) required to boot image.
   Defaults to 0. Changing this creates a new Image.

* `min_ram_mb` - (Optional, Int, ForceNew) Amount of ram (in MB) required to boot image.
  Defauts to 0. Changing this creates a new Image.

* `name` - (Required, String) The name of the image.

* `protected` - (Optional, Bool, ForceNew) If true, image will not be deletable.
  Defaults to **false**. Changing this creates a new Image.

* `tags` - (Optional, List) The tags of the image. It must be a list of strings.
  At this time, it is not possible to delete all tags of an image.

* `visibility` - (Optional, String) The visibility of the image. Must be one of
  "public", "private", "community", or "shared". The ability to set the
  visibility depends upon the configuration of the FlexibleEngine cloud.

Note: The `properties` attribute handling in the gophercloud library is currently buggy
and needs to be fixed before being implemented in this resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - A unique ID assigned by Glance.

* `checksum` - The checksum of the data associated with the image.

* `file` - the trailing path after the glance endpoint that represent the location of the image or the path
  to retrieve it.

* `metadata` - The metadata associated with the image.
  Image metadata allow for meaningfully define the image properties and tags.

* `owner` - The id of the flexibleengine user who owns the image.

* `schema` - The path to the JSON-schema that represent the image or image

* `size_bytes` - The size in bytes of the data associated with the image.

* `status` - The status of the image. It can be "queued", "active" or "saving".

* `created_at` - The date the image was created.

* `updated_at` - The date the image was last updated.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 30 minutes.

## Import

Images can be imported using the `id`, e.g.

```shell
terraform import flexibleengine_images_image_v2.rancheros 89c60255-9bd6-460c-822a-e2b959ede9d2
```
