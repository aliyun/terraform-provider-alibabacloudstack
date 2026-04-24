package alibabacloudstack

import (
	"context"
	"regexp"
	"slices"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
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
			"shared": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Description: "Whether to query resources shared from other organizations. If set to true, shared resources will be included in the results.",
			},
			"output_file": {
				Type:       schema.TypeString,
				Optional:   true,
				Deprecated: "The 'output_file' field has been deprecated and is scheduled for removal in version 3.19.0. To write content to a file, use the 'local_file' provider instead.",
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
	ossService := OssService{client}
	endpointMap, err := ossService.GetBucketEndpointMap()
	if err != nil || len(endpointMap) < 1 {
		endpoint, err := ossService.GetDefaultOssEndpoint()
		if err != nil {
			return errmsgs.WrapError(err)
		}
		endpointMap = map[string]string{
			"defaultCluster": endpoint,
		}
	}
	var buckets []oss.BucketProperties
	names := make([]string, 0)
	for _, endpoint := range endpointMap {
		ossclietn, err := ossService.GetOssClient(endpoint)
		if err != nil {
			return errmsgs.WrapError(err)
		}
		for {
			request := &oss.ListBucketsRequest{}
			lsRes, err := ossclietn.ListBuckets(context.TODO(), request)
			if err != nil {
				if errmsgs.IsHostNotFound(err) || errmsgs.IsExpectedErrors(err, "With part response body Bad Gateway: [Errno 61] Connection refused.") {
					break
				}
				return errmsgs.WrapError(err)
			}
			for _, bucket := range lsRes.Buckets {
				if !slices.Contains(names, *bucket.Name) {
					buckets = append(buckets, bucket)
					names = append(names, *bucket.Name)
				}
			}

			if !lsRes.IsTruncated {
				break
			}
			request.Marker = lsRes.NextMarker
		}
	}

	if len(buckets) == 0 {
		d.SetId(dataResourceIdHash([]string{}))
		return nil
	}
	var filteredBucketsTemp []oss.BucketProperties
	idsMap := getIdsStringFilter(d)
	nameRegex := d.Get("name_regex")
	var r *regexp.Regexp
	if nameRegex != "" {
		r = regexp.MustCompile(nameRegex.(string))
	}
	for _, bucket := range buckets {
		if bucket.Name == nil {
			continue
		}
		if len(idsMap) > 0 {
			if _, ok := idsMap[*bucket.Name]; !ok {
				continue
			}
		}
		if r != nil && !r.MatchString(*bucket.Name) {
			continue
		}
		filteredBucketsTemp = append(filteredBucketsTemp, bucket)
	}
	return bucketsDescriptionAttributes(d, filteredBucketsTemp, meta)
}

func bucketsDescriptionAttributes(d *schema.ResourceData, buckets []oss.BucketProperties, meta interface{}) error {
	var ids []string
	var s []map[string]interface{}
	var names []string

	// Sort buckets by Name
	slices.SortFunc(buckets, func(a, b oss.BucketProperties) int {
		if a.Name == nil {
			return 1
		}
		if b.Name == nil {
			return -1
		}
		if *a.Name < *b.Name {
			return -1
		}
		if *a.Name > *b.Name {
			return 1
		}
		return 0
	})

	for _, bucket := range buckets {
		creationDate, _ := bucket.CreationDate.MarshalText()
		creationDateStr := string(creationDate)
		mapping := map[string]interface{}{
			"id":            bucket.Name,
			"name":          bucket.Name,
			"location":      bucket.Location,
			"storage_class": bucket.StorageClass,
			"creation_date": creationDateStr,
			// "extranet_endpoint": bucket.Extranetendpoint,
			// "intranet_endpoint": bucket.Intranetendpoint,
		}
		ids = append(ids, *bucket.Name)
		s = append(s, mapping)
		names = append(names, *bucket.Name)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("buckets", s); err != nil {
		return errmsgs.WrapError(err)
	}
	if err := d.Set("names", names); err != nil {
		return errmsgs.WrapError(err)
	}
	if err := d.Set("ids", ids); err != nil {
		return errmsgs.WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		if err := writeToFile(output.(string), s); err != nil {
			return err
		}
	}
	return nil
}
