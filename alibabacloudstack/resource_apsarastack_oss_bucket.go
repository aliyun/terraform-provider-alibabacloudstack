package alibabacloudstack

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/helper/hashcode"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlibabacloudStackOssBucket() *schema.Resource {
	resource := &schema.Resource{
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
				Default:      oss.ACLPrivate,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"private", "public-read", "public-read-write"}, false),
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
				Optional: true,
				Computed: true,
			},
			"owner": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"storage_class": {
				Type:     schema.TypeString,
				Default:  oss.StorageStandard,
				Optional: true,
				ForceNew: true,
			},
			"vpclist": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"bucket_sync": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
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
	det, err := ossService.DescribeOssBucket(bucketName)
	log.Printf("======================== det:%#v", det)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_oss_bucket", "IsBucketExist", errmsgs.AlibabacloudStackOssGoSdk)
	}
	acl := d.Get("acl").(string)
	storageClass := d.Get("storage_class").(string)
	if storageClass == "" {
		storageClass = "Standard"
	}
	storage_capacity := d.Get("storage_capacity").(int)
	// If not present, Create Bucket
	if det.BucketInfo.Name == "" {
		ossEndpointData, err := ossService.GetOssEndpointList()
		if err != nil {
			return errmsgs.WrapError(err)
		}
		request := client.NewCommonRequest("GET", "OneRouter", "2018-12-12", "DoApi", "")
		request.QueryParams["AppAction"] = "BucketCreate"
		request.QueryParams["AppName"] = "one-console-app-oss"
		params := map[string]interface{}{
			"region":        client.RegionId,
			"BucketName":    bucketName,
			"Department":    client.Department,
			"ResourceGroup": client.ResourceGroup,
			"BucketQuota":   storage_capacity,
			"Size":          1,
			"params": map[string]interface{}{
				"bucketName":         bucketName,
				"department":         client.Department,
				"storageClass":       storageClass,
				"bucketQuota":        storage_capacity,
				"size":               1,
				"xOssAcl":            acl,
				"ossEndpoint":        ossEndpointData["oss-endpoint"],
				"location":           ossEndpointData["location"],
				"dualClusterEnabled": false,
			},
		}
		if querybytes, err := json.Marshal(params); err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "json Marshal", "CreateBucket", errmsgs.AlibabacloudStackOssGoSdk)
		} else {
			request.QueryParams["Params"] = string(querybytes)
		}

		bresponse, err := client.ProcessCommonRequest(request)
		addDebug("CreateBucketInfo", bresponse, request, request.QueryParams)
		if err != nil {
			if bresponse == nil {
				return errmsgs.WrapErrorf(err, "Process Common Request Failed")
			}
			if ossNotFoundError(err) {
				return errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackOssGoSdk)
			}
			errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, bucketName, "CreateBucketInfo", errmsgs.AlibabacloudStackOssGoSdk, errmsg)
		}

		if bresponse.GetHttpStatus() != 200 {
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_oss_bucket", "CreateBucket", errmsgs.AlibabacloudStackOssGoSdk)
		}

		err = resource.Retry(3*time.Minute, func() *resource.RetryError {
			det, err := ossService.DescribeOssBucket(bucketName)
			if err != nil {
				return resource.NonRetryableError(err)
			}
			if det.BucketInfo.Name == "" {
				return resource.RetryableError(errmsgs.Error("Trying to ensure new OSS bucket %#v has been created successfully.", bucketName))
			}
			return nil
		})
		if err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_oss_bucket", "Bucket Not Found", errmsgs.AlibabacloudStackOssGoSdk)
		}
	}
	// Assign the bucket name as the resource ID
	d.SetId(bucketName)
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
	logging, err := resourceAlibabacloudStackOssBucketLoggingDescribe(client, d)
	log.Printf("read describe logging %v", logging)
	d.Set("bucket", d.Id())
	if object.BucketInfo.Name == "" {
		log.Print("read: BucketInfo fail!!!!!!")
	}
	d.Set("creation_date", object.BucketInfo.CreationDate.Format("2006-01-02"))
	d.Set("extranet_endpoint", object.BucketInfo.ExtranetEndpoint)
	d.Set("intranet_endpoint", object.BucketInfo.IntranetEndpoint)
	d.Set("location", object.BucketInfo.Location)
	d.Set("owner", object.BucketInfo.Owner.ID)
	d.Set("storage_class", object.BucketInfo.StorageClass)
	var list []map[string]interface{}
	desclog := logging.Data.BucketLoggingStatus.LoggingEnabled
	list = append(list, map[string]interface{}{"target_bucket": desclog.TargetBucket, "target_prefix": desclog.TargetPrefix})
	if err = d.Set("logging", list); err != nil {
		return errmsgs.WrapError(err)
	}
	bvclient := meta.(*connectivity.AlibabacloudStackClient)
	bvserver := BucketVpcService{bvclient}
	vpclist, binderr := bvserver.BucketVpcList(d.Get("bucket").(string))
	if binderr != nil {
		return errmsgs.WrapError(binderr)
	}
	var vlist []string
	if len(vpclist.VpcList) > 0 {
		for _, v := range vpclist.VpcList {
			vpc := v.(map[string]interface{})
			vlist = append(vlist, vpc["vpcId"].(string))
		}
	}
	d.Set("vpclist", vlist)

	bucketName := d.Get("bucket").(string)

	bucketSync, err := ossService.GetBucketSync(bucketName)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	d.Set("bucket_sync", false)
	for _, rule := range bucketSync.Data.ReplicationConfiguration.Rule {
		if rule.Status == "doing" && rule.SrcLocation == "" {
			// Disaster recovery relationships appear in pairs
			d.Set("bucket_sync", true)
			d.Set("dual_sync_role", rule.SyncRole)
			dual_kms_key, ok := rule.EncryptionConfiguration["ReplicaKmsKeyID"]
			if ok {
				d.Set("dual_kms_key", dual_kms_key)
			}
			break
		}
	}

	// Get ACL information
	request := client.NewCommonRequest("GET", "OneRouter", "2018-12-12", "DoOpenApi", "")
	request.QueryParams["OpenApiAction"] = "GetBucketAcl"
	request.QueryParams["ProductName"] = "oss"
	request.QueryParams["Params"] = fmt.Sprintf("{\"BucketName\":\"%s\", \"acl\":\"acl\"}", bucketName)

	bresponse, err := client.ProcessCommonRequest(request)
	if err != nil {
		if bresponse == nil {
			return errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		if ossNotFoundError(err) {
			return errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackOssGoSdk)
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, bucketName, "CreateBucketInfo", errmsgs.AlibabacloudStackOssGoSdk, errmsg)
	}
	bucketAcl := BucketAclResponse{}
	json.Unmarshal([]byte(bresponse.GetHttpContentString()), &bucketAcl)
	d.Set("acl", bucketAcl.Data.AccessControlPolicy.AccessControlList.Grant)

	// Get storage capacity information
	request = client.NewCommonRequest("GET", "OneRouter", "2018-12-12", "DoOpenApi", "")
	request.QueryParams["OpenApiAction"] = "GetBucketStorageCapacity"
	request.QueryParams["ProductName"] = "oss"
	request.QueryParams["Params"] = fmt.Sprintf("{\"BucketName\":\"%s\"}", bucketName)

	bresponse, err = client.ProcessCommonRequest(request)
	if err != nil {
		if bresponse == nil {
			return errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		if ossNotFoundError(err) {
			return errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackOssGoSdk)
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, bucketName, "GetBucketStorageCapacity", errmsgs.AlibabacloudStackOssGoSdk, errmsg)
	}
	storageCapacity := BucketStorageCapacityResponse{}
	json.Unmarshal([]byte(bresponse.GetHttpContentString()), &storageCapacity)
	if v, err := strconv.Atoi(storageCapacity.Data.BucketUserQos.StorageCapacity); err == nil {
		d.Set("storage_capacity", v)
	} else {
		return errmsgs.WrapErrorf(err, "Get storage capacity failed")
	}

	// Get encryption information
	request = client.NewCommonRequest("GET", "OneRouter", "2018-12-12", "DoOpenApi", "")
	request.QueryParams["OpenApiAction"] = "GetBucketEncryption"
	request.QueryParams["ProductName"] = "oss"
	request.QueryParams["Params"] = fmt.Sprintf("{\"BucketName\":\"%s\"}", bucketName)

	bresponse, err = client.ProcessCommonRequest(request)
	if err != nil {
		if bresponse == nil {
			return errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		if ossNotFoundError(err) {
			return errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackOssGoSdk)
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, bucketName, "GetBucketEncryption", errmsgs.AlibabacloudStackOssGoSdk, errmsg)
	}
	storageEncryption := BucketEncryptionResponse{}
	json.Unmarshal([]byte(bresponse.GetHttpContentString()), &storageEncryption)
	if storageEncryption.Code == "NoSuchServerSideEncryptionRule" {
		d.Set("sse_algorithm", "")
	} else {
		d.Set("sse_algorithm", storageEncryption.Data.ServerSideEncryptionRule.ApplyServerSideEncryptionByDefault.SSEAlgorithm)
		if storageEncryption.Data.ServerSideEncryptionRule.ApplyServerSideEncryptionByDefault.SSEAlgorithm == "KMS" {
			d.Set("kms_key_id", storageEncryption.Data.ServerSideEncryptionRule.ApplyServerSideEncryptionByDefault.KMSMasterKeyID)
		}
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
	var err error
	if !d.IsNewResource() && d.HasChange("tags") {
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

	if (d.IsNewResource() && d.Get("bucket_sync").(bool)) || (!d.IsNewResource() && d.HasChange("bucket_sync")) {
		request := client.NewCommonRequest("GET", "OneRouter", "2018-12-12", "DoOpenApi", "")
		target := ""
		process := "closing"
		failed := "starting"
		if v := d.Get("bucket_sync").(bool); v {
			request.QueryParams["OpenApiAction"] = "PutBucketSync"
			dual_kms_key := d.Get("dual_kms_key").(string)
			dual_sync_role := d.Get("dual_sync_role").(string)
			if dual_kms_key != "" && dual_sync_role != "" {
				content := `<ReplicationConfiguration><Rule><SyncRole>%s</SyncRole><SourceSelectionCriteria><SseKmsEncryptedObjects><Status>Enabled</Status></SseKmsEncryptedObjects></SourceSelectionCriteria><EncryptionConfiguration><ReplicaKmsKeyID>%s</ReplicaKmsKeyID></EncryptionConfiguration></Rule></ReplicationConfiguration>`
				request.QueryParams["Content"] = fmt.Sprintf(content, dual_sync_role, dual_kms_key)
			} else if (dual_kms_key != "" && dual_sync_role == "") || (dual_kms_key == "" && dual_sync_role != "") {
				return fmt.Errorf("dual_kms_key and dual_sync_role must be set at the same time")
			}
			target = "doing"
			process = "starting"
			failed = "closing"

		} else {
			request.QueryParams["OpenApiAction"] = "DeleteBucketSync"
		}
		request.QueryParams["ProductName"] = "oss"
		request.QueryParams["Params"] = fmt.Sprintf("{\"BucketName\":\"%s\"}", bucketName)

		bresponse, err := client.ProcessCommonRequest(request)
		addDebug(request.GetActionName(), bresponse, request, request.QueryParams)
		if err != nil {
			if bresponse == nil {
				return errmsgs.WrapErrorf(err, "Process Common Request Failed")
			}
			if ossNotFoundError(err) {
				return errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackOssGoSdk)
			}
			errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, bucketName, "CreateBucketInfo", errmsgs.AlibabacloudStackOssGoSdk, errmsg)
		}
		stateConf := BuildStateConf([]string{process}, []string{target}, d.Timeout(schema.TimeoutCreate), 2*time.Minute, ossService.OssBucketSyncStateRefreshFunc(bucketName, []string{failed}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
		}
	}

	if (d.IsNewResource() && d.Get("storage_capacity").(int) != -1) || (!d.IsNewResource() && d.HasChange("storage_capacity")) {
		storageCapacity := d.Get("storage_capacity").(int)
		request := client.NewCommonRequest("GET", "OneRouter", "2018-12-12", "DoOpenApi", "")
		request.QueryParams["OpenApiAction"] = "SetBucketStorageCapacity"
		request.QueryParams["ProductName"] = "oss"
		request.QueryParams["Params"] = fmt.Sprintf("{\"BucketName\":\"%s\", \"StorageCapacity\":%d}", bucketName, storageCapacity)
		request.QueryParams["Content"] = fmt.Sprintf("<BucketUserQos><StorageCapacity>%d</StorageCapacity></BucketUserQos>", storageCapacity)

		bresponse, err := client.ProcessCommonRequest(request)
		if err != nil {
			if bresponse == nil {
				return errmsgs.WrapErrorf(err, "Process Common Request Failed")
			}
			if ossNotFoundError(err) {
				return errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackOssGoSdk)
			}
			errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, bucketName, "PutBucketACL", errmsgs.AlibabacloudStackOssGoSdk, errmsg)
		}
	}

	if d.HasChange("acl") {
		acl := d.Get("acl").(string)
		request := client.NewCommonRequest("GET", "OneRouter", "2018-12-12", "DoOpenApi", "")
		request.QueryParams["OpenApiAction"] = "PutBucketACL"
		request.QueryParams["ProductName"] = "oss"
		request.QueryParams["Params"] = fmt.Sprintf("{\"BucketName\":\"%s\", \"x-oss-acl\":\"%s\"}", bucketName, acl)

		bresponse, err := client.ProcessCommonRequest(request)
		if err != nil {
			if bresponse == nil {
				return errmsgs.WrapErrorf(err, "Process Common Request Failed")
			}
			if ossNotFoundError(err) {
				return errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackOssGoSdk)
			}
			errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, bucketName, "PutBucketACL", errmsgs.AlibabacloudStackOssGoSdk, errmsg)
		}
	}

	if d.HasChanges("sse_algorithm", "kms_key_id") {
		if d.Get("sse_algorithm").(string) == "" {
			request := client.NewCommonRequest("GET", "OneRouter", "2018-12-12", "DoOpenApi", "")
			request.QueryParams["OpenApiAction"] = "DeleteBucketEncryption"
			request.QueryParams["ProductName"] = "oss"
			request.QueryParams["Params"] = fmt.Sprintf("{\"BucketName\":\"%s\"}", bucketName)

			bresponse, err := client.ProcessCommonRequest(request)
			if err != nil {
				if bresponse == nil {
					return errmsgs.WrapErrorf(err, "Process Common Request Failed")
				}
				if ossNotFoundError(err) {
					return errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackOssGoSdk)
				}
				errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
				return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, bucketName, "DeleteBucketEncryption", errmsgs.AlibabacloudStackOssGoSdk, errmsg)
			}
		} else {
			sse_algorithm := d.Get("sse_algorithm").(string)
			kms_key_id := d.Get("kms_key_id").(string)
			request := client.NewCommonRequest("GET", "OneRouter", "2018-12-12", "DoOpenApi", "")
			request.QueryParams["OpenApiAction"] = "PutBucketEncryption"
			request.QueryParams["ProductName"] = "oss"
			request.QueryParams["Params"] = fmt.Sprintf("{\"BucketName\":\"%s\"}", bucketName)
			if sse_algorithm == "KMS" {
				request.QueryParams["Content"] = fmt.Sprintf("<ServerSideEncryptionRule><ApplyServerSideEncryptionByDefault><SSEAlgorithm>KMS</SSEAlgorithm><KMSMasterKeyID>%s</KMSMasterKeyID></ApplyServerSideEncryptionByDefault></ServerSideEncryptionRule>", kms_key_id)
			} else {
				request.QueryParams["Content"] = fmt.Sprintf("<ServerSideEncryptionRule><ApplyServerSideEncryptionByDefault><SSEAlgorithm>%s</SSEAlgorithm></ApplyServerSideEncryptionByDefault></ServerSideEncryptionRule>", sse_algorithm)
			}

			bresponse, err := client.ProcessCommonRequest(request)
			if err != nil {
				if bresponse == nil {
					return errmsgs.WrapErrorf(err, "Process Common Request Failed")
				}
				if ossNotFoundError(err) {
					return errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackOssGoSdk)
				}
				errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
				return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, bucketName, "PutBucketEncryption", errmsgs.AlibabacloudStackOssGoSdk, errmsg)
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
	if d.HasChange("vpclist") {
		// FIXME: Calling this interface will add a rule that denies all permissions
		o, n := d.GetChange("vpclist")
		oldlist := o.([]interface{})
		newlist := n.([]interface{})
		vpc_err := checkVpcListChange(oldlist, newlist, d, meta)
		if vpc_err != nil {
			return errmsgs.WrapError(vpc_err)
		}
	}

	return nil
}

