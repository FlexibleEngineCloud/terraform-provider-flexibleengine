---
layout: "flexibleengine"
page_title: "FlexibleEngine: flexibleengine_identity_role_v3"
sidebar_current: "docs-flexibleengine-resource-identity-role-v3"
description: |-
  custom role management in FlexibleEngine
---

# flexibleengine\_identity\_role\_v3

custom role management in FlexibleEngine

## Example Usage

### Role

```hcl
resource "flexibleengine_identity_role_v3" "role" {
  name        = "custom_role"
  description = "a custom role"
  scope       = "domain"

  policy {
    effect = "Allow"
    action = ["ecs:*:list*"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` -
  (Required)
  Specify the name of a role. The value cannot exceed 64 characters.

* `description` -
  (Required)
  Specify the description of a role. The value cannot exceed 256 characters.

* `scope` -
  (Required)
  Specify the scope layer of a role. The value supports:
  - domain - A role is displayed at the domain layer.
  - project - A role is displayed at the project layer.

* `policy` -
  (Required)
  The policy field contains the `effect` and `action` elements.
  Effect indicates whether the policy allows or denies access.
  Action indicates authorization items. The number of policy
  cannot exceed 8. Structure is documented below.

The `policy` block supports:

* `action` -
  (Required)
  Permission set, which specifies the operation permissions on
  resources. The number of permission sets cannot exceed 100.
  Format:  The value format is Service name:Resource type:Action,
  for example, vpc:ports:create.  Service name: indicates the
  product name, such as ecs, evs, or vpc. Only lowercase letters
  are allowed.  Resource type and Action: The values are
  case-insensitive, and the wildcard (*) are allowed. A wildcard
  (*) can represent all or part of information about resource
  types and actions for the specific service.

* `effect` -
  (Required)
  The value can be Allow and Deny. If both Allow and Deny are
  found in statements, the policy evaluation starts with Deny.

- - -

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:

* `catalog` -
  Directory where a role locates

* `domain_id` -
  ID of the domain to which a role belongs

## Import

Role can be imported using the following format:

```
$ terraform import flexibleengine_identity_role_v3.default {{ resource id}}
```
