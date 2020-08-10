package flexibleengine

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/openstack/sfs/v2/shares"
)

func TestAccSFSAccessRuleV2_basic(t *testing.T) {
	var rule shares.AccessRight
	shareName := fmt.Sprintf("sfs-acc-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSFSAccessRuleV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: configAccSFSAccessRuleV2_basic(shareName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSFSAccessRuleV2Exists("flexibleengine_sfs_access_rule_v2.rule_1", &rule),
					resource.TestCheckResourceAttr(
						"flexibleengine_sfs_access_rule_v2.rule_1", "access_to", OS_VPC_ID),
					resource.TestCheckResourceAttr(
						"flexibleengine_sfs_access_rule_v2.rule_1", "access_level", "rw"),
					resource.TestCheckResourceAttr(
						"flexibleengine_sfs_access_rule_v2.rule_1", "status", "active"),
				),
			},
			{
				Config: configAccSFSAccessRuleV2_ipAuth(shareName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSFSAccessRuleV2Exists("flexibleengine_sfs_access_rule_v2.rule_1", &rule),
					resource.TestCheckResourceAttr(
						"flexibleengine_sfs_access_rule_v2.rule_1", "access_to",
						strings.Join([]string{OS_VPC_ID, "192.168.10.0/24", "0", "no_all_squash,no_root_squash"}, "#")),
					resource.TestCheckResourceAttr(
						"flexibleengine_sfs_access_rule_v2.rule_1", "status", "active"),
				),
			},
		},
	})
}

func testAccCheckSFSAccessRuleV2Destroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	sfsClient, err := config.sfsV2Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating Flexibleengine sfs client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "flexibleengine_sfs_access_rule_v2" {
			continue
		}

		sfsID := rs.Primary.Attributes["sfs_id"]
		if sfsID == "" {
			return fmt.Errorf("No SFSID is set in flexibleengine_sfs_access_rule_v2")
		}
		rules, err := shares.ListAccessRights(sfsClient, sfsID).ExtractAccessRights()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return nil
			}

			return err
		}

		for _, v := range rules {
			if v.ID == rs.Primary.ID {
				return fmt.Errorf("resource flexibleengine_sfs_access_rule_v2 still exists")
			}
		}
	}

	return nil
}

func testAccCheckSFSAccessRuleV2Exists(n string, rule *shares.AccessRight) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set in %s", n)
		}

		config := testAccProvider.Meta().(*Config)
		sfsClient, err := config.sfsV2Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating Flexibleengine sfs client: %s", err)
		}

		sfsID := rs.Primary.Attributes["sfs_id"]
		if sfsID == "" {
			return fmt.Errorf("No SFSID is set in %s", n)
		}

		rules, err := shares.ListAccessRights(sfsClient, sfsID).ExtractAccessRights()
		if err != nil {
			return err
		}

		for _, v := range rules {
			if v.ID == rs.Primary.ID {
				*rule = v
				return nil
			}
		}

		return fmt.Errorf("sfs access rule %s was not found", n)
	}
}

func configAccSFSAccessRuleV2_basic(sfsName string) string {
	return fmt.Sprintf(`
resource "flexibleengine_sfs_file_system_v2" "sfs_1" {
  share_proto = "NFS"
  size        = 10
  name        = "%s"
  description = "sfs file system created by terraform testacc"
}

resource "flexibleengine_sfs_access_rule_v2" "rule_1" {
  sfs_id = flexibleengine_sfs_file_system_v2.sfs_1.id
  access_to = "%s"
}`, sfsName, OS_VPC_ID)
}

func configAccSFSAccessRuleV2_ipAuth(sfsName string) string {
	return fmt.Sprintf(`
resource "flexibleengine_sfs_file_system_v2" "sfs_1" {
  share_proto = "NFS"
  size        = 10
  name        = "%s"
  description = "sfs file system created by terraform testacc"
}

resource "flexibleengine_sfs_access_rule_v2" "rule_1" {
  sfs_id = flexibleengine_sfs_file_system_v2.sfs_1.id
  access_to = join("#", ["%s", "192.168.10.0/24", "0", "no_all_squash,no_root_squash"])
}`, sfsName, OS_VPC_ID)
}
