---
subcategory: "Cloud Container Engine (CCE)"
---

# flexibleengine_cce_addon_v3

Provides a CCE addon resource within FlexibleEngine.

## Example Usage

```hcl
variable "cluster_id" {}

resource "flexibleengine_cce_addon_v3" "addon_test" {
  cluster_id    = var.cluster_id
  template_name = "metrics-server"
  version       = "1.0.6"
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required, String, ForceNew) ID of the cluster. Changing this parameter will create a new resource.

* `template_name` - (Required, String, ForceNew) Name of the addon template.
  Changing this parameter will create a new resource.

* `version` - (Required, String, ForceNew) Version of the addon. Changing this parameter will create a new resource.

* `values` - (Optional, List, ForceNew) Add-on template installation parameters.
  These parameters vary depending on the add-on. Changing this parameter will create a new resource.

* The `values` block supports:

* `basic` - (Required, String, ForceNew) The basic parameters in json string fomart.
  Changing this will create a new resource.

* `custom` - (Optional, String, ForceNew) The custom parameters in json string fomart.
  Changing this will create a new resource.

* `flavor` - (Optional, String, ForceNew) The flavor parameters in json string fomart.
  Changing this will create a new resource.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the addon instance.
* `status` - Addon status information.
* `description` - Description of addon instance.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minute.
* `delete` - Default is 3 minute.

## Import

CCE addon can be imported using the cluster ID and addon ID separated by a slash, e.g.:

```
$ terraform import flexibleengine_cce_addon_v3.my_addon bb6923e4-b16e-11eb-b0cd-0255ac101da1/c7ecb230-b16f-11eb-b3b6-0255ac1015a3
```
