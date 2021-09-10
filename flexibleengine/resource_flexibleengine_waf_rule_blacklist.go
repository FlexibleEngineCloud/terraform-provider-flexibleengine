package flexibleengine

import (
	"fmt"
	"log"

	rules "github.com/chnsz/golangsdk/openstack/waf/v1/whiteblackip_rules"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceWafRuleBlackList() *schema.Resource {
	return &schema.Resource{
		Create: resourceWafRuleBlackListCreate,
		Read:   resourceWafRuleBlackListRead,
		Update: resourceWafRuleBlackListUpdate,
		Delete: resourceWafRuleBlackListDelete,
		Importer: &schema.ResourceImporter{
			State: resourceWafRulesImport,
		},

		Schema: map[string]*schema.Schema{
			"policy_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"address": {
				Type:     schema.TypeString,
				Required: true,
			},
			"action": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      0,
				ValidateFunc: validation.IntInSlice([]int{0, 1}),
			},
		},
	}
}

func resourceWafRuleBlackListCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	wafClient, err := config.WafV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("error creating Flexibleengine WAF Client: %s", err)
	}

	policyID := d.Get("policy_id").(string)
	createOpts := rules.CreateOpts{
		Addr:  d.Get("address").(string),
		White: d.Get("action").(int),
	}

	log.Printf("[DEBUG] WAF black list rule creating opts: %#v", createOpts)
	rule, err := rules.Create(wafClient, policyID, createOpts).Extract()
	if err != nil {
		return fmt.Errorf("error creating WAF black list rule: %s", err)
	}

	log.Printf("[DEBUG] WAF black list rule created: %#v", rule)
	d.SetId(rule.Id)

	return resourceWafRuleBlackListRead(d, meta)
}

func resourceWafRuleBlackListRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	wafClient, err := config.WafV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("error creating Flexibleengine WAF client: %s", err)
	}

	policyID := d.Get("policy_id").(string)
	n, err := rules.Get(wafClient, policyID, d.Id()).Extract()
	if err != nil {
		return CheckDeleted(d, err, "WAF Black List Rule")
	}
	log.Printf("[DEBUG] fetching WAF black list rule: %#v", n)

	d.SetId(n.Id)
	d.Set("policy_id", n.PolicyID)
	d.Set("address", n.Addr)
	d.Set("action", n.White)

	return nil
}

func resourceWafRuleBlackListUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	wafClient, err := config.WafV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("error creating Flexibleengine WAF Client: %s", err)
	}

	if d.HasChanges("address", "action") {
		white := d.Get("action").(int)
		updateOpts := rules.UpdateOpts{
			Addr:  d.Get("address").(string),
			White: &white,
		}
		log.Printf("[DEBUG] updateOpts: %#v", updateOpts)

		policyID := d.Get("policy_id").(string)
		_, err = rules.Update(wafClient, policyID, d.Id(), updateOpts).Extract()
		if err != nil {
			return fmt.Errorf("error updating Flexibleengine WAF WhiteBlackIP Rule: %s", err)
		}
	}

	return resourceWafRuleBlackListRead(d, meta)
}

func resourceWafRuleBlackListDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	wafClient, err := config.WafV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("error creating Flexibleengine WAF client: %s", err)
	}

	policyID := d.Get("policy_id").(string)
	err = rules.Delete(wafClient, policyID, d.Id()).ExtractErr()
	if err != nil {
		return fmt.Errorf("error deleting Flexibleengine WAF WhiteBlackIP Rule: %s", err)
	}

	d.SetId("")
	return nil
}
