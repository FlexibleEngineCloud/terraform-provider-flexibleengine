---
subcategory: "Identity and Access Management (IAM)"
---

# flexibleengine\_identity\_group_membership_v3

Manages a User Group Membership resource within FlexibleEngine IAM service.

Note: You _must_ have admin privileges in your FlexibleEngine cloud to use
this resource.

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

