---
subcategory: "Image Management Service (IMS)"
---

# flexibleengine_images_image_share_accepter

Use this resource to accept an IMS image share from other users within FlexibleEngine.

## Example Usage

```hcl
variable "image_id" {}

resource "flexibleengine_images_image_share_accepter" "test" {
  image_id = var.image_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `image_id` - (Required, String, ForceNew) Specifies the ID of the image.

  Changing this parameter will create a new resource.

* `vault_id` - (Optional, String, ForceNew) Specifies the ID of a vault. This parameter is mandatory if you want
  to accept a shared full-ECS image created from a CBR backup.

  Changing this parameter will create a new resource.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
