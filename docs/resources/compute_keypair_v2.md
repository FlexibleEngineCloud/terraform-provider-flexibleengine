---
subcategory: "Elastic Cloud Server (ECS)"
description: ""
page_title: "flexibleengine_compute_keypair_v2"
---

# flexibleengine_compute_keypair_v2

Manages a V2 keypair resource within FlexibleEngine.

## Example Usage

### Create Key Pair

```hcl
resource "flexibleengine_compute_keypair_v2" "new" {
  name = "my-keypair"
}
```

### Import Key Pair

```hcl
resource "flexibleengine_compute_keypair_v2" "import" {
  name       = "eixst-keypair"
  public_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDAjpC1hwiOCCmKEWxJ4qzTTsJbKzndLotBCz5PcwtUnflmU+gHJtWMZKpuEGVi29h0A/+ydKek1O18k10Ff+4tyFjiHDQAnOfgWf7+b1yK+qDip3X1C0UPMbwHlTfSGWLGZqd9LvEFx9k3h/M+VtMvwR1lJ9LUyTAImnNjWG7TaIPmui30HvM2UiFEmqkr4ijq45MyX2+fLIePLRIF61p4whjHAQYufqyno3BS48icQb4p6iVEZPo4AE2o9oIyQvj2mx4dk5Y8CgSETOZTYDOR3rU2fZTRDRgPJDH9FWvQjF5tA0p3d9CoWWd2s6GKKbfoUIi8R/Db1BSPJwkqB"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the keypair resource.
    If omitted, the `region` argument of the provider is used.
    Changing this creates a new keypair.

* `name` - (Required, String, ForceNew) Specifies a unique name for the keypair.
    Changing this creates a new keypair.

* `public_key` - (Optional, String, ForceNew) Specifies a imported OpenSSH-formatted public key.
    Changing this creates a new keypair.

* `private_key_path` - (Optional, String, ForceNew) Specifies the path of the created private key.
    The private key file (**.pem**) is created only after the resource is created.
    By default, the private key file will be created in the same folder as the work directory.
    If you need to create it in another folder, please specify the path for `private_key_path`.
    Changing this creates a new keypair.

    ->**NOTE:** The private key file will be removed after the keypair is deleted.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID which as same as keypair name.

## Import

Keypairs can be imported using the `name`, e.g.

```
$ terraform import flexibleengine_compute_keypair_v2.my-keypair test-keypair
```
