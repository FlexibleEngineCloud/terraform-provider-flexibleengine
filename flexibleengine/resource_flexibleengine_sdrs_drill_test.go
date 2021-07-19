package flexibleengine

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/huaweicloud/golangsdk/openstack/sdrs/v1/drill"
)

func TestAccSdrsDrillV1_basic(t *testing.T) {
	var repDrill drill.Drill

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckSdrs(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSdrsDrillV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSdrsDrillV1_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSdrsDrillV1Exists("flexibleengine_sdrs_drill_v1.drill_1", &repDrill),
					resource.TestCheckResourceAttr(
						"flexibleengine_sdrs_drill_v1.drill_1", "name", "drill_1"),
					resource.TestCheckResourceAttr(
						"flexibleengine_sdrs_drill_v1.drill_1", "status", "available"),
				),
			},
			{
				Config: testAccSdrsDrillV1_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSdrsDrillV1Exists("flexibleengine_sdrs_drill_v1.drill_1", &repDrill),
					resource.TestCheckResourceAttr(
						"flexibleengine_sdrs_drill_v1.drill_1", "name", "drill_updated"),
				),
			},
		},
	})
}

func testAccCheckSdrsDrillV1Destroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	sdrsClient, err := config.sdrsV1Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine SDRS client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "flexibleengine_sdrs_drill_v1" {
			continue
		}

		_, err := drill.Get(sdrsClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("SDRS DR drill still exists")
		}
	}

	return nil
}

func testAccCheckSdrsDrillV1Exists(n string, drdrill *drill.Drill) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		sdrsClient, err := config.sdrsV1Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating FlexibleEngine SDRS client: %s", err)
		}

		found, err := drill.Get(sdrsClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.Id != rs.Primary.ID {
			return fmt.Errorf("SDRS DR drill not found")
		}

		*drdrill = *found
		return nil
	}
}

var testAccSdrsDrillV1_basic = fmt.Sprintf(`
data "flexibleengine_sdrs_domain_v1" "domain_1" {
  name = "SDRS_HypeDomain01"
}

# create vpc and subnet for drill
resource "flexibleengine_vpc_v1" "vpc_drill" {
  name = "vpc_drill"
  cidr = "192.168.0.0/16"
}
resource "flexibleengine_vpc_subnet_v1" "subnet_1"{
  name = "subnet_1"
  cidr = "192.168.0.0/24"
  gateway_ip = "192.168.0.1"
  vpc_id = flexibleengine_vpc_v1.vpc_drill.id
}

# create ecs server for protection 
resource "flexibleengine_compute_instance_v2" "server_1" {
  name = "server_1"
  security_groups = ["default"]
  availability_zone = "eu-west-0a"
  network {
    uuid = "%s"
  }
}

resource "flexibleengine_sdrs_protectiongroup_v1" "group_1" {
  name = "group_1"
  description = "test description"
  source_availability_zone = "eu-west-0a"
  target_availability_zone = "eu-west-0b"
  domain_id = data.flexibleengine_sdrs_domain_v1.domain_1.id
  source_vpc_id = "%s"
  dr_type = "migration"
}
resource "flexibleengine_sdrs_protectedinstance_v1" "instance_1" {
  name        = "instance_1"
  description = "test description"
  group_id    = flexibleengine_sdrs_protectiongroup_v1.group_1.id
  server_id   = flexibleengine_compute_instance_v2.server_1.id
  delete_target_server = true
  delete_target_eip = true
}
resource "flexibleengine_sdrs_drill_v1" "drill_1" {
  name = "drill_1"
  group_id = flexibleengine_sdrs_protectiongroup_v1.group_1.id
  drill_vpc_id = flexibleengine_vpc_v1.vpc_drill.id

  depends_on = [
    flexibleengine_vpc_subnet_v1.subnet_1,
    flexibleengine_sdrs_protectedinstance_v1.instance_1,
  ]
}`, OS_NETWORK_ID, OS_VPC_ID)

var testAccSdrsDrillV1_update = fmt.Sprintf(`
data "flexibleengine_sdrs_domain_v1" "domain_1" {
  name = "SDRS_HypeDomain01"
}

# create vpc and subnet for drill
resource "flexibleengine_vpc_v1" "vpc_drill" {
  name = "vpc_drill"
  cidr = "192.168.0.0/16"
}
resource "flexibleengine_vpc_subnet_v1" "subnet_1"{
  name = "subnet_1"
  cidr = "192.168.0.0/24"
  gateway_ip = "192.168.0.1"
  vpc_id = flexibleengine_vpc_v1.vpc_drill.id
}

# create ecs server for protection 
resource "flexibleengine_compute_instance_v2" "server_1" {
  name = "server_1"
  security_groups = ["default"]
  availability_zone = "eu-west-0a"
  network {
    uuid = "%s"
  }
}

resource "flexibleengine_sdrs_protectiongroup_v1" "group_1" {
  name = "group_1"
  description = "test description"
  source_availability_zone = "eu-west-0a"
  target_availability_zone = "eu-west-0b"
  domain_id = data.flexibleengine_sdrs_domain_v1.domain_1.id
  source_vpc_id = "%s"
  dr_type = "migration"
}
resource "flexibleengine_sdrs_protectedinstance_v1" "instance_1" {
  name        = "instance_1"
  description = "test description"
  group_id    = flexibleengine_sdrs_protectiongroup_v1.group_1.id
  server_id   = flexibleengine_compute_instance_v2.server_1.id
  delete_target_server = true
  delete_target_eip = true
}

resource "flexibleengine_sdrs_drill_v1" "drill_1" {
  name = "drill_updated"
  group_id = flexibleengine_sdrs_protectiongroup_v1.group_1.id
  drill_vpc_id = flexibleengine_vpc_v1.vpc_drill.id

  depends_on = [
    flexibleengine_vpc_subnet_v1.subnet_1,
    flexibleengine_sdrs_protectedinstance_v1.instance_1,
  ]
}`, OS_NETWORK_ID, OS_VPC_ID)
