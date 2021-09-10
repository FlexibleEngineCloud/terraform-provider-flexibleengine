package flexibleengine

import (
	"fmt"
	"log"
	"time"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/mls/v1/instances"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceMlsInstanceV1() *schema.Resource {
	return &schema.Resource{
		Create: resourceMlsInstanceCreate,
		Read:   resourceMlsInstanceRead,
		Delete: resourceMlsInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
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
				ForceNew: true,
			},

			"version": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"network": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"vpc_id": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"subnet_id": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"security_group": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},
						"available_zone": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"public_ip": {
							Type:     schema.TypeList,
							Required: true,
							ForceNew: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"bind_type": {
										Type:     schema.TypeString,
										Required: true,
										ForceNew: true,
									},
									"eip_id": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},

			"flavor": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"agency": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"mrs_cluster": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"user_name": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"user_password": {
							Type:      schema.TypeString,
							Sensitive: true,
							Optional:  true,
							Computed:  true,
						},
					},
				},
			},

			"status": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"inner_endpoint": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"public_endpoint": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"created": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"updated": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceMlsInstanceNetwork(d *schema.ResourceData) instances.NetworkOpts {
	var network instances.NetworkOpts
	networkRaw := d.Get("network").([]interface{})
	log.Printf("[DEBUG] networkRaw: %+v", networkRaw)
	if len(networkRaw) == 1 {
		net := networkRaw[0].(map[string]interface{})
		publicipRaw := net["public_ip"].([]interface{})
		var publicip instances.PublicIPOpts
		publicip.BindType = publicipRaw[0].(map[string]interface{})["bind_type"].(string)

		network.VpcId = net["vpc_id"].(string)
		network.SubnetId = net["subnet_id"].(string)
		network.SecurityGroupId = net["security_group"].(string)
		network.AvailableZone = net["available_zone"].(string)
		network.PublicIP = publicip
	}
	log.Printf("[DEBUG] network: %+v", network)
	return network
}

func resourceMlsInstanceMrsCluster(d *schema.ResourceData) instances.MrsClusterOpts {
	var mrsCluster instances.MrsClusterOpts
	MrsClusterRaw := d.Get("mrs_cluster").([]interface{})
	log.Printf("[DEBUG] MrsClusterRaw: %+v", MrsClusterRaw)
	if len(MrsClusterRaw) == 1 {
		mrsCluster.Id = MrsClusterRaw[0].(map[string]interface{})["id"].(string)
		mrsCluster.UserName = MrsClusterRaw[0].(map[string]interface{})["user_name"].(string)
		mrsCluster.UserPassword = MrsClusterRaw[0].(map[string]interface{})["user_password"].(string)
	}
	log.Printf("[DEBUG] mrsCluster: %+v", mrsCluster)
	return mrsCluster
}

func MlsInstanceStateRefreshFunc(client *golangsdk.ServiceClient, instanceID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		instance, err := instances.Get(client, instanceID).Extract()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return instance, "DELETED", nil
			}
			return nil, "", err
		}

		return instance, instance.Status, nil
	}
}

func resourceMlsInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client, err := config.MlsV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine mls client: %s ", err)
	}

	createOpts := instances.CreateOpts{
		Name:       d.Get("name").(string),
		Version:    d.Get("version").(string),
		Network:    resourceMlsInstanceNetwork(d),
		Agency:     d.Get("agency").(string),
		FlavorRef:  d.Get("flavor").(string),
		MrsCluster: resourceMlsInstanceMrsCluster(d),
	}
	log.Printf("[DEBUG] Create Options: %#v", createOpts)

	instance, err := instances.Create(client, createOpts).Extract()
	if err != nil {
		return fmt.Errorf("Error getting instance from result: %s ", err)
	}
	log.Printf("[DEBUG] Create : instance %s: %#v", instance.ID, instance)

	d.SetId(instance.ID)
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"CREATING"},
		Target:     []string{"AVAILABLE"},
		Refresh:    MlsInstanceStateRefreshFunc(client, instance.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      30 * time.Second,
		MinTimeout: 15 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf(
			"Error waiting for instance (%s) to become ready: %s ",
			instance.ID, err)
	}

	if instance.ID != "" {
		return resourceMlsInstanceRead(d, meta)
	}
	return fmt.Errorf("Unexpected conversion error in resourceMlsInstanceCreate. ")
}

func resourceMlsInstanceRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client, err := config.MlsV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine mrs client: %s", err)
	}

	instanceID := d.Id()
	instance, err := instances.Get(client, instanceID).Extract()
	if err != nil {
		return CheckDeleted(d, err, "instance")
	}

	log.Printf("[DEBUG] Retrieved instance %s: %#v", instanceID, instance)

	d.Set("name", instance.Name)
	d.Set("version", instance.Version)
	d.Set("flavor", instance.FlavorRef)
	d.Set("status", instance.Status)

	publicIp := []map[string]interface{}{
		{
			"bind_type": instance.Network.PublicIP.BindType,
			"eip_id":    instance.Network.PublicIP.EipId,
		},
	}

	network := []map[string]interface{}{
		{
			"vpc_id":         instance.Network.VpcId,
			"subnet_id":      instance.Network.SubnetId,
			"security_group": instance.Network.SecurityGroupId,
			"available_zone": instance.Network.AvailableZone,
			"public_ip":      publicIp,
		},
	}
	log.Printf("[DEBUG] network: %+v", network)
	if err := d.Set("network", network); err != nil {
		return fmt.Errorf("[DEBUG] Error saving network to MLS instance (%s): %s", d.Id(), err)
	}

	mrsCluster := []map[string]interface{}{
		{
			"id": instance.MrsCluster.Id,
		},
	}
	log.Printf("[DEBUG] mrsCluster: %+v", mrsCluster)
	if err := d.Set("mrs_cluster", mrsCluster); err != nil {
		return fmt.Errorf("[DEBUG] Error saving mrs_cluster to MLS instance (%s): %s", d.Id(), err)
	}

	d.Set("created", instance.Created)
	d.Set("updated", instance.Updated)
	d.Set("inner_endpoint", instance.InnerEndPoint)
	d.Set("public_endpoint", instance.PublicEndPoint)

	d.Set("region", GetRegion(d, config))
	return nil
}

func resourceMlsInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client, err := config.MlsV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine mls client: %s ", err)
	}

	log.Printf("[DEBUG] Deleting Instance %s", d.Id())

	id := d.Id()
	result := instances.Delete(client, id)
	if result.Err != nil {
		return err
	}
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"AVAILABLE"},
		Target:     []string{"DELETED"},
		Refresh:    MlsInstanceStateRefreshFunc(client, id),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      15 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf(
			"Error waiting for instance (%s) to be deleted: %s ",
			id, err)
	}
	log.Printf("[DEBUG] Successfully deleted instance %s", id)
	return nil
}
