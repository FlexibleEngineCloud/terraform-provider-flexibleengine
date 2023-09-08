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

func getDdmInstanceResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := OS_REGION_NAME
	// getInstance: Query DDM instance
	var (
		getInstanceHttpUrl = "v1/{project_id}/instances/{instance_id}"
		getInstanceProduct = "ddm"
	)
	getInstanceClient, err := cfg.NewServiceClient(getInstanceProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating DDM client: %s", err)
	}

	getInstancePath := getInstanceClient.Endpoint + getInstanceHttpUrl
	getInstancePath = strings.ReplaceAll(getInstancePath, "{project_id}", getInstanceClient.ProjectID)
	getInstancePath = strings.ReplaceAll(getInstancePath, "{instance_id}", fmt.Sprintf("%v", state.Primary.ID))

	getInstanceOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	getInstanceResp, err := getInstanceClient.Request("GET", getInstancePath, &getInstanceOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving DdmInstance: %s", err)
	}

	getInstanceRespBody, err := utils.FlattenResponse(getInstanceResp)
	if err != nil {
		return nil, err
	}

	status := utils.PathSearch("status", getInstanceRespBody, nil)
	if status == "DELETED" {
		return nil, fmt.Errorf("error get DDM instance")
	}
	return getInstanceRespBody, nil
}

func TestAccDdmInstance_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceNameWithDash()
	updateName := acceptance.RandomAccResourceNameWithDash()
	rName := "flexibleengine_ddm_instance.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getDdmInstanceResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testDdmInstance_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "node_num", "2"),
					resource.TestCheckResourceAttr(rName, "admin_user", "test_user_1"),
					resource.TestCheckResourceAttrPair(rName, "flavor_id",
						"data.flexibleengine_ddm_flavors.test", "flavors.0.id"),
					resource.TestCheckResourceAttrPair(rName, "engine_id",
						"data.flexibleengine_ddm_engines.test", "engines.0.id"),
					resource.TestCheckResourceAttrPair(rName, "vpc_id",
						"flexibleengine_vpc_v1.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "subnet_id",
						"flexibleengine_vpc_subnet_v1.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "security_group_id",
						"flexibleengine_networking_secgroup_v2.test", "id"),
				),
			},
			{
				Config: testDdmInstance_basic_update(name, updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", updateName),
					resource.TestCheckResourceAttr(rName, "node_num", "4"),
					resource.TestCheckResourceAttr(rName, "admin_user", "test_user_1"),
					resource.TestCheckResourceAttrPair(rName, "flavor_id",
						"data.flexibleengine_ddm_flavors.test", "flavors.0.id"),
					resource.TestCheckResourceAttrPair(rName, "engine_id",
						"data.flexibleengine_ddm_engines.test", "engines.0.id"),
					resource.TestCheckResourceAttrPair(rName, "vpc_id",
						"flexibleengine_vpc_v1.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "subnet_id",
						"flexibleengine_vpc_subnet_v1.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "security_group_id",
						"flexibleengine_networking_secgroup_v2.test_update", "id"),
				),
			},
			{
				Config: testDdmInstance_basic_update_reduce(name, updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", updateName),
					resource.TestCheckResourceAttr(rName, "node_num", "2"),
					resource.TestCheckResourceAttr(rName, "admin_user", "test_user_1"),
					resource.TestCheckResourceAttrPair(rName, "flavor_id",
						"data.flexibleengine_ddm_flavors.test", "flavors.0.id"),
					resource.TestCheckResourceAttrPair(rName, "engine_id",
						"data.flexibleengine_ddm_engines.test", "engines.0.id"),
					resource.TestCheckResourceAttrPair(rName, "vpc_id",
						"flexibleengine_vpc_v1.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "subnet_id",
						"flexibleengine_vpc_subnet_v1.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "security_group_id",
						"flexibleengine_networking_secgroup_v2.test_update", "id"),
				),
			},
			{
				ResourceName:            rName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"admin_password", "engine_id", "flavor_id"},
			},
		},
	})
}

func testDdmInstance_base(name string) string {
	return fmt.Sprintf(`
%s

data "flexibleengine_availability_zones" "test" {}

data "flexibleengine_ddm_engines" test {
  version = "3.0.8.5"
}

data "flexibleengine_ddm_flavors" test {
  engine_id = data.flexibleengine_ddm_engines.test.engines[0].id
  cpu_arch  = "X86"
}
`, testBaseNetwork(name))
}

func testDdmInstance_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "flexibleengine_ddm_instance" "test" {
  name              = "%[2]s"
  flavor_id         = data.flexibleengine_ddm_flavors.test.flavors[0].id
  node_num          = 2
  engine_id         = data.flexibleengine_ddm_engines.test.engines[0].id
  vpc_id            = flexibleengine_vpc_v1.test.id
  subnet_id         = flexibleengine_vpc_subnet_v1.test.id
  security_group_id = flexibleengine_networking_secgroup_v2.test.id
  admin_user        = "test_user_1"
  admin_password    = "test_password_123"

  availability_zones = [
    data.flexibleengine_availability_zones.test.names[0]
  ]
}
`, testDdmInstance_base(name), name)
}

func testDdmInstance_basic_update(name, updateName string) string {
	return fmt.Sprintf(`
%[1]s

resource "flexibleengine_networking_secgroup_v2" "test_update" {
  name = "%[2]s"
}

resource "flexibleengine_ddm_instance" "test" {
  name              = "%[2]s"
  flavor_id         = data.flexibleengine_ddm_flavors.test.flavors[0].id
  node_num          = 4
  engine_id         = data.flexibleengine_ddm_engines.test.engines[0].id
  vpc_id            = flexibleengine_vpc_v1.test.id
  subnet_id         = flexibleengine_vpc_subnet_v1.test.id
  security_group_id = flexibleengine_networking_secgroup_v2.test_update.id
  admin_user        = "test_user_1"
  admin_password    = "test_password_123"

  availability_zones = [
    data.flexibleengine_availability_zones.test.names[0]
  ]
}
`, testDdmInstance_base(name), updateName)
}

func testDdmInstance_basic_update_reduce(name, updateName string) string {
	return fmt.Sprintf(`
%s

resource "flexibleengine_networking_secgroup_v2" "test_update" {
  name = "%[2]s"
}

resource "flexibleengine_ddm_instance" "test" {
  name              = "%[2]s"
  flavor_id         = data.flexibleengine_ddm_flavors.test.flavors[0].id
  node_num          = 2
  engine_id         = data.flexibleengine_ddm_engines.test.engines[0].id
  vpc_id            = flexibleengine_vpc_v1.test.id
  subnet_id         = flexibleengine_vpc_subnet_v1.test.id
  security_group_id = flexibleengine_networking_secgroup_v2.test_update.id

  availability_zones = [
    data.flexibleengine_availability_zones.test.names[0]
  ]
}
`, testDdmInstance_base(name), updateName)
}
