package flexibleengine

import (
	"fmt"
	"log"
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/huaweicloud/golangsdk/openstack/common/tags"
	"github.com/huaweicloud/golangsdk/openstack/dli/v1/queues"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var regexp4Name = regexp.MustCompile(`^[a-z0-9_]+$`)

func ResourceDliQueueV1() *schema.Resource {
	return &schema.Resource{
		Create: resourceDliQueueCreate,
		Read:   resourceDliQueueV1Read,
		Delete: resourceDliQueueV1Delete,

		Schema: map[string]*schema.Schema{
			"queue_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(string)
					if !regexp4Name.MatchString(v) {
						errs = append(errs, fmt.Errorf("%q can contain only digits, lower letters, and underscores (_) ", key))
					}
					return
				},
			},

			"queue_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "sql",
			},

			"description": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
				ForceNew: true,
			},

			"cu_count": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},

			"charging_mode": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
				Default:  1,
			},

			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "0",
				StateFunc: func(val interface{}) string {
					if val.(string) == "" {
						return "0"
					} else {
						return val.(string)
					}
				},
			},

			"platform": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "x86_64",
			},

			"resource_mode": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},

			"labels": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"feature": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "basic",
			},

			"tags": {
				Type: schema.TypeMap,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
				ForceNew: true,
			},

			"cidr_in_vpc": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceDliQueueCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	region := GetRegion(d, config)

	dliClient, err := config.Config.DliV1Client(region)
	if err != nil {
		return fmt.Errorf("creating dli client failed: %s", err)
	}

	queueName := d.Get("queue_name").(string)

	log.Printf("[DEBUG] create dli queues queueName: %s", queueName)
	createOpts := queues.CreateOpts{
		QueueName:           queueName,
		QueueType:           d.Get("queue_type").(string),
		Description:         d.Get("description").(string),
		CuCount:             d.Get("cu_count").(int),
		ChargingMode:        d.Get("charging_mode").(int),
		EnterpriseProjectId: d.Get("enterprise_project_id").(string),
		Platform:            d.Get("platform").(string),
		ResourceMode:        d.Get("resource_mode").(int),
		Feature:             d.Get("feature").(string),
		Labels:              assembleMapFromRecource("Labels", d),
		Tags:                assembleTagsFromRecource("tags", d),
	}

	log.Printf("[DEBUG] create dli queues using paramaters: %+v", createOpts)
	createResult := queues.Create(dliClient, createOpts)
	if createResult.Err != nil {
		return fmt.Errorf("create dli queues failed: %s", createResult.Err)
	}

	//query queue detail,trriger read to refresh the state
	d.SetId(queueName)

	return resourceDliQueueV1Read(d, meta)
}

func assembleMapFromRecource(key string, d *schema.ResourceData) map[string]string {
	m := make(map[string]string)

	if v, ok := d.GetOk(key); ok {
		for key, val := range v.(map[string]interface{}) {
			m[key] = val.(string)
		}
	}

	return m
}

func assembleTagsFromRecource(key string, d *schema.ResourceData) []tags.ResourceTag {
	if v, ok := d.GetOk(key); ok {
		tagRaw := v.(map[string]interface{})
		taglist := utils.ExpandResourceTags(tagRaw)
		return taglist
	}
	return nil

}

func resourceDliQueueV1Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	region := GetRegion(d, config)

	dliClient, err := config.Config.DliV1Client(region)

	if err != nil {
		return fmt.Errorf("creating sdk client failed, err=%s", err)
	}

	queueName := d.Get("queue_name").(string)

	queryOpts := queues.ListOpts{
		QueueType: d.Get("queue_type").(string),
	}

	log.Printf("[DEBUG] create dli queues using paramaters: %+v", queryOpts)

	queryAllResult := queues.List(dliClient, queryOpts)
	if queryAllResult.Err != nil {
		return fmt.Errorf("query queues failed: %s", queryAllResult.Err)
	}

	//filter by queue_name
	queueDetail, err := filterByQueueName(queryAllResult.Body, queueName)
	if err != nil {
		return err
	}

	if queueDetail != nil {
		log.Printf("[debug]The detail of queue from SDK:%+v", queueDetail)

		d.Set("queue_name", queueDetail.QueueName)
		d.Set("queue_type", queueDetail.QueueType)
		d.Set("description", queueDetail.Description)
		d.Set("cu_count", queueDetail.CuCount)
		d.Set("charging_mode", queueDetail.ChargingMode)
		if queueDetail.EnterpriseProjectId != "" {
			d.Set("enterprise_project_id", queueDetail.EnterpriseProjectId)
		}

		d.Set("platform", queueDetail.Platform)
		d.Set("resource_mode", queueDetail.ResourceMode)
		d.Set("feature", queueDetail.Feature)
	}

	return nil
}

func filterByQueueName(body interface{}, queueName string) (r *queues.Queue, err error) {
	if queueList, ok := body.(*queues.ListResult); ok {
		log.Printf("[debug]The list of queue from SDK:%+v", queueList)

		for _, v := range queueList.Queues {
			if v.QueueName == queueName {
				return &v, nil
			}
		}
		return nil, nil

	} else {
		return nil, fmt.Errorf("sdk-client response type is wrong, expect type:*queues.ListResult,acutal Type:%T",
			body)
	}

}

func resourceDliQueueV1Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	region := GetRegion(d, config)

	client, err := config.Config.DliV1Client(region)
	if err != nil {
		return fmt.Errorf("creating sdk client failed, err=%s", err)
	}

	queueName := d.Get("queue_name").(string)
	log.Printf("[DEBUG] Deleting dli Queue %q", d.Id())

	result := queues.Delete(client, queueName)
	if result.Err != nil {
		return fmt.Errorf("Error deleting dli Queue %q, err=%s", d.Id(), result.Err)
	}

	return nil
}
