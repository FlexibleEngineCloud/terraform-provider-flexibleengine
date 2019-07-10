package flexibleengine

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/openstack/csbs/v1/policies"
)

func resourceCSBSBackupPolicyV1() *schema.Resource {
	return &schema.Resource{
		Create: resourceCSBSBackupPolicyCreate,
		Read:   resourceCSBSBackupPolicyRead,
		Update: resourceCSBSBackupPolicyUpdate,
		Delete: resourceCSBSBackupPolicyDelete,
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
				ForceNew: true,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"provider_id": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "fc4d5750-22e7-4798-8a46-f48f62c4c1da",
				ForceNew: true,
			},
			"common": {
				Type:     schema.TypeMap,
				Optional: true,
			},
			"scheduled_operation": {
				Type:     schema.TypeSet,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "default",
						},
						"description": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"enabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
						},
						"max_backups": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"retention_duration_days": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"permanent": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"trigger_pattern": {
							Type:     schema.TypeString,
							Required: true,
						},
						"operation_type": {
							Type:     schema.TypeString,
							Required: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"trigger_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"trigger_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"trigger_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"resource": {
				Type:     schema.TypeSet,
				Required: true,
				ForceNew: false,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"type": {
							Type:     schema.TypeString,
							Required: true,
						},
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
		},
	}

}

func resourceCSBSBackupPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	policyClient, err := config.csbsV1Client(GetRegion(d, config))

	if err != nil {
		return fmt.Errorf("Error creating backup policy Client: %s", err)
	}

	createOpts := policies.CreateOpts{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		ProviderId:  d.Get("provider_id").(string),
		Parameters: policies.PolicyParam{
			Common: resourceCSBSCommonParamsV1(d),
		},
		ScheduledOperations: resourceCSBSScheduleV1(d),

		Resources: resourceCSBSResourceV1(d),
	}

	backupPolicy, err := policies.Create(policyClient, createOpts).Extract()

	if err != nil {
		return fmt.Errorf("Error creating Backup Policy : %s", err)
	}

	d.SetId(backupPolicy.ID)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"creating"},
		Target:     []string{"suspended"},
		Refresh:    waitForCSBSBackupPolicyActive(policyClient, backupPolicy.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, StateErr := stateConf.WaitForState()
	if StateErr != nil {
		return fmt.Errorf("Error waiting for Backup Policy (%s) to become available: %s", backupPolicy.ID, StateErr)
	}

	return resourceCSBSBackupPolicyRead(d, meta)

}

func resourceCSBSBackupPolicyRead(d *schema.ResourceData, meta interface{}) error {

	config := meta.(*Config)
	policyClient, err := config.csbsV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating csbs client: %s", err)
	}

	backupPolicy, err := policies.Get(policyClient, d.Id()).Extract()
	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			log.Printf("[WARN] Removing backup policy %s as it's already gone", d.Id())
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error retrieving backup policy: %s", err)
	}

	if err := d.Set("resource", flattenCSBSPolicyResources(*backupPolicy)); err != nil {
		return err
	}

	if err := d.Set("scheduled_operation", flattenCSBSScheduledOperations(*backupPolicy)); err != nil {
		return err
	}

	d.Set("name", backupPolicy.Name)
	d.Set("common", backupPolicy.Parameters.Common)
	d.Set("status", backupPolicy.Status)
	d.Set("description", backupPolicy.Description)
	d.Set("provider_id", backupPolicy.ProviderId)

	d.Set("region", GetRegion(d, config))

	return nil
}

func resourceCSBSBackupPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	policyClient, err := config.csbsV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating csbs client: %s", err)
	}
	var updateOpts policies.UpdateOpts
	if d.HasChange("name") {
		updateOpts.Name = d.Get("name").(string)
	}
	if d.HasChange("description") {
		updateOpts.Description = d.Get("description").(string)
	}

	updateOpts.Parameters.Common = resourceCSBSCommonParamsV1(d)

	if d.HasChange("resource") {
		updateOpts.Resources = resourceCSBSResourceV1(d)
	}
	if d.HasChange("scheduled_operation") {
		updateOpts.ScheduledOperations = resourceCSBScheduleUpdateV1(d)
	}

	_, err = policies.Update(policyClient, d.Id(), updateOpts).Extract()
	if err != nil {
		return fmt.Errorf("Error updating Backup Policy: %s", err)
	}

	return resourceCSBSBackupPolicyRead(d, meta)
}

func resourceCSBSBackupPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	policyClient, err := config.csbsV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating csbs client: %s", err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"available"},
		Target:     []string{"deleted"},
		Refresh:    waitForVBSPolicyDelete(policyClient, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Error deleting Backup Policy: %s", err)
	}

	d.SetId("")
	return nil
}

func waitForCSBSBackupPolicyActive(policyClient *golangsdk.ServiceClient, policyID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		n, err := policies.Get(policyClient, policyID).Extract()
		if err != nil {
			return nil, "", err
		}

		if n.Status == "error" {
			return n, n.Status, nil
		}
		return n, n.Status, nil
	}
}

