package flexibleengine

import (
	"fmt"
	"log"

	"github.com/chnsz/golangsdk/openstack/common/tags"
	"github.com/chnsz/golangsdk/openstack/compute/v2/servers"
	"github.com/chnsz/golangsdk/openstack/ecs/v1/cloudservers"
	"github.com/chnsz/golangsdk/openstack/networking/v2/ports"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceComputeInstance() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceComputeInstanceRead,

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
			"availability_zone": {
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
			"system_disk_id": {
				Type:     schema.TypeString,
				Computed: true,
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
			"block_device": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"uuid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"boot_index": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"pci_address": {
							Type:     schema.TypeString,
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
	}
}

func dataSourceComputeInstanceRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	ecsClient, err := config.ComputeV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine ECS client: %s", err)
	}

	listOpts := &cloudservers.ListOpts{
		Name:   d.Get("name").(string),
		Flavor: d.Get("flavor_id").(string),
		IP:     d.Get("fixed_ip_v4").(string),
	}

	pages, err := cloudservers.List(ecsClient, listOpts).AllPages()
	if err != nil {
		return err
	}

	allServers, err := cloudservers.ExtractServers(pages)
	if err != nil {
		return fmt.Errorf("Unable to retrieve cloud servers: %s ", err)
	}

	if len(allServers) < 1 {
		return fmt.Errorf("Your query returned no results. " +
			"Please change your search criteria and try again.")
	}
	if len(allServers) > 1 {
		return fmt.Errorf("Your query returned more than one result. " +
			"Please try a more specific search criteria.")
	}

	server := allServers[0]
	log.Printf("[DEBUG] fetching the ecs instance: %#v", server)

	d.SetId(server.ID)
	d.Set("region", GetRegion(d, config))
	d.Set("availability_zone", server.AvailabilityZone)
	d.Set("name", server.Name)
	d.Set("status", server.Status)

	flavorInfo := server.Flavor
	d.Set("flavor_id", flavorInfo.ID)
	d.Set("flavor_name", flavorInfo.Name)
	d.Set("image_id", server.Image.ID)

	metaData := server.Metadata
	if metaData.ImageName != "" {
		d.Set("image_name", metaData.ImageName)
	}
	if server.KeyName != "" {
		d.Set("key_pair", server.KeyName)
	}
	if server.UserData != "" {
		d.Set("user_data", server.UserData)
	}

	// set security groups
	secGrpNames := make([]string, len(server.SecurityGroups))
	for i, sg := range server.SecurityGroups {
		secGrpNames[i] = sg.Name
	}
	d.Set("security_groups", secGrpNames)

	// set os:scheduler_hints
	osHints := server.OsSchedulerHints
	if len(osHints.Group) > 0 {
		schedulerHints := make([]map[string]interface{}, len(osHints.Group))
		for i, v := range osHints.Group {
			schedulerHints[i] = map[string]interface{}{
				"group": v,
			}
		}
		d.Set("scheduler_hints", schedulerHints)
	}

	// Set the instance network and address information
	networks, eip := flattenComputeNetworks(d, meta, &server)
	if err := d.Set("network", networks); err != nil {
		log.Printf("[WARN] Error setting network of ecs instance %s: %s", d.Id(), err)
	}
	if eip != "" {
		d.Set("floating_ip", eip)
	}

	// Set volume attached
	if len(server.VolumeAttached) > 0 {
		volumes, rootID, err := flattenInstanceVolumeAttached(d, config, &server)
		if err != nil {
			return nil
		}
		d.Set("block_device", volumes)
		d.Set("system_disk_id", rootID)
	}

	// Set instance tags
	resourceTags, err := tags.Get(ecsClient, "servers", d.Id()).Extract()
	if err == nil {
		tagmap := tagsToMap(resourceTags.Tags)
		d.Set("tags", tagmap)
	} else {
		log.Printf("[WARN] Error fetching tags of ecs instance %s: %s", d.Id(), err)
	}

	// Set meta
	computeClient, err := config.ComputeV2Client(GetRegion(d, config))
	if err == nil {
		novaResp, err := servers.Get(computeClient, d.Id()).Extract()
		if err == nil {
			d.Set("metadata", novaResp.Metadata)
		} else {
			log.Printf("[WARN] Error fetching metadata of ecs instance %s: %s", d.Id(), err)
		}
	}

	return nil
}

// flattenComputeNetworks collects instance network information
func flattenComputeNetworks(
	d *schema.ResourceData, meta interface{}, server *cloudservers.CloudServer) ([]map[string]interface{}, string) {

	config := meta.(*Config)
	networkingClient, err := config.NetworkingV2Client(GetRegion(d, config))
	if err != nil {
		log.Printf("[ERROR] failed to create FlexibleEngine networking client: %s", err)
		return nil, ""
	}

	publicIP := ""
	networks := []map[string]interface{}{}

	for _, addresses := range server.Addresses {
		for _, addr := range addresses {
			if addr.Type == "floating" {
				publicIP = addr.Addr
				continue
			}

			// get networkID
			var networkID string
			p, err := ports.Get(networkingClient, addr.PortID).Extract()
			if err != nil {
				networkID = ""
				log.Printf("[DEBUG] failed to fetch port %s", addr.PortID)
			} else {
				networkID = p.NetworkID
			}

			v := map[string]interface{}{
				"uuid": networkID,
				"port": addr.PortID,
				"mac":  addr.MacAddr,
			}
			if addr.Version == "6" {
				v["fixed_ip_v6"] = addr.Addr
			} else {
				v["fixed_ip_v4"] = addr.Addr
			}

			networks = append(networks, v)
		}
	}

	log.Printf("[DEBUG] flatten Instance Networks: %#v", networks)
	return networks, publicIP
}
