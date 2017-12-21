package orangecloud

import (
	"fmt"
	"log"
	"time"

	// "github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/elbaas"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/elbaas/backendmember"
)

const loadbalancerActiveTimeoutSeconds = 300
const loadbalancerDeleteTimeoutSeconds = 300

func resourceBackend() *schema.Resource {
	return &schema.Resource{
		Create: resourceBackendCreate,
		Read:   resourceBackendRead,
		Delete: resourceBackendDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"listener_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"server_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"address": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceBackendCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client, err := config.otcV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomCloud networking client: %s", err)
	}

	addOpts := backendmember.AddOpts{
		ServerId: d.Get("server_id").(string),
		Address:  d.Get("address").(string),
	}
	log.Printf("[DEBUG] Create Options: %#v", addOpts)

	listener_id := d.Get("listener_id").(string)
	job, err := backendmember.Add(client, listener_id, addOpts).ExtractJobResponse()
	if err != nil {
		return err
	}

	log.Printf("Waiting for backend to become active, job=%#v", job)

	j := &elbaas.Job{Uri: job.URI, JobId: job.JobID}
	timeout := d.Timeout(schema.TimeoutCreate)
	jobInfo, err := waitForELBJobSuccess(client, j, timeout)
	if err != nil {
		return err
	}

	entity, ok := jobInfo.Entities["members"]
	if !ok {
		return fmt.Errorf("Error getting the entity from job info")
	}
	log.Printf("[DEBUG] get job entity: %#v", entity)

	if members, ok := entity.([]interface{}); ok {
		if len(members) > 0 {
			vmember := members[0]
			if member, ok := vmember.(map[string]interface{}); ok {
				if vid, ok := member["id"]; ok {
					if id, ok := vid.(string); ok {
						d.SetId(id)
						return resourceBackendRead(d, meta)
					}
				}
			}
		}
	}
	return fmt.Errorf("Unexpected conversion error in resourceBackendCreate.")
}

func resourceBackendRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client, err := config.otcV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomCloud networking client: %s", err)
	}

	listener_id := d.Get("listener_id").(string)
	b, err := backendmember.Get(client, listener_id, d.Id()).Extract()
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Retrieved backend: %#v", b)

	backend := b[0]
	d.Set("server_address", backend.ServerAddress)
	d.Set("address", backend.Address)
	d.Set("status", backend.Status)
	d.Set("health_status", backend.HealthStatus)
	d.Set("update_time", backend.UpdateTime)
	d.Set("create_time", backend.CreateTime)
	d.Set("server_name", backend.ServerName)
	d.Set("server_id", backend.ServerID)
	d.Set("listeners", backend.Listeners)

	return nil
}

func resourceBackendDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client, err := config.otcV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomCloud networking client: %s", err)
	}

	log.Printf("[DEBUG] Deleting backend member %s", d.Id())
	listener_id := d.Get("listener_id").(string)
	id := d.Id()
	job, err := backendmember.Remove(client, listener_id, id).ExtractJobResponse()
	if err != nil {
		return err
	}

	log.Printf("Waiting for backend member %s to delete", id)

	if err := gophercloud.WaitForJobSuccess(client, job.URI, loadbalancerActiveTimeoutSeconds); err != nil {
		return err
	}

	log.Printf("Successfully deleted backend member %s", id)
	return nil
}
