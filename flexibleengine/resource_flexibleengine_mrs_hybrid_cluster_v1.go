package flexibleengine

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/chnsz/golangsdk/openstack/mrs/v1/cluster"
	"github.com/chnsz/golangsdk/openstack/networking/v1/subnets"
	"github.com/chnsz/golangsdk/openstack/networking/v1/vpcs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceMRSHybridClusterV1() *schema.Resource {
	return &schema.Resource{
		Create: resourceHybridClusterV1Create,
		Read:   resourceHybridClusterV1Read,
		Delete: resourceHybridClusterV1Delete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"available_zone": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"cluster_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"cluster_version": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"safe_mode": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				Default:      1,
				ValidateFunc: validation.IntBetween(0, 1),
			},
			"cluster_admin_secret": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
				ForceNew:  true,
			},
			"master_node_key_pair": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
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
				Optional: true,
				ForceNew: true,
			},
			"log_collection": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
				Default:  1,
			},
			"master_nodes": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				MaxItems: 1,
				Elem:     nodeSchemaResource(1, 2),
			},
			"analysis_core_nodes": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				MaxItems: 1,
				Elem:     nodeSchemaResource(1, 500),
			},
			"streaming_core_nodes": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				MaxItems: 1,
				Elem:     nodeSchemaResource(1, 500),
			},
			"analysis_task_nodes": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem:     nodeSchemaResource(1, 500),
			},
			"streaming_task_nodes": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem:     nodeSchemaResource(1, 500),
			},
			"component_list": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"total_node_number": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"master_node_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"external_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"private_ip_first": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"internal_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"external_alternate_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"billing_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vnc": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"update_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"charging_start_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"components": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"component_name": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"component_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"component_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"component_desc": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func nodeSchemaResource(min, max int) *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"flavor": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"node_number": {
				Type:         schema.TypeInt,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntBetween(min, max),
			},
			"data_volume_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"data_volume_size": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"data_volume_count": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"root_volume_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"root_volume_size": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func resourceHybridClusterV1Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	region := GetRegion(d, config)
	client, err := config.MrsV1Client(region)
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine MRS client: %s", err)
	}
	vpcClient, err := config.NetworkingV1Client(region)
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine Vpc client: %s", err)
	}

	// Get vpc name
	vpc, err := vpcs.Get(vpcClient, d.Get("vpc_id").(string)).Extract()
	if err != nil {
		return fmt.Errorf("Error retrieving FlexibleEngine Vpc: %s", err)
	}
	// Get subnet name
	subnet, err := subnets.Get(vpcClient, d.Get("subnet_id").(string)).Extract()
	if err != nil {
		return fmt.Errorf("Error retrieving FlexibleEngine Subnet: %s", err)
	}

	createOpts := &cluster.CreateOpts{
		BillingType:        12,
		ClusterType:        2,
		LoginMode:          1,
		DataCenter:         region,
		AvailableZoneID:    d.Get("available_zone").(string),
		ClusterName:        d.Get("cluster_name").(string),
		ClusterVersion:     d.Get("cluster_version").(string),
		NodePublicCertName: d.Get("master_node_key_pair").(string),
		SafeMode:           d.Get("safe_mode").(int),
		ClusterAdminSecret: d.Get("cluster_admin_secret").(string),
		LogCollection:      d.Get("log_collection").(int),
		VpcID:              d.Get("vpc_id").(string),
		SubnetID:           d.Get("subnet_id").(string),
		SecurityGroupsID:   d.Get("security_group_id").(string),
		Vpc:                vpc.Name,
		SubnetName:         subnet.Name,
		NodeGroups:         getHybridClusterNodeGroups(d),
		ComponentList:      getHybridClusterComponents(d),
	}

	log.Printf("[DEBUG] Create Options: %#v", createOpts)

	clusterCreate, err := cluster.Create(client, createOpts).Extract()
	if err != nil {
		return fmt.Errorf("Error creating Cluster: %s", err)
	}

	d.SetId(clusterCreate.ClusterID)
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"starting"},
		Target:       []string{"running"},
		Refresh:      ClusterStateRefreshFunc(client, clusterCreate.ClusterID),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        600 * time.Second,
		PollInterval: 20 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		//the system will recyle the cluster when creating failed
		d.SetId("")
		return fmt.Errorf(
			"Error waiting for cluster (%s) to become ready: %s ",
			clusterCreate.ClusterID, err)
	}

	return resourceHybridClusterV1Read(d, meta)
}

