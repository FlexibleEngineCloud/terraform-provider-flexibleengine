package flexibleengine

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/networking/v2/extensions/security/groups"
	"github.com/chnsz/golangsdk/openstack/networking/v2/extensions/security/rules"
)

func resourceNetworkingSecGroupV2() *schema.Resource {
	return &schema.Resource{
		Create: resourceNetworkingSecGroupV2Create,
		Read:   resourceNetworkingSecGroupV2Read,
		Update: resourceNetworkingSecGroupV2Update,
		Delete: resourceNetworkingSecGroupV2Delete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
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
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"delete_default_rules": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},

			"tenant_id": {
				Type:       schema.TypeString,
				Optional:   true,
				ForceNew:   true,
				Computed:   true,
				Deprecated: "tenant_id is deprecated",
			},
		},
	}
}

func resourceNetworkingSecGroupV2Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	networkingClient, err := config.NetworkingV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine networking client: %s", err)
	}

	opts := groups.CreateOpts{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		TenantID:    d.Get("tenant_id").(string),
	}

	log.Printf("[DEBUG] Create FlexibleEngine Security Group: %#v", opts)

	securityGroup, err := groups.Create(networkingClient, opts).Extract()
	if err != nil {
		return err
	}

	// Delete the default security group rules if it has been requested.
	deleteDefaultRules := d.Get("delete_default_rules").(bool)
	if deleteDefaultRules {
		securityGroup, err := groups.Get(networkingClient, securityGroup.ID).Extract()
		if err != nil {
			return err
		}
		for _, rule := range securityGroup.Rules {
			if err := rules.Delete(networkingClient, rule.ID).ExtractErr(); err != nil {
				return fmt.Errorf(
					"There was a problem deleting a default security group rule: %s", err)
			}
		}
	}

	log.Printf("[DEBUG] FlexibleEngine Security Group created: %#v", securityGroup)
	d.SetId(securityGroup.ID)

	return resourceNetworkingSecGroupV2Read(d, meta)
}

func resourceNetworkingSecGroupV2Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	region := GetRegion(d, config)
	networkingClient, err := config.NetworkingV2Client(region)
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine networking client: %s", err)
	}

	securityGroup, err := groups.Get(networkingClient, d.Id()).Extract()
	if err != nil {
		return CheckDeleted(d, err, "FlexibleEngine Security group")
	}

	d.Set("region", region)
	d.Set("name", securityGroup.Name)
	d.Set("description", securityGroup.Description)

	return nil
}

func resourceNetworkingSecGroupV2Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	networkingClient, err := config.NetworkingV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine networking client: %s", err)
	}

	var update bool
	var updateOpts groups.UpdateOpts

	if d.HasChange("name") {
		update = true
		updateOpts.Name = d.Get("name").(string)
	}

	if d.HasChange("description") {
		update = true
		description := d.Get("description").(string)
		updateOpts.Description = &description
	}

	if update {
		log.Printf("[DEBUG] Updating SecGroup %s with options: %#v", d.Id(), updateOpts)
		_, err = groups.Update(networkingClient, d.Id(), updateOpts).Extract()
		if err != nil {
			return fmt.Errorf("Error updating FlexibleEngine SecGroup: %s", err)
		}
	}

	return resourceNetworkingSecGroupV2Read(d, meta)
}

func resourceNetworkingSecGroupV2Delete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Destroy security group: %s", d.Id())

	config := meta.(*Config)
	networkingClient, err := config.NetworkingV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine networking client: %s", err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"ACTIVE"},
		Target:     []string{"DELETED"},
		Refresh:    waitForSecGroupDelete(networkingClient, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Error deleting FlexibleEngine Security Group: %s", err)
	}

	d.SetId("")
	return err
}

func waitForSecGroupDelete(networkingClient *golangsdk.ServiceClient, secGroupId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		log.Printf("[DEBUG] Attempting to delete FlexibleEngine Security Group %s.\n", secGroupId)

		r, err := groups.Get(networkingClient, secGroupId).Extract()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				log.Printf("[DEBUG] Successfully deleted FlexibleEngine Security Group %s", secGroupId)
				return r, "DELETED", nil
			}
			return r, "ACTIVE", err
		}

		err = groups.Delete(networkingClient, secGroupId).ExtractErr()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				log.Printf("[DEBUG] Successfully deleted FlexibleEngine Security Group %s", secGroupId)
				return r, "DELETED", nil
			}
			if errCode, ok := err.(golangsdk.ErrUnexpectedResponseCode); ok {
				if errCode.Actual == 409 {
					return r, "ACTIVE", nil
				}
			}
			return r, "ACTIVE", err
		}

		log.Printf("[DEBUG] FlexibleEngine Security Group %s still active.\n", secGroupId)
		return r, "ACTIVE", nil
	}
}
