package alibabacloudstack

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ons"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
)

type OnsService struct {
	client *connectivity.AlibabacloudStackClient
}

func (s *OnsService) DescribeOnsInstance(instanceid string) (response *OnsInstance, err error) {
	var requestInfo *ons.Client

	request := s.client.NewCommonRequest("POST", "Ons-inner", "2018-02-05", "ConsoleInstanceBaseInfo", "")
	request.QueryParams["OnsRegionId"] = s.client.RegionId
	request.QueryParams["PreventCache"] = ""
	request.QueryParams["InstanceId"] = instanceid

	var resp = &OnsInstance{}
	bresponse, err := s.client.ProcessCommonRequest(request)
	addDebug(request.GetActionName(), bresponse, request, request.QueryParams)
	log.Printf(" response of raw ConsoleInstanceBaseInfo : %s", bresponse)
	if err != nil {
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return resp, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "ConsoleInstanceBaseInfo", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug("ConsoleInstanceBaseInfo", response, requestInfo, request)
	err = json.Unmarshal(bresponse.GetHttpContentBytes(), resp)
	if err != nil {
		return resp, errmsgs.WrapError(err)
	}

	if bresponse != nil && !resp.Success {
		return resp, errmsgs.WrapError(err)
	}

	return resp, nil
}

type TopicStruct struct {
	Data      string `json:"Data"`
	Message   string `json:"Message"`
	RequestID string `json:"RequestId"`
	Success   bool   `json:"Success"`
	Code      int    `json:"Code"`
}

func (s *OnsService) DescribeOnsTopic(id string) (response *Topic, err error) {
	did, err := ParseResourceId(id, 2)
	if err != nil {
		return response, errmsgs.WrapError(err)
	}
	TopicId := did[0]
	InstanceId := did[1]

	request := s.client.NewCommonRequest("POST", "Ons-inner", "2018-02-05", "ConsoleTopicList", "")
	request.QueryParams["Topic"] = TopicId
	request.QueryParams["OnsRegionId"] = s.client.RegionId
	request.QueryParams["PreventCache"] = ""
	request.QueryParams["InstanceId"] = InstanceId

	var resp = &Topic{}
	bresponse, err := s.client.ProcessCommonRequest(request)
	addDebug(request.GetActionName(), bresponse, request, request.QueryParams)
	log.Printf(" response of raw DescribeProjectMeta : %s", bresponse)
	if err != nil {
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return resp, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "ConsoleTopicList", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	err = json.Unmarshal(bresponse.GetHttpContentBytes(), resp)
	if err != nil {
		return resp, errmsgs.WrapError(err)
	}

	if len(resp.Data) < 1 || resp.Code == 200 {
		return resp, errmsgs.WrapError(err)
	}

	return resp, nil
}

func (s *OnsService) DescribeOnsGroup(id string) (response *OnsGroup, err error) {
	did, err := ParseResourceId(id, 2)
	if err != nil {
		return response, errmsgs.WrapError(err)
	}
	GroupId := did[0]
	InstanceId := did[1]

	request := s.client.NewCommonRequest("POST", "Ons-inner", "2018-02-05", "ConsoleGroupList", "")
	request.QueryParams["GroupId"] = GroupId
	request.QueryParams["OnsRegionId"] = s.client.RegionId
	request.QueryParams["PreventCache"] = ""
	request.QueryParams["InstanceId"] = InstanceId

	var resp = &OnsGroup{}
	bresponse, err := s.client.ProcessCommonRequest(request)
	addDebug(request.GetActionName(), bresponse, request, request.QueryParams)
	log.Printf(" response of raw ConsoleGroupList : %s", bresponse)
	if err != nil {
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return resp, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "ConsoleGroupList", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	err = json.Unmarshal(bresponse.GetHttpContentBytes(), resp)
	if err != nil {
		return resp, errmsgs.WrapError(err)
	}

	if len(resp.Data) < 1 || resp.Code == 200 {
		return resp, errmsgs.WrapError(err)
	}

	return resp, nil
}

func (s *OnsService) DescribeMqttInstance(id string) (map[string]interface{}, error) {
	reqQuery := map[string]interface{}{
		"Platform":       "onsConsole",
		"OnsRegionId":    s.client.RegionId,
		"MqttInstanceId": id,
		"PreventCache":   time.Now().UnixNano() / 1e6,
	}

	response, err := s.client.DoTeaRequest("POST", "Ons-inner", "2018-02-05", "ConsoleMqttInstanceBaseInfo", "", nil, reqQuery, nil)
	if err != nil {
		if errmsgs.IsExpectedErrors(err, "No matched instances are found in this account") {
			return nil, errmsgs.GetNotFoundErrorFromString(fmt.Sprintf("MqttInstance %s not found", id))
		}
		return nil, err
	}
	if data, ok := response["Data"]; !ok || data == nil {
		return nil, errmsgs.GetNotFoundErrorFromString(fmt.Sprintf("MqttInstance %s not found", id))
	} else {
		return data.(map[string]interface{}), nil
	}
}

func (s *OnsService) DescribeOnsMqttTopic(id string) (map[string]interface{}, error) {
	parts := strings.SplitN(id, ":", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid id format, expected {InstanceId}:{Topic}")
	}
	instanceId := parts[0]
	topic := parts[1]

	reqQuery := map[string]interface{}{
		"PreventCache":   time.Now().UnixNano() / 1e6,
		"Dauth_url_hash": "mqtt%2Fconsole%2Finstances%2FinstanceDetail",
		"OnsRegionId":    s.client.RegionId,
		"InstanceId":     instanceId,
		"CurrentPage":    1,
		"PageSize":       100,
		"isFuzzy":        false,
	}

	response, err := s.client.DoTeaRequest("GET", "Ons-inner", "2018-02-05", "ConsoleTopicListInPage", "", nil, reqQuery, nil)
	if err != nil {
		if errmsgs.IsExpectedErrors(err, "The specified instance does not exist.") {
			return nil, errmsgs.GetNotFoundErrorFromString(fmt.Sprintf("Mqtt Topic %s not found", id))
		}
		return nil, err
	}
	if data, ok := response["Data"].([]interface{}); ok {
		for _, item := range data {
			if topicItem, ok := item.(map[string]interface{}); ok {
				if topicItem["topic"].(string) == topic {
					return topicItem, nil
				}
			}
		}
	}

	return nil, errmsgs.GetNotFoundErrorFromString(fmt.Sprintf("Mqtt Topic %s not found", id))
}

func (s *OnsService) DescribeOnsMqttGroup(id string) (map[string]interface{}, error) {
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return nil, errmsgs.WrapError(err)
	}
	instanceId := parts[0]
	groupId := parts[1]

	// Prepare parameters for the API call
	query := map[string]interface{}{
		"MqttInstanceId": instanceId,
		"currentPage":    1,
		"pageSize":       100,
		"Platform":       "onsConsole",
		"OnsRegionId":    s.client.RegionId,
		"PreventCache":   time.Now().UnixNano() / 1e6, // milliseconds
		"Dauth_url_hash": "mqtt%2Fconsole%2Finstances%2FinstanceDetail",
	}

	resp, err := s.client.DoTeaRequest("POST", "Ons-inner", "2018-02-05", "ConsoleMqttListGroupIdInPage", "", nil, query, nil)
	if err != nil {
		if errmsgs.IsExpectedErrors(err, "No matched instances are found in this account.") {
			return nil, errmsgs.GetNotFoundErrorFromString("Resource Mqtt GroupId not found")
		}
		return nil, errmsgs.WrapError(err)
	}

	dataList, ok := resp["Data"].([]interface{})
	if !ok || len(dataList) == 0 {
		return nil, errmsgs.GetNotFoundErrorFromString("Resource Mqtt GroupId not found")
	}

	var targetGroup map[string]interface{}
	for _, item := range dataList {
		if groupInfo, ok := item.(map[string]interface{}); ok {
			if groupInfo["groupId"] == groupId {
				targetGroup = groupInfo
				break
			}
		}
	}
	if targetGroup == nil {
		return nil, errmsgs.GetNotFoundErrorFromString("Resource Mqtt GroupId not found")
	}

	return targetGroup, nil
}
