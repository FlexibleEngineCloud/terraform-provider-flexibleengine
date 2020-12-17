package flexibleengine

import (
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/openstack/cce/v3/clusters"
	"github.com/huaweicloud/golangsdk/openstack/cce/v3/nodes"
	"github.com/huaweicloud/golangsdk/openstack/common/tags"
)

func resourceCCENodeV3() *schema.Resource {
	return &schema.Resource{
		Create: resourceCCENodeV3Create,
		Read:   resourceCCENodeV3Read,
		Update: resourceCCENodeV3Update,
		Delete: resourceCCENodeV3Delete,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				arr := strings.Split(d.Id(), "/")
				cluster_id := arr[0]
				node_id := arr[1]
				d.Set("cluster_id", cluster_id)
				d.SetId(node_id)

				return []*schema.ResourceData{d}, nil
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
			"cluster_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"labels": { //(k8s_tags)
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
			},
			"tags": {
				Type:     schema.TypeMap,
				Optional: true,
			},
			"server_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"annotations": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
			},
			"flavor_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"availability_zone": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"os": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"key_pair": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"root_volume": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"size": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"volumetype": {
							Type:     schema.TypeString,
							Required: true,
						},
						"extend_param": {
							Type:     schema.TypeString,
							Optional: true,
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
						},
						"volumetype": {
							Type:     schema.TypeString,
							Required: true,
						},
						"extend_param": {
							Type:     schema.TypeString,
							Optional: true,
						},
					}},
			},
			"taints": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
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
			"eip_ids": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
				Computed: true,
			},
			"eip_count": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"iptype": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"bandwidth_charge_mode": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"sharetype": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"bandwidth_size": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"billing_mode": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"extend_param_charging_mode": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"ecs_performance_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"order_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"product_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"max_pods": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"public_key": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
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
			"private_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"public_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceCCENodeAnnotationsV2(d *schema.ResourceData) map[string]string {
	m := make(map[string]string)
	for key, val := range d.Get("annotations").(map[string]interface{}) {
		m[key] = val.(string)
	}
	return m
}
func resourceCCEDataVolume(d *schema.ResourceData) []nodes.VolumeSpec {
	volumeRaw := d.Get("data_volumes").([]interface{})
	volumes := make([]nodes.VolumeSpec, len(volumeRaw))
	for i, raw := range volumeRaw {
		rawMap := raw.(map[string]interface{})
		volumes[i] = nodes.VolumeSpec{
			Size:        rawMap["size"].(int),
			VolumeType:  rawMap["volumetype"].(string),
			ExtendParam: rawMap["extend_param"].(string),
		}
	}
	return volumes
}

func resourceCCETaint(d *schema.ResourceData) []nodes.TaintSpec {
	taintRaw := d.Get("taints").([]interface{})
	taints := make([]nodes.TaintSpec, len(taintRaw))
	for i, raw := range taintRaw {
		rawMap := raw.(map[string]interface{})
		taints[i] = nodes.TaintSpec{
			Key:    rawMap["key"].(string),
			Value:  rawMap["value"].(string),
			Effect: rawMap["effect"].(string),
		}
	}
	return taints
}

func resourceCCERootVolume(d *schema.ResourceData) nodes.VolumeSpec {
	var nics nodes.VolumeSpec
	nicsRaw := d.Get("root_volume").([]interface{})
	if len(nicsRaw) == 1 {
		nics.Size = nicsRaw[0].(map[string]interface{})["size"].(int)
		nics.VolumeType = nicsRaw[0].(map[string]interface{})["volumetype"].(string)
		nics.ExtendParam = nicsRaw[0].(map[string]interface{})["extend_param"].(string)
	}
	return nics
}
func resourceCCEEipIDs(d *schema.ResourceData) []string {
	rawID := d.Get("eip_ids").(*schema.Set)
	id := make([]string, rawID.Len())
	for i, raw := range rawID.List() {
		id[i] = raw.(string)
	}
	return id
}

func resourceCCENodeK8sTags(d *schema.ResourceData) map[string]string {
	m := make(map[string]string)
	for key, val := range d.Get("labels").(map[string]interface{}) {
		m[key] = val.(string)
	}
	return m
}

