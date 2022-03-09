package flexibleengine

import (
	"fmt"
	"strconv"

	"github.com/chnsz/golangsdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceRdsFlavorV3() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceRdsFlavorV3Read,

		Schema: map[string]*schema.Schema{
			"db_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"db_version": {
				Type:     schema.TypeString,
				Required: true,
			},
			"instance_mode": {
				Type:     schema.TypeString,
				Required: true,
			},
			"vcpus": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"memory": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"flavors": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"vcpus": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"memory": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"mode": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceRdsFlavorV3Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client, err := config.RdsV3Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine rds client: %s", err)
	}

	link := fmt.Sprintf("flavors/%s?version_name=%s", d.Get("db_type").(string), d.Get("db_version").(string))
	url := client.ServiceURL(link)

	r, err := sendRdsFlavorV3ListRequest(client, url)
	if err != nil {
		return err
	}

	mode := d.Get("instance_mode").(string)
	cpu := d.Get("vcpus").(int)
	mem := d.Get("memory").(int)

	flavors := make([]interface{}, 0, len(r.([]interface{})))
	for _, item := range r.([]interface{}) {
		val := item.(map[string]interface{})

		vcpu, _ := strconv.Atoi(val["vcpus"].(string))
		if cpu > 0 && vcpu != cpu {
			continue
		}

		if mem > 0 && int(val["ram"].(float64)) != mem {
			continue
		}

		if val["instance_mode"].(string) != mode {
			continue
		}

		flavors = append(flavors, map[string]interface{}{
			"vcpus":  val["vcpus"],
			"memory": val["ram"],
			"name":   val["spec_code"],
			"mode":   val["instance_mode"],
		})
	}

	d.SetId("flavors")
	return d.Set("flavors", flavors)
}

func sendRdsFlavorV3ListRequest(client *golangsdk.ServiceClient, url string) (interface{}, error) {
	r := golangsdk.Result{}
	_, r.Err = client.Get(url, &r.Body, &golangsdk.RequestOpts{
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
			"X-Language":   "en-us",
		}})
	if r.Err != nil {
		return nil, fmt.Errorf("Error fetching flavors for rds v3, error: %s", r.Err)
	}

	v, err := navigateValue(r.Body, []string{"flavors"}, nil)
	if err != nil {
		return nil, err
	}
	return v, nil
}
