package flexibleengine

import (
	"fmt"
	"log"
	"time"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/sdrs/v1/protectedinstances"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceSdrsProtectedInstanceV1() *schema.Resource {
	return &schema.Resource{
		Create: resourceSdrsProtectedInstanceV1Create,
		Read:   resourceSdrsProtectedInstanceV1Read,
		Update: resourceSdrsProtectedInstanceV1Update,
		Delete: resourceSdrsProtectedInstanceV1Delete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"group_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"server_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
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
			"cluster_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"primary_subnet_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"primary_ip_address": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"delete_target_server": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"delete_target_eip": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"target_server": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceSdrsProtectedInstanceV1Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	sdrsClient, err := sdrsV1Client(config, GetRegion(d, config))

	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine SDRS Client: %s", err)
	}

	createOpts := protectedinstances.CreateOpts{
		GroupID:     d.Get("group_id").(string),
		ServerID:    d.Get("server_id").(string),
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		ClusterID:   d.Get("cluster_id").(string),
		SubnetID:    d.Get("primary_subnet_id").(string),
		IpAddress:   d.Get("primary_ip_address").(string),
	}
	log.Printf("[DEBUG] CreateOpts: %#v", createOpts)

	n, err := protectedinstances.Create(sdrsClient, createOpts).ExtractJobResponse()
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine SDRS Protectedinstance: %s", err)
	}

	if err := protectedinstances.WaitForJobSuccess(sdrsClient, int(d.Timeout(schema.TimeoutCreate)/time.Second), n.JobID); err != nil {
		return err
	}

	entity, err := protectedinstances.GetJobEntity(sdrsClient, n.JobID, "protected_instance_id")
	if err != nil {
		return err
	}

	if id, ok := entity.(string); ok {
		d.SetId(id)
		return resourceSdrsProtectedInstanceV1Read(d, meta)
	}

	return fmt.Errorf("Unexpected conversion error in resourceSdrsProtectedInstanceV1Create.")
}

func resourceSdrsProtectedInstanceV1Read(d *schema.ResourceData, meta interface{}) error {

	config := meta.(*Config)
	sdrsClient, err := sdrsV1Client(config, GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine SDRS client: %s", err)
	}
	n, err := protectedinstances.Get(sdrsClient, d.Id()).Extract()

	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error retrieving FlexibleEngine SDRS ProtectedInstance: %s", err)
	}

	d.Set("name", n.Name)
	d.Set("description", n.Description)
	d.Set("group_id", n.GroupID)
	d.Set("target_server", n.TargetServer)

	return nil
}

func resourceSdrsProtectedInstanceV1Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	sdrsClient, err := sdrsV1Client(config, GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine SDRS Client: %s", err)
	}
	var updateOpts protectedinstances.UpdateOpts

	if d.HasChange("name") {
		updateOpts.Name = d.Get("name").(string)
	}
	log.Printf("[DEBUG] updateOpts: %#v", updateOpts)

	_, err = protectedinstances.Update(sdrsClient, d.Id(), updateOpts).Extract()
	if err != nil {
		return fmt.Errorf("Error updating FlexibleEngine SDRS ProtectedInstance: %s", err)
	}
	return resourceSdrsProtectedInstanceV1Read(d, meta)
}

func resourceSdrsProtectedInstanceV1Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	sdrsClient, err := sdrsV1Client(config, GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine SDRS client: %s", err)
	}

	deleteOpts := protectedinstances.DeleteOpts{
		DeleteTargetServer: d.Get("delete_target_server").(bool),
		DeleteTargetEip:    d.Get("delete_target_eip").(bool),
	}
	log.Printf("[DEBUG] CreateOpts: %#v", deleteOpts)

	n, err := protectedinstances.Delete(sdrsClient, d.Id(), deleteOpts).ExtractJobResponse()
	if err != nil {
		return fmt.Errorf("Error deleting FlexibleEngine SDRS ProtectedInstance: %s", err)
	}
	if err := protectedinstances.WaitForJobSuccess(sdrsClient, int(d.Timeout(schema.TimeoutDelete)/time.Second), n.JobID); err != nil {
		return err
	}

	d.SetId("")
	return nil
}
