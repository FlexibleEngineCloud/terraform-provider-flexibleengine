---
subcategory: "Domain Name Service (DNS)"
---

# flexibleengine\_dns\_zone\_v2

Use this data source to get the ID of an available FlexibleEngine DNS zone.

## Example Usage

```hcl
data "flexibleengine_dns_zone_v2" "zone_1" {
  name = "example.com"
}
```

## Argument Reference

* `region` - (Optional) The region in which to obtain the V2 DNS client.
  A DNS client is needed to retrieve zone ids. If omitted, the
  `region` argument of the provider is used.

* `name` - (Optional) The name of the zone.

* `description` - (Optional) A description of the zone.

* `email` - (Optional) The email contact for the zone record.

* `status` - (Optional) The zone's status.

* `ttl` - (Optional) The time to live (TTL) of the zone.

* `zone_type` - (Optional) The type of the zone. Can either be `public` or `private`.

## Attributes Reference

`id` is set to the ID of the found zone. In addition, the following attributes
are exported:

* `region` - See Argument Reference above.
* `name` - See Argument Reference above.
* `email` - See Argument Reference above.
* `zone_type` - See Argument Reference above.
* `ttl` - See Argument Reference above.
* `description` - See Argument Reference above.
* `status` - See Argument Reference above.
* `masters` - An array of master DNS servers.
* `serial` - The serial number of the zone.
* `pool_id` - The ID of the pool hosting the zone.
* `project_id` - The project ID that owns the zone.
