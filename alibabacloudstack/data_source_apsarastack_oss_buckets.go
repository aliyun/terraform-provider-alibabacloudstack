package alibabacloudstack

import (
	"encoding/json"
	"encoding/xml"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"regexp"
	"strings"
)

func dataSourceAlibabacloudStackOssBuckets() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackOssBucketsRead,

		Schema: map[string]*schema.Schema{
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
				ForceNew:     true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},

			// Computed values
			"names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"buckets": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"acl": {
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
							Computed: true,
						},
						"creation_date": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"cors_rules": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"allowed_headers": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"allowed_methods": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"allowed_origins": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"expose_headers": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"max_age_seconds": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},

						"website": {
							Type: schema.TypeList,
							//Computed: true,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"index_document": {
										Type:     schema.TypeString,
										Computed: true,
									},

									"error_document": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
							MaxItems: 1,
						},

						"logging": {
							Type: schema.TypeList,
							//Computed: true,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"target_bucket": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"target_prefix": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
							MaxItems: 1,
						},

						"referer_config": {
							Type: schema.TypeList,
							//Computed: true,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"allow_empty": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"referers": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
							MaxItems: 1,
						},

						"lifecycle_rule": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"prefix": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"enabled": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"expiration": {
										Type: schema.TypeList,
										//Computed: true,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"date": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"days": {
													Type:     schema.TypeInt,
													Optional: true,
												},
											},
										},
										MaxItems: 1,
									},
								},
							},
						},

						"policy": {
							Type:     schema.TypeString,
							Optional: true,
						},

						"server_side_encryption_rule": {
							Type: schema.TypeList,
							//Computed: true,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"sse_algorithm": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
							MaxItems: 1,
						},

						"tags": tagsSchemaComputed(),

						"versioning": {
							Type: schema.TypeList,
							//Computed: true,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"status": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
							MaxItems: 1,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlibabacloudStackOssBucketsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	var requestInfo *oss.Client
	var allBuckets []oss.BucketProperties
	nextMarker := ""
	for {
		var options []oss.Option
		if nextMarker != "" {
			options = append(options, oss.Marker(nextMarker))
		}
		request := requests.NewCommonRequest()

		if client.Config.Insecure {
			request.SetHTTPSInsecure(client.Config.Insecure)
		}
		request.QueryParams = map[string]string{

			"AccessKeySecret": client.SecretKey,
			"Product":         "OneRouter",
			//"Department":       client.Department,
			//"ResourceGroup":    client.ResourceGroup,
			"RegionId":         client.RegionId,
			"Action":           "DoOpenApi",
			"AccountInfo":      "123456",
			"Version":          "2018-12-12",
			"SignatureVersion": "1.0",
			"OpenApiAction":    "GetService",
			"ProductName":      "oss",
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

		var bucketList = &BucketList{}
		raw, err := client.WithOssNewClient(func(ossClient *ecs.Client) (interface{}, error) {

			return ossClient.ProcessCommonRequest(request)
		})

		if err != nil {
			if ossNotFoundError(err) {
				return WrapErrorf(err, NotFoundMsg, AlibabacloudStackOssGoSdk)
			}
			return WrapErrorf(err, DefaultErrorMsg, "GetBucketInfo", AlibabacloudStackOssGoSdk)
		}
		addDebug("GetBucketInfo", raw, requestInfo, request)
		bresponse, _ := raw.(*responses.CommonResponse)

		err = json.Unmarshal(bresponse.GetHttpContentBytes(), bucketList)
		if err != nil {
			return WrapError(err)
		}
		if bucketList.Code != "200" || len(bucketList.Data.ListAllMyBucketsResult.Buckets.Bucket) < 1 {
			return WrapError(err)
		}

		//var found = false
		//for _, j := range bucketList.Data.ListAllMyBucketsResult.Buckets.Bucket {
		//	if j.Name == id {
		//		response.BucketInfo.Name = j.Name
		//		response.BucketInfo.StorageClass = j.StorageClass
		//		response.BucketInfo.ExtranetEndpoint = j.ExtranetEndpoint
		//		response.BucketInfo.IntranetEndpoint = j.IntranetEndpoint
		//		response.BucketInfo.Owner.ID = fmt.Sprint(j.ResourceGroupName)
		//		//response.BucketInfo.CreationDate=fmt.Sprint(j.CreationDate.
		//		response.BucketInfo.Location = j.Location
		//		found = true
		//		break
		//	}
		//}
		//if !found{
		//	break
		//}

		//raw, err := client.WithOssClient(func(ossClient *oss.Client) (interface{}, error) {
		//	requestInfo = ossClient
		//	return ossClient.ListBuckets(options...)
		//})
		//if err != nil {
		//	return WrapErrorf(err, DataDefaultErrorMsg, "alibabacloudstack_oss_bucket", "CreateBucket", AlibabacloudStackOssGoSdk)
		//}
		if debugOn() {
			addDebug("ListBuckets", raw, requestInfo, map[string]interface{}{"options": options})
		}
		for _, k := range bucketList.Data.ListAllMyBucketsResult.Buckets.Bucket {
			allBuckets = append(allBuckets, oss.BucketProperties{
				XMLName:  xml.Name{},
				Name:     k.Name,
				Location: k.Location,
				//CreationDate: ,
				StorageClass:     k.StorageClass,
				CreationDate:     k.CreationDate,
				Extranetendpoint: k.ExtranetEndpoint,
				Intranetendpoint: k.IntranetEndpoint,
			})
		}
		break
	}

	var filteredBucketsTemp []oss.BucketProperties
	nameRegex, ok := d.GetOk("name_regex")
	if ok && nameRegex.(string) != "" {
		var r *regexp.Regexp
		if nameRegex != "" {
			r = regexp.MustCompile(nameRegex.(string))
		}
		for _, bucket := range allBuckets {
			if r != nil && !r.MatchString(bucket.Name) {
				continue
			}
			filteredBucketsTemp = append(filteredBucketsTemp, bucket)
		}
	} else {
		filteredBucketsTemp = allBuckets
	}
	return bucketsDescriptionAttributes(d, filteredBucketsTemp, meta)
}

func bucketsDescriptionAttributes(d *schema.ResourceData, buckets []oss.BucketProperties, meta interface{}) error {
	//client := meta.(*connectivity.AlibabacloudStackClient)

	var ids []string
	var s []map[string]interface{}
	var names []string
	//var requestInfo *oss.Client
	for _, bucket := range buckets {
		mapping := map[string]interface{}{
			"name":              bucket.Name,
			"location":          bucket.Location,
			"storage_class":     bucket.StorageClass,
			"creation_date":     bucket.CreationDate,
			"extranet_endpoint": bucket.Extranetendpoint,
			"intranet_endpoint": bucket.Intranetendpoint,
			//"creation_date": bucket.CreationDate.Format("2006-01-02"),
		}
		ids = append(ids, bucket.Name)
		s = append(s, mapping)
		names = append(names, bucket.Name)
		// Add additional information
		//raw, err := client.WithOssClient(func(ossClient *oss.Client) (interface{}, error) {
		//	requestInfo = ossClient
		//	return ossClient.GetBucketInfo(bucket.Name)
		//})
		//if err == nil {
		//	if debugOn() {
		//		addDebug("GetBucketInfo", raw, requestInfo, map[string]string{"bucketName": bucket.Name})
		//	}
		//	response, _ := raw.(oss.GetBucketInfoResult)
		//	mapping["acl"] = response.BucketInfo.ACL
		//	mapping["extranet_endpoint"] = response.BucketInfo.ExtranetEndpoint
		//	mapping["intranet_endpoint"] = response.BucketInfo.IntranetEndpoint
		//	mapping["owner"] = response.BucketInfo.Owner.ID
		//
		//	//Add ServerSideEncryption information
		//	var sseconfig []map[string]interface{}
		//	if &response.BucketInfo.SseRule != nil {
		//		if len(response.BucketInfo.SseRule.SSEAlgorithm) > 0 && response.BucketInfo.SseRule.SSEAlgorithm != "None" {
		//			data := map[string]interface{}{
		//				"sse_algorithm": response.BucketInfo.SseRule.SSEAlgorithm,
		//			}
		//			sseconfig = make([]map[string]interface{}, 0)
		//			sseconfig = append(sseconfig, data)
		//		}
		//	}
		//	mapping["server_side_encryption_rule"] = sseconfig
		//
		//	//Add versioning information
		//	var versioning []map[string]interface{}
		//	if response.BucketInfo.Versioning != "" {
		//		data := map[string]interface{}{
		//			"status": response.BucketInfo.Versioning,
		//		}
		//		versioning = make([]map[string]interface{}, 0)
		//		versioning = append(versioning, data)
		//	}
		//	mapping["versioning"] = versioning
		//
		//} else {
		//	log.Printf("[WARN] Unable to get additional information for the bucket %s: %v", bucket.Name, err)
		//}
		//
		//// Add CORS rule information
		//var ruleMappings []map[string]interface{}
		//raw, err = client.WithOssClient(func(ossClient *oss.Client) (interface{}, error) {
		//	requestInfo = ossClient
		//	return ossClient.GetBucketCORS(bucket.Name)
		//})
		//if err == nil {
		//	if debugOn() {
		//		addDebug("GetBucketCORS", raw, requestInfo, map[string]string{"bucketName": bucket.Name})
		//	}
		//	cors, _ := raw.(oss.GetBucketCORSResult)
		//	if cors.CORSRules != nil {
		//		for _, rule := range cors.CORSRules {
		//			ruleMapping := make(map[string]interface{})
		//			ruleMapping["allowed_headers"] = rule.AllowedHeader
		//			ruleMapping["allowed_methods"] = rule.AllowedMethod
		//			ruleMapping["allowed_origins"] = rule.AllowedOrigin
		//			ruleMapping["expose_headers"] = rule.ExposeHeader
		//			ruleMapping["max_age_seconds"] = rule.MaxAgeSeconds
		//			ruleMappings = append(ruleMappings, ruleMapping)
		//		}
		//	}
		//} else if !IsExpectedErrors(err, []string{"NoSuchCORSConfiguration"}) {
		//	log.Printf("[WARN] Unable to get CORS information for the bucket %s: %v", bucket.Name, err)
		//}
		//mapping["cors_rules"] = ruleMappings
		//
		//// Add website configuration
		//var websiteMappings []map[string]interface{}
		//raw, err = client.WithOssClient(func(ossClient *oss.Client) (interface{}, error) {
		//	requestInfo = ossClient
		//	return ossClient.GetBucketWebsite(bucket.Name)
		//})
		//if err == nil {
		//	if debugOn() {
		//		addDebug("GetBucketWebsite", raw, requestInfo, map[string]string{"bucketName": bucket.Name})
		//	}
		//	ws, _ := raw.(oss.GetBucketWebsiteResult)
		//	websiteMapping := make(map[string]interface{})
		//	if v := &ws.IndexDocument; v != nil {
		//		websiteMapping["index_document"] = v.Suffix
		//	}
		//	if v := &ws.ErrorDocument; v != nil {
		//		websiteMapping["error_document"] = v.Key
		//	}
		//	websiteMappings = append(websiteMappings, websiteMapping)
		//} else if !IsExpectedErrors(err, []string{"NoSuchWebsiteConfiguration"}) {
		//	log.Printf("[WARN] Unable to get website information for the bucket %s: %v", bucket.Name, err)
		//}
		//mapping["website"] = websiteMappings
		//
		//// Add logging information
		//var loggingMappings []map[string]interface{}
		//raw, err = client.WithOssClient(func(ossClient *oss.Client) (interface{}, error) {
		//	return ossClient.GetBucketLogging(bucket.Name)
		//})
		//if err == nil {
		//	addDebug("GetBucketLogging", raw)
		//	logging, _ := raw.(oss.GetBucketLoggingResult)
		//	if logging.LoggingEnabled.TargetBucket != "" || logging.LoggingEnabled.TargetPrefix != "" {
		//		loggingMapping := map[string]interface{}{
		//			"target_bucket": logging.LoggingEnabled.TargetBucket,
		//			"target_prefix": logging.LoggingEnabled.TargetPrefix,
		//		}
		//		loggingMappings = append(loggingMappings, loggingMapping)
		//	}
		//} else {
		//	log.Printf("[WARN] Unable to get logging information for the bucket %s: %v", bucket.Name, err)
		//}
		//mapping["logging"] = loggingMappings
		//
		//// Add referer information
		//var refererMappings []map[string]interface{}
		//raw, err = client.WithOssClient(func(ossClient *oss.Client) (interface{}, error) {
		//	requestInfo = ossClient
		//	return ossClient.GetBucketReferer(bucket.Name)
		//})
		//if err == nil {
		//	if debugOn() {
		//		addDebug("GetBucketReferer", raw, requestInfo, map[string]string{"bucketName": bucket.Name})
		//	}
		//	referer, _ := raw.(oss.GetBucketRefererResult)
		//	refererMapping := map[string]interface{}{
		//		"allow_empty": referer.AllowEmptyReferer,
		//		"referers":    referer.RefererList,
		//	}
		//	refererMappings = append(refererMappings, refererMapping)
		//} else {
		//	log.Printf("[WARN] Unable to get referer information for the bucket %s: %v", bucket.Name, err)
		//}
		//mapping["referer_config"] = refererMappings
		//
		//// Add lifecycle information
		//var lifecycleRuleMappings []map[string]interface{}
		//raw, err = client.WithOssClient(func(ossClient *oss.Client) (interface{}, error) {
		//	requestInfo = ossClient
		//	return ossClient.GetBucketLifecycle(bucket.Name)
		//})
		//if err == nil {
		//	if debugOn() {
		//		addDebug("GetBucketLifecycle", raw, requestInfo, map[string]string{"bucketName": bucket.Name})
		//	}
		//	lifecycle, _ := raw.(oss.GetBucketLifecycleResult)
		//	if len(lifecycle.Rules) > 0 {
		//		for _, lifecycleRule := range lifecycle.Rules {
		//			ruleMapping := make(map[string]interface{})
		//			ruleMapping["id"] = lifecycleRule.ID
		//			ruleMapping["prefix"] = lifecycleRule.Prefix
		//			if LifecycleRuleStatus(lifecycleRule.Status) == ExpirationStatusEnabled {
		//				ruleMapping["enabled"] = true
		//			} else {
		//				ruleMapping["enabled"] = false
		//			}
		//
		//			// Expiration
		//			expirationMapping := make(map[string]interface{})
		//			if lifecycleRule.Expiration.Date != "" {
		//				t, err := time.Parse("2006-01-02T15:04:05.000Z", lifecycleRule.Expiration.Date)
		//				if err != nil {
		//					return WrapError(err)
		//				}
		//				expirationMapping["date"] = t.Format("2006-01-02")
		//			}
		//			if &lifecycleRule.Expiration.Days != nil {
		//				expirationMapping["days"] = int(lifecycleRule.Expiration.Days)
		//			}
		//			ruleMapping["expiration"] = []map[string]interface{}{expirationMapping}
		//			lifecycleRuleMappings = append(lifecycleRuleMappings, ruleMapping)
		//		}
		//	}
		//} else {
		//	log.Printf("[WARN] Unable to get lifecycle information for the bucket %s: %v", bucket.Name, err)
		//}
		//mapping["lifecycle_rule"] = lifecycleRuleMappings
		//
		//// Add policy information
		//var policy string
		//raw, err = client.WithOssClient(func(ossClient *oss.Client) (interface{}, error) {
		//	requestInfo = ossClient
		//	params := map[string]interface{}{}
		//	params["policy"] = nil
		//	return ossClient.Conn.Do("GET", bucket.Name, "", params, nil, nil, 0, nil)
		//})
		//
		//if err == nil {
		//	if debugOn() {
		//		addDebug("GetPolicyByConn", raw, requestInfo, map[string]string{"bucketName": bucket.Name})
		//	}
		//	rawResp := raw.(*oss.Response)
		//	rawData, err := ioutil.ReadAll(rawResp.Body)
		//	if err != nil {
		//		return WrapError(err)
		//	}
		//	policy = string(rawData)
		//} else {
		//	log.Printf("[WARN] Unable to get policy information for the bucket %s: %v", bucket.Name, err)
		//}
		//mapping["policy"] = policy
		//
		//// Add tags information
		//tagsMap := make(map[string]string)
		//raw, err = client.WithOssClient(func(ossClient *oss.Client) (interface{}, error) {
		//	requestInfo = ossClient
		//	return ossClient.GetBucketTagging(bucket.Name)
		//})
		//if err == nil {
		//	if debugOn() {
		//		addDebug("GetBucketTagging", raw, requestInfo, map[string]string{"bucketName": bucket.Name})
		//	}
		//	tagging, _ := raw.(oss.GetBucketTaggingResult)
		//	for _, t := range tagging.Tags {
		//		tagsMap[t.Key] = t.Value
		//	}
		//} else {
		//	log.Printf("[WARN] Unable to get tagging information for the bucket %s: %v", bucket.Name, err)
		//}
		//mapping["tags"] = tagsMap

		//ids = append(ids, bucket.Name)
		//s = append(s, mapping)
		//names = append(names, bucket.Name)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("buckets", s); err != nil {
		return WrapError(err)
	}
	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}
	// create a json file in current directory and write data source to it.
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
