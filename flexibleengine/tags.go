package flexibleengine

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/openstack/common/tags"
)

// tagsSchema returns the schema to use for tags.
func tagsSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeMap,
		Optional: true,
	}
}

// UpdateResourceTags is a helper to update the tags for a resource.
// It expects the tags field to be named "tags"
func UpdateResourceTags(conn *golangsdk.ServiceClient, d *schema.ResourceData, resourceType string) error {
	if d.HasChange("tags") {
		oRaw, nRaw := d.GetChange("tags")
		oMap := oRaw.(map[string]interface{})
		nMap := nRaw.(map[string]interface{})

		// remove old tags
		if len(oMap) > 0 {
			taglist := expandResourceTags(oMap)
			err := tags.Delete(conn, resourceType, d.Id(), taglist).ExtractErr()
			if err != nil {
				return err
			}
		}

		// set new tags
		if len(nMap) > 0 {
			taglist := expandResourceTags(nMap)
			err := tags.Create(conn, resourceType, d.Id(), taglist).ExtractErr()
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// tagsToMap returns the list of tags into a map.
func tagsToMap(tags []tags.ResourceTag) map[string]string {
	result := make(map[string]string)
	for _, val := range tags {
		result[val.Key] = val.Value
	}

	return result
}

// expandResourceTags returns the tags for the given map of data.
func expandResourceTags(tagmap map[string]interface{}) []tags.ResourceTag {
	var taglist []tags.ResourceTag

	for k, v := range tagmap {
		tag := tags.ResourceTag{
			Key:   k,
			Value: v.(string),
		}
		taglist = append(taglist, tag)
	}

	return taglist
}
