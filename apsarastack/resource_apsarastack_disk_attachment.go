package apsarastack

import (
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-alibabacloudstack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceApsaraStackDiskAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceApsaraStackDiskAttachmentCreate,
		Read:   resourceApsaraStackDiskAttachmentRead,
		Delete: resourceApsaraStackDiskAttachmentDelete,

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
}

func resourceApsaraStackDiskAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	ecsService := EcsService{client}

	diskID := d.Get("disk_id").(string)
	instanceID := d.Get("instance_id").(string)
	oldDisk, err := ecsService.DescribeDisk(diskID)
	if err != nil {
		return WrapError(err)
	}
	request := ecs.CreateAttachDiskRequest()
	request.RegionId = client.RegionId
	request.InstanceId = instanceID
	if strings.ToLower(client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "ecs", "Department": client.Department, "ResourceGroup": client.ResourceGroup}
	request.DiskId = diskID

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.AttachDisk(request)
		})

		if err != nil {
			if IsExpectedErrors(err, DiskInvalidOperation) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "apsarastack_disk_attachment", request.GetActionName(), ApsaraStackSdkGoERROR)
	}
	d.SetId(request.DiskId + ":" + request.InstanceId)

	if err := ecsService.WaitForDiskAttachment(d.Id(), DiskInUse, DefaultTimeout); err != nil {
		return WrapError(err)
	}
	newDisk, err := ecsService.DescribeDisk(diskID)
	if err != nil {
		return WrapError(err)
	}
	if newDisk.DeleteAutoSnapshot != oldDisk.DeleteAutoSnapshot {
		request := ecs.CreateModifyDiskAttributeRequest()
		if strings.ToLower(client.Config.Protocol) == "https" {
			request.Scheme = "https"
		} else {
			request.Scheme = "http"
		}
		request.RegionId = client.RegionId
		request.Headers = map[string]string{"RegionId": client.RegionId}
		request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "ecs", "Department": client.Department, "ResourceGroup": client.ResourceGroup}
		request.DiskId = diskID
		request.DeleteAutoSnapshot = requests.NewBoolean(oldDisk.DeleteAutoSnapshot)
		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.ModifyDiskAttribute(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), ApsaraStackSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	}
	return resourceApsaraStackDiskAttachmentRead(d, meta)
}

func resourceApsaraStackDiskAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	wiatSecondsIfWithTest(1)
	client := meta.(*connectivity.ApsaraStackClient)
	ecsService := EcsService{client}
	disk, err := ecsService.DescribeDiskAttachment(d.Id())

	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("instance_id", disk.InstanceId)
	d.Set("disk_id", disk.DiskId)

	if strings.HasPrefix(disk.Device, "/dev/x") {
		disk.Device = "/dev/" + disk.Device[len("/dev/x"):]
	}
	d.Set("device_name", disk.Device)

	return nil
}

func resourceApsaraStackDiskAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	ecsService := EcsService{client}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	request := ecs.CreateDetachDiskRequest()
	if strings.ToLower(client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.RegionId = client.RegionId
	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "ecs", "Department": client.Department, "ResourceGroup": client.ResourceGroup}
	request.InstanceId = parts[1]
	request.DiskId = parts[0]

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.DetachDisk(request)
		})
		if err != nil {
			if IsExpectedErrors(err, DiskInvalidOperation) {
				time.Sleep(3 * time.Second)
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
	return WrapError(ecsService.WaitForDiskAttachment(d.Id(), Deleted, DefaultTimeout))
}
