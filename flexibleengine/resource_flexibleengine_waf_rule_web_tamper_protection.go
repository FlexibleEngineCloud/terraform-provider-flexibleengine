package flexibleengine

import (
	"fmt"
	"log"

	rules "github.com/chnsz/golangsdk/openstack/waf/v1/webtamperprotection_rules"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceWafRuleWebTamperProtection() *schema.Resource {
	return &schema.Resource{
		Create: resourceWafRuleWebTamperProtectionCreate,
		Read:   resourceWafRuleWebTamperProtectionRead,
		Delete: resourceWafRuleWebTamperProtectionDelete,
		Importer: &schema.ResourceImporter{
			State: resourceWafRulesImport,
		},

		Schema: map[string]*schema.Schema{
			"policy_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"domain": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"path": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceWafRuleWebTamperProtectionCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	wafClient, err := config.WafV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("error creating Flexibleengine WAF Client: %s", err)
	}

	createOpts := rules.CreateOpts{
		Hostname: d.Get("domain").(string),
		Url:      d.Get("path").(string),
	}

	policyID := d.Get("policy_id").(string)
	rule, err := rules.Create(wafClient, policyID, createOpts).Extract()
	if err != nil {
		return fmt.Errorf("error creating Flexibleengine WAF Web Tamper Protection Rule: %s", err)
	}

	log.Printf("[DEBUG] WAF web tamper protection rule created: %#v", rule)
	d.SetId(rule.Id)

	return resourceWafRuleWebTamperProtectionRead(d, meta)
}

func resourceWafRuleWebTamperProtectionRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	wafClient, err := config.WafV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("error creating Flexibleengine WAF client: %s", err)
	}

	policyID := d.Get("policy_id").(string)
	n, err := rules.Get(wafClient, policyID, d.Id()).Extract()
	if err != nil {
		return CheckDeleted(d, err, "WAF Web Tamper Protection Rule")
	}

	d.SetId(n.Id)
	d.Set("policy_id", n.PolicyID)
	d.Set("domain", n.Hostname)
	d.Set("path", n.Url)

	return nil
}

func resourceWafRuleWebTamperProtectionDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	wafClient, err := config.WafV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("error creating Flexibleengine WAF client: %s", err)
	}

	policyID := d.Get("policy_id").(string)
	err = rules.Delete(wafClient, policyID, d.Id()).ExtractErr()
	if err != nil {
		return fmt.Errorf("error deleting Flexibleengine WAF Web Tamper Protection Rule: %s", err)
	}

	d.SetId("")
	return nil
}
