
package main

import (
	//"strconv"
	"strings"
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
	return &schema.Provider{ // Source https://github.com/hashicorp/terraform/blob/v0.6.6/helper/schema/provider.go#L20-L43
		Schema:        providerSchema(),
		ResourcesMap:  providerResources(),
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

func providerResources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"cloudatcost_machine": &schema.Resource{
			SchemaVersion: 1,
			Create:        createFunc,
			Read:          readFunc,
			Update:        updateFunc,
			Delete:        deleteFunc,
			Schema: map[string]*schema.Schema{ // List of supported configuration fields for your resource
				"storage": &schema.Schema{
					Type:     schema.TypeString,
					Required: true,
				},
				"datacenter": &schema.Schema{
					Type:     schema.TypeString,
					Required: true,
				},
				"os": &schema.Schema{
					Type:     schema.TypeString,
					Required: true,
				},
				"cpu": &schema.Schema{
					Type:     schema.TypeString,
					Required: true,
				},
				"ram": &schema.Schema{
					Type:     schema.TypeString,
					Required: true,
				},
				"runmode": &schema.Schema{
					Type:     schema.TypeString,
					Optional: true,
				},
				"label": &schema.Schema{
					Type:     schema.TypeString,
					Optional: true,
				},
				"ip": &schema.Schema{
					Type:     schema.TypeString,
					Computed: true,
				},
				"password": &schema.Schema{
					Type:     schema.TypeString,
					Computed: true,
				},
			},
		},
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {


	return  cloudatcost.NewClient(&cloudatcost.Option{Login: d.Get("login").(string), Key: d.Get("api_key").(string)})
}

func createFunc(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*cloudatcost.Client)
	_, _, err := client.CloudProService.Create(&cloudatcost.CreateServerOptions{
		Cpu: d.Get("cpu").(string),
		Ram:  d.Get("ram").(string),
		Storage:  d.Get("storage").(string),
		OS:  d.Get("os").(string),
		Datacenter: d.Get("datacenter").(string) },
	)

	if err != nil {
		return err
	}

	listservers,_, err := client.ServersService.List()

	serverLength := len(listservers)
	//need fix both servers are created at the same time
	//impossible to know which server is which one
	server := listservers[serverLength-1]
	d.SetId(server.Sid)
	if err != nil {
		return err
	}
	_,_,error := client.RunModeService.Mode(server.Sid,strings.ToLower(d.Get("runmode").(string)))
	if error != nil {
		return error
	}

	if d.Get("label") != nil && d.Get("label").(string) != "" {
		_,_, errr := client.ServersService.Rename(server.Sid, d.Get("label").(string))
		if errr != nil {
			return errr
		}
	}
	d.Set("ip",server.IP)
	d.Set("password",server.Rootpass)
	//d.SetId(strconv.Itoa(server.Sid))
	return nil
}

func readFunc(d *schema.ResourceData, meta interface{}) error {

	return nil
}

func updateFunc(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*cloudatcost.Client)
	d.Partial(true)

	if d.HasChange("cpu") == true  || d.HasChange("ram") == true || d.HasChange("storage") == true || d.HasChange("os") == true{
		deleteFunc(d,meta)
		createFunc(d,meta)
	}

	if d.HasChange("runmode"){
		d.SetPartial("runmode")
		_,_,err := client.RunModeService.Mode(d.Id(),strings.ToLower(d.Get("runmode").(string)))
		if err != nil {
			return err
		}
	}
	if d.HasChange("label") && d.Get("label").(string) != "" {
		d.SetPartial("label")
		_, _, errr := client.ServersService.Rename(d.Id(), d.Get("label").(string))
		if errr != nil {
			return errr
		}
	}


	d.Partial(false)
	return readFunc(d,meta)
}

func deleteFunc(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*cloudatcost.Client)

	_, _, err := client.CloudProService.Delete(d.Id())

	if err != nil {
		return err
	}

	return nil
}
