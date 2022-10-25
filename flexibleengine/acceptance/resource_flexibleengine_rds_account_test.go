package acceptance

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	model "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/rds/v3/model"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccRdsAccount_basic(t *testing.T) {
	rName := acceptance.RandomAccResourceName()
	resourceName := "flexibleengine_rds_account.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: TestAccProviderFactories,
		CheckDestroy:      testAccCheckRdsAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testRdsAccount_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRdsAccountExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "password", "Test@12345678"),
				),
			},
			{
				Config: testRdsAccount_update(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRdsAccountExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "password", "Test@123456789"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password"},
			},
		},
	})
}

func testAccCheckRdsAccountDestroy(s *terraform.State) error {
	c := testAccProvider.Meta().(*config.Config)
	client, err := c.HcRdsV3Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating RDS client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "flexibleengine_rds_account" {
			continue
		}

		// Split instance_id and user from resource id
		parts := strings.SplitN(rs.Primary.ID, "/", 2)
		if len(parts) != 2 {
			return fmt.Errorf("invalid id format, must be <instance_id>/<user>")
		}
		instanceId := parts[0]
		userName := parts[1]
		// items on every page, [1, 100]
		limit := int32(100)
		// List all db users
		request := &model.ListDbUsersRequest{
			InstanceId: instanceId,
			Limit:      limit,
			Page:       int32(1),
		}

		for {
			response, err := client.ListDbUsers(request)
			if err != nil {
				return nil
			}
			users := *response.Users
			if len(users) == 0 {
				break
			}
			request.Page += 1
			for _, user := range users {
				if user.Name == userName {
					return fmt.Errorf("Rds account still exists")
				}
			}
		}
	}

	return nil
}

func testAccCheckRdsAccountExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		c := testAccProvider.Meta().(*config.Config)
		client, err := c.HcRdsV3Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating RDS client: %s", err)
		}

		// Split instance_id and user from resource id
		parts := strings.SplitN(rs.Primary.ID, "/", 2)
		if len(parts) != 2 {
			return fmt.Errorf("invalid id format, must be <instance_id>/<user>")
		}
		instanceId := parts[0]
		userName := parts[1]
		// items on every page, [1, 100]
		limit := int32(100)
		// List all db users
		request := &model.ListDbUsersRequest{
			InstanceId: instanceId,
			Limit:      limit,
			Page:       int32(1),
		}

		for {
			response, err := client.ListDbUsers(request)
			if err != nil {
				return fmt.Errorf("error listing RDS db users: %s", err)
			}
			users := *response.Users
			if len(users) == 0 {
				break
			} else {
				request.Page += 1
				for _, user := range users {
					if user.Name == userName {
						return nil
					}
				}
			}
		}

		return fmt.Errorf("rds account not found")
	}
}

func testRdsAccount_base(rName string) string {
	return fmt.Sprintf(`
data "flexibleengine_availability_zones" "test" {}

resource "flexibleengine_vpc_v1" "test" {
  name = "%[1]s"
  cidr = "192.168.0.0/16"
}

resource "flexibleengine_vpc_subnet_v1" "test" {
  name       = "%[1]s"
  cidr       = "192.168.0.0/24"
  gateway_ip = "192.168.0.1"
  vpc_id     = flexibleengine_vpc_v1.test.id

  timeouts {
	delete = "20m"
  }
}

resource "flexibleengine_networking_secgroup_v2" "test" {
  name = "%[1]s"
}

resource "flexibleengine_rds_instance_v3" "test" {
  name                = "%[1]s"
  flavor              = "rds.mysql.s3.large.2.ha"
  security_group_id   = flexibleengine_networking_secgroup_v2.test.id
  subnet_id           = flexibleengine_vpc_subnet_v1.test.id
  vpc_id              = flexibleengine_vpc_v1.test.id
  fixed_ip            = "192.168.0.58"
  ha_replication_mode = "semisync"

  availability_zone = [
    data.flexibleengine_availability_zones.test.names[0],
    data.flexibleengine_availability_zones.test.names[1],
  ]

  db {
    password = "Huangwei!120521"
    type     = "MySQL"
    version  = "5.7"
    port     = 3306
  }

  volume {
    type = "ULTRAHIGH"
    size = 50
  }
}
`, rName)
}

func testRdsAccount_basic(rName string) string {
	return fmt.Sprintf(`
%s

resource "flexibleengine_rds_account" "test" {
  instance_id = flexibleengine_rds_instance_v3.test.id
  name        = "%s"
  password    = "Test@12345678"
}
`, testRdsAccount_base(rName), rName)
}

func testRdsAccount_update(rName string) string {
	return fmt.Sprintf(`
%s

resource "flexibleengine_rds_account" "test" {
  instance_id = flexibleengine_rds_instance_v3.test.id
  name        = "%s"
  password    = "Test@123456789"
}
`, testRdsAccount_base(rName), rName)
}
