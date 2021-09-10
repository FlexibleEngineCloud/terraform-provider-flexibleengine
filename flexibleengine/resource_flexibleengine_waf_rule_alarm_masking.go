package flexibleengine

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	rules "github.com/chnsz/golangsdk/openstack/waf/v1/falsealarmmasking_rules"
)

func resourceWafRuleAlarmMasking() *schema.Resource {
	return &schema.Resource{
		Create: resourceWafRuleAlarmMaskingCreate,
		Read:   resourceWafRuleAlarmMaskingRead,
		Update: resourceWafRuleAlarmMaskingUpdate,
		Delete: resourceWafRuleAlarmMaskingDelete,
		Importer: &schema.ResourceImporter{
			State: resourceWafRulesImport,
		},

		Schema: map[string]*schema.Schema{
			"policy_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"path": {
				Type:     schema.TypeString,
				Required: true,
			},
			"event_id": {
				Type:     schema.TypeString,
				Required: true,
			},

			"event_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceWafRuleAlarmMaskingCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	wafClient, err := config.WafV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("error creating Flexibleengine WAF Client: %s", err)
	}

	policyID := d.Get("policy_id").(string)
	createOpts := rules.CreateOpts{
		Path:    d.Get("path").(string),
		EventID: d.Get("event_id").(string),
	}

	log.Printf("[DEBUG] WAF Alarm Masking Rule creating opts: %#v", createOpts)
	rule, err := rules.Create(wafClient, policyID, createOpts).Extract()
	if err != nil {
		return fmt.Errorf("error creating Flexibleengine WAF Alarm Masking Rule: %s", err)
	}

	log.Printf("[DEBUG] WAF alarm masking rule created: %#v", rule)
	d.SetId(rule.Id)

	return resourceWafRuleAlarmMaskingRead(d, meta)
}

func resourceWafRuleAlarmMaskingRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	wafClient, err := config.WafV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("error creating Flexibleengine WAF client: %s", err)
	}

	policyID := d.Get("policy_id").(string)
	n, err := rules.Get(wafClient, policyID, d.Id()).Extract()
	if err != nil {
		return CheckDeleted(d, err, "WAF Alarm Masking Rule")
	}
	log.Printf("[DEBUG] fetching WAF alarm masking rule created: %#v", n)

	d.SetId(n.Id)
	d.Set("policy_id", n.PolicyID)
	d.Set("path", n.Path)
	d.Set("event_id", n.EventID)
	d.Set("event_type", n.EventType)

	return nil
}

func resourceWafRuleAlarmMaskingUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	wafClient, err := config.WafV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("error creating Flexibleengine WAF Client: %s", err)
	}

	if d.HasChanges("path", "event_id") {
		policyID := d.Get("policy_id").(string)
		updateOpts := rules.UpdateOpts{
			Path:    d.Get("path").(string),
			EventID: d.Get("event_id").(string),
		}

		log.Printf("[DEBUG] WAF Alarm Masking Rule updating opts: %#v", updateOpts)
		_, err = rules.Update(wafClient, policyID, d.Id(), updateOpts).Extract()
		if err != nil {
			return fmt.Errorf("error updating WAF Alarm Masking Rule: %s", err)
		}
	}

	return resourceWafRuleAlarmMaskingRead(d, meta)
}

func resourceWafRuleAlarmMaskingDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	wafClient, err := config.WafV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("error creating Flexibleengine WAF client: %s", err)
	}

	policyID := d.Get("policy_id").(string)
	err = rules.Delete(wafClient, policyID, d.Id()).ExtractErr()
	if err != nil {
		return fmt.Errorf("error deleting WAF Alarm Masking Rule: %s", err)
	}

	d.SetId("")
	return nil
}
