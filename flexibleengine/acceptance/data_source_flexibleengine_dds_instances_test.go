package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDatasourceDdsInstance_basic(t *testing.T) {
	rName := "data.flexibleengine_dds_instances.test"
	name := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(rName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceDdsInstance_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "instances.0.name", name),
					resource.TestCheckResourceAttr(rName, "instances.0.mode", "Sharding"),
				),
			},
		},
	})
}

func testAccDatasourceDdsInstance_base(rName string) string {
	return fmt.Sprintf(`
%s

data "flexibleengine_availability_zones" "test" {}

resource "flexibleengine_dds_instance_v3" "test" {
  name              = "%s"
  availability_zone = data.flexibleengine_availability_zones.test.names[0]
  vpc_id            = flexibleengine_vpc_v1.test.id
  subnet_id         = flexibleengine_vpc_subnet_v1.test.id
  security_group_id = flexibleengine_networking_secgroup_v2.test.id
  password          = "Terraform@123"
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
}`, testBaseNetwork(rName), rName)
}

func testAccDatasourceDdsInstance_basic(name string) string {
	return fmt.Sprintf(`
%s

data "flexibleengine_dds_instances" "test" {
  name = flexibleengine_dds_instance_v3.test.name
}
`, testAccDatasourceDdsInstance_base(name))
}
