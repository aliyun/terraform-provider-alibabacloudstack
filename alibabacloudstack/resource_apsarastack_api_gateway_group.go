package alibabacloudstack

import (
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/cloudapi"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackApigatewayGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlibabacloudStackApigatewayGroupCreate,
		Read:   resourceAlibabacloudStackApigatewayGroupRead,
		Update: resourceAlibabacloudStackApigatewayGroupUpdate,
		Delete: resourceAlibabacloudStackApigatewayGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"description": {
				Type:     schema.TypeString,
				Required: true,
			},
			"sub_domain": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vpc_domain": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAlibabacloudStackApigatewayGroupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	request := cloudapi.CreateCreateApiGroupRequest()
	request.Headers = map[string]string{
		"RegionId": client.RegionId,
	}
	request.QueryParams = map[string]string{
		"AccessKeySecret": client.SecretKey,
		"AccessKeyId":     client.AccessKey,
		"Product":         "CloudAPI",
		"RegionId":        client.RegionId,
		"Department":      client.Department,
		"ResourceGroup":   client.ResourceGroup,
		"Action":          "CreateApiGroup",
		"Version":         "2016-07-14",
	}
	request.RegionId = client.RegionId
	request.GroupName = d.Get("name").(string)
	request.Description = d.Get("description").(string)

	if err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := client.WithCloudApiClient(func(cloudApiClient *cloudapi.Client) (interface{}, error) {
			return cloudApiClient.CreateApiGroup(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{"RepeatedCommit"}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ := raw.(*cloudapi.CreateApiGroupResponse)
		d.SetId(response.GroupId)
		return nil
	}); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alibabacloudstack_api_gateway_group", request.GetActionName(), AlibabacloudStackSdkGoERROR)
	}

	return resourceAlibabacloudStackApigatewayGroupRead(d, meta)
}

func resourceAlibabacloudStackApigatewayGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	cloudApiService := CloudApiService{client}
	apiGroup, err := cloudApiService.DescribeApiGatewayGroup(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("name", apiGroup.GroupName)
	d.Set("description", apiGroup.Description)
	d.Set("sub_domain", apiGroup.SubDomain)
	d.Set("vpc_domain", apiGroup.VpcDomain)

	return nil
}

func resourceAlibabacloudStackApigatewayGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	if d.HasChange("name") || d.HasChange("description") {
		request := cloudapi.CreateModifyApiGroupRequest()
		request.RegionId = client.RegionId
		request.Description = d.Get("description").(string)
		request.GroupName = d.Get("name").(string)
		request.GroupId = d.Id()
		request.Headers = map[string]string{
			"RegionId": client.RegionId,
		}
		request.QueryParams = map[string]string{
			"AccessKeySecret": client.SecretKey,
			"AccessKeyId":     client.AccessKey,
			"Product":         "CloudAPI",
			"RegionId":        client.RegionId,
			"Department":      client.Department,
			"ResourceGroup":   client.ResourceGroup,
			"Action":          "ModifyApiGroup",
			"Version":         "2016-07-14",
		}
		raw, err := client.WithCloudApiClient(func(cloudApiClient *cloudapi.Client) (interface{}, error) {
			return cloudApiClient.ModifyApiGroup(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabacloudStackSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	}
	return resourceAlibabacloudStackApigatewayGroupRead(d, meta)
}

func resourceAlibabacloudStackApigatewayGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	cloudApiService := CloudApiService{client}
	request := cloudapi.CreateDeleteApiGroupRequest()
	request.RegionId = client.RegionId
	request.GroupId = d.Id()
	request.Headers = map[string]string{
		"RegionId": client.RegionId,
	}
	request.QueryParams = map[string]string{
		"AccessKeySecret": client.SecretKey,
		"AccessKeyId":     client.AccessKey,
		"Product":         "CloudAPI",
		"RegionId":        client.RegionId,
		"Department":      client.Department,
		"ResourceGroup":   client.ResourceGroup,
		"Action":          "DeleteApiGroup",
		"Version":         "2016-07-14",
	}
	raw, err := client.WithCloudApiClient(func(cloudApiClient *cloudapi.Client) (interface{}, error) {
		return cloudApiClient.DeleteApiGroup(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabacloudStackSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	return WrapError(cloudApiService.WaitForApiGatewayGroup(d.Id(), Deleted, DefaultTimeout))

}
