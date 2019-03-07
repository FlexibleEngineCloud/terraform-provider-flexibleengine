package flexibleengine

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"

	"testing"

	"github.com/huaweicloud/golangsdk/openstack/vbs/v2/policies"
)

func TestAccVBSBackupPolicyV2_basic(t *testing.T) {
	var policy policies.Policy

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckRequiredEnvVars(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccVBSBackupPolicyV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVBSBackupPolicyV2_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccVBSBackupPolicyV2Exists("flexibleengine_vbs_backup_policy_v2.vbs", &policy),
					resource.TestCheckResourceAttr(
						"flexibleengine_vbs_backup_policy_v2.vbs", "name", "policy_001"),
					resource.TestCheckResourceAttr(
						"flexibleengine_vbs_backup_policy_v2.vbs", "status", "ON"),
				),
			},
			{
				Config: testAccVBSBackupPolicyV2_update,
				Check: resource.ComposeTestCheckFunc(
					testAccVBSBackupPolicyV2Exists("flexibleengine_vbs_backup_policy_v2.vbs", &policy),
					resource.TestCheckResourceAttr(
						"flexibleengine_vbs_backup_policy_v2.vbs", "name", "policy_002"),
					resource.TestCheckResourceAttr(
						"flexibleengine_vbs_backup_policy_v2.vbs", "status", "ON"),
				),
			},
		},
	})
}

func testAccVBSBackupPolicyV2Destroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	vbsClient, err := config.vbsV2Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating sfs client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "flexibleengine_vbs_backup_policy_v2" {
			continue
		}

		_, err := policies.List(vbsClient, policies.ListOpts{ID: rs.Primary.ID})
		if err != nil {
			return fmt.Errorf("Backup Policy still exists")
		}
	}

	return nil
}

func testAccVBSBackupPolicyV2Exists(n string, policy *policies.Policy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		vbsClient, err := config.vbsV2Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating vbs client: %s", err)
		}

		policyList, err := policies.List(vbsClient, policies.ListOpts{ID: rs.Primary.ID})
		if err != nil {
			return err
		}
		found := policyList[0]
		if found.ID != rs.Primary.ID {
			return fmt.Errorf("backup policy not found")
		}

		*policy = found

		return nil
	}
}

var testAccVBSBackupPolicyV2_basic = fmt.Sprintf(`
resource "flexibleengine_vbs_backup_policy_v2" "vbs" {
  name = "policy_001"
  start_time  = "12:00"
  status  = "ON"
  retain_first_backup = "N"
  rentention_num = 2
  frequency = 1     
}
`)

var testAccVBSBackupPolicyV2_update = fmt.Sprintf(`
resource "flexibleengine_vbs_backup_policy_v2" "vbs" {
  name = "policy_002"
  start_time  = "12:00"
  status  = "ON"
  retain_first_backup = "N"
  rentention_num = 2
  frequency = 1      
}
`)
