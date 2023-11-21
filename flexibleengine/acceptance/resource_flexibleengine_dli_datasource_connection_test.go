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

func getDatasourceConnectionResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := OS_REGION_NAME
	// getDatasourceConnection: Query the DLI instance
	var (
		getDatasourceConnectionHttpUrl = "v2.0/{project_id}/datasource/enhanced-connections/{id}"
		getDatasourceConnectionProduct = "dli"
	)
	getDatasourceConnectionClient, err := cfg.NewServiceClient(getDatasourceConnectionProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating DLI Client: %s", err)
	}

	getDatasourceConnectionPath := getDatasourceConnectionClient.Endpoint + getDatasourceConnectionHttpUrl
	getDatasourceConnectionPath = strings.ReplaceAll(getDatasourceConnectionPath, "{project_id}",
		getDatasourceConnectionClient.ProjectID)
	getDatasourceConnectionPath = strings.ReplaceAll(getDatasourceConnectionPath, "{id}", state.Primary.ID)

	getDatasourceConnectionOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	getDatasourceConnectionResp, err := getDatasourceConnectionClient.Request("GET", getDatasourceConnectionPath,
		&getDatasourceConnectionOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving DatasourceConnection: %s", err)
	}

	getDatasourceConnectionRespBody, err := utils.FlattenResponse(getDatasourceConnectionResp)
	if err != nil {
		return nil, err
	}

	if utils.PathSearch("status", getDatasourceConnectionRespBody, "") == "DELETED" {
		return nil, golangsdk.ErrDefault404{}
	}

	return getDatasourceConnectionRespBody, nil
}

func TestAccDatasourceConnection_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "flexibleengine_dli_datasource_connection.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getDatasourceConnectionResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testDatasourceConnection_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttrPair(rName, "vpc_id", "flexibleengine_vpc_v1.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "subnet_id", "flexibleengine_vpc_subnet_v1.test", "id"),
					resource.TestCheckResourceAttr(rName, "tags.foo", "bar"),
				),
			},
			{
				Config: testDatasourceConnection_basic_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "queues.0", name),
					resource.TestCheckResourceAttr(rName, "hosts.0.ip", "172.0.0.2"),
					resource.TestCheckResourceAttr(rName, "hosts.0.name", "test.test.com"),
				),
			},
			{
				ResourceName:            rName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"tags"},
			},
		},
	})
}

func testDatasourceConnectionbase(name string) string {
	return fmt.Sprintf(`
resource "flexibleengine_vpc_v1" "test" {
  name = "%s"
  cidr = "192.168.0.0/16"
}

resource "flexibleengine_vpc_subnet_v1" "test" {
  name       = "%s"
  vpc_id     = flexibleengine_vpc_v1.test.id
  cidr       = "192.168.0.0/24"
  gateway_ip = "192.168.0.1"
}
`, name, name)
}

func testDatasourceConnection_basic(name string) string {
	return fmt.Sprintf(`
%s

resource "flexibleengine_dli_datasource_connection" "test" {
  name      = "%s"
  vpc_id    = flexibleengine_vpc_v1.test.id
  subnet_id = flexibleengine_vpc_subnet_v1.test.id

  tags = {
    foo = "bar"
  }
}
`, testDatasourceConnectionbase(name), name)
}

func testDatasourceConnection_basic_update(name string) string {
	return fmt.Sprintf(`
%s

resource "flexibleengine_dli_queue" "test" {
  name          = "%s"
  cu_count      = 16
  resource_mode = 1
}

resource "flexibleengine_dli_datasource_connection" "test" {
  name      = "%s"
  vpc_id    = flexibleengine_vpc_v1.test.id
  subnet_id = flexibleengine_vpc_subnet_v1.test.id

  queues = [flexibleengine_dli_queue.test.name]

  hosts {
    ip   = "172.0.0.2"
    name = "test.test.com"
  }
}
`, testDatasourceConnectionbase(name), name, name)
}
