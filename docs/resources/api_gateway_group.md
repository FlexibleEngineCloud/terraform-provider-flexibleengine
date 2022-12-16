---
subcategory: "API Gateway"
description: ""
page_title: "flexibleengine_api_gateway_group"
---

# flexibleengine_api_gateway_group

Provides an API gateway group resource within FlexibleEngine.

## Example Usage

```hcl
resource "flexibleengine_api_gateway_group" "apigw_group" {
  name        = "apigw_group"
  description = "your descpiption"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the API gateway group resource. If omitted, the
  provider-level region will be used. Changing this creates a new gateway group resource.

* `name` - (Required, String) Specifies the name of the API group. An API group name consists of 3–64 characters,
  starting with a letter. Only letters, digits, and underscores (_) are allowed.

* `description` - (Optional, String) Specifies the description of the API group. The description cannot exceed 255
  characters.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the API group.
* `status` - Status of the API group.

## Import

API groups can be imported using the `id`, e.g.

```
$ terraform import flexibleengine_api_gateway_group.apigw_group "c8738f7c-a4b0-4c5f-a202-bda7dc4018a4"
```
