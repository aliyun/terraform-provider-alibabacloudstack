package apsarastack

import (
	"encoding/json"
	"log"
	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/aliyun/terraform-provider-alibabaCloudStack/apsarastack/connectivity"
)

type BucketVpcService struct {
	client *connectivity.ApsaraStackClient
}

type VpcListResult struct {
	Api             string        `json:"api"`
	AsapiRequestId  string        `json:"asapiRequestId"`
	AsapiSuccess    bool          `json:"asapiSuccess"`
	HttpOk          bool          `json:"httpOk"`
	Success         bool          `json:"success"`
	Code            int64         `json:"code"`
	Domain          string        `json:"domain"`
	Message         string        `json:"message"`
	ServerRole      string        `json:"serverRole"`
	EagleEyeTraceId string        `json:"eagleEyeTraceId"`
	VpcList         []interface{} `json:"data"`
	PageModel       interface{}   `json:"pageModel"`
}

func (s *BucketVpcService) BucketVpcList(bucketName string) (vpclist *VpcListResult, err error) {
	var requestInfo *oss.Client
	request := requests.NewCommonRequest()
	if s.client.Config.Insecure {
		request.SetHTTPSInsecure(s.client.Config.Insecure)
	}
	request.QueryParams = map[string]string{
		"Department":       s.client.Department,
		"ResourceGroup":    s.client.ResourceGroup,
		"RegionId":         s.client.RegionId,
		"Action":           "ListBucketVpc",
		"AccountInfo":      "123456",
		"Version":          "2019-05-10",
		"SignatureVersion": "2.1",
		"OpenApiAction":    "ListBucketVpc",
		"Product":          "Ascm",
		"BucketName":       bucketName,
		// "Params":           fmt.Sprintf("{\"%s\":%s,\"%s\":%s,\"%s\":\"%s\",\"%s\":\"%s\",}", "action", "ListBucketVpc", "product", "Ascm", "region", s.client.RegionId, "params", "{\"BucketName\":\""+bucketName+"\"}"),
	}
	request.Method = "POST"  // Set request method
	request.Product = "Ascm" // Specify product
	if strings.ToLower(s.client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	} // Set request scheme. Default: http
	request.ApiName = "ListBucketVpc"
	request.Version = "2019-05-10" // Specify product version
	request.Headers = map[string]string{"RegionId": s.client.RegionId}
	raw, err := s.client.WithEcsClient(func(ossClient *ecs.Client) (interface{}, error) {

		return ossClient.ProcessCommonRequest(request)
	})
	log.Printf("Response of ListBucketVpc: %s", raw)
	log.Printf("Bresponse ListBucketVpc before error")
	if err != nil {
		if ossNotFoundError(err) {
			return vpclist, WrapErrorf(err, NotFoundMsg, ApsaraStackOssGoSdk)
		}
		return vpclist, WrapErrorf(err, DefaultErrorMsg, bucketName, "ListBucketVpc", ApsaraStackOssGoSdk)
	}
	log.Printf("Bresponse ListBucketVpc after error")
	addDebug("ListBucketVpc", raw, requestInfo, request)
	bresponse, _ := raw.(*responses.CommonResponse)
	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &vpclist)
	if err != nil {
		return vpclist, WrapError(err)
	}
	if !vpclist.Success {
		return vpclist, WrapError(err)
	}
	return vpclist, nil
}

