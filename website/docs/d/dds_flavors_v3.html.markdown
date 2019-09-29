---
layout: "flexibleengine"
page_title: "FlexibleEngine: flexibleengine_dds_flavors_v3"
sidebar_current: "docs-flexibleengine-datasource-dds-flavors-v3"
description: |-
  Get the flavor information on FlexibleEngine DDS service.
---

# flexibleengine\_dds\_flavors\_v3

Use this data source to get the ID of an available FlexibleEngine dds flavor.

## Example Usage

```hcl
data "flexibleengine_dds_flavors_v3" "flavor" {
    region = "eu-west-0"
    engine_name = "DDS-Community"
}
```

## Argument Reference

* `region` - (Optional) The region in which to obtain the V3 dds client.

* `engine_name` - (Optional) The engine name of the dds, now only DDS-Community is supported.

* `speccode` - (Optional) The spec code of a dds flavor.

## Available value for attributes

engine_name | type | vcpus | ram | speccode
---- | --- | ---
DDS-Community | mongos | 1 | 4 | dds.mongodb.s3.medium.4.mongos
DDS-Community | mongos | 2 | 8 | dds.mongodb.s3.large.4.mongos
DDS-Community | mongos | 4 | 16 | dds.mongodb.s3.xlarge.4.mongos
DDS-Community | mongos | 8 | 32 | dds.mongodb.s3.2xlarge.4.mongos
DDS-Community | mongos | 16 | 64 | dds.mongodb.s3.4xlarge.4.mongos
DDS-Community | shard | 1 | 4 | dds.mongodb.s3.medium.4.shard
DDS-Community | shard | 2 | 8 | dds.mongodb.s3.large.4.shard
DDS-Community | shard | 4 | 16 | dds.mongodb.s3.xlarge.4.shard
DDS-Community | shard | 8 | 32 | dds.mongodb.s3.2xlarge.4.shard
DDS-Community | shard | 16 | 64 | dds.mongodb.s3.4xlarge.4.shard
DDS-Community | config | 2 | 4 | dds.mongodb.s3.large.2.config
DDS-Community | replica | 1 | 4 | dds.mongodb.s3.medium.4.repset
DDS-Community | replica | 2 | 8 | dds.mongodb.s3.large.4.repset
DDS-Community | replica | 4 | 16 | dds.mongodb.s3.xlarge.4.repset
DDS-Community | replica | 8 | 32 | dds.mongodb.s3.2xlarge.4.repset
DDS-Community | replica | 16 | 64 | dds.mongodb.s3.4xlarge.4.repset


## Attributes Reference

* `region` - See Argument Reference above.
* `engine_name` - See Argument Reference above.
* `speccode` - See Argument Reference above.
* `type` - The type of the dds flavor.
* `vcpus` - The vcpus of the dds flavor.
* `ram` - The ram of the dds flavor.
