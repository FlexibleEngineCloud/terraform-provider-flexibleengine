---
subcategory: "Cloud Container Engine (CCE)"
---

# flexibleengine_cce_node_ids_v3

`flexibleengine_cce_node_ids_v3` provides a list of node ids for a CCE cluster.
This data source can be useful for getting back a list of node ids for a CCE cluster.

## Example Usage

```hcl
variable "cluster_id" {}

data "flexibleengine_cce_node_ids_v3" "node_ids" {
  cluster_id = var.cluster_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) The region in which to query the data source. If omitted, the provider-level region
  will be used.

* `cluster_id` (Required, String) - Specifies the CCE cluster ID used as the query filter.

## Attribute Reference

The following attributes are exported:

* `ids` - A list of all the node ids found. This data source will fail if none are found.
