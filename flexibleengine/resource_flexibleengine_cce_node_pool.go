package flexibleengine

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/cce/v3/nodepools"
	"github.com/chnsz/golangsdk/openstack/cce/v3/nodes"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceCCENodePool() *schema.Resource {
	return &schema.Resource{
		Create: resourceCCENodePoolCreate,
		Read:   resourceCCENodePoolRead,
		Update: resourceCCENodePoolUpdate,
		Delete: resourceCCENodePoolDelete,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				var err error = nil
				arr := strings.Split(d.Id(), "/")
				if len(arr) == 2 {
					cluster_id := arr[0]
					node_pool_id := arr[1]
					d.Set("cluster_id", cluster_id)
					d.SetId(node_pool_id)
				} else {
					err = fmt.Errorf("[ERROR] Missing argument, must be of the form: 'cluster_id/node_pool_id'")
				}

				return []*schema.ResourceData{d}, err
			},
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
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
			"initial_node_count": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"cluster_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"flavor_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"root_volume": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"size": {
							Type:     schema.TypeInt,
							Required: true,
							ForceNew: true,
						},
						"volumetype": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"extend_params": {
							Type:     schema.TypeMap,
							Optional: true,
							ForceNew: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					}},
			},
			"data_volumes": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"size": {
							Type:     schema.TypeInt,
							Required: true,
							ForceNew: true,
						},
						"volumetype": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"kms_key_id": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"extend_params": {
							Type:     schema.TypeMap,
							Optional: true,
							ForceNew: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					}},
			},
			"availability_zone": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "random",
			},
			"os": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"key_pair": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ExactlyOneOf: []string{"password", "key_pair"},
			},
			"password": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Sensitive:    true,
				ExactlyOneOf: []string{"password", "key_pair"},
			},
			"labels": { //(k8s_tags)
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"taints": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:     schema.TypeString,
							Required: true,
						},
						"value": {
							Type:     schema.TypeString,
							Required: true,
						},
						"effect": {
							Type:     schema.TypeString,
							Required: true,
						},
					}},
			},
			"tags": tagsSchema(),
			"max_pods": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"preinstall": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				StateFunc: func(v interface{}) string {
					switch v.(type) {
					case string:
						return installScriptHashSum(v.(string))
					default:
						return ""
					}
				},
			},
			"postinstall": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				StateFunc: func(v interface{}) string {
					switch v.(type) {
					case string:
						return installScriptHashSum(v.(string))
					default:
						return ""
					}
				},
			},
			"extend_param": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"subnet_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"scall_enable": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"min_node_count": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"max_node_count": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"scale_down_cooldown_time": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"priority": {
				Type:     schema.TypeInt,
				Optional: true,
			},

			"billing_mode": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func buildCCENodePoolLoginSpec(d *schema.ResourceData) nodes.LoginSpec {
	var loginSpec nodes.LoginSpec

	if v1, ok := d.GetOk("key_pair"); ok {
		loginSpec.SshKey = v1.(string)
	} else if v2, ok := d.GetOk("password"); ok {
		loginSpec.UserPassword = nodes.UserPassword{
			Username: "root",
			Password: v2.(string),
		}
	}

	return loginSpec
}

func resourceCCENodePoolCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	nodePoolClient, err := config.CceV3Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating Flexibleengine CCE Node Pool client: %s", err)
	}

	// wait for the cce cluster to become available
	clusterid := d.Get("cluster_id").(string)
	stateCluster := &resource.StateChangeConf{
		Target:     []string{"Available"},
		Refresh:    waitForClusterAvailable(nodePoolClient, clusterid),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      5 * time.Second,
		MinTimeout: 5 * time.Second,
	}
	if _, err = stateCluster.WaitForState(); err != nil {
		return fmt.Errorf("CCE Cluster %s is inactive: %s", clusterid, err)
	}

	initialNodeCount := d.Get("initial_node_count").(int)
	createOpts := nodepools.CreateOpts{
		Kind:       "NodePool",
		ApiVersion: "v3",
		Metadata: nodepools.CreateMetaData{
			Name: d.Get("name").(string),
		},
		Spec: nodepools.CreateSpec{
			Type: d.Get("type").(string),
			NodeTemplate: nodes.Spec{
				Flavor:      d.Get("flavor_id").(string),
				Az:          d.Get("availability_zone").(string),
				Os:          d.Get("os").(string),
				RootVolume:  resourceCCERootVolume(d),
				DataVolumes: resourceCCEDataVolume(d),
				K8sTags:     resourceCCENodeK8sTags(d),
				BillingMode: 0,
				Count:       1,
				NodeNicSpec: nodes.NodeNicSpec{
					PrimaryNic: nodes.PrimaryNic{
						SubnetId: d.Get("subnet_id").(string),
					},
				},
				ExtendParam: resourceCCEExtendParam(d),
				Taints:      resourceCCETaint(d),
				UserTags:    resourceCCENodeUserTags(d),
			},
			Autoscaling: nodepools.AutoscalingSpec{
				Enable:                d.Get("scall_enable").(bool),
				MinNodeCount:          d.Get("min_node_count").(int),
				MaxNodeCount:          d.Get("max_node_count").(int),
				ScaleDownCooldownTime: d.Get("scale_down_cooldown_time").(int),
				Priority:              d.Get("priority").(int),
			},
			InitialNodeCount: &initialNodeCount,
		},
	}

	log.Printf("[DEBUG] Create Options: %#v", createOpts)
	// Add loginSpec here so it wouldn't go in the above log entry
	createOpts.Spec.NodeTemplate.Login = buildCCENodePoolLoginSpec(d)

	s, err := nodepools.Create(nodePoolClient, clusterid, createOpts).Extract()
	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault403); ok {
			retryNode, err := recursiveNodePoolCreate(nodePoolClient, createOpts, clusterid, 403)
			if err == "fail" {
				return fmt.Errorf("Error creating Flexibleengine Node Pool")
			}
			s = retryNode
		} else {
			return fmt.Errorf("Error creating Flexibleengine Node Pool: %s", err)
		}
	}

	if len(s.Metadata.Id) == 0 {
		return fmt.Errorf("Error fetching CreateNodePool id")
	}

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"Synchronizing", "Synchronized"},
		Target:       []string{""},
		Refresh:      waitForCceNodePoolActive(nodePoolClient, clusterid, s.Metadata.Id),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        120 * time.Second,
		PollInterval: 20 * time.Second,
	}
	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Error creating Flexibleengine CCE Node Pool: %s", err)
	}

	log.Printf("[DEBUG] Create node pool: %v", s)

	d.SetId(s.Metadata.Id)
	return resourceCCENodePoolRead(d, meta)
}

func resourceCCENodePoolRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	nodePoolClient, err := config.CceV3Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating Flexibleengine CCE Node Pool client: %s", err)
	}
	clusterid := d.Get("cluster_id").(string)
	s, err := nodepools.Get(nodePoolClient, clusterid, d.Id()).Extract()

	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error retrieving Flexibleengine Node Pool: %s", err)
	}

	mErr := multierror.Append(
		d.Set("name", s.Metadata.Name),
		d.Set("flavor_id", s.Spec.NodeTemplate.Flavor),
		d.Set("availability_zone", s.Spec.NodeTemplate.Az),
		d.Set("os", s.Spec.NodeTemplate.Os),
		d.Set("billing_mode", s.Spec.NodeTemplate.BillingMode),
		d.Set("key_pair", s.Spec.NodeTemplate.Login.SshKey),
		d.Set("initial_node_count", s.Spec.InitialNodeCount),
		d.Set("scall_enable", s.Spec.Autoscaling.Enable),
		d.Set("min_node_count", s.Spec.Autoscaling.MinNodeCount),
		d.Set("max_node_count", s.Spec.Autoscaling.MaxNodeCount),
		d.Set("scale_down_cooldown_time", s.Spec.Autoscaling.ScaleDownCooldownTime),
		d.Set("priority", s.Spec.Autoscaling.Priority),
		d.Set("type", s.Spec.Type),
		d.Set("status", s.Status.Phase),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return err
	}

	rootVolume := expandResourceCCERootVolume(s.Spec.NodeTemplate)
	d.Set("root_volume", rootVolume)
	if err := d.Set("root_volume", rootVolume); err != nil {
		return fmt.Errorf("Error saving root volume of cce node pool %s: %s", d.Id(), err)
	}

	volumes := expandResourceCCEDataVolumes(s.Spec.NodeTemplate)
	if err := d.Set("data_volumes", volumes); err != nil {
		return fmt.Errorf("Error saving data volumes of cce node pool %s: %s", d.Id(), err)
	}

	nodeTaints := expandResourceCCETaints(s.Spec.NodeTemplate)
	if err := d.Set("taints", nodeTaints); err != nil {
		return fmt.Errorf("Error saving taints of cce node %s: %s", d.Id(), err)
	}

	labels := expandResourceCCEK8sTags(s.Spec.NodeTemplate)
	if err := d.Set("labels", labels); err != nil {
		return fmt.Errorf("Error saving labels/k8stags of cce node pool %s: %s", d.Id(), err)
	}

	tagmap := tagsToMap(s.Spec.NodeTemplate.UserTags)
	// ignore "CCE-Dynamic-Provisioning-Node"
	delete(tagmap, "CCE-Dynamic-Provisioning-Node")
	if err := d.Set("tags", tagmap); err != nil {
		return fmt.Errorf("Error saving tags to state for CCE Node Pool(%s): %s", d.Id(), err)
	}

	return nil
}

func resourceCCENodePoolUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	nodePoolClient, err := config.CceV3Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating Flexibleengine CCE client: %s", err)
	}

	initialNodeCount := d.Get("initial_node_count").(int)
	updateOpts := nodepools.UpdateOpts{
		Kind:       "NodePool",
		ApiVersion: "v3",
		Metadata: nodepools.UpdateMetaData{
			Name: d.Get("name").(string),
		},
		Spec: nodepools.UpdateSpec{
			InitialNodeCount: &initialNodeCount,
			Autoscaling: nodepools.AutoscalingSpec{
				Enable:                d.Get("scall_enable").(bool),
				MinNodeCount:          d.Get("min_node_count").(int),
				MaxNodeCount:          d.Get("max_node_count").(int),
				ScaleDownCooldownTime: d.Get("scale_down_cooldown_time").(int),
				Priority:              d.Get("priority").(int),
			},
			NodeTemplate: nodes.Spec{
				Flavor:      d.Get("flavor_id").(string),
				Az:          d.Get("availability_zone").(string),
				Login:       buildCCENodePoolLoginSpec(d),
				RootVolume:  resourceCCERootVolume(d),
				DataVolumes: resourceCCEDataVolume(d),
				Count:       1,
				K8sTags:     resourceCCENodeK8sTags(d),
				UserTags:    resourceCCENodeUserTags(d),
				Taints:      resourceCCETaint(d),
				ExtendParam: resourceCCEExtendParam(d),
			},
			Type: d.Get("type").(string),
		},
	}

	clusterid := d.Get("cluster_id").(string)
	_, err = nodepools.Update(nodePoolClient, clusterid, d.Id(), updateOpts).Extract()
	if err != nil {
		return fmt.Errorf("Error updating Flexibleengine Node Node Pool: %s", err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"Synchronizing", "Synchronized"},
		Target:     []string{""},
		Refresh:    waitForCceNodePoolActive(nodePoolClient, clusterid, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      60 * time.Second,
		MinTimeout: 10 * time.Second,
	}
	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Error creating Flexibleengine CCE Node Pool: %s", err)
	}

	return resourceCCENodePoolRead(d, meta)
}

func resourceCCENodePoolDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	nodePoolClient, err := config.CceV3Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating Flexibleengine CCE client: %s", err)
	}
	clusterid := d.Get("cluster_id").(string)
	err = nodepools.Delete(nodePoolClient, clusterid, d.Id()).ExtractErr()
	if err != nil {
		return fmt.Errorf("Error deleting Flexibleengine CCE Node Pool: %s", err)
	}
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"Deleting"},
		Target:       []string{"Deleted"},
		Refresh:      waitForCceNodePoolDelete(nodePoolClient, clusterid, d.Id()),
		Timeout:      d.Timeout(schema.TimeoutDelete),
		Delay:        60 * time.Second,
		PollInterval: 20 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Error deleting Flexibleengine CCE Node Pool: %s", err)
	}

	d.SetId("")
	return nil
}

func waitForCceNodePoolActive(cceClient *golangsdk.ServiceClient, clusterId, nodePoolId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		n, err := nodepools.Get(cceClient, clusterId, nodePoolId).Extract()
		if err != nil {
			return nil, "", err
		}
		return n, n.Status.Phase, nil
	}
}

func waitForCceNodePoolDelete(cceClient *golangsdk.ServiceClient, clusterId, nodePoolId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		log.Printf("[DEBUG] Attempting to delete Flexibleengine CCE Node Pool %s.\n", nodePoolId)

		r, err := nodepools.Get(cceClient, clusterId, nodePoolId).Extract()

		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				log.Printf("[DEBUG] Successfully deleted Flexibleengine CCE Node Pool %s", nodePoolId)
				return r, "Deleted", nil
			}
			return r, "Deleting", err
		}

		log.Printf("[DEBUG] Flexibleengine CCE Node Pool %s still available.\n", nodePoolId)
		return r, r.Status.Phase, nil
	}
}

func recursiveNodePoolCreate(cceClient *golangsdk.ServiceClient, opts nodepools.CreateOptsBuilder, ClusterID string, errCode int) (*nodepools.NodePool, string) {
	if errCode == 403 {
		stateCluster := &resource.StateChangeConf{
			Target:     []string{"Available"},
			Refresh:    waitForClusterAvailable(cceClient, ClusterID),
			Timeout:    15 * time.Minute,
			Delay:      15 * time.Second,
			MinTimeout: 3 * time.Second,
		}
		_, stateErr := stateCluster.WaitForState()
		if stateErr != nil {
			log.Printf("[INFO] Cluster Unavailable %s.\n", stateErr)
		}
		s, err := nodepools.Create(cceClient, ClusterID, opts).Extract()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault403); ok {
				return recursiveNodePoolCreate(cceClient, opts, ClusterID, 403)
			} else {
				return s, "fail"
			}
		} else {
			return s, "success"
		}
	}
	return nil, "fail"
}
