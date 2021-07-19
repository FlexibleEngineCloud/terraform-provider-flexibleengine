package flexibleengine

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccELBV2LoadbalancerDataSource_basic(t *testing.T) {
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLBV2LoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccELBV2LoadbalancerDataSource_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckELBV2LoadbalancerDataSourceID("data.flexibleengine_lb_loadbalancer_v2.test_by_name"),
					testAccCheckELBV2LoadbalancerDataSourceID("data.flexibleengine_lb_loadbalancer_v2.test_by_id"),
					resource.TestCheckResourceAttr(
						"data.flexibleengine_lb_loadbalancer_v2.test_by_name", "name", rName),
					resource.TestCheckResourceAttr(
						"data.flexibleengine_lb_loadbalancer_v2.test_by_id", "name", rName),
				),
			},
		},
	})
}

func testAccCheckELBV2LoadbalancerDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find elb load balancer data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("load balancer data source ID not set")
		}

		return nil
	}
}

func testAccELBV2LoadbalancerDataSource_basic(rName string) string {
	return fmt.Sprintf(`
resource "flexibleengine_lb_loadbalancer_v2" "test" {
  name          = "%s"
  description   = "resource for load balancer data source"
  vip_subnet_id = "%s"

  tags = {
    key   = "value"
    owner = "terraform"
  }
}

data "flexibleengine_lb_loadbalancer_v2" "test_by_name" {
  name = flexibleengine_lb_loadbalancer_v2.test.name
}

data "flexibleengine_lb_loadbalancer_v2" "test_by_id" {
  id = flexibleengine_lb_loadbalancer_v2.test.id
}
`, rName, OS_SUBNET_ID)
}
