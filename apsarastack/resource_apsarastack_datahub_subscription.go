package apsarastack

import (
	"encoding/json"
	//"encoding/json"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/aliyun/aliyun-datahub-sdk-go/datahub"

	"github.com/apsara-stack/terraform-provider-apsarastack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceApsaraStackDatahubSubscription() *schema.Resource {
	return &schema.Resource{
		Create: resourceApsaraStackDatahubSubscriptionCreate,
		Read:   resourceApsaraStackDatahubSubscriptionRead,
		Update: resourceApsaraStackDatahubSubscriptionUpdate,
		Delete: resourceApsaraStackDatahubSubscriptionDelete,
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
				//Default:      "subscription added by terraform",
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
				Type:     schema.TypeString, //uint64 value from sdk
				Computed: true,
			},
			"last_modify_time": {
				Type:     schema.TypeString, //uint64 value from sdk
				Computed: true,
			},
		},
	}
}

func resourceApsaraStackDatahubSubscriptionCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	subscription := &datahub.SubscriptionCreate{}
	projectName := d.Get("project_name").(string)
	topicName := d.Get("topic_name").(string)
	subComment := d.Get("comment").(string)

	request := requests.NewCommonRequest()
	request.Method = "GET"
	request.Product = "datahub"
	request.Domain = client.Domain
	request.Version = "2019-11-20"
	if strings.ToLower(client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.ApiName = "CreateSubscription"
	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.QueryParams = map[string]string{
		"AccessKeySecret": client.SecretKey,
		"AccessKeyId":     client.AccessKey,
		"Product":         "datahub",
		"RegionId":        client.RegionId,
		"Department":      client.Department,
		"ResourceGroup":   client.ResourceGroup,
		"Action":          "CreateSubscription",
		"Version":         "2019-11-20",
		"ProjectName":     projectName,
		"TopicName":       topicName,
		"Application":     "CreateSubscription",
		"Comment":         subComment,
	}

	raw, err := client.WithEcsClient(func(dataHubClient *ecs.Client) (interface{}, error) {
		return dataHubClient.ProcessCommonRequest(request)
	})
	var requestInfo *datahub.DataHub
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "apsarastack_datahub_subscription", "CreateSubscription", ApsaraStackDatahubSdkGo)
	}
	bresponse := raw.(*responses.CommonResponse)
		err = json.Unmarshal(bresponse.GetHttpContentBytes(), subscription)

	if debugOn() {
		requestMap := make(map[string]string)
		requestMap["ProjectName"] = projectName
		requestMap["TopicName"] = topicName
		requestMap["SubComment"] = subComment
		addDebug("CreateSubscription", raw, requestInfo, requestMap)
	}
	d.SetId(fmt.Sprintf("%s%s%s%s%s", strings.ToLower(projectName), COLON_SEPARATED, strings.ToLower(topicName), COLON_SEPARATED,subscription.SubId))
	return resourceApsaraStackDatahubSubscriptionRead(d, meta)
}

func resourceApsaraStackDatahubSubscriptionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	datahubService := DatahubService{client}
	parts, err := ParseResourceId(d.Id(), 3)
	if err != nil {
		return WrapError(err)
	}
	projectName := parts[0]
	TopicName := parts[1]
	SubId := parts[2]

	object, err := datahubService.DescribeDatahubSubscription(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
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

func resourceApsaraStackDatahubSubscriptionUpdate(d *schema.ResourceData, meta interface{}) error {
	//parts, err := ParseResourceId(d.Id(), 3)
	//if err != nil {
	//	return WrapError(err)
	//}
	//projectName, topicName, subId := parts[0], parts[1], parts[2]
	//client := meta.(*connectivity.ApsaraStackClient)
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
	//		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "UpdateSubscription", ApsaraStackDatahubSdkGo)
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

	return resourceApsaraStackDatahubSubscriptionRead(d, meta)
}

func resourceApsaraStackDatahubSubscriptionDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	datahubService := DatahubService{client}

	parts, err := ParseResourceId(d.Id(), 3)
	if err != nil {
		return WrapError(err)
	}
	projectName, topicName, subId := parts[0], parts[1], parts[2]

	request := requests.NewCommonRequest()
	request.Method = "GET"
	request.Product = "datahub"
	request.Domain = client.Domain
	request.Version = "2019-11-20"
	if strings.ToLower(client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.ApiName = "DeleteSubscription"
	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.QueryParams = map[string]string{
		"AccessKeySecret": client.SecretKey,
		"AccessKeyId":     client.AccessKey,
		"Product":         "datahub",
		"RegionId":        client.RegionId,
		"Department":      client.Department,
		"ResourceGroup":   client.ResourceGroup,
		"Action":          "DeleteSubscription",
		"Version":         "2019-11-20",
		"ProjectName":     projectName,
		"TopicName":         topicName,
		"SubscriptionId":         subId,
	}



	var requestInfo *datahub.DataHub

	err = resource.Retry(3*time.Minute, func() *resource.RetryError {
		raw, err := client.WithEcsClient(func(dataHubClient *ecs.Client) (interface{}, error) {
			return dataHubClient.ProcessCommonRequest(request)
		})
		if err != nil {
			if isRetryableDatahubError(err) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
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
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "DeleteSubscription", ApsaraStackDatahubSdkGo)
	}
	return WrapError(datahubService.WaitForDatahubSubscription(d.Id(), Deleted, DefaultTimeout))
}
