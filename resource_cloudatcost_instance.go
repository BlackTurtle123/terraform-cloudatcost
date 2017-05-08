package main

import (
	"bytes"
	"strings"
	"github.com/hashicorp/terraform/helper/schema"
	"strconv"
	"github.com/BlackTurtle123/go-cloudatcost/cloudatcost"
)

func resourceCloudInstance() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,
		Create:        resourceCloudInstanceCreate,
		Read:          resourceCloudInstanceRead,
		Update:        resourceCloudInstanceUpdate,
		Delete:        resourceCloudInstanceDelete,
		Schema: map[string]*schema.Schema{// List of supported configuration fields for your resource
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
	res, _, err := client.CloudProService.Resources()
	uCPU, _ := strconv.Atoi(res.Used.CPU)
	tCPU, _ := strconv.Atoi(res.Total.CPU)
	uStorage, _ := strconv.Atoi(res.Used.Storage)
	tStorage, _ := strconv.Atoi(res.Total.Storage)
	uRam, _ := strconv.Atoi(res.Used.Ram)
	tRam, _ := strconv.Atoi(res.Total.Ram)
	if tCPU != 0 {

		remainingCpu := tCPU - uCPU
		remainingRam := tRam - uRam
		remainingStorage := tStorage - uStorage
		ram, _ := strconv.Atoi(d.Get("ram").(string))
		storage, _ := strconv.Atoi(d.Get("storage").(string))
		cpu, _ := strconv.Atoi(d.Get("cpu").(string))
		if remainingRam < ram || remainingCpu < cpu || remainingStorage < storage {
			return &notEnoughResources{d.Get("cpu").(string),
				d.Get("ram").(string),
				d.Get("storage").(string),
				strconv.Itoa(remainingCpu),
				strconv.Itoa(remainingRam),
				strconv.Itoa(remainingStorage),
			}
		}
	} else {
		return &notEnoughResources{d.Get("cpu").(string),
			d.Get("ram").(string),
			d.Get("storage").(string),
			"0",
			"0",
			"0",
		}
	}
	if err != nil {
		return err
	}
	imageID, err := resourceCloudMapImageToId(d, meta)

	if err != nil {
		return err
	}
	_, _, err = client.CloudProService.Create(&cloudatcost.CreateServerOptions{
		Cpu: d.Get("cpu").(string),
		Ram:  d.Get("ram").(string),
		Storage:  d.Get("storage").(string),
		OS:  imageID,
		Datacenter: d.Get("datacenter").(string) },
	)

	if err != nil {
		return err
	}

	listservers, _, err := client.ServersService.List()

	serverLength := len(listservers)
	//need fix both servers are created at the same time
	//impossible to know which server is which one
	server := listservers[serverLength - 1]
	d.SetId(server.Sid)
	if err != nil {
		return err
	}
	_, _, err = client.RunModeService.Mode(server.Sid, strings.ToLower(d.Get("runmode").(string)))
	if err != nil {
		return err
	} else {
		d.Set("runmode", strings.ToLower(d.Get("runmode").(string)))
	}

	if d.Get("label") != nil && d.Get("label").(string) != "" {
		_, _, err = client.ServersService.Rename(server.Sid, d.Get("label").(string))
		if err != nil {
			return err
		}
	}

	d.Set("ip", server.IP)
	d.Set("password", server.Rootpass)
	d.Set("status", server.Status)
	return nil
}

func resourceCloudInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*cloudatcost.Client)
	var server cloudatcost.ListServer
	listservers, _, _ := client.ServersService.List()

	for i := 0; i < len(listservers); i++ {
		if listservers[i].Sid == d.Id() {
			server = listservers[i]
			break
		}
	}

	d.Set("storage", server.Storage)
	d.Set("os", server.Template)
	d.Set("cpu", server.CPU)
	d.Set("ram", server.RAM)
	d.Set("runmode", strings.ToLower(server.Mode))
	d.Set("label", server.Lable)
	d.Set("ip", server.IP)
	d.Set("password", server.Rootpass)
	d.Set("status", server.Status)

	return nil
}

func resourceCloudInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*cloudatcost.Client)
	d.Partial(true)

	if d.HasChange("cpu") == true || d.HasChange("ram") == true || d.HasChange("storage") == true || d.HasChange("os") == true {
		resourceCloudInstanceDelete(d, meta)
		resourceCloudInstanceCreate(d, meta)
	}

	if d.HasChange("runmode") {
		d.SetPartial("runmode")
		_, _, err := client.RunModeService.Mode(d.Id(), strings.ToLower(d.Get("runmode").(string)))
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
	return resourceCloudInstanceRead(d, meta)
}

func resourceCloudInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*cloudatcost.Client)

	_, _, err := client.CloudProService.Delete(d.Id())

	if err != nil {
		return err
	}

	return nil
}

func resourceCloudMapImageToId(d *schema.ResourceData, meta interface{}) (string, error) {
	client := meta.(*cloudatcost.Client)
	var s []cloudatcost.ListTemplate
	s, _, _ = client.ListTemplatesService.ListTemplates()
	osImage := d.Get("os").(string)
	for i := 0; i < len(s); i++ {
		if s[i].Name == osImage {
			return s[i].Ce_id, nil
			break
		}
	}
	var buffer bytes.Buffer
	buffer.WriteString("'")
	for i := 0; i < len(s); i++ {

		buffer.WriteString(s[i].Name)
		buffer.WriteString("',")
	}

	return "", &osImageError{osImage, buffer.String()}
}
