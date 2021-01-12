package flexibleengine

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/huaweicloud/golangsdk/openstack/networking/v1/vpcs"
)

func TestAccFlexibleEngineVpcV1_basic(t *testing.T) {
	var vpc vpcs.Vpc

	resourceName := "flexibleengine_vpc_v1.vpc_1"
	rName := fmt.Sprintf("vpc-acc-test-%s", acctest.RandString(5))
	rNameUpdate := rName + "-updated"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckFlexibleEngineVpcV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcV1_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFlexibleEngineVpcV1Exists(resourceName, &vpc),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "cidr", "192.168.0.0/16"),
					resource.TestCheckResourceAttr(resourceName, "status", "OK"),
					resource.TestCheckResourceAttr(resourceName, "shared", "false"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
				),
			},
			{
				Config: testAccVpcV1_update(rNameUpdate),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFlexibleEngineVpcV1Exists(resourceName, &vpc),
					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdate),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value_updated"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckFlexibleEngineVpcV1Destroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	vpcClient, err := config.networkingV1Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine vpc client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "flexibleengine_vpc_v1" {
			continue
		}

		_, err := vpcs.Get(vpcClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Vpc still exists")
		}
	}

	return nil
}

func testAccCheckFlexibleEngineVpcV1Exists(n string, vpc *vpcs.Vpc) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		vpcClient, err := config.networkingV1Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating FlexibleEngine vpc client: %s", err)
		}

		found, err := vpcs.Get(vpcClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("vpc not found")
		}

		*vpc = *found

		return nil
	}
}

func testAccVpcV1_basic(rName string) string {
	return fmt.Sprintf(`
resource "flexibleengine_vpc_v1" "vpc_1" {
  name = "%s"
  cidr="192.168.0.0/16"

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, rName)
}

func testAccVpcV1_update(rName string) string {
	return fmt.Sprintf(`
resource "flexibleengine_vpc_v1" "vpc_1" {
  name = "%s"
  cidr="192.168.0.0/16"

  tags = {
    foo = "bar"
    key = "value_updated"
  }
}
`, rName)
}