func waitForVBSPolicyDelete(policyClient *golangsdk.ServiceClient, policyID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {

		r, err := policies.Get(policyClient, policyID).Extract()

		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				log.Printf("[INFO] Successfully deleted Backup Policy %s", policyID)
				return r, "deleted", nil
			}
			return r, "available", err
		}

		policy := policies.Delete(policyClient, policyID)
		err = policy.Err
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				log.Printf("[INFO] Successfully deleted Backup Policy %s", policyID)
				return r, "deleted", nil
			}
			if errCode, ok := err.(golangsdk.ErrUnexpectedResponseCode); ok {
				if errCode.Actual == 409 {
					return r, "available", nil
				}
			}
			return r, "available", err
		}

		return r, "deleted", nil
	}
}

func resourceCSBSScheduleV1(d *schema.ResourceData) []policies.ScheduledOperation {
	scheduledOperations := d.Get("scheduled_operation").(*schema.Set).List()
	so := make([]policies.ScheduledOperation, len(scheduledOperations))
	for i, raw := range scheduledOperations {
		rawMap := raw.(map[string]interface{})
		so[i] = policies.ScheduledOperation{
			Name:          rawMap["name"].(string),
			Description:   rawMap["description"].(string),
			Enabled:       rawMap["enabled"].(bool),
			OperationType: rawMap["operation_type"].(string),
			Trigger: policies.Trigger{
				Properties: policies.TriggerProperties{
					Pattern: rawMap["trigger_pattern"].(string),
				},
			},
			OperationDefinition: policies.OperationDefinition{
				MaxBackups:            rawMap["max_backups"].(int),
				RetentionDurationDays: rawMap["retention_duration_days"].(int),
				Permanent:             rawMap["permanent"].(bool),
			},
		}
	}

	log.Printf("[DEBUG] Printing so create options %#v", so)
	return so
}

func resourceCSBSResourceV1(d *schema.ResourceData) []policies.Resource {
	resources := d.Get("resource").(*schema.Set).List()
	res := make([]policies.Resource, len(resources))
	for i, raw := range resources {
		rawMap := raw.(map[string]interface{})
		res[i] = policies.Resource{
			Name: rawMap["name"].(string),
			Id:   rawMap["id"].(string),
			Type: rawMap["type"].(string),
		}
	}
	return res
}

func resourceCSBScheduleUpdateV1(d *schema.ResourceData) []policies.ScheduledOperationToUpdate {

	oldSORaw, newSORaw := d.GetChange("scheduled_operation")
	oldSOList := oldSORaw.(*schema.Set).List()
	newSOSetList := newSORaw.(*schema.Set).List()

	schedule := make([]policies.ScheduledOperationToUpdate, len(newSOSetList))
	for i, raw := range newSOSetList {
		rawNewMap := raw.(map[string]interface{})
		rawOldMap := oldSOList[i].(map[string]interface{})
		log.Printf("[DEBUG] Printing update options before %#v", rawNewMap)
		schedule[i] = policies.ScheduledOperationToUpdate{
			Id:          rawOldMap["id"].(string),
			Name:        rawNewMap["name"].(string),
			Description: rawNewMap["description"].(string),
			Enabled:     rawNewMap["enabled"].(bool),
			Trigger: policies.Trigger{
				Properties: policies.TriggerProperties{
					Pattern: rawNewMap["trigger_pattern"].(string),
				},
			},
			OperationDefinition: policies.OperationDefinition{
				MaxBackups:            rawNewMap["max_backups"].(int),
				RetentionDurationDays: rawNewMap["retention_duration_days"].(int),
				Permanent:             rawNewMap["permanent"].(bool),
			},
		}
	}

	log.Printf("[DEBUG] Printing update options after %#v", schedule)
	return schedule
}

func resourceCSBSCommonParamsV1(d *schema.ResourceData) map[string]string {
	m := make(map[string]string)
	for key, val := range d.Get("common").(map[string]interface{}) {
		m[key] = val.(string)
	}
	return m
}

func flattenCSBSScheduledOperations(backupPolicy policies.BackupPolicy) []map[string]interface{} {
	var scheduledOperationList []map[string]interface{}
	for _, schedule := range backupPolicy.ScheduledOperations {
		mapping := map[string]interface{}{
			"enabled":                 schedule.Enabled,
			"trigger_id":              schedule.TriggerID,
			"name":                    schedule.Name,
			"description":             schedule.Description,
			"operation_type":          schedule.OperationType,
			"max_backups":             schedule.OperationDefinition.MaxBackups,
			"retention_duration_days": schedule.OperationDefinition.RetentionDurationDays,
			"permanent":               schedule.OperationDefinition.Permanent,
			"trigger_name":            schedule.Trigger.Name,
			"trigger_type":            schedule.Trigger.Type,
			"trigger_pattern":         schedule.Trigger.Properties.Pattern,
			"id":                      schedule.ID,
		}
		scheduledOperationList = append(scheduledOperationList, mapping)
	}

	log.Printf("[DEBUG] Flatten List %#v", scheduledOperationList)

	return scheduledOperationList
}

func flattenCSBSPolicyResources(backupPolicy policies.BackupPolicy) []map[string]interface{} {
	var resourceList []map[string]interface{}
	for _, resources := range backupPolicy.Resources {
		mapping := map[string]interface{}{
			"id":   resources.Id,
			"type": resources.Type,
			"name": resources.Name,
		}
		resourceList = append(resourceList, mapping)
	}

	return resourceList
}
