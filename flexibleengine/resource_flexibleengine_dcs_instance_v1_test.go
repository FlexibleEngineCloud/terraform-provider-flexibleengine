package flexibleengine

import (
	"fmt"
	"testing"

	"github.com/chnsz/golangsdk/openstack/dcs/v1/instances"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
)

func TestAccDcsInstancesV1_basic(t *testing.T) {
	var instance instances.Instance
	var randName = fmt.Sprintf("acc_test_%s", acctest.RandString(5))
	resourceName := "flexibleengine_dcs_instance_v1.instance_1"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDcsV1InstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDcsV1Instance_basic(randName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDcsV1InstanceExists(resourceName, instance),
					resource.TestCheckResourceAttr(resourceName, "name", randName),
					resource.TestCheckResourceAttr(resourceName, "engine", "Redis"),
					resource.TestCheckResourceAttr(resourceName, "status", "RUNNING"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"password", "backup_type", "backup_at", "begin_at", "save_days", "period_type",
				},
			},
		},
	})
}

func TestAccDcsInstancesV1_redisV5(t *testing.T) {
	var instance instances.Instance
	var randName = fmt.Sprintf("acc_test_%s", acctest.RandString(5))
	resourceName := "flexibleengine_dcs_instance_v1.instance_1"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDcsV1InstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDcsV1Instance_redisV5(randName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDcsV1InstanceExists(resourceName, instance),
					resource.TestCheckResourceAttr(resourceName, "name", randName),
					resource.TestCheckResourceAttr(resourceName, "engine", "Redis"),
					resource.TestCheckResourceAttr(resourceName, "status", "RUNNING"),
					resource.TestCheckResourceAttr(resourceName, "port", "9999"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"password", "backup_type", "backup_at", "begin_at", "save_days", "period_type",
				},
			},
		},
	})
}

func TestAccDcsInstancesV1_memcached(t *testing.T) {
	var instance instances.Instance
	var randName = fmt.Sprintf("acc_test_%s", acctest.RandString(5))
	resourceName := "flexibleengine_dcs_instance_v1.instance_1"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDcsV1InstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDcsV1Instance_memcached(randName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDcsV1InstanceExists(resourceName, instance),
					resource.TestCheckResourceAttr(resourceName, "name", randName),
					resource.TestCheckResourceAttr(resourceName, "engine", "Memcached"),
					resource.TestCheckResourceAttr(resourceName, "status", "RUNNING"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"password", "backup_type", "backup_at", "begin_at", "save_days", "period_type",
				},
			},
		},
	})
}

func testAccCheckDcsV1InstanceDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	dcsClient, err := config.DcsV1Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating Flexibleengine instance client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "flexibleengine_dcs_instance_v1" {
			continue
		}

		_, err := instances.Get(dcsClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("The Dcs instance still exists")
		}
	}
	return nil
}

func testAccCheckDcsV1InstanceExists(n string, instance instances.Instance) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		dcsClient, err := config.DcsV1Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating Flexibleengine instance client: %s", err)
		}

		v, err := instances.Get(dcsClient, rs.Primary.ID).Extract()
		if err != nil {
			return fmt.Errorf("Error getting Flexibleengine instance: %s, err: %s", rs.Primary.ID, err)
		}

		if v.InstanceID != rs.Primary.ID {
			return fmt.Errorf("The DCS instance not found")
		}
		instance = *v
		return nil
	}
}

func testAccDcsV1Instance_network(rName string) string {
	return fmt.Sprintf(`
resource "flexibleengine_vpc_v1" "vpc_1" {
  name = "%[1]s"
  cidr = "192.168.0.0/16"
}

resource "flexibleengine_vpc_subnet_v1" "subnet_1" {
  name       = "%[1]s"
  cidr       = "192.168.0.0/24"
  gateway_ip = "192.168.0.1"
  vpc_id     = flexibleengine_vpc_v1.vpc_1.id
}

resource "flexibleengine_networking_secgroup_v2" "secgroup_1" {
  name        = "%[1]s"
  description = "secgroup_1"
}
`, rName)
}

func testAccDcsV1Instance_basic(rName string) string {
	return fmt.Sprintf(`
%s

resource "flexibleengine_dcs_instance_v1" "instance_1" {
  name              = "%s"
  engine            = "Redis"
  engine_version    = "3.0"
  password          = "Huawei_test"
  product_id        = "dcs.master_standby-h"
  capacity          = 2
  vpc_id            = flexibleengine_vpc_v1.vpc_1.id
  network_id        = flexibleengine_vpc_subnet_v1.subnet_1.id
  security_group_id = flexibleengine_networking_secgroup_v2.secgroup_1.id
  available_zones   = ["eu-west-0a"]

  save_days   = 1
  backup_type = "manual"
  begin_at    = "00:00-01:00"
  period_type = "weekly"
  backup_at   = [1]
}
`, testAccDcsV1Instance_network(rName), rName)
}

func testAccDcsV1Instance_redisV5(rName string) string {
	return fmt.Sprintf(`
data "flexibleengine_dcs_product_v1" "product_ha" {
  engine         = "Redis"
  engine_version = "4.0;5.0"
  cache_mode     = "ha"
  capacity       = 0.5
  replica_count  = 2
}

%s

resource "flexibleengine_dcs_instance_v1" "instance_1" {
  name            = "%s"
  engine          = "Redis"
  engine_version  = "5.0"
  password        = "Huawei_test"
  product_id      = data.flexibleengine_dcs_product_v1.product_ha.id
  capacity        = 0.5
  vpc_id          = flexibleengine_vpc_v1.vpc_1.id
  network_id      = flexibleengine_vpc_subnet_v1.subnet_1.id
  available_zones = ["eu-west-0a", "eu-west-0c"]
  port            = "9999"

  save_days   = 1
  backup_type = "manual"
  begin_at    = "00:00-01:00"
  period_type = "weekly"
  backup_at   = [1]
}
`, testAccDcsV1Instance_network(rName), rName)
}

func testAccDcsV1Instance_memcached(rName string) string {
	return fmt.Sprintf(`
%s

resource "flexibleengine_dcs_instance_v1" "instance_1" {
  name              = "%s"
  engine            = "Memcached"
  access_user       = "admin"
  password          = "Huawei_test"
  product_id        = "dcs.memcached.master_standby-h"
  capacity          = 2
  vpc_id            = flexibleengine_vpc_v1.vpc_1.id
  network_id        = flexibleengine_vpc_subnet_v1.subnet_1.id
  security_group_id = flexibleengine_networking_secgroup_v2.secgroup_1.id
  available_zones   = ["eu-west-0a"]

  save_days   = 1
  backup_type = "manual"
  begin_at    = "00:00-01:00"
  period_type = "weekly"
  backup_at   = [1]
}
`, testAccDcsV1Instance_network(rName), rName)
}
