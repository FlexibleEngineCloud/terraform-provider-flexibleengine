package flexibleengine

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/openstack/sfs/v2/shares"
)

func resourceSFSFileSystemV2() *schema.Resource {
	return &schema.Resource{
		Create: resourceSFSFileSystemV2Create,
		Read:   resourceSFSFileSystemV2Read,
		Update: resourceSFSFileSystemV2Update,
		Delete: resourceSFSFileSystemV2Delete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"share_proto": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "NFS",
			},
			"size": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"is_public": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"metadata": &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
			},
			"availability_zone": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"access_level": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"access_type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "cert",
			},
			"access_to": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"share_access_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"access_rule_status": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"host": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"volume_type": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"export_location": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceSFSFileSystemV2Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	sfsClient, err := config.sfsV2Client(GetRegion(d, config))

	if err != nil {
		return fmt.Errorf("Error creating Flexibleengine File Share Client: %s", err)
	}

	createOpts := shares.CreateOpts{
		ShareProto:       d.Get("share_proto").(string),
		Size:             d.Get("size").(int),
		Name:             d.Get("name").(string),
		Description:      d.Get("description").(string),
		IsPublic:         d.Get("is_public").(bool),
		Metadata:         resourceSFSMetadataV2(d),
		AvailabilityZone: d.Get("availability_zone").(string),
	}

	create, err := shares.Create(sfsClient, createOpts).Extract()

	if err != nil {
		return fmt.Errorf("Error creating Flexibleengine File Share: %s", err)
	}
	d.SetId(create.ID)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"creating"},
		Target:     []string{"available"},
		Refresh:    waitForSFSFileActive(sfsClient, create.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	_, StateErr := stateConf.WaitForState()
	if StateErr != nil {
		return fmt.Errorf("Error waiting for Share File (%s) to become ready: %s ", d.Id(), StateErr)
	}

	grantAccessOpts := shares.GrantAccessOpts{
		AccessLevel: d.Get("access_level").(string),
		AccessType:  d.Get("access_type").(string),
		AccessTo:    d.Get("access_to").(string),
	}

	grant, accessErr := shares.GrantAccess(sfsClient, d.Id(), grantAccessOpts).ExtractAccess()

	if accessErr != nil {
		return fmt.Errorf("Error applying access rules to share file : %s", accessErr)
	}

	log.Printf("[DEBUG] Applied access rule (%s) to share file %s", grant.ID, d.Id())

	return resourceSFSFileSystemV2Read(d, meta)

}

func resourceSFSFileSystemV2Read(d *schema.ResourceData, meta interface{}) error {

	config := meta.(*Config)
	sfsClient, err := config.sfsV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating Flexibleengine File Share: %s", err)
	}

	n, err := shares.Get(sfsClient, d.Id()).Extract()
	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error retrieving Flexibleengine Shares: %s", err)
	}

	d.Set("name", n.Name)
	d.Set("share_proto", n.ShareProto)
	d.Set("status", n.Status)
	d.Set("size", n.Size)
	d.Set("description", n.Description)
	d.Set("share_type", n.ShareType)
	d.Set("volume_type", n.VolumeType)
	d.Set("is_public", n.IsPublic)
	d.Set("availability_zone", n.AvailabilityZone)
	d.Set("region", GetRegion(d, config))
	d.Set("export_location", n.ExportLocation)
	d.Set("host", n.Host)
	d.Set("links", n.Links)

	// NOTE: This tries to remove system metadata.
	md := make(map[string]string)
	var sys_keys = [2]string{"enterprise_project_id", "share_used"}

OUTER:
	for key, val := range n.Metadata {
		if strings.HasPrefix(key, "#sfs") {
			continue
		}
		for i := range sys_keys {
			if key == sys_keys[i] {
				continue OUTER
			}
		}
		md[key] = val
	}
	d.Set("metadata", md)

	rules, err := shares.ListAccessRights(sfsClient, d.Id()).ExtractAccessRights()

	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error retrieving Flexibleengine Shares: %s", err)
	}

	if len(rules) > 0 {
		rule := rules[0]
		d.Set("share_access_id", rule.ID)
		d.Set("access_rule_status", rule.State)
		d.Set("access_to", rule.AccessTo)
		d.Set("access_type", rule.AccessType)
		d.Set("access_level", rule.AccessLevel)
	}
	return nil
}

