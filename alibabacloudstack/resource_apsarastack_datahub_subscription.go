package alibabacloudstack

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/aliyun/aliyun-datahub-sdk-go/datahub"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/helper/sdk_patch/datahub_patch"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlibabacloudStackDatahubSubscription() *schema.Resource {
	resource := &schema.Resource{
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
			"application_name":{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`^[a-z][a-z0-9_]{0,31}$`), "application_name must be '^[a-z][a-z0-9_]{0,31}$'"),
			},
			"comment": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
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
	setResourceFunc(resource, resourceAlibabacloudStackDatahubSubscriptionCreate, 
		resourceAlibabacloudStackDatahubSubscriptionRead, nil, resourceAlibabacloudStackDatahubSubscriptionDelete)
	return resource
}

func resourceAlibabacloudStackDatahubSubscriptionCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	projectName := d.Get("project_name").(string)
	topicName := d.Get("topic_name").(string)
	subComment := d.Get("comment").(string)

	request := client.NewCommonRequest("GET", "datahub", "2019-11-20", "CreateSubscription", "")
	request.QueryParams["ProjectName"] = projectName
	request.QueryParams["TopicName"] = topicName
	request.QueryParams["Application"] = d.Get("application_name").(string)
	request.QueryParams["Comment"] = subComment

	bresponse, err := client.ProcessCommonRequest(request)
	if debugOn() {
		requestMap := make(map[string]string)
		requestMap["ProjectName"] = projectName
		requestMap["TopicName"] = topicName
		requestMap["SubComment"] = subComment
		addDebug("CreateSubscription", bresponse, nil, requestMap)
	}
	var subscription *datahub_patch.SubscriptionCreate
	if err != nil {
		if bresponse == nil {
			return errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_datahub_subscription", "CreateSubscription", errmsgs.AlibabacloudStackDatahubSdkGo, errmsg)
	}
	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &subscription)

	d.SetId(fmt.Sprintf("%s:%s:%s", strings.ToLower(projectName), strings.ToLower(topicName), subscription.SubId))
	return nil
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

	d.Set("project_name", projectName)
	d.Set("topic_name", TopicName)
	d.Set("sub_id", SubId)
	d.Set("comment", object.Comment)
	d.Set("application_name", object.Application)
	return nil
}

func resourceAlibabacloudStackDatahubSubscriptionDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

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

	bresponse, err := client.ProcessCommonRequest(request)
	if debugOn() {
		requestMap := make(map[string]string)
		requestMap["ProjectName"] = projectName
		requestMap["TopicName"] = topicName
		requestMap["SubId"] = subId
		addDebug("DeleteSubscription", bresponse, requestInfo, requestMap)
	}
	if err != nil {
		if isDatahubNotExistError(err) {
			return nil
		}
		if bresponse == nil {
			return errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), "DeleteSubscription", errmsgs.AlibabacloudStackDatahubSdkGo, errmsg)
	}
	return nil
}
