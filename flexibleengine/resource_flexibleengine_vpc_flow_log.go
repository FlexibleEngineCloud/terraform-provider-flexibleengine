package flexibleengine

import (
	"fmt"
	"log"
	"time"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/networking/v1/flowlogs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceVpcFlowLogV1() *schema.Resource {
	return &schema.Resource{
		Create: resourceVpcFlowLogV1Create,
		Read:   resourceVpcFlowLogV1Read,
		Update: resourceVpcFlowLogV1Update,
		Delete: resourceVpcFlowLogV1Delete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"resource_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "port",
				ValidateFunc: validation.StringInSlice([]string{
					"port", "vpc", "network",
				}, true),
			},
			"resource_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"traffic_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "all",
				ValidateFunc: validation.StringInSlice([]string{
					"all", "accept", "reject",
				}, true),
			},
			"log_group_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"log_topic_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceVpcFlowLogV1Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	vpcClient, err := config.NetworkingV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("error creating FlexibleEngine vpc client: %s", err)
	}

	createOpts := flowlogs.CreateOpts{
		Name:         d.Get("name").(string),
		Description:  d.Get("description").(string),
		ResourceType: d.Get("resource_type").(string),
		ResourceID:   d.Get("resource_id").(string),
		TrafficType:  d.Get("traffic_type").(string),
		LogGroupID:   d.Get("log_group_id").(string),
		LogTopicID:   d.Get("log_topic_id").(string),
	}

	log.Printf("[DEBUG] Create VPC Flow Log Options: %#v", createOpts)
	fl, err := flowlogs.Create(vpcClient, createOpts).Extract()
	if err != nil {
		return fmt.Errorf("error creating FlexibleEngine VPC flow log: %s", err)
	}

	d.SetId(fl.ID)
	return resourceVpcFlowLogV1Read(d, config)
}

func resourceVpcFlowLogV1Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	vpcClient, err := config.NetworkingV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("error creating FlexibleEngine vpc client: %s", err)
	}

	fl, err := flowlogs.Get(vpcClient, d.Id()).Extract()
	if err != nil {
		return CheckDeleted(d, err, "error retrieving flowlog")
	}

	d.Set("name", fl.Name)
	d.Set("description", fl.Description)
	d.Set("resource_type", fl.ResourceType)
	d.Set("resource_id", fl.ResourceID)
	d.Set("traffic_type", fl.TrafficType)
	d.Set("log_group_id", fl.LogGroupID)
	d.Set("log_topic_id", fl.LogTopicID)
	d.Set("status", fl.Status)
	d.Set("region", GetRegion(d, config))

	return nil
}

func resourceVpcFlowLogV1Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	vpcClient, err := config.NetworkingV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("error creating FlexibleEngine vpc client: %s", err)
	}

	if d.HasChanges("name", "description") {
		updateOpts := flowlogs.UpdateOpts{
			Name:        d.Get("name").(string),
			Description: d.Get("description").(string),
		}

		_, err = flowlogs.Update(vpcClient, d.Id(), updateOpts).Extract()
		if err != nil {
			return fmt.Errorf("error updating FlexibleEngine VPC flow log: %s", err)
		}
	}

	return resourceVpcFlowLogV1Read(d, meta)
}

func resourceVpcFlowLogV1Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	vpcClient, err := config.NetworkingV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("error creating FlexibleEngine vpc client: %s", err)
	}

	err = flowlogs.Delete(vpcClient, d.Id()).ExtractErr()
	if err != nil {
		// ignore ErrDefault404
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			log.Printf("[INFO] Successfully deleted FlexibleEngine vpc flow log %s", d.Id())
			return nil
		}
		return err
	}

	d.SetId("")
	return nil
}
