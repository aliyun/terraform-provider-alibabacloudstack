package alibabacloudstack

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ots"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackOtsInstanceAttachment() *schema.Resource {
	resource := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"instance_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"vpc_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vswitch_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
	setResourceFunc(resource, resourceAliyunOtsInstanceAttachmentCreate, resourceAliyunOtsInstanceAttachmentRead, nil, resourceAliyunOtsInstanceAttachmentDelete)
	return resource
}

func resourceAliyunOtsInstanceAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	vpcService := VpcService{client}

	request := ots.CreateBindInstance2VpcRequest()
	client.InitRpcRequest(*request.RpcRequest)
	request.InstanceName = d.Get("instance_name").(string)
	request.InstanceVpcName = d.Get("vpc_name").(string)
	request.VirtualSwitchId = d.Get("vswitch_id").(string)

	if vsw, err := vpcService.DescribeVSwitch(d.Get("vswitch_id").(string)); err != nil {
		return errmsgs.WrapError(err)
	} else {
		request.VpcId = vsw.VpcId
	}

	raw, err := client.WithOtsClient(func(otsClient *ots.Client) (interface{}, error) {
		return otsClient.BindInstance2Vpc(request)
	})
	if err != nil {
		errmsg := ""
		if response, ok := raw.(*ots.BindInstance2VpcResponse); ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_ots_instance_attachment", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	d.SetId(request.InstanceName)
	return nil
}

func resourceAliyunOtsInstanceAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	otsService := OtsService{client}
	object, err := otsService.DescribeOtsInstanceAttachment(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}
	// There is a bug that inst does not contain instance name and vswitch ID, so this resource does not support import function.
	//d.Set("instance_name", inst.InstanceName)
	d.Set("vpc_name", object.InstanceVpcName)
	d.Set("vpc_id", object.VpcId)
	return nil
}

func resourceAliyunOtsInstanceAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	otsService := OtsService{client}
	object, err := otsService.DescribeOtsInstanceAttachment(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			return nil
		}
		return errmsgs.WrapError(err)
	}
	request := ots.CreateUnbindInstance2VpcRequest()
	client.InitRpcRequest(*request.RpcRequest)
	request.InstanceName = d.Id()
	request.InstanceVpcName = object.InstanceVpcName

	raw, err := client.WithOtsClient(func(otsClient *ots.Client) (interface{}, error) {
		return otsClient.UnbindInstance2Vpc(request)
	})
	if err != nil {
		if errmsgs.NotFoundError(err) {
			return nil
		}
		errmsg := ""
		if response, ok := raw.(*ots.UnbindInstance2VpcResponse); ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	return errmsgs.WrapError(otsService.WaitForOtsInstanceVpc(d.Id(), Deleted, DefaultTimeout))
}
