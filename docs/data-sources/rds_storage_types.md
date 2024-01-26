---
subcategory: "Relational Database Service (RDS)"
---

# flexibleengine_rds_storage_types

Use this data source to get the list of RDS storage types.

## Example Usage

```hcl
variable "instance_id" {}

data "flexibleengine_rds_storage_types" "test" {
  db_type    = "MySQL"
  db_version = "8.0"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `db_type` - (Required, String) Specifies the DB engine type. Its value can be any of the following and
  is case-insensitive: **MySQL**, **PostgreSQL**, **SQLServer** and **MariaDB**.

* `db_version` - (Required, String) Specifies the database version. For details about how to obtain the database
  version, see section [Querying Version Information About a DB Engine](https://docs.prod-cloud-ocb.orange-business.com/en-us/api/rds/rds_06_0001.html).

* `instance_mode` - (Optional, String) Specifies the HA mode. The value options are as
  follows: **single**, **ha**, **replica**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `storage_types` - Indicates the DB instance specifications information list. For details, see Data structure of
  the storage_type field. The [storage_types](#Storagetype_storageType) structure is documented below.

<a name="Storagetype_storageType"></a>
The `storage_types` block supports:

* `name` - Indicates the storage type. Its value can be any of the following:
  - **COMMON**: Indicates the SATA type.
  - **ULTRAHIGH**: Indicates the SSD type.

* `az_status` - The status details of the AZs to which the specification belongs.
  Key indicates the AZ ID, and value indicates the specification status in the AZ.
  The options of value are as follows:
    - **normal**: The specifications in the AZ are available.
    - **unsupported**: The specifications are not supported by the AZ.
    - **sellout**: The specifications in the AZ are sold out.

* `support_compute_group_type` - Performance specifications.
  The options are as follows:
    - **normal**: General-enhanced.
    - **normal2**: General-enhanced II.
    - **armFlavors**: Kunpeng general-enhanced.
    - **dedicicatenormal**: Exclusive x86.
    - **armlocalssd**: Standard Kunpeng.
    - **normallocalssd**: Standard x86.
    - **general**: General-purpose.
    - **dedicated**: Dedicated, which is only supported for cloud SSDs.
    - **rapid**: Dedicated, which is only supported for extreme SSDs.
    - **bigmen**: Large-memory.
