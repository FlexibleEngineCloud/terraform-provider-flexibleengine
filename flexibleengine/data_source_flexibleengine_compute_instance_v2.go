package flexibleengine

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/huaweicloud/golangsdk/openstack/common/tags"
	"github.com/huaweicloud/golangsdk/openstack/ecs/v1/block_devices"
	"github.com/huaweicloud/golangsdk/openstack/ecs/v1/cloudservers"
	"github.com/huaweicloud/golangsdk/openstack/networking/v2/ports"
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
	ecsClient, err := config.computeV1Client(GetRegion(d, config))
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
	instanceVolumes := []map[string]interface{}{}
	if len(server.VolumeAttached) > 0 {
		for _, b := range server.VolumeAttached {
			va, err := block_devices.Get(ecsClient, d.Id(), b.ID).Extract()
			if err != nil {
				return err
			}
			log.Printf("[DEBUG] Retrieved block device %s: %#v", d.Id(), va)
			v := map[string]interface{}{
				"uuid":        b.ID,
				"boot_index":  va.BootIndex,
				"size":        va.Size,
				"pci_address": va.PciAddress,
			}
			instanceVolumes = append(instanceVolumes, v)
			if va.BootIndex == 0 {
				d.Set("system_disk_id", b.ID)
			}
		}
		d.Set("block_device", instanceVolumes)
	}

	// Set instance tags
	resourceTags, err := tags.Get(ecsClient, "servers", d.Id()).Extract()
	if err == nil {
		tagmap := tagsToMap(resourceTags.Tags)
		d.Set("tags", tagmap)
	} else {
		log.Printf("[WARN] Error fetching tags of ecs instance %s: %s", d.Id(), err)
	}

	return nil
}

// flattenComputeNetworks collects instance network information
func flattenComputeNetworks(
	d *schema.ResourceData, meta interface{}, server *cloudservers.CloudServer) ([]map[string]interface{}, string) {

	config := meta.(*Config)
	networkingClient, err := config.networkingV2Client(GetRegion(d, config))
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
