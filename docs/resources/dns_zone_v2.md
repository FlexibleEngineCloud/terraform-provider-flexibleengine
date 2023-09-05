---
subcategory: "Domain Name Service (DNS)"
description: ""
page_title: "flexibleengine_dns_zone_v2"
---

# flexibleengine_dns_zone_v2

Manages a DNS zone in the FlexibleEngine DNS Service.

## Example Usage

### Create a public DNS zone

```hcl
resource "flexibleengine_dns_zone_v2" "my_public_zone" {
  name        = "example.com."
  email       = "jdoe@example.com"
  description = "my public zone"
  ttl         = 3000
}
```

### Create a private DNS zone

```hcl
resource "flexibleengine_dns_zone_v2" "my_private_zone" {
  name        = "1.example.com."
  email       = "jdoe@example.com"
  description = "my private zone"
  ttl         = 3000
  zone_type   = "private"

  router {
    router_id     = "2c1fe4bd-ebad-44ca-ae9d-e94e63847b75"
    router_region = "eu-west-0"
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the DNS zone.
  If omitted, the `region` argument of the provider is used.
  Changing this creates a new DNS zone.

* `name` - (Required, String, ForceNew) The name of the zone. Note the `.` at the end of the name.
  Changing this creates a new DNS zone.

* `email` - (Optional, String) The email contact for the zone record.

* `zone_type` - (Optional, ForceNew) The type of zone. Can either be `public` or `private`.
  Default is `public`. Changing this creates a new DNS zone.

* `router` - (Optional, List) Router configuration block which is required if zone_type is private.
  The router structure is documented below.

* `ttl` - (Optional, Int) The time to live (TTL) of the zone. TTL ranges from 1 to 2147483647 seconds.
  Default is  `300`.

* `description` - (Optional, String) A description of the zone. Max length is `255` characters.

* `tags` - (Optional, Map) The key/value pairs to associate with the zone.

* `value_specs` - (Optional, ForceNew) Map of additional options.
  Changing this creates a new DNS zone.

The `router` block supports:

* `router_id` - (Required, String) The VPC UUID.

* `router_region` - (Optional, String) The region of the VPC. Defaults to the `region`.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Specifies a resource ID in UUID format.

* `masters` - An array of master DNS servers.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.
* `update` - Default is 10 minutes.
* `delete` - Default is 10 minutes.

## Import

This resource can be imported by specifying the zone ID:

```shell
terraform import flexibleengine_dns_zone_v2.zone_1 <zone_id>
```
