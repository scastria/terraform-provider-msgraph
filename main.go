package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/scastria/terraform-provider-msgraph/msgraph"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: msgraph.Provider,
	})
}
