---
subcategory: "Domain Name Service (DNS)"
description: ""
page_title: "flexibleengine_dns_recordset_v2"
---

# flexibleengine_dns_recordset_v2

Manages a DNS record set in the FlexibleEngine DNS Service.

## Example Usage

### Automatically detect the correct network

```hcl
resource "flexibleengine_dns_zone_v2" "example_zone" {
  name        = "example.com."
  email       = "email2@example.com"
  description = "a zone"
  ttl         = 6000
  zone_type   = "public"
}

resource "flexibleengine_dns_recordset_v2" "rs_example_com" {
  zone_id     = flexibleengine_dns_zone_v2.example_zone.id
  name        = "rs.example.com."
  description = "An example record set"
  ttl         = 3000
  type        = "A"
  records     = ["10.0.0.1"]
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the DNS record set.
  If omitted, the `region` argument of the provider is used.
  Changing this creates a new DNS record set.

* `zone_id` - (Required, String, ForceNew) The ID of the zone in which to create the record set.
  Changing this creates a new DNS record set.

* `name` - (Required, String, ForceNew) The name of the record set. Note the `.` at the end of the name.
  Changing this creates a new DNS record set.

* `type` - (Required, String, ForceNew) The type of record set. The options include `A`, `AAAA`, `MX`,
  `CNAME`, `TXT`, `NS`, `SRV`, `CAA`, and `PTR`.
  Changing this creates a new DNS record set.

* `records` - (Required, List) An array of DNS records.

* `ttl` - (Optional, Int) The time to live (TTL) of the record set (in seconds). The value
  range is 300–2147483647. The default value is `300`.

* `description` - (Optional, String) A description of the record set. Max length is `255` characters.

* `tags` - (Optional, Map) The key/value pairs to associate with the record set.

* `value_specs` - (Optional, ForceNew) Map of additional options.
  Changing this creates a new DNS record set.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.
* `update` - Default is 10 minutes.
* `delete` - Default is 10 minutes.

## Import

This resource can be imported by specifying the zone ID (`zone_id`) and recordset ID (`id`),
separated by a forward slash `/`.

```shell
terraform import flexibleengine_dns_recordset_v2.recordset_1 <zone_id>/<recordset_id>
```
