package flexibleengine

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/huaweicloud/golangsdk/openstack/sfs/v2/shares"
)

func TestAccSFSFileSystemV2_basic(t *testing.T) {
	var share shares.Share
	resourceName := "flexibleengine_sfs_file_system_v2.sfs_1"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSFSFileSystemV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSFSFileSystemV2_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSFSFileSystemV2Exists(resourceName, &share),
					resource.TestCheckResourceAttr(resourceName, "name", "sfs-test1"),
					resource.TestCheckResourceAttr(resourceName, "share_proto", "NFS"),
					resource.TestCheckResourceAttr(resourceName, "status", "available"),
					resource.TestCheckResourceAttr(resourceName, "size", "10"),
					resource.TestCheckResourceAttr(resourceName, "access_level", "rw"),
				),
			},
			{
				Config: testAccSFSFileSystemV2_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSFSFileSystemV2Exists(resourceName, &share),
					resource.TestCheckResourceAttr(resourceName, "name", "sfs-test2"),
					resource.TestCheckResourceAttr(resourceName, "share_proto", "NFS"),
					resource.TestCheckResourceAttr(resourceName, "status", "available"),
					resource.TestCheckResourceAttr(resourceName, "size", "20"),
					resource.TestCheckResourceAttr(resourceName, "access_level", "rw"),
				),
			},
		},
	})
}

func TestAccSFSFileSystemV2_without_rule(t *testing.T) {
	var share shares.Share
	resourceName := "flexibleengine_sfs_file_system_v2.sfs_1"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSFSFileSystemV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSFSFileSystemV2_without_rule,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSFSFileSystemV2Exists(resourceName, &share),
					resource.TestCheckResourceAttr(resourceName, "name", "sfs-test-no-rules"),
					resource.TestCheckResourceAttr(resourceName, "share_proto", "NFS"),
					resource.TestCheckResourceAttr(resourceName, "status", "unavailable"),
					resource.TestCheckResourceAttr(resourceName, "size", "10"),
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

func testAccCheckSFSFileSystemV2Destroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	sfsClient, err := config.sfsV2Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating Flexibleengine sfs client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "flexibleengine_sfs_file_system_v2" {
			continue
		}

		_, err := shares.Get(sfsClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Share File still exists")
		}
	}

	return nil
}

func testAccCheckSFSFileSystemV2Exists(n string, share *shares.Share) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		sfsClient, err := config.sfsV2Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating Flexibleengine sfs client: %s", err)
		}

		found, err := shares.Get(sfsClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("share file not found")
		}

		*share = *found

		return nil
	}
}

var testAccSFSFileSystemV2_basic = fmt.Sprintf(`
resource "flexibleengine_sfs_file_system_v2" "sfs_1" {
  share_proto  = "NFS"
  size         = 10
  name         = "sfs-test1"
  description  = "sfs_c2c_test-file"
  access_to    = "%s"
  access_level = "rw"
  availability_zone = "%s"
}
`, OS_VPC_ID, OS_AVAILABILITY_ZONE)

var testAccSFSFileSystemV2_update = fmt.Sprintf(`
resource "flexibleengine_sfs_file_system_v2" "sfs_1" {
  share_proto = "NFS"
  size        = 20
  name        = "sfs-test2"
  description = "sfs_c2c_test-file"
  access_to   = "%s"
  access_type = "cert"
  availability_zone = "%s"
}
`, OS_VPC_ID, OS_AVAILABILITY_ZONE)

const testAccSFSFileSystemV2_without_rule = `
resource "flexibleengine_sfs_file_system_v2" "sfs_1" {
  share_proto = "NFS"
  size        = 10
  name        = "sfs-test-no-rules"
  description = "sfs_c2c_test-file"
}
`
