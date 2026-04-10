package alibabacloudstack

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/credentials"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// OssService *connectivity.AlibabacloudStackClient
type OssService struct {
	client *connectivity.AlibabacloudStackClient
}

type BucketSyncRule struct {
	Status      string `xml:"Status" json:"Status"`
	Destination struct {
		Bucket   string `xml:"Bucket"`
		Location string `xml:"Location"`
	} `xml:"Destination" json:"Destination"`
	Action                  string `xml:"Action" json:"Action"`
	ID                      string `xml:"ID" json:"ID"`
	SyncRole                string `xml:"SyncRole" json:"SyncRole"`
	SrcLocation             string `xml:"SrcLocation" json:"SrcLocation"`
	EncryptionConfiguration struct {
		ReplicaKmsKeyID string `xml:"ReplicaKmsKeyID"`
	} `xml:"EncryptionConfiguration" json:"EncryptionConfiguration"`
	HistoricalObjectReplication string `xml:"HistoricalObjectReplication" json:"HistoricalObjectReplication"`
}

// BucketSyncXMLResponse represents the XML response from GetBucketSync (OSS native API)
type BucketSyncXMLResponse struct {
	Rule []BucketSyncRule `xml:"Rule"`
}

// BucketSyncResponse wraps the XML response for compatibility with existing code
type BucketSyncResponse struct {
	Data struct {
		ReplicationConfiguration BucketSyncXMLResponse
	}
}

type BucketStorageCapacityResponse struct {
	RequestID string `json:"requestId"`
	Data      struct {
		BucketUserQos struct {
			StorageCapacity string `json:"StorageCapacity"`
		} `json:"BucketUserQos"`
	} `json:"data"`
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
		name := *object.Name
		if name != "" && status != Deleted {
			return nil
		}
		if name == "" && status == Deleted {
			return nil
		}
		if time.Now().After(deadline) {
			return errmsgs.WrapErrorf(err, errmsgs.WaitTimeoutMsg, id, GetFunc(1), timeout, object.Name, status, errmsgs.ProviderERROR)
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

func (s *OssService) WaitForOssBucketObject(bucketName string, id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		err := s.HeadOssBucketObject(bucketName, id)
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

func (s *OssService) PutOssBucketTags(bucketId string, tags []OssTags) error {
	ossClient, err := s.GetBucketClient(bucketId)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	var bucket_name string
	parts := strings.Split(bucketId, COLON_SEPARATED)
	if len(parts) > 1 {
		bucket_name = parts[0]
	} else {
		bucket_name = bucketId
	}
	ossTags := make([]oss.Tag, 0, len(tags))
	for _, tag := range tags {
		key := tag.Key
		val := tag.Value
		ossTags = append(ossTags, oss.Tag{
			Key:   &key,
			Value: &val,
		})
	}

	request := &oss.PutBucketTagsRequest{
		Bucket: &bucket_name,
		Tagging: &oss.Tagging{
			TagSet: &oss.TagSet{
				Tags: ossTags,
			},
		},
	}
	_, err = ossClient.PutBucketTags(context.Background(), request)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, bucketId, "PutBucketTags", errmsgs.AlibabacloudStackOssGoSdk)
	}
	return nil
}

func (s *OssService) GetBucketTags(bucketId string) (tags []interface{}, err error) {
	tags = make([]interface{}, 0)
	ossClient, err := s.GetBucketClient(bucketId)
	if err != nil {
		return nil, errmsgs.WrapError(err)
	}
	var bucket_name string
	parts := strings.Split(bucketId, COLON_SEPARATED)
	if len(parts) > 1 {
		bucket_name = parts[0]
	} else {
		bucket_name = bucketId
	}
	result, err := ossClient.GetBucketTags(context.Background(), &oss.GetBucketTagsRequest{
		Bucket: &bucket_name,
	})
	if err != nil {
		return nil, errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, bucketId, "GetBucketTags", errmsgs.AlibabacloudStackOssGoSdk)
	}

	if result.Tagging != nil && result.Tagging.TagSet != nil {
		for _, t := range result.Tagging.TagSet.Tags {
			tagMap := map[string]interface{}{}
			if t.Key != nil {
				tagMap["Key"] = *t.Key
			}
			if t.Value != nil {
				tagMap["Value"] = *t.Value
			}
			tags = append(tags, tagMap)
		}
	}
	return tags, nil
}

