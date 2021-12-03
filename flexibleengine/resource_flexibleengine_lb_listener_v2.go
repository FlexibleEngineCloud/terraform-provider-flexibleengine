package flexibleengine

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk/openstack/common/tags"
	"github.com/chnsz/golangsdk/openstack/networking/v2/extensions/lbaas_v2/listeners"
)

func resourceListenerV2() *schema.Resource {
	return &schema.Resource{
		Create: resourceListenerV2Create,
		Read:   resourceListenerV2Read,
		Update: resourceListenerV2Update,
		Delete: resourceListenerV2Delete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"loadbalancer_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"protocol": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"TCP", "UDP", "HTTP", "TERMINATED_HTTPS",
				}, false),
			},
			"protocol_port": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},

			"default_pool_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"http2_enable": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"default_tls_container_ref": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"sni_container_refs": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"tls_ciphers_policy": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"tags": tagsSchema(),

			"tenant_id": {
				Type:       schema.TypeString,
				Optional:   true,
				Computed:   true,
				Deprecated: "tenant_id is deprecated",
			},
			"admin_state_up": {
				Type:       schema.TypeBool,
				Default:    true,
				Optional:   true,
				Deprecated: "admin_state_up is deprecated",
			},
		},
	}
}

func resourceListenerV2Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	lbClient, err := config.ElbV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine ELB v2.0 client: %s", err)
	}

	http2Enable := d.Get("http2_enable").(bool)
	var sniContainerRefs []string
	if raw, ok := d.GetOk("sni_container_refs"); ok {
		for _, v := range raw.([]interface{}) {
			sniContainerRefs = append(sniContainerRefs, v.(string))
		}
	}
	createOpts := listeners.CreateOpts{
		Protocol:               listeners.Protocol(d.Get("protocol").(string)),
		ProtocolPort:           d.Get("protocol_port").(int),
		LoadbalancerID:         d.Get("loadbalancer_id").(string),
		Name:                   d.Get("name").(string),
		DefaultPoolID:          d.Get("default_pool_id").(string),
		Description:            d.Get("description").(string),
		DefaultTlsContainerRef: d.Get("default_tls_container_ref").(string),
		SniContainerRefs:       sniContainerRefs,
		TlsCiphersPolicy:       d.Get("tls_ciphers_policy").(string),
		Http2Enable:            &http2Enable,
	}

	log.Printf("[DEBUG] Create Options: %#v", createOpts)

	// Wait for LoadBalancer to become active before continuing
	lbID := createOpts.LoadbalancerID
	timeout := d.Timeout(schema.TimeoutCreate)
	err = waitForLBV2LoadBalancer(lbClient, lbID, "ACTIVE", nil, timeout)
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Attempting to create listener")
	var listener *listeners.Listener
	err = resource.Retry(timeout, func() *resource.RetryError {
		listener, err = listeners.Create(lbClient, createOpts).Extract()
		if err != nil {
			return checkForRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("Error creating listener: %s", err)
	}

	// Wait for LoadBalancer to become active again before continuing
	err = waitForLBV2LoadBalancer(lbClient, lbID, "ACTIVE", nil, timeout)
	if err != nil {
		return err
	}

	d.SetId(listener.ID)

	//set tags
	tagRaw := d.Get("tags").(map[string]interface{})
	if len(tagRaw) > 0 {
		taglist := expandResourceTags(tagRaw)
		if tagErr := tags.Create(lbClient, "listeners", listener.ID, taglist).ExtractErr(); tagErr != nil {
			return fmt.Errorf("Error setting tags of elb listener %s: %s", listener.ID, tagErr)
		}
	}

	return resourceListenerV2Read(d, meta)
}

func resourceListenerV2Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	region := GetRegion(d, config)
	lbClient, err := config.ElbV2Client(region)
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine ELB v2.0 client: %s", err)
	}

	listener, err := listeners.Get(lbClient, d.Id()).Extract()
	if err != nil {
		return CheckDeleted(d, err, "listener")
	}

	log.Printf("[DEBUG] Retrieved listener %s: %#v", d.Id(), listener)

	d.SetId(listener.ID)
	d.Set("region", region)
	d.Set("name", listener.Name)
	d.Set("protocol", listener.Protocol)
	d.Set("description", listener.Description)
	d.Set("protocol_port", listener.ProtocolPort)
	d.Set("http2_enable", listener.Http2Enable)
	d.Set("default_pool_id", listener.DefaultPoolID)
	if err := d.Set("sni_container_refs", listener.SniContainerRefs); err != nil {
		return fmt.Errorf("[DEBUG] Error saving sni_container_refs to state for FlexibleEngine listener (%s): %s", d.Id(), err)
	}
	d.Set("tls_ciphers_policy", listener.TlsCiphersPolicy)
	d.Set("default_tls_container_ref", listener.DefaultTlsContainerRef)

	// fetch tags
	if resourceTags, err := tags.Get(lbClient, "listeners", d.Id()).Extract(); err == nil {
		tagmap := tagsToMap(resourceTags.Tags)
		d.Set("tags", tagmap)
	} else {
		log.Printf("[WARN] fetching tags of elb listener failed: %s", err)
	}

	return nil
}

