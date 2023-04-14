package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/geminidb/v3/instances"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func getNosqlInstance(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.GeminiDBV3Client(OS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("Error creating GaussDB client: %s", err)
	}

	found, err := instances.GetInstanceByID(client, state.Primary.ID)
	if err != nil {
		return nil, err
	}
	if found.Id == "" {
		return nil, fmt.Errorf("Instance <%s> not found.", state.Primary.ID)
	}

	return &found, nil
}

func TestAccGaussInfluxInstance_basic(t *testing.T) {
	var instance instances.GeminiDBInstance
	rName := acceptance.RandomAccResourceName()
	resourceName := "flexibleengine_gaussdb_influx_instance.test"
	password := acceptance.RandomPassword()
	newPassword := acceptance.RandomPassword()
	rc := acceptance.InitResourceCheck(
		resourceName,
		&instance,
		getNosqlInstance,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccGaussInfluxInstanceConfig_basic(rName, password),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "password", password),
					resource.TestCheckResourceAttr(resourceName, "node_num", "3"),
					resource.TestCheckResourceAttr(resourceName, "volume_size", "100"),
					resource.TestCheckResourceAttr(resourceName, "status", "normal"),
				),
			},
			{
				Config: testAccGaussInfluxInstanceConfig_update(rName, newPassword),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("%s-update", rName)),
					resource.TestCheckResourceAttr(resourceName, "password", newPassword),
					resource.TestCheckResourceAttr(resourceName, "node_num", "4"),
					resource.TestCheckResourceAttr(resourceName, "volume_size", "200"),
					resource.TestCheckResourceAttr(resourceName, "status", "normal"),
				),
			},
		},
	})
}

func testAccGaussInfluxInstanceConfig_basic(rName, password string) string {
	return fmt.Sprintf(`
%s

data "flexibleengine_availability_zones" "test" {}

data "flexibleengine_gaussdb_nosql_flavors" "test" {
  vcpus             = 4
  engine            = "influxdb"
}

resource "flexibleengine_gaussdb_influx_instance" "test" {
  name        = "%s"
  password    = "%s"
  flavor      = data.flexibleengine_gaussdb_nosql_flavors.test.flavors[0].name
  volume_size = 100
  vpc_id      = flexibleengine_vpc_v1.test.id
  subnet_id   = flexibleengine_vpc_subnet_v1.test.id
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

func testAccGaussInfluxInstanceConfig_update(rName, password string) string {
	return fmt.Sprintf(`
%s

data "flexibleengine_availability_zones" "test" {}

data "flexibleengine_gaussdb_nosql_flavors" "test" {
  vcpus             = 4
  engine            = "influxdb"
}

resource "flexibleengine_gaussdb_influx_instance" "test" {
  name        = "%s-update"
  password    = "%s"
  flavor      = data.flexibleengine_gaussdb_nosql_flavors.test.flavors[0].name
  volume_size = 200
  vpc_id      = flexibleengine_vpc_v1.test.id
  subnet_id   = flexibleengine_vpc_subnet_v1.test.id
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
