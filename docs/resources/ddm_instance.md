---
subcategory: "Distributed Database Middleware (DDM)"
---

# flexibleengine_ddm_instance

Manages DDM instance resource within FlexibleEngine.

## Example Usage

```hcl
resource "flexibleengine_vpc_v1" "test" {
  name = "test_vpc"
  cidr = "192.168.0.0/24"
}

resource "flexibleengine_vpc_subnet_v1" "test" {
  name       = "test_subnet"
  cidr       = "192.168.0.0/24"
  gateway_ip = "192.168.0.1"
  vpc_id     = flexibleengine_vpc_v1.test.id
}

resource "flexibleengine_networking_secgroup_v2" "test" {
  name = "test_secgroup"
}

data "flexibleengine_availability_zones" "test" {}

data "flexibleengine_ddm_engines" test {
  version = "3.0.8.5"
}

data "flexibleengine_ddm_flavors" test {
  engine_id = data.flexibleengine_ddm_engines.test.engines[0].id
  cpu_arch  = "X86"
}

resource "flexibleengine_ddm_instance" "test" {
  name              = "ddm_test"
  flavor_id         = data.flexibleengine_ddm_flavors.test.flavors[0].id
  node_num          = 2
  engine_id         = data.flexibleengine_ddm_engines.test.engines[0].id
  vpc_id            = flexibleengine_vpc_v1.test.id
  subnet_id         = flexibleengine_vpc_subnet_v1.test.id
  security_group_id = flexibleengine_networking_secgroup_v2.test.id
  
  availability_zones = [
    data.flexibleengine_availability_zones.test.names[0]
  ]
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `name` - (Required, String) Specifies the name of the DDM instance.
  An instance name starts with a letter, consists of 4 to 64 characters, and can contain only letters,
  digits, and hyphens (-).

* `flavor_id` - (Required, String, ForceNew) Specifies the ID of a product.

  Changing this parameter will create a new resource.

* `node_num` - (Required, Int) Specifies the number of nodes.

* `engine_id` - (Required, String, ForceNew) Specifies the ID of an Engine.

  Changing this parameter will create a new resource.

* `availability_zones` - (Required, List, ForceNew) Specifies the list of availability zones.

  Changing this parameter will create a new resource.

* `vpc_id` - (Required, String, ForceNew) Specifies the ID of a VPC.

  Changing this parameter will create a new resource.

* `security_group_id` - (Required, String) Specifies the ID of a security group.

* `subnet_id` - (Required, String, ForceNew) Specifies the ID of a subnet.

  Changing this parameter will create a new resource.

* `param_group_id` - (Optional, String, ForceNew) Specifies the ID of parameter group.

  Changing this parameter will create a new resource.

* `time_zone` - (Optional, String, ForceNew) Specifies the time zone.

  Changing this parameter will create a new resource.

* `admin_user` - (Optional, String, ForceNew) Specifies the username of the administrator.
  The user name starts with a letter, consists of 1 to 32 characters, and can contain only letters,
  digits, and underscores (_).

  Changing this parameter will create a new resource.

* `admin_password` - (Optional, String, ForceNew) Specifies the password of the administrator.
  The password consists of 8 to 32 characters, and must be a combination of uppercase letters,
  lowercase letters, digits, and the following special characters: ~!@#%^*-_=+?.

  Changing this parameter will create a new resource.

* `delete_rds_data` - (Optional, String) Specifies whether data stored on the associated DB instances is deleted.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `status` - Indicates the status of the DDM instance.

* `access_ip` - Indicates the address for accessing the DDM instance.

* `access_port` - Indicates the port for accessing the DDM instance.

* `engine_version` - Indicates the engine version.

* `nodes` - Indicates the node information.
  The [NodeInfoRef](#DdmInstance_NodeInfoRef) structure is documented below.

<a name="DdmInstance_NodeInfoRef"></a>
The `NodeInfoRef` block supports:

* `status` - Indicates the status of the DDM instance node.

* `port` - Indicates the port of the DDM instance node.

* `ip` - Indicates the IP address of the DDM instance node.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 15 minutes.
* `update` - Default is 60 minutes.
* `delete` - Default is 10 minutes.

## Import

The ddm instance can be imported using the `id`, e.g.

```bash
terraform import flexibleengine_ddm_instance.test 4bc36477c36642479acf2d90751c8c29in09
```
