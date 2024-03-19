---
subcategory: "Object Storage Service (OBS)"
description: ""
page_title: "flexibleengine_obs_bucket"
---

# flexibleengine_obs_bucket

Manages an OBS bucket resource within FlexibleEngine.

## Example Usage

### Basic Bucket

```hcl
resource "flexibleengine_obs_bucket" "b" {
  bucket = "my-tf-test-bucket"
  acl    = "private"
  storage_class = "STANDARD"
}
```

### Enable versioning

```hcl
resource "flexibleengine_obs_bucket" "b" {
  bucket     = "my-tf-test-bucket"
  acl        = "private"
  versioning = true
}
```

### Enable Logging

```hcl
resource "flexibleengine_obs_bucket" "log_bucket" {
  bucket = "my-tf-log-bucket"
  acl    = "log-delivery-write"
}

resource "flexibleengine_obs_bucket" "b" {
  bucket = "my-tf-test-bucket"
  acl    = "private"

  logging {
    target_bucket = flexibleengine_obs_bucket.log_bucket.id
    target_prefix = "log/"
  }
}
```

### Static Website Hosting

```hcl
resource "flexibleengine_obs_bucket" "b" {
  bucket = "obs-website-test.hashicorp.com"
  acl    = "public-read"

  website {
    index_document = "index.html"
    error_document = "error.html"

    routing_rules = <<EOF
[{
    "Condition": {
        "KeyPrefixEquals": "docs/"
    },
    "Redirect": {
        "ReplaceKeyPrefixWith": "documents/"
    }
}]
EOF
  }
}
```

### Using CORS

```hcl
resource "flexibleengine_obs_bucket" "b" {
  bucket = "obs-website-test.hashicorp.com"
  acl    = "public-read"

  cors_rule {
    allowed_origins = ["https://obs-website-test.hashicorp.com"]
    allowed_methods = ["PUT", "POST"]
    allowed_headers = ["*"]
    expose_headers  = ["ETag"]
    max_age_seconds = 3000
  }
}
```

### Using object lifecycle

```hcl
resource "flexibleengine_obs_bucket" "bucket" {
  bucket     = "my-bucket"
  acl        = "private"
  versioning = true

  lifecycle_rule {
    name    = "log"
    prefix  = "log/"
    enabled = true

    expiration {
      days = 365
    }
    transition {
      days = 60
      storage_class = "STANDARD_IA"
    }
    transition {
      days = 180
      storage_class = "GLACIER"
    }
  }

  lifecycle_rule {
    name    = "tmp"
    prefix  = "tmp/"
    enabled = true

    noncurrent_version_expiration {
      days = 180
    }
    noncurrent_version_transition {
      days = 30
      storage_class = "STANDARD_IA"
    }
    noncurrent_version_transition {
      days = 60
      storage_class = "GLACIER"
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the bucket resource.
  If omitted, the provider-level region will be used. Changing this will create a new bucket resource.

* `bucket` - (Required, String, ForceNew) Specifies the name of the bucket. Changing this parameter will create a new
  resource. A bucket must be named according to the globally applied DNS naming regulations as follows:
  + The name must be globally unique in OBS.
  + The name must contain 3 to 63 characters. Only lowercase letters, digits, hyphens (-), and periods (.) are allowed.
  + The name cannot start or end with a period (.) or hyphen (-), and cannot contain two consecutive periods (.) or
    contain a period (.) and a hyphen (-) adjacent to each other.
  + The name cannot be an IP address.
  + If the name contains any periods (.), a security certificate verification message may appear when you access
    the bucket or its objects by entering a domain name.

* `storage_class` - (Optional, String) Specifies the storage class of the bucket. OBS provides three storage classes:
  "STANDARD", "STANDARD_IA" (Infrequent Access) and "GLACIER" (Archive). Defaults to `STANDARD`.

* `acl` - (Optional, String) Specifies the ACL policy for a bucket. The predefined common policies are as follows:
  "private", "public-read", "public-read-write" and "log-delivery-write". Defaults to `private`.

* `versioning` - (Optional, Bool) Whether enable versioning. Once you version-enable a bucket,
  it can never return to an unversioned state. You can, however, suspend versioning on that bucket.

* `encryption` - (Optional, Bool) Whether enable default server-side encryption of the bucket in SSE-KMS mode.

* `kms_key_id` - (Optional, String) Specifies the ID of a KMS key. If omitted, the default master key will be used.

* `kms_key_project_id` - (Optional, String) Specifies the project ID to which the KMS key belongs. This field is valid
  only when `kms_key_id` is specified.

* `logging` - (Optional, List) A settings of bucket logging. The [logging](#obs_logging) object structure is documented
  below.

* `website` - (Optional, List) A website object. The [website](#obs_website) object structure is documented below.

* `cors_rule` - (Optional, List) A rule of Cross-Origin Resource Sharing. The [cors_rule](#obs_cors_rule) object
  structure is documented below.

* `lifecycle_rule` - (Optional, List) A configuration of object lifecycle management. The [lifecycle_rule](#obs_lifecycle_rule)
  object structure is documented below.

* `force_destroy` - (Optional, Bool) A boolean that indicates all objects should be deleted from the bucket so that
  the bucket can be destroyed without error. Default to `false`.

* `multi_az` - (Optional, Bool, ForceNew) Whether enable the multi-AZ mode for the bucket. When the multi-AZ mode is
  enabled, data in the bucket is duplicated and stored in multiple AZs. Changing this creates a new bucket.

* `parallel_fs` - (Optional, Bool, ForceNew) Whether enable a bucket as a parallel file system. Changing this will
  create a new bucket.

<a name="obs_logging"></a>
The `logging` object supports:

* `target_bucket` - (Required, String) The name of the bucket that will receive the log objects.
  The acl policy of the target bucket should be `log-delivery-write`.
* `target_prefix` - (Optional, String) To specify a key prefix for log objects.

<a name="obs_website"></a>
The `website` object supports:

* `index_document` - (Required, String) Specifies the default homepage of the static website, only HTML web pages are
  supported. It is **Optional** if using `redirect_all_requests_to`.
  OBS only allows files such as `index.html` in the root directory of a bucket to function as the default homepage.
  That is to say, do not set the default homepage with a multi-level directory structure (for example, /page/index.html).

* `error_document` - (Optional, String) Specifies the error page returned when an error occurs during static website
  access. Only HTML, JPG, PNG, BMP, and WEBP files under the root directory are supported.

* `redirect_all_requests_to` - (Optional, String) A hostname to redirect all website requests for this bucket to.
  Hostname can optionally be prefixed with a protocol (`http://` or `https://`) to use when redirecting requests.
  The default is the protocol that is used in the original request.

