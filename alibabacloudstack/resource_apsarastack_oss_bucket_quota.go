package alibabacloudstack

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/PaesslerAG/jsonpath"
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
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_oss_bucket", "IsBucketExist", errmsgs.AlibabacloudStackLogGoSdkERROR)
	}
	quota := d.Get("quota").(int)

	if det.Name == bucketName {
		request := client.NewCommonRequest("GET", "OneRouter", "2018-12-12", "DoOpenApi", "")
		request.QueryParams["OpenApiAction"] = "SetBucketStorageCapacity"
		request.QueryParams["ProductName"] = "oss"
		request.QueryParams["Params"] = fmt.Sprintf("{\"%s\":\"%s\",\"%s\":%d}", "BucketName", bucketName, "StorageCapacity", quota)
		request.QueryParams["Content"] = fmt.Sprintf("%s%d%s", "<BucketUserQos><StorageCapacity>", quota, "</StorageCapacity></BucketUserQos>")
		bresponse, err := client.ProcessCommonRequest(request)
		addDebug(request.GetActionName(), bresponse, request, request.QueryParams)
		log.Printf(" response of raw SetBucketStorageCapacity : %s", bresponse)
		if err != nil {
			if bresponse == nil {
				return errmsgs.WrapErrorf(err, "Process Common Request Failed")
			}
			errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_oss_bucket", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
		log.Printf("Enter for logging")
	}
	d.SetId(bucketName)

	return nil
}

func resourceAlibabacloudStackOssBucketQuotaRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	ossService := OssService{client}
	bucketName := d.Id()
	_, err := ossService.DescribeOssBucket(bucketName)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_oss_bucket", "IsBucketExist", errmsgs.AlibabacloudStackLogGoSdkERROR)
	}
	request := client.NewCommonRequest("GET", "OneRouter", "2018-12-12", "DoOpenApi", "")
	request.QueryParams["OpenApiAction"] = "GetBucketStorageCapacity"
	request.QueryParams["ProductName"] = "oss"
	request.QueryParams["Params"] = fmt.Sprintf("{\"%s\":\"%s\"}", "BucketName", bucketName)

	bresponse, err := client.ProcessCommonRequest(request)
	addDebug(request.GetActionName(), bresponse, request, request.QueryParams)
	log.Printf(" response of raw GetBucketStorageCapacity : %s", bresponse)
	if err != nil {
		if bresponse == nil {
			return errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_oss_bucket", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	resp := make(map[string]interface{})
	json.Unmarshal(bresponse.GetHttpContentBytes(), &resp)

	storageCapacity, err := jsonpath.Get("$.Data.BucketUserQos.StorageCapacity", resp)
	if err != nil {
		return errmsgs.WrapErrorf(err, "Process Common Request Failed")
	}
	quota, err := strconv.Atoi(storageCapacity.(string))
	if err != nil {
		return errmsgs.WrapErrorf(err, "quota value type error")
	}
	d.Set("quota", quota)
	d.Set("bucket", bucketName)

	return nil
}

func resourceAlibabacloudStackOssBucketQuotaDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	ossService := OssService{client}
	bucketName := d.Get("bucket").(string)
	_, err := ossService.DescribeOssBucket(bucketName)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_oss_bucket", "IsBucketExist", errmsgs.AlibabacloudStackLogGoSdkERROR)
	}

	request := client.NewCommonRequest("GET", "OneRouter", "2018-12-12", "DoOpenApi", "")
	request.QueryParams["OpenApiAction"] = "SetBucketStorageCapacity"
	request.QueryParams["ProductName"] = "oss"
	request.QueryParams["Params"] = fmt.Sprintf("{\"%s\":\"%s\",\"%s\":%d}", "BucketName", bucketName, "StorageCapacity", -1)
	request.QueryParams["Content"] = fmt.Sprintf("%s%d%s", "<BucketUserQos><StorageCapacity>", -1, "</StorageCapacity></BucketUserQos>")

	bresponse, err := client.ProcessCommonRequest(request)
	addDebug(request.GetActionName(), bresponse, request, request.QueryParams)
	log.Printf(" response of raw SetBucketStorageCapacity : %s", bresponse)
	if err != nil {
		if bresponse == nil {
			return errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_oss_bucket", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	log.Printf("Enter for logging")
	d.SetId(bucketName)

	return nil
}
