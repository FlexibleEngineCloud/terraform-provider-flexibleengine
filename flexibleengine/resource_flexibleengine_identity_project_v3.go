package flexibleengine

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/identity/v3/projects"
)

func resourceIdentityProjectV3() *schema.Resource {
	return &schema.Resource{
		Create: resourceIdentityProjectV3Create,
		Read:   resourceIdentityProjectV3Read,
		Update: resourceIdentityProjectV3Update,
		Delete: resourceIdentityProjectV3Delete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
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
				Computed: true,
			},
			"parent_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_domain": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func resourceIdentityProjectV3Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	identityClient, err := config.identityV3Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenStack identity client: %s", err)
	}

	createOpts := projects.CreateOpts{
		DomainID:    d.Get("domain_id").(string),
		Name:        d.Get("name").(string),
		ParentID:    d.Get("parent_id").(string),
		Description: d.Get("description").(string),
	}

	log.Printf("[DEBUG] Create Options: %#v", createOpts)

	project, err := projects.Create(identityClient, createOpts).Extract()

	if err != nil {
		return fmt.Errorf("Error creating project")
	}

	d.SetId(project.ID)
	log.Printf("[INFO] Project ID: %s", project.ID)

	return resourceIdentityProjectV3Read(d, meta)
}

func resourceIdentityProjectV3Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	identityClient, err := config.identityV3Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenStack identity client: %s", err)
	}

	project, err := projects.Get(identityClient, d.Id()).Extract()
	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error retrieving Project: %s", err)
	}

	log.Printf("[DEBUG] Retrieved FlexibleEngine Project: %#v", project)

	d.Set("description", project.Description)
	d.Set("domain_id", project.DomainID)
	d.Set("name", project.Name)
	d.Set("parent_id", project.ParentID)
	d.Set("domain_id", project.DomainID)
	d.Set("enabled", project.Enabled)

	return nil
}

func resourceIdentityProjectV3Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	identityClient, err := config.identityV3Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenStack identity client: %s", err)
	}

	var hasChange bool
	var updateOpts projects.UpdateOpts

	if d.HasChange("description") {
		hasChange = true
		updateOpts.Description = d.Get("description").(string)
	}

	if d.HasChange("name") {
		hasChange = true
		updateOpts.Description = d.Get("name").(string)
	}

	if hasChange {
		log.Printf("[DEBUG] Update Options: %#v", updateOpts)

		_, err := projects.Update(identityClient, d.Id(), updateOpts).Extract()
		if err != nil {
			return fmt.Errorf("Error updating FlexibleEngine project: %s", err)
		}
	}

	return resourceIdentityProjectV3Read(d, meta)
}

func resourceIdentityProjectV3Delete(d *schema.ResourceData, meta interface{}) error {
	log.Println("[WARN] Project deletion is not supported by FlexibleEngine API")

	return nil

	// config := meta.(*Config)
	// identityClient, err := config.identityV3Client(GetRegion(d, config))
	// if err != nil {
	// 	return fmt.Errorf("Error creating OpenStack identity client: %s", err)
	// }

	// err = projects.Delete(identityClient, d.Id()).ExtractErr()
	// if err != nil {
	// 	return fmt.Errorf("Error deleting FlexibleEngine project: %s", err)
	// }

	// return nil

}
