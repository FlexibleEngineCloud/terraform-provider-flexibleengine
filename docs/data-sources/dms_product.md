---
subcategory: "Distributed Message Service (DMS)"
---

# flexibleengine_dms_product

Use this data source to get details about an available FlexibleEngine DMS product.

## Example Usage

```hcl
data "flexibleengine_dms_product" "product1" {
  engine    = "kafka"
  bandwidth = "300MB"
}
```

## Argument Reference

* `region` - (Optional, String) Specifies the region in which to obtain the DMS products.
  If omitted, the provider-level region will be used.

* `bandwidth` - (Required, String) Specifies the bandwidth of a DMS instance.
  The valid values are **100MB**, **300MB**, **600MB** and **1200MB**.

* `engine` - (Optional, String) Specifies the name of a message engine. Only **kafka** is supported.

* `version` - (Optional, String) Specifies the version of a message engine. The default value is **2.3.0**.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The DMS product ID.

* `availability_zones` - The list of availability zones where there are available resources.

* `spec_code` - The DMS product specification, for example, dms.instance.kafka.cluster.c3.small.2.

* `cpu_arch` - The CPU architecture of a DMS instance.

* `ecs_flavor_id` - The flavor of the corresponding ECS.

* `partition_num` - The maximum number of topics in a Kafka instance.

* `storage_space` - The minimum storage capacity of the DMS product.

* `storage_spec_codes` - The list of supported storage specification.
  The item of the list can be one of **dms.physical.storage.ultra** and **dms.physical.storage.high**.

* `max_tps` - The maximum number of messages per unit time.