func resourceHybridClusterV1Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client, err := config.MrsV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine MRS client: %s", err)
	}

	clusterGet, err := cluster.Get(client, d.Id()).Extract()
	if err != nil {
		return CheckDeleted(d, err, "Cluster")
	}

	// ignore the terminated cluster
	if clusterGet.Clusterstate == "terminated" {
		log.Printf("[INFO] Retrieved Cluster %s, but it was terminated, abort it", d.Id())
		d.SetId("")
		return nil
	}

	log.Printf("[DEBUG] Retrieved Cluster %s: %#v", d.Id(), clusterGet)
	d.SetId(clusterGet.Clusterid)
	d.Set("region", clusterGet.Datacenter)
	d.Set("available_zone", clusterGet.Azname)
	d.Set("cluster_name", clusterGet.Clustername)
	d.Set("cluster_version", clusterGet.Clusterversion)
	d.Set("state", clusterGet.Clusterstate)
	d.Set("safe_mode", clusterGet.Safemode)
	d.Set("log_collection", clusterGet.LogCollection)
	d.Set("master_node_key_pair", clusterGet.Nodepubliccertname)
	d.Set("vpc_id", clusterGet.Vpcid)
	d.Set("subnet_id", clusterGet.Subnetid)
	d.Set("security_groups_id", clusterGet.Securitygroupsid)

	d.Set("billing_type", clusterGet.Billingtype)
	d.Set("vnc", clusterGet.Vnc)
	d.Set("master_node_ip", clusterGet.Masternodeip)
	d.Set("external_ip", clusterGet.Externalip)
	d.Set("private_ip_first", clusterGet.Privateipfirst)
	d.Set("internal_ip", clusterGet.Internalip)
	d.Set("slave_security_groups_id", clusterGet.Slavesecuritygroupsid)
	d.Set("external_alternate_ip", clusterGet.Externalalternateip)

	totalNodes, err := strconv.Atoi(clusterGet.Totalnodenum)
	if err == nil {
		d.Set("total_node_number", totalNodes)
	}

	updateAt, err := strconv.ParseInt(clusterGet.Updateat, 10, 64)
	if err == nil {
		updateAtTm := time.Unix(updateAt, 0)
		d.Set("update_at", updateAtTm.Format(time.RFC1123))
	}

	createAt, err := strconv.ParseInt(clusterGet.Createat, 10, 64)
	if err == nil {
		createAtTm := time.Unix(createAt, 0)
		d.Set("create_at", createAtTm.Format(time.RFC1123))
	}

	chargingStartTime, err := strconv.ParseInt(clusterGet.Chargingstarttime, 10, 64)
	if err == nil {
		chargingStartTimeTm := time.Unix(chargingStartTime, 0)
		d.Set("charging_start_time", chargingStartTimeTm.Format(time.RFC1123))
	}

	components := make([]map[string]interface{}, len(clusterGet.Componentlist))
	for i, attachment := range clusterGet.Componentlist {
		components[i] = make(map[string]interface{})
		components[i]["component_id"] = attachment.Componentid
		components[i]["component_name"] = attachment.Componentname
		components[i]["component_version"] = attachment.Componentversion
		components[i]["component_desc"] = attachment.Componentdesc
		log.Printf("[DEBUG] components: %v", components)
	}
	d.Set("components", components)

	return setHybridClusterNodeGroups(d, clusterGet)
}

