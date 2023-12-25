---
subcategory: "API Gateway (Dedicated APIG)"
description: ""
page_title: "flexibleengine_apig_vpc_channel"
---

# flexibleengine_apig_vpc_channel

Manages a VPC channel resource within FlexibleEngine.

## Example Usage

```hcl
variable "instance_id" {}
variable "channel_name" {}
variable "ecs_id1" {}
variable "ecs_id2" {}

resource "flexibleengine_apig_vpc_channel" "test" {
  instance_id = var.instance_id
  name        = var.app_name
  port        = 8080
  protocol    = "HTTPS"
  path        = "/"
  http_code   = "201,202,203"

  members {
    id     = var.ecs_id1
    weight = 30
  }

  members {
    id     = var.ecs_id2
    weight = 70
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the VPC channel resource.
  If omitted, the provider-level region will be used.
  Changing this will create a new VPC channel resource.

* `instance_id` - (Required, String, ForceNew) Specifies an ID of the APIG dedicated instance to which the APIG
  vpc channel belongs to.
  Changing this will create a new VPC channel resource.

* `name` - (Required, String) Specifies the name of the VPC channel.
  The channel name consists of 3 to 64 characters, starting with a letter.
  Only letters, digits and underscores (_) are allowed.
  Chinese characters must be in UTF-8 or Unicode format.

* `port` - (Required, Int) Specifies the host port of the VPC channel.
  The valid value is range from 1 to 65535.

* `members` - (Required, List) Specifies an array of one or more backend server IDs or IP addresses that bind the VPC
  channel. The object structure is documented below.

* `member_type` - (Optional, String) Specifies the type of the backend service.
  The valid types are *ECS* and *EIP*, default to *ECS*.

* `algorithm` - (Optional, String) Specifies the type of the backend service.
  The valid types are *WRR*, *WLC*, *SH* and *URI hashing*, default to *WRR*.

* `protocol` - (Optional, String) Specifies the protocol for performing health checks on backend servers in the VPC
  channel.
  The valid values are *TCP*, *HTTP* and *HTTPS*, default to *TCP*.

* `path` - (Optional, String) Specifies the destination path for health checks.
  Required if `protocol` is *HTTP* or *HTTPS*.

* `healthy_threshold` - (Optional, Int) Specifies the healthy threshold, which refers to the number of consecutive
  successful checks required for a backend server to be considered healthy.
  The valid value is range from 2 to 10, default to 2.

* `unhealthy_threshold` - (Optional, Int) Specifies the unhealthy threshold, which refers to the number of consecutive
  failed checks required for a backend server to be considered unhealthy.
  The valid value is range from 2 to 10, default to 5.

* `timeout` - (Optional, Int) Specifies the timeout for determining whether a health check fails, in second.
  The value must be less than the value of time_interval.
  The valid value is range from 2 to 30, default to 5.

* `interval` - (Optional, Int) Specifies the interval between consecutive checks, in second.
  The valid value is range from 5 to 300, default to 10.

* `http_code` - (Optional, String) Specifies the response codes for determining a successful HTTP response.  
  The valid value ranges from `100` to `599` and the valid formats are as follows:
  + The multiple values, for example, **200,201,202**.
  + The range, for example, **200-299**.
  + Both multiple values and ranges, for example, **201,202,210-299**.

  It is Required if the `protocol` is **HTTP**.

The `members` block supports:

* `id` - (Optional, String) Specifies the ECS ID for each backend servers.
  Required if `member_type` is *ECS*.
  This parameter and `ip_address` are alternative.

* `ip_address` - (Optional, String) Specifies the IP address each backend servers.
  Required if `member_type` is *EIP*.

* `weight` - (Optional, Int) Specifies the backend server weight.
  The valid values are range from 1 to 100, default to 1.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the VPC channel.
* `create_at` - Time when the channel created, in UTC format.
* `status` - The status of VPC channel, supports *Normal* and *Abnormal*.

## Import

VPC Channels can be imported using the ID of the APIG dedicated instance to which the channel
belongs and `name`, separated by a slash, e.g.

```shell
terraform import flexibleengine_apig_vpc_channel.test <instance_id>/<name>
```
