package alibabacloudstack

import (
	"context"
	"log"
	"regexp"
	"time"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceAlibabacloudStackOssBucketObjects() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackOssBucketObjectsRead,

		Schema: map[string]*schema.Schema{
			"bucket_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"oss_cluster": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"key_prefix": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"key_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
				ForceNew:     true,
			},
			"output_file": {
				Type:       schema.TypeString,
				Optional:   true,
				Deprecated: "The 'output_file' field has been deprecated and is scheduled for removal in version 3.19.0. To write content to a file, use the 'local_file' provider instead.",
			},

			// Computed values
			"objects": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"acl": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"content_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cache_control": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"content_disposition": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"content_encoding": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"content_md5": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"expires": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"server_side_encryption": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"sse_kms_key_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"storage_class": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"last_modification_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlibabacloudStackOssBucketObjectsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	bucketName := d.Get("bucket_name").(string)
	ossService := OssService{client}
	ossClient, err := ossService.GetOssClientForCluster(d.Get("oss_cluster").(string))
	if err != nil {
		return err
	}

	// List bucket objects
	var allObjects []oss.ObjectProperties
	nextMarker := ""
	prefix := ""
	if v, ok := d.GetOk("key_prefix"); ok && v.(string) != "" {
		prefix = v.(string)
	}
	for {
		input := &oss.ListObjectsRequest{
			Bucket: &bucketName,
			Prefix: &prefix,
		}
		if nextMarker != "" {
			input.Marker = &nextMarker
		}

		response, err := ossClient.ListObjects(context.Background(), input)
		if err != nil {
			return err
		}
		if len(response.Contents) < 1 {
			break
		}

		allObjects = append(allObjects, response.Contents...)

		if response.NextMarker == nil || *response.NextMarker == "" {
			break
		}
		nextMarker = *response.NextMarker
	}

	var filteredObjectsTemp []oss.ObjectProperties
	keyRegex, ok := d.GetOk("key_regex")
	if ok && keyRegex.(string) != "" {
		var r *regexp.Regexp
		if keyRegex != "" {
			r = regexp.MustCompile(keyRegex.(string))
		}
		for _, object := range allObjects {
			if r != nil && !r.MatchString(*object.Key) {
				continue
			}
			filteredObjectsTemp = append(filteredObjectsTemp, object)
		}
	} else {
		filteredObjectsTemp = allObjects
	}

	return bucketObjectsDescriptionAttributes(d, bucketName, filteredObjectsTemp, meta)
}

func bucketObjectsDescriptionAttributes(d *schema.ResourceData, bucketName string, objects []oss.ObjectProperties, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	var ids []string
	var s []map[string]interface{}
	ossService := OssService{client}
	ossClient, err := ossService.GetOssClientForCluster(d.Get("oss_cluster").(string))
	if err != nil {
		return err
	}
	for _, object := range objects {
		key := ""
		if object.Key != nil {
			key = *object.Key
		}
		storageClass := ""
		if object.StorageClass != nil {
			storageClass = *object.StorageClass
		}
		lastModified := ""
		if object.LastModified != nil {
			lastModified = object.LastModified.Format(time.RFC3339)
		}
		mapping := map[string]interface{}{
			"key":                    key,
			"storage_class":          storageClass,
			"last_modification_time": lastModified,
		}

		// Add metadata information
		headReq := &oss.HeadObjectRequest{
			Bucket: &bucketName,
			Key:    object.Key,
		}
		objectHeader, err := ossClient.HeadObject(context.Background(), headReq)
		if err != nil {
			log.Printf("[ERROR] Unable to get metadata for the object %s: %v", key, err)
		} else {
			if objectHeader.ContentType != nil {
				mapping["content_type"] = *objectHeader.ContentType
			}
			if objectHeader.CacheControl != nil {
				mapping["cache_control"] = *objectHeader.CacheControl
			}
			if objectHeader.ContentDisposition != nil {
				mapping["content_disposition"] = *objectHeader.ContentDisposition
			}
			if objectHeader.ContentEncoding != nil {
				mapping["content_encoding"] = *objectHeader.ContentEncoding
			}
			if objectHeader.ContentMD5 != nil {
				mapping["content_md5"] = *objectHeader.ContentMD5
			}
			if objectHeader.Expires != nil {
				mapping["expires"] = *objectHeader.Expires
			}
			if objectHeader.ServerSideEncryption != nil {
				mapping["server_side_encryption"] = *objectHeader.ServerSideEncryption
			}
			if objectHeader.SSEKMSKeyId != nil {
				mapping["sse_kms_key_id"] = *objectHeader.SSEKMSKeyId
			}
		}
		// Add ACL information
		aclReq := &oss.GetObjectAclRequest{
			Bucket: &bucketName,
			Key:    object.Key,
		}
		objectACL, err := ossClient.GetObjectAcl(context.Background(), aclReq)
		if err != nil {
			log.Printf("[ERROR] Unable to get ACL for the object %s: %v", key, err)
		} else if objectACL.ACL != nil {
			mapping["acl"] = *objectACL.ACL
		}

		ids = append(ids, key)
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("objects", s); err != nil {
		return errmsgs.WrapError(err)
	}
	// create a json file in current directory and write data source to it.
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		if err := writeToFile(output.(string), s); err != nil {
			return err
		}
	}
	return nil
}
