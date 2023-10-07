---
subcategory: "Resource Template Service (RTS)"
description: ""
page_title: "flexibleengine_rts_software_config_v1"
---

# flexibleengine_rts_software_config_v1

Provides an RTS software config resource.

## Example Usage

 ```hcl
variable "config_name" {}
 
resource "flexibleengine_rts_software_config_v1" "myconfig" {
  name = var.config_name
}
 ```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the RTS software resource.
  If omitted, the provider-level region will be used. Changing this will create a new RTS software resource.

* `name` - (Required, String, ForceNew) The name of the software configuration. Changing this will create a new RTS
  software resource.

* `group` - (Optional, String, ForceNew) The namespace that groups this software configuration by when it is delivered
  to a server. Changing this will create a new RTS software resource.

* `input_values` - (Optional, List, ForceNew) A list of software configuration inputs. Changing this will create a new
  RTS software resource.

* `output_values` - (Optional, List, ForceNew) A list of software configuration outputs. Changing this will create a
  new RTS software resource.

* `config` - (Optional, String, ForceNew) The software configuration code. Changing this will create a new RTS software
  resource.

* `options` - (Optional, Map, ForceNew) The software configuration options. Changing this will create a new RTS software
  resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The id of the software config.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.
* `delete` - Default is 3 minutes.

## Import

Software Config can be imported using the `config id`, e.g.

```shell
terraform import flexibleengine_rts_software_config_v1 4779ab1c-7c1a-44b1-a02e-93dfc361b32d
```
