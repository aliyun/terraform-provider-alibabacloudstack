package alibabacloudstack

import (
	"encoding/json"
	"encoding/xml"
	"regexp"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceAlibabacloudStackOssBuckets() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackOssBucketsRead,

		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
				ForceNew: true,
				MinItems: 1,
			},
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
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
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

		request := client.NewCommonRequest("POST", "OneRouter", "2018-12-12", "DoOpenApi", "")
		mergeMaps(request.QueryParams, map[string]string{
			"AccountInfo":      "123456",
			"SignatureVersion": "1.0",
			"OpenApiAction":    "GetService",
			"ProductName":      "oss",
		})

		var bucketList = &BucketList{}
		raw, err := client.WithOssNewClient(func(ossClient *ecs.Client) (interface{}, error) {
			return ossClient.ProcessCommonRequest(request)
		})

		bresponse, ok := raw.(*responses.CommonResponse)
		if err != nil {
			errmsg := ""
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			}
			if ossNotFoundError(err) {
				return errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackLogGoSdkERROR)
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "GetBucketInfo", errmsgs.AlibabacloudStackLogGoSdkERROR, errmsg)
		}
		addDebug("GetBucketInfo", raw, requestInfo, request)

		err = json.Unmarshal(bresponse.GetHttpContentBytes(), bucketList)
		if err != nil {
			return errmsgs.WrapError(err)
		}
		if bucketList.Code != "200" || len(bucketList.Data.ListAllMyBucketsResult.Buckets.Bucket) < 1 {
			return errmsgs.WrapError(err)
		}

		for _, k := range bucketList.Data.ListAllMyBucketsResult.Buckets.Bucket {
			allBuckets = append(allBuckets, oss.BucketProperties{
				XMLName:          xml.Name{},
				Name:             k.Name,
				Location:         k.Location,
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
	var ids []string
	var s []map[string]interface{}
	var names []string
	for _, bucket := range buckets {
		mapping := map[string]interface{}{
			"id":                bucket.Name,
			"name":              bucket.Name,
			"location":          bucket.Location,
			"storage_class":     bucket.StorageClass,
			"creation_date":     bucket.CreationDate,
			"extranet_endpoint": bucket.Extranetendpoint,
			"intranet_endpoint": bucket.Intranetendpoint,
		}
		ids = append(ids, bucket.Name)
		s = append(s, mapping)
		names = append(names, bucket.Name)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("buckets", s); err != nil {
		return errmsgs.WrapError(err)
	}
	if err := d.Set("names", names); err != nil {
		return errmsgs.WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		if err := writeToFile(output.(string), s); err != nil {
			return err
		}
	}
	return nil
}
