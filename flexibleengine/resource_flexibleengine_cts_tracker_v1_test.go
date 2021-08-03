package flexibleengine

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/openstack/cts/v1/tracker"
)

func TestAccCTSTrackerV1_basic(t *testing.T) {
	var tracker tracker.Tracker
	var bucketName = fmt.Sprintf("terra-test-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCTSTrackerV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCTSTrackerV1_basic(bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCTSTrackerV1Exists("flexibleengine_cts_tracker_v1.tracker_v1", &tracker),
					resource.TestCheckResourceAttr(
						"flexibleengine_cts_tracker_v1.tracker_v1", "bucket_name", bucketName),
					resource.TestCheckResourceAttr(
						"flexibleengine_cts_tracker_v1.tracker_v1", "file_prefix_name", "yO8Q"),
				),
			},
		},
	})
}

func TestAccCTSTrackerV1_timeout(t *testing.T) {
	var tracker tracker.Tracker
	var bucketName = fmt.Sprintf("terra-test-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCTSTrackerV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCTSTrackerV1_timeout(bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCTSTrackerV1Exists("flexibleengine_cts_tracker_v1.tracker_v1", &tracker),
				),
			},
		},
	})
}

func testAccCheckCTSTrackerV1Destroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	ctsClient, err := config.ctsV1Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating cts client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "flexibleengine_cts_tracker_v1" {
			continue
		}

		_, err := tracker.List(ctsClient, tracker.ListOpts{TrackerName: rs.Primary.ID})
		if err != nil {
			return fmt.Errorf("cts tracker still exists.")
		}
		if _, ok := err.(golangsdk.ErrDefault404); !ok {
			return err
		}
	}

	return nil
}

func testAccCheckCTSTrackerV1Exists(n string, trackers *tracker.Tracker) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		ctsClient, err := config.ctsV1Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating cts client: %s", err)
		}

		trackerList, err := tracker.List(ctsClient, tracker.ListOpts{TrackerName: rs.Primary.ID})
		if err != nil {
			return err
		}
		found := trackerList[0]
		if found.TrackerName != rs.Primary.ID {
			return fmt.Errorf("cts tracker not found")
		}

		*trackers = found

		return nil
	}
}

func testAccCTSTrackerV1_basic(bucketName string) string {
	return fmt.Sprintf(`
resource "flexibleengine_s3_bucket" "bucket" {
  bucket = "%s"
  acl = "public-read"
  force_destroy = true
}

resource "flexibleengine_cts_tracker_v1" "tracker_v1" {
  bucket_name      = "${flexibleengine_s3_bucket.bucket.bucket}"
  file_prefix_name      = "yO8Q"
}
`, bucketName)
}

func testAccCTSTrackerV1_timeout(bucketName string) string {
	return fmt.Sprintf(`
resource "flexibleengine_s3_bucket" "bucket" {
  bucket = "%s"
  acl = "public-read"
  force_destroy = true
}

resource "flexibleengine_cts_tracker_v1" "tracker_v1" {
  bucket_name      = "${flexibleengine_s3_bucket.bucket.bucket}"
  file_prefix_name      = "yO8Q"

timeouts {
    create = "5m"
    delete = "5m"
  }
}
`, bucketName)
}
