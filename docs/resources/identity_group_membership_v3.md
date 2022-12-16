---
subcategory: "Identity and Access Management (IAM)"
description: ""
page_title: "flexibleengine_identity_group_membership_v3"
---

# flexibleengine_identity_group_membership_v3

Manages a User Group Membership resource within FlexibleEngine IAM service.

-> You *must* have admin privileges in your FlexibleEngine cloud to use this resource.

## Example Usage

```hcl
resource "flexibleengine_identity_group_v3" "group_1" {
  name        = "group1"
  description = "This is a test group"
}

resource "flexibleengine_identity_user_v3" "user_1" {
  name     = "user1"
  enabled  = true
  password = "password12345!"
}

resource "flexibleengine_identity_user_v3" "user_2" {
  name     = "user2"
  enabled  = true
  password = "password12345!"
}

resource "flexibleengine_identity_group_membership_v3" "membership_1" {
  group = flexibleengine_identity_group_v3.group_1.id
  users = [
    flexibleengine_identity_user_v3.user_1.id,
    flexibleengine_identity_user_v3.user_2.id
  ]
}
```

## Argument Reference

The following arguments are supported:

* `group` - (Required) The group ID of this membership.

* `users` - (Required) A List of user IDs to associate to the group.

## Attributes Reference

The following attributes are exported:

* `group` - See Argument Reference above.

* `users` - See Argument Reference above.

## Import

A group membership can be imported using the `id`, e.g.

```shell
terraform import flexibleengine_identity_group_membership_v3.membership_1 5e6f7e8d-9c0b-4a1a-8b2c-3d4e5f6a7b8c
```
