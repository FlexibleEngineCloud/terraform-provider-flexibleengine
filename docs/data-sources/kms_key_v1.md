---
subcategory: "Key Management Service (KMS)"
---

# flexibleengine_kms_key_v1

Use this data source to get the ID of an available FlexibleEngine KMS key.

## Example Usage

```hcl
data "flexibleengine_kms_key_v1" "key_1" {
  key_alias = "test_key"
}
```

## Argument Reference

* `key_id` - (Optional) The globally unique identifier for the key. Changing this gets the new key.

* `key_alias` - (Optional) The alias in which to create the key. It is required when
    we create a new key. Changing this gets the new key.

* `key_description` - (Optional) The description of the key as viewed in FlexibleEngine console.
    Changing this gets a new key.

* `key_state` - (Optional) The state of a key. "2" indicates that the key is enabled.
    "3" indicates that the key is disabled. "4" indicates that the key is scheduled for deletion.
    Changing this gets a new key.

* `default_key_flag` - (Optional) Identification of a Master Key. The value "1" indicates a Default
    Master Key, and the value "0" indicates a key. Changing this gets a new key.

* `domain_id` - (Optional) ID of a user domain for the key. Changing this gets a new key.

* `origin` - (Optional) Origin of a key. such as: kms. Changing this gets a new key.

* `realm` - (Optional) Region where a key resides. Changing this gets a new key.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.
* `creation_date` - Creation time (time stamp) of a key.
* `scheduled_deletion_date` - Scheduled deletion time (time stamp) of a key.
* `rotation_enabled` - Indicates whether the key rotation is enabled or not.
* `rotation_interval` - The key rotation interval. It's valid when rotation is enabled.
* `rotation_number` - The total number of key rotations. It's valid when rotation is enabled.
