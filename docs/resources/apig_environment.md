---
subcategory: "API Gateway (Dedicated APIG)"
description: ""
page_title: "flexibleengine_apig_environment"
---

# flexibleengine_apig_environment

Manages an APIG environment resource within Flexibleengine.

## Example Usage

```hcl
variable "instance_id" {}
variable "environment_name" {}
variable "description" {}

resource "flexibleengine_apig_environment" "test" {
  instance_id = var.instance_id
  name        = var.environment_name
  description = var.description
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the APIG environment resource. If
  omitted, the provider-level region will be used. Changing this will create a new APIG environment resource.

* `instance_id` - (Required, String, ForceNew) Specifies an ID of the APIG dedicated instance to which the API
  environment belongs to. Changing this will create a new APIG environment resource.

* `name` - (Required, String) Specifies the name of the API environment. The API environment name consists of 3 to 64
  characters, starting with a letter. Only letters, digits and underscores (_) are allowed.

* `description` - (Optional, String) Specifies the description about the API environment. The description contain a
  maximum of 255 characters and the angle brackets (< and >) are not allowed. Chinese characters must be in UTF-8 or
  Unicode format.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the APIG environment.
* `create_at` - Time when the APIG environment was created, in RFC-3339 format.

## Import

Environments can be imported using the ID of the APIG instance to which the environment belongs and environment ID,
separated by a slash, e.g.

```shell
terraform import flexibleengine_apig_environment.test <instance_id>/<id>
```
