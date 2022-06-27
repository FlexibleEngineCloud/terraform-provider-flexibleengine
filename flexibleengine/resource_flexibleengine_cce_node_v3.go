package flexibleengine

import (
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/cce/v3/clusters"
	"github.com/chnsz/golangsdk/openstack/cce/v3/nodes"
	"github.com/chnsz/golangsdk/openstack/common/tags"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceCCENodeV3() *schema.Resource {
	return &schema.Resource{
		Create: resourceCCENodeV3Create,
		Read:   resourceCCENodeV3Read,
		Update: resourceCCENodeV3Update,
		Delete: resourceCCENodeV3Delete,
		Importer: &schema.ResourceImporter{
			State: resourceCCENodeV3Import,
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
			"annotations": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
			},
			"tags": tagsSchema(),

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
			"taints": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"value": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"effect": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
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
			"ecs_performance_type": {
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
			"extend_param": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"private_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"public_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"server_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},

			// Deprecated
			"billing_mode": {
				Type:       schema.TypeInt,
				Optional:   true,
				ForceNew:   true,
				Deprecated: "will be removed later",
			},
			"extend_param_charging_mode": {
				Type:       schema.TypeInt,
				Optional:   true,
				ForceNew:   true,
				Deprecated: "will be removed later",
			},
			"order_id": {
				Type:       schema.TypeString,
				Optional:   true,
				ForceNew:   true,
				Deprecated: "will be removed later",
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
			ExtendParam: rawMap["extend_params"].(map[string]interface{}),
		}

		if rawMap["kms_key_id"].(string) != "" {
			metadata := nodes.VolumeMetadata{
				SystemEncrypted: "1",
				SystemCmkid:     rawMap["kms_key_id"].(string),
			}
			volumes[i].Metadata = &metadata
		}
	}
	return volumes
}

func resourceCCERootVolume(d *schema.ResourceData) nodes.VolumeSpec {
	var root nodes.VolumeSpec
	rootRaw := d.Get("root_volume").([]interface{})
	if len(rootRaw) == 1 {
		rawMap := rootRaw[0].(map[string]interface{})
		root.Size = rawMap["size"].(int)
		root.VolumeType = rawMap["volumetype"].(string)
		root.ExtendParam = rawMap["extend_params"].(map[string]interface{})
	}
	return root
}

func resourceCCEExtendParam(d *schema.ResourceData) map[string]interface{} {
	extendParam := make(map[string]interface{})
	if v, ok := d.GetOk("extend_param"); ok {
		for key, val := range v.(map[string]interface{}) {
			extendParam[key] = val.(string)
		}
		if v, ok := extendParam["periodNum"]; ok {
			periodNum, err := strconv.Atoi(v.(string))
			if err != nil {
				log.Printf("[WARNING] PeriodNum %s invalid, Type conversion error: %s", v.(string), err)
			}
			extendParam["periodNum"] = periodNum
		}
	}
	if v, ok := d.GetOk("extend_param_charging_mode"); ok {
		extendParam["chargingMode"] = v.(int)
	}
	if v, ok := d.GetOk("ecs_performance_type"); ok {
		extendParam["ecs:performancetype"] = v.(string)
	}
	if v, ok := d.GetOk("max_pods"); ok {
		extendParam["maxPods"] = v.(int)
	}
	if v, ok := d.GetOk("order_id"); ok {
		extendParam["orderID"] = v.(string)
	}
	if v, ok := d.GetOk("product_id"); ok {
		extendParam["productID"] = v.(string)
	}
	if v, ok := d.GetOk("public_key"); ok {
		extendParam["publicKey"] = v.(string)
	}
	if v, ok := d.GetOk("preinstall"); ok {
		extendParam["alpha.cce/preInstall"] = installScriptEncode(v.(string))
	}
	if v, ok := d.GetOk("postinstall"); ok {
		extendParam["alpha.cce/postInstall"] = installScriptEncode(v.(string))
	}

	return extendParam
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
	nodeClient, err := config.CceV3Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating flexibleengine CCE Node client: %s", err)
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
			ExtendParam: resourceCCEExtendParam(d),
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
		},
	}

	clusterid := d.Get("cluster_id").(string)
	stateCluster := &resource.StateChangeConf{
		Target:       []string{"Available"},
		Refresh:      waitForClusterAvailable(nodeClient, clusterid),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        5 * time.Second,
		PollInterval: 5 * time.Second,
	}
	_, err = stateCluster.WaitForState()

	log.Printf("[DEBUG] Create Options: %#v", createOpts)
	s, err := nodes.Create(nodeClient, clusterid, createOpts).Extract()
	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault403); ok {
			retryNode, err := recursiveCreate(nodeClient, createOpts, clusterid, 403)
			if err == "fail" {
				return fmt.Errorf("Error creating flexibleengine Node: %s", err)
			}
			s = retryNode
		} else {
			return fmt.Errorf("Error creating flexibleengine Node: %s", err)
		}
	}

	nodeID, err := getResourceIDFromJob(nodeClient, s.Status.JobID, "CreateNode", "CreateNodeVM",
		d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return err
	}
	d.SetId(nodeID)

	log.Printf("[DEBUG] Waiting for CCE Node (%s) to become available", s.Metadata.Name)
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"Build", "Installing"},
		Target:       []string{"Active"},
		Refresh:      waitForCceNodeActive(nodeClient, clusterid, nodeID),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        5 * time.Second,
		PollInterval: 5 * time.Second,
	}
	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Error creating flexibleengine CCE Node: %s", err)
	}

	return resourceCCENodeV3Read(d, meta)
}

