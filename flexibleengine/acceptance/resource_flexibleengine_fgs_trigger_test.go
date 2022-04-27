package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/fgs/v2/trigger"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func getTriggerResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	c, err := conf.FgsV2Client(OS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating FunctionGraph v2 client: %s", err)
	}
	return trigger.Get(c, state.Primary.Attributes["function_urn"], state.Primary.Attributes["type"],
		state.Primary.ID).Extract()
}

func TestAccFunctionGraphTrigger_basic(t *testing.T) {
	var (
		randName     = acceptance.RandomAccResourceName()
		resourceName = "flexibleengine_fgs_trigger.test"
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccFunctionGraphTimingTrigger_basic(randName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "type", "TIMER"),
					resource.TestCheckResourceAttr(resourceName, "timer.0.name", randName),
					resource.TestCheckResourceAttr(resourceName, "timer.0.schedule_type", "Rate"),
					resource.TestCheckResourceAttr(resourceName, "timer.0.schedule", "3d"),
					resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
					resource.TestCheckResourceAttrPair(resourceName, "function_urn",
						"flexibleengine_fgs_function.test", "urn"),
				),
			},
			{
				Config: testAccFunctionGraphTimingTrigger_update(randName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "type", "TIMER"),
					resource.TestCheckResourceAttr(resourceName, "timer.0.name", randName),
					resource.TestCheckResourceAttr(resourceName, "timer.0.schedule_type", "Rate"),
					resource.TestCheckResourceAttr(resourceName, "timer.0.schedule", "3d"),
					resource.TestCheckResourceAttr(resourceName, "status", "DISABLED"),
					resource.TestCheckResourceAttrPair(resourceName, "function_urn",
						"flexibleengine_fgs_function.test", "urn"),
				),
			},
		},
	})
}

func TestAccFunctionGraphTrigger_cronTimer(t *testing.T) {
	var (
		randName     = acceptance.RandomAccResourceName()
		resourceName = "flexibleengine_fgs_trigger.test"
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccFunctionGraphTimingTrigger_cron(randName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "type", "TIMER"),
					resource.TestCheckResourceAttr(resourceName, "timer.0.name", randName),
					resource.TestCheckResourceAttr(resourceName, "timer.0.schedule_type", "Cron"),
					resource.TestCheckResourceAttr(resourceName, "timer.0.schedule", "@every 1h30m"),
					resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
					resource.TestCheckResourceAttrPair(resourceName, "function_urn",
						"flexibleengine_fgs_function.test", "urn"),
				),
			},
			{
				Config: testAccFunctionGraphTimingTrigger_cronUpdate(randName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "type", "TIMER"),
					resource.TestCheckResourceAttr(resourceName, "timer.0.name", randName),
					resource.TestCheckResourceAttr(resourceName, "timer.0.schedule_type", "Cron"),
					resource.TestCheckResourceAttr(resourceName, "timer.0.schedule", "@every 1h30m"),
					resource.TestCheckResourceAttr(resourceName, "status", "DISABLED"),
					resource.TestCheckResourceAttrPair(resourceName, "function_urn",
						"flexibleengine_fgs_function.test", "urn"),
				),
			},
		},
	})
}

// OBS trigger does not suppport status updation.
func TestAccFunctionGraphTrigger_obs(t *testing.T) {
	var (
		// The underscores (_) are not allowed.
		randName     = fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
		resourceName = "flexibleengine_fgs_trigger.test"
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckOBS(t)
		},
		ProviderFactories: TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccFunctionGraphObsTrigger_basic(randName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "type", "OBS"),
					resource.TestCheckResourceAttr(resourceName, "obs.0.bucket_name", randName),
					resource.TestCheckResourceAttr(resourceName, "obs.0.event_notification_name", randName),
					resource.TestCheckResourceAttr(resourceName, "obs.0.suffix", ".json"),
					resource.TestCheckResourceAttr(resourceName, "obs.0.events.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
					resource.TestCheckResourceAttrPair(resourceName, "function_urn",
						"flexibleengine_fgs_function.test", "urn"),
				),
			},
		},
	})
}

