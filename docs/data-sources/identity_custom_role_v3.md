---
subcategory: "Identity and Access Management (IAM)"
---

# flexibleengine_identity_custom_role_v3

Use this data source to get the ID of an IAM **custom policy**.

## Example Usage

```hcl
data "flexibleengine_identity_custom_role_v3" "role" {
  name = "custom_role"
}
```

## Argument Reference

* `name` - (Optional) Name of the custom policy.

* `id` - (Optional) ID of the custom policy.

* `domain_id` - (Optional) The domain the policy belongs to.

* `references` - (Optional) The number of citations for the custom policy.

* `description` - (Optional) Description of the custom policy.

* `type` - (Optional) Display mode. Valid options are AX: Account level and XA: Project level.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `policy` - Document of the custom policy.

* `catalog` - The catalog of the custom policy.
