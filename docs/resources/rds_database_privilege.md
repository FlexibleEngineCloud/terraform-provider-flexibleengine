---
subcategory: "Relational Database Service (RDS)"
description: ""
page_title: "flexibleengine_rds_database_privilege"
---

# flexibleengine_rds_database_privilege

Manages RDS Mysql database privilege resource within FlexibleEngine.

## Example Usage

```hcl
variable "instance_id" {}
variable "db_name" {}
variable "user_name_1" {}
variable "user_name_2" {}

resource "flexibleengine_rds_database_privilege" "test" {
  instance_id = var.instance_id
  db_name     = var.db_name

  users {
    name     = var.user_name_1
    readonly = true
  }

  users {
    name     = var.user_name_2
    readonly = false
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the RDS database privilege resource.
  If omitted, the provider-level region will be used. Changing this will create a new RDS database privilege resource.

* `instance_id` - (Required, String, ForceNew) Specifies the RDS instance ID. Changing this will create a new resource.

* `db_name` - (Required, String, ForceNew) Specifies the database name. Changing this creates a new resource.

* `users` - (Required, List, ForceNew) Specifies the account that associated with the database. This parameter supports
  a maximum of 50 elements. The [users](#rds_users) object structure is documented below.
  Changing this creates a new resource.

<a name="rds_users"></a>
The `users` block supports:

* `name` - (Required, String) Specifies the username of the database account.

* `readonly` - (Optional, Bool) Specifies the read-only permission. The value can be:
  + **true**: indicates the read-only permission.
  + **false**: indicates the read and write permission.

  The default value is **false**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID of database privilege which is formatted `<instance_id>/<database_name>`.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 30 minutes.
* `update` - Default is 30 minutes.
* `delete` - Default is 30 minutes.

## Import

RDS database privilege can be imported using the `instance id` and `database name`, e.g.

```shell
terraform import flexibleengine_rds_database_privilege.test instance_id/database_name
```
