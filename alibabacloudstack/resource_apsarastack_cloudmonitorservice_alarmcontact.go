package alibabacloudstack

import (
	"fmt"
	"log"
	"reflect"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/cms"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
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
				Type:         schema.TypeString,
				Optional:     true,
				Deprecated:   "Field 'channels_aliim' is deprecated and will be removed in a future release. Please use 'channels_ali_im' instead.",
			},
			"channels_ali_im": {
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

	request := client.NewCommonRequest("POST", "Cms", "2019-01-01", "PutContact", "")
	request.Headers["Content-Type"] = "application/json; charset=UTF-8"
	request.QueryParams["ContactName"] = d.Get("alarm_contact_name").(string)
	request.QueryParams["Describe"] = d.Get("describe").(string)
	if v, err := connectivity.GetResourceData(d, reflect.TypeOf(""), "channels_ali_im", "channels_aliim"); err == nil {
		request.QueryParams["Channels.AliIM"] = v.(string)
	} else {
		return err
	}
	if v, ok := d.GetOk("channels_ding_web_hook"); ok {
		request.QueryParams["Channels.DingWebHook"] = v.(string)
	}
	if v, ok := d.GetOk("channels_mail"); ok {
		request.QueryParams["ChannelsMail"] = v.(string)
	}
	if v, ok := d.GetOk("channels_sms"); ok {
		request.QueryParams["Channels.SMS"] = v.(string)
	}
	if v, ok := d.GetOk("lang"); ok {
		request.QueryParams["Lang"] = v.(string)
	}

	raw, err := client.WithCmsClient(func(cmsClient *cms.Client) (interface{}, error) {
		return cmsClient.ProcessCommonRequest(request)
	})
	bresponse, ok := raw.(*cms.PutContactResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_cms_alarm_contact", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug(request.GetActionName(), raw)

	if bresponse.Code != "200" {
		return errmsgs.WrapError(errmsgs.Error("PutContact failed for " + bresponse.Message))
	}
	d.SetId(fmt.Sprintf("%v", d.Get("alarm_contact_name").(string)))

	return resourceAlibabacloudstackCmsAlarmContactRead(d, meta)
}

func resourceAlibabacloudstackCmsAlarmContactRead(d *schema.ResourceData, meta interface{}) error {
	waitSecondsIfWithTest(1)
	client := meta.(*connectivity.AlibabacloudStackClient)
	cmsService := CmsService{client}
	object, err := cmsService.DescribeCmsAlarmContact(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alibabacloudstack_cloud_monitor_service_alarm_contact cmsService.DescribeCmsAlarmContact Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}

	d.Set("alarm_contact_name", d.Id())
	connectivity.SetResourceData(d, object.Channels.AliIM, "channels_ali_im", "channels_aliim")
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
	client.InitRpcRequest(*request.RpcRequest)
	request.ContactName = d.Id()
	if d.HasChange("channels_ali_im") || d.HasChange("channels_aliim") {
		update = true
	}
	if v, err := connectivity.GetResourceData(d, reflect.TypeOf(""), "channels_ali_im", "channels_aliim"); err == nil {
		request.ChannelsAliIM = v.(string)
	} else {
		return err
	}
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
		bresponse, ok := raw.(*cms.PutContactResponse)
		addDebug(request.GetActionName(), raw)

		if bresponse.Code != "200" {
			return errmsgs.WrapError(errmsgs.Error("PutContact failed for " + bresponse.Message))
		}
		if err != nil {
			errmsg := ""
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
	}
	return resourceAlibabacloudstackCmsAlarmContactRead(d, meta)
}

func resourceAlibabacloudstackCmsAlarmContactDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	request := cms.CreateDeleteContactRequest()
	client.InitRpcRequest(*request.RpcRequest)
	request.ContactName = d.Id()

	raw, err := client.WithCmsClient(func(cmsClient *cms.Client) (interface{}, error) {
		return cmsClient.DeleteContact(request)
	})
	bresponse, ok := raw.(*cms.DeleteContactResponse)
	addDebug(request.GetActionName(), raw)

	if bresponse.Code != "200" {
		return errmsgs.WrapError(errmsgs.Error("DeleteContact failed for " + bresponse.Message))
	}
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		if errmsgs.IsExpectedErrors(err, []string{"400", "403", "404", "ContactNotExists"}) {
			return nil
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	return nil
}
