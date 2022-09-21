package flexibleengine

import (
	"fmt"
	"log"

	"github.com/chnsz/golangsdk/openstack/cce/v3/clusters"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceCCEClusterV3() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceCCEClusterV3Read,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"flavor_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cluster_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cluster_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"billing_mode": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"subnet_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"highway_subnet_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"container_network_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"container_network_cidr": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"service_network_cidr": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"internal_endpoint": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"external_endpoint": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"external_apig_endpoint": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"security_group_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"authentication_mode": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"masters": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"availability_zone": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceCCEClusterV3Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	region := GetRegion(d, config)
	cceClient, err := config.CceV3Client(region)
	if err != nil {
		return fmt.Errorf("Unable to create flexibleengine CCE client : %s", err)
	}

	listOpts := clusters.ListOpts{
		ID:    d.Get("id").(string),
		Name:  d.Get("name").(string),
		Type:  d.Get("cluster_type").(string),
		Phase: d.Get("status").(string),
		VpcID: d.Get("vpc_id").(string),
	}

	refinedClusters, err := clusters.List(cceClient, listOpts)
	log.Printf("[DEBUG] Value of allClusters: %#v", refinedClusters)
	if err != nil {
		return fmt.Errorf("Unable to retrieve clusters: %s", err)
	}

	if len(refinedClusters) < 1 {
		return fmt.Errorf("Your query returned no results. " +
			"Please change your search criteria and try again.")
	}

	if len(refinedClusters) > 1 {
		return fmt.Errorf("Your query returned more than one result." +
			" Please try a more specific search criteria")
	}

	clusterInfo := refinedClusters[0]
	log.Printf("[DEBUG] Retrieved Clusters using given filter %s: %+v", clusterInfo.Metadata.Id, clusterInfo)

	d.SetId(clusterInfo.Metadata.Id)
	d.Set("region", region)
	d.Set("name", clusterInfo.Metadata.Name)
	d.Set("flavor_id", clusterInfo.Spec.Flavor)
	d.Set("description", clusterInfo.Spec.Description)
	d.Set("cluster_version", clusterInfo.Spec.Version)
	d.Set("cluster_type", clusterInfo.Spec.Type)
	d.Set("billing_mode", clusterInfo.Spec.BillingMode)
	d.Set("vpc_id", clusterInfo.Spec.HostNetwork.VpcId)
	d.Set("subnet_id", clusterInfo.Spec.HostNetwork.SubnetId)
	d.Set("security_group_id", clusterInfo.Spec.HostNetwork.SecurityGroup)
	d.Set("highway_subnet_id", clusterInfo.Spec.HostNetwork.HighwaySubnet)
	d.Set("container_network_cidr", clusterInfo.Spec.ContainerNetwork.Cidr)
	d.Set("container_network_type", clusterInfo.Spec.ContainerNetwork.Mode)
	d.Set("service_network_cidr", clusterInfo.Spec.KubernetesSvcIPRange)
	d.Set("authentication_mode", clusterInfo.Spec.Authentication.Mode)
	d.Set("status", clusterInfo.Status.Phase)

	// Set masters
	var masterList []map[string]interface{}
	for _, masterObj := range clusterInfo.Spec.Masters {
		master := make(map[string]interface{})
		master["availability_zone"] = masterObj.MasterAZ
		masterList = append(masterList, master)
	}
	d.Set("masters", masterList)

	// Set endpoint
	var internalEP, externalEP string
	for _, ep := range clusterInfo.Status.Endpoints {
		if ep.Type == "Internal" {
			internalEP = ep.Url
		} else if ep.Type == "External" {
			externalEP = ep.Url
		}
	}
	d.Set("internal_endpoint", internalEP)
	d.Set("external_endpoint", externalEP)
	// the value is always empty, keep compatibility
	d.Set("external_apig_endpoint", clusterInfo.Status.Endpoints[0].ExternalOTC)

	return nil
}
