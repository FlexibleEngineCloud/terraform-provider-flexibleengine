---
subcategory: "Resource Template Service (RTS)"
---

# flexibleengine_rts_stack_resource_v1

The FlexibleEngine RTS Stack Resource data source allows access to stack resource metadata.

## Example Usage

```hcl
variable "stack_name" { }
variable "resource_name" { }

data "flexibleengine_rts_stack_resource_v1" "stackresource" {
  stack_name    = var.stack_name
  resource_name = var.resource_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) The region in which to query the data source. If omitted, the provider-level region
  will be used.

* `stack_name` - (Required, String) The unique stack name.

* `resource_name` - (Optional, String) The name of a resource in the stack.

* `physical_resource_id` - (Optional, String) The physical resource ID.

* `resource_type` - (Optional, String) The resource type.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `logical_resource_id` - The logical resource ID.

* `resource_status` - The status of the resource.

* `resource_status_reason` - The resource operation reason.

* `required_by` - Specifies the resource dependency.
