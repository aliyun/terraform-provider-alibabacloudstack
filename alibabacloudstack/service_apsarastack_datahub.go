package alibabacloudstack

import (
	"fmt"
	"strings"
	"time"

	"encoding/json"

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

type DatahubProject struct {
	Comment     string `json:"Comment"`
	ProjectName string `json:"ProjectName"`
	CreateTime  int64  `json:"CreateTime"`
	UpdateTime  int64  `json:"UpdateTime"`
	Creator     string `json:"Creator"`
}

type ListProjectResult struct {
	TotalCount int `json:"TotalCount"`
	List       struct {
		Project []DatahubProject `json:"Project"`
	} `json:"List"`
}

func (s *DatahubService) DescribeDatahubProject(id string) (*datahub.GetProjectResult, error) {
	resp := &datahub.GetProjectResult{}

	request := s.client.NewCommonRequest("GET", "datahub", "2019-11-20", "GetProject", "")
	request.QueryParams["ProjectName"] = id

	bresponse, err := s.client.ProcessCommonRequest(request)
	if err != nil {
		if isDatahubNotExistError(err) {
			return resp, errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackDatahubSdkGo)
		}
		if bresponse == nil {
			return resp, errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return resp, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, "GetProject", errmsgs.AlibabacloudStackDatahubSdkGo, errmsg)
	}
	addDebug("GetProject", bresponse, nil, request)

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

type SubscriptionEntry struct {
	SubId       string `json:"SubscriptionId"`
	Application string `json:"Application"`
	Creator     string `json:"Creator"`
	Type        string `json:"Type"`
	State       int    `json:"State,omitempty"`
	Comment     string `json:"Comment,omitempty"`
}

type ListSubscriptionResult struct {
	TotalCount int `json:"TotalCount"`
	List       struct {
		Subscriptions []SubscriptionEntry `json:"Subscription"`
	} `json:"List"`
}

func (s *DatahubService) DescribeDatahubSubscription(id string) (*SubscriptionEntry, error) {
	subscriptions := &ListSubscriptionResult{}
	subscription := &SubscriptionEntry{}
	parts, err := ParseResourceId(id, 3)
	if err != nil {
		return subscription, errmsgs.WrapError(err)
	}
	projectName, topicName, subId := parts[0], parts[1], parts[2]

	request := s.client.NewCommonRequest("GET", "datahub", "2019-11-20", "ListSubscriptions", "")
	request.QueryParams["ProjectName"] = projectName
	request.QueryParams["TopicName"] = topicName
	request.QueryParams["Keyword"] = subId
	request.QueryParams["PageNumber"] = "1"
	request.QueryParams["PageSize"] = "100"

	bresponse, err := s.client.ProcessCommonRequest(request)
	if err != nil {
		if isDatahubNotExistError(err) {
			return subscription, errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackDatahubSdkGo)
		}
		if bresponse == nil {
			return subscription, errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return subscription, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, "GetSubscription", errmsgs.AlibabacloudStackDatahubSdkGo, errmsg)
	}
	if debugOn() {
		requestMap := make(map[string]string)
		requestMap["ProjectName"] = projectName
		requestMap["TopicName"] = topicName
		requestMap["SubId"] = subId
		addDebug("GetProject", bresponse, nil, requestMap)
	}
	err = json.Unmarshal(bresponse.GetHttpContentBytes(), subscriptions)
	if err != nil {
		return nil, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("DatahubProject", id)), errmsgs.NotFoundMsg, errmsgs.ProviderERROR)
	}
	for _, sub := range subscriptions.List.Subscriptions {
		if sub.SubId == subId {
			return &sub, nil
		}
	}
	return subscription, errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackDatahubSdkGo)
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

	bresponse, err := s.client.ProcessCommonRequest(request)
	if err != nil {
		if isDatahubNotExistError(err) {
			return topic, errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackDatahubSdkGo)
		}
		if bresponse == nil {
			return topic, errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return topic, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, "GetTopic", errmsgs.AlibabacloudStackDatahubSdkGo, errmsg)
	}
	if debugOn() {
		requestMap := make(map[string]string)
		requestMap["ProjectName"] = projectName
		requestMap["TopicName"] = topicName
		addDebug("GetTopic", bresponse, nil, requestMap)
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
	return errmsgs.IsExpectedErrors(err, datahub.NoSuchProject, datahub.NoSuchTopic, datahub.NoSuchShard, datahub.NoSuchSubscription, DoesNotExist)
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

// The following code has been adapted to align with existing naming conventions and structure.
// Original function names and structures have been modified where necessary to maintain consistency.

func (s *DatahubService) DescribeDatahubKafkaGroup(id string) (map[string]interface{}, error) {
	// Split the id to get ProjectName and GroupName
	parts := strings.SplitN(id, ":", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid kafka group id: %s", id)
	}
	projectName := parts[0]
	groupName := parts[1]

	// Prepare request query parameters
	reqQuery := map[string]interface{}{
		"ProjectName": projectName,
		"PageSize":    10,
		"PageNumber":  1,
		"Keyword":     groupName,
	}

	// Call ListKafkaGroup API to find the target group
	response, err := s.client.DoTeaRequest("GET", "datahub", "2019-11-20", "ListKafkaGroup", "", nil, reqQuery, nil)
	if err != nil {
		return nil, err
	}

	// Check if the response contains the expected data structure

	list, ok := response["List"].([]interface{})
	if !ok || len(list) == 0 {
		return nil, errmsgs.GetNotFoundErrorFromString(fmt.Sprintf("kafka group %s not found", groupName))
	}

	// Iterate over the list to find the matching GroupName
	for _, item := range list {
		kafkaGroup, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		if kafkaGroup["GroupName"].(string) == groupName {
			return kafkaGroup, nil
		}
	}

	// If no match is found, return a not found error
	return nil, errmsgs.GetNotFoundErrorFromString(fmt.Sprintf("kafka group %s not found", groupName))
}
