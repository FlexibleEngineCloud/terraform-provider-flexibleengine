package acceptance

import (
	"fmt"
	"testing"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/workspace/v2/services"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func getServiceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.WorkspaceV2Client(OS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating Workspace v2 client: %s", err)
	}
	resp, err := services.Get(client)
	if resp.Status == "CLOSED" {
		return nil, golangsdk.ErrDefault404{}
	}
	return resp, err
}

func TestAccService_basic(t *testing.T) {
	var (
		service      services.Service
		resourceName = "flexibleengine_workspace_service.test"
		rName        = acceptance.RandomAccResourceNameWithDash()
	)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&service,
		getServiceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccService_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(resourceName, "vpc_id", "flexibleengine_vpc_v1.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "network_ids.0",
						"flexibleengine_vpc_subnet_v1.master", "id"),
					resource.TestCheckResourceAttr(resourceName, "auth_type", "LITE_AS"),
					resource.TestCheckResourceAttr(resourceName, "access_mode", "INTERNET"),
					resource.TestCheckResourceAttrSet(resourceName, "management_subnet_cidr"),
					resource.TestCheckResourceAttrSet(resourceName, "infrastructure_security_group.0.name"),
					resource.TestCheckResourceAttrSet(resourceName, "infrastructure_security_group.0.id"),
					resource.TestCheckResourceAttrSet(resourceName, "desktop_security_group.0.name"),
					resource.TestCheckResourceAttrSet(resourceName, "desktop_security_group.0.id"),
					resource.TestCheckResourceAttrSet(resourceName, "internet_access_port"),
					resource.TestCheckResourceAttrSet(resourceName, "internet_access_address"),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
				),
			},
			{
				Config: testAccService_update(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(resourceName, "network_ids.0",
						"flexibleengine_vpc_subnet_v1.master", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "network_ids.1",
						"flexibleengine_vpc_subnet_v1.standby", "id"),
					resource.TestCheckResourceAttr(resourceName, "internet_access_port", "9001"),
					resource.TestCheckResourceAttrSet(resourceName, "internet_access_address"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_id", rName),
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

func TestAccService_localAD(t *testing.T) {
	var (
		service      services.Service
		resourceName = "flexibleengine_workspace_service.test"
		rName        = acceptance.RandomAccResourceNameWithDash()
	)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&service,
		getServiceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			TestAccPreCheckWorkspaceAD(t)
		},
		ProviderFactories: TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccService_localAD_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "auth_type", "LOCAL_AD"),
					resource.TestCheckResourceAttr(resourceName, "ad_domain.0.name", OS_WORKSPACE_AD_DOMAIN_NAME),
					resource.TestCheckResourceAttr(resourceName, "ad_domain.0.admin_account", "Administrator"),
					resource.TestCheckResourceAttr(resourceName, "ad_domain.0.password", OS_WORKSPACE_AD_SERVER_PWD),
					resource.TestCheckResourceAttr(resourceName, "ad_domain.0.active_domain_ip", OS_WORKSPACE_AD_DOMAIN_IP),
					resource.TestCheckResourceAttr(resourceName, "ad_domain.0.active_domain_name",
						fmt.Sprintf("server.%s", OS_WORKSPACE_AD_DOMAIN_NAME)),
					resource.TestCheckResourceAttr(resourceName, "ad_domain.0.active_dns_ip", OS_WORKSPACE_AD_DOMAIN_IP),
					resource.TestCheckResourceAttr(resourceName, "access_mode", "INTERNET"),
					resource.TestCheckResourceAttr(resourceName, "vpc_id", OS_WORKSPACE_AD_VPC_ID),
					resource.TestCheckResourceAttr(resourceName, "network_ids.0", OS_WORKSPACE_AD_NETWORK_ID),
					resource.TestCheckResourceAttrSet(resourceName, "infrastructure_security_group.0.name"),
					resource.TestCheckResourceAttrSet(resourceName, "infrastructure_security_group.0.id"),
					resource.TestCheckResourceAttrSet(resourceName, "desktop_security_group.0.name"),
					resource.TestCheckResourceAttrSet(resourceName, "desktop_security_group.0.id"),
					resource.TestCheckResourceAttrSet(resourceName, "internet_access_port"),
					resource.TestCheckResourceAttrSet(resourceName, "internet_access_address"),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
				),
			},
			{
				Config: testAccService_localAD_update(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "network_ids.0", OS_WORKSPACE_AD_NETWORK_ID),
					resource.TestCheckResourceAttrPair(resourceName, "network_ids.1",
						"flexibleengine_vpc_subnet_v1.master", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "network_ids.2",
						"flexibleengine_vpc_subnet_v1.standby", "id"),
					resource.TestCheckResourceAttr(resourceName, "internet_access_port", "9001"),
					resource.TestCheckResourceAttrSet(resourceName, "internet_access_address"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_id", rName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"ad_domain.0.password",
				},
			},
		},
	})
}

