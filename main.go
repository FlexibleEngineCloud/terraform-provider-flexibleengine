package main

import (
	"log"

	"github.com/FlexibleEngineCloud/terraform-provider-flexibleengine/flexibleengine"
	"github.com/chnsz/golangsdk/openstack/obs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
)

func main() {
	// Remove any date and time prefix in log package function output to
	// prevent duplicate timestamp and incorrect log level setting
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))

	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: flexibleengine.Provider})

	obs.CloseLog()
}
