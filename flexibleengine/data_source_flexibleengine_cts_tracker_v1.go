package flexibleengine

import (
	"fmt"
	"log"

	"github.com/chnsz/golangsdk/openstack/cts/v1/tracker"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceCTSTrackerV1() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceCTSTrackerV1Read,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"bucket_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"file_prefix_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"tracker_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func dataSourceCTSTrackerV1Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	region := GetRegion(d, config)
	trackerClient, err := config.CtsV1Client(region)
	if err != nil {
		return fmt.Errorf("Error creating CTS client: %s", err)
	}

	listOpts := tracker.ListOpts{
		TrackerName:    d.Get("tracker_name").(string),
		BucketName:     d.Get("bucket_name").(string),
		FilePrefixName: d.Get("file_prefix_name").(string),
		Status:         d.Get("status").(string),
	}

	refinedTrackers, err := tracker.List(trackerClient, listOpts)
	if err != nil {
		return fmt.Errorf("Unable to retrieve cts tracker: %s", err)
	}

	if len(refinedTrackers) < 1 {
		return fmt.Errorf("Your query returned no results. " +
			"Please change your search criteria and try again.")
	}

	if len(refinedTrackers) > 1 {
		return fmt.Errorf("Your query returned more than one result." +
			" Please try a more specific search criteria")
	}

	trackers := refinedTrackers[0]
	log.Printf("[INFO] Retrieved cts tracker %s using given filter", trackers.TrackerName)

	d.SetId(trackers.TrackerName)

	d.Set("region", region)
	d.Set("tracker_name", trackers.TrackerName)
	d.Set("bucket_name", trackers.BucketName)
	d.Set("file_prefix_name", trackers.FilePrefixName)
	d.Set("status", trackers.Status)

	return nil
}
