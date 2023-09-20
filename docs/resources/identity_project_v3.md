---
subcategory: "Identity and Access Management (IAM)"
description: ""
page_title: "flexibleengine_identity_project_v3"
---

# flexibleengine_identity_project_v3

Manages a Project resource within FlexibleEngine IAM service.

-> You *must* have admin privileges in your FlexibleEngine cloud to use this resource.

!> Project deletion is not supported by FlexibleEngine API

## Example Usage

```hcl
resource "flexibleengine_identity_project_v3" "project_1" {
  name        = "eu-west-0_project_1"
  description = "A ACC test project"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, String) The name of the project. The length is less than or equal
     to 64 bytes. Name mut be prefixed with a valid region name (eg. eu-west-0_project_1).

* `description` - (Optional, String) A description of the project.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - A resource ID in UUID format.

* `parent_id` - The parent of this project.

* `enabled` - Enabling status of this project.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.
* `update` - Default is 10 minutes.
* `delete` - Default is 10 minutes.

## Import

Projects can be imported using the `id`, e.g.

```shell
terraform import flexibleengine_identity_project_v3.project_1 <ID>
```
