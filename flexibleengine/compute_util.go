package flexibleengine

import (
	"fmt"
	"log"
	"os"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/blockstorage/v2/volumes"
	bms "github.com/chnsz/golangsdk/openstack/bms/v2/servers"
	"github.com/chnsz/golangsdk/openstack/compute/v2/flavors"
	"github.com/chnsz/golangsdk/openstack/compute/v2/servers"
	"github.com/chnsz/golangsdk/openstack/ecs/v1/block_devices"
	"github.com/chnsz/golangsdk/openstack/ecs/v1/cloudservers"
	"github.com/chnsz/golangsdk/openstack/imageservice/v2/images"
	"github.com/chnsz/golangsdk/openstack/networking/v2/networks"
	"github.com/chnsz/golangsdk/openstack/networking/v2/ports"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceComputeSecGroupsV2(d *schema.ResourceData) []string {
	rawSecGroups := d.Get("security_groups").(*schema.Set).List()
	secgroups := make([]string, len(rawSecGroups))
	for i, raw := range rawSecGroups {
		secgroups[i] = raw.(string)
	}
	return secgroups
}

func resourceComputeMetadataV2(d *schema.ResourceData) map[string]string {
	m := make(map[string]string)
	for key, val := range d.Get("metadata").(map[string]interface{}) {
		m[key] = val.(string)
	}
	return m
}

func checkBlockDeviceConfig(d *schema.ResourceData) error {
	if vL, ok := d.GetOk("block_device"); ok {
		for _, v := range vL.([]interface{}) {
			vM := v.(map[string]interface{})

			if vM["source_type"] != "blank" && vM["uuid"] == "" {
				return fmt.Errorf("You must specify a uuid for %s block device types", vM["source_type"])
			}

			if vM["source_type"] == "image" && vM["destination_type"] == "volume" {
				if vM["volume_size"] == 0 {
					return fmt.Errorf("You must specify a volume_size when creating a volume from an image")
				}
			}

			if vM["source_type"] == "blank" && vM["destination_type"] == "local" {
				if vM["volume_size"] == 0 {
					return fmt.Errorf("You must specify a volume_size when creating a blank block device")
				}
			}
		}
	}

	return nil
}

func getComputeFlavorID(client *golangsdk.ServiceClient, d *schema.ResourceData) (string, error) {
	flavorID := d.Get("flavor_id").(string)

	if flavorID != "" {
		return flavorID, nil
	}

	flavorName := d.Get("flavor_name").(string)
	return flavors.IDFromName(client, flavorName)
}

// getInstanceImageID determines the Image ID using the following rules:
// If a bootable block_device was specified, ignore the image altogether.
// If an image_id was specified, use it.
// If an image_name was specified, look up the image ID, report if error.
func getInstanceImageID(client *golangsdk.ServiceClient, d *schema.ResourceData) (string, error) {

	if vL, ok := d.GetOk("block_device"); ok {
		needImage := false
		for _, v := range vL.([]interface{}) {
			vM := v.(map[string]interface{})
			if vM["source_type"] == "image" && vM["destination_type"] == "local" {
				needImage = true
			}
		}
		if !needImage {
			return "", nil
		}
	}

	if imageID := d.Get("image_id").(string); imageID != "" {
		return imageID, nil
	}
	// try the OS_IMAGE_ID environment variable
	if v := os.Getenv("OS_IMAGE_ID"); v != "" {
		return v, nil
	}

	imageName := d.Get("image_name").(string)
	if imageName == "" {
		// try the OS_IMAGE_NAME environment variable
		if v := os.Getenv("OS_IMAGE_NAME"); v != "" {
			imageName = v
		}
	}

	if imageName != "" {
		img, err := getImage(client, "", imageName)
		if err != nil {
			return "", err
		}
		return img.ID, nil
	}

	return "", fmt.Errorf("neither a boot device, image ID, or image name were able to be determined")
}

// getBMSImageID determines the Image ID using the following rules:
// If an image_id was specified, use it.
// If an image_name was specified, look up the image ID, report if error.
func getBMSImageID(client *golangsdk.ServiceClient, d *schema.ResourceData) (string, error) {

	if imageID := d.Get("image_id").(string); imageID != "" {
		return imageID, nil
	}
	// try the OS_IMAGE_ID environment variable
	if v := os.Getenv("OS_IMAGE_ID"); v != "" {
		return v, nil
	}

	imageName := d.Get("image_name").(string)
	if imageName == "" {
		// try the OS_IMAGE_NAME environment variable
		if v := os.Getenv("OS_IMAGE_NAME"); v != "" {
			imageName = v
		}
	}

	if imageName != "" {
		img, err := getImage(client, "", imageName)
		if err != nil {
			return "", err
		}
		return img.ID, nil
	}

	return "", fmt.Errorf("neither a image ID, or image name were able to be determined")
}

func getImage(client *golangsdk.ServiceClient, id, name string) (*images.Image, error) {
	listOpts := &images.ListOpts{
		ID:    id,
		Name:  name,
		Limit: 1,
	}
	allPages, err := images.List(client, listOpts).AllPages()
	if err != nil {
		return nil, fmt.Errorf("Unable to query images: %s", err)
	}

	allImages, err := images.ExtractImages(allPages)
	if err != nil {
		return nil, fmt.Errorf("Unable to retrieve images: %s", err)
	}

	if len(allImages) < 1 {
		return nil, fmt.Errorf("Unable to find images %s: Maybe not existed", id)
	}

	img := allImages[0]
	if id != "" && img.ID != id {
		return nil, fmt.Errorf("Unexpected images ID")
	}
	if name != "" && img.Name != name {
		return nil, fmt.Errorf("Unexpected images Name")
	}

	log.Printf("[DEBUG] Retrieved Image %s: %#v", id, img)
	return &img, nil
}

func setInstanceImageInfo(d *schema.ResourceData, client *golangsdk.ServiceClient, imageID string) error {
	// If block_device was used, an Image does not need to be specified, unless an image/local
	// combination was used. This emulates normal boot behavior. Otherwise, ignore the image altogether.
	if vL, ok := d.GetOk("block_device"); ok {
		needImage := false
		for _, v := range vL.([]interface{}) {
			vM := v.(map[string]interface{})
			if vM["source_type"] == "image" && vM["destination_type"] == "local" {
				needImage = true
			}
		}
		if !needImage {
			d.Set("image_id", "Attempt to boot from volume - no image supplied")
			return nil
		}
	}

	if imageID != "" {
		d.Set("image_id", imageID)
		image, err := images.Get(client, imageID).Extract()
		if err != nil {
			// If the image name can't be found, set the value to "Image not found".
			// The most likely scenario is that the image no longer exists in the Image Service
			// but the instance still has a record from when it existed.
			d.Set("image_name", "Image not found")
			return nil
		}
		d.Set("image_name", image.Name)
	}

	return nil
}

func setBMSImageInfo(client *golangsdk.ServiceClient, server *bms.Server, d *schema.ResourceData) error {
	imageID := server.Image.ID
	if imageID != "" {
		d.Set("image_id", imageID)
		image, err := images.Get(client, imageID).Extract()
		if err != nil {
			// If the image name can't be found, set the value to "Image not found".
			// The most likely scenario is that the image no longer exists in the Image Service
			// but the instance still has a record from when it existed.
			d.Set("image_name", "Image not found")
			return nil
		}
		d.Set("image_name", image.Name)
	}

	return nil
}

// computeV2StateRefreshFunc returns a resource.StateRefreshFunc that is used to watch
// an FlexibleEngine instance.
func computeV2StateRefreshFunc(client *golangsdk.ServiceClient, instanceID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		s, err := servers.Get(client, instanceID).Extract()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return s, "DELETED", nil
			}
			return nil, "", err
		}

		// get fault message when status is ERROR
		if s.Status == "ERROR" {
			fault := fmt.Errorf("[error code: %d, message: %s]", s.Fault.Code, s.Fault.Message)
			return s, "ERROR", fault
		}
		return s, s.Status, nil
	}
}

