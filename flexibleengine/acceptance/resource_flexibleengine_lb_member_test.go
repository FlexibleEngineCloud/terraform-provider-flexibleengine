package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/elb/v2/pools"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
)

func TestAccLBV2Member_basic(t *testing.T) {
	var member_1 pools.Member
	var member_2 pools.Member
	resourceName1 := "flexibleengine_lb_member.member_1"
	resourceName2 := "flexibleengine_lb_member.member_2"
	rName := acceptance.RandomAccResourceNameWithDash()

	rc1 := acceptance.InitResourceCheck(
		resourceName1,
		&member_1,
		getMemberResourceFunc,
	)
	rc2 := acceptance.InitResourceCheck(
		resourceName2,
		&member_2,
		getMemberResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: TestAccProviderFactories,
		CheckDestroy:      testAccCheckLBV2MemberDestroy,
		Steps: []resource.TestStep{
			{
				Config:             testAccLBV2MemberConfig_basic(rName),
				ExpectNonEmptyPlan: true, // Because admin_state_up remains false.
				Check: resource.ComposeTestCheckFunc(
					rc1.CheckResourceExists(),
					rc2.CheckResourceExists(),
				),
			},
			{
				Config:             testAccLBV2MemberConfig_update(rName),
				ExpectNonEmptyPlan: true, // Because admin_state_up remains false.
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("flexibleengine_lb_member.member_1", "weight", "10"),
					resource.TestCheckResourceAttr("flexibleengine_lb_member.member_2", "weight", "15"),
				),
			},
			{
				ResourceName:      "flexibleengine_lb_member.member_1",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccLBMemberImportStateIdFunc(),
			},
		},
	})
}

func testAccCheckLBV2MemberDestroy(s *terraform.State) error {
	config := acceptance.TestAccProvider.Meta().(*config.Config)
	elbClient, err := config.LoadBalancerClient(OS_REGION_NAME)
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud elb client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "flexibleengine_lb_member" {
			continue
		}

		poolId := rs.Primary.Attributes["pool_id"]
		_, err := pools.GetMember(elbClient, poolId, rs.Primary.ID).Extract()
		if err == nil {
			return fmtp.Errorf("Member still exists: %s", rs.Primary.ID)
		}
	}

	return nil
}

func testAccLBMemberImportStateIdFunc() resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		pool, ok := s.RootModule().Resources["flexibleengine_lb_pool.pool_1"]
		if !ok {
			return "", fmt.Errorf("pool not found: %s", pool)
		}
		member, ok := s.RootModule().Resources["flexibleengine_lb_member.member_1"]
		if !ok {
			return "", fmt.Errorf("member not found: %s", member)
		}
		if pool.Primary.ID == "" || member.Primary.ID == "" {
			return "", fmt.Errorf("resource not found: %s/%s", pool.Primary.ID, member.Primary.ID)
		}
		return fmt.Sprintf("%s/%s", pool.Primary.ID, member.Primary.ID), nil
	}
}

func testAccLBV2MemberConfig_basic(rName string) string {
	return fmt.Sprintf(`
resource "flexibleengine_vpc_v1" "test" {
  name = "%[1]s"
  cidr = "192.168.0.0/16"
}

resource "flexibleengine_vpc_subnet_v1" "test" {
  name       = "%[1]s"
  cidr       = "192.168.0.0/24"
  vpc_id     = flexibleengine_vpc_v1.test.id
  gateway_ip = "192.168.0.1"
}

resource "flexibleengine_lb_loadbalancer" "loadbalancer_1" {
  name          = "%[1]s"
  vip_subnet_id = flexibleengine_vpc_subnet_v1.test.ipv4_subnet_id

}

resource "flexibleengine_lb_listener" "listener_1" {
  name            = "%[1]s"
  protocol        = "HTTP"
  protocol_port   = 8080
  loadbalancer_id = flexibleengine_lb_loadbalancer.loadbalancer_1.id
}

resource "flexibleengine_lb_pool" "pool_1" {
  name        = "%[1]s"
  protocol    = "HTTP"
  lb_method   = "ROUND_ROBIN"
  listener_id = flexibleengine_lb_listener.listener_1.id
}

resource "flexibleengine_lb_member" "member_1" {
  address       = "192.168.0.10"
  protocol_port = 8080
  pool_id       = flexibleengine_lb_pool.pool_1.id
  subnet_id     = flexibleengine_vpc_subnet_v1.test.ipv4_subnet_id

  timeouts {
    create = "5m"
    update = "5m"
    delete = "5m"
  }
}

resource "flexibleengine_lb_member" "member_2" {
  address       = "192.168.0.11"
  protocol_port = 8080
  pool_id       = flexibleengine_lb_pool.pool_1.id
  subnet_id     = flexibleengine_vpc_subnet_v1.test.ipv4_subnet_id

  timeouts {
    create = "5m"
    update = "5m"
    delete = "5m"
  }
}
`, rName)
}

func testAccLBV2MemberConfig_update(rName string) string {
	return fmt.Sprintf(`
resource "flexibleengine_vpc_v1" "test" {
  name = "%[1]s"
  cidr = "192.168.0.0/16"
}

resource "flexibleengine_vpc_subnet_v1" "test" {
  name       = "%[1]s"
  cidr       = "192.168.0.0/24"
  vpc_id     = flexibleengine_vpc_v1.test.id
  gateway_ip = "192.168.0.1"
}

resource "flexibleengine_lb_loadbalancer" "loadbalancer_1" {
  name          = "%[1]s"
  vip_subnet_id = flexibleengine_vpc_subnet_v1.test.ipv4_subnet_id
}

resource "flexibleengine_lb_listener" "listener_1" {
  name            = "%[1]s"
  protocol        = "HTTP"
  protocol_port   = 8080
  loadbalancer_id = flexibleengine_lb_loadbalancer.loadbalancer_1.id
}

resource "flexibleengine_lb_pool" "pool_1" {
  name        = "%[1]s"
  protocol    = "HTTP"
  lb_method   = "ROUND_ROBIN"
  listener_id = flexibleengine_lb_listener.listener_1.id
}

resource "flexibleengine_lb_member" "member_1" {
  address        = "192.168.0.10"
  protocol_port  = 8080
  weight         = 10
  admin_state_up = "true"
  pool_id        = flexibleengine_lb_pool.pool_1.id
  subnet_id      = flexibleengine_vpc_subnet_v1.test.ipv4_subnet_id

  timeouts {
    create = "5m"
    update = "5m"
    delete = "5m"
  }
}

resource "flexibleengine_lb_member" "member_2" {
  address        = "192.168.0.11"
  protocol_port  = 8080
  weight         = 15
  admin_state_up = "true"
  pool_id        = flexibleengine_lb_pool.pool_1.id
  subnet_id      = flexibleengine_vpc_subnet_v1.test.ipv4_subnet_id

  timeouts {
    create = "5m"
    update = "5m"
    delete = "5m"
  }
}
`, rName)
}
