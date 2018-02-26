package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/huaweicloud/terraform-provider-flexibleengine/flexibleengine" // TODO: Revert path when merge
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: flexibleengine.Provider})
}
