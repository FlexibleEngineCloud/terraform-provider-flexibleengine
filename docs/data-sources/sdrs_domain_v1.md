---
subcategory: "Storage Disaster Recovery Service (SDRS)"
---

# flexibleengine\_sdrs\_domain_v1

Use this data source to get the ID of an available FlexibleEngine SDRS domain.

## Example Usage

```hcl

data "flexibleengine_sdrs_domain_v1" "dom_1" {
  name = "SDRS_HypeDomain01"
}

```

## Argument Reference

* `name` - (Optional) Specifies the name of an active-active domain. Currently only support SDRS_HypeDomain01.

## Attributes Reference

`id` is set to the ID of the active-active domain. In addition, the following attributes
are exported:

* `name` - See Argument Reference above.
* `description` - Specifies the description of an active-active domain.
