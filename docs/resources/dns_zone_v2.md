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

* `region` - (Optional) The region in which to create the DNS zone.
    If omitted, the `region` argument of the provider is used.
    Changing this creates a new DNS zone.

* `name` - (Required) The name of the zone. Note the `.` at the end of the name.
  Changing this creates a new DNS zone.

* `email` - (Optional) The email contact for the zone record.

* `zone_type` - (Optional) The type of zone. Can either be `public` or `private`.
  Changing this creates a new DNS zone.

* `router` - (Optional) Router configuration block which is required if zone_type is private.
  The router structure is documented below.

* `ttl` - (Optional) The time to live (TTL) of the zone.

* `description` - (Optional) A description of the zone.

* `tags` - (Optional, Map) The key/value pairs to associate with the zone.

* `value_specs` - (Optional) Map of additional options. Changing this creates a
  new DNS zone.

The `router` block supports:

* `router_id` - (Required) The VPC UUID.

* `router_region` - (Optional) The region of the VPC. Defaults to the `region`.

## Attributes Reference

The following attributes are exported:

* `region` - See Argument Reference above.
* `name` - See Argument Reference above.
* `email` - See Argument Reference above.
* `zone_type` - See Argument Reference above.
* `ttl` - See Argument Reference above.
* `description` - See Argument Reference above.
* `masters` - An array of master DNS servers.
* `value_specs` - See Argument Reference above.

## Import

This resource can be imported by specifying the zone ID:

```
$ terraform import flexibleengine_dns_zone_v2.zone_1 <zone_id>
```
