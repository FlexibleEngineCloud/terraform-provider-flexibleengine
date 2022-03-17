package flexibleengine

import (
	"fmt"
	"log"

	"github.com/chnsz/golangsdk/openstack/identity/v3/groups"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceIdentityGroupV3() *schema.Resource {
	return &schema.Resource{
		Create: resourceIdentityGroupV3Create,
		Read:   resourceIdentityGroupV3Read,
		Update: resourceIdentityGroupV3Update,
		Delete: resourceIdentityGroupV3Delete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"domain_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceIdentityGroupV3Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	identityClient, err := config.IdentityV3Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine identity client: %s", err)
	}

	createOpts := groups.CreateOpts{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		DomainID:    d.Get("domain_id").(string),
	}

	log.Printf("[DEBUG] Create Options: %#v", createOpts)

	group, err := groups.Create(identityClient, createOpts).Extract()
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine Group: %s", err)
	}

	d.SetId(group.ID)

	return resourceIdentityGroupV3Read(d, meta)
}

func resourceIdentityGroupV3Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	identityClient, err := config.IdentityV3Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine identity client: %s", err)
	}

	group, err := groups.Get(identityClient, d.Id()).Extract()
	if err != nil {
		return CheckDeleted(d, err, "group")
	}

	log.Printf("[DEBUG] Retrieved FlexibleEngine Group: %#v", group)

	d.Set("domain_id", group.DomainID)
	d.Set("description", group.Description)
	d.Set("name", group.Name)

	return nil
}

func resourceIdentityGroupV3Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	identityClient, err := config.IdentityV3Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine identity client: %s", err)
	}

	var hasChange bool
	var updateOpts groups.UpdateOpts

	if d.HasChange("description") {
		hasChange = true
		updateOpts.Description = d.Get("description").(string)
	}

	if d.HasChange("domain_id") {
		hasChange = true
		updateOpts.DomainID = d.Get("domain_id").(string)
	}

	if d.HasChange("name") {
		hasChange = true
		updateOpts.Name = d.Get("name").(string)
	}

	if hasChange {
		log.Printf("[DEBUG] Update Options: %#v", updateOpts)
	}

	if hasChange {
		_, err := groups.Update(identityClient, d.Id(), updateOpts).Extract()
		if err != nil {
			return fmt.Errorf("Error updating FlexibleEngine group: %s", err)
		}
	}

	return resourceIdentityGroupV3Read(d, meta)
}

func resourceIdentityGroupV3Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	identityClient, err := config.IdentityV3Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine identity client: %s", err)
	}

	err = groups.Delete(identityClient, d.Id()).ExtractErr()
	if err != nil {
		return fmt.Errorf("Error deleting FlexibleEngine group: %s", err)
	}

	return nil
}