func resourceCCENodeUserTags(d *schema.ResourceData) []tags.ResourceTag {
	tagRaw := d.Get("tags").(map[string]interface{})
	return expandResourceTags(tagRaw)
}

func resourceCCENodeV3Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	nodeClient, err := config.cceV3Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating flexibleengine CCE Node client: %s", err)
	}

	var base64PreInstall, base64PostInstall string
	if v, ok := d.GetOk("preinstall"); ok {
		base64PreInstall = installScriptEncode(v.(string))
	}
	if v, ok := d.GetOk("postinstall"); ok {
		base64PostInstall = installScriptEncode(v.(string))
	}

	createOpts := nodes.CreateOpts{
		Kind:       "Node",
		ApiVersion: "v3",
		Metadata: nodes.CreateMetaData{
			Name:        d.Get("name").(string),
			Annotations: resourceCCENodeAnnotationsV2(d),
		},
		Spec: nodes.Spec{
			Flavor:      d.Get("flavor_id").(string),
			Az:          d.Get("availability_zone").(string),
			Os:          d.Get("os").(string),
			Login:       nodes.LoginSpec{SshKey: d.Get("key_pair").(string)},
			RootVolume:  resourceCCERootVolume(d),
			DataVolumes: resourceCCEDataVolume(d),
			UserTags:    resourceCCENodeUserTags(d),
			K8sTags:     resourceCCENodeK8sTags(d),
			Taints:      resourceCCETaint(d),
			PublicIP: nodes.PublicIPSpec{
				Ids:   resourceCCEEipIDs(d),
				Count: d.Get("eip_count").(int),
				Eip: nodes.EipSpec{
					IpType: d.Get("iptype").(string),
					Bandwidth: nodes.BandwidthOpts{
						ChargeMode: d.Get("bandwidth_charge_mode").(string),
						Size:       d.Get("bandwidth_size").(int),
						ShareType:  d.Get("sharetype").(string),
					},
				},
			},
			BillingMode: d.Get("billing_mode").(int),
			Count:       1,
			ExtendParam: nodes.ExtendParam{
				ChargingMode:       d.Get("extend_param_charging_mode").(int),
				EcsPerformanceType: d.Get("ecs_performance_type").(string),
				MaxPods:            d.Get("max_pods").(int),
				OrderID:            d.Get("order_id").(string),
				ProductID:          d.Get("product_id").(string),
				PublicKey:          d.Get("public_key").(string),
				PreInstall:         base64PreInstall,
				PostInstall:        base64PostInstall,
			},
		},
	}

	clusterid := d.Get("cluster_id").(string)
	stateCluster := &resource.StateChangeConf{
		Target:     []string{"Available"},
		Refresh:    waitForClusterAvailable(nodeClient, clusterid),
		Timeout:    15 * time.Minute,
		Delay:      15 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	_, err = stateCluster.WaitForState()

	log.Printf("[DEBUG] Create Options: %#v", createOpts)
	s, err := nodes.Create(nodeClient, clusterid, createOpts).Extract()
	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault403); ok {
			retryNode, err := recursiveCreate(nodeClient, createOpts, clusterid, 403)
			if err == "fail" {
				return fmt.Errorf("Error creating flexibleengine Node")
			}
			s = retryNode
		} else {
			return fmt.Errorf("Error creating flexibleengine Node: %s", err)
		}
	}

	job, err := nodes.GetJobDetails(nodeClient, s.Status.JobID).ExtractJob()
	if err != nil {
		return fmt.Errorf("Error fetching flexibleengine Job Details: %s", err)
	}
	jobResorceId := job.Spec.SubJobs[0].Metadata.ID

	subjob, err := nodes.GetJobDetails(nodeClient, jobResorceId).ExtractJob()
	if err != nil {
		return fmt.Errorf("Error fetching flexibleengine Job Details: %s", err)
	}

	var nodeid string
	for _, s := range subjob.Spec.SubJobs {
		if s.Spec.Type == "CreateNodeVM" {
			nodeid = s.Spec.ResourceID
			break
		}
	}

	log.Printf("[DEBUG] Waiting for CCE Node (%s) to become available", s.Metadata.Name)
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"Build", "Installing"},
		Target:     []string{"Active"},
		Refresh:    waitForCceNodeActive(nodeClient, clusterid, nodeid),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Error creating flexibleengine CCE Node: %s", err)
	}

	node, err := nodes.Get(nodeClient, clusterid, nodeid).Extract()
	d.SetId(node.Metadata.Id)
	d.Set("iptype", s.Spec.PublicIP.Eip.IpType)
	d.Set("bandwidth_charge_mode", s.Spec.PublicIP.Eip.Bandwidth.ChargeMode)
	d.Set("bandwidth_size", s.Spec.PublicIP.Eip.Bandwidth.Size)
	d.Set("sharetype", s.Spec.PublicIP.Eip.Bandwidth.ShareType)

	return resourceCCENodeV3Read(d, meta)
}

