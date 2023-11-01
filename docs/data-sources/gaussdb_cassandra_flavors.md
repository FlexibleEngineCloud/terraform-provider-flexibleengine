---
subcategory: "GaussDB NoSQL"
---

# flexibleengine_gaussdb_cassandra_flavors

Use this data source to get available FlexibleEngine gaussdb cassandra flavors.

## Example Usage

```hcl
data "flexibleengine_gaussdb_cassandra_flavors" "flavors" {
}
```

## Argument Reference

* `region` - (Optional, String) The region in which to obtain the flavors. If omitted, the provider-level region will be
  used.

* `vcpus` - (Optional, String) Specifies the count of vcpus of the flavors.

* `memory` - (Optional, String) Specifies the memory size of the flavors.

* `version` - (Optional, String) Specifies the engine version of the flavors.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Specifies the data source ID.

* `flavors` - Indicates the flavors information. The [flavors](#gaussdb_flavors) object structure is documented below.

<a name="gaussdb_flavors"></a>
The `flavors` block supports:

* `name` - Indicates the spec code of the flavor.

* `vcpus` - Indicates the CPU size.

* `memory` - Indicates the memory size in GB.

* `version` - Indicates the database version.

* `az_status` - Indicates the flavor status in each availability zone.
