package alibabacloudstack

import (
	"fmt"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/cloudapi"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlibabacloudStackApigatewayAppAttachment() *schema.Resource {
	resource := &schema.Resource{
		Schema: map[string]*schema.Schema{

			"app_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"group_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"api_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"stage_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"PRE", "RELEASE", "TEST"}, false),
			},
		},
	}
	setResourceFunc(resource, 
		resourceAlibabacloudStackApigatewayAppAttachmentCreate,
		resourceAlibabacloudStackApigatewayAppAttachmentRead,
		nil,
		resourceAlibabacloudStackApigatewayAppAttachmentDelete)
	return resource
}

func resourceAlibabacloudStackApigatewayAppAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	cloudApiService := CloudApiService{client}

	apiId := d.Get("api_id").(string)
	groupId := d.Get("group_id").(string)
	stageName := d.Get("stage_name").(string)
	appId := d.Get("app_id").(string)

	request := cloudapi.CreateSetAppsAuthoritiesRequest()
	client.InitRpcRequest(*request.RpcRequest)
	request.GroupId = groupId
	request.ApiId = apiId
	request.AppIds = appId
	request.StageName = stageName

	raw, err := client.WithCloudApiClient(func(cloudApiClient *cloudapi.Client) (interface{}, error) {
		return cloudApiClient.SetAppsAuthorities(request)
	})
	var response *cloudapi.SetAppsAuthoritiesResponse
	var ok bool
	if err != nil {
		if raw != nil {
			response, ok = raw.(*cloudapi.SetAppsAuthoritiesResponse)
			if ok {
				errmsg := errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
				return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_apigateway_app_attachment", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
			}
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_apigateway_app_attachment", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, "")
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	id := fmt.Sprintf("%s%s%s%s%s%s%s", groupId, COLON_SEPARATED, apiId, COLON_SEPARATED, appId, COLON_SEPARATED, stageName)

	err = cloudApiService.WaitForApiGatewayAppAttachment(id, Normal, DefaultTimeout)
	if err != nil {
		return errmsgs.WrapError(err)
	}

	d.SetId(id)
	return nil
}

func resourceAlibabacloudStackApigatewayAppAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	cloudApiService := CloudApiService{client}

	_, err := cloudApiService.DescribeApiGatewayAppAttachment(d.Id())
	if err != nil {
		return errmsgs.WrapError(err)
	}

	parts, err := ParseResourceId(d.Id(), 4)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	d.Set("group_id", parts[0])
	d.Set("api_id", parts[1])
	d.Set("app_id", parts[2])
	d.Set("stage_name", parts[3])

	return nil
}

func resourceAlibabacloudStackApigatewayAppAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	cloudApiService := CloudApiService{client}

	request := cloudapi.CreateRemoveAppsAuthoritiesRequest()
	client.InitRpcRequest(*request.RpcRequest)
	parts, err := ParseResourceId(d.Id(), 4)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	request.GroupId = parts[0]
	request.ApiId = parts[1]
	request.AppIds = parts[2]
	request.StageName = parts[3]

	raw, err := client.WithCloudApiClient(func(cloudApiClient *cloudapi.Client) (interface{}, error) {
		return cloudApiClient.RemoveAppsAuthorities(request)
	})
	var response *cloudapi.RemoveAppsAuthoritiesResponse
	var ok bool
	if err != nil {
		if errmsgs.IsExpectedErrors(err, []string{"NotFoundAuthorization"}) {
			return nil
		}
		if raw != nil {
			response, ok = raw.(*cloudapi.RemoveAppsAuthoritiesResponse)
			if ok {
				errmsg := errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
				return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
			}
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, "")
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	return errmsgs.WrapError(cloudApiService.WaitForApiGatewayAppAttachment(d.Id(), Deleted, DefaultLongTimeout))
}
