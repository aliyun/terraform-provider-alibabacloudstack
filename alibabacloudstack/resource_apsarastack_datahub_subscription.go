package alibabacloudstack

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/aliyun-datahub-sdk-go/datahub"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlibabacloudStackDatahubSubscription() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlibabacloudStackDatahubSubscriptionCreate,
		Read:   resourceAlibabacloudStackDatahubSubscriptionRead,
		Update: resourceAlibabacloudStackDatahubSubscriptionUpdate,
		Delete: resourceAlibabacloudStackDatahubSubscriptionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"project_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(3, 32),
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return strings.ToLower(new) == strings.ToLower(old)
				},
			},
			"topic_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 128),
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return strings.ToLower(new) == strings.ToLower(old)
				},
			},
			"comment": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(0, 255),
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return strings.ToLower(new) == strings.ToLower(old)
				},
			},
			"sub_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"last_modify_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAlibabacloudStackDatahubSubscriptionCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	projectName := d.Get("project_name").(string)
	topicName := d.Get("topic_name").(string)
	subComment := d.Get("comment").(string)

	request := client.NewCommonRequest("GET", "datahub", "2019-11-20", "CreateSubscription", "")
	request.QueryParams["ProjectName"] = projectName
	request.QueryParams["TopicName"] = topicName
	request.QueryParams["Application"] = "CreateSubscription"
	request.QueryParams["Comment"] = subComment

	raw, err := client.WithEcsClient(func(dataHubClient *ecs.Client) (interface{}, error) {
		return dataHubClient.ProcessCommonRequest(request)
	})
	var subscription *datahub.SubscriptionCreate
	bresponse, ok := raw.(*responses.CommonResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_datahub_subscription", "CreateSubscription", errmsgs.AlibabacloudStackDatahubSdkGo, errmsg)
	}
	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &subscription)

	if debugOn() {
		requestMap := make(map[string]string)
		requestMap["ProjectName"] = projectName
		requestMap["TopicName"] = topicName
		requestMap["SubComment"] = subComment
		addDebug("CreateSubscription", raw, nil, requestMap)
	}
	d.SetId(fmt.Sprintf("%s%s%s%s%s", strings.ToLower(projectName), COLON_SEPARATED, strings.ToLower(topicName), COLON_SEPARATED, subscription.SubId))
	return resourceAlibabacloudStackDatahubSubscriptionRead(d, meta)
}

func resourceAlibabacloudStackDatahubSubscriptionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	datahubService := DatahubService{client}
	parts, err := ParseResourceId(d.Id(), 3)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	projectName := parts[0]
	TopicName := parts[1]
	SubId := parts[2]

	object, err := datahubService.DescribeDatahubSubscription(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}
	d.SetId(fmt.Sprintf("%s%s%s%s%s", strings.ToLower(projectName), COLON_SEPARATED, strings.ToLower(TopicName), COLON_SEPARATED, SubId))

	d.Set("project_name", projectName)
	d.Set("topic_name", TopicName)
	d.Set("sub_id", SubId)
	d.Set("comment", object.Comment)
	d.Set("create_time", strconv.FormatInt(object.CreateTime, 10))
	d.Set("last_modify_time", strconv.FormatInt(object.LastModifyTime, 10))
	return nil
}

func resourceAlibabacloudStackDatahubSubscriptionUpdate(d *schema.ResourceData, meta interface{}) error {
	//parts, err := ParseResourceId(d.Id(), 3)
	//if err != nil {
	//	return errmsgs.WrapError(err)
	//}
	//projectName, topicName, subId := parts[0], parts[1], parts[2]
	//client := meta.(*connectivity.AlibabacloudStackClient)
	//
	//if d.HasChange("comment") {
	//	subComment := d.Get("comment").(string)
	//
	//	var requestInfo *datahub.DataHub
	//
	//	raw, err := client.WithDataHubClient(func(dataHubClient datahub.DataHubApi) (interface{}, error) {
	//		requestInfo = dataHubClient.(*datahub.DataHub)
	//		return dataHubClient.UpdateSubscription(projectName, topicName, subId, subComment)
	//	})
	//	if err != nil {
	//		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), "UpdateSubscription", errmsgs.AlibabacloudStackDatahubSdkGo)
	//	}
	//	if debugOn() {
	//		requestMap := make(map[string]string)
	//		requestMap["ProjectName"] = projectName
	//		requestMap["TopicName"] = topicName
	//		requestMap["SubId"] = subId
	//		requestMap["SubComment"] = subComment
	//		addDebug("UpdateSubscription", raw, requestInfo, requestMap)
	//	}
	//}

	return resourceAlibabacloudStackDatahubSubscriptionRead(d, meta)
}

func resourceAlibabacloudStackDatahubSubscriptionDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	datahubService := DatahubService{client}

	parts, err := ParseResourceId(d.Id(), 3)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	projectName, topicName, subId := parts[0], parts[1], parts[2]

	request := client.NewCommonRequest("GET", "datahub", "2019-11-20", "DeleteSubscription", "")
	request.QueryParams["ProjectName"] = projectName
	request.QueryParams["TopicName"] = topicName
	request.QueryParams["SubscriptionId"] = subId

	var requestInfo *datahub.DataHub

	err = resource.Retry(3*time.Minute, func() *resource.RetryError {
		raw, err := client.WithEcsClient(func(dataHubClient *ecs.Client) (interface{}, error) {
			return dataHubClient.ProcessCommonRequest(request)
		})
		bresponse, ok := raw.(*responses.CommonResponse)
		if err != nil {
			errmsg := ""
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			}
			if isRetryableDatahubError(err) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), "DeleteSubscription", errmsgs.AlibabacloudStackDatahubSdkGo, errmsg))
		}
		if debugOn() {
			requestMap := make(map[string]string)
			requestMap["ProjectName"] = projectName
			requestMap["TopicName"] = topicName
			requestMap["SubId"] = subId
			addDebug("DeleteSubscription", raw, requestInfo, requestMap)
		}
		return nil
	})
	if err != nil {
		if isDatahubNotExistError(err) {
			return nil
		}
		return err
	}
	return errmsgs.WrapError(datahubService.WaitForDatahubSubscription(d.Id(), Deleted, DefaultTimeout))
}
