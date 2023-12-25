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

func getDdsAuditLogPolicyResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := OS_REGION_NAME
	// getAuditLog: Query DDS audit log
	var (
		getAuditLogPolicyHttpUrl = "v3/{project_id}/instances/{instance_id}/auditlog-policy"
		getAuditLogPolicyProduct = "dds"
	)
	getAuditLogPolicyClient, err := cfg.NewServiceClient(getAuditLogPolicyProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating DDS Client: %s", err)
	}

	instanceID := state.Primary.ID
	getAuditLogPolicyPath := getAuditLogPolicyClient.Endpoint + getAuditLogPolicyHttpUrl
	getAuditLogPolicyPath = strings.ReplaceAll(getAuditLogPolicyPath, "{project_id}",
		getAuditLogPolicyClient.ProjectID)
	getAuditLogPolicyPath = strings.ReplaceAll(getAuditLogPolicyPath, "{instance_id}", instanceID)

	getAuditLogPolicyOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}
	getAuditLogPolicyResp, err := getAuditLogPolicyClient.Request("GET", getAuditLogPolicyPath, &getAuditLogPolicyOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving DDS audit log policy: %s", err)
	}

	getAuditLogPolicyRespBody, err := utils.FlattenResponse(getAuditLogPolicyResp)
	if err != nil {
		return nil, err
	}

	keepDays := utils.PathSearch("keep_days", getAuditLogPolicyRespBody, 0)
	if keepDays.(float64) == 0 {
		return nil, fmt.Errorf("the instance %s has no audit log policy", instanceID)
	}

	return getAuditLogPolicyRespBody, nil
}

func TestAccDdsAuditLogPolicy_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "flexibleengine_dds_audit_log_policy.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getDdsAuditLogPolicyResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testDdsAuditLogPolicy_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "keep_days", "7"),
				),
			},
			{
				Config: testDdsAuditLogPolicy_basic_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "keep_days", "15"),
				),
			},
			{
				ResourceName:            rName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"instance_id"},
			},
		},
	})
}

func TestAccDdsAuditLogPolicy_audit_types(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "flexibleengine_dds_audit_log_policy.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getDdsAuditLogPolicyResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDdsAuditLogPolicy_audit_types(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "keep_days", "7"),
					resource.TestCheckResourceAttr(rName, "audit_types.#", "6"),
					resource.TestCheckResourceAttrSet(rName, "audit_types.#"),
				),
			},
		},
	})
}

func testAccResourceDdsAuditLogPolicy_base(rName string) string {
	return fmt.Sprintf(`
%s

data "flexibleengine_availability_zones" "test" {}

resource "flexibleengine_dds_instance_v3" "test" {
  name              = "%s"
  availability_zone = data.flexibleengine_availability_zones.test.names[0]
  vpc_id            = flexibleengine_vpc_v1.test.id
  subnet_id         = flexibleengine_vpc_subnet_v1.test.id
  security_group_id = flexibleengine_networking_secgroup_v2.test.id
  password          = "Terraform@123"
  mode              = "Sharding"

  datastore {
    type           = "DDS-Community"
    version        = "3.4"
    storage_engine = "wiredTiger"
  }

  flavor {
    type      = "mongos"
    num       = 2
    spec_code = "dds.mongodb.s3.medium.4.mongos"
  }

  flavor {
    type      = "shard"
    num       = 2
    storage   = "ULTRAHIGH"
    size      = 20
    spec_code = "dds.mongodb.s3.medium.4.shard"
  }

  flavor {
    type      = "config"
    num       = 1
    storage   = "ULTRAHIGH"
    size      = 20
    spec_code = "dds.mongodb.s3.large.2.config"
  }

  backup_strategy {
    start_time = "08:00-09:00"
    keep_days  = "8"
  }

  tags = {
    foo   = "bar"
    owner = "terraform"
  }
}`, testBaseNetwork(rName), rName)
}

func testDdsAuditLogPolicy_basic(name string) string {
	return fmt.Sprintf(`
%s

resource "flexibleengine_dds_audit_log_policy" "test" {
  instance_id       = flexibleengine_dds_instance_v3.test.id
  keep_days         = 7
  reserve_auditlogs = true
}
`, testAccResourceDdsAuditLogPolicy_base(name))
}

func testDdsAuditLogPolicy_basic_update(name string) string {
	return fmt.Sprintf(`
%s

resource "flexibleengine_dds_audit_log_policy" "test" {
  instance_id       = flexibleengine_dds_instance_v3.test.id
  keep_days         = 15
  reserve_auditlogs = false
}
`, testAccResourceDdsAuditLogPolicy_base(name))
}

func testAccDdsAuditLogPolicy_audit_types(name string) string {
	return fmt.Sprintf(`
%s

resource "flexibleengine_dds_audit_log_policy" "test" {
  instance_id = flexibleengine_dds_instance_v3.test.id
  keep_days   = 7

  audit_types = [
    "delete", "insert", "update", "query", "auth", "command"
  ]
}
`, testAccResourceDdsAuditLogPolicy_base(name))
}
