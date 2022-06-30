package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/elb/v3/ipgroups"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccElbV3IpGroup_basic(t *testing.T) {
	var c ipgroups.IpGroup
	name := fmt.Sprintf("tf-acc-%s", acctest.RandString(5))
	resourceName := "flexibleengine_elb_ipgroup.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: TestAccProviderFactories,
		CheckDestroy:      testAccCheckElbV3IpGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccElbV3IpGroupConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckElbV3IpGroupExists(resourceName, &c),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "description", "terraform test"),
					resource.TestCheckResourceAttr(resourceName, "ip_list.#", "1"),
				),
			},
			{
				Config: testAccElbV3IpGroupConfig_update(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("%s_updated", name)),
					resource.TestCheckResourceAttr(resourceName, "description", "terraform test updated"),
					resource.TestCheckResourceAttr(resourceName, "ip_list.#", "2"),
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

func testAccCheckElbV3IpGroupDestroy(s *terraform.State) error {
	conf := acceptance.TestAccProvider.Meta().(*config.Config)
	elbClient, err := conf.ElbV3Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating ELB v3 client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "flexibleengine_elb_ipgroup" {
			continue
		}

		_, err := ipgroups.Get(elbClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("IpGroup still exists: %s", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckElbV3IpGroupExists(
	n string, c *ipgroups.IpGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		conf := acceptance.TestAccProvider.Meta().(*config.Config)
		elbClient, err := conf.ElbV3Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating ELB v3 client: %s", err)
		}

		found, err := ipgroups.Get(elbClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("IpGroup not found")
		}

		*c = *found

		return nil
	}
}

func testAccElbV3IpGroupConfig_basic(name string) string {
	return fmt.Sprintf(`
resource "flexibleengine_elb_ipgroup" "test"{
  name        = "%s"
  description = "terraform test"

  ip_list {
    ip = "192.168.10.10"
    description = "ECS01"
  }
}
`, name)
}

func testAccElbV3IpGroupConfig_update(name string) string {
	return fmt.Sprintf(`
resource "flexibleengine_elb_ipgroup" "test"{
  name        = "%s_updated"
  description = "terraform test updated"

  ip_list {
    ip          = "192.168.10.10"
    description = "ECS01"
  }

  ip_list {
    ip          = "192.168.10.11"
    description = "ECS02"
  }
}
`, name)
}
