package alibabacloudstack

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
)

// OssService *connectivity.AlibabacloudStackClient
type OssService struct {
	client *connectivity.AlibabacloudStackClient
}

func (s *OssService) DescribeOssBucket(id string) (response oss.GetBucketInfoResult, err error) {
	request := s.client.NewCommonRequest("POST", "OneRouter", "2018-12-12", "DoOpenApi", "")
	request.QueryParams["OpenApiAction"]=    "GetService"
	request.QueryParams["ProductName"]=      "oss"
	var bucketList = &BucketList{}
	raw, err := s.client.WithOssNewClient(func(ossClient *ecs.Client) (interface{}, error) {
		return ossClient.ProcessCommonRequest(request)
	})

	bresponse, ok := raw.(*responses.CommonResponse)
	if err != nil {
		if ossNotFoundError(err) {
			return response, errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackOssGoSdk)
		}
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		return response, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, "GetBucketInfo", errmsgs.AlibabacloudStackOssGoSdk, errmsg)
	}
	addDebug("GetBucketInfo", raw, request)

	err = json.Unmarshal(bresponse.GetHttpContentBytes(), bucketList)
	if err != nil {
		return response, errmsgs.WrapError(err)
	}
	if bucketList.Code != "200" || len(bucketList.Data.ListAllMyBucketsResult.Buckets.Bucket) < 1 {
		return response, errmsgs.WrapError(err)
	}

	var found = false
	for _, j := range bucketList.Data.ListAllMyBucketsResult.Buckets.Bucket {
		if j.Name == id {
			response.BucketInfo.Name = j.Name
			response.BucketInfo.StorageClass = j.StorageClass
			response.BucketInfo.ExtranetEndpoint = j.ExtranetEndpoint
			response.BucketInfo.IntranetEndpoint = j.IntranetEndpoint
			response.BucketInfo.Owner.ID = fmt.Sprint(j.ResourceGroupName)
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

func (s *OssService) ListOssBucket() (response BucketList, err error) {
	request := s.client.NewCommonRequest("POST", "OneRouter", "2018-12-12", "DoOpenApi", "")
	mergeMaps(request.QueryParams, map[string]string{
		"AccountInfo":      "",
		"SignatureVersion": "1.0",
		"OpenApiAction":    "GetService",
		"ProductName":      "oss",
	})
	bucketList := BucketList{}
	raw, err := s.client.WithOssNewClient(func(ossClient *ecs.Client) (interface{}, error) {
		return ossClient.ProcessCommonRequest(request)
	})

	bresponse, ok := raw.(*responses.CommonResponse)
	if err != nil {
		if ossNotFoundError(err) {
			return response, errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackOssGoSdk)
		}
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		return response, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "GetBucketInfo", errmsgs.AlibabacloudStackOssGoSdk, errmsg)
	}
	addDebug("GetBucketInfo", raw, request)

	err = json.Unmarshal(bresponse.GetHttpContentBytes(), bucketList)
	if err != nil {
		return response, errmsgs.WrapError(err)
	}
	if bucketList.Code != "200" || len(bucketList.Data.ListAllMyBucketsResult.Buckets.Bucket) < 1 {
		return response, errmsgs.WrapError(err)
	}
	return bucketList, nil
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
			if errmsgs.NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return errmsgs.WrapError(err)
			}
		}

		if object.BucketInfo.Name != "" && status != Deleted {
			return nil
		}
		if object.BucketInfo.Name == "" && status == Deleted {
			return nil
		}
		if time.Now().After(deadline) {
			return errmsgs.WrapErrorf(err, errmsgs.WaitTimeoutMsg, id, GetFunc(1), timeout, object.BucketInfo.Name, status, errmsgs.ProviderERROR)
		}
	}
}

func (s *OssService) HeadOssBucketObject(bucketName string, objectName string) error {
	request := s.client.NewCommonRequest("POST", "OneRouter", "2018-12-12", "DoApi", "")
	mergeMaps(request.QueryParams, map[string]string{
		"AppAction":   "HeadObject",
		"AppName":     "one-console-app-oss",
		"Params":      "{\"region\":\"" + s.client.RegionId + "\",\"params\":{\"bucketName\":\"" + bucketName + "\",\"objectName\":\"" + objectName + "\"}}",
		"AccountInfo": "",
	})
	request.Headers["x-acs-instanceid"] = bucketName

	raw, err := s.client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.ProcessCommonRequest(request)
	})

	bresponse, ok := raw.(*responses.CommonResponse)
	if err != nil || bresponse.GetHttpStatus() != 200 {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, objectName, "HeadObject", errmsgs.AlibabacloudStackOssGoSdk, errmsg)
	}

	addDebug("HeadObject", raw, request, bresponse.GetHttpContentString())

	resp := make(map[string]interface{})
	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &resp)
	if err != nil {
		return errmsgs.WrapError(err)
	}

	if resp["asapiSuccess"] == false && (resp["Message"] == "Not Found" || resp["Code"] == "NoSuchKey") {
		return errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("OssObject", objectName)), errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
	}

	return nil
}

func (s *OssService) WaitForOssBucketObject(bucket *oss.Bucket, id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		err := s.HeadOssBucketObject(bucket.BucketName, id)
		if err != nil {
			if errmsgs.NotFoundError(err) {
				return nil
			}
			return err
		}

		if time.Now().After(deadline) {
			return errmsgs.WrapErrorf(err, errmsgs.WaitTimeoutMsg, id, GetFunc(1), timeout, strconv.FormatBool(true), status, errmsgs.ProviderERROR)
		}
	}
}
