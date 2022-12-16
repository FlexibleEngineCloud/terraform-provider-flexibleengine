package flexibleengine

import (
	"fmt"
	"testing"

	"github.com/chnsz/golangsdk/openstack/lts/huawei/loggroups"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccLTSGroupV2_basic(t *testing.T) {
	var group loggroups.LogGroup
	groupName := fmt.Sprintf("acc-group-%s", acctest.RandString(5))
	resourceName := "flexibleengine_lts_group.testacc_group"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLTSGroupV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLTSGroupV2_basic(groupName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLTSGroupV2Exists(resourceName, &group),
					resource.TestCheckResourceAttr(resourceName, "group_name", groupName),
					resource.TestCheckResourceAttr(resourceName, "ttl_in_days", "7"),
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

func testAccCheckLTSGroupV2Destroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	ltsclient, err := config.LtsV2Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine LTS client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "flexibleengine_lts_group" {
			continue
		}

		groups, err := loggroups.List(ltsclient).Extract()
		if err != nil {
			return fmt.Errorf("Log group get list err: %s", err.Error())
		}
		for _, group := range groups.LogGroups {
			if group.ID == rs.Primary.ID {
				return fmt.Errorf("Log group (%s) still exists.", rs.Primary.ID)
			}
		}
	}
	return nil
}

func testAccCheckLTSGroupV2Exists(n string, group *loggroups.LogGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		ltsclient, err := config.LtsV2Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating FlexibleEngine LTS client: %s", err)
		}

		var founds *loggroups.LogGroups
		founds, err = loggroups.List(ltsclient).Extract()
		if err != nil {
			return err
		}
		for _, loggroup := range founds.LogGroups {
			if rs.Primary.ID == loggroup.ID {
				*group = loggroup
				return nil
			}
		}

		return fmt.Errorf("resource not found: lts group %s", rs.Primary.ID)
	}
}

func testAccLTSGroupV2_basic(name string) string {
	return fmt.Sprintf(`
resource "flexibleengine_lts_group" "testacc_group" {
  group_name = "%s"
}
`, name)
}
