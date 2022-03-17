package flexibleengine

import (
	"fmt"
	"log"
	"strconv"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/ecs/v1/auto_recovery"
	"github.com/chnsz/golangsdk/openstack/ecs/v1/tags"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceECSAutoRecoveryV1Read(d *schema.ResourceData, meta interface{}, instanceID string) (bool, error) {
	config := meta.(*Config)
	client, err := config.ComputeV1Client(GetRegion(d, config))
	if err != nil {
		return false, fmt.Errorf("Error creating FlexibleEngine client: %s", err)
	}

	rId := instanceID

	r, err := auto_recovery.Get(client, rId).Extract()
	if err != nil {
		return false, err
	}
	log.Printf("[DEBUG] Retrieved ECS-AutoRecovery:%#v of instance:%s", rId, r)
	return strconv.ParseBool(r.SupportAutoRecovery)
}

func setAutoRecoveryForInstance(d *schema.ResourceData, meta interface{}, instanceID string, ar bool) error {
	config := meta.(*Config)
	client, err := config.ComputeV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine client: %s", err)
	}

	rId := instanceID

	updateOpts := auto_recovery.UpdateOpts{SupportAutoRecovery: strconv.FormatBool(ar)}

	timeout := d.Timeout(schema.TimeoutUpdate)

	log.Printf("[DEBUG] Setting ECS-AutoRecovery for instance:%s with options: %#v", rId, updateOpts)
	err = resource.Retry(timeout, func() *resource.RetryError {
		err := auto_recovery.Update(client, rId, updateOpts)
		if err != nil {
			return checkForRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("Error setting ECS-AutoRecovery for instance%s: %s", rId, err)
	}
	return nil
}

func resourceECSTagsV1Read(d *schema.ResourceData, meta interface{}, instanceID string) (map[string]string, error) {
	tagmap := make(map[string]string)

	config := meta.(*Config)
	client, err := config.ComputeV1Client(GetRegion(d, config))
	if err != nil {
		return tagmap, fmt.Errorf("Error creating FlexibleEngine compute v1 client: %s", err)
	}

	ecsTaglist, err := tags.Get(client, instanceID).Extract()
	if err != nil {
		return tagmap, fmt.Errorf("Error fetching ECS instance tags: %s", err)
	}

	for _, val := range ecsTaglist.Tags {
		tagmap[val.Key] = val.Value
	}
	return tagmap, nil
}

func setTagsForInstance(client *golangsdk.ServiceClient, instanceID string, tagmap map[string]interface{}) error {
	taglist := expandInstanceTags(tagmap)

	createOpts := tags.BatchOpts{Action: tags.ActionCreate, Tags: taglist}
	createTags := tags.BatchAction(client, instanceID, createOpts)
	if createTags.Err != nil {
		return fmt.Errorf("Error creating ecs instance v1 tags: %s", createTags.Err)
	}

	return nil
}

func deleteTagsForInstance(client *golangsdk.ServiceClient, instanceID string, tagmap map[string]interface{}) error {
	taglist := expandInstanceTags(tagmap)

	deleteOpts := tags.BatchOpts{Action: tags.ActionDelete, Tags: taglist}
	deleteTags := tags.BatchAction(client, instanceID, deleteOpts)
	if deleteTags.Err != nil {
		return fmt.Errorf("Error deleting ecs instance v1 tags: %s", deleteTags.Err)
	}

	return nil
}

// expandResourceTags returns the tags for the given map of data.
func expandInstanceTags(tagmap map[string]interface{}) []tags.Tag {
	var taglist []tags.Tag

	for k, v := range tagmap {
		tag := tags.Tag{
			Key:   k,
			Value: v.(string),
		}
		taglist = append(taglist, tag)
	}

	return taglist
}
