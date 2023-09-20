---
subcategory: "Identity and Access Management (IAM)"
---

# flexibleengine_identity_project_v3

Use this data source to get the ID of a FlexibleEngine project.

## Example Usage

```hcl
data "flexibleengine_identity_project_v3" "project_1" {
  name = "eu-west-0"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Optional, String) The name of the project.

* `domain_id` - (Optional, String) The domain this project belongs to.

* `parent_id` - (Optional, String) The parent of this project.

* `enabled` - (Optional, Bool) The enabling status of this project.

* `is_domain` - (Optional, Bool) Indicates whether the user calling the API is a tenant.

## Attribute Reference

`id` is set to the ID of the found project. In addition, the following attributes
are exported:

* `description` - The description of the project.
