package flexibleengine

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/openstack/rds/v3/configurations"
)

func resourceRdsConfigurationV3() *schema.Resource {
	return &schema.Resource{
		Create: resourceRdsConfigurationV3Create,
		Read:   resourceRdsConfigurationV3Read,
		Update: resourceRdsConfigurationV3Update,
		Delete: resourceRdsConfigurationV3Delete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
			},
			"values": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: false,
			},
			"datastore": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Required: true,
						},
						"version": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"configuration_parameters": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"value": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"restart_required": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"readonly": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"value_range": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func getValues(d *schema.ResourceData) map[string]string {
	m := make(map[string]string)
	for key, val := range d.Get("values").(map[string]interface{}) {
		m[key] = val.(string)
	}
	return m
}

func getDatastore(d *schema.ResourceData) configurations.DataStore {
	datastoreRaw := d.Get("datastore").([]interface{})
	rawMap := datastoreRaw[0].(map[string]interface{})

	datastore := configurations.DataStore{
		Type:    rawMap["type"].(string),
		Version: rawMap["version"].(string),
	}

	log.Printf("[DEBUG] getDatastore: %#v", datastore)
	return datastore
}

func resourceRdsConfigurationV3Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	rdsClient, err := config.rdsV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine RDS Client: %s", err)
	}
	rdsClient.Endpoint = strings.Replace(rdsClient.Endpoint, "/rds/v1/", "/v3/", 1)

	createOpts := configurations.CreateOpts{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Values:      getValues(d),
		DataStore:   getDatastore(d),
	}
	log.Printf("[DEBUG] CreateOpts: %#v", createOpts)

	configuration, err := configurations.Create(rdsClient, createOpts).Extract()
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine RDS Configuration: %s", err)
	}

	log.Printf("[DEBUG] RDS configuration created: %#v", configuration)
	d.SetId(configuration.Id)

	return resourceRdsConfigurationV3Read(d, meta)
}

func resourceRdsConfigurationV3Read(d *schema.ResourceData, meta interface{}) error {

	config := meta.(*Config)
	rdsClient, err := config.rdsV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine RDS client: %s", err)
	}
	rdsClient.Endpoint = strings.Replace(rdsClient.Endpoint, "/rds/v1/", "/v3/", 1)
	n, err := configurations.Get(rdsClient, d.Id()).Extract()

	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error retrieving FlexibleEngine RDS Configuration: %s", err)
	}

	d.SetId(n.Id)
	d.Set("name", n.Name)
	d.Set("description", n.Description)

	datastore := []map[string]string{
		{
			"type":    n.DatastoreName,
			"version": n.DatastoreVersionName,
		},
	}
	d.Set("datastore", datastore)

	parameters := make([]map[string]interface{}, len(n.Parameters))
	for i, parameter := range n.Parameters {
		parameters[i] = make(map[string]interface{})
		parameters[i]["name"] = parameter.Name
		parameters[i]["value"] = parameter.Value
		parameters[i]["restart_required"] = parameter.RestartRequired
		parameters[i]["readonly"] = parameter.ReadOnly
		parameters[i]["value_range"] = parameter.ValueRange
		parameters[i]["type"] = parameter.Type
		parameters[i]["description"] = parameter.Description
	}
	d.Set("configuration_parameters", parameters)
	return nil
}

func resourceRdsConfigurationV3Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	rdsClient, err := config.rdsV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine RDS Client: %s", err)
	}
	rdsClient.Endpoint = strings.Replace(rdsClient.Endpoint, "/rds/v1/", "/v3/", 1)
	var updateOpts configurations.UpdateOpts

	if d.HasChange("name") {
		updateOpts.Name = d.Get("name").(string)
	}
	if d.HasChange("description") {
		updateOpts.Description = d.Get("description").(string)
	}
	if d.HasChange("values") {
		updateOpts.Values = getValues(d)
	}
	log.Printf("[DEBUG] updateOpts: %#v", updateOpts)

	err = configurations.Update(rdsClient, d.Id(), updateOpts).ExtractErr()
	if err != nil {
		return fmt.Errorf("Error updating FlexibleEngine RDS Configuration: %s", err)
	}
	return resourceRdsConfigurationV3Read(d, meta)
}

func resourceRdsConfigurationV3Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	rdsClient, err := config.rdsV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine RDS client: %s", err)
	}
	rdsClient.Endpoint = strings.Replace(rdsClient.Endpoint, "/rds/v1/", "/v3/", 1)

	err = configurations.Delete(rdsClient, d.Id()).ExtractErr()
	if err != nil {
		return fmt.Errorf("Error deleting FlexibleEngine RDS Configuration: %s", err)
	}

	d.SetId("")
	return nil
}
