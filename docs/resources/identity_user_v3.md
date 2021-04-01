---
subcategory: "Identity and Access Management (IAM)"
---

# flexibleengine\_identity\_user_v3

Manages a User resource within FlexibleEngine IAM service.

Note: You _must_ have admin privileges in your FlexibleEngine cloud to use
this resource.

## Example Usage

```hcl
resource "flexibleengine_identity_user_v3" "user_1" {
  name        = "user_1"
  description = "A user"
  password    = "password123!"

}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the user. The user name consists of 5 to 32
     characters. It can contain only uppercase letters, lowercase letters, 
     digits, spaces, and special characters (-_) and cannot start with a digit.

* `description` - (Optional) A description of the user.

* `default_project_id` - (Optional) The default project this user belongs to.

* `domain_id` - (Optional) The domain this user belongs to.

* `enabled` - (Optional) Whether the user is enabled or disabled. Valid
    values are `true` and `false`.

* `password` - (Optional) The password for the user. It must contain at least 
     two of the following character types: uppercase letters, lowercase letters, 
     digits, and special characters.

## Attributes Reference

The following attributes are exported:

* `domain_id` - See Argument Reference above.

## Import

Users can be imported using the `id`, e.g.

```
$ terraform import flexibleengine_identity_user_v3.user_1 89c60255-9bd6-460c-822a-e2b959ede9d2
```
