---
subcategory: "Identity and Access Management (IAM)"
---

# flexibleengine\_identity\_role\_v3

custom role management in FlexibleEngine

## Example Usage

### Role

```hcl
resource "flexibleengine_identity_role_v3" "role" {
  name        = "test"
  description = "created by terraform"
  type        = "AX"

  policy = <<EOF
{
  "Version": "1.1",
  "Statement": [
    {
      "Action": [
        "obs:bucket:GetBucketAcl"
      ],
      "Effect": "Allow",
      "Resource": [
        "obs:*:*:bucket:*"
      ],
      "Condition": {
        "StringStartWith": {
          "g:ProjectName": [
            "eu-west-0"
          ]
        }
      }
    }
  ]
}
EOF
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, String) Name of the custom policy.

* `description` - (Required, String) Description of the custom policy.

* `type` - (Required, String) Display mode. Valid options are AX: Account level and XA: Project level.

* `policy` - (Required, String) Document of the custom policy.

- - -

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:

* `id` - The role id.

* `domain_id` - The account id.

* `references` - The number of references.

## Import

Role can be imported using the following format:

```
$ terraform import flexibleengine_identity_role_v3.default {{ resource id}}
```
