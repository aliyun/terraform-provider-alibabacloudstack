package alibabacloudstack

import (
	"encoding/json"
	"errors"
	"strconv"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ons"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlibabacloudStackOnsTopic() *schema.Resource {
	resource := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"topic": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 128),
			},
			"message_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"remark": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 128),
			},
			"perm": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntInSlice([]int{2, 4, 6}),
			},
		},
	}
	setResourceFunc(resource, resourceAlibabacloudStackOnsTopicCreate, resourceAlibabacloudStackOnsTopicRead, resourceAlibabacloudStackOnsTopicUpdate, resourceAlibabacloudStackOnsTopicDelete)
	return resource
}

func resourceAlibabacloudStackOnsTopicCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	ordertype := d.Get("message_type").(string)
	instanceId := d.Get("instance_id").(string)
	remark := d.Get("remark").(string)
	topic := d.Get("topic").(string)

	request := client.NewCommonRequest("POST", "Ons-inner", "2018-02-05", "ConsoleTopicCreate", "")
	mergeMaps(request.QueryParams, map[string]string{
		"ProductName":  "Ons-inner",
		"PreventCache": "",
		"OrderType":    ordertype,
		"Topic":        topic,
		"Remark":       remark,
		"OnsRegionId":  client.RegionId,
		"InstanceId":   instanceId,
	})
	request.ServiceCode = "Ons-inner"
	request.Domain = client.Domain

	response := TopicStruct{}

	raw, err := client.WithOnsClient(func(onsClient *ons.Client) (interface{}, error) {
		return onsClient.ProcessCommonRequest(request)
	})
	bresponse, ok := raw.(*responses.CommonResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_ascm_ons_topic", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	if !bresponse.IsSuccess() {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		return errmsgs.WrapErrorf(errors.New(bresponse.GetHttpContentString()), errmsgs.RequestV1ErrorMsg, "alibabacloudstack_ascm_ons_topic", "ConsoleTopicCreate", errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	_ = json.Unmarshal(bresponse.GetHttpContentBytes(), &response)
	if !response.Success {
		return errmsgs.WrapErrorf(errors.New(response.Message), errmsgs.DefaultErrorMsg, "alibabacloudstack_ascm_ons_topic", "ConsoleTopicCreate", errmsgs.AlibabacloudStackSdkGoERROR)
	}
	d.SetId(topic + COLON_SEPARATED + instanceId)

	return nil
}

func resourceAlibabacloudStackOnsTopicRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	onsService := OnsService{client}

	object, err := onsService.DescribeOnsTopic(d.Id())
	if err != nil {
		// Handle exceptions
		if errmsgs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}

	d.Set("instance_id", object.Data[0].NamespaceID)
	d.Set("topic", object.Data[0].Topic)
	d.Set("message_type", strconv.Itoa(object.Data[0].OrderType))
	d.Set("remark", object.Data[0].Remark)

	return nil
}

func resourceAlibabacloudStackOnsTopicUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceAlibabacloudStackOnsTopicDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	onsService := OnsService{client}
	var requestInfo *ons.Client

	check, err := onsService.DescribeOnsTopic(d.Id())
	parts, err := ParseResourceId(d.Id(), 2)

	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), "IsTopicExist", errmsgs.AlibabacloudStackSdkGoERROR)
	}
	addDebug("IsTopicExist", check, requestInfo, map[string]string{"Topic": d.Id()})
	err = resource.Retry(2*time.Minute, func() *resource.RetryError {

		request := client.NewCommonRequest("POST", "Ons-inner", "2018-02-05", "ConsoleTopicDelete", "")
		mergeMaps(request.QueryParams, map[string]string{
			"ProductName":  "Ons-inner",
			"PreventCache": "",
			"Topic":        parts[0],
			"OnsRegionId":  client.RegionId,
			"InstanceId":   parts[1],
		})
		request.ServiceCode = "Ons-inner"
		request.Domain = client.Domain

		raw, err := client.WithOnsClient(func(onsClient *ons.Client) (interface{}, error) {
			return onsClient.ProcessCommonRequest(request)
		})
		bresponse, ok := raw.(*responses.CommonResponse)
		if err != nil {
			errmsg := ""
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			}
			return resource.RetryableError(errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_ascm_ons_topic", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg))
		}
		check, err = onsService.DescribeOnsTopic(d.Id())

		if err != nil {
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_ascm_ons_topic", "ConsoleTopicDelete", errmsgs.AlibabacloudStackSdkGoERROR)
	}
	return nil
}