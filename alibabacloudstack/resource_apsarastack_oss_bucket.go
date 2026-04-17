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
					if k == "logging.#" && old == "1" && new == "0" {
						loggings := d.Get("logging").([]interface{})
						logging := loggings[0].(map[string]interface{})
						if logging["target_bucket"] == "" && logging["target_prefix"] == "" {
							return true
						}
					}
					return false
				},
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
	bucketName := d.Get("bucket").(string)
	ossCluster := d.Get("oss_cluster").(string)
	det, err := ossService.DescribeOssBucket(bucketName)

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

		ossclient, err := ossService.GetOssClientForCluster(ossCluster)
		if err != nil {
			return errmsgs.WrapError(err)
		}

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
			det, err := ossService.DescribeOssBucket(bucketName)
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
	d.SetId(bucketName)
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
		err = ossService.PutOssBucketTags(bucketName, tag_objs)
		if err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, bucketName, "PutBucketTags", errmsgs.AlibabacloudStackOssGoSdk, err.Error()) // nolint
		}
	}
	return nil
}

func resourceAlibabacloudStackOssBucketRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	ossService := OssService{client}
	object, err := ossService.DescribeOssBucket(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}
	bucketClient, err := ossService.GetBucketClient(d.Id())
	if err != nil {
		return errmsgs.WrapError(err)
	}
	logging, err := ossService.DescribeOssBucketLogging(d.Id())
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), "GetBucketLogging", errmsgs.AlibabacloudStackOssGoSdk)
	}
	ossEndpointData, err := ossService.GetOssEndpointList()
	if err != nil {
		return errmsgs.WrapError(err)
	}
	log.Printf("read describe logging %v", logging)
	d.Set("bucket", object.Name)
	if *object.Name == "" {
		log.Print("read: BucketInfo fail!!!!!!")
	}
	d.Set("creation_date", object.CreationDate.Format("2006-01-02 15:04:05 +0000 UTC"))
	d.Set("extranet_endpoint", *object.ExtranetEndpoint)
	d.Set("intranet_endpoint", *object.IntranetEndpoint)
	d.Set("location", *object.Location)
	d.Set("storage_class", *object.StorageClass)
	for _, v := range ossEndpointData {
		endpoint := v.(map[string]interface{})
		if endpoint["oss-endpoint"].(string) == *object.IntranetEndpoint {
			d.Set("oss_cluster", endpoint["cluster"])
			break
		}
	}
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

	// policy, err := ossService.DescribeOssBucketPolicy(d.Id())
	// if err != nil {
	// 	log.Printf("==========================================================%#v", policy)
	// }
	// bvclient := meta.(*connectivity.AlibabacloudStackClient)
	// bvserver := BucketVpcService{bvclient}
	// vpclist, binderr := bvserver.BucketVpcList(d.Get("bucket").(string))
	// if binderr != nil {
	// 	return errmsgs.WrapError(binderr)
	// }
	// var vlist []interface{}
	// if len(vpclist.VpcList) > 0 {
	// 	for _, v := range vpclist.VpcList {
	// 		vpc := v.(map[string]interface{})
	// 		vlist = append(vlist, vpc["vpcId"].(string))
	// 	}
	// }
	// d.Set("vpclist", schema.NewSet(schema.HashString, vlist))

	bucketSync, err := ossService.GetBucketSync(d.Id())
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
	bucketName := d.Id()
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
	tags, _ := ossService.GetBucketTags(bucketName)
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
	bucketClient, err := ossService.GetBucketClient(bucketName)
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
		stateConf := BuildStateConf([]string{process}, []string{target}, d.Timeout(schema.TimeoutCreate), 2*time.Minute, ossService.OssBucketSyncStateRefreshFunc(bucketName, []string{failed}))
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

	if d.HasChange("logging") {
		log.Print("changes in logging")
		err := resourceAlibabacloudStackOssBucketLoggingCreate(client, d)
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
			err = ossService.DeleteBucketTags(bucketName)
		} else {
			var tag_objs []OssTags
			for k, v := range tags {
				tag_objs = append(tag_objs, OssTags{
					Key:   k,
					Value: v.(string),
				})
			}
			err = ossService.PutOssBucketTags(bucketName, tag_objs)
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
	// ossendpoint, err := ossService.GetOssEndpointForCluster(d.Get("cluster").(string))
	// if err != nil {
	// 	return errmsgs.WrapError(err)
	// }
	// bvserver := BucketVpcService{client}
	// vpclist, binderr := bvserver.BucketVpcList(d.Id())
	// if binderr != nil {
	// 	return errmsgs.WrapError(binderr)
	// }
	// var vlist []string
	// if len(vpclist.VpcList) > 0 {
	// 	for _, v := range vpclist.VpcList {
	// 		vpc := v.(map[string]interface{})
	// 		client2 := meta.(*connectivity.AlibabacloudStackClient)
	// 		bvserver := BucketVpcService{client2}
	// 		binderr := bvserver.UnBindBucket(vpc["vpcId"].(string), d.Id(), ossendpoint)
	// 		if binderr != nil {
	// 			return errmsgs.WrapError(binderr)
	// 		}
	// 	}
	// }
	det, err := ossService.DescribeOssBucket(d.Id())
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), "IsBucketExist", errmsgs.AlibabacloudStackOssGoSdk)
	}
	addDebug("IsBucketExist", det, nil, map[string]string{"bucketName": d.Id()})
	if *det.Name == "" {
		return nil
	}

	if err := ossService.DeleteBucket(d.Id()); err != nil {
		return errmsgs.WrapError(err)
	}
	return errmsgs.WrapError(ossService.WaitForOssBucket(d.Id(), Deleted, DefaultTimeoutMedium))
}

