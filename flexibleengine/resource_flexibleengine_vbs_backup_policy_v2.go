package flexibleengine

import (
	"fmt"
	"log"
	"time"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/vbs/v2/policies"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceVBSBackupPolicyV2() *schema.Resource {
	return &schema.Resource{
		Create: resourceVBSBackupPolicyV2Create,
		Read:   resourceVBSBackupPolicyV2Read,
		Update: resourceVBSBackupPolicyV2Update,
		Delete: resourceVBSBackupPolicyV2Delete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateVBSPolicyName,
			},

			"resources": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"start_time": {
				Type:     schema.TypeString,
				Required: true,
			},
			"frequency": {
				Type:          schema.TypeInt,
				Optional:      true,
				ConflictsWith: []string{"week_frequency"},
				ValidateFunc:  validateVBSPolicyFrequency,
			},
			"week_frequency": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 7,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"rentention_num": {
				Type:          schema.TypeInt,
				Optional:      true,
				ConflictsWith: []string{"rentention_day"},
				ValidateFunc:  validateVBSPolicyRetentionNum,
			},
			"rentention_day": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validateVBSPolicyRetentionNum,
			},
			"retain_first_backup": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateVBSPolicyRetainBackup,
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "ON",
				ValidateFunc: validateVBSPolicyStatus,
			},
			"policy_resource_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func resourceVBSBackupPolicyV2Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	vbsClient, err := config.VbsV2Client(GetRegion(d, config))

	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine VBS client: %s", err)
	}

	_, isExist1 := d.GetOk("frequency")
	_, isExist2 := d.GetOk("week_frequency")
	if !isExist1 && !isExist2 {
		return fmt.Errorf("either frequency or week_frequency must be specified")
	}

	_, isExist1 = d.GetOk("rentention_num")
	_, isExist2 = d.GetOk("rentention_day")
	if !isExist1 && !isExist2 {
		return fmt.Errorf("either rentention_num or rentention_day must be specified")
	}

	weeks, err := buildWeekFrequencyResource(d)
	if err != nil {
		return err
	}

	createOpts := policies.CreateOpts{
		Name: d.Get("name").(string),
		ScheduledPolicy: policies.ScheduledPolicy{
			StartTime:         d.Get("start_time").(string),
			Frequency:         d.Get("frequency").(int),
			WeekFrequency:     weeks,
			RententionNum:     d.Get("rentention_num").(int),
			RententionDay:     d.Get("rentention_day").(int),
			RemainFirstBackup: d.Get("retain_first_backup").(string),
			Status:            d.Get("status").(string),
		},
	}

	create, err := policies.Create(vbsClient, createOpts).Extract()

	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine Backup Policy: %s", err)
	}
	d.SetId(create.ID)

	// associate volumes to backup policy
	resources := buildAssociateResource(d.Get("resources").([]interface{}))
	if len(resources) > 0 {
		opts := policies.AssociateOpts{
			PolicyID:  d.Id(),
			Resources: resources,
		}

		_, err := policies.Associate(vbsClient, opts).ExtractResource()
		if err != nil {
			return fmt.Errorf("Error associate volumes to VBS backup policy %s: %s",
				d.Id(), err)
		}
	}

	return resourceVBSBackupPolicyV2Read(d, meta)

}

func resourceVBSBackupPolicyV2Read(d *schema.ResourceData, meta interface{}) error {

	config := meta.(*Config)
	vbsClient, err := config.VbsV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine VBS client: %s", err)
	}

	PolicyOpts := policies.ListOpts{ID: d.Id()}
	policies, err := policies.List(vbsClient, PolicyOpts)
	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error retrieving FlexibleEngine Backup Policy: %s", err)
	}

	n := policies[0]

	d.Set("name", n.Name)
	d.Set("start_time", n.ScheduledPolicy.StartTime)
	d.Set("frequency", n.ScheduledPolicy.Frequency)
	d.Set("week_frequency", n.ScheduledPolicy.WeekFrequency)
	d.Set("rentention_num", n.ScheduledPolicy.RententionNum)
	d.Set("rentention_day", n.ScheduledPolicy.RententionDay)
	d.Set("retain_first_backup", n.ScheduledPolicy.RemainFirstBackup)
	d.Set("status", n.ScheduledPolicy.Status)
	d.Set("policy_resource_count", n.ResourceCount)

	return nil
}

