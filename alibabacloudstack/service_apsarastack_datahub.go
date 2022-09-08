package alibabacloudstack

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"strings"
	"time"

	"encoding/json"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"

	"github.com/aliyun/aliyun-datahub-sdk-go/datahub"
)

type DatahubService struct {
	client *connectivity.AlibabacloudStackClient
}

func (s *DatahubService) DescribeDatahubProject(id string) (*datahub.GetProjectResult, error) {
	var requestInfo *ecs.Client
	request := requests.NewCommonRequest()
	resp := &datahub.GetProjectResult{}

	request.Method = "GET"
	request.Product = "datahub"
	request.Version = "2019-11-20"
	if strings.ToLower(s.client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.ApiName = "GetProject"
	request.Headers = map[string]string{
		"RegionId":              s.client.RegionId,
		"x-acs-resourcegroupid": s.client.ResourceGroup,
		"x-acs-regionid":        s.client.RegionId,
		"x-acs-organizationid":  s.client.Department,
	}
	request.QueryParams = map[string]string{
		"AccessKeySecret": s.client.SecretKey,
		"AccessKeyId":     s.client.AccessKey,
		"Product":         "datahub",
		"Department":      s.client.Department,
		"ResourceGroup":   s.client.ResourceGroup,
		"RegionId":        s.client.RegionId,
		"Action":          "GetProject",
		"Version":         "2019-11-20",
		"ProjectName":     id,
	}

	raw, err := s.client.WithEcsClient(func(dataHubClient *ecs.Client) (interface{}, error) {
		return dataHubClient.ProcessCommonRequest(request)
	})
	if err != nil {
		if isDatahubNotExistError(err) {
			return resp, WrapErrorf(err, NotFoundMsg, AlibabacloudStackDatahubSdkGo)
		}
		return resp, WrapErrorf(err, DefaultErrorMsg, id, "GetProject", AlibabacloudStackDatahubSdkGo)
	}
	addDebug("GetProject", raw, requestInfo, request)

	bresponse, _ := raw.(*responses.CommonResponse)
	err = json.Unmarshal(bresponse.GetHttpContentBytes(), resp)
	if err != nil {
		return nil, WrapErrorf(Error(GetNotFoundMessage("DatahubProject", id)), NotFoundMsg, ProviderERROR)
	}

	return resp, nil
}

func (s *DatahubService) WaitForDatahubProject(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)

	for {
		object, err := s.DescribeDatahubProject(id)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return WrapError(err)
			}
		}

		if time.Now().After(deadline) {
			objstringfy, err := convertArrayObjectToJsonString(object)
			if err != nil {
				return WrapError(err)
			}
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, objstringfy, id, ProviderERROR)
		}

	}
}

func (s *DatahubService) DescribeDatahubSubscription(id string) (*datahub.GetSubscriptionResult, error) {
	subscription := &datahub.GetSubscriptionResult{}
	var requestInfo *ecs.Client
	request := requests.NewCommonRequest()
	parts, err := ParseResourceId(id, 3)
	if err != nil {
		return subscription, WrapError(err)
	}
	projectName, topicName, subId := parts[0], parts[1], parts[2]

	request.Method = "GET"
	request.Product = "datahub"
	request.Version = "2019-11-20"
	if strings.ToLower(s.client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.ApiName = "GetSubscriptionOffset"
	request.Headers = map[string]string{
		"RegionId":              s.client.RegionId,
		"x-acs-resourcegroupid": s.client.ResourceGroup,
		"x-acs-regionid":        s.client.RegionId,
		"x-acs-organizationid":  s.client.Department,
	}
	request.QueryParams = map[string]string{
		"AccessKeySecret":  s.client.SecretKey,
		"AccessKeyId":      s.client.AccessKey,
		"Product":          "datahub",
		"Department":       s.client.Department,
		"ResourceGroup":    s.client.ResourceGroup,
		"RegionId":         s.client.RegionId,
		"Action":           "GetSubscriptionOffset",
		"Version":          "2019-11-20",
		"ProjectName":      projectName,
		"TopicName":        topicName,
		"SubscriptionId":   subId,
		"SignatureMethod":  "HMAC-SHA256",
		"Format":           "JSON",
		"SignatureVersion": "2.1",
	}

	raw, err := s.client.WithEcsClient(func(dataHubClient *ecs.Client) (interface{}, error) {
		return dataHubClient.ProcessCommonRequest(request)
	})

	if err != nil {
		if isDatahubNotExistError(err) {
			return subscription, WrapErrorf(err, NotFoundMsg, AlibabacloudStackDatahubSdkGo)
		}
		return subscription, WrapErrorf(err, DefaultErrorMsg, id, "GetSubscription", AlibabacloudStackDatahubSdkGo)
	}
	if debugOn() {
		requestMap := make(map[string]string)
		requestMap["ProjectName"] = projectName
		requestMap["TopicName"] = topicName
		requestMap["SubId"] = subId
		addDebug("GetProject", raw, requestInfo, requestMap)
	}
	bresponse, _ := raw.(*responses.CommonResponse)
	err = json.Unmarshal(bresponse.GetHttpContentBytes(), subscription)
	if err != nil {
		return nil, WrapErrorf(Error(GetNotFoundMessage("DatahubProject", id)), NotFoundMsg, ProviderERROR)
	}
	return subscription, nil
}

func (s *DatahubService) WaitForDatahubSubscription(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	parts, err := ParseResourceId(id, 3)
	if err != nil {
		return WrapError(err)
	}
	topicName, subId := parts[1], parts[2]
	//for {
	object, err := s.DescribeDatahubSubscription(id)
	if err != nil {
		if NotFoundError(err) {
			if status == Deleted {
				return nil
			}
		} else {
			return WrapError(err)
		}
	}
	if object.TopicName == topicName && object.SubId == subId && status != Deleted {
		return nil
	}
	if time.Now().After(deadline) {
		return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, object.TopicName+":"+object.SubId, parts[1]+":"+parts[2], ProviderERROR)
	}

	return nil
}

