---
subcategory: "Distributed Cache Service (DCS)"
---

# flexibleengine_dcs_flavors

Use this data source to get a list of available DCS flavors.

## Example Usage

```hcl
data "flexibleengine_dcs_flavors" "flavors" {
  capacity = "4"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to obtain the DCS flavors.
  If omitted, the provider-level region will be used.

* `capacity` - (Required, Float) The total memory of the cache, in GB.
  + **Redis4.0, Redis5.0 and Redis6.0**: Stand-alone and active/standby type instance values:
    `0.125`, `0.25`, `0.5`, `1`, `2`, `4`, `8`, `16`, `32` and `64`.
    Cluster instance specifications support `4`,`8`,`16`,`24`, `32`, `48`, `64`, `96`, `128`, `192`,
    `256`, `384`, `512`, `768` and `1024`.
  + **Redis3.0**: Stand-alone and active/standby type instance values: `2`, `4`, `8`, `16`, `32` and `64`.
    Proxy cluster instance specifications support `64`, `128`, `256`, `512`, and `1024`.
  + **Memcached**: Stand-alone and active/standby type instance values: `2`, `4`, `8`, `16`, `32` and `64`.

* `engine` - (Optional, String) The engine of the cache instance. Valid values are *Redis* and *Memcached*.
  Default value is *Redis*.

* `engine_version` - (Optional, String) The version of a cache engine.
  It is mandatory when the engine is *Redis*, the value can be `3.0`, `4.0`, `5.0`, or `6.0`.

* `cache_mode` - (Optional, String) The mode of a cache engine. The valid values are as follows:
  + `single` - Single-node.
  + `ha` - Master/Standby.
  + `cluster` - Redis Cluster.
  + `proxy` - Proxy Cluster. Redis6.0 not support this mode.
  + `ha_rw_split` - Read/Write splitting. Redis6.0 not support this mode.
  
* `name` - (Optional, String) The flavor name of the cache instance.

* `cpu_architecture` - (Optional, String) The CPU architecture of cache instance.
  Valid values *x86_64* and *aarch64*.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `flavors` - A list of DCS flavors.

The `flavors` block supports:

* `name` - The flavor name of the cache instance.

* `cache_mode` - The mode of a cache instance.

* `engine` - The engine of the cache instance. Value is *redis* or *memcached*.

* `engine_versions` - Supported versions of the specification.

* `cpu_architecture` - The CPU architecture of cache instance. Value is *x86_64* or *aarch64*.

* `capacity` - The total memory of the cache, in GB.

* `available_zones` - An array of available zones where the cache specification can be used.

* `ip_count` - Number of IP addresses corresponding to the specifications.
