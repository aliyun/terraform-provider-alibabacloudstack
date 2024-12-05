package alibabacloudstack

import (
	"encoding/json"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"strconv"
	"strings"
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
)

// OssService *connectivity.AlibabacloudStackClient
type OssService struct {
	client *connectivity.AlibabacloudStackClient
}

func (s *OssService) DescribeOssBucket(id string) (response oss.GetBucketInfoResult, err error) {
	//request := map[string]string{"bucketName": id, "Department": s.client.Department, "ResourceGroup": s.client.ResourceGroup}
	var requestInfo *oss.Client

	request := requests.NewCommonRequest()
	if s.client.Config.Insecure {
		request.SetHTTPSInsecure(s.client.Config.Insecure)
	}
	request.QueryParams = map[string]string{

		
		"Product":          "OneRouter",
		"Department":       s.client.Department,
		"ResourceGroup":    s.client.ResourceGroup,
		"RegionId":         s.client.RegionId,
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
	if strings.ToLower(s.client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	} // Set request scheme. Default: http
	request.ApiName = "DoOpenApi"
	request.Headers = map[string]string{"RegionId": s.client.RegionId}

	var bucketList = &BucketList{}
	raw, err := s.client.WithOssNewClient(func(ossClient *ecs.Client) (interface{}, error) {

		return ossClient.ProcessCommonRequest(request)
	})

	if err != nil {
		if ossNotFoundError(err) {
			return response, WrapErrorf(err, NotFoundMsg, AlibabacloudStackOssGoSdk)
		}
		return response, WrapErrorf(err, DefaultErrorMsg, id, "GetBucketInfo", AlibabacloudStackOssGoSdk)
	}
	addDebug("GetBucketInfo", raw, requestInfo, request)
	bresponse, _ := raw.(*responses.CommonResponse)

	err = json.Unmarshal(bresponse.GetHttpContentBytes(), bucketList)
	if err != nil {
		return response, WrapError(err)
	}
	if bucketList.Code != "200" || len(bucketList.Data.ListAllMyBucketsResult.Buckets.Bucket) < 1 {
		return response, WrapError(err)
	}

	var found = false
	for _, j := range bucketList.Data.ListAllMyBucketsResult.Buckets.Bucket {
		if j.Name == id {
			response.BucketInfo.Name = j.Name
			response.BucketInfo.StorageClass = j.StorageClass
			response.BucketInfo.ExtranetEndpoint = j.ExtranetEndpoint
			response.BucketInfo.IntranetEndpoint = j.IntranetEndpoint
			response.BucketInfo.Owner.ID = fmt.Sprint(j.ResourceGroupName)
			//response.BucketInfo.CreationDate=fmt.Sprint(j.CreationDate.
			response.BucketInfo.Location = j.Location
			found = true
			break
		}
	}
	if !found {
		response.BucketInfo.Name = ""
	}
	return
}

type BucketList struct {
	Data struct {
		ListAllMyBucketsResult struct {
			Buckets struct {
				Bucket []struct {
					Comment           string `json:"Comment"`
					CreationDate      string `json:"CreationDate"`
					Department        int64  `json:"Department"`
					DepartmentName    string `json:"DepartmentName"`
					ExtranetEndpoint  string `json:"ExtranetEndpoint"`
					IntranetEndpoint  string `json:"IntranetEndpoint"`
					Location          string `json:"Location"`
					Name              string `json:"Name"`
					ResourceGroup     int64  `json:"ResourceGroup"`
					ResourceGroupName string `json:"ResourceGroupName"`
					StorageClass      string `json:"StorageClass"`
				} `json:"Bucket"`
			} `json:"Buckets"`
			Owner struct{} `json:"Owner"`
		} `json:"ListAllMyBucketsResult"`
	} `json:"Data"`
	Code         string `json:"code"`
	Cost         int64  `json:"cost"`
	Message      string `json:"message"`
	PureListData bool   `json:"pureListData"`
	Redirect     bool   `json:"redirect"`
	Success      bool   `json:"success"`
}

func (s *OssService) WaitForOssBucket(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeOssBucket(id)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return WrapError(err)
			}
		}

		if object.BucketInfo.Name != "" && status != Deleted {
			return nil
		}
		if object.BucketInfo.Name == "" && status == Deleted {
			return nil
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, object.BucketInfo.Name, status, ProviderERROR)
		}
	}
}

func (s *OssService) HeadOssBucketObject(bucketName string, objectName string) error {
	client := s.client
	var requestInfo *oss.Client
	request := requests.NewCommonRequest()
	if client.Config.Insecure {
		request.SetHTTPSInsecure(client.Config.Insecure)
	}
	request.QueryParams = map[string]string{
		"RegionId":        client.RegionId,
		
		
		"Product":         "OneRouter",
		"Action":          "DoApi",
		"AppAction":       "HeadObject",
		"AppName":         "one-console-app-oss",
		"Version":         "2018-12-12",
		"Params":          "{\"region\":\"" + client.RegionId + "\",\"params\":{\"bucketName\":\"" + bucketName + "\",\"objectName\":\"" + objectName + "\"}}",
		"AccountInfo":     "",
	}
	request.Method = "POST"
	request.Product = "OneRouter"
	request.Version = "2018-12-12"
	request.Domain = client.Domain
	if strings.ToLower(client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.ApiName = "DoApi"
	request.RegionId = client.RegionId
	request.Headers = map[string]string{"RegionId": client.RegionId, "x-acs-instanceid": bucketName}

	raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.ProcessCommonRequest(request)
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, objectName, "HeadObject", AlibabacloudStackOssGoSdk)
	}

	bresponse, _ := raw.(*responses.CommonResponse)
	if bresponse.GetHttpStatus() != 200 {
		return WrapErrorf(err, DefaultErrorMsg, objectName, "HeadObject", AlibabacloudStackOssGoSdk)
	}
	addDebug("HeadObject", raw, requestInfo, bresponse.GetHttpContentString())

	resp := make(map[string]interface{})
	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &resp)
	if err != nil {
		return WrapError(err)
	}

	if resp["asapiSuccess"] == false && (resp["Message"] == "Not Found" || resp["Code"] == "NoSuchKey") {
		return WrapErrorf(Error(GetNotFoundMessage("OssObject", objectName)), NotFoundMsg, AlibabacloudStackSdkGoERROR)
	}

	return nil
}

func (s *OssService) WaitForOssBucketObject(bucket *oss.Bucket, id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		err := s.HeadOssBucketObject(bucket.BucketName, id)
		if err != nil {
			if NotFoundError(err) {
				return nil
			}
			return WrapErrorf(err, DefaultErrorMsg, id, "IsObjectExist", AlibabacloudStackOssGoSdk)
		}

		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, strconv.FormatBool(true), status, ProviderERROR)
		}
	}
}
