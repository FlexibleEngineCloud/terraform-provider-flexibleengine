package flexibleengine

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/cce/v3/addons"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceCCEAddon() *schema.Resource {
	return &schema.Resource{
		Create: resourceCCEAddonCreate,
		Read:   resourceCCEAddonRead,
		Delete: resourceCCEAddonDelete,

		Importer: &schema.ResourceImporter{
			State: resourceCCEAddonImport,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(3 * time.Minute),
		},

		Schema: map[string]*schema.Schema{ // request and response parameters
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"cluster_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"version": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"template_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"values": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"basic": {
							Type:         schema.TypeString,
							Required:     true,
							ForceNew:     true,
							ValidateFunc: validation.StringIsJSON,
						},
						"custom": {
							Type:         schema.TypeString,
							Optional:     true,
							ForceNew:     true,
							ValidateFunc: validation.StringIsJSON,
						},
						"flavor": {
							Type:         schema.TypeString,
							Optional:     true,
							ForceNew:     true,
							ValidateFunc: validation.StringIsJSON,
						},
					},
				},
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func getValuesValues(d *schema.ResourceData) (basic, custom, flavor map[string]interface{}, err error) {
	values := d.Get("values").([]interface{})
	if len(values) == 0 {
		basic = map[string]interface{}{}
		return
	}

	valuesMap := values[0].(map[string]interface{})

	if basicRaw := valuesMap["basic"].(string); basicRaw != "" {
		err = json.Unmarshal([]byte(basicRaw), &basic)
		if err != nil {
			err = fmt.Errorf("Error unmarshalling basic json: %s", err)
			return
		}
	}

	if customRaw := valuesMap["custom"].(string); customRaw != "" {
		err = json.Unmarshal([]byte(customRaw), &custom)
		if err != nil {
			err = fmt.Errorf("Error unmarshalling custom json: %s", err)
			return
		}
	}

	if flavorRaw := valuesMap["flavor"].(string); flavorRaw != "" {
		err = json.Unmarshal([]byte(flavorRaw), &flavor)
		if err != nil {
			err = fmt.Errorf("Error unmarshalling flavor json %s", err)
			return
		}
	}

	return
}

func resourceCCEAddonCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	cceClient, err := config.CceAddonV3Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Unable to create FlexibleEngine CCE client : %s", err)
	}

	var clusterID = d.Get("cluster_id").(string)

	basic, custom, flavor, err := getValuesValues(d)
	if err != nil {
		return fmt.Errorf("error getting values for CCE addon: %s", err)
	}

	createOpts := addons.CreateOpts{
		Kind:       "Addon",
		ApiVersion: "v3",
		Metadata: addons.CreateMetadata{
			Anno: addons.Annotations{
				AddonInstallType: "install",
			},
		},
		Spec: addons.RequestSpec{
			Version:           d.Get("version").(string),
			ClusterID:         clusterID,
			AddonTemplateName: d.Get("template_name").(string),
			Values: addons.Values{
				Basic:  basic,
				Custom: custom,
				Flavor: flavor,
			},
		},
	}

	create, err := addons.Create(cceClient, createOpts, clusterID).Extract()
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine CCEAddon: %s", err)
	}

	d.SetId(create.Metadata.Id)

	log.Printf("[DEBUG] Waiting for FlexibleEngine CCEAddon (%s) to become available", create.Metadata.Id)
	stateConf := &resource.StateChangeConf{
		Pending:                   []string{"installing"},
		Target:                    []string{"running", "available", "abnormal"},
		Refresh:                   waitForCCEAddonActive(cceClient, create.Metadata.Id, clusterID),
		Timeout:                   d.Timeout(schema.TimeoutCreate),
		Delay:                     10 * time.Second,
		PollInterval:              10 * time.Second,
		ContinuousTargetOccurence: 3,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Error installing FlexibleEngine CCEAddon: %s", err)
	}

	return resourceCCEAddonRead(d, meta)
}

func resourceCCEAddonRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	cceClient, err := config.CceAddonV3Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine CCE client: %s", err)
	}

	var clusterID = d.Get("cluster_id").(string)
	n, err := addons.Get(cceClient, d.Id(), clusterID).Extract()
	if err != nil {
		return CheckDeleted(d, err, "CCE Addon")
	}

	d.Set("cluster_id", n.Spec.ClusterID)
	d.Set("version", n.Spec.Version)
	d.Set("template_name", n.Spec.AddonTemplateName)
	d.Set("status", n.Status.Status)
	d.Set("description", n.Spec.Description)

	return nil
}

func resourceCCEAddonDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	cceClient, err := config.CceAddonV3Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine CCEAddon Client: %s", err)
	}

	var clusterID = d.Get("cluster_id").(string)
	err = addons.Delete(cceClient, d.Id(), clusterID).ExtractErr()
	if err != nil {
		return fmt.Errorf("Error deleting FlexibleEngine CCE Addon: %s", err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"Deleting", "Available", "Unavailable"},
		Target:       []string{"Deleted"},
		Refresh:      waitForCCEAddonDelete(cceClient, d.Id(), clusterID),
		Timeout:      d.Timeout(schema.TimeoutDelete),
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}

	_, err = stateConf.WaitForState()

	if err != nil {
		return fmt.Errorf("Error deleting FlexibleEngine CCE Addon: %s", err)
	}

	d.SetId("")
	return nil
}

func waitForCCEAddonActive(cceAddonV3Client *golangsdk.ServiceClient, id, clusterID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		n, err := addons.Get(cceAddonV3Client, id, clusterID).Extract()
		if err != nil {
			return nil, "", err
		}

		return n, n.Status.Status, nil
	}
}

func waitForCCEAddonDelete(cceClient *golangsdk.ServiceClient, id, clusterID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		log.Printf("[DEBUG] Attempting to delete FlexibleEngine CCE Addon %s", id)

		r, err := addons.Get(cceClient, id, clusterID).Extract()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				log.Printf("[DEBUG] Successfully deleted FlexibleEngine CCE Addon %s", id)
				return r, "Deleted", nil
			}
		}
		if r.Status.Status == "Deleting" {
			return r, "Deleting", nil
		}
		log.Printf("[DEBUG] FlexibleEngine CCE Addon %s still available", id)
		return r, "Available", nil
	}
}

func resourceCCEAddonImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	parts := strings.SplitN(d.Id(), "/", 2)
	if len(parts) != 2 {
		err := fmt.Errorf("Invalid format specified for CCE Addon. Format must be <cluster id>/<addon id>")
		return nil, err
	}

	clusterID := parts[0]
	addonID := parts[1]

	d.SetId(addonID)
	d.Set("cluster_id", clusterID)

	return []*schema.ResourceData{d}, nil
}
