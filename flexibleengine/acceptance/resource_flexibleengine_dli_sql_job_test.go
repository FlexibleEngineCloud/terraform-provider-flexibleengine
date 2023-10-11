package acceptance

import (
	"fmt"
	"testing"

	"github.com/chnsz/golangsdk/openstack/dli/v1/sqljob"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func getDliSqlJobResourceFunc(config *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := config.DliV1Client(OS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating Dli v1 client, err=%s", err)
	}
	return sqljob.Status(client, state.Primary.ID)
}

// check the DDL sql
func TestAccResourceDliSqlJob_basic(t *testing.T) {
	var sqlJobObj sqljob.SqlJobOpts
	resourceName := "flexibleengine_dli_sql_job.test"
	name := acceptance.RandomAccResourceName()
	obsBucketName := acceptance.RandomAccResourceNameWithDash()

	rc := acceptance.InitResourceCheck(
		resourceName,
		&sqlJobObj,
		getDliSqlJobResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: TestAccProviderFactories,
		CheckDestroy:      testAccCheckDliSqlJobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSqlJobBaseResource_basic(name, obsBucketName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "sql", fmt.Sprint("DESC ", name)),
					resource.TestCheckResourceAttr(resourceName, "database_name", name),
					resource.TestCheckResourceAttr(resourceName, "job_type", "DDL"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"rows", "schema"},
			},
		},
	})
}

func TestAccResourceDliSqlJob_query(t *testing.T) {
	var sqlJobObj sqljob.SqlJobOpts
	resourceName := "flexibleengine_dli_sql_job.test"
	name := acceptance.RandomAccResourceName()
	obsBucketName := acceptance.RandomAccResourceNameWithDash()

	rc := acceptance.InitResourceCheck(
		resourceName,
		&sqlJobObj,
		getDliSqlJobResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: TestAccProviderFactories,
		CheckDestroy:      testAccCheckDliSqlJobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSqlJobBaseResource_query(name, obsBucketName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "sql", fmt.Sprint("SELECT * FROM ", name)),
					resource.TestCheckResourceAttr(resourceName, "database_name", name),
					resource.TestCheckResourceAttr(resourceName, "queue_name", "default"),
					resource.TestCheckResourceAttr(resourceName, "job_type", "QUERY"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"rows", "schema"},
			},
		},
	})
}

func TestAccResourceDliSqlJob_async(t *testing.T) {
	var sqlJobObj sqljob.SqlJobOpts
	resourceName := "flexibleengine_dli_sql_job.test"
	name := acceptance.RandomAccResourceName()
	obsBucketName := acceptance.RandomAccResourceNameWithDash()

	rc := acceptance.InitResourceCheck(
		resourceName,
		&sqlJobObj,
		getDliSqlJobResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: TestAccProviderFactories,
		CheckDestroy:      testAccCheckDliSqlJobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSqlJobResource_aync(name, obsBucketName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "sql", fmt.Sprint("SELECT * FROM ", name)),
					resource.TestCheckResourceAttr(resourceName, "database_name", name),
					resource.TestCheckResourceAttr(resourceName, "queue_name", "default"),
					resource.TestCheckResourceAttr(resourceName, "job_type", "QUERY"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"rows", "schema", "conf", "duration", "status"},
			},
		},
	})
}

func testAccSqlJobBaseResource(name, obsBucketName string) string {
	return fmt.Sprintf(`
resource "flexibleengine_obs_bucket" "test" {
  bucket = "%s"
  acl    = "private"
}

resource "flexibleengine_obs_bucket_object" "test" {
  bucket       = flexibleengine_obs_bucket.test.bucket
  key          = "user/data/user.csv"
  content      = "Jason,Tokyo"
  content_type = "text/plain"
}

resource "flexibleengine_dli_database" "test" {
  name        = "%s"
  description = "For terraform acc test"
}

resource "flexibleengine_dli_table" "test" {
  database_name   = flexibleengine_dli_database.test.name
  name            = "%s"
  data_location   = "OBS"
  description     = "dli table test"
  data_format     = "csv"
  bucket_location = "obs://${flexibleengine_obs_bucket_object.test.bucket}/user/data"

  columns {
    name        = "name"
    type        = "string"
    description = "person name"
  }

  columns {
    name        = "addrss"
    type        = "string"
    description = "home address"
  }
}
`, obsBucketName, name, name)
}

func testAccSqlJobBaseResource_basic(name, obsBucketName string) string {
	return fmt.Sprintf(`
%s

resource "flexibleengine_dli_sql_job" "test" {
  sql           = "DESC ${flexibleengine_dli_table.test.name}"
  database_name = flexibleengine_dli_database.test.name
}
`, testAccSqlJobBaseResource(name, obsBucketName))
}

func testAccSqlJobBaseResource_query(name, obsBucketName string) string {
	return fmt.Sprintf(`
%s

resource "flexibleengine_dli_sql_job" "test" {
  sql           = "SELECT * FROM ${flexibleengine_dli_table.test.name}"
  database_name = flexibleengine_dli_database.test.name
}
`, testAccSqlJobBaseResource(name, obsBucketName))
}

func testAccSqlJobResource_aync(name, obsBucketName string) string {
	return fmt.Sprintf(`
%s

resource "flexibleengine_dli_sql_job" "test" {
  sql           = "SELECT * FROM ${flexibleengine_dli_table.test.name}"
  database_name = flexibleengine_dli_database.test.name

  conf {
    dli_sql_sqlasync_enabled = true
  }
}
`, testAccSqlJobBaseResource(name, obsBucketName))
}

func testAccCheckDliSqlJobDestroy(s *terraform.State) error {
	config := acceptance.TestAccProvider.Meta().(*config.Config)
	client, err := config.DliV1Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating Dli client, err=%s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "flexibleengine_dli_sql_job" {
			continue
		}

		res, err := sqljob.Status(client, rs.Primary.ID)
		if err == nil && res != nil && (res.Status != sqljob.JobStatusCancelled &&
			res.Status != sqljob.JobStatusFinished && res.Status != sqljob.JobStatusFailed) {
			return fmt.Errorf("flexibleengine_dli_sql_job still exists:%s,%+v,%+v", rs.Primary.ID, err, res)
		}
	}

	return nil
}
