package flexibleengine

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/openstack/blockstorage/v2/volumes"
	"github.com/huaweicloud/golangsdk/openstack/compute/v2/extensions/startstop"
	"github.com/huaweicloud/golangsdk/openstack/drs/v2/replications"
)

// resourceReplication defines the schema of replication
func resourceReplication() *schema.Resource {
	return &schema.Resource{
		Create: resourceReplicationCreate,
		Read:   resourceReplicationRead,
		Delete: resourceReplicationDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			// volume_ids maybe list[Request] or string[Response]
			"volume_ids": &schema.Schema{
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"priority_station": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			// Default:  "hypermetro",
			"replication_model": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "hypermetro",
			},
			"status": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"replication_consistency_group_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"replication_status": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"progress": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"failure_detail": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"record_metadata": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"fault_level": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

// resourceVolumeIDsFromSchema returns volume ids
func resourceVolumeIDsFromSchema(d *schema.ResourceData) []string {
	rawVolumeIDs := d.Get("volume_ids").([]interface{})
	volumeids := make([]string, len(rawVolumeIDs))
	for i, raw := range rawVolumeIDs {
		volumeids[i] = raw.(string)
	}
	return volumeids
}

// resourceShutdownInstance shutdowns the ecs the volume is attached to in disaster recovery environment
func resourceShutdownInstance(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	blockStorageClient, err := config.blockStorageV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine block storage client: %s", err)
	}

	volumeids := resourceVolumeIDsFromSchema(d)
	prioritystation := d.Get("priority_station").(string)

	for _, volumeid := range volumeids {
		v, err := volumes.Get(blockStorageClient, volumeid).Extract()
		if err != nil {
			return CheckDeleted(d, err, "volume")
		}

		if v.AvailabilityZone != prioritystation {
			log.Printf("[DEBUG] Get disaster recovery volume (%s): %#v", volumeid, v)

			computeClient, err := config.computeV2Client(GetRegion(d, config))
			if err != nil {
				return fmt.Errorf("Error creating FlexibleEngine compute client: %s", err)
			}

			for _, attachment := range v.Attachments {
				log.Printf("[DEBUG] Get disaster recovery instance %#v", attachment)
				err = startstop.Stop(computeClient, attachment.ServerID).ExtractErr()
				if err != nil {
					log.Printf("[WARN] Error stopping FlexibleEngine instance: %s", err)
				} else {
					stopStateConf := &resource.StateChangeConf{
						Pending:    []string{"ACTIVE"},
						Target:     []string{"SHUTOFF"},
						Refresh:    ServerV2StateRefreshFunc(computeClient, attachment.ServerID),
						Timeout:    3 * time.Minute,
						Delay:      10 * time.Second,
						MinTimeout: 3 * time.Second,
					}
					log.Printf("[DEBUG] Waiting for instance (%s) to stop", attachment.ServerID)
					_, err = stopStateConf.WaitForState()
					if err != nil {
						return fmt.Errorf("Error waiting for instance (%s) to stop: %s", attachment.ServerID, err)
					}
				}
			}
		}
	}
	return nil
}

// resourceReplicationCreate creates a replication resource
func resourceReplicationCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client, err := config.drsV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine drs client: %s", err)
	}

	createOpts := replications.CreateOps{
		Name:             d.Get("name").(string),
		Description:      d.Get("description").(string),
		VolumeIDs:        resourceVolumeIDsFromSchema(d),
		PriorityStation:  d.Get("priority_station").(string),
		ReplicationModel: d.Get("replication_model").(string),
	}
	log.Printf("[DEBUG] Create replication Options: %#v", createOpts)

	// Shutdown the ecs the volume is attached to in disaster recovery environment
	err = resourceShutdownInstance(d, meta)
	if err != nil {
		return fmt.Errorf("Error shutdown instance: %s", err)
	}

	replication, err := replications.Create(client, createOpts).Extract()
	if err != nil {
		return fmt.Errorf("Error getting replication from result: %s", err)
	}

	log.Printf("[DEBUG] Waiting for replication (%s) to become available", replication.ID)
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"creating"},
		Target:     []string{"available"},
		Refresh:    replicationStateRefreshFunc(client, replication.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Error waiting for replication (%s) to create: %s", replication.ID, err)
	}

	d.SetId(replication.ID)
	log.Printf("[DEBUG] Created replication (%s): %#v", replication.ID, replication)
	return resourceReplicationRead(d, meta)
}

// resourceVolumeIDsFromString returns volume ids
func resourceVolumeIDsFromString(VolumeIDs string) []string {
	volumeids := []string{}
	ids := strings.Split(VolumeIDs, ",")
	for _, id := range ids {
		volumeids = append(volumeids, strings.TrimSpace(id))
	}
	return volumeids
}

// resourceReplicationRead returns a replication resource
func resourceReplicationRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client, err := config.drsV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine drs client: %s", err)
	}

	id := d.Id()
	replication, err := replications.Get(client, id).Extract()
	if err != nil {
		return CheckDeleted(d, err, "replication")
	}

	log.Printf("[DEBUG] Read replication (%s): %#v", id, replication)

	d.Set("id", replication.ID)
	d.Set("name", replication.Name)
	d.Set("description", replication.Description)
	d.Set("status", replication.Status)
	d.Set("replication_consistency_group_id", replication.ReplicationConsistencyGroupID)
	// String => TypeList
	d.Set("volume_ids", resourceVolumeIDsFromString(replication.VolumeIDs))
	d.Set("priority_station", replication.PriorityStation)
	d.Set("created_at", replication.CreatedAt)
	d.Set("updated_at", replication.UpdatedAt)
	d.Set("replication_model", replication.ReplicationModel)
	d.Set("replication_status", replication.ReplicationStatus)
	d.Set("progress", replication.Progress)
	d.Set("failure_detail", replication.FailureDetail)
	// TypeMap => NoChange
	// RecordMetadata includes volume_type and multiattach currently.
	d.Set("record_metadata", string(replication.RecordMetadata))
	d.Set("fault_level", replication.FaultLevel)

	return nil
}

// replicationStateRefreshFunc is used to watch a replication state.
func replicationStateRefreshFunc(client *golangsdk.ServiceClient, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		r, err := replications.Get(client, id).Extract()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return r, "deleted", nil
			}
			return nil, "", err
		}

		log.Printf("[DEBUG] replication (%s) current status: %s", r.ID, r.Status)
		return r, r.Status, nil
	}
}

// resourceReplicationDelete deletes a replication resource
func resourceReplicationDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client, err := config.drsV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine drs client: %s", err)
	}

	id := d.Id()
	log.Printf("[DEBUG] Deleting replication (%s)", id)

	result := replications.Delete(client, id)
	if result.Err != nil {
		log.Printf("[DEBUG] Error deleting replication %s", result.Err)
		return result.Err
	}

	// Wait for the replication to delete before moving on.
	log.Printf("[DEBUG] Waiting for replication (%s) to delete", id)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"deleting", "available"},
		Target:     []string{"deleted"},
		Refresh:    replicationStateRefreshFunc(client, id),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Error waiting for replication (%s) to delete: %s", id, err)
	}

	d.SetId("")
	log.Printf("[DEBUG] Successfully deleted replication (%s)", id)
	return nil
}
