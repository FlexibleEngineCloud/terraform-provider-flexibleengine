---
subcategory: "Deprecated"
---

# flexibleengine_dcs_az_v1

Use this data source to get the ID of an available Flexibleengine dcs az.

!> **Warning:** It has been deprecated, you can use the availability zone code directly, e.g. eu-west-0b.

## Example Usage

```hcl

data "flexibleengine_dcs_az_v1" "az1" {
  name = "AZ1"
  port = "8004"
  code = "sa-chile-1a"
}
```

## Argument Reference

* `name` - (Optional) Indicates the name of an AZ.

* `code` - (Optional) Indicates the code of an AZ.

* `port` - (Optional) Indicates the port number of an AZ.

## Attributes Reference

`id` is set to the ID of the found az. In addition, the following attributes
are exported:

* `name` - See Argument Reference above.
* `code` - See Argument Reference above.
* `port` - See Argument Reference above.
