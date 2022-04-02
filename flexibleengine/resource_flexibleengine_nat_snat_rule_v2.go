package flexibleengine

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/networking/v2/extensions/snatrules"
)

func resourceNatSnatRuleV2() *schema.Resource {
	return &schema.Resource{
		Create: resourceNatSnatRuleV2Create,
		Read:   resourceNatSnatRuleV2Read,
		Delete: resourceNatSnatRuleV2Delete,

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
			"nat_gateway_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"floating_ip_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"source_type": {
				Type:         schema.TypeInt,
				ValidateFunc: validation.IntBetween(0, 1),
				Optional:     true,
				ForceNew:     true,
			},
			"subnet_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ExactlyOneOf: []string{"cidr", "network_id"},
			},
			"cidr": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ExactlyOneOf: []string{"subnet_id", "network_id"},
			},

			"floating_ip_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},

			// deprecated
			"network_id": {
				Type:       schema.TypeString,
				Optional:   true,
				ForceNew:   true,
				Deprecated: "use subnet_id instead",
			},
		},
	}
}

func resourceNatSnatRuleV2Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	natClient, err := config.NatV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine nat client: %s", err)
	}

	var subnetID string
	if v, ok := d.GetOk("subnet_id"); ok {
		subnetID = v.(string)
	} else {
		subnetID = d.Get("network_id").(string)
	}

	sourceType := d.Get("source_type").(int)
	if sourceType == 1 && subnetID != "" {
		return fmt.Errorf("source_type and subnet_id is incompatible in the Direct Connect scenario (source_type=1)")
	}

	createOpts := &snatrules.CreateOpts{
		NatGatewayID: d.Get("nat_gateway_id").(string),
		FloatingIPID: d.Get("floating_ip_id").(string),
		Cidr:         d.Get("cidr").(string),
		NetworkID:    subnetID,
		SourceType:   sourceType,
	}

	log.Printf("[DEBUG] Create Options: %#v", createOpts)
	snatRule, err := snatrules.Create(natClient, createOpts).Extract()
	if err != nil {
		return fmt.Errorf("Error creatting Snat Rule: %s", err)
	}

	log.Printf("[DEBUG] Waiting for FlexibleEngine Snat Rule (%s) to become available.", snatRule.ID)

	stateConf := &resource.StateChangeConf{
		Target:     []string{"ACTIVE"},
		Refresh:    waitForSnatRuleActive(natClient, snatRule.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine Snat Rule: %s", err)
	}

	d.SetId(snatRule.ID)

	return resourceNatSnatRuleV2Read(d, meta)
}

func resourceNatSnatRuleV2Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	natClient, err := config.NatV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine nat client: %s", err)
	}

	snatRule, err := snatrules.Get(natClient, d.Id()).Extract()
	if err != nil {
		return CheckDeleted(d, err, "Snat Rule")
	}

	d.Set("region", GetRegion(d, config))
	d.Set("nat_gateway_id", snatRule.NatGatewayID)
	d.Set("floating_ip_id", snatRule.FloatingIPID)
	d.Set("floating_ip_address", snatRule.FloatingIPAddress)
	d.Set("subnet_id", snatRule.NetworkID)
	d.Set("cidr", snatRule.Cidr)
	d.Set("status", snatRule.Status)

	sourceType, _ := strconv.Atoi(snatRule.SourceType)
	d.Set("source_type", sourceType)

	return nil
}

func resourceNatSnatRuleV2Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	natClient, err := config.NatV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine nat client: %s", err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"ACTIVE"},
		Target:     []string{"DELETED"},
		Refresh:    waitForSnatRuleDelete(natClient, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Error deleting FlexibleEngine Snat Rule: %s", err)
	}

	d.SetId("")
	return nil
}

func waitForSnatRuleActive(natClient *golangsdk.ServiceClient, nId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		n, err := snatrules.Get(natClient, nId).Extract()
		if err != nil {
			return nil, "", err
		}

		log.Printf("[DEBUG] FlexibleEngine Snat Rule: %+v", n)
		if n.Status == "ACTIVE" {
			return n, "ACTIVE", nil
		}

		return n, "", nil
	}
}

func waitForSnatRuleDelete(natClient *golangsdk.ServiceClient, nId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		log.Printf("[DEBUG] Attempting to delete FlexibleEngine Snat Rule %s.\n", nId)

		n, err := snatrules.Get(natClient, nId).Extract()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				log.Printf("[DEBUG] Successfully deleted FlexibleEngine Snat Rule %s", nId)
				return n, "DELETED", nil
			}
			return n, "ACTIVE", err
		}

		err = snatrules.Delete(natClient, nId).ExtractErr()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				log.Printf("[DEBUG] Successfully deleted FlexibleEngine Snat Rule %s", nId)
				return n, "DELETED", nil
			}
			return n, "ACTIVE", err
		}

		log.Printf("[DEBUG] FlexibleEngine Snat Rule %s still active.\n", nId)
		return n, "ACTIVE", nil
	}
}
