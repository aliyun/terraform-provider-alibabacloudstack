package alibabacloudstack

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
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
	request.RegionId = client.RegionId
	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.QueryParams = map[string]string{ "Product": "ecs", "Department": client.Department, "ResourceGroup": client.ResourceGroup}
	request.DiskId = d.Get("disk_id").(string)
	request.ClientToken = buildClientToken(request.GetActionName())
	if strings.ToLower(client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	if name, ok := d.GetOk("name"); ok {
		request.SnapshotName = name.(string)
	}
	if description, ok := d.GetOk("description"); ok {
		request.Description = description.(string)
	}

	raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.CreateSnapshot(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultDebugMsg, "alibabacloudstack_snapshot", request.GetActionName(), AlibabacloudStackSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response := raw.(*ecs.CreateSnapshotResponse)
	d.SetId(response.SnapshotId)

	ecsService := EcsService{client}

	stateConf := BuildStateConf([]string{}, []string{string(SnapshotCreatingAccomplished)}, d.Timeout(schema.TimeoutCreate), 2*time.Minute,
		ecsService.SnapshotStateRefreshFunc(d.Id(), []string{string(SnapshotCreatingFailed)}))

	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAlibabacloudStackSnapshotUpdate(d, meta)
}

func resourceAlibabacloudStackSnapshotRead(d *schema.ResourceData, meta interface{}) error {
	waitSecondsIfWithTest(1)
	client := meta.(*connectivity.AlibabacloudStackClient)
	ecsService := EcsService{client}
	snapshot, err := ecsService.DescribeSnapshot(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("name", snapshot.SnapshotName)
	d.Set("disk_id", snapshot.SourceDiskId)
	d.Set("description", snapshot.Description)

	tags, err := ecsService.DescribeTags(d.Id(), TagResourceSnapshot)
	if err != nil && !NotFoundError(err) {
		return WrapError(err)
	}
	if len(tags) > 0 {
		d.Set("tags", tagsToMap(tags))
	}

	return nil
}

func resourceAlibabacloudStackSnapshotUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	if err := setTags(client, TagResourceSnapshot, d); err != nil {
		return WrapError(err)
	}
	return resourceAlibabacloudStackSnapshotRead(d, meta)
}

func resourceAlibabacloudStackSnapshotDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	ecsService := EcsService{client}

	request := ecs.CreateDeleteSnapshotRequest()
	request.RegionId = client.RegionId
	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.QueryParams = map[string]string{ "Product": "ecs", "Department": client.Department, "ResourceGroup": client.ResourceGroup}
	request.SnapshotId = d.Id()
	if strings.ToLower(client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	var raw interface{}
	var err error
	err = resource.Retry(DefaultTimeout*time.Second, func() *resource.RetryError {
		raw, err = client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.DeleteSnapshot(request)
		})
		if err != nil {
			if IsExpectedErrors(err, SnapshotInvalidOperations) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}

		return nil
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidSnapshotId.NotFound"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabacloudStackSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 0,
		ecsService.SnapshotStateRefreshFunc(d.Id(), []string{string(SnapshotCreatingFailed)}))

	if _, err = stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil

}
