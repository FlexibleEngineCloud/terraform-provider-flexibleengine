package flexibleengine

import (
	"fmt"
	"log"

	"github.com/chnsz/golangsdk/openstack/lts/huawei/loggroups"
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
	region := GetRegion(d, config)
	client, err := config.LtsV2Client(region)
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine LTS client: %s", err)
	}

	groups, err := loggroups.List(client).Extract()
	if err != nil {
		return fmt.Errorf("Error querying log group list: %s", err)
	}

	resourceID := d.Id()
	for _, group := range groups.LogGroups {
		if group.ID == resourceID {
			d.SetId(group.ID)
			d.Set("region", region)
			d.Set("group_name", group.Name)
			d.Set("ttl_in_days", group.TTLinDays)
			return nil
		}
	}

	log.Printf("[WARN] log group %s: resource is gone and will be removed in Terraform state", resourceID)
	d.SetId("")
	return nil
}

func resourceLTSGroupV2Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client, err := config.LtsV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine LTS client: %s", err)
	}

	err = loggroups.Delete(client, d.Id()).ExtractErr()
	if err != nil {
		return CheckDeleted(d, err, "Error deleting log group")
	}

	d.SetId("")
	return nil
}
