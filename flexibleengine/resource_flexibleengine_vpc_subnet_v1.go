package flexibleengine

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/huaweicloud/golangsdk/openstack/common/tags"
	"github.com/huaweicloud/golangsdk/openstack/networking/v1/subnets"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/huaweicloud/golangsdk"
)

func resourceSubnetDNSListV1(d *schema.ResourceData) []string {
	rawDNSN := d.Get("dns_list").([]interface{})
	dnsn := make([]string, len(rawDNSN))
	for i, raw := range rawDNSN {
		dnsn[i] = raw.(string)
	}
	return dnsn
}

func resourceVpcSubnetV1() *schema.Resource {
	return &schema.Resource{
		Create: resourceVpcSubnetV1Create,
		Read:   resourceVpcSubnetV1Read,
		Update: resourceVpcSubnetV1Update,
		Delete: resourceVpcSubnetV1Delete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{ //request and response parameters
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     false,
				ValidateFunc: validateName,
			},
			"cidr": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateCIDR,
			},
			"dns_list": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validateIP,
				},
			},
			"gateway_ip": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateIP,
			},
			"dhcp_enable": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
				ForceNew: false,
			},
			"primary_dns": {
				Type:         schema.TypeString,
				ForceNew:     false,
				Optional:     true,
				ValidateFunc: validateIP,
				Computed:     true,
			},
			"secondary_dns": {
				Type:         schema.TypeString,
				ForceNew:     false,
				Optional:     true,
				ValidateFunc: validateIP,
				Computed:     true,
			},
			"availability_zone": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
				Computed: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"subnet_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": tagsSchema(),
		},
	}
}

func resourceVpcSubnetV1Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	subnetClient, err := config.networkingV1Client(GetRegion(d, config))

	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine networking client: %s", err)
	}

	createOpts := subnets.CreateOpts{
		Name:             d.Get("name").(string),
		CIDR:             d.Get("cidr").(string),
		AvailabilityZone: d.Get("availability_zone").(string),
		GatewayIP:        d.Get("gateway_ip").(string),
		EnableDHCP:       d.Get("dhcp_enable").(bool),
		VPC_ID:           d.Get("vpc_id").(string),
		PRIMARY_DNS:      d.Get("primary_dns").(string),
		SECONDARY_DNS:    d.Get("secondary_dns").(string),
		DnsList:          resourceSubnetDNSListV1(d),
	}

	n, err := subnets.Create(subnetClient, createOpts).Extract()
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine VPC subnet: %s", err)
	}

	d.SetId(n.ID)
	log.Printf("[INFO] Vpc Subnet ID: %s", n.ID)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"CREATING"},
		Target:     []string{"ACTIVE"},
		Refresh:    waitForVpcSubnetActive(subnetClient, n.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, stateErr := stateConf.WaitForState()
	if stateErr != nil {
		return fmt.Errorf(
			"Error waiting for Subnet (%s) to become ACTIVE: %s",
			n.ID, stateErr)
	}

	//set tags
	tagRaw := d.Get("tags").(map[string]interface{})
	if len(tagRaw) > 0 {
		vpcV2Client, err := config.networkingV2Client(GetRegion(d, config))
		if err != nil {
			return fmt.Errorf("Error creating FlexibleEngine VPC client: %s", err)
		}
		taglist := expandResourceTags(tagRaw)
		if tagErr := tags.Create(vpcV2Client, "subnets", n.ID, taglist).ExtractErr(); tagErr != nil {
			return fmt.Errorf("Error setting tags of VPC Subnet %s: %s", n.ID, tagErr)
		}
	}

	return resourceVpcSubnetV1Read(d, config)

}

func resourceVpcSubnetV1Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	subnetClient, err := config.networkingV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine networking client: %s", err)
	}

	n, err := subnets.Get(subnetClient, d.Id()).Extract()
	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error retrieving FlexibleEngine Subnets: %s", err)
	}

	d.Set("name", n.Name)
	d.Set("cidr", n.CIDR)
	d.Set("dns_list", n.DnsList)
	d.Set("gateway_ip", n.GatewayIP)
	d.Set("dhcp_enable", n.EnableDHCP)
	d.Set("primary_dns", n.PRIMARY_DNS)
	d.Set("secondary_dns", n.SECONDARY_DNS)
	d.Set("availability_zone", n.AvailabilityZone)
	d.Set("vpc_id", n.VPC_ID)
	d.Set("subnet_id", n.SubnetId)
	d.Set("region", GetRegion(d, config))

	// save VpcSubnet tags
	vpcV2Client, err := config.networkingV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine VPC client: %s", err)
	}
	resourceTags, err := tags.Get(vpcV2Client, "subnets", d.Id()).Extract()
	if err == nil {
		tagmap := tagsToMap(resourceTags.Tags)
		d.Set("tags", tagmap)
	} else {
		log.Printf("[WARN] fetching VPC Subnet %s tags failed: %s", d.Id(), err)
	}

	return nil
}

