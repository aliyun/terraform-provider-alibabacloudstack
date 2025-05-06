package alibabacloudstack

import (
	"fmt"
	"log"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cms"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudstackCmsAlarmContactGroup() *schema.Resource {
	resource := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"alarm_contact_group_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"contacts": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"describe": {
				Type:     schema.TypeString,
				Required: true,
			},
			"enable_subscribed": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
		},
	}
	setResourceFunc(resource, resourceAlibabacloudstackCmsAlarmContactGroupCreate, resourceAlibabacloudstackCmsAlarmContactGroupRead, resourceAlibabacloudstackCmsAlarmContactGroupUpdate, resourceAlibabacloudstackCmsAlarmContactGroupDelete)
	return resource
}

func resourceAlibabacloudstackCmsAlarmContactGroupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	request := cms.CreatePutContactGroupRequest()
	client.InitRpcRequest(*request.RpcRequest)
	request.ContactGroupName = d.Get("alarm_contact_group_name").(string)
	if v, ok := d.GetOk("contacts"); ok {
		contactNames := expandStringList(v.(*schema.Set).List())
		request.ContactNames = &contactNames
	}

	if v, ok := d.GetOk("describe"); ok {
		request.Describe = v.(string)
	}

	if v, ok := d.GetOkExists("enable_subscribed"); ok {
		request.EnableSubscribed = requests.NewBoolean(v.(bool))
	}

	raw, err := client.WithCmsClient(func(cmsClient *cms.Client) (interface{}, error) {
		return cmsClient.PutContactGroup(request)
	})
	addDebug(request.GetActionName(), raw)
	response, ok := raw.(*cms.PutContactGroupResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_cms_alarm_contact_group", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	if response.Code != "200" {
		return errmsgs.WrapError(errmsgs.Error("PutContactGroup failed for " + response.Message))
	}
	d.SetId(fmt.Sprintf("%v", request.ContactGroupName))

	return nil
}

func resourceAlibabacloudstackCmsAlarmContactGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	cmsService := CmsService{client}
	object, err := cmsService.DescribeCmsAlarmContactGroup(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alibabacloudstack_cloud_monitor_service_alarm_contact_group cmsService.DescribeCmsAlarmContactGroup Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}

	d.Set("alarm_contact_group_name", d.Id())
	d.Set("contacts", object.Contacts.Contact)
	d.Set("describe", object.Describe)
	d.Set("enable_subscribed", object.EnableSubscribed)
	return nil
}

func resourceAlibabacloudstackCmsAlarmContactGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	update := false
	request := cms.CreatePutContactGroupRequest()
	client.InitRpcRequest(*request.RpcRequest)
	request.ContactGroupName = d.Id()
	if d.HasChange("contacts") {
		update = true
		contactNames := expandStringList(d.Get("contacts").(*schema.Set).List())
		request.ContactNames = &contactNames
	}
	if d.HasChange("describe") {
		update = true
		request.Describe = d.Get("describe").(string)
	}
	if d.HasChange("enable_subscribed") {
		update = true
		request.EnableSubscribed = requests.NewBoolean(d.Get("enable_subscribed").(bool))
	}
	if update {
		raw, err := client.WithCmsClient(func(cmsClient *cms.Client) (interface{}, error) {
			return cmsClient.PutContactGroup(request)
		})
		addDebug(request.GetActionName(), raw)
		response, ok := raw.(*cms.PutContactGroupResponse)
		if err != nil {
			errmsg := ""
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}

		if response.Code != "200" {
			return errmsgs.WrapError(errmsgs.Error("PutContactGroup failed for " + response.Message))
		}
	}
	return nil
}

func resourceAlibabacloudstackCmsAlarmContactGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	request := cms.CreateDeleteContactGroupRequest()
	client.InitRpcRequest(*request.RpcRequest)
	request.ContactGroupName = d.Id()

	raw, err := client.WithCmsClient(func(cmsClient *cms.Client) (interface{}, error) {
		return cmsClient.DeleteContactGroup(request)
	})
	addDebug(request.GetActionName(), raw)
	response, ok := raw.(*cms.DeleteContactGroupResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
		}
		if errmsgs.IsExpectedErrors(err, []string{"400", "403", "404", "ContactNotExists"}) {
			return nil
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	if response.Code != "200" {
		return errmsgs.WrapError(errmsgs.Error("DeleteContactGroup failed for " + response.Message))
	}
	return nil
}