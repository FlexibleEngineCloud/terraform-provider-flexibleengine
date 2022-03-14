package flexibleengine

import (
	"fmt"
	"log"
	"time"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/common/tags"
	"github.com/chnsz/golangsdk/openstack/networking/v1/bandwidths"
	"github.com/chnsz/golangsdk/openstack/networking/v1/eips"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceVpcEIPV1() *schema.Resource {
	return &schema.Resource{
		Create: resourceVpcEIPV1Create,
		Read:   resourceVpcEIPV1Read,
		Update: resourceVpcEIPV1Update,
		Delete: resourceVpcEIPV1Delete,
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
			"publicip": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"ip_address": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
							Computed: true,
						},
						"port_id": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"bandwidth": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: false,
						},
						"size": {
							Type:     schema.TypeInt,
							Required: true,
							ForceNew: false,
						},
						"share_type": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"charge_mode": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
							Computed: true,
						},
					},
				},
			},
			"tags": tagsSchema(),

			"address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceVpcEIPV1Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	networkingClient, err := config.NetworkingV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating networking client: %s", err)
	}

	createOpts := eips.ApplyOpts{
		IP:        resourcePublicIP(d),
		Bandwidth: resourceBandWidth(d),
	}

	log.Printf("[DEBUG] Create Options: %#v", createOpts)
	eIP, err := eips.Apply(networkingClient, createOpts).Extract()
	if err != nil {
		return fmt.Errorf("Error allocating EIP: %s", err)
	}

	d.SetId(eIP.ID)
	log.Printf("[DEBUG] Waiting for EIP %#v to become available.", eIP)

	timeout := d.Timeout(schema.TimeoutCreate)
	err = waitForEIPActive(networkingClient, eIP.ID, timeout)
	if err != nil {
		return fmt.Errorf(
			"Error waiting for EIP (%s) to become ready: %s",
			eIP.ID, err)
	}

	err = bindToPort(d, eIP.ID, networkingClient, timeout)
	if err != nil {
		return fmt.Errorf("Error binding eip:%s to port: %s", eIP.ID, err)
	}

	//set tags
	tagRaw := d.Get("tags").(map[string]interface{})
	if len(tagRaw) > 0 {
		vpcV2Client, err := config.NetworkingV2Client(GetRegion(d, config))
		if err != nil {
			return fmt.Errorf("Error creating FlexibleEngine vpc client: %s", err)
		}
		taglist := expandResourceTags(tagRaw)
		if tagErr := tags.Create(vpcV2Client, "publicips", eIP.ID, taglist).ExtractErr(); tagErr != nil {
			return fmt.Errorf("Error setting tags of EIP %s: %s", eIP.ID, tagErr)
		}
	}

	return resourceVpcEIPV1Read(d, meta)
}

func resourceVpcEIPV1Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	networkingClient, err := config.NetworkingV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating networking client: %s", err)
	}

	eIP, err := eips.Get(networkingClient, d.Id()).Extract()
	if err != nil {
		return CheckDeleted(d, err, "eIP")
	}
	bandWidth, err := bandwidths.Get(networkingClient, eIP.BandwidthID).Extract()
	if err != nil {
		return fmt.Errorf("Error fetching bandwidth: %s", err)
	}

	// Set public ip
	publicIP := []map[string]string{
		{
			"type":       eIP.Type,
			"ip_address": eIP.PublicAddress,
			"port_id":    eIP.PortID,
		},
	}
	d.Set("publicip", publicIP)

	// Set bandwidth
	bW := []map[string]interface{}{
		{
			"name":        bandWidth.Name,
			"size":        eIP.BandwidthSize,
			"share_type":  eIP.BandwidthShareType,
			"charge_mode": bandWidth.ChargeMode,
		},
	}
	d.Set("bandwidth", bW)
	d.Set("region", GetRegion(d, config))
	d.Set("address", eIP.PublicAddress)
	d.Set("status", normalizeEIPStatus(eIP.Status))

	// save tags
	vpcV2Client, err := config.NetworkingV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine vpc client: %s", err)
	}
	resourceTags, err := tags.Get(vpcV2Client, "publicips", d.Id()).Extract()
	if err == nil {
		tagmap := tagsToMap(resourceTags.Tags)
		d.Set("tags", tagmap)
	} else {
		log.Printf("[WARN] fetching EIP %s tags failed: %s", d.Id(), err)
	}

	return nil
}

func resourceVpcEIPV1Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	networkingClient, err := config.NetworkingV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating networking client: %s", err)
	}

	// Update bandwidth change
	if d.HasChange("bandwidth") {
		var updateOpts bandwidths.UpdateOpts

		newBWList := d.Get("bandwidth").([]interface{})
		newMap := newBWList[0].(map[string]interface{})
		updateOpts.Size = newMap["size"].(int)
		updateOpts.Name = newMap["name"].(string)

		log.Printf("[DEBUG] Bandwidth Update Options: %#v", updateOpts)

		eIP, err := eips.Get(networkingClient, d.Id()).Extract()
		if err != nil {
			return CheckDeleted(d, err, "eIP")
		}
		_, err = bandwidths.Update(networkingClient, eIP.BandwidthID, updateOpts).Extract()
		if err != nil {
			return fmt.Errorf("Error updating bandwidth: %s", err)
		}

	}

	// Update publicip change
	if d.HasChange("publicip") {
		var updateOpts eips.UpdateOpts

		newIPList := d.Get("publicip").([]interface{})
		newMap := newIPList[0].(map[string]interface{})
		updateOpts.PortID = newMap["port_id"].(string)

		log.Printf("[DEBUG] PublicIP Update Options: %#v", updateOpts)
		_, err = eips.Update(networkingClient, d.Id(), updateOpts).Extract()
		if err != nil {
			return fmt.Errorf("Error updating publicip: %s", err)
		}

	}

	//update tags
	if d.HasChange("tags") {
		vpcV2Client, err := config.NetworkingV2Client(GetRegion(d, config))
		if err != nil {
			return fmt.Errorf("Error creating FlexibleEngine vpc client: %s", err)
		}

		tagErr := UpdateResourceTags(vpcV2Client, d, "publicips", d.Id())
		if tagErr != nil {
			return fmt.Errorf("Error updating tags of EIP %s: %s", d.Id(), tagErr)
		}
	}

	return resourceVpcEIPV1Read(d, meta)
}

