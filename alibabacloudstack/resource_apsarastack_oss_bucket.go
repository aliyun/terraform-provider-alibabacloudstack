package alibabacloudstack

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/helper/hashcode"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackOssBucket() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlibabacloudStackOssBucketCreate,
		Read:   resourceAlibabacloudStackOssBucketRead,
		Update: resourceAlibabacloudStackOssBucketUpdate,
		Delete: resourceAlibabacloudStackOssBucketDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

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

			"cors_rule": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"allowed_headers": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"allowed_methods": {
							Type:     schema.TypeList,
							Required: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"allowed_origins": {
							Type:     schema.TypeList,
							Required: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"expose_headers": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"max_age_seconds": {
							Type:     schema.TypeInt,
							Optional: true,
						},
					},
				},
				MaxItems: 10,
			},

			"website": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"index_document": {
							Type:     schema.TypeString,
							Required: true,
						},

						"error_document": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
				MaxItems: 1,
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

			"referer_config": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"allow_empty": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"referers": {
							Type:     schema.TypeList,
							Required: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
				MaxItems: 1,
			},

			"lifecycle_rule": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.StringLenBetween(0, 255),
						},
						"prefix": {
							Type:     schema.TypeString,
							Required: true,
						},
						"enabled": {
							Type:     schema.TypeBool,
							Required: true,
						},
						"expiration": {
							Type:     schema.TypeSet,
							Optional: true,
							Set:      expirationHash,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"date": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validateOssBucketDateTimestamp,
									},
									"days": {
										Type:     schema.TypeInt,
										Optional: true,
									},
								},
							},
						},
						//"transitions": {
						//	Type:     schema.TypeSet,
						//	Optional: true,
						//	Set:      transitionsHash,
						//	Elem: &schema.Resource{
						//		Schema: map[string]*schema.Schema{
						//			"created_before_date": {
						//				Type:         schema.TypeString,
						//				Optional:     true,
						//				ValidateFunc: validateOssBucketDateTimestamp,
						//			},
						//			"days": {
						//				Type:     schema.TypeInt,
						//				Optional: true,
						//			},
						//			"storage_class": {
						//				Type:     schema.TypeString,
						//				Default:  oss.StorageStandard,
						//				Optional: true,
						//				ValidateFunc: validation.StringInSlice([]string{
						//					string(oss.StorageStandard),
						//					string(oss.StorageIA),
						//					string(oss.StorageArchive),
						//				}, false),
						//			},
						//		},
						//	},
						//},
					},
				},
				MaxItems: 1000,
			},

			"policy": {
				Type:     schema.TypeString,
				Optional: true,
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
				Default:  oss.StorageStandard,
				Optional: true,
				ForceNew: true,
			},
			"sse_algorithm": {
				Type:     schema.TypeString,
				Default:  "",
				Optional: true,
				ForceNew: true,
			},
			"server_side_encryption_rule": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"sse_algorithm": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								ServerSideEncryptionAes256,
								ServerSideEncryptionKMS,
							}, false),
						},
					},
				},
				MaxItems: 1,
			},

			"tags": tagsSchema(),

			"force_destroy": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"versioning": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"status": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"Enabled",
								"Suspended",
							}, false),
						},
					},
				},
				MaxItems: 1,
			},
			"vpclist": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceAlibabacloudStackOssBucketCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	ossService := OssService{client}
	request := map[string]string{"bucketName": d.Get("bucket").(string), "Department": client.Department, "ResourceGroup": client.ResourceGroup}
	var requestInfo *oss.Client
	bucketName := d.Get("bucket").(string)
	det, err := ossService.DescribeOssBucket(bucketName)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alibabacloudstack_oss_bucket", "IsBucketExist", AlibabacloudStackOssGoSdk)
	}
	acl := d.Get("acl").(string)
	storageClass := d.Get("storage_class")
	if storageClass == "" {
		storageClass = "Standard"
	}
	if acl == "" {
		acl = "private"
	}

	sse_algo := d.Get("storage_class")
	// If not present, Create Bucket
	if det.BucketInfo.Name == "" {
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
			"OpenApiAction":    "PutBucket",
			"ProductName":      "oss",
			"Params":           fmt.Sprintf("{\"%s\":%s,\"%s\":%s,\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":\"%s\"}", "Department", client.Department, "ResourceGroup", client.ResourceGroup, "RegionId", client.RegionId, "asVersion", "enterprise", "asArchitechture", "x86", "haAlibabacloudStack", "true", "Language", "en", "BucketName", bucketName, "StorageClass", storageClass, "x-oss-acl", acl, "SSEAlgorithm", sse_algo, "x-one-console-endpoint", "http://oss-cn-neimeng-env30-d01-a.intra.env30.shuguang.com"), //,"x-one-console-endpoint","http://oss-cn-neimeng-env30-d01-a.intra.env30.shuguang.com"),
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
		log.Printf("Response of Create Bucket: %s", raw)
		log.Printf("Bresponse ossbucket before error")
		if err != nil {
			if ossNotFoundError(err) {
				return WrapErrorf(err, NotFoundMsg, AlibabacloudStackOssGoSdk)
			}
			return WrapErrorf(err, DefaultErrorMsg, bucketName, "CreateBucketInfo", AlibabacloudStackOssGoSdk)
		}
		log.Printf("Bresponse ossbucket after error")
		addDebug("CreateBucketInfo", raw, requestInfo, request)
		log.Printf("Bresponse ossbucket check")
		bresponse, _ := raw.(*responses.CommonResponse)
		log.Printf("Bresponse ossbucket %s", bresponse)

		if bresponse.GetHttpStatus() != 200 {
			return WrapErrorf(err, DefaultErrorMsg, "alibabacloudstack_oss_bucket", "CreateBucket", AlibabacloudStackOssGoSdk)
		}
		//logging:= make(map[string]interface{})
		log.Printf("Enter for logging")

		//addDebug("CreateBucket", raw, requestInfo, bresponse.GetHttpContentString())

	}

	err = resource.Retry(3*time.Minute, func() *resource.RetryError {
		det, err := ossService.DescribeOssBucket(bucketName)
		if err != nil {
			return resource.NonRetryableError(err)
		}
		if det.BucketInfo.Name == "" {
			return resource.RetryableError(Error("Trying to ensure new OSS bucket %#v has been created successfully.", request["bucketName"]))
		}
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alibabacloudstack_oss_bucket", "Bucket Not Found", AlibabacloudStackOssGoSdk)
	}

	// Assign the bucket name as the resource ID
	d.SetId(bucketName)
	if v := d.Get("logging"); v != nil {
		log.Printf("Enter for logging condition passed")

		err = resourceAlibabacloudStackOssBucketLoggingCreate(client, d)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "alibabacloudstack_oss_bucket", "Logging Failed", AlibabacloudStackOssGoSdk)
		}
	}
	newlist := d.Get("vpclist").([]interface{})
	bvclient := meta.(*connectivity.AlibabacloudStackClient)
	bvserver := BucketVpcService{bvclient}
	vpclist, binderr := bvserver.BucketVpcList(bucketName)
	if binderr != nil {
		return WrapError(binderr)
	}
	oldlist := vpclist.VpcList
	vpc_err := checkVpcListChange(oldlist, newlist, d, meta)
	if vpc_err != nil {
		return WrapError(vpc_err)
	}
	return resourceAlibabacloudStackOssBucketRead(d, meta)
}

