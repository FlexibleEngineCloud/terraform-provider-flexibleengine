package flexibleengine

import (
	"fmt"
	"log"
	"time"

	"github.com/chnsz/golangsdk/openstack/cts/v1/tracker"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
			"status": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"tracker_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}

}

func resourceCTSTrackerCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	ctsClient, err := config.CtsV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating cts Client: %s", err)
	}

	createOpts := tracker.CreateOpts{
		BucketName:     d.Get("bucket_name").(string),
		FilePrefixName: d.Get("file_prefix_name").(string),
	}

	log.Printf("[DEBUG] CTS tracker creating options: %#v", createOpts)
	trackers, err := tracker.Create(ctsClient, createOpts).Extract()
	if err != nil {
		return fmt.Errorf("Error creating CTS tracker : %s", err)
	}

	d.SetId(trackers.TrackerName)

	time.Sleep(5 * time.Second)
	return resourceCTSTrackerRead(d, meta)
}

func resourceCTSTrackerRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	region := GetRegion(d, config)
	ctsClient, err := config.CtsV1Client(region)
	if err != nil {
		return fmt.Errorf("Error creating cts Client: %s", err)
	}

	trackerName := d.Id()
	listOpts := tracker.ListOpts{
		TrackerName: trackerName,
	}
	trackers, err := tracker.List(ctsClient, listOpts)
	if err != nil {
		return CheckDeleted(d, err, "CTS tracker")
	}

	if len(trackers) == 0 {
		return fmt.Errorf("cannot find CTS tracker %s", trackerName)
	}

	ctsTracker := trackers[0]
	log.Printf("[DEBUG] fetching CTS tracker: %#v", ctsTracker)

	d.Set("region", region)
	d.Set("bucket_name", ctsTracker.BucketName)
	d.Set("file_prefix_name", ctsTracker.FilePrefixName)
	d.Set("tracker_name", ctsTracker.TrackerName)
	d.Set("status", ctsTracker.Status)

	time.Sleep(5 * time.Second)
	return nil
}

func resourceCTSTrackerUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	ctsClient, err := config.CtsV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating cts Client: %s", err)
	}
	var updateOpts tracker.UpdateOpts

	// bucket_name is mandatory while updating tracker
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

	time.Sleep(10 * time.Second)
	return resourceCTSTrackerRead(d, meta)
}

func resourceCTSTrackerDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	ctsClient, err := config.CtsV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating cts Client: %s", err)
	}

	result := tracker.Delete(ctsClient)
	if result.Err != nil {
		return err
	}

	time.Sleep(10 * time.Second)
	log.Printf("[DEBUG] Successfully deleted cts tracker %s", d.Id())

	return nil
}
