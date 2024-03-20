---
subcategory: "Identity and Access Management (IAM)"
---

# flexibleengine_identity_group_role_assignment

Manages an IAM user group role assignment within FlexibleEngine IAM Service.
This is an alternative to `flexibleengine_identity_role_assignment_v3`

-> **NOTE:** 1. You *must* have admin privileges to use this resource.
  <br/>2. When the resource is created, the permissions will take effect after 15 to 30 minutes.

## Example Usage

### Assign role with project

```hcl
variable "project_id" {}

data "flexibleengine_identity_role_v3" "test" {
  # RDS Administrator
  name = "rds_adm"
}

resource "flexibleengine_identity_group_v3" "test" {
  name = "group_1"
}

resource "flexibleengine_identity_group_role_assignment" "test" {
  group_id   = flexibleengine_identity_group_v3.test.id
  role_id    = data.flexibleengine_identity_role_v3.test.id
  project_id = var.project_id
}
```

### Assign role with all projects

```hcl
data "flexibleengine_identity_role_v3" "test" {
  # RDS Administrator
  name = "rds_adm"
}

resource "flexibleengine_identity_group_v3" "test" {
  name = "group_1"
}

resource "flexibleengine_identity_group_role_assignment" "all" {
  group_id   = flexibleengine_identity_group_v3.test.id
  role_id    = data.flexibleengine_identity_role_v3.test.id
  project_id = "all"
}
```

### Assign role with domain

```hcl
variable "domain_id" {}

data "flexibleengine_identity_role_v3" "test" {
  # OBS Administrator
  name = "obs_adm"
}

resource "flexibleengine_identity_group_v3" "test" {
  name = "group_1"
}

resource "flexibleengine_identity_group_role_assignment" "test" {
  group_id  = flexibleengine_identity_group_v3.test.id
  role_id   = data.flexibleengine_identity_role_v3.test.id
  domain_id = var.domain_id
}
```

### Assign role with enterprise project

```hcl
variable "enterprise_project_id" {}

data "flexibleengine_identity_role_v3" "test" {
  # RDS Administrator
  name = "rds_adm"
}

resource "flexibleengine_identity_group_v3" "test" {
  name = "group_1"
}

resource "flexibleengine_identity_group_role_assignment" "test" {
  group_id              = flexibleengine_identity_group_v3.test.id
  role_id               = data.flexibleengine_identity_role_v3.test.id
  enterprise_project_id = var.enterprise_project_id
}
```

## Argument Reference

The following arguments are supported:

* `group_id` - (Required, String, ForceNew) Specifies the group to assign the role to.
  Changing this parameter will create a new resource.

* `role_id` - (Required, String, ForceNew) Specifies the role to assign.
  Changing this parameter will create a new resource.

* `domain_id` - (Optional, String, ForceNew) Specifies the domain to assign the role in.
  Changing this parameter will create a new resource.

* `project_id` - (Optional, String, ForceNew) Specifies the project to assign the role in.
  If `project_id` is set to **all**, it means that the specified user group will be able to use all projects,
  including existing and future projects.

  Changing this parameter will create a new resource.

* `enterprise_project_id` - (Optional, String, ForceNew) Specifies the enterprise project to assign the role in.
  Changing this parameter will create a new resource.

  ~> Exactly one of `domain_id`, `project_id` or `enterprise_project_id` must be specified.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID. When assign in domain, the format is `<group_id>/<role_id>/<domain_id>`;
  when assign in project, the format is `<group_id>/<role_id>/<project_id>`;
  when assign in enterprise project, the format is `<group_id>/<role_id>/<enterprise_project_id>`;

## Import

The role assignments can be imported using the `group_id`, `role_id` and  `domain_id`, `project_id`,
  `enterprise_project_id`, e.g.

```bash
$ terraform import flexibleengine_identity_group_role_assignment.test <group_id>/<role_id>/<domain_id>
```

or

```bash
$ terraform import flexibleengine_identity_group_role_assignment.test <group_id>/<role_id>/<project_id>
```

or

```bash
$ terraform import flexibleengine_identity_group_role_assignment.test <group_id>/<role_id>/all
```

or

```bash
$ terraform import flexibleengine_identity_group_role_assignment.test <group_id>/<role_id>/<enterprise_project_id>
```
