package flexibleengine

import (
	"context"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"log"
	"strings"

	"github.com/chnsz/golangsdk/openstack/obs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceObsBucketNotification() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceObsBucketNotificationCreate,
		UpdateContext: resourceObsBucketNotificationCreate,
		ReadContext:   resourceObsBucketNotificationRead,
		DeleteContext: resourceObsBucketNotificationDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceObsBucketNotificationImport,
		},

		Schema: map[string]*schema.Schema{
			"bucket": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"events": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"prefix": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"suffix": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"topic_urn": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceObsBucketNotificationCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
		return diag.Errorf("Error setting Notification Configuration of OBS bucket: %s, err: %s", bucket, err)
	}
	return nil
}

func buildTopicConfiguration(d *schema.ResourceData) []obs.TopicConfiguration {
	events := d.Get("events").([]interface{})
	eventTypes := make([]obs.EventType, 0, len(events))
	for _, value := range events {
		eventTypes = append(eventTypes, obs.EventType(value.(string)))
	}

	filterRules := make([]obs.FilterRule, 0, 2)
	m := map[string]string{
		"prefix": d.Get("prefix").(string),
		"suffix": d.Get("suffix").(string),
	}
	for k, v := range m {
		if len(v) == 0 {
			continue
		}
		filterRule := obs.FilterRule{
			Name:  k,
			Value: v,
		}
		filterRules = append(filterRules, filterRule)
	}
	return []obs.TopicConfiguration{
		{
			ID:          d.Get("name").(string),
			Topic:       d.Get("topic_urn").(string),
			Events:      eventTypes,
			FilterRules: filterRules,
		},
	}
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

	mErr := &multierror.Error{}
	for _, config := range output.TopicConfigurations {
		if config.ID != d.Get("name").(string) {
			continue
		}
		events := make([]string, 0, len(config.Events))
		for _, v := range config.Events {
			events = append(events, string(v))
		}

		for _, v := range config.FilterRules {
			mErr = multierror.Append(mErr, d.Set(v.Name, v.Value))
		}

		mErr = multierror.Append(mErr,
			d.Set("name", config.ID),
			d.Set("events", events),
			d.Set("topic_urn", config.Topic))
		break
	}
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

func resourceObsBucketNotificationImport(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.SplitN(d.Id(), "/", 2)
	if len(parts) != 2 {
		err := fmtp.Errorf("Invalid format specified for OBS notification. Format must be <bucket>/<name>")
		return nil, err
	}
	bucket := parts[0]
	name := parts[1]

	d.SetId(bucket)
	if err := d.Set("name", name); err != nil {
		return nil, fmtp.Errorf("Error setting OBS notification name, %s", err)
	}

	return []*schema.ResourceData{d}, nil
}