func (s *DatahubService) DescribeDatahubTopic(id string) (*GetTopicResult, error) {
	topic := &GetTopicResult{}
	var requestInfo *ecs.Client
	request := requests.NewCommonRequest()
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return topic, WrapError(err)
	}
	projectName, topicName := parts[0], parts[1]
	request.Method = "GET"
	request.Product = "datahub"
	request.Version = "2019-11-20"
	if strings.ToLower(s.client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.ApiName = "GetTopic"
	request.Headers = map[string]string{
		"RegionId":              s.client.RegionId,
		"x-acs-resourcegroupid": s.client.ResourceGroup,
		"x-acs-regionid":        s.client.RegionId,
		"x-acs-organizationid":  s.client.Department,
	}
	request.QueryParams = map[string]string{
		"AccessKeySecret": s.client.SecretKey,
		"AccessKeyId":     s.client.AccessKey,
		"Product":         "datahub",
		"Department":      s.client.Department,
		"ResourceGroup":   s.client.ResourceGroup,
		"RegionId":        s.client.RegionId,
		"Action":          "GetTopic",
		"Version":         "2019-11-20",
		"ProjectName":     projectName,
		"TopicName":       topicName,
	}

	raw, err := s.client.WithEcsClient(func(dataHubClient *ecs.Client) (interface{}, error) {
		return dataHubClient.ProcessCommonRequest(request)
	})

	if err != nil {
		if isDatahubNotExistError(err) {
			return topic, WrapErrorf(err, NotFoundMsg, AlibabacloudStackDatahubSdkGo)
		}
		return topic, WrapErrorf(err, DefaultErrorMsg, id, "GetTopic", AlibabacloudStackDatahubSdkGo)
	}
	if debugOn() {
		requestMap := make(map[string]string)
		requestMap["ProjectName"] = projectName
		requestMap["TopicName"] = topicName
		addDebug("GetTopic", raw, requestInfo, requestMap)
	}
	bresponse, _ := raw.(*responses.CommonResponse)
	err = json.Unmarshal(bresponse.GetHttpContentBytes(), topic)
	if err != nil {
		return nil, WrapErrorf(Error(GetNotFoundMessage("DatahubProject", id)), NotFoundMsg, ProviderERROR)
	}
	return topic, nil
}

func (s *DatahubService) WaitForDatahubTopic(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return WrapError(err)
	}
	projectName, topicName := parts[0], parts[1]
	for {
		object, err := s.DescribeDatahubTopic(id)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return WrapError(err)
			}
		}
		if object.ProjectName == projectName && object.TopicName == topicName && status != Deleted {
			return nil
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, object.ProjectName+":"+object.TopicName, id, ProviderERROR)
		}

	}
}

func convUint64ToDate(t uint64) string {
	return time.Unix(int64(t), 0).Format("2006-01-02 15:04:05")
}

func getNow() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func getRecordSchema(typeMap map[string]interface{}) (recordSchema *datahub.RecordSchema) {
	recordSchema = datahub.NewRecordSchema()

	for k, v := range typeMap {
		recordSchema.AddField(datahub.Field{Name: string(k), Type: datahub.FieldType(v.(string))})
	}

	return recordSchema
}

func isRetryableDatahubError(err error) bool {
	if e, ok := err.(*datahub.DatahubClientError); ok && e.StatusCode >= 500 {
		return true
	}

	return false
}

// It is proactive defense to the case that SDK extends new datahub objects.
const (
	DoesNotExist = "does not exist"
)

func isDatahubNotExistError(err error) bool {
	return IsExpectedErrors(err, []string{datahub.NoSuchProject, datahub.NoSuchTopic, datahub.NoSuchShard, datahub.NoSuchSubscription, DoesNotExist})
}

func isTerraformTestingDatahubObject(name string) bool {
	prefixes := []string{
		"tf_testAcc",
		"tf_test_",
		"testAcc",
	}

	for _, prefix := range prefixes {
		if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
			return true
		}
	}

	return false
}

func getDefaultRecordSchemainMap() map[string]interface{} {

	return map[string]interface{}{
		"string_field": "STRING",
	}
}

func recordSchemaToMap(fields []datahub.Field) map[string]string {
	result := make(map[string]string)
	for _, f := range fields {
		result[f.Name] = f.Type.String()
	}

	return result
}
