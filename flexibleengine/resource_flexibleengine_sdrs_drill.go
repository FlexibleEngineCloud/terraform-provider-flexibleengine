package flexibleengine

import (
	"fmt"
	"log"
	"time"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/sdrs/v1/drill"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceSdrsDrillV1() *schema.Resource {
	return &schema.Resource{
		Create: resourceSdrsDrillV1Create,
		Read:   resourceSdrsDrillV1Read,
		Update: resourceSdrsDrillV1Update,
		Delete: resourceSdrsDrillV1Delete,
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
			"group_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"drill_vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceSdrsDrillV1Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	sdrsClient, err := sdrsV1Client(config, GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine SDRS Client: %s", err)
	}

	createOpts := drill.CreateOpts{
		Name:       d.Get("name").(string),
		GroupID:    d.Get("group_id").(string),
		DrillVpcID: d.Get("drill_vpc_id").(string),
	}
	log.Printf("[DEBUG] CreateOpts: %#v", createOpts)

	n, err := drill.Create(sdrsClient, createOpts).ExtractJobResponse()
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine SDRS DR drill: %s", err)
	}

	if err := drill.WaitForJobSuccess(sdrsClient, int(d.Timeout(schema.TimeoutCreate)/time.Second), n.JobID); err != nil {
		return err
	}

	drillId, err := drill.GetJobEntity(sdrsClient, n.JobID, "disaster_recovery_drill_id")
	if err != nil {
		return err
	}

	if id, ok := drillId.(string); ok {
		d.SetId(id)
		return resourceSdrsDrillV1Read(d, meta)
	}

	return fmt.Errorf("Unexpected conversion error in resourceSdrsDrillV1Create.")
}

func resourceSdrsDrillV1Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	sdrsClient, err := sdrsV1Client(config, GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine SDRS client: %s", err)
	}

	n, err := drill.Get(sdrsClient, d.Id()).Extract()
	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error retrieving FlexibleEngine SDRS DR drill: %s", err)
	}

	d.Set("name", n.Name)
	d.Set("group_id", n.GroupID)
	d.Set("drill_vpc_id", n.DrillVpcID)
	d.Set("status", n.Status)

	return nil
}

func resourceSdrsDrillV1Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	sdrsClient, err := sdrsV1Client(config, GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine SDRS Client: %s", err)
	}
	var updateOpts drill.UpdateOpts

	if d.HasChange("name") {
		updateOpts.Name = d.Get("name").(string)
	}
	log.Printf("[DEBUG] updateOpts: %#v", updateOpts)

	_, err = drill.Update(sdrsClient, d.Id(), updateOpts).Extract()
	if err != nil {
		return fmt.Errorf("Error updating FlexibleEngine SDRS DR drill: %s", err)
	}
	return resourceSdrsDrillV1Read(d, meta)
}

func resourceSdrsDrillV1Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	sdrsClient, err := sdrsV1Client(config, GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine SDRS client: %s", err)
	}

	n, err := drill.Delete(sdrsClient, d.Id()).ExtractJobResponse()
	if err != nil {
		return fmt.Errorf("Error deleting FlexibleEngine SDRS DR drill: %s", err)
	}

	if err := drill.WaitForJobSuccess(sdrsClient, int(d.Timeout(schema.TimeoutDelete)/time.Second), n.JobID); err != nil {
		return err
	}

	d.SetId("")
	return nil
}
