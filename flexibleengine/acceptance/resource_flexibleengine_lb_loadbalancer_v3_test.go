package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/elb/v3/loadbalancers"
	"github.com/chnsz/golangsdk/openstack/networking/v1/eips"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func getELBResourceFunc(c *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := c.ElbV3Client(OS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("Error creating ELB v3 client: %s", err)
	}

	eipID := state.Primary.Attributes["ipv4_eip_id"]
	eipType := state.Primary.Attributes["iptype"]
	if eipType != "" && eipID != "" {
		eipClient, err := c.NetworkingV1Client(OS_REGION_NAME)
		if err != nil {
			return nil, fmt.Errorf("Error creating VPC v1 client: %s", err)
		}

		if _, err := eips.Get(eipClient, eipID).Extract(); err != nil {
			return nil, err
		}
	}

	return loadbalancers.Get(client, state.Primary.ID).Extract()
}

func TestAccElbV3LoadBalancer_basic(t *testing.T) {
	var lb loadbalancers.LoadBalancer
	rName := acceptance.RandomAccResourceNameWithDash()
	rNameUpdate := acceptance.RandomAccResourceNameWithDash()
	resourceName := "flexibleengine_lb_loadbalancer_v3.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&lb,
		getELBResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccElbV3LoadBalancerConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "cross_vpc_backend", "false"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
					resource.TestCheckResourceAttr(resourceName, "tags.owner", "terraform"),
				),
			},
			{
				Config: testAccElbV3LoadBalancerConfig_update(rName, rNameUpdate),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdate),
					resource.TestCheckResourceAttr(resourceName, "cross_vpc_backend", "true"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1"),
					resource.TestCheckResourceAttr(resourceName, "tags.owner", "terraform_update"),
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

func TestAccElbV3LoadBalancer_withEIP(t *testing.T) {
	var lb loadbalancers.LoadBalancer
	rName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "flexibleengine_lb_loadbalancer_v3.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&lb,
		getELBResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccElbV3LoadBalancerConfig_withEIP(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "iptype", "5_bgp"),
					resource.TestCheckResourceAttrSet(resourceName, "ipv4_eip_id"),
				),
			},
		},
	})
}

func testAccElbV3Config_base(rName string) string {
	return fmt.Sprintf(`
data "flexibleengine_availability_zones" "test" {}

resource "flexibleengine_vpc_v1" "test" {
  name = "vpc_%s"
  cidr = "192.168.0.0/16"
}

resource "flexibleengine_vpc_subnet_v1" "test" {
  name        = "subnet_%s"
  cidr        = "192.168.0.0/24"
  gateway_ip  = "192.168.0.1"
  vpc_id      = flexibleengine_vpc_v1.test.id
  ipv6_enable = true
}
`, rName, rName)
}

func testAccElbV3LoadBalancerConfig_basic(rName string) string {
	return fmt.Sprintf(`
%s

resource "flexibleengine_lb_loadbalancer_v3" "test" {
  name           		 = "%s"
  ipv4_subnet_id  	     = flexibleengine_vpc_subnet_v1.test.subnet_id
  ipv6_network_id        = flexibleengine_vpc_subnet_v1.test.id
  enterprise_project_id  = 0  

  availability_zone = [
    data.flexibleengine_availability_zones.test.names[0]
  ]

  tags = {
    key   = "value"
    owner = "terraform"
  }
}
`, testAccElbV3Config_base(rName), rName)
}

func testAccElbV3LoadBalancerConfig_update(rName, rNameUpdate string) string {
	return fmt.Sprintf(`
%s

resource "flexibleengine_lb_loadbalancer_v3" "test" {
  name              	 = "%s"
  cross_vpc_backend 	 = true
  ipv4_subnet_id   	     = flexibleengine_vpc_subnet_v1.test.subnet_id
  ipv6_network_id   	 = flexibleengine_vpc_subnet_v1.test.id
  enterprise_project_id  = 0 

  availability_zone = [
    data.flexibleengine_availability_zones.test.names[0]
  ]

  tags = {
    key1  = "value1"
    owner = "terraform_update"
  }
}
`, testAccElbV3Config_base(rName), rNameUpdate)
}

func testAccElbV3LoadBalancerConfig_withEIP(rName string) string {
	return fmt.Sprintf(`
%s

resource "flexibleengine_lb_loadbalancer_v3" "test" {
  name              = "%s"
  ipv4_subnet_id    = flexibleengine_vpc_subnet_v1.test.subnet_id
  availability_zone = [data.flexibleengine_availability_zones.test.names[0]]

  iptype                = "5_bgp"
  bandwidth_charge_mode = "traffic"
  sharetype             = "PER"
  bandwidth_size        = 5
}
`, testAccElbV3Config_base(rName), rName)
}
