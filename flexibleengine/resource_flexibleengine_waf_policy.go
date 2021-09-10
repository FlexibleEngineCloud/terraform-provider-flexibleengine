package flexibleengine

import (
	"fmt"
	"log"
	"time"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/waf/v1/policies"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceWafPolicyV1() *schema.Resource {
	return &schema.Resource{
		Create: resourceWafPolicyV1Create,
		Read:   resourceWafPolicyV1Read,
		Update: resourceWafPolicyV1Update,
		Delete: resourceWafPolicyV1Delete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"domains": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"protection_mode": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "log",
				ValidateFunc: validation.StringInSlice([]string{
					"log", "block",
				}, false),
			},
			"level": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      2,
				ValidateFunc: validation.IntBetween(0, 3),
			},
			"full_detection": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"protection_status": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"basic_web_protection": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"general_check": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"crawler_engine": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"crawler_scanner": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"crawler_script": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"crawler_other": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"webshell": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"cc_protection": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"precise_protection": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"blacklist": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"data_masking": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"false_alarm_masking": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"web_tamper_protection": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func buildWafPolicyOptions(d *schema.ResourceData) *policies.Options {
	optionsRaw := d.Get("protection_status").([]interface{})
	if len(optionsRaw) == 0 {
		return nil
	}

	rawMap := optionsRaw[0].(map[string]interface{})
	webattack := rawMap["basic_web_protection"].(bool)
	common := rawMap["general_check"].(bool)
	crawlerEngine := rawMap["crawler_engine"].(bool)
	crawlerScanner := rawMap["crawler_scanner"].(bool)
	crawlerScript := rawMap["crawler_script"].(bool)
	crawlerOther := rawMap["crawler_other"].(bool)
	webshell := rawMap["webshell"].(bool)

	cc := rawMap["cc_protection"].(bool)
	custom := rawMap["precise_protection"].(bool)
	whiteblackip := rawMap["blacklist"].(bool)
	privacy := rawMap["data_masking"].(bool)
	ignore := rawMap["false_alarm_masking"].(bool)
	antitamper := rawMap["web_tamper_protection"].(bool)

	options := &policies.Options{
		WebAttack:      &webattack,
		Common:         &common,
		CrawlerEngine:  &crawlerEngine,
		CrawlerScanner: &crawlerScanner,
		CrawlerScript:  &crawlerScript,
		CrawlerOther:   &crawlerOther,
		WebShell:       &webshell,
		Cc:             &cc,
		Custom:         &custom,
		WhiteblackIp:   &whiteblackip,
		Privacy:        &privacy,
		Ignore:         &ignore,
		AntiTamper:     &antitamper,
	}

	log.Printf("[DEBUG] get WAF policy options: %#v", options)
	return options
}

func resourceWafPolicyV1Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	wafClient, err := config.WafV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating Flexibleengine WAF client: %s", err)
	}

	createOpts := policies.CreateOpts{
		Name: d.Get("name").(string),
	}
	policy, err := policies.Create(wafClient, createOpts).Extract()
	if err != nil {
		return fmt.Errorf("Error creating waf policy: %s", err)
	}

	log.Printf("[DEBUG] WAF policy created: %#v", policy)
	d.SetId(policy.Id)

	// Update the policy as POST API only supports Name argument
	var updateOpts policies.UpdateOpts

	if v, ok := d.GetOk("level"); ok && v.(int) != policy.Level {
		updateOpts.Level = v.(int)
	}

	if v, ok := d.GetOk("protection_mode"); ok && v.(string) != policy.Action.Category {
		updateOpts.Action = &policies.Action{
			Category: v.(string),
		}
	}
	if v, ok := d.GetOk("full_detection"); ok && v.(bool) != policy.FullDetection {
		detectionMode := v.(bool)
		updateOpts.FullDetection = &detectionMode
	}

	if _, ok := d.GetOk("protection_status"); ok {
		updateOpts.Options = buildWafPolicyOptions(d)
	}

	if updateOpts != (policies.UpdateOpts{}) {
		log.Printf("[DEBUG] updateOpts: %#v", updateOpts)
		_, err = policies.Update(wafClient, d.Id(), updateOpts).Extract()
		if err != nil {
			return fmt.Errorf("error updating WAF policy: %s", err)
		}
	}

	return resourceWafPolicyV1Read(d, meta)
}

