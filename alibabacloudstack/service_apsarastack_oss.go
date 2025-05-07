package alibabacloudstack

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
)

// OssService *connectivity.AlibabacloudStackClient
type OssService struct {
	client *connectivity.AlibabacloudStackClient
}

type BucketSyncResponse struct {
	RequestID string `json:"requestId"`
	Code      string `json:"code"`
	Data      struct {
		ReplicationConfiguration struct {
			Rule []struct {
				Status                      string            `json:"Status"`
				Destination                 map[string]string `json:"Destination"`
				Action                      string            `json:"Action"`
				ID                          string            `json:"ID"`
				HistoricalObjectReplication string            `json:"HistoricalObjectReplication"`
			} `json:"Rule"`
		} `json:"ReplicationConfiguration"`
	} `json:"data"`
	Cost            int    `json:"cost"`
	APICost         int    `json:"apiCost"`
	EagleEyeTraceID string `json:"eagleEyeTraceId"`
	AscmCode        bool   `json:"ascmCode"`
	SuccessResponse bool   `json:"successResponse"`
}

type BucketAclResponse struct {
	RequestID string `json:"requestId"`
	Code      string `json:"code"`
	Data      struct {
		AccessControlPolicy struct {
			AccessControlList struct {
				Grant string `json:"Grant"`
			} `json:"AccessControlList"`
			Owner struct {
				DisplayName string `json:"DisplayName"`
				ID          string `json:"ID"`
			} `json:"Owner"`
		} `json:"AccessControlPolicy"`
	} `json:"data"`
	Cost            int    `json:"cost"`
	APICost         int    `json:"apiCost"`
	EagleEyeTraceID string `json:"eagleEyeTraceId"`
	AscmCode        bool   `json:"ascmCode"`
	SuccessResponse bool   `json:"successResponse"`
}

type BucketStorageCapacityResponse struct {
	RequestID string `json:"requestId"`
	Data      struct {
		BucketUserQos struct {
			StorageCapacity string `json:"StorageCapacity"`
		} `json:"BucketUserQos"`
	} `json:"data"`
}

type BucketEncryptionResponse struct {
	RequestID string `json:"requestId"`
	Code      string `json:"code"`
	Data      struct {
		ServerSideEncryptionRule struct {
			ApplyServerSideEncryptionByDefault struct {
				SSEAlgorithm   string `json:"SSEAlgorithm"`
				KMSMasterKeyID string `json:"KMSMasterKeyID"`
			} `json:"ApplyServerSideEncryptionByDefault"`
		} `json:"ServerSideEncryptionRule"`
	} `json:"data"`
}

