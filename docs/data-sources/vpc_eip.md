---
subcategory: "Elastic IP (EIP)"
---

# flexibleengine_vpc_eip

Use this data source to get the details of an available EIP.

## Example Usage

```hcl
data "flexibleengine_vpc_eip" "by_address" {
  public_ip = "123.60.208.163"
}
```

## Argument Reference

* `public_ip` - (Optional, String) The public ip address of the EIP.

* `port_id` - (Optional, String) The port id of the EIP.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `status` - The status of the EIP.

* `type` - The type of the EIP.

* `private_ip` - The private ip of the EIP.

* `bandwidth_id` - The bandwidth id of the EIP.

* `bandwidth_size` - The bandwidth size of the EIP.

* `bandwidth_share_type` - The bandwidth share type of the EIP.