func resourceWafPolicyV1Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	wafClient, err := config.WafV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("error creating Flexibleengine WAF client: %s", err)
	}

	n, err := policies.Get(wafClient, d.Id()).Extract()
	if err != nil {
		return CheckDeleted(d, err, "WAF policy")
	}

	log.Printf("[DEBUG] fetching WAF policy %s: %#v", d.Id(), n)
	d.Set("region", GetRegion(d, config))
	d.Set("name", n.Name)
	d.Set("level", n.Level)
	d.Set("protection_mode", n.Action.Category)
	d.Set("domains", n.Hosts)
	d.Set("full_detection", n.FullDetection)

	options := []map[string]interface{}{
		{
			"basic_web_protection":  *n.Options.WebAttack,
			"general_check":         *n.Options.Common,
			"crawler_engine":        *n.Options.CrawlerEngine,
			"crawler_scanner":       *n.Options.CrawlerScanner,
			"crawler_script":        *n.Options.CrawlerScript,
			"crawler_other":         *n.Options.CrawlerOther,
			"webshell":              *n.Options.WebShell,
			"cc_protection":         *n.Options.Cc,
			"precise_protection":    *n.Options.Custom,
			"blacklist":             *n.Options.WhiteblackIp,
			"data_masking":          *n.Options.Privacy,
			"false_alarm_masking":   *n.Options.Ignore,
			"web_tamper_protection": *n.Options.AntiTamper,
		},
	}
	d.Set("protection_status", options)

	return nil
}

func resourceWafPolicyV1Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	wafClient, err := config.WafV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("error creating Flexibleengine WAF Client: %s", err)
	}

	var updateOpts policies.UpdateOpts
	var changed bool

	if d.HasChange("name") {
		changed = true
		updateOpts.Name = d.Get("name").(string)
	}
	if d.HasChanges("level", "protection_mode", "full_detection") {
		changed = true
		updateOpts.Level = d.Get("level").(int)
		updateOpts.Action = &policies.Action{
			Category: d.Get("protection_mode").(string),
		}

		detectionMode := d.Get("full_detection").(bool)
		updateOpts.FullDetection = &detectionMode
	}
	if d.HasChange("protection_status") {
		changed = true
		updateOpts.Options = buildWafPolicyOptions(d)
	}

	if changed {
		log.Printf("[DEBUG] updateOpts: %#v", updateOpts)
		_, err = policies.Update(wafClient, d.Id(), updateOpts).Extract()
		if err != nil {
			return fmt.Errorf("error updating WAF policy: %s", err)
		}
	}

	if d.HasChange("domains") {
		v := d.Get("domains").([]interface{})
		hosts := make([]string, len(v))
		for i, v := range v {
			hosts[i] = v.(string)
		}

		updateHostsOpts := policies.UpdateHostsOpts{
			Hosts: hosts,
		}
		_, err = policies.UpdateHosts(wafClient, d.Id(), updateHostsOpts).Extract()
		if err != nil {
			return fmt.Errorf("error binding WAF policy to domain %v: %s", hosts, err)
		}
	}
	return resourceWafPolicyV1Read(d, meta)
}

func resourceWafPolicyV1Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	wafClient, err := config.WafV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("error creating Flexibleengine WAF client: %s", err)
	}

	if hosts, ok := d.GetOk("domains"); ok {
		log.Printf("[DEBUG] Policy already used by domain %#v, should unbind it", hosts)
		var updateHostsOpts policies.UpdateHostsOpts
		updateHostsOpts.Hosts = make([]string, 0)

		_, err = policies.UpdateHosts(wafClient, d.Id(), updateHostsOpts).Extract()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				d.SetId("")
				return nil
			}
			return fmt.Errorf("error unbinding WAF policy domain: %s", err)
		}
	}

	err = policies.Delete(wafClient, d.Id()).ExtractErr()
	if err != nil {
		return fmt.Errorf("error deleting WAF policy: %s", err)
	}

	d.SetId("")
	return nil
}
