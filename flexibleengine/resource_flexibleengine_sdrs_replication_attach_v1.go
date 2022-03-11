package flexibleengine

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/sdrs/v1/attachreplication"
	"github.com/chnsz/golangsdk/openstack/sdrs/v1/protectedinstances"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ExtractAttachId(s string) (string, string) {
	rgs := strings.Split(s, ":")
	if len(rgs) >= 2 {
		log.Printf("[DEBUG] ExtractAttachId: %s:%s from (%s)", rgs[0], rgs[1], s)
		return rgs[0], rgs[1]
	}
	log.Printf("[DEBUG] ExtractAttachId: length of string < 2")
	return "", ""
}

func FormatAttachId(insId string, id string) string {
	return fmt.Sprintf("%s:%s", insId, id)
}

func resourceSdrsReplicationAttachV1() *schema.Resource {
	return &schema.Resource{
		Create: resourceSdrsReplicationAttachV1Create,
		Read:   resourceSdrsReplicationAttachV1Read,
		Delete: resourceSdrsReplicationAttachV1Delete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"replication_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"device": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceSdrsReplicationAttachV1Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	sdrsClient, err := sdrsV1Client(config, GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine SDRS client: %s", err)
	}

	instanceID := d.Get("instance_id").(string)
	replicationID := d.Get("replication_id").(string)

	attachOpts := attachreplication.CreateOpts{
		ReplicationID: replicationID,
		Device:        d.Get("device").(string),
	}

	log.Printf("[DEBUG] Creating replication attachment: %#v", attachOpts)

	n, err := attachreplication.Create(sdrsClient, instanceID, attachOpts).ExtractJobResponse()
	if err != nil {
		return err
	}

	if err := attachreplication.WaitForJobSuccess(sdrsClient, int(d.Timeout(schema.TimeoutCreate)/time.Second), n.JobID); err != nil {
		return err
	}

	d.SetId(FormatAttachId(instanceID, replicationID))

	return resourceSdrsReplicationAttachV1Read(d, meta)
}

func resourceSdrsReplicationAttachV1Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	sdrsClient, err := sdrsV1Client(config, GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine SDRS client: %s", err)
	}
	instId, replicaId := ExtractAttachId(d.Id())
	n, err := protectedinstances.Get(sdrsClient, instId).Extract()

	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error retrieving FlexibleEngine SDRS ProtectedInstance: %s", err)
	}

	find := false
	var attach protectedinstances.Attachment
	for _, attach = range n.Attachment {
		if attach.Replication == replicaId {
			find = true
			break
		}
	}
	if find == false {
		d.SetId("")
		return nil
	}
	log.Printf("[DEBUG] Retrieved replication attachment: %#v", attach)

	d.Set("device", attach.Device)
	d.Set("replication_id", attach.Replication)
	d.Set("status", n.Status)
	return nil
}

func resourceSdrsReplicationAttachV1Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	sdrsClient, err := sdrsV1Client(config, GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine SDRS client: %s", err)
	}

	instId, replicaId := ExtractAttachId(d.Id())
	n, err := attachreplication.Delete(sdrsClient, instId, replicaId).ExtractJobResponse()
	if err != nil {
		return err
	}

	if err := attachreplication.WaitForJobSuccess(sdrsClient, int(d.Timeout(schema.TimeoutCreate)/time.Second), n.JobID); err != nil {
		return err
	}

	return nil
}
