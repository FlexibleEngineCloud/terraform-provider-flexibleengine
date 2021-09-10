package flexibleengine

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/networking/v2/extensions/natgateways"
)

func resourceNatGatewayV2() *schema.Resource {
	return &schema.Resource{
		Create: resourceNatGatewayV2Create,
		Read:   resourceNatGatewayV2Read,
		Update: resourceNatGatewayV2Update,
		Delete: resourceNatGatewayV2Delete,

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
				ForceNew: false,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: false,
			},
			"spec": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     false,
				ValidateFunc: resourceNatGatewayV2ValidateSpec,
			},
			"tenant_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"router_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"internal_network_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceNatGatewayV2Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	natV2Client, err := config.natV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine nat client: %s", err)
	}

	createOpts := &natgateways.CreateOpts{
		Name:              d.Get("name").(string),
		Description:       d.Get("description").(string),
		Spec:              d.Get("spec").(string),
		TenantID:          d.Get("tenant_id").(string),
		RouterID:          d.Get("router_id").(string),
		InternalNetworkID: d.Get("internal_network_id").(string),
	}

	log.Printf("[DEBUG] Create Options: %#v", createOpts)
	natGateway, err := natgateways.Create(natV2Client, createOpts).Extract()
	if err != nil {
		return fmt.Errorf("Error creatting Nat Gateway: %s", err)
	}

	log.Printf("[DEBUG] Waiting for FlexibleEngine Nat Gateway (%s) to become available.", natGateway.ID)

	stateConf := &resource.StateChangeConf{
		Target:     []string{"ACTIVE"},
		Refresh:    waitForNatGatewayActive(natV2Client, natGateway.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine Nat Gateway: %s", err)
	}

	d.SetId(natGateway.ID)

	return resourceNatGatewayV2Read(d, meta)
}

func resourceNatGatewayV2Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	natV2Client, err := config.natV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine nat client: %s", err)
	}

	natGateway, err := natgateways.Get(natV2Client, d.Id()).Extract()
	if err != nil {
		return CheckDeleted(d, err, "Nat Gateway")
	}

	d.Set("name", natGateway.Name)
	d.Set("description", natGateway.Description)
	d.Set("spec", natGateway.Spec)
	d.Set("router_id", natGateway.RouterID)
	d.Set("internal_network_id", natGateway.InternalNetworkID)
	d.Set("tenant_id", natGateway.TenantID)

	d.Set("region", GetRegion(d, config))

	return nil
}

func resourceNatGatewayV2Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	natV2Client, err := config.natV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine nat client: %s", err)
	}

	var updateOpts natgateways.UpdateOpts

	if d.HasChange("name") {
		updateOpts.Name = d.Get("name").(string)
	}
	if d.HasChange("description") {
		updateOpts.Description = d.Get("description").(string)
	}
	if d.HasChange("spec") {
		updateOpts.Spec = d.Get("spec").(string)
	}

	log.Printf("[DEBUG] Update Options: %#v", updateOpts)

	_, err = natgateways.Update(natV2Client, d.Id(), updateOpts).Extract()
	if err != nil {
		return fmt.Errorf("Error updating Nat Gateway: %s", err)
	}

	return resourceNatGatewayV2Read(d, meta)
}

func resourceNatGatewayV2Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	natV2Client, err := config.natV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine nat client: %s", err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"ACTIVE"},
		Target:     []string{"DELETED"},
		Refresh:    waitForNatGatewayDelete(natV2Client, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Error deleting FlexibleEngine Nat Gateway: %s", err)
	}

	d.SetId("")
	return nil
}

func waitForNatGatewayActive(natV2Client *golangsdk.ServiceClient, nId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		n, err := natgateways.Get(natV2Client, nId).Extract()
		if err != nil {
			return nil, "", err
		}

		log.Printf("[DEBUG] FlexibleEngine Nat Gateway: %+v", n)
		if n.Status == "ACTIVE" {
			return n, "ACTIVE", nil
		}

		return n, "", nil
	}
}

func waitForNatGatewayDelete(natV2Client *golangsdk.ServiceClient, nId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		log.Printf("[DEBUG] Attempting to delete FlexibleEngine Nat Gateway %s.\n", nId)

		n, err := natgateways.Get(natV2Client, nId).Extract()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				log.Printf("[DEBUG] Successfully deleted FlexibleEngine Nat gateway %s", nId)
				return n, "DELETED", nil
			}
			return n, "ACTIVE", err
		}

		err = natgateways.Delete(natV2Client, nId).ExtractErr()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				log.Printf("[DEBUG] Successfully deleted FlexibleEngine Nat Gateway %s", nId)
				return n, "DELETED", nil
			}
			return n, "ACTIVE", err
		}

		log.Printf("[DEBUG] FlexibleEngine Nat Gateway %s still active.\n", nId)
		return n, "ACTIVE", nil
	}
}

var Specs = [4]string{"1", "2", "3", "4"}

func resourceNatGatewayV2ValidateSpec(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	for i := range Specs {
		if value == Specs[i] {
			return
		}
	}
	errors = append(errors, fmt.Errorf("%q must be one of %v", k, Specs))
	return
}
