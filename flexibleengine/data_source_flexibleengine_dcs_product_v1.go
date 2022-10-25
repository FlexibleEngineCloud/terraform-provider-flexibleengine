package flexibleengine

import (
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk/openstack/dcs/v2/flavors"
	hw_utils "github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func dataSourceDcsProductV1() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceDcsProductRead,

		Schema: map[string]*schema.Schema{
			"engine": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "redis",
				ValidateFunc: validation.StringInSlice([]string{
					"redis", "memcached",
				}, true),
			},
			"engine_version": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"3.0", "4.0;5.0",
				}, false),
				Computed: true,
			},
			"cache_mode": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"capacity": {
				Type:     schema.TypeFloat,
				Optional: true,
				Computed: true,
			},
			"replica_count": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"spec_code": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"cpu_architecture": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceDcsProductRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	dcsV2Client, err := config.DcsV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error get DCS v2 client: %s", err)
	}

	var capacity string
	if v, ok := d.GetOk("capacity"); ok {
		capacity = strconv.FormatFloat(v.(float64), 'f', -1, 64)
	}

	// build a list options
	opts := flavors.ListOpts{
		CacheMode:     d.Get("cache_mode").(string),
		Engine:        d.Get("engine").(string),
		EngineVersion: d.Get("engine_version").(string),
		Capacity:      capacity,
		SpecCode:      d.Get("spec_code").(string),
	}
	log.Printf("[DEBUG] The options of list DCS flavors: %#v", opts)

	list, err := flavors.List(dcsV2Client, opts).Extract()
	if err != nil {
		return fmt.Errorf("Error while filtering data: %s", err)
	}
	if len(list) < 1 {
		return fmt.Errorf("Your query returned no results. " +
			"Please change your search criteria and try again.")
	}

	var pd flavors.Flavor
	if v, ok := d.GetOk("replica_count"); ok {
		filteredPds, err := hw_utils.FilterSliceWithField(list, map[string]interface{}{
			"ReplicaCount": v.(int),
		})
		if err != nil {
			return fmt.Errorf("Error while filtering data: %s", err)
		}

		if len(filteredPds) < 1 {
			return fmt.Errorf("Your query returned no results with replica_count=%d. "+
				"Please change your search criteria and try again.", v)
		}
		pd = filteredPds[0].(flavors.Flavor)
	} else {
		pd = list[0]
	}

	log.Printf("[DEBUG] querying DCS product: %+v", pd)
	productID := pd.SpecCode + "-h"
	d.SetId(productID)

	cap, _ := strconv.ParseFloat(pd.Capacity[0], 64)
	mErr := multierror.Append(nil,
		d.Set("capacity", cap),
		d.Set("spec_code", pd.SpecCode),
		d.Set("engine", pd.Engine),
		d.Set("engine_version", pd.EngineVersion),
		d.Set("cache_mode", pd.CacheMode),
		d.Set("replica_count", pd.ReplicaCount),
		d.Set("cpu_architecture", pd.CPUType),
	)

	if mErr.ErrorOrNil() != nil {
		return fmt.Errorf("error setting DCS product attributes: %s", mErr)
	}

	return nil
}
