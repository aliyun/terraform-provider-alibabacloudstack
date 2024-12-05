package alibabacloudstack

import (
	"fmt"
	"log"
	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
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
		return WrapErrorf(err, DefaultErrorMsg, "alibabacloudstack_oss_bucket", "IsBucketExist", AlibabacloudStackOssGoSdk)
	}
	quota := d.Get("quota").(int)

	if det.BucketInfo.Name == bucketName {
		request := requests.NewCommonRequest()
		if client.Config.Insecure {
			request.SetHTTPSInsecure(client.Config.Insecure)
		}
		request.QueryParams = map[string]string{
			
			"Product":          "OneRouter",
			"Department":       client.Department,
			"ResourceGroup":    client.ResourceGroup,
			"RegionId":         client.RegionId,
			"Action":           "DoOpenApi",
			"AccountInfo":      "123456",
			"Version":          "2018-12-12",
			"SignatureVersion": "1.0",
			"OpenApiAction":    "SetBucketStorageCapacity",
			"ProductName":      "oss",
			"Params":           fmt.Sprintf("{\"%s\":\"%s\",\"%s\":%d}", "BucketName", bucketName, "StorageCapacity", quota),
			"Content":          fmt.Sprintf("%s%d%s", "<BucketUserQos><StorageCapacity>", quota, "</StorageCapacity></BucketUserQos>"),
		}
		request.Method = "POST"        // Set request method
		request.Product = "OneRouter"  // Specify product
		request.Version = "2018-12-12" // Specify product version
		request.ServiceCode = "OneRouter"
		if strings.ToLower(client.Config.Protocol) == "https" {
			request.Scheme = "https"
		} else {
			request.Scheme = "http"
		} // Set request scheme. Default: http
		request.ApiName = "DoOpenApi"
		request.Headers = map[string]string{"RegionId": client.RegionId}

		raw, err := client.WithEcsClient(func(ossClient *ecs.Client) (interface{}, error) {
			return ossClient.ProcessCommonRequest(request)
		})
		log.Printf("Response of SetBucketStorageCapacity: %s", raw)
		log.Printf("Bresponse ossbucket before error")
		if err != nil {
			if ossNotFoundError(err) {
				return WrapErrorf(err, NotFoundMsg, AlibabacloudStackOssGoSdk)
			}
			return WrapErrorf(err, DefaultErrorMsg, bucketName, "SetBucketStorageCapacity", AlibabacloudStackOssGoSdk)
		}
		log.Printf("Bresponse ossbucket after error")
		addDebug("SetBucketStorageCapacity", raw, requestInfo, request)
		log.Printf("Bresponse ossbucket check")
		bresponse, _ := raw.(*responses.CommonResponse)
		log.Printf("Bresponse ossbucket %s", bresponse)

		if bresponse.GetHttpStatus() != 200 {
			return WrapErrorf(err, DefaultErrorMsg, "alibabacloudstack_oss_bucket", "SetBucketStorageCapacity", AlibabacloudStackOssGoSdk)
		}
		//logging:= make(map[string]interface{})
		log.Printf("Enter for logging")
	}
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alibabacloudstack_oss_bucket", "Bucket Not Found", AlibabacloudStackOssGoSdk)
	}
	// Assign the bucket name as the resource ID
	d.SetId(bucketName)

	return resourceAlibabacloudStackOssBucketKmsRead(d, meta)
}

