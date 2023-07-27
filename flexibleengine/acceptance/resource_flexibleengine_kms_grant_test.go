package acceptance

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/jmespath/go-jmespath"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getKmsGrantResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := OS_REGION_NAME
	// getGrant: Query the KMS manual Grant
	var (
		getGrantHttpUrl = "v1.0/{project_id}/kms/list-grants"
		getGrantProduct = "kms"
	)
	getGrantClient, err := cfg.NewServiceClient(getGrantProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating KMS Client: %s", err)
	}

	getGrantPath := getGrantClient.Endpoint + getGrantHttpUrl
	getGrantPath = strings.ReplaceAll(getGrantPath, "{project_id}", getGrantClient.ProjectID)

	getGrantOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
		JSONBody: map[string]interface{}{
			"key_id": state.Primary.Attributes["key_id"],
			"limit":  100,
		},
	}
	getGrantResp, err := getGrantClient.Request("POST", getGrantPath, &getGrantOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving KMS grant: %s", err)
	}

	grantResponseBody, err := utils.FlattenResponse(getGrantResp)
	if err != nil {
		return nil, fmt.Errorf("error retrieving KMS grant: %s", err)
	}

	searchPath := fmt.Sprintf("grants[?grant_id=='%s']|[0]", state.Primary.ID)
	r, err := jmespath.Search(searchPath, grantResponseBody)
	if err != nil || r == nil {
		return nil, fmt.Errorf("error retrieving KMS grant: %s", err)
	}

	return r, nil
}
func TestAccKmsGrant_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "flexibleengine_kms_grant.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getKmsGrantResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckAdminOnly(t)
			testAccPreCheck(t)
		},
		ProviderFactories: TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKmsGrant_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "operations.#", "2"),
					resource.TestCheckResourceAttrPair(rName, "key_id", "flexibleengine_kms_key_v1.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "grantee_principal", "flexibleengine_identity_user_v3.test", "id"),
					resource.TestCheckResourceAttrSet(rName, "creator"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testKmsGrantImportState(rName),
			},
		},
	})
}

func testKmsGrant_basic(name string) string {
	return fmt.Sprintf(`

resource "flexibleengine_kms_key_v1" "test" {
  key_alias    = "%s"
  pending_days = "7"
}

resource "flexibleengine_identity_user_v3" "test" {
  name        = "%s"
  password    = "password123@!"
  enabled     = true
  description = "tested by terraform"
}

resource "flexibleengine_kms_grant" "test" {
  key_id            = flexibleengine_kms_key_v1.test.id
  grantee_principal = flexibleengine_identity_user_v3.test.id
  operations        = ["create-datakey", "encrypt-datakey"]
}

`, name, name)
}

func testKmsGrantImportState(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("Resource (%s) not found: %s", name, rs)
		}
		if rs.Primary.Attributes["key_id"] == "" {
			return "", fmt.Errorf("Attribute (key_id) of Resource (%s) not found: %s", name, rs)
		}
		if rs.Primary.ID == "" {
			return "", fmt.Errorf("Attribute (ID) of Resource (%s) not found: %s", name, rs)
		}

		return rs.Primary.Attributes["key_id"] + "/" +
			rs.Primary.ID, nil
	}
}
