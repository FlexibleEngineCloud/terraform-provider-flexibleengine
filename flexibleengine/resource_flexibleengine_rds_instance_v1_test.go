package flexibleengine

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/huaweicloud/golangsdk/openstack/rds/v1/instances"
)

func TestAccRDSV1Instance_basic(t *testing.T) {
	var instance instances.Instance

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRDSV1InstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: TestAccSInstanceV1Config_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRDSV1InstanceExists("flexibleengine_rds_instance_v1.instance", &instance),
					resource.TestCheckResourceAttr(
						"flexibleengine_rds_instance_v1.instance", "status", "ACTIVE"),
				),
			},
		},
	})
}

func TestAccRDSV1Instance_PostgreSQL(t *testing.T) {
	var instance instances.Instance

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRDSV1InstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: TestAccSInstanceV1Config_PostgreSQL,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRDSV1InstanceExists("flexibleengine_rds_instance_v1.instance", &instance),
					resource.TestCheckResourceAttr(
						"flexibleengine_rds_instance_v1.instance", "status", "ACTIVE"),
				),
			},
		},
	})
}

func testAccCheckRDSV1InstanceDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	rdsClient, err := config.rdsV1Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine rds: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "flexibleengine_rds_instance_v1" {
			continue
		}

		_, err := instances.Get(rdsClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Instance still exists. ")
		}
	}

	return nil
}

func testAccCheckRDSV1InstanceExists(n string, instance *instances.Instance) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s. ", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set. ")
		}

		config := testAccProvider.Meta().(*Config)
		rdsClient, err := config.rdsV1Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating FlexibleEngine rds client: %s ", err)
		}

		found, err := instances.Get(rdsClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Instance not found. ")
		}

		*instance = *found
		time.Sleep(30 * time.Second)

		return nil
	}
}

var TestAccSInstanceV1Config_basic = fmt.Sprintf(`
data "flexibleengine_rds_flavors_v1" "flavor" {
    region = "eu-west-0"
    datastore_name = "MySQL"
    datastore_version = "5.6.30"
    speccode = "rds.mysql.s1.medium.ha"
}

resource "flexibleengine_networking_secgroup_v2" "secgrp_rds" {
  name        = "secgrp-rds-instance"
  description = "Rds Security Group"
}

resource "flexibleengine_rds_instance_v1" "instance" {
  name = "rds-instance"
  datastore {
    type = "MySQL"
    version = "5.6.30"
  }
  flavorref = "${data.flexibleengine_rds_flavors_v1.flavor.id}"
  volume {
    type = "COMMON"
    size = 100
  }
  region = "eu-west-0"
  availabilityzone = "eu-west-0a"
  vpc = "%s"
  nics {
    subnetid = "%s"
  }
  securitygroup {
    id = "${flexibleengine_networking_secgroup_v2.secgrp_rds.id}"
  }
  dbport = "8635"
  backupstrategy {
    starttime = "00:00:00"
    keepdays = 0
  }
  dbrtpd = "Huangwei-120521"
  ha {
    enable = true
    replicationmode = "async"
  }
  depends_on = ["flexibleengine_networking_secgroup_v2.secgrp_rds"]
}`, OS_VPC_ID, OS_NETWORK_ID)

var TestAccSInstanceV1Config_PostgreSQL = fmt.Sprintf(`
data "flexibleengine_rds_flavors_v1" "flavor" {
    region = "eu-west-0"
    datastore_name = "PostgreSQL"
    datastore_version = "9.5.5"
    speccode = "rds.pg.s1.medium.ha"
}

resource "flexibleengine_networking_secgroup_v2" "secgrp_rds" {
  name        = "secgrp-rds-instance"
  description = "Rds Security Group"
}

resource "flexibleengine_rds_instance_v1" "instance" {
  name = "rds-instance"
  datastore {
    type = "PostgreSQL"
    version = "9.5.5"
  }
  flavorref = "${data.flexibleengine_rds_flavors_v1.flavor.id}"
  volume {
    type = "COMMON"
    size = 100
  }
  region = "eu-west-0"
  availabilityzone = "eu-west-0a"
  vpc = "%s"
  nics {
    subnetid = "%s"
  }
  securitygroup {
    id = "${flexibleengine_networking_secgroup_v2.secgrp_rds.id}"
  }
  dbport = "8635"
  backupstrategy {
    starttime = "00:00:00"
    keepdays = 0
  }
  dbrtpd = "Huangwei-120521"
  ha {
    enable = true
    replicationmode = "async"
  }
  depends_on = ["flexibleengine_networking_secgroup_v2.secgrp_rds"]
}`, OS_VPC_ID, OS_NETWORK_ID)
