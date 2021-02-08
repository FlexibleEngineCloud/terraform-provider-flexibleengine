package flexibleengine

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/huaweicloud/golangsdk/openstack/compute/v2/extensions/servergroups"
)

func resourceComputeServerGroupV2() *schema.Resource {
	return &schema.Resource{
		Create: resourceComputeServerGroupV2Create,
		Read:   resourceComputeServerGroupV2Read,
		Update: nil,
		Delete: resourceComputeServerGroupV2Delete,
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

			"name": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"policies": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"members": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"value_specs": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceComputeServerGroupV2Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	computeClient, err := config.computeV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine compute client: %s", err)
	}

	groupPolicies := resourceServerGroupPoliciesV2(d)
	for _, p := range groupPolicies {
		if p != "anti-affinity" {
			return fmt.Errorf("Only the anti-affinity policy is supported")
		}
	}

	createOpts := ServerGroupCreateOpts{
		servergroups.CreateOpts{
			Name:     d.Get("name").(string),
			Policies: groupPolicies,
		},
		MapValueSpecs(d),
	}

	log.Printf("[DEBUG] Create Options: %#v", createOpts)
	newSG, err := servergroups.Create(computeClient, createOpts).Extract()
	if err != nil {
		return fmt.Errorf("Error creating ServerGroup: %s", err)
	}

	d.SetId(newSG.ID)

	return resourceComputeServerGroupV2Read(d, meta)
}

func resourceComputeServerGroupV2Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	computeClient, err := config.computeV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine compute client: %s", err)
	}

	sg, err := servergroups.Get(computeClient, d.Id()).Extract()
	if err != nil {
		return CheckDeleted(d, err, "server group")
	}

	log.Printf("[DEBUG] Retrieved ServerGroup %s: %+v", d.Id(), sg)

	// Set the name
	d.Set("name", sg.Name)

	// Set the policies
	policies := []string{}
	for _, p := range sg.Policies {
		policies = append(policies, p)
	}
	d.Set("policies", policies)

	// Set the members
	members := []string{}
	for _, m := range sg.Members {
		members = append(members, m)
	}
	d.Set("members", members)

	d.Set("region", GetRegion(d, config))

	return nil
}

func resourceComputeServerGroupV2Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	computeClient, err := config.computeV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine compute client: %s", err)
	}

	log.Printf("[DEBUG] Deleting ServerGroup %s", d.Id())
	if err := servergroups.Delete(computeClient, d.Id()).ExtractErr(); err != nil {
		return fmt.Errorf("Error deleting ServerGroup: %s", err)
	}

	return nil
}

func resourceServerGroupPoliciesV2(d *schema.ResourceData) []string {
	rawPolicies := d.Get("policies").([]interface{})
	policies := make([]string, len(rawPolicies))
	for i, raw := range rawPolicies {
		policies[i] = raw.(string)
	}
	return policies
}
