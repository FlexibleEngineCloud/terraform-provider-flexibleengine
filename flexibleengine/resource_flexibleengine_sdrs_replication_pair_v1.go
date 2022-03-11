package flexibleengine

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/sdrs/v1/replications"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceSdrsReplicationPairV1() *schema.Resource {
	return &schema.Resource{
		Create: resourceSdrsReplicationPairV1Create,
		Read:   resourceSdrsReplicationPairV1Read,
		Update: resourceSdrsReplicationPairV1Update,
		Delete: resourceSdrsReplicationPairV1Delete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"group_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"volume_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			// used for delete
			"delete_target_volume": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
				ForceNew: false,
			},
			// the following attributes are computed
			"replication_model": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"fault_level": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"target_volume_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceSdrsReplicationPairV1Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	sdrsClient, err := sdrsV1Client(config, GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine SDRS Client: %s", err)
	}

	createOpts := replications.CreateOpts{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		GroupID:     d.Get("group_id").(string),
		VolumeID:    d.Get("volume_id").(string),
	}
	log.Printf("[DEBUG] CreateOpts: %#v", createOpts)

	n, err := replications.Create(sdrsClient, createOpts).ExtractJobResponse()
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine SDRS Replication pair: %s", err)
	}

	if err := replications.WaitForJobSuccess(sdrsClient, int(d.Timeout(schema.TimeoutCreate)/time.Second), n.JobID); err != nil {
		return err
	}

	entity, err := replications.GetJobEntity(sdrsClient, n.JobID, "replication_pair_id")
	if err != nil {
		return err
	}

	if id, ok := entity.(string); ok {
		d.SetId(id)
		return resourceSdrsReplicationPairV1Read(d, meta)
	}

	return fmt.Errorf("Unexpected conversion error in resourceSdrsReplicationPairV1Create.")
}

func resourceSdrsReplicationPairV1Read(d *schema.ResourceData, meta interface{}) error {

	config := meta.(*Config)
	sdrsClient, err := sdrsV1Client(config, GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine SDRS client: %s", err)
	}
	n, err := replications.Get(sdrsClient, d.Id()).Extract()

	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error retrieving FlexibleEngine SDRS Replication pair: %s", err)
	}

	d.Set("name", n.Name)
	d.Set("description", n.Description)
	d.Set("group_id", n.GroupID)
	d.Set("status", n.Status)
	d.Set("replication_model", n.ReplicaModel)
	d.Set("fault_level", n.FaultLevel)

	// set "volume_id" and "target_volume_id" from VolumeIDs
	volumes := strings.Split(n.VolumeIDs, ",")
	if len(volumes) < 2 {
		return fmt.Errorf("Error retrieving VolumeIDs od FlexibleEngine SDRS Replication pair: Invalid format.")
	}
	d.Set("volume_id", volumes[0])
	d.Set("target_volume_id", volumes[1])

	return nil
}

func resourceSdrsReplicationPairV1Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	sdrsClient, err := sdrsV1Client(config, GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine SDRS Client: %s", err)
	}
	var updateOpts replications.UpdateOpts

	if d.HasChange("name") {
		updateOpts.Name = d.Get("name").(string)
	}
	log.Printf("[DEBUG] updateOpts: %#v", updateOpts)

	_, err = replications.Update(sdrsClient, d.Id(), updateOpts).Extract()
	if err != nil {
		return fmt.Errorf("Error updating FlexibleEngine SDRS Replication pair: %s", err)
	}
	return resourceSdrsReplicationPairV1Read(d, meta)
}

func resourceSdrsReplicationPairV1Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	sdrsClient, err := sdrsV1Client(config, GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine SDRS client: %s", err)
	}

	deleteOpts := replications.DeleteOpts{
		GroupID:      d.Get("group_id").(string),
		DeleteVolume: d.Get("delete_target_volume").(bool),
	}
	n, err := replications.Delete(sdrsClient, d.Id(), deleteOpts).ExtractJobResponse()
	if err != nil {
		return fmt.Errorf("Error deleting FlexibleEngine SDRS Replication pair: %s", err)
	}

	if err := replications.WaitForJobSuccess(sdrsClient, int(d.Timeout(schema.TimeoutDelete)/time.Second), n.JobID); err != nil {
		return err
	}

	d.SetId("")
	return nil
}
