package flexibleengine

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/compute/v2/extensions/keypairs"
)

func TestAccComputeV2Keypair_basic(t *testing.T) {
	var keypair keypairs.KeyPair
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "flexibleengine_compute_keypair_v2.import"
	publicKey, _, _ := acctest.RandSSHKeyPair("Generated-by-AccTest")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeV2KeypairDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2Keypair_import(rName, publicKey),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2KeypairExists(resourceName, &keypair),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeV2Keypair_create(t *testing.T) {
	var keypair keypairs.KeyPair
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "flexibleengine_compute_keypair_v2.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeV2KeypairDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2Keypair_create(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2KeypairExists(resourceName, &keypair),
					resource.TestCheckResourceAttrSet(resourceName, "private_key_path"),
				),
			},
		},
	})
}

func testAccCheckComputeV2KeypairDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	computeClient, err := config.computeV2Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine compute client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "flexibleengine_compute_keypair_v2" {
			continue
		}

		_, err := keypairs.Get(computeClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Keypair still exists")
		}

		privateKey := rs.Primary.Attributes["private_key_path"]
		if privateKey != "" {
			if _, err := os.Stat(privateKey); err == nil {
				return fmt.Errorf("private key file still exists")
			}
		}
	}

	return nil
}

func testAccCheckComputeV2KeypairExists(n string, kp *keypairs.KeyPair) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		computeClient, err := config.computeV2Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating FlexibleEngine compute client: %s", err)
		}

		found, err := keypairs.Get(computeClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.Name != rs.Primary.ID {
			return fmt.Errorf("Keypair not found")
		}

		privateKey := rs.Primary.Attributes["private_key_path"]
		if privateKey != "" {
			if _, err := os.Stat(privateKey); err != nil {
				return fmt.Errorf("private key file not found: %s", err)
			}
		}

		*kp = *found

		return nil
	}
}

func testAccComputeV2Keypair_import(rName, keypair string) string {
	return fmt.Sprintf(`
resource "flexibleengine_compute_keypair_v2" "import" {
  name       = "%s"
  public_key = "%s"
}
`, rName, keypair)
}

func testAccComputeV2Keypair_create(rName string) string {
	return fmt.Sprintf(`
resource "flexibleengine_compute_keypair_v2" "test" {
  name = "%s"
}
`, rName)
}
