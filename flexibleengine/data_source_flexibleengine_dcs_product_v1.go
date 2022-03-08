package flexibleengine

import (
	"fmt"
	"log"

	"github.com/chnsz/golangsdk/openstack/dcs/v1/products"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	hw_utils "github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func dataSourceDcsProductV1() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceDcsProductV1Read,

		Schema: map[string]*schema.Schema{
			"spec_code": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"engine": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "redis",
				ValidateFunc: validation.StringInSlice([]string{
					"redis", "memcached",
				}, false),
			},
			"engine_version": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"cache_mode": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func dataSourceDcsProductV1Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	dcsV1Client, err := config.DcsV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error get dcs product client: %s", err)
	}

	v, err := products.Get(dcsV1Client).Extract()
	if err != nil {
		return err
	}
	log.Printf("[DEBUG] get %d DCS products", len(v.Products))

	filteredPds, err := hw_utils.FilterSliceWithField(v.Products, map[string]interface{}{
		"SpecCode":      d.Get("spec_code").(string),
		"Engine":        d.Get("engine").(string),
		"EngineVersion": d.Get("engine_version").(string),
		"CacheMode":     d.Get("cache_mode").(string),
	})
	if err != nil {
		return fmt.Errorf("Error while filtering data : %s", err)
	}

	if len(filteredPds) < 1 {
		return fmt.Errorf("Your query returned no results. Please change your filters and try again")
	}

	pd := filteredPds[0].(products.Product)
	log.Printf("[DEBUG] DCS product : %+v", pd)

	d.SetId(pd.ProductID)
	d.Set("spec_code", pd.SpecCode)
	d.Set("engine", pd.Engine)
	d.Set("engine_version", pd.EngineVersion)
	d.Set("cache_mode", pd.CacheMode)

	return nil
}
