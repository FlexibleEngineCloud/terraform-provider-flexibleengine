package acceptance

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccSDRSDomainDataSource_basic(t *testing.T) {
	rName := "data.flexibleengine_sdrs_domain_v1.test"
	dc := acceptance.InitDataSourceCheck(rName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckSDRSDomainDataSource_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(rName, "name"),
					resource.TestCheckResourceAttrSet(rName, "description"),
				),
			},
			{
				Config:      testAccCheckSDRSDomainDataSource_checkError,
				ExpectError: regexp.MustCompile(`your query returned no results. Please change your search criteria and try again`),
			},
		},
	})
}

const testAccCheckSDRSDomainDataSource_basic = `
data "flexibleengine_sdrs_domain_v1" "test" {}
`

const testAccCheckSDRSDomainDataSource_checkError = `
data "flexibleengine_sdrs_domain_v1" "test" {
  name = "error_check"
}`