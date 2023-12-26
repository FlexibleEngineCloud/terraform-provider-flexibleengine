---
subcategory: "API Gateway (Dedicated APIG)"
description: ""
page_title: "flexibleengine_apig_group"
---

# flexibleengine_apig_group

Manages an APIG (API) group resource within FlexibleEngine.

## Example Usage

```hcl
variable "instance_id" {}
variable "group_name" {}
variable "description" {}
variable "environment_id" {}

resource "flexibleengine_apig_group" "test" {
  instance_id = var.instance_id
  name        = var.group_name
  description = var.description

  environment {
    variable {
      name  = "TERRAFORM"
      value = "/stage/terraform"
    }
    environment_id = var.environment_id
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the API group resource. If omitted,
  the provider-level region will be used. Changing this will create a new API group resource.

* `instance_id` - (Required, String, ForceNew) Specifies an ID of the APIG dedicated instance to which the API group
  belongs to. Changing this will create a new API group resource.

* `name` - (Required, String) Specifies the name of the API group. The API group name consists of 3 to 64 characters,
  starting with a letter. Only letters, digits and underscores (_) are allowed. Chinese characters must be in UTF-8 or
  Unicode format.

* `description` - (Optional, String) Specifies the description about the API group. The description contain a maximum of
  255 characters and the angle brackets (< and >) are not allowed. Chinese characters must be in UTF-8 or Unicode
  format.

* `environment` - (Optional, List) Specifies an array of one or more APIG environments of the associated APIG group. The
  object structure is documented below.

The `environment` block supports:

* `variable` - (Required, List) Specifies an array of one or more APIG environment variables. The object structure is
  documented below. The environment variables of different groups are isolated in the same environment.

* `environment_id` - (Required, String) Specifies the APIG environment ID of the associated APIG group.

The `variable` block supports:

* `name` - (Required, String) Specifies the variable name, which can contains of 3 to 32 characters, starting with a
  letter. Only letters, digits, hyphens (-), and underscores (_) are allowed. In the definition of an API, `name` (
  case-sensitive) indicates a variable, such as #Name#. It is replaced by the actual value when the API is published in
  an environment. The variable names are not allowed to be repeated for an API group.

* `value` - (Required, String) Specifies the environment ariable value, which can contains of 1 to 255 characters. Only
  letters, digits and special characters (_-/.:) are allowed.

  -> **NOTE:** The variable value will be displayed in plain text on the console.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the API group.
* `registration_at` - Registration time, in RFC-3339 format.
* `update_at` - Time when the API group was last modified, in RFC-3339 format.
* `environment` - The array of one or more environments of the associated group.  
  The [object](#group_environment_attr) structure is documented below.

<a name="group_environment_attr"></a>
The `environment` block supports:

* `variable` - The array of one or more environment variables.  
  The [object](#group_environment_variable_attr) structure is documented below.

<a name="group_environment_variable_attr"></a>
The `variable` block supports:

* `id` - The variable ID.

## Import

API groups of the APIG can be imported using the ID of the APIG instance to which the group belongs and API group ID,
separated by a slash, e.g.

```shell
terraform import flexibleengine_apig_group.test <instance_id>/<id>
```
