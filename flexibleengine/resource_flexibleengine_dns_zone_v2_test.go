package flexibleengine

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/golangsdk/openstack/dns/v2/zones"
)

func TestAccDNSV2Zone_basic(t *testing.T) {
	var zone zones.Zone
	var zoneName = fmt.Sprintf("acpttest%s.com.", acctest.RandString(5))
	resourceName := "flexibleengine_dns_zone_v2.zone_1"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDNSV2ZoneDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDNSV2Zone_basic(zoneName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDNSV2ZoneExists(resourceName, &zone),
					resource.TestCheckResourceAttr(resourceName, "name", zoneName),
					resource.TestCheckResourceAttr(resourceName, "zone_type", "public"),
					resource.TestCheckResourceAttr(resourceName, "description", "a zone"),
					resource.TestCheckResourceAttr(resourceName, "email", "email1@example.com"),
					resource.TestCheckResourceAttr(resourceName, "ttl", "300"),
					resource.TestCheckResourceAttr(resourceName, "tags.zone_type", "public"),
					resource.TestCheckResourceAttr(resourceName, "tags.owner", "terraform"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccDNSV2Zone_update(zoneName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "description", "an updated zone"),
					resource.TestCheckResourceAttr(resourceName, "email", "email2@example.com"),
					resource.TestCheckResourceAttr(resourceName, "ttl", "600"),
					resource.TestCheckResourceAttr(resourceName, "tags.zone_type", "public"),
					resource.TestCheckResourceAttr(resourceName, "tags.owner", "tf-acc"),
				),
			},
		},
	})
}

func TestAccDNSV2Zone_private(t *testing.T) {
	var zone zones.Zone
	var zoneName = fmt.Sprintf("acpttest%s.com.", acctest.RandString(5))
	resourceName := "flexibleengine_dns_zone_v2.zone_1"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDNSV2ZoneDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDNSV2Zone_private(zoneName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDNSV2ZoneExists(resourceName, &zone),
					resource.TestCheckResourceAttr(resourceName, "name", zoneName),
					resource.TestCheckResourceAttr(resourceName, "zone_type", "private"),
					resource.TestCheckResourceAttr(resourceName, "description", "a private zone"),
					resource.TestCheckResourceAttr(resourceName, "email", "email@example.com"),
					resource.TestCheckResourceAttr(resourceName, "ttl", "300"),
					resource.TestCheckResourceAttr(resourceName, "tags.zone_type", "private"),
					resource.TestCheckResourceAttr(resourceName, "tags.owner", "terraform"),
				),
			},
		},
	})
}

func TestAccDNSV2Zone_readTTL(t *testing.T) {
	var zone zones.Zone
	var zoneName = fmt.Sprintf("acpttest%s.com.", acctest.RandString(5))
	resourceName := "flexibleengine_dns_zone_v2.zone_1"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDNSV2ZoneDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDNSV2Zone_readTTL(zoneName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDNSV2ZoneExists(resourceName, &zone),
					resource.TestMatchResourceAttr(resourceName, "ttl", regexp.MustCompile("^[0-9]+$")),
				),
			},
		},
	})
}

func testAccCheckDNSV2ZoneDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	dnsClient, err := config.dnsV2Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine DNS client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "flexibleengine_dns_zone_v2" {
			continue
		}

		_, err := zones.Get(dnsClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Zone still exists")
		}
	}

	return nil
}

func testAccCheckDNSV2ZoneExists(n string, zone *zones.Zone) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		dnsClient, err := config.dnsV2Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating FlexibleEngine DNS client: %s", err)
		}

		found, err := zones.Get(dnsClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Zone not found")
		}

		*zone = *found

		return nil
	}
}

func testAccDNSV2Zone_basic(zoneName string) string {
	return fmt.Sprintf(`
resource "flexibleengine_dns_zone_v2" "zone_1" {
  name        = "%s"
  email       = "email1@example.com"
  description = "a zone"
  ttl         = 300

  tags = {
    zone_type = "public"
    owner     = "terraform"
  }
}
	`, zoneName)
}

func testAccDNSV2Zone_update(zoneName string) string {
	return fmt.Sprintf(`
resource "flexibleengine_dns_zone_v2" "zone_1" {
  name        = "%s"
  email       = "email2@example.com"
  description = "an updated zone"
  ttl         = 600

  tags = {
    zone_type = "public"
    owner     = "tf-acc"
  }
}
	`, zoneName)
}

func testAccDNSV2Zone_readTTL(zoneName string) string {
	return fmt.Sprintf(`
resource "flexibleengine_dns_zone_v2" "zone_1" {
  name  = "%s"
  email = "email1@example.com"
}
	`, zoneName)
}

func testAccDNSV2Zone_private(rName string) string {
	return fmt.Sprintf(`
resource "flexibleengine_vpc_v1" "vpc_1" {
  name = "%s"
  cidr = "192.168.0.0/16"
}

resource "flexibleengine_dns_zone_v2" "zone_1" {
  name        = "%s"
  email       = "email@example.com"
  description = "a private zone"
  zone_type   = "private"

  router {
    router_id = flexibleengine_vpc_v1.vpc_1.id
  }
  tags = {
    zone_type = "private"
    owner     = "terraform"
  }
}
	`, rName, rName)
}