func resourceCCENodeV3Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	nodeClient, err := config.CceV3Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating flexibleengine CCE Node client: %s", err)
	}
	computeClient, err := config.ComputeV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating Flexibleengine compute client: %s", err)
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

	mErr := multierror.Append(
		d.Set("region", GetRegion(d, config)),
		d.Set("name", s.Metadata.Name),
		d.Set("flavor_id", s.Spec.Flavor),
		d.Set("availability_zone", s.Spec.Az),
		d.Set("os", s.Spec.Os),
		d.Set("key_pair", s.Spec.Login.SshKey),
		d.Set("ecs_performance_type", s.Spec.ExtendParam["ecs:performancetype"]),
		d.Set("product_id", s.Spec.ExtendParam["productID"]),
		d.Set("public_key", s.Spec.ExtendParam["publicKey"]),
		d.Set("max_pods", s.Spec.ExtendParam["maxPods"]),
		d.Set("private_ip", s.Status.PrivateIP),
		d.Set("public_ip", s.Status.PublicIP),
		d.Set("status", s.Status.Phase),
		d.Set("server_id", s.Status.ServerID),

		d.Set("eip_ids", s.Spec.PublicIP.Ids),
		d.Set("eip_count", s.Spec.PublicIP.Count),
		d.Set("iptype", s.Spec.PublicIP.Eip.IpType),
		d.Set("bandwidth_charge_mode", s.Spec.PublicIP.Eip.Bandwidth.ChargeMode),
		d.Set("bandwidth_size", s.Spec.PublicIP.Eip.Bandwidth.Size),
		d.Set("sharetype", s.Spec.PublicIP.Eip.Bandwidth.ShareType),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return err
	}

	rootVolume := expandResourceCCERootVolume(s.Spec)
	d.Set("root_volume", rootVolume)
	if err := d.Set("root_volume", rootVolume); err != nil {
		return fmt.Errorf("Error saving root volume of cce node %s: %s", d.Id(), err)
	}

	volumes := expandResourceCCEDataVolumes(s.Spec)
	if err := d.Set("data_volumes", volumes); err != nil {
		return fmt.Errorf("Error saving data volumes of cce node %s: %s", d.Id(), err)
	}

	nodeTaints := expandResourceCCETaints(s.Spec)
	if err := d.Set("taints", nodeTaints); err != nil {
		return fmt.Errorf("Error saving taints of cce node %s: %s", d.Id(), err)
	}

	labels := expandResourceCCEK8sTags(s.Spec)
	if err := d.Set("labels", labels); err != nil {
		return fmt.Errorf("Error saving labels/k8stags of cce node %s: %s", d.Id(), err)
	}

	// fetch tags from ECS instance as Spec.UserTags is empty
	if tagmap, err := expandResourceCCETagsByServer(computeClient, s.Status.ServerID); err == nil {
		d.Set("tags", tagmap)
	} else {
		return err
	}

	return nil
}

func resourceCCENodeV3Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	nodeClient, err := config.CceV3Client(GetRegion(d, config))
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
		computeClient, err := config.ComputeV1Client(GetRegion(d, config))
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
	nodeClient, err := config.CceV3Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating flexibleengine CCE client: %s", err)
	}
	clusterid := d.Get("cluster_id").(string)
	err = nodes.Delete(nodeClient, clusterid, d.Id()).ExtractErr()
	if err != nil {
		return fmt.Errorf("Error deleting flexibleengine CCE Cluster: %s", err)
	}
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"Deleting"},
		Target:       []string{"Deleted"},
		Refresh:      waitForCceNodeDelete(nodeClient, clusterid, d.Id()),
		Timeout:      d.Timeout(schema.TimeoutDelete),
		Delay:        60 * time.Second,
		PollInterval: 15 * time.Second,
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
			Target:       []string{"Available"},
			Refresh:      waitForClusterAvailable(cceClient, ClusterID),
			Timeout:      15 * time.Minute,
			Delay:        5 * time.Second,
			PollInterval: 5 * time.Second,
		}
		_, stateErr := stateCluster.WaitForState()
		if stateErr != nil {
			log.Printf("[INFO] Cluster Unavailable %s.\n", stateErr)
		}
		s, err := nodes.Create(cceClient, ClusterID, opts).Extract()
		if err != nil {
			return s, "fail"
		}
		return s, "success"
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

