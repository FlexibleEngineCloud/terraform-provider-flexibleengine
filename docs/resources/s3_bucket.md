---
subcategory: "Object Storage Service (OBS)"
description: ""
page_title: "flexibleengine_s3_bucket"
---

# flexibleengine_s3_bucket

Provides a S3 bucket resource.

## Example Usage

### Private Bucket w/ Tags

```hcl
resource "flexibleengine_s3_bucket" "b" {
  bucket = "my-tf-test-bucket"
  acl    = "private"
}
```

### Static Website Hosting

```hcl
resource "flexibleengine_s3_bucket" "b" {
  bucket = "s3-website-test.hashicorp.com"
  acl    = "public-read"
  policy = file("policy.json")

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
resource "flexibleengine_s3_bucket" "b" {
  bucket = "s3-website-test.hashicorp.com"
  acl    = "public-read"

  cors_rule {
    allowed_headers = ["*"]
    allowed_methods = ["PUT", "POST"]
    allowed_origins = ["https://s3-website-test.hashicorp.com"]
    expose_headers  = ["ETag"]
    max_age_seconds = 3000
  }
}
```

### Using versioning

```hcl
resource "flexibleengine_s3_bucket" "b" {
  bucket = "my-tf-test-bucket"
  acl    = "private"

  versioning {
    enabled = true
  }
}
```

### Enable Logging

```hcl
resource "flexibleengine_s3_bucket" "log_bucket" {
  bucket = "my-tf-log-bucket"
  acl    = "log-delivery-write"
}

resource "flexibleengine_s3_bucket" "b" {
  bucket = "my-tf-test-bucket"
  acl    = "private"

  logging {
    target_bucket = flexibleengine_s3_bucket.log_bucket.id
    target_prefix = "log/"
  }
}
```

### Using object lifecycle

```hcl
resource "flexibleengine_s3_bucket" "bucket" {
  bucket = "my-bucket"
  acl    = "private"

  lifecycle_rule {
    id      = "log"
    enabled = true

    prefix  = "log/"

    expiration {
      days = 90
    }
  }

  lifecycle_rule {
    id      = "tmp"
    prefix  = "tmp/"
    enabled = true

    expiration {
      date = "2016-01-12"
    }
  }
}

resource "flexibleengine_s3_bucket" "versioning_bucket" {
  bucket = "my-versioning-bucket"
  acl    = "private"

  versioning {
    enabled = true
  }

  lifecycle_rule {
    prefix  = "config/"
    enabled = true
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the s3 bucket resource.
  If omitted, the provider-level region will be used. Changing this will create a new s3 bucket resource.

* `bucket` - (Optional, String, ForceNew) The name of the bucket. If omitted, Terraform will assign a random,
  unique name. Changing this will create a new resource.

* `bucket_prefix` - (Optional, String, ForceNew) Creates a unique bucket name beginning with the specified prefix.
  Conflicts with `bucket`. Changing this will create a new resource.

* `acl` - (Optional, String) The [canned ACL](https://docs.aws.amazon.com/AmazonS3/latest/dev/acl-overview.html#canned-acl)
  to apply. Defaults to **private**.

* `policy` - (Optional, String) A valid [bucket policy](https://docs.aws.amazon.com/AmazonS3/latest/dev/example-bucket-policies.html)
  JSON document. Note that if the policy document is not specific enough (but still valid), Terraform may view the
  policy as constantly changing in a `terraform plan`. In this case, please make sure you use the verbose/specific
  version of the policy.

* `force_destroy` - (Optional, Bool) A boolean that indicates all objects should be deleted from the bucket
  so that the bucket can be destroyed without error. These objects are *not* recoverable. Default to **false**.

* `website` - (Optional, List) A website object.
  The [website](#obs_website) object structure is documented below.

* `cors_rule` - (Optional, List) A rule of [Cross-Origin Resource Sharing](https://docs.aws.amazon.com/AmazonS3/latest/dev/cors.html).
  The [cors_rule](#obs_cors_rule) object structure is documented below.

* `versioning` - (Optional, List) A state of [versioning](https://docs.aws.amazon.com/AmazonS3/latest/dev/Versioning.html).
  The [versioning](#obs_versioning) object structure is documented below.

* `logging` - (Optional, List) A settings of [bucket logging](https://docs.aws.amazon.com/AmazonS3/latest/UG/ManagingBucketLogging.html).
  The [logging](#obs_logging) object structure is documented below.

* `lifecycle_rule` - (Optional, List) A configuration of [object lifecycle management](http://docs.aws.amazon.com/AmazonS3/latest/dev/object-lifecycle-mgmt.html)
  The [lifecycle_rule](#obs_lifecycle_rule) object structure is documented below.

<a name="obs_website"></a>
The `website` object supports:

* `index_document` - (Required, String) Amazon S3 returns this index document when requests are made to the root domain
  or any of the subfolders. It is **Optional** if using `redirect_all_requests_to`.

* `error_document` - (Optional, String) An absolute path to the document to return in case of a 4XX error.

* `redirect_all_requests_to` - (Optional, String) A hostname to redirect all website requests for this bucket to.
  Hostname can optionally be prefixed with a protocol (`http://` or `https://`) to use when redirecting requests.
  The default is the protocol that is used in the original request.

