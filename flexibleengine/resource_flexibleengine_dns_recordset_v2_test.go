package flexibleengine

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/dns/v2/recordsets"
)

func randomZoneName() string {
	return fmt.Sprintf("acpttest-zone-%s.com.", acctest.RandString(5))
}

func TestAccDNSV2RecordSet_basic(t *testing.T) {
	var recordset recordsets.RecordSet
	zoneName := randomZoneName()
	resourceName := "flexibleengine_dns_recordset_v2.recordset_1"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDNSV2RecordSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDNSV2RecordSet_basic(zoneName, 3000),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDNSV2RecordSetExists(resourceName, &recordset),
					resource.TestCheckResourceAttr(resourceName, "name", zoneName),
					resource.TestCheckResourceAttr(resourceName, "description", "a record set"),
					resource.TestCheckResourceAttr(resourceName, "type", "A"),
					resource.TestCheckResourceAttr(resourceName, "ttl", "3000"),
					resource.TestCheckResourceAttr(resourceName, "records.#", "2"),
				),
			},
			{
				// only update tags
				Config: testAccDNSV2RecordSet_tags(zoneName, 3000),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", zoneName),
					resource.TestCheckResourceAttr(resourceName, "description", "a record set"),
					resource.TestCheckResourceAttr(resourceName, "ttl", "3000"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
				),
			},
			{
				// only update ttl
				Config: testAccDNSV2RecordSet_tags(zoneName, 6000),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "ttl", "6000"),
					resource.TestCheckResourceAttr(resourceName, "records.#", "2"),
				),
			},
			{
				// update ttl, description, records and tags
				Config: testAccDNSV2RecordSet_update(zoneName, 5000),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "description", "an updated record set"),
					resource.TestCheckResourceAttr(resourceName, "ttl", "5000"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value_updated"),
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

func TestAccDNSV2RecordSet_readTTL(t *testing.T) {
	var recordset recordsets.RecordSet
	zoneName := randomZoneName()
	resourceName := "flexibleengine_dns_recordset_v2.recordset_1"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDNSV2RecordSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDNSV2RecordSet_readTTL(zoneName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDNSV2RecordSetExists(resourceName, &recordset),
					resource.TestMatchResourceAttr(resourceName, "ttl", regexp.MustCompile("^[0-9]+$")),
				),
			},
		},
	})
}

func TestAccDNSV2RecordSet_private(t *testing.T) {
	var recordset recordsets.RecordSet
	zoneName := randomZoneName()
	rName := fmt.Sprintf("acpttest-%s", acctest.RandString(5))
	resourceName := "flexibleengine_dns_recordset_v2.recordset_1"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDNSV2RecordSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDNSV2RecordSet_private(rName, zoneName, 3000),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDNSV2RecordSetExists(resourceName, &recordset),
					resource.TestCheckResourceAttr(resourceName, "name", zoneName),
					resource.TestCheckResourceAttr(resourceName, "description", "a private record set"),
					resource.TestCheckResourceAttr(resourceName, "type", "A"),
					resource.TestCheckResourceAttr(resourceName, "ttl", "3000"),
					resource.TestCheckResourceAttr(resourceName, "records.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
				),
			},
			{
				Config: testAccDNSV2RecordSet_private(rName, zoneName, 600),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "ttl", "600"),
					resource.TestCheckResourceAttr(resourceName, "records.#", "3"),
				),
			},
		},
	})
}

func testAccCheckDNSV2RecordSetDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	dnsClient, err := config.DnsV2Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine DNS client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "flexibleengine_dns_recordset_v2" {
			continue
		}

		zoneID, recordsetID, err := parseDNSV2RecordSetId(rs.Primary.ID)
		if err != nil {
			return err
		}

		_, err = recordsets.Get(dnsClient, zoneID, recordsetID).Extract()
		if err == nil {
			return fmt.Errorf("Record set still exists")
		}
	}

	return nil
}

func testAccCheckDNSV2RecordSetExists(n string, recordset *recordsets.RecordSet) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		dnsClient, err := config.DnsV2Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating FlexibleEngine DNS client: %s", err)
		}

		zoneID, recordsetID, err := parseDNSV2RecordSetId(rs.Primary.ID)
		if err != nil {
			return err
		}

		found, err := recordsets.Get(dnsClient, zoneID, recordsetID).Extract()
		if err != nil {
			return err
		}

		if found.ID != recordsetID {
			return fmt.Errorf("Record set not found")
		}

		*recordset = *found

		return nil
	}
}

func testAccDNSV2RecordSet_base(zoneName string) string {
	return fmt.Sprintf(`
resource "flexibleengine_dns_zone_v2" "zone_1" {
  name        = "%s"
  email       = "email@example.com"
  description = "a zone for acc test"
  ttl         = 6000
}
`, zoneName)
}

func testAccDNSV2RecordSet_basic(zoneName string, ttl int) string {
	return fmt.Sprintf(`
%s

resource "flexibleengine_dns_recordset_v2" "recordset_1" {
  zone_id     = flexibleengine_dns_zone_v2.zone_1.id
  name        = "%s"
  type        = "A"
  description = "a record set"
  ttl         = %d
  records     = ["10.1.0.0", "10.1.0.1"]
}
`, testAccDNSV2RecordSet_base(zoneName), zoneName, ttl)
}

func testAccDNSV2RecordSet_tags(zoneName string, ttl int) string {
	return fmt.Sprintf(`
%s

resource "flexibleengine_dns_recordset_v2" "recordset_1" {
  zone_id     = flexibleengine_dns_zone_v2.zone_1.id
  name        = "%s"
  type        = "A"
  description = "a record set"
  ttl         = %d
  records     = ["10.1.0.0", "10.1.0.1"]

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, testAccDNSV2RecordSet_base(zoneName), zoneName, ttl)
}

func testAccDNSV2RecordSet_update(zoneName string, ttl int) string {
	return fmt.Sprintf(`
%s

resource "flexibleengine_dns_recordset_v2" "recordset_1" {
  zone_id     = flexibleengine_dns_zone_v2.zone_1.id
  name        = "%s"
  type        = "A"
  description = "an updated record set"
  ttl         = %d
  records     = ["10.1.0.2", "10.1.0.1"]

  tags = {
    foo = "bar"
    key = "value_updated"
  }
}
`, testAccDNSV2RecordSet_base(zoneName), zoneName, ttl)
}

func testAccDNSV2RecordSet_readTTL(zoneName string) string {
	return fmt.Sprintf(`
%s

resource "flexibleengine_dns_recordset_v2" "recordset_1" {
  zone_id = flexibleengine_dns_zone_v2.zone_1.id
  name    = "%s"
  type    = "A"
  records = [ "10.1.0.1", "10.1.0.2"]
}
`, testAccDNSV2RecordSet_base(zoneName), zoneName)
}

func testAccDNSV2RecordSet_private(rName, zoneName string, ttl int) string {
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
}

resource "flexibleengine_dns_recordset_v2" "recordset_1" {
  zone_id     = flexibleengine_dns_zone_v2.zone_1.id
  name        = "%s"
  type        = "A"
  description = "a private record set"
  ttl         = %d
  records     = ["10.1.0.3", "10.1.0.2", "10.1.0.1"]

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, rName, zoneName, zoneName, ttl)
}
