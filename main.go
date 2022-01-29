package main

import (
	"github.com/chnsz/golangsdk/openstack/obs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/terraform-providers/terraform-provider-flexibleengine/flexibleengine"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: flexibleengine.Provider})

	obs.CloseLog()
}
