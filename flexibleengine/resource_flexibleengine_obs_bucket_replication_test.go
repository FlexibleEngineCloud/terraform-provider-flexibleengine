package flexibleengine

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccObsBucketReplication_basic(t *testing.T) {
	rName := acctest.RandString(4)
	resourceName := "flexibleengine_obs_bucket_replication.replica"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckOBSReplication(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckObsBucketReplicationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccObsBucketReplication_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckObsBucketReplicationExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "destination_bucket", OS_DESTINATION_BUCKET),
					resource.TestCheckResourceAttr(resourceName, "rule.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "rule.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "rule.0.prefix", "abc"),
					resource.TestCheckResourceAttrSet(resourceName, "rule.0.id"),
				),
			},
			{
				Config: testAccObsBucketReplication_update(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "rule.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "rule.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "rule.1.enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "rule.1.prefix", "terraform"),
					resource.TestCheckResourceAttr(resourceName, "rule.1.storage_class", "COLD"),
					resource.TestCheckResourceAttrSet(resourceName, "rule.1.id"),
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

func testAccCheckObsBucketReplicationDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	obsClient, err := config.objectStorageClientWithSignature(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine OBS client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "flexibleengine_obs_bucket_replication" {
			continue
		}

		_, err := obsClient.GetBucketReplication(rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("OBS Bucket %s Replication configuration still exists", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckObsBucketReplicationExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		obsClient, err := config.objectStorageClientWithSignature(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating FlexibleEngine OBS client: %s", err)
		}

		_, err = obsClient.GetBucketReplication(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("OBS Bucket Replication configuration not found: %v", err)
		}
		return nil
	}
}

func testAccObsBucketReplication_base(rName string) string {
	return fmt.Sprintf(`
resource "flexibleengine_identity_agency_v3" "agency_1" {
  name                   = "agency-%s"
  description            = "This is a iam agency for obs bucket replication"
  duration               = "ONEDAY"
  delegated_service_name = "op_svc_obs"

  domain_roles = [
    "OBS FullAccess",
  ]
}

resource "flexibleengine_obs_bucket" "source" {
  bucket        = "tf-test-bucket-%s"
  storage_class = "STANDARD"
  acl           = "private"
}
`, rName, rName)
}

func testAccObsBucketReplication_basic(rName string) string {
	return fmt.Sprintf(`
%s

resource "flexibleengine_obs_bucket_replication" "replica" {
  bucket             = flexibleengine_obs_bucket.source.id
  destination_bucket = "%s"
  agency             = flexibleengine_identity_agency_v3.agency_1.name

  rule {
    prefix = "abc"
  }
}
`, testAccObsBucketReplication_base(rName), OS_DESTINATION_BUCKET)
}

func testAccObsBucketReplication_update(rName string) string {
	return fmt.Sprintf(`
%s

resource "flexibleengine_obs_bucket_replication" "replica" {
  bucket             = flexibleengine_obs_bucket.source.id
  destination_bucket = "%s"
  agency             = flexibleengine_identity_agency_v3.agency_1.name

  rule {
    prefix = "abc"
  }
  rule {
	enabled       = false
    prefix        = "terraform"
	storage_class = "COLD"
  }
}
`, testAccObsBucketReplication_base(rName), OS_DESTINATION_BUCKET)
}
