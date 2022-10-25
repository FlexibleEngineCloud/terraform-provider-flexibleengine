package flexibleengine

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDcsProductDataSource_basic(t *testing.T) {
	dataSourceName := "data.flexibleengine_dcs_product_v1.product_spec"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDcsProductDataSource_multi,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDcsProductDataSourceID(dataSourceName),
					resource.TestCheckResourceAttr(dataSourceName, "spec_code", "dcs.single_node"),
					resource.TestCheckResourceAttr(dataSourceName, "engine", "redis"),
					resource.TestCheckResourceAttrSet("data.flexibleengine_dcs_product_v1.product1", "spec_code"),
					resource.TestCheckResourceAttrSet("data.flexibleengine_dcs_product_v1.product2", "spec_code"),
					resource.TestCheckResourceAttrSet("data.flexibleengine_dcs_product_v1.product3", "spec_code"),
					resource.TestCheckResourceAttrSet("data.flexibleengine_dcs_product_v1.product4", "spec_code"),
					resource.TestCheckResourceAttrSet("data.flexibleengine_dcs_product_v1.product5", "spec_code"),
					resource.TestCheckResourceAttrSet("data.flexibleengine_dcs_product_v1.product6", "spec_code"),
				),
			},
		},
	})
}

func testAccCheckDcsProductDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find DCS product data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Dcs product data source ID not set")
		}

		return nil
	}
}

var testAccDcsProductDataSource_multi = fmt.Sprintf(`
# product of Redis 4.0/5.0 with Redis Cluster type
data "flexibleengine_dcs_product_v1" "product1" {
  engine         = "redis"
  engine_version = "4.0;5.0"
  cache_mode     = "cluster"
  capacity       = 8
  replica_count  = 2
}

# product of Redis 4.0/5.0 with Master/Standby type
data "flexibleengine_dcs_product_v1" "product2" {
  engine         = "redis"
  engine_version = "4.0;5.0"
  cache_mode     = "ha"
  capacity       = 0.125
  replica_count  = 2
}

# product of Redis 4.0/5.0 with Single-node type
data "flexibleengine_dcs_product_v1" "product3" {
  engine         = "redis"
  engine_version = "4.0;5.0"
  cache_mode     = "single"
  capacity       = 1
}

# product of Redis 4.0/5.0 with Proxy Cluster type
data "flexibleengine_dcs_product_v1" "product4" {
  engine         = "redis"
  engine_version = "4.0;5.0"
  cache_mode     = "proxy"
  capacity       = 4
}

# product of Redis 3.0 instance
data "flexibleengine_dcs_product_v1" "product5" {
  engine         = "redis"
  engine_version = "3.0"
  cache_mode     = "ha"
}

# product of Memcached instance
data "flexibleengine_dcs_product_v1" "product6" {
  engine     = "memcached"
  cache_mode = "single"
}

# product with spec_code
data "flexibleengine_dcs_product_v1" "product_spec" {
  spec_code = "dcs.single_node"
}
`)
