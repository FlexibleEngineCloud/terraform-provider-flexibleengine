package flexibleengine

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/networking/v2/extensions/security/rules"
)

func resourceNetworkingSecGroupRuleV2() *schema.Resource {
	return &schema.Resource{
		Create: resourceNetworkingSecGroupRuleV2Create,
		Read:   resourceNetworkingSecGroupRuleV2Read,
		Delete: resourceNetworkingSecGroupRuleV2Delete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"security_group_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"direction": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"ethertype": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"port_range_min": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"port_range_max": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"protocol": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"remote_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"remote_ip_prefix": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ValidateFunc: validateCIDR,
				StateFunc: func(v interface{}) string {
					return strings.ToLower(v.(string))
				},
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"tenant_id": {
				Type:       schema.TypeString,
				Optional:   true,
				ForceNew:   true,
				Computed:   true,
				Deprecated: "tenant_id is deprecated",
			},
		},
	}
}

func resourceNetworkingSecGroupRuleV2Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	networkingClient, err := config.NetworkingV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine networking client: %s", err)
	}

	portRangeMin := d.Get("port_range_min").(int)
	portRangeMax := d.Get("port_range_max").(int)
	protocol := d.Get("protocol").(string)

	if protocol == "" {
		if portRangeMin != 0 || portRangeMax != 0 {
			return fmt.Errorf("A protocol must be specified when using port_range_min and port_range_max")
		}
	}

	opts := rules.CreateOpts{
		SecGroupID:     d.Get("security_group_id").(string),
		PortRangeMin:   d.Get("port_range_min").(int),
		PortRangeMax:   d.Get("port_range_max").(int),
		RemoteGroupID:  d.Get("remote_group_id").(string),
		RemoteIPPrefix: d.Get("remote_ip_prefix").(string),
		Description:    d.Get("description").(string),
		TenantID:       d.Get("tenant_id").(string),
	}

	if v, ok := d.GetOk("direction"); ok {
		direction := resourceNetworkingSecGroupRuleV2DetermineDirection(v.(string))
		opts.Direction = direction
	}

	if v, ok := d.GetOk("ethertype"); ok {
		ethertype := resourceNetworkingSecGroupRuleV2DetermineEtherType(v.(string))
		opts.EtherType = ethertype
	}

	if v, ok := d.GetOk("protocol"); ok {
		protocol := resourceNetworkingSecGroupRuleV2DetermineProtocol(v.(string))
		opts.Protocol = protocol
	}

	log.Printf("[DEBUG] Create FlexibleEngine security group: %#v", opts)

	sgRule, err := rules.Create(networkingClient, opts).Extract()
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] FlexibleEngine Security Group Rule created: %#v", sgRule)
	d.SetId(sgRule.ID)

	return resourceNetworkingSecGroupRuleV2Read(d, meta)
}

func resourceNetworkingSecGroupRuleV2Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	region := GetRegion(d, config)
	networkingClient, err := config.NetworkingV2Client(region)
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine networking client: %s", err)
	}

	sgRule, err := rules.Get(networkingClient, d.Id()).Extract()
	if err != nil {
		return CheckDeleted(d, err, "FlexibleEngine Security Group Rule")
	}

	d.Set("region", region)
	d.Set("direction", sgRule.Direction)
	d.Set("ethertype", sgRule.EtherType)
	d.Set("protocol", sgRule.Protocol)
	d.Set("port_range_min", sgRule.PortRangeMin)
	d.Set("port_range_max", sgRule.PortRangeMax)
	d.Set("remote_group_id", sgRule.RemoteGroupID)
	d.Set("remote_ip_prefix", sgRule.RemoteIPPrefix)
	d.Set("security_group_id", sgRule.SecGroupID)
	d.Set("description", sgRule.Description)

	return nil
}