func resourceVBSBackupPolicyV2Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	vbsClient, err := config.VbsV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error updating FlexibleEngine VBS client: %s", err)
	}

	_, isExist1 := d.GetOk("frequency")
	_, isExist2 := d.GetOk("week_frequency")
	if !isExist1 && !isExist2 {
		return fmt.Errorf("either frequency or week_frequency must be specified")
	}

	_, isExist1 = d.GetOk("rentention_num")
	_, isExist2 = d.GetOk("rentention_day")
	if !isExist1 && !isExist2 {
		return fmt.Errorf("either rentention_num or rentention_day must be specified")
	}

	frequency := d.Get("frequency").(int)
	weeks, err := buildWeekFrequencyResource(d)
	if err != nil {
		return err
	}

	var updateOpts policies.UpdateOpts
	if frequency != 0 {
		updateOpts.ScheduledPolicy.Frequency = frequency
	} else {
		updateOpts.ScheduledPolicy.WeekFrequency = weeks
	}

	if d.HasChange("name") || d.HasChange("start_time") || d.HasChange("retain_first_backup") ||
		d.HasChange("rentention_num") || d.HasChange("rentention_day") || d.HasChange("status") ||
		d.HasChange("frequency") || d.HasChange("week_frequency") {
		if d.HasChange("name") {
			updateOpts.Name = d.Get("name").(string)
		}
		if d.HasChange("start_time") {
			updateOpts.ScheduledPolicy.StartTime = d.Get("start_time").(string)
		}
		if d.HasChange("rentention_num") {
			updateOpts.ScheduledPolicy.RententionNum = d.Get("rentention_num").(int)
		}
		if d.HasChange("rentention_day") {
			updateOpts.ScheduledPolicy.RententionDay = d.Get("rentention_day").(int)
		}
		if d.HasChange("retain_first_backup") {
			updateOpts.ScheduledPolicy.RemainFirstBackup = d.Get("retain_first_backup").(string)
		}
		if d.HasChange("status") {
			updateOpts.ScheduledPolicy.Status = d.Get("status").(string)
		}

		_, err = policies.Update(vbsClient, d.Id(), updateOpts).Extract()
		if err != nil {
			return fmt.Errorf("Error updating FlexibleEngine backup policy: %s", err)
		}
	}

	if d.HasChange("resources") {
		old, new := d.GetChange("resources")

		// disassociate old volumes from backup policy
		removeResources := buildDisassociateResource(old.([]interface{}))
		if len(removeResources) > 0 {
			opts := policies.DisassociateOpts{
				Resources: removeResources,
			}

			_, err := policies.Disassociate(vbsClient, d.Id(), opts).ExtractResource()
			if err != nil {
				return fmt.Errorf("Error disassociate volumes from VBS backup policy %s: %s",
					d.Id(), err)
			}
		}

		// associate new volumes to backup policy
		addResources := buildAssociateResource(new.([]interface{}))
		if len(addResources) > 0 {
			opts := policies.AssociateOpts{
				PolicyID:  d.Id(),
				Resources: addResources,
			}

			_, err := policies.Associate(vbsClient, opts).ExtractResource()
			if err != nil {
				return fmt.Errorf("Error associate volumes to VBS backup policy %s: %s",
					d.Id(), err)
			}
		}
	}

	return resourceVBSBackupPolicyV2Read(d, meta)
}

func resourceVBSBackupPolicyV2Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	vbsClient, err := config.VbsV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine VBS client: %s", err)
	}

	delete := policies.Delete(vbsClient, d.Id())
	if delete.Err != nil {
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			log.Printf("[INFO] Successfully deleted FlexibleEngine VBS Backup Policy %s", d.Id())

		}
		if errCode, ok := err.(golangsdk.ErrUnexpectedResponseCode); ok {
			if errCode.Actual == 409 {
				log.Printf("[INFO] Error deleting FlexibleEngine VBS Backup Policy %s", d.Id())
			}
		}
		log.Printf("[INFO] Successfully deleted FlexibleEngine VBS Backup Policy %s", d.Id())
	}

	d.SetId("")
	return nil
}

func buildAssociateResource(raw []interface{}) []policies.AssociateResource {
	resources := make([]policies.AssociateResource, len(raw))
	for i, v := range raw {
		resources[i] = policies.AssociateResource{
			ResourceID:   v.(string),
			ResourceType: "volume",
		}
	}
	return resources
}

func buildDisassociateResource(raw []interface{}) []policies.DisassociateResource {
	resources := make([]policies.DisassociateResource, len(raw))
	for i, v := range raw {
		resources[i] = policies.DisassociateResource{
			ResourceID: v.(string),
		}
	}
	return resources
}

func buildWeekFrequencyResource(d *schema.ResourceData) ([]string, error) {
	validateList := []string{"SUN", "MON", "TUE", "WED", "THU", "FRI", "SAT"}
	weeks := []string{}

	weekRaws := d.Get("week_frequency").([]interface{})
	for _, wf := range weekRaws {
		found := false
		for _, value := range validateList {
			if wf.(string) == value {
				found = true
				break
			}
		}

		if found {
			weeks = append(weeks, wf.(string))
		} else {
			return nil, fmt.Errorf("expected item of week_frequency to be one of %v, got %s",
				validateList, wf.(string))
		}
	}
	return weeks, nil
}
