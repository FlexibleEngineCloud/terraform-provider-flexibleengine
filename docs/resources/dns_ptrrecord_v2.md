---
subcategory: "Domain Name Service (DNS)"
description: ""
page_title: "flexibleengine_dns_ptrrecord_v2"
---

# flexibleengine_dns_ptrrecord_v2

Manages a DNS PTR record in the FlexibleEngine DNS Service.

## Example Usage

```hcl
resource "flexibleengine_vpc_eip" "eip_1" {
  publicip {
    type = "5_bgp"
  }
  bandwidth {
    name        = "test"
    size        = 5
    share_type  = "PER"
    charge_mode = "traffic"
  }
}

resource "flexibleengine_dns_ptrrecord_v2" "ptr_1" {
  name          = "ptr.example.com."
  description   = "An example PTR record"
  floatingip_id = flexibleengine_vpc_eip.eip_1.id
  ttl           = 3000

  tags = {
    foo = "bar"
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional) The region in which to create the PTR record.
    If omitted, the `region` argument of the provider is used.
    Changing this creates a new PTR record.

* `name` - (Required) Domain name of the PTR record. A domain name is case insensitive.
  Uppercase letters will also be converted into lowercase letters.

* `description` - (Optional) Description of the PTR record.

* `floatingip_id` - (Required) The ID of the FloatingIP/EIP.
  Changing this creates a new PTR record.

* `ttl` - (Optional) The time to live (TTL) of the record set (in seconds). The value
  range is 300–2147483647. The default value is 300.

* `tags` - (Optional) Tags key/value pairs to associate with the PTR record.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` -  The PTR record ID, which is in {region}:{floatingip_id} format.

* `address` - The address of the FloatingIP/EIP.

## Import

PTR records can be imported using region and floatingip/eip ID, separated by a colon(:), e.g.

```
$ terraform import flexibleengine_dns_ptrrecord_v2.ptr_1 eu-west-0:d90ce693-5ccf-4136-a0ed-152ce412b6b9
```