func resourceAlibabacloudStackOssBucketDelete(d *schema.ResourceData, meta interface{}) error {
	bvclient := meta.(*connectivity.AlibabacloudStackClient)
	bvserver := BucketVpcService{bvclient}
	vpclist, binderr := bvserver.BucketVpcList(d.Id())
	if binderr != nil {
		return errmsgs.WrapError(binderr)
	}
	var vlist []string
	if len(vpclist.VpcList) > 0 {
		for _, v := range vpclist.VpcList {
			vpc := v.(map[string]interface{})
			client2 := meta.(*connectivity.AlibabacloudStackClient)
			bvserver := BucketVpcService{client2}
			binderr := bvserver.UnBindBucket(vpc["vpcId"].(string), d.Id())
			if binderr != nil {
				return errmsgs.WrapError(binderr)
			}
		}
	}
	d.Set("vpclist", vlist)
	client := meta.(*connectivity.AlibabacloudStackClient)
	ossService := OssService{client}
	var requestInfo *oss.Client
	det, err := ossService.DescribeOssBucket(d.Id())
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), "IsBucketExist", errmsgs.AlibabacloudStackOssGoSdk)
	}
	addDebug("IsBucketExist", det.BucketInfo, requestInfo, map[string]string{"bucketName": d.Id()})
	if det.BucketInfo.Name == "" {
		return nil
	}

	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		request := client.NewCommonRequest("GET", "OneRouter", "2018-12-12", "DoOpenApi", "")
		request.QueryParams["OpenApiAction"] = "DeleteBucket"
		request.QueryParams["ProductName"] = "oss"
		request.QueryParams["Params"] = fmt.Sprintf("{\"%s\":%s,\"%s\":%s,\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":\"%s\"}", "Department", client.Department, "ResourceGroup", client.ResourceGroup, "RegionId", client.RegionId, "asVersion", "enterprise", "asArchitechture", "x86", "haAlibabacloudStack", "true", "Language", "en", "BucketName", d.Id(), "StorageClass", "Standard")

		bresponse, err := client.ProcessCommonRequest(request)

		if err != nil {
			if bresponse == nil {
				return resource.RetryableError(errmsgs.WrapErrorf(err, "Process Common Request Failed"))
			}
			if ossNotFoundError(err) {
				return resource.NonRetryableError(errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackOssGoSdk))
			}
			errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			return resource.RetryableError(errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), "DeleteBucket", errmsgs.AlibabacloudStackOssGoSdk, errmsg))
		}
		det, err := ossService.DescribeOssBucket(d.Id())
		if err != nil {
			return resource.NonRetryableError(errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), "IsBucketExist", errmsgs.AlibabacloudStackOssGoSdk))
		}
		if det.BucketInfo.Name != "" {
			return resource.RetryableError(errmsgs.Error("Trying to delete OSS bucket %#v successfully.", d.Id()))
		}
		return nil
	})
	log.Print(err)
	return errmsgs.WrapError(ossService.WaitForOssBucket(d.Id(), Deleted, DefaultTimeoutMedium))
}

