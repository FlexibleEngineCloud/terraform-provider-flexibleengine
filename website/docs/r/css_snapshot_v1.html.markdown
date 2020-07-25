---
layout: "flexibleengine"
page_title: "FlexibleEngine: flexibleengine_css_snapshot_v1"
sidebar_current: "docs-flexibleengine-resource-css-snapshot-v1"
description: |-
 CSS cluster snapshot management
---

# flexibleengine\_css\_cluster\_v1

CSS cluster snapshot management

## Example Usage

### create a snapshot

```hcl
resource "flexibleengine_css_snapshot_v1" "snapshot" {
  name       = "terraform_test_cluster"
  cluster_id = var.css_cluster_id
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the snapshot name. The snapshot name must
  start with a letter and contains 4 to 64 characters consisting of only
  lowercase letters, digits, hyphens (-), and underscores (_).
  Changing this parameter will create a new resource.

* `cluster_id` - (Required) Specifies ID of the CSS cluster where index data is to be backed up.
  Changing this parameter will create a new resource.

* `indices` - (Optional) Specifies the name of the index to be backed up. Multiple index names
  are separated by commas (,). By default, data of all indices is backed up. You can use the
  asterisk (*) to back up data of certain indices. For example, if you enter 2018-06*, then
  data of indices with the name prefix of 2018-06 will be backed up.
  The value contains 0 to 1024 characters. Uppercase letters, spaces, and certain special
  characters (including "\<|>/?) are not allowed.
  Changing this parameter will create a new resource.

* `description` - (Optional) Specifies the description of a snapshot. The value contains 0 to 256 characters, and angle brackets (<) and (>) are not allowed.
  Changing this parameter will create a new resource.


## Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:

* `status` - Indicates the snapshot status.

* `cluster_name` - Indicates the CSS cluster name.

* `backup_type` - Indicates the snapshot creation mode, the value should be "manual" or "automated".


## Import

This resource can be imported by specifying the CSS cluster ID and snapshot ID
separated by a slash, e.g.:

```
$ terraform import flexibleengine_css_snapshot_v1.snapshot_1 <cluster_id>/<snapshot_id>
```
