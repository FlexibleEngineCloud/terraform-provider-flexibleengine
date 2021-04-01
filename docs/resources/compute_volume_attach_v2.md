---
subcategory: "Elastic Cloud Server (ECS)"
---

# flexibleengine\_compute\_volume_attach_v2

Attaches a Block Storage Volume to an Instance using the FlexibleEngine
Compute (Nova) v2 API.

## Example Usage

### Basic attachment of a single volume to a single instance

```hcl
resource "flexibleengine_blockstorage_volume_v2" "volume_1" {
  name = "volume_1"
  size = 1
}

resource "flexibleengine_compute_instance_v2" "instance_1" {
  name            = "instance_1"
  security_groups = ["default"]
}

resource "flexibleengine_compute_volume_attach_v2" "va_1" {
  instance_id = "${flexibleengine_compute_instance_v2.instance_1.id}"
  volume_id   = "${flexibleengine_blockstorage_volume_v2.volume_1.id}"
}
```

### Attaching multiple volumes to a single instance

```hcl
resource "flexibleengine_blockstorage_volume_v2" "volumes" {
  count = 2
  name  = "${format("vol-%02d", count.index + 1)}"
  size  = 1
}

resource "flexibleengine_compute_instance_v2" "instance_1" {
  name            = "instance_1"
  security_groups = ["default"]
}

resource "flexibleengine_compute_volume_attach_v2" "attachments" {
  count       = 2
  instance_id = "${flexibleengine_compute_instance_v2.instance_1.id}"
  volume_id   = "${element(flexibleengine_blockstorage_volume_v2.volumes.*.id, count.index)}"
}

output "volume devices" {
  value = "${flexibleengine_compute_volume_attach_v2.attachments.*.device}"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional) The region in which to obtain the V2 Compute client.
    A Compute client is needed to create a volume attachment. If omitted, the
    `region` argument of the provider is used. Changing this creates a
    new volume attachment.

* `instance_id` - (Required) The ID of the Instance to attach the Volume to.

* `volume_id` - (Required) The ID of the Volume to attach to an Instance.

* `device` - (Optional) The device of the volume attachment (ex: `/dev/vdc`).
  _NOTE_: Being able to specify a device is dependent upon the hypervisor in
  use. There is a chance that the device specified in Terraform will not be
  the same device the hypervisor chose. If this happens, Terraform will wish
  to update the device upon subsequent applying which will cause the volume
  to be detached and reattached indefinitely. Please use with caution.

## Attributes Reference

The following attributes are exported:

* `region` - See Argument Reference above.
* `instance_id` - See Argument Reference above.
* `volume_id` - See Argument Reference above.
* `device` - See Argument Reference above. _NOTE_: The correctness of this
  information is dependent upon the hypervisor in use. In some cases, this
  should not be used as an authoritative piece of information.

## Import

Volume Attachments can be imported using the Instance ID and Volume ID
separated by a slash, e.g.

```
$ terraform import flexibleengine_compute_volume_attach_v2.va_1 89c60255-9bd6-460c-822a-e2b959ede9d2/45670584-225f-46c3-b33e-6707b589b666
```
