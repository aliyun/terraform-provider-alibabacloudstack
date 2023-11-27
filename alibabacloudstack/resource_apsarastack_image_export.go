package alibabacloudstack

import (
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackImageExport() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlibabacloudStackImageExportCreate,
		Read:   resourceAlibabacloudStackImageExportRead,
		Delete: resourceAlibabacloudStackImageExportDelete,
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
}

func resourceAlibabacloudStackImageExportCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	ecsService := EcsService{client: client}

	request := ecs.CreateExportImageRequest()
	request.RegionId = client.RegionId
	if strings.ToLower(client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "ecs", "Department": client.Department, "ResourceGroup": client.ResourceGroup}
	request.ImageId = d.Get("image_id").(string)
	request.OSSBucket = d.Get("oss_bucket").(string)
	request.OSSPrefix = d.Get("oss_prefix").(string)
	raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.ExportImage(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alibabacloudstack_image_export", request.GetActionName(), AlibabacloudStackSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	response := raw.(*ecs.ExportImageResponse)
	taskId := response.TaskId
	d.SetId(request.ImageId)
	stateConf := BuildStateConf([]string{"ExportImage", "Waiting", "Processing"}, []string{"Finished"}, d.Timeout(schema.TimeoutCreate), 2*time.Minute, ecsService.TaskStateRefreshFunc(taskId, []string{"CreateFailed", "UnAvailable"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return resourceAlibabacloudStackImageExportRead(d, meta)

}

func resourceAlibabacloudStackImageExportRead(d *schema.ResourceData, meta interface{}) error {
	waitSecondsIfWithTest(1)
	client := meta.(*connectivity.AlibabacloudStackClient)
	ecsService := EcsService{client: client}

	object, err := ecsService.DescribeImageById(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("image_id", object.ImageId)
	return WrapError(err)
}

func resourceAlibabacloudStackImageExportDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	ossService := OssService{client}
	var requestInfo *oss.Client
	raw, err := client.WithOssClient(func(ossClient *oss.Client) (interface{}, error) {
		requestInfo = ossClient
		return ossClient.Bucket(d.Get("oss_bucket").(string))
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Get("oss_bucket").(string), "OSS Bucket", AlibabacloudStackOssGoSdk)
	}
	addDebug("OSS Bucket", raw, requestInfo, map[string]string{"bucketName": d.Get("oss_bucket").(string)})
	bucket, _ := raw.(*oss.Bucket)
	objectName := d.Id() + "_system.raw.tar.gz"
	if d.Get("oss_prefix").(string) != "" {
		objectName = d.Get("oss_prefix").(string) + "_" + objectName
	}

	request := requests.NewCommonRequest()
	if client.Config.Insecure {
		request.SetHTTPSInsecure(client.Config.Insecure)
	}
	request.QueryParams = map[string]string{
		"RegionId":        client.RegionId,
		"AccessKeySecret": client.SecretKey,
		"AccessKeyId":     client.AccessKey,
		"Product":         "OneRouter",
		"Action":          "DoApi",
		"AppAction":       "DeleteObjects",
		"AppName":         "one-console-app-oss",
		"Version":         "2018-12-12",
		"Params":          "{\"region\":\"" + client.RegionId + "\",\"params\":{\"bucketName\":\"" + d.Get("oss_bucket").(string) + "\",\"objects\":[\"" + objectName + "\"]}}",
		"AccountInfo":     "",
	}
	request.Method = "POST"
	request.Product = "OneRouter"
	request.Version = "2018-12-12"
	request.Domain = client.Domain
	if strings.ToLower(client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.ApiName = "DoApi"
	request.RegionId = client.RegionId
	request.Headers = map[string]string{"RegionId": client.RegionId, "x-acs-instanceid": d.Get("oss_bucket").(string)}

	raw, err = client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.ProcessCommonRequest(request)
	})

	if err != nil {
		WrapErrorf(err, DefaultErrorMsg, d.Id(), "DeleteObject", raw)
	}

	addDebug("DeleteObjects", raw, requestInfo, request)

	bresponse, _ := raw.(*responses.CommonResponse)
	if bresponse.GetHttpStatus() != 200 {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "DeleteObject", AlibabacloudStackOssGoSdk)
	}
	addDebug("DeleteObjects", raw, requestInfo, bresponse.GetHttpContentString())

	return WrapError(ossService.WaitForOssBucketObject(bucket, d.Id(), Deleted, DefaultTimeoutMedium))
}
