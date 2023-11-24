---
subcategory: "Identity and Access Management (IAM)"
---

# flexibleengine_identity_projects

Use this data source to query the IAM project list within FlexibleEngine.

-> **NOTE:** You *must* have IAM read privileges to use this data source.

## Example Usage

### Obtain project information by name

```hcl
data "flexibleengine_identity_projects" "test" {
  name = "eu-west-0_demo"
}
```

### Obtain special project information by name

```hcl
data "flexibleengine_identity_projects" "test" {
  name = "MOS" // The project for OBS Billing
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Optional, String) Specifies the IAM project name to query.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `projects` - The details of the query projects. The [projects](#iam_projects) object structure is documented below.

<a name="iam_projects"></a>
The `projects` block supports:

* `id` - The IAM project ID.

* `name` - The IAM project name.

* `enabled` - Whether the IAM project is enabled.
