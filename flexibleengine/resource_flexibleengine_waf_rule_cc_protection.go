package flexibleengine

import (
	"fmt"
	"log"

	rules "github.com/chnsz/golangsdk/openstack/waf/v1/ccattackprotection_rules"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceWafRuleCCAttackProtection() *schema.Resource {
	return &schema.Resource{
		Create: resourceWafRuleCCAttackProtectionCreate,
		Read:   resourceWafRuleCCAttackProtectionRead,
		Update: resourceWafRuleCCAttackProtectionUpdate,
		Delete: resourceWafRuleCCAttackProtectionDelete,
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
			"mode": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"ip", "cookie", "other",
				}, false),
			},
			"action": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"block", "captcha",
				}, false),
			},
			"limit_num": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"limit_period": {
				Type:     schema.TypeInt,
				Required: true,
			},

			"cookie": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"content"},
			},
			"content": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"block_time": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"block_page_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"application/json", "text/html", "text/xml",
				}, false),
			},
			"block_page_content": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				RequiredWith: []string{"block_page_type"},
			},
		},
	}
}

func buildWafRuleTagCondition(d *schema.ResourceData) *rules.TagCondition {
	if v, ok := d.GetOk("content"); ok {
		condition := rules.TagCondition{
			Category: "Referer",
			Contents: []string{v.(string)},
		}
		return &condition
	}

	return nil
}

func buildWafRuleCcAction(d *schema.ResourceData) rules.Action {
	action := rules.Action{
		Category: d.Get("action").(string),
	}

	if v, ok := d.GetOk("block_page_content"); ok {
		response := rules.Response{
			ContentType: d.Get("block_page_type").(string),
			Content:     v.(string),
		}
		action.Detail = &rules.Detail{
			Response: response,
		}
	}

	log.Printf("[DEBUG] Waf Rule CC Action opts: %#v", action)
	return action
}

func resourceWafRuleCCAttackProtectionCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	wafClient, err := config.WafV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("error creating Flexibleengine WAF Client: %s", err)
	}

	policyID := d.Get("policy_id").(string)
	limitNum := d.Get("limit_num").(int)
	limitPeriod := d.Get("limit_period").(int)
	blockTime := d.Get("block_time").(int)
	createOpts := rules.CreateOpts{
		LimitNum:     &limitNum,
		LimitPeriod:  &limitPeriod,
		LockTime:     &blockTime,
		Url:          d.Get("path").(string),
		TagType:      d.Get("mode").(string),
		TagIndex:     d.Get("cookie").(string),
		TagCondition: buildWafRuleTagCondition(d),
		Action:       buildWafRuleCcAction(d),
	}

	log.Printf("[DEBUG] WAF CC attack protection rule creating opts: %#v", createOpts)
	rule, err := rules.Create(wafClient, policyID, createOpts).Extract()
	if err != nil {
		return fmt.Errorf("error creating Flexibleengine WAF CC Attack Protection Rule: %s", err)
	}

	log.Printf("[DEBUG] Waf CC attack protection rule created: %#v", rule)
	d.SetId(rule.Id)

	return resourceWafRuleCCAttackProtectionRead(d, meta)
}

func resourceWafRuleCCAttackProtectionRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	wafClient, err := config.WafV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("error creating Flexibleengine WAF client: %s", err)
	}

	policyID := d.Get("policy_id").(string)
	n, err := rules.Get(wafClient, policyID, d.Id()).Extract()
	if err != nil {
		return CheckDeleted(d, err, "WAF CC Attack Protection Rule")
	}

	mErr := multierror.Append(nil,
		d.Set("policy_id", n.PolicyID),
		d.Set("path", n.Url),
		d.Set("limit_num", n.LimitNum),
		d.Set("limit_period", n.LimitPeriod),
		d.Set("block_time", n.LockTime),
		d.Set("mode", n.TagType),
		d.Set("cookie", n.TagIndex),
		d.Set("action", n.Action.Category),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return fmt.Errorf("error setting WAF fields: %s", err)
	}

	contens := n.TagCondition.Contents
	if len(contens) > 0 {
		d.Set("content", contens[0])
	}
	blockDetail := n.Action.Detail
	if blockDetail != nil {
		d.Set("block_page_type", blockDetail.Response.ContentType)
		d.Set("block_page_content", blockDetail.Response.Content)
	}

	return nil
}

func resourceWafRuleCCAttackProtectionUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	wafClient, err := config.WafV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("error creating Flexibleengine WAF Client: %s", err)
	}

	policyID := d.Get("policy_id").(string)
	limitNum := d.Get("limit_num").(int)
	limitPeriod := d.Get("limit_period").(int)
	blockTime := d.Get("block_time").(int)
	updateOpts := rules.CreateOpts{
		LimitNum:     &limitNum,
		LimitPeriod:  &limitPeriod,
		LockTime:     &blockTime,
		Url:          d.Get("path").(string),
		TagType:      d.Get("mode").(string),
		TagIndex:     d.Get("cookie").(string),
		TagCondition: buildWafRuleTagCondition(d),
		Action:       buildWafRuleCcAction(d),
	}

	log.Printf("[DEBUG] WAF CC attack protection rule updating opts: %#v", updateOpts)
	_, err = rules.Update(wafClient, policyID, d.Id(), updateOpts).Extract()
	if err != nil {
		return fmt.Errorf("error updating Flexibleengine WAF CC Attack Protection Rule: %s", err)
	}

	return resourceWafRuleCCAttackProtectionRead(d, meta)
}

func resourceWafRuleCCAttackProtectionDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	wafClient, err := config.WafV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("error creating Flexibleengine WAF client: %s", err)
	}

	policyID := d.Get("policy_id").(string)
	err = rules.Delete(wafClient, policyID, d.Id()).ExtractErr()
	if err != nil {
		return fmt.Errorf("error deleting Flexibleengine WAF CC Attack Protection Rule: %s", err)
	}

	d.SetId("")
	return nil
}
