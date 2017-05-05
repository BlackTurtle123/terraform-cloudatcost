package main

import (
	"strings"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/BlackTurtle123/go-cloudatcost/cloudatcost"
)

func resourceCloudInstance() *schema.Resource{
	return &schema.Resource{
		SchemaVersion: 1,
		Create:        resourceCloudInstanceCreate,
		Read:          resourceCloudInstanceRead,
		Update:        resourceCloudInstanceUpdate,
		Delete:        resourceCloudInstanceDelete,
		Schema: map[string]*schema.Schema{ // List of supported configuration fields for your resource
			"storage": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"datacenter": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
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
			"status": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceCloudInstanceCreate(d *schema.ResourceData, meta interface{}) error {
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
	}else{
		d.Set("runmode",strings.ToLower(d.Get("runmode").(string)))
	}

	if d.Get("label") != nil && d.Get("label").(string) != "" {
		_,_, errr := client.ServersService.Rename(server.Sid, d.Get("label").(string))
		if errr != nil {
			return errr
		}
	}

	d.Set("ip",server.IP)
	d.Set("password",server.Rootpass)
	d.Set("status",server.Status)
	//d.SetId(strconv.Itoa(server.Sid))
	return nil
}

func resourceCloudInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*cloudatcost.Client)
	var server cloudatcost.ListServer
	listservers,_, _ := client.ServersService.List()

	for i := 0; i< len(listservers); i++{
		if listservers[i].Sid == d.Id(){
			server = listservers[i]
			break
		}
	}

	d.Set("storage",server.Storage)
	d.Set("os",server.Packageid)
	d.Set("cpu",server.CPU)
	d.Set("ram",server.RAM)
	d.Set("runmode",strings.ToLower(server.Mode))
	d.Set("label",server.Lable)
	d.Set("ip",server.IP)
	d.Set("password",server.Rootpass)
	d.Set("status",server.Status)

	return nil
}

func resourceCloudInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*cloudatcost.Client)
	d.Partial(true)

	if d.HasChange("cpu") == true  || d.HasChange("ram") == true || d.HasChange("storage") == true || d.HasChange("os") == true{
		resourceCloudInstanceDelete(d,meta)
		resourceCloudInstanceCreate(d,meta)
	}

	if d.HasChange("runmode") {
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
	return resourceCloudInstanceRead(d,meta)
}

func resourceCloudInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*cloudatcost.Client)

	_, _, err := client.CloudProService.Delete(d.Id())

	if err != nil {
		return err
	}

	return nil
}