func (s *OssService) DeleteBucketTags(bucketId string) error {
	ossClient, err := s.GetBucketClient(bucketId)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	var bucket_name string
	parts := strings.Split(bucketId, COLON_SEPARATED)
	if len(parts) > 1 {
		bucket_name = parts[0]
	} else {
		bucket_name = bucketId
	}
	_, err = ossClient.DeleteBucketTags(context.Background(), &oss.DeleteBucketTagsRequest{
		Bucket: &bucket_name,
	})
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, bucketId, "DeleteBucketTags", errmsgs.AlibabacloudStackOssGoSdk)
	}
	return nil
}

func (s *OssService) DeleteBucket(bucketId string) error {
	ossClient, err := s.GetBucketClient(bucketId)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	var bucket_name string
	parts := strings.Split(bucketId, COLON_SEPARATED)
	if len(parts) > 1 {
		bucket_name = parts[0]
	} else {
		bucket_name = bucketId
	}
	return resource.Retry(1*time.Minute, func() *resource.RetryError {
		_, err := ossClient.DeleteBucket(context.Background(), &oss.DeleteBucketRequest{
			Bucket: &bucket_name,
		})
		if err != nil {
			if ossNotFoundError(err) {
				return resource.NonRetryableError(errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackOssGoSdk))
			}
			return resource.RetryableError(errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, bucketId, "DeleteBucket", errmsgs.AlibabacloudStackOssGoSdk))
		}
		det, err := s.DescribeOssBucket(bucketId)
		if err != nil {
			return resource.NonRetryableError(errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, bucketId, "IsBucketExist", errmsgs.AlibabacloudStackOssGoSdk))
		}
		if *det.Name != "" {
			return resource.RetryableError(errmsgs.Error("Trying to delete OSS bucket %#v failed.", bucketId))
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
	ossService := OssSdkService{s.client}
	ossClient, err := ossService.GetBucketClient(bucketName)
	if err != nil {
		return nil, errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, bucketName, "GetBucketSync", errmsgs.AlibabacloudStackOssGoSdk)
	}
	input := &oss.OperationInput{
		OpName:     "GetBucketSync",
		Method:     "GET",
		Bucket:     oss.Ptr(bucketName),
		Parameters: map[string]string{"replication": ""},
	}
	output, err := ossClient.InvokeOperation(context.Background(), input)
	addDebug("GetBucketSync", output, input, nil)
	if err != nil {
		return nil, errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, bucketName, "GetBucketSync", errmsgs.AlibabacloudStackOssGoSdk)
	}
	defer output.Body.Close()
	body, err := io.ReadAll(output.Body)
	if err != nil {
		return nil, errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, bucketName, "GetBucketSync", errmsgs.AlibabacloudStackOssGoSdk)
	}
	var xmlResp BucketSyncXMLResponse
	if err = xml.Unmarshal(body, &xmlResp); err != nil {
		return nil, errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, bucketName, "GetBucketSync", errmsgs.AlibabacloudStackOssGoSdk)
	}
	result := &BucketSyncResponse{}
	result.Data.ReplicationConfiguration = xmlResp
	return result, nil
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

func (s OssService) GetBucketClient(bucketName string) (*oss.Client, error) {
	bucketInfo, err := s.DescribeOssBucket(bucketName)
	if err != nil {
		return nil, errmsgs.WrapError(err)
	}
	if bucketInfo.Name == "" {
		return nil, errmsgs.GetNotFoundErrorFromString("Bucket " + bucketName + " Not Found")
	}

	bucketEndpoint := bucketInfo.ExtranetEndpoint
	schma := strings.ToLower(s.client.Config.Protocol)
	if !strings.HasPrefix(bucketEndpoint, "http") {
		bucketEndpoint = fmt.Sprintf("%s://%s", schma, bucketEndpoint)
	}

	cfg := oss.LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			s.client.Config.AccessKey,
			s.client.Config.SecretKey,
			s.client.Config.SecurityToken,
		)).
		WithEndpoint(bucketEndpoint).
		WithRegion(s.client.RegionId)

	client := oss.NewClient(cfg)
	return client, nil
}

