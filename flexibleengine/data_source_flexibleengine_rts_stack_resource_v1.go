package flexibleengine

import (
	"fmt"

	"github.com/chnsz/golangsdk/openstack/rts/v1/stackresources"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceRTSStackResourcesV1() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceRTSStackResourcesV1Read,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"stack_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"resource_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"logical_resource_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"required_by": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},
			"resource_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"resource_status_reason": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"physical_resource_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"resource_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func dataSourceRTSStackResourcesV1Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	orchestrationClient, err := orchestrationV1Client(config, GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating RTS client: %s", err)
	}

	listOpts := stackresources.ListOpts{
		Name:       d.Get("resource_name").(string),
		PhysicalID: d.Get("physical_resource_id").(string),
		Type:       d.Get("resource_type").(string),
	}

	refinedResources, err := stackresources.List(orchestrationClient, d.Get("stack_name").(string), listOpts)
	if err != nil {
		return fmt.Errorf("Unable to retrieve Stack Resources: %s", err)
	}

	if len(refinedResources) < 1 {
		return fmt.Errorf("No matching resource found. " +
			"Please change your search criteria and try again.")
	}

	if len(refinedResources) > 1 {
		return fmt.Errorf("Multiple resources matched; use additional constraints to reduce matches to a single resource")
	}

	stackResource := refinedResources[0]
	d.SetId(stackResource.PhysicalID)

	d.Set("resource_name", stackResource.Name)
	d.Set("resource_status", stackResource.Status)
	d.Set("logical_resource_id", stackResource.LogicalID)
	d.Set("physical_resource_id", stackResource.PhysicalID)
	d.Set("required_by", stackResource.RequiredBy)
	d.Set("resource_status_reason", stackResource.StatusReason)
	d.Set("resource_type", stackResource.Type)
	d.Set("region", GetRegion(d, config))
	return nil
}
