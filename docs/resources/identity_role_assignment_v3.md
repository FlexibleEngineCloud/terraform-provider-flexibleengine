---
subcategory: "Identity and Access Management (IAM)"
---

# flexibleengine_identity_role_assignment_v3

Manages a V3 Role assignment within group on FlexibleEngine IAM Service.

-> You *must* have admin privileges in your FlexibleEngine cloud to use this resource.

## Example Usage: Assign Role On Project Level

```hcl
data "flexibleengine_identity_project_v3" "project_1" {
  name = "eu-west-0_project_1"
}

data "flexibleengine_identity_role_v3" "role_1" {
  name = "system_all_1"
}

resource "flexibleengine_identity_group_v3" "group_1" {
  name = "group_1"
}

resource "flexibleengine_identity_role_assignment_v3" "role_assignment_1" {
  group_id   = flexibleengine_identity_group_v3.group_1.id
  project_id = data.flexibleengine_identity_project_v3.project_1.id
  role_id    = data.flexibleengine_identity_role_v3.role_1.id
}
```

## Example Usage: Assign Role On Domain Level

```hcl

variable "domain_id" {
    default = "01aafcf63744d988ebef2b1e04c5c34"
    description = "this is the domain id"
}

resource "flexibleengine_identity_group_v3" "group_1" {
  name = "group_1"
}

data "flexibleengine_identity_role_v3" "role_1" {
  name = "secu_admin" #security admin
}

resource "flexibleengine_identity_role_assignment_v3" "role_assignment_1" {
  group_id  = flexibleengine_identity_group_v3.group_1.id
  domain_id = var.domain_id
  role_id   = data.flexibleengine_identity_role_v3.role_1.id
} 

```

## Argument Reference

The following arguments are supported:

* `role_id` - (Required) The role to assign.

* `group_id` - (Required) The group to assign the role in.

* `domain_id` - (Optional; Required if `project_id` is empty) The domain to assign the role in.

* `project_id` - (Optional; Required if `domain_id` is empty) The project to assign the role in.

## Attributes Reference

The following attributes are exported:

* `group_id` - See Argument Reference above.
* `role_id` - See Argument Reference above.
* `domain_id` - See Argument Reference above.
* `project_id` - See Argument Reference above.
