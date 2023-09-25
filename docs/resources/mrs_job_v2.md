---
subcategory: "MapReduce Service (MRS)"
description: ""
page_title: "flexibleengine_mrs_job_v2"
---

# flexibleengine_mrs_job_v2

Manage a job resource within FlexibleEngine MRS.

## Example Usage

```hcl
variable "cluster_id" {}
variable "job_name" {}
variable "program_path" {}
variable "access_key" {}
variable "secret_key" {}

resource "flexibleengine_mrs_job_v2" "test" {
  cluster_id   = var.cluster_id
  type         = "SparkSubmit"
  name         = var.job_name
  program_path = var.program_path
  parameters   = "${var.access_key} ${var.secret_key} 1 s3a://obs-demo-analysis/input s3a://obs-demo-analysis/output"

  program_parameters = {
    "--class" = "com.orange.bigdata.spark.examples.DriverBehavior"
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the MRS job resource.
  If omitted, the provider-level region will be used. Changing this will create a new MRS job resource.

* `cluster_id` - (Required, String, ForceNew) Specifies an ID of the MRS cluster to which the job belongs to.
  Changing this will create a new MRS job resource.

* `name` - (Required, String, ForceNew) Specifies the name of the MRS job. The name can contain 1 to 64
  characters, which may consist of letters, digits, underscores (_) and hyphens (-).
  Changing this will create a new MRS job resource.

* `type` - (Required, String, ForceNew) Specifies the job type. The valid values are **MapReduce**,
  **Flink**, **HiveSql**, **HiveScript**, **SparkSubmit**, **SparkSql** and **SparkScript**.

  Changing this will create a new MRS job resource.

  -> Spark, Hive, and Flink jobs can be added to only clusters that include Spark, Hive, and Flink components.

* `program_path` - (Optional, String, ForceNew) Specifies the .jar package path or .py file path for program execution.
  The parameter must meet the following requirements:
  + Contains a maximum of 1023 characters, excluding special characters such as `;|&><'$`.
  + The address cannot be empty or full of spaces.
  + The program support OBS or DHFS to storage program file or package. For OBS, starts with (OBS:) **s3a://** and end
      with **.jar** or **.py**. For DHFS, starts with (DHFS:) **/user**.

  Required if `type` is **MapReduce** or **SparkSubmit**. Changing this will create a new MRS job resource.

* `parameters` - (Optional, String, ForceNew) Specifies the parameters for the MRS job. Add an at sign (@) before
  each parameter can prevent the parameters being saved in plaintext format. Each parameters are separated with spaces.
  This parameter can be set when `type` is **Flink**, **MRS** or **SparkSubmit**. Changing this will create a new
  MRS job resource.

* `program_parameters` - (Optional, Map, ForceNew) Specifies the the key/value pairs of the program parameters, such as
  thread, memory, and vCPUs, are used to optimize resource usage and improve job execution performance. This parameter
  can be set when `type` is **Flink**, **SparkSubmit**, **SparkSql**, **SparkScript**, **HiveSql** or
  **HiveScript**. Changing this will create a new MRS job resource.

* `service_parameters` - (Optional, Map, ForceNew) Specifies the key/value pairs used to modify service configuration.
  Parameter configurations of services are available on the Service Configuration tab page of MRS Manager.
  Changing this will create a new MRS job resource.

* `sql` - (Optional, String, ForceNew) Specifies the SQL command or file path. Only required if `type` is **HiveSql**
  or **SparkSql**. Changing this will create a new MRS job resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the MRS job in UUID format.

* `status` - Status of the MRS job.

* `start_time` - The creation time of the MRS job.

* `submit_time` - The submission time of the MRS job.

* `finish_time` - The completion time of the MRS job.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 20 minutes.

## Import

MRS jobs can be imported using their `id` and the IDs of the MRS cluster to which the job belongs, separated
by a slash, e.g.

```shell
terraform import flexibleengine_mrs_job_v2.test <cluster_id>/<id>
```
