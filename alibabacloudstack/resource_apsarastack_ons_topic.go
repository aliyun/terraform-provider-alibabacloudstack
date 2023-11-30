package alibabacloudstack

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackOnsTopic() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlibabacloudStackOnsTopicCreate,
		Read:   resourceAlibabacloudStackOnsTopicRead,
		Update: resourceAlibabacloudStackOnsTopicUpdate,
		Delete: resourceAlibabacloudStackOnsTopicDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

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
}

func resourceAlibabacloudStackOnsTopicCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	ordertype := d.Get("message_type").(string)
	instanceId := d.Get("instance_id").(string)
	remark := d.Get("remark").(string)
	topic := d.Get("topic").(string)
	request := requests.NewCommonRequest()
	request.QueryParams = map[string]string{
		"RegionId":        client.RegionId,
		"AccessKeySecret": client.SecretKey,
		"Department":      client.Department,
		"ResourceGroup":   client.ResourceGroup,
		"Product":         "Ons-inner",
		"Action":          "ConsoleTopicCreate",
		"Version":         "2018-02-05",
		"ProductName":     "Ons-inner",
		"PreventCache":    "",
		"OrderType":       ordertype,
		"Topic":           topic,
		"Remark":          remark,
		"OnsRegionId":     client.RegionId,
		"InstanceId":      instanceId,
	}
	request.Method = "POST"
	request.Product = "Ons-inner"
	request.Version = "2018-02-05"
	request.ServiceCode = "Ons-inner"
	request.Domain = client.Domain
	if strings.ToLower(client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.ApiName = "ConsoleTopicCreate"
	request.RegionId = client.RegionId
	request.Headers = map[string]string{"RegionId": client.RegionId}
	response := TopicStruct{}

	raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.ProcessCommonRequest(request)
	})
	if err != nil {
		return WrapErrorf(err, DataDefaultErrorMsg, "alibabacloudstack_ascm_ons_topic", request.GetActionName(), AlibabacloudStackSdkGoERROR)
	}

	bresponse, _ := raw.(*responses.CommonResponse)
	if bresponse.IsSuccess() != true {
		return WrapErrorf(err, DefaultErrorMsg, "alibabacloudstack_ascm_ons_topic", "ConsoleTopicCreate", AlibabacloudStackSdkGoERROR)
	}
	_ = json.Unmarshal(bresponse.GetHttpContentBytes(), &response)
	if response.Success != true {
		return WrapErrorf(errors.New(response.Message), DefaultErrorMsg, "alibabacloudstack_ascm_ons_topic", "ConsoleTopicCreate", AlibabacloudStackSdkGoERROR)
	}
	d.SetId(topic + COLON_SEPARATED + instanceId)

	return resourceAlibabacloudStackOnsTopicRead(d, meta)
}

func resourceAlibabacloudStackOnsTopicRead(d *schema.ResourceData, meta interface{}) error {
	waitSecondsIfWithTest(1)
	client := meta.(*connectivity.AlibabacloudStackClient)
	onsService := OnsService{client}

	object, err := onsService.DescribeOnsTopic(d.Id())
	if err != nil {
		// Handle exceptions
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("instance_id", object.Data[0].NamespaceID)
	d.Set("topic", object.Data[0].Topic)
	d.Set("message_type", strconv.Itoa(object.Data[0].OrderType))
	d.Set("remark", object.Data[0].Remark)

	return nil
}

func resourceAlibabacloudStackOnsTopicUpdate(d *schema.ResourceData, meta interface{}) error {
	return resourceAlibabacloudStackOnsTopicRead(d, meta)
}

func resourceAlibabacloudStackOnsTopicDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	onsService := OnsService{client}
	var requestInfo *ecs.Client

	check, err := onsService.DescribeOnsTopic(d.Id())
	parts, err := ParseResourceId(d.Id(), 2)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "IsTopicExist", AlibabacloudStackSdkGoERROR)
	}
	addDebug("IsTopicExist", check, requestInfo, map[string]string{"Topic": d.Id()})
	err = resource.Retry(2*time.Minute, func() *resource.RetryError {

		request := requests.NewCommonRequest()
		request.QueryParams = map[string]string{
			"RegionId":        client.RegionId,
			"AccessKeySecret": client.SecretKey,
			"Department":      client.Department,
			"ResourceGroup":   client.ResourceGroup,
			"Product":         "Ons-inner",
			"Action":          "ConsoleTopicDelete",
			"Version":         "2018-02-05",
			"ProductName":     "Ons-inner",
			"PreventCache":    "",
			"Topic":           parts[0],
			"OnsRegionId":     client.RegionId,
			"InstanceId":      parts[1],
		}

		request.Method = "POST"
		request.Product = "Ons-inner"
		request.Version = "2018-02-05"
		request.ServiceCode = "Ons-inner"
		request.Domain = client.Domain
		if strings.ToLower(client.Config.Protocol) == "https" {
			request.Scheme = "https"
		} else {
			request.Scheme = "http"
		}
		request.ApiName = "ConsoleTopicDelete"
		request.Headers = map[string]string{"RegionId": client.RegionId}
		request.RegionId = client.RegionId

		_, err := client.WithEcsClient(func(csClient *ecs.Client) (interface{}, error) {
			return csClient.ProcessCommonRequest(request)
		})
		if err != nil {
			return resource.RetryableError(err)
		}
		check, err = onsService.DescribeOnsTopic(d.Id())

		if err != nil {
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alibabacloudstack_ascm_ons_topic", "ConsoleTopicDelete", AlibabacloudStackSdkGoERROR)
	}
	return nil
}
