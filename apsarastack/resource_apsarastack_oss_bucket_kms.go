package apsarastack

import (
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"log"
	"strings"
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/aliyun/terraform-provider-alibabaCloudStack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceApsaraStackOssBucketKms() *schema.Resource {
	return &schema.Resource{
		Create: resourceApsaraStackOssBucketKmsCreate,
		Read:   resourceApsaraStackOssBucketKmsRead,
		Update: resourceApsaraStackOssBucketKmsCreate,
		Delete: resourceApsaraStackOssBucketKmsDelete,

		Schema: map[string]*schema.Schema{
			"bucket": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"sse_algorithm": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"kms_data_encryption": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"kms_master_key_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"content3": {
				Type:     schema.TypeString,
				Optional: true,
				//ConflictsWith: []string{"source3"},
			},

			"acl3": {
				Type:         schema.TypeString,
				Default:      oss.ACLPrivate,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"private", "public-read", "public-read-write"}, false),
			},

			"content_type3": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceApsaraStackOssBucketKmsCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	ossService := OssService{client}
	var requestInfo *oss.Client
	bucketName := d.Get("bucket").(string)
	det, err := ossService.DescribeOssBucket(bucketName)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "apsarastack_oss_bucket", "IsBucketExist", ApsaraStackOssGoSdk)
	}
	sseAlgorithm := d.Get("sse_algorithm").(string)
	kmsDateEncryption := ""
	kmsMasterKeyID := ""
	if sseAlgorithm == "KMS" {
		kmsDateEncryption = d.Get("kms_data_encryption").(string)
		kmsMasterKeyID = d.Get("kms_master_key_id").(string)
	}

	if det.BucketInfo.Name == bucketName {
		request := requests.NewCommonRequest()
		if client.Config.Insecure {
			request.SetHTTPSInsecure(client.Config.Insecure)
		}
		request.QueryParams = map[string]string{
			"AccessKeySecret":  client.SecretKey,
			"Product":          "OneRouter",
			"Department":       client.Department,
			"ResourceGroup":    client.ResourceGroup,
			"RegionId":         client.RegionId,
			"Action":           "DoOpenApi",
			"AccountInfo":      "123456",
			"Version":          "2018-12-12",
			"SignatureVersion": "1.0",
			"OpenApiAction":    "PutBucketEncryption",
			"ProductName":      "oss",
			"Params":           fmt.Sprintf("{\"%s\":\"%s\"}", "BucketName", bucketName),
			"Content":          fmt.Sprintf("%s%s%s%s%s%s%s", "<ServerSideEncryptionRule><ApplyServerSideEncryptionByDefault><SSEAlgorithm>", sseAlgorithm, "</SSEAlgorithm><KMSDataEncryption>", kmsDateEncryption, "</KMSDataEncryption><KMSMasterKeyID>", kmsMasterKeyID, "</KMSMasterKeyID></ApplyServerSideEncryptionByDefault></ServerSideEncryptionRule>"),
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
		log.Printf("Response of PutBucketEncryption: %s", raw)
		log.Printf("Bresponse ossbucket before error")
		if err != nil {
			if ossNotFoundError(err) {
				return WrapErrorf(err, NotFoundMsg, ApsaraStackOssGoSdk)
			}
			return WrapErrorf(err, DefaultErrorMsg, bucketName, "PutBucketEncryption", ApsaraStackOssGoSdk)
		}
		log.Printf("Bresponse ossbucket after error")
		addDebug("PutBucketEncryption", raw, requestInfo, request)
		log.Printf("Bresponse ossbucket check")
		bresponse, _ := raw.(*responses.CommonResponse)
		log.Printf("Bresponse ossbucket %s", bresponse)

		if bresponse.GetHttpStatus() != 200 {
			return WrapErrorf(err, DefaultErrorMsg, "apsarastack_oss_bucket", "PutBucketEncryption", ApsaraStackOssGoSdk)
		}
		//logging:= make(map[string]interface{})
		log.Printf("Enter for logging")
	}
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "apsarastack_oss_bucket", "Bucket Not Found", ApsaraStackOssGoSdk)
	}
	// Assign the bucket name as the resource ID
	d.SetId(bucketName)

	return resourceApsaraStackOssBucketKmsRead(d, meta)
}

