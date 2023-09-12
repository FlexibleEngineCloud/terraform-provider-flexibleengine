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

* `datastore_type` - (Optional, String) The type of datastore. The options are as follows:
  - **dws**: OLAP, elastic scaling, unlimited scaling of compute and storage capacity.
  - **hybrid**: a single data warehouse used for transaction and analytics workloads,
    in single-node or cluster mode.
  - **stream**: built-in time series operators; up to 40:1 compression ratio; applicable to IoT services.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Indicates a data source ID.

* `flavors` - Indicates the flavors information. The [flavors](#dws_flavors) object structure is documented below.

<a name="dws_flavors"></a>
The `flavors` block supports:

* `flavor_id` - The name of the DWS node flavor. It is referenced by **node_type** in `flexibleengine_dws_cluster_v1`.

* `vcpus` - Indicates the vcpus of the DWS node flavor.

* `memory` - Indicates the ram of the DWS node flavor in GB.

* `volumetype` - Indicates Disk type.

* `size` - Indicates the Disk size in GB.

* `availability_zones` - Indicates the availability zone where the node resides.

* `datastore_type` - The type of datastore.The options are as follows:
  - **dws**: OLAP, elastic scaling, unlimited scaling of compute and storage capacity.
  - **hybrid**: a single data warehouse used for transaction and analytics workloads,
    in single-node or cluster mode.
  - **stream**: built-in time series operators; up to 40:1 compression ratio; applicable to IoT services.

* `elastic_volume_specs` - The [elastic_volume_specs](#dws_elastic_volume_specs) object structure is documented below.

<a name="dws_elastic_volume_specs"></a>
The `elastic_volume_specs` block supports:

* `step` - Disk size increment step.

* `min_size` - Minimum disk size.

* `max_size` - Maximum disk size.
