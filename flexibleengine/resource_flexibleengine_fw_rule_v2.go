package flexibleengine

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/openstack/networking/v2/extensions/fwaas_v2/policies"
	"github.com/huaweicloud/golangsdk/openstack/networking/v2/extensions/fwaas_v2/rules"
)

func resourceFWRuleV2() *schema.Resource {
	return &schema.Resource{
		Create: resourceFWRuleV2Create,
		Read:   resourceFWRuleV2Read,
		Update: resourceFWRuleV2Update,
		Delete: resourceFWRuleV2Delete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"region": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"protocol": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"action": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"ip_version": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default:  4,
			},
			"source_ip_address": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"destination_ip_address": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"source_port": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"destination_port": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"enabled": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"tenant_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"value_specs": &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceFWRuleV2Create(d *schema.ResourceData, meta interface{}) error {

	config := meta.(*Config)
	networkingClient, err := config.hwNetworkV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine networking client: %s", err)
	}

	enabled := d.Get("enabled").(bool)
	ipVersion := resourceFWRuleV2DetermineIPVersion(d.Get("ip_version").(int))
	protocol := resourceFWRuleV2DetermineProtocol(d.Get("protocol").(string))

	ruleConfiguration := RuleCreateOpts{
		rules.CreateOpts{
			Name:                 d.Get("name").(string),
			Description:          d.Get("description").(string),
			Protocol:             protocol,
			Action:               d.Get("action").(string),
			IPVersion:            ipVersion,
			SourceIPAddress:      d.Get("source_ip_address").(string),
			DestinationIPAddress: d.Get("destination_ip_address").(string),
			SourcePort:           d.Get("source_port").(string),
			DestinationPort:      d.Get("destination_port").(string),
			Enabled:              &enabled,
			TenantID:             d.Get("tenant_id").(string),
		},
		MapValueSpecs(d),
	}

	log.Printf("[DEBUG] Create firewall rule: %#v", ruleConfiguration)

	rule, err := rules.Create(networkingClient, ruleConfiguration).Extract()

	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Firewall rule with id %s : %#v", rule.ID, rule)

	d.SetId(rule.ID)

	return resourceFWRuleV2Read(d, meta)
}

func resourceFWRuleV2Read(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Retrieve information about firewall rule: %s", d.Id())

	config := meta.(*Config)
	networkingClient, err := config.hwNetworkV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine networking client: %s", err)
	}

	rule, err := rules.Get(networkingClient, d.Id()).Extract()
	if err != nil {
		return CheckDeleted(d, err, "FW rule")
	}

	log.Printf("[DEBUG] Read FlexibleEngine Firewall Rule %s: %#v", d.Id(), rule)

	d.Set("action", rule.Action)
	d.Set("name", rule.Name)
	d.Set("description", rule.Description)
	d.Set("ip_version", rule.IPVersion)
	d.Set("source_ip_address", rule.SourceIPAddress)
	d.Set("destination_ip_address", rule.DestinationIPAddress)
	d.Set("source_port", rule.SourcePort)
	d.Set("destination_port", rule.DestinationPort)
	d.Set("enabled", rule.Enabled)

	if rule.Protocol == "" {
		d.Set("protocol", "any")
	} else {
		d.Set("protocol", rule.Protocol)
	}

	d.Set("region", GetRegion(d, config))

	return nil
}

func resourceFWRuleV2Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	networkingClient, err := config.hwNetworkV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine networking client: %s", err)
	}

	var updateOpts rules.UpdateOpts
	if d.HasChange("name") {
		name := d.Get("name").(string)
		updateOpts.Name = &name
	}
	if d.HasChange("description") {
		description := d.Get("description").(string)
		updateOpts.Description = &description
	}
	if d.HasChange("protocol") {
		protocol := d.Get("protocol").(string)
		updateOpts.Protocol = &protocol
	}
	if d.HasChange("action") {
		action := d.Get("action").(string)
		updateOpts.Action = &action
	}
	if d.HasChange("ip_version") {
		ipVersion := resourceFWRuleV2DetermineIPVersion(d.Get("ip_version").(int))
		updateOpts.IPVersion = &ipVersion
	}
	if d.HasChange("source_ip_address") {
		sourceIPAddress := d.Get("source_ip_address").(string)
		updateOpts.SourceIPAddress = &sourceIPAddress
	}
	if d.HasChange("source_port") {
		sourcePort := d.Get("source_port").(string)
		updateOpts.SourcePort = &sourcePort
	}
	if d.HasChange("destination_ip_address") {
		destinationIPAddress := d.Get("destination_ip_address").(string)
		updateOpts.DestinationIPAddress = &destinationIPAddress
	}
	if d.HasChange("destination_port") {
		destinationPort := d.Get("destination_port").(string)
		updateOpts.DestinationPort = &destinationPort
	}
	if d.HasChange("enabled") {
		enabled := d.Get("enabled").(bool)
		updateOpts.Enabled = &enabled
	}

	log.Printf("[DEBUG] Updating firewall rules: %#v", updateOpts)
	err = rules.Update(networkingClient, d.Id(), updateOpts).Err
	if err != nil {
		return err
	}

	return resourceFWRuleV2Read(d, meta)
}

func resourceFWRuleV2Delete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Destroy firewall rule: %s", d.Id())

	config := meta.(*Config)
	networkingClient, err := config.hwNetworkV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine networking client: %s", err)
	}

	rule, err := rules.Get(networkingClient, d.Id()).Extract()
	if err != nil {
		return err
	}

	if rule.PolicyID != "" {
		_, err := policies.RemoveRule(networkingClient, rule.PolicyID, rule.ID).Extract()
		if err != nil {
			return err
		}
	}

	return rules.Delete(networkingClient, d.Id()).Err
}

func resourceFWRuleV2DetermineIPVersion(ipv int) golangsdk.IPVersion {
	// Determine the IP Version
	var ipVersion golangsdk.IPVersion
	switch ipv {
	case 4:
		ipVersion = golangsdk.IPv4
	case 6:
		ipVersion = golangsdk.IPv6
	}

	return ipVersion
}

func resourceFWRuleV2DetermineProtocol(p string) rules.Protocol {
	var protocol rules.Protocol
	switch p {
	case "any":
		protocol = rules.ProtocolAny
	case "icmp":
		protocol = rules.ProtocolICMP
	case "tcp":
		protocol = rules.ProtocolTCP
	case "udp":
		protocol = rules.ProtocolUDP
	}

	return protocol
}