// getInstanceNetworkInfo will query for network information in order to make
// an accurate determination of a network's name and a network's ID.
func getInstanceNetworkInfo(
	d *schema.ResourceData, meta interface{}, queryType, queryTerm string) (map[string]string, error) {

	config := meta.(*Config)
	networkClient, err := config.NetworkingV2Client(GetRegion(d, config))
	if err != nil {
		return nil, fmt.Errorf("Error creating FlexibleEngine networking client: %s", err)
	}

	// If a port was specified, convert it to the network ID
	// and then query the network as if a network ID was originally used.
	if queryType == "port" {
		portID := queryTerm
		port, err := ports.Get(networkClient, portID).Extract()
		if err != nil {
			return nil, fmt.Errorf("Could not find any matching port for %s", portID)
		}

		queryType = "id"
		queryTerm = port.NetworkID
	}

	listOpts := networks.ListOpts{
		Status: "ACTIVE",
	}

	switch queryType {
	case "name":
		listOpts.Name = queryTerm
	default:
		listOpts.ID = queryTerm
	}

	allPages, err := networks.List(networkClient, listOpts).AllPages()
	if err != nil {
		return nil, fmt.Errorf("Unable to retrieve networks from the Network API: %s", err)
	}

	allNetworks, err := networks.ExtractNetworks(allPages)
	if err != nil {
		return nil, fmt.Errorf("Unable to retrieve networks from the Network API: %s", err)
	}

	var network networks.Network
	switch len(allNetworks) {
	case 0:
		return nil, fmt.Errorf("Could not find any matching network for %s %s", queryType, queryTerm)
	case 1:
		network = allNetworks[0]
	default:
		// may happened when querying by "name"
		return nil, fmt.Errorf("More than one network found for %s %s", queryType, queryTerm)
	}

	networkInfo := map[string]string{
		"uuid": network.ID,
		"name": network.Name,
	}

	log.Printf("[DEBUG] getInstanceNetworkInfo: %#v", networkInfo)
	return networkInfo, nil
}

