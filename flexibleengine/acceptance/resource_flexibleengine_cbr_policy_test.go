package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/cbr/v3/policies"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func TestAccCBRV3Policy_basic(t *testing.T) {
	var cbrPolicy policies.Policy
	randName := fmt.Sprintf("tf_acc_test_%s", acctest.RandString(5))
	resourceName := "flexibleengine_cbr_policy.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: TestAccProviderFactories,
		CheckDestroy:      testAccCheckCBRPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testCBRV3Policy_basic(randName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCBRPolicyExists(resourceName, &cbrPolicy),
					resource.TestCheckResourceAttr(resourceName, "name", randName),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "type", "backup"),
					resource.TestCheckResourceAttr(resourceName, "time_period", "20"),
					resource.TestCheckResourceAttr(resourceName, "backup_cycle.0.days", "MO,TU"),
					resource.TestCheckResourceAttr(resourceName, "backup_cycle.0.execution_times.0", "06:00"),
					resource.TestCheckResourceAttr(resourceName, "backup_cycle.0.execution_times.1", "18:00"),
				),
			},
			{
				Config: testCBRV3Policy_update(randName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", randName+"-update"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "type", "backup"),
					resource.TestCheckResourceAttr(resourceName, "backup_quantity", "5"),
					resource.TestCheckResourceAttr(resourceName, "backup_cycle.0.days", "SA,SU"),
					resource.TestCheckResourceAttr(resourceName, "backup_cycle.0.execution_times.0", "08:00"),
					resource.TestCheckResourceAttr(resourceName, "backup_cycle.0.execution_times.1", "20:00"),
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

func TestAccCBRV3Policy_retention(t *testing.T) {
	var cbrPolicy policies.Policy
	randName := fmt.Sprintf("tf_acc_test_%s", acctest.RandString(5))
	resourceName := "flexibleengine_cbr_policy.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: TestAccProviderFactories,
		CheckDestroy:      testAccCheckCBRPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testCBRV3Policy_retention(randName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCBRPolicyExists(resourceName, &cbrPolicy),
					resource.TestCheckResourceAttr(resourceName, "name", randName),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "type", "backup"),
					resource.TestCheckResourceAttr(resourceName, "backup_quantity", "15"),
					resource.TestCheckResourceAttr(resourceName, "time_zone", "UTC+08:00"),
					resource.TestCheckResourceAttr(resourceName, "long_term_retention.0.daily", "10"),
					resource.TestCheckResourceAttr(resourceName, "long_term_retention.0.weekly", "10"),
					resource.TestCheckResourceAttr(resourceName, "long_term_retention.0.monthly", "1"),
					resource.TestCheckResourceAttr(resourceName, "long_term_retention.0.full_backup_interval", "-1"),
					resource.TestCheckResourceAttr(resourceName, "backup_cycle.0.days", "SA,SU"),
					resource.TestCheckResourceAttr(resourceName, "backup_cycle.0.execution_times.0", "08:00"),
					resource.TestCheckResourceAttr(resourceName, "backup_cycle.0.execution_times.1", "20:00"),
				),
			},
			{
				Config: testCBRV3Policy_retention_update(randName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", randName),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "type", "backup"),
					resource.TestCheckResourceAttr(resourceName, "backup_quantity", "35"),
					resource.TestCheckResourceAttr(resourceName, "time_zone", "UTC+08:00"),
					resource.TestCheckResourceAttr(resourceName, "long_term_retention.0.daily", "20"),
					resource.TestCheckResourceAttr(resourceName, "long_term_retention.0.weekly", "20"),
					resource.TestCheckResourceAttr(resourceName, "long_term_retention.0.monthly", "6"),
					resource.TestCheckResourceAttr(resourceName, "long_term_retention.0.yearly", "1"),
					resource.TestCheckResourceAttr(resourceName, "long_term_retention.0.full_backup_interval", "5"),
					resource.TestCheckResourceAttr(resourceName, "backup_cycle.0.days", "SA,SU"),
					resource.TestCheckResourceAttr(resourceName, "backup_cycle.0.execution_times.0", "08:00"),
					resource.TestCheckResourceAttr(resourceName, "backup_cycle.0.execution_times.1", "20:00"),
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

func TestAccCBRV3Policy_replication(t *testing.T) {
	var cbrPolicy policies.Policy
	randName := fmt.Sprintf("tf_acc_test_%s", acctest.RandString(5))
	resourceName := "flexibleengine_cbr_policy.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckReplication(t)
		},
		ProviderFactories: TestAccProviderFactories,
		CheckDestroy:      testAccCheckCBRPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testCBRV3Policy_replication(randName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCBRPolicyExists(resourceName, &cbrPolicy),
					resource.TestCheckResourceAttr(resourceName, "name", randName),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "type", "replication"),
					resource.TestCheckResourceAttr(resourceName, "destination_region", OS_DEST_REGION),
					resource.TestCheckResourceAttr(resourceName, "destination_project_id", OS_DEST_PROJECT_ID),
					resource.TestCheckResourceAttr(resourceName, "time_period", "20"),
					resource.TestCheckResourceAttr(resourceName, "enable_acceleration", "true"),
					resource.TestCheckResourceAttr(resourceName, "backup_cycle.0.interval", "5"),
					resource.TestCheckResourceAttr(resourceName, "backup_cycle.0.execution_times.0", "06:00"),
					resource.TestCheckResourceAttr(resourceName, "backup_cycle.0.execution_times.1", "18:00"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"enable_acceleration",
				},
			},
		},
	})
}

