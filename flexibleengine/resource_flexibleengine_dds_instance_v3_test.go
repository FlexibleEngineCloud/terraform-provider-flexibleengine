package flexibleengine

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/huaweicloud/golangsdk/openstack/dds/v3/instances"
)

func TestAccDDSV3Instance_basic(t *testing.T) {
	var instance instances.Instance
	var instanceName = fmt.Sprintf("dds-%s", acctest.RandString(5))
	resourceName := "flexibleengine_dds_instance_v3.instance"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDDSV3InstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDDSInstanceV3Config_basic(instanceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDDSV3InstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "name", instanceName),
					resource.TestCheckResourceAttr(resourceName, "port", "8635"),
					resource.TestCheckResourceAttr(resourceName, "ssl", "true"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "tags.owner", "terraform"),
				),
			},
		},
	})
}

func testAccCheckDDSV3InstanceDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	client, err := config.ddsV3Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine DDS client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "flexibleengine_dds_instance_v3" {
			continue
		}

		opts := instances.ListInstanceOpts{
			Id: rs.Primary.ID,
		}
		allPages, err := instances.List(client, &opts).AllPages()
		if err != nil {
			return err
		}
		instances, err := instances.ExtractInstances(allPages)
		if err != nil {
			return err
		}

		if instances.TotalCount > 0 {
			return fmt.Errorf("Instance still exists. ")
		}
	}

	return nil
}

func testAccCheckDDSV3InstanceExists(n string, instance *instances.Instance) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s. ", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set. ")
		}

		config := testAccProvider.Meta().(*Config)
		client, err := config.ddsV3Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating FlexibleEngine DDS client: %s ", err)
		}

		opts := instances.ListInstanceOpts{
			Id: rs.Primary.ID,
		}
		allPages, err := instances.List(client, &opts).AllPages()
		if err != nil {
			return err
		}
		instances, err := instances.ExtractInstances(allPages)
		if err != nil {
			return err
		}
		if instances.TotalCount == 0 {
			return fmt.Errorf("Instance not found. ")
		}

		return nil
	}
}

func testAccDDSInstanceV3Config_basic(name string) string {
	return fmt.Sprintf(`
resource "flexibleengine_networking_secgroup_v2" "secgroup_1" {
  name = "secgroup_dds"
}

resource "flexibleengine_dds_instance_v3" "instance" {
  name              = "%s"
  region            = "%s"
  availability_zone = "%s"
  vpc_id            = "%s"
  subnet_id         = "%s"
  security_group_id = flexibleengine_networking_secgroup_v2.secgroup_1.id
  password          = "Test@123"
  mode              = "Sharding"

  datastore {
    type           = "DDS-Community"
    version        = "3.4"
    storage_engine = "wiredTiger"
  }
  flavor {
    type      = "mongos"
    num       = 2
    spec_code = "dds.mongodb.s3.medium.4.mongos"
  }
  flavor {
    type      = "shard"
    num       = 2
    storage   = "ULTRAHIGH"
    size      = 20
    spec_code = "dds.mongodb.s3.medium.4.shard"
  }
  flavor {
    type      = "config"
    num       = 1
    storage   = "ULTRAHIGH"
    size      = 20
    spec_code = "dds.mongodb.s3.large.2.config"
  }
  backup_strategy {
    start_time = "08:00-09:00"
    keep_days  = "8"
  }
  tags = {
    foo   = "bar"
    owner = "terraform"
  }
}`, name, OS_REGION_NAME, OS_AVAILABILITY_ZONE, OS_VPC_ID, OS_NETWORK_ID)
}
