---
subcategory: "Data Sources"
description: ""
page_title: "flexibleengine_availability_zones"
---

# flexibleengine_availability_zones

Use this data source to get a list of availability zones from FlexibleEngine.

## Example Usage

```hcl
data "flexibleengine_availability_zones" "zones" {}
```

## Argument Reference

* `region` - (Optional, String) The region in which to query the data source. If omitted, the provider-level region
  will be used.

* `state` - (Optional, String) The `state` of the availability zones to match, default ("available").

## Attribute Reference

`id` is set to hash of the returned zone list. In addition, the following attributes are exported:

* `names` - The names of the availability zones, ordered alphanumerically, that match the queried `state`.
