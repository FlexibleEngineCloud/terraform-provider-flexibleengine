---
subcategory: "GaussDB NoSQL"
---

# flexibleengine_gaussdb_nosql_flavors

Use this data source to get available FlexibleEngine GaussDB (for NoSQL) flavors.

## Example Usage

```hcl
data "flexibleengine_gaussdb_nosql_flavors" "flavors" {
  vcpus  = 4
  memory = 16
}
```

## Argument Reference

* `region` - (Optional, String) Specifies the region in which to obtain the GaussDB specifications.
  If omitted, the provider-level region will be used.

* `engine` - (Optional, String) Specifies the type of the database engine. The valid values are as follows:
  + **cassandra**: The default value and means to query GaussDB (for Cassandra) instance specifications.
  + **influxdb**: Means to query GaussDB (for Influx) instance specifications.

* `engine_version` - (Optional, String) Specifies the version of the database engine.

* `vcpus` - (Optional, Int) Specifies the number of vCPUs.

* `memory` - (Optional, Int) Specifies the memory size in gigabytes (GB).

* `availability_zone` - (Optional, String) Specifies the availability zone (AZ) of the GaussDB specifications.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Data source ID.

* `flavors` - The information of the GaussDB specifications. Structure is documented below.

The `flavors` block contains:

* `name` - The spec code of the flavor.

* `vcpus` - The number of vCPUs.

* `memory` - The memory size, in GB.

* `engine` - The type of the database engine.

* `engine_version` - The version of the database engine.

* `availability_zones` - All available zones (on sale) for current flavor.
