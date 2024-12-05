package alibabacloudstack

import (
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/cloudapi"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
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
	client.InitRpcRequest(*request.RpcRequest)
	request.AppName = d.Get("name").(string)
	if v, exist := d.GetOk("description"); exist {
		request.Description = v.(string)
	}
	request.Description = d.Get("description").(string)
	request.QueryParams["Product"] = "CloudAPI"
	request.QueryParams["Action"] = "CreateApp"
	request.QueryParams["Version"] = "2016-07-14"

	if err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := client.WithCloudApiClient(func(cloudApiClient *cloudapi.Client) (interface{}, error) {
			return cloudApiClient.CreateApp(request)
		})
		bresponse, ok := raw.(*cloudapi.CreateAppResponse)
		if err != nil {
			if errmsgs.IsExpectedErrors(err, []string{"RepeatedCommit"}) {
				return resource.RetryableError(err)
			}
			errmsg := ""
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			}
			return resource.NonRetryableError(errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_apigateway_app", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg))
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		d.SetId(strconv.FormatInt(bresponse.AppId, 10))
		return nil
	}); err != nil {
		return err
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
			if errmsgs.IsExpectedErrors(err, []string{"NotFoundResourceId"}) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		d.Set("tags", cloudApiService.tagsToMap(tags))
		return nil
	}); err != nil {
		return errmsgs.WrapError(err)
	}

	if err := resource.Retry(3*time.Second, func() *resource.RetryError {
		object, err := cloudApiService.DescribeApiGatewayApp(d.Id())
		if err != nil {
			if errmsgs.NotFoundError(err) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}

		d.Set("name", object.AppName)
		d.Set("description", object.Description)
		return nil
	}); err != nil {
		return errmsgs.WrapError(err)
	}
	return nil
}

func resourceAlibabacloudStackApigatewayAppUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	cloudApiService := CloudApiService{client}
	if err := cloudApiService.setInstanceTags(d, TagResourceApp); err != nil {
		return errmsgs.WrapError(err)
	}
	if d.IsNewResource() {
		d.Partial(false)
		return resourceAlibabacloudStackApigatewayAppRead(d, meta)
	}
	if d.HasChange("name") || d.HasChange("description") {
		request := cloudapi.CreateModifyAppRequest()
		client.InitRpcRequest(*request.RpcRequest)
		request.AppId = requests.Integer(d.Id())
		request.AppName = d.Get("name").(string)
		if v, exist := d.GetOk("description"); exist {
			request.Description = v.(string)
		}
		request.QueryParams["Product"] = "CloudAPI"
		request.QueryParams["Action"] = "ModifyApp"
		request.QueryParams["Version"] = "2016-07-14"

		raw, err := client.WithCloudApiClient(func(cloudApiClient *cloudapi.Client) (interface{}, error) {
			return cloudApiClient.ModifyApp(request)
		})
		bresponse, ok := raw.(*cloudapi.ModifyAppResponse)
		if err != nil {
			errmsg := ""
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
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
	client.InitRpcRequest(*request.RpcRequest)
	request.AppId = requests.Integer(d.Id())
	request.QueryParams["Product"] = "CloudAPI"
	request.QueryParams["Action"] = "DeleteApp"
	request.QueryParams["Version"] = "2016-07-14"

	raw, err := client.WithCloudApiClient(func(cloudApiClient *cloudapi.Client) (interface{}, error) {
		return cloudApiClient.DeleteApp(request)
	})
	bresponse, ok := raw.(*cloudapi.DeleteAppResponse)
	if err != nil {
		if errmsgs.IsExpectedErrors(err, []string{"NotFoundApp"}) {
			return nil
		}
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	return errmsgs.WrapError(cloudApiService.WaitForApiGatewayApp(d.Id(), Deleted, DefaultTimeout))
}