func (s *BucketVpcService) BindBucket(vpcId string, vpcName string, vLan string, bucket string) error {
	var requestInfo *oss.Client
	request := requests.NewCommonRequest()
	if s.client.Config.Insecure {
		request.SetHTTPSInsecure(s.client.Config.Insecure)
	}
	request.QueryParams = map[string]string{
		"AccessKeySecret":  s.client.SecretKey,
		"AccessKeyId":      s.client.AccessKey,
		"Product":          "ascm",
		"Department":       s.client.Department,
		"ResourceGroup":    s.client.ResourceGroup,
		"RegionId":         s.client.RegionId,
		"Version":          "2019-05-10",
		"SignatureVersion": "2.1",
		"Action":           "BindBucketPolicy",
		"BucketName":       bucket,
		"VpcName":          vpcName,
		"VLan":             vLan,
		"VpcId":            vpcId,
		// "Params":           fmt.Sprintf("{\"%s\":%s,\"%s\":%s,\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":\"%s\",}", "BucketName", bucket, "VpcName", vpcName, "VLan", vLan, "VpcId", vpcId, "AccessKeyId", s.client.AccessKey, "AccessKeySecret", s.client.SecretKey, "RegionId", s.client.RegionId),
	}
	request.Method = "POST"        // Set request method
	request.Product = "ascm"       // Specify product
	request.Version = "2019-05-10" // Specify product version
	if strings.ToLower(s.client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	} // Set request scheme. Default: http
	request.ApiName = "BindBucketPolicy"
	request.Headers = map[string]string{"RegionId": s.client.RegionId}

	raw, err := s.client.WithEcsClient(func(ossClient *ecs.Client) (interface{}, error) {

		return ossClient.ProcessCommonRequest(request)
	})
	log.Printf("Response of BindBucketPolicy: %s", raw)
	log.Printf("Bresponse BindBucketPolicy before error")
	if err != nil {
		if ossNotFoundError(err) {
			return WrapErrorf(err, NotFoundMsg, ApsaraStackOssGoSdk)
		}
		return WrapErrorf(err, DefaultErrorMsg, bucket, "BindBucketPolicy", ApsaraStackOssGoSdk)
	}
	log.Printf("Bresponse BindBucketPolicy after error")
	addDebug("CreateBucketInfo", raw, requestInfo, request)
	log.Printf("Bresponse BindBucketPolicy check")
	bresponse, _ := raw.(*responses.CommonResponse)
	log.Printf("Bresponse BindBucketPolicy %s", bresponse)

	if bresponse.GetHttpStatus() != 200 {
		return WrapErrorf(err, DefaultErrorMsg, "apsarastack_oss_bucket_vpc", "BindBucketPolicy", ApsaraStackOssGoSdk)
	}
	return nil
}

func (s *BucketVpcService) UnBindBucket(vpcId string, bucket string) error {
	var requestInfo *oss.Client
	request := requests.NewCommonRequest()
	if s.client.Config.Insecure {
		request.SetHTTPSInsecure(s.client.Config.Insecure)
	}
	request.QueryParams = map[string]string{
		"AccessKeySecret":  s.client.SecretKey,
		"Product":          "Ascm",
		"Department":       s.client.Department,
		"ResourceGroup":    s.client.ResourceGroup,
		"RegionId":         s.client.RegionId,
		"Action":           "UnBindBucketPolicy",
		"Version":          "2019-05-10",
		"SignatureVersion": "2.1",
		"BucketName":       bucket,
		"vpcId":            vpcId,
		// "Params":           fmt.Sprintf("{\"%s\":%s,\"%s\":%s,\"%s\":\"%s\",\"%s\":\"%s\",}", "action", "UnBindBucketPolicy", "product", "Ascm", "region", s.client.Region, "params", "{\"BucketName\":"+bucket+",\"endpoint\":"+""+",\"vpcId\":"+vpcId+"\",}"),
	}
	request.Method = "POST"                // Set request method
	request.Product = "Ascm"               // Specify product
	request.Version = "2019-05-10"         // Specify product version
	request.ApiName = "UnBindBucketPolicy" // Specify ApiName
	if strings.ToLower(s.client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	} // Set request scheme. Default: http
	request.Headers = map[string]string{"RegionId": s.client.RegionId}

	raw, err := s.client.WithEcsClient(func(ossClient *ecs.Client) (interface{}, error) {

		return ossClient.ProcessCommonRequest(request)
	})
	log.Printf("Response of UnBindBucketPolicy: %s", raw)
	log.Printf("Bresponse UnBindBucketPolicy before error")
	if err != nil {
		if ossNotFoundError(err) {
			return WrapErrorf(err, NotFoundMsg, ApsaraStackOssGoSdk)
		}
		return WrapErrorf(err, DefaultErrorMsg, bucket, "UnBindBucketPolicy", ApsaraStackOssGoSdk)
	}
	log.Printf("Bresponse UnBindBucketPolicy after error")
	addDebug("CreateBucketInfo", raw, requestInfo, request)
	log.Printf("Bresponse UnBindBucketPolicy check")
	bresponse, _ := raw.(*responses.CommonResponse)
	log.Printf("Bresponse UnBindBucketPolicy %s", bresponse)

	if bresponse.GetHttpStatus() != 200 {
		return WrapErrorf(err, DefaultErrorMsg, "apsarastack_oss_bucket_vpc", "UnBindBucketPolicy", ApsaraStackOssGoSdk)
	}
	return nil
}
