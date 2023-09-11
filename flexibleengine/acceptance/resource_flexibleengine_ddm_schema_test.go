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

func getDdmSchemaResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := OS_REGION_NAME
	// getSchema: Query DDM schema
	var (
		getSchemaHttpUrl = "v1/{project_id}/instances/{instance_id}/databases/{ddm_dbname}"
		getSchemaProduct = "ddm"
	)
	getSchemaClient, err := cfg.NewServiceClient(getSchemaProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating DDM client: %s", err)
	}

	parts := strings.SplitN(state.Primary.ID, "/", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid id format, must be <instance_id>/<db_name>")
	}
	instanceID := parts[0]
	schemaName := parts[1]
	getSchemaPath := getSchemaClient.Endpoint + getSchemaHttpUrl
	getSchemaPath = strings.ReplaceAll(getSchemaPath, "{project_id}", getSchemaClient.ProjectID)
	getSchemaPath = strings.ReplaceAll(getSchemaPath, "{instance_id}", fmt.Sprintf("%v", instanceID))
	getSchemaPath = strings.ReplaceAll(getSchemaPath, "{ddm_dbname}", fmt.Sprintf("%v", schemaName))

	getSchemaOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	getSchemaResp, err := getSchemaClient.Request("GET", getSchemaPath, &getSchemaOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving DdmSchema: %s", err)
	}
	getSchemaRespBody, err := utils.FlattenResponse(getSchemaResp)
	if err != nil {
		return nil, err
	}

	return utils.PathSearch("database", getSchemaRespBody, nil), nil
}

func TestAccDdmSchema_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	instanceName := strings.ReplaceAll(name, "_", "-")
	rName := "flexibleengine_ddm_schema.test"
	dbPwd := "Test@12345678"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getDdmSchemaResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testDdmSchema_basic(instanceName, name, dbPwd),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "shard_mode", "single"),
					resource.TestCheckResourceAttr(rName, "shard_number", "1"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{"instance_id", "data_nodes.0.admin_user",
					"data_nodes.0.admin_password", "delete_rds_data"},
			},
		},
	})
}

func testDdmSchema_base(name, dbPwd string) string {
	return fmt.Sprintf(`
%[1]s

resource "flexibleengine_networking_secgroup_v2" "test" {
  name = "%[2]s"
}

resource "flexibleengine_networking_secgroup_rule_v2" "test" {
  security_group_id = flexibleengine_networking_secgroup_v2.test.id
  direction         = "ingress"
  ethertype         = "IPv4"
  remote_group_id   = flexibleengine_networking_secgroup_v2.test.id
}

data "flexibleengine_availability_zones" "test" {}

data "flexibleengine_ddm_engines" test {
  version = "3.0.8.5"
}

data "flexibleengine_ddm_flavors" test {
 engine_id = data.flexibleengine_ddm_engines.test.engines[0].id
 cpu_arch  = "X86"
}

resource "flexibleengine_ddm_instance" "test" {
 name              = "%[2]s"
 flavor_id         = data.flexibleengine_ddm_flavors.test.flavors[0].id
 node_num          = 2
 engine_id         = data.flexibleengine_ddm_engines.test.engines[0].id
 vpc_id            = flexibleengine_vpc_v1.test.id
 subnet_id         = flexibleengine_vpc_subnet_v1.test.id
 security_group_id = flexibleengine_networking_secgroup_v2.test.id

 availability_zones = [
   data.flexibleengine_availability_zones.test.names[0]
 ]
}

data "flexibleengine_rds_flavors_v3" "test" {
  db_type       = "MySQL"
  db_version    = "5.7"
  instance_mode = "single"
  vcpus         = 2
  memory        = 4
}

resource "flexibleengine_rds_instance_v3" "test" {
  name               = "%[2]s"
  flavor             = data.flexibleengine_rds_flavors_v3.test.flavors[0].name
  security_group_id  = flexibleengine_networking_secgroup_v2.test.id
  subnet_id          = flexibleengine_vpc_subnet_v1.test.id
  vpc_id             = flexibleengine_vpc_v1.test.id

 availability_zone = [
   data.flexibleengine_availability_zones.test.names[0]
 ]

 db {
   password = "%[3]s"
   type     = "MySQL"
   version  = "5.7"
   port     = 3306
 }

 volume {
   type = "ULTRAHIGH"
   size = 40
 }
}
`, testVpc(name), name, dbPwd)
}

func testDdmSchema_basic(instanceName, name, dbPwd string) string {
	return fmt.Sprintf(`
%[1]s

resource "flexibleengine_ddm_schema" "test" {
  instance_id  = flexibleengine_ddm_instance.test.id
  name         = "%[2]s"
  shard_mode   = "single"
  shard_number = "1"

  data_nodes {
    id             = flexibleengine_rds_instance_v3.test.id
    admin_user     = "root"
    admin_password = "%[3]s"
  }

  delete_rds_data = "true"

  lifecycle {
    ignore_changes = [
      data_nodes,
    ]
  }
}
`, testDdmSchema_base(instanceName, dbPwd), name, dbPwd)
}
