---
subcategory: "Enterprise Project Management Service (EPS)"
description: ""
page_title: "flexibleengine_enterprise_project"
---

# flexibleengine_enterprise_project

Use this resource to manage an enterprise project within FlexibleEngine.

-> Deleting enterprise projects is not support. If you destroy a resource of enterprise project,
  the project is only disabled and removed from the state, but it remains in the cloud.
  Please set `insecure = true` in provider block to ignore SSL certificate verification when you got an x509 error.

## Example Usage

```hcl
resource "flexibleengine_enterprise_project" "test" {
  name        = "test"
  description = "example project"
}
```

## Argument Reference

* `name` - (Optional, String) Specifies the name of the enterprise project.
  This parameter can contain 1 to 64 characters. Only letters, digits, underscores (_), and hyphens (-) are allowed.
  The name must be unique in the domain and cannot include any form of the word "default" ("deFaulT", for instance).

* `description` - (Optional, String) Specifies the description of the enterprise project.

* `enable` - (Optional, Bool) Specifies whether to enable the enterprise project. Default to *true*.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `status` - Indicates the status of an enterprise project.
  + 1 indicates Enabled.
  + 2 indicates Disabled.

* `type` - Indicates the type of the enterprise project.

* `created_at` - Indicates the UTC time when the enterprise project was created. Example: 2018-05-18T06:49:06Z

* `updated_at` - Indicates the UTC time when the enterprise project was modified. Example: 2018-05-28T02:21:36Z

## Import

Enterprise projects can be imported using the `id`, e.g.

```shell
terraform import flexibleengine_enterprise_project.test 88f889c7-270e-4e77-8230-bf7db08d9b0e
```
