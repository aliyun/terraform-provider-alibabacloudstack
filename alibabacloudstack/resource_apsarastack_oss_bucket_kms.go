package alibabacloudstack

import (
	"context"
	"log"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
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

	ossService2 := OssSdkService{client}
	ossClient, err := ossService2.GetBucketClient(bucketName)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	applyDefault := &oss.ApplyServerSideEncryptionByDefault{
		SSEAlgorithm: oss.Ptr(sseAlgorithm),
	}
	if sseAlgorithm == "KMS" {
		applyDefault.KMSMasterKeyID = oss.Ptr(kmsMasterKeyID)
	}
	putResult, err := ossClient.PutBucketEncryption(context.Background(), &oss.PutBucketEncryptionRequest{
		Bucket: oss.Ptr(bucketName),
		ServerSideEncryptionRule: &oss.ServerSideEncryptionRule{
			ApplyServerSideEncryptionByDefault: applyDefault,
		},
	})
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, bucketName, "PutBucketEncryption", errmsgs.AlibabacloudStackLogGoSdkERROR)
	}
	addDebug("PutBucketEncryption", putResult, nil, map[string]string{"bucketName": bucketName})
	log.Printf("Enter for logging")
	d.SetId(bucketName)

	return nil
}

func resourceAlibabacloudStackOssBucketKmsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	ossService := OssService{client}
	bucketName := d.Id()
	_, err := ossService.DescribeOssBucket(bucketName)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_oss_bucket", "IsBucketExist", errmsgs.AlibabacloudStackLogGoSdkERROR)
	}
	ossService2 := OssSdkService{client}
	ossClient, err := ossService2.GetBucketClient(bucketName)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	getResult, err := ossClient.GetBucketEncryption(context.Background(), &oss.GetBucketEncryptionRequest{
		Bucket: oss.Ptr(bucketName),
	})
	log.Printf("Response of GetBucketEncryption: %v", getResult)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, bucketName, "GetBucketEncryption", errmsgs.AlibabacloudStackLogGoSdkERROR)
	}
	addDebug("BucketEncryption", getResult, nil, map[string]string{"bucketName": bucketName})
	log.Printf("Enter for logging")
	var sseAlgorithmVal string
	var kmsMasterKeyIDVal string
	if getResult.ServerSideEncryptionRule != nil && getResult.ServerSideEncryptionRule.ApplyServerSideEncryptionByDefault != nil {
		apply := getResult.ServerSideEncryptionRule.ApplyServerSideEncryptionByDefault
		if apply.SSEAlgorithm != nil {
			sseAlgorithmVal = *apply.SSEAlgorithm
		}
		if apply.KMSMasterKeyID != nil {
			kmsMasterKeyIDVal = *apply.KMSMasterKeyID
		}
	}
	d.Set("bucket", bucketName)
	d.Set("sse_algorithm", sseAlgorithmVal)
	d.Set("kms_master_key_id", kmsMasterKeyIDVal)

	return nil
}

func resourceAlibabacloudStackOssBucketKmsDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	ossService := OssService{client}
	det, err := ossService.DescribeOssBucket(d.Id())
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), "IsBucketExist", errmsgs.AlibabacloudStackLogGoSdkERROR)
	}
	addDebug("IsBucketExist", det, nil, map[string]string{"bucketName": d.Id()})
	if det.Name == "" {
		return nil
	}

	ossService2 := OssSdkService{client}
	ossClient, err := ossService2.GetBucketClient(d.Id())
	if err != nil {
		return errmsgs.WrapError(err)
	}
	delResult, err := ossClient.DeleteBucketEncryption(context.Background(), &oss.DeleteBucketEncryptionRequest{
		Bucket: oss.Ptr(d.Id()),
	})
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), "DeleteBucketEncryption", errmsgs.AlibabacloudStackSdkGoERROR)
	}
	addDebug("DeleteBucketEncryption", delResult, nil, map[string]string{"bucketName": d.Id()})
	return nil
}
