---
subcategory: "Key Management Service (KMS)"
---

# flexibleengine_kms_grant

Users can create authorizations for other IAM users or accounts,
granting them permission to use their own master key (CMK),
and a maximum of 100 authorizations can be created under one master key.

## Example Usage

```hcl
variable "key_id" {}
variable "user_id" {}

resource "flexibleengine_kms_grant" "test" {
  key_id            = var.key_id
  type              = "user"
  grantee_principal = var.user_id
  operations        = ["create-datakey", "encrypt-datakey"]
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `key_id` - (Required, String, ForceNew) Key ID.

  Changing this parameter will create a new resource.

* `grantee_principal` - (Required, String, ForceNew) The ID of the authorized user or account.  

  Changing this parameter will create a new resource.

* `operations` - (Required, List, ForceNew) List of granted operations.
  The options are: **create-datakey**, **create-datakey-without-plaintext**, **encrypt-datakey**,
  **decrypt-datakey**, **describe-key**, **create-grant**, **retire-grant**, **encrypt-data**, **decrypt-data**
  A value containing only **create-grant** is invalid.

  Changing this parameter will create a new resource.

* `name` - (Optional, String, ForceNew) Grant name.  
  It must be 1 to 255 characters long, start with a letter, and contain only letters (case-sensitive),
  digits, hyphens (-), underscores (_), and slash(/).

  Changing this parameter will create a new resource.

* `type` - (Optional, String, ForceNew) Authorization type.
  The options are: **user**, **domain**. The default value is **user**.  

  Changing this parameter will create a new resource.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `creator` - The ID of the creator.  

## Import

The kms grant can be imported using
`key_id`, `grant_id`, separated by slashes, e.g.

```bash
terraform import flexibleengine_kms_grant.test <key_id>/<grant_id>
```