func TestAccFunctionGraphTrigger_dis(t *testing.T) {
	var (
		randName     = acceptance.RandomAccResourceName()
		resourceName = "flexibleengine_fgs_trigger.test"
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccFunctionGraphDisTrigger_basic(randName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "type", "DIS"),
					resource.TestCheckResourceAttr(resourceName, "dis.0.stream_name", randName),
					resource.TestCheckResourceAttr(resourceName, "dis.0.starting_position", "TRIM_HORIZON"),
					resource.TestCheckResourceAttr(resourceName, "dis.0.max_fetch_bytes", "2097152"),
					resource.TestCheckResourceAttr(resourceName, "dis.0.pull_period", "30000"),
					resource.TestCheckResourceAttr(resourceName, "dis.0.serial_enable", "true"),
					resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
					resource.TestCheckResourceAttrPair(resourceName, "function_urn",
						"flexibleengine_fgs_function.test", "urn"),
				),
			},
			{
				Config: testAccFunctionGraphDisTrigger_update(randName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "type", "DIS"),
					resource.TestCheckResourceAttr(resourceName, "dis.0.stream_name", randName),
					resource.TestCheckResourceAttr(resourceName, "dis.0.starting_position", "TRIM_HORIZON"),
					resource.TestCheckResourceAttr(resourceName, "dis.0.max_fetch_bytes", "2097152"),
					resource.TestCheckResourceAttr(resourceName, "dis.0.pull_period", "30000"),
					resource.TestCheckResourceAttr(resourceName, "dis.0.serial_enable", "true"),
					resource.TestCheckResourceAttr(resourceName, "status", "DISABLED"),
					resource.TestCheckResourceAttrPair(resourceName, "function_urn",
						"flexibleengine_fgs_function.test", "urn"),
				),
			},
		},
	})
}

func TestAccFunctionGraphTrigger_smn(t *testing.T) {
	var (
		randName     = acceptance.RandomAccResourceName()
		resourceName = "flexibleengine_fgs_trigger.test"
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccFunctionGraphSmnTrigger_basic(randName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "type", "SMN"),
					resource.TestCheckResourceAttrSet(resourceName, "smn.0.topic_urn"),
					resource.TestCheckResourceAttrPair(resourceName, "function_urn",
						"flexibleengine_fgs_function.test", "urn"),
				),
			},
		},
	})
}

/*
func TestAccFunctionGraphTrigger_kafka(t *testing.T) {
	var (
		randName  = acceptance.RandomAccResourceName()
		adminPass = fmt.Sprintf("%s%s%d", acctest.RandString(5), acctest.RandStringFromCharSet(2, "#$"),
			acctest.RandIntRange(100, 999))
		resourceName = "flexibleengine_fgs_trigger.test"
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    acceptance.TestAccProviders,
		//CheckDestroy: rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccFunctionGraphKafkaTrigger_basic(randName, adminPass),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "type", "KAFKA"),
					resource.TestCheckResourceAttr(resourceName, "kafka.0.batch_size", "100"),
					resource.TestCheckResourceAttr(resourceName, "kafka.0.topic_ids.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
					resource.TestCheckResourceAttrPair(resourceName, "function_urn",
						"flexibleengine_fgs_function.test", "urn"),
					resource.TestCheckResourceAttrPair(resourceName, "kafka.0.instance_id",
						"flexibleengine_dms_kafka_instance.test", "id"),
				),
			},
			{
				Config: testAccFunctionGraphKafkaTrigger_update(randName, adminPass),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "type", "KAFKA"),
					resource.TestCheckResourceAttr(resourceName, "kafka.0.batch_size", "100"),
					resource.TestCheckResourceAttr(resourceName, "kafka.0.topic_ids.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "status", "DISABLED"),
					resource.TestCheckResourceAttrPair(resourceName, "function_urn",
						"flexibleengine_fgs_function.test", "urn"),
					resource.TestCheckResourceAttrPair(resourceName, "kafka.0.instance_id",
						"flexibleengine_dms_kafka_instance.test", "id"),
				),
			},
		},
	})
}
*/

