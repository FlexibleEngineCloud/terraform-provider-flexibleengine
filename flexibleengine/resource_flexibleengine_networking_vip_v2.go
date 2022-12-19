package flexibleengine

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/networking/v1/ports"
	"github.com/chnsz/golangsdk/openstack/networking/v1/subnets"
)

func resourceNetworkingVIPV2() *schema.Resource {
	return &schema.Resource{
		Create: resourceNetworkingVIPCreate,
		Read:   resourceNetworkingVIPRead,
		Update: resourceNetworkingVIPUpdate,
		Delete: resourceNetworkingVIPDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(2 * time.Minute),
			Delete: schema.DefaultTimeout(2 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"network_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"ip_version": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ValidateFunc: validation.IntInSlice([]int{4, 6}),
			},
			"ip_address": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			// Computed
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"mac_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"device_owner": {
				Type:     schema.TypeString,
				Computed: true,
			},

			// Deprecated
			"subnet_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				Computed:      true,
				ConflictsWith: []string{"ip_version"},
				Deprecated:    "use ip_version instead",
			},
		},
	}
}

func resourceNetworkingVIPCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client, err := config.NetworkingV1Client(config.GetRegion(d))
	if err != nil {
		return fmt.Errorf("Error creating VPC v1 client: %s", err)
	}

	networkId := d.Get("network_id").(string)
	n, err := subnets.Get(client, networkId).Extract()
	if err != nil {
		return fmt.Errorf("Error retrieving subnet by network ID (%s): %s", networkId, err)
	}

	// Check whether the subnet ID entered by the user belongs to the same subnet as the network ID.
	subnetId := d.Get("subnet_id").(string)
	if subnetId != "" && subnetId != n.SubnetId && subnetId != n.IPv6SubnetId {
		return fmt.Errorf("The subnet ID does not belong to the subnet where the network ID is located.")
	}

	// Pre-check for subnet network, the virtual IP of IPv6 must be established on the basis that the subnet supports
	// IPv6.
	if d.Get("ip_version").(int) == 6 {
		if n.IPv6SubnetId == "" {
			return fmt.Errorf("The subnet does not support IPv6, please enable IPv6 first.")
		}
		subnetId = n.IPv6SubnetId
	} else {
		subnetId = n.SubnetId
	}

	opts := ports.CreateOpts{
		Name:        d.Get("name").(string),
		DeviceOwner: "neutron:VIP_PORT",
		NetworkId:   networkId,
		FixedIps: []ports.FixedIp{
			{
				SubnetId:  subnetId,
				IpAddress: d.Get("ip_address").(string),
			},
		},
	}

	log.Printf("[DEBUG] Creating network VIP (%s) with options: %#v", d.Id(), opts)
	vip, err := ports.Create(client, opts)
	if err != nil {
		return fmt.Errorf("Error creating network VIP: %s", err)
	}

	d.SetId(vip.ID)
	log.Printf("[DEBUG] Waiting for network VIP (%s) to become available.", vip.ID)
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"BUILD"},
		Target:     []string{"DOWN", "ACTIVE"},
		Refresh:    waitForNetworkVIPActive(client, vip.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      3 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return err
	}

	return resourceNetworkingVIPRead(d, meta)
}

func resourceNetworkingVIPRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	region := config.GetRegion(d)
	client, err := config.NetworkingV1Client(region)
	if err != nil {
		return fmt.Errorf("Error creating VPC v1 client: %s", err)
	}

	vip, err := ports.Get(client, d.Id())
	if err != nil {
		return CheckDeleted(d, err, "VPC network VIP")
	}

	log.Printf("[DEBUG] Retrieved VIP %s: %+v", d.Id(), vip)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("name", vip.Name),
		d.Set("status", normalizeNetworkVipStatus(vip.Status)),
		d.Set("device_owner", vip.DeviceOwner),
		d.Set("mac_address", vip.MacAddress),
		d.Set("network_id", vip.NetworkId),
		setVipNetworkParams(d, vip),
	)

	return mErr.ErrorOrNil()
}

func setVipNetworkParams(d *schema.ResourceData, port *ports.Port) error {
	if len(port.FixedIps) > 0 {
		addr := port.FixedIps[0].IpAddress
		var ipVersion int
		if isIPv4Address(addr) {
			ipVersion = 4
		} else {
			ipVersion = 6
		}
		mErr := multierror.Append(nil,
			d.Set("ip_version", ipVersion),
			d.Set("ip_address", addr),
			d.Set("subnet_id", port.FixedIps[0].SubnetId),
		)
		return mErr.ErrorOrNil()
	}
	return nil
}

func resourceNetworkingVIPUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client, err := config.NetworkingV1Client(config.GetRegion(d))
	if err != nil {
		return fmt.Errorf("Error creating VPC v1 client: %s", err)
	}

	opts := ports.UpdateOpts{
		Name: d.Get("name").(string),
	}

	log.Printf("[DEBUG] Updating network VIP (%s) with options: %#v", d.Id(), opts)
	_, err = ports.Update(client, d.Id(), opts)
	if err != nil {
		return fmt.Errorf("Error updating networking VIP: %s", err)
	}

	return resourceNetworkingVIPRead(d, meta)
}

func resourceNetworkingVIPDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client, err := config.NetworkingV1Client(config.GetRegion(d))
	if err != nil {
		return fmt.Errorf("Error creating VPC v1 client: %s", err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"ACTIVE"},
		Target:     []string{"DELETED"},
		Refresh:    waitForNetworkVIPDelete(client, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Error deleting FlexibleEngine networking VIP: %s", err)
	}

	d.SetId("")
	return nil
}

func waitForNetworkVIPActive(vpcClient *golangsdk.ServiceClient, vipid string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		p, err := ports.Get(vpcClient, vipid)
		if err != nil {
			return nil, "", err
		}

		log.Printf("[DEBUG] FlexibleEngine networking VIP port: %+v", p)
		if p.Status == "DOWN" || p.Status == "ACTIVE" {
			return p, "ACTIVE", nil
		}

		return p, p.Status, nil
	}
}

func waitForNetworkVIPDelete(vpcClient *golangsdk.ServiceClient, vipid string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		log.Printf("[DEBUG] Attempting to delete FlexibleEngine networking VIP port %s", vipid)

		p, err := ports.Get(vpcClient, vipid)
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				log.Printf("[DEBUG] Successfully deleted FlexibleEngine VIP %s", vipid)
				return p, "DELETED", nil
			}
			return p, "ACTIVE", err
		}

		err = ports.Delete(vpcClient, vipid).ExtractErr()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				log.Printf("[DEBUG] Successfully deleted FlexibleEngine VIP %s", vipid)
				return p, "DELETED", nil
			}
			return p, "ACTIVE", err
		}

		log.Printf("[DEBUG] FlexibleEngine VIP %s still active.\n", vipid)
		return p, "ACTIVE", nil
	}
}

// isIPv4Address is used to check whether the addr string is IPv4 format
func isIPv4Address(addr string) bool {
	pattern := "^((25[0-5]|2[0-4]\\d|(1\\d{2}|[1-9]?\\d))\\.){3}(25[0-5]|2[0-4]\\d|(1\\d{2}|[1-9]?\\d))$"
	matched, _ := regexp.MatchString(pattern, addr)
	return matched
}

// For VIP ports, the status will always be 'DOWN'.
func normalizeNetworkVipStatus(status string) string {
	if status == "DOWN" {
		return "ACTIVE"
	}
	return status
}
