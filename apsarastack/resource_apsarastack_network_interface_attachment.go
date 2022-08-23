package apsarastack

import (
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	//"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/aliyun/terraform-provider-alibabaCloudStack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceNetworkInterfaceAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceApsaraStackNetworkInterfaceAttachmentCreate,
		Read:   resourceApsaraStackNetworkInterfaceAttachmentRead,
		Delete: resourceApsaraStackNetworkInterfaceAttachmentDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

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
}

func resourceApsaraStackNetworkInterfaceAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	ecsService := EcsService{client}

	eniId := d.Get("network_interface_id").(string)
	instanceId := d.Get("instance_id").(string)

	request := ecs.CreateAttachNetworkInterfaceRequest()
	request.RegionId = client.RegionId
	if strings.ToLower(client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "ecs", "Department": client.Department, "ResourceGroup": client.ResourceGroup}
	request.InstanceId = instanceId
	request.NetworkInterfaceId = eniId

	err := resource.Retry(DefaultTimeout*time.Minute, func() *resource.RetryError {
		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.AttachNetworkInterface(request)
		})
		if err != nil {
			if IsExpectedErrors(err, NetworkInterfaceInvalidOperations) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "apsarastack_network_interface_attachment", request.GetActionName(), ApsaraStackSdkGoERROR)
	}
	d.SetId(eniId + COLON_SEPARATED + instanceId)
	if err = ecsService.WaitForNetworkInterface(eniId, InUse, DefaultTimeout); err != nil {
		return WrapError(err)
	}
	return resourceApsaraStackNetworkInterfaceAttachmentRead(d, meta)
}

func resourceApsaraStackNetworkInterfaceAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	wiatSecondsIfWithTest(1)
	client := meta.(*connectivity.ApsaraStackClient)
	ecsService := EcsService{client}

	object, err := ecsService.DescribeNetworkInterfaceAttachment(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("instance_id", object.InstanceId)
	d.Set("network_interface_id", object.NetworkInterfaceId)

	return nil
}

func resourceApsaraStackNetworkInterfaceAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	ecsService := EcsService{client}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	eniId, instanceId := parts[0], parts[1]

	_, err = ecsService.DescribeNetworkInterfaceAttachment(d.Id())
	if err != nil {
		if NotFoundError(err) {
			return nil
		}
		return WrapError(err)
	}

	request := ecs.CreateDetachNetworkInterfaceRequest()
	request.RegionId = client.RegionId
	if strings.ToLower(client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.InstanceId = instanceId
	request.NetworkInterfaceId = eniId
	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "ecs", "Department": client.Department, "ResourceGroup": client.ResourceGroup}

	err = resource.Retry(DefaultTimeout*time.Second, func() *resource.RetryError {
		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.DetachNetworkInterface(request)
		})
		if err != nil {
			if IsExpectedErrors(err, NetworkInterfaceInvalidOperations) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), ApsaraStackSdkGoERROR)
	}
	return WrapError(ecsService.WaitForNetworkInterface(eniId, Available, DefaultTimeoutMedium))
}