func testAccFunctionGraphTimingTrigger_base(rName string) string {
	return fmt.Sprintf(`
resource "flexibleengine_fgs_function" "test" {
  name        = "%s"
  app         = "default"
  handler     = "index.handler"
  memory_size = 128
  timeout     = 10
  runtime     = "Python2.7"
  code_type   = "inline"
  func_code   = "aW1wb3J0IGpzb24KZGVmIGhhbmRsZXIgKGZW50LCBjb250ZXh0KToKICAgIG91dHB1dCA9ICdIZWxsbyBtZXNzYWdlOiAnICsganNvbi5kdW1wcyhldmVudCkKICAgIHJldHVybiBvdXRwdXQ="
}`, rName)
}

func testAccFunctionGraphTimingTrigger_basic(rName string) string {
	return fmt.Sprintf(`
%s

resource "flexibleengine_fgs_trigger" "test" {
  function_urn = flexibleengine_fgs_function.test.urn
  type         = "TIMER"

  timer {
    name          = "%s"
    schedule_type = "Rate"
    schedule      = "3d"
  }
}
`, testAccFunctionGraphTimingTrigger_base(rName), rName)
}

func testAccFunctionGraphTimingTrigger_update(rName string) string {
	return fmt.Sprintf(`
%s

resource "flexibleengine_fgs_trigger" "test" {
  function_urn = flexibleengine_fgs_function.test.urn
  type         = "TIMER"
  status       = "DISABLED"

  timer {
	name          = "%s"
	schedule_type = "Rate"
	schedule      = "3d"
  }
}
`, testAccFunctionGraphTimingTrigger_base(rName), rName)
}

func testAccFunctionGraphTimingTrigger_cron(rName string) string {
	return fmt.Sprintf(`
%s

resource "flexibleengine_fgs_trigger" "test" {
  function_urn = flexibleengine_fgs_function.test.urn
  type         = "TIMER"

  timer {
    name          = "%s"
    schedule_type = "Cron"
    schedule      = "@every 1h30m"
  }
}
`, testAccFunctionGraphTimingTrigger_base(rName), rName)
}

func testAccFunctionGraphTimingTrigger_cronUpdate(rName string) string {
	return fmt.Sprintf(`
%s

resource "flexibleengine_fgs_trigger" "test" {
  function_urn = flexibleengine_fgs_function.test.urn
  type         = "TIMER"
  status       = "DISABLED"

  timer {
	name          = "%s"
	schedule_type = "Cron"
	schedule      = "@every 1h30m"
  }
}
`, testAccFunctionGraphTimingTrigger_base(rName), rName)
}

func testAccFunctionGraphObsTrigger_base(rName string) string {
	return fmt.Sprintf(`
resource "flexibleengine_obs_bucket" "test" {
  bucket = "%s"
  acl    = "private"
}

resource "flexibleengine_identity_agency_v3" "test" {
  name                   = "%s"
  delegated_service_name = "op_svc_cff"

  project_role {
    project = "MOS"

    roles = [
      "OBS OperateAccess",
    ]
  }

  domain_roles = [
    "OBS OperateAccess",
  ]
}

data "flexibleengine_fgs_dependencies" "test" {
  runtime = "Python2.7"
  name    = "esdk_obs_python-3x"
}

resource "flexibleengine_fgs_function" "test" {
  name        = "%s"
  app         = "default"
  agency      = flexibleengine_identity_agency_v3.test.name
  handler     = "index.handler"
  memory_size = 256
  timeout     = 15
  runtime     = "Python2.7"
  code_type   = "inline"
  func_code   = "aW1wb3J0IGpzb24KZGVmIGhhbmRsZXIgKGZW50LCBjb250ZXh0KToKICAgIG91dHB1dCA9ICdIZWxsbyBtZXNzYWdlOiAnICsganNvbi5kdW1wcyhldmVudCkKICAgIHJldHVybiBvdXRwdXQ="
  depend_list = [data.flexibleengine_fgs_dependencies.test.packages[0].id]
}`, rName, rName, rName)
}

