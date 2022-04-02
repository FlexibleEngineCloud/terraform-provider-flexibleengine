package main

import (
	"github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/flexibleengine"
	"github.com/chnsz/golangsdk/openstack/obs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: flexibleengine.Provider})

	obs.CloseLog()
}