// VpcipEntry represents a single VPC IP entry in the ListVpcip XML response
type VpcipEntry struct {
	Cluster string `xml:"Cluster"`
	VpcId   string `xml:"VpcId"`
	Vip     string `xml:"Vip"`
	Label   string `xml:"Label"`
	Shared  int    `xml:"Shared"`
}

// ListVpcipXMLResponse represents the XML response from ListVpcip
type ListVpcipXMLResponse struct {
	Vpcip []VpcipEntry `xml:"Vpcip"`
}

func (s *OssService) listVpcipEntries() ([]VpcipEntry, error) {
	ossService := OssSdkService{s.client}
	ossClient, err := ossService.GetOssClient()
	if err != nil {
		return nil, errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "ListVpcip", "GetOssClient", errmsgs.AlibabacloudStackOssGoSdk)
	}
	input := &oss.OperationInput{
		OpName:     "ListVpcip",
		Method:     "GET",
		Parameters: map[string]string{"vpcip": ""},
	}
	output, err := ossClient.InvokeOperation(context.Background(), input)
	addDebug("ListVpcip", output, input, nil)
	if err != nil {
		return nil, errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "ListVpcip", "InvokeOperation", errmsgs.AlibabacloudStackOssGoSdk)
	}
	defer output.Body.Close()
	body, err := io.ReadAll(output.Body)
	if err != nil {
		return nil, errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "ListVpcip", "ReadBody", errmsgs.AlibabacloudStackOssGoSdk)
	}
	var xmlResp ListVpcipXMLResponse
	if err = xml.Unmarshal(body, &xmlResp); err != nil {
		return nil, errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "ListVpcip", "XMLUnmarshal", errmsgs.AlibabacloudStackOssGoSdk)
	}
	return xmlResp.Vpcip, nil
}

func (s *OssService) DescribeOssSingleTunnel(id string) (map[string]interface{}, error) {
	parts := strings.Split(id, ":")
	entries, err := s.listVpcipEntries()
	if err != nil {
		return nil, err
	}
	for _, entry := range entries {
		if entry.Cluster != parts[0] {
			continue
		}
		if entry.VpcId != parts[1] {
			continue
		}
		if entry.Vip != parts[2] {
			continue
		}
		data := map[string]interface{}{
			"Cluster": entry.Cluster,
			"VpcId":   entry.VpcId,
			"Vip":     entry.Vip,
			"Label":   entry.Label,
			"shared":  entry.Shared,
		}
		return data, nil
	}

	return nil, errmsgs.GetNotFoundErrorFromString(fmt.Sprintf("OSS Single Tunnel not found with id: %s", id))
}

func (s OssService) GetOssClient(endpoint string) (*oss.Client, error) {
	// OSS SDK v2 uses a different configuration approach
	cfg := oss.LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			s.client.Config.AccessKey,
			s.client.Config.SecretKey,
			s.client.Config.SecurityToken,
		)).
		WithEndpoint(endpoint).
		WithRegion(s.client.RegionId)

	client := oss.NewClient(cfg)
	return client, nil
}

func (s OssService) GetBucketEndpointMap() (map[string]string, error) {
	schma := strings.ToLower(s.client.Config.Protocol)
	var ossEndpointMap map[string]string
	if len(ossEndpointMap) > 0 {
		return ossEndpointMap, nil
	} else {
		endpoints, err := s.GetOssEndpointList()
		if err != nil {
			return nil, errmsgs.WrapError(err)
		}
		for _, v := range endpoints {
			endpointData := v.(map[string]interface{})
			endpoint := endpointData["oss-endpoint"].(string)
			if !strings.HasPrefix(endpoint, "http") {
				endpoint = fmt.Sprintf("%s://%s", schma, endpoint)
			}
			ossEndpointMap[endpointData["cluster"].(string)] = endpoint
		}
	}
	return ossEndpointMap, nil
}