func checkVpcListChange(oldlist []interface{}, newlist []interface{}, d *schema.ResourceData, meta interface{}) error {
	vpclist := []string{}
	for _, ovpcid := range oldlist {
		isdelete := true
		for _, nvpcid := range newlist {
			if ovpcid == nvpcid {
				isdelete = false
			}
		}
		if isdelete {
			client2 := meta.(*connectivity.AlibabacloudStackClient)
			bvserver := BucketVpcService{client2}
			binderr := bvserver.UnBindBucket(ovpcid.(string), d.Id())
			if binderr != nil {
				return errmsgs.WrapError(binderr)
			}
		}
	}
	for _, nvpcid := range newlist {
		iscreate := true
		vpclist = append(vpclist, nvpcid.(string))
		for _, ovpcid := range oldlist {
			if ovpcid == nvpcid {
				iscreate = false
			}
		}
		if iscreate {
			client := meta.(*connectivity.AlibabacloudStackClient)
			vpcServer := VpcService{client}
			vpcdata, err := vpcServer.DescribeVpc(nvpcid.(string))
			if err != nil {
				return errmsgs.WrapError(err)
			}
			client2 := meta.(*connectivity.AlibabacloudStackClient)
			bvserver := BucketVpcService{client2}
			binderr := bvserver.BindBucket(vpcdata.VpcId, vpcdata.VpcName, vpcdata.CidrBlock, d.Id())
			if binderr != nil {
				return errmsgs.WrapError(binderr)
			}
		}
	}
	return nil
}

func expirationHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})
	if v, ok := m["date"]; ok {
		buf.WriteString(fmt.Sprintf("%s-", v.(string)))
	}
	if v, ok := m["days"]; ok {
		buf.WriteString(fmt.Sprintf("%d-", v.(int)))
	}
	return hashcode.String(buf.String())
}

func transitionsHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})
	if v, ok := m["created_before_date"]; ok {
		buf.WriteString(fmt.Sprintf("%s-", v.(string)))
	}
	if v, ok := m["storage_class"]; ok {
		buf.WriteString(fmt.Sprintf("%s-", v.(string)))
	}
	if v, ok := m["days"]; ok {
		buf.WriteString(fmt.Sprintf("%d-", v.(int)))
	}
	return hashcode.String(buf.String())
}

func resourceAlibabacloudStackOssBucketLoggingCreate(client *connectivity.AlibabacloudStackClient, d *schema.ResourceData) error {
	describelogging, err := resourceAlibabacloudStackOssBucketLoggingDescribe(client, d)

	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), "GetBucketLogging", errmsgs.AlibabacloudStackOssGoSdk)
	}
	var check Logcheck
	if describelogging.Data.BucketLoggingStatus.LoggingEnabled != check {
		log.Printf("logging is not null %v", d.Get("logging"))
		if _, v := d.GetOk("logging"); v == false {
			log.Print("logging is being disabled")
			logrequest := client.NewCommonRequest("GET", "OneRouter", "2018-12-12", "DoOpenApi", "")
			logrequest.QueryParams["OpenApiAction"] = "PutBucketLogging"
			logrequest.QueryParams["ProductName"] = "oss"
			logrequest.QueryParams["Content"] = fmt.Sprint("<BucketLoggingStatus></BucketLoggingStatus>")
			logrequest.QueryParams["Params"] = fmt.Sprintf("{\"%s\":\"%s\"}", "BucketName", d.Id())

			bresponse, err := client.ProcessCommonRequest(logrequest)
			log.Printf("Response of Logging Bucket: %s", bresponse)
			if err != nil {
				if bresponse == nil {
					return errmsgs.WrapErrorf(err, "Process Common Request Failed")
				}
				if ossNotFoundError(err) {
					return errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackOssGoSdk)
				}
				errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
				return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), "CreateBucketInfo", errmsgs.AlibabacloudStackOssGoSdk, errmsg)
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
					bucket := fmt.Sprint(logging["target_bucket"])
					log.Printf("checking bucket %v", bucket)
					ossService := OssService{client}
					_, err := ossService.DescribeOssBucket(bucket)

					if err != nil {
						return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_oss_bucket", "DescribeBucket")
					}
					logrequest := client.NewCommonRequest("GET", "OneRouter", "2018-12-12", "DoOpenApi", "")
					logrequest.QueryParams["OpenApiAction"] = "PutBucketLogging"
					logrequest.QueryParams["ProductName"] = "oss"
					logrequest.QueryParams["Content"] = fmt.Sprint("<BucketLoggingStatus><LoggingEnabled><TargetBucket>", logging["target_bucket"], "</TargetBucket><TargetPrefix>", logging["target_prefix"], "</TargetPrefix></LoggingEnabled></BucketLoggingStatus>")
					logrequest.QueryParams["Params"] = fmt.Sprintf("{\"%s\":\"%s\"}", "BucketName", d.Id())

					bresponse, err := client.ProcessCommonRequest(logrequest)
					log.Printf("Response of Logging Bucket: %s", bresponse)
					if err != nil {
						if bresponse == nil {
							return errmsgs.WrapErrorf(err, "Process Common Request Failed")
						}

						if ossNotFoundError(err) {
							return errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackOssGoSdk)
						}
						errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
						return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), "CreateBucketInfo", errmsgs.AlibabacloudStackOssGoSdk, errmsg)
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
				bucket := fmt.Sprint(logging["target_bucket"])
				log.Printf("checking bucket %v", bucket)
				ossService := OssService{client}
				_, err := ossService.DescribeOssBucket(bucket)

				if err != nil {
					return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_oss_bucket", "DescribeBucket")
				}
				logrequest := client.NewCommonRequest("GET", "OneRouter", "2018-12-12", "DoOpenApi", "")
				logrequest.QueryParams["OpenApiAction"] = "PutBucketLogging"
				logrequest.QueryParams["ProductName"] = "oss"
				logrequest.QueryParams["Content"] = fmt.Sprint("<BucketLoggingStatus><LoggingEnabled><TargetBucket>", logging["target_bucket"], "</TargetBucket><TargetPrefix>", logging["target_prefix"], "</TargetPrefix></LoggingEnabled></BucketLoggingStatus>")
				logrequest.QueryParams["Params"] = fmt.Sprintf("{\"%s\":\"%s\"}", "BucketName", d.Id())

				bresponse, err := client.ProcessCommonRequest(logrequest)
				log.Printf("Response of Logging Bucket: %s", bresponse)
				if err != nil {
					if bresponse == nil {
						return errmsgs.WrapErrorf(err, "Process Common Request Failed")
					}

					if ossNotFoundError(err) {
						return errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackOssGoSdk)
					}
					errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
					return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), "CreateBucketInfo", errmsgs.AlibabacloudStackOssGoSdk, errmsg)
				}
				log.Printf("logging oss done")
			}
		}
	}

	return nil
}

