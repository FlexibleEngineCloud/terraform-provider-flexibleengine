---
subcategory: "VPC Endpoint (VPCEP)"
---

# flexibleengine_vpcep_endpoints

Use this data source to get VPC endpoints.

## Example Usage

```hcl
data "flexibleengine_vpcep_endpoints" "all_endpoints" {
}

data "flexibleengine_vpcep_endpoints" "dns_endpoints" {
  endpoint_service_name = "dns"
}
```

## Argument Reference

* `service_name` - (Optional, String) Specifies the name of the VPC endpoint service.
    The value is not case-sensitive and supports fuzzy match.

* `endpoint_id` - (Optional, String) Specifies the unique ID of the VPC endpoint.

* `vpc_id` - (Optional, String) Specifies the unique ID of the vpc holding the VPC endpoint service.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Specifies a data source ID in UUID format.

* `endpoints` - Indicates the public VPC endpoints information. Structure is documented below.

The `endpoints` block contains:

* `id` - The unique ID of the public VPC endpoint service.
* `status` - The connection status of the VPC endpoint.
* `service_id` - The ID of the VPC endpoint service.
* `service_name` - The name of the VPC endpoint service.
* `service_type` - The type of the VPC endpoint.
* `vpc_id` - The ID of the VPC holding the VPC endpoint service.
* `network_id` - The ID of the subnet holding the VPC endpoint.
* `ip_address` - The IP of the VPC endpoint.
* `packet_id` - The marker id of the VPC endpoint.
* `enable_dns` - Flag indicating dns has been enabled for the VPC endpoint.
* `enable_whitelist` - Flag indicating access control have been enabled on this VPC endpoint.
* `whitelist` - List of IP or CIDR block which can access the VPC endpoint.
* `private_domain_name` - DNS name pointing to the VPC endpoint ip.
* `tags` - The key/value pairs to associate with the VPC endpoint.
    + `key` - The tag key. Each tag key contains a maximum of 127 unicode characters but cannot be left blank.
    + `value` - The tag value list. Each value contains a maximum of 255 Unicode characters.
      Before using values, delete SBC spaces before and after the value.
* `project_id` - The ID of the project holding the VPC endpoint.
* `created_at` - Creation date of the VPC endpoint.
* `updated_at` - Last update date of the VPC endpoint.
