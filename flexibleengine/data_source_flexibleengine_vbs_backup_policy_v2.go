package flexibleengine

import (
	"fmt"
	"log"

	"github.com/chnsz/golangsdk/openstack/vbs/v2/policies"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceVBSBackupPolicyV2() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceVBSPolicyV2Read,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"frequency": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"remain_first_backup": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"rentention_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"policy_resource_count": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"start_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceVBSPolicyV2Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	vbsClient, err := config.VbsV2Client(GetRegion(d, config))

	listOpts := policies.ListOpts{
		ID:     d.Get("id").(string),
		Name:   d.Get("name").(string),
		Status: d.Get("status").(string),
	}

	refinedPolicies, err := policies.List(vbsClient, listOpts)
	if err != nil {
		return fmt.Errorf("Unable to retrieve policies: %s", err)
	}

	if len(refinedPolicies) < 1 {
		return fmt.Errorf("Your query returned no results. " +
			"Please change your search criteria and try again.")
	}

	if len(refinedPolicies) > 1 {
		return fmt.Errorf("Your query returned more than one result." +
			" Please try a more specific search criteria")
	}

	Policy := refinedPolicies[0]

	log.Printf("[INFO] Retrieved Policy using given filter %s: %+v", Policy.ID, Policy)
	d.SetId(Policy.ID)

	d.Set("name", Policy.Name)
	d.Set("policy_resource_count", Policy.ResourceCount)
	d.Set("frequency", Policy.ScheduledPolicy.Frequency)
	d.Set("remain_first_backup", Policy.ScheduledPolicy.RemainFirstBackup)
	d.Set("rentention_num", Policy.ScheduledPolicy.RententionNum)
	d.Set("start_time", Policy.ScheduledPolicy.StartTime)
	d.Set("status", Policy.ScheduledPolicy.Status)
	d.Set("region", GetRegion(d, config))

	return nil
}
