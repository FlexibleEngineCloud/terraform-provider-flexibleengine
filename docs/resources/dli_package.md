---
subcategory: "Data Lake Insight (DLI)"
---

# flexibleengine_dli_package

Manages DLI package resource within Flexibleengine

## Example Usage

### Upload the specified python script as a resource package

```hcl
variable "group_name" {}
variable "access_domain_name" {}

resource "flexibleengine_dli_package" "queue" {
  group_name  = var.group_name
  object_path = "https://${var.access_domain_name}/dli/packages/object_file.py"
  type        = "pyFile"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to upload packages.
  If omitted, the provider-level region will be used.
  Changing this parameter will delete the current package and upload a new package.

* `group_name` - (Required, String, ForceNew) Specifies the group name which the package belongs to.
  Changing this parameter will delete the current package and upload a new package.

* `type` - (Required, String, ForceNew) Specifies the package type.
  + **jar**: `.jar` or jar related files.
  + **pyFile**: `.py` or python related files.
  + **file**: Other user files.

  Changing this parameter will delete the current package and upload a new package.

* `object_path` - (Required, String, ForceNew) Specifies the OBS storage path where the package is located.
  For example, `https://{bucket_name}.oss.{region}.prod-cloud-ocb.orange-business.com/dli/packages/object_file.py`.
  Changing this parameter will delete the current package and upload a new package.

* `is_async` - (Optional, Bool, ForceNew) Specifies whether to upload resource packages in asynchronous mode.
  The default value is **false**. Changing this parameter will delete the current package and upload a new package.

* `owner` - (Optional, String) Specifies the name of the package owner. The owner must be IAM user.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Resource ID. The ID is constructed from the `group_name` and `object_name`, separated by slash.

* `object_name` - The package name.

* `status` - Status of a package group to be uploaded.

* `created_at` - Time when a queue is created.

* `updated_at` - The last time when the package configuration update has complated.