func resourceVpcSubnetV1Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	subnetClient, err := config.networkingV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine networking client: %s", err)
	}

	var updateOpts subnets.UpdateOpts

	//as name is mandatory while updating subnet
	updateOpts.Name = d.Get("name").(string)

	if d.HasChange("primary_dns") {
		updateOpts.PRIMARY_DNS = d.Get("primary_dns").(string)
	}
	if d.HasChange("secondary_dns") {
		updateOpts.SECONDARY_DNS = d.Get("secondary_dns").(string)
	}
	if d.HasChange("dns_list") {
		updateOpts.DnsList = resourceSubnetDNSListV1(d)
	}
	if d.HasChange("dhcp_enable") {
		updateOpts.EnableDHCP = d.Get("dhcp_enable").(bool)

	} else if d.Get("dhcp_enable").(bool) { //maintaining dhcp to be true if it was true earlier as default update option for dhcp bool is always going to be false in golangsdk
		updateOpts.EnableDHCP = true
	}

	vpc_id := d.Get("vpc_id").(string)

	_, err = subnets.Update(subnetClient, vpc_id, d.Id(), updateOpts).Extract()
	if err != nil {
		return fmt.Errorf("Error updating FlexibleEngine VPC Subnet: %s", err)
	}

	//update tags
	if d.HasChange("tags") {
		vpcV2Client, err := config.networkingV2Client(GetRegion(d, config))
		if err != nil {
			return fmt.Errorf("Error creating FlexibleEngine VPC client: %s", err)
		}

		tagErr := UpdateResourceTags(vpcV2Client, d, "subnets", d.Id())
		if tagErr != nil {
			return fmt.Errorf("Error updating tags of VPC subnet %s: %s", d.Id(), tagErr)
		}
	}
	return resourceVpcSubnetV1Read(d, meta)
}

func resourceVpcSubnetV1Delete(d *schema.ResourceData, meta interface{}) error {

	config := meta.(*Config)
	subnetClient, err := config.networkingV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine networking client: %s", err)
	}
	vpc_id := d.Get("vpc_id").(string)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"ACTIVE"},
		Target:     []string{"DELETED"},
		Refresh:    waitForVpcSubnetDelete(subnetClient, vpc_id, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Error deleting FlexibleEngine Subnet: %s", err)
	}

	d.SetId("")
	return nil
}

func waitForVpcSubnetActive(subnetClient *golangsdk.ServiceClient, vpcId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		n, err := subnets.Get(subnetClient, vpcId).Extract()
		if err != nil {
			return nil, "", err
		}

		if n.Status == "ACTIVE" {
			return n, "ACTIVE", nil
		}

		//If subnet status is other than Active, send error
		if n.Status == "DOWN" || n.Status == "ERROR" {
			return nil, "", fmt.Errorf("Subnet status: '%s'", n.Status)
		}

		return n, "CREATING", nil
	}
}

func waitForVpcSubnetDelete(subnetClient *golangsdk.ServiceClient, vpcId string, subnetId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {

		r, err := subnets.Get(subnetClient, subnetId).Extract()

		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				log.Printf("[INFO] Successfully deleted FlexibleEngine subnet %s", subnetId)
				return r, "DELETED", nil
			}
			return r, "ACTIVE", err
		}
		err = subnets.Delete(subnetClient, vpcId, subnetId).ExtractErr()

		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				log.Printf("[INFO] Successfully deleted FlexibleEngine subnet %s", subnetId)
				return r, "DELETED", nil
			}
			if errCode, ok := err.(golangsdk.ErrUnexpectedResponseCode); ok {
				if errCode.Actual == 409 {
					return r, "ACTIVE", nil
				}
			}
			return r, "ACTIVE", err
		}

		return r, "ACTIVE", nil
	}
}
