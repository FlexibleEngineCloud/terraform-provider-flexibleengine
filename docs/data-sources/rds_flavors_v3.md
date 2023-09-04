---
subcategory: "Relational Database Service (RDS)"
---

# flexibleengine_rds_flavors_v3

Use this data source to get available FlexibleEngine RDS flavors.

## Example Usage

```hcl
data "flexibleengine_rds_flavors_v3" "flavor" {
  db_type       = "PostgreSQL"
  db_version    = "12"
  instance_mode = "ha"
  vcpus         = 4
}
```

## Argument Reference

* `db_type` - (Required, String) Specifies the DB engine. Value: MySQL, PostgreSQL, SQLServer.

* `db_version` - (Optional, String) Specifies the database version. MySQL databases support MySQL 5.6
  and 5.7. PostgreSQL databases support PostgreSQL 9.5 and 9.6. Microsoft SQL Server databases support
  2014_SE, 2016_SE, and 2016_EE.

* `instance_mode` - (Optional, String) The mode of instance. Value: *ha*(indicates primary/standby instance),
  *single*(indicates single instance) and *replica*(indicates read replicas).

* `vcpus` - (Optional, Int) Specifies the number of vCPUs in the RDS flavor.

* `memory` - (Optional, Int) Specifies the memory size(GB) in the RDS flavor.

* `group_type` - (Optional, String) Specifies the performance specification, the valid values are as follows:
  + **normal**: General enhanced.
  + **normal2**: General enhanced type II.
  + **dedicatedNormal**: (dedicatedNormalLocalssd): Dedicated for x86.
  + **normalLocalssd**: x86 general type.
  + **general**: General type.
  + **bigmem**: Large memory type.

* `availability_zone` - (Optional, String) Specifies the availability zone which the RDS flavor belongs to.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `flavors` - Indicates the flavors information. Structure is documented below.

The `flavors` block contains:

* `id` - The ID of the rds flavor.
* `name` - The name of the rds flavor.
* `vcpus` - The CPU size.
* `memory` - The memory size in GB.
* `group_type` - The performance specification.
* `instance_mode` - The mode of instance.
* `availability_zones` - The availability zones which the RDS flavor belongs to.
* `db_versions` - The Available versions of the database.
