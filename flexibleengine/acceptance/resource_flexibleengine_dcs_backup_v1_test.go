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

func getDcsBackupResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := OS_REGION_NAME
	// getBackup: Query DCS backup
	var (
		getBackupHttpUrl = "v2/{project_id}/instances/{instance_id}/backups"
		getBackupProduct = "dcs"
	)
	getBackupClient, err := cfg.NewServiceClient(getBackupProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating DCS Client: %s", err)
	}

	parts := strings.SplitN(state.Primary.ID, "/", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid id format, must be <instance_id>/<backup_id>")
	}
	instanceID := parts[0]
	backupId := parts[1]
	getBackupPath := getBackupClient.Endpoint + getBackupHttpUrl
	getBackupPath = strings.ReplaceAll(getBackupPath, "{project_id}", getBackupClient.ProjectID)
	getBackupPath = strings.ReplaceAll(getBackupPath, "{instance_id}", instanceID)

	getDdmSchemasOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}

	var currentTotal int
	getBackupPath += buildGetDcsBackupQueryParams(currentTotal)

	for {
		getBackupResp, err := getBackupClient.Request("GET", getBackupPath, &getDdmSchemasOpt)
		if err != nil {
			return nil, fmt.Errorf("error retrieving DcsBackup")
		}
		getBackupRespBody, err := utils.FlattenResponse(getBackupResp)
		if err != nil {
			return nil, err
		}
		backups := utils.PathSearch("backup_record_response", getBackupRespBody, make([]interface{}, 0)).([]interface{})
		total := utils.PathSearch("total_num", getBackupRespBody, 0)
		for _, backup := range backups {
			id := utils.PathSearch("backup_id", backup, "")
			if id != backupId {
				continue
			}
			status := utils.PathSearch("status", backup, "")
			if status == "deleted" {
				return nil, fmt.Errorf("error get DCS backup by backup_id (%s)", backupId)
			}
			return backup, nil
		}
		currentTotal += len(backups)
		if currentTotal == total {
			break
		}
		getBackupPath = updatePathOffset(getBackupPath, currentTotal)
	}
	return nil, fmt.Errorf("error get DCS backup by backup_id (%s)", state.Primary.ID)
}

func buildGetDcsBackupQueryParams(offset int) string {
	return fmt.Sprintf("?limit=10&offset=%v", offset)
}

func updatePathOffset(path string, offset int) string {
	index := strings.Index(path, "offset")
	return fmt.Sprintf("%soffset=%v", path[:index], offset)
}

func TestAccDcsBackup_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "flexibleengine_dcs_backup_v1.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getDcsBackupResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testDcsBackup_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "instance_id",
						"flexibleengine_dcs_instance_v1.instance_1", "id"),
					resource.TestCheckResourceAttr(rName, "type", "manual"),
					resource.TestCheckResourceAttr(rName, "status", "succeed"),
					resource.TestCheckResourceAttr(rName, "description", "test DCS backup remark"),
					resource.TestCheckResourceAttr(rName, "backup_format", "rdb"),
					resource.TestCheckResourceAttrSet(rName, "name"),
					resource.TestCheckResourceAttrSet(rName, "size"),
					resource.TestCheckResourceAttrSet(rName, "type"),
					resource.TestCheckResourceAttrSet(rName, "begin_time"),
					resource.TestCheckResourceAttrSet(rName, "end_time"),
					resource.TestCheckResourceAttrSet(rName, "status"),
					resource.TestCheckResourceAttrSet(rName, "is_support_restore"),
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

func testDcsBackup_basic(name string) string {
	return fmt.Sprintf(`
%s

resource "flexibleengine_dcs_backup_v1" "test" {
  instance_id   = flexibleengine_dcs_instance_v1.instance_1.id
  description   = "test DCS backup remark"
  backup_format = "rdb"
}
`, testAccDcsV1Instance_basic(name))
}
