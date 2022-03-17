package flexibleengine

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/chnsz/golangsdk/openstack/identity/v3.0/policies"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceIdentityRoleV3() *schema.Resource {
	return &schema.Resource{
		Create: resourceIdentityRoleCreate,
		Read:   resourceIdentityRoleRead,
		Update: resourceIdentityRoleUpdate,
		Delete: resourceIdentityRoleDelete,
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
				Required: true,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"policy": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsJSON,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					equal, _ := compareJsonTemplateAreEquivalent(old, new)
					return equal
				},
			},
			"references": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"domain_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceIdentityRoleCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	identityClient, err := config.IAMV3Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine IAM client: %s", err)
	}

	policy := policies.Policy{}
	policyDoc := d.Get("policy").(string)
	err = json.Unmarshal([]byte(policyDoc), &policy)
	if err != nil {
		return fmt.Errorf("Error unmarshalling policy, please check the format of the policy document: %s", err)
	}
	createOpts := policies.CreateOpts{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Type:        d.Get("type").(string),
		Policy:      policy,
	}

	log.Printf("[DEBUG] Create Options: %#v", createOpts)

	role, err := policies.Create(identityClient, createOpts).Extract()
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine Role: %s", err)
	}

	d.SetId(role.ID)

	return resourceIdentityRoleRead(d, meta)
}

func resourceIdentityRoleRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	identityClient, err := config.IAMV3Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine IAM client: %s", err)
	}

	role, err := policies.Get(identityClient, d.Id()).Extract()
	if err != nil {
		return CheckDeleted(d, err, "role")
	}

	log.Printf("[DEBUG] Retrieved FlexibleEngine Role: %#v", role)

	d.Set("name", role.Name)
	d.Set("description", role.Description)
	d.Set("type", role.Type)
	d.Set("references", role.References)
	d.Set("domain_id", role.DomainId)

	policy, err := json.Marshal(role.Policy)
	if err != nil {
		return fmt.Errorf("Error marshalling policy: %s", err)
	}
	d.Set("policy", string(policy))

	return nil
}

func resourceIdentityRoleUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	identityClient, err := config.IAMV3Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine IAM client: %s", err)
	}

	policy := policies.Policy{}
	policyDoc := d.Get("policy").(string)
	err = json.Unmarshal([]byte(policyDoc), &policy)
	if err != nil {
		return fmt.Errorf("Error unmarshalling policy, please check the format of the policy document: %s", err)
	}
	createOpts := policies.CreateOpts{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Type:        d.Get("type").(string),
		Policy:      policy,
	}

	log.Printf("[DEBUG] Create Options: %#v", createOpts)

	_, err = policies.Update(identityClient, d.Id(), createOpts).Extract()
	if err != nil {
		return fmt.Errorf("Error updating FlexibleEngine Role: %s", err)
	}

	return resourceIdentityRoleRead(d, meta)
}

func resourceIdentityRoleDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	identityClient, err := config.IAMV3Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine IAM client: %s", err)
	}

	err = policies.Delete(identityClient, d.Id()).ExtractErr()
	if err != nil {
		return fmt.Errorf("Error deleting FlexibleEngine Role: %s", err)
	}

	return nil
}
