package flexibleengine

import (
	"context"
	"log"
	"strings"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/hashcode"

	"github.com/chnsz/golangsdk/openstack/ecs/v1/cloudservers"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceComputeInstances() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceComputeInstancesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"flavor_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"fixed_ip_v4": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"flavor_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"image_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"availability_zone": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"key_pair": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"instances": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"availability_zone": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"flavor_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"flavor_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"image_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"image_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"key_pair": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"user_data": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"security_groups": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"floating_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"network": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"uuid": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"port": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"fixed_ip_v4": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"fixed_ip_v6": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"mac": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"volume_attached": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"volume_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"is_sys_volume": {
										Type:     schema.TypeBool,
										Computed: true,
									},
								},
							},
						},
						"scheduler_hints": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"group": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"metadata": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"tags": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func filterCloudServers(d *schema.ResourceData, servers []cloudservers.CloudServer) ([]cloudservers.CloudServer,
	[]string) {
	result := make([]cloudservers.CloudServer, 0, len(servers))
	ids := make([]string, 0, len(servers))

	for _, server := range servers {
		if flavorName, ok := d.GetOk("flavor_name"); ok && flavorName != server.Flavor.Name {
			continue
		}
		if iamgeId, ok := d.GetOk("image_id"); ok && iamgeId != server.Image.ID {
			continue
		}
		if az, ok := d.GetOk("availability_zone"); ok && az != server.AvailabilityZone {
			continue
		}
		if keypair, ok := d.GetOk("key_pair"); ok && keypair != server.KeyName {
			continue
		}
		result = append(result, server)
		ids = append(ids, server.ID)
	}

	return result, ids
}

func isSystemVolume(index string) bool {
	if index == "0" {
		return true
	}
	return false
}

func parseEcsInstanceVolumeAttachedInfo(attachList []cloudservers.VolumeAttached) []map[string]interface{} {
	result := make([]map[string]interface{}, len(attachList))

	for i, volume := range attachList {
		result[i] = map[string]interface{}{
			"volume_id":     volume.ID,
			"is_sys_volume": isSystemVolume(volume.BootIndex),
		}
	}

	return result
}

func setComputeInstancesParams(d *schema.ResourceData, servers []cloudservers.CloudServer, meta interface{}) diag.Diagnostics {
	result := make([]map[string]interface{}, len(servers))

	for i, val := range servers {
		server := map[string]interface{}{
			"id":                val.ID,
			"user_data":         val.UserData,
			"name":              val.Name,
			"flavor_name":       val.Flavor.Name,
			"status":            val.Status,
			"flavor_id":         val.Flavor.ID,
			"image_id":          val.Image.ID,
			"image_name":        val.Metadata.ImageName,
			"availability_zone": val.AvailabilityZone,
			"key_pair":          val.KeyName,
		}

		server["security_groups"] = parseEcsInstanceSecurityGroupIds(val.SecurityGroups)

		if len(val.OsSchedulerHints.Group) > 0 {
			server["scheduler_hints"] = parseEcsInstanceSchedulerHintInfo(val.OsSchedulerHints)
		}

		// Set the instance network and address information
		networks, eip := flattenComputeNetworks(d, meta, &val)
		server["network"] = networks
		if eip != "" {
			server["floating_ip"] = eip
		}

		if len(val.VolumeAttached) > 0 {
			server["volume_attached"] = parseEcsInstanceVolumeAttachedInfo(val.VolumeAttached)
		}

		if len(val.Tags) > 0 {
			server["tags"] = parseEcsInstanceTagInfo(val.Tags)
		}

		result[i] = server
	}
	if err := d.Set("instances", result); err != nil {
		return diag.Errorf("Error setting cloud server list: %s", err)
	}

	return nil
}

func parseEcsInstanceSecurityGroupIds(groups []cloudservers.SecurityGroups) []string {
	result := make([]string, len(groups))

	for i, sg := range groups {
		result[i] = sg.ID
	}

	return result
}

func parseEcsInstanceSchedulerHintInfo(hints cloudservers.OsSchedulerHints) []map[string]interface{} {
	result := make([]map[string]interface{}, len(hints.Group))

	for i, val := range hints.Group {
		result[i] = map[string]interface{}{
			"group": val,
		}
	}

	return result
}

func parseEcsInstanceTagInfo(tags []string) map[string]interface{} {
	result := map[string]interface{}{}
	for _, tag := range tags {
		kv := strings.SplitN(tag, "=", 2)
		if len(kv) != 2 {
			log.Printf("[WARN] Invalid key/value format of tag: %s", tag)
			continue
		}
		result[kv[0]] = kv[1]
	}

	return result
}

func dataSourceComputeInstancesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	ecsClient, err := config.ComputeV1Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating FlexibleEngine ECS client: %s", err)
	}

	listOpts := &cloudservers.ListOpts{
		Name:   d.Get("name").(string),
		Flavor: d.Get("flavor_id").(string),
		IP:     d.Get("fixed_ip_v4").(string),
		Status: d.Get("status").(string),
	}

	pages, err := cloudservers.List(ecsClient, listOpts).AllPages()
	if err != nil {
		return diag.FromErr(err)
	}

	allServers, err := cloudservers.ExtractServers(pages)
	if err != nil {
		return diag.Errorf("Unable to retrieve cloud servers: %s ", err)
	}

	if len(allServers) < 1 {
		return diag.Errorf("Your query returned no results. " +
			"Please change your search criteria and try again.")
	}

	log.Printf("[DEBUG] fetching %d ecs instances.", len(allServers))
	servers, ids := filterCloudServers(d, allServers)
	d.SetId(hashcode.Strings(ids))

	return setComputeInstancesParams(d, servers, meta)
}
