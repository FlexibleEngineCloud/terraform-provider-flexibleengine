---
subcategory: "Document Database Service (DDS)"
---

# flexibleengine_dds_parameter_template

Manages a DDS parameter template resource within FlexibleEngine.

## Example Usage

```hcl
variable "name" {}
variable "parameter_values" {}
variable "node_type" {}
variable "node_version" {}

resource "flexibleengine_dds_parameter_template" "test"{
  name             = var.name
  parameter_values = var.parameter_values
  node_type        = var.node_type
  node_version     = var.node_version
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `name` - (Required, String) Specifies the parameter template name.
  The value must be 1 to 64 characters in length and start with a letter (from A to Z or from a to z).
  It is case-sensitive and can contain only letters, digits (from 0 to 9), hyphens (-), and underscores (_).

* `node_type` - (Required, String, ForceNew) Specifies the node type of parameter template node_type. Valid value:
  + **mongos**: the mongos node type.
  + **shard**: the shard node type.
  + **config**: the config node type.
  + **replica**: the replica node type.
  + **single**: the single node type.

  Changing this parameter will create a new resource.

* `node_version` - (Required, String, ForceNew) Specifies the database version.
  The value can be **4.2**, **4.0** or **3.4**.

  Changing this parameter will create a new resource.

* `parameter_values` - (Optional, Map) Specifies the mapping between parameter names and parameter values.
  You can customize parameter values based on the parameters in the default parameter template.

* `description` - (Optional, String) Specifies the parameter template description.
  The description must consist of a maximum of 256 characters and cannot contain the carriage
  return character or the following special characters: >!<"&'=.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `parameters` - Indicates the parameters defined by users based on the default parameter templates.
  The [parameters](#DdsParameterTemplate_Parameter) structure is documented below.

<a name="DdsParameterTemplate_Parameter"></a>
The `parameters` block supports:

* `name` - Indicates the parameter name.

* `value` - Indicates the parameter value.

* `description` - Indicates the parameter description.

* `type` - Indicates the parameter type. The value can be integer, string, boolean, float, or list.

* `value_range` - Indicates the value range.

* `restart_required` - Indicates whether the instance needs to be restarted.
  + If the value is **true**, restart is required.
  + If the value is **false**, restart is not required.

* `readonly` - Indicates whether the parameter is read-only.
  + If the value is **true**, the parameter is read-only.
  + If the value is **false**, the parameter is not read-only.

## Import

The DDS parameter template can be imported using the `id`, e.g.

```shell
terraform import flexibleengine_dds_parameter_template.test <tempalate_id>
```
