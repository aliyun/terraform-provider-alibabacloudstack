package alibabacloudstack

import (
	"fmt"
	"log"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackOssBucketQuota() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlibabacloudStackOssBucketQuotaCreate,
		Read:   resourceAlibabacloudStackOssBucketQuotaRead,
		//Update: resourceAlibabacloudStackOssBucketQuotaCreate,
		Delete: resourceAlibabacloudStackOssBucketQuotaDelete,

		Schema: map[string]*schema.Schema{
			"bucket": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"quota": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAlibabacloudStackOssBucketQuotaCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	ossService := OssService{client}
	var requestInfo *oss.Client
	bucketName := d.Get("bucket").(string)
	det, err := ossService.DescribeOssBucket(bucketName)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_oss_bucket", "IsBucketExist", errmsgs.AlibabacloudStackLogGoSdkERROR)
	}
	quota := d.Get("quota").(int)

	if det.BucketInfo.Name == bucketName {
		request := client.NewCommonRequest("POST", "OneRouter", "2018-12-12", "DoOpenApi", "")
		request.QueryParams["AccountInfo"] = "123456"
		request.QueryParams["SignatureVersion"] = "1.0"
		request.QueryParams["OpenApiAction"] = "SetBucketStorageCapacity"
		request.QueryParams["ProductName"] = "oss"
		request.QueryParams["Params"] = fmt.Sprintf("{\"%s\":\"%s\",\"%s\":%d}", "BucketName", bucketName, "StorageCapacity", quota)
		request.QueryParams["Content"] = fmt.Sprintf("%s%d%s", "<BucketUserQos><StorageCapacity>", quota, "</StorageCapacity></BucketUserQos>")

		raw, err := client.WithEcsClient(func(ossClient *ecs.Client) (interface{}, error) {
			return ossClient.ProcessCommonRequest(request)
		})
		log.Printf("Response of SetBucketStorageCapacity: %s", raw)
		log.Printf("Bresponse ossbucket before error")
		if err != nil {
			errmsg := ""
			if bresponse, ok := raw.(*responses.CommonResponse); ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			}
			if ossNotFoundError(err) {
				return errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackLogGoSdkERROR)
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, bucketName, "SetBucketStorageCapacity", errmsgs.AlibabacloudStackLogGoSdkERROR, errmsg)
		}
		log.Printf("Bresponse ossbucket after error")
		addDebug("SetBucketStorageCapacity", raw, requestInfo, request)
		log.Printf("Bresponse ossbucket check")
		bresponse, ok := raw.(*responses.CommonResponse)
		log.Printf("Bresponse ossbucket %s", bresponse)

		if !ok || bresponse.GetHttpStatus() != 200 {
			errmsg := ""
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_oss_bucket", "SetBucketStorageCapacity", errmsgs.AlibabacloudStackLogGoSdkERROR, errmsg)
		}
		log.Printf("Enter for logging")
	}
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_oss_bucket", "Bucket Not Found", errmsgs.AlibabacloudStackLogGoSdkERROR)
	}
	d.SetId(bucketName)

	return resourceAlibabacloudStackOssBucketKmsRead(d, meta)
}

