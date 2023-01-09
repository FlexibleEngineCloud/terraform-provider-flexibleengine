package flexibleengine

import (
	"fmt"
	"testing"

	"github.com/chnsz/golangsdk/openstack/dms/v1/instances"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDmsKafkaInstance_basic(t *testing.T) {
	var instanceName = fmt.Sprintf("ACCPTTEST-%s", acctest.RandString(5))
	var instanceUpdate = instanceName + "-update"
	resourceName := "flexibleengine_dms_kafka_instance.instance_1"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: TestAccProviderFactories,
		CheckDestroy:      testAccCheckDmsKafkaInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDmsKafkaInstance_basic(instanceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDmsKafkaInstanceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", instanceName),
					resource.TestCheckResourceAttr(resourceName, "engine", "kafka"),
					resource.TestCheckResourceAttr(resourceName, "engine_version", "2.3.0"),
					resource.TestCheckResourceAttr(resourceName, "engine_type", "cluster"),
					resource.TestCheckResourceAttr(resourceName, "status", "RUNNING"),
					resource.TestCheckResourceAttr(resourceName, "ssl_enable", "true"),
					resource.TestCheckResourceAttrSet(resourceName, "product_spec_code"),
					resource.TestCheckResourceAttrSet(resourceName, "manegement_connect_address"),
					resource.TestCheckResourceAttrSet(resourceName, "port"),
				),
			},
			{
				Config: testAccDmsKafkaInstance_update(instanceName, instanceUpdate),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", instanceUpdate),
					resource.TestCheckResourceAttr(resourceName, "description", "instance update description"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"access_user", "password", "manager_user", "manager_password",
				},
			},
		},
	})
}

func testAccCheckDmsKafkaInstanceDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	client, err := config.DmsV1Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating DMS client, err=%s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "flexibleengine_dms_kafka_instance" {
			continue
		}

		_, err := instances.Get(client, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("DMS kafka instance %s still exists", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckDmsKafkaInstanceExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		client, err := config.DmsV1Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating DMS client, err=%s", err)
		}

		_, err = instances.Get(client, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		return nil
	}
}

func testAccDmsKafkaInstance_base(resName string) string {
	return fmt.Sprintf(`
data "flexibleengine_dms_product" "product_1" {
  bandwidth = "300MB"
}

data "flexibleengine_availability_zones" "zones" {}

resource "flexibleengine_vpc_v1" "vpc_1" {
  name = "%s"
  cidr = "192.168.0.0/16"
}

resource "flexibleengine_vpc_subnet_v1" "vpc_subnet_1" {
  name       = "%s"
  cidr       = "192.168.10.0/24"
  gateway_ip = "192.168.10.1"
  vpc_id     = flexibleengine_vpc_v1.vpc_1.id
}

resource "flexibleengine_networking_secgroup_v2" "secgroup_1" {
  name        = "%s"
  description = "secgroup for DMS"
}`, resName, resName, resName)
}

func testAccDmsKafkaInstance_basic(resName string) string {
	return fmt.Sprintf(`
%s

resource "flexibleengine_dms_kafka_instance" "instance_1" {
  name               = "%s"
  manager_user       = "admin"
  manager_password   = "Dmstest@123"
  access_user        = "user"
  password           = "Dmstest@123"
  vpc_id             = flexibleengine_vpc_v1.vpc_1.id
  network_id         = flexibleengine_vpc_subnet_v1.vpc_subnet_1.id
  security_group_id  = flexibleengine_networking_secgroup_v2.secgroup_1.id
  availability_zones = data.flexibleengine_availability_zones.zones.names
  bandwidth          = data.flexibleengine_dms_product.product_1.bandwidth
  product_id         = data.flexibleengine_dms_product.product_1.id
  storage_space      = data.flexibleengine_dms_product.product_1.storage_space
  engine_version     = data.flexibleengine_dms_product.product_1.engine_version
}`, testAccDmsKafkaInstance_base(resName), resName)
}

func testAccDmsKafkaInstance_update(resName, resUpdate string) string {
	return fmt.Sprintf(`
%s

resource "flexibleengine_dms_kafka_instance" "instance_1" {
  name               = "%s"
  description        = "instance update description"
  manager_user       = "admin"
  manager_password   = "Dmstest@123"
  access_user        = "user"
  password           = "Dmstest@123"
  vpc_id             = flexibleengine_vpc_v1.vpc_1.id
  network_id         = flexibleengine_vpc_subnet_v1.vpc_subnet_1.id
  security_group_id  = flexibleengine_networking_secgroup_v2.secgroup_1.id
  availability_zones = data.flexibleengine_dms_product.product_1.availability_zones
  bandwidth          = data.flexibleengine_dms_product.product_1.bandwidth
  product_id         = data.flexibleengine_dms_product.product_1.id
  storage_space      = data.flexibleengine_dms_product.product_1.storage_space
  engine_version     = data.flexibleengine_dms_product.product_1.engine_version
  }`, testAccDmsKafkaInstance_base(resName), resUpdate)
}
