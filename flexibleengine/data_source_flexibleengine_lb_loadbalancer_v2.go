package flexibleengine

import (
	"fmt"
	"log"

	"github.com/chnsz/golangsdk/openstack/common/tags"
	"github.com/chnsz/golangsdk/openstack/networking/v2/extensions/lbaas_v2/loadbalancers"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceELBV2Loadbalancer() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceELBV2LoadbalancerRead,
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
			"vip_subnet_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"vip_address": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"vip_port_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceELBV2LoadbalancerRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	region := GetRegion(d, config)
	lbClient, err := config.ElbV2Client(region)
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine ELB v2.0 client: %s", err)
	}

	listOpts := loadbalancers.ListOpts{
		Name:        d.Get("name").(string),
		ID:          d.Get("id").(string),
		Description: d.Get("description").(string),
		VipSubnetID: d.Get("vip_subnet_id").(string),
		VipAddress:  d.Get("vip_address").(string),
	}
	pages, err := loadbalancers.List(lbClient, listOpts).AllPages()
	if err != nil {
		return fmt.Errorf("Unable to retrieve loadbalancers: %s", err)
	}
	lbList, err := loadbalancers.ExtractLoadBalancers(pages)
	if err != nil {
		return fmt.Errorf("Unable to extract loadbalancers: %s", err)
	}

	if len(lbList) < 1 {
		return fmt.Errorf("Your query returned no results, Please change your search criteria and try again")
	}

	if len(lbList) > 1 {
		return fmt.Errorf("Your query returned more than one result, Please try a more specific search criteria")
	}

	lb := lbList[0]
	d.SetId(lb.ID)

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("name", lb.Name),
		d.Set("description", lb.Description),
		d.Set("status", lb.OperatingStatus),
		d.Set("vip_address", lb.VipAddress),
		d.Set("vip_subnet_id", lb.VipSubnetID),
		d.Set("vip_port_id", lb.VipPortID),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return fmt.Errorf("Error setting elb loadbalancer fields: %s", err)
	}

	// Get tags
	if resourceTags, err := tags.Get(lbClient, "loadbalancers", d.Id()).Extract(); err == nil {
		tagmap := tagsToMap(resourceTags.Tags)
		d.Set("tags", tagmap)
	} else {
		log.Printf("[WARN] fetching tags of elb loadbalancer failed: %s", err)
	}

	return nil
}