func checkVpcListChange(d *schema.ResourceData, meta interface{}) error {
	// FIXME: Calling this interface will add a rule that denies all permissions
	client := meta.(*connectivity.AlibabacloudStackClient)
	bvserver := BucketVpcService{client}
	ossService := OssService{client}
	vpcServer := VpcService{client}

	ossendpoint, err := ossService.GetOssEndpointForCluster(d.Get("oss_cluster").(string))
	if err != nil {
		return errmsgs.WrapError(err)
	}

	// Get old and new VPC lists as Sets
	oldVal, newVal := d.GetChange("vpclist")
	oldVpcs := oldVal.(*schema.Set)
	newVpcs := newVal.(*schema.Set)

	// Calculate VPCs to remove (in old but not in new)
	toRemove := oldVpcs.Difference(newVpcs)
	for _, vpcId := range toRemove.List() {
		binderr := bvserver.UnBindBucket(vpcId.(string), d.Id(), ossendpoint)
		if binderr != nil {
			return errmsgs.WrapError(binderr)
		}
	}

	// Calculate VPCs to add (in new but not in old)
	toAdd := newVpcs.Difference(oldVpcs)
	for _, vpcId := range toAdd.List() {
		vpcdata, err := vpcServer.DescribeVpc(vpcId.(string))
		if err != nil {
			return errmsgs.WrapError(err)
		}
		binderr := bvserver.BindBucket(vpcdata.VpcId, vpcdata.VpcName, vpcdata.CidrBlock, d.Id(), ossendpoint)
		if binderr != nil {
			return errmsgs.WrapError(binderr)
		}
	}

	return nil
}

