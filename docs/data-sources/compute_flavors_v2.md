---
subcategory: "Elastic Cloud Server (ECS)"
---

# flexibleengine_compute_flavors_v2

Use this data source to get the available Compute Flavors.

## Example Usage

```hcl
data "flexibleengine_compute_flavors_v2" "flavors" {
  availability_zone = "eu-west-0a"
  performance_type  = "normal"
  cpu_core          = 2
  memory_size       = 4
}

# Create ECS instance with the first matched flavor
resource "flexibleengine_compute_instance_v2" "instance" {
  flavor_id = data.flexibleengine_compute_flavors_v2.flavors.flavors[0]
  ...
}
```

## Argument Reference

* `region` - (Optional, String) The region in which to obtain the flavors.
  If omitted, the provider-level region will be used.

* `availability_zone` - (Optional, String) Specifies the AZ name.

* `performance_type` - (Optional, String) Specifies the ECS flavor type.

* `generation` - (Optional, String) Specifies the generation of an ECS type.

* `cpu_core` - (Optional, Int) Specifies the number of vCPUs in the ECS flavor.

* `memory_size` - (Optional, Int) Specifies the memory size(GB) in the ECS flavor.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Specifies a data source ID in UUID format.

* `flavors` - A list of flavors.
