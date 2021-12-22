package flexibleengine

import (
	"context"
	"log"
	"strconv"

	"github.com/chnsz/golangsdk/openstack/dms/v1/maintainwindows"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func DataSourceDmsMaintainWindow() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDmsMaintainWindowRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
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
			"default": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func dataSourceDmsMaintainWindowRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	region := GetRegion(d, config)

	dmsV1Client, err := config.DmsV1Client(region)
	if err != nil {
		return diag.Errorf("Error creating FlexibleEngine DMS client V1: %s", err)
	}

	maintainWindows, err := maintainwindows.Get(dmsV1Client).Extract()
	if err != nil {
		return diag.FromErr(err)
	}
	if len(maintainWindows.MaintainWindows) < 1 {
		return diag.Errorf("Your query returned no results. Please change your filters and try again.")
	}
	filter := make(map[string]interface{})
	if v, ok := d.GetOk("seq"); ok {
		filter["ID"] = v.(int)
	}
	if v, ok := d.GetOk("begin"); ok {
		filter["Begin"] = v.(string)
	}
	if v, ok := d.GetOk("end"); ok {
		filter["End"] = v.(string)
	}
	if v, ok := d.GetOk("default"); ok {
		filter["Default"] = v.(bool)
	}

	data, err := utils.FilterSliceWithZeroField(maintainWindows.MaintainWindows, filter)
	if err != nil {
		return diag.Errorf("Error filtering DMS maintain window data, %s", err)
	}
	if len(data) < 1 {
		return diag.Errorf("Your query returned no results. Please change your filters and try again.")
	}

	mw := data[0].(maintainwindows.MaintainWindow)
	log.Printf("[DEBUG] DMS MaintainWindow : %#v", mw)

	d.SetId(strconv.Itoa(mw.ID))
	mErr := multierror.Append(
		d.Set("seq", mw.ID),
		d.Set("begin", mw.Begin),
		d.Set("end", mw.End),
		d.Set("default", mw.Default),
	)
	if mErr.ErrorOrNil() != nil {
		return diag.Errorf("error setting DMS maintain window attributes: %s", mErr)
	}

	return nil
}
