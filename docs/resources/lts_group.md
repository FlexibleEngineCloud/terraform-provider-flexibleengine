---
subcategory: "Log Tank Service (LTS)"
---

# flexibleengine_lts_group

Manages a log group resource within FlexibleEngine.

## Example Usage

```hcl
resource "flexibleengine_lts_group" "group_1" {
  group_name = "log_group1"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the log group resource.
  If omitted, the provider-level region will be used. Changing this creates a new log group resource.

* `group_name` - (Required, String, ForceNew) Specifies the log group name.
  Changing this parameter will create a new resource.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The log group ID in UUID format.

* `ttl_in_days` - Indicates the log expiration time. The value is fixed to 7 days.

## Import

Log group can be imported using the `id`, e.g.

```sh
terraform import flexibleengine_lts_group.group_1 6e728c21-e3b6-11eb-b081-286ed488cb76
```