func testAccFunctionGraphObsTrigger_basic(rName string) string {
	return fmt.Sprintf(`
%s

resource "flexibleengine_fgs_trigger" "test" {
  function_urn = flexibleengine_fgs_function.test.urn
  type         = "OBS"
  status       = "ACTIVE"

  obs {
    bucket_name             = flexibleengine_obs_bucket.test.bucket
    event_notification_name = "%s"
    suffix                  = ".json"

    events = ["ObjectCreated"]
  }
}`, testAccFunctionGraphObsTrigger_base(rName), rName)
}

func testAccFunctionGraphDisTrigger_base(rName string) string {
	return fmt.Sprintf(`
resource "flexibleengine_dis_stream" "test" {
  name            = "%s"
  type            = "COMMON"
  partition_count = 3
}

resource "flexibleengine_identity_agency_v3" "test" {
  name                   = "%s"
  delegated_service_name = "op_svc_cff"

  project_role {
    project = "%s"

    roles = [
      "DIS Administrator",
    ]
  }
}

resource "flexibleengine_fgs_function" "test" {
  name        = "%s"
  app         = "default"
  agency      = flexibleengine_identity_agency_v3.test.name
  handler     = "index.handler"
  memory_size = 128
  timeout     = 10
  runtime     = "Python2.7"
  code_type   = "inline"
  func_code   = "aW1wb3J0IGpzb24KZGVmIGhhbmRsZXIgKGZW50LCBjb250ZXh0KToKICAgIG91dHB1dCA9ICdIZWxsbyBtZXNzYWdlOiAnICsganNvbi5kdW1wcyhldmVudCkKICAgIHJldHVybiBvdXRwdXQ="
}`, rName, rName, OS_REGION_NAME, rName)
}

func testAccFunctionGraphDisTrigger_basic(rName string) string {
	return fmt.Sprintf(`
%s

resource "flexibleengine_fgs_trigger" "test" {
  function_urn = flexibleengine_fgs_function.test.urn
  type         = "DIS"
  status       = "ACTIVE"

  dis {
    stream_name       = flexibleengine_dis_stream.test.name
    starting_position = "TRIM_HORIZON"
	max_fetch_bytes   = 2097152
    pull_period       = 30000
    serial_enable     = true
  }
}`, testAccFunctionGraphDisTrigger_base(rName))
}

func testAccFunctionGraphDisTrigger_update(rName string) string {
	return fmt.Sprintf(`
%s

resource "flexibleengine_fgs_trigger" "test" {
  function_urn = flexibleengine_fgs_function.test.urn
  type         = "DIS"
  status       = "DISABLED"

  dis {
    stream_name       = flexibleengine_dis_stream.test.name
    starting_position = "TRIM_HORIZON"
	max_fetch_bytes   = 2097152
    pull_period       = 30000
    serial_enable     = true
  }
}`, testAccFunctionGraphDisTrigger_base(rName))
}

func testAccFunctionGraphSmnTrigger_basic(rName string) string {
	return fmt.Sprintf(`
%s

resource "flexibleengine_smn_topic_v2" "test" {
  name = "%s"
}

resource "flexibleengine_fgs_trigger" "test" {
  function_urn = flexibleengine_fgs_function.test.urn
  type         = "SMN"

  smn {
    topic_urn = flexibleengine_smn_topic_v2.test.topic_urn
  }
}`, testAccFunctionGraphTimingTrigger_base(rName), rName)
}

func testAccNetwork_config(rName string) string {
	return fmt.Sprintf(`
data "flexibleengine_availability_zones" "test" {}

resource "flexibleengine_vpc_v1" "test" {
  name = "%s"
  cidr = "192.168.128.0/20"
}

resource "flexibleengine_vpc_subnet_v1" "test" {
  name       = "%s"
  vpc_id     = flexibleengine_vpc_v1.test.id
  cidr       = "192.168.128.0/24"
  gateway_ip = "192.168.128.1"
}

resource "flexibleengine_networking_secgroup_v2" "test" {
  name = "%s"
}

resource "flexibleengine_networking_secgroup_rule" "test" {
  security_group_id = flexibleengine_networking_secgroup_v2.test.id
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  port_range_min    = 9092
  port_range_max    = 9092
  remote_ip_prefix  = "0.0.0.0/0"
}`, rName, rName, rName)
}

