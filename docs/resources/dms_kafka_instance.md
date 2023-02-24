---
subcategory: "Distributed Message Service (DMS)"
description: ""
page_title: "flexibleengine_dms_kafka_instance"
---

# flexibleengine_dms_kafka_instance

Manage a DMS Kafka instance resources within FlexibleEngine.

## Example Usage

```hcl
data "flexibleengine_dms_product" "test" {
  bandwidth = "100MB"
}

resource "flexibleengine_vpc_v1" "example_vpc" {
  name = "example-vpc"
  cidr = "192.168.0.0/16"
}

resource "flexibleengine_vpc_subnet_v1" "example_subnet" {
  name       = "example-vpc-subnet"
  cidr       = "192.168.0.0/24"
  gateway_ip = "192.168.0.1"
  vpc_id     = flexibleengine_vpc_v1.example_vpc.id
}

resource "flexibleengine_networking_secgroup_v2" "example_secgroup" {
  name        = "example-secgroup"
  description = "My neutron security group"
}

resource "flexibleengine_dms_kafka_instance" "product_1" {
  name               = "instance_1"
  engine_version     = "2.3.0"
  bandwidth          = "100MB"
  availability_zones = data.flexibleengine_dms_product.product_1.availability_zones
  product_id         = data.flexibleengine_dms_product.test.id
  storage_space      = data.flexibleengine_dms_product.test.storage_space
  storage_spec_code  = "dms.physical.storage.ultra"

  vpc_id            = flexibleengine_vpc_v1.example_vpc.id
  network_id        = flexibleengine_vpc_subnet_v1.example_subnet.id
  security_group_id = flexibleengine_networking_secgroup_v2.example_secgroup.id

  manager_user     = "admin"
  manager_password = "AdminTest@123"
  access_user      = "user"
  password         = "Kafkatest@123"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the DMS Kafka instance resource.
  If omitted, the provider-level region will be used. Changing this creates a new instance resource.

* `name` - (Required, String) Specifies the name of the DMS Kafka instance. An instance name starts with a letter,
  consists of 4 to 64 characters, and supports only letters, digits, hyphens (-) and underscores (_).

* `bandwidth` - (Required, String, ForceNew) The baseline bandwidth of a Kafka instance, that is, the maximum amount of
  data transferred per unit time. The valid values are **100MB**, **300MB**, **600MB** and **1200MB**.
  Changing this creates a new instance resource.

* `product_id` - (Required, String, ForceNew) Specifies a product ID. You can get the value from id of
  [flexibleengine_dms_product](https://registry.terraform.io/providers/FlexibleEngineCloud/flexibleengine/latest/docs/data-sources/dms_product)
  data source. Changing this creates a new instance resource.

* `storage_space` - (Required, Int, ForceNew) Specifies the message storage capacity, the unit is GB. Value range:
  + When bandwidth is **100MB**: 600–90,000 GB
  + When bandwidth is **300MB**: 1,200–90,000 GB
  + When bandwidth is **600MB**: 2,400–90,000 GB
  + When bandwidth is **1,200MB**: 4,800–90,000 GB

  Changing this creates a new instance resource.

* `availability_zones` - (Required, List, ForceNew) The names of the AZ where the Kafka instance resides.
  Changing this creates a new instance resource.

  -> **NOTE:** Deploy one availability zone or at least three availability zones. Do not select two availability zones.
  Deploy to more availability zones, the better the reliability and SLA coverage.

* `vpc_id` - (Required, String, ForceNew) Specifies the ID of a VPC.
  Changing this creates a new instance resource.

* `network_id` - (Required, String, ForceNew) Specifies the ID of a VPC subnet.
  Changing this creates a new instance resource.

* `security_group_id` - (Required, String) Specifies the ID of a security group.

* `description` - (Optional, String) Specifies the description of the DMS Kafka instance.
  It is a character string containing not more than 1,024 characters.

* `engine_version` - (Optional, String, ForceNew) Specifies the version of the Kafka engine. Valid values are "1.1.0"
  and "2.3.0". Defaults to **2.3.0**. Changing this creates a new instance resource.

* `storage_spec_code` - (Optional, String, ForceNew) Specifies the storage I/O specification. Value range:
  + When bandwidth is **100MB**: dms.physical.storage.high or dms.physical.storage.ultra
  + When bandwidth is **300MB**: dms.physical.storage.high or dms.physical.storage.ultra
  + When bandwidth is **600MB**: dms.physical.storage.ultra
  + When bandwidth is **1,200MB**: dms.physical.storage.ultra

  Defaults to **dms.physical.storage.ultra**. Changing this creates a new instance resource.

* `manager_user` - (Optional, String, ForceNew) Specifies the username for logging in to the Kafka Manager.
  The username consists of 4 to 64 characters and can contain letters, digits, hyphens (-), and underscores (_).
  Changing this creates a new instance resource.

* `manager_password` - (Optional, String, ForceNew) Specifies the password for logging in to the Kafka Manager. The
  password must meet the following complexity requirements: Must be 8 to 32 characters long. Must contain at least 2 of
  the following character types: lowercase letters, uppercase letters, digits, and special characters (`~!@#$%^&*()-_
  =+\\|[{}]:'",<.>/?). Changing this creates a new instance resource.

