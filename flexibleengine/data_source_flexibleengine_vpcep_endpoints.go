package flexibleengine

import (
	"fmt"

	"github.com/chnsz/golangsdk/openstack/vpcep/v1/endpoints"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceVPCEPEndpoints() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceVpcepEndpointsRead,

		Schema: map[string]*schema.Schema{
			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"service_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"endpoint_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"endpoints": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"service_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"service_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"service_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vpc_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"network_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ip_address": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"enable_dns": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"enable_whitelist": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"whitelist": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"packet_id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"private_domain_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tags": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"value": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"project_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"created_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"updated_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceVpcepEndpointsRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	region := GetRegion(d, config)
	vpcepClient, err := config.VPCEPClient(region)
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine VPC endpoint client: %s", err)
	}

	listOpts := endpoints.ListOpts{
		ServiceName: d.Get("service_name").(string),
		VPCID:       d.Get("vpc_id").(string),
		ID:          d.Get("endpoint_id").(string),
	}

	allEndpoints, err := endpoints.List(vpcepClient, listOpts)
	if err != nil {
		return fmt.Errorf("Unable to retrieve vpc endpoints: %s", err)
	}

	if len(allEndpoints) < 1 {
		return fmt.Errorf("Your query returned no results. " +
			"Please change your search criteria and try again.")
	}

	d.SetId(allEndpoints[0].ID)
	endpoints := make([]map[string]interface{}, len(allEndpoints))
	for i, v := range allEndpoints {
		var tag []map[string]interface{}
		for _, tagContent := range v.Tags {
			mapping := map[string]interface{}{
				"key":   tagContent.Key,
				"value": tagContent.Value,
			}
			tag = append(tag, mapping)
		}
		privateDomainName := ""
		if len(v.DNSNames) > 0 {
			privateDomainName = v.DNSNames[0]
		}
		endpoints[i] = map[string]interface{}{
			"id":                  v.ID,
			"status":              v.Status,
			"service_id":          v.ServiceID,
			"service_name":        v.ServiceName,
			"service_type":        v.ServiceType,
			"vpc_id":              v.VpcID,
			"network_id":          v.SubnetID,
			"ip_address":          v.IPAddr,
			"enable_dns":          v.EnableDNS,
			"enable_whitelist":    v.EnableWhitelist,
			"whitelist":           v.Whitelist,
			"packet_id":           v.MarkerID,
			"private_domain_name": privateDomainName,
			"tags":                tag,
			"project_id":          v.ProjectID,
			"created_at":          v.Created,
			"updated_at":          v.Updated,
		}
	}
	if err := d.Set("endpoints", endpoints); err != nil {
		return err
	}

	return nil
}
