package alibabacloudstack

import (
	"context"
	"log"
	"strings"

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
			"oss_cluster": {
				Type:     schema.TypeString,
				Optional: true,
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
	id := d.Get("oss_cluster").(string) + ":" + bucketName
	_, err := ossService.DescribeOssBucket(id)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_oss_bucket", "IsBucketExist", errmsgs.AlibabacloudStackOssGoSdk)
	}
	sseAlgorithm := d.Get("sse_algorithm").(string)
	kmsMasterKeyID := ""
	if sseAlgorithm == "KMS" {
		kmsMasterKeyID = d.Get("kms_master_key_id").(string)
	}
	ossClient, err := ossService.GetOssClientForCluster(d.Get("oss_cluster").(string))
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
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, bucketName, "PutBucketEncryption", errmsgs.AlibabacloudStackOssGoSdk)
	}
	addDebug("PutBucketEncryption", putResult, nil, map[string]string{"bucketName": bucketName})
	log.Printf("Enter for logging")
	d.SetId(id)

	return nil
}

func resourceAlibabacloudStackOssBucketKmsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	ossService := OssService{client}
	if !strings.Contains(d.Id(), ":") {
		d.SetId(d.Get("oss_cluster").(string) + ":" + d.Id())
	}
	apply, err := ossService.DescribeOssBucketKms(d.Id())
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_oss_bucket", "DescribeOssBucketKms", errmsgs.AlibabacloudStackOssGoSdk)
	}
	d.Set("bucket", strings.Split(d.Id(), ":")[1])
	d.Set("oss_cluster", strings.Split(d.Id(), ":")[0])
	if apply != nil {
		if apply.SSEAlgorithm != nil {
			d.Set("sse_algorithm", *apply.SSEAlgorithm)
		}
		if apply.KMSMasterKeyID != nil {
			d.Set("kms_master_key_id", *apply.KMSMasterKeyID)
		}
	}

	return nil
}

func resourceAlibabacloudStackOssBucketKmsDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	ossService := OssService{client}
	ossClient, err := ossService.GetOssClientForCluster(d.Get("oss_cluster").(string))
	if err != nil {
		return errmsgs.WrapError(err)
	}
	delResult, err := ossClient.DeleteBucketEncryption(context.Background(), &oss.DeleteBucketEncryptionRequest{
		Bucket: oss.Ptr(d.Get("bucket").(string)),
	})
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), "DeleteBucketEncryption", errmsgs.AlibabacloudStackOssGoSdk)
	}
	addDebug("DeleteBucketEncryption", delResult, nil, map[string]string{"bucketName": d.Id()})
	return nil
}
