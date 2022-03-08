package flexibleengine

import (
	"fmt"
	"log"
	"strconv"

	"github.com/chnsz/golangsdk/openstack/dcs/v1/maintainwindows"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceDcsMaintainWindowV1() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceDcsMaintainWindowV1Read,

		Schema: map[string]*schema.Schema{
			"default": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"seq": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"begin": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"end": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func dataSourceDcsMaintainWindowV1Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	dcsV1Client, err := config.DcsV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating dcs key client: %s", err)
	}

	seq := d.Get("seq").(int)
	begin := d.Get("begin").(string)
	end := d.Get("end").(string)

	df := d.Get("default").(bool)
	if seq == 0 && begin == "" && end == "" {
		df = true
	}

	v, err := maintainwindows.Get(dcsV1Client).Extract()
	if err != nil {
		return err
	}

	maintainWindows := v.MaintainWindows
	var filteredMVs []maintainwindows.MaintainWindow
	for _, mv := range maintainWindows {
		if seq != 0 && mv.ID != seq {
			continue
		}
		if begin != "" && mv.Begin != begin {
			continue
		}
		if end != "" && mv.End != end {
			continue
		}

		if mv.Default != df {
			continue
		}
		filteredMVs = append(filteredMVs, mv)
	}
	if len(filteredMVs) < 1 {
		return fmt.Errorf("Your query returned no results. Please change your filters and try again.")
	}
	mw := filteredMVs[0]
	d.SetId(strconv.Itoa(mw.ID))
	d.Set("begin", mw.Begin)
	d.Set("end", mw.End)
	d.Set("default", mw.Default)
	log.Printf("[DEBUG] Dcs MaintainWindow : %+v", mw)

	return nil
}