func (s OssService) GetBucketClient(bucketId string) (*oss.Client, error) {
	bucketInfo, err := s.DescribeOssBucket(bucketId)
	if err != nil {
		return nil, errmsgs.WrapError(err)
	}
	if *bucketInfo.Name == "" {
		return nil, errmsgs.GetNotFoundErrorFromString("Bucket " + bucketId + " Not Found")
	}

	bucketEndpoint := *bucketInfo.ExtranetEndpoint
	schma := strings.ToLower(s.client.Config.Protocol)
	if !strings.HasPrefix(bucketEndpoint, "http") {
		endpoint := fmt.Sprintf("%s://%s", schma, bucketEndpoint)
		bucketEndpoint = endpoint
	}

	// OSS SDK v2 uses a different configuration approach
	cfg := oss.LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			s.client.Config.AccessKey,
			s.client.Config.SecretKey,
			s.client.Config.SecurityToken,
		)).
		WithEndpoint(bucketEndpoint).
		WithRegion(s.client.RegionId)

	client := oss.NewClient(cfg)

	return client, nil
}

func (s OssService) DescribeOssBucket(bucketId string) (*oss.BucketInfo, error) {
	var bucket_name, cluster_name, endpoint string
	parts := strings.Split(bucketId, COLON_SEPARATED)
	if len(parts) > 1 {
		bucket_name = parts[0]
		cluster_name = parts[1]
	} else {
		bucket_name = bucketId
	}
	endpoints, err := s.GetBucketEndpointMap()
	if err != nil {
		return nil, errmsgs.WrapError(err)
	}
	if cluster_name != "" {
		endpoint = endpoints[cluster_name]
	} else {
		for _, v := range endpoints {
			endpoint = v
			break
		}
	}
	client, err := s.GetOssClient(endpoint)
	if err != nil {
		return nil, errmsgs.WrapError(err)
	}
	var request *oss.GetBucketInfoRequest
	request.Bucket = &bucket_name
	bucketResult, err := client.GetBucketInfo(context.Background(), request)
	if err != nil {
		return nil, errmsgs.WrapError(err)
	}
	if bucketResult.StatusCode == 404 || &bucketResult.BucketInfo == nil {
		return nil, errmsgs.GetNotFoundErrorFromString("Bucket " + bucketId + " Not Found")
	}
	return &bucketResult.BucketInfo, nil
}

func (s OssService) DescribeOssBucketKms(bucketId string) (*oss.ApplyServerSideEncryptionByDefault, error) {
	client, err := s.GetBucketClient(bucketId)
	var bucket_name string
	if err != nil {
		return nil, errmsgs.WrapError(err)
	}
	parts := strings.Split(bucketId, COLON_SEPARATED)
	if len(parts) > 1 {
		bucket_name = parts[0]
	} else {
		bucket_name = bucketId
	}
	getResult, err := client.GetBucketEncryption(context.Background(), &oss.GetBucketEncryptionRequest{
		Bucket: oss.Ptr(bucket_name),
	})
	addDebug("BucketEncryption", getResult, nil, map[string]string{"bucketName": bucketId})
	if err != nil {
		return nil, errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, bucketId, "GetBucketEncryption", errmsgs.AlibabacloudStackLogGoSdkERROR)
	}
	if getResult.ServerSideEncryptionRule != nil && getResult.ServerSideEncryptionRule.ApplyServerSideEncryptionByDefault != nil {
		apply := getResult.ServerSideEncryptionRule.ApplyServerSideEncryptionByDefault
		return apply, nil
	}
	return nil, errmsgs.GetNotFoundErrorFromString("Bucket Encryption data Not Found")
}

func (s OssService) DescribeOssBucketLogging(bucketId string) (*oss.GetBucketLoggingResult, error) {
	ossClient, err := s.GetBucketClient(bucketId)
	if err != nil {
		return nil, errmsgs.WrapError(err)
	}
	var bucket_name string
	parts := strings.Split(bucketId, COLON_SEPARATED)
	if len(parts) > 1 {
		bucket_name = parts[0]
	} else {
		bucket_name = bucketId
	}
	result, err := ossClient.GetBucketLogging(context.Background(), &oss.GetBucketLoggingRequest{
		Bucket: &bucket_name,
	})
	log.Printf("GetBucketLogging result: %v", result)
	if err != nil {
		return nil, errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, bucketId, "GetBucketLogging", errmsgs.AlibabacloudStackOssGoSdk)
	}

	return result, nil
}
