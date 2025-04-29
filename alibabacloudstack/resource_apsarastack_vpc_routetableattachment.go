package alibabacloudstack

import (
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackRouteTableAttachment() *schema.Resource {
	resource := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"route_table_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"vswitch_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
	setResourceFunc(resource, resourceAliyunRouteTableAttachmentCreate, resourceAliyunRouteTableAttachmentRead, nil, resourceAliyunRouteTableAttachmentDelete)
	return resource
}

func resourceAliyunRouteTableAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	vpcService := VpcService{client}

	request := vpc.CreateAssociateRouteTableRequest()
	client.InitRpcRequest(*request.RpcRequest)
	request.RouteTableId = Trim(d.Get("route_table_id").(string))
	request.VSwitchId = Trim(d.Get("vswitch_id").(string))
	request.ClientToken = buildClientToken(request.GetActionName())
	if err := resource.Retry(3*time.Minute, func() *resource.RetryError {
		args := *request
		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.AssociateRouteTable(&args)
		})
		if err != nil {
			if errmsgs.IsExpectedErrors(err, []string{"TaskConflict"}) {
				return resource.RetryableError(err)
			}
			errmsg := ""
			if raw != nil {
				response, ok := raw.(*vpc.AssociateRouteTableResponse)
				if ok {
					errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
				}
			}
			return resource.NonRetryableError(errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_route_table_attachment", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg))
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	}); err != nil {
		return err
	}
	d.SetId(request.RouteTableId + COLON_SEPARATED + request.VSwitchId)
	err := vpcService.WaitForRouteTableAttachment(d.Id(), Available, DefaultTimeout)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	if err := vpcService.WaitForVSwitch(request.VSwitchId, Available, DefaultTimeoutMedium); err != nil {
		return errmsgs.WrapError(err)
	}
	return nil
}

func resourceAliyunRouteTableAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	vpcService := VpcService{client}
	object, err := vpcService.DescribeRouteTableAttachment(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	d.Set("route_table_id", object.RouteTableId)
	d.Set("vswitch_id", parts[1])
	return nil
}

func resourceAliyunRouteTableAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	vpcService := VpcService{client}

	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return errmsgs.WrapError(err)
	}

	request := vpc.CreateUnassociateRouteTableRequest()
	client.InitRpcRequest(*request.RpcRequest)
	request.RouteTableId = parts[0]
	request.VSwitchId = parts[1]
	request.ClientToken = buildClientToken(request.GetActionName())
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		args := *request
		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.UnassociateRouteTable(&args)
		})

		if err != nil {
			if errmsgs.IsExpectedErrors(err, []string{"TaskConflict"}) {
				return resource.RetryableError(err)
			}
			errmsg := ""
			if raw != nil {
				response, ok := raw.(*vpc.UnassociateRouteTableResponse)
				if ok {
					errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
				}
			}
			return resource.NonRetryableError(errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg))
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})
	if err != nil {
		return err
	}
	return errmsgs.WrapError(vpcService.WaitForRouteTableAttachment(d.Id(), Deleted, DefaultTimeoutMedium))
}
