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

* `default` - (Optional) Specifies whether a maintenance time window is set to the default time segment.

* `seq` - (Optional) Specifies the sequential number of a maintenance time window.

* `begin` - (Optional) Specifies the time at which a maintenance time window starts.

* `end` - (Optional) Specifies the time at which a maintenance time window ends.

## Attributes Reference

`id` is set to the ID of the found maintainwindow. In addition, the following attributes
are exported:

* `begin` - See Argument Reference above.
* `end` - See Argument Reference above.
* `default` - See Argument Reference above.
