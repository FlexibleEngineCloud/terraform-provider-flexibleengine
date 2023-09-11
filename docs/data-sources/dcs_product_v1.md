---
subcategory: "Distributed Cache Service (DCS)"
---

# flexibleengine_dcs_product_v1

Use this data source to get the ID of an available Flexibleengine DCS product.

## Example Usage

```hcl
# product of Redis 4.0/5.0 with Redis Cluster type
data "flexibleengine_dcs_product_v1" "product1" {
  engine         = "redis"
  engine_version = "4.0;5.0"
  cache_mode     = "cluster"
  capacity       = 8
  replica_count  = 2
}

# product of Redis 4.0/5.0 with Master/Standby type
data "flexibleengine_dcs_product_v1" "product2" {
  engine         = "redis"
  engine_version = "4.0;5.0"
  cache_mode     = "ha"
  capacity       = 0.125
  replica_count  = 2
}

# product of Redis 4.0/5.0 with Single-node type
data "flexibleengine_dcs_product_v1" "product3" {
  engine         = "redis"
  engine_version = "4.0;5.0"
  cache_mode     = "single"
  capacity       = 1
}

# product of Redis 4.0/5.0 with Proxy Cluster type
data "flexibleengine_dcs_product_v1" "product4" {
  engine         = "redis"
  engine_version = "4.0;5.0"
  cache_mode     = "proxy"
  capacity       = 4
}

# product of Redis 3.0 instance
data "flexibleengine_dcs_product_v1" "product5" {
  engine         = "redis"
  engine_version = "3.0"
  cache_mode     = "ha"
}

# product of Memcached instance
data "flexibleengine_dcs_product_v1" "product6" {
  engine     = "memcached"
  cache_mode = "single"
}
```

## Argument Reference

* `engine` - (Optional, String) The engine of the cache instance. Valid values are *redis* and *memcached*.
  Default value is *redis*.

* `engine_version` - (Optional, String) The version of a cache engine.
  It is valid when the engine is *redis*, the value can be `3.0`or `4.0;5.0`.

* `cache_mode` - (Optional, String) The mode of a cache engine. The valid values are as follows:
  + `single` - Single-node.
  + `ha` - Master/Standby.
  + `cluster` - Redis Cluster, it is valid when the engine is *redis*.
  + `proxy` - Proxy Cluster, it is valid when the engine is *redis*.

* `capacity` - (Optional, Float) The total memory of the cache, in GB.
  It is valid when the engine is redis 4.0/5.0.
  + Single-node and Master/Standby instances support:
    `0.125`, `0.25`, `0.5`, `1`, `2`, `4`, `8`, `16`, `24`, `32`, `48` and `64`.
  + Redis Cluster and Proxy Cluster instances support:
    `4`, `8`, `16`, `24`, `32`, `48`, `64`, `96`, `128`, `192`, `256`, `384`, `512`, `768` and `1024`.

* `replica_count` - (Optional, Int) The number of replicas includes the master.
  It is valid when the engine is redis 4.0/5.0 with **Master/Standby** or **Redis Cluster** type.

* `spec_code` - (Optional, String) Specifies the DCS instance specification code. You can log in to the DCS console,
  click *Buy DCS Instance*, and find the corresponding instance specification.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The found product ID.
* `cpu_architecture` - The CPU architecture of DCS instance.
