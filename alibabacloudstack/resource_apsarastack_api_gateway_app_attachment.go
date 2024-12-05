package alibabacloudstack

import (
	"fmt"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/cloudapi"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAliyunApigatewayAppAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliyunApigatewayAppAttachmentCreate,
		Read:   resourceAliyunApigatewayAppAttachmentRead,
		Delete: resourceAliyunApigatewayAppAttachmentDelete,

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
}

func resourceAliyunApigatewayAppAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	cloudApiService := CloudApiService{client}

	apiId := d.Get("api_id").(string)
	groupId := d.Get("group_id").(string)
	stageName := d.Get("stage_name").(string)
	appId := d.Get("app_id").(string)

	request := cloudapi.CreateSetAppsAuthoritiesRequest()
	request.RegionId = client.RegionId
	request.GroupId = groupId
	request.ApiId = apiId
	request.AppIds = appId
	request.StageName = stageName
	request.Headers = map[string]string{
		"RegionId": client.RegionId,
	}
	request.QueryParams = map[string]string{
		
		
		"Product":         "CloudAPI",
		"RegionId":        client.RegionId,
		"Department":      client.Department,
		"ResourceGroup":   client.ResourceGroup,
		"Action":          "SetAppsAuthorities",
		"Version":         "2016-07-14",
	}
	raw, err := client.WithCloudApiClient(func(cloudApiClient *cloudapi.Client) (interface{}, error) {
		return cloudApiClient.SetAppsAuthorities(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alibabacloudstack_apigateway_app_attachment", request.GetActionName(), AlibabacloudStackSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	id := fmt.Sprintf("%s%s%s%s%s%s%s", groupId, COLON_SEPARATED, apiId, COLON_SEPARATED, appId, COLON_SEPARATED, stageName)

	err = cloudApiService.WaitForApiGatewayAppAttachment(id, Normal, DefaultTimeout)
	if err != nil {
		return WrapError(err)
	}

	d.SetId(id)
	return resourceAliyunApigatewayAppAttachmentRead(d, meta)
}

func resourceAliyunApigatewayAppAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	cloudApiService := CloudApiService{client}

	_, err := cloudApiService.DescribeApiGatewayAppAttachment(d.Id())
	if err != nil {
		return WrapError(err)
	}

	parts, err := ParseResourceId(d.Id(), 4)
	if err != nil {
		return WrapError(err)
	}
	d.Set("group_id", parts[0])
	d.Set("api_id", parts[1])
	d.Set("app_id", parts[2])
	d.Set("stage_name", parts[3])

	return nil
}

func resourceAliyunApigatewayAppAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	cloudApiService := CloudApiService{client}
	request := cloudapi.CreateRemoveAppsAuthoritiesRequest()
	request.RegionId = client.RegionId
	parts, err := ParseResourceId(d.Id(), 4)
	if err != nil {
		return WrapError(err)
	}
	request.GroupId = parts[0]
	request.ApiId = parts[1]
	request.AppIds = parts[2]
	request.StageName = parts[3]
	request.Headers = map[string]string{
		"RegionId": client.RegionId,
	}
	request.QueryParams = map[string]string{
		
		
		"Product":         "CloudAPI",
		"RegionId":        client.RegionId,
		"Department":      client.Department,
		"ResourceGroup":   client.ResourceGroup,
		"Action":          "RemoveAppsAuthorities",
		"Version":         "2016-07-14",
	}
	raw, err := client.WithCloudApiClient(func(cloudApiClient *cloudapi.Client) (interface{}, error) {
		return cloudApiClient.RemoveAppsAuthorities(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"NotFoundAuthorization"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabacloudStackSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	return WrapError(cloudApiService.WaitForApiGatewayAppAttachment(d.Id(), Deleted, DefaultLongTimeout))
}
