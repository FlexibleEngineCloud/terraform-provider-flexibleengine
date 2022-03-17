package flexibleengine

import (
	"fmt"
	"log"

	"github.com/chnsz/golangsdk/openstack/bms/v2/flavors"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceBMSFlavorV2() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceBMSFlavorV2Read,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"vcpus": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntAtLeast(0),
			},
			"min_ram": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"min_disk": {
				Type:     schema.TypeInt,
				Optional: true,
			},

			"ram": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"disk": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"swap": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"rx_tx_factor": {
				Type:     schema.TypeFloat,
				Computed: true,
			},
			"sort_key": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "id",
			},
			"sort_dir": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "asc",
				ValidateFunc: dataSourceImagesImageV2SortDirection,
			},
		},
	}
}

func dataSourceBMSFlavorV2Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	flavorClient, err := config.ComputeV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating Flexibleengine bms client: %s", err)
	}

	listOpts := flavors.ListOpts{
		MinDisk: d.Get("min_disk").(int),
		MinRAM:  d.Get("min_ram").(int),
		Name:    d.Get("name").(string),
		ID:      d.Get("id").(string),
		SortKey: d.Get("sort_key").(string),
		SortDir: d.Get("sort_dir").(string),
	}
	var flavor flavors.Flavor
	allFlavors, err := flavors.List(flavorClient, listOpts)
	if err != nil {
		return fmt.Errorf("Unable to retrieve flavors: %s", err)
	}

	if len(allFlavors) < 1 {
		return fmt.Errorf("Your query returned no results. " +
			"Please change your search criteria and try again.")
	}

	vcpus := d.Get("vcpus").(int)
	index := findBMSFlavorsByCPUs(allFlavors, vcpus)
	if index == -1 {
		return fmt.Errorf("Your query returned no results by %d vcpus. "+
			"Please change your search criteria and try again.", vcpus)
	}

	flavor = allFlavors[index]
	log.Printf("[DEBUG] Retrieve BMS flavor: %#v", flavor)

	d.SetId(flavor.ID)
	d.Set("name", flavor.Name)
	d.Set("vcpus", flavor.VCPUs)
	d.Set("disk", flavor.Disk)
	d.Set("min_disk", flavor.MinDisk)
	d.Set("min_ram", flavor.MinRAM)
	d.Set("ram", flavor.RAM)
	d.Set("rx_tx_factor", flavor.RxTxFactor)
	d.Set("swap", flavor.Swap)

	return nil
}

func findBMSFlavorsByCPUs(all []flavors.Flavor, vcpus int) int {
	if vcpus == 0 {
		return 0
	}

	for i, item := range all {
		if item.VCPUs == vcpus {
			return i
		}
	}

	return -1
}
