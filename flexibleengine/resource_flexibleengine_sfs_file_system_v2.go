package flexibleengine

import (
	"fmt"
	"log"
	"time"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/sfs/v2/shares"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"share_proto": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "NFS",
			},
			"size": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"is_public": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"metadata": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
			},
			"availability_zone": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"access_level": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"access_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"access_to": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"share_access_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"access_rule_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"volume_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"export_location": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"access_rules": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"access_rule_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"access_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"access_to": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"access_level": {
							Type:     schema.TypeString,
							Computed: true,
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

func resourceSFSFileSystemV2Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	sfsClient, err := config.SfsV2Client(GetRegion(d, config))

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

	// specified the "access_to" field, apply first access rule to share file
	if _, ok := d.GetOk("access_to"); ok {
		grantAccessOpts := shares.GrantAccessOpts{
			AccessTo: d.Get("access_to").(string),
		}

		if v, ok := d.GetOk("access_level"); ok {
			grantAccessOpts.AccessLevel = v.(string)
		} else {
			grantAccessOpts.AccessLevel = "rw"
		}

		if v, ok := d.GetOk("access_type"); ok {
			grantAccessOpts.AccessType = v.(string)
		} else {
			grantAccessOpts.AccessType = "cert"
		}

		grant, accessErr := shares.GrantAccess(sfsClient, d.Id(), grantAccessOpts).ExtractAccess()
		if accessErr != nil {
			return fmt.Errorf("Error applying access rule to share file : %s", accessErr)
		}

		log.Printf("[DEBUG] Applied access rule (%s) to share file %s", grant.ID, d.Id())
		d.Set("share_access_id", grant.ID)
	}

	return resourceSFSFileSystemV2Read(d, meta)

}

func resourceSFSFileSystemV2Read(d *schema.ResourceData, meta interface{}) error {

	config := meta.(*Config)
	sfsClient, err := config.SfsV2Client(GetRegion(d, config))
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
	d.Set("size", n.Size)
	d.Set("description", n.Description)
	d.Set("volume_type", n.VolumeType)
	d.Set("is_public", n.IsPublic)
	d.Set("availability_zone", n.AvailabilityZone)
	d.Set("region", GetRegion(d, config))
	d.Set("export_location", n.ExportLocation)

	// NOTE: only support the following metadata key
	var metaKeys = [3]string{"#sfs_crypt_key_id", "#sfs_crypt_domain_id", "#sfs_crypt_alias"}
	md := make(map[string]string)

	for key, val := range n.Metadata {
		for i := range metaKeys {
			if key == metaKeys[i] {
				md[key] = val
				break
			}
		}
	}
	d.Set("metadata", md)

	// list access rules
	rules, err := shares.ListAccessRights(sfsClient, d.Id()).ExtractAccessRights()
	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error retrieving Flexibleengine Shares: %s", err)
	}

	var ruleExist bool
	accessID := d.Get("share_access_id").(string)
	allAccessRules := make([]map[string]interface{}, 0, len(rules))
	for _, rule := range rules {
		acessRule := map[string]interface{}{
			"access_rule_id": rule.ID,
			"access_level":   rule.AccessLevel,
			"access_type":    rule.AccessType,
			"access_to":      rule.AccessTo,
			"status":         rule.State,
		}
		allAccessRules = append(allAccessRules, acessRule)

		// find share_access_id
		if accessID != "" && rule.ID == accessID {
			d.Set("access_rule_status", rule.State)
			d.Set("access_to", rule.AccessTo)
			d.Set("access_type", rule.AccessType)
			d.Set("access_level", rule.AccessLevel)
			ruleExist = true
		}
	}

	if accessID != "" && !ruleExist {
		log.Printf("[WARN] access rule (%s) of share file %s was not exist!", accessID, d.Id())
		d.Set("share_access_id", "")
	}
	d.Set("access_rules", allAccessRules)

	if len(rules) != 0 {
		d.Set("status", n.Status)
	} else {
		// The file system is not bind with any VPC.
		d.Set("status", "unavailable")
	}

	return nil
}

func resourceSFSFileSystemV2Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	sfsClient, err := config.SfsV2Client(GetRegion(d, config))
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
		ruleID := d.Get("share_access_id").(string)
		if ruleID != "" {
			deleteAccessOpts := shares.DeleteAccessOpts{AccessID: ruleID}
			deny := shares.DeleteAccess(sfsClient, d.Id(), deleteAccessOpts)
			if deny.Err != nil {
				return fmt.Errorf("Error changing access rules for share file : %s", deny.Err)
			}
			d.Set("share_access_id", "")
		}

		if v, ok := d.GetOk("access_to"); ok {
			grantAccessOpts := shares.GrantAccessOpts{
				AccessTo: v.(string),
			}

			if v, ok := d.GetOk("access_level"); ok {
				grantAccessOpts.AccessLevel = v.(string)
			} else {
				grantAccessOpts.AccessLevel = "rw"
			}
			if v, ok := d.GetOk("access_type"); ok {
				grantAccessOpts.AccessType = v.(string)
			} else {
				grantAccessOpts.AccessType = "cert"
			}

			log.Printf("[DEBUG] Grant Access Rules: %#v", grantAccessOpts)
			grant, accessErr := shares.GrantAccess(sfsClient, d.Id(), grantAccessOpts).ExtractAccess()
			if accessErr != nil {
				return fmt.Errorf("Error changing access rules for share file : %s", accessErr)
			}
			d.Set("share_access_id", grant.ID)
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

		stateConf := &resource.StateChangeConf{
			Pending:    []string{"shrinking", "extending"},
			Target:     []string{"available"},
			Refresh:    waitForSFSFileActive(sfsClient, d.Id()),
			Timeout:    d.Timeout(schema.TimeoutCreate),
			Delay:      5 * time.Second,
			MinTimeout: 3 * time.Second,
		}
		_, StateErr := stateConf.WaitForState()
		if StateErr != nil {
			return fmt.Errorf("Error waiting for Share File (%s) to become ready: %s ", d.Id(), StateErr)
		}
	}

	return resourceSFSFileSystemV2Read(d, meta)
}

func resourceSFSFileSystemV2Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	sfsClient, err := config.SfsV2Client(GetRegion(d, config))
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