func resourceVpcEIPV1Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	networkingClient, err := config.NetworkingV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating VPC client: %s", err)
	}

	timeout := d.Timeout(schema.TimeoutDelete)
	err = unbindToPort(d, d.Id(), networkingClient, timeout)
	if err != nil {
		return fmt.Errorf("Error unbinding eip:%s to port: %s", d.Id(), err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"ACTIVE"},
		Target:     []string{"DELETED"},
		Refresh:    waitForEIPDelete(networkingClient, d.Id()),
		Timeout:    timeout,
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Error deleting EIP: %s", err)
	}

	d.SetId("")

	return nil
}

func resourcePublicIP(d *schema.ResourceData) eips.PublicIpOpts {
	publicIPRaw := d.Get("publicip").([]interface{})
	rawMap := publicIPRaw[0].(map[string]interface{})

	publicip := eips.PublicIpOpts{
		Type:    rawMap["type"].(string),
		Address: rawMap["ip_address"].(string),
	}
	return publicip
}

func resourceBandWidth(d *schema.ResourceData) eips.BandwidthOpts {
	bandwidthRaw := d.Get("bandwidth").([]interface{})
	rawMap := bandwidthRaw[0].(map[string]interface{})

	bandwidth := eips.BandwidthOpts{
		Name:       rawMap["name"].(string),
		Size:       rawMap["size"].(int),
		ShareType:  rawMap["share_type"].(string),
		ChargeMode: rawMap["charge_mode"].(string),
	}
	return bandwidth
}

func bindToPort(d *schema.ResourceData, eipID string, networkingClient *golangsdk.ServiceClient, timeout time.Duration) error {
	publicIPRaw := d.Get("publicip").([]interface{})
	rawMap := publicIPRaw[0].(map[string]interface{})
	portID, ok := rawMap["port_id"]
	if !ok || portID == "" {
		return nil
	}

	pd := portID.(string)
	log.Printf("[DEBUG] Bind eip:%s to port: %s", eipID, pd)

	updateOpts := eips.UpdateOpts{PortID: pd}
	_, err := eips.Update(networkingClient, eipID, updateOpts).Extract()
	if err != nil {
		return err
	}
	return waitForEIPActive(networkingClient, eipID, timeout)
}

func unbindToPort(d *schema.ResourceData, eipID string, networkingClient *golangsdk.ServiceClient, timeout time.Duration) error {
	publicIPRaw := d.Get("publicip").([]interface{})
	rawMap := publicIPRaw[0].(map[string]interface{})
	portID, ok := rawMap["port_id"]
	if !ok || portID == "" {
		return nil
	}

	pd := portID.(string)
	log.Printf("[DEBUG] Unbind eip:%s to port: %s", eipID, pd)

	updateOpts := eips.UpdateOpts{PortID: ""}
	_, err := eips.Update(networkingClient, eipID, updateOpts).Extract()
	if err != nil {
		return err
	}
	return waitForEIPActive(networkingClient, eipID, timeout)
}

func getEIPStatus(networkingClient *golangsdk.ServiceClient, eId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		e, err := eips.Get(networkingClient, eId).Extract()
		if err != nil {
			return nil, "", err
		}

		log.Printf("[DEBUG] EIP: %+v", e)
		if e.Status == "DOWN" || e.Status == "ACTIVE" {
			return e, "ACTIVE", nil
		}

		return e, "", nil
	}
}

func waitForEIPActive(networkingClient *golangsdk.ServiceClient, eipID string, timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Target:     []string{"ACTIVE"},
		Refresh:    getEIPStatus(networkingClient, eipID),
		Timeout:    timeout,
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err := stateConf.WaitForState()
	return err
}

func normalizeEIPStatus(status string) string {
	var ret string = status

	// "DOWN" means the eip is active but unbound
	if status == "DOWN" {
		ret = "UNBOUND"
	} else if status == "ACTIVE" {
		ret = "BOUND"
	}

	return ret
}

func waitForEIPDelete(networkingClient *golangsdk.ServiceClient, eId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		log.Printf("[DEBUG] Attempting to delete EIP %s.\n", eId)

		e, err := eips.Get(networkingClient, eId).Extract()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				log.Printf("[DEBUG] Successfully deleted EIP %s", eId)
				return e, "DELETED", nil
			}
			return e, "ACTIVE", err
		}

		err = eips.Delete(networkingClient, eId).ExtractErr()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				log.Printf("[DEBUG] Successfully deleted EIP %s", eId)
				return e, "DELETED", nil
			}
			return e, "ACTIVE", err
		}

		log.Printf("[DEBUG] EIP %s still active.\n", eId)
		return e, "ACTIVE", nil
	}
}
