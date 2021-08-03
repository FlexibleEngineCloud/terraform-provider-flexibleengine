package flexibleengine

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/huaweicloud/golangsdk/openstack/compute/v2/servers"
)

func TestAccComputeInstanceDataSource_basic(t *testing.T) {
	resourceName := "data.flexibleengine_compute_instance_v2.this"
	var instance servers.Server

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeV2InstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstanceDataSource_basic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2InstanceExists("flexibleengine_compute_instance_v2.instance_1", &instance),
					testAccCheckComputeInstanceDataSourceID(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", "instance_1"),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
					resource.TestCheckResourceAttrSet(resourceName, "system_disk_id"),
					resource.TestCheckResourceAttrSet(resourceName, "security_groups.#"),
					resource.TestCheckResourceAttrSet(resourceName, "network.#"),
					resource.TestCheckResourceAttrSet(resourceName, "block_device.#"),
				),
			},
		},
	})
}

func testAccCheckComputeInstanceDataSourceID(n string) resource.TestCheckFunc {
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

func testAccComputeInstanceDataSource_basic() string {
	return fmt.Sprintf(`
%s

data "flexibleengine_compute_instance_v2" "this" {
  name = flexibleengine_compute_instance_v2.instance_1.name
}
`, testAccComputeV2Instance_basic)
}