func testAccDmsKafka_config(rName, password string) string {
	return fmt.Sprintf(`
%s

data "flexibleengine_dms_az" "test" {}

data "flexibleengine_dms_product" "test" {
  engine            = "kafka"
  version           = "1.1.0"
  instance_type     = "cluster"
  partition_num     = 300
  storage           = 600
  storage_spec_code = "dms.physical.storage.high"
}

resource "flexibleengine_dms_kafka_instance" "test" {
  name              = "%s"
  vpc_id            = flexibleengine_vpc_v1.test.id
  network_id        = flexibleengine_vpc_subnet_v1.test.id
  security_group_id = flexibleengine_networking_secgroup_v2.test.id
  available_zones   = [data.flexibleengine_dms_az.test.id]
  product_id        = data.flexibleengine_dms_product.test.id
  engine_version    = data.flexibleengine_dms_product.test.version
  bandwidth         = data.flexibleengine_dms_product.test.bandwidth
  storage_space     = data.flexibleengine_dms_product.test.storage
  storage_spec_code = data.flexibleengine_dms_product.test.storage_spec_code
  manager_user      = "%s"
  manager_password  = "%s"
}

resource "flexibleengine_dms_kafka_topic" "test" {
  instance_id = flexibleengine_dms_kafka_instance.test.id
  name        = "%s"
  partitions  = 20
}`, testAccNetwork_config(rName), rName, rName, password, rName)
}

func testAccFunctionGraphKafkaTrigger_base(rName, password string) string {
	return fmt.Sprintf(`
%s

resource "flexibleengine_identity_agency_v3" "test" {
  name                   = "%s"
  delegated_service_name = "op_svc_cff"

  project_role {
	project = "%s"

    roles = [
      "DMS Administrator",
      "VPC FullAccess",
    ]
  }
}

resource "flexibleengine_fgs_function" "test" {
  name        = "%s"
  app         = "default"
  handler     = "index.handler"
  memory_size = 128
  timeout     = 10
  runtime     = "Python2.7"
  code_type   = "inline"
  agency      = flexibleengine_identity_agency_v3.test.name
  vpc_id      = flexibleengine_vpc_v1.test.id
  network_id  = flexibleengine_vpc_subnet_v1.test.id
  func_code   = "aW1wb3J0IGpzb24KZGVmIGhhbmRsZXIgKGZW50LCBjb250ZXh0KToKICAgIG91dHB1dCA9ICdIZWxsbyBtZXNzYWdlOiAnICsganNvbi5kdW1wcyhldmVudCkKICAgIHJldHVybiBvdXRwdXQ="
}`, testAccDmsKafka_config(rName, password), rName, OS_REGION_NAME, rName)
}

func testAccFunctionGraphKafkaTrigger_basic(rName, password string) string {
	return fmt.Sprintf(`
%s

resource "flexibleengine_fgs_trigger" "test" {
  function_urn = flexibleengine_fgs_function.test.urn
  type         = "KAFKA"

  kafka {
    instance_id = flexibleengine_dms_kafka_instance.test.id
    batch_size  = 100

    topic_ids = [
      flexibleengine_dms_kafka_topic.test.id
    ]
  }
}`, testAccFunctionGraphKafkaTrigger_base(rName, password))
}

func testAccFunctionGraphKafkaTrigger_update(rName, password string) string {
	return fmt.Sprintf(`
%s

resource "flexibleengine_fgs_trigger" "test" {
  function_urn = flexibleengine_fgs_function.test.urn
  type         = "KAFKA"
  status       = "DISABLED"
	
  kafka {
    instance_id = flexibleengine_dms_kafka_instance.test.id
    batch_size  = 100
  
    topic_ids = [
      flexibleengine_dms_kafka_topic.test.id
    ]
  }
}`, testAccFunctionGraphKafkaTrigger_base(rName, password))
}
