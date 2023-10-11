package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataTopics_basic(t *testing.T) {
	dataSourceName := "data.flexibleengine_smn_topics.test"
	resourceName := "flexibleengine_smn_topic_v2.topic_1"
	dc := acceptance.InitDataSourceCheck(dataSourceName)
	rName := acceptance.RandomAccResourceNameWithDash()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataTopicsConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "name", rName),
					resource.TestCheckResourceAttrPair(dataSourceName, "topics.0.id", resourceName, "id"),
					resource.TestCheckResourceAttrPair(dataSourceName, "topics.0.topic_urn", resourceName, "topic_urn"),
					resource.TestCheckResourceAttrPair(dataSourceName, "topics.0.display_name", resourceName, "display_name"),
				),
			},
		},
	})
}

func testAccDataTopicsConfig_basic(rName string) string {
	return fmt.Sprintf(`
resource "flexibleengine_smn_topic_v2" "topic_1" {
  name                  = "%[1]s"
  display_name          = "The display name of %[1]s"
  enterprise_project_id = "0"
}

data "flexibleengine_smn_topics" "test" {
  name = "%[1]s"

  depends_on = [
    flexibleengine_smn_topic_v2.topic_1
  ]
}
`, rName)
}
