---
subcategory: "Bare Metal Server (BMS)"
---

# flexibleengine_compute_bms_keypairs_v2

`flexibleengine_compute_bms_keypairs_v2` used to query SSH key pairs.

## Example Usage

```hcl
variable "keypair_name" {}

data "flexibleengine_compute_bms_keypairs_v2" "keypair" {
  name = var.keypair_name
}
```

## Argument Reference

The arguments of this data source act as filters for querying the BMSs details.

* `region` - (Optional, String) The region in which to query the data source. If omitted, the provider-level region
  will be used.

* `name` - (Required, String) - It is the key pair name.

## Attribute Reference

All of the argument attributes are also exported as result attributes.

* `public_key` - It gives the information about the public key in the key pair.

* `fingerprint` - It is the fingerprint information about the key pair.
