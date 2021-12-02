---
subcategory: "Distributed Cache Service (DCS)"
---

# flexibleengine_dcs_product_v1

Use this data source to get the ID of an available Flexibleengine dcs product.

## Example Usage

```hcl
data "flexibleengine_dcs_product_v1" "product1" {
  engine = "redis"
}

data "flexibleengine_dcs_product_v1" "product2" {
  spec_code = "redis.cluster.xu1.large.r1.8"
}
```

## Argument Reference

* `engine` - (Optional, String) The engine of the cache instance. Valid values are *redis* and *memcached*.
  Default value is *redis*.

* `engine_version` - (Optional, String) The version of a cache engine.
  It is valid when the engine is *redis*, the value can be `3.0`or `4.0;5.0`.

* `spec_code` - (Optional, String) Specifies the DCS instance specification code. You can log in to the DCS console,
  click *Buy DCS Instance*, and find the corresponding instance specification.

* `cache_mode` - (Optional, String) The mode of a cache engine. The valid values are as follows:
  + `single` - Single-node.
  + `ha` - Master/Standby.
  + `cluster` - Redis Cluster.
  + `proxy` - Proxy Cluster.


## Attributes Reference

`id` is set to the ID of the found product. In addition, the following attributes
are exported:

* `engine` - See Argument Reference above.
* `engine_version` - See Argument Reference above.
* `spec_code` - See Argument Reference above.
* `cache_mode` - See Argument Reference above.
