package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/cbr/v3/vaults"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/cbr"
)

func TestAccCBRV3Vault_BasicServer(t *testing.T) {
	var vault vaults.Vault
	randName := acceptance.RandomAccResourceName()
	resourceName := "flexibleengine_cbr_vault.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: TestAccProviderFactories,
		CheckDestroy:      testAccCheckCBRVaultDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCBRV3Vault_serverBasic(randName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCBRVaultExists(resourceName, &vault),
					resource.TestCheckResourceAttr(resourceName, "name", randName),
					resource.TestCheckResourceAttr(resourceName, "consistent_level", "crash_consistent"),
					resource.TestCheckResourceAttr(resourceName, "type", cbr.VaultTypeServer),
					resource.TestCheckResourceAttr(resourceName, "protection_type", "backup"),
					resource.TestCheckResourceAttr(resourceName, "size", "200"),
					resource.TestCheckResourceAttr(resourceName, "resources.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
					resource.TestCheckResourceAttrPair(resourceName, "resources.0.server_id", "flexibleengine_compute_instance_v2.test", "id"),
				),
			},
			{
				Config: testAccCBRV3Vault_serverUpdate(randName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", randName+"-update"),
					resource.TestCheckResourceAttr(resourceName, "consistent_level", "crash_consistent"),
					resource.TestCheckResourceAttr(resourceName, "type", cbr.VaultTypeServer),
					resource.TestCheckResourceAttr(resourceName, "protection_type", "backup"),
					resource.TestCheckResourceAttr(resourceName, "size", "300"),
					resource.TestCheckResourceAttr(resourceName, "resources.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo1", "bar"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value_update"),
					resource.TestCheckResourceAttrPair(resourceName, "policy_id", "flexibleengine_cbr_policy.test", "id"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccCBRV3Vault_ReplicaServer(t *testing.T) {
	var vault vaults.Vault
	randName := acceptance.RandomAccResourceName()
	resourceName := "flexibleengine_cbr_vault.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: TestAccProviderFactories,
		CheckDestroy:      testAccCheckCBRVaultDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCBRV3Vault_serverReplication(randName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCBRVaultExists(resourceName, &vault),
					resource.TestCheckResourceAttr(resourceName, "name", randName),
					resource.TestCheckResourceAttr(resourceName, "consistent_level", "crash_consistent"),
					resource.TestCheckResourceAttr(resourceName, "type", cbr.VaultTypeServer),
					resource.TestCheckResourceAttr(resourceName, "protection_type", "replication"),
					resource.TestCheckResourceAttr(resourceName, "size", "200"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccCBRV3Vault_BasicVolume(t *testing.T) {
	var vault vaults.Vault
	randName := acceptance.RandomAccResourceName()
	resourceName := "flexibleengine_cbr_vault.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: TestAccProviderFactories,
		CheckDestroy:      testAccCheckCBRVaultDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCBRV3Vault_volumeBasic(randName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCBRVaultExists(resourceName, &vault),
					resource.TestCheckResourceAttr(resourceName, "name", randName),
					resource.TestCheckResourceAttr(resourceName, "consistent_level", "crash_consistent"),
					resource.TestCheckResourceAttr(resourceName, "type", cbr.VaultTypeDisk),
					resource.TestCheckResourceAttr(resourceName, "protection_type", "backup"),
					resource.TestCheckResourceAttr(resourceName, "size", "50"),
					resource.TestCheckResourceAttr(resourceName, "auto_expand", "false"),
					resource.TestCheckResourceAttr(resourceName, "resources.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "resources.0.includes.#", "2"),
				),
			},
			{
				Config: testAccCBRV3Vault_volumeUpdate(randName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", randName+"-update"),
					resource.TestCheckResourceAttr(resourceName, "consistent_level", "crash_consistent"),
					resource.TestCheckResourceAttr(resourceName, "type", cbr.VaultTypeDisk),
					resource.TestCheckResourceAttr(resourceName, "protection_type", "backup"),
					resource.TestCheckResourceAttr(resourceName, "size", "100"),
					resource.TestCheckResourceAttr(resourceName, "auto_expand", "true"),
					resource.TestCheckResourceAttrSet(resourceName, "policy_id"),
					resource.TestCheckResourceAttr(resourceName, "resources.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "resources.0.includes.#", "2"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccCBRV3Vault_BasicTurbo(t *testing.T) {
	var vault vaults.Vault
	randName := acceptance.RandomAccResourceName()
	resourceName := "flexibleengine_cbr_vault.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: TestAccProviderFactories,
		CheckDestroy:      testAccCheckCBRVaultDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCBRV3Vault_turboBasic(randName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCBRVaultExists(resourceName, &vault),
					resource.TestCheckResourceAttr(resourceName, "name", randName),
					resource.TestCheckResourceAttr(resourceName, "consistent_level", "crash_consistent"),
					resource.TestCheckResourceAttr(resourceName, "type", cbr.VaultTypeTurbo),
					resource.TestCheckResourceAttr(resourceName, "protection_type", "backup"),
					resource.TestCheckResourceAttr(resourceName, "size", "800"),
					resource.TestCheckResourceAttr(resourceName, "resources.#", "1"),
				),
			},
			{
				Config: testAccCBRV3Vault_turboUpdate(randName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", randName+"-update"),
					resource.TestCheckResourceAttr(resourceName, "consistent_level", "crash_consistent"),
					resource.TestCheckResourceAttr(resourceName, "type", cbr.VaultTypeTurbo),
					resource.TestCheckResourceAttr(resourceName, "protection_type", "backup"),
					resource.TestCheckResourceAttr(resourceName, "size", "1000"),
					resource.TestCheckResourceAttrSet(resourceName, "policy_id"),
					resource.TestCheckResourceAttr(resourceName, "resources.#", "1"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccCBRV3Vault_ReplicaTurbo(t *testing.T) {
	var vault vaults.Vault
	randName := acceptance.RandomAccResourceName()
	resourceName := "flexibleengine_cbr_vault.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: TestAccProviderFactories,
		CheckDestroy:      testAccCheckCBRVaultDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCBRV3Vault_turboReplication(randName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCBRVaultExists(resourceName, &vault),
					resource.TestCheckResourceAttr(resourceName, "name", randName),
					resource.TestCheckResourceAttr(resourceName, "consistent_level", "crash_consistent"),
					resource.TestCheckResourceAttr(resourceName, "type", cbr.VaultTypeTurbo),
					resource.TestCheckResourceAttr(resourceName, "protection_type", "replication"),
					resource.TestCheckResourceAttr(resourceName, "size", "1000"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckCBRVaultDestroy(s *terraform.State) error {
	conf := testAccProvider.Meta().(*config.Config)
	client, err := conf.CbrV3Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine CBR client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "flexibleengine_cbr_vault" {
			continue
		}

		_, err := vaults.Get(client, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("CBR vault still exists")
		}
	}

	return nil
}

func testAccCheckCBRVaultExists(n string, vault *vaults.Vault) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		conf := testAccProvider.Meta().(*config.Config)
		client, err := conf.CbrV3Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating FlexibleEngine CBR client: %s", err)
		}

		found, err := vaults.Get(client, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("CBR vault not found")
		}

		*vault = *found

		return nil
	}
}

func testAccCBRV3Vault_policy(rName string) string {
	return fmt.Sprintf(`
resource "flexibleengine_cbr_policy" "test" {
  name        = "%s"
  type        = "backup"
  time_period = 20

  backup_cycle {
    days            = "MO,TU"
    execution_times = ["06:00", "18:00"]
  }
}
`, rName)
}

func testAccEvsVolumeConfiguration_basic() string {
	return fmt.Sprintf(`
variable "volume_configuration" {
  type = list(object({
    volume_type = string
    size        = number
  }))
  default = [
    {volume_type = "SSD", size = 100},
    {volume_type = "SSD", size = 100},
  ]
}`)
}

func testAccEvsVolumeConfiguration_update() string {
	return fmt.Sprintf(`
variable "volume_configuration" {
  type = list(object({
    volume_type = string
    size        = number
  }))
  default = [
    {volume_type = "SAS", size = 100},
    {volume_type = "SAS", size = 100},
  ]
}`)
}

func testAccCBRV3VaultBasicConfiguration(config, rName string) string {
	return fmt.Sprintf(`
%s

data "flexibleengine_availability_zones" "test" {}

data "flexibleengine_compute_flavors_v2" "test" {
  performance_type = "normal"
  cpu_core         = 2
  memory_size      = 4
}

data "flexibleengine_images_image_v2" "test" {
  name        = "OBS Ubuntu 18.04"
  most_recent = true
}

resource "flexibleengine_vpc_v1" "test" {
  name = "%s"
  cidr = "192.168.0.0/20"
}

resource "flexibleengine_vpc_subnet_v1" "test" {
  name       = "%s"
  cidr       = "192.168.0.0/24"
  vpc_id     = flexibleengine_vpc_v1.test.id
  gateway_ip = "192.168.0.1"
}

resource "flexibleengine_networking_secgroup_v2" "test" {
  name = "%s"
}

resource "flexibleengine_compute_keypair_v2" "test" {
  name = "%s"
  lifecycle {
    ignore_changes = [
      public_key,
    ]
  }
}

resource "flexibleengine_compute_instance_v2" "test" {
  availability_zone = data.flexibleengine_availability_zones.test.names[0]
  name              = "%s"
  image_id          = data.flexibleengine_images_image_v2.test.id
  flavor_id         = data.flexibleengine_compute_flavors_v2.test.flavors[0]
  key_pair          = flexibleengine_compute_keypair_v2.test.name

  security_groups = [
    flexibleengine_networking_secgroup_v2.test.name
  ]

  network {
    uuid = flexibleengine_vpc_subnet_v1.test.id
  }
}

resource "flexibleengine_blockstorage_volume_v2" "test" {
  count = length(var.volume_configuration)

  availability_zone = data.flexibleengine_availability_zones.test.names[0]
  volume_type       = var.volume_configuration[count.index].volume_type
  name              = "%s_${tostring(count.index)}"
  size              = var.volume_configuration[count.index].size
}

resource "flexibleengine_compute_volume_attach_v2" "test" {
  count = length(flexibleengine_blockstorage_volume_v2.test)

  instance_id = flexibleengine_compute_instance_v2.test.id
  volume_id   = flexibleengine_blockstorage_volume_v2.test[count.index].id
}`, config, rName, rName, rName, rName, rName, rName)
}

func testAccCBRV3Vault_serverBasic(rName string) string {
	return fmt.Sprintf(`
%s

resource "flexibleengine_cbr_vault" "test" {
  name             = "%s"
  type             = "server"
  consistent_level = "crash_consistent"
  protection_type  = "backup"
  size             = 200

  resources {
    server_id = flexibleengine_compute_instance_v2.test.id
  }

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, testAccCBRV3VaultBasicConfiguration(testAccEvsVolumeConfiguration_basic(), rName),
		rName)
}

func testAccCBRV3Vault_serverUpdate(rName string) string {
	return fmt.Sprintf(`
%s

%s

resource "flexibleengine_cbr_vault" "test" {
  name             = "%s-update"
  type             = "server"
  consistent_level = "crash_consistent"
  protection_type  = "backup"
  size             = 300
  policy_id        = flexibleengine_cbr_policy.test.id

  resources {
    server_id = flexibleengine_compute_instance_v2.test.id
  }

  tags = {
    foo1 = "bar"
    key  = "value_update"
  }
}
`, testAccCBRV3VaultBasicConfiguration(testAccEvsVolumeConfiguration_update(), rName), testAccCBRV3Vault_policy(rName),
		rName)
}

func testAccCBRV3Vault_serverReplication(rName string) string {
	return fmt.Sprintf(`
resource "flexibleengine_cbr_vault" "test" {
  name             = "%s"
  type             = "server"
  consistent_level = "crash_consistent"
  protection_type  = "replication"
  size             = 200
}
`, rName)
}

func testAccCBRV3Vault_volumeBasic(rName string) string {
	return fmt.Sprintf(`
%s

resource "flexibleengine_cbr_vault" "test" {
  name            = "%s"
  type            = "disk"
  protection_type = "backup"
  size            = 50

  resources {
    includes = flexibleengine_compute_volume_attach_v2.test[*].volume_id
  }
}
`, testAccCBRV3VaultBasicConfiguration(testAccEvsVolumeConfiguration_basic(), rName),
		rName)
}

func testAccCBRV3Vault_volumeUpdate(rName string) string {
	return fmt.Sprintf(`
%s

%s

resource "flexibleengine_cbr_vault" "test" {
  name            = "%s-update"
  type            = "disk"
  protection_type = "backup"
  size            = 100
  auto_expand     = true
  policy_id       = flexibleengine_cbr_policy.test.id

  resources {
    includes = flexibleengine_compute_volume_attach_v2.test[*].volume_id
  }
}
`, testAccCBRV3VaultBasicConfiguration(testAccEvsVolumeConfiguration_basic(), rName),
		testAccCBRV3Vault_policy(rName), rName)
}

// Vaults of type 'turbo'
func testAccCBRV3Vault_turboBase(rName string) string {
	return fmt.Sprintf(`
data "flexibleengine_availability_zones" "test" {}

resource "flexibleengine_vpc_v1" "test" {
  name = "%s"
  cidr = "192.168.0.0/20"
}

resource "flexibleengine_vpc_subnet_v1" "test" {
  name       = "%s"
  cidr       = "192.168.0.0/22"
  gateway_ip = "192.168.0.1"
  vpc_id     = flexibleengine_vpc_v1.test.id
}

resource "flexibleengine_networking_secgroup_v2" "test" {
  name = "%s"
}

resource "flexibleengine_sfs_turbo" "test1" {
  name              = "%s-1"
  size              = 500
  share_proto       = "NFS"
  vpc_id            = flexibleengine_vpc_v1.test.id
  subnet_id         = flexibleengine_vpc_subnet_v1.test.id
  security_group_id = flexibleengine_networking_secgroup_v2.test.id
  availability_zone = data.flexibleengine_availability_zones.test.names[0]
}`, rName, rName, rName, rName)
}

func testAccCBRV3Vault_turboBasic(rName string) string {
	return fmt.Sprintf(`
%s

resource "flexibleengine_cbr_vault" "test" {
  name            = "%s"
  type            = "turbo"
  protection_type = "backup"
  size            = 800

  resources {
    includes = [
      flexibleengine_sfs_turbo.test1.id
    ]
  }
}
`, testAccCBRV3Vault_turboBase(rName), rName)
}

func testAccCBRV3Vault_turboUpdate(rName string) string {
	return fmt.Sprintf(`
%s

%s

resource "flexibleengine_sfs_turbo" "test2" {
  name              = "%s-2"
  size              = 500
  share_proto       = "NFS"
  vpc_id            = flexibleengine_vpc_v1.test.id
  subnet_id         = flexibleengine_vpc_subnet_v1.test.id
  security_group_id = flexibleengine_networking_secgroup_v2.test.id
  availability_zone = data.flexibleengine_availability_zones.test.names[0]
}

resource "flexibleengine_cbr_vault" "test" {
  name            = "%s-update"
  type            = "turbo"
  protection_type = "backup"
  size            = 1000
  policy_id       = flexibleengine_cbr_policy.test.id

  resources {
    includes = [
      flexibleengine_sfs_turbo.test1.id,
      flexibleengine_sfs_turbo.test2.id
    ]
  }
}
`, testAccCBRV3Vault_turboBase(rName), testAccCBRV3Vault_policy(rName), rName, rName)
}

func testAccCBRV3Vault_turboReplication(rName string) string {
	return fmt.Sprintf(`
resource "flexibleengine_cbr_vault" "test" {
  name            = "%s"
  type            = "turbo"
  protection_type = "replication"
  size            = 1000
}
`, rName)
}
