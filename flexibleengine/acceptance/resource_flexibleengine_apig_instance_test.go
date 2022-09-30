package acceptance

import (
	"fmt"
	"testing"

	"github.com/chnsz/golangsdk/openstack/apigw/dedicated/v2/instances"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
)

func TestAccApigInstanceV2_basic(t *testing.T) {
	var resourceName = "flexibleengine_apig_instance.test"
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	var instance instances.Instance

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: TestAccProviderFactories,
		CheckDestroy:      testAccCheckApigInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccApigInstance_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApigInstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "edition", "BASIC"),
					resource.TestCheckResourceAttr(resourceName, "maintain_begin", "14:00:00"),
					resource.TestCheckResourceAttr(resourceName, "maintain_end", "18:00:00"),
					resource.TestCheckResourceAttr(resourceName, "description", "created by acc test"),
					resource.TestCheckResourceAttrSet(resourceName, "vpc_ingress_address"),
				),
			},
			{
				Config: testAccApigInstance_update(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApigInstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "name", rName+"-update"),
					resource.TestCheckResourceAttr(resourceName, "edition", "BASIC"),
					resource.TestCheckResourceAttr(resourceName, "maintain_begin", "18:00:00"),
					resource.TestCheckResourceAttr(resourceName, "maintain_end", "22:00:00"),
					resource.TestCheckResourceAttr(resourceName, "description", "updated by acc test"),
					resource.TestCheckResourceAttrSet(resourceName, "vpc_ingress_address"),
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

func TestAccApigInstanceV2_egress(t *testing.T) {
	var resourceName = "flexibleengine_apig_instance.test"
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	var instance instances.Instance

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: TestAccProviderFactories,
		CheckDestroy:      testAccCheckApigInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccApigInstance_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApigInstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "edition", "BASIC"),
					resource.TestCheckResourceAttr(resourceName, "maintain_begin", "14:00:00"),
					resource.TestCheckResourceAttr(resourceName, "description", "created by acc test"),
					resource.TestCheckResourceAttrSet(resourceName, "vpc_ingress_address"),
					resource.TestCheckResourceAttr(resourceName, "bandwidth_size", "0"),
				),
			},
			{
				Config: testAccApigInstance_egress(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApigInstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttrSet(resourceName, "vpc_ingress_address"),
					resource.TestCheckResourceAttr(resourceName, "bandwidth_size", "3"),
					resource.TestCheckResourceAttrSet(resourceName, "egress_address"),
				),
			},
			{
				Config: testAccApigInstance_egressUpdate(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApigInstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttrSet(resourceName, "vpc_ingress_address"),
					resource.TestCheckResourceAttr(resourceName, "bandwidth_size", "5"),
					resource.TestCheckResourceAttrSet(resourceName, "egress_address"),
				),
			},
			{
				Config: testAccApigInstance_basic(rName), // Unbind egress nat
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApigInstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttrSet(resourceName, "vpc_ingress_address"),
					resource.TestCheckResourceAttr(resourceName, "bandwidth_size", "0"),
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

func TestAccApigInstanceV2_ingress(t *testing.T) {
	var resourceName = "flexibleengine_apig_instance.test"
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	var instance instances.Instance

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: TestAccProviderFactories,
		CheckDestroy:      testAccCheckApigInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccApigInstance_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApigInstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "edition", "BASIC"),
					resource.TestCheckResourceAttr(resourceName, "maintain_begin", "14:00:00"),
					resource.TestCheckResourceAttr(resourceName, "description", "created by acc test"),
					resource.TestCheckResourceAttrSet(resourceName, "vpc_ingress_address"),
				),
			},
			{
				Config: testAccApigInstance_ingress(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApigInstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttrSet(resourceName, "vpc_ingress_address"),
					resource.TestCheckResourceAttrSet(resourceName, "eip_id"),
					resource.TestCheckResourceAttrSet(resourceName, "ingress_address"),
				),
			},
			{
				Config: testAccApigInstance_ingressUpdate(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApigInstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttrSet(resourceName, "vpc_ingress_address"),
					resource.TestCheckResourceAttrSet(resourceName, "eip_id"),
					resource.TestCheckResourceAttrSet(resourceName, "ingress_address"),
				),
			},
			{
				Config: testAccApigInstance_basic(rName), // Unbind ingress eip
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApigInstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttrSet(resourceName, "vpc_ingress_address"),
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

func testAccCheckApigInstanceDestroy(s *terraform.State) error {
	config := acceptance.TestAccProvider.Meta().(*config.Config)
	client, err := config.ApigV2Client(OS_REGION_NAME)
	if err != nil {
		return fmtp.Errorf("Error creating Flexibleengine APIG v2 client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "flexibleengine_apig_instance" {
			continue
		}
		_, err := instances.Get(client, rs.Primary.ID).Extract()
		if err == nil {
			return fmtp.Errorf("APIG v2 instance (%s) is still exists", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckApigInstanceExists(n string, instance *instances.Instance) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmtp.Errorf("Resource %s not found", n)
		}
		if rs.Primary.ID == "" {
			return fmtp.Errorf("No ID is set")
		}

		config := acceptance.TestAccProvider.Meta().(*config.Config)
		client, err := config.ApigV2Client(OS_REGION_NAME)
		if err != nil {
			return fmtp.Errorf("Error creating Flexibleengine APIG v2 client: %s", err)
		}

		found, err := instances.Get(client, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}
		*instance = *found
		return nil
	}
}

func testAccApigInstance_base(rName string) string {
	return fmt.Sprintf(`
data "flexibleengine_availability_zones" "test" {}

resource "flexibleengine_vpc_v1" "test" {
  name = "%s"
  cidr = "192.168.0.0/16"
}

resource "flexibleengine_vpc_subnet_v1" "test" {
  name       = "%s"
  vpc_id     = flexibleengine_vpc_v1.test.id
  gateway_ip = "192.168.0.1"
  cidr       = "192.168.0.0/24"
}
`, rName, rName)
}

func testAccApigInstance_basic(rName string) string {
	return fmt.Sprintf(`
%s

resource "flexibleengine_networking_secgroup_v2" "test" {
  name = "%s"
}

resource "flexibleengine_apig_instance" "test" {
  name                  = "%s"
  edition               = "BASIC"
  vpc_id                = flexibleengine_vpc_v1.test.id
  subnet_id             = flexibleengine_vpc_subnet_v1.test.id
  security_group_id     = flexibleengine_networking_secgroup_v2.test.id
  maintain_begin        = "14:00:00"
  description           = "created by acc test"
  enterprise_project_id = "0"

  available_zones = [
    data.flexibleengine_availability_zones.test.names[0],
  ]
}
`, testAccApigInstance_base(rName), rName, rName)
}

func testAccApigInstance_update(rName string) string {
	return fmt.Sprintf(`
%s

resource "flexibleengine_networking_secgroup_v2" "test" {
  name = "%s"
}

resource "flexibleengine_networking_secgroup_v2" "update" {
  name = "%s-update"
}

resource "flexibleengine_apig_instance" "test" {
  name                  = "%s-update"
  edition               = "BASIC"
  vpc_id                = flexibleengine_vpc_v1.test.id
  subnet_id             = flexibleengine_vpc_subnet_v1.test.id
  security_group_id     = flexibleengine_networking_secgroup_v2.update.id
  maintain_begin        = "18:00:00"
  description           = "updated by acc test"
  enterprise_project_id = "0"

  available_zones = [
    data.flexibleengine_availability_zones.test.names[0],
  ]
}
`, testAccApigInstance_base(rName), rName, rName, rName)
}

func testAccApigInstance_egress(rName string) string {
	return fmt.Sprintf(`
%s

resource "flexibleengine_networking_secgroup_v2" "test" {
  name = "%s"
}

resource "flexibleengine_apig_instance" "test" {
  name                  = "%s"
  edition               = "BASIC"
  vpc_id                = flexibleengine_vpc_v1.test.id
  subnet_id             = flexibleengine_vpc_subnet_v1.test.id
  security_group_id     = flexibleengine_networking_secgroup_v2.test.id
  maintain_begin        = "14:00:00"
  description           = "created by acc test"
  bandwidth_size        = 3
  enterprise_project_id = "0"

  available_zones = [
    data.flexibleengine_availability_zones.test.names[0],
  ]
}
`, testAccApigInstance_base(rName), rName, rName)
}

func testAccApigInstance_egressUpdate(rName string) string {
	return fmt.Sprintf(`
%s

resource "flexibleengine_networking_secgroup_v2" "test" {
  name = "%s"
}

resource "flexibleengine_apig_instance" "test" {
  name                  = "%s"
  edition               = "BASIC"
  vpc_id                = flexibleengine_vpc_v1.test.id
  subnet_id             = flexibleengine_vpc_subnet_v1.test.id
  security_group_id     = flexibleengine_networking_secgroup_v2.test.id
  maintain_begin        = "14:00:00"
  description           = "created by acc test"
  bandwidth_size        = 5
  enterprise_project_id = "0"

  available_zones = [
    data.flexibleengine_availability_zones.test.names[0],
  ]
}
`, testAccApigInstance_base(rName), rName, rName)
}

func testAccApigInstance_ingress(rName string) string {
	return fmt.Sprintf(`
%s

resource "flexibleengine_vpc_eip" "test" {
  publicip {
    type = "5_bgp"
  }
  bandwidth {
    name        = "%s"
    size        = 3
    share_type  = "PER"
    charge_mode = "traffic"
  }
}

resource "flexibleengine_networking_secgroup_v2" "test" {
  name = "%s"
}

resource "flexibleengine_apig_instance" "test" {
  name                  = "%s"
  edition               = "BASIC"
  vpc_id                = flexibleengine_vpc_v1.test.id
  subnet_id             = flexibleengine_vpc_subnet_v1.test.id
  security_group_id     = flexibleengine_networking_secgroup_v2.test.id
  maintain_begin        = "14:00:00"
  description           = "created by acc test"
  eip_id                = flexibleengine_vpc_eip.test.id
  enterprise_project_id = "0"

  available_zones = [
    data.flexibleengine_availability_zones.test.names[0],
  ]
}
`, testAccApigInstance_base(rName), rName, rName, rName)
}

func testAccApigInstance_ingressUpdate(rName string) string {
	return fmt.Sprintf(`
%s

resource "flexibleengine_vpc_eip" "update" {
  publicip {
    type = "5_bgp"
  }
  bandwidth {
    name        = "%s-update"
    size        = 4
    share_type  = "PER"
    charge_mode = "traffic"
  }
}

resource "flexibleengine_networking_secgroup_v2" "test" {
  name = "%s"
}

resource "flexibleengine_apig_instance" "test" {
  name                  = "%s"
  edition               = "BASIC"
  vpc_id                = flexibleengine_vpc_v1.test.id
  subnet_id             = flexibleengine_vpc_subnet_v1.test.id
  security_group_id     = flexibleengine_networking_secgroup_v2.test.id
  maintain_begin        = "14:00:00"
  description           = "created by acc test"
  eip_id                = flexibleengine_vpc_eip.update.id
  enterprise_project_id = "0"

  available_zones = [
    data.flexibleengine_availability_zones.test.names[0],
  ]
}
`, testAccApigInstance_base(rName), rName, rName, rName)
}
