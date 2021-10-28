---
subcategory: "Document Database Service (DDS)"
---

# flexibleengine_dds_flavors_v3

Use this data source to get the details of available DDS flavors.

## Example Usage

```hcl
data "flexibleengine_dds_flavors_v3" "flavor" {
  engine_name = "DDS-Community"
  vcpus       = 8
  memory      = 32
}
```

## Argument Reference

* `engine_name` - (Optional, String) Specifies the engine name of the dds, the default value is
  "DDS-Community".

* `type` - (Optional, String) Specifies the type of the dds falvor. "mongos", "shard", "config",
  "replica" and "single" are supported.

* `vcpus` - (Optional, String) Specifies the vcpus of the dds flavor.

* `memory` - (Optional, String) Specifies the ram of the dds flavor in GB.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Specifies a data source ID.

* `flavors` - Indicates the flavors information. Structure is documented below.

The `flavors` block contains:

* `spec_code` - The name of the dds flavor.
* `type` - See `type` above.
* `vcpus` - See `vcpus` above.
* `memory` - See `memory` above.
