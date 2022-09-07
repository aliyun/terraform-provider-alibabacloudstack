package alibabacloudstack

import (
	"fmt"
	"log"
	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/cms"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlibabacloudstackCmsAlarmContact() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlibabacloudstackCmsAlarmContactCreate,
		Read:   resourceAlibabacloudstackCmsAlarmContactRead,
		Update: resourceAlibabacloudstackCmsAlarmContactUpdate,
		Delete: resourceAlibabacloudstackCmsAlarmContactDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"alarm_contact_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"channels_aliim": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"channels_ding_web_hook": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"channels_mail": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"channels_sms": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"describe": {
				Type:     schema.TypeString,
				Required: true,
			},
			"lang": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"en", "zh-cn"}, false),
			},
		},
	}
}

func resourceAlibabacloudstackCmsAlarmContactCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	request := cms.CreatePutContactRequest()
	if strings.ToLower(client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.RegionId = client.RegionId
	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "Cms", "Department": client.Department, "ResourceGroup": client.ResourceGroup}
	request.ContactName = d.Get("alarm_contact_name").(string)
	if v, ok := d.GetOk("channels_aliim"); ok {
		request.ChannelsAliIM = v.(string)
	}

	if v, ok := d.GetOk("channels_ding_web_hook"); ok {
		request.ChannelsDingWebHook = v.(string)
	}

	if v, ok := d.GetOk("channels_mail"); ok {
		request.ChannelsMail = v.(string)
	}

	if v, ok := d.GetOk("channels_sms"); ok {
		request.ChannelsSMS = v.(string)
	}

	request.Describe = d.Get("describe").(string)
	if v, ok := d.GetOk("lang"); ok {
		request.Lang = v.(string)
	}

	raw, err := client.WithCmsClient(func(cmsClient *cms.Client) (interface{}, error) {
		return cmsClient.PutContact(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alibabacloudstack_cms_alarm_contact", request.GetActionName(), AlibabacloudStackSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw)
	response, _ := raw.(*cms.PutContactResponse)

	if response.Code != "200" {
		return WrapError(Error("PutContact failed for " + response.Message))
	}
	d.SetId(fmt.Sprintf("%v", request.ContactName))

	return resourceAlibabacloudstackCmsAlarmContactRead(d, meta)
}
func resourceAlibabacloudstackCmsAlarmContactRead(d *schema.ResourceData, meta interface{}) error {
	wiatSecondsIfWithTest(1)
	client := meta.(*connectivity.AlibabacloudStackClient)
	cmsService := CmsService{client}
	object, err := cmsService.DescribeCmsAlarmContact(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alibabacloudstack_cloud_monitor_service_alarm_contact cmsService.DescribeCmsAlarmContact Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("alarm_contact_name", d.Id())
	d.Set("channels_aliim", object.Channels.AliIM)
	d.Set("channels_ding_web_hook", object.Channels.DingWebHook)
	d.Set("channels_mail", object.Channels.Mail)
	d.Set("channels_sms", object.Channels.SMS)
	d.Set("describe", object.Desc)
	d.Set("lang", object.Lang)
	return nil
}

func resourceAlibabacloudstackCmsAlarmContactUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	update := false
	request := cms.CreatePutContactRequest()
	request.ContactName = d.Id()
	if d.HasChange("channels_aliim") {
		update = true
	}
	request.ChannelsAliIM = d.Get("channels_aliim").(string)
	if d.HasChange("channels_ding_web_hook") {
		update = true
	}
	request.ChannelsDingWebHook = d.Get("channels_ding_web_hook").(string)
	if d.HasChange("channels_mail") {
		update = true
	}
	request.ChannelsMail = d.Get("channels_mail").(string)
	if d.HasChange("channels_sms") {
		update = true
	}
	request.ChannelsSMS = d.Get("channels_sms").(string)
	if d.HasChange("describe") {
		update = true
	}
	request.Describe = d.Get("describe").(string)
	if d.HasChange("lang") {
		update = true
		request.Lang = d.Get("lang").(string)
	}
	if update {
		raw, err := client.WithCmsClient(func(cmsClient *cms.Client) (interface{}, error) {
			return cmsClient.PutContact(request)
		})
		addDebug(request.GetActionName(), raw)
		response, _ := raw.(*cms.PutContactResponse)

		if response.Code != "200" {
			return WrapError(Error("PutContact failed for " + response.Message))
		}
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabacloudStackSdkGoERROR)
		}
	}
	return resourceAlibabacloudstackCmsAlarmContactRead(d, meta)
}
func resourceAlibabacloudstackCmsAlarmContactDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	request := cms.CreateDeleteContactRequest()
	request.ContactName = d.Id()
	raw, err := client.WithCmsClient(func(cmsClient *cms.Client) (interface{}, error) {
		return cmsClient.DeleteContact(request)
	})
	addDebug(request.GetActionName(), raw)
	response, _ := raw.(*cms.DeleteContactResponse)

	if response.Code != "200" {
		return WrapError(Error("DeleteContact failed for " + response.Message))
	}
	if err != nil {
		if IsExpectedErrors(err, []string{"400", "403", "404", "ContactNotExists"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabacloudStackSdkGoERROR)
	}
	return nil
}
