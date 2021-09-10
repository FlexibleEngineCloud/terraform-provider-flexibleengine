package flexibleengine

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/chnsz/golangsdk/openstack/waf/v1/domains"
	"github.com/chnsz/golangsdk/openstack/waf/v1/policies"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceWafDomainV1() *schema.Resource {
	return &schema.Resource{
		Create: resourceWafDomainV1Create,
		Read:   resourceWafDomainV1Read,
		Update: resourceWafDomainV1Update,
		Delete: resourceWafDomainV1Delete,
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
			"domain": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"server": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: false,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"client_protocol": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"server_protocol": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"address": {
							Type:     schema.TypeString,
							Required: true,
						},
						"port": {
							Type:         schema.TypeInt,
							ValidateFunc: validation.IntBetween(0, 65535),
							Required:     true,
						},
					},
				},
			},
			"certificate_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"policy_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"proxy": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"sip_header_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"default", "cloudflare", "akamai", "custom"}, true),
			},
			"sip_header_list": {
				Type:         schema.TypeList,
				Optional:     true,
				RequiredWith: []string{"sip_header_name"},
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"keep_policy": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},

			"cname": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"txt_code": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"sub_domain": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"protect_status": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"access_status": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"protocol": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func buildWafDomainServers(d *schema.ResourceData) []domains.ServerOpts {
	servers := d.Get("server").([]interface{})

	serverOpts := make([]domains.ServerOpts, len(servers))
	for i, v := range servers {
		server := v.(map[string]interface{})
		serverOpts[i] = domains.ServerOpts{
			ClientProtocol: server["client_protocol"].(string),
			ServerProtocol: server["server_protocol"].(string),
			Address:        server["address"].(string),
			Port:           strconv.Itoa(server["port"].(int)),
		}
	}

	log.Printf("[DEBUG] build WAF domain ServerOpts: %#v", serverOpts)
	return serverOpts
}

func resourceWafDomainV1Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	wafClient, err := config.WafV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("error creating Flexibleengine WAF Client: %s", err)
	}

	var hosts []string
	if v, ok := d.GetOk("policy_id"); ok {
		policyID := v.(string)
		policy, err := policies.Get(wafClient, policyID).Extract()
		if err != nil {
			return fmt.Errorf("error retrieving Waf Policy %s: %s", policyID, err)
		}
		hosts = append(hosts, policy.Hosts...)
	}

	proxy := d.Get("proxy").(bool)
	if _, ok := d.GetOk("sip_header_name"); ok && !proxy {
		return fmt.Errorf("sip_header_name is only valid when proxy is true")
	}
	v := d.Get("sip_header_list").([]interface{})
	headers := make([]string, len(v))
	for i, v := range v {
		headers[i] = v.(string)
	}

	createOpts := domains.CreateOpts{
		HostName:      d.Get("domain").(string),
		CertificateId: d.Get("certificate_id").(string),
		Servers:       buildWafDomainServers(d),
		Proxy:         &proxy,
		SipHeaderName: d.Get("sip_header_name").(string),
		SipHeaderList: headers,
	}
	log.Printf("[DEBUG] CreateOpts: %#v", createOpts)

	domain, err := domains.Create(wafClient, createOpts).Extract()
	if err != nil {
		return fmt.Errorf("error creating WAF Domain: %s", err)
	}

	log.Printf("[DEBUG] Waf domain created: %#v", domain)
	d.SetId(domain.Id)

	if v, ok := d.GetOk("policy_id"); ok {
		policyID := v.(string)
		hosts = append(hosts, d.Id())
		updateHostsOpts := policies.UpdateHostsOpts{
			Hosts: hosts,
		}

		log.Printf("[DEBUG] Bind Waf domain %s to policy %s", d.Id(), policyID)
		_, err = policies.UpdateHosts(wafClient, policyID, updateHostsOpts).Extract()
		if err != nil {
			return fmt.Errorf("error updating WAF Policy Hosts: %s", err)
		}

		// delete the policy that was auto-created by domain
		err = policies.Delete(wafClient, domain.PolicyID).ExtractErr()
		if err != nil {
			log.Printf("[WARN] error deleting WAF Policy %s: %s", domain.PolicyID, err)
		}
	}

	return resourceWafDomainV1Read(d, meta)
}

func resourceWafDomainV1Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	wafClient, err := config.WafV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("error creating Flexibleengine WAF client: %s", err)
	}

	n, err := domains.Get(wafClient, d.Id()).Extract()
	if err != nil {
		return CheckDeleted(d, err, "WAF Domain")
	}

	mErr := multierror.Append(nil,
		d.Set("region", GetRegion(d, config)),
		d.Set("domain", n.HostName),
		d.Set("certificate_id", n.CertificateId),
		d.Set("policy_id", n.PolicyID),
		d.Set("proxy", n.Proxy),
		d.Set("sip_header_name", n.SipHeaderName),
		d.Set("sip_header_list", n.SipHeaderList),
		d.Set("cname", n.Cname),
		d.Set("txt_code", n.TxtCode),
		d.Set("sub_domain", n.SubDomain),
		d.Set("protect_status", n.ProtectStatus),
		d.Set("access_status", n.AccessStatus),
		d.Set("protocol", n.Protocol),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return fmt.Errorf("error setting WAF fields: %s", err)
	}

	servers := make([]map[string]interface{}, len(n.Servers))
	for i, server := range n.Servers {
		servers[i] = map[string]interface{}{
			"client_protocol": server.ClientProtocol,
			"server_protocol": server.ServerProtocol,
			"address":         server.Address,
			"port":            server.Port,
		}
	}
	d.Set("server", servers)

	return nil
}

func resourceWafDomainV1Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	wafClient, err := config.WafV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("error creating Flexibleengine WAF Client: %s", err)
	}

	if d.HasChanges("certificate_id", "server", "proxy", "sip_header_name", "sip_header_list") {
		proxy := d.Get("proxy").(bool)
		headerName := d.Get("sip_header_name").(string)
		if headerName != "" && !proxy {
			return fmt.Errorf("sip_header_name is only valid when proxy is true")
		}

		updateOpts := domains.UpdateOpts{
			CertificateId: d.Get("certificate_id").(string),
			Servers:       buildWafDomainServers(d),
			Proxy:         &proxy,
			SipHeaderName: headerName,
		}

		v := d.Get("sip_header_list").([]interface{})
		headers := make([]string, len(v))
		for i, v := range v {
			headers[i] = v.(string)
		}
		updateOpts.SipHeaderList = headers

		log.Printf("[DEBUG] updateOpts: %#v", updateOpts)

		_, err = domains.Update(wafClient, d.Id(), updateOpts).Extract()
		if err != nil {
			return fmt.Errorf("error updating WAF Domain: %s", err)
		}
	}
	return resourceWafDomainV1Read(d, meta)
}

func resourceWafDomainV1Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	wafClient, err := config.WafV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("error creating Flexibleengine WAF client: %s", err)
	}

	delOpts := domains.DeleteOpts{
		KeepPolicy: d.Get("keep_policy").(bool),
	}

	err = domains.Delete(wafClient, d.Id(), delOpts).ExtractErr()
	if err != nil {
		return fmt.Errorf("error deleting WAF Domain: %s", err)
	}

	d.SetId("")
	return nil
}
