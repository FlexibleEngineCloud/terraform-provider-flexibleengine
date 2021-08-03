package flexibleengine

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/golangsdk/openstack/waf/v1/domains"
)

func TestAccWafDomainV1_basic(t *testing.T) {
	var domain domains.Domain
	resourceName := "flexibleengine_waf_domain.domain_1"
	randName := acctest.RandString(10)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckWafDomainV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccWafDomainV1_basic(randName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWafDomainV1Exists(resourceName, &domain),
					resource.TestCheckResourceAttr(resourceName, "domain", fmt.Sprintf("www.%s.com", randName)),
					resource.TestCheckResourceAttr(resourceName, "proxy", "false"),
					resource.TestCheckResourceAttr(resourceName, "server.0.client_protocol", "HTTPS"),
					resource.TestCheckResourceAttr(resourceName, "server.0.server_protocol", "HTTP"),
					resource.TestCheckResourceAttr(resourceName, "server.0.port", "8080"),
					resource.TestCheckResourceAttrSet(resourceName, "policy_id"),
					resource.TestCheckResourceAttrSet(resourceName, "cname"),
				),
			},
			{
				Config: testAccWafDomainV1_update(randName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWafDomainV1Exists(resourceName, &domain),
					resource.TestCheckResourceAttr(resourceName, "server.0.port", "8443"),
					resource.TestCheckResourceAttr(resourceName, "proxy", "true"),
					resource.TestCheckResourceAttr(resourceName, "sip_header_name", "default"),
					resource.TestCheckResourceAttrSet(resourceName, "sub_domain"),
					resource.TestCheckResourceAttrSet(resourceName, "txt_code"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"keep_policy"},
			},
		},
	})
}

func TestAccWafDomainV1_policy(t *testing.T) {
	var domain domains.Domain
	resourceName := "flexibleengine_waf_domain.domain_1"
	randName := acctest.RandString(10)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckWafDomainV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccWafDomainV1_policy(randName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWafDomainV1Exists(resourceName, &domain),
					resource.TestCheckResourceAttr(resourceName, "domain", fmt.Sprintf("www.%s.com", randName)),
					resource.TestCheckResourceAttr(resourceName, "proxy", "true"),
					resource.TestCheckResourceAttr(resourceName, "sip_header_name", "default"),
					resource.TestCheckResourceAttr(resourceName, "server.0.client_protocol", "HTTPS"),
					resource.TestCheckResourceAttr(resourceName, "server.0.server_protocol", "HTTP"),
					resource.TestCheckResourceAttr(resourceName, "server.0.port", "8080"),
					resource.TestCheckResourceAttrSet(resourceName, "policy_id"),
					resource.TestCheckResourceAttrSet(resourceName, "cname"),
					resource.TestCheckResourceAttrSet(resourceName, "sub_domain"),
				),
			},
		},
	})
}

func testAccCheckWafDomainV1Destroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	wafClient, err := config.WafV1Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating Flexibleengine WAF client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "flexibleengine_waf_domain" {
			continue
		}

		_, err := domains.Get(wafClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Waf domain still exists")
		}
	}

	return nil
}

func testAccCheckWafDomainV1Exists(n string, domain *domains.Domain) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		wafClient, err := config.WafV1Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating Flexibleengine WAF client: %s", err)
		}

		found, err := domains.Get(wafClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.Id != rs.Primary.ID {
			return fmt.Errorf("Waf domain not found")
		}

		*domain = *found

		return nil
	}
}

func testAccWafDomainV1_basic(name string) string {
	return fmt.Sprintf(`
resource "flexibleengine_networking_floatingip_v2" "fip_1" {
}

%s

resource "flexibleengine_waf_domain" "domain_1" {
  domain         = "www.%s.com"
  certificate_id = flexibleengine_waf_certificate.certificate_1.id
  keep_policy    = false

  server {
    client_protocol = "HTTPS"
    server_protocol = "HTTP"
    address         = flexibleengine_networking_floatingip_v2.fip_1.address
    port            = 8080
  }
}
`, testAccWafCertificateV1_conf(name), name)
}

func testAccWafDomainV1_update(name string) string {
	return fmt.Sprintf(`
resource "flexibleengine_networking_floatingip_v2" "fip_1" {
}

%s

resource "flexibleengine_waf_domain" "domain_1" {
  domain          = "www.%s.com"
  certificate_id  = flexibleengine_waf_certificate.certificate_1.id
  keep_policy     = false
  proxy           = true
  sip_header_name = "default"
  sip_header_list = ["X-Forwarded-For"]

  server {
    client_protocol = "HTTPS"
    server_protocol = "HTTP"
    address         = flexibleengine_networking_floatingip_v2.fip_1.address
    port            = 8443
  }
}
`, testAccWafCertificateV1_conf(name), name)
}

func testAccWafDomainV1_policy(name string) string {
	return fmt.Sprintf(`
resource "flexibleengine_networking_floatingip_v2" "fip_1" {
}

%s

resource "flexibleengine_waf_policy" "policy_1" {
  name = "policy_%s"
}

resource "flexibleengine_waf_domain" "domain_1" {
  domain          = "www.%s.com"
  certificate_id  = flexibleengine_waf_certificate.certificate_1.id
  policy_id       = flexibleengine_waf_policy.policy_1.id
  proxy           = true
  sip_header_name = "default"
  sip_header_list = ["X-Forwarded-For"]

  server {
    client_protocol = "HTTPS"
    server_protocol = "HTTP"
    address         = flexibleengine_networking_floatingip_v2.fip_1.address
    port            = 8080
  }
}
`, testAccWafCertificateV1_conf(name), name, name)
}
