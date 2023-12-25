---
subcategory: "Relational Database Service (RDS)"
description: ""
page_title: "flexibleengine_rds_database"
---

# flexibleengine_rds_database

Manages RDS Mysql database resource within FlexibleEngine.

## Example Usage

```hcl
variable "instance_id" {}

resource "flexibleengine_rds_database" "test" {
  instance_id   = var.instance_id
  name          = "test"
  character_set = "utf8"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the RDS database resource.
  If omitted, the provider-level region will be used. Changing this will create a new RDS database resource.

* `instance_id` - (Required, String, ForceNew) Specifies the RDS instance ID. Changing this will create a new resource.

* `name` - (Required, String, ForceNew) Specifies the database name. The database name contains **1** to **64**
  characters. The name can only consist of lowercase letters, digits, hyphens (-), underscores (_) and dollar signs
  ($). The total number of hyphens (-) and dollar signs ($) cannot exceed **10**. RDS for **MySQL 8.0** does not
  support dollar signs ($). Changing this will create a new resource.

* `character_set` - (Required, String, ForceNew) Specifies the character set used by the database, For example **utf8**,
  **gbk**, **ascii**, etc. Changing this will create a new resource.

* `description` - (Optional, String) Specifies the database description. The value can contain **0** to **512**
  characters. This parameter takes effect only for DB instances whose kernel versions are at least **5.6.51.3**,
  **5.7.33.1**, or **8.0.21.4**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID of database which is formatted `<instance_id>/<database_name>`.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.
* `update` - Default is 10 minutes.
* `delete` - Default is 10 minutes.

## Import

RDS database can be imported using the `instance id` and `database name`, e.g.

```shell
terraform import flexibleengine_rds_database.database_1 instance_id/database_name
```
