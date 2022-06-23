# flexibleengine_availability_zones

Use this data source to get a list of availability zones from FlexibleEngine.

## Example Usage

```hcl
data "flexibleengine_availability_zones" "zones" {}
```

## Argument Reference

* `region` - (Optional) The `region` to fetch availability zones from, defaults to the provider's `region`.

* `state` - (Optional) The `state` of the availability zones to match, default ("available").

## Attributes Reference

`id` is set to hash of the returned zone list. In addition, the following attributes are exported:

* `names` - The names of the availability zones, ordered alphanumerically, that match the queried `state`.
