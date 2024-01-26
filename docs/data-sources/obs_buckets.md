---
subcategory: "Object Storage Service (OBS)"
---

# flexibleengine_obs_buckets

Use this data source to get all OBS buckets.

```hcl
data "flexibleengine_obs_buckets" "buckets" {
  bucket = "your-bucket-name"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) The region in which to obtain the OBS bucket.
  If omitted, the provider-level region will be used.

* `bucket` - (Optional, String) The name of the OBS bucket.

* `enterprise_project_id` - (Optional, String) The enterprise project id of the OBS bucket.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID of the OBS bucket list.

* `buckets` - A list of OBS buckets. The [buckets](#obs_buckets) object structure is documented below.

<a name="obs_buckets"></a>
The `buckets` block supports:

* `region` - The region where the OBS bucket belongs.

* `bucket` - The name of the OBS bucket.

* `enterprise_project_id` - The enterprise project id of the OBS bucket.

* `storage_class` - The storage class of the OBS bucket.

* `created_at` - The date when the OBS bucket was created.
