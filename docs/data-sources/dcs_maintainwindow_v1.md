---
subcategory: "Distributed Cache Service (DCS)"
---

# flexibleengine_dcs_maintainwindow_v1

Use this data source to get the ID of an available Flexibleengine DCS maintainwindow.

## Example Usage

```hcl

data "flexibleengine_dcs_maintainwindow_v1" "maintainwindow1" {
  default = true
}

```

## Argument Reference

* `region` - (Optional, String) The region in which to obtain the dcs maintenance windows. If omitted, the provider-level
  region will be used.

* `seq` - (Optional, Int) Specifies the sequential number of a maintenance time window.

* `begin` - (Optional, String) Specifies the time at which a maintenance time window starts.

* `end` - (Optional, String) Specifies the time at which a maintenance time window ends.

* `default` - (Optional, Bool) Specifies whether a maintenance time window is set to the default time segment.

## Attribute Reference

`id` is set to the ID of the found maintainwindow. In addition, the following attributes
are exported:

* `begin` - See Argument Reference above.
* `end` - See Argument Reference above.
* `default` - See Argument Reference above.
