---
layout: "flexibleengine"
page_title: "FlexibleEngine: flexibleengine_identity_project_v3"
sidebar_current: "docs-flexibleengine-datasource-identity-project-v3"
description: |-
  Get information on a FlexibleEngine Project.
---

# flexibleengine\_identity\_project_v3

Use this data source to get the ID of a FlexibleEngine project.

## Example Usage

```hcl
data "flexibleengine_identity_project_v3" "project_1" {
  name = "eu-west-0"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Optional) The name of the project.

* `domain_id` - (Optional) The domain this project belongs to.

* `parent_id` - (Optional) The parent of this project.

## Attributes Reference

`id` is set to the ID of the found project. In addition, the following attributes
are exported:

* `name` - See Argument Reference above.
* `domain_id` - See Argument Reference above.
* `parent_id` - See Argument Reference above.
* `description` - The description of the project.
* `enabled` - Whether the project is available.
* `is_domain` - Whether this project is a domain.
