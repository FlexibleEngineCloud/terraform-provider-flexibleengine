---
subcategory: "FunctionGraph"
description: ""
page_title: "flexibleengine_fgs_dependency"
---

# flexibleengine_fgs_dependency

Manages a custom dependency package within FlexibleEngine FunctionGraph.

## Example Usage

### Create a custom dependency package using a OBS bucket path where the zip file is located

```hcl
variable "package_name"
variable "package_location"
variable "dependency_name"

resource "flexibleengine_obs_bucket" "test" {
  ...
}

resource "flexibleengine_obs_bucket_object" "test" {
  bucket = flexibleengine_obs_bucket.test.bucket
  key    = format("terraform_dependencies/%s", var.package_name)
  source = var.package_location
}

resource "flexibleengine_fgs_dependency" "test" {
  name    = var.dependency_name
  runtime = "Python3.6"
  link    = format("https://%s/%s", flexibleengine_obs_bucket.test.bucket_domain_name, flexibleengine_obs_bucket_object.test.key)
}
```

## Argument Reference

* `region` - (Optional, String, ForceNew) Specifies the region in which to create a custom dependency package.
  If omitted, the provider-level region will be used. Changing this will create a new resource.

* `runtime` - (Required, String) Specifies the dependency package runtime.
  The valid values are **Java8**, **Node.js6.10**, **Node.js8.10**, **Node.js10.16**, **Node.js12.13**, **Python2.7**,
  **Python3.6**, **Go1.8**, **Go1.x**, **C#(.NET Core 2.0)**, **C#(.NET Core 2.1)**, **C#(.NET Core 3.1)** and
  **PHP7.3**.

* `name` - (Required, String) Specifies the dependeny name.
  The name can contain a maximum of 96 characters and must start with a letter and end with a letter or digit.
  Only letters, digits, underscores (_), periods (.), and hyphens (-) are allowed.

* `link` - (Required, String) Specifies the OBS bucket path where the dependency package is located. The OBS object URL
  must be in zip format, such as 'https://obs-terraform.oss.eu-west-0.prod-cloud-ocb.orange-business.com/dependencies/sdkcore.zip'.

-> A link can only be used to create at most one dependency package.

* `description` - (Optional, String) Specifies the dependency description.
  The description can contain a maximum of 512 characters.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The dependency ID in UUID format.

* `owner` - The base64 encoded digest of the dependency after encryption by MD5.

* `etag` - The unique ID of the dependency package.

* `size` - The dependency package size in bytes.

## Import

Dependencies can be imported using the `id`, e.g.:

```
$ terraform import flexibleengine_fgs_dependency.test 795e722f-0c23-41b6-a189-dcd56f889cf6
```