func resourceAlibabacloudStackOssBucketLoggingCreate(client *connectivity.AlibabacloudStackClient, d *schema.ResourceData) error {
	bucket_name := d.Id()
	ossService := OssService{client}
	ossClient, err := ossService.GetBucketClient(d.Id())
	if err != nil {
		return errmsgs.WrapError(err)
	}

	existingLogging, err := ossService.DescribeOssBucketLogging(d.Id())
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), "GetBucketLogging", errmsgs.AlibabacloudStackOssGoSdk)
	}

	loggingEnabled := existingLogging != nil && existingLogging.BucketLoggingStatus != nil && existingLogging.BucketLoggingStatus.LoggingEnabled != nil

	if loggingEnabled {
		log.Printf("logging is not null %v", d.Get("logging"))
		if _, v := d.GetOk("logging"); v == false {
			log.Print("logging is being disabled")
			_, err = ossClient.DeleteBucketLogging(context.TODO(), &oss.DeleteBucketLoggingRequest{
				Bucket: &bucket_name,
			})
			if err != nil {
				return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), "DeleteBucketLogging", errmsgs.AlibabacloudStackOssGoSdk)
			}
			log.Printf("deleting logs oss done")
		} else {
			logging := make(map[string]interface{})
			log.Print("logging to be updated")
			if v := d.Get("logging"); v != nil {
				log.Print("logging is being enabled")
				all, ok := v.([]interface{})
				if ok {
					log.Printf("printall %v", all)
					for _, a := range all {
						logging, _ = a.(map[string]interface{})
						log.Printf("check target_bucket %v", logging["target_bucket"])
						log.Printf("check target_prefix %v", logging["target_prefix"])
					}
					targetBucketName := fmt.Sprint(logging["target_bucket"])
					log.Printf("checking bucket %v", targetBucketName)
					oldOssService := OssService{client}
					_, err := oldOssService.DescribeOssBucket(targetBucketName)
					if err != nil {
						return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_oss_bucket", "DescribeBucket")
					}
					targetBucket := fmt.Sprint(logging["target_bucket"])
					targetPrefix := fmt.Sprint(logging["target_prefix"])
					_, err = ossClient.PutBucketLogging(context.TODO(), &oss.PutBucketLoggingRequest{
						Bucket: &bucket_name,
						BucketLoggingStatus: &oss.BucketLoggingStatus{
							LoggingEnabled: &oss.LoggingEnabled{
								TargetBucket: &targetBucket,
								TargetPrefix: &targetPrefix,
							},
						},
					})
					if err != nil {
						return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), "PutBucketLogging", errmsgs.AlibabacloudStackOssGoSdk)
					}
					log.Printf("logging oss done")
				}
			}
		}
	} else {
		logging := make(map[string]interface{})
		log.Print("logging is  null")
		if v := d.Get("logging"); v != nil {
			log.Print("logging is being enabled")
			all, ok := v.([]interface{})
			if ok {
				log.Printf("printall %v", all)
				for _, a := range all {
					logging, _ = a.(map[string]interface{})
					log.Printf("check target_bucket %v", logging["target_bucket"])
					log.Printf("check target_prefix %v", logging["target_prefix"])
				}
				targetBucketName := fmt.Sprint(logging["target_bucket"])
				log.Printf("checking bucket %v", targetBucketName)
				oldOssService := OssService{client}
				_, err := oldOssService.DescribeOssBucket(targetBucketName)
				if err != nil {
					return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_oss_bucket", "DescribeBucket")
				}
				targetBucket := fmt.Sprint(logging["target_bucket"])
				targetPrefix := fmt.Sprint(logging["target_prefix"])
				_, err = ossClient.PutBucketLogging(context.TODO(), &oss.PutBucketLoggingRequest{
					Bucket: &bucket_name,
					BucketLoggingStatus: &oss.BucketLoggingStatus{
						LoggingEnabled: &oss.LoggingEnabled{
							TargetBucket: &targetBucket,
							TargetPrefix: &targetPrefix,
						},
					},
				})
				if err != nil {
					return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), "PutBucketLogging", errmsgs.AlibabacloudStackOssGoSdk)
				}
				log.Printf("logging oss done")
			}
		}
	}

	return nil
}

type OssTags struct {
	Value string `name:"Value"`
	Key   string `name:"Key"`
}
