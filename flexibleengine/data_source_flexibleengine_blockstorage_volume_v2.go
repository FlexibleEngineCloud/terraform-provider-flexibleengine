package flexibleengine

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/openstack/blockstorage/v2/volumes"
)

func dataSourceBlockStorageVolumeV2() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceBlockStorageVolumeV2Read,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func dataSourceBlockStorageVolumeV2Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	blockStorageClient, err := config.BlockStorageV2Client(GetRegion(d, config))

	listOpts := volumes.ListOpts{
		Name:   d.Get("name").(string),
		Status: d.Get("status").(string),
	}

	pages, err := volumes.List(blockStorageClient, listOpts).AllPages()
	allVolumes, err := volumes.ExtractVolumes(pages)
	if err != nil {
		return fmt.Errorf("Unable to retrieve volumes: %s", err)
	}

	if len(allVolumes) < 1 {
		return fmt.Errorf("Your query returned no results. " +
			"Please change your search criteria and try again.")
	}

	if len(allVolumes) > 1 {
		return fmt.Errorf("Your query returned more than one result." +
			" Please try a more specific search criteria")
	}

	volume := allVolumes[0]

	log.Printf("[DEBUG] Retrieved Volume %s: %+v", volume.ID, volume)
	d.SetId(volume.ID)

	d.Set("name", volume.Name)
	d.Set("status", volume.Status)
	d.Set("region", GetRegion(d, config))

	return nil
}
