---
subcategory: "Object Storage Service (OBS)"
description: ""
page_title: "flexibleengine_s3_bucket_policy"
---

# flexibleengine_s3_bucket_policy

Attaches a policy to an S3 bucket resource.

## Example Usage

### Basic Usage

```hcl
resource "flexibleengine_s3_bucket" "b" {
  bucket = "my_tf_test_bucket"
}

resource "flexibleengine_s3_bucket_policy" "b" {
  bucket = flexibleengine_s3_bucket.b.id
  policy = <<POLICY
{
  "Version": "2012-10-17",
  "Id": "MYBUCKETPOLICY",
  "Statement": [
    {
      "Sid": "IPAllow",
      "Effect": "Deny",
      "Principal": "*",
      "Action": "s3:*",
      "Resource": "arn:aws:s3:::my_tf_test_bucket/*",
      "Condition": {
         "IpAddress": {"aws:SourceIp": "8.8.8.8/32"}
      } 
    } 
  ]
}
POLICY
}
```

## Argument Reference

The following arguments are supported:

* `bucket` - (Required, String, ForceNew) The name of the bucket to which to apply the policy.

* `policy` - (Required, String) The text of the policy.

## Attribute Reference

All the arguments above can also be exported attributes.
