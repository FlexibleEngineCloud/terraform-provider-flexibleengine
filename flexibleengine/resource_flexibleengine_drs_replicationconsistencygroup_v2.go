package flexibleengine

import (
	"fmt"
	"log"
	"time"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/drs/v2/replicationconsistencygroups"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// resourceReplicationConsistencyGroup defines the schema of replication consistency group
func resourceReplicationConsistencyGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceReplicationConsistencyGroupCreate,
		Read:   resourceReplicationConsistencyGroupRead,
		Delete: resourceReplicationConsistencyGroupDelete,
		Update: resourceReplicationConsistencyGroupUpdate,

		DeprecationMessage: "It has been deprecated",
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"replication_ids": {
				Type:     schema.TypeSet,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},
			"priority_station": {
				Type:     schema.TypeString,
				Required: true,
			},
			// Default:  "hypermetro",
			"replication_model": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "hypermetro",
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"replication_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"failure_detail": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"fault_level": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

// resourceReplicationIDsFromSchema returns replication ids
func resourceReplicationIDsFromSchema(d *schema.ResourceData) []string {
	rawReplicationIDs := d.Get("replication_ids").(*schema.Set)
	replicationids := make([]string, rawReplicationIDs.Len())
	for i, raw := range rawReplicationIDs.List() {
		replicationids[i] = raw.(string)
	}
	return replicationids
}

// resourceReplicationConsistencyGroupCreate creates a replication consistency group resource
func resourceReplicationConsistencyGroupCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client, err := drsV2Client(config, GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine drs client: %s", err)
	}

	createOpts := replicationconsistencygroups.CreateOps{
		Name:             d.Get("name").(string),
		Description:      d.Get("description").(string),
		ReplicationIDs:   resourceReplicationIDsFromSchema(d),
		PriorityStation:  d.Get("priority_station").(string),
		ReplicationModel: d.Get("replication_model").(string),
	}
	log.Printf("[DEBUG] Create replication Options: %#v", createOpts)

	rcg, err := replicationconsistencygroups.Create(client, createOpts).Extract()
	if err != nil {
		return fmt.Errorf("Error getting replicationconsistencygroup from result: %s", err)
	}

	log.Printf("[DEBUG] Waiting for replicationconsistencygroup (%s) to become available", rcg.ID)
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"creating"},
		Target:     []string{"available"},
		Refresh:    replicationconsistencygroupStateRefreshFunc(client, rcg.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Error waiting for replicationconsistencygroup (%s) to create: %s", rcg.ID, err)
	}

	d.SetId(rcg.ID)
	log.Printf("[DEBUG] Created replicationconsistencygroup (%s): %#v", rcg.ID, rcg)

	// Sync
	log.Printf("[DEBUG] Syncing replicationconsistencygroup (%s)", rcg.ID)
	syncResult := replicationconsistencygroups.Sync(client, rcg.ID)
	if syncResult.Err != nil {
		log.Printf("[DEBUG] Error syncing replicationconsistencygroup %s", syncResult.Err)
		return syncResult.Err
	}

	return resourceReplicationConsistencyGroupRead(d, meta)
}

// resourceReplicationConsistencyGroupRead returns a replication consistency group resource
func resourceReplicationConsistencyGroupRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client, err := drsV2Client(config, GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine drs client: %s", err)
	}

	id := d.Id()
	rcg, err := replicationconsistencygroups.Get(client, id).Extract()
	if err != nil {
		return CheckDeleted(d, err, "replicationconsistencygroup")
	}

	log.Printf("[DEBUG] Read replicationconsistencygroup (%s): %#v", id, rcg)

	d.SetId(rcg.ID)
	d.Set("name", rcg.Name)
	d.Set("description", rcg.Description)
	d.Set("status", rcg.Status)
	d.Set("priority_station", rcg.PriorityStation)
	d.Set("replication_model", rcg.ReplicationModel)
	d.Set("replication_status", rcg.ReplicationStatus)
	d.Set("replication_ids", rcg.ReplicationIDs)

	d.Set("created_at", rcg.CreatedAt)
	d.Set("updated_at", rcg.UpdatedAt)
	d.Set("failure_detail", rcg.FailureDetail)
	d.Set("fault_level", rcg.FaultLevel)

	return nil
}

