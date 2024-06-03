package alibabacloudstack

import (
	"fmt"
	"log"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cms"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudstackCmsAlarmContactGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlibabacloudstackCmsAlarmContactGroupCreate,
		Read:   resourceAlibabacloudstackCmsAlarmContactGroupRead,
		Update: resourceAlibabacloudstackCmsAlarmContactGroupUpdate,
		Delete: resourceAlibabacloudstackCmsAlarmContactGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
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
}

func resourceAlibabacloudstackCmsAlarmContactGroupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	request := cms.CreatePutContactGroupRequest()
	request.Headers["x-ascm-product-name"] = "Cms"
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
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alibabacloudstack_cms_alarm_contact_group", request.GetActionName(), AlibabacloudStackSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw)
	response, _ := raw.(*cms.PutContactGroupResponse)

	if response.Code != "200" {
		return WrapError(Error("PutContactGroup failed for " + response.Message))
	}
	d.SetId(fmt.Sprintf("%v", request.ContactGroupName))

	return resourceAlibabacloudstackCmsAlarmContactGroupRead(d, meta)
}
func resourceAlibabacloudstackCmsAlarmContactGroupRead(d *schema.ResourceData, meta interface{}) error {
	waitSecondsIfWithTest(1)
	client := meta.(*connectivity.AlibabacloudStackClient)
	cmsService := CmsService{client}
	object, err := cmsService.DescribeCmsAlarmContactGroup(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alibabacloudstack_cloud_monitor_service_alarm_contact_group cmsService.DescribeCmsAlarmContactGroup Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
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
		response, _ := raw.(*cms.PutContactGroupResponse)

		if response.Code != "200" {
			return WrapError(Error("PutContactGroup failed for " + response.Message))
		}
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabacloudStackSdkGoERROR)
		}
	}
	return resourceAlibabacloudstackCmsAlarmContactGroupRead(d, meta)
}
func resourceAlibabacloudstackCmsAlarmContactGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	request := cms.CreateDeleteContactGroupRequest()
	request.ContactGroupName = d.Id()
	raw, err := client.WithCmsClient(func(cmsClient *cms.Client) (interface{}, error) {
		return cmsClient.DeleteContactGroup(request)
	})
	addDebug(request.GetActionName(), raw)
	response, _ := raw.(*cms.DeleteContactGroupResponse)

	if response.Code != "200" {
		return WrapError(Error("DeleteContactGroup failed for " + response.Message))
	}
	if err != nil {
		if IsExpectedErrors(err, []string{"400", "403", "404", "ContactNotExists"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabacloudStackSdkGoERROR)
	}
	return nil
}
