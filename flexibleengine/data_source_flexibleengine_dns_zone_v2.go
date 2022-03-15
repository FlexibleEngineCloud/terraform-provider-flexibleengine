package flexibleengine

import (
	"fmt"
	"log"

	"github.com/chnsz/golangsdk/openstack/dns/v2/zones"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceDNSZoneV2() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceDNSZoneV2Read,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"pool_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"project_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"email": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"status": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"zone_type": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"ttl": {
				Type:     schema.TypeInt,
				Optional: true,
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"masters": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"serial": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func dataSourceDNSZoneV2Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	dnsClient, err := config.DnsV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine DNS client: %s", err)
	}

	listOpts := zones.ListOpts{}

	if v, ok := d.GetOk("name"); ok {
		listOpts.Name = v.(string)
	}

	if v, ok := d.GetOk("description"); ok {
		listOpts.Description = v.(string)
	}

	if v, ok := d.GetOk("email"); ok {
		listOpts.Email = v.(string)
	}

	if v, ok := d.GetOk("status"); ok {
		listOpts.Status = v.(string)
	}

	if v, ok := d.GetOk("ttl"); ok {
		listOpts.TTL = v.(int)
	}

	if v, ok := d.GetOk("zone_type"); ok {
		listOpts.Type = v.(string)
	}

	pages, err := zones.List(dnsClient, listOpts).AllPages()
	if err != nil {
		return fmt.Errorf("Unable to retrieve zones: %s", err)
	}

	allZones, err := zones.ExtractZones(pages)
	if err != nil {
		return fmt.Errorf("Unable to extract zones: %s", err)
	}

	if len(allZones) < 1 {
		return fmt.Errorf("Your query returned no results. " +
			"Please change your search criteria and try again.")
	}

	if len(allZones) > 1 {
		return fmt.Errorf("Your query returned more than one result." +
			" Please try a more specific search criteria")
	}

	zone := allZones[0]

	log.Printf("[DEBUG] Retrieved DNS Zone %s: %+v", zone.ID, zone)
	d.SetId(zone.ID)

	// strings
	d.Set("name", zone.Name)
	d.Set("pool_id", zone.PoolID)
	d.Set("project_id", zone.ProjectID)
	d.Set("email", zone.Email)
	d.Set("description", zone.Description)
	d.Set("status", zone.Status)
	d.Set("zone_type", zone.ZoneType)
	d.Set("region", GetRegion(d, config))

	// ints
	d.Set("ttl", zone.TTL)
	d.Set("serial", zone.Serial)

	// slices
	err = d.Set("masters", zone.Masters)
	if err != nil {
		log.Printf("[DEBUG] Unable to set masters: %s", err)
		return err
	}

	return nil
}
