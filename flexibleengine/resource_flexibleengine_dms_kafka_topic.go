package flexibleengine

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/chnsz/golangsdk/openstack/dms/v1/topics"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// resourceDmsKafkaTopic implements the resource of "flexibleengine_dms_kafka_topic"
func resourceDmsKafkaTopic() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDmsKafkaTopicCreate,
		ReadContext:   resourceDmsKafkaTopicRead,
		DeleteContext: resourceDmsKafkaTopicDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceDmsKafkaTopicImport,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"partitions": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"replicas": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"aging_time": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"sync_replication": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"sync_flushing": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
		},
	}
}

func resourceDmsKafkaTopicCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	dmsv1Client, err := config.DmsV1Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating FlexibleEngine DMS client: %s", err)
	}

	createOpts := &topics.CreateOps{
		Name:             d.Get("name").(string),
		Partition:        d.Get("partitions").(int),
		Replication:      d.Get("replicas").(int),
		RetentionTime:    d.Get("aging_time").(int),
		SyncReplication:  d.Get("sync_replication").(bool),
		SyncMessageFlush: d.Get("sync_flushing").(bool),
	}

	log.Printf("[DEBUG] Create Options: %#v", createOpts)
	instanceID := d.Get("instance_id").(string)
	v, err := topics.Create(dmsv1Client, instanceID, createOpts).Extract()
	if err != nil {
		return diag.Errorf("error creating FlexibleEngine DMS kafka topic: %s", err)
	}

	// use topic name as the resource ID
	d.SetId(v.Name)
	return resourceDmsKafkaTopicRead(ctx, d, meta)
}

func resourceDmsKafkaTopicRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	dmsv1Client, err := config.DmsV1Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating FlexibleEngine DMS client: %s", err)
	}

	instanceID := d.Get("instance_id").(string)
	allTopics, err := topics.List(dmsv1Client, instanceID).Extract()
	if err != nil {
		return CheckDeletedDiag(d, err, "DMS kafka topic")
	}

	topicID := d.Id()
	var found *topics.Topic
	for _, item := range allTopics {
		if item.Name == topicID {
			found = &item
			break
		}
	}

	if found == nil {
		log.Printf("[WARN] the DMS kafka topic %s does not exist in instance %s", d.Id(), instanceID)
		d.SetId("")
		return nil
	}

	log.Printf("[DEBUG] DMS kafka topic %s: %+v", d.Id(), found)

	mErr := multierror.Append(nil,
		d.Set("region", config.GetRegion(d)),
		d.Set("name", found.Name),
		d.Set("partitions", found.Partition),
		d.Set("replicas", found.Replication),
		d.Set("aging_time", found.RetentionTime),
		d.Set("sync_replication", found.SyncReplication),
		d.Set("sync_flushing", found.SyncMessageFlush),
	)
	if mErr.ErrorOrNil() != nil {
		return diag.FromErr(mErr)
	}

	return nil
}

func resourceDmsKafkaTopicDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	dmsv1Client, err := config.DmsV1Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating FlexibleEngine DMS client: %s", err)
	}

	topicID := d.Id()
	instanceID := d.Get("instance_id").(string)
	response, err := topics.Delete(dmsv1Client, instanceID, []string{topicID}).Extract()
	if err != nil {
		return diag.Errorf("error deleting DMS kafka topic: %s", err)
	}

	var success bool
	for _, item := range response {
		if item.Name == topicID {
			success = item.Success
			break
		}
	}
	if !success {
		return diag.Errorf("error deleting DMS kafka topic")
	}

	d.SetId("")
	return nil
}

// resourceDmsKafkaTopicImport query the rules from FlexibleEngine and imports them to Terraform.
// It is a common function in waf and is also called by other rule resources.
func resourceDmsKafkaTopicImport(_ context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData,
	error) {
	parts := strings.SplitN(d.Id(), "/", 2)
	if len(parts) != 2 {
		err := fmt.Errorf("Invalid format specified for DMS kafka topic. Format must be <instance id>/<topic name>")
		return nil, err
	}

	instanceID := parts[0]
	topicID := parts[1]

	d.SetId(topicID)
	d.Set("instance_id", instanceID)

	return []*schema.ResourceData{d}, nil
}
