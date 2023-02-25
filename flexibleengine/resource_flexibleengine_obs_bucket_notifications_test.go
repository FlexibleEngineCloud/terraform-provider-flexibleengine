package flexibleengine

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccObsBucket_notifications(t *testing.T) {
	rInt := acctest.RandInt()
	obsName := "flexibleengine_obs_bucket.bucket"
	resourceName := "flexibleengine_obs_bucket_notifications.notification"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckS3(t)
			testAccPreCheckOBSNotification(t)
		},
		ProviderFactories: TestAccProviderFactories,
		CheckDestroy:      testAccCheckObsBucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccObsBucketConfigWithNotification(rInt, OS_OBS_URN_SMN),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckObsBucketExists(obsName),
					resource.TestCheckResourceAttr(resourceName, "notifications.0.name", "001"),
					resource.TestCheckResourceAttr(resourceName, "notifications.0.events.0", "ObjectCreated:*"),
					resource.TestCheckResourceAttr(resourceName, "notifications.0.prefix", "tf"),
					resource.TestCheckResourceAttr(resourceName, "notifications.0.suffix", ".jpg"),
					resource.TestCheckResourceAttr(resourceName, "notifications.0.topic_urn", OS_OBS_URN_SMN),
					resource.TestCheckResourceAttr(resourceName, "notifications.1.events.0", "ObjectCreated:Post"),
					resource.TestCheckResourceAttr(resourceName, "notifications.1.prefix", "iac"),
					resource.TestCheckResourceAttr(resourceName, "notifications.1.suffix", ".txt"),
					resource.TestCheckResourceAttr(resourceName, "notifications.1.topic_urn", OS_OBS_URN_SMN),
					resource.TestCheckResourceAttrSet(resourceName, "notifications.1.name"),
				),
			},
			{
				Config: testAccObsBucketConfigWithNotificationUpdate(rInt, OS_OBS_URN_SMN),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckObsBucketExists(obsName),
					resource.TestCheckResourceAttr(resourceName, "notifications.0.name", "003"),
					resource.TestCheckResourceAttr(resourceName, "notifications.0.events.0", "ObjectRemoved:*"),
					resource.TestCheckResourceAttr(resourceName, "notifications.0.prefix", "tf_update"),
					resource.TestCheckResourceAttr(resourceName, "notifications.0.suffix", ".png"),
					resource.TestCheckResourceAttr(resourceName, "notifications.0.topic_urn", OS_OBS_URN_SMN),
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

func testAccObsBucketConfigWithNotification(randInt int, urnSmn string) string {
	return fmt.Sprintf(`
resource "flexibleengine_obs_bucket" "bucket" {
  bucket = "tf-test-bucket-%d"
  acl    = "public-read"
}

resource "flexibleengine_obs_bucket_notifications" notification {
  bucket    = flexibleengine_obs_bucket.bucket.bucket

  notifications {
    name      = "001"
    events    = ["ObjectCreated:*"]
    prefix    = "tf"
    suffix    = ".jpg"
    topic_urn = "%s"
  }

  notifications {
    events    = ["ObjectCreated:Post"]
    prefix    = "iac"
    suffix    = ".txt"
    topic_urn = "%s"
  }
}

`, randInt, urnSmn, urnSmn)
}

func testAccObsBucketConfigWithNotificationUpdate(randInt int, urnSmn string) string {
	return fmt.Sprintf(`
resource "flexibleengine_obs_bucket" "bucket" {
  bucket = "tf-test-bucket-%d"
  acl    = "public-read"
}

resource "flexibleengine_obs_bucket_notifications" notification {
  bucket = flexibleengine_obs_bucket.bucket.bucket

  notifications {
    name      = "003"
    events    = ["ObjectRemoved:*"]
    prefix    = "tf_update"
    suffix    = ".png"
    topic_urn = "%s"
  }
}

`, randInt, urnSmn)
}
