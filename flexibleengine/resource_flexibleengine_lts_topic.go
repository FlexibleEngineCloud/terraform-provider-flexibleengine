package flexibleengine

import (
	"fmt"
	"log"
	"strings"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/lts/huawei/logstreams"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceLTSTopicV2() *schema.Resource {
	return &schema.Resource{
		Create: resourceLTSTopicV2Create,
		Read:   resourceLTSTopicV2Read,
		Delete: resourceLTSTopicV2Delete,
		Importer: &schema.ResourceImporter{
			State: resourceLTSTopicV2Import,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"group_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"topic_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"filter_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"index_enabled": {
				Type:       schema.TypeBool,
				Computed:   true,
				Deprecated: "it's deprecated",
			},
		},
	}
}

func resourceLTSTopicV2Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client, err := config.LtsV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine LTS client: %s", err)
	}

	groupID := d.Get("group_id").(string)
	createOpts := &logstreams.CreateOpts{
		LogStreamName: d.Get("topic_name").(string),
	}

	log.Printf("[DEBUG] Create Options: %#v", createOpts)
	topicCreate, err := logstreams.Create(client, groupID, createOpts).Extract()
	if err != nil {
		return fmt.Errorf("Error creating log topic: %s", err)
	}

	d.SetId(topicCreate.ID)
	return resourceLTSTopicV2Read(d, meta)
}

func resourceLTSTopicV2Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	region := GetRegion(d, config)
	client, err := config.LtsV2Client(region)
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine LTS client: %s", err)
	}

	topicID := d.Id()
	groupID := d.Get("group_id").(string)
	streams, err := logstreams.List(client, groupID).Extract()
	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault400); ok {
			log.Printf("[WARN] log group topic %s: the log group %s is gone", topicID, groupID)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error querying log topic %s: %s", topicID, err)
	}

	for _, stream := range streams.LogStreams {
		if stream.ID == topicID {
			log.Printf("[DEBUG] Retrieved log topic %s: %#v", topicID, stream)
			d.SetId(stream.ID)
			d.Set("region", region)
			d.Set("topic_name", stream.Name)
			d.Set("filter_count", stream.FilterCount)
			return nil
		}
	}

	log.Printf("[WARN] log group topic %s: resource is gone and will be removed in Terraform state", topicID)
	d.SetId("")

	return nil
}

func resourceLTSTopicV2Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client, err := config.LtsV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine LTS client: %s", err)
	}

	groupID := d.Get("group_id").(string)
	err = logstreams.Delete(client, groupID, d.Id()).ExtractErr()
	if err != nil {
		return CheckDeleted(d, err, "Error deleting log topic")
	}

	d.SetId("")
	return nil
}

func resourceLTSTopicV2Import(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	parts := strings.SplitN(d.Id(), "/", 2)
	if len(parts) != 2 {
		err := fmt.Errorf("Invalid format specified for LTS topic. Format must be <group id>/<topic id>")
		return nil, err
	}

	groupID := parts[0]
	topicID := parts[1]

	d.SetId(topicID)
	d.Set("group_id", groupID)

	return []*schema.ResourceData{d}, nil
}