func resourceApsaraStackOssBucketKmsRead(d *schema.ResourceData, meta interface{}) error {
	wiatSecondsIfWithTest(1)
	client := meta.(*connectivity.ApsaraStackClient)
	ossService := OssService{client}
	var requestInfo *oss.Client
	bucketName := d.Get("bucket").(string)
	det, err := ossService.DescribeOssBucket(bucketName)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "apsarastack_oss_bucket", "IsBucketExist", ApsaraStackOssGoSdk)
	}
	if det.BucketInfo.Name == bucketName {
		request := requests.NewCommonRequest()
		if client.Config.Insecure {
			request.SetHTTPSInsecure(client.Config.Insecure)
		}
		request.QueryParams = map[string]string{
			"AccessKeySecret":  client.SecretKey,
			"Product":          "OneRouter",
			"Department":       client.Department,
			"ResourceGroup":    client.ResourceGroup,
			"RegionId":         client.RegionId,
			"Action":           "DoOpenApi",
			"AccountInfo":      "123456",
			"Version":          "2018-12-12",
			"SignatureVersion": "1.0",
			"OpenApiAction":    "GetBucketEncryption",
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
		log.Printf("Response of GetBucketEncryption: %s", raw)
		if err != nil {
			if ossNotFoundError(err) {
				return WrapErrorf(err, NotFoundMsg, ApsaraStackOssGoSdk)
			}
			return WrapErrorf(err, DefaultErrorMsg, bucketName, "GetBucketEncryption", ApsaraStackOssGoSdk)
		}
		addDebug("BucketEncryption", raw, requestInfo, request)
		log.Printf("Bresponse ossbucket check")
		bresponse, _ := raw.(*responses.CommonResponse)
		log.Printf("Bresponse ossbucket %s", bresponse)

		if bresponse.GetHttpStatus() != 200 {
			return WrapErrorf(err, DefaultErrorMsg, "apsarastack_oss_bucket", "GetBucketEncryption", ApsaraStackOssGoSdk)
		}
		//logging:= make(map[string]interface{})
		log.Printf("Enter for logging")
		//d.Set("SSEAlgorithm",bresponse.GetHttpContentString())
	}
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "apsarastack_oss_bucket", "Bucket Not Found", ApsaraStackOssGoSdk)
	}

	return nil
}

func resourceApsaraStackOssBucketKmsDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	ossService := OssService{client}
	var requestInfo *oss.Client
	det, err := ossService.DescribeOssBucket(d.Id())
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "IsBucketExist", ApsaraStackOssGoSdk)
	}
	addDebug("IsBucketExist", det.BucketInfo, requestInfo, map[string]string{"bucketName": d.Id()})
	if det.BucketInfo.Name == "" {
		return nil
	}

	err = resource.Retry(3*time.Second, func() *resource.RetryError {
		request := requests.NewCommonRequest()
		if client.Config.Insecure {
			request.SetHTTPSInsecure(client.Config.Insecure)
		}
		request.QueryParams = map[string]string{

			"AccessKeySecret":  client.SecretKey,
			"Product":          "OneRouter",
			"Department":       client.Department,
			"ResourceGroup":    client.ResourceGroup,
			"RegionId":         client.RegionId,
			"Action":           "DoOpenApi",
			"AccountInfo":      "123456",
			"Version":          "2018-12-12",
			"SignatureVersion": "1.0",
			"OpenApiAction":    "DeleteBucketEncryption",
			"ProductName":      "oss",
			"Params":           fmt.Sprintf("{\"%s\":\"%s\"}", "BucketName", d.Id()),
		}
		request.Method = "DELETE"      // Set request method
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

		_, err := client.WithOssNewClient(func(ossClient *ecs.Client) (interface{}, error) {

			return ossClient.ProcessCommonRequest(request)
		})

		if err != nil {
			if ossNotFoundError(err) {
				return resource.NonRetryableError(err)
			}
			return resource.RetryableError(err)
		}
		return nil
	})
	return nil
	//return WrapError(ossService.WaitForOssBucket(d.Id(), Deleted, DefaultTimeoutMedium))

}