func testAccCheckCBRPolicyDestroy(s *terraform.State) error {
	conf := testAccProvider.Meta().(*config.Config)
	client, err := conf.CbrV3Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine CBR client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "flexibleengine_cbr_policy" {
			continue
		}

		_, err := policies.Get(client, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("CBR policy still exists")
		}
	}

	return nil
}

func testAccCheckCBRPolicyExists(n string, policy *policies.Policy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		conf := testAccProvider.Meta().(*config.Config)
		client, err := conf.CbrV3Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating FlexibleEngine CBR client: %s", err)
		}

		found, err := policies.Get(client, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("CBR policy not found")
		}

		*policy = *found

		return nil
	}
}

func testCBRV3Policy_basic(rName string) string {
	return fmt.Sprintf(`
resource "flexibleengine_cbr_policy" "test" {
  name        = "%s"
  type        = "backup"
  time_period = 20

  backup_cycle {
    days            = "MO,TU"
    execution_times = ["06:00", "18:00"]
  }
}
`, rName)
}

func testCBRV3Policy_update(rName string) string {
	return fmt.Sprintf(`
resource "flexibleengine_cbr_policy" "test" {
  name            = "%s-update"
  type            = "backup"
  backup_quantity = 5

  backup_cycle {
    days            = "SA,SU"
    execution_times = ["08:00", "20:00"]
  }
}
`, rName)
}

func testCBRV3Policy_retention(rName string) string {
	return fmt.Sprintf(`
resource "flexibleengine_cbr_policy" "test" {
  name            = "%s"
  type            = "backup"
  backup_quantity = 15

  time_zone       = "UTC+08:00"
  long_term_retention {
    daily                = 10
    weekly               = 10
    monthly              = 1
    full_backup_interval = -1	
  }

  backup_cycle {
    days            = "SA,SU"
    execution_times = ["08:00", "20:00"]
  }
}
`, rName)
}

func testCBRV3Policy_retention_update(rName string) string {
	return fmt.Sprintf(`
resource "flexibleengine_cbr_policy" "test" {
  name            = "%s"
  type            = "backup"
  backup_quantity = 35

  time_zone       = "UTC+08:00"
  long_term_retention {
    daily                = 20
    weekly               = 20
    monthly              = 6
    yearly               = 1
    full_backup_interval = 5
  }

  backup_cycle {
    days            = "SA,SU"
    execution_times = ["08:00", "20:00"]
  }
}
`, rName)
}

func testCBRV3Policy_replication(rName string) string {
	return fmt.Sprintf(`
resource "flexibleengine_cbr_policy" "test" {
  name                   = "%s"
  type                   = "replication"
  destination_region     = "%s"
  destination_project_id = "%s"
  time_period            = 20
  enable_acceleration    = true

  backup_cycle {
    interval        = 5
    execution_times = ["06:00", "18:00"]
  }
}
`, rName, OS_DEST_REGION, OS_DEST_PROJECT_ID)
}
