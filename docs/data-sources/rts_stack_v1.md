---
subcategory: "Resource Template Service (RTS)"
---

# flexibleengine_rts_stack_v1

The FlexibleEngine RTS Stack data source allows access to stack outputs and other useful data including the template body.

## Example Usage

```hcl
variable "stack_name" {}

data "flexibleengine_rts_stack_v1" "mystack" {
  name = var.stack_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) The region in which to query the data source. If omitted, the provider-level region
  will be used.

* `name` - (Required, String) The name of the stack.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - A unique identifier of the stack.

* `capabilities` - List of stack capabilities for stack.

* `notification_topics` - List of notification topics for stack.

* `status` - Specifies the stack status.

* `disable_rollback` - Whether the rollback of the stack is disabled when stack creation fails.

* `outputs` - A list of stack outputs.

* `parameters` - A map of parameters that specify input parameters for the stack.

* `template_body` - Structure containing the template body.

* `timeout_mins` - Specifies the timeout duration.
