package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/elb/v3/pools"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func getMemberResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	c, err := conf.LoadBalancerClient(OS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating ELB v3 client: %s", err)
	}
	poolId := state.Primary.Attributes["pool_id"]
	resp, err := pools.GetMember(c, poolId, state.Primary.ID).Extract()
	if resp == nil && err == nil {
		return resp, fmt.Errorf("unable to find the member (%s)", state.Primary.ID)
	}
	return resp, err
}

func TestAccElbV3Member_basic(t *testing.T) {
	var member_1 pools.Member
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "flexibleengine_lb_member_v3.member_1"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&member_1,
		getMemberResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccElbV3MemberConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					testAccCheckElbV3MemberExists(resourceName, &member_1),
				),
			},
			{
				Config: testAccElbV3MemberConfig_update(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "weight", "10"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccELBMemberImportStateIdFunc(),
			},
		},
	})
}

func TestAccElbV3Member_crossVpcBackend(t *testing.T) {
	var member_1 pools.Member
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "flexibleengine_lb_member_v3.member_1"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&member_1,
		getMemberResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccElbV3MemberConfig_crossVpcBackend_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					testAccCheckElbV3MemberExists(resourceName, &member_1),
				),
			},
			{
				Config: testAccElbV3MemberConfig_crossVpcBackend_update(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "weight", "10"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccELBMemberImportStateIdFunc(),
			},
		},
	})
}

func testAccCheckElbV3MemberExists(n string, member *pools.Member) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		config := acceptance.TestAccProvider.Meta().(*config.Config)
		elbClient, err := config.ElbV3Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating ELB v3 client: %s", err)
		}

		poolId := rs.Primary.Attributes["pool_id"]
		found, err := pools.GetMember(elbClient, poolId, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("member not found")
		}

		*member = *found

		return nil
	}
}

func testAccELBMemberImportStateIdFunc() resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		pool, ok := s.RootModule().Resources["flexibleengine_lb_pool_v3.test"]
		if !ok {
			return "", fmt.Errorf("pool not found: %s", pool)
		}
		member, ok := s.RootModule().Resources["flexibleengine_lb_member_v3.member_1"]
		if !ok {
			return "", fmt.Errorf("member not found: %s", member)
		}
		if pool.Primary.ID == "" || member.Primary.ID == "" {
			return "", fmt.Errorf("resource not found: %s/%s", pool.Primary.ID, member.Primary.ID)
		}
		return fmt.Sprintf("%s/%s", pool.Primary.ID, member.Primary.ID), nil
	}
}

func testAccMemberV3Config_base(rName string) string {
	return fmt.Sprintf(`
resource "flexibleengine_vpc_v1" "test" {
  name = "vpc_%[1]s"
  cidr = "192.168.0.0/16"
}

resource "flexibleengine_vpc_subnet_v1" "test" {
  name        = "subnet_%[1]s"
  cidr        = "192.168.0.0/24"
  gateway_ip  = "192.168.0.1"
  vpc_id      = flexibleengine_vpc_v1.test.id
  ipv6_enable = true
}

resource "flexibleengine_lb_loadbalancer_v3" "test" {
  name            = "%[1]s"
  ipv4_subnet_id  = flexibleengine_vpc_subnet_v1.test.subnet_id
  ipv6_network_id = flexibleengine_vpc_subnet_v1.test.id
  availability_zone = [
    "%[2]s"
  ]
}

resource "flexibleengine_lb_listener_v3" "test" {
  name            = "%[1]s"
  description     = "test description"
  protocol        = "HTTP"
  protocol_port   = 8080
  loadbalancer_id = flexibleengine_lb_loadbalancer_v3.test.id
  forward_eip = true
  idle_timeout     = 60
  request_timeout  = 60
  response_timeout = 60
}

resource "flexibleengine_lb_pool_v3" "test" {
  name        = "%[1]s"
  protocol    = "HTTP"
  lb_method   = "ROUND_ROBIN"
  listener_id = flexibleengine_lb_listener_v3.test.id
}
`, rName, OS_AVAILABILITY_ZONE)
}

func testAccElbV3MemberConfig_basic(rName string) string {
	return fmt.Sprintf(`
%s

resource "flexibleengine_lb_member_v3" "member_1" {
  address        = "192.168.0.10"
  protocol_port  = 8080
  weight         = 10
  pool_id        = flexibleengine_lb_pool_v3.test.id
  subnet_id      = flexibleengine_vpc_subnet_v1.test.subnet_id
}
`, testAccMemberV3Config_base(rName))
}

func testAccElbV3MemberConfig_update(rName string) string {
	return fmt.Sprintf(`
%s

resource "flexibleengine_lb_member_v3" "member_1" {
  address        = "192.168.0.10"
  protocol_port  = 8080
  weight         = 10
  pool_id        = flexibleengine_lb_pool_v3.test.id
  subnet_id      = flexibleengine_vpc_subnet_v1.test.subnet_id
}
`, testAccMemberV3Config_base(rName))
}

func testAccMemberV3Config_crossVpcBackend_base(rName string) string {
	return fmt.Sprintf(`
resource "flexibleengine_vpc_v1" "test" {
  name = "vpc_%[1]s"
  cidr = "192.168.0.0/16"
}

resource "flexibleengine_vpc_subnet_v1" "test" {
  name        = "subnet_%[1]s"
  cidr        = "192.168.0.0/24"
  gateway_ip  = "192.168.0.1"
  vpc_id      = flexibleengine_vpc_v1.test.id
  ipv6_enable = true
}

resource "flexibleengine_lb_loadbalancer_v3" "test" {
  name              = "%[1]s"
  cross_vpc_backend = true
  ipv4_subnet_id    = flexibleengine_vpc_subnet_v1.test.subnet_id
  ipv6_network_id   = flexibleengine_vpc_subnet_v1.test.id
  availability_zone = [
    "%[2]s"
  ]
}

resource "flexibleengine_lb_listener_v3" "test" {
  name            = "%[1]s"
  description     = "test description"
  protocol        = "HTTP"
  protocol_port   = 8080
  loadbalancer_id = flexibleengine_lb_loadbalancer_v3.test.id
  forward_eip = true
  idle_timeout     = 60
  request_timeout  = 60
  response_timeout = 60
}

resource "flexibleengine_lb_pool_v3" "test" {
  name        = "%[1]s"
  protocol    = "HTTP"
  lb_method   = "ROUND_ROBIN"
  listener_id = flexibleengine_lb_listener_v3.test.id
}
`, rName, OS_AVAILABILITY_ZONE)
}

func testAccElbV3MemberConfig_crossVpcBackend_basic(rName string) string {
	return fmt.Sprintf(`
%s

resource "flexibleengine_lb_member_v3" "member_1" {
  address        = "121.121.0.120"
  protocol_port  = 8080
  pool_id        = flexibleengine_lb_pool_v3.test.id
}
`, testAccMemberV3Config_crossVpcBackend_base(rName))
}

func testAccElbV3MemberConfig_crossVpcBackend_update(rName string) string {
	return fmt.Sprintf(`
%s

resource "flexibleengine_lb_member_v3" "member_1" {
  address        = "121.121.0.120"
  protocol_port  = 8080
  weight         = 10
  pool_id        = flexibleengine_lb_pool_v3.test.id
}
`, testAccMemberV3Config_crossVpcBackend_base(rName))
}
