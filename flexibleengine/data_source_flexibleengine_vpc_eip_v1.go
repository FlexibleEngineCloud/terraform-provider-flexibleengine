package flexibleengine

import (
	"context"

	"github.com/chnsz/golangsdk/openstack/networking/v1/eips"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceVpcEipV1() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceVpcEipRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"public_ip": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"port_id": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"private_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"bandwidth_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"bandwidth_size": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"bandwidth_share_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceVpcEipRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	networkingClient, err := config.NetworkingV1Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("Error creating networking client: %s", err)
	}

	var listOpts eips.ListOpts
	if portId, ok := d.GetOk("port_id"); ok {
		listOpts.PortId = []string{portId.(string)}
	}

	if publicIp, ok := d.GetOk("public_ip"); ok {
		listOpts.PublicIp = []string{publicIp.(string)}
	}

	pages, err := eips.List(networkingClient, listOpts).AllPages()
	if err != nil {
		return diag.FromErr(err)
	}

	allEips, err := eips.ExtractPublicIPs(pages)
	if err != nil {
		return diag.Errorf("Unable to retrieve eips: %s ", err)
	}

	if len(allEips) < 1 {
		return diag.Errorf("Your query returned no results. " +
			"Please change your search criteria and try again.")
	}

	if len(allEips) > 1 {
		return diag.Errorf("Your query returned more than one result." +
			" Please try a more specific search criteria")
	}

	Eip := allEips[0]

	d.SetId(Eip.ID)
	d.Set("region", config.GetRegion(d))
	d.Set("public_ip", Eip.PublicAddress)
	d.Set("port_id", Eip.PortID)
	d.Set("status", normalizeEIPStatus(Eip.Status))
	d.Set("type", Eip.Type)
	d.Set("private_ip", Eip.PrivateAddress)
	d.Set("bandwidth_id", Eip.BandwidthID)
	d.Set("bandwidth_size", Eip.BandwidthSize)
	d.Set("bandwidth_share_type", Eip.BandwidthShareType)

	return nil
}
