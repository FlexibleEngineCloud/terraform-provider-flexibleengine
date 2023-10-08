---
subcategory: "Relational Database Service (RDS)"
description: ""
page_title: "flexibleengine_rds_parametergroup_v3"
---

# flexibleengine_rds_parametergroup_v3

Manages a V3 RDS parametergroup resource.

## Example Usage

```hcl

resource "flexibleengine_rds_parametergroup_v3" "pg_1" {
  name        = "pg_1"
  description = "description_1"

  values = {
    max_connections = "10"
    autocommit      = "OFF"
  }
  datastore {
    type    = "mysql"
    version = "5.6"
  }
}

```

## Argument Reference

The following arguments are supported:

* `name` - (Required, String) The parameter group name. It contains a maximum of 64 characters.

* `description` - (Optional, String) The parameter group description. It contains a maximum of 256 characters and
  cannot contain the following special characters:>!<"&'= the value is left blank by default.

* `values` - (Optional, Map) Parameter group values key/value pairs defined by users based on the default
  parameter groups.

* `datastore` - (Required, List, ForceNew) Database object. The [datastore](#rds_datastore) object structure is
  documented below. Changing this creates a new parameter group.

<a name="rds_datastore"></a>
The `datastore` block supports:

* `type` - (Required, String) The DB engine. Currently, MySQL, PostgreSQL, and Microsoft SQL Server are supported.
  The value is case-insensitive and can be mysql, postgresql, or sqlserver.

* `version` - (Required, String) Specifies the database version.

  + MySQL databases support MySQL 5.6 and 5.7. Example value: 5.7.
  + PostgreSQL databases support PostgreSQL 9.5 and 9.6. Example value: 9.5.
  + Microsoft SQL Server databases support 2014 SE, 2016 SE, and 2016 EE. Example value: 2014_SE.

## Attribute Reference

The following attributes are exported:

* `id` -  ID of the parameter group.

* `configuration_parameters` - Indicates the parameter configuration defined by users based on the default
  parameters groups. The [configuration_parameters](#rds_configuration_parameters) object structure is documented below.

<a name="rds_configuration_parameters"></a>
The `configuration_parameters` block supports:

* `name` - Indicates the parameter name.

* `value` - Indicates the parameter value.

* `restart_required` - Indicates whether a restart is required.

* `readonly` - Indicates whether the parameter is read-only.

* `value_range` - Indicates the parameter value range.

* `type` - Indicates the parameter type.

* `description` - Indicates the parameter description.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.
* `delete` - Default is 10 minutes.

## Import

Parameter groups can be imported using the `id`, e.g.

```shell
terraform import flexibleengine_rds_parametergroup_v3.pg_1 7117d38e-4c8f-4624-a505-bd96b97d024c
```