// replicationconsistencygroupStateRefreshFunc is used to watch a replication consistency group state.
func replicationconsistencygroupStateRefreshFunc(client *golangsdk.ServiceClient, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		rcg, err := replicationconsistencygroups.Get(client, id).Extract()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return rcg, "deleted", nil
			}
			return nil, "", err
		}

		log.Printf("[DEBUG] replicationconsistencygroup (%s) current status: %s", rcg.ID, rcg.Status)
		return rcg, rcg.Status, nil
	}
}

// resourceReplicationConsistencyGroupDelete deletes a replication consistency group resource
func resourceReplicationConsistencyGroupDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client, err := drsV2Client(config, GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine drs client: %s", err)
	}

	// Delete Workflow: Stop => Delete
	id := d.Id()
	log.Printf("[DEBUG] Stopping replicationconsistencygroup (%s)", id)

	// Stop
	stopResult := replicationconsistencygroups.Stop(client, id)
	if stopResult.Err != nil {
		log.Printf("[DEBUG] Error stopping replicationconsistencygroup %s", stopResult.Err)
		return stopResult.Err
	}

	log.Printf("[DEBUG] Deleting replicationconsistencygroup (%s)", id)

	// Remove replications in replicationconsistencygroups
	rcg, err := replicationconsistencygroups.Get(client, id).Extract()
	if err != nil {
		return CheckDeleted(d, err, "replicationconsistencygroup")
	}

	// ReplicationIDs is more than 0
	if len(rcg.ReplicationIDs) > 0 {
		removeOpts := replicationconsistencygroups.UpdateOpts{
			RemoveReplicationIDs: rcg.ReplicationIDs,
		}

		log.Printf("[DEBUG] Removing replicationconsistencygroup (%s): %#v", id, removeOpts)
		_, err = replicationconsistencygroups.Update(client, id, removeOpts).Extract()
		if err != nil {
			return fmt.Errorf("Error Removing replicationconsistencygroup from result: %s", err)
		}

		// Wait for the replicationconsistencygroup to remove before moving on.
		log.Printf("[DEBUG] Waiting for replicationconsistencygroup (%s) to remove", id)

		stateRemoveConf := &resource.StateChangeConf{
			Pending:    []string{"updating"},
			Target:     []string{"available"},
			Refresh:    replicationconsistencygroupStateRefreshFunc(client, id),
			Timeout:    d.Timeout(schema.TimeoutUpdate),
			Delay:      5 * time.Second,
			MinTimeout: 3 * time.Second,
		}

		_, err = stateRemoveConf.WaitForState()
		if err != nil {
			return fmt.Errorf("Error waiting for replicationconsistencygroup (%s) to remove: %s", id, err)
		}
	}

	// Delete
	result := replicationconsistencygroups.Delete(client, id)
	if result.Err != nil {
		log.Printf("[DEBUG] Error deleting replicationconsistencygroup %s", result.Err)
		return result.Err
	}

	// Wait for the replicationconsistencygroup to delete before moving on.
	log.Printf("[DEBUG] Waiting for replicationconsistencygroup (%s) to delete", id)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"deleting", "available"},
		Target:     []string{"deleted"},
		Refresh:    replicationconsistencygroupStateRefreshFunc(client, id),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Error waiting for replicationconsistencygroup (%s) to delete: %s", id, err)
	}

	d.SetId("")
	log.Printf("[DEBUG] Successfully deleted replicationconsistencygroup (%s)", id)
	return nil
}

