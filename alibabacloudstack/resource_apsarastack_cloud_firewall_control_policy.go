package alibabacloudstack

import (
	"fmt"
	"log"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlibabacloudStackCloudFirewallControlPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlibabacloudStackCloudFirewallControlPolicyCreate,
		Read:   resourceAlibabacloudStackCloudFirewallControlPolicyRead,
		Update: resourceAlibabacloudStackCloudFirewallControlPolicyUpdate,
		Delete: resourceAlibabacloudStackCloudFirewallControlPolicyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
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
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if v, ok := d.GetOk("dest_port_type"); ok && v.(string) == "port" {
						return false
					}
					return true
				},
			},
			"dest_port_group": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if v, ok := d.GetOk("dest_port_type"); ok && v.(string) == "group" {
						return false
					}
					return true
				},
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
}

func resourceAlibabacloudStackCloudFirewallControlPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	var response map[string]interface{}
	action := "AddControlPolicy"
	request := make(map[string]interface{})
	conn, err := client.NewCloudfwClient()
	if err != nil {
		return WrapError(err)
	}
	request["AclAction"] = d.Get("acl_action")
	request["ApplicationName"] = d.Get("application_name")
	request["Description"] = d.Get("description")
	request["RegionId"] = client.RegionId
	request["Product"] = "Cloudfw"
	request["OrganizationId"] = client.Department
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
	// order属性不透出
	request["NewOrder"] = "-1"
	request["Proto"] = d.Get("proto")
	request["Source"] = d.Get("source")
	if v, ok := d.GetOk("source_ip"); ok {
		request["SourceIp"] = v
	}
	request["SourceType"] = d.Get("source_type")
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-12-07"), StringPointer("AK"), nil, request, &util.RuntimeOptions{IgnoreSSL: tea.Bool(client.Config.Insecure)})
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alibabacloudstack_cloud_firewall_control_policy", action, AlibabacloudStackSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["AclUuid"], ":", request["Direction"]))

	return resourceAlibabacloudStackCloudFirewallControlPolicyRead(d, meta)
}
func resourceAlibabacloudStackCloudFirewallControlPolicyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	cloudfwService := CloudfwService{client}
	object, err := cloudfwService.DescribeCloudFirewallControlPolicy(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alibabacloudstack_cloud_firewall_control_policy cloudfwService.DescribeCloudFirewallControlPolicy Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
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
	var response map[string]interface{}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	update := false
	request := map[string]interface{}{
		"AclUuid":   parts[0],
		"Direction": parts[1],
	}
	request["RegionId"] = client.RegionId
	request["Product"] = "Cloudfw"
	request["OrganizationId"] = client.Department
	if d.HasChange("acl_action") {
		update = true
	}
	request["AclAction"] = d.Get("acl_action")
	if d.HasChange("application_name") {
		update = true
	}
	request["ApplicationName"] = d.Get("application_name")
	if d.HasChange("description") {
		update = true
	}
	request["Description"] = d.Get("description")
	if d.HasChange("destination") {
		update = true
	}
	request["Destination"] = d.Get("destination")
	if d.HasChange("destination_type") {
		update = true
	}
	request["DestinationType"] = d.Get("destination_type")
	if d.HasChange("proto") {
		update = true
	}
	request["Proto"] = d.Get("proto")
	if d.HasChange("source") {
		update = true
	}
	request["Source"] = d.Get("source")
	if d.HasChange("source_type") {
		update = true
	}
	request["SourceType"] = d.Get("source_type")
	if d.HasChange("dest_port") {
		update = true
	}
	request["DestPort"] = d.Get("dest_port")
	if d.HasChange("dest_port_group") {
		update = true
		request["DestPortGroup"] = d.Get("dest_port_group")
	}
	if d.HasChange("dest_port_type") {
		update = true
	}
	request["DestPortType"] = d.Get("dest_port_type")
	if d.HasChange("lang") {
		update = true
		request["Lang"] = d.Get("lang")
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
		conn, err := client.NewCloudfwClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-12-07"), StringPointer("AK"), nil, request, &util.RuntimeOptions{IgnoreSSL: tea.Bool(client.Config.Insecure)})
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, request)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabacloudStackSdkGoERROR)
		}
	}
	return resourceAlibabacloudStackCloudFirewallControlPolicyRead(d, meta)
}
func resourceAlibabacloudStackCloudFirewallControlPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	action := "DeleteControlPolicy"
	var response map[string]interface{}
	conn, err := client.NewCloudfwClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"AclUuid":   parts[0],
		"Direction": parts[1],
	}

	if v, ok := d.GetOk("source_ip"); ok {
		request["SourceIp"] = v
	}
	request["RegionId"] = client.RegionId
	request["Product"] = "Cloudfw"
	request["OrganizationId"] = client.Department
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-12-07"), StringPointer("AK"), nil, request, &util.RuntimeOptions{IgnoreSSL: tea.Bool(client.Config.Insecure)})
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabacloudStackSdkGoERROR)
	}
	return nil
}
