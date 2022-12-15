---
subcategory: "Object Storage Service (OSS)"
---

# flexibleengine_obs_bucket_replication

Manages an OBS bucket **Cross-Region Replication** resource within FlexibleEngine.

[Cross-Region replication](https://docs.prod-cloud-ocb.orange-business.com/usermanual/obs/obs_03_0002.html)
provides the capability for data disaster recovery across regions, catering to your needs for off-site data backup.

## Example Usage

### Replicate all objects

```hcl
resource "flexibleengine_obs_bucket_replication" "replica" {
  bucket             = "my-source-bucket"
  destination_bucket = "my-target-bucket"
  agency             = "obs-fullaccess"
}
```

### Replicate objects matched by prefix

```hcl
resource "flexibleengine_obs_bucket_replication" "replica" {
  bucket             = "my-source-bucket"
  destination_bucket = "my-target-bucket"
  agency             = "obs-fullaccess"

  rule {
    enabled = true
    prefix  = "log"
  }

  rule {
    enabled = false
    prefix  = "imgs/"
  }
}
```

## Argument Reference

The following arguments are supported:

* `bucket` - (Required) Specifies the name of the source bucket. Changing this parameter will create a new resource.

* `destination_bucket` - (Required) Specifies the name of the destination bucket.

  -> The destination bucket cannot be in the region where the source bucket resides.

* `agency` - (Required) Specifies the IAM agency applied to the cross-region replication function.

  -> The IAM agency is a cloud service agency of OBS. The OBS project must have the **OBS FullAccess** permissions.

* `rule` - (Optional) A configuration of object cross-region replication management. The object supports the following:

  + `enabled` - (Optional) Specifies cross-region replication rule status. Defaults to `true`.

  + `prefix` - (Optional) Specifies the object key prefix identifying one or more objects to which the rule applies and
    duplicated prefixes are not supported. If omitted, all objects in the bucket will be managed by the lifecycle rule.
    To copy a folder, end the prefix with a slash (/), for example, imgs/.

  + `storage_class` - (Optional) Specifies the storage class for replicated objects. Valid values are "STANDARD",
    "WARM" (Infrequent Access) and "COLD" (Archive).
    If omitted, the storage class of object copies is the same as that of objects in the source bucket.

## Attributes Reference

The following attributes are exported:

* `id` - The name of the bucket.
* `rule/id` - The ID of a rule in UUID format.

## Import

OBS bucket cross-region replication can be imported using the *source bucket name*, e.g.

```
$ terraform import flexibleengine_obs_bucket_replication.replica my-source-bucket
```
