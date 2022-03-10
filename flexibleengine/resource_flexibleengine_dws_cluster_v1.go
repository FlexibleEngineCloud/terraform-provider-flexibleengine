package flexibleengine

import (
	"fmt"
	"log"
	"time"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/dws/v1/cluster"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceDWSClusterV1() *schema.Resource {
	return &schema.Resource{
		Create: resourceDWSClusterV1Create,
		Read:   resourceDWSClusterV1Read,
		Delete: resourceDWSClusterV1Delete,
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
			},

			"availability_zone": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"node_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"number_of_node": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},

			"user_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"user_pwd": {
				Type:      schema.TypeString,
				Required:  true,
				ForceNew:  true,
				Sensitive: true,
			},

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
			"security_group_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"port": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"public_ip": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"eip_id": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"public_bind_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"endpoints": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"connect_info": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"jdbc_url": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"public_endpoints": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"jdbc_url": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"public_connect_info": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"version": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"sub_status": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"task_status": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"created": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"updated": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func getPublicIP(d *schema.ResourceData) *cluster.PublicIpOpts {
	ip, ok := d.Get("public_ip").([]interface{})
	if !ok || len(ip) == 0 {
		return nil
	}
	info := ip[0].(map[string]interface{})
	return &cluster.PublicIpOpts{
		EipID:          info["eip_id"].(string),
		PublicBindType: info["public_bind_type"].(string),
	}
}

func resourceDWSClusterV1Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client, err := config.DwsV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine client: %s", err)
	}

	opts := cluster.CreateOpts{
		Name:             d.Get("name").(string),
		NumberOfNode:     d.Get("number_of_node").(int),
		AvailabilityZone: d.Get("availability_zone").(string),
		SubnetID:         d.Get("subnet_id").(string),
		UserPwd:          d.Get("user_pwd").(string),
		SecurityGroupID:  d.Get("security_group_id").(string),
		PublicIp:         getPublicIP(d),
		NodeType:         d.Get("node_type").(string),
		VpcID:            d.Get("vpc_id").(string),
		UserName:         d.Get("user_name").(string),
		Port:             d.Get("port").(int),
	}
	log.Printf("[DEBUG] Create DWS-Cluster Options: %#v", opts)

	c, err := cluster.Create(client, opts)
	if err != nil {
		return fmt.Errorf("Error creating DWS-Cluster: %s", err)
	}

	clusterID := c.Cluster.Id
	d.SetId(clusterID)

	// Wait for Cluster to become active before continuing
	stateConf := &resource.StateChangeConf{
		Target:     []string{"AVAILABLE"},
		Pending:    []string{"CREATING"},
		Refresh:    getDWSCluster(client, clusterID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Error waiting for DWS cluster %s(%s) to become AVAILABLE, error=%s",
			opts.Name, clusterID, err)
	}

	return resourceDWSClusterV1Read(d, meta)
}

func resourceDWSClusterV1Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client, err := config.DwsV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine client: %s", err)
	}

	r, err := cluster.Get(client, d.Id())
	if err != nil {
		return CheckDeleted(d, err, "DWS-Cluster")
	}
	log.Printf("[DEBUG] Retrieved DWS-Cluster %s: %#v", d.Id(), r)

	m, err := convertStructToMap(r, map[string]string{"endPoints": "endpoints"})
	if err != nil {
		return fmt.Errorf("Error converting struct to map, err=%s", err)
	}

	d.Set("region", GetRegion(d, config))

	d.Set("status", r.Status)
	d.Set("sub_status", r.SubStatus)
	d.Set("updated", r.Updated)
	d.Set("endpoints", []interface{}{m["endpoints"]})
	d.Set("name", r.Name)
	d.Set("number_of_node", r.NumberOfNode)
	d.Set("availability_zone", r.AvailabilityZone)
	d.Set("subnet_id", r.SubnetID)
	d.Set("public_endpoints", []interface{}{m["public_endpoints"]})
	d.Set("created", r.Created)
	d.Set("security_group_id", r.SecurityGroupID)
	d.Set("port", r.Port)
	d.Set("node_type", r.NodeType)
	d.Set("version", r.Version)
	d.Set("public_ip", []interface{}{m["public_ip"]})
	d.Set("vpc_id", r.VpcID)
	d.Set("task_status", r.TaskStatus)
	d.Set("user_name", r.UserName)

	return nil
}

func resourceDWSClusterV1Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client, err := config.DwsV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine client: %s", err)
	}

	rID := d.Id()
	log.Printf("[DEBUG] Deleting DWS-Cluster %s", rID)

	timeout := d.Timeout(schema.TimeoutDelete)
	err = resource.Retry(timeout, func() *resource.RetryError {
		err := cluster.Delete(client, rID).ExtractErr()
		if err != nil {
			return checkForRetryableError(err)
		}
		return nil
	})
	if err != nil {
		if isResourceNotFound(err) {
			log.Printf("[INFO] deleting an unavailable DWS-Cluster: %s", rID)
			return nil
		}
		return fmt.Errorf("Error deleting DWS-Cluster %s: %s", rID, err)
	}

	return nil
}

func getDWSCluster(client *golangsdk.ServiceClient, clusterID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		r, err := cluster.Get(client, clusterID)
		if err != nil {
			return nil, "", err
		}
		if r.FailedReasons != nil {
			return nil, r.Status, fmt.Errorf(r.FailedReasons.ErrorMsg)
		}

		return r, r.Status, nil
	}
}
