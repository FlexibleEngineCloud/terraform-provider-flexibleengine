package flexibleengine

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/identity/federatedauth/providers"
)

func TestAccIdentityProvider_basic(t *testing.T) {
	var provider providers.Provider
	var name = fmt.Sprintf("idp-ACCPTTEST-%s", acctest.RandString(5))
	resourceName := "flexibleengine_identity_provider.provider_1"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
		},
		ProviderFactories: TestAccProviderFactories,
		CheckDestroy:      testAccCheckIdentityProviderDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityProvider_saml(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIdentityProviderExists(resourceName, &provider),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "protocol", "saml"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
				),
			},
			{
				Config: testAccIdentityProvider_saml_update(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "protocol", "saml"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "false"),
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

func TestAccIdentityProvider_oidc(t *testing.T) {
	var provider providers.Provider
	var name = fmt.Sprintf("idp-ACCPTTEST-%s", acctest.RandString(5))
	resourceName := "flexibleengine_identity_provider.provider_1"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
		},
		ProviderFactories: TestAccProviderFactories,
		CheckDestroy:      testAccCheckIdentityProviderDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityProvider_oidc(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIdentityProviderExists(resourceName, &provider),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "protocol", "oidc"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "openid_connect_config.0.access_type", "program_console"),
					resource.TestCheckResourceAttr(resourceName, "openid_connect_config.0.client_id", "client_id_example"),
				),
			},
			{
				Config: testAccIdentityProvider_oidc_update(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "protocol", "oidc"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "openid_connect_config.0.access_type", "program"),
					resource.TestCheckResourceAttr(resourceName, "openid_connect_config.0.client_id", "client_id_demo"),
				),
			},
		},
	})
}

func testAccCheckIdentityProviderDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	identityClient, err := config.IAMNoVersionClient(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine identity client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "flexibleengine_identity_provider" {
			continue
		}

		_, err := providers.Get(identityClient, rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("Identity Provider still exists")
		}
	}

	return nil
}

func testAccCheckIdentityProviderExists(n string, idp *providers.Provider) resource.TestCheckFunc {
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

		found, err := providers.Get(identityClient, rs.Primary.ID)
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Identity Provider not found")
		}

		*idp = *found

		return nil
	}
}

func testAccIdentityProvider_saml(name string) string {
	return fmt.Sprintf(`
resource "flexibleengine_identity_provider" "provider_1" {
  name     = "%s"
  protocol = "saml"
}
`, name)
}

func testAccIdentityProvider_saml_update(name string) string {
	return fmt.Sprintf(`
resource "flexibleengine_identity_provider" "provider_1" {
  name     = "%s"
  protocol = "saml"
  enabled  = false
}
`, name)
}

func testAccIdentityProvider_oidc(name string) string {
	return fmt.Sprintf(`
resource "flexibleengine_identity_provider" "provider_1" {
  name        = "%s"
  protocol    = "oidc"
  description = "unit test"

  openid_connect_config {
    access_type            = "program_console"
    provider_url           = "https://accounts.example.com"
    client_id              = "client_id_example"
    authorization_endpoint = "https://accounts.example.com/o/oauth2/v2/auth"
    scopes                 = ["openid"]
    signing_key            = jsonencode(
    {
      keys = [
        {
          alg = "RS256"
          e   = "AQAB"
          kid = "d05ef20c4512645vv1..."
          kty = "RSA"
          n   = "cws_cnjiwsbvweolwn_-vnl..."
          use = "sig"
        },
      ]
    }
    )
  }
}
`, name)
}

func testAccIdentityProvider_oidc_update(name string) string {
	return fmt.Sprintf(`
resource "flexibleengine_identity_provider" "provider_1" {
  name        = "%s"
  protocol    = "oidc"
  enabled     = false
  description = "acceptance test"

  openid_connect_config {
    access_type  = "program"
    provider_url = "https://accounts.example.com"
    client_id    = "client_id_demo"
    signing_key  = jsonencode(
    {
      keys = [
        {
          alg = "RS256"
          e   = "AQAB"
          kid = "d05ef20c4512645vv1..."
          kty = "RSA"
          n   = "cws_cnjiwsbvweolwn_-vnl..."
          use = "sig"
        },
      ]
    }
    )
  }
}
`, name)
}
