package flexibleengine

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/huaweicloud/golangsdk"
)

func chooseCESClient(d *schema.ResourceData, config *Config) (*golangsdk.ServiceClient, error) {
	return config.loadCESClient(GetRegion(d, config))
}
