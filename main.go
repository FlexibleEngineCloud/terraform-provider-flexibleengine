package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/khdegraaf/terraform-provider-orangecloud/orangecloud" // TODO: Revert path when merge
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: orangecloud.Provider})
}
