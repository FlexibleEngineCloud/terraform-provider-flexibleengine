---
subcategory: "Relational Database Service (RDS)"
---

# flexibleengine_rds_flavors_v3

Use this data source to get available FlexibleEngine rds flavors.

## Example Usage

```hcl
data "flexibleengine_rds_flavors_v3" "flavor" {
    db_type = "PostgreSQL"
    db_version = "9.5"
    instance_mode = "ha"
}
```

## Argument Reference

* `db_type` - (Required, String) Specifies the DB engine. Value: MySQL, PostgreSQL, SQLServer.

* `db_version` - (Required, String) Specifies the database version. MySQL databases support MySQL 5.6
  and 5.7. PostgreSQL databases support PostgreSQL 9.5 and 9.6. Microsoft SQL Server databases support
  2014_SE, 2016_SE, and 2016_EE.

* `instance_mode` - (Required, String) The mode of instance. Value: ha(indicates primary/standby instance), single(indicates single instance)

* `vcpus` - (Optional, Int) Specifies the number of vCPUs in the RDS flavor.

* `memory` - (Optional, Int) Specifies the memory size(GB) in the RDS flavor.

## Attributes Reference

In addition, the following attributes are exported:

* `flavors` - Indicates the flavors information. Structure is documented below.

The `flavors` block contains:

* `name` - The name of the rds flavor.
* `vcpus` - Indicates the CPU size.
* `memory` - Indicates the memory size in GB.
* `mode` - See 'instance_mode' above.