func resourceHybridClusterV1Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client, err := config.MrsV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine MRS client: %s", err)
	}

	rId := d.Id()
	clusterGet, err := cluster.Get(client, d.Id()).Extract()
	if err != nil {
		if isResourceNotFound(err) {
			log.Printf("[INFO] getting an unavailable Cluster: %s", rId)
			return nil
		}
		return fmt.Errorf("Error getting Cluster %s: %s", rId, err)
	}

	if clusterGet.Clusterstate == "terminated" {
		log.Printf("[DEBUG] The Cluster %s has been terminated.", rId)
		d.SetId("")
		return nil
	}

	log.Printf("[DEBUG] Deleting Cluster %s", rId)

	err = cluster.Delete(client, rId).ExtractErr()
	if err != nil {
		return fmt.Errorf("Error deleting FlexibleEngine Cluster: %s", err)
	}

	log.Printf("[DEBUG] Waiting for Cluster (%s) to be terminated", rId)

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"running", "terminating"},
		Target:       []string{"terminated"},
		Refresh:      ClusterStateRefreshFunc(client, rId),
		Timeout:      d.Timeout(schema.TimeoutDelete),
		Delay:        40 * time.Second,
		PollInterval: 10 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf(
			"Error waiting for Cluster (%s) to be terminated: %s",
			d.Id(), err)
	}

	d.SetId("")
	return nil
}

func getHybridClusterComponents(d *schema.ResourceData) []cluster.ComponentOpts {
	components := d.Get("component_list").([]interface{})
	componentOpts := make([]cluster.ComponentOpts, len(components))

	for i, v := range components {
		opts := cluster.ComponentOpts{
			ComponentName: v.(string),
		}
		componentOpts[i] = opts
	}

	return componentOpts
}

func getHybridClusterNodeGroups(d *schema.ResourceData) []cluster.NodeGroupOpts {
	var groupOpts []cluster.NodeGroupOpts
	var groupMap = map[string]string{
		"master_nodes":         "master_node_default_group",
		"analysis_core_nodes":  "core_node_analysis_group",
		"streaming_core_nodes": "core_node_streaming_group",
		"analysis_task_nodes":  "task_node_analysis_group",
		"streaming_task_nodes": "task_node_streaming_group",
	}
	for k, v := range groupMap {
		if opts := getNodeGroupOpts(d, k, v); opts != nil {
			groupOpts = append(groupOpts, *opts)
		}
	}
	return groupOpts
}

func getNodeGroupOpts(d *schema.ResourceData, key, name string) *cluster.NodeGroupOpts {
	optsRaw := d.Get(key).([]interface{})
	if len(optsRaw) == 1 {
		opts := optsRaw[0].(map[string]interface{})
		return &cluster.NodeGroupOpts{
			GroupName:       name,
			NodeSize:        opts["flavor"].(string),
			NodeNum:         opts["node_number"].(int),
			DataVolumeType:  opts["data_volume_type"].(string),
			DataVolumeSize:  opts["data_volume_size"].(int),
			DataVolumeCount: opts["data_volume_count"].(int),
			RootVolumeType:  "SATA",
			RootVolumeSize:  40,
		}
	}

	return nil
}

func setHybridClusterNodeGroups(d *schema.ResourceData, data *cluster.Cluster) error {
	var groupMap = map[string]string{
		"master_node_default_group": "master_nodes",
		"core_node_analysis_group":  "analysis_core_nodes",
		"core_node_streaming_group": "streaming_core_nodes",
		"task_node_analysis_group":  "analysis_task_nodes",
		"task_node_streaming_group": "streaming_task_nodes",
	}

	allGroups := append(data.NodeGroups, data.TaskNodeGroups...)
	for _, group := range allGroups {
		if key, ok := groupMap[group.GroupName]; ok {
			if err := d.Set(key, flattenClusterNodeGroup(group)); err != nil {
				return fmt.Errorf("Error saving %s: %s", key, err)
			}
		} else {
			log.Printf("[DEBUG] %s is not in the resource data", group.GroupName)
		}
	}

	return nil
}

func flattenClusterNodeGroup(group cluster.NodeGroup) []map[string]interface{} {
	return []map[string]interface{}{
		{
			"node_number":       group.NodeNum,
			"flavor":            group.NodeSize,
			"data_volume_type":  group.DataVolumeType,
			"data_volume_size":  group.DataVolumeSize,
			"data_volume_count": group.DataVolumeCount,
			"root_volume_type":  group.RootVolumeType,
			"root_volume_size":  group.RootVolumeSize,
		},
	}
}
