package acceptance

import (
	"fmt"
	"testing"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccVpcsDataSource_basic(t *testing.T) {
	randName := acceptance.RandomAccResourceName()
	randCidr := acceptance.RandomCidr()
	dataSourceName := "data.flexibleengine_vpcs.test"

	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceVpcs_basic(randName, randCidr),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "vpcs.0.cidr", randCidr),
					resource.TestCheckResourceAttr(dataSourceName, "vpcs.0.name", randName),
					resource.TestCheckResourceAttr(dataSourceName, "vpcs.0.status", "OK"),
					acceptance.TestCheckResourceAttrWithVariable(dataSourceName, "vpcs.0.id",
						"${flexibleengine_vpc_v1.test.id}"),
				),
			},
		},
	})
}

func testAccDataSourceVpcs_base(rName, cidr string) string {
	return fmt.Sprintf(`
resource "flexibleengine_vpc_v1" "test" {
  name = "%s"
  cidr = "%s"
}
`, rName, cidr)
}

func testAccDataSourceVpcs_basic(rName, cidr string) string {
	return fmt.Sprintf(`
%s

data "flexibleengine_vpcs" "test" {
  id = flexibleengine_vpc_v1.test.id
}
`, testAccDataSourceVpcs_base(rName, cidr))
}

func TestAccVpcsDataSource_byCidr(t *testing.T) {
	randName := acceptance.RandomAccResourceName()
	randCidr := acceptance.RandomCidr()
	dataSourceName := "data.flexibleengine_vpcs.test"

	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceVpcs_byCidr(randName, randCidr),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "cidr", randCidr),
					resource.TestCheckResourceAttr(dataSourceName, "vpcs.0.cidr", randCidr),
					resource.TestCheckResourceAttr(dataSourceName, "vpcs.0.name", randName),
					resource.TestCheckResourceAttr(dataSourceName, "vpcs.0.status", "OK"),
				),
			},
		},
	})
}

func testAccDataSourceVpcs_byCidr(rName, cidr string) string {
	return fmt.Sprintf(`
%s

data "flexibleengine_vpcs" "test" {
  cidr = flexibleengine_vpc_v1.test.cidr

  depends_on = [
    flexibleengine_vpc_v1.test
  ]
}
`, testAccDataSourceVpcs_base(rName, cidr))
}

func TestAccVpcsDataSource_byName(t *testing.T) {
	randName := acceptance.RandomAccResourceName()
	randCidr := acceptance.RandomCidr()
	dataSourceName := "data.flexibleengine_vpcs.test"

	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceVpcs_byName(randName, randCidr),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "name", randName),
					resource.TestCheckResourceAttr(dataSourceName, "vpcs.0.cidr", randCidr),
					resource.TestCheckResourceAttr(dataSourceName, "vpcs.0.name", randName),
					resource.TestCheckResourceAttr(dataSourceName, "vpcs.0.status", "OK"),
					acceptance.TestCheckResourceAttrWithVariable(dataSourceName, "vpcs.0.id",
						"${flexibleengine_vpc_v1.test.id}"),
				),
			},
		},
	})
}

func testAccDataSourceVpcs_byName(rName, cidr string) string {
	return fmt.Sprintf(`
%s

data "flexibleengine_vpcs" "test" {
  name = flexibleengine_vpc_v1.test.name

  depends_on = [
    flexibleengine_vpc_v1.test
  ]
}
`, testAccDataSourceVpcs_base(rName, cidr))
}

func TestAccVpcsDataSource_byAll(t *testing.T) {
	randName := acceptance.RandomAccResourceName()
	randCidr := acceptance.RandomCidr()
	dataSourceName := "data.flexibleengine_vpcs.test"

	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckEpsID(t)
		},
		ProviderFactories: TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceVpcs_byAll(randName, randCidr),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "name", randName),
					resource.TestCheckResourceAttr(dataSourceName, "vpcs.0.cidr", randCidr),
					resource.TestCheckResourceAttr(dataSourceName, "vpcs.0.name", randName),
					resource.TestCheckResourceAttr(dataSourceName, "vpcs.0.status", "OK"),
					acceptance.TestCheckResourceAttrWithVariable(dataSourceName, "vpcs.0.id",
						"${flexibleengine_vpc_v1.test.id}"),
				),
			},
		},
	})
}

func testAccDataSourceVpcs_byAll(rName, cidr string) string {
	return fmt.Sprintf(`
%s

data "flexibleengine_vpcs" "test" {
  id     = flexibleengine_vpc_v1.test.id
  name   = flexibleengine_vpc_v1.test.name
  cidr   = flexibleengine_vpc_v1.test.cidr
  status = "OK"

  depends_on = [
    flexibleengine_vpc_v1.test
  ]
}
`, testAccDataSourceVpcs_base(rName, cidr))
}

func TestAccVpcsDataSource_tags(t *testing.T) {
	randName1 := acceptance.RandomAccResourceName()
	randName2 := acceptance.RandomAccResourceName()
	dataSourceName := "data.flexibleengine_vpcs.test"

	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceVpcs_tags(randName1, randName2),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "tags.foo", randName1),
					resource.TestCheckResourceAttr(dataSourceName, "vpcs.0.name", randName1),
					resource.TestCheckResourceAttr(dataSourceName, "vpcs.0.status", "OK"),
				),
			},
		},
	})
}

func testAccDataSourceVpcs_tags(rName1, rName2 string) string {
	return fmt.Sprintf(`
resource "flexibleengine_vpc_v1" "test1" {
  name = "%s"
  cidr = "172.16.0.0/24"
  tags = {
    foo = "%s"
  }
}

resource "flexibleengine_vpc_v1" "test2" {
  name = "%s"
  cidr = "10.12.2.0/24"
  tags = {
    foo = "%s"
  }
}

data "flexibleengine_vpcs" "test" {
  tags = {
    foo = "%s"
  }
  depends_on = [
    flexibleengine_vpc_v1.test1,
    flexibleengine_vpc_v1.test2,
  ]
}
`, rName1, rName1, rName2, rName2, rName1)
}