func testAccService_base(rName string) string {
	return fmt.Sprintf(`
resource "flexibleengine_vpc_v1" "test" {
  name = "%[1]s"
  cidr = "192.168.0.0/20"
}

resource "flexibleengine_vpc_subnet_v1" "master" {
  vpc_id = flexibleengine_vpc_v1.test.id

  name       = "%[1]s-master"
  cidr       = cidrsubnet(flexibleengine_vpc_v1.test.cidr, 4, 1)
  gateway_ip = cidrhost(cidrsubnet(flexibleengine_vpc_v1.test.cidr, 4, 1), 1)
}

resource "flexibleengine_vpc_subnet_v1" "standby" {
  vpc_id = flexibleengine_vpc_v1.test.id

  name       = "%[1]s-standby"
  cidr       = cidrsubnet(flexibleengine_vpc_v1.test.cidr, 4, 2)
  gateway_ip = cidrhost(cidrsubnet(flexibleengine_vpc_v1.test.cidr, 4, 2), 1)
}
`, rName)
}

func testAccService_basic(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "flexibleengine_workspace_service" "test" {
  access_mode = "INTERNET"
  vpc_id      = flexibleengine_vpc_v1.test.id
  network_ids = [
    flexibleengine_vpc_subnet_v1.master.id,
  ]
}
`, testAccService_base(rName))
}

func testAccService_update(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "flexibleengine_workspace_service" "test" {
  access_mode = "INTERNET"
  vpc_id      = flexibleengine_vpc_v1.test.id
  network_ids = [
    flexibleengine_vpc_subnet_v1.master.id,
    flexibleengine_vpc_subnet_v1.standby.id,
  ]

  internet_access_port = 9001
  enterprise_id        = "%[2]s"
}
`, testAccService_base(rName), rName)
}

func testAccService_localAD_base(rName string) string {
	return fmt.Sprintf(`
data "flexibleengine_vpc_v1" "test" {
  id = "%[1]s"
}

resource "flexibleengine_vpc_subnet_v1" "master" {
  vpc_id = "%[1]s"

  name       = "%[2]s-master"
  cidr       = cidrsubnet(data.flexibleengine_vpc_v1.test.cidr, 4, 1)
  gateway_ip = cidrhost(cidrsubnet(data.flexibleengine_vpc_v1.test.cidr, 4, 1), 1)
}

resource "flexibleengine_vpc_subnet_v1" "standby" {
  vpc_id = "%[1]s"

  name       = "%[2]s-standby"
  cidr       = cidrsubnet(data.flexibleengine_vpc_v1.test.cidr, 4, 2)
  gateway_ip = cidrhost(cidrsubnet(data.flexibleengine_vpc_v1.test.cidr, 4, 2), 1)
}
`, OS_WORKSPACE_AD_VPC_ID, rName)
}

func testAccService_localAD_basic(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "flexibleengine_workspace_service" "test" {
  ad_domain {
    name               = "%[2]s"
    admin_account      = "Administrator"
    password           = "%[3]s"
    active_domain_ip   = "%[4]s"
    active_domain_name = "server.%[2]s"
    active_dns_ip      = "%[4]s"
  }

  auth_type   = "LOCAL_AD"
  access_mode = "INTERNET"
  vpc_id      = "%[5]s"
  network_ids = ["%[6]s"]
}
`, testAccService_localAD_base(rName), OS_WORKSPACE_AD_DOMAIN_NAME, OS_WORKSPACE_AD_SERVER_PWD,
		OS_WORKSPACE_AD_DOMAIN_IP, OS_WORKSPACE_AD_VPC_ID, OS_WORKSPACE_AD_NETWORK_ID)
}

func testAccService_localAD_update(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "flexibleengine_workspace_service" "test" {
  depends_on = [
    flexibleengine_vpc_subnet_v1.master,
	flexibleengine_vpc_subnet_v1.standby,
  ]

  ad_domain {
    name               = "%[2]s"
    admin_account      = "Administrator"
    password           = "%[3]s"
    active_domain_ip   = "%[4]s"
    active_domain_name = "server.%[2]s"
    active_dns_ip      = "%[4]s"
  }

  auth_type   = "LOCAL_AD"
  access_mode = "INTERNET"
  vpc_id      = "%[5]s"
  network_ids = [
    "%[6]s",
    flexibleengine_vpc_subnet_v1.master.id,
    flexibleengine_vpc_subnet_v1.standby.id,
  ]

  internet_access_port = 9001
  enterprise_id        = "%[7]s"
}
`, testAccService_localAD_base(rName), OS_WORKSPACE_AD_DOMAIN_NAME, OS_WORKSPACE_AD_SERVER_PWD,
		OS_WORKSPACE_AD_DOMAIN_IP, OS_WORKSPACE_AD_VPC_ID, OS_WORKSPACE_AD_NETWORK_ID,
		rName)
}
