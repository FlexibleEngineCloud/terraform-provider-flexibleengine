package flexibleengine

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/identity/federatedauth/mappings"
)

func TestAccIdentityProviderConversion_basic(t *testing.T) {
	var mappingRules mappings.IdentityMapping
	var name = fmt.Sprintf("idp-ACCPTTEST-%s", acctest.RandString(5))
	resourceName := "flexibleengine_identity_provider_conversion.conversion"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
		},
		ProviderFactories: TestAccProviderFactories,
		CheckDestroy:      testAccCheckIdentityProviderConversionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityProviderConversion_conf(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIdentityProviderConversionExists(resourceName, &mappingRules),
					resource.TestCheckResourceAttr(resourceName, "conversion_rules.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "conversion_rules.0.local.0.username", "Tom"),
				),
			},
			{
				Config: testAccIdentityProviderConversion_update(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "conversion_rules.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "conversion_rules.0.local.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "conversion_rules.0.local.0.username", "Tom"),
					resource.TestCheckResourceAttr(resourceName, "conversion_rules.1.remote.0.value.#", "2"),
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

func testAccCheckIdentityProviderConversionDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	identityClient, err := config.IAMNoVersionClient(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine identity client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "flexibleengine_identity_provider_conversion" {
			continue
		}

		_, err := mappings.Get(identityClient, rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("Identity Provider Conversion still exists")
		}
	}

	return nil
}

func testAccCheckIdentityProviderConversionExists(n string, rules *mappings.IdentityMapping) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		identityClient, err := config.IAMNoVersionClient(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating FlexibleEngine identity client: %s", err)
		}

		found, err := mappings.Get(identityClient, rs.Primary.ID)
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Identity Provider Conversion not found")
		}

		*rules = *found

		return nil
	}
}

func testAccIdentityProviderConversion_conf(name string) string {
	return fmt.Sprintf(`
resource "flexibleengine_identity_provider" "provider_1" {
  name     = "%s"
  protocol = "oidc"
}

resource "flexibleengine_identity_provider_conversion" "conversion" {
  provider_id = flexibleengine_identity_provider.provider_1.id

  conversion_rules {
    local {
      username = "Tom"
    }
    remote {
      attribute = "Tom"
    }
  }
}
`, name)
}

func testAccIdentityProviderConversion_update(name string) string {
	return fmt.Sprintf(`
resource "flexibleengine_identity_provider" "provider_1" {
  name     = "%s"
  protocol = "oidc"
}

resource "flexibleengine_identity_provider_conversion" "conversion" {
  provider_id = flexibleengine_identity_provider.provider_1.id

  conversion_rules {
    local {
      username = "Tom"
    }
    local {
      username = "federateduser"
    }
    remote {
      attribute = "Tom"
    }
    remote {
      attribute = "federatedgroup"
    }
  }

  conversion_rules {
    local {
      username = "Jams"
    }
    remote {
      attribute = "username"
      condition = "any_one_of"
      value     = ["Tom", "Jerry"]
    }
  }
}
`, name)
}
