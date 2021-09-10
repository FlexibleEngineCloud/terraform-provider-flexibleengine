package flexibleengine

import (
	"fmt"
	"log"
	"strings"

	"github.com/chnsz/golangsdk/openstack/lts/v2/logtopics"
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
			"index_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
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
	client.ResourceBase = strings.Replace(client.ResourceBase, "/v2/", "/v2.0/", 1)

	groupID := d.Get("group_id").(string)
	createOpts := &logtopics.CreateOpts{
		LogTopicName: d.Get("topic_name").(string),
	}

	log.Printf("[DEBUG] Create Options: %#v", createOpts)
	topicCreate, err := logtopics.Create(client, groupID, createOpts).Extract()
	if err != nil {
		return fmt.Errorf("Error creating log topic: %s", err)
	}

	d.SetId(topicCreate.ID)
	return resourceLTSTopicV2Read(d, meta)
}

func resourceLTSTopicV2Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client, err := config.LtsV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine LTS client: %s", err)
	}
	client.ResourceBase = strings.Replace(client.ResourceBase, "/v2/", "/v2.0/", 1)

	groupID := d.Get("group_id").(string)
	topic, err := logtopics.Get(client, groupID, d.Id()).Extract()
	if err != nil {
		return CheckDeleted(d, err, "Error querying log topic")
	}

	log.Printf("[DEBUG] Retrieved log topic %s: %#v", d.Id(), topic)
	d.Set("topic_name", topic.Name)
	d.Set("index_enabled", topic.IndexEnabled)
	d.Set("region", GetRegion(d, config))

	return nil
}

func resourceLTSTopicV2Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client, err := config.LtsV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine LTS client: %s", err)
	}
	client.ResourceBase = strings.Replace(client.ResourceBase, "/v2/", "/v2.0/", 1)

	groupID := d.Get("group_id").(string)
	err = logtopics.Delete(client, groupID, d.Id()).ExtractErr()
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
