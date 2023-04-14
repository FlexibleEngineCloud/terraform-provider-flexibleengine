package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/geminidb/v3/instances"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccGeminiDBInstance_basic(t *testing.T) {
	var instance instances.GeminiDBInstance

	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "flexibleengine_gaussdb_cassandra_instance.test"
	password := acceptance.RandomPassword()
	newPassword := acceptance.RandomPassword()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: TestAccProviderFactories,
		CheckDestroy:      testAccCheckGeminiDBInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccGeminiDBInstanceConfig_basic(rName, password),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGeminiDBInstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "password", password),
					resource.TestCheckResourceAttr(resourceName, "node_num", "3"),
					resource.TestCheckResourceAttr(resourceName, "volume_size", "100"),
					resource.TestCheckResourceAttr(resourceName, "status", "normal"),
				),
			},
			{
				Config: testAccGeminiDBInstanceConfig_update(rName, newPassword),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGeminiDBInstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "name", rName+"-update"),
					resource.TestCheckResourceAttr(resourceName, "password", newPassword),
					resource.TestCheckResourceAttr(resourceName, "node_num", "4"),
					resource.TestCheckResourceAttr(resourceName, "volume_size", "200"),
					resource.TestCheckResourceAttr(resourceName, "status", "normal"),
				),
			},
		},
	})
}

func testAccCheckGeminiDBInstanceDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*config.Config)
	client, err := config.GeminiDBV3Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating GeminiDB client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "flexibleengine_gaussdb_cassandra_instance" {
			continue
		}

		found, err := instances.GetInstanceByID(client, rs.Primary.ID)
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return nil
			}
			return err
		}
		if found.Id != "" {
			return fmt.Errorf("Instance <%s> still exists.", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckGeminiDBInstanceExists(n string, instance *instances.GeminiDBInstance) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s.", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set.")
		}

		config := testAccProvider.Meta().(*config.Config)
		client, err := config.GeminiDBV3Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating GeminiDB client: %s", err)
		}

		found, err := instances.GetInstanceByID(client, rs.Primary.ID)
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return fmt.Errorf("Instance <%s> not found.", rs.Primary.ID)
			}
			return err
		}
		if found.Id == "" {
			return fmt.Errorf("Instance <%s> not found.", rs.Primary.ID)
		}
		instance = &found

		return nil
	}
}

func testAccGaussDBNosql_base(rName string) string {
	return fmt.Sprintf(`
resource "flexibleengine_vpc_v1" "test" {
  name = "%[1]s"
  cidr = "192.168.0.0/16"
}

resource "flexibleengine_vpc_subnet_v1" "test" {
  name       = "%[1]s"
  vpc_id     = flexibleengine_vpc_v1.test.id
  cidr       = "192.168.0.0/24"
  gateway_ip = "192.168.0.1"
}

resource "flexibleengine_networking_secgroup_v2" "test" {
  name = "%[1]s"
}
`, rName)
}

func testAccGeminiDBInstanceConfig_basic(rName, password string) string {
	return fmt.Sprintf(`
%s

data "flexibleengine_availability_zones" "test" {}

data "flexibleengine_gaussdb_nosql_flavors" "test" {
  vcpus             = 4
  engine            = "cassandra"
}

resource "flexibleengine_gaussdb_cassandra_instance" "test" {
  name        = "%s"
  password    = "%s"
  flavor      = data.flexibleengine_gaussdb_nosql_flavors.test.flavors[0].name
  volume_size = 100
  vpc_id      = flexibleengine_vpc_v1.test.id
  subnet_id   = flexibleengine_vpc_subnet_v1.test.id
  ssl         = true
  node_num    = 3

  security_group_id = flexibleengine_networking_secgroup_v2.test.id
  availability_zone = join(",", data.flexibleengine_availability_zones.test.names)

  backup_strategy {
    start_time = "03:00-04:00"
    keep_days  = 14
  }

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, testAccGaussDBNosql_base(rName), rName, password)
}

func testAccGeminiDBInstanceConfig_update(rName, password string) string {
	return fmt.Sprintf(`
%s

data "flexibleengine_availability_zones" "test" {}

data "flexibleengine_gaussdb_nosql_flavors" "test" {
  vcpus             = 4
  engine            = "cassandra"
}

resource "flexibleengine_gaussdb_cassandra_instance" "test" {
  name        = "%s-update"
  password    = "%s"
  flavor      = data.flexibleengine_gaussdb_nosql_flavors.test.flavors[0].name
  volume_size = 200
  vpc_id      = flexibleengine_vpc_v1.test.id
  subnet_id   = flexibleengine_vpc_subnet_v1.test.id
  ssl         = true
  node_num    = 4

  security_group_id = flexibleengine_networking_secgroup_v2.test.id
  availability_zone = join(",", data.flexibleengine_availability_zones.test.names)

  backup_strategy {
    start_time = "03:00-04:00"
    keep_days  = 14
  }

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, testAccGaussDBNosql_base(rName), rName, password)
}
