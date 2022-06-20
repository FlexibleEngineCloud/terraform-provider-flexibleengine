---
subcategory: "Software Repository for Container (SWR)"
---

# flexibleengine_swr_repository

Manages an SWR repository resource within FlexibleEngine.

## Example Usage

```hcl
variable "organization_name" {} 

resource "flexibleengine_swr_repository" "test" {
  organization = var.organization_name
  name         = "%s"
  description  = "Test repository"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource. If omitted, the
  provider-level region will be used. Changing this creates a new resource.

* `organization` - (Required, String, ForceNew) Specifies the name of the organization (namespace) the repository belongs.
  Changing this creates a new resource.

* `name` - (Required, String, ForceNew) Specifies the name of the repository. Changing this creates a new resource.

* `is_public` - (Optional, Bool) Specifies whether the repository is public. Default is `false`.
  + `true` - Indicates the repository is *public*.
  + `false` - Indicates the repository is *private*.

* `description` - (Optional, String) Specifies the description of the repository.

* `category` - (Optional, String) Specifies the category of the repository.
  The value can be `app_server`, `linux`, `framework_app`, `database`, `lang`, `other`, `windows`, `arm`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the repository. The value is the name of the repository.

* `repository_id` - Numeric ID of the repository

* `path` - Image address for docker pull.

* `internal_path` - Intra-cluster image address for docker pull.

* `num_images` - Number of image tags in a repository.

* `size` - Repository size.

## Import

Repository can be imported using the organization name and repository name separated by a slash, e.g.:

```
$ terraform import flexibleengine_swr_repository.test org-name/repo-name
```
