package main

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/plugin"
	"github.com/hashicorp/terraform/terraform"
	"github.com/BlackTurtle123/go-cloudatcost/cloudatcost"
)

func main() {
	opts := plugin.ServeOpts{
		ProviderFunc: Provider,
	}
	plugin.Serve(&opts)
}

func Provider() terraform.ResourceProvider {
	return &schema.Provider{// Source https://github.com/hashicorp/terraform/blob/v0.6.6/helper/schema/provider.go#L20-L43
		Schema:        providerSchema(),
		ResourcesMap:  map[string]*schema.Resource{
			"cloudatcost_instance":                resourceCloudInstance(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"api_key": &schema.Schema{
			Type:        schema.TypeString,
			Required:    true,
			Description: "API Key used to authenticate with the service provider",
		},
		"login": &schema.Schema{
			Type:        schema.TypeString,
			Required:    true,
			Description: "login used to authenticate with the service provider",
		},
	}
}
func providerConfigure(d *schema.ResourceData) (interface{}, error) {

	return cloudatcost.NewClient(&cloudatcost.Option{Login: d.Get("login").(string), Key: d.Get("api_key").(string)})
}