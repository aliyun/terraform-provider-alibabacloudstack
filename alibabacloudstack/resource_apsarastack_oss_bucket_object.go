package alibabacloudstack

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mitchellh/go-homedir"
)

func resourceAlibabacloudStackOssBucketObject() *schema.Resource {
	resource := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"bucket": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"key": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"source": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"content"},
			},

			"content": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"source"},
			},

			"acl": {
				Type:         schema.TypeString,
				Default:      "private",
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"private", "public-read", "public-read-write"}, false),
			},

			"content_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"cache_control": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"content_disposition": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"content_encoding": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"content_md5": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"expires": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"server_side_encryption": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{string(ServerSideEncryptionKMS), string(ServerSideEncryptionAes256)}, false),
			},

			"kms_key_id": {
				Type:     schema.TypeString,
				Optional: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return ServerSideEncryptionKMS != d.Get("server_side_encryption").(string)
				},
				Computed: true,
			},

			"version_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
	setResourceFunc(resource, resourceAlibabacloudStackOssBucketObjectCreate,
		resourceAlibabacloudStackOssBucketObjectRead,
		resourceAlibabacloudStackOssBucketObjectUpdate,
		resourceAlibabacloudStackOssBucketObjectDelete)
	return resource
}

func resourceAlibabacloudStackOssBucketObjectCreate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceAlibabacloudStackOssBucketObjectUpdate(d *schema.ResourceData, meta interface{}) error {
	return resourceAlibabacloudStackOssBucketObjectPut(d, meta)
}

func resourceAlibabacloudStackOssBucketObjectPut(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	bucketName := d.Get("bucket").(string)
	ossService := OssService{client}
	ossClient, err := ossService.GetBucketClient(bucketName)
	if err != nil {
		return err
	}

	key := d.Get("key").(string)

	// Build PutObjectRequest
	putReq := &oss.PutObjectRequest{
		Bucket: &bucketName,
		Key:    &key,
	}

	// Set optional headers
	if v, ok := d.GetOk("acl"); ok {
		putReq.Acl = oss.ObjectACLType(v.(string))
	}
	if v, ok := d.GetOk("content_type"); ok {
		s := v.(string)
		putReq.ContentType = &s
	}
	if v, ok := d.GetOk("cache_control"); ok {
		s := v.(string)
		putReq.CacheControl = &s
	}
	if v, ok := d.GetOk("content_disposition"); ok {
		s := v.(string)
		putReq.ContentDisposition = &s
	}
	if v, ok := d.GetOk("content_encoding"); ok {
		s := v.(string)
		putReq.ContentEncoding = &s
	}
	if v, ok := d.GetOk("content_md5"); ok {
		s := v.(string)
		putReq.ContentMD5 = &s
	}
	if v, ok := d.GetOk("expires"); ok {
		s := v.(string)
		if _, err := time.Parse(time.RFC1123, s); err != nil {
			return fmt.Errorf("expires format must respect the RFC1123 standard (current value: %s)", s)
		}
		putReq.Expires = &s
	}
	if v, ok := d.GetOk("server_side_encryption"); ok {
		s := v.(string)
		putReq.ServerSideEncryption = &s
		if s == ServerSideEncryptionKMS {
			if v, ok := d.GetOk("kms_key_id"); ok {
				kmsKey := v.(string)
				putReq.SSEKMSKeyId = &kmsKey
			}
		}
	}

	// Set body
	if v, ok := d.GetOk("source"); ok {
		source := v.(string)
		path, err := homedir.Expand(source)
		if err != nil {
			return errmsgs.WrapError(err)
		}
		file, err := os.Open(path)
		if err != nil {
			return errmsgs.WrapError(err)
		}
		defer file.Close()
		putReq.Body = file
	} else if v, ok := d.GetOk("content"); ok {
		content := v.(string)
		putReq.Body = bytes.NewReader([]byte(content))
	} else {
		return errmsgs.WrapError(errmsgs.Error("[ERROR] Must specify \"source\" or \"content\" field"))
	}

	_, err = ossClient.PutObject(context.Background(), putReq)
	if err != nil {
		return errmsgs.WrapError(errmsgs.Error("Error putting object in Oss bucket (%s): %s", bucketName, err))
	}

	d.SetId(fmt.Sprintf("%s:%s", bucketName, key))
	return nil
}

func resourceAlibabacloudStackOssBucketObjectRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	var bucketName, key string
	if id_info := strings.SplitN(d.Id(), ":", 2); len(id_info) == 1 {
		// Compatible with old ID d.SetId(key)
		bucketName = d.Get("bucket").(string)
		key = d.Get("key").(string)
		d.SetId(fmt.Sprintf("%s:%s", bucketName, key))
	} else {
		bucketName = id_info[0]
		key = id_info[1]
	}
	ossService := OssService{client}
	ossClient, err := ossService.GetBucketClient(bucketName)
	if err != nil {
		return err
	}

	headReq := &oss.HeadObjectRequest{
		Bucket: &bucketName,
		Key:    &key,
	}
	object, err := ossClient.HeadObject(context.Background(), headReq)
	if err != nil {
		if errmsgs.IsExpectedErrors(err, "404 Not Found") {
			d.SetId("")
			return errmsgs.WrapError(errmsgs.Error("To get the Object: %#v but it is not exist in the specified bucket %s.", key, bucketName))
		}
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), "HeadObject", errmsgs.AlibabacloudStackLogGoSdkERROR)
	}
	addDebug("HeadObject", object, nil, map[string]interface{}{
		"objectKey": key,
	})

	// ACL - requires special permissions, may fail. Do not overwrite attribute when it fails.
	aclReq := &oss.GetObjectAclRequest{
		Bucket: &bucketName,
		Key:    &key,
	}
	if acl, err := ossClient.GetObjectAcl(context.Background(), aclReq); err == nil {
		d.Set("acl", acl.ACL)
	}

	d.Set("bucket", bucketName)
	d.Set("key", key)
	if object.ContentType != nil {
		d.Set("content_type", *object.ContentType)
	}
	if object.ContentMD5 != nil {
		d.Set("content_md5", *object.ContentMD5)
	}
	if object.ServerSideEncryption != nil {
		d.Set("server_side_encryption", *object.ServerSideEncryption)
		if *object.ServerSideEncryption == ServerSideEncryptionKMS && object.SSEKMSKeyId != nil {
			d.Set("kms_key_id", *object.SSEKMSKeyId)
		}
	}
	if object.ContentDisposition != nil {
		d.Set("content_disposition", *object.ContentDisposition)
	}
	if object.ContentEncoding != nil {
		d.Set("content_encoding", *object.ContentEncoding)
	}
	if object.Expires != nil {
		d.Set("expires", *object.Expires)
	}
	if object.VersionId != nil {
		d.Set("version_id", *object.VersionId)
	}

	return nil
}

func resourceAlibabacloudStackOssBucketObjectDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	ossService := OssService{client}

	bucketName := d.Get("bucket").(string)
	key := d.Get("key").(string)
	ossClient, err := ossService.GetBucketClient(bucketName)
	if err != nil {
		return err
	}

	delReq := &oss.DeleteObjectRequest{
		Bucket: &bucketName,
		Key:    &key,
	}
	_, err = ossClient.DeleteObject(context.Background(), delReq)
	if err != nil {
		if errmsgs.IsExpectedErrors(err, "No Content", "Not Found") {
			return nil
		}
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), "DeleteObject", errmsgs.AlibabacloudStackLogGoSdkERROR)
	}

	return errmsgs.WrapError(ossService.WaitForOssBucketObject(bucketName, d.Id(), Deleted, DefaultTimeoutMedium))
}
