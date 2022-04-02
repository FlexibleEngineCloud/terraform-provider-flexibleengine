package flexibleengine

import (
	"fmt"
	"log"

	"github.com/chnsz/golangsdk/openstack/networking/v2/extensions/natgateways"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceNatGatewayV2() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNatGatewayV2Read,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"spec": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"subnet_id": {
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

func dataSourceNatGatewayV2Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	natClient, err := config.NatV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine nat client: %s", err)
	}

	listOpts := natgateways.ListOpts{
		ID:                d.Get("id").(string),
		Name:              d.Get("name").(string),
		Description:       d.Get("description").(string),
		Spec:              d.Get("spec").(string),
		RouterID:          d.Get("vpc_id").(string),
		InternalNetworkID: d.Get("subnet_id").(string),
		Status:            d.Get("status").(string),
	}

	pages, err := natgateways.List(natClient, listOpts).AllPages()
	if err != nil {
		return err
	}

	allNatGateways, err := natgateways.ExtractNatGateways(pages)

	if err != nil {
		return fmt.Errorf("Unable to retrieve natgateways: %s", err)
	}

	if len(allNatGateways) < 1 {
		return fmt.Errorf("Your query returned no results. " +
			"Please change your search criteria and try again.")
	}

	if len(allNatGateways) > 1 {
		return fmt.Errorf("Your query returned more than one result." +
			" Please try a more specific search criteria")
	}

	natgateway := allNatGateways[0]

	log.Printf("[DEBUG] Retrieved Natgateway %s: %+v", natgateway.ID, natgateway)

	d.SetId(natgateway.ID)
	d.Set("name", natgateway.Name)
	d.Set("description", natgateway.Description)
	d.Set("spec", natgateway.Spec)
	d.Set("vpc_id", natgateway.RouterID)
	d.Set("subnet_id", natgateway.InternalNetworkID)
	d.Set("status", natgateway.Status)

	return nil
}