func (s *OssService) DescribeOssBucket(id string) (response oss.GetBucketInfoResult, err error) {
	request := s.client.NewCommonRequest("POST", "OneRouter", "2018-12-12", "DoOpenApi", "")
	request.QueryParams["OpenApiAction"] = "GetService"
	request.QueryParams["ProductName"] = "oss"
	bresponse, err := s.client.ProcessCommonRequest(request)
	addDebug("GetBucketInfo", bresponse, request)
	if err != nil {
		if bresponse == nil {
			return response, errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		if ossNotFoundError(err) {
			return response, errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackOssGoSdk)
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return response, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, "GetBucketInfo", errmsgs.AlibabacloudStackOssGoSdk, errmsg)
	}

	buckets, err := getBucketListResponseBuckets(bresponse)

	var found = false
	for _, j := range buckets {
		if j.Name == id {
			response.BucketInfo.Name = j.Name
			response.BucketInfo.StorageClass = j.StorageClass
			response.BucketInfo.ExtranetEndpoint = j.ExtranetEndpoint
			response.BucketInfo.IntranetEndpoint = j.IntranetEndpoint
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

func (s *OssService) ListOssBucket() (response []BucketListBucket, err error) {
	request := s.client.NewCommonRequest("POST", "OneRouter", "2018-12-12", "DoOpenApi", "")
	request.QueryParams["OpenApiAction"] = "GetService"
	request.QueryParams["ProductName"] = "oss"
	bresponse, err := s.client.ProcessCommonRequest(request)
	if err != nil {
		if bresponse == nil {
			return response, errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		if ossNotFoundError(err) {
			return response, errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackOssGoSdk)
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return response, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "GetBucketInfo", errmsgs.AlibabacloudStackOssGoSdk, errmsg)
	}
	addDebug("GetBucketInfo", bresponse, request)

	buckets, err := getBucketListResponseBuckets(bresponse)
	if err != nil {
		return buckets, errmsgs.WrapError(err)
	}
	return buckets, nil
}
func getBucketListResponseBuckets(response *responses.CommonResponse) ([]BucketListBucket, error) {
	var buckets []BucketListBucket

	var bucketList BucketList
	err := json.Unmarshal(response.GetHttpContentBytes(), &bucketList)
	// 3.16.2 会发返回，但3.18.x不返回
	if err != nil || (bucketList.Code != "" && bucketList.Code != "200") {
		return buckets, errmsgs.WrapError(err)
	}

	if _, ok := bucketList.Data.ListAllMyBucketsResult.Buckets.(string); ok {
		return buckets, errmsgs.WrapErrorf(err, "Not Found: Oss Bucket")
	}

	var bucketInterface interface{}
	if v, ok := bucketList.Data.ListAllMyBucketsResult.Buckets.(map[string]interface{}); !ok {
		return buckets, errmsgs.WrapErrorf(err, "Error Response Format")
	} else {
		bucketInterface = v["Bucket"]
	}

	switch v := bucketInterface.(type) {
	case map[string]interface{}:
		// 单个 Bucket 结构体
		bucket := BucketListBucket{
			Comment:          v["Comment"].(string),
			CreationDate:     v["CreationDate"].(string),
			ExtranetEndpoint: v["ExtranetEndpoint"].(string),
			IntranetEndpoint: v["IntranetEndpoint"].(string),
			Location:         v["Location"].(string),
			Name:             v["Name"].(string),
			StorageClass:     v["StorageClass"].(string),
		}
		buckets = append(buckets, bucket)
	case []interface{}:
		// 多个 Bucket 结构体
		for _, vv := range v {
			vvv, ok := vv.(map[string]interface{})
			if !ok {
				return buckets, errmsgs.WrapErrorf(err, "Error Response Format")
			}

			bucket := BucketListBucket{
				Comment:          vvv["Comment"].(string),
				CreationDate:     vvv["CreationDate"].(string),
				ExtranetEndpoint: vvv["ExtranetEndpoint"].(string),
				IntranetEndpoint: vvv["IntranetEndpoint"].(string),
				Location:         vvv["Location"].(string),
				Name:             vvv["Name"].(string),
				StorageClass:     vvv["StorageClass"].(string),
			}
			buckets = append(buckets, bucket)
		}
	default:
		return buckets, errmsgs.WrapErrorf(err, "Error Response Format")
	}
	return buckets, nil
}

type BucketListBucket struct {
	Comment          string `json:"Comment"`
	CreationDate     string `json:"CreationDate"`
	ExtranetEndpoint string `json:"ExtranetEndpoint"`
	IntranetEndpoint string `json:"IntranetEndpoint"`
	Location         string `json:"Location"`
	Name             string `json:"Name"`
	StorageClass     string `json:"StorageClass"`
}

type BucketList struct {
	Data struct {
		ListAllMyBucketsResult struct {
			Buckets interface{} `json:"Buckets"`
			Owner   struct{}    `json:"Owner"`
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
		"AppAction": "HeadObject",
		"AppName":   "one-console-app-oss",
		"Params":    "{\"region\":\"" + s.client.RegionId + "\",\"params\":{\"bucketName\":\"" + bucketName + "\",\"objectName\":\"" + objectName + "\"}}",
	})
	request.Headers["x-acs-instanceid"] = bucketName

	bresponse, err := s.client.ProcessCommonRequest(request)

	if err != nil || bresponse.GetHttpStatus() != 200 {
		if bresponse == nil {
			return errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, objectName, "HeadObject", errmsgs.AlibabacloudStackOssGoSdk, errmsg)
	}

	addDebug("HeadObject", bresponse, request, bresponse.GetHttpContentString())

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

func (s *OssService) PutOssBucketTags(bucketName string, tags []OssTags) error {
	osstags := ""
	if len(tags) > 0 {
		for _, tag := range tags {
			tag := fmt.Sprintf(`<Tag><Key>%s</Key><Value>%s</Value></Tag>`, tag.Key, tag.Value)
			osstags = osstags + tag
		}
	} else {
		osstags = "<Tag></Tag>"
	}

	content := fmt.Sprintf(`<Tagging><TagSet>%s</TagSet></Tagging>`, osstags)
	request := s.client.NewCommonRequest("POST", "OneRouter", "2018-12-12", "DoOpenApi", "")
	mergeMaps(request.QueryParams, map[string]string{
		"OpenApiAction": "PutBucketTags",
		"ProductName":   "oss",
		"Content":       content,
		"Params":        "{\"BucketName\":\"" + bucketName + "\"}",
	})
	request.Headers["x-acs-instanceid"] = bucketName

	bresponse, err := s.client.ProcessCommonRequest(request)
	addDebug("PutBucketTags", bresponse, request, bresponse.GetHttpContentString())
	if err != nil || bresponse.GetHttpStatus() != 200 {
		if bresponse == nil {
			return errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "PutBucketTags", errmsgs.AlibabacloudStackOssGoSdk, errmsg)
	}
	resp := make(map[string]interface{})
	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &resp)
	if err != nil {
		return errmsgs.WrapError(err)
	}

	if resp["asapiSuccess"].(bool) == false {
		return errmsgs.WrapError(errmsgs.Error(fmt.Sprintf("put bucket tags error %#v", resp)))
	}

	return nil
}

func (s *OssService) GetBucketTags(bucketName string) (tags []interface{}, err error) {
	tags = make([]interface{}, 0)
	request := s.client.NewCommonRequest("POST", "OneRouter", "2018-12-12", "DoOpenApi", "")
	mergeMaps(request.QueryParams, map[string]string{
		"OpenApiAction": "GetBucketTags",
		"ProductName":   "oss",
		"Params":        "{\"BucketName\":\"" + bucketName + "\"}",
	})
	request.Headers["x-acs-instanceid"] = bucketName

	bresponse, err := s.client.ProcessCommonRequest(request)
	addDebug("GetBucketTags", bresponse, request, bresponse.GetHttpContentString())
	if err != nil || bresponse.GetHttpStatus() != 200 {
		if bresponse == nil {
			return nil, errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return nil, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "GetBucketTags", errmsgs.AlibabacloudStackOssGoSdk, errmsg)
	}
	response := make(map[string]interface{})
	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &response)
	if err != nil {
		return nil, errmsgs.WrapError(err)
	}
	tags_data, err := jsonpath.Get("$.Data.Tagging.TagSet.Tag", response)
	if tags_data != nil {
		tags = tags_data.([]interface{})
	}
	return tags, err
}

func (s *OssService) DeleteBucketTags(bucketName string) error {
	request := s.client.NewCommonRequest("POST", "OneRouter", "2018-12-12", "DoOpenApi", "")
	mergeMaps(request.QueryParams, map[string]string{
		"OpenApiAction": "DeleteBucketTags",
		"ProductName":   "oss",
		"Params":        "{\"BucketName\":\"" + bucketName + "\"}",
	})
	request.Headers["x-acs-instanceid"] = bucketName

	bresponse, err := s.client.ProcessCommonRequest(request)
	addDebug("DeleteBucketTags", bresponse, request, bresponse.GetHttpContentString())
	if err != nil || bresponse.GetHttpStatus() != 200 {
		if bresponse == nil {
			return errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "GetBucketTags", errmsgs.AlibabacloudStackOssGoSdk, errmsg)
	}
	return nil
}
