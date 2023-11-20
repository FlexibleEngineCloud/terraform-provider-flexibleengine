package acceptance

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getSecurityPoliciesV3ResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := OS_REGION_NAME
	// getSecurityPolicy: Query the ELB security policy
	var (
		getSecurityPolicyHttpUrl = "v3/{project_id}/elb/security-policies/{security_policy_id}"
		getSecurityPolicyProduct = "elb"
	)
	getSecurityPolicyClient, err := cfg.NewServiceClient(getSecurityPolicyProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating SecurityPolicies Client: %s", err)
	}

	getSecurityPolicyPath := getSecurityPolicyClient.Endpoint + getSecurityPolicyHttpUrl
	getSecurityPolicyPath = strings.ReplaceAll(getSecurityPolicyPath, "{project_id}", getSecurityPolicyClient.ProjectID)
	getSecurityPolicyPath = strings.ReplaceAll(getSecurityPolicyPath, "{security_policy_id}", fmt.Sprintf("%v", state.Primary.ID))

	getSecurityPolicyOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	getSecurityPolicyResp, err := getSecurityPolicyClient.Request("GET", getSecurityPolicyPath, &getSecurityPolicyOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving SecurityPolicies: %s", err)
	}
	return utils.FlattenResponse(getSecurityPolicyResp)
}

func TestAccSecurityPoliciesV3_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "flexibleengine_lb_security_policy_v3.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getSecurityPoliciesV3ResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testSecurityPoliciesV3_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "protocols.0", "TLSv1"),
					resource.TestCheckResourceAttr(rName, "protocols.1", "TLSv1.1"),
					resource.TestCheckResourceAttr(rName, "ciphers.0", "ECDHE-RSA-AES256-GCM-SHA384"),
					resource.TestCheckResourceAttr(rName, "ciphers.1", "ECDHE-ECDSA-AES128-SHA"),
				),
			},
			{
				Config: testSecurityPoliciesV3_basic_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "protocols.0", "TLSv1.2"),
					resource.TestCheckResourceAttr(rName, "ciphers.0", "ECDHE-ECDSA-AES128-SHA"),
					resource.TestCheckResourceAttr(rName, "name", name),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testSecurityPoliciesV3_basic(name string) string {
	return fmt.Sprintf(`
resource "flexibleengine_lb_security_policy_v3" "test" {
  protocols = [
    "TLSv1",
    "TLSv1.1",
    "TLSv1.2",
    "TLSv1.3",
  ]
  ciphers = [
    "ECDHE-RSA-AES256-GCM-SHA384",
    "ECDHE-ECDSA-AES128-SHA",
    "TLS_AES_128_CCM_8_SHA256",
    "ECDHE-RSA-AES128-GCM-SHA256",
  ]
  name = "%s"
}
`, name)
}

func testSecurityPoliciesV3_basic_update(name string) string {
	return fmt.Sprintf(`
resource "flexibleengine_lb_security_policy_v3" "test" {
  protocols = [
    "TLSv1.2",
  ]
  ciphers = [
    "ECDHE-ECDSA-AES128-SHA"
  ]
  name = "%s"
}
`, name)
}
