package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDatasourceSmnMessageTemplates_basic(t *testing.T) {
	rName := "data.flexibleengine_smn_message_templates.test"
	name := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(rName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceSmnMessageTemplates_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(rName, "templates.0.id"),
					resource.TestCheckResourceAttrSet(rName, "templates.0.name"),
					resource.TestCheckResourceAttrSet(rName, "templates.0.protocol"),
					resource.TestCheckResourceAttrSet(rName, "templates.0.tag_names.#"),
					resource.TestCheckResourceAttrSet(rName, "templates.0.created_at"),
					resource.TestCheckResourceAttrSet(rName, "templates.0.updated_at"),
					resource.TestCheckOutput("name_filter_is_useful", "true"),
					resource.TestCheckOutput("protocol_filter_is_useful", "true"),
					resource.TestCheckOutput("template_id_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testAccDatasourceSmnMessageTemplates_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "flexibleengine_smn_message_templates" "test" {
  depends_on = [flexibleengine_smn_message_template.test]
}

data "flexibleengine_smn_message_templates" "name_filter" {
  name = "%[2]s"

  depends_on = [flexibleengine_smn_message_template.test]
}

output "name_filter_is_useful" {
  value = length(data.flexibleengine_smn_message_templates.name_filter.templates) > 0 && alltrue(
    [for v in data.flexibleengine_smn_message_templates.name_filter.templates[*].name : v == "%[2]s"]
  )  
}

data "flexibleengine_smn_message_templates" "protocol_filter" {
  protocol = "default"

  depends_on = [flexibleengine_smn_message_template.test]
}
output "protocol_filter_is_useful" {
  value = length(data.flexibleengine_smn_message_templates.protocol_filter.templates) > 0 && alltrue(
    [for v in data.flexibleengine_smn_message_templates.protocol_filter.templates[*].protocol : v == "default"]
  )  
}

locals {
  template_id = flexibleengine_smn_message_template.test.id
}

data "flexibleengine_smn_message_templates" "template_id_filter" {
  template_id = local.template_id
}

output "template_id_filter_is_useful" {
  value = length(data.flexibleengine_smn_message_templates.template_id_filter.templates) > 0 && alltrue(
    [for v in data.flexibleengine_smn_message_templates.template_id_filter.templates[*].id : v == local.template_id]
  )  
}
`, testSmnMessageTemplate_basic(name), name)
}
