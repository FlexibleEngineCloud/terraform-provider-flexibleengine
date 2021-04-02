---
subcategory: "Elastic Load Balance (ELB)"
---

# flexibleengine\_lb\_certificate\_v2

flexibleengine_lb_certificate_v2 provides details about a specific Certificate.

## Example Usage

The following example shows how one might accept a certificate name as a variable to fetch this data source.

```hcl

variable "cert_name" {}

data "flexibleengine_lb_certificate_v2" "by_name" {
  name = "${var.cert_name}"
}

```

## Argument Reference

The arguments of this data source act as filters for querying the available Certificates in the current region. The given filters must match exactly one Certificate whose data will be exported as attributes.

* `id` - (Optional) The id of the specific Certificate to retrieve.

* `name` - (Optional) Human-readable name for the Certificate. Does not have
    to be unique.

* `description` - (Optional) Human-readable description for the Certificate.

* `domain` - (Optional) The domain of the Certificate.


## Attributes Reference

The following attributes are exported:

* `name` - See Argument Reference above.
* `description` - See Argument Reference above.
* `domain` - See Argument Reference above.
* `private_key` - The private encrypted key of the Certificate, PEM format.
* `certificate` - The public encrypted key of the Certificate, PEM format.
* `update_time` - Indicates the update time.
* `create_time` - Indicates the creation time.