* `access_user` - (Optional, String, ForceNew) Specifies a username who can accesse the instance with
  SASL authentication. A username consists of 4 to 64 characters and supports only letters, digits, and hyphens (-).
  Changing this creates a new instance resource.

* `password` - (Optional, String, ForceNew) Specifies the password of the access user. A password must meet the
  following complexity requirements: Must be 8 to 32 characters long. Must contain at least 2 of the following character
  types: lowercase letters, uppercase letters, digits, and special characters (`~!@#$%^&*()-_=+\\|[{}]:'",<.>/?).
  Changing this creates a new instance resource.

  -> **NOTE:** If `access_user` and `password` are specified, Kafka SASL_SSL will be automatically enabled.

* `maintain_begin` - (Optional, String) Specifies the time at which a maintenance time window starts. Format: HH:mm:ss.
  The start time must be set to 22:00:00, 02:00:00, 06:00:00, 10:00:00, 14:00:00, or 18:00:00.
  The system automatically allocates the default start time 02:00:00.

* `maintain_end` - (Optional, String) Specifies the time at which a maintenance time window ends. Format: HH:mm:ss.
  The end time is four hours later than the start time. For example, if the start time is 22:00:00, the end time is
  02:00:00. The system automatically allocates the default end time 06:00:00.

  -> **NOTE:**  The start time and end time of a maintenance time window must be set in pairs.

* `enable_auto_topic` - (Optional, Bool, ForceNew) Specifies whether to enable automatic topic creation. If automatic
  topic creation is enabled, a topic will be automatically created with 3 partitions and 3 replicas when a message is
  produced to or consumed from a topic that does not exist. Changing this creates a new instance resource.

* `tags` - (Optional, Map) Specifies the key/value pairs to associate with the DMS Kafka instance.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Specifies a resource ID in UUID format.
* `status` - Indicates the status of the DMS Kafka instance.
* `engine` - Indicates the message engine, the value is "kafka".
* `engine_type` - Indicates the DMS Kafka instance type, the value is "cluster".
* `product_spec_code` - Indicates the DMS Kafka instance specification.
* `partition_num` - Indicates the maximum number of topics in the DMS Kafka instance.
* `used_storage_space` - Indicates the used message storage space. Unit: GB
* `vpc_name` - Indicates the name of a vpc.
* `subnet_name` - Indicates the name of a subnet.
* `security_group_name` - Indicates the name of a security group.
* `node_num` - Indicates the count of ECS instances.
* `manegement_connect_address` - Indicates the connection address of the Kafka Manager of a Kafka instance.
* `connect_address` - Indicates the IP addresses of the DMS Kafka instance.
* `port` - Indicates the port number of the DMS Kafka instance.
* `ssl_enable` - Indicates whether the Kafka SASL_SSL is enabled.
* `created_at` - Indicates the creation time of the DMS Kafka instance.

## Import

DMS Kafka instance can be imported using the instance id, e.g.

```shell
terraform import flexibleengine_dms_kafka_instance.instance_1 8d3c7938-dc47-4937-a30f-c80de381c5e3
```

Note that the imported state may not be identical to your resource definition, because of `access_user`, `password`,
`manager_user` and `manager_password` are missing from the API response due to security reason.
It is generally recommended running `terraform plan` after importing a DMS Kafka instance.
You can then decide if changes should be applied to the instance, or the resource
definition should be updated to align with the instance. Also you can ignore changes as below.

```hcl
resource "flexibleengine_dms_kafka_instance" "instance_1" {
    ...

  lifecycle {
    ignore_changes = [
      access_user, password, manager_user, manager_password,
    ]
  }
}
```
