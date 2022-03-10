package flexibleengine

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/networking/v2/extensions/elbaas/loadbalancer_elbs"
)

func resourceELoadBalancer() *schema.Resource {
	return &schema.Resource{
		Create: resourceELoadBalancerCreate,
		Read:   resourceELoadBalancerRead,
		Update: resourceELoadBalancerUpdate,
		Delete: resourceELoadBalancerDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},

		DeprecationMessage: "It has been deprecated, using enhanced load balancer instead",
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"bandwidth": {
				Type:     schema.TypeInt,
				Optional: true,
			},

			"vip_subnet_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"az": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"security_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"vip_address": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"admin_state_up": {
				Type:     schema.TypeBool,
				Default:  true,
				Optional: true,
			},

			"tenantid": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
		},
	}
}

func resourceELoadBalancerCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client, err := otcV1Client(config, GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine networking client: %s", err)
	}

	adminStateUp := d.Get("admin_state_up").(bool)
	createOpts := loadbalancer_elbs.CreateOpts{
		Name:            d.Get("name").(string),
		Description:     d.Get("description").(string),
		VpcID:           d.Get("vpc_id").(string),
		Bandwidth:       d.Get("bandwidth").(int),
		Type:            d.Get("type").(string),
		AdminStateUp:    &adminStateUp,
		VipSubnetID:     d.Get("vip_subnet_id").(string),
		AZ:              d.Get("az").(string),
		SecurityGroupID: d.Get("security_group_id").(string),
		VipAddress:      d.Get("vip_address").(string),
		TenantID:        d.Get("tenantid").(string),
	}

	log.Printf("[DEBUG] Create Options: %#v", createOpts)

	job, err := loadbalancer_elbs.Create(client, createOpts).ExtractJobResponse()
	if err != nil {
		return err
	}

	if err := golangsdk.WaitForJobSuccess(client, job.URI, int(d.Timeout(schema.TimeoutCreate)/time.Second)); err != nil {
		return err
	}

	entity, err := golangsdk.GetJobEntity(client, job.URI, "elb")

	if mlb, ok := entity.(map[string]interface{}); ok {
		if vid, ok := mlb["id"]; ok {
			if id, ok := vid.(string); ok {
				// If all has been successful, set the ID on the resource, return Read of it
				d.SetId(id)
				return resourceELoadBalancerRead(d, meta)
			}
		}
	}

	return fmt.Errorf("Unexpected conversion error in resourceELoadBalancerCreate.")
}

func resourceELoadBalancerRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client, err := otcV1Client(config, GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine networking client: %s", err)
	}

	lb, err := loadbalancer_elbs.Get(client, d.Id()).Extract()
	if err != nil {
		return CheckDeleted(d, err, "loadbalancer")
	}

	log.Printf("[DEBUG] Retrieved loadbalancer %s: %#v", d.Id(), lb)

	d.Set("name", lb.Name)
	d.Set("description", lb.Description)
	d.Set("vpc_id", lb.VpcID)
	d.Set("bandwidth", lb.Bandwidth)
	d.Set("type", lb.Type)
	basu := false
	// Can be 0 (not up) or 2 (frozen)
	if lb.AdminStateUp == 1 {
		basu = true
	}
	d.Set("admin_state_up", basu)
	d.Set("vip_subnet_id", lb.VipSubnetID)
	d.Set("vip_address", lb.VipAddress)
	secgroup_Id := lb.SecurityGroupID
	// flexibleengine will return empty string of security_group_id when lb type is External
	if lb.Type == "External" {
		if secgroup_Id == "" {
			secgroup_Id = d.Get("security_group_id").(string)
		}
	}
	d.Set("security_group_id", secgroup_Id)
	d.Set("region", GetRegion(d, config))

	return nil
}

func resourceELoadBalancerUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client, err := otcV1Client(config, GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine networking client: %s", err)
	}

	var updateOpts loadbalancer_elbs.UpdateOpts
	if d.HasChange("name") {
		updateOpts.Name = d.Get("name").(string)
	}
	if d.HasChange("description") {
		updateOpts.Description = d.Get("description").(string)
	}
	if d.HasChange("bandwidth") {
		updateOpts.Bandwidth = d.Get("bandwidth").(int)
	}
	if d.HasChange("admin_state_up") {
		asu := d.Get("admin_state_up").(bool)
		updateOpts.AdminStateUp = &asu
	}

	log.Printf("[DEBUG] Updating loadbalancer %s with options: %#v", d.Id(), updateOpts)
	job, err := loadbalancer_elbs.Update(client, d.Id(), updateOpts).ExtractJobResponse()
	if err := golangsdk.WaitForJobSuccess(client, job.URI, int(d.Timeout(schema.TimeoutUpdate)/time.Second)); err != nil {
		return err
	}

	return resourceELoadBalancerRead(d, meta)
}

func resourceELoadBalancerDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client, err := otcV1Client(config, GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine networking client: %s", err)
	}

	id := d.Id()
	log.Printf("[DEBUG] Deleting loadbalancer %s", d.Id())
	job, err := loadbalancer_elbs.Delete(client, id, false).ExtractJobResponse()
	if err != nil {
		return err
	}

	log.Printf("Waiting for loadbalancer %s to delete", id)

	if err := golangsdk.WaitForJobSuccess(client, job.URI, int(d.Timeout(schema.TimeoutDelete)/time.Second)); err != nil {
		return err
	}

	log.Printf("Successfully deleted loadbalancer %s", id)
	return nil
}
