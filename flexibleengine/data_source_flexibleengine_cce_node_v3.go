package flexibleengine

import (
	"fmt"
	"log"

	"github.com/chnsz/golangsdk/openstack/cce/v3/nodes"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceCceNodesV3() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceCceNodesV3Read,

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
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"node_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"flavor_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"availability_zone": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"key_pair": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"bandwidth_size": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"share_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"charge_mode": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ip_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"disk_size": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"volume_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"data_volumes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"disk_size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"volume_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"billing_mode": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"server_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"public_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"private_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"eip_ids": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceCceNodesV3Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	cceClient, err := config.CceV3Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Unable to create flexibleengine CCE client : %s", err)
	}

	listOpts := nodes.ListOpts{
		Uid:   d.Get("node_id").(string),
		Name:  d.Get("name").(string),
		Phase: d.Get("status").(string),
	}

	refinedNodes, err := nodes.List(cceClient, d.Get("cluster_id").(string), listOpts)
	if err != nil {
		return fmt.Errorf("Unable to retrieve Nodes: %s", err)
	}

	if len(refinedNodes) < 1 {
		return fmt.Errorf("Your query returned no results. " +
			"Please change your search criteria and try again.")
	}

	if len(refinedNodes) > 1 {
		return fmt.Errorf("Your query returned more than one result." +
			" Please try a more specific search criteria")
	}

	Node := refinedNodes[0]
	log.Printf("[DEBUG] Retrieved Node %s using given filter: %+v", Node.Metadata.Id, Node)

	d.SetId(Node.Metadata.Id)
	d.Set("node_id", Node.Metadata.Id)
	d.Set("name", Node.Metadata.Name)
	d.Set("flavor_id", Node.Spec.Flavor)
	d.Set("availability_zone", Node.Spec.Az)
	d.Set("status", Node.Status.Phase)
	d.Set("disk_size", Node.Spec.RootVolume.Size)
	d.Set("volume_type", Node.Spec.RootVolume.VolumeType)
	d.Set("key_pair", Node.Spec.Login.SshKey)
	d.Set("server_id", Node.Status.ServerID)
	d.Set("public_ip", Node.Status.PublicIP)
	d.Set("private_ip", Node.Status.PrivateIP)
	d.Set("billing_mode", Node.Spec.BillingMode)

	d.Set("eip_ids", Node.Spec.PublicIP.Ids)
	d.Set("ip_type", Node.Spec.PublicIP.Eip.IpType)
	d.Set("charge_mode", Node.Spec.PublicIP.Eip.Bandwidth.ChargeMode)
	d.Set("share_type", Node.Spec.PublicIP.Eip.Bandwidth.ShareType)
	if bandwidthSize := Node.Spec.PublicIP.Eip.Bandwidth.Size; bandwidthSize > 0 {
		d.Set("bandwidth_size", bandwidthSize)
	}

	var v []map[string]interface{}
	for _, volume := range Node.Spec.DataVolumes {
		mapping := map[string]interface{}{
			"disk_size":   volume.Size,
			"volume_type": volume.VolumeType,
		}
		v = append(v, mapping)
	}
	d.Set("data_volumes", v)

	return nil
}
