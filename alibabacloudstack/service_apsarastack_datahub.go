package alibabacloudstack

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"strings"
	"time"

	"encoding/json"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"

	"github.com/aliyun/aliyun-datahub-sdk-go/datahub"
)

type DatahubService struct {
	client *connectivity.AlibabacloudStackClient
}

func (s *DatahubService) DoDatahubGetkafkagroupRequest(id string) (*datahub.GetProjectResult, error) { 
	return s.DescribeDatahubProject(id)
}

func (s *DatahubService) DescribeDatahubProject(id string) (*datahub.GetProjectResult, error) {
	resp := &datahub.GetProjectResult{}

	request := s.client.NewCommonRequest("GET", "datahub", "2019-11-20", "GetProject", "")
	request.QueryParams["ProjectName"] = id

	raw, err := s.client.WithEcsClient(func(dataHubClient *ecs.Client) (interface{}, error) {
		return dataHubClient.ProcessCommonRequest(request)
	})
	bresponse, ok := raw.(*responses.CommonResponse)
	if err != nil {
		if isDatahubNotExistError(err) {
			return resp, errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackDatahubSdkGo)
		}
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		return resp, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, "GetProject", errmsgs.AlibabacloudStackDatahubSdkGo, errmsg)
	}
	addDebug("GetProject", raw, nil, request)

	err = json.Unmarshal(bresponse.GetHttpContentBytes(), resp)
	if err != nil {
		return nil, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("DatahubProject", id)), errmsgs.NotFoundMsg, errmsgs.ProviderERROR)
	}

	return resp, nil
}

func (s *DatahubService) WaitForDatahubProject(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)

	for {
		object, err := s.DescribeDatahubProject(id)
		if err != nil {
			if errmsgs.NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return errmsgs.WrapError(err)
			}
		}

		if time.Now().After(deadline) {
			objstringfy, err := convertArrayObjectToJsonString(object)
			if err != nil {
				return errmsgs.WrapError(err)
			}
			return errmsgs.WrapErrorf(err, errmsgs.WaitTimeoutMsg, id, GetFunc(1), timeout, objstringfy, id, errmsgs.ProviderERROR)
		}

	}
}

func (s *DatahubService) DoDatahubGetsubscriptionoffsetRequest(id string) (*datahub.GetSubscriptionResult, error) {
    return s.DescribeDatahubSubscription(id)
}
func (s *DatahubService) DescribeDatahubSubscription(id string) (*datahub.GetSubscriptionResult, error) {
	subscription := &datahub.GetSubscriptionResult{}
	parts, err := ParseResourceId(id, 3)
	if err != nil {
		return subscription, errmsgs.WrapError(err)
	}
	projectName, topicName, subId := parts[0], parts[1], parts[2]

	request := s.client.NewCommonRequest("GET", "datahub", "2019-11-20", "GetSubscriptionOffset", "")
	request.QueryParams["ProjectName"] = projectName
	request.QueryParams["TopicName"] = topicName
	request.QueryParams["SubscriptionId"] = subId
	request.QueryParams["SignatureMethod"] = "HMAC-SHA256"
	request.QueryParams["Format"] = "JSON"
	request.QueryParams["SignatureVersion"] = "2.1"

	raw, err := s.client.WithEcsClient(func(dataHubClient *ecs.Client) (interface{}, error) {
		return dataHubClient.ProcessCommonRequest(request)
	})
	bresponse, ok := raw.(*responses.CommonResponse)
	if err != nil {
		if isDatahubNotExistError(err) {
			return subscription, errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackDatahubSdkGo)
		}
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		return subscription, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, "GetSubscription", errmsgs.AlibabacloudStackDatahubSdkGo, errmsg)
	}
	if debugOn() {
		requestMap := make(map[string]string)
		requestMap["ProjectName"] = projectName
		requestMap["TopicName"] = topicName
		requestMap["SubId"] = subId
		addDebug("GetProject", raw, nil, requestMap)
	}
	err = json.Unmarshal(bresponse.GetHttpContentBytes(), subscription)
	if err != nil {
		return nil, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("DatahubProject", id)), errmsgs.NotFoundMsg, errmsgs.ProviderERROR)
	}
	return subscription, nil
}

func (s *DatahubService) WaitForDatahubSubscription(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	parts, err := ParseResourceId(id, 3)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	topicName, subId := parts[1], parts[2]
	//for {
	object, err := s.DescribeDatahubSubscription(id)
	if err != nil {
		if errmsgs.NotFoundError(err) {
			if status == Deleted {
				return nil
			}
		} else {
			return errmsgs.WrapError(err)
		}
	}
	if object.TopicName == topicName && object.SubId == subId && status != Deleted {
		return nil
	}
	if time.Now().After(deadline) {
		return errmsgs.WrapErrorf(err, errmsgs.WaitTimeoutMsg, id, GetFunc(1), timeout, object.TopicName+":"+object.SubId, parts[1]+":"+parts[2], errmsgs.ProviderERROR)
	}

	return nil
}

func (s *DatahubService) DoDatahubGettopicRequest(id string) (*GetTopicResult, error) {
    return s.DescribeDatahubTopic(id)
}
func (s *DatahubService) DescribeDatahubTopic(id string) (*GetTopicResult, error) {
	topic := &GetTopicResult{}
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return topic, errmsgs.WrapError(err)
	}
	projectName, topicName := parts[0], parts[1]

	request := s.client.NewCommonRequest("GET", "datahub", "2019-11-20", "GetTopic", "")
	request.QueryParams["ProjectName"] = projectName
	request.QueryParams["TopicName"] = topicName

	raw, err := s.client.WithEcsClient(func(dataHubClient *ecs.Client) (interface{}, error) {
		return dataHubClient.ProcessCommonRequest(request)
	})
	bresponse, ok := raw.(*responses.CommonResponse)
	if err != nil {
		if isDatahubNotExistError(err) {
			return topic, errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackDatahubSdkGo)
		}
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		return topic, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, "GetTopic", errmsgs.AlibabacloudStackDatahubSdkGo, errmsg)
	}
	if debugOn() {
		requestMap := make(map[string]string)
		requestMap["ProjectName"] = projectName
		requestMap["TopicName"] = topicName
		addDebug("GetTopic", raw, nil, requestMap)
	}
	err = json.Unmarshal(bresponse.GetHttpContentBytes(), topic)
	if err != nil {
		return nil, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("DatahubProject", id)), errmsgs.NotFoundMsg, errmsgs.ProviderERROR)
	}
	return topic, nil
}

func (s *DatahubService) WaitForDatahubTopic(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	projectName, topicName := parts[0], parts[1]
	for {
		object, err := s.DescribeDatahubTopic(id)
		if err != nil {
			if errmsgs.NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return errmsgs.WrapError(err)
			}
		}
		if object.ProjectName == projectName && object.TopicName == topicName && status != Deleted {
			return nil
		}
		if time.Now().After(deadline) {
			return errmsgs.WrapErrorf(err, errmsgs.WaitTimeoutMsg, id, GetFunc(1), timeout, object.ProjectName+":"+object.TopicName, id, errmsgs.ProviderERROR)
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
	return errmsgs.IsExpectedErrors(err, []string{datahub.NoSuchProject, datahub.NoSuchTopic, datahub.NoSuchShard, datahub.NoSuchSubscription, DoesNotExist})
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
