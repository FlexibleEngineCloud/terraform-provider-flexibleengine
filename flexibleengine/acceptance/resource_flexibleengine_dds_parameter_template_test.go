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

func getDdsParameterTemplateResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := OS_REGION_NAME
	// getParameterTemplate: Query DDS parameter template
	var (
		getParameterTemplateHttpUrl = "v3/{project_id}/configurations/{config_id}"
		getParameterTemplateProduct = "dds"
	)
	getParameterTemplateClient, err := cfg.NewServiceClient(getParameterTemplateProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating DDS Client: %s", err)
	}

	getParameterTemplatePath := getParameterTemplateClient.Endpoint + getParameterTemplateHttpUrl
	getParameterTemplatePath = strings.ReplaceAll(getParameterTemplatePath, "{project_id}",
		getParameterTemplateClient.ProjectID)
	getParameterTemplatePath = strings.ReplaceAll(getParameterTemplatePath, "{config_id}", state.Primary.ID)

	getParameterTemplateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}
	getParameterTemplateResp, err := getParameterTemplateClient.Request("GET",
		getParameterTemplatePath, &getParameterTemplateOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving DDS parameter template: %s", err)
	}
	return utils.FlattenResponse(getParameterTemplateResp)
}

func TestAccDdsParameterTemplate_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	updateName := acceptance.RandomAccResourceName()
	rName := "flexibleengine_dds_parameter_template.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getDdsParameterTemplateResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testDdsParameterTemplate_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "description", "test description"),
					resource.TestCheckResourceAttr(rName, "node_version", "4.0"),
					resource.TestCheckResourceAttr(rName, "parameters.0.name",
						"connPoolMaxConnsPerHost"),
					resource.TestCheckResourceAttr(rName, "parameters.0.value", "800"),
					resource.TestCheckResourceAttr(rName, "parameters.1.name",
						"connPoolMaxShardedConnsPerHost"),
					resource.TestCheckResourceAttr(rName, "parameters.1.value", "800"),
				),
			},
			{
				Config: testDdsParameterTemplate_basic_update(updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", updateName),
					resource.TestCheckResourceAttr(rName, "description", "test description update"),
					resource.TestCheckResourceAttr(rName, "node_type", "mongos"),
					resource.TestCheckResourceAttr(rName, "node_version", "4.0"),
					resource.TestCheckResourceAttr(rName, "parameters.0.name",
						"connPoolMaxConnsPerHost"),
					resource.TestCheckResourceAttr(rName, "parameters.0.value", "500"),
					resource.TestCheckResourceAttr(rName, "parameters.1.name",
						"connPoolMaxShardedConnsPerHost"),
					resource.TestCheckResourceAttr(rName, "parameters.1.value", "500"),
				),
			},
			{
				ResourceName:            rName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"node_type", "parameter_values"},
			},
		},
	})
}

func TestAccDdsParameterTemplate_shared_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "flexibleengine_dds_parameter_template.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getDdsParameterTemplateResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testDdsParameterTemplate_shared_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "description", "test description shared node_type"),
					resource.TestCheckResourceAttr(rName, "node_type", "shard"),
					resource.TestCheckResourceAttr(rName, "node_version", "4.2"),
				),
			},
		},
	})
}

func TestAccDdsParameterTemplate_config_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "flexibleengine_dds_parameter_template.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getDdsParameterTemplateResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testDdsParameterTemplate_config_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "description", "test description config node_type"),
					resource.TestCheckResourceAttr(rName, "node_type", "config"),
					resource.TestCheckResourceAttr(rName, "node_version", "3.4"),
				),
			},
		},
	})
}

func TestAccDdsParameterTemplate_replica_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "flexibleengine_dds_parameter_template.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getDdsParameterTemplateResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testDdsParameterTemplate_replica_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "description", "test description replica node_type"),
					resource.TestCheckResourceAttr(rName, "node_type", "replica"),
					resource.TestCheckResourceAttr(rName, "node_version", "4.0"),
				),
			},
		},
	})
}

func TestAccDdsParameterTemplate_single_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "flexibleengine_dds_parameter_template.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getDdsParameterTemplateResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testDdsParameterTemplate_single_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "description", "test description single node_type"),
					resource.TestCheckResourceAttr(rName, "node_type", "single"),
					resource.TestCheckResourceAttr(rName, "node_version", "4.0"),
				),
			},
		},
	})
}

func testDdsParameterTemplate_basic(name string) string {
	return fmt.Sprintf(`
resource "flexibleengine_dds_parameter_template" "test" {
  name         = "%s"
  description  = "test description"
  node_type    = "mongos"
  node_version = "4.0"

  parameter_values = {
    connPoolMaxConnsPerHost        = 800
    connPoolMaxShardedConnsPerHost = 800
  }
}
`, name)
}

func testDdsParameterTemplate_basic_update(name string) string {
	return fmt.Sprintf(`
resource "flexibleengine_dds_parameter_template" "test" {
  name         = "%s"
  description  = "test description update"
  node_type    = "mongos"
  node_version = "4.0"

  parameter_values = {
    connPoolMaxConnsPerHost        = 500
    connPoolMaxShardedConnsPerHost = 500
  }
}
`, name)
}

func testDdsParameterTemplate_shared_basic(name string) string {
	return fmt.Sprintf(`
resource "flexibleengine_dds_parameter_template" "test" {
  name         = "%s"
  description  = "test description shared node_type"
  node_type    = "shard"
  node_version = "4.2"
}
`, name)
}

func testDdsParameterTemplate_config_basic(name string) string {
	return fmt.Sprintf(`
resource "flexibleengine_dds_parameter_template" "test" {
  name         = "%s"
  description  = "test description config node_type"
  node_type    = "config"
  node_version = "3.4"
}
`, name)
}

func testDdsParameterTemplate_replica_basic(name string) string {
	return fmt.Sprintf(`
resource "flexibleengine_dds_parameter_template" "test" {
  name         = "%s"
  description  = "test description replica node_type"
  node_type    = "replica"
  node_version = "4.0"
}
`, name)
}

func testDdsParameterTemplate_single_basic(name string) string {
	return fmt.Sprintf(`
resource "flexibleengine_dds_parameter_template" "test" {
  name         = "%s"
  description  = "test description single node_type"
  node_type    = "single"
  node_version = "4.0"
}
`, name)
}