func resourceSFSFileSystemV2Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	sfsClient, err := config.sfsV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error updating Flexibleengine Share File: %s", err)
	}
	var updateOpts shares.UpdateOpts

	if d.HasChange("description") || d.HasChange("name") {
		updateOpts.DisplayName = d.Get("name").(string)
		updateOpts.DisplayDescription = d.Get("description").(string)

		_, err = shares.Update(sfsClient, d.Id(), updateOpts).Extract()
		if err != nil {
			return fmt.Errorf("Error updating Flexibleengine Share File: %s", err)
		}
	}
	if d.HasChange("access_to") || d.HasChange("access_level") || d.HasChange("access_type") {
		deleteAccessOpts := shares.DeleteAccessOpts{AccessID: d.Get("share_access_id").(string)}
		deny := shares.DeleteAccess(sfsClient, d.Id(), deleteAccessOpts)
		if deny.Err != nil {
			return fmt.Errorf("Error changing access rules for share file : %s", deny.Err)
		}

		grantAccessOpts := shares.GrantAccessOpts{
			AccessLevel: d.Get("access_level").(string),
			AccessType:  d.Get("access_type").(string),
			AccessTo:    d.Get("access_to").(string),
		}

		log.Printf("[DEBUG] Grant Access Rules: %#v", grantAccessOpts)
		_, accessErr := shares.GrantAccess(sfsClient, d.Id(), grantAccessOpts).ExtractAccess()

		if accessErr != nil {
			return fmt.Errorf("Error changing access rules for share file : %s", accessErr)
		}
	}

	if d.HasChange("size") {
		old, newsize := d.GetChange("size")
		if old.(int) < newsize.(int) {
			expandOpts := shares.ExpandOpts{OSExtend: shares.OSExtendOpts{NewSize: newsize.(int)}}
			expand := shares.Expand(sfsClient, d.Id(), expandOpts)
			if expand.Err != nil {
				return fmt.Errorf("Error Expanding Flexibleengine Share File size: %s", expand.Err)
			}
		} else {
			shrinkOpts := shares.ShrinkOpts{OSShrink: shares.OSShrinkOpts{NewSize: newsize.(int)}}
			shrink := shares.Shrink(sfsClient, d.Id(), shrinkOpts)
			if shrink.Err != nil {
				return fmt.Errorf("Error Shrinking Flexibleengine Share File size: %s", shrink.Err)
			}
		}
	}

	return resourceSFSFileSystemV2Read(d, meta)
}

func resourceSFSFileSystemV2Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	sfsClient, err := config.sfsV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating Flexibleengine Shared File: %s", err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"available", "deleting"},
		Target:     []string{"deleted"},
		Refresh:    waitForSFSFileDelete(sfsClient, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Error deleting Flexibleengine Share File: %s", err)
	}

	d.SetId("")
	return nil
}

func waitForSFSFileActive(sfsClient *golangsdk.ServiceClient, shareID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		n, err := shares.Get(sfsClient, shareID).Extract()
		if err != nil {
			return nil, "", err
		}

		if n.Status == "error" {
			return n, n.Status, nil
		}
		return n, n.Status, nil
	}
}

func waitForSFSFileDelete(sfsClient *golangsdk.ServiceClient, shareId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {

		r, err := shares.Get(sfsClient, shareId).Extract()

		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				log.Printf("[INFO] Successfully deleted Flexibleengine shared File %s", shareId)
				return r, "deleted", nil
			}
			return r, "available", err
		}
		err = shares.Delete(sfsClient, shareId).ExtractErr()

		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				log.Printf("[INFO] Successfully deleted Flexibleengine shared File %s", shareId)
				return r, "deleted", nil
			}
			if errCode, ok := err.(golangsdk.ErrUnexpectedResponseCode); ok {
				if errCode.Actual == 409 {
					return r, "available", nil
				}
			}
			return r, "available", err
		}

		return r, r.Status, nil
	}
}

func resourceSFSMetadataV2(d *schema.ResourceData) map[string]string {
	m := make(map[string]string)
	for key, val := range d.Get("metadata").(map[string]interface{}) {
		m[key] = val.(string)
	}
	return m
}
