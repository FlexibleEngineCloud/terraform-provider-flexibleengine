---
subcategory: "Elastic Load Balance (ELB)"
---

# flexibleengine_lb_certificate

Use this data source to get the certificate details in FlexibleEngine Elastic Load Balance (ELB).

## Example Usage

The following example shows how one might accept a certificate name as a variable to fetch this data source.

```hcl
variable "cert_name" {}

data "flexibleengine_lb_certificate" "by_name" {
  name = var.cert_name
}
```

## Argument Reference

The arguments of this data source act as filters for querying the available Certificates in the current region.
The given filters must match exactly one Certificate whose data will be exported as attributes.

* `id` - (Optional, String) The id of the specific Certificate to retrieve.

* `name` - (Optional, String) Human-readable name for the Certificate. Does not have to be unique.

* `description` - (Optional, String) Human-readable description for the LB Certificate.

* `domain` - (Optional, String) The domain of the Certificate.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `private_key` - The private encrypted key of the Certificate, PEM format.

* `certificate` - The public encrypted key of the Certificate, PEM format.

* `update_time` - Indicates the update time.

* `create_time` - Indicates the creation time.
