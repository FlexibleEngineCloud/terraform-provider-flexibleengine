package flexibleengine

import (
	"fmt"
	"log"
	"reflect"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/common/tags"
	listeners_v3 "github.com/chnsz/golangsdk/openstack/elb/v3/listeners"
	"github.com/chnsz/golangsdk/openstack/networking/v2/extensions/lbaas_v2/listeners"
)

func resourceListenerV2() *schema.Resource {
	return &schema.Resource{
		Create: resourceListenerCreate,
		Read:   resourceListenerRead,
		Update: resourceListenerUpdate,
		Delete: resourceListenerDelete,

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
			},

			"transparent_client_ip_enable": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
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

			"idle_timeout": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(0, 4000),
			},
			"request_timeout": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(1, 300),
			},
			"response_timeout": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(1, 300),
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

func resourceListenerCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	region := GetRegion(d, config)
	lbClient, err := config.ElbV2Client(region)
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine ELB v2.0 client: %s", err)
	}
	v3Client, err := config.ElbV3Client(region)
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine ELB v3 client: %s", err)
	}

	http2Enable := d.Get("http2_enable").(bool)
	var sniContainerRefs []string
	if raw, ok := d.GetOk("sni_container_refs"); ok {
		for _, v := range raw.([]interface{}) {
			sniContainerRefs = append(sniContainerRefs, v.(string))
		}
	}
	createOpts := listeners_v3.CreateOpts{
		Protocol:               listeners_v3.Protocol(d.Get("protocol").(string)),
		ProtocolPort:           d.Get("protocol_port").(int),
		LoadbalancerID:         d.Get("loadbalancer_id").(string),
		Name:                   d.Get("name").(string),
		DefaultPoolID:          d.Get("default_pool_id").(string),
		Description:            d.Get("description").(string),
		DefaultTlsContainerRef: d.Get("default_tls_container_ref").(string),
		TlsCiphersPolicy:       d.Get("tls_ciphers_policy").(string),
		SniContainerRefs:       sniContainerRefs,
		Http2Enable:            &http2Enable,
	}

	if transparentIP := d.Get("transparent_client_ip_enable").(bool); transparentIP {
		createOpts.TransparentClientIP = &transparentIP
	}
	if v1, ok := d.GetOk("idle_timeout"); ok {
		createOpts.KeepaliveTimeout = golangsdk.IntToPointer(v1.(int))
	}
	if v2, ok := d.GetOk("request_timeout"); ok {
		createOpts.ClientTimeout = golangsdk.IntToPointer(v2.(int))
	}
	if v3, ok := d.GetOk("response_timeout"); ok {
		createOpts.MemberTimeout = golangsdk.IntToPointer(v3.(int))
	}

	log.Printf("[DEBUG] Create v3 Options: %#v", createOpts)

	// Wait for LoadBalancer to become active before continuing
	lbID := createOpts.LoadbalancerID
	timeout := d.Timeout(schema.TimeoutCreate)
	err = waitForLBV2LoadBalancer(lbClient, lbID, "ACTIVE", nil, timeout)
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Attempting to create listener with v3 API")
	var listener *listeners_v3.Listener
	err = resource.Retry(timeout, func() *resource.RetryError {
		listener, err = listeners_v3.Create(v3Client, createOpts).Extract()
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

	return resourceListenerRead(d, meta)
}

func resourceListenerRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	region := GetRegion(d, config)
	lbClient, err := config.ElbV2Client(region)
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine ELB v2.0 client: %s", err)
	}
	v3Client, err := config.ElbV3Client(region)
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine ELB v3 client: %s", err)
	}

	listener, err := listeners_v3.Get(v3Client, d.Id()).Extract()
	if err != nil {
		return CheckDeleted(d, err, "listener")
	}

	log.Printf("[DEBUG] Retrieved listener %s: %#v", d.Id(), listener)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("name", listener.Name),
		d.Set("description", listener.Description),
		d.Set("protocol", listener.Protocol),
		d.Set("protocol_port", listener.ProtocolPort),
		d.Set("http2_enable", listener.Http2Enable),
		d.Set("transparent_client_ip_enable", listener.TransparentClientIP),
		d.Set("default_pool_id", listener.DefaultPoolID),
		d.Set("sni_container_refs", listener.SniContainerRefs),
		d.Set("tls_ciphers_policy", listener.TlsCiphersPolicy),
		d.Set("default_tls_container_ref", listener.DefaultTlsContainerRef),
		d.Set("idle_timeout", listener.KeepaliveTimeout),
		d.Set("request_timeout", listener.ClientTimeout),
		d.Set("response_timeout", listener.MemberTimeout),
	)
	if mErr.ErrorOrNil() != nil {
		return mErr
	}

	// fetch tags
	if resourceTags, err := tags.Get(lbClient, "listeners", d.Id()).Extract(); err == nil {
		tagmap := tagsToMap(resourceTags.Tags)
		d.Set("tags", tagmap)
	} else {
		log.Printf("[WARN] fetching tags of elb listener failed: %s", err)
	}

	return nil
}

func resourceListenerUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	region := GetRegion(d, config)
	lbClient, err := config.ElbV2Client(region)
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine ELB v2.0 client: %s", err)
	}
	v3Client, err := config.ElbV3Client(region)
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine ELB v3 client: %s", err)
	}

	// Wait for LoadBalancer to become active before continuing
	lbID := d.Get("loadbalancer_id").(string)
	timeout := d.Timeout(schema.TimeoutUpdate)
	err = waitForLBV2LoadBalancer(lbClient, lbID, "ACTIVE", nil, timeout)
	if err != nil {
		return err
	}

	var updateOpts listeners_v3.UpdateOpts
	if d.HasChange("name") {
		updateOpts.Name = d.Get("name").(string)
	}
	if d.HasChange("description") {
		desc := d.Get("description").(string)
		updateOpts.Description = &desc
	}
	if d.HasChange("sni_container_refs") {
		var sniContainerRefs []string
		if raw, ok := d.GetOk("sni_container_refs"); ok {
			for _, v := range raw.([]interface{}) {
				sniContainerRefs = append(sniContainerRefs, v.(string))
			}
		}
		updateOpts.SniContainerRefs = &sniContainerRefs
	}
	if d.HasChange("default_tls_container_ref") {
		tlsContainerRef := d.Get("default_tls_container_ref").(string)
		updateOpts.DefaultTlsContainerRef = &tlsContainerRef
	}
	if d.HasChange("tls_ciphers_policy") {
		tlsPolicy := d.Get("tls_ciphers_policy").(string)
		updateOpts.TlsCiphersPolicy = &tlsPolicy
	}
	if d.HasChange("http2_enable") {
		http2 := d.Get("http2_enable").(bool)
		updateOpts.Http2Enable = &http2
	}
	if d.HasChange("transparent_client_ip_enable") {
		transparentIPEnable := d.Get("transparent_client_ip_enable").(bool)
		updateOpts.TransparentClientIP = &transparentIPEnable
	}

	if d.HasChange("idle_timeout") {
		updateOpts.KeepaliveTimeout = golangsdk.IntToPointer(d.Get("idle_timeout").(int))
	}
	if d.HasChanges("request_timeout", "response_timeout") {
		updateOpts.ClientTimeout = golangsdk.IntToPointer(d.Get("request_timeout").(int))
		updateOpts.MemberTimeout = golangsdk.IntToPointer(d.Get("response_timeout").(int))
	}

	if !reflect.DeepEqual(updateOpts, listeners_v3.UpdateOpts{}) {
		log.Printf("[DEBUG] Updating listener %s with options: %#v", d.Id(), updateOpts)

		err = resource.Retry(timeout, func() *resource.RetryError {
			_, err = listeners_v3.Update(v3Client, d.Id(), updateOpts).Extract()
			if err != nil {
				return checkForRetryableError(err)
			}
			return nil
		})

		if err != nil {
			return fmt.Errorf("Error updating listener %s: %s", d.Id(), err)
		}

		// Wait for LoadBalancer to become active again before continuing
		err = waitForLBV2LoadBalancer(lbClient, lbID, "ACTIVE", nil, timeout)
		if err != nil {
			return err
		}
	}

	// update tags
	if d.HasChange("tags") {
		tagErr := UpdateResourceTags(lbClient, d, "listeners", d.Id())
		if tagErr != nil {
			return fmt.Errorf("Error updating tags of elb listener:%s, err:%s", d.Id(), tagErr)
		}
	}

	return resourceListenerRead(d, meta)

}

func resourceListenerDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	lbClient, err := config.ElbV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine ELB v2.0 client: %s", err)
	}

	// Wait for LoadBalancer to become active before continuing
	lbID := d.Get("loadbalancer_id").(string)
	timeout := d.Timeout(schema.TimeoutDelete)
	err = waitForLBV2LoadBalancer(lbClient, lbID, "ACTIVE", nil, timeout)
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Deleting listener %s", d.Id())
	err = resource.Retry(timeout, func() *resource.RetryError {
		err = listeners.Delete(lbClient, d.Id()).ExtractErr()
		if err != nil {
			return checkForRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("Error deleting listener %s: %s", d.Id(), err)
	}

	// Wait for LoadBalancer to become active again before continuing
	err = waitForLBV2LoadBalancer(lbClient, lbID, "ACTIVE", nil, timeout)
	if err != nil {
		return err
	}

	// Wait for Listener to delete
	err = waitForLBV2Listener(lbClient, d.Id(), "DELETED", nil, timeout)
	if err != nil {
		return err
	}

	return nil
}
