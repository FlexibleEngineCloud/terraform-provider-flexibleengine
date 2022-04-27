package flexibleengine

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/chnsz/golangsdk/openstack/common/tags"
	"github.com/chnsz/golangsdk/openstack/networking/v1/vpcs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// resourceVirtualPrivateCloudV1 is not used any more since v1.29.0
// flexibleengine_vpc_v1 can be imported directly
func resourceVirtualPrivateCloudV1() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVirtualPrivateCloudV1Create,
		ReadContext:   resourceVirtualPrivateCloudV1Read,
		UpdateContext: resourceVirtualPrivateCloudV1Update,
		DeleteContext: resourceVirtualPrivateCloudV1Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(3 * time.Minute),
		},

		Schema: map[string]*schema.Schema{ //request and response parameters
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
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
				ForceNew:     false,
				ValidateFunc: validateCIDR,
			},
			"tags": tagsSchema(),
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"shared": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func resourceVirtualPrivateCloudV1Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	vpcClient, err := config.NetworkingV1Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating FlexibleEngine vpc client: %s", err)
	}

	createOpts := vpcs.CreateOpts{
		Name: d.Get("name").(string),
		CIDR: d.Get("cidr").(string),
	}

	n, err := vpcs.Create(vpcClient, createOpts).Extract()
	if err != nil {
		return diag.Errorf("Error creating FlexibleEngine VPC: %s", err)
	}
	d.SetId(n.ID)

	log.Printf("[INFO] Vpc ID: %s", n.ID)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"CREATING"},
		Target:     []string{"ACTIVE"},
		Refresh:    waitForVpcActive(vpcClient, n.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, stateErr := stateConf.WaitForState()
	if stateErr != nil {
		return diag.Errorf(
			"Error waiting for Vpc (%s) to become ACTIVE: %s",
			n.ID, stateErr)
	}

	//set tags
	tagRaw := d.Get("tags").(map[string]interface{})
	if len(tagRaw) > 0 {
		vpcV2Client, err := config.NetworkingV2Client(GetRegion(d, config))
		if err != nil {
			return diag.Errorf("Error creating FlexibleEngine vpc client: %s", err)
		}
		taglist := expandResourceTags(tagRaw)
		if tagErr := tags.Create(vpcV2Client, "vpcs", n.ID, taglist).ExtractErr(); tagErr != nil {
			return diag.Errorf("Error setting tags of VPC %s: %s", n.ID, tagErr)
		}
	}

	return resourceVirtualPrivateCloudV1Read(ctx, d, meta)

}

func resourceVirtualPrivateCloudV1Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	vpcClient, err := config.NetworkingV1Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating FlexibleEngine Vpc client: %s", err)
	}

	n, err := vpcs.Get(vpcClient, d.Id()).Extract()
	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			d.SetId("")
			return nil
		}

		return diag.Errorf("Error retrieving FlexibleEngine Vpc: %s", err)
	}

	d.SetId(n.ID)
	d.Set("name", n.Name)
	d.Set("cidr", n.CIDR)
	d.Set("status", n.Status)
	d.Set("shared", n.EnableSharedSnat)
	d.Set("region", GetRegion(d, config))

	// save tags
	vpcV2Client, err := config.NetworkingV2Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating FlexibleEngine vpc client: %s", err)
	}
	resourceTags, err := tags.Get(vpcV2Client, "vpcs", d.Id()).Extract()
	if err == nil {
		tagmap := tagsToMap(resourceTags.Tags)
		d.Set("tags", tagmap)
	} else {
		log.Printf("[WARN] fetching VPC %s tags failed: %s", d.Id(), err)
	}

	return nil
}

func resourceVirtualPrivateCloudV1Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	vpcClient, err := config.NetworkingV1Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating FlexibleEngine Vpc: %s", err)
	}

	if d.HasChanges("name", "cidr") {
		updateOpts := vpcs.UpdateOpts{
			Name: d.Get("name").(string),
			CIDR: d.Get("cidr").(string),
		}

		_, err = vpcs.Update(vpcClient, d.Id(), updateOpts).Extract()
		if err != nil {
			return diag.Errorf("Error updating FlexibleEngine Vpc: %s", err)
		}
	}

	//update tags
	if d.HasChange("tags") {
		vpcV2Client, err := config.NetworkingV2Client(GetRegion(d, config))
		if err != nil {
			return diag.Errorf("Error creating FlexibleEngine vpc client: %s", err)
		}

		tagErr := UpdateResourceTags(vpcV2Client, d, "vpcs", d.Id())
		if tagErr != nil {
			return diag.Errorf("Error updating tags of VPC %s: %s", d.Id(), tagErr)
		}
	}

	return resourceVirtualPrivateCloudV1Read(ctx, d, meta)
}

func resourceVirtualPrivateCloudV1Delete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	vpcClient, err := config.NetworkingV1Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating FlexibleEngine vpc: %s", err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"ACTIVE"},
		Target:     []string{"DELETED"},
		Refresh:    waitForVpcDelete(vpcClient, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return diag.Errorf("Error deleting FlexibleEngine Vpc: %s", err)
	}

	d.SetId("")
	return nil
}

func waitForVpcActive(vpcClient *golangsdk.ServiceClient, vpcId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		n, err := vpcs.Get(vpcClient, vpcId).Extract()
		if err != nil {
			return nil, "", err
		}

		if n.Status == "OK" {
			return n, "ACTIVE", nil
		}

		//If vpc status is other than Ok, send error
		if n.Status == "DOWN" {
			return nil, "", fmt.Errorf("Vpc status: '%s'", n.Status)
		}

		return n, n.Status, nil
	}
}

func waitForVpcDelete(vpcClient *golangsdk.ServiceClient, vpcId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {

		r, err := vpcs.Get(vpcClient, vpcId).Extract()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				log.Printf("[INFO] Successfully deleted FlexibleEngine vpc %s", vpcId)
				return r, "DELETED", nil
			}
			return r, "ACTIVE", err
		}

		err = vpcs.Delete(vpcClient, vpcId).ExtractErr()

		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				log.Printf("[INFO] Successfully deleted FlexibleEngine vpc %s", vpcId)
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
