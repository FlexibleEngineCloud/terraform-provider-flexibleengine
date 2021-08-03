package flexibleengine

import (
	"time"

	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/openstack/cts/v1/tracker"
)

func resourceCTSTrackerV1() *schema.Resource {
	return &schema.Resource{
		Create: resourceCTSTrackerCreate,
		Read:   resourceCTSTrackerRead,
		Update: resourceCTSTrackerUpdate,
		Delete: resourceCTSTrackerDelete,
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
				Optional: true,
				Computed: true,
			},
			"tracker_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"bucket_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"file_prefix_name": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validateName,
			},
		},
	}

}

func resourceCTSTrackerCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	ctsClient, err := config.ctsV1Client(GetRegion(d, config))

	if err != nil {
		return fmt.Errorf("Error creating cts Client: %s", err)
	}

	createOpts := tracker.CreateOpts{
		BucketName:     d.Get("bucket_name").(string),
		FilePrefixName: d.Get("file_prefix_name").(string),
	}

	trackers, err := tracker.Create(ctsClient, createOpts).Extract()
	log.Printf("[DEBUG]trackers %#v", trackers)
	if err != nil {
		return fmt.Errorf("Error creating CTS tracker : %s", err)
	}

	d.SetId(trackers.TrackerName)

	time.Sleep(20 * time.Second)
	return resourceCTSTrackerRead(d, meta)
}

func resourceCTSTrackerRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	ctsClient, err := config.ctsV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating cts Client: %s", err)
	}
	listOpts := tracker.ListOpts{
		TrackerName:    d.Get("tracker_name").(string),
		BucketName:     d.Get("bucket_name").(string),
		FilePrefixName: d.Get("file_prefix_name").(string),
		Status:         d.Get("status").(string),
	}
	trackers, err := tracker.List(ctsClient, listOpts)
	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			log.Printf("[WARN] Removing cts tracker %s as it's already gone", d.Id())
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error retrieving cts tracker: %s", err)
	}

	ctsTracker := trackers[0]

	d.Set("tracker_name", ctsTracker.TrackerName)
	d.Set("bucket_name", ctsTracker.BucketName)
	d.Set("status", ctsTracker.Status)
	d.Set("file_prefix_name", ctsTracker.FilePrefixName)

	d.Set("region", GetRegion(d, config))
	time.Sleep(20 * time.Second)

	return nil
}

func resourceCTSTrackerUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	ctsClient, err := config.ctsV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating cts Client: %s", err)
	}
	var updateOpts tracker.UpdateOpts

	//as bucket_name is mandatory while updating tracker
	updateOpts.BucketName = d.Get("bucket_name").(string)

	if d.HasChange("file_prefix_name") {
		updateOpts.FilePrefixName = d.Get("file_prefix_name").(string)
	}
	if d.HasChange("status") {
		updateOpts.Status = d.Get("status").(string)
	}

	_, err = tracker.Update(ctsClient, updateOpts).Extract()
	if err != nil {
		return fmt.Errorf("Error updating cts tracker: %s", err)
	}
	time.Sleep(20 * time.Second)
	return resourceCTSTrackerRead(d, meta)
}

func resourceCTSTrackerDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	ctsClient, err := config.ctsV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating cts Client: %s", err)
	}

	result := tracker.Delete(ctsClient)
	if result.Err != nil {
		return err
	}

	time.Sleep(20 * time.Second)
	log.Printf("[DEBUG] Successfully deleted cts tracker %s", d.Id())

	return nil
}
