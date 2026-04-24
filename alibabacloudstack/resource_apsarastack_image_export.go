package alibabacloudstack

import (
	"context"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackImageExport() *schema.Resource {
	resource := &schema.Resource{
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"image_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"oss_bucket": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"oss_cluster": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"oss_prefix": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"oss_object": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
	setResourceFunc(resource, resourceAlibabacloudStackImageExportCreate, resourceAlibabacloudStackImageExportRead, nil, resourceAlibabacloudStackImageExportDelete)
	resource.Importer = nil
	return resource
}

func resourceAlibabacloudStackImageExportCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	ecsService := EcsService{client: client}

	request := ecs.CreateExportImageRequest()
	client.InitRpcRequest(*request.RpcRequest)
	request.ImageId = d.Get("image_id").(string)
	request.OSSBucket = d.Get("oss_bucket").(string)
	request.OSSPrefix = d.Get("oss_prefix").(string)
	raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.ExportImage(request)
	})
	if err != nil {
		errmsg := ""
		if raw != nil {
			response, ok := raw.(*ecs.ExportImageResponse)
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
			}
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_image_export", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	response := raw.(*ecs.ExportImageResponse)
	taskId := response.TaskId
	d.SetId(request.ImageId)
	stateConf := BuildStateConf([]string{"ExportImage", "Waiting", "Processing"}, []string{"Finished"}, d.Timeout(schema.TimeoutCreate), 2*time.Minute, ecsService.TaskStateRefreshFunc(taskId, []string{"CreateFailed", "UnAvailable"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
	}
	return nil
}

func resourceAlibabacloudStackImageExportRead(d *schema.ResourceData, meta interface{}) error {
	objectName := d.Get("image_id").(string) + "_system.raw.tar.gz"
	if d.Get("oss_prefix").(string) != "" {
		objectName = d.Get("oss_prefix").(string) + "_" + objectName
	}
	d.Set("oss_object", objectName)
	return nil
}

func resourceAlibabacloudStackImageExportDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	ossService := OssService{client: client}
	ossCluster := d.Get("oss_cluster").(string)
	bucket, err := ossService.GetOssClientForCluster(ossCluster)
	if err != nil {
		return err
	}
	objectName := d.Get("oss_object").(string)
	bucketName := d.Get("oss_bucket").(string)
	delReq := &oss.DeleteObjectRequest{
		Bucket: &bucketName,
		Key:    &objectName,
	}
	_, err = bucket.DeleteObject(context.Background(), delReq)
	if err != nil {
		if errmsgs.IsExpectedErrors(err, "No Content", "Not Found") {
			return nil
		}
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, objectName, "DeleteObject", errmsgs.AlibabacloudStackLogGoSdkERROR)
	}
	id := ossCluster + ":" + bucketName + ":" + objectName
	return errmsgs.WrapError(ossService.WaitForOssBucketObject(id, Deleted, DefaultTimeoutMedium))
}
