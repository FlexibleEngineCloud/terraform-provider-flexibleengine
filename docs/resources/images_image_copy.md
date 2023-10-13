---
subcategory: "Image Management Service (IMS)"
---

# flexibleengine_images_image_copy

Use this resource to copy IMS images from one region to another within FlexibleEngine.

## Example Usage

### Copy image within region

```hcl
variable "source_image_id" {}
variable "name" {}
variable "kms_key_id" {}

resource "flexibleengine_images_image_copy" "test" {
  source_image_id = var.source_image_id
  name            = var.name
  kms_key_id      = var.kms_key_id
}
```

### Copy image cross region

```hcl
variable "source_image_id" {}
variable "name" {}
variable "target_region" {}
variable "agency_name" {}

resource "flexibleengine_images_image_copy" "test" {
  source_image_id = var.source_image_id
  name            = var.name
  target_region   = var.target_region
  agency_name     = var.agency_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region to which the source image belongs.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `source_image_id` - (Required, String, ForceNew) Specifies the ID of the copied image.
  Changing this parameter will create a new resource.

* `name` - (Required, String) Specifies the name of the copy image. The name can contain `1` to `128` characters,
  only Chinese and English letters, digits, underscore (_), hyphens (-), dots (.) and space are
  allowed, but it cannot start or end with a space.

* `target_region` - (Optional, String, ForceNew) Specifies the target region name.
  If specified, it means cross-region replication. Changing this parameter will create a new resource.

* `description` - (Optional, String) Specifies the description of the copy image.

* `kms_key_id` - (Optional, String, ForceNew) Specifies the master key used for encrypting an image.
  Only copying scene within a region is supported. Changing this parameter will create a new resource.

* `enterprise_project_id` - (Optional, String, ForceNew) Specifies the enterprise project id of the image.
  Only copying scene within a region is supported. Changing this parameter will create a new resource.

* `agency_name` - (Optional, String, ForceNew) Specifies the agency name. It is required in the cross-region scene.
  Changing this parameter will create a new resource.

* `vault_id` - (Optional, String, ForceNew) Specifies the ID of the vault. It is used in the cross-region scene,
  and it is mandatory if you are replicating a full-ECS image.
  Changing this parameter will create a new resource.

* `tags` - (Optional, Map) Specifies the key/value pairs to associate with the copy image.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `instance_id` - Indicates the ID of the ECS that needs to be converted into an image.

* `os_version` - Indicates the OS version.

* `visibility` - Indicates whether the image is visible to other tenants.

* `data_origin` - Indicates the image resource.
  The pattern can be 'instance,**instance_id**' or 'file,**image_url**'.

* `disk_format` - Indicates the image file format.
  The value can be `vhd`, `zvhd`, `raw`, `zvhd2`, or `qcow2`.

* `image_size` - Indicates the size(bytes) of the image file format.

* `checksum` - Indicates the checksum of the data associated with the image.

* `status` - Indicates the status of the image.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 60 minutes.
* `delete` - Default is 3 minutes.