func resourceAlibabacloudStackOssBucketQuotaRead(d *schema.ResourceData, meta interface{}) error {
	waitSecondsIfWithTest(1)
	client := meta.(*connectivity.AlibabacloudStackClient)
	ossService := OssService{client}
	var requestInfo *oss.Client
	bucketName := d.Get("bucket").(string)
	det, err := ossService.DescribeOssBucket(bucketName)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alibabacloudstack_oss_bucket", "IsBucketExist", AlibabacloudStackOssGoSdk)
	}
	if det.BucketInfo.Name == bucketName {
		request := requests.NewCommonRequest()
		if client.Config.Insecure {
			request.SetHTTPSInsecure(client.Config.Insecure)
		}
		request.QueryParams = map[string]string{
			
			"Product":          "OneRouter",
			"Department":       client.Department,
			"ResourceGroup":    client.ResourceGroup,
			"RegionId":         client.RegionId,
			"Action":           "DoOpenApi",
			"AccountInfo":      "123456",
			"Version":          "2018-12-12",
			"SignatureVersion": "1.0",
			"OpenApiAction":    "GetBucketStorageCapacity",
			"ProductName":      "oss",
			"Params":           fmt.Sprintf("{\"%s\":\"%s\"}", "BucketName", bucketName),
		}
		request.Method = "GET"         // Set request method
		request.Product = "OneRouter"  // Specify product
		request.Version = "2018-12-12" // Specify product version
		request.ServiceCode = "OneRouter"
		if strings.ToLower(client.Config.Protocol) == "https" {
			request.Scheme = "https"
		} else {
			request.Scheme = "http"
		} // Set request scheme. Default: http
		request.ApiName = "DoOpenApi"
		request.Headers = map[string]string{"RegionId": client.RegionId}

		raw, err := client.WithEcsClient(func(ossClient *ecs.Client) (interface{}, error) {
			return ossClient.ProcessCommonRequest(request)
		})
		log.Printf("Response of GetBucketStorageCapacity: %s", raw)
		if err != nil {
			if ossNotFoundError(err) {
				return WrapErrorf(err, NotFoundMsg, AlibabacloudStackOssGoSdk)
			}
			return WrapErrorf(err, DefaultErrorMsg, bucketName, "GetBucketStorageCapacity", AlibabacloudStackOssGoSdk)
		}
		addDebug("GetBucketStorageCapacity", raw, requestInfo, request)
		log.Printf("Bresponse ossbucket check")
		bresponse, _ := raw.(*responses.CommonResponse)
		log.Printf("Bresponse ossbucket %s", bresponse)

		if bresponse.GetHttpStatus() != 200 {
			return WrapErrorf(err, DefaultErrorMsg, "alibabacloudstack_oss_bucket", "GetBucketStorageCapacity", AlibabacloudStackOssGoSdk)
		}
		//logging:= make(map[string]interface{})
		log.Printf("Enter for logging")
		//d.Set("SSEAlgorithm",bresponse.GetHttpContentString())
	}
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alibabacloudstack_oss_bucket", "Bucket Not Found", AlibabacloudStackOssGoSdk)
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
		return WrapErrorf(err, DefaultErrorMsg, "alibabacloudstack_oss_bucket", "IsBucketExist", AlibabacloudStackOssGoSdk)
	}
	//quota := d.Get("quota").(int)

	if det.BucketInfo.Name == bucketName {
		request := requests.NewCommonRequest()
		if client.Config.Insecure {
			request.SetHTTPSInsecure(client.Config.Insecure)
		}
		request.QueryParams = map[string]string{
			
			"Product":          "OneRouter",
			"Department":       client.Department,
			"ResourceGroup":    client.ResourceGroup,
			"RegionId":         client.RegionId,
			"Action":           "DoOpenApi",
			"AccountInfo":      "123456",
			"Version":          "2018-12-12",
			"SignatureVersion": "1.0",
			"OpenApiAction":    "SetBucketStorageCapacity",
			"ProductName":      "oss",
			"Params":           fmt.Sprintf("{\"%s\":\"%s\",\"%s\":%d}", "BucketName", bucketName, "StorageCapacity", -1),
			"Content":          fmt.Sprintf("%s%d%s", "<BucketUserQos><StorageCapacity>", -1, "</StorageCapacity></BucketUserQos>"),
		}
		request.Method = "POST"        // Set request method
		request.Product = "OneRouter"  // Specify product
		request.Version = "2018-12-12" // Specify product version
		request.ServiceCode = "OneRouter"
		if strings.ToLower(client.Config.Protocol) == "https" {
			request.Scheme = "https"
		} else {
			request.Scheme = "http"
		} // Set request scheme. Default: http
		request.ApiName = "DoOpenApi"
		request.Headers = map[string]string{"RegionId": client.RegionId}

		raw, err := client.WithEcsClient(func(ossClient *ecs.Client) (interface{}, error) {
			return ossClient.ProcessCommonRequest(request)
		})
		log.Printf("Response of SetBucketStorageCapacity: %s", raw)
		log.Printf("Bresponse ossbucket before error")
		if err != nil {
			if ossNotFoundError(err) {
				return WrapErrorf(err, NotFoundMsg, AlibabacloudStackOssGoSdk)
			}
			return WrapErrorf(err, DefaultErrorMsg, bucketName, "SetBucketStorageCapacity", AlibabacloudStackOssGoSdk)
		}
		log.Printf("Bresponse ossbucket after error")
		addDebug("SetBucketStorageCapacity", raw, requestInfo, request)
		log.Printf("Bresponse ossbucket check")
		bresponse, _ := raw.(*responses.CommonResponse)
		log.Printf("Bresponse ossbucket %s", bresponse)

		if bresponse.GetHttpStatus() != 200 {
			return WrapErrorf(err, DefaultErrorMsg, "alibabacloudstack_oss_bucket", "SetBucketStorageCapacity", AlibabacloudStackOssGoSdk)
		}
		//logging:= make(map[string]interface{})
		log.Printf("Enter for logging")
	}
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alibabacloudstack_oss_bucket", "Bucket Not Found", AlibabacloudStackOssGoSdk)
	}
	// Assign the bucket name as the resource ID
	d.SetId(bucketName)

	return resourceAlibabacloudStackOssBucketKmsRead(d, meta)
}
