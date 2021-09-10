package flexibleengine

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	rules "github.com/chnsz/golangsdk/openstack/waf/v1/preciseprotection_rules"
)

func resourceWafRulePreciseProtection() *schema.Resource {
	return &schema.Resource{
		Create: resourceWafRulePreciseProtectionCreate,
		Read:   resourceWafRulePreciseProtectionRead,
		Update: resourceWafRulePreciseProtectionUpdate,
		Delete: resourceWafRulePreciseProtectionDelete,
		Importer: &schema.ResourceImporter{
			State: resourceWafRulesImport,
		},

		Schema: map[string]*schema.Schema{
			"policy_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"action": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "block",
				ValidateFunc: validation.StringInSlice([]string{
					"block", "pass",
				}, false),
			},
			"priority": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntBetween(0, 65535),
			},
			"conditions": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 30,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"field": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"path", "user-agent", "ip", "params", "cookie", "referer", "header",
							}, false),
						},
						"subfield": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"logic": {
							Type:     schema.TypeString,
							Required: true,
						},
						"content": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"start_time": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: IsRFC3339Time,
			},
			"end_time": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: IsRFC3339Time,
			},
		},
	}
}

func buildWafRuleProtectionConditions(d *schema.ResourceData) []rules.Condition {
	conditions := d.Get("conditions").([]interface{})
	conditionOpts := make([]rules.Condition, len(conditions))

	for i, v := range conditions {
		cond := v.(map[string]interface{})
		conditionOpts[i] = rules.Condition{
			Category: cond["field"].(string),
			Index:    cond["subfield"].(string),
			Logic:    cond["logic"].(string),
			Contents: []string{cond["content"].(string)},
		}
	}

	log.Printf("[DEBUG] build Protection Condition Rule opts: %#v", conditionOpts)
	return conditionOpts
}

func buildWafRuleProtectionAction(d *schema.ResourceData) rules.Action {
	action := rules.Action{
		Category: d.Get("action").(string),
	}
	return action
}

func convertTimeToUnix(date string) (int64, error) {
	var err error
	if t, err := time.Parse(RFC3339ZNoTNoZ, date); err == nil {
		return t.Unix(), nil
	}
	return 0, fmt.Errorf("[%s] is not an expected date: %+v", date, err)
}

func buildWafRuleProtectionOpts(d *schema.ResourceData) (*rules.CreateOpts, error) {
	priority := d.Get("priority").(int)
	opts := rules.CreateOpts{
		Name:       d.Get("name").(string),
		Priority:   &priority,
		Conditions: buildWafRuleProtectionConditions(d),
		Action:     buildWafRuleProtectionAction(d),
	}

	if v, ok := d.GetOk("start_time"); ok {
		start, err := convertTimeToUnix(v.(string))
		if err != nil {
			return nil, fmt.Errorf("error converting start time: %s", err)
		}
		opts.Time = true
		opts.Start = start
	}
	if v, ok := d.GetOk("end_time"); ok {
		end, err := convertTimeToUnix(v.(string))
		if err != nil {
			return nil, fmt.Errorf("error converting end time: %s", err)
		}
		opts.Time = true
		opts.End = end
	}

	return &opts, nil
}

func resourceWafRulePreciseProtectionCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	wafClient, err := config.WafV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("error creating Flexibleengine WAF Client: %s", err)
	}

	createOpts, err := buildWafRuleProtectionOpts(d)
	if err != nil {
		return fmt.Errorf("error building WAF Precise Protection Rule opts: %s", err)
	}
	log.Printf("[DEBUG] WAF precise protection rule creating opts: %#v", createOpts)

	policyID := d.Get("policy_id").(string)
	rule, err := rules.Create(wafClient, policyID, createOpts).Extract()
	if err != nil {
		return fmt.Errorf("error creating Flexibleengine WAF Precise Protection Rule: %s", err)
	}

	log.Printf("[DEBUG] WAF precise protection rule created: %#v", rule)
	d.SetId(rule.Id)

	return resourceWafRulePreciseProtectionRead(d, meta)
}

func resourceWafRulePreciseProtectionUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	wafClient, err := config.WafV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("error creating Flexibleengine WAF Client: %s", err)
	}

	updateOpts, err := buildWafRuleProtectionOpts(d)
	if err != nil {
		return fmt.Errorf("error building WAF Precise Protection Rule opts: %s", err)
	}
	log.Printf("[DEBUG] WAF precise protection rule updating opts: %#v", updateOpts)

	policyID := d.Get("policy_id").(string)
	_, err = rules.Update(wafClient, policyID, d.Id(), updateOpts).Extract()
	if err != nil {
		return fmt.Errorf("error updating Flexibleengine WAF Precise Protection Rule: %s", err)
	}

	return resourceWafRulePreciseProtectionRead(d, meta)
}

func resourceWafRulePreciseProtectionRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	wafClient, err := config.WafV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("error creating Flexibleengine WAF client: %s", err)
	}

	policyID := d.Get("policy_id").(string)
	n, err := rules.Get(wafClient, policyID, d.Id()).Extract()
	if err != nil {
		return CheckDeleted(d, err, "WAF Precise Protection Rule")
	}

	d.SetId(n.Id)
	d.Set("policy_id", n.PolicyID)
	d.Set("name", n.Name)
	d.Set("action", n.Action.Category)
	d.Set("priority", n.Priority)

	conditions := make([]map[string]interface{}, len(n.Conditions))
	for i, condition := range n.Conditions {
		conditions[i] = make(map[string]interface{})
		conditions[i]["field"] = condition.Category
		conditions[i]["subfield"] = condition.Index
		conditions[i]["logic"] = condition.Logic
		if len(condition.Contents) > 0 {
			conditions[i]["content"] = condition.Contents[0]
		}
	}
	d.Set("conditions", conditions)

	if n.Start != 0 {
		startTime := time.Unix(n.Start, 0).UTC().Format(RFC3339ZNoTNoZ)
		d.Set("start_time", startTime)
	}
	if n.End != 0 {
		endTime := time.Unix(n.End, 0).UTC().Format(RFC3339ZNoTNoZ)
		d.Set("end_time", endTime)
	}

	return nil
}

func resourceWafRulePreciseProtectionDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	wafClient, err := config.WafV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("error creating Flexibleengine WAF client: %s", err)
	}

	policyID := d.Get("policy_id").(string)
	err = rules.Delete(wafClient, policyID, d.Id()).ExtractErr()
	if err != nil {
		return fmt.Errorf("error deleting Flexibleengine WAF Precise Protection Rule: %s", err)
	}

	d.SetId("")
	return nil
}
