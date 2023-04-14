package flexibleengine

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/smn/v2/subscriptions"
)

func resourceSubscription() *schema.Resource {
	return &schema.Resource{
		Create: resourceSubscriptionCreate,
		Read:   resourceSubscriptionRead,
		Delete: resourceSubscriptionDelete,

		Schema: map[string]*schema.Schema{
			"topic_urn": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"endpoint": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"protocol": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					return ValidateStringList(v, k, []string{"email", "sms", "http", "https"})
				},
			},
			"remark": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"subscription_urn": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"owner": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
		},
	}
}

func resourceSubscriptionCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client, err := config.SmnV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine smn client: %s", err)
	}
	topicUrn := d.Get("topic_urn").(string)
	createOpts := subscriptions.CreateOps{
		Endpoint: d.Get("endpoint").(string),
		Protocol: d.Get("protocol").(string),
		Remark:   d.Get("remark").(string),
	}
	log.Printf("[DEBUG] Create Options: %#v", createOpts)

	subscription, err := subscriptions.Create(client, createOpts, topicUrn).Extract()
	if err != nil {
		return fmt.Errorf("Error getting subscription from result: %s", err)
	}
	log.Printf("[DEBUG] Create : subscription.SubscriptionUrn %s", subscription.SubscriptionUrn)
	if subscription.SubscriptionUrn != "" {
		d.SetId(subscription.SubscriptionUrn)
		d.Set("subscription_urn", subscription.SubscriptionUrn)
		return resourceSubscriptionRead(d, meta)
	}

	return fmt.Errorf("Unexpected conversion error in resourceSubscriptionCreate.")
}

func resourceSubscriptionDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client, err := config.SmnV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine smn client: %s", err)
	}

	log.Printf("[DEBUG] Deleting subscription %s", d.Id())

	id := d.Id()
	result := subscriptions.Delete(client, id)
	if result.Err != nil {
		return err
	}

	log.Printf("[DEBUG] Successfully deleted subscription %s", id)
	return nil
}

func resourceSubscriptionRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client, err := config.SmnV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine smn client: %s", err)
	}

	subscriptionslist, err := subscriptions.List(client).Extract()
	if err != nil {
		return fmt.Errorf("error fetching the list of subscriptions: %s", err)
	}

	var targetSubscription *subscriptions.SubscriptionGet
	id := d.Id()
	for i := range subscriptionslist {
		if subscriptionslist[i].SubscriptionUrn == id {
			targetSubscription = &subscriptionslist[i]
			break
		}
	}

	if targetSubscription == nil {
		return CheckDeleted(d, golangsdk.ErrDefault404{}, "subscription")
	}

	log.Printf("[DEBUG] fetching subscription: %#v", targetSubscription)
	d.Set("topic_urn", targetSubscription.TopicUrn)
	d.Set("endpoint", targetSubscription.Endpoint)
	d.Set("protocol", targetSubscription.Protocol)
	d.Set("subscription_urn", targetSubscription.SubscriptionUrn)
	d.Set("owner", targetSubscription.Owner)
	d.Set("remark", targetSubscription.Remark)
	d.Set("status", targetSubscription.Status)
	return nil
}
