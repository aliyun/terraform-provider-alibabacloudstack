package alibabacloudstack

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/signer"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlibabacloudStackOssBucket() *schema.Resource {
	resource := &schema.Resource{
		SchemaVersion: 1,
		Schema: map[string]*schema.Schema{
			"bucket": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(3, 63),
				Default:      resource.PrefixedUniqueId("tf-oss-bucket-"),
			},
			"acl": {
				Type:         schema.TypeString,
				Default:      "private",
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"private", "public-read", "public-read-write"}, false),
			},
			"oss_cluster": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"logging": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"target_bucket": {
							Type:     schema.TypeString,
							Required: true,
						},
						"target_prefix": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
				MaxItems: 1,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					// Treat undefined and empty list as equivalent (no logging).
					// Suppress diff when both old and new logging are effectively empty.
					oldRaw, newRaw := d.GetChange("logging")
					oldList := oldRaw.([]interface{})
					newList := newRaw.([]interface{})
					oldEmpty := len(oldList) == 0 || (len(oldList) == 1 && oldList[0].(map[string]interface{})["target_bucket"] == "")
					newEmpty := len(newList) == 0 || (len(newList) == 1 && newList[0].(map[string]interface{})["target_bucket"] == "")
					return oldEmpty && newEmpty || old == new
				},
				DiffSuppressOnRefresh: true,
			},
			"creation_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"extranet_endpoint": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"intranet_endpoint": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"location": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"owner": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"storage_class": {
				Type:     schema.TypeString,
				Default:  "Standard",
				Optional: true,
				ForceNew: true,
			},
			"vpclist": {
				Type:       schema.TypeSet,
				Optional:   true,
				Elem:       &schema.Schema{Type: schema.TypeString},
				Deprecated: "`Vpclist` is not available in the latest versions, and is scheduled for removal in version 3.21.0",
			},
			"bucket_sync": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"storage_capacity": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      1024,
				ValidateFunc: validation.IntBetween(1, 2048000000),
			},
			"sse_algorithm": {
				Type:         schema.TypeString,
				Default:      "",
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"", "AES256", "SM4", "KMS"}, false),
			},
			"kms_key_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"dual_kms_key": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"dual_sync_role": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tags": tagsSchema(),
		},
		// Use CustomizeDiff to add conditional validation
		CustomizeDiff: func(ctx context.Context, diff *schema.ResourceDiff, v interface{}) error {
			sseAlgorithm := diff.Get("sse_algorithm").(string)
			kmsID := diff.Get("kms_key_id").(string)

			if sseAlgorithm == "KMS" && kmsID == "" {
				return fmt.Errorf("kms_key_id must be set when sse_algorithm is KMS")
			}

			return nil
		},
	}
	setResourceFunc(resource, resourceAlibabacloudStackOssBucketCreate, resourceAlibabacloudStackOssBucketRead, resourceAlibabacloudStackOssBucketUpdate, resourceAlibabacloudStackOssBucketDelete)
	return resource
}

func resourceAlibabacloudStackOssBucketCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	ossService := OssService{client}

	ossCluster := d.Get("oss_cluster").(string)
	if endpoints, err := ossService.GetBucketEndpointMap(); err != nil {
		return err
	} else {
		if ossCluster != "" {
			if _, ok := endpoints[ossCluster]; !ok {
				return errmsgs.GetNotFoundErrorFromString(fmt.Sprintf("OssEndpoint not found for cluster: %s", ossCluster))
			}

		} else {
			if len(endpoints) > 1 {
				return errmsgs.Error("The OssCluster in the current region is greater than 1, the `oss_cluster` attribute must be set.")
			}
		}
	}
	ossclient, err := ossService.GetOssClientForCluster(ossCluster)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	bucketName := d.Get("bucket").(string)
	id := ossCluster + ":" + bucketName
	det, err := ossService.DescribeOssBucket(id)

	log.Printf("======================== det:%#v", det)
	if err != nil && !errmsgs.NotFoundError(err) {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_oss_bucket", "DescribeOssBucket", errmsgs.AlibabacloudStackOssGoSdk)
	}
	acl := d.Get("acl").(string)
	storageClass := d.Get("storage_class").(string)
	if storageClass == "" {
		storageClass = "Standard"
	}
	// storage_capacity := d.Get("storage_capacity").(int)
	// If not present, Create Bucket
	if det == nil || *det.Name == "" {
		// Use OSS SDK low-level API (InvokeOperation) to support dualClusterEnabled parameter
		bucketSync := false
		// Do not enable bucketSync during creation, preparing for enabling Role and KMS key later

		// Build XML body for PutBucket request
		xmlBody := fmt.Sprintf(`<CreateBucketConfiguration><StorageClass>%s</StorageClass><DataRedundancyType>LRS</DataRedundancyType><DualClusterEnabled>%t</DualClusterEnabled></CreateBucketConfiguration>`, storageClass, bucketSync)

		input := &oss.OperationInput{
			OpName: "PutBucket",
			Method: "PUT",
			Bucket: oss.Ptr(bucketName),
			Headers: map[string]string{
				"Content-Type": "application/xml",
				"x-oss-acl":    acl,
			},
			Body: strings.NewReader(xmlBody),
		}

		output, err := ossclient.InvokeOperation(context.TODO(), input)
		if err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, bucketName, "PutBucket", errmsgs.AlibabacloudStackOssGoSdk)
		}
		if output.Body != nil {
			output.Body.Close()
		}
		addDebug("PutBucket", output, input, map[string]interface{}{"bucket": bucketName, "acl": acl, "storage_class": storageClass, "bucket_sync": bucketSync})

		err = resource.Retry(3*time.Minute, func() *resource.RetryError {
			det, err := ossService.DescribeOssBucket(id)
			if err != nil {
				return resource.NonRetryableError(err)
			}
			if *det.Name == "" {
				return resource.RetryableError(errmsgs.Error("Trying to ensure new OSS bucket %#v has been created successfully.", bucketName))
			}
			return nil
		})
		if err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_oss_bucket", "Bucket Not Found", errmsgs.AlibabacloudStackOssGoSdk)
		}

	} else {
		return fmt.Errorf("OSS bucket %#v already exists.", bucketName)
	}
	d.SetId(ossCluster + ":" + bucketName)
	ascmService := AscmService{client}
	err = resource.Retry(8*time.Minute, func() *resource.RetryError {
		err = ascmService.ReBindResourceGroup("oss_instance", bucketName)
		if err != nil {
			// Retry on temporary errors
			return resource.RetryableError(err)
		}
		return nil
	})
	if err != nil {
		// ReBindResourceGroup failed after retries. The bucket has already been created.
		// To avoid resource leakage, delete the created bucket before returning the error.
		log.Printf("[WARN] ReBindResourceGroup failed for bucket %s, cleaning up created bucket to avoid resource leakage", bucketName)
		if deleteErr := ossService.DeleteBucket(ossCluster, bucketName); deleteErr != nil {
			return errmsgs.WrapError(fmt.Errorf("failed to move resource group and cleanup bucket: %v; cleanup error: %v", err, deleteErr))
		}
		if waitErr := ossService.WaitForOssBucket(ossCluster, bucketName, Deleted, DefaultTimeoutMedium); waitErr != nil {
			return errmsgs.WrapError(fmt.Errorf("failed to move resource group and cleanup bucket: %v; wait for deletion error: %v", err, waitErr))
		}
		log.Printf("[INFO] Successfully cleaned up bucket %s after ReBindResourceGroup failure", bucketName)
		return errmsgs.WrapError(err)
	}
	tags := d.Get("tags").(map[string]interface{})
	if len(tags) > 0 {
		var tag_objs []OssTags
		for k, v := range tags {
			tag_objs = append(tag_objs, OssTags{
				Key:   k,
				Value: v.(string),
			})
		}
		err = ossService.PutOssBucketTags(ossCluster, bucketName, tag_objs)
		if err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, bucketName, "PutBucketTags", errmsgs.AlibabacloudStackOssGoSdk, err.Error()) // nolint
		}
	}
	return nil
}

func resourceAlibabacloudStackOssBucketRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	ossService := OssService{client}
	if !strings.Contains(d.Id(), ":") {
		d.SetId(d.Get("oss_cluster").(string) + ":" + d.Id())
	}
	ossCluster := strings.Split(d.Id(), ":")[0]
	bucketName := strings.Split(d.Id(), ":")[1]
	object, err := ossService.DescribeOssBucket(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}
	bucketClient, err := ossService.GetOssClientForCluster(ossCluster)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	logging, err := ossService.DescribeOssBucketLogging(ossCluster, bucketName)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, bucketName, "GetBucketLogging", errmsgs.AlibabacloudStackOssGoSdk)
	}
	log.Printf("read describe logging %v", logging)
	d.Set("bucket", bucketName)
	d.Set("oss_cluster", ossCluster)
	if *object.Name == "" {
		log.Print("read: BucketInfo fail!!!!!!")
	}
	log.Printf("========================================================%#v", object.CreationDate)
	log.Printf("========================================================%s", object.CreationDate.Format("2006-01-02T15:04:05.000Z"))

	d.Set("creation_date", object.CreationDate.Format("2006-01-02T15:04:05.000Z"))
	d.Set("extranet_endpoint", *object.ExtranetEndpoint)
	d.Set("intranet_endpoint", *object.IntranetEndpoint)
	d.Set("location", *object.Location)
	d.Set("storage_class", *object.StorageClass)
	var list []map[string]interface{}
	if logging != nil && logging.BucketLoggingStatus != nil && logging.BucketLoggingStatus.LoggingEnabled != nil {
		desclog := logging.BucketLoggingStatus.LoggingEnabled
		var targetBucket, targetPrefix string
		if desclog.TargetBucket != nil {
			targetBucket = *desclog.TargetBucket
		}
		if desclog.TargetPrefix != nil {
			targetPrefix = *desclog.TargetPrefix
		}
		list = append(list, map[string]interface{}{"target_bucket": targetBucket, "target_prefix": targetPrefix})
	}
	if err = d.Set("logging", list); err != nil {
		return errmsgs.WrapError(err)
	}

	bucketSync, err := ossService.GetBucketSync(ossCluster, bucketName)
	d.Set("bucket_sync", false)
	if err != nil {
		if !errmsgs.NotFoundError(err) {
			return errmsgs.WrapError(err)
		}
	} else {
		for _, rule := range bucketSync.Data.ReplicationConfiguration.Rule {
			if rule.Status == "doing" && rule.SrcLocation == "" {
				// Disaster recovery relationships appear in pairs
				d.Set("bucket_sync", true)
				d.Set("dual_sync_role", rule.SyncRole)
				if rule.EncryptionConfiguration.ReplicaKmsKeyID != "" {
					d.Set("dual_kms_key", rule.EncryptionConfiguration.ReplicaKmsKeyID)
				}
				break
			}
		}
	}
	aclResult, err := bucketClient.GetBucketAcl(context.TODO(), &oss.GetBucketAclRequest{
		Bucket: &bucketName,
	})
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, bucketName, "GetBucketAcl", errmsgs.AlibabacloudStackOssGoSdk)
	}
	if aclResult.ACL != nil {
		d.Set("acl", string(*aclResult.ACL))
	}

	scInput := &oss.OperationInput{
		OpName:     "GetBucketStorageCapacity",
		Method:     "GET",
		Bucket:     oss.Ptr(bucketName),
		Parameters: map[string]string{"qos": ""},
	}
	scInput.OpMetadata.Set(signer.SubResource, []string{"qos"})
	scOutput, err := bucketClient.InvokeOperation(context.TODO(), scInput)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, bucketName, "GetBucketStorageCapacity", errmsgs.AlibabacloudStackOssGoSdk)
	}
	scBody, err := io.ReadAll(scOutput.Body)
	if err != nil {
		return errmsgs.WrapErrorf(err, "Read GetBucketStorageCapacity response body failed")
	}
	type BucketUserQosXML struct {
		StorageCapacity string `xml:"StorageCapacity"`
	}
	var qos BucketUserQosXML
	if xmlErr := xml.Unmarshal(scBody, &qos); xmlErr != nil {
		return errmsgs.WrapErrorf(xmlErr, "Parse GetBucketStorageCapacity XML failed")
	}
	if v, convErr := strconv.Atoi(qos.StorageCapacity); convErr == nil {
		d.Set("storage_capacity", v)
	} else {
		return errmsgs.WrapErrorf(convErr, "Get storage capacity failed")
	}

	// Get encryption information
	encResult, err := bucketClient.GetBucketEncryption(context.TODO(), &oss.GetBucketEncryptionRequest{
		Bucket: &bucketName,
	})
	if err != nil {
		// NoSuchServerSideEncryptionRule means no encryption configured
		if ossNotFoundError(err) {
			d.Set("sse_algorithm", "")
		} else {
			return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, bucketName, "GetBucketEncryption", errmsgs.AlibabacloudStackOssGoSdk)
		}
	} else if encResult.ServerSideEncryptionRule != nil && encResult.ServerSideEncryptionRule.ApplyServerSideEncryptionByDefault != nil {
		apply := encResult.ServerSideEncryptionRule.ApplyServerSideEncryptionByDefault
		if apply.SSEAlgorithm != nil {
			d.Set("sse_algorithm", *apply.SSEAlgorithm)
			if *apply.SSEAlgorithm == "KMS" && apply.KMSMasterKeyID != nil {
				d.Set("kms_key_id", *apply.KMSMasterKeyID)
			}
		}
	} else {
		d.Set("sse_algorithm", "")
	}
	tags_map := make(map[string]string)
	tags, _ := ossService.GetBucketTags(ossCluster, bucketName)
	if len(tags) > 0 {
		for _, tag := range tags {
			tagmap := tag.(map[string]interface{})
			if !ossService.ossTagIgnored(tagmap) {
				tags_map[tagmap["Key"].(string)] = tagmap["Value"].(string)
			}
		}
	}
	d.Set("tags", tags_map)
	return nil
}

func resourceAlibabacloudStackOssBucketUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	ossService := OssService{client}
	bucketName := d.Get("bucket").(string)
	ossCluster := d.Get("oss_cluster").(string)
	bucketClient, err := ossService.GetOssClientForCluster(ossCluster)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	if (d.IsNewResource() && d.Get("bucket_sync").(bool)) || (!d.IsNewResource() && d.HasChange("bucket_sync")) {
		target := ""
		process := "closing"
		failed := "starting"
		var syncInput *oss.OperationInput
		if v := d.Get("bucket_sync").(bool); v {
			dual_kms_key := d.Get("dual_kms_key").(string)
			dual_sync_role := d.Get("dual_sync_role").(string)
			var syncXmlBody string
			if dual_kms_key != "" && dual_sync_role != "" {
				content := `<ReplicationConfiguration><Rule><SyncRole>%s</SyncRole><SourceSelectionCriteria><SseKmsEncryptedObjects><Status>Enabled</Status></SseKmsEncryptedObjects></SourceSelectionCriteria><EncryptionConfiguration><ReplicaKmsKeyID>%s</ReplicaKmsKeyID></EncryptionConfiguration></Rule></ReplicationConfiguration>`
				syncXmlBody = fmt.Sprintf(content, dual_sync_role, dual_kms_key)
			} else if (dual_kms_key != "" && dual_sync_role == "") || (dual_kms_key == "" && dual_sync_role != "") {
				return fmt.Errorf("dual_kms_key and dual_sync_role must be set at the same time")
			}
			syncInput = &oss.OperationInput{
				OpName:     "PutBucketSync",
				Method:     "PUT",
				Bucket:     oss.Ptr(bucketName),
				Parameters: map[string]string{"syncinternal": ""},
				Headers:    map[string]string{"Content-Type": "application/xml"},
				Body:       strings.NewReader(syncXmlBody),
			}
			if syncXmlBody != "" {
				syncInput.Body = strings.NewReader(syncXmlBody)
			}
			target = "doing"
			process = "starting"
			failed = "closing"
		} else {
			syncInput = &oss.OperationInput{
				OpName:     "DeleteBucketSync",
				Method:     "DELETE",
				Bucket:     oss.Ptr(bucketName),
				Parameters: map[string]string{"syncinternal": ""},
			}
		}
		// syncInput.OpMetadata.Set(signer.SubResource, []string{"syncinternal"})
		_, err = bucketClient.InvokeOperation(context.TODO(), syncInput)
		if err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, bucketName, syncInput.OpName, errmsgs.AlibabacloudStackOssGoSdk)
		}
		stateConf := BuildStateConf([]string{process}, []string{target}, d.Timeout(schema.TimeoutCreate), 2*time.Minute, ossService.OssBucketSyncStateRefreshFunc(ossCluster, bucketName, []string{failed}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
		}
	}

	if d.HasChange("storage_capacity") {
		storageCapacity := d.Get("storage_capacity").(int)
		xmlBody := fmt.Sprintf("<BucketUserQos><StorageCapacity>%d</StorageCapacity></BucketUserQos>", storageCapacity)
		scInput2 := &oss.OperationInput{
			OpName:     "SetBucketStorageCapacity",
			Method:     "PUT",
			Bucket:     oss.Ptr(bucketName),
			Parameters: map[string]string{"qos": ""},
			Headers:    map[string]string{"Content-Type": "application/xml"},
			Body:       strings.NewReader(xmlBody),
		}
		scInput2.OpMetadata.Set(signer.SubResource, []string{"qos"})
		_, err = bucketClient.InvokeOperation(context.TODO(), scInput2)
		if err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, bucketName, "SetBucketStorageCapacity", errmsgs.AlibabacloudStackOssGoSdk)
		}
	}

	if d.HasChanges("sse_algorithm", "kms_key_id") {
		if d.Get("sse_algorithm").(string) == "" {
			_, err := bucketClient.DeleteBucketEncryption(context.TODO(), &oss.DeleteBucketEncryptionRequest{
				Bucket: &bucketName,
			})
			if err != nil {
				return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, bucketName, "DeleteBucketEncryption", errmsgs.AlibabacloudStackOssGoSdk)
			}
		} else {
			sse_algorithm := d.Get("sse_algorithm").(string)
			kms_key_id := d.Get("kms_key_id").(string)
			applyDefault := &oss.ApplyServerSideEncryptionByDefault{
				SSEAlgorithm: &sse_algorithm,
			}
			if sse_algorithm == "KMS" {
				applyDefault.KMSMasterKeyID = &kms_key_id
			}
			_, err := bucketClient.PutBucketEncryption(context.TODO(), &oss.PutBucketEncryptionRequest{
				Bucket: &bucketName,
				ServerSideEncryptionRule: &oss.ServerSideEncryptionRule{
					ApplyServerSideEncryptionByDefault: applyDefault,
				},
			})
			if err != nil {
				return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, bucketName, "PutBucketEncryption", errmsgs.AlibabacloudStackOssGoSdk)
			}
		}
	}

	needLoggingUpdate := false
	if _, ok := d.GetOk("logging"); ok && d.IsNewResource() {
		// For new resources, compare existing server-side logging with desired config.
		existingLogging, err := ossService.DescribeOssBucketLogging(ossCluster, bucketName)
		if err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, bucketName, "GetBucketLogging", errmsgs.AlibabacloudStackOssGoSdk)
		}
		desiredLogging := d.Get("logging").([]interface{})
		needLoggingUpdate = !isLoggingEqual(existingLogging, desiredLogging)
	} else {
		// For existing resources, just check if logging config changed.
		needLoggingUpdate = d.HasChange("logging")
	}
	if needLoggingUpdate {
		err := resourceAlibabacloudStackOssBucketLoggingUpdate(client, d)
		if err != nil {
			return errmsgs.WrapError(err)
		}
	}
	// if d.HasChange("vpclist") {
	// 	vpc_err := checkVpcListChange(d, meta)
	// 	if vpc_err != nil {
	// 		return errmsgs.WrapError(vpc_err)
	// 	}
	// }
	if d.HasChange("acl") {
		acl := d.Get("acl").(string)
		_, err = bucketClient.PutBucketAcl(context.TODO(), &oss.PutBucketAclRequest{
			Bucket: &bucketName,
			Acl:    oss.BucketACLType(acl),
		})
		if err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, bucketName, "PutBucketAcl", errmsgs.AlibabacloudStackOssGoSdk)
		}
	}
	if d.IsNewResource() {
		return nil
	}
	if d.HasChange("tags") {
		tags := d.Get("tags").(map[string]interface{})
		var tag_objs []OssTags
		for k, v := range tags {
			tag_objs = append(tag_objs, OssTags{
				Key:   k,
				Value: v.(string),
			})
		}
		if len(tag_objs) <= 0 {
			err = ossService.DeleteBucketTags(ossCluster, bucketName)
		} else {
			var tag_objs []OssTags
			for k, v := range tags {
				tag_objs = append(tag_objs, OssTags{
					Key:   k,
					Value: v.(string),
				})
			}
			err = ossService.PutOssBucketTags(ossCluster, bucketName, tag_objs)
		}
		if err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, bucketName, "PutBucketTags", errmsgs.AlibabacloudStackOssGoSdk, err.Error()) // nolint
		}
	}
	return nil
}

func resourceAlibabacloudStackOssBucketDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	ossService := OssService{client}
	ossCluster := d.Get("oss_cluster").(string)
	bucketName := d.Get("bucket").(string)
	det, err := ossService.DescribeOssBucket(d.Id())
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), "IsBucketExist", errmsgs.AlibabacloudStackOssGoSdk)
	}
	addDebug("IsBucketExist", det, nil, map[string]string{"bucketName": bucketName})
	if *det.Name == "" {
		return nil
	}

	if err := ossService.DeleteBucket(ossCluster, bucketName); err != nil {
		return errmsgs.WrapError(err)
	}
	return errmsgs.WrapError(ossService.WaitForOssBucket(ossCluster, bucketName, Deleted, DefaultTimeoutMedium))
}

func resourceAlibabacloudStackOssBucketLoggingUpdate(client *connectivity.AlibabacloudStackClient, d *schema.ResourceData) error {
	ossService := OssService{client}
	ossCluster := d.Get("oss_cluster").(string)
	bucketName := d.Get("bucket").(string)
	ossClient, err := ossService.GetOssClientForCluster(ossCluster)
	if err != nil {
		return errmsgs.WrapError(err)
	}

	loggingList := d.Get("logging").([]interface{})

	// If logging is empty or not set, delete existing logging configuration.
	if len(loggingList) == 0 {
		log.Printf("[DEBUG] Deleting bucket logging for %s", bucketName)
		_, err = ossClient.DeleteBucketLogging(context.TODO(), &oss.DeleteBucketLoggingRequest{
			Bucket: &bucketName,
		})
		if err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, bucketName, "DeleteBucketLogging", errmsgs.AlibabacloudStackOssGoSdk)
		}
		return nil
	}

	// Put logging configuration.
	logging := loggingList[0].(map[string]interface{})
	targetBucket := fmt.Sprint(logging["target_bucket"])
	targetPrefix := fmt.Sprint(logging["target_prefix"])

	// Verify the target bucket exists.
	_, err = ossService.DescribeOssBucket(d.Id())
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_oss_bucket", "DescribeBucket")
	}

	log.Printf("[DEBUG] Putting bucket logging for %s, target_bucket: %s, target_prefix: %s", bucketName, targetBucket, targetPrefix)
	_, err = ossClient.PutBucketLogging(context.TODO(), &oss.PutBucketLoggingRequest{
		Bucket: &bucketName,
		BucketLoggingStatus: &oss.BucketLoggingStatus{
			LoggingEnabled: &oss.LoggingEnabled{
				TargetBucket: &targetBucket,
				TargetPrefix: &targetPrefix,
			},
		},
	})
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, bucketName, "PutBucketLogging", errmsgs.AlibabacloudStackOssGoSdk)
	}

	return nil
}

// isLoggingEqual compares server-side logging config with Terraform desired config.
func isLoggingEqual(existing *oss.GetBucketLoggingResult, desired []interface{}) bool {
	// Both empty means equal.
	hasExisting := existing != nil && existing.BucketLoggingStatus != nil && existing.BucketLoggingStatus.LoggingEnabled != nil
	if !hasExisting && len(desired) == 0 {
		return true
	}
	// One has logging, the other doesn't.
	if !hasExisting || len(desired) == 0 {
		return false
	}

	existingLogging := existing.BucketLoggingStatus.LoggingEnabled
	desiredMap := desired[0].(map[string]interface{})

	existingBucket := ""
	if existingLogging.TargetBucket != nil {
		existingBucket = *existingLogging.TargetBucket
	}
	existingPrefix := ""
	if existingLogging.TargetPrefix != nil {
		existingPrefix = *existingLogging.TargetPrefix
	}

	return existingBucket == fmt.Sprint(desiredMap["target_bucket"]) &&
		existingPrefix == fmt.Sprint(desiredMap["target_prefix"])
}

type OssTags struct {
	Value string `name:"Value"`
	Key   string `name:"Key"`
}
