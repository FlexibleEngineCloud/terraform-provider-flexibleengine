package flexibleengine

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/openstack/networking/v2/extensions/lbaas_v2/certificates"
)

func dataSourceCertificateV2() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceCertificateV2Read,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"domain": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"private_key": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"certificate": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"update_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceCertificateV2Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	lbClient, err := config.ElbV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine ELB v2.0 client: %s", err)
	}

	listOpts := certificates.ListOpts{
		ID:          d.Get("id").(string),
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Type:        d.Get("type").(string),
		Domain:      d.Get("domain").(string),
	}
	allPages, err := certificates.List(lbClient, listOpts).AllPages()
	if err != nil {
		return fmt.Errorf("Error retrieving flexibleengine_lb_certificate_v2: %s", err)
	}
	certs, err := certificates.ExtractCertificates(allPages)
	if err != nil {
		return fmt.Errorf("Error extracting flexibleengine_lb_certificate_v2 from response: %s", err)
	}

	if len(certs) < 1 {
		return fmt.Errorf("Your query returned no results. " +
			"Please change your search criteria and try again.")
	}

	if len(certs) > 1 {
		return fmt.Errorf("Your query returned more than one result." +
			" Please try a more specific search criteria")
	}

	Cert := certs[0]

	d.SetId(Cert.ID)
	d.Set("name", Cert.Name)
	d.Set("description", Cert.Description)
	d.Set("type", Cert.Type)
	d.Set("domain", Cert.Domain)
	d.Set("certificate", Cert.Certificate)
	d.Set("private_key", Cert.PrivateKey)
	d.Set("create_time", Cert.CreateTime)
	d.Set("update_time", Cert.UpdateTime)

	return nil
}
