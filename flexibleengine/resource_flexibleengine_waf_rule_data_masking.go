package flexibleengine

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	rules "github.com/chnsz/golangsdk/openstack/waf/v1/datamasking_rules"
)

func resourceWafRuleDataMasking() *schema.Resource {
	return &schema.Resource{
		Create: resourceWafRuleDataMaskingCreate,
		Read:   resourceWafRuleDataMaskingRead,
		Update: resourceWafRuleDataMaskingUpdate,
		Delete: resourceWafRuleDataMaskingDelete,
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
			"field": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"params", "header",
				}, false),
			},
			"subfield": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceWafRuleDataMaskingCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	wafClient, err := config.WafV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("error creating Flexibleengine WAF Client: %s", err)
	}

	policyID := d.Get("policy_id").(string)
	createOpts := rules.CreateOpts{
		Path:     d.Get("path").(string),
		Category: d.Get("field").(string),
		Index:    d.Get("subfield").(string),
	}

	log.Printf("[DEBUG] WAF Data Masking Rule creating opts: %#v", createOpts)
	rule, err := rules.Create(wafClient, policyID, createOpts).Extract()
	if err != nil {
		return fmt.Errorf("error creating Flexibleengine WAF Data Masking Rule: %s", err)
	}

	log.Printf("[DEBUG] WAF data masking rule created: %#v", rule)
	d.SetId(rule.Id)

	return resourceWafRuleDataMaskingRead(d, meta)
}

func resourceWafRuleDataMaskingRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	wafClient, err := config.WafV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("error creating Flexibleengine WAF client: %s", err)
	}

	policyID := d.Get("policy_id").(string)
	n, err := rules.Get(wafClient, policyID, d.Id()).Extract()
	if err != nil {
		return CheckDeleted(d, err, "WAF Data Masking Rule")
	}
	log.Printf("[DEBUG] fetching WAF data masking rule: %#v", n)

	d.SetId(n.Id)
	d.Set("policy_id", n.PolicyID)
	d.Set("path", n.Path)
	d.Set("field", n.Category)
	d.Set("subfield", n.Index)

	return nil
}

func resourceWafRuleDataMaskingUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	wafClient, err := config.WafV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("error creating Flexibleengine WAF Client: %s", err)
	}

	if d.HasChanges("path", "field", "subfield") {
		policyID := d.Get("policy_id").(string)
		updateOpts := rules.UpdateOpts{
			Path:     d.Get("path").(string),
			Category: d.Get("field").(string),
			Index:    d.Get("subfield").(string),
		}

		log.Printf("[DEBUG] WAF Data Masking Rule updating opts: %#v", updateOpts)
		_, err = rules.Update(wafClient, policyID, d.Id(), updateOpts).Extract()
		if err != nil {
			return fmt.Errorf("error updating Flexibleengine WAF Data Masking Rule: %s", err)
		}
	}

	return resourceWafRuleDataMaskingRead(d, meta)
}

func resourceWafRuleDataMaskingDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	wafClient, err := config.WafV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("error creating Flexibleengine WAF client: %s", err)
	}

	policyID := d.Get("policy_id").(string)
	err = rules.Delete(wafClient, policyID, d.Id()).ExtractErr()
	if err != nil {
		return fmt.Errorf("error deleting Flexibleengine WAF Data Masking Rule: %s", err)
	}

	d.SetId("")
	return nil
}

func resourceWafRulesImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	parts := strings.SplitN(d.Id(), "/", 2)
	if len(parts) != 2 {
		err := fmt.Errorf("Invalid format specified for WAF rule. Format must be <policy id>/<rule id>")
		return nil, err
	}

	policyID := parts[0]
	ruleID := parts[1]

	d.SetId(ruleID)
	d.Set("policy_id", policyID)

	return []*schema.ResourceData{d}, nil
}
