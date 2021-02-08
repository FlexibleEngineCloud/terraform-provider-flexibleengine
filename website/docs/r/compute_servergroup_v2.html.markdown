---
layout: "flexibleengine"
page_title: "FlexibleEngine: flexibleengine_compute_servergroup_v2"
sidebar_current: "docs-flexibleengine-resource-compute-servergroup-v2"
description: |-
  Manages a V2 Server Group resource within FlexibleEngine.
---

# flexibleengine\_compute\_servergroup_v2

Manages a V2 Server Group resource within FlexibleEngine.

## Example Usage

```hcl
resource "flexibleengine_compute_servergroup_v2" "test-sg" {
  name     = "my-sg"
  policies = ["anti-affinity"]
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional) The region in which to obtain the V2 Compute client.
    If omitted, the `region` argument of the provider is used. Changing
    this creates a new server group.

* `name` - (Required) A unique name for the server group. Changing this creates
    a new server group.

* `policies` - (Required) The set of policies for the server group. Only **anti-affinity**
    policy is supported right now, which menas all servers in this group must be
    deployed on different hosts. Changing this creates a new server group.

* `value_specs` - (Optional) Map of additional options.

## Attributes Reference

The following attributes are exported:

* `region` - See Argument Reference above.
* `name` - See Argument Reference above.
* `policies` - See Argument Reference above.
* `members` - The instances that are part of this server group.

## Import

Server Groups can be imported using the `id`, e.g.

```
$ terraform import flexibleengine_compute_servergroup_v2.test-sg 1bc30ee9-9d5b-4c30-bdd5-7f1e663f5edf
```
