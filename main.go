package main

import (
	"github.com/Karajan-project/terraform-provider-orangecloud/orangecloud" // TODO: Revert path when merge
	"github.com/hashicorp/terraform/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: orangecloud.Provider})
}
