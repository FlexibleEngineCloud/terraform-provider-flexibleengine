package flexibleengine

import (
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/huaweicloud/golangsdk/openstack/dcs/v1/maintainwindows"
)

func dataSourceDcsMaintainWindowV1() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceDcsMaintainWindowV1Read,

		Schema: map[string]*schema.Schema{
			"seq": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"begin": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"end": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"default": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func dataSourceDcsMaintainWindowV1Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	dcsV1Client, err := config.dcsV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating dcs key client: %s", err)
	}

	v, err := maintainwindows.Get(dcsV1Client).Extract()
	if err != nil {
		return err
	}

	maintainWindows := v.MaintainWindows
	var filteredMVs []maintainwindows.MaintainWindow
	for _, mv := range maintainWindows {
		seq := d.Get("seq").(int)
		if seq != 0 && mv.ID != seq {
			continue
		}

		begin := d.Get("begin").(string)
		if begin != "" && mv.Begin != begin {
			continue
		}
		end := d.Get("end").(string)
		if end != "" && mv.End != end {
			continue
		}

		df := d.Get("default").(bool)
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
