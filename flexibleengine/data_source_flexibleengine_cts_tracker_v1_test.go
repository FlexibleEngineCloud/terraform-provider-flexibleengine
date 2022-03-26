package flexibleengine

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccCTSTrackerV1DataSource_basic(t *testing.T) {
	var bucketName = fmt.Sprintf("acc-test-%s", acctest.RandString(5))
	dataName := "data.flexibleengine_cts_tracker_v1.tracker_v1"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCTSTrackerV1DataSource_basic(bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCTSTrackerV1DataSourceID(dataName),
					resource.TestCheckResourceAttr(dataName, "bucket_name", bucketName),
					resource.TestCheckResourceAttr(dataName, "status", "enabled"),
				),
			},
		},
	})
}

func testAccCheckCTSTrackerV1DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find cts tracker data source: %s ", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("tracker data source not set ")
		}

		return nil
	}
}

func testAccCTSTrackerV1DataSource_basic(bucketName string) string {
	return fmt.Sprintf(`
resource "flexibleengine_obs_bucket" "bucket" {
  bucket        = "%s"
  acl           = "public-read"
  force_destroy = true
}

resource "flexibleengine_cts_tracker_v1" "tracker_v1" {
  bucket_name      = flexibleengine_obs_bucket.bucket.bucket
}

data "flexibleengine_cts_tracker_v1" "tracker_v1" {  
  tracker_name = flexibleengine_cts_tracker_v1.tracker_v1.id
}
`, bucketName)
}
