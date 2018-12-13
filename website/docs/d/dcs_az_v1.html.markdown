---
layout: "flexibleengine"
page_title: "Flexibleengine: flexibleengine_dcs_az_v1"
sidebar_current: "docs-flexibleengine-datasource-dcs-az-v1"
description: |-
  Get information on an Flexibleengine dcs az.
---

# flexibleengine\_dcs\_az_v1

Use this data source to get the ID of an available Flexibleengine dcs az.

## Example Usage

```hcl

data "flexibleengine_dcs_az_v1" "az1" {
  name = "AZ1"
  port = "8004"
  code = "sa-chile-1a"
}
```

## Argument Reference

* `name` - (Required) Indicates the name of an AZ.

* `code` - (Optional) Indicates the code of an AZ.

* `port` - (Required) Indicates the port number of an AZ.


## Attributes Reference

`id` is set to the ID of the found az. In addition, the following attributes
are exported:

* `name` - See Argument Reference above.
* `code` - See Argument Reference above.
* `port` - See Argument Reference above.
