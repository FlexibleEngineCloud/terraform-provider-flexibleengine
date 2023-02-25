package flexibleengine

import (
	"context"
	"log"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk/openstack/obs"
)

func resourceObsBucketNotifications() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceObsBucketNotificationCreate,
		UpdateContext: resourceObsBucketNotificationCreate,
		ReadContext:   resourceObsBucketNotificationRead,
		DeleteContext: resourceObsBucketNotificationDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"bucket": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"notifications": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"topic_urn": {
							Type:     schema.TypeString,
							Required: true,
						},
						"events": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
								ValidateFunc: validation.StringInSlice([]string{
									"ObjectCreated:*", "ObjectCreated:Put", "ObjectCreated:Post", "ObjectCreated:Copy",
									"ObjectCreated:CompleteMultipartUpload", "ObjectRemoved:*", "ObjectRemoved:Delete",
									"ObjectRemoved:DeleteMarkerCreated",
								}, false),
							},
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"prefix": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"suffix": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func resourceObsBucketNotificationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	obsClient, err := config.ObjectStorageClient(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OBS client: %s", err)
	}

	bucket := d.Get("bucket").(string)
	notificationOpt := obs.SetBucketNotificationInput{}
	notificationOpt.Bucket = bucket

	// set notification
	configurations := buildTopicConfiguration(d)
	notificationOpt.TopicConfigurations = configurations
	_, err = obsClient.SetBucketNotification(&notificationOpt)
	if err != nil {
		return diag.Errorf("Error setting Notification Configuration of OBS bucket %s, err: %s", bucket, err)
	}
	d.SetId(bucket)
	return resourceObsBucketNotificationRead(ctx, d, meta)
}

func buildTopicConfiguration(d *schema.ResourceData) []obs.TopicConfiguration {
	notifications := d.Get("notifications").([]interface{})

	configurations := make([]obs.TopicConfiguration, 0, len(notifications))
	for _, notification := range notifications {
		notificationMap := notification.(map[string]interface{})
		configuration := obs.TopicConfiguration{
			ID:          notificationMap["name"].(string),
			Topic:       notificationMap["topic_urn"].(string),
			Events:      buildEvents(notificationMap["events"].([]interface{})),
			FilterRules: buildFilterRules(notificationMap),
		}
		configurations = append(configurations, configuration)
	}
	return configurations
}

func buildFilterRules(notificationMap map[string]interface{}) []obs.FilterRule {
	var filterRules []obs.FilterRule
	for k, v := range notificationMap {
		if k == "prefix" || k == "suffix" {
			filterRule := obs.FilterRule{
				Name:  k,
				Value: v.(string),
			}
			filterRules = append(filterRules, filterRule)
		}
	}
	return filterRules
}

func buildEvents(events []interface{}) []obs.EventType {
	eventTypes := make([]obs.EventType, 0, len(events))
	for _, value := range events {
		eventTypes = append(eventTypes, obs.EventType(value.(string)))
	}
	return eventTypes
}

func resourceObsBucketNotificationRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	obsClient, err := config.ObjectStorageClient(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OBS client: %s", err)
	}

	bucket := d.Id()
	output, err := obsClient.GetBucketNotification(bucket)
	if err != nil {
		return diag.Errorf("Error getting OBS Notification Configuration: %s", err)
	}

	mErr := multierror.Append(nil, d.Set("bucket", bucket))
	notifications := make([]map[string]interface{}, 0, len(output.TopicConfigurations))
	for _, config := range output.TopicConfigurations {
		events := make([]string, 0, len(config.Events))
		for _, v := range config.Events {
			events = append(events, string(v))
		}

		notificationMap := make(map[string]interface{})
		notificationMap["name"] = config.ID
		notificationMap["topic_urn"] = config.Topic
		notificationMap["events"] = events
		for _, v := range config.FilterRules {
			if v.Name == "prefix" || v.Name == "suffix" {
				notificationMap[v.Name] = v.Value
			}
		}
		notifications = append(notifications, notificationMap)
	}
	mErr = multierror.Append(mErr, d.Set("notifications", notifications))
	if mErr.ErrorOrNil() != nil {
		return diag.Errorf("Error saving bucket notification %s: %s", d.Id(), mErr)
	}
	return nil
}

func resourceObsBucketNotificationDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	obsClient, err := config.ObjectStorageClient(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OBS client: %s", err)
	}

	bucket := d.Id()
	log.Printf("[DEBUG] delete Notification Configuration of OBS bucket %s", bucket)

	notificationOpt := obs.SetBucketNotificationInput{}
	notificationOpt.Bucket = bucket
	_, err = obsClient.SetBucketNotification(&notificationOpt)
	if err != nil {
		return diag.Errorf("Error deleting Notification Configuration of OBS bucket: %s, err: %s", bucket, err)
	}
	return nil
}
