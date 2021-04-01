---
subcategory: "Identity and Access Management (IAM)"
---

# flexibleengine\_identity\_group_v3

Manages a User Group resource within FlexibleEngine IAM service.

Note: You _must_ have admin privileges in your FlexibleEngine cloud to use
this resource.

## Example Usage

```hcl
resource "flexibleengine_identity_group_v3" "group_1" {
  name        = "group_1"
  description = "This is a test group"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the group.The length is less than or equal 
     to 64 bytes 

* `description` - (Optional) A description of the group.

* `domain_id` - (Optional) The domain this group belongs to.

## Attributes Reference

The following attributes are exported:

* `domain_id` - See Argument Reference above.

## Import

Groups can be imported using the `id`, e.g.

```
$ terraform import flexibleengine_identity_group_v3.group_1 89c60255-9bd6-460c-822a-e2b959ede9d2
```
