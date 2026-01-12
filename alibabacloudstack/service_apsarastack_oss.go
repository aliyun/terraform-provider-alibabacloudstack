package alibabacloudstack

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
)

// OssService *connectivity.AlibabacloudStackClient
type OssService struct {
	client *connectivity.AlibabacloudStackClient
}

type BucketSyncRule struct {
	Status                      string            `json:"Status"`
	Destination                 map[string]string `json:"Destination"`
	Action                      string            `json:"Action"`
	ID                          string            `json:"ID"`
	SyncRole                    string            `json:"SyncRole"`
	SrcLocation                 string            `json:"SrcLocation"`
	EncryptionConfiguration     map[string]string `json:"EncryptionConfiguration"`
	HistoricalObjectReplication string            `json:"HistoricalObjectReplication"`
}

type BucketSyncResponse struct {
	RequestID string `json:"requestId"`
	Code      string `json:"code"`
	Data      struct {
		ReplicationConfiguration struct {
			Rule []BucketSyncRule `json:"Rule"`
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

func (s *OssService) ListOssBucket() (response []BucketListBucket, err error) {
	request := s.client.NewCommonRequest("GET", "OneRouter", "2018-12-12", "DoOpenApi", "")
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
	// 3.16.2 will return, but 3.18.x does not return
	if err != nil || (bucketList.Code != "" && bucketList.Code != "200") {
		return buckets, errmsgs.WrapError(err)
	}

	if _, ok := bucketList.Data.ListAllMyBucketsResult.Buckets.(string); ok {
		return buckets, errmsgs.GetNotFoundErrorFromString("Not Found: Oss Bucket")
	}

	var bucketInterface interface{}
	if v, ok := bucketList.Data.ListAllMyBucketsResult.Buckets.(map[string]interface{}); !ok {
		return buckets, errmsgs.WrapErrorf(err, "Error Response Format")
	} else {
		bucketInterface = v["Bucket"]
	}

	switch v := bucketInterface.(type) {
	case map[string]interface{}:
		// Single Bucket structure
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
		// Multiple Bucket structures
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

func (s *OssService) DescribeOssBucket(id string) (response oss.GetBucketInfoResult, err error) {

	response.BucketInfo.Name = ""
	if buckets, err := s.ListOssBucket(); err == nil {
		for _, j := range buckets {
			if j.Name == id {
				response.BucketInfo.Name = j.Name
				response.BucketInfo.StorageClass = j.StorageClass
				response.BucketInfo.ExtranetEndpoint = j.ExtranetEndpoint
				response.BucketInfo.IntranetEndpoint = j.IntranetEndpoint
				response.BucketInfo.Location = j.Location
				break
			}
		}
	}

	return response, err
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
	request := s.client.NewCommonRequest("GET", "OneRouter", "2018-12-12", "DoApi", "")
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
	request := s.client.NewCommonRequest("GET", "OneRouter", "2018-12-12", "DoOpenApi", "")
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

	if v, ok := resp["asapiSuccess"]; ok && !v.(bool) {
		return errmsgs.WrapError(errmsgs.Error(fmt.Sprintf("put bucket tags error %#v", resp)))
	}

	return nil
}

func (s *OssService) GetBucketTags(bucketName string) (tags []interface{}, err error) {
	tags = make([]interface{}, 0)
	request := s.client.NewCommonRequest("GET", "OneRouter", "2018-12-12", "DoOpenApi", "")
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
		ok := true
		tags, ok = tags_data.([]interface{})
		if !ok {
			ts, ok := tags_data.(map[string]interface{})
			if !ok {
				return nil, errmsgs.WrapErrorf(err, "GetBucketTags", errmsgs.AlibabacloudStackOssGoSdk, fmt.Sprintf("GetBucketTags error : %#v", tags_data))
			}
			tags = []interface{}{ts}
		}
	}
	return tags, err
}

func (s *OssService) DeleteBucketTags(bucketName string) error {
	request := s.client.NewCommonRequest("GET", "OneRouter", "2018-12-12", "DoOpenApi", "")
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

func (s *OssService) DeleteBucket(bucketName string) error {
	return resource.Retry(1*time.Minute, func() *resource.RetryError {
		request := s.client.NewCommonRequest("GET", "OneRouter", "2018-12-12", "DoOpenApi", "")
		request.QueryParams["OpenApiAction"] = "DeleteBucket"
		request.QueryParams["ProductName"] = "oss"

		params := map[string]string{
			"Department":          s.client.Department,
			"ResourceGroup":       s.client.ResourceGroup,
			"RegionId":            s.client.RegionId,
			"asVersion":           "enterprise",
			"asArchitechture":     "x86",
			"haAlibabacloudStack": "true",
			"Language":            "en",
			"BucketName":          bucketName,
			"StorageClass":        "Standard",
		}

		if content, err := json.Marshal(params); err != nil {
			return resource.NonRetryableError(err)
		} else {
			request.QueryParams["Params"] = string(content)
		}

		bresponse, err := s.client.ProcessCommonRequest(request)

		if err != nil {
			if bresponse == nil {
				return resource.RetryableError(errmsgs.WrapErrorf(err, "Process Common Request Failed"))
			}
			if ossNotFoundError(err) {
				return resource.NonRetryableError(errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackOssGoSdk))
			}
			errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			return resource.RetryableError(errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, bucketName, "DeleteBucket", errmsgs.AlibabacloudStackOssGoSdk, errmsg))
		}
		det, err := s.DescribeOssBucket(bucketName)
		if err != nil {
			return resource.NonRetryableError(errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, bucketName, "IsBucketExist", errmsgs.AlibabacloudStackOssGoSdk))
		}
		if det.BucketInfo.Name != "" {
			return resource.RetryableError(errmsgs.Error("Trying to delete OSS bucket %#v successfully.", bucketName))
		}
		return nil
	})
}

func (s *OssService) GetOssEndpointList() ([]interface{}, error) {
	request := s.client.NewCommonRequest("GET", "OneRouter", "2018-12-12", "DoApi", "")
	request.QueryParams["AppAction"] = "GetOssEndpointList"
	request.QueryParams["AppName"] = "one-console-app-oss"
	request.QueryParams["Params"] = fmt.Sprintf("{\"params\":{\"region\":\"%s\"}}", s.client.RegionId)
	bresponse, err := s.client.ProcessCommonRequest(request)
	addDebug("GetOssEndpointList", bresponse, request, request.QueryParams)
	if err != nil {
		if bresponse == nil {
			return nil, errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return nil, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "GetOssEndpointList", errmsgs.AlibabacloudStackOssGoSdk, errmsg)
	}
	result := make(map[string]interface{})
	_ = json.Unmarshal(bresponse.GetHttpContentBytes(), &result)
	data, ok := result["Data"]
	if !ok || len(data.([]interface{})) == 0 {
		return nil, errmsgs.Error(fmt.Sprintf("GetOssEndpointList Failed! region: %s \n %#v", s.client.RegionId, bresponse.GetHttpContentString()))
	}
	return data.([]interface{}), nil
}

func (s *OssService) ossTagIgnored(t map[string]interface{}) bool {
	filter := []string{"^aliyun", "^acs:", "^ascm:", "^http://", "^https://"}
	for _, v := range filter {
		ok, _ := regexp.MatchString(v, t["Key"].(string))
		if ok {
			return true
		}
	}
	return false
}

func (s *OssService) GetBucketSync(bucketName string) (object *BucketSyncResponse, err error) {
	request := s.client.NewCommonRequest("GET", "OneRouter", "2018-12-12", "DoOpenApi", "")
	request.QueryParams["OpenApiAction"] = "GetBucketSync"
	request.QueryParams["ProductName"] = "oss"
	request.QueryParams["Params"] = fmt.Sprintf("{\"BucketName\":\"%s\"}", bucketName)

	bresponse, err := s.client.ProcessCommonRequest(request)
	addDebug(request.GetActionName(), bresponse, request, request.QueryParams)
	if err != nil {
		if bresponse == nil {
			return nil, errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return nil, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, bucketName, "GetBucketSync", errmsgs.AlibabacloudStackOssGoSdk, errmsg)
	}
	bucketSync := BucketSyncResponse{}
	err = json.Unmarshal([]byte(bresponse.GetHttpContentString()), &bucketSync)
	if err != nil {
		return nil, errmsgs.WrapErrorf(err, "Process Common Request Failed")
	}
	return &bucketSync, nil
}

func (s *OssService) OssBucketSyncStateRefreshFunc(bucketName string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		result, err := s.GetBucketSync(bucketName)
		if err != nil {
			if errmsgs.NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", errmsgs.WrapError(err)
		}
		var object BucketSyncRule
		for _, rule := range result.Data.ReplicationConfiguration.Rule {
			if rule.SrcLocation == "" {
				object = rule
				break
			}
		}
		for _, failState := range failStates {
			if object.Status == failState {
				return object, object.Status, errmsgs.WrapError(errmsgs.Error(errmsgs.FailedToReachTargetStatus, object.Status))
			}
		}
		return object, object.Status, nil
	}
}

func (s OssService) GetBucketClient(bucketName string) (*oss.Bucket, error) {
	bucketInfo, err := s.DescribeOssBucket(bucketName)
	if bucketInfo.BucketInfo.Name == "" {
		return nil, errmsgs.GetNotFoundErrorFromString("Bucket " + bucketName + " Not Found")
	}
	var ossconn *oss.Client

	bucketEndpoint := bucketInfo.BucketInfo.ExtranetEndpoint
	schma := strings.ToLower(s.client.Config.Protocol)
	if !strings.HasPrefix(bucketEndpoint, "http") {
		bucketEndpoint = fmt.Sprintf("%s://%s", schma, bucketEndpoint)
	}

	clientOptions := []oss.ClientOption{oss.UserAgent(s.client.GetUserAgent()),
		oss.SecurityToken(s.client.Config.SecurityToken)}
	if s.client.Config.Proxy != "" {
		clientOptions = append(clientOptions, oss.Proxy(s.client.Config.Proxy))
	}

	clientOptions = append(clientOptions, oss.UseCname(false))

	if ossconn, err = oss.New(bucketEndpoint, s.client.Config.AccessKey, s.client.Config.SecretKey, clientOptions...); err != nil {
		return nil, fmt.Errorf("unable to initialize the OSS client: %#v", err)
	}

	bucket, err := ossconn.Bucket(bucketName)

	if err != nil {
		return nil, fmt.Errorf("unable to get the bucket %s: %#v", bucketName, err)
	} else {
		return bucket, nil
	}

}
func (s *OssService) DescribeOssSingleTunnel(id string) (map[string]interface{}, error) {
	request := map[string]interface{}{}

	response, err := s.client.DoTeaRequest("GET", "oss", "2019-09-01", "ListVpcip", "", nil, request, nil)
	if err != nil {
		return nil, err
	}

	if vpcipList, ok := response["ListVpcipResult"].(map[string]interface{})["Vpcip"].([]interface{}); ok {
		for _, item := range vpcipList {
			vpcip := item.(map[string]interface{})

			if vip, exists := vpcip["Vip"]; exists && vip == id {
				return vpcip, nil
			}
		}
	}

	return nil, errmsgs.GetNotFoundErrorFromString(fmt.Sprintf("OSS Single Tunnel not found with id: %s", id))
}
