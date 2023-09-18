---
subcategory: "API Gateway (Dedicated APIG)"
description: ""
page_title: "flexibleengine_apig_custom_authorizer"
---

# flexibleengine_apig_custom_authorizer

Manages an APIG custom authorizer resource within Flexibleengine.

## Example Usage

```hcl
variable "instance_id" {}
variable "authorizer_name" {}
variable "function_urn" {}

resource "flexibleengine_apig_custom_authorizer" "test" {
  instance_id  = var.instance_id
  name         = var.authorizer_name
  function_urn = var.function_urn
  type         = "FRONTEND"
  cache_age    = 60

  identity {
    name     = "user_name"
    location = "QUERY"
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the custom authorizer resource.
  If omitted, the provider-level region will be used.
  Changing this will create a new custom authorizer resource.

* `instance_id` - (Required, String, ForceNew) Specifies an ID of the APIG dedicated instance to which the
  custom authorizer belongs to.
  Changing this will create a new custom authorizer resource.

* `name` - (Required, String) Specifies the name of the custom authorizer.
  The custom authorizer name consists of 3 to 64 characters, starting with a letter.
  Only letters, digits and underscores (_) are allowed.

* `type` - (Optional, String, ForceNew) Specifies the custom authoriz type.
  The valid values are *FRONTEND* and *BACKEND*. Changing this will create a new custom authorizer resource.

* `function_urn` - (Required, String) Specifies the uniform function URN of the function graph resource.

* `is_body_send` - (Optional, Bool) Specifies whether to send the body.

* `cache_age` - (Optional, Int) Specifies the maximum cache age.

* `user_data` - (Optional, String) Specifies the user data, which can contain a maximum of 2,048 characters.
  The user data is used by APIG to invoke the specified authentication function when accessing the backend service.

  -> **NOTE:** The user data will be displayed in plain text on the console.

* `identity` - (Optional, List) Specifies an array of one or more parameter identities of the custom authorizer.
  The object structure is documented below.

The `identity` block supports:

* `name` - (Required, String) Specifies the name of the parameter to be verified.
  The parameter includes front-end and back-end parameters.

* `location` - (Required, String) Specifies the parameter location, which support 'HEADER' and 'QUERY'.

* `validation` - (Optional, String) Specifies the parameter verification expression.
  If omitted, the custom authorizer will not perform verification.
  The valid value is range form 1 to 2,048.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the custom authorizer.
* `create_at` - Time when the APIG custom authorizer was created.

## Import

Custom Authorizers of the APIG can be imported using the ID of the APIG instance to which the group belongs and
Custom Authorizer `name`, separated by a slash, e.g.

```shell
terraform import flexibleengine_apig_custom_authorizer.test <instance_id>/<name>
```