func resourceAlibabacloudStackOssBucketRead(d *schema.ResourceData, meta interface{}) error {
	wiatSecondsIfWithTest(1)
	client := meta.(*connectivity.AlibabacloudStackClient)
	ossService := OssService{client}
	object, err := ossService.DescribeOssBucket(d.Id())
	acl := d.Get("acl").(string)
	if acl == "" {
		acl = "private"
	}
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	logging, err := resourceAlibabacloudStackOssBucketLoggingDescribe(client, d)
	log.Printf("read describe logging %v", logging)
	log.Printf("read describe error %v", err)
	d.Set("bucket", d.Id())
	if object.BucketInfo.Name == "" {
		log.Print("read: BucketInfo fail!!!!!!")
	}
	d.Set("acl", acl)
	d.Set("creation_date", object.BucketInfo.CreationDate.Format("2006-01-02"))
	d.Set("extranet_endpoint", object.BucketInfo.ExtranetEndpoint)
	d.Set("intranet_endpoint", object.BucketInfo.IntranetEndpoint)
	d.Set("location", object.BucketInfo.Location)
	d.Set("owner", object.BucketInfo.Owner.ID)
	d.Set("storage_class", object.BucketInfo.StorageClass)
	var list []map[string]interface{}
	desclog := logging.Data.BucketLoggingStatus.LoggingEnabled
	list = append(list, map[string]interface{}{"target_bucket": desclog.TargetBucket, "target_prefix": desclog.TargetPrefix})
	//d.Set("logging",list)
	if err = d.Set("logging", list); err != nil {
		return WrapError(err)
	}
	bvclient := meta.(*connectivity.AlibabacloudStackClient)
	bvserver := BucketVpcService{bvclient}
	vpclist, binderr := bvserver.BucketVpcList(d.Get("bucket").(string))
	if binderr != nil {
		return WrapError(binderr)
	}
	var vlist []string
	if len(vpclist.VpcList) > 0 {
		for _, v := range vpclist.VpcList {
			vpc := v.(map[string]interface{})
			vlist = append(vlist, vpc["vpcId"].(string))
		}
	}
	d.Set("vpclist", vlist)
	//request := map[string]string{"bucketName": d.Id()}
	//var requestInfo *oss.Client
	//
	//raw, err := client.WithOssClient(func(ossClient *oss.Client) (interface{}, error) {
	//	return ossClient.GetBucketLogging(d.Id())
	//})
	//if err != nil {
	//	return WrapErrorf(err, DefaultErrorMsg, d.Id(), "GetBucketLogging", AlibabacloudStackOssGoSdk)
	//}
	//
	//addDebug("GetBucketLogging", raw, requestInfo, request)
	//logging, _ := raw.(oss.GetBucketLoggingResult)

	//if &logging != nil {
	//	enable := logging.LoggingEnabled
	//	if &enable != nil {
	//		lgs := make([]map[string]interface{}, 0)
	//		tb := logging.LoggingEnabled.TargetBucket
	//		tp := logging.LoggingEnabled.TargetPrefix
	//		if tb != "" || tp != "" {
	//			lgs = append(lgs, map[string]interface{}{
	//				"target_bucket": tb,
	//				"target_prefix": tp,
	//			})
	//		}
	//		if err := d.Set("logging", lgs); err != nil {
	//			return WrapError(err)
	//		}
	//	}
	//}

	//if &object.BucketInfo.SseRule != nil {
	//	if len(object.BucketInfo.SseRule.SSEAlgorithm) > 0 && object.BucketInfo.SseRule.SSEAlgorithm != "None" {
	//		rule := make(map[string]interface{})
	//		rule["sse_algorithm"] = object.BucketInfo.SseRule.SSEAlgorithm
	//		data := make([]map[string]interface{}, 0)
	//		data = append(data, rule)
	//		d.Set("server_side_encryption_rule", data)
	//	}
	//}
	//
	//if object.BucketInfo.Versioning != "" {
	//	data := map[string]interface{}{
	//		"status": object.BucketInfo.Versioning,
	//	}
	//	versioning := make([]map[string]interface{}, 0)
	//	versioning = append(versioning, data)
	//	d.Set("versioning", versioning)
	//}
	//request := map[string]string{"bucketName": d.Id(), "Department": client.Department, "ResourceGroup": client.ResourceGroup}
	//var requestInfo *oss.Client
	//
	//// Read the CORS
	//raw, err := client.WithOssClient(func(ossClient *oss.Client) (interface{}, error) {
	//	requestInfo = ossClient
	//	return ossClient.GetBucketCORS(request["bucketName"])
	//})
	//if err != nil && !IsExpectedErrors(err, []string{"NoSuchCORSConfiguration"}) {
	//	return WrapErrorf(err, DefaultErrorMsg, d.Id(), "GetBucketCORS", AlibabacloudStackOssGoSdk)
	//}
	//addDebug("GetBucketCORS", raw, requestInfo, request)
	//cors, _ := raw.(oss.GetBucketCORSResult)
	//rules := make([]map[string]interface{}, 0, len(cors.CORSRules))
	//for _, r := range cors.CORSRules {
	//	rule := make(map[string]interface{})
	//	rule["allowed_headers"] = r.AllowedHeader
	//	rule["allowed_methods"] = r.AllowedMethod
	//	rule["allowed_origins"] = r.AllowedOrigin
	//	rule["expose_headers"] = r.ExposeHeader
	//	rule["max_age_seconds"] = r.MaxAgeSeconds
	//
	//	rules = append(rules, rule)
	//}
	//if err := d.Set("cors_rule", rules); err != nil {
	//	return WrapError(err)
	//}
	//
	//// Read the website configuration
	//raw, err = client.WithOssClient(func(ossClient *oss.Client) (interface{}, error) {
	//	return ossClient.GetBucketWebsite(d.Id())
	//})
	//if err != nil && !IsExpectedErrors(err, []string{"NoSuchWebsiteConfiguration"}) {
	//	return WrapErrorf(err, DefaultErrorMsg, d.Id(), "GetBucketWebsite", AlibabacloudStackOssGoSdk)
	//}
	//addDebug("GetBucketWebsite", raw, requestInfo, request)
	//ws, _ := raw.(oss.GetBucketWebsiteResult)
	//websites := make([]map[string]interface{}, 0)
	//if err == nil && &ws != nil {
	//	w := make(map[string]interface{})
	//
	//	if v := &ws.IndexDocument; v != nil {
	//		w["index_document"] = v.Suffix
	//	}
	//
	//	if v := &ws.ErrorDocument; v != nil {
	//		w["error_document"] = v.Key
	//	}
	//	websites = append(websites, w)
	//}
	//if err := d.Set("website", websites); err != nil {
	//	return WrapError(err)
	//}
	//
	//// Read the logging configuration
	//raw, err = client.WithOssClient(func(ossClient *oss.Client) (interface{}, error) {
	//	return ossClient.GetBucketLogging(d.Id())
	//})
	//if err != nil {
	//	return WrapErrorf(err, DefaultErrorMsg, d.Id(), "GetBucketLogging", AlibabacloudStackOssGoSdk)
	//}
	//addDebug("GetBucketLogging", raw, requestInfo, request)
	//logging, _ := raw.(oss.GetBucketLoggingResult)
	//
	//if &logging != nil {
	//	enable := logging.LoggingEnabled
	//	if &enable != nil {
	//		lgs := make([]map[string]interface{}, 0)
	//		tb := logging.LoggingEnabled.TargetBucket
	//		tp := logging.LoggingEnabled.TargetPrefix
	//		if tb != "" || tp != "" {
	//			lgs = append(lgs, map[string]interface{}{
	//				"target_bucket": tb,
	//				"target_prefix": tp,
	//			})
	//		}
	//		if err := d.Set("logging", lgs); err != nil {
	//			return WrapError(err)
	//		}
	//	}
	//}
	//
	//// Read the bucket referer
	//raw, err = client.WithOssClient(func(ossClient *oss.Client) (interface{}, error) {
	//	return ossClient.GetBucketReferer(d.Id())
	//})
	//if err != nil {
	//	return WrapErrorf(err, DefaultErrorMsg, d.Id(), "GetBucketReferer", AlibabacloudStackOssGoSdk)
	//}
	//addDebug("GetBucketReferer", raw, requestInfo, request)
	//referers := make([]map[string]interface{}, 0)
	//referer, _ := raw.(oss.GetBucketRefererResult)
	//if len(referer.RefererList) > 0 {
	//	referers = append(referers, map[string]interface{}{
	//		"allow_empty": referer.AllowEmptyReferer,
	//		"referers":    referer.RefererList,
	//	})
	//	if err := d.Set("referer_config", referers); err != nil {
	//		return WrapError(err)
	//	}
	//}
	//
	//// Read the lifecycle rule configuration
	//raw, err = client.WithOssClient(func(ossClient *oss.Client) (interface{}, error) {
	//	return ossClient.GetBucketLifecycle(d.Id())
	//})
	//if err != nil && !ossNotFoundError(err) {
	//	return WrapErrorf(err, DefaultErrorMsg, d.Id(), "GetBucketLifecycle", AlibabacloudStackOssGoSdk)
	//}
	//addDebug("GetBucketLifecycle", raw, requestInfo, request)
	//lrules := make([]map[string]interface{}, 0)
	//lifecycle, _ := raw.(oss.GetBucketLifecycleResult)
	//for _, lifecycleRule := range lifecycle.Rules {
	//	rule := make(map[string]interface{})
	//	rule["id"] = lifecycleRule.ID
	//	rule["prefix"] = lifecycleRule.Prefix
	//	if LifecycleRuleStatus(lifecycleRule.Status) == ExpirationStatusEnabled {
	//		rule["enabled"] = true
	//	} else {
	//		rule["enabled"] = false
	//	}
	//	// expiration
	//	if lifecycleRule.Expiration != nil {
	//		e := make(map[string]interface{})
	//		if lifecycleRule.Expiration.Date != "" {
	//			t, err := time.Parse("2006-01-02T15:04:05.000Z", lifecycleRule.Expiration.Date)
	//			if err != nil {
	//				return WrapError(err)
	//			}
	//			e["date"] = t.Format("2006-01-02")
	//		}
	//		e["days"] = int(lifecycleRule.Expiration.Days)
	//		rule["expiration"] = schema.NewSet(expirationHash, []interface{}{e})
	//	}
	//	// transitions
	//	//if len(lifecycleRule.Transitions) != 0 {
	//	//	var eSli []interface{}
	//	//	for _, transition := range lifecycleRule.Transitions {
	//	//		e := make(map[string]interface{})
	//	//		if transition.CreatedBeforeDate != "" {
	//	//			t, err := time.Parse("2006-01-02T15:04:05.000Z", transition.CreatedBeforeDate)
	//	//			if err != nil {
	//	//				return WrapError(err)
	//	//			}
	//	//			e["created_before_date"] = t.Format("2006-01-02")
	//	//		}
	//	//		e["days"] = transition.Days
	//	//		e["storage_class"] = string(transition.StorageClass)
	//	//		eSli = append(eSli, e)
	//	//	}
	//	//	rule["transitions"] = schema.NewSet(transitionsHash, eSli)
	//	//}
	//
	//	lrules = append(lrules, rule)
	//}
	//
	//if err := d.Set("lifecycle_rule", lrules); err != nil {
	//	return WrapError(err)
	//}
	//
	//// Read Policy
	//raw, err = client.WithOssClient(func(ossClient *oss.Client) (interface{}, error) {
	//	params := map[string]interface{}{}
	//	params["policy"] = nil
	//	return ossClient.Conn.Do("GET", d.Id(), "", params, nil, nil, 0, nil)
	//})
	//
	//if err != nil && !ossNotFoundError(err) {
	//	return WrapErrorf(err, DefaultErrorMsg, d.Id(), "GetPolicyByConn", AlibabacloudStackOssGoSdk)
	//}
	//addDebug("GetPolicyByConn", raw, requestInfo, request)
	//policy := ""
	//if err == nil {
	//	rawResp := raw.(*oss.Response)
	//	defer rawResp.Body.Close()
	//	rawData, err := ioutil.ReadAll(rawResp.Body)
	//	if err != nil {
	//		return WrapError(err)
	//	}
	//	policy = string(rawData)
	//}
	//
	//if err := d.Set("policy", policy); err != nil {
	//	return WrapError(err)
	//}
	//
	//// Read tags
	//raw, err = client.WithOssClient(func(ossClient *oss.Client) (interface{}, error) {
	//	return ossClient.GetBucketTagging(d.Id())
	//})
	//
	//if err != nil {
	//	return WrapErrorf(err, DefaultErrorMsg, d.Id(), "GetBucketTagging", AlibabacloudStackOssGoSdk)
	//}
	//addDebug("GetBucketTagging", raw, requestInfo, request)
	//tagging, _ := raw.(oss.GetBucketTaggingResult)
	//tagsMap := make(map[string]string)
	//if len(tagging.Tags) > 0 {
	//	for _, t := range tagging.Tags {
	//		tagsMap[t.Key] = t.Value
	//	}
	//}
	//if err := d.Set("tags", tagsMap); err != nil {
	//	return WrapError(err)
	//}

	return nil
}

func resourceAlibabacloudStackOssBucketUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	d.Partial(true)
	if d.HasChange("logging") {
		//if err := resourceAlibabacloudStackOssBucketLoggingUpdate(client, d); err != nil {
		//	return WrapError(err)
		//}
		////d.SetPartial("logging")
		log.Print("changes in logging")
		err := resourceAlibabacloudStackOssBucketLoggingCreate(client, d)
		if err != nil {
			return WrapError(err)
		}
	}
	if d.HasChange("vpclist") {
		o, n := d.GetChange("vpclist")
		oldlist := o.([]interface{})
		newlist := n.([]interface{})
		vpc_err := checkVpcListChange(oldlist, newlist, d, meta)
		if vpc_err != nil {
			return WrapError(vpc_err)
		}
	}
	d.Partial(false)
	return resourceAlibabacloudStackOssBucketRead(d, meta)
}

func resourceAlibabacloudStackOssBucketDelete(d *schema.ResourceData, meta interface{}) error {
	bvclient := meta.(*connectivity.AlibabacloudStackClient)
	bvserver := BucketVpcService{bvclient}
	vpclist, binderr := bvserver.BucketVpcList(d.Id())
	if binderr != nil {
		return WrapError(binderr)
	}
	var vlist []string
	if len(vpclist.VpcList) > 0 {
		for _, v := range vpclist.VpcList {
			vpc := v.(map[string]interface{})
			client2 := meta.(*connectivity.AlibabacloudStackClient)
			bvserver := BucketVpcService{client2}
			binderr := bvserver.UnBindBucket(vpc["vpcId"].(string), d.Id())
			if binderr != nil {
				return WrapError(binderr)
			}
		}
	}
	d.Set("vpclist", vlist)
	client := meta.(*connectivity.AlibabacloudStackClient)
	ossService := OssService{client}
	var requestInfo *oss.Client
	det, err := ossService.DescribeOssBucket(d.Id())
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "IsBucketExist", AlibabacloudStackOssGoSdk)
	}
	addDebug("IsBucketExist", det.BucketInfo, requestInfo, map[string]string{"bucketName": d.Id()})
	if det.BucketInfo.Name == "" {
		return nil
	}

	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
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
			"OpenApiAction":    "DeleteBucket",
			"ProductName":      "oss",
			"Params":           fmt.Sprintf("{\"%s\":%s,\"%s\":%s,\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":\"%s\"}", "Department", client.Department, "ResourceGroup", client.ResourceGroup, "RegionId", client.RegionId, "asVersion", "enterprise", "asArchitechture", "x86", "haAlibabacloudStack", "true", "Language", "en", "BucketName", d.Id(), "StorageClass", "Standard"), //,"x-one-console-endpoint","http://oss-cn-neimeng-env30-d01-a.intra.env30.shuguang.com"),

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

		_, err := client.WithOssNewClient(func(ossClient *ecs.Client) (interface{}, error) {

			return ossClient.ProcessCommonRequest(request)
		})

		if err != nil {
			if ossNotFoundError(err) {
				return resource.NonRetryableError(err)
			}
			return resource.RetryableError(err)
		}
		det, err := ossService.DescribeOssBucket(d.Id())
		if err != nil {
			return resource.NonRetryableError(err)
		}
		if det.BucketInfo.Name != "" {
			return resource.RetryableError(Error("Trying to delete OSS bucket %#v successfully.", d.Id()))
		}
		return nil
	})
	log.Print(err)
	return WrapError(ossService.WaitForOssBucket(d.Id(), Deleted, DefaultTimeoutMedium))
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
				return WrapError(binderr)
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
				return WrapError(err)
			}
			client2 := meta.(*connectivity.AlibabacloudStackClient)
			bvserver := BucketVpcService{client2}
			binderr := bvserver.BindBucket(vpcdata.VpcId, vpcdata.VpcName, vpcdata.CidrBlock, d.Id())
			if binderr != nil {
				return WrapError(binderr)
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
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "GetBucketLogging", AlibabacloudStackOssGoSdk)
	}
	var check Logcheck
	if describelogging.Data.BucketLoggingStatus.LoggingEnabled != check {
		log.Printf("logging is not null %v", d.Get("logging"))
		if _, v := d.GetOk("logging"); v == false {
			log.Print("logging is being disabled")
			logrequest := requests.NewCommonRequest()
			if client.Config.Insecure {
				logrequest.SetHTTPSInsecure(client.Config.Insecure)
			}
			logrequest.QueryParams = map[string]string{

				"AccessKeySecret":  client.SecretKey,
				"Product":          "OneRouter",
				"Department":       client.Department,
				"ResourceGroup":    client.ResourceGroup,
				"RegionId":         client.RegionId,
				"Action":           "DoOpenApi",
				"AccountInfo":      "123456",
				"Version":          "2018-12-12",
				"SignatureVersion": "1.0",
				"OpenApiAction":    "PutBucketLogging",
				"ProductName":      "oss",
				"Content":          fmt.Sprint("<BucketLoggingStatus></BucketLoggingStatus>"),
				//"Content": oss-accesslog/",
				"Params": fmt.Sprintf("{\"%s\":\"%s\"}", "BucketName", d.Id()),
				//"Params": "{\"BucketName\":\"source-sample-bucket\"}",

				//"Params":           fmt.Sprintf("{\"%s\":%s,\"%s\":%s,\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":\"%s\"}", "Department", client.Department, "ResourceGroup", client.ResourceGroup, "RegionId", client.RegionId, "asVersion", "enterprise", "asArchitechture", "x86", "haAlibabacloudStack", "true", "Language", "en", "BucketName", bucketName, "StorageClass", storageClass, "x-oss-acl", acl, "SSEAlgorithm", sse_algo), //,"x-one-console-endpoint","http://oss-cn-neimeng-env30-d01-a.intra.env30.shuguang.com"),

			}
			logrequest.Method = "POST"        // Set request method
			logrequest.Product = "OneRouter"  // Specify product
			logrequest.Version = "2018-12-12" // Specify product version
			logrequest.ServiceCode = "OneRouter"
			if strings.ToLower(client.Config.Protocol) == "https" {
				logrequest.Scheme = "https"
			} else {
				logrequest.Scheme = "http"
			} // Set request scheme. Default: http
			logrequest.ApiName = "DoOpenApi"
			logrequest.Headers = map[string]string{"RegionId": client.RegionId}

			raw, err := client.WithEcsClient(func(ossClient *ecs.Client) (interface{}, error) {

				return ossClient.ProcessCommonRequest(logrequest)
			})
			log.Printf("Response of Logging Bucket: %s", raw)
			if err != nil {
				if ossNotFoundError(err) {
					return WrapErrorf(err, NotFoundMsg, AlibabacloudStackOssGoSdk)
				}
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), "CreateBucketInfo", AlibabacloudStackOssGoSdk)
			}
			log.Printf("deleting logs oss done")
			//}

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
					//b, _ :=json.Marshal(logging)
					//log.Printf("checking b %v",b)
					//bucket:= bytes.NewBuffer(b).String()
					//log.Printf("Checking buckets %v",bucket)
					ossService := OssService{client}
					_, err := ossService.DescribeOssBucket(bucket)

					if err != nil {
						return WrapErrorf(err, DefaultErrorMsg, "alibabacloudstack_oss_bucket", "DescribeBucket")
					}
					logrequest := requests.NewCommonRequest()
					if client.Config.Insecure {
						logrequest.SetHTTPSInsecure(client.Config.Insecure)
					}
					logrequest.QueryParams = map[string]string{

						"AccessKeySecret":  client.SecretKey,
						"Product":          "OneRouter",
						"Department":       client.Department,
						"ResourceGroup":    client.ResourceGroup,
						"RegionId":         client.RegionId,
						"Action":           "DoOpenApi",
						"AccountInfo":      "123456",
						"Version":          "2018-12-12",
						"SignatureVersion": "1.0",
						"OpenApiAction":    "PutBucketLogging",
						"ProductName":      "oss",
						"Content":          fmt.Sprint("<BucketLoggingStatus><LoggingEnabled><TargetBucket>", logging["target_bucket"], "</TargetBucket><TargetPrefix>", logging["target_prefix"], "</TargetPrefix></LoggingEnabled></BucketLoggingStatus>"),
						//"Content": oss-accesslog/",
						"Params": fmt.Sprintf("{\"%s\":\"%s\"}", "BucketName", d.Id()),
						//"Params": "{\"BucketName\":\"source-sample-bucket\"}",

						//"Params":           fmt.Sprintf("{\"%s\":%s,\"%s\":%s,\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":\"%s\"}", "Department", client.Department, "ResourceGroup", client.ResourceGroup, "RegionId", client.RegionId, "asVersion", "enterprise", "asArchitechture", "x86", "haAlibabacloudStack", "true", "Language", "en", "BucketName", bucketName, "StorageClass", storageClass, "x-oss-acl", acl, "SSEAlgorithm", sse_algo), //,"x-one-console-endpoint","http://oss-cn-neimeng-env30-d01-a.intra.env30.shuguang.com"),

					}
					logrequest.Method = "POST"        // Set request method
					logrequest.Product = "OneRouter"  // Specify product
					logrequest.Version = "2018-12-12" // Specify product version
					logrequest.ServiceCode = "OneRouter"
					if strings.ToLower(client.Config.Protocol) == "https" {
						logrequest.Scheme = "https"
					} else {
						logrequest.Scheme = "http"
					} // Set request scheme. Default: http
					logrequest.ApiName = "DoOpenApi"
					logrequest.Headers = map[string]string{"RegionId": client.RegionId}

					raw, err := client.WithEcsClient(func(ossClient *ecs.Client) (interface{}, error) {

						return ossClient.ProcessCommonRequest(logrequest)
					})
					log.Printf("Response of Logging Bucket: %s", raw)
					if err != nil {
						if ossNotFoundError(err) {
							return WrapErrorf(err, NotFoundMsg, AlibabacloudStackOssGoSdk)
						}
						return WrapErrorf(err, DefaultErrorMsg, d.Id(), "CreateBucketInfo", AlibabacloudStackOssGoSdk)
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
				//b, _ :=json.Marshal(logging)
				//log.Printf("checking b %v",b)
				//bucket:= bytes.NewBuffer(b).String()
				//log.Printf("Checking buckets %v",bucket)
				ossService := OssService{client}
				_, err := ossService.DescribeOssBucket(bucket)

				if err != nil {
					return WrapErrorf(err, DefaultErrorMsg, "alibabacloudstack_oss_bucket", "DescribeBucket")
				}
				logrequest := requests.NewCommonRequest()
				if client.Config.Insecure {
					logrequest.SetHTTPSInsecure(client.Config.Insecure)
				}
				logrequest.QueryParams = map[string]string{

					"AccessKeySecret":  client.SecretKey,
					"Product":          "OneRouter",
					"Department":       client.Department,
					"ResourceGroup":    client.ResourceGroup,
					"RegionId":         client.RegionId,
					"Action":           "DoOpenApi",
					"AccountInfo":      "123456",
					"Version":          "2018-12-12",
					"SignatureVersion": "1.0",
					"OpenApiAction":    "PutBucketLogging",
					"ProductName":      "oss",
					"Content":          fmt.Sprint("<BucketLoggingStatus><LoggingEnabled><TargetBucket>", logging["target_bucket"], "</TargetBucket><TargetPrefix>", logging["target_prefix"], "</TargetPrefix></LoggingEnabled></BucketLoggingStatus>"),
					//"Content": oss-accesslog/",
					"Params": fmt.Sprintf("{\"%s\":\"%s\"}", "BucketName", d.Id()),
					//"Params": "{\"BucketName\":\"source-sample-bucket\"}",

					//"Params":           fmt.Sprintf("{\"%s\":%s,\"%s\":%s,\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":\"%s\"}", "Department", client.Department, "ResourceGroup", client.ResourceGroup, "RegionId", client.RegionId, "asVersion", "enterprise", "asArchitechture", "x86", "haAlibabacloudStack", "true", "Language", "en", "BucketName", bucketName, "StorageClass", storageClass, "x-oss-acl", acl, "SSEAlgorithm", sse_algo), //,"x-one-console-endpoint","http://oss-cn-neimeng-env30-d01-a.intra.env30.shuguang.com"),

				}
				logrequest.Method = "POST"        // Set request method
				logrequest.Product = "OneRouter"  // Specify product
				logrequest.Version = "2018-12-12" // Specify product version
				logrequest.ServiceCode = "OneRouter"
				if strings.ToLower(client.Config.Protocol) == "https" {
					logrequest.Scheme = "https"
				} else {
					logrequest.Scheme = "http"
				} // Set request scheme. Default: http
				logrequest.ApiName = "DoOpenApi"
				logrequest.Headers = map[string]string{"RegionId": client.RegionId}

				raw, err := client.WithEcsClient(func(ossClient *ecs.Client) (interface{}, error) {

					return ossClient.ProcessCommonRequest(logrequest)
				})
				log.Printf("Response of Logging Bucket: %s", raw)
				if err != nil {
					if ossNotFoundError(err) {
						return WrapErrorf(err, NotFoundMsg, AlibabacloudStackOssGoSdk)
					}
					return WrapErrorf(err, DefaultErrorMsg, d.Id(), "CreateBucketInfo", AlibabacloudStackOssGoSdk)
				}
				log.Printf("logging oss done")
			}

		}
	}

	return nil
}
func resourceAlibabacloudStackOssBucketLoggingDescribe(client *connectivity.AlibabacloudStackClient, d *schema.ResourceData) (*Logging, error) {

	logdescribe := requests.NewCommonRequest()
	if client.Config.Insecure {
		logdescribe.SetHTTPSInsecure(client.Config.Insecure)
	}
	describelogging := Logging{}
	logdescribe.QueryParams = map[string]string{

		"AccessKeySecret":   client.SecretKey,
		"Product":           "OneRouter",
		"Department":        client.Department,
		"ResourceGroup":     client.ResourceGroup,
		"RegionId":          client.RegionId,
		"Action":            "DoOpenApi",
		"AccountInfo":       "123456",
		"Forwardedregionid": client.RegionId,
		"Version":           "2018-12-12",
		"SignatureVersion":  "1.0",
		"OpenApiAction":     "GetBucketLogging",
		"ProductName":       "oss",
		"Params":            fmt.Sprintf("{\"%s\":\"%s\"}", "BucketName", d.Id()),
	}
	logdescribe.Method = "POST"        // Set request method
	logdescribe.Product = "OneRouter"  // Specify product
	logdescribe.Version = "2018-12-12" // Specify product version
	logdescribe.ServiceCode = "OneRouter"
	if strings.ToLower(client.Config.Protocol) == "https" {
		logdescribe.Scheme = "https"
	} else {
		logdescribe.Scheme = "http"
	} // Set request scheme. Default: http
	logdescribe.ApiName = "DoOpenApi"
	logdescribe.Headers = map[string]string{"RegionId": client.RegionId}

	lograw, err := client.WithEcsClient(func(ossClient *ecs.Client) (interface{}, error) {

		return ossClient.ProcessCommonRequest(logdescribe)
	})
	log.Printf("Response of Logging Bucket: %s", lograw)
	if err != nil {
		return &describelogging, WrapErrorf(err, DefaultErrorMsg, d.Id(), "GetBucketLogging", AlibabacloudStackOssGoSdk)
	}

	osslog, _ := lograw.(*responses.CommonResponse)
	_ = json.Unmarshal(osslog.GetHttpContentBytes(), &describelogging)
	log.Printf("describerawlogging %v", osslog)
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
