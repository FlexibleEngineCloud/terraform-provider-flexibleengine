package flexibleengine

import (
	"fmt"
	"testing"

	"github.com/chnsz/golangsdk/openstack/compute/v2/servers"
	// "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccComputeInstancesDataSource_basic(t *testing.T) {
	dataSourceName := "data.flexibleengine_compute_instances.this"
	var instance servers.Server

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeV2InstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstancesDataSource_basic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2InstanceExists("flexibleengine_compute_instance_v2.instance_1", &instance),
					testAccCheckComputeInstanceDataSourceID(dataSourceName),
					resource.TestCheckResourceAttr(dataSourceName, "name", "instance_1"),
					resource.TestCheckResourceAttr(dataSourceName, "instances.#", "1"),
					resource.TestCheckResourceAttrSet(dataSourceName, "instances.0.id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "instances.0.image_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "instances.0.flavor_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "instances.0.flavor_name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "instances.0.status"),
					resource.TestCheckResourceAttrSet(dataSourceName, "instances.0.availability_zone"),
					resource.TestCheckResourceAttr(dataSourceName, "instances.0.security_groups.#", "1"),
				),
			},
		},
	})
}

func testAccCheckComputeInstancesDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find compute instance data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Compute instance data source ID not set")
		}

		return nil
	}
}

func testAccComputeInstancesDataSource_basic() string {
	return fmt.Sprintf(`
%s
resource "flexibleengine_compute_instance_v2" "instance_2" {
  name               = "instance_2"
  security_groups    = ["default"]
  availability_zone  = "%s"
  network {
    uuid = "%s"
  }
}
data "flexibleengine_compute_instances" "this" {
  name = flexibleengine_compute_instance_v2.instance_1.name
}
`, testAccComputeV2Instance_basic, OS_AVAILABILITY_ZONE, OS_NETWORK_ID)
}
