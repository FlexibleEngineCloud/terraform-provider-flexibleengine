package main

import (
	"github.com/Karajan-project/terraform-provider-flexibleengine/flexibleengine" // TODO: Revert path when merge
	"github.com/hashicorp/terraform/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: flexibleengine.Provider})
}
