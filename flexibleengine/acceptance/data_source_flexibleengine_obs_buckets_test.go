package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceObsBuckets_basic(t *testing.T) {
	dataSourceName := "data.flexibleengine_obs_buckets.buckets"
	name := acceptance.RandomAccResourceNameWithDash()
	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckOBS(t)
		},
		ProviderFactories: TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccObsBuckets_conf(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "buckets.0.bucket", name),
				),
			},
		},
	})
}

func testAccObsBuckets_conf(name string) string {
	return fmt.Sprintf(`
resource "flexibleengine_obs_bucket" "bucket" {
  bucket        = "%s"
  storage_class = "STANDARD"
  acl           = "private"
}

data "flexibleengine_obs_buckets" "buckets" {
  bucket = "%s"

  depends_on = [flexibleengine_obs_bucket.bucket]
}
`, name, name)
}
