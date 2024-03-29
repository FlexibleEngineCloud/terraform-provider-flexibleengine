package flexibleengine

import (
	"fmt"

	"github.com/chnsz/golangsdk/openstack/elb/v3/flavors"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceElbFlavorsV3() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceElbFlavorsV3Read,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"max_connections": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"cps": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"qps": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"bandwidth": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			// Computed values.
			"flavors": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"max_connections": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"cps": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"qps": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"bandwidth": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceElbFlavorsV3Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	region := GetRegion(d, config)
	elbClient, err := config.ElbV3Client(region)
	if err != nil {
		return fmt.Errorf("Error creating ELB v3 client: %s", err)
	}

	listOpts := flavors.ListOpts{}
	if v, ok := d.GetOk("type"); ok {
		listOpts.Type = []string{v.(string)}
	}

	pages, err := flavors.List(elbClient, listOpts).AllPages()
	if err != nil {
		return err
	}

	allFlavors, err := flavors.ExtractFlavors(pages)
	if err != nil {
		return fmt.Errorf("Unable to retrieve flavors: %s", err)
	}

	maxConnections := d.Get("max_connections").(int)
	cps := d.Get("cps").(int)
	qps := d.Get("qps").(int)
	bandwidth := d.Get("bandwidth").(int)

	var ids []string
	var s []map[string]interface{}
	for _, flavor := range allFlavors {
		if flavor.SoldOut {
			continue
		}

		if maxConnections > 0 && flavor.Info.Connection != maxConnections {
			continue
		}

		if cps > 0 && flavor.Info.Cps != cps {
			continue
		}

		if qps > 0 && flavor.Info.Qps != qps {
			continue
		}

		if bandwidth > 0 && flavor.Info.Bandwidth != bandwidth*1000 {
			continue
		}

		ids = append(ids, flavor.ID)
		mapping := map[string]interface{}{
			"id":              flavor.ID,
			"name":            flavor.Name,
			"type":            flavor.Type,
			"max_connections": flavor.Info.Connection,
			"cps":             flavor.Info.Cps,
			"qps":             flavor.Info.Qps,
			"bandwidth":       int(flavor.Info.Bandwidth / 1000),
		}
		s = append(s, mapping)
	}

	if len(ids) < 1 {
		return fmt.Errorf("Your query returned no results. " +
			"Please change your search criteria and try again.")
	}

	d.SetId(HashStrings(ids))

	d.Set("region", region)
	d.Set("ids", ids)
	if err := d.Set("flavors", s); err != nil {
		return err
	}

	return nil
}
