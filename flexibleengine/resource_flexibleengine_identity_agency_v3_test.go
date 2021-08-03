package flexibleengine

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/golangsdk/openstack/identity/v3/agency"
)

func TestAccIdentityV3Agency_basic(t *testing.T) {
	var agency agency.Agency
	rName := fmt.Sprintf("acc-agency-%s", acctest.RandString(5))
	resourceName := "flexibleengine_identity_agency_v3.agency_1"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccIdentityV3AgencyPreCheck(t)
			testAccPreCheckAdminOnly(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIdentityV3AgencyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityAgency_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIdentityV3AgencyExists(resourceName, &agency),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "description", "This is a test service agency"),
					resource.TestCheckResourceAttr(resourceName, "delegated_service_name", "op_svc_obs"),
					resource.TestCheckResourceAttr(resourceName, "duration", "FOREVER"),
					resource.TestCheckResourceAttr(resourceName, "domain_roles.#", "1"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccIdentityAgency_update(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIdentityV3AgencyExists(resourceName, &agency),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "description", "This is a updated test service agency"),
					resource.TestCheckResourceAttr(resourceName, "delegated_service_name", "op_svc_evs"),
					resource.TestCheckResourceAttr(resourceName, "duration", "FOREVER"),
					resource.TestCheckResourceAttr(resourceName, "domain_roles.#", "1"),
				),
			},
		},
	})
}

func TestAccIdentityV3Agency_domain(t *testing.T) {
	var agency agency.Agency
	rName := fmt.Sprintf("acc-agency-%s", acctest.RandString(5))
	resourceName := "flexibleengine_identity_agency_v3.agency_1"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccIdentityV3AgencyPreCheck(t)
			testAccPreCheckAdminOnly(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIdentityV3AgencyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityAgency_domain(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIdentityV3AgencyExists(resourceName, &agency),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "description", "This is a test agency"),
					resource.TestCheckResourceAttr(resourceName, "delegated_domain_name", OS_DELEGATED_DOMAIN_NAME),
					resource.TestCheckResourceAttr(resourceName, "duration", "FOREVER"),
					resource.TestCheckResourceAttr(resourceName, "domain_roles.#", "1"),
				),
			},
			{
				Config: testAccIdentityAgency_domainUpdate(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIdentityV3AgencyExists(resourceName, &agency),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "description", "This is a updated test agency"),
					resource.TestCheckResourceAttr(resourceName, "delegated_domain_name", OS_DELEGATED_DOMAIN_NAME),
					resource.TestCheckResourceAttr(resourceName, "duration", "ONEDAY"),
					resource.TestCheckResourceAttr(resourceName, "domain_roles.#", "1"),
				),
			},
		},
	})
}

func testAccCheckIdentityV3AgencyDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	identityClient, err := config.identityV3Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine identity client: %s", err)
	}
	identityClient.Endpoint = strings.Replace(identityClient.Endpoint, "v3", "v3.0", 1)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "flexibleengine_identity_agency_v3" {
			continue
		}

		_, err := agency.Get(identityClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Agency still exists")
		}
	}

	return nil
}

func testAccCheckIdentityV3AgencyExists(n string, a *agency.Agency) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		identityClient, err := config.identityV3Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating FlexibleEngine identity client: %s", err)
		}
		identityClient.Endpoint = strings.Replace(identityClient.Endpoint, "v3", "v3.0", 1)
		found, err := agency.Get(identityClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Agency not found")
		}

		*a = *found

		return nil
	}
}

func testAccIdentityAgency_basic(rName string) string {
	return fmt.Sprintf(`
resource "flexibleengine_identity_agency_v3" "agency_1" {
  name                   = "%s"
  description            = "This is a test service agency"
  delegated_service_name = "op_svc_obs"

  domain_roles = [
    "OBS OperateAccess",
  ]
}
`, rName)
}

func testAccIdentityAgency_update(rName string) string {
	return fmt.Sprintf(`
resource "flexibleengine_identity_agency_v3" "agency_1" {
  name                   = "%s"
  description            = "This is a updated test service agency"
  delegated_service_name = "op_svc_evs"

  domain_roles = [
    "Anti-DDoS Administrator",
  ]
}
`, rName)
}

func testAccIdentityAgency_domain(rName string) string {
	return fmt.Sprintf(`
resource "flexibleengine_identity_agency_v3" "agency_1" {
  name                  = "%s"
  description           = "This is a test agency"
  delegated_domain_name = "%s"

  domain_roles = [
    "Anti-DDoS Administrator",
  ]
}
`, rName, OS_DELEGATED_DOMAIN_NAME)
}

func testAccIdentityAgency_domainUpdate(rName string) string {
	return fmt.Sprintf(`
resource "flexibleengine_identity_agency_v3" "agency_1" {
  name                  = "%s"
  description           = "This is a updated test agency"
  delegated_domain_name = "%s"
  duration              = "ONEDAY"

  domain_roles = [
    "Anti-DDoS Administrator",
  ]
}
`, rName, OS_DELEGATED_DOMAIN_NAME)
}