// resourceReplicationConsistencyGroupUpdate updates a replication consistency group resource
func resourceReplicationConsistencyGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client, err := drsV2Client(config, GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine drs client: %s", err)
	}

	// Update Workflow: Stop => Update => Sync
	id := d.Id()
	log.Printf("[DEBUG] Stopping replicationconsistencygroup (%s)", id)

	// Stop
	stopResult := replicationconsistencygroups.Stop(client, id)
	if stopResult.Err != nil {
		log.Printf("[DEBUG] Error stopping replicationconsistencygroup %s", stopResult.Err)
		return stopResult.Err
	}

	log.Printf("[DEBUG] Updating replicationconsistencygroup %s", id)

	// Update
	updateOpts := replicationconsistencygroups.UpdateOpts{}
	if d.HasChange("name") {
		updateOpts.Name = d.Get("name").(string)
	}
	if d.HasChange("description") {
		updateOpts.Description = d.Get("description").(string)
	}
	if d.HasChange("replication_model") {
		updateOpts.ReplicationModel = d.Get("replication_model").(string)
	}
	if d.HasChange("replication_ids") {
		addresults, removeresults, errors := resourceGetReplicationIDs(d, meta)
		if errors != nil {
			return fmt.Errorf("Error getting replicationconsistencygroup replicationids: %s", err)
		}
		if len(addresults) > 0 {
			updateOpts.AddReplicationIDs = addresults
		}
		if len(removeresults) > 0 {
			updateOpts.RemoveReplicationIDs = removeresults
		}
	}
	log.Printf("[DEBUG] Updating replicationconsistencygroup (%s): %#v", id, updateOpts)

	_, err = replicationconsistencygroups.Update(client, id, updateOpts).Extract()
	if err != nil {
		return fmt.Errorf("Error updating replicationconsistencygroup from result: %s", err)
	}

	// Wait for the replicationconsistencygroup to update before moving on.
	log.Printf("[DEBUG] Waiting for replicationconsistencygroup (%s) to update", id)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"updating", "creating"},
		Target:     []string{"available"},
		Refresh:    replicationconsistencygroupStateRefreshFunc(client, id),
		Timeout:    d.Timeout(schema.TimeoutUpdate),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Error waiting for replicationconsistencygroup (%s) to update: %s", id, err)
	}

	log.Printf("[DEBUG] Syncing replicationconsistencygroup (%s)", id)

	// Sync
	syncResult := replicationconsistencygroups.Sync(client, id)
	if syncResult.Err != nil {
		log.Printf("[DEBUG] Error syncing replicationconsistencygroup %s", syncResult.Err)
		return syncResult.Err
	}

	return resourceReplicationConsistencyGroupRead(d, meta)
}

// resourceGetReplicationIDs returns add and remove replication ids
func resourceGetReplicationIDs(d *schema.ResourceData, meta interface{}) ([]string, []string, error) {
	config := meta.(*Config)
	client, err := drsV2Client(config, GetRegion(d, config))
	if err != nil {
		return nil, nil, fmt.Errorf("Error creating FlexibleEngine drs client: %s", err)
	}

	// get replication ids from database
	id := d.Id()
	rcg, err := replicationconsistencygroups.Get(client, id).Extract()
	if err != nil {
		return nil, nil, CheckDeleted(d, err, "replicationconsistencygroup")
	}

	// compare replication ids between database and cache
	cacheReplicationIDs := resourceReplicationIDsFromSchema(d)

	// get add replication ids
	addreplicationids := make(map[string]string)
	for _, cache := range cacheReplicationIDs {
		// Check if cache replication is exist in database
		found := false
		for _, raw := range rcg.ReplicationIDs {
			if cache == raw {
				found = true
				break
			}
		}

		// If cache replication is not exist in database
		if !found {
			addreplicationids[cache] = cache
		}
	}

	// convert add results from map to array
	addresults := make([]string, len(addreplicationids))
	var addindex = 0
	for addvalue := range addreplicationids {
		addresults[addindex] = addvalue
		addindex++
	}

	// get remove replication ids
	removereplicationids := make(map[string]string)
	for _, cache := range rcg.ReplicationIDs {
		// Check if database replication is exist in cache
		found := false
		for _, raw := range cacheReplicationIDs {
			if cache == raw {
				found = true
				break
			}
		}

		// If database replication is not exist in cache
		if !found {
			removereplicationids[cache] = cache
		}
	}

	// convert remove results from map to array
	removeresults := make([]string, len(removereplicationids))
	var removeindex = 0
	for removevalue := range removereplicationids {
		removeresults[removeindex] = removevalue
		removeindex++
	}

	return addresults, removeresults, nil
}
