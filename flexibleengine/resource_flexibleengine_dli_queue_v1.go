package flexibleengine

import (
	"fmt"
	"log"
	"math"
	"regexp"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk/openstack/common/tags"
	"github.com/chnsz/golangsdk/openstack/dli/v1/queues"
)

var regexp4Name = regexp.MustCompile(`^[a-z0-9_]{1,128}$`)

const CU_16 = 16
const RESOURCE_MODE_SHARED, RESOURCE_MODE_EXCLUSIVE = 0, 1
const QUEUE_TYPE_SQL, QUEUE_TYPE_GENERAL = "sql", "general"

const (
	actionRestart  = "restart"
	actionScaleOut = "scale_out"
	actionScaleIn  = "scale_in"
)

func ResourceDliQueueV1() *schema.Resource {
	return &schema.Resource{
		Create: resourceDliQueueCreate,
		Read:   resourceDliQueueRead,
		Update: resourceDliQueueUpdate,
		Delete: resourceDliQueueDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Default: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringMatch(regexp4Name,
					"only contain digits, lower letters, and underscores (_)"),
			},

			"queue_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      "sql",
				ValidateFunc: validation.StringInSlice([]string{QUEUE_TYPE_SQL, QUEUE_TYPE_GENERAL}, false),
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"cu_count": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validCuCount,
			},

			"resource_mode": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntInSlice([]int{RESOURCE_MODE_SHARED, RESOURCE_MODE_EXCLUSIVE}),
			},

			"tags": {
				Type: schema.TypeMap,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
				ForceNew: true,
			},

			"create_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func resourceDliQueueCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	region := GetRegion(d, config)

	dliClient, err := config.DliV1Client(region)
	if err != nil {
		return fmt.Errorf("creating dli client failed: %s", err)
	}

	queueName := d.Get("name").(string)

	log.Printf("[DEBUG] create dli queues queueName: %s", queueName)
	createOpts := queues.CreateOpts{
		QueueName:    queueName,
		QueueType:    d.Get("queue_type").(string),
		Description:  d.Get("description").(string),
		CuCount:      d.Get("cu_count").(int),
		ResourceMode: d.Get("resource_mode").(int),
		Tags:         assembleTagsFromRecource("tags", d),
	}

	log.Printf("[DEBUG] create dli queues using paramaters: %+v", createOpts)
	createResult := queues.Create(dliClient, createOpts)
	if createResult.Err != nil {
		return fmt.Errorf("create dli queues failed: %s", createResult.Err)
	}

	//query queue detail,trriger read to refresh the state
	d.SetId(queueName)
	// This is a workaround to avoid issue: the queue is assigning, which is not available
	time.Sleep(120 * time.Second) //lintignore:R018

	return resourceDliQueueRead(d, meta)
}

func assembleTagsFromRecource(key string, d *schema.ResourceData) []tags.ResourceTag {
	if v, ok := d.GetOk(key); ok {
		tagRaw := v.(map[string]interface{})
		taglist := expandResourceTags(tagRaw)
		return taglist
	}
	return nil

}

func resourceDliQueueRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	region := GetRegion(d, config)

	dliClient, err := config.DliV1Client(region)

	if err != nil {
		return fmt.Errorf("creating dli client failed, err=%s", err)
	}

	queueName := d.Id()

	queryOpts := queues.ListOpts{
		QueueType: d.Get("queue_type").(string),
	}

	log.Printf("[DEBUG] query dli queues using paramaters: %+v", queryOpts)

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
		log.Printf("[DEBUG]The detail of queue from SDK:%+v", queueDetail)

		d.Set("name", queueDetail.QueueName)
		d.Set("queue_type", queueDetail.QueueType)
		d.Set("description", queueDetail.Description)
		d.Set("cu_count", queueDetail.CuCount)
		d.Set("resource_mode", queueDetail.ResourceMode)
		d.Set("create_time", queueDetail.CreateTime)
	}

	return nil
}

func filterByQueueName(body interface{}, queueName string) (r *queues.Queue, err error) {
	if queueList, ok := body.(*queues.ListResult); ok {
		log.Printf("[DEBUG]The list of queue from SDK:%+v", queueList)

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

func resourceDliQueueDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	region := GetRegion(d, config)

	client, err := config.DliV1Client(region)
	if err != nil {
		return fmt.Errorf("deleting dli client failed, err=%s", err)
	}

	queueName := d.Get("name").(string)
	log.Printf("[DEBUG] Deleting dli Queue %q", d.Id())

	result := queues.Delete(client, queueName)
	if result.Err != nil {
		return fmt.Errorf("error deleting dli Queue %q, err=%s", d.Id(), result.Err)
	}

	return nil
}

// support cu_count scaling
func resourceDliQueueUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client, err := config.DliV1Client(config.GetRegion(d))
	if err != nil {
		return fmt.Errorf("error creating DliV1Client: %s", err)
	}
	opt := queues.ActionOpts{
		QueueName: d.Id(),
	}
	if d.HasChange("cu_count") {
		oldValue, newValue := d.GetChange("cu_count")
		cuChange := newValue.(int) - oldValue.(int)

		opt.CuCount = int(math.Abs(float64(cuChange)))
		opt.Action = buildScaleActionParam(oldValue.(int), newValue.(int))

		log.Printf("[DEBUG]DLI queue Update Option: %#v", opt)
		result := queues.ScaleOrRestart(client, opt)
		if result.Err != nil {
			return fmt.Errorf("update dli queues failed,queueName=%s,error:%s", opt.QueueName, result.Err)
		}

		updateStateConf := &resource.StateChangeConf{
			Pending: []string{fmt.Sprintf("%d", oldValue)},
			Target:  []string{fmt.Sprintf("%d", newValue)},
			Refresh: func() (interface{}, string, error) {
				getResult := queues.Get(client, d.Id())
				queueDetail := getResult.Body.(*queues.Queue4Get)
				return getResult, fmt.Sprintf("%d", queueDetail.CuCount), nil
			},
			Timeout:      d.Timeout(schema.TimeoutUpdate),
			Delay:        30 * time.Second,
			PollInterval: 20 * time.Second,
		}
		_, err = updateStateConf.WaitForState()
		if err != nil {
			return fmt.Errorf("error waiting for dli.queue (%s) to be scale: %s", d.Id(), err)
		}

	}

	return resourceDliQueueRead(d, meta)
}

func buildScaleActionParam(oldValue, newValue int) string {
	if oldValue > newValue {
		return actionScaleIn
	} else {
		return actionScaleOut
	}
}

func validCuCount(val interface{}, key string) (warns []string, errs []error) {
	diviNum := 16
	warns, errs = validation.IntAtLeast(diviNum)(val, key)
	if len(errs) > 0 {
		return warns, errs
	}
	return validation.IntDivisibleBy(diviNum)(val, key)
}
