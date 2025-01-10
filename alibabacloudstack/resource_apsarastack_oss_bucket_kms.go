package alibabacloudstack

import (
	"fmt"
	"log"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackOssBucketKms() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlibabacloudStackOssBucketKmsCreate,
		Read:   resourceAlibabacloudStackOssBucketKmsRead,
		Delete: resourceAlibabacloudStackOssBucketKmsDelete,

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
			// 			"content3": {
			// 				Type:     schema.TypeString,
			// 				Optional: true,
			// 				ConflictsWith: []string{"source3"},
			// 			},
			//
			// 			"acl3": {
			// 				Type:         schema.TypeString,
			// 				Default:      oss.ACLPrivate,
			// 				Optional:     true,
			// 				ValidateFunc: validation.StringInSlice([]string{"private", "public-read", "public-read-write"}, false),
			// 			},
			//
			// 			"content_type3": {
			// 				Type:     schema.TypeString,
			// 				Optional: true,
			// 				Computed: true,
			// 			},
		},
	}
}

func resourceAlibabacloudStackOssBucketKmsCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	ossService := OssService{client}
	var requestInfo *oss.Client
	bucketName := d.Get("bucket").(string)
	det, err := ossService.DescribeOssBucket(bucketName)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_oss_bucket", "IsBucketExist", errmsgs.AlibabacloudStackLogGoSdkERROR)
	}
	sseAlgorithm := d.Get("sse_algorithm").(string)
	kmsDateEncryption := ""
	kmsMasterKeyID := ""
	if sseAlgorithm == "KMS" {
		kmsDateEncryption = d.Get("kms_data_encryption").(string)
		kmsMasterKeyID = d.Get("kms_master_key_id").(string)
	}

	if det.BucketInfo.Name == bucketName {
		request := client.NewCommonRequest("POST", "OneRouter", "2018-12-12", "DoOpenApi", "")
		mergeMaps(request.QueryParams, map[string]string{
			"OpenApiAction": "PutBucketEncryption",
			"ProductName":   "oss",
			"Params":        fmt.Sprintf("{\"%s\":\"%s\"}", "BucketName", bucketName),
			"Content":       fmt.Sprintf("%s%s%s%s%s%s%s", "<ServerSideEncryptionRule><ApplyServerSideEncryptionByDefault><SSEAlgorithm>", sseAlgorithm, "</SSEAlgorithm><KMSDataEncryption>", kmsDateEncryption, "</KMSDataEncryption><KMSMasterKeyID>", kmsMasterKeyID, "</KMSMasterKeyID></ApplyServerSideEncryptionByDefault></ServerSideEncryptionRule>"),
		})

		bresponse, err := client.ProcessCommonRequest(request)
		if err != nil {
			if bresponse == nil {
				return errmsgs.WrapErrorf(err, "Process Common Request Failed")
			}
			if ossNotFoundError(err) {
				return errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackLogGoSdkERROR)
			}
			errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, bucketName, "PutBucketEncryption", errmsgs.AlibabacloudStackLogGoSdkERROR, errmsg)
		}
		addDebug("PutBucketEncryption", bresponse, requestInfo, request)

		if bresponse.GetHttpStatus() != 200 {
			errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_oss_bucket", "PutBucketEncryption", errmsgs.AlibabacloudStackLogGoSdkERROR, errmsg)
		}
		log.Printf("Enter for logging")
	}
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_oss_bucket", "Bucket Not Found", errmsgs.AlibabacloudStackLogGoSdkERROR)
	}
	d.SetId(bucketName)

	return resourceAlibabacloudStackOssBucketKmsRead(d, meta)
}

func resourceAlibabacloudStackOssBucketKmsRead(d *schema.ResourceData, meta interface{}) error {
	waitSecondsIfWithTest(1)
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
		mergeMaps(request.QueryParams, map[string]string{
			"AccountInfo":      "123456",
			"SignatureVersion": "1.0",
			"OpenApiAction":    "GetBucketEncryption",
			"ProductName":      "oss",
			"Params":           fmt.Sprintf("{\"%s\":\"%s\"}", "BucketName", bucketName),
		})

		bresponse, err := client.ProcessCommonRequest(request)
		log.Printf("Response of GetBucketEncryption: %s", bresponse)
		if err != nil {
			if bresponse == nil {
				return errmsgs.WrapErrorf(err, "Process Common Request Failed")
			}
			if ossNotFoundError(err) {
				return errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackLogGoSdkERROR)
			}
			errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, bucketName, "GetBucketEncryption", errmsgs.AlibabacloudStackLogGoSdkERROR, errmsg)
		}
		addDebug("BucketEncryption", bresponse, requestInfo, request)
		log.Printf("Bresponse ossbucket check")
		log.Printf("Bresponse ossbucket %s", bresponse)

		if bresponse.GetHttpStatus() != 200 {
			errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_oss_bucket", "GetBucketEncryption", errmsgs.AlibabacloudStackLogGoSdkERROR, errmsg)
		}
		log.Printf("Enter for logging")
	}
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_oss_bucket", "Bucket Not Found", errmsgs.AlibabacloudStackLogGoSdkERROR)
	}

	return nil
}

func resourceAlibabacloudStackOssBucketKmsDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	ossService := OssService{client}
	var requestInfo *oss.Client
	det, err := ossService.DescribeOssBucket(d.Id())
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), "IsBucketExist", errmsgs.AlibabacloudStackLogGoSdkERROR)
	}
	addDebug("IsBucketExist", det.BucketInfo, requestInfo, map[string]string{"bucketName": d.Id()})
	if det.BucketInfo.Name == "" {
		return nil
	}

	request := client.NewCommonRequest("DELETE", "OneRouter", "2018-12-12", "DoOpenApi", "")
	mergeMaps(request.QueryParams, map[string]string{
		"OpenApiAction": "DeleteBucketEncryption",
		"ProductName":   "oss",
		"Params":        fmt.Sprintf("{\"%s\":\"%s\"}", "BucketName", d.Id()),
	})

	bresponse, err := client.ProcessCommonRequest(request)
	if err != nil {
		if bresponse == nil {
			return errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		if ossNotFoundError(err) {
			return nil
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		err = errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "OssBucketKms", "DeleteBucketEncryption", errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		return err
	}
	return nil
}
