---
subcategory: "Resource Template Service (RTS)"
---

# flexibleengine_rts_software_config_v1

The RTS Software Config data source provides details about a specific RTS Software Config.

## Example Usage

```hcl
variable "config_name" {}
variable "server_id" {}

data "flexibleengine_rts_software_config_v1" "myconfig" {
  id = var.config_name
}

resource "flexibleengine_rts_software_deployment_v1" "mydeployment" {
  config_id = data.flexibleengine_rts_software_config_v1.myconfig.id
  server_id = var.server_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) The region in which to query the data source. If omitted, the provider-level region
  will be used.

* `id` - (Optional, String) The id of the software configuration.

* `name` - (Optional, String) The name of the software configuration.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `group` - The namespace that groups this software configuration by when it is delivered to a server.

* `input_values` -  A list of software configuration inputs.

* `output_values` - A list of software configuration outputs.

* `config` - The software configuration code.

* `options` - The software configuration options.
