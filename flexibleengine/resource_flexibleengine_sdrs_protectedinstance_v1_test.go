package flexibleengine

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/golangsdk/openstack/sdrs/v1/protectedinstances"
)

func TestAccSdrsProtectedInstanceV1_basic(t *testing.T) {
	var instance protectedinstances.Instance

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckSdrs(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSdrsProtectedInstanceV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSdrsProtectedInstanceV1_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSdrsProtectedInstanceV1Exists("flexibleengine_sdrs_protectedinstance_v1.instance_1", &instance),
					resource.TestCheckResourceAttr(
						"flexibleengine_sdrs_protectedinstance_v1.instance_1", "name", "instance_1"),
				),
			},
			{
				Config: testAccSdrsProtectedInstanceV1_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSdrsProtectedInstanceV1Exists("flexibleengine_sdrs_protectedinstance_v1.instance_1", &instance),
					resource.TestCheckResourceAttr(
						"flexibleengine_sdrs_protectedinstance_v1.instance_1", "name", "instance_updated"),
				),
			},
		},
	})
}

func testAccCheckSdrsProtectedInstanceV1Destroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	sdrsClient, err := config.sdrsV1Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine SDRS client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "flexibleengine_sdrs_protectedinstance_v1" {
			continue
		}

		_, err := protectedinstances.Get(sdrsClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("SDRS protectedinstance still exists")
		}
	}

	return nil
}

func testAccCheckSdrsProtectedInstanceV1Exists(n string, instance *protectedinstances.Instance) resource.TestCheckFunc {
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

		found, err := protectedinstances.Get(sdrsClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.Id != rs.Primary.ID {
			return fmt.Errorf("SDRS protectedinstance not found")
		}

		*instance = *found

		return nil
	}
}

var testAccSdrsProtectedInstanceV1_basic = fmt.Sprintf(`
data "flexibleengine_sdrs_domain_v1" "domain_1" {
	name = "SDRS_HypeDomain01"
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

resource "flexibleengine_compute_instance_v2" "instance_1" {
  name = "instance_1"
  security_groups = ["default"]
  availability_zone = "eu-west-0a"
  network {
    uuid = "%s"
  }
}

resource "flexibleengine_sdrs_protectedinstance_v1" "instance_1" {
	group_id = flexibleengine_sdrs_protectiongroup_v1.group_1.id
	server_id = flexibleengine_compute_instance_v2.instance_1.id
	name = "instance_1"
	description = "test description"
    delete_target_server = true
    delete_target_eip = true
}`, OS_VPC_ID, OS_NETWORK_ID)

var testAccSdrsProtectedInstanceV1_update = fmt.Sprintf(`
data "flexibleengine_sdrs_domain_v1" "domain_1" {
	name = "SDRS_HypeDomain01"
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

resource "flexibleengine_compute_instance_v2" "instance_1" {
  name = "instance_1"
  security_groups = ["default"]
  availability_zone = "eu-west-0a"
  network {
    uuid = "%s"
  }
}

resource "flexibleengine_sdrs_protectedinstance_v1" "instance_1" {
	group_id = flexibleengine_sdrs_protectiongroup_v1.group_1.id
	server_id = flexibleengine_compute_instance_v2.instance_1.id
	name = "instance_updated"
	description = "test description"
    delete_target_server = true
    delete_target_eip = true
}`, OS_VPC_ID, OS_NETWORK_ID)