func resourceCCENodeV3Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	nodeClient, err := config.cceV3Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating flexibleengine CCE Node client: %s", err)
	}
	clusterid := d.Get("cluster_id").(string)
	s, err := nodes.Get(nodeClient, clusterid, d.Id()).Extract()

	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error retrieving flexibleengine Node: %s", err)
	}

	d.Set("name", s.Metadata.Name)
	d.Set("flavor_id", s.Spec.Flavor)
	d.Set("availability_zone", s.Spec.Az)
	d.Set("os", s.Spec.Os)
	d.Set("billing_mode", s.Spec.BillingMode)
	d.Set("extend_param_charging_mode", s.Spec.ExtendParam.ChargingMode)
	d.Set("ecs:performance_type", s.Spec.ExtendParam.PublicKey)
	d.Set("order_id", s.Spec.ExtendParam.OrderID)
	d.Set("product_id", s.Spec.ExtendParam.ProductID)
	d.Set("max_pods", s.Spec.ExtendParam.MaxPods)
	d.Set("ecs_performance_type", s.Spec.ExtendParam.EcsPerformanceType)
	d.Set("key_pair", s.Spec.Login.SshKey)

	var volumes []map[string]interface{}
	for _, pairObject := range s.Spec.DataVolumes {
		volume := make(map[string]interface{})
		volume["size"] = pairObject.Size
		volume["volumetype"] = pairObject.VolumeType
		volume["extend_param"] = pairObject.ExtendParam
		volumes = append(volumes, volume)
	}
	if err := d.Set("data_volumes", volumes); err != nil {
		return fmt.Errorf("[DEBUG] Error saving dataVolumes to state for flexibleengine Node (%s): %s", d.Id(), err)
	}

	rootVolume := []map[string]interface{}{
		{
			"size":         s.Spec.RootVolume.Size,
			"volumetype":   s.Spec.RootVolume.VolumeType,
			"extend_param": s.Spec.RootVolume.ExtendParam,
		},
	}
	d.Set("root_volume", rootVolume)
	if err := d.Set("root_volume", rootVolume); err != nil {
		return fmt.Errorf("[DEBUG] Error saving root Volume to state for flexibleengine Node (%s): %s", d.Id(), err)
	}

	d.Set("eip_ids", s.Spec.PublicIP.Ids)
	d.Set("eip_count", s.Spec.PublicIP.Count)
	d.Set("region", GetRegion(d, config))
	d.Set("private_ip", s.Status.PrivateIP)
	d.Set("public_ip", s.Status.PublicIP)

	serverId := s.Status.ServerID
	d.Set("server_id", serverId)

	// fetch tags from ECS instance
	computeClient, err := config.computeV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating Flexibleengine compute client: %s", err)
	}

	resourceTags, err := tags.Get(computeClient, "servers", serverId).Extract()
	if err != nil {
		return fmt.Errorf("Error fetching Flexibleengine instance tags: %s", err)
	}

	tagmap := tagsToMap(resourceTags.Tags)
	//ignore "CCE-Dynamic-Provisioning-Node"
	delete(tagmap, "CCE-Dynamic-Provisioning-Node")
	if err := d.Set("tags", tagmap); err != nil {
		return fmt.Errorf("Error saving tags of cce node: %s", err)
	}

	return nil
}

