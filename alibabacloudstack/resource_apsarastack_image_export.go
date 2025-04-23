package alibabacloudstack

import (
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
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
			"oss_prefix": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
		},
	}
	setResourceFunc(resource, resourceAlibabacloudStackImageExportCreate, resourceAlibabacloudStackImageExportRead, nil, resourceAlibabacloudStackImageExportDelete)
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
	client := meta.(*connectivity.AlibabacloudStackClient)
	ecsService := EcsService{client: client}

	object, err := ecsService.DescribeImageById(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}
	d.Set("image_id", object.ImageId)
	return nil
}

func resourceAlibabacloudStackImageExportDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	ossService := OssService{client: client}
	var requestInfo *oss.Client
	raw, err := client.WithOssDataClient(func(ossClient *oss.Client) (interface{}, error) {
		requestInfo = ossClient
		return ossClient.Bucket(d.Get("oss_bucket").(string))
	})
	if err != nil {
		errmsg := ""
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Get("oss_bucket").(string), "OSS Bucket", errmsgs.AlibabacloudStackLogGoSdkERROR, errmsg)
	}
	addDebug("OSS Bucket", raw, requestInfo, map[string]string{"bucketName": d.Get("oss_bucket").(string)})
	bucket, _ := raw.(*oss.Bucket)
	objectName := d.Id() + "_system.raw.tar.gz"
	if d.Get("oss_prefix").(string) != "" {
		objectName = d.Get("oss_prefix").(string) + "_" + objectName
	}
	err = bucket.DeleteObject(objectName)
	if err != nil {
		if errmsgs.IsExpectedErrors(err, []string{"No Content", "Not Found"}) {
			return nil
		}
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, objectName, "DeleteObject", errmsgs.AlibabacloudStackLogGoSdkERROR)
	}

	return errmsgs.WrapError(ossService.WaitForOssBucketObject(bucket, objectName, Deleted, DefaultTimeoutMedium))
}
