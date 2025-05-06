package alibabacloudstack

import (
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceNetworkInterfaceAttachment() *schema.Resource {
	resource := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"network_interface_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
	setResourceFunc(resource, resourceAlibabacloudStackNetworkInterfaceAttachmentCreate, resourceAlibabacloudStackNetworkInterfaceAttachmentRead, nil, resourceAlibabacloudStackNetworkInterfaceAttachmentDelete)
	return resource
}

func resourceAlibabacloudStackNetworkInterfaceAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	ecsService := EcsService{client}

	eniId := d.Get("network_interface_id").(string)
	instanceId := d.Get("instance_id").(string)

	request := ecs.CreateAttachNetworkInterfaceRequest()
	client.InitRpcRequest(*request.RpcRequest)
	request.InstanceId = instanceId
	request.NetworkInterfaceId = eniId

	err := resource.Retry(DefaultTimeout*time.Minute, func() *resource.RetryError {
		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.AttachNetworkInterface(request)
		})
		if err != nil {
			if errmsgs.IsExpectedErrors(err, errmsgs.NetworkInterfaceInvalidOperations) {
				return resource.RetryableError(err)
			}
			errmsg := ""
			if raw != nil {
				response, ok := raw.(*ecs.AttachNetworkInterfaceResponse)
				if ok {
					errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
				}
			}
			return resource.NonRetryableError(errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_network_interface_attachment", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg))
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})
	if err != nil {
		return err
	}
	d.SetId(eniId + COLON_SEPARATED + instanceId)
	if err = ecsService.WaitForNetworkInterface(eniId, InUse, DefaultTimeout); err != nil {
		return errmsgs.WrapError(err)
	}
	return nil
}

func resourceAlibabacloudStackNetworkInterfaceAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	ecsService := EcsService{client}

	object, err := ecsService.DescribeNetworkInterfaceAttachment(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}

	d.Set("instance_id", object.InstanceId)
	d.Set("network_interface_id", object.NetworkInterfaceId)

	return nil
}

func resourceAlibabacloudStackNetworkInterfaceAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	ecsService := EcsService{client}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	eniId, instanceId := parts[0], parts[1]

	_, err = ecsService.DescribeNetworkInterfaceAttachment(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			return nil
		}
		return errmsgs.WrapError(err)
	}

	request := ecs.CreateDetachNetworkInterfaceRequest()
	client.InitRpcRequest(*request.RpcRequest)
	request.InstanceId = instanceId
	request.NetworkInterfaceId = eniId

	err = resource.Retry(DefaultTimeout*time.Second, func() *resource.RetryError {
		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.DetachNetworkInterface(request)
		})
		if err != nil {
			if errmsgs.IsExpectedErrors(err, errmsgs.NetworkInterfaceInvalidOperations) {
				return resource.RetryableError(err)
			}
			errmsg := ""
			if raw != nil {
				response, ok := raw.(*ecs.DetachNetworkInterfaceResponse)
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
	return errmsgs.WrapError(ecsService.WaitForNetworkInterface(eniId, Available, DefaultTimeoutMedium))
}