func expandResourceCCERootVolume(spec nodes.Spec) []map[string]interface{} {
	rootVolume := []map[string]interface{}{
		{
			"size":          spec.RootVolume.Size,
			"volumetype":    spec.RootVolume.VolumeType,
			"extend_params": spec.RootVolume.ExtendParam,
		},
	}

	return rootVolume
}

func expandResourceCCEDataVolumes(spec nodes.Spec) []map[string]interface{} {
	volumes := make([]map[string]interface{}, len(spec.DataVolumes))
	for i, item := range spec.DataVolumes {
		var kmsID string
		if item.Metadata != nil {
			kmsID = item.Metadata.SystemCmkid
		}

		volumes[i] = map[string]interface{}{
			"size":          item.Size,
			"volumetype":    item.VolumeType,
			"extend_params": item.ExtendParam,
			"kms_key_id":    kmsID,
		}
	}
	return volumes
}

func expandResourceCCETaints(spec nodes.Spec) []map[string]interface{} {
	taints := make([]map[string]interface{}, len(spec.Taints))
	for i, item := range spec.Taints {
		taints[i] = map[string]interface{}{
			"key":    item.Key,
			"value":  item.Value,
			"effect": item.Effect,
		}
	}
	return taints
}

func expandResourceCCEK8sTags(spec nodes.Spec) map[string]string {
	labels := map[string]string{}
	for key, val := range spec.K8sTags {
		if strings.Contains(key, "cce.cloud.com") {
			continue
		}
		labels[key] = val
	}
	return labels
}

// expandResourceCCETagsByServer fetch tags from compute instance
// we have to call it as Spec.UserTags is empty in nodes.Get response
func expandResourceCCETagsByServer(client *golangsdk.ServiceClient, serverID string) (map[string]string, error) {
	resourceTags, err := tags.Get(client, "servers", serverID).Extract()
	if err != nil {
		return nil, fmt.Errorf("Error fetching compute instance tags: %s", err)
	}

	tagmap := tagsToMap(resourceTags.Tags)
	//ignore "CCE-Dynamic-Provisioning-Node"
	delete(tagmap, "CCE-Dynamic-Provisioning-Node")

	return tagmap, nil
}

func resourceCCENodeV3Import(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	parts := strings.SplitN(d.Id(), "/", 2)
	if len(parts) != 2 {
		err := fmt.Errorf("Invalid format specified for CCE Node. Format must be <cluster id>/<node id>")
		return nil, err
	}

	clusterID := parts[0]
	nodeID := parts[1]

	d.SetId(nodeID)
	d.Set("cluster_id", clusterID)

	return []*schema.ResourceData{d}, nil
}

func waitForJobStatus(cceClient *golangsdk.ServiceClient, jobID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		job, err := nodes.GetJobDetails(cceClient, jobID).ExtractJob()
		if err != nil {
			return nil, "", err
		}

		return job, job.Status.Phase, nil
	}
}

func getResourceIDFromJob(client *golangsdk.ServiceClient, jobID, jobType, subJobType string,
	timeout time.Duration) (string, error) {

	stateJob := &resource.StateChangeConf{
		Pending:      []string{"Initializing", "Running"},
		Target:       []string{"Success"},
		Refresh:      waitForJobStatus(client, jobID),
		Timeout:      timeout,
		Delay:        120 * time.Second,
		PollInterval: 20 * time.Second,
	}

	v, err := stateJob.WaitForState()
	if err != nil {
		if job, ok := v.(*nodes.Job); ok {
			return "", fmt.Errorf("Error waiting for job (%s) to become success: %s, reason: %s",
				jobID, err, job.Status.Reason)
		}

		return "", fmt.Errorf("Error waiting for job (%s) to become success: %s", jobID, err)
	}

	job := v.(*nodes.Job)
	if len(job.Spec.SubJobs) == 0 {
		return "", fmt.Errorf("Error fetching sub jobs from %s", jobID)
	}

	var subJobID string
	var refreshJob bool
	for _, s := range job.Spec.SubJobs {
		// postPaid: should get details of sub job ID
		if s.Spec.Type == jobType {
			subJobID = s.Metadata.ID
			refreshJob = true
			break
		}
	}

	if refreshJob {
		job, err = nodes.GetJobDetails(client, subJobID).ExtractJob()
		if err != nil {
			return "", fmt.Errorf("Error fetching sub Job %s: %s", subJobID, err)
		}
	}

	var nodeid string
	for _, s := range job.Spec.SubJobs {
		if s.Spec.Type == subJobType {
			nodeid = s.Spec.ResourceID
			break
		}
	}
	if nodeid == "" {
		return "", fmt.Errorf("Error fetching %s Job resource id", subJobType)
	}
	return nodeid, nil
}
