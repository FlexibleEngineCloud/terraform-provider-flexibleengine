package flexibleengine

import (
	"fmt"
	"log"
	"regexp"

	"github.com/chnsz/golangsdk/openstack/dis/v2/streams"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var regexpStreamName = regexp.MustCompile(`^[a-zA-Z0-9_-]{1,64}$`)

func resourceDisStreamV2() *schema.Resource {
	return &schema.Resource{
		Create: resourceDisStreamCreate,
		Read:   resourceDisStreamRead,
		Update: resourceDisStreamUpdate,
		Delete: resourceDisStreamDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringMatch(regexpStreamName,
					"1 to 64 in length, only letters, digits, hyphens (-), and underscores (_) are allowed."),
			},
			"type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      streams.StreamTypeCommon,
				ValidateFunc: validation.StringInSlice([]string{streams.StreamTypeCommon, streams.StreamTypeAdvanced}, false),
			},
			"partition_count": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"data_duration": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  24,
			},
			"data_schema": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceDisStreamCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	region := GetRegion(d, config)

	disClient, err := config.DisV2Client(region)
	if err != nil {
		return fmt.Errorf("Creating dis client failed, err=%s", err)
	}

	streamName := d.Get("name").(string)
	log.Printf("[DEBUG] Create dis stream streamName: %s", streamName)
	createOpts := streams.CreateOpts{
		StreamName:     streamName,
		StreamType:     d.Get("type").(string),
		PartitionCount: d.Get("partition_count").(int),
		DataDuration:   d.Get("data_duration").(int),
		DataSchema:     d.Get("data_schema").(string),
	}

	log.Printf("[DEBUG] Create dis stream using parameters: %+v", createOpts)
	_, err = streams.Create(disClient, createOpts)
	if err != nil {
		return fmt.Errorf("Create dis stream failed: %s", err)
	}

	d.SetId(streamName)

	return resourceDisStreamRead(d, meta)
}

func resourceDisStreamRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	region := GetRegion(d, config)

	disClient, err := config.DisV2Client(region)
	if err != nil {
		return fmt.Errorf("Creating dis client failed, err=%s", err)
	}

	streamName := d.Id()

	log.Printf("[DEBUG] Query dis stream using name: %s", streamName)

	getOpts := streams.GetOpts{}
	streamDetail, err := streams.Get(disClient, streamName, getOpts)
	if err != nil {
		return err
	}

	if streamDetail != nil {
		d.Set("name", streamDetail.StreamName)
		// d.Set("type", streamDetail.StreamType)
		// d.Set("partition_count", streamDetail.WritablePartitionCount) // TODO: Check that this match the actual partition_count (there is no partition_count attr)
		d.Set("data_duration", streamDetail.RetentionPeriod)
		d.Set("data_schema", streamDetail.DataSchema)
	}

	return nil
}

func resourceDisStreamDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	region := GetRegion(d, config)

	disClient, err := config.DisV2Client(region)
	if err != nil {
		return fmt.Errorf("Creating dis client failed, err=%s", err)
	}

	streamName := d.Get("name").(string)
	log.Printf("[DEBUG] Deleting dis stream streamName: %q", d.Id())

	result := streams.Delete(disClient, streamName)
	if result.Err != nil {
		return fmt.Errorf("Error deleting dis Stream %q, err=%s", d.Id(), result.Err)
	}

	return nil
}

func resourceDisStreamUpdate(d *schema.ResourceData, meta interface{}) error {
	return resourceDisStreamRead(d, meta)
}