func resourceListenerV2Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	elbV2Client, err := config.ElbV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine ELB v2.0 client: %s", err)
	}

	var updateOpts listeners.UpdateOpts
	if d.HasChange("name") {
		updateOpts.Name = d.Get("name").(string)
	}
	if d.HasChange("description") {
		updateOpts.Description = d.Get("description").(string)
	}
	if d.HasChange("default_tls_container_ref") {
		updateOpts.DefaultTlsContainerRef = d.Get("default_tls_container_ref").(string)
	}
	if d.HasChange("sni_container_refs") {
		var sniContainerRefs []string
		if raw, ok := d.GetOk("sni_container_refs"); ok {
			for _, v := range raw.([]interface{}) {
				sniContainerRefs = append(sniContainerRefs, v.(string))
			}
		}
		updateOpts.SniContainerRefs = sniContainerRefs
	}
	if d.HasChange("tls_ciphers_policy") {
		updateOpts.TlsCiphersPolicy = d.Get("tls_ciphers_policy").(string)
	}
	if d.HasChange("http2_enable") {
		http2 := d.Get("http2_enable").(bool)
		updateOpts.Http2Enable = &http2
	}

	// Wait for LoadBalancer to become active before continuing
	lbID := d.Get("loadbalancer_id").(string)
	timeout := d.Timeout(schema.TimeoutUpdate)
	err = waitForLBV2LoadBalancer(elbV2Client, lbID, "ACTIVE", nil, timeout)
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Updating listener %s with options: %#v", d.Id(), updateOpts)
	err = resource.Retry(timeout, func() *resource.RetryError {
		_, err = listeners.Update(elbV2Client, d.Id(), updateOpts).Extract()
		if err != nil {
			return checkForRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("Error updating listener %s: %s", d.Id(), err)
	}

	// Wait for LoadBalancer to become active again before continuing
	err = waitForLBV2LoadBalancer(elbV2Client, lbID, "ACTIVE", nil, timeout)
	if err != nil {
		return err
	}

	// update tags
	if d.HasChange("tags") {
		tagErr := UpdateResourceTags(elbV2Client, d, "listeners", d.Id())
		if tagErr != nil {
			return fmt.Errorf("Error updating tags of elb listener:%s, err:%s", d.Id(), tagErr)
		}
	}

	return resourceListenerV2Read(d, meta)

}

func resourceListenerV2Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	elbV2Client, err := config.ElbV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine ELB v2.0 client: %s", err)
	}

	// Wait for LoadBalancer to become active before continuing
	lbID := d.Get("loadbalancer_id").(string)
	timeout := d.Timeout(schema.TimeoutDelete)
	err = waitForLBV2LoadBalancer(elbV2Client, lbID, "ACTIVE", nil, timeout)
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Deleting listener %s", d.Id())
	err = resource.Retry(timeout, func() *resource.RetryError {
		err = listeners.Delete(elbV2Client, d.Id()).ExtractErr()
		if err != nil {
			return checkForRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("Error deleting listener %s: %s", d.Id(), err)
	}

	// Wait for LoadBalancer to become active again before continuing
	err = waitForLBV2LoadBalancer(elbV2Client, lbID, "ACTIVE", nil, timeout)
	if err != nil {
		return err
	}

	// Wait for Listener to delete
	err = waitForLBV2Listener(elbV2Client, d.Id(), "DELETED", nil, timeout)
	if err != nil {
		return err
	}

	return nil
}
