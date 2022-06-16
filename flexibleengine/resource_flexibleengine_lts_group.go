package flexibleengine

import (
	"fmt"
	"log"
	"strings"

	"github.com/chnsz/golangsdk/openstack/lts/v2/loggroups"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceLTSGroupV2() *schema.Resource {
	return &schema.Resource{
		Create: resourceLTSGroupV2Create,
		Read:   resourceLTSGroupV2Read,
		Delete: resourceLTSGroupV2Delete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"group_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"ttl_in_days": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func resourceLTSGroupV2Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client, err := config.LtsV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine LTS client: %s", err)
	}
	client.ResourceBase = strings.Replace(client.ResourceBase, "/v2/", "/v2.0/", 1)

	createOpts := &loggroups.CreateOpts{
		LogGroupName: d.Get("group_name").(string),
		TTL:          7, // fixed to 7days
	}

	log.Printf("[DEBUG] Create Options: %#v", createOpts)
	groupCreate, err := loggroups.Create(client, createOpts).Extract()
	if err != nil {
		return fmt.Errorf("Error creating log group: %s", err)
	}

	d.SetId(groupCreate.ID)
	return resourceLTSGroupV2Read(d, meta)
}

func resourceLTSGroupV2Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client, err := config.LtsV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine LTS client: %s", err)
	}
	client.ResourceBase = strings.Replace(client.ResourceBase, "/v2/", "/v2.0/", 1)

	group, err := loggroups.Get(client, d.Id()).Extract()
	if err != nil {
		return CheckDeleted(d, err, "Error querying log group")
	}

	d.Set("group_name", group.Name)
	d.Set("ttl_in_days", group.TTLinDays)
	d.Set("region", GetRegion(d, config))

	return nil

}

func resourceLTSGroupV2Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client, err := config.LtsV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine LTS client: %s", err)
	}
	client.ResourceBase = strings.Replace(client.ResourceBase, "/v2/", "/v2.0/", 1)

	err = loggroups.Delete(client, d.Id()).ExtractErr()
	if err != nil {
		return CheckDeleted(d, err, "Error deleting log group")
	}

	d.SetId("")
	return nil
}