func resourceAlibabacloudStackOssBucketLoggingDescribe(client *connectivity.AlibabacloudStackClient, d *schema.ResourceData) (*Logging, error) {
	logdescribe := client.NewCommonRequest("GET", "OneRouter", "2018-12-12", "DoOpenApi", "")
	logdescribe.QueryParams["Forwardedregionid"] = client.RegionId
	logdescribe.QueryParams["OpenApiAction"] = "GetBucketLogging"
	logdescribe.QueryParams["ProductName"] = "oss"
	logdescribe.QueryParams["Params"] = fmt.Sprintf("{\"%s\":\"%s\"}", "BucketName", d.Id())

	bresponse, err := client.ProcessCommonRequest(logdescribe)
	log.Printf("Response of Logging Bucket: %s", bresponse)
	if err != nil {
		if bresponse == nil {
			return &Logging{}, errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return &Logging{}, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), "GetBucketLogging", errmsgs.AlibabacloudStackOssGoSdk, errmsg)
	}

	var describelogging Logging
	_ = json.Unmarshal(bresponse.GetHttpContentBytes(), &describelogging)
	log.Printf("describerawlogging %v", bresponse)
	log.Printf("describelogging %v", describelogging)

	return &describelogging, nil
}

type Logcheck struct {
	TargetPrefix string `json:"TargetPrefix"`
	TargetBucket string `json:"TargetBucket"`
}

type Logging struct {
	Data struct {
		BucketLoggingStatus struct {
			LoggingEnabled struct {
				TargetPrefix string `json:"TargetPrefix"`
				TargetBucket string `json:"TargetBucket"`
			} `json:"LoggingEnabled"`
		} `json:"BucketLoggingStatus"`
	} `json:"Data"`
	API string `json:"api"`
}

type OssTags struct {
	Value string `name:"Value"`
	Key   string `name:"Key"`
}
