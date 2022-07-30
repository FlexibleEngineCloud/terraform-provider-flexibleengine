package acceptance

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/cce/v1/namespaces"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func getNamespaceResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	c, err := conf.CceV1Client(OS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating Flexibleengine CCE v1 client: %s", err)
	}
	resp, err := namespaces.Get(c, state.Primary.Attributes["cluster_id"],
		state.Primary.Attributes["name"]).Extract()
	if resp == nil && err == nil {
		return resp, fmt.Errorf("Unable to find the namespace (%s)", state.Primary.ID)
	}
	return resp, err
}

func TestAccCCENamespaceV1_basic(t *testing.T) {
	var namespace namespaces.Namespace
	resourceName := "flexibleengine_cce_namespace.test"
	randName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))

	rc := acceptance.InitResourceCheck(
		resourceName,
		&namespace,
		getNamespaceResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccCCENamespaceV1_basic(randName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(resourceName, "cluster_id",
						"flexibleengine_cce_cluster_v3.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "name", randName),
					resource.TestCheckResourceAttr(resourceName, "status", "Active"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccCCENamespaceImportStateIdFunc(randName),
			},
		},
	})
}

func TestAccCCENamespaceV1_generateName(t *testing.T) {
	var namespace namespaces.Namespace
	resourceName := "flexibleengine_cce_namespace.test"
	randName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))

	rc := acceptance.InitResourceCheck(
		resourceName,
		&namespace,
		getNamespaceResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccCCENamespaceV1_generateName(randName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(resourceName, "cluster_id",
						"flexibleengine_cce_cluster_v3.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "prefix", randName),
					resource.TestCheckResourceAttr(resourceName, "status", "Active"),
					resource.TestMatchResourceAttr(resourceName, "name", regexp.MustCompile(fmt.Sprintf(`^%s[a-z0-9-]*`, randName))),
				),
			},
		},
	})
}

func testAccCCENamespaceImportStateIdFunc(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		var clusterID string
		for _, rs := range s.RootModule().Resources {
			if rs.Type == "flexibleengine_cce_cluster_v3" {
				clusterID = rs.Primary.ID
			}
		}
		if clusterID == "" || name == "" {
			return "", fmt.Errorf("resource not found: %s/%s", clusterID, name)
		}
		return fmt.Sprintf("%s/%s", clusterID, name), nil

	}
}

func testAccCCENamespaceV1_basic(rName string) string {
	return fmt.Sprintf(`
%s

resource "flexibleengine_cce_namespace" "test" {
  cluster_id = flexibleengine_cce_cluster_v3.test.id
  name       = "%s"
}
`, testAccCceCluster_config(rName), rName)
}

func testAccCCENamespaceV1_generateName(rName string) string {
	return fmt.Sprintf(`
%s

resource "flexibleengine_cce_namespace" "test" {
  cluster_id = flexibleengine_cce_cluster_v3.test.id
  prefix     = "%s"
}
`, testAccCceCluster_config(rName), rName)
}
