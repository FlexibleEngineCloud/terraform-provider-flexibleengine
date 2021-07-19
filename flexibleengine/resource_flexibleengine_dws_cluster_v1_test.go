package flexibleengine

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/huaweicloud/golangsdk/openstack/dws/cluster"
)

func TestDWSClusterBasic(t *testing.T) {
	var ar cluster.Cluster

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testDWSClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testDWSClusterBasic,
				Check: resource.ComposeTestCheckFunc(
					testDWSClusterExists("flexibleengine_dws_cluster_v1.cluster", &ar),
				),
			},
		},
	})
}

func testDWSClusterDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	client, err := config.dwsV1Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine DWS client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "flexibleengine_dws_cluster_v1" {
			continue
		}

		id := rs.Primary.ID
		_, err := cluster.Get(client, id).Extract()
		if err == nil {
			return fmt.Errorf("Cluster still exists")
		}
	}

	return nil
}

func testDWSClusterExists(n string, ar *cluster.Cluster) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		client, err := config.dwsV1Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating FlexibleEngine DWS client: %s", err)
		}

		id := rs.Primary.ID
		found, err := cluster.Get(client, id).Extract()
		if err != nil {
			return err
		}

		*ar = *found

		return nil
	}
}

var testDWSClusterBasic = fmt.Sprintf(`
resource "flexibleengine_networking_secgroup_v2" "secgroup" {
  name = "terraform_security_group_test"
  description = "terraform security group acceptance test"
}

resource "flexibleengine_dws_cluster_v1" "cluster" {
  "node_type" = "dws.d1.xlarge"
  "number_of_node" = 3
  "subnet_id" = "%s"
  "vpc_id" = "%s"
  "security_group_id" = "${flexibleengine_networking_secgroup_v2.secgroup.id}"
  "availability_zone" = "%s"
  "name" = "terraform_dws_cluster_test"
  "user_name" = "test_cluster_admin"
  "user_pwd" = "cluster123@!"

  timeouts {
    create = "30m"
    delete = "30m"
  }
}
`, OS_NETWORK_ID, OS_VPC_ID, OS_AVAILABILITY_ZONE)
