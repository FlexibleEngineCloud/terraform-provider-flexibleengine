package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/cbr"
)

func TestAccDataCBRVaults_BasicServer(t *testing.T) {
	randName := acceptance.RandomAccResourceNameWithDash()
	dataSourceName := "data.flexibleengine_cbr_vaults.test"

	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: TestAccProviderFactories,
		CheckDestroy:      dc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDataCBRVaults_serverBasic(randName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "name", randName),
					resource.TestCheckResourceAttr(dataSourceName, "vaults.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "vaults.0.consistent_level", "crash_consistent"),
					resource.TestCheckResourceAttr(dataSourceName, "vaults.0.type", cbr.VaultTypeServer),
					resource.TestCheckResourceAttr(dataSourceName, "vaults.0.protection_type", "backup"),
					resource.TestCheckResourceAttr(dataSourceName, "vaults.0.size", "200"),
					resource.TestCheckResourceAttr(dataSourceName, "vaults.0.resources.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "vaults.0.tags.foo", "bar"),
					resource.TestCheckResourceAttr(dataSourceName, "vaults.0.tags.key", "value"),
				),
			},
		},
	})
}

func TestAccDataCBRVaults_ReplicaServer(t *testing.T) {
	randName := acceptance.RandomAccResourceName()
	dataSourceName := "data.flexibleengine_cbr_vaults.test"

	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: TestAccProviderFactories,
		CheckDestroy:      dc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDataCBRVaults_serverReplication(randName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "name", randName),
					resource.TestCheckResourceAttr(dataSourceName, "vaults.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "vaults.0.consistent_level", "crash_consistent"),
					resource.TestCheckResourceAttr(dataSourceName, "vaults.0.type", cbr.VaultTypeServer),
					resource.TestCheckResourceAttr(dataSourceName, "vaults.0.protection_type", "replication"),
					resource.TestCheckResourceAttr(dataSourceName, "vaults.0.size", "200"),
				),
			},
		},
	})
}

func TestAccDataCBRVaults_BasicVolume(t *testing.T) {
	randName := acceptance.RandomAccResourceName()
	dataSourceName := "data.flexibleengine_cbr_vaults.test"

	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: TestAccProviderFactories,
		CheckDestroy:      dc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDataCBRVaults_volumeBasic(randName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "name", randName),
					resource.TestCheckResourceAttr(dataSourceName, "vaults.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "vaults.0.type", cbr.VaultTypeDisk),
					resource.TestCheckResourceAttr(dataSourceName, "vaults.0.protection_type", "backup"),
					resource.TestCheckResourceAttr(dataSourceName, "vaults.0.size", "50"),
					resource.TestCheckResourceAttr(dataSourceName, "vaults.0.resources.#", "1"),
				),
			},
		},
	})
}

func TestAccDataCBRVaults_BasicTurbo(t *testing.T) {
	randName := acceptance.RandomAccResourceName()
	dataSourceName := "data.flexibleengine_cbr_vaults.test"

	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: TestAccProviderFactories,
		CheckDestroy:      dc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDataCBRVaults_turboBasic(randName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "name", randName),
					resource.TestCheckResourceAttr(dataSourceName, "vaults.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "vaults.0.type", cbr.VaultTypeTurbo),
					resource.TestCheckResourceAttr(dataSourceName, "vaults.0.protection_type", "backup"),
					resource.TestCheckResourceAttr(dataSourceName, "vaults.0.size", "800"),
					resource.TestCheckResourceAttr(dataSourceName, "vaults.0.resources.#", "1"),
				),
			},
		},
	})
}

func TestAccDataCBRVaults_ReplicaTurbo(t *testing.T) {
	randName := acceptance.RandomAccResourceName()
	dataSourceName := "data.flexibleengine_cbr_vaults.test"

	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: TestAccProviderFactories,
		CheckDestroy:      dc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDataCBRVaults_turboReplication(randName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "name", randName),
					resource.TestCheckResourceAttr(dataSourceName, "vaults.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "vaults.0.type", cbr.VaultTypeTurbo),
					resource.TestCheckResourceAttr(dataSourceName, "vaults.0.protection_type", "replication"),
					resource.TestCheckResourceAttr(dataSourceName, "vaults.0.size", "1000"),
				),
			},
		},
	})
}

func testAccDataCBRVaults_serverBasic(rName string) string {
	return fmt.Sprintf(`
%s

data "flexibleengine_cbr_vaults" "test" {
  name = flexibleengine_cbr_vault.test.name
}
`, testAccCBRV3Vault_serverBasic(rName))
}

func testAccDataCBRVaults_serverReplication(rName string) string {
	return fmt.Sprintf(`
%s

data "flexibleengine_cbr_vaults" "test" {
  name = flexibleengine_cbr_vault.test.name
}
`, testAccCBRV3Vault_serverReplication(rName))
}

func testAccDataCBRVaults_volumeBasic(rName string) string {
	return fmt.Sprintf(`
%s

data "flexibleengine_cbr_vaults" "test" {
  name = flexibleengine_cbr_vault.test.name
}
`, testAccCBRV3Vault_volumeBasic(rName))
}

func testAccDataCBRVaults_turboBasic(rName string) string {
	return fmt.Sprintf(`
%s

data "flexibleengine_cbr_vaults" "test" {
  name = flexibleengine_cbr_vault.test.name
}
`, testAccCBRV3Vault_turboBasic(rName))
}

func testAccDataCBRVaults_turboReplication(rName string) string {
	return fmt.Sprintf(`
%s

data "flexibleengine_cbr_vaults" "test" {
  name = flexibleengine_cbr_vault.test.name
}
`, testAccCBRV3Vault_turboReplication(rName))
}