* `routing_rules` - (Optional, String) A json array containing [routing rules](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-s3-websiteconfiguration-routingrules.html)
  describing redirect behavior and when redirects are applied.

<a name="obs_cors_rule"></a>
The `cors_rule` object supports:

* `allowed_headers` (Optional, List) Specifies which headers are allowed.

* `allowed_methods` (Required, List) Specifies which methods are allowed. Can be `GET`, `PUT`, `POST`, `DELETE`
  or `HEAD`.

* `allowed_origins` (Required, List) Specifies which origins are allowed.

* `expose_headers` (Optional, List) Specifies expose header in the response.

* `max_age_seconds` (Optional, Int) Specifies time in seconds that browser can cache the response for a preflight
  request.

<a name="obs_versioning"></a>
The `versioning` object supports:

* `enabled` - (Optional, Bool) Enable versioning. Once you version-enable a bucket, it can never return to an
  unversioned state. You can, however, suspend versioning on that bucket.

* `mfa_delete` - (Optional, Bool) Enable MFA delete for either change the versioning state of your bucket or
  permanently delete an object version. Default is `false`.

<a name="obs_logging"></a>
The `logging` object supports:

* `target_bucket` - (Required, String) The name of the bucket that will receive the log objects.

* `target_prefix` - (Optional, String) To specify a key prefix for log objects.

<a name="obs_lifecycle_rule"></a>
The `lifecycle_rule` object supports:

* `id` - (Optional, String) Unique identifier for the rule.

* `prefix` - (Optional, String) Object key prefix identifying one or more objects to which the rule applies.

* `enabled` - (Required, Bool) Specifies lifecycle rule status.

* `abort_incomplete_multipart_upload_days` (Optional, Int) Specifies the number of days after initiating a multipart
  upload when the multipart upload must be completed.

* `expiration` - (Optional, List) Specifies a period in the object's expire.
  The [expiration](#obs_expiration) object structure is documented below.

* `noncurrent_version_expiration` - (Optional, List) Specifies when noncurrent object versions expire.
  The [noncurrent_version_expiration](#obs_noncurrent_version_expiration) object structure is documented below.

At least one of `expiration`, `noncurrent_version_expiration` must be specified.

<a name="obs_expiration"></a>
The `expiration` object supports:

* `date` (Optional, String) Specifies the date after which you want the corresponding action to take effect.

* `days` (Optional, Int) Specifies the number of days after object creation when the specific rule action takes effect.

* `expired_object_delete_marker` (Optional, Bool) On a versioned bucket (versioning-enabled or versioning-suspended
  bucket), you can add this element in the lifecycle configuration to direct Amazon S3 to delete expired object delete
  markers.

<a name="obs_noncurrent_version_expiration"></a>
The `noncurrent_version_expiration` object supports:

* `days` (Required, Int) Specifies the number of days an object is noncurrent object versions expire.

## Attribute Reference

The following attributes are exported:

* `id` - The name of the bucket.

* `arn` - The ARN of the bucket. Will be of format `arn:aws:s3:::bucketname`.

* `bucket_domain_name` - The bucket domain name. Will be of format `bucketname.s3.amazonaws.com`.

* `hosted_zone_id` - The [Route 53 Hosted Zone ID](https://docs.aws.amazon.com/general/latest/gr/rande.html#s3_website_region_endpoints)
  for this bucket's region.

* `website_endpoint` - The website endpoint, if the bucket is configured with a website. If not, this will be an
  empty string.

* `website_domain` - The domain of the website endpoint, if the bucket is configured with a website.
  If not, this will be an empty string. This is used to create Route 53 alias records.

## Import

S3 bucket can be imported using the `bucket`, e.g.

```shell
terraform import flexibleengine_s3_bucket.bucket bucket-name
```
