package flexibleengine

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccRdsFlavorV1DataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccRdsFlavorV1DataSource_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRdsFlavorV1DataSourceID("data.flexibleengine_rds_flavors_v1.flavor"),
					resource.TestCheckResourceAttrSet(
						"data.flexibleengine_rds_flavors_v1.flavor", "id"),
					resource.TestCheckResourceAttrSet(
						"data.flexibleengine_rds_flavors_v1.flavor", "speccode"),
				),
			},
		},
	})
}

func TestAccRdsFlavorV1DataSource_speccode(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccRdsFlavorV1DataSource_speccode,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingNetworkV2DataSourceID("data.flexibleengine_rds_flavors_v1.flavor"),
					resource.TestCheckResourceAttr(
						"data.flexibleengine_rds_flavors_v1.flavor", "speccode", "rds.mysql.s1.medium"),
				),
			},
		},
	})
}

func testAccCheckRdsFlavorV1DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find rds data source: %s ", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Rds data source ID not set ")
		}

		return nil
	}
}

var testAccRdsFlavorV1DataSource_basic = `

data "flexibleengine_rds_flavors_v1" "flavor" {
    region = "eu-west-0"
	datastore_name = "MySQL"
    datastore_version = "5.6.30"
}
`

var testAccRdsFlavorV1DataSource_speccode = `

data "flexibleengine_rds_flavors_v1" "flavor" {
    region = "eu-west-0"
	datastore_name = "MySQL"
    datastore_version = "5.6.30"
    speccode = "rds.mysql.s1.medium"
}
`
