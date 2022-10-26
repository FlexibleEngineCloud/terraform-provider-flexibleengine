package acceptance

import (
	"fmt"
	"testing"

	"github.com/chnsz/golangsdk/openstack/dds/v3/roles"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func getDatabaseRoleFunc(c *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := c.DdsV3Client(OS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating DDS v3 client: %s ", err)
	}

	instanceId := state.Primary.Attributes["instance_id"]
	name := state.Primary.Attributes["name"]
	opts := roles.ListOpts{
		Name:   state.Primary.Attributes["name"],
		DbName: state.Primary.Attributes["db_name"],
	}
	resp, err := roles.List(client, instanceId, opts)
	if err != nil {
		return nil, fmt.Errorf("error getting role (%s) from DDS instance (%s): %v", name, instanceId, err)
	}
	if len(resp) < 1 {
		return nil, fmt.Errorf("unable to find role (%s) from DDS instance (%s)", name, instanceId)
	}
	role := resp[0]
	return &role, nil
}

func TestAccDatabaseRole_basic(t *testing.T) {
	var role roles.Role
	rName := acceptance.RandomAccResourceName()
	resourceName := "flexibleengine_dds_database_role.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&role,
		getDatabaseRoleFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDatabaseRole_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttrPair(resourceName, "roles.0.name",
						"flexibleengine_dds_database_role.base", "name"),
					resource.TestCheckResourceAttrPair(resourceName, "inherited_privileges",
						"flexibleengine_dds_database_role.base", "inherited_privileges"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccDatabaseRoleImportStateIdFunc(),
			},
		},
	})
}

func testAccDatabaseRoleImportStateIdFunc() resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		var instanceId, dbName, name string
		for _, rs := range s.RootModule().Resources {
			if rs.Type == "flexibleengine_dds_database_role" {
				instanceId = rs.Primary.Attributes["instance_id"]
				dbName = rs.Primary.Attributes["db_name"]
				name = rs.Primary.Attributes["name"]
			}
		}
		if instanceId == "" || dbName == "" || name == "" {
			return "", fmt.Errorf("resource not found: %s/%s/%s", instanceId, dbName, name)
		}
		return fmt.Sprintf("%s/%s/%s", instanceId, dbName, name), nil
	}
}

func testAccDatabaseRole_base(rName string) string {
	randCidr := acceptance.RandomCidr()

	return fmt.Sprintf(`
data "flexibleengine_availability_zones" "test" {}

resource "flexibleengine_vpc_v1" "test" {
  name = "%[1]s"
  cidr = "%[2]s"
}

resource "flexibleengine_vpc_subnet_v1" "test" {
  vpc_id = flexibleengine_vpc_v1.test.id

  name       = "%[1]s"
  cidr       = cidrsubnet(flexibleengine_vpc_v1.test.cidr, 4, 1)
  gateway_ip = cidrhost(cidrsubnet(flexibleengine_vpc_v1.test.cidr, 4, 1), 1)

  timeouts {
    delete = "20m"
  }
}

resource "flexibleengine_networking_secgroup_v2" "test" {
  name = "%[1]s"
}

resource "flexibleengine_dds_instance_v3" "test" {
  availability_zone = data.flexibleengine_availability_zones.test.names[0]
  vpc_id            = flexibleengine_vpc_v1.test.id
  subnet_id         = flexibleengine_vpc_subnet_v1.test.id
  security_group_id = flexibleengine_networking_secgroup_v2.test.id

  name     = "%[1]s"
  mode     = "Sharding"
  password = "Test@12345678"

  datastore {
    type           = "DDS-Community"
    version        = "3.4"
    storage_engine = "wiredTiger"
  }

  flavor {
    type      = "mongos"
    num       = 2
    spec_code = "dds.mongodb.c6.large.2.mongos"
  }
  flavor {
    type      = "shard"
    num       = 2
    storage   = "ULTRAHIGH"
    size      = 20
    spec_code = "dds.mongodb.c6.large.2.shard"
  }
  flavor {
    type      = "config"
    num       = 1
    storage   = "ULTRAHIGH"
    size      = 20
    spec_code = "dds.mongodb.c6.large.2.config"
  }
}
`, rName, randCidr)
}

func testAccDatabaseRole_basic(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "flexibleengine_dds_database_role" "base" {
  instance_id = flexibleengine_dds_instance_v3.test.id

  name    = "%[2]s-base"
  db_name = "admin"

  timeouts {
    create = "10m"
    delete = "10m"
  }
}

resource "flexibleengine_dds_database_role" "test" {
  instance_id = flexibleengine_dds_instance_v3.test.id

  name    = "%[2]s"
  db_name = "admin"

  roles {
    name    = flexibleengine_dds_database_role.base.name
    db_name = "admin"
  }

  timeouts {
    create = "10m"
    delete = "10m"
  }
}
`, testAccDatabaseRole_base(rName), rName)
}
