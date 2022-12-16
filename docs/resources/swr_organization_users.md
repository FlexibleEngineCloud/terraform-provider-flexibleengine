---
subcategory: "Software Repository for Container (SWR)"
description: ""
page_title: "flexibleengine_swr_organization_users"
---

# flexibleengine_swr_organization_users

Manages user permissions for the SWR organization resource within FlexibleEngine.

## Example Usage

```hcl
variable "organization_name" {}
variable "user_1" {}
variable "user_2" {}

resource "flexibleengine_swr_organization_users" "test" {
  organization = var.organization_name

  users {
    user_name  = var.user_1.name
    user_id    = var.user_1.id
    permission = "Read"
  }

  users {
    user_name  = var.user_2.name
    user_id    = var.user_2.id
    permission = "Read"
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource. If omitted, the
  provider-level region will be used. Changing this creates a new resource.

* `organization` - (Required, String, ForceNew) Specifies the name of the organization (namespace) to be accessed.
  Changing this creates a new resource.

* `users` - (Required, List) Specifies the users to access to the organization (namespace).
  Structure is documented below.

The `users` block supports:

* `permission` - (Required, String) Specifies the permission of the existing IAM user.
  The values can be **Manage**, **Write** and **Read**.

* `user_id` - (Required, String) Specifies the ID of the existing IAM user.

* `user_name` - (Optional, String) Specifies the name of the existing IAM user.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource. The value is the name of the organization.

* `creator` - The creator user name of the organization.

* `self_permission` - The permission informations of current user.

The `self_permission` block supports:

* `user_name` - The name of current user.

* `user_id` - The ID of current user.

* `permission` - The permission of current user.

## Import

Organization Permissions can be imported using the `id` (organization name), e.g.

```shell
terraform import flexibleengine_swr_organization_users.test org-test
```
