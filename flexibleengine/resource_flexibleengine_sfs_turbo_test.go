package flexibleengine

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/huaweicloud/golangsdk/openstack/sfs_turbo/v1/shares"
)

func TestAccSFSTurbo_basic(t *testing.T) {
	randSuffix := acctest.RandString(5)
	turboName := fmt.Sprintf("sfs-turbo-acc-%s", randSuffix)
	resourceName := "flexibleengine_sfs_turbo.sfs-turbo1"
	var turbo shares.Turbo

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSFSTurboDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSFSTurbo_basic(randSuffix),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSFSTurboExists(resourceName, &turbo),
					resource.TestCheckResourceAttr(resourceName, "name", turboName),
					resource.TestCheckResourceAttr(resourceName, "share_proto", "NFS"),
					resource.TestCheckResourceAttr(resourceName, "share_type", "STANDARD"),
					resource.TestCheckResourceAttr(resourceName, "size", "500"),
					resource.TestCheckResourceAttr(resourceName, "status", "200"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccSFSTurbo_update(randSuffix),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSFSTurboExists(resourceName, &turbo),
					resource.TestCheckResourceAttr(resourceName, "size", "600"),
					resource.TestCheckResourceAttr(resourceName, "status", "221"),
				),
			},
		},
	})
}

func TestAccSFSTurbo_crypt(t *testing.T) {
	randSuffix := acctest.RandString(5)
	turboName := fmt.Sprintf("sfs-turbo-acc-%s", randSuffix)
	resourceName := "flexibleengine_sfs_turbo.sfs-turbo1"
	var turbo shares.Turbo

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSFSTurboDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSFSTurbo_crypt(randSuffix),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSFSTurboExists(resourceName, &turbo),
					resource.TestCheckResourceAttr(resourceName, "name", turboName),
					resource.TestCheckResourceAttr(resourceName, "share_proto", "NFS"),
					resource.TestCheckResourceAttr(resourceName, "share_type", "STANDARD"),
					resource.TestCheckResourceAttr(resourceName, "size", "500"),
					resource.TestCheckResourceAttr(resourceName, "status", "200"),
					resource.TestCheckResourceAttrSet(resourceName, "crypt_key_id"),
				),
			},
		},
	})
}

func testAccCheckSFSTurboDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	sfsClient, err := config.sfsV1Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine sfs turbo client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "flexibleengine_sfs_turbo" {
			continue
		}

		_, err := shares.Get(sfsClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("SFS Turbo still exists")
		}
	}

	return nil
}

func testAccCheckSFSTurboExists(n string, share *shares.Turbo) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		sfsClient, err := config.sfsV1Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating FlexibleEngine sfs turbo client: %s", err)
		}

		found, err := shares.Get(sfsClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("sfs turbo not found")
		}

		*share = *found
		return nil
	}
}

func testAccNetworkPreConditions(suffix string) string {
	return fmt.Sprintf(`
resource "flexibleengine_vpc_v1" "test" {
  name = "tf-acc-vpc-%s"
  cidr = "192.168.0.0/16"
}

resource "flexibleengine_vpc_subnet_v1" "test" {
  name       = "tf-acc-subnet-%s"
  cidr       = "192.168.0.0/16"
  gateway_ip = "192.168.0.1"
  vpc_id     = flexibleengine_vpc_v1.test.id
}

resource "flexibleengine_networking_secgroup_v2" "secgroup" {
  name        = "tf-acc-sg-%s"
  description = "terraform security group for sfs turbo acceptance test"
}
`, suffix, suffix, suffix)
}

func testAccSFSTurbo_basic(suffix string) string {
	return fmt.Sprintf(`
%s

resource "flexibleengine_sfs_turbo" "sfs-turbo1" {
  name        = "sfs-turbo-acc-%s"
  size        = 500
  share_proto = "NFS"
  vpc_id      = flexibleengine_vpc_v1.test.id
  subnet_id   = flexibleengine_vpc_subnet_v1.test.id
  security_group_id = flexibleengine_networking_secgroup_v2.secgroup.id
  availability_zone = "%s"
}
`, testAccNetworkPreConditions(suffix), suffix, OS_AVAILABILITY_ZONE)
}

func testAccSFSTurbo_update(suffix string) string {
	return fmt.Sprintf(`
%s

resource "flexibleengine_sfs_turbo" "sfs-turbo1" {
  name        = "sfs-turbo-acc-%s"
  size        = 600
  share_proto = "NFS"
  vpc_id      = flexibleengine_vpc_v1.test.id
  subnet_id   = flexibleengine_vpc_subnet_v1.test.id
  security_group_id = flexibleengine_networking_secgroup_v2.secgroup.id
  availability_zone = "%s"
}
`, testAccNetworkPreConditions(suffix), suffix, OS_AVAILABILITY_ZONE)
}

func testAccSFSTurbo_crypt(suffix string) string {
	return fmt.Sprintf(`
%s
resource "flexibleengine_kms_key_v1" "key_1" {
  key_alias = "kms-acc-%s"
  pending_days = "7"
}

resource "flexibleengine_sfs_turbo" "sfs-turbo1" {
  name         = "sfs-turbo-acc-%s"
  size         = 500
  share_proto  = "NFS"
  crypt_key_id = flexibleengine_kms_key_v1.key_1.id
  vpc_id       = flexibleengine_vpc_v1.test.id
  subnet_id    = flexibleengine_vpc_subnet_v1.test.id
  security_group_id = flexibleengine_networking_secgroup_v2.secgroup.id
  availability_zone = "%s"
}
`, testAccNetworkPreConditions(suffix), suffix, suffix, OS_AVAILABILITY_ZONE)
}
