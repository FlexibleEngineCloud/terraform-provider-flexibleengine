package flexibleengine

import (
	"context"
	"log"
	"strconv"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk/openstack/dms/v2/products"
)

func dataSourceDmsProduct() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDmsProductRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"bandwidth": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice(
					[]string{"100MB", "300MB", "600MB", "1200MB"},
					false,
				),
			},
			"engine": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "kafka",
				ValidateFunc: validation.StringInSlice([]string{"kafka"}, false),
			},
			"engine_version": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "2.3.0",
			},

			"availability_zones": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"spec_code": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cpu_arch": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ecs_flavor_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"partition_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"storage_space": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"storage_spec_codes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"max_tps": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func dataSourceDmsProductRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	dmsV1Client, err := config.DmsV1Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error get FlexibleEngine DMS product client V1: %s", err)
	}

	bandwidth := d.Get("bandwidth").(string)
	engine := d.Get("engine").(string)
	version := d.Get("engine_version").(string)

	v, err := products.Get(dmsV1Client, engine)
	if err != nil {
		return diag.FromErr(err)
	}

	var filteredProduct *products.Detail
	for _, pd := range v.Hourly {
		if pd.Version != version {
			continue
		}

		for _, value := range pd.Values {
			for _, detail := range value.Details {
				if detail.Bandwidth == bandwidth {
					filteredProduct = &detail
					break
				}
			}
			if filteredProduct != nil {
				break
			}
		}
		if filteredProduct != nil {
			break
		}
	}

	if filteredProduct == nil {
		return diag.Errorf("Your query returned no results. Please change your filters and try again.")
	}

	log.Printf("[DEBUG] DMS product detail : %#v", filteredProduct)
	d.SetId(filteredProduct.ProductID)

	maxTps, _ := strconv.Atoi(filteredProduct.Tps)
	partitionNum, _ := strconv.Atoi(filteredProduct.PartitionNum)
	volumeSize, _ := strconv.Atoi(filteredProduct.Storage)
	storageSpecCodes := make([]string, 0, len(filteredProduct.IOs))
	for _, v := range filteredProduct.IOs {
		storageSpecCodes = append(storageSpecCodes, v.StorageSpecCode)
	}

	var mErr *multierror.Error
	mErr = multierror.Append(err,
		d.Set("bandwidth", filteredProduct.Bandwidth),
		d.Set("ecs_flavor_id", filteredProduct.EcsFlavorId),
		d.Set("availability_zones", filteredProduct.AvailableZones),
		d.Set("cpu_arch", filteredProduct.ArchType),
		d.Set("spec_code", filteredProduct.SpecCode),
		d.Set("partition_num", partitionNum),
		d.Set("max_tps", maxTps),
		d.Set("storage_space", volumeSize),
		d.Set("storage_spec_codes", storageSpecCodes),
	)

	if mErr.ErrorOrNil() != nil {
		return diag.Errorf("Error setting DMS product attributes: %s", mErr)
	}

	return nil
}
