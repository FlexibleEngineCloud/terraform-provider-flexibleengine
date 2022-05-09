package flexibleengine

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/networking/v2/extensions/security/groups"
)

func TestAccNetworkingV2SecGroup_basic(t *testing.T) {
	var secGroup groups.SecGroup
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	updateName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "flexibleengine_networking_secgroup_v2.secgroup_1"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkingV2SecGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2SecGroup_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2SecGroupExists(resourceName, &secGroup),
					testAccCheckNetworkingV2SecGroupRuleCount(&secGroup, 2),
					resource.TestCheckResourceAttr(resourceName, "name", "sg-"+rName),
				),
			},
			{
				Config: testAccNetworkingV2SecGroup_update(updateName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPtr(resourceName, "id", &secGroup.ID),
					resource.TestCheckResourceAttr(resourceName, "name", "sg-"+updateName),
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

func TestAccNetworkingV2SecGroup_noDefaultRules(t *testing.T) {
	var secGroup groups.SecGroup
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "flexibleengine_networking_secgroup_v2.secgroup_1"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkingV2SecGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2SecGroup_noDefaultRules(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2SecGroupExists(resourceName, &secGroup),
					testAccCheckNetworkingV2SecGroupRuleCount(&secGroup, 0),
					resource.TestCheckResourceAttr(resourceName, "name", "sg-"+rName),
				),
			},
		},
	})
}

func testAccCheckNetworkingV2SecGroupDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	networkingClient, err := config.NetworkingV2Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine networking client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "flexibleengine_networking_secgroup_v2" {
			continue
		}

		_, err := groups.Get(networkingClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Security group still exists")
		}
	}

	return nil
}

func testAccCheckNetworkingV2SecGroupExists(n string, secGroup *groups.SecGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		networkingClient, err := config.NetworkingV2Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating FlexibleEngine networking client: %s", err)
		}

		found, err := groups.Get(networkingClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Security group not found")
		}

		*secGroup = *found

		return nil
	}
}

func testAccCheckNetworkingV2SecGroupRuleCount(
	sg *groups.SecGroup, count int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if len(sg.Rules) == count {
			return nil
		}

		return fmt.Errorf("Unexpected number of rules in group %s. Expected %d, got %d",
			sg.ID, count, len(sg.Rules))
	}
}

func testAccNetworkingV2SecGroup_basic(rName string) string {
	return fmt.Sprintf(`
resource "flexibleengine_networking_secgroup_v2" "secgroup_1" {
  name        = "sg-%s"
  description = "terraform security group acceptance test"
}
`, rName)
}

func testAccNetworkingV2SecGroup_update(rName string) string {
	return fmt.Sprintf(`
resource "flexibleengine_networking_secgroup_v2" "secgroup_1" {
  name        = "sg-%s"
  description = "terraform security group acceptance test"
}
`, rName)
}

func testAccNetworkingV2SecGroup_noDefaultRules(rName string) string {
	return fmt.Sprintf(`
resource "flexibleengine_networking_secgroup_v2" "secgroup_1" {
  name                 = "sg-%s"
  description          = "terraform security group acceptance test"
  delete_default_rules = true
}
`, rName)
}