func resourceNetworkingSecGroupRuleV2Delete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Destroy security group rule: %s", d.Id())

	config := meta.(*Config)
	networkingClient, err := config.NetworkingV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine networking client: %s", err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"ACTIVE"},
		Target:     []string{"DELETED"},
		Refresh:    waitForSecGroupRuleDelete(networkingClient, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Error deleting FlexibleEngine Security Group Rule: %s", err)
	}

	d.SetId("")
	return err
}

func resourceNetworkingSecGroupRuleV2DetermineDirection(v string) rules.RuleDirection {
	var direction rules.RuleDirection
	switch v {
	case "ingress":
		direction = rules.DirIngress
	case "egress":
		direction = rules.DirEgress
	}

	return direction
}

func resourceNetworkingSecGroupRuleV2DetermineEtherType(v string) rules.RuleEtherType {
	var etherType rules.RuleEtherType
	switch v {
	case "IPv4":
		etherType = rules.EtherType4
	case "IPv6":
		etherType = rules.EtherType6
	}

	return etherType
}

func resourceNetworkingSecGroupRuleV2DetermineProtocol(v string) rules.RuleProtocol {
	var protocol rules.RuleProtocol

	// Check and see if the requested protocol matched a list of known protocol names.
	switch v {
	case "tcp":
		protocol = rules.ProtocolTCP
	case "udp":
		protocol = rules.ProtocolUDP
	case "icmp":
		protocol = rules.ProtocolICMP
	case "ah":
		protocol = rules.ProtocolAH
	case "dccp":
		protocol = rules.ProtocolDCCP
	case "egp":
		protocol = rules.ProtocolEGP
	case "esp":
		protocol = rules.ProtocolESP
	case "gre":
		protocol = rules.ProtocolGRE
	case "igmp":
		protocol = rules.ProtocolIGMP
	case "ipv6-encap":
		protocol = rules.ProtocolIPv6Encap
	case "ipv6-frag":
		protocol = rules.ProtocolIPv6Frag
	case "ipv6-icmp":
		protocol = rules.ProtocolIPv6ICMP
	case "ipv6-nonxt":
		protocol = rules.ProtocolIPv6NoNxt
	case "ipv6-opts":
		protocol = rules.ProtocolIPv6Opts
	case "ipv6-route":
		protocol = rules.ProtocolIPv6Route
	case "ospf":
		protocol = rules.ProtocolOSPF
	case "pgm":
		protocol = rules.ProtocolPGM
	case "rsvp":
		protocol = rules.ProtocolRSVP
	case "sctp":
		protocol = rules.ProtocolSCTP
	case "udplite":
		protocol = rules.ProtocolUDPLite
	case "vrrp":
		protocol = rules.ProtocolVRRP
	}

	// If the protocol wasn't matched above, see if it's an integer.
	if protocol == "" {
		_, err := strconv.Atoi(v)
		if err == nil {
			protocol = rules.RuleProtocol(v)
		}
	}

	return protocol
}

func waitForSecGroupRuleDelete(networkingClient *golangsdk.ServiceClient, secGroupRuleId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		log.Printf("[DEBUG] Attempting to delete FlexibleEngine Security Group Rule %s.\n", secGroupRuleId)

		r, err := rules.Get(networkingClient, secGroupRuleId).Extract()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				log.Printf("[DEBUG] Successfully deleted FlexibleEngine Security Group Rule %s", secGroupRuleId)
				return r, "DELETED", nil
			}
			return r, "ACTIVE", err
		}

		err = rules.Delete(networkingClient, secGroupRuleId).ExtractErr()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				log.Printf("[DEBUG] Successfully deleted FlexibleEngine Security Group Rule %s", secGroupRuleId)
				return r, "DELETED", nil
			}
			return r, "ACTIVE", err
		}

		log.Printf("[DEBUG] FlexibleEngine Security Group Rule %s still active.\n", secGroupRuleId)
		return r, "ACTIVE", nil
	}
}