func flattenInstanceVolumeAttached(
	d *schema.ResourceData, meta interface{}, server *cloudservers.CloudServer) ([]map[string]interface{}, string, error) {

	config := meta.(*Config)
	ecsClient, err := config.ComputeV1Client(GetRegion(d, config))
	blockStorageClient, err := config.BlockStorageV2Client(GetRegion(d, config))
	if err != nil {
		return nil, "", fmt.Errorf("Error creating FlexibleEngine client: %s", err)
	}

	var systemDiskID string = ""
	bds := make([]map[string]interface{}, len(server.VolumeAttached))
	for i, b := range server.VolumeAttached {
		// retrieve volume `size` and `type`
		volumeInfo, err := volumes.Get(blockStorageClient, b.ID).Extract()
		if err != nil {
			return nil, "", err
		}
		log.Printf("[DEBUG] Retrieved volume %s: %#v", b.ID, volumeInfo)

		// retrieve volume `pci_address`
		va, err := block_devices.Get(ecsClient, server.ID, b.ID).Extract()
		if err != nil {
			return nil, "", err
		}
		log.Printf("[DEBUG] Retrieved block device %s: %#v", b.ID, va)

		bds[i] = map[string]interface{}{
			"uuid":        b.ID,
			"size":        volumeInfo.Size,
			"type":        volumeInfo.VolumeType,
			"boot_index":  va.BootIndex,
			"pci_address": va.PciAddress,
		}

		if va.BootIndex == 0 {
			systemDiskID = b.ID
		}
	}
	return bds, systemDiskID, nil
}
