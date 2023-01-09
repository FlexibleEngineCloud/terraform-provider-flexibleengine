---
subcategory: "Data Warehouse Service (DWS)"
---

# flexibleengine_dws_flavors

Use this data source to get available flavors of FlexibleEngine DWS cluster node.

## Example Usage

```hcl
data "flexibleengine_dws_flavors" "flavor" {
  availability_zone = "eu-west-0a"
}
```

## Argument Reference

* `region` - (Optional, String) Specifies the region in which to obtain the DWS cluster client.
  If omitted, the provider-level region will be used.

* `availability_zone` - (Optional, String) Specifies the availability zone name.

* `vcpus` - (Optional, String) Specifies the vcpus of the DWS node flavor.

* `memory` - (Optional, String) Specifies the ram of the DWS node flavor in GB.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Indicates a data source ID.

* `flavors` - Indicates the flavors information. Structure is documented below.

The `flavors` block contains:

* `flavor_id` - The name of the DWS node flavor. It is referenced by **node_type** in `flexibleengine_dws_cluster_v1`.
* `vcpus` - Indicates the vcpus of the DWS node flavor.
* `memory` - Indicates the ram of the DWS node flavor in GB.
* `volumetype` - Indicates Disk type.
* `size` - Indicates the Disk size in GB.
* `availability_zone` - Indicates the availability zone where the node resides.
