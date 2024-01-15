---
subcategory: "Storage Disaster Recovery Service (SDRS)"
---

# flexibleengine_sdrs_domain_v1

Use this data source to get an available SDRS domain.

## Example Usage

```hcl

data "flexibleengine_sdrs_domain_v1" "dom_1" {
  name = "SDRS_HypeDomain01"
}

```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `name` - (Optional, String) Specifies the name of an available SDRS domain.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `description` - Indicates the description of the SDRS domain.
