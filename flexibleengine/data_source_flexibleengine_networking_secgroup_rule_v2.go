package flexibleengine

import (
	"fmt"
	"log"
	"strings"

	"github.com/chnsz/golangsdk/openstack/networking/v2/extensions/security/rules"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceNetworkingSecGroupRuleV2() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNetworkingSecGroupRuleV2Read,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"security_group_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"direction": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ethertype": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"port_range_min": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"port_range_max": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"protocol": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"remote_group_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"remote_ip_prefix": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateCIDR,
				StateFunc: func(v interface{}) string {
					return strings.ToLower(v.(string))
				},
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tenant_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
		},
	}
}

func dataSourceNetworkingSecGroupRuleV2Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	networkingClient, err := config.NetworkingV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine networking v2 client: %s", err)
	}

	listOpts := rules.ListOpts{
		ID:             d.Get("id").(string),
		SecGroupID:     d.Get("security_group_id").(string),
		RemoteIPPrefix: d.Get("remote_ip_prefix").(string),
		RemoteGroupID:  d.Get("remote_group_id").(string),
		Protocol:       d.Get("protocol").(string),
		Direction:      d.Get("direction").(string),
		EtherType:      d.Get("ethertype").(string),
		PortRangeMax:   d.Get("port_range_max").(int),
		PortRangeMin:   d.Get("port_range_min").(int),
		TenantID:       d.Get("tenant_id").(string),
	}

	pages, err := rules.List(networkingClient, listOpts).AllPages()
	allSecGroupRules, err := rules.ExtractRules(pages)
	if err != nil {
		return fmt.Errorf("Unable to retrieve Security Group Rule: %s", err)
	}

	if len(allSecGroupRules) < 1 {
		return fmt.Errorf("No Security Group Rule found")
	}

	if len(allSecGroupRules) > 1 {
		return fmt.Errorf("More than one Security Group Rule found")
	}

	secGroupRule := allSecGroupRules[0]

	log.Printf("[DEBUG] Retrieved Security Group Rule %s: %+v", secGroupRule.ID, secGroupRule)
	d.SetId(secGroupRule.ID)
	d.Set("description", secGroupRule.Description)
	d.Set("tenant_id", secGroupRule.TenantID)
	d.Set("region", GetRegion(d, config))
	d.Set("direction", secGroupRule.Direction)
	d.Set("ethertype", secGroupRule.EtherType)
	d.Set("port_range_max", secGroupRule.PortRangeMax)
	d.Set("port_range_min", secGroupRule.PortRangeMin)
	d.Set("protocol", secGroupRule.Protocol)
	d.Set("remote_group_id", secGroupRule.RemoteGroupID)
	d.Set("remote_ip_prefix", secGroupRule.RemoteIPPrefix)
	d.Set("security_group_id", secGroupRule.SecGroupID)

	return nil
}
