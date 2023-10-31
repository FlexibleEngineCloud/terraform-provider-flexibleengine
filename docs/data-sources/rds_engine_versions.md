---
subcategory: "Relational Database Service (RDS)"
---

# flexibleengine_rds_engine_versions

Use this data source to obtain all version information of the specified engine type of FlexibleEngine.

## Example Usage

```hcl
data "flexibleengine_rds_engine_versions" "test" {
  type = "SQLServer"
}
```

## Argument Reference

* `region` - (Optional, String) The region in which to obtain the RDS engine versions.
  If omitted, the provider-level region will be used.

* `type` - (Optional, String) Specifies the RDS engine type.
  The valid values are **MySQL**, **PostgreSQL**, **SQLServer** and **MariaDB**, default to **MySQL**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Data source ID in hashcode format.

* `versions` - Indicates the list of database versions. The [versions](#rds_versions) object structure is
  documented below.

<a name="rds_versions"></a>
The `versions` block supports:

* `id` - Indicates the database version ID. Its value is unique.

* `name` - Indicates the database version number. Only the major version number (two digits) is returned.
  For example, if the version number is MySQL 5.6.X, only 5.6 is returned.
