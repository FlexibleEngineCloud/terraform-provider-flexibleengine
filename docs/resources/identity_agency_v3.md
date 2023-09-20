---
subcategory: "Identity and Access Management (IAM)"
description: ""
page_title: "flexibleengine_identity_agency_v3"
---

# flexibleengine_identity_agency_v3

Manages an agency resource within FlexibleEngine.

## Example Usage

### Delegate another account to perform operations on your resources

```hcl
resource "flexibleengine_identity_agency_v3" "agency" {
  name                  = "test_agency"
  description           = "this is a domain test agency"
  delegated_domain_name = "***"

  project_role {
    project = "eu-west-0"
    roles = [
      "KMS Administrator",
    ]
  }
  domain_roles = [
    "Anti-DDoS Administrator",
  ]
}
```

### Delegate a cloud service to access your resources in other cloud services

```hcl
resource "flexibleengine_identity_agency_v3" "agency" {
  name                   = "test_agency"
  description            = "this is a service test agency"
  delegated_service_name = "op_svc_evs"

  project_role {
    project = "eu-west-0"
    roles = [
      "Tenant Administrator",
    ]
  }
  domain_roles = [
    "OBS OperateAccess",
  ]
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, String) Specifies the name of agency. The name is a string of 1 to 64 characters.
    Changing this will create a new agency.

* `description` - (Optional, String) Specifies the supplementary information about the agency.
    The value is a string of 0 to 255 characters.

* `delegated_domain_name` - (Optional, String) Specifies the name of delegated user domain.
    This parameter and `delegated_service_name` are alternative.

* `delegated_service_name` - (Optional, String) Specifies the name of delegated cloud service.
    This parameter and `delegated_domain_name` are alternative.

* `duration` - (Optional, String) Specifies the validity period of an agency.
    The valid value are *ONEDAY* and *FOREVER*, defaults to *FOREVER*.

* `project_role` - (Optional, List) Specifies an array of one or more roles and projects which are used to grant
    permissions to agency on project. The [project_role](#identity_project_role) object structure is documented below.

* `domain_roles` - (Optional, List) Specifies an array of one or more role names which stand for the permissionis to
    be granted to agency on domain.

<a name="identity_project_role"></a>
The `project_role` block supports:

* `project` - (Required, String) Specifies the name of project.

* `roles` - (Required, List) Specifies an array of role names.

-> **NOTE**
    - At least one of `project_role` and `domain_roles` must be specified when creating an agency.
    - We can get all **System-Defined Roles** from
[FlexibleEngine](https://docs.prod-cloud-ocb.orange-business.com/permissions/index.html) or
[data.flexibleengine_identity_role_v3](https://registry.terraform.io/providers/FlexibleEngineCloud/flexibleengine/latest/docs/data-sources/identity_role_v3).

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The agency ID.
* `expire_time` - The expiration time of agency.
* `create_time` - The time when the agency was created.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.
* `update` - Default is 10 minutes.
* `delete` - Default is 5 minutes.

## Import

Agencies can be imported using the `id`, e.g.

```shell
terraform import flexibleengine_identity_agency_v3.agency 0b97661f9900f23f4fc2c00971ea4dc0
```
