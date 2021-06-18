package flexibleengine

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/openstack/waf/v1/policies"
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

			"options": {
				Type:     schema.TypeList,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"webattack": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"common": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"crawler": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"crawler_engine": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"crawler_scanner": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"crawler_script": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"crawler_other": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"webshell": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"cc": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"custom": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"whiteblackip": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"privacy": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"ignore": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"antitamper": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},

			"full_detection": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
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

	log.Printf("[DEBUG] Waf policy created: %#v", policy)
	d.SetId(policy.Id)
	d.Set("name", policy.Name)

	return resourceWafPolicyV1Update(d, meta)
}

func resourceWafPolicyV1Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	wafClient, err := config.WafV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("error creating Flexibleengine WAF client: %s", err)
	}

	n, err := policies.Get(wafClient, d.Id()).Extract()
	if err != nil {
		return CheckDeleted(d, err, "Waf Policy")
	}

	d.Set("region", GetRegion(d, config))
	d.Set("name", n.Name)
	d.Set("level", n.Level)
	d.Set("protection_mode", n.Action.Category)
	d.Set("domains", n.Hosts)
	d.Set("full_detection", n.FullDetection)

	options := []map[string]interface{}{
		{
			"webattack":       *n.Options.WebAttack,
			"common":          *n.Options.Common,
			"crawler":         *n.Options.Crawler,
			"crawler_engine":  *n.Options.CrawlerEngine,
			"crawler_scanner": *n.Options.CrawlerScanner,
			"crawler_script":  *n.Options.CrawlerScript,
			"crawler_other":   *n.Options.CrawlerOther,
			"webshell":        *n.Options.WebShell,
			"cc":              *n.Options.Cc,
			"custom":          *n.Options.Custom,
			"whiteblackip":    *n.Options.WhiteblackIp,
			"privacy":         *n.Options.Privacy,
			"ignore":          *n.Options.Ignore,
			"antitamper":      *n.Options.AntiTamper,
		},
	}
	d.Set("options", options)
	return nil
}

func resourceWafPolicyV1Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	wafClient, err := config.WafV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("error creating Flexibleengine WAF Client: %s", err)
	}

	if d.HasChanges("name", "level", "protection_mode") {
		updateOpts := policies.UpdateOpts{
			Name:  d.Get("name").(string),
			Level: d.Get("level").(int),
			Action: &policies.Action{
				Category: d.Get("protection_mode").(string),
			},
		}

		log.Printf("[DEBUG] updateOpts: %#v", updateOpts)
		_, err = policies.Update(wafClient, d.Id(), updateOpts).Extract()
		if err != nil {
			return fmt.Errorf("error updating WAF Policy: %s", err)
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
			return fmt.Errorf("error unbinding WAF Policy domain: %s", err)
		}
	}

	err = policies.Delete(wafClient, d.Id()).ExtractErr()
	if err != nil {
		return fmt.Errorf("error deleting WAF Policy: %s", err)
	}

	d.SetId("")
	return nil
}
