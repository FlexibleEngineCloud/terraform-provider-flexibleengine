package flexibleengine

import (
	"fmt"
	"testing"

	"github.com/chnsz/golangsdk/openstack/mls/v1/instances"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccMLSV1Instance_basic(t *testing.T) {
	var instance instances.Instance

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckMrs(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMLSV1InstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: TestAccMLSInstanceV1Config_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMLSV1InstanceExists("flexibleengine_mls_instance_v1.instance", &instance),
					resource.TestCheckResourceAttr(
						"flexibleengine_mls_instance_v1.instance", "status", "AVAILABLE"),
				),
			},
		},
	})
}

func testAccCheckMLSV1InstanceDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	mlsClient, err := config.MlsV1Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine mls: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "flexibleengine_mls_instance_v1" {
			continue
		}

		_, err := instances.Get(mlsClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Instance still exists. ")
		}
	}

	return nil
}

func testAccCheckMLSV1InstanceExists(n string, instance *instances.Instance) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s. ", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set. ")
		}

		config := testAccProvider.Meta().(*Config)
		mlsClient, err := config.MlsV1Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating FlexibleEngine mls client: %s ", err)
		}

		found, err := instances.Get(mlsClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Instance not found. ")
		}

		*instance = *found
		return nil
	}
}

var TestAccMLSInstanceV1Config_basic = fmt.Sprintf(`
%s

resource "flexibleengine_mls_instance_v1" "instance" {
  name    = "terraform-mls-instance"
  version = "1.2.0"
  flavor  = "mls.c2.2xlarge.common"
  network {
    vpc_id         = "%s"
    subnet_id      = "%s"
	available_zone = "%s"
	public_ip {
	  bind_type = "not_use"
	}
  }
  mrs_cluster {
    id = flexibleengine_mrs_cluster_v1.cluster1.id
  }
}`, testAccMRSV1ClusterConfig_basic, OS_VPC_ID, OS_NETWORK_ID, OS_AVAILABILITY_ZONE)
