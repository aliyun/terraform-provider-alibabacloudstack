package alibabacloudstack

import (
	"fmt"
	"log"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlibabacloudStackCloudFirewallControlPolicy() *schema.Resource {
	resource := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"acl_action": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"accept", "drop", "log"}, false),
			},
			"acl_uuid": {
				Type:     schema.TypeString,
				Computed: true,
				ForceNew: true,
			},
			"application_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"ANY", "HTTP", "HTTPS", "MQTT", "Memcache", "MongoDB", "MySQL", "RDP", "Redis", "SMTP", "SMTPS", "SSH", "SSL", "VNC"}, false),
			},
			"description": {
				Type:     schema.TypeString,
				Required: true,
			},
			"dest_port": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				DiffSuppressFunc: suppressIfDestPortTypeIsNotPort,
			},
			"dest_port_group": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				DiffSuppressFunc: suppressIfDestPortTypeIsNotGroup,
			},
			"dest_port_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"group", "port"}, false),
			},
			"destination": {
				Type:     schema.TypeString,
				Required: true,
			},
			"destination_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"group", "location", "net", "domain"}, false),
			},
			"direction": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"in", "out"}, false),
			},
			"ip_version": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"lang": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"en", "zh"}, false),
			},
			"proto": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"ANY", "TCP", "UDP", "ICMP"}, false),
			},
			"release": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"source": {
				Type:     schema.TypeString,
				Required: true,
			},
			"source_ip": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"source_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"group", "location", "net"}, false),
			},
		},
	}
	setResourceFunc(resource, resourceAlibabacloudStackCloudFirewallControlPolicyCreate,
		resourceAlibabacloudStackCloudFirewallControlPolicyRead, 
		resourceAlibabacloudStackCloudFirewallControlPolicyUpdate, 
		resourceAlibabacloudStackCloudFirewallControlPolicyDelete)
	return resource
}

func suppressIfDestPortTypeIsNotPort(k, old, new string, d *schema.ResourceData) bool {
	if v, ok := d.GetOk("dest_port_type"); ok && v.(string) == "port" {
		return false
	}
	return true
}

func suppressIfDestPortTypeIsNotGroup(k, old, new string, d *schema.ResourceData) bool {
	if v, ok := d.GetOk("dest_port_type"); ok && v.(string) == "group" {
		return false
	}
	return true
}

func resourceAlibabacloudStackCloudFirewallControlPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	var response map[string]interface{}
	action := "AddControlPolicy"
	request := make(map[string]interface{})
	request["AclAction"] = d.Get("acl_action")
	request["ApplicationName"] = d.Get("application_name")
	request["Description"] = d.Get("description")
	if v, ok := d.GetOk("dest_port"); ok {
		request["DestPort"] = v
	}
	if v, ok := d.GetOk("dest_port_group"); ok {
		request["DestPortGroup"] = v
	}
	if v, ok := d.GetOk("dest_port_type"); ok {
		request["DestPortType"] = v
	}
	request["Destination"] = d.Get("destination")
	request["DestinationType"] = d.Get("destination_type")
	request["Direction"] = d.Get("direction")
	if v, ok := d.GetOk("ip_version"); ok {
		request["IpVersion"] = v
	}
	if v, ok := d.GetOk("lang"); ok {
		request["Lang"] = v
	}
	request["NewOrder"] = "-1"
	request["Proto"] = d.Get("proto")
	request["Source"] = d.Get("source")
	if v, ok := d.GetOk("source_ip"); ok {
		request["SourceIp"] = v
	}
	request["SourceType"] = d.Get("source_type")

	response, err := client.DoTeaRequest("POST", "Cloudfw", "2017-12-07", action, "", nil, nil, request)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprint(response["AclUuid"], ":", request["Direction"]))

	return nil
}

func resourceAlibabacloudStackCloudFirewallControlPolicyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	cloudfwService := CloudfwService{client}
	object, err := cloudfwService.DescribeCloudFirewallControlPolicy(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alibabacloudstack_cloud_firewall_control_policy cloudfwService.DescribeCloudFirewallControlPolicy Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return errmsgs.WrapError(err)
	}

	d.Set("acl_uuid", parts[0])
	d.Set("direction", parts[1])
	d.Set("acl_action", object["AclAction"])
	d.Set("application_name", object["ApplicationName"])
	d.Set("description", object["Description"])
	d.Set("dest_port", object["DestPort"])
	d.Set("dest_port_group", object["DestPortGroup"])
	d.Set("dest_port_type", object["DestPortType"])
	d.Set("destination", object["Destination"])
	d.Set("destination_type", object["DestinationType"])
	d.Set("direction", object["Direction"])
	d.Set("proto", object["Proto"])
	d.Set("release", object["Release"])
	d.Set("source", object["Source"])
	d.Set("source_type", object["SourceType"])
	return nil
}

func resourceAlibabacloudStackCloudFirewallControlPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	update := false
	request := map[string]interface{}{
		"AclUuid":   parts[0],
		"Direction": parts[1],
	}
	request["AclAction"] = d.Get("acl_action")
	request["ApplicationName"] = d.Get("application_name")
	request["Description"] = d.Get("description")
	request["Destination"] = d.Get("destination")
	request["DestinationType"] = d.Get("destination_type")
	request["Proto"] = d.Get("proto")
	request["Source"] = d.Get("source")
	request["SourceType"] = d.Get("source_type")
	request["DestPort"] = d.Get("dest_port")
	request["DestPortType"] = d.Get("dest_port_type")
	request["Lang"] = d.Get("lang")
	if d.HasChanges("acl_action","application_name","description","destination","destination_type","proto","source","dest_port","dest_port_type","lang") {
		update = true
	}
	if d.HasChange("dest_port_group") {
		update = true
		request["DestPortGroup"] = d.Get("dest_port_group")
	}
	if d.HasChange("release") || d.IsNewResource() {
		update = true
		request["Release"] = d.Get("release")
	}
	if update {
		if v, ok := d.GetOk("source_ip"); ok {
			request["SourceIp"] = v
		}
		action := "ModifyControlPolicy"
		_, err = client.DoTeaRequest("POST", "Cloudfw", "2017-12-07", action, "", nil, nil, request)
		if err != nil {
			return err
		}
	}
	return nil
}

func resourceAlibabacloudStackCloudFirewallControlPolicyDelete(d *schema.ResourceData, meta interface{}) (err error) {
	client := meta.(*connectivity.AlibabacloudStackClient)
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	action := "DeleteControlPolicy"
	request := map[string]interface{}{
		"AclUuid":   parts[0],
		"Direction": parts[1],
	}
	if v, ok := d.GetOk("source_ip"); ok {
		request["SourceIp"] = v
	}

	_, err = client.DoTeaRequest("POST", "Cloudfw", "2017-12-07", action, "", nil, nil, request)
	if err != nil {
		return err
	}
	return nil
}
