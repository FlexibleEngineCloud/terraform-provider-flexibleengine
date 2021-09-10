package flexibleengine

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/chnsz/golangsdk/openstack/waf/v1/certificates"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceWafCertificateV1() *schema.Resource {
	return &schema.Resource{
		Create: resourceWafCertificateV1Create,
		Read:   resourceWafCertificateV1Read,
		Update: resourceWafCertificateV1Update,
		Delete: resourceWafCertificateV1Delete,
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
				Computed: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`^[\w-]{1,256}$`),
					"The maximum length is 256 characters. Only digits, letters, underscores (_), and hyphens (-) are allowed"),
			},
			"certificate": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"private_key": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"expiration": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceWafCertificateV1Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	wafClient, err := config.WafV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("error creating Flexibleengine WAF Client: %s", err)
	}

	createOpts := certificates.CreateOpts{
		Name:    d.Get("name").(string),
		Content: strings.TrimSpace(d.Get("certificate").(string)),
		Key:     strings.TrimSpace(d.Get("private_key").(string)),
	}

	certificate, err := certificates.Create(wafClient, createOpts).Extract()
	if err != nil {
		return fmt.Errorf("error creating WAF Certificate: %w", err)
	}

	log.Printf("[DEBUG] Waf certificate created: %#v", certificate)
	d.SetId(certificate.Id)

	return resourceWafCertificateV1Read(d, meta)
}

func resourceWafCertificateV1Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	wafClient, err := config.WafV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("error creating Flexibleengine WAF Client: %s", err)
	}

	n, err := certificates.Get(wafClient, d.Id()).Extract()
	if err != nil {
		return CheckDeleted(d, err, "Waf Certificate")
	}

	expires := time.Unix(int64(n.ExpireTime/1000), 0).UTC().Format("2006-01-02 15:04:05 MST")

	d.Set("region", GetRegion(d, config))
	d.Set("name", n.Name)
	d.Set("expiration", expires)

	return nil
}

func resourceWafCertificateV1Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	wafClient, err := config.WafV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("error creating Flexibleengine WAF Client: %s", err)
	}

	updateOpts := certificates.UpdateOpts{
		Name: d.Get("name").(string),
	}

	_, err = certificates.Update(wafClient, d.Id(), updateOpts).Extract()
	if err != nil {
		return fmt.Errorf("error updating WAF Certificate: %w", err)
	}
	return resourceWafCertificateV1Read(d, meta)
}

func resourceWafCertificateV1Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	wafClient, err := config.WafV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("error creating Flexibleengine WAF Client: %s", err)
	}

	err = certificates.Delete(wafClient, d.Id()).ExtractErr()
	if err != nil {
		return fmt.Errorf("error deleting WAF Certificate: %s", err)
	}

	d.SetId("")
	return nil
}
