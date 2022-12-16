---
subcategory: "Key Management Service (KMS)"
description: ""
page_title: "flexibleengine_kms_key_v1"
---

# flexibleengine_kms_key_v1

Manages a V1 key resource within KMS.

## Example Usage

```hcl
resource "flexibleengine_kms_key_v1" "key_1" {
  key_alias       = "key_1"
  pending_days    = "7"
  key_description = "first test key"
  realm           = "cn-north-1"
  is_enabled      = true
}
```

## Argument Reference

The following arguments are supported:

* `key_alias` - (Required) Specifies the name of a KMS key.

* `key_description` - (Optional) Specifies the description of a KMS key.

* `realm` - (Optional) Region where a key resides. Changing this creates a new key.

* `pending_days` - (Optional) Specifies the duration in days after which the key is deleted
    after destruction of the resource, must be between 7 and 1096 days. Defaults to 7.
    It only be used when delete a key.

* `is_enabled` - (Optional) Specifies whether the key is enabled. Defaults to true.

* `rotation_enabled` - (Optional) Specifies whether the key rotation is enabled. Defaults to false.

* `rotation_interval` - (Optional) Specifies the key rotation interval. The valid value is range from 30 to 365,
  defaults to 365.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The globally unique identifier for the key.
* `default_key_flag` - Identification of a Master Key. The value 1 indicates a Default
    Master Key, and the value 0 indicates a key.
* `origin` - Origin of a key. The default value is kms.
* `domain_id` - ID of a user domain for the key.
* `creation_date` - Creation time (time stamp) of a key.
* `rotation_number` - The total number of key rotations.

## Import

KMS Keys can be imported using the `id`, e.g.

```
$ terraform import flexibleengine_kms_key_v1.key_1 7056d636-ac60-4663-8a6c-82d3c32c1c64
```