func resourceAlibabacloudStackOssBucketQuotaRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	ossService := OssService{client}
	var requestInfo *oss.Client
	bucketName := d.Get("bucket").(string)
	det, err := ossService.DescribeOssBucket(bucketName)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_oss_bucket", "IsBucketExist", errmsgs.AlibabacloudStackLogGoSdkERROR)
	}
	if det.BucketInfo.Name == bucketName {
		request := client.NewCommonRequest("GET", "OneRouter", "2018-12-12", "DoOpenApi", "")
		request.QueryParams["AccountInfo"] = "123456"
		request.QueryParams["SignatureVersion"] = "1.0"
		request.QueryParams["OpenApiAction"] = "GetBucketStorageCapacity"
		request.QueryParams["ProductName"] = "oss"
		request.QueryParams["Params"] = fmt.Sprintf("{\"%s\":\"%s\"}", "BucketName", bucketName)

		raw, err := client.WithEcsClient(func(ossClient *ecs.Client) (interface{}, error) {
			return ossClient.ProcessCommonRequest(request)
		})
		log.Printf("Response of GetBucketStorageCapacity: %s", raw)
		if err != nil {
			errmsg := ""
			if bresponse, ok := raw.(*responses.CommonResponse); ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			}
			if ossNotFoundError(err) {
				return errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackLogGoSdkERROR)
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, bucketName, "GetBucketStorageCapacity", errmsgs.AlibabacloudStackLogGoSdkERROR, errmsg)
		}
		addDebug("GetBucketStorageCapacity", raw, requestInfo, request)
		log.Printf("Bresponse ossbucket check")
		bresponse, ok := raw.(*responses.CommonResponse)
		log.Printf("Bresponse ossbucket %s", bresponse)

		if !ok || bresponse.GetHttpStatus() != 200 {
			errmsg := ""
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_oss_bucket", "GetBucketStorageCapacity", errmsgs.AlibabacloudStackLogGoSdkERROR, errmsg)
		}
		log.Printf("Enter for logging")
	}
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_oss_bucket", "Bucket Not Found", errmsgs.AlibabacloudStackLogGoSdkERROR)
	}

	return nil
}

func resourceAlibabacloudStackOssBucketQuotaDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	ossService := OssService{client}
	var requestInfo *oss.Client
	bucketName := d.Get("bucket").(string)
	det, err := ossService.DescribeOssBucket(bucketName)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_oss_bucket", "IsBucketExist", errmsgs.AlibabacloudStackLogGoSdkERROR)
	}

	if det.BucketInfo.Name == bucketName {
		request := client.NewCommonRequest("POST", "OneRouter", "2018-12-12", "DoOpenApi", "")
		request.QueryParams["AccountInfo"] = "123456"
		request.QueryParams["SignatureVersion"] = "1.0"
		request.QueryParams["OpenApiAction"] = "SetBucketStorageCapacity"
		request.QueryParams["ProductName"] = "oss"
		request.QueryParams["Params"] = fmt.Sprintf("{\"%s\":\"%s\",\"%s\":%d}", "BucketName", bucketName, "StorageCapacity", -1)
		request.QueryParams["Content"] = fmt.Sprintf("%s%d%s", "<BucketUserQos><StorageCapacity>", -1, "</StorageCapacity></BucketUserQos>")

		raw, err := client.WithEcsClient(func(ossClient *ecs.Client) (interface{}, error) {
			return ossClient.ProcessCommonRequest(request)
		})
		log.Printf("Response of SetBucketStorageCapacity: %s", raw)
		log.Printf("Bresponse ossbucket before error")
		if err != nil {
			errmsg := ""
			if bresponse, ok := raw.(*responses.CommonResponse); ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			}
			if ossNotFoundError(err) {
				return errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackLogGoSdkERROR)
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, bucketName, "SetBucketStorageCapacity", errmsgs.AlibabacloudStackLogGoSdkERROR, errmsg)
		}
		log.Printf("Bresponse ossbucket after error")
		addDebug("SetBucketStorageCapacity", raw, requestInfo, request)
		log.Printf("Bresponse ossbucket check")
		bresponse, ok := raw.(*responses.CommonResponse)
		log.Printf("Bresponse ossbucket %s", bresponse)

		if !ok || bresponse.GetHttpStatus() != 200 {
			errmsg := ""
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_oss_bucket", "SetBucketStorageCapacity", errmsgs.AlibabacloudStackLogGoSdkERROR, errmsg)
		}
		log.Printf("Enter for logging")
	}
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_oss_bucket", "Bucket Not Found", errmsgs.AlibabacloudStackLogGoSdkERROR)
	}
	d.SetId(bucketName)

	return resourceAlibabacloudStackOssBucketKmsRead(d, meta)
}
