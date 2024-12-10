package alibabacloudstack

import (
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/cloudapi"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
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
				Type:         schema.TypeString,
				Optional:      true,
				Computed:      true,
				Deprecated:   "Field 'name' is deprecated and will be removed in a future release. Please use new field 'api_group_name' instead.",
				ConflictsWith: []string{"api_group_name"},
			},
			"api_group_name": {
				Type:     schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"name"},
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
	client.InitRpcRequest(*request.RpcRequest)
	request.GroupName = connectivity.GetResourceData(d, "api_group_name", "name").(string)
	if err := errmsgs.CheckEmpty(request.GroupName, schema.TypeString, "api_name", "name"); err != nil {
		return errmsgs.WrapError(err)
	}
	request.Description = d.Get("description").(string)

	if err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := client.WithCloudApiClient(func(cloudApiClient *cloudapi.Client) (interface{}, error) {
			return cloudApiClient.CreateApiGroup(request)
		})
		if err != nil {
			if errmsgs.IsExpectedErrors(err, []string{"RepeatedCommit"}) {
				return resource.RetryableError(err)
			}
			errmsg := ""
			if raw != nil {
				response, ok := raw.(*cloudapi.CreateApiGroupResponse)
				if ok {
					errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
				}
			}
			return resource.NonRetryableError(errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_api_gateway_group", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg))
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ := raw.(*cloudapi.CreateApiGroupResponse)
		d.SetId(response.GroupId)
		return nil
	}); err != nil {
		return err
	}

	return resourceAlibabacloudStackApigatewayGroupRead(d, meta)
}

func resourceAlibabacloudStackApigatewayGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	cloudApiService := CloudApiService{client}
	apiGroup, err := cloudApiService.DescribeApiGatewayGroup(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}

	connectivity.SetResourceData(d, apiGroup.GroupName, "api_group_name", "name")
	d.Set("description", apiGroup.Description)
	d.Set("sub_domain", apiGroup.SubDomain)
	d.Set("vpc_domain", apiGroup.VpcDomain)

	return nil
}

func resourceAlibabacloudStackApigatewayGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	if d.HasChanges("api_group_name", "description") {
		request := cloudapi.CreateModifyApiGroupRequest()
		client.InitRpcRequest(*request.RpcRequest)
		request.GroupId = d.Id()
		request.GroupName = connectivity.GetResourceData(d, "api_group_name", "name").(string)
		request.Description = d.Get("description").(string)

		raw, err := client.WithCloudApiClient(func(cloudApiClient *cloudapi.Client) (interface{}, error) {
			return cloudApiClient.ModifyApiGroup(request)
		})
		if err != nil {
			errmsg := ""
			if raw != nil {
				response, ok := raw.(*cloudapi.ModifyApiGroupResponse)
				if ok {
					errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
				}
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	}
	return resourceAlibabacloudStackApigatewayGroupRead(d, meta)
}

func resourceAlibabacloudStackApigatewayGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	cloudApiService := CloudApiService{client}
	request := cloudapi.CreateDeleteApiGroupRequest()
	client.InitRpcRequest(*request.RpcRequest)
	request.GroupId = d.Id()

	raw, err := client.WithCloudApiClient(func(cloudApiClient *cloudapi.Client) (interface{}, error) {
		return cloudApiClient.DeleteApiGroup(request)
	})
	if err != nil {
		errmsg := ""
		if raw != nil {
			response, ok := raw.(*cloudapi.DeleteApiGroupResponse)
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
			}
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	return errmsgs.WrapError(cloudApiService.WaitForApiGatewayGroup(d.Id(), Deleted, DefaultTimeout))
}