* `routing_rules` - (Optional, String) A JSON or XML format containing routing rules describing redirect behavior and
  when redirects are applied.
  Each rule contains a `Condition` and a `Redirect` as shown in the following table:

  Parameter | Key
  --- | ---
  Condition | KeyPrefixEquals, HttpErrorCodeReturnedEquals
  Redirect | Protocol, HostName, ReplaceKeyPrefixWith, ReplaceKeyWith, HttpRedirectCode

<a name="obs_cors_rule"></a>
The `cors_rule` object supports:

* `allowed_origins` (Required, List) Requests from this origin can access the bucket. Multiple matching rules are
  allowed. One rule occupies one line, and allows one wildcard character (*) at most.

* `allowed_methods` (Required, List) Specifies the acceptable operation type of buckets and objects.
  The methods include `GET`, `PUT`, `POST`, `DELETE` or `HEAD`.

* `allowed_headers` (Optional, List) Specifies the allowed header of cross-origin requests.
  Only CORS requests matching the allowed header are valid.

* `expose_headers` (Optional, List) Specifies the exposed header in CORS responses, providing additional information
  for clients.

* `max_age_seconds` (Optional, Int) Specifies the duration that your browser can cache CORS responses, expressed in
  seconds. The default value is 100.

<a name="obs_lifecycle_rule"></a>
The `lifecycle_rule` object supports:

* `name` - (Required, String) Unique identifier for lifecycle rules. The Rule Name contains a maximum of 255 characters.

* `enabled` - (Required, Bool) Specifies lifecycle rule status.

* `prefix` - (Optional, String) Object key prefix identifying one or more objects to which the rule applies.
  If omitted, all objects in the bucket will be managed by the lifecycle rule.
  The prefix cannot start or end with a slash (/), cannot have consecutive slashes (/),
  and cannot contain the following special characters: \:*?"<>|.

* `expiration` - (Optional, List) Specifies a period when objects that have been last updated are automatically deleted.
  The [expiration](#obs_expiration) object structure is documented below.

* `transition` - (Optional, List) Specifies a period when objects that have been last updated are automatically
  transitioned to `STANDARD_IA` or `GLACIER` storage class.
  The [transition](#obs_transition) object structure is documented below.

* `noncurrent_version_expiration` - (Optional, List) Specifies a period when noncurrent object versions are automatically
  deleted. The [noncurrent_version_expiration](#obs_noncurrent_version_expiration) object structure is documented below.

* `noncurrent_version_transition` - (Optional, List) Specifies a period when noncurrent object versions are automatically
  transitioned to `STANDARD_IA` or `GLACIER` storage class.
  The [noncurrent_version_transition](#obs_noncurrent_version_transition) object structure is documented below.

At least one of `expiration`, `transition`, `noncurrent_version_expiration`, `noncurrent_version_transition` must
be specified.

<a name="obs_expiration"></a>
The `expiration` object supports:

* `days` (Required, Int) Specifies the number of days when objects that have been last updated are automatically deleted.
  The expiration time must be greater than the transition times.

<a name="obs_transition"></a>
The `transition` object supports:

* `days` (Required, Int) Specifies the number of days when objects that have been last updated are automatically
  transitioned to the specified storage class.
* `storage_class` - (Required, String) The class of storage used to store the object. Only "STANDARD_IA" and "GLACIER"
  are supported.

<a name="obs_noncurrent_version_expiration"></a>
The `noncurrent_version_expiration` object supports:

* `days` (Required, Int) Specifies the number of days when noncurrent object versions are automatically deleted.

<a name="obs_noncurrent_version_transition"></a>
The `noncurrent_version_transition` object supports:

* `days` (Required, Int) Specifies the number of days when noncurrent object versions are automatically
  transitioned to the specified storage class.
* `storage_class` - (Required, String) The class of storage used to store the object. Only "STANDARD_IA" and "GLACIER"
  are supported.

## Attribute Reference

The following attributes are exported:

* `id` - The name of the bucket.

* `bucket_domain_name` - The bucket domain name. Will be of format `bucketname.oss.region.prod-cloud-ocb.orange-business.com`.

## Import

OBS bucket can be imported using the `bucket`, e.g.

```shell
terraform import flexibleengine_obs_bucket.bucket bucket-name
```

Note that the imported state may not be identical to your resource definition, due to some attributes
missing from the API response. The missing attributes include `acl` and `force_destroy`.
It is generally recommended running `terraform plan` after importing an OBS bucket.
Also, you can ignore changes as below.

```hcl
resource "flexibleengine_obs_bucket" "bucket" {
    ...

  lifecycle {
    ignore_changes = [
      acl, force_destroy,
    ]
  }
}
```
