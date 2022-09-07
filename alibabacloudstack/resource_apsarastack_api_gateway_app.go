package alibabacloudstack

import (
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/cloudapi"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackApigatewayApp() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlibabacloudStackApigatewayAppCreate,
		Read:   resourceAlibabacloudStackApigatewayAppRead,
		Update: resourceAlibabacloudStackApigatewayAppUpdate,
		Delete: resourceAlibabacloudStackApigatewayAppDelete,
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
				Optional: true,
			},
			"tags": tagsSchema(),
		},
	}
}

func resourceAlibabacloudStackApigatewayAppCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	request := cloudapi.CreateCreateAppRequest()
	request.RegionId = client.RegionId
	request.AppName = d.Get("name").(string)
	if v, exist := d.GetOk("description"); exist {
		request.Description = v.(string)
	}
	request.Description = d.Get("description").(string)
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
		"Action":          "CreateApp",
		"Version":         "2016-07-14",
	}
	if err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := client.WithCloudApiClient(func(cloudApiClient *cloudapi.Client) (interface{}, error) {
			return cloudApiClient.CreateApp(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{"RepeatedCommit"}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ := raw.(*cloudapi.CreateAppResponse)
		d.SetId(strconv.FormatInt(response.AppId, 10))
		return nil
	}); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alibabacloudstack_apigateway_app", request.GetActionName(), AlibabacloudStackSdkGoERROR)
	}
	return resourceAlibabacloudStackApigatewayAppUpdate(d, meta)
}

func resourceAlibabacloudStackApigatewayAppRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	cloudApiService := CloudApiService{client}

	wait := incrementalWait(3*time.Second, 5*time.Second)
	if err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		tags, err := cloudApiService.DescribeTags(d.Id(), nil, TagResourceApp)
		if err != nil {
			if IsExpectedErrors(err, []string{"NotFoundResourceId"}) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		d.Set("tags", cloudApiService.tagsToMap(tags))
		return nil
	}); err != nil {
		return WrapError(err)
	}

	if err := resource.Retry(3*time.Second, func() *resource.RetryError {
		object, err := cloudApiService.DescribeApiGatewayApp(d.Id())
		if err != nil {
			if NotFoundError(err) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}

		d.Set("name", object.AppName)
		d.Set("description", object.Description)
		return nil
	}); err != nil {
		return WrapError(err)
	}
	return nil
}

func resourceAlibabacloudStackApigatewayAppUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	cloudApiService := CloudApiService{client}
	if err := cloudApiService.setInstanceTags(d, TagResourceApp); err != nil {
		return WrapError(err)
	}
	if d.IsNewResource() {
		d.Partial(false)
		return resourceAlibabacloudStackApigatewayAppRead(d, meta)
	}
	if d.HasChange("name") || d.HasChange("description") {
		request := cloudapi.CreateModifyAppRequest()
		request.RegionId = client.RegionId
		request.AppId = requests.Integer(d.Id())
		request.AppName = d.Get("name").(string)
		if v, exist := d.GetOk("description"); exist {
			request.Description = v.(string)
		}
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
			"Action":          "ModifyApp",
			"Version":         "2016-07-14",
		}
		raw, err := client.WithCloudApiClient(func(cloudApiClient *cloudapi.Client) (interface{}, error) {
			return cloudApiClient.ModifyApp(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabacloudStackSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	}
	time.Sleep(3 * time.Second)
	return resourceAlibabacloudStackApigatewayAppRead(d, meta)
}

func resourceAlibabacloudStackApigatewayAppDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	cloudApiService := CloudApiService{client}

	request := cloudapi.CreateDeleteAppRequest()
	request.RegionId = client.RegionId
	request.AppId = requests.Integer(d.Id())
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
		"Action":          "DeleteApp",
		"Version":         "2016-07-14",
	}
	raw, err := client.WithCloudApiClient(func(cloudApiClient *cloudapi.Client) (interface{}, error) {
		return cloudApiClient.DeleteApp(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"NotFoundApp"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabacloudStackSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	return WrapError(cloudApiService.WaitForApiGatewayApp(d.Id(), Deleted, DefaultTimeout))
}
