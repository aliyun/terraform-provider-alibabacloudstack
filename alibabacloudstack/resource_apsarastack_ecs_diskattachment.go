package alibabacloudstack

import (
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackDiskAttachment() *schema.Resource {
	resource := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"disk_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"device_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
		},
	}
	setResourceFunc(resource, resourceAlibabacloudStackDiskAttachmentCreate,
		resourceAlibabacloudStackDiskAttachmentRead, nil, resourceAlibabacloudStackDiskAttachmentDelete)
	return resource
}

func resourceAlibabacloudStackDiskAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	ecsService := EcsService{client}

	diskID := d.Get("disk_id").(string)
	instanceID := d.Get("instance_id").(string)
	oldDisk, err := ecsService.DescribeDisk(diskID)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	request := ecs.CreateAttachDiskRequest()
	client.InitRpcRequest(*request.RpcRequest)
	request.InstanceId = instanceID
	request.DiskId = diskID

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.AttachDisk(request)
		})

		if err != nil {
			if errmsgs.IsExpectedErrors(err, errmsgs.DiskInvalidOperation) {
				return resource.RetryableError(err)
			}
			errmsg := ""
			if raw != nil {
				response, ok := raw.(*ecs.AttachDiskResponse)
				if ok {
					errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
				}
			}
			return resource.NonRetryableError(errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_disk_attachment", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg))
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})
	if err != nil {
		return err
	}
	d.SetId(request.DiskId + ":" + request.InstanceId)

	if err := ecsService.WaitForDiskAttachment(d.Id(), DiskInUse, DefaultTimeout); err != nil {
		return errmsgs.WrapError(err)
	}
	newDisk, err := ecsService.DescribeDisk(diskID)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	if newDisk.DeleteAutoSnapshot != oldDisk.DeleteAutoSnapshot {
		request := ecs.CreateModifyDiskAttributeRequest()
		client.InitRpcRequest(*request.RpcRequest)
		request.DiskId = diskID
		request.DeleteAutoSnapshot = requests.NewBoolean(oldDisk.DeleteAutoSnapshot)
		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.ModifyDiskAttribute(request)
		})
		if err != nil {
			errmsg := ""
			if raw != nil {
				response, ok := raw.(*ecs.ModifyDiskAttributeResponse)
				if ok {
					errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
				}
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	}
	return nil
}

func resourceAlibabacloudStackDiskAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	ecsService := EcsService{client}
	disk, err := ecsService.DescribeDiskAttachment(d.Id())

	if err != nil {
		if errmsgs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}

	d.Set("instance_id", disk.InstanceId)
	d.Set("disk_id", disk.DiskId)

	if strings.HasPrefix(disk.Device, "/dev/x") {
		disk.Device = "/dev/" + disk.Device[len("/dev/x"):]
	}
	d.Set("device_name", disk.Device)

	return nil
}

func resourceAlibabacloudStackDiskAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	ecsService := EcsService{client}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	request := ecs.CreateDetachDiskRequest()
	client.InitRpcRequest(*request.RpcRequest)
	request.InstanceId = parts[1]
	request.DiskId = parts[0]

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.DetachDisk(request)
		})
		if err != nil {
			if errmsgs.IsExpectedErrors(err, errmsgs.DiskInvalidOperation) {
				time.Sleep(3 * time.Second)
				return resource.RetryableError(err)
			}
			errmsg := ""
			if raw != nil {
				response, ok := raw.(*ecs.DetachDiskResponse)
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
	return errmsgs.WrapError(ecsService.WaitForDiskAttachment(d.Id(), Deleted, DefaultTimeout))
}
