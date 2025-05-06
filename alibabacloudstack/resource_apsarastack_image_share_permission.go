package alibabacloudstack

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackImageSharePermission() *schema.Resource {
	resource := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"image_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"account_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
	setResourceFunc(resource, resourceAlibabacloudStackImageSharePermissionCreate,
		resourceAlibabacloudStackImageSharePermissionRead, nil, resourceAlibabacloudStackImageSharePermissionDelete)
	return resource
}

func resourceAlibabacloudStackImageSharePermissionCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	imageId := d.Get("image_id").(string)
	accountId := d.Get("account_id").(string)
	request := ecs.CreateModifyImageSharePermissionRequest()
	client.InitRpcRequest(*request.RpcRequest)
	request.ImageId = imageId
	accountSli := []string{accountId}
	request.AddAccount = &accountSli
	raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.ModifyImageSharePermission(request)
	})
	response, ok := raw.(*ecs.ModifyImageSharePermissionResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_image_share_permission", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	d.SetId(imageId + ":" + accountId)
	return nil
}

func resourceAlibabacloudStackImageSharePermissionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	ecsService := EcsService{client: client}
	object, err := ecsService.DescribeImageShareByImageId(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}
	parts, err := ParseResourceId(d.Id(), 2)
	d.Set("image_id", object.ImageId)
	d.Set("account_id", parts[1])
	return errmsgs.WrapError(err)
}

func resourceAlibabacloudStackImageSharePermissionDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	request := ecs.CreateModifyImageSharePermissionRequest()
	client.InitRpcRequest(*request.RpcRequest)
	parts, err := ParseResourceId(d.Id(), 2)
	request.ImageId = parts[0]
	accountSli := []string{parts[1]}
	request.RemoveAccount = &accountSli
	raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.ModifyImageSharePermission(request)
	})
	response, ok := raw.(*ecs.ModifyImageSharePermissionResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_image_share_permission", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	return nil
}
