---
subcategory: "Software Repository for Container (SWR)"
description: ""
page_title: "flexibleengine_swr_organization"
---

# flexibleengine_swr_organization

Manages an SWR organization resource within FlexibleEngine.

## Example Usage

```hcl
resource "flexibleengine_swr_organization" "test" {
  name = "org-test"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource. If omitted, the
  provider-level region will be used. Changing this creates a new resource.

* `name` - (Required, String, ForceNew) Specifies the name of the organization. The organization name must be globally
  unique. Changing this creates a new resource.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the organization.

* `creator` - The creator user name of the organization.

* `permission` - The permission of the organization, the value can be Manage, Write, and Read.

* `login_server` - The URL that can be used to log into the container registry.

## Import

Organizations can be imported using the `name`, e.g.

```
$ terraform import flexibleengine_swr_organization.test org-name
```