func resourceCCENodeV3Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	nodeClient, err := config.cceV3Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating flexibleengine CCE client: %s", err)
	}

	var updateOpts nodes.UpdateOpts

	if d.HasChange("name") {
		updateOpts.Metadata.Name = d.Get("name").(string)

		clusterid := d.Get("cluster_id").(string)
		_, err = nodes.Update(nodeClient, clusterid, d.Id(), updateOpts).Extract()
		if err != nil {
			return fmt.Errorf("Error updating flexibleengine Node: %s", err)
		}
	}

	// update tags
	if d.HasChange("tags") {
		computeClient, err := config.computeV1Client(GetRegion(d, config))
		if err != nil {
			return fmt.Errorf("Error creating Flexibleengine compute client: %s", err)
		}

		serverId := d.Get("server_id").(string)
		tagErr := UpdateResourceTags(computeClient, d, "servers", serverId)
		if tagErr != nil {
			return fmt.Errorf("Error updateing tags of cce node %s: %s", d.Id(), tagErr)
		}
	}
	return resourceCCENodeV3Read(d, meta)
}

func resourceCCENodeV3Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	nodeClient, err := config.cceV3Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating flexibleengine CCE client: %s", err)
	}
	clusterid := d.Get("cluster_id").(string)
	err = nodes.Delete(nodeClient, clusterid, d.Id()).ExtractErr()
	if err != nil {
		return fmt.Errorf("Error deleting flexibleengine CCE Cluster: %s", err)
	}
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"Deleting"},
		Target:     []string{"Deleted"},
		Refresh:    waitForCceNodeDelete(nodeClient, clusterid, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Error deleting flexibleengine CCE Node: %s", err)
	}

	d.SetId("")
	return nil
}

func waitForCceNodeActive(cceClient *golangsdk.ServiceClient, clusterId, nodeId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		n, err := nodes.Get(cceClient, clusterId, nodeId).Extract()
		if err != nil {
			return nil, "", err
		}

		return n, n.Status.Phase, nil
	}
}

func waitForCceNodeDelete(cceClient *golangsdk.ServiceClient, clusterId, nodeId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		log.Printf("[DEBUG] Attempting to delete flexibleengine CCE Node %s.\n", nodeId)

		r, err := nodes.Get(cceClient, clusterId, nodeId).Extract()

		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				log.Printf("[DEBUG] Successfully deleted flexibleengine CCE Node %s", nodeId)
				return r, "Deleted", nil
			}
			return r, "Deleting", err
		}

		log.Printf("[DEBUG] flexibleengine CCE Node %s still available.\n", nodeId)
		return r, r.Status.Phase, nil
	}
}

func waitForClusterAvailable(cceClient *golangsdk.ServiceClient, clusterId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		log.Printf("[INFO] Waiting for flexibleengine Cluster to be available %s.\n", clusterId)
		n, err := clusters.Get(cceClient, clusterId).Extract()

		if err != nil {
			return nil, "", err
		}

		return n, n.Status.Phase, nil
	}
}

func recursiveCreate(cceClient *golangsdk.ServiceClient, opts nodes.CreateOptsBuilder, ClusterID string, errCode int) (*nodes.Nodes, string) {
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
		s, err := nodes.Create(cceClient, ClusterID, opts).Extract()
		if err != nil {
			//if err.(golangsdk.ErrUnexpectedResponseCode).Actual == 403 {
			if _, ok := err.(golangsdk.ErrDefault403); ok {
				return recursiveCreate(cceClient, opts, ClusterID, 403)
			} else {
				return s, "fail"
			}
		} else {
			return s, "success"
		}
	}
	return nil, "fail"
}

func installScriptHashSum(script string) string {
	// Check whether the preinstall/postinstall is not Base64 encoded.
	// Always calculate hash of base64 decoded value since we
	// check against double-encoding when setting it
	v, base64DecodeError := base64.StdEncoding.DecodeString(script)
	if base64DecodeError != nil {
		v = []byte(script)
	}

	hash := sha1.Sum(v)
	return hex.EncodeToString(hash[:])
}

func installScriptEncode(script string) string {
	if _, err := base64.StdEncoding.DecodeString(script); err != nil {
		return base64.StdEncoding.EncodeToString([]byte(script))
	}
	return script
}
