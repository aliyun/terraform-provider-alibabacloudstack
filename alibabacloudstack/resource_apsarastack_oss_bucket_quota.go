package alibabacloudstack

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"strings"

	oss "github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/signer"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackOssBucketQuota() *schema.Resource {
	resource := &schema.Resource{
		DeprecationMessage: "oss_bucket already includes corresponding functions and is scheduled for removal in version 3.21.0",
		Schema: map[string]*schema.Schema{
			"bucket": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"quota": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
		},
	}
	setResourceFunc(resource, resourceAlibabacloudStackOssBucketQuotaCreate, resourceAlibabacloudStackOssBucketQuotaRead, nil, resourceAlibabacloudStackOssBucketQuotaDelete)
	return resource
}

func resourceAlibabacloudStackOssBucketQuotaCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	ossService := OssService{client}
	bucketName := d.Get("bucket").(string)
	det, err := ossService.DescribeOssBucket(bucketName)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_oss_bucket", "IsBucketExist", errmsgs.AlibabacloudStackOssGoSdk)
	}
	quota := d.Get("quota").(int)

	if det != nil && *det.Name == bucketName {
		ossClient, err := ossService.GetBucketClient(bucketName)
		if err != nil {
			return errmsgs.WrapError(err)
		}

		xmlBody := fmt.Sprintf("<BucketUserQos><StorageCapacity>%d</StorageCapacity></BucketUserQos>", quota)
		input := &oss.OperationInput{
			OpName: "SetBucketStorageCapacity",
			Method: "PUT",
			Bucket: oss.Ptr(bucketName),
			Parameters: map[string]string{
				"qos": "",
			},
			Headers: map[string]string{
				"Content-Type": "application/xml",
			},
			Body: strings.NewReader(xmlBody),
		}
		input.OpMetadata.Set(signer.SubResource, []string{"qos"})
		output, err := ossClient.InvokeOperation(context.Background(), input)
		if err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, bucketName, "SetBucketStorageCapacity", errmsgs.AlibabacloudStackOssGoSdk)
		}
		if output.Body != nil {
			output.Body.Close()
		}
		addDebug("SetBucketStorageCapacity", output, input, map[string]interface{}{"bucket": bucketName, "quota": quota})
		log.Printf(" response of SetBucketStorageCapacity for bucket %s", bucketName)
		log.Printf("Enter for logging")
	}
	d.SetId(bucketName)

	return nil
}

// bucketUserQosXML is used to parse the GetBucketStorageCapacity response
type bucketUserQosXML struct {
	StorageCapacity int64 `xml:"StorageCapacity"`
}

func resourceAlibabacloudStackOssBucketQuotaRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	ossService := OssService{client}
	bucketName := d.Id()
	_, err := ossService.DescribeOssBucket(bucketName)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_oss_bucket", "IsBucketExist", errmsgs.AlibabacloudStackOssGoSdk)
	}

	ossClient, err := ossService.GetBucketClient(bucketName)
	if err != nil {
		return errmsgs.WrapError(err)
	}

	input := &oss.OperationInput{
		OpName: "GetBucketStorageCapacity",
		Method: "GET",
		Bucket: oss.Ptr(bucketName),
		Parameters: map[string]string{
			"qos": "",
		},
	}
	input.OpMetadata.Set(signer.SubResource, []string{"qos"})
	output, err := ossClient.InvokeOperation(context.Background(), input)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, bucketName, "GetBucketStorageCapacity", errmsgs.AlibabacloudStackOssGoSdk)
	}
	addDebug("GetBucketStorageCapacity", output, input, map[string]interface{}{"bucket": bucketName})
	log.Printf(" response of GetBucketStorageCapacity for bucket %s", bucketName)

	defer output.Body.Close()
	responseBody, err := io.ReadAll(output.Body)
	if err != nil {
		return errmsgs.WrapErrorf(err, "read GetBucketStorageCapacity response failed")
	}

	var qos bucketUserQosXML
	if err := xml.Unmarshal(responseBody, &qos); err != nil {
		return errmsgs.WrapErrorf(err, "unmarshal GetBucketStorageCapacity response failed")
	}

	d.Set("quota", int(qos.StorageCapacity))
	d.Set("bucket", bucketName)

	return nil
}

func resourceAlibabacloudStackOssBucketQuotaDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	ossService := OssService{client}
	bucketName := d.Get("bucket").(string)
	_, err := ossService.DescribeOssBucket(bucketName)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_oss_bucket", "IsBucketExist", errmsgs.AlibabacloudStackOssGoSdk)
	}
	ossClient, err := ossService.GetBucketClient(bucketName)
	if err != nil {
		return errmsgs.WrapError(err)
	}

	xmlBody := fmt.Sprintf("<BucketUserQos><StorageCapacity>%d</StorageCapacity></BucketUserQos>", -1)
	input := &oss.OperationInput{
		OpName: "SetBucketStorageCapacity",
		Method: "PUT",
		Bucket: oss.Ptr(bucketName),
		Parameters: map[string]string{
			"qos": "",
		},
		Headers: map[string]string{
			"Content-Type": "application/xml",
		},
		Body: strings.NewReader(xmlBody),
	}
	input.OpMetadata.Set(signer.SubResource, []string{"qos"})
	output, err := ossClient.InvokeOperation(context.Background(), input)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, bucketName, "SetBucketStorageCapacity", errmsgs.AlibabacloudStackOssGoSdk)
	}
	if output.Body != nil {
		output.Body.Close()
	}
	addDebug("SetBucketStorageCapacity", output, input, map[string]interface{}{"bucket": bucketName, "quota": -1})
	log.Printf(" response of SetBucketStorageCapacity for bucket %s", bucketName)
	log.Printf("Enter for logging")
	d.SetId(bucketName)

	return nil
}
