package alibabacloudstack

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"time"
	"reflect"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackSnapshot() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlibabacloudStackSnapshotCreate,
		Read:   resourceAlibabacloudStackSnapshotRead,
		Update: resourceAlibabacloudStackSnapshotUpdate,
		Delete: resourceAlibabacloudStackSnapshotDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(DefaultTimeout * time.Second),
		},

		Schema: map[string]*schema.Schema{
			"disk_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(2, 128),
				Deprecated:   "Field 'name' is deprecated and will be removed in a future release. Please use 'snapshot_name' instead.",
			},
			"snapshot_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(2, 128),
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"tags": tagsSchema(),
		},
	}
}

func resourceAlibabacloudStackSnapshotCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	request := ecs.CreateCreateSnapshotRequest()
	client.InitRpcRequest(*request.RpcRequest)
	request.DiskId = d.Get("disk_id").(string)
	request.ClientToken = buildClientToken(request.GetActionName())

	if v, err := connectivity.GetResourceData(d, reflect.TypeOf(""), "snapshot_name", "name"); err == nil {
		request.SnapshotName = v.(string)
	} else {
		return err
	}

	if description, ok := d.GetOk("description"); ok {
		request.Description = description.(string)
	}

	raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.CreateSnapshot(request)
	})
	if err != nil {
		errmsg := ""
		if response, ok := raw.(*ecs.CreateSnapshotResponse); ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_snapshot", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response := raw.(*ecs.CreateSnapshotResponse)
	d.SetId(response.SnapshotId)

	ecsService := EcsService{client}

	stateConf := BuildStateConf([]string{}, []string{string(SnapshotCreatingAccomplished)}, d.Timeout(schema.TimeoutCreate), 2*time.Minute,
		ecsService.SnapshotStateRefreshFunc(d.Id(), []string{string(SnapshotCreatingFailed)}))

	if _, err := stateConf.WaitForState(); err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
	}

	return resourceAlibabacloudStackSnapshotUpdate(d, meta)
}

func resourceAlibabacloudStackSnapshotRead(d *schema.ResourceData, meta interface{}) error {
	waitSecondsIfWithTest(1)
	client := meta.(*connectivity.AlibabacloudStackClient)
	ecsService := EcsService{client}
	snapshot, err := ecsService.DescribeSnapshot(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}
	connectivity.SetResourceData(d, snapshot.SnapshotName, "snapshot_name", "name")
	d.Set("disk_id", snapshot.SourceDiskId)
	d.Set("description", snapshot.Description)

	tags, err := ecsService.DescribeTags(d.Id(), TagResourceSnapshot)
	if err != nil && !errmsgs.NotFoundError(err) {
		return errmsgs.WrapError(err)
	}
	if len(tags) > 0 {
		d.Set("tags", tagsToMap(tags))
	}

	return nil
}

func resourceAlibabacloudStackSnapshotUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	if err := setTags(client, TagResourceSnapshot, d); err != nil {
		return errmsgs.WrapError(err)
	}
	return resourceAlibabacloudStackSnapshotRead(d, meta)
}

func resourceAlibabacloudStackSnapshotDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	ecsService := EcsService{client}

	request := ecs.CreateDeleteSnapshotRequest()
	client.InitRpcRequest(*request.RpcRequest)
	request.SnapshotId = d.Id()

	var raw interface{}
	var err error
	err = resource.Retry(DefaultTimeout*time.Second, func() *resource.RetryError {
		raw, err = client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.DeleteSnapshot(request)
		})
		if err != nil {
			if errmsgs.IsExpectedErrors(err, errmsgs.SnapshotInvalidOperations) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}

		return nil
	})
	if err != nil {
		if errmsgs.IsExpectedErrors(err, []string{"InvalidSnapshotId.NotFound"}) {
			return nil
		}
		errmsg := ""
		if response, ok := raw.(*ecs.DeleteSnapshotResponse); ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 0,
		ecsService.SnapshotStateRefreshFunc(d.Id(), []string{string(SnapshotCreatingFailed)}))

	if _, err = stateConf.WaitForState(); err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
	}
	return nil
}
