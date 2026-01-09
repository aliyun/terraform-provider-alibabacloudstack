package alibabacloudstack

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackOssBucketKms() *schema.Resource {
	resource := &schema.Resource{
		DeprecationMessage: "oss_bucket already includes corresponding functions, and is scheduled for removal in version 3.21.0",
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
	setResourceFunc(resource, resourceAlibabacloudStackOssBucketKmsCreate, resourceAlibabacloudStackOssBucketKmsRead, nil, resourceAlibabacloudStackOssBucketKmsDelete)
	return resource
}

func resourceAlibabacloudStackOssBucketKmsCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	ossService := OssService{client}
	var requestInfo *oss.Client
	bucketName := d.Get("bucket").(string)
	_, err := ossService.DescribeOssBucket(bucketName)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_oss_bucket", "IsBucketExist", errmsgs.AlibabacloudStackLogGoSdkERROR)
	}
	sseAlgorithm := d.Get("sse_algorithm").(string)
	kmsMasterKeyID := ""
	if sseAlgorithm == "KMS" {
		kmsMasterKeyID = d.Get("kms_master_key_id").(string)
	}

	request := client.NewCommonRequest("GET", "OneRouter", "2018-12-12", "DoOpenApi", "")
	mergeMaps(request.QueryParams, map[string]string{
		"OpenApiAction": "PutBucketEncryption",
		"ProductName":   "oss",
		"Params":        fmt.Sprintf("{\"%s\":\"%s\"}", "BucketName", bucketName),
	})
	if sseAlgorithm == "KMS" {
		request.QueryParams["Content"] = fmt.Sprintf("<ServerSideEncryptionRule><ApplyServerSideEncryptionByDefault><SSEAlgorithm>KMS</SSEAlgorithm><KMSMasterKeyID>%s</KMSMasterKeyID></ApplyServerSideEncryptionByDefault></ServerSideEncryptionRule>", kmsMasterKeyID)
	} else {
		request.QueryParams["Content"] = fmt.Sprintf("<ServerSideEncryptionRule><ApplyServerSideEncryptionByDefault><SSEAlgorithm>%s</SSEAlgorithm></ApplyServerSideEncryptionByDefault></ServerSideEncryptionRule>", sseAlgorithm)
	}

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
	d.SetId(bucketName)

	return nil
}

func resourceAlibabacloudStackOssBucketKmsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	ossService := OssService{client}
	var requestInfo *oss.Client
	bucketName := d.Id()
	_, err := ossService.DescribeOssBucket(bucketName)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_oss_bucket", "IsBucketExist", errmsgs.AlibabacloudStackLogGoSdkERROR)
	}
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
	resp := make(map[string]interface{})
	json.Unmarshal(bresponse.GetHttpContentBytes(), &resp)
	encryption_data, err := jsonpath.Get("$.Data.ServerSideEncryptionRule.ApplyServerSideEncryptionByDefault", resp)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_oss_bucket", "Bucket Not Found", errmsgs.AlibabacloudStackLogGoSdkERROR)
	}
	encryption := encryption_data.(map[string]interface{})
	d.Set("bucket", bucketName)
	d.Set("sse_algorithm", encryption["SSEAlgorithm"].(string))
	d.Set("kms_master_key_id", encryption["KMSMasterKeyID"])

